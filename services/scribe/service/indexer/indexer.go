package indexer

import (
	"context"
	"errors"
	"fmt"
	"github.com/synapsecns/sanguine/services/scribe/backend"
	scribeTypes "github.com/synapsecns/sanguine/services/scribe/types"

	"github.com/synapsecns/sanguine/services/scribe/logger"
	"math/big"
	"time"

	"github.com/lmittmann/w3"
	"github.com/lmittmann/w3/module/eth"
	"github.com/lmittmann/w3/w3types"
	"github.com/synapsecns/sanguine/core/mapmutex"
	"github.com/synapsecns/sanguine/core/metrics"
	"go.opentelemetry.io/otel/attribute"
	otelMetrics "go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	lru "github.com/hashicorp/golang-lru"
	"github.com/jpillora/backoff"
	"github.com/synapsecns/sanguine/services/scribe/config"
	"github.com/synapsecns/sanguine/services/scribe/db"
	"golang.org/x/sync/errgroup"
)

// Indexer is a backfiller that fetches logs for a specific contract.
type Indexer struct {
	// indexerConfig holds all the metadata needed for logging and indexing.
	indexerConfig scribeTypes.IndexerConfig
	// eventDB is the database to store event data in.
	eventDB db.EventDB
	// client is the client for filtering.
	client []backend.ScribeBackend
	// cache is a cache for txHashes.
	cache *lru.Cache
	// mux is the mutex used to prevent double inserting logs from the same tx
	mux mapmutex.StringerMapMutex
	// handler is the metrics handler for the scribe.
	handler metrics.Handler
	// blockMeter is an otel historgram for doing metrics on block heights by chain
	blockMeter otelMetrics.Int64Histogram
	// refreshRate is the rate at which the indexer will refresh when livefilling.
	refreshRate uint64
	// toTip is a boolean signifying if the indexer is livefilling to the tip.
	toTip bool
}

// retryTolerance is the number of times to retry a failed operation before rerunning the entire Backfill function.
const retryTolerance = 20

// txNotSupportedError is for handling the legacy Arbitrum tx type.
const txNotSupportedError = "transaction type not supported"

// invalidTxVRSError is for handling Aurora VRS error.
const invalidTxVRSError = "invalid transaction v, r, s values"

// txNotFoundError is for handling omniRPC errors for BSx.
const txNotFoundError = "not found"

// txData returns the transaction data for a given transaction hash.
type txData struct {
	receipt     types.Receipt
	transaction types.Transaction
	blockHeader types.Header
	success     bool
}

var errNoContinue = errors.New("encountered unreconcilable error, will not attempt to store tx")

// errNoTx indicates a tx cannot be parsed, this is only returned when the tx doesn't match our data model.
var errNoTx = errors.New("tx is not supported by the client")

// NewIndexer creates a new backfiller for a contract.
func NewIndexer(chainConfig config.ChainConfig, addresses []common.Address, eventDB db.EventDB, client []backend.ScribeBackend, handler metrics.Handler, blockMeter otelMetrics.Int64Histogram, toTip bool) (*Indexer, error) {
	cache, err := lru.New(500)
	if err != nil {
		return nil, fmt.Errorf("could not initialize cache: %w", err)
	}

	refreshRate := uint64(1)
	if len(addresses) > 1 || len(addresses) == 0 { // livefill settings
		chainConfig.GetLogsRange = chainConfig.LivefillRange
		chainConfig.GetLogsBatchAmount = 1
	} else {
		for i := range chainConfig.Contracts { // get the refresh rate for the contract
			contract := chainConfig.Contracts[i]
			// Refresh rate for more than one contract is 1 second, the refresh rate set in the config is used when it is the only contract.
			if contract.Address == addresses[0].String() && contract.RefreshRate > 0 {
				refreshRate = contract.RefreshRate
				break
			}
		}
	}

	indexerConfig := scribeTypes.IndexerConfig{
		Addresses:          addresses,
		GetLogsRange:       chainConfig.GetLogsRange,
		GetLogsBatchAmount: chainConfig.GetLogsBatchAmount,
		StoreConcurrency:   chainConfig.StoreConcurrency,
		ChainID:            chainConfig.ChainID,
	}

	return &Indexer{
		indexerConfig: indexerConfig,
		eventDB:       eventDB,
		client:        client,
		cache:         cache,
		mux:           mapmutex.NewStringerMapMutex(),
		handler:       handler,
		blockMeter:    blockMeter,
		refreshRate:   refreshRate,
		toTip:         toTip,
	}, nil
}

// UpdateAddress updates the address arrays for the indexer.
func (x *Indexer) UpdateAddress(addresses []common.Address) {
	x.indexerConfig.Addresses = addresses
}

// GetIndexerConfig returns the indexer config.
func (x *Indexer) GetIndexerConfig() scribeTypes.IndexerConfig {
	return x.indexerConfig
}

// RefreshRate returns the refresh rate for the indexer.
func (x *Indexer) RefreshRate() uint64 {
	return x.refreshRate
}

// Index retrieves logs, receipts, and transactions for a contract from a given range and does so in the following manner.
// 1. Get logs for the contract in chunks of batch requests.
// 2. Iterate through each log's Tx Hash and performs the following
//   - Get the receipt for each log and store it and all of its logs.
//   - Get the transaction for each log and store it.
//
//nolint:gocognit, cyclop
func (x *Indexer) Index(parentCtx context.Context, startHeight uint64, endHeight uint64) (err error) {
	ctx, span := x.handler.Tracer().Start(parentCtx, "contract.Backfill", trace.WithAttributes(
		attribute.Int("chain", int(x.indexerConfig.ChainID)),
		attribute.String("address", x.addressesToString(x.indexerConfig.Addresses)),
		attribute.Int("start", int(startHeight)),
		attribute.Int("end", int(endHeight)),
	))

	defer func() {
		metrics.EndSpanWithErr(span, err)
	}()

	g, groupCtx := errgroup.WithContext(ctx)

	// For logging
	x.indexerConfig.StartHeight = startHeight
	x.indexerConfig.EndHeight = endHeight

	// logsChain and errChan are used to pass logs from rangeFilter onto the next stage of the backfill process.
	logsChan, errChan := x.getLogs(groupCtx, startHeight, endHeight)

	// Reads from the local logsChan and stores the logs and associated receipts / txs.
	g.Go(func() error {
		concurrentCalls := 0
		gS, storeCtx := errgroup.WithContext(ctx)
		// could change this to for - range
		for {
			select {
			case <-groupCtx.Done():
				logger.ReportIndexerError(ctx.Err(), x.indexerConfig, logger.ContextCancelled)
				return fmt.Errorf("context canceled while storing and retrieving logs: %w", groupCtx.Err())
			case log, ok := <-logsChan: // empty log passed when ok is false.
				if !ok {
					return nil
				}
				concurrentCalls++
				gS.Go(func() error {
					// another goroutine is already storing this receipt
					locker, ok := x.mux.TryLock(log.TxHash)
					if !ok {
						return nil
					}
					defer locker.Unlock()

					// Check if the txHash has already been stored in the cache.
					if _, ok := x.cache.Get(log.TxHash); ok {
						return nil
					}

					err := x.store(storeCtx, log)
					if err != nil {
						logger.ReportIndexerError(err, x.indexerConfig, logger.StoreError)

						return fmt.Errorf("could not store log: %w", err)
					}

					return nil
				})

				// Stop spawning store threads and wait
				if concurrentCalls >= x.indexerConfig.StoreConcurrency || x.indexerConfig.ConcurrencyThreshold > endHeight-log.BlockNumber {
					if err = gS.Wait(); err != nil {
						return fmt.Errorf("error waiting for go routines: %w", err)
					}

					// Reset context TODO make this better
					gS, storeCtx = errgroup.WithContext(ctx)
					concurrentCalls = 0
					err = x.eventDB.StoreLastIndexedMultiple(ctx, x.indexerConfig.Addresses, x.indexerConfig.ChainID, log.BlockNumber)
					if err != nil {
						logger.ReportIndexerError(err, x.indexerConfig, logger.StoreError)
						return fmt.Errorf("could not store last indexed block: %w", err)
					}

					x.blockMeter.Record(ctx, int64(log.BlockNumber), otelMetrics.WithAttributeSet(
						attribute.NewSet(attribute.Int64("start_block", int64(startHeight)), attribute.Int64("chain_id", int64(x.indexerConfig.ChainID)))),
					)
				}

			case errFromChan := <-errChan:
				logger.ReportIndexerError(fmt.Errorf("errChan returned an err %s", errFromChan), x.indexerConfig, logger.GetLogsError)
				return fmt.Errorf("errChan returned an err %s", errFromChan)
			}
		}
	})

	err = g.Wait()

	if err != nil {
		return fmt.Errorf("could not backfill contract: %w \nChain: %d\nLog 's Contract Address: %s\n ", err, x.indexerConfig.ChainID, x.indexerConfig.Addresses)
	}

	err = x.eventDB.StoreLastIndexedMultiple(ctx, x.indexerConfig.Addresses, x.indexerConfig.ChainID, endHeight)
	if err != nil {
		return fmt.Errorf("could not store last indexed block: %w", err)
	}
	x.blockMeter.Record(ctx, int64(endHeight), otelMetrics.WithAttributeSet(
		attribute.NewSet(attribute.Int64("start_block", int64(startHeight)), attribute.Int64("chain_id", int64(x.indexerConfig.ChainID)))),
	)
	// LogEvent(InfoLevel, "Finished backfilling contract", LogData{"cid": x.indexerConfig.ChainID, "ca": x.addressesToString(x.indexerConfig.Addresses)})

	return nil
}

// TODO split two goroutines into sep functions for maintainability
// store stores the logs, receipts, and transactions for a tx hash.
//
//nolint:cyclop,gocognit,maintidx
func (x *Indexer) store(parentCtx context.Context, log types.Log) (err error) {
	ctx, span := x.handler.Tracer().Start(parentCtx, "store", trace.WithAttributes(
		attribute.String("contract", x.addressesToString(x.indexerConfig.Addresses)),
		attribute.String("tx", log.TxHash.Hex()),
		attribute.String("block", fmt.Sprintf("%d", log.BlockNumber)),
	))

	defer func() {
		metrics.EndSpanWithErr(span, err)
	}()

	b := &backoff.Backoff{
		Factor: 2,
		Jitter: true,
		Min:    3 * time.Millisecond,
		Max:    2 * time.Second,
	}

	timeout := time.Duration(0)
	tryCount := 0

	var tx *txData
	hasTX := true

OUTER:
	for {
		select {
		case <-ctx.Done():
			// LogEvent(ErrorLevel, "Context canceled while storing logs/receipts", LogData{"cid": x.indexerConfig.ChainID, "bn": log.BlockNumber, "tx": log.TxHash.Hex(), "la": log.Address.String(), "ca": x.addressesToString(x.indexerConfig.Addresses), "e": ctx.Err()})

			return fmt.Errorf("context canceled while storing logs/receipts: %w", ctx.Err())
		case <-time.After(timeout):
			tryCount++

			tx, err = x.fetchEventData(ctx, log.TxHash, log.BlockNumber)
			if err != nil {
				if errors.Is(err, errNoContinue) {
					return nil
				}

				if errors.Is(err, errNoTx) {
					hasTX = false
					break OUTER
				}

				if tryCount > retryTolerance {
					return fmt.Errorf("retry tolerance exceeded: %w", err)
				}

				timeout = b.Duration()
				continue
			}

			break OUTER
		}
	}

	g, groupCtx := errgroup.WithContext(ctx)
	g.Go(func() error {
		// Store receipt in the EventDB.
		if x.toTip {
			err = x.eventDB.StoreReceiptAtHead(groupCtx, x.indexerConfig.ChainID, tx.receipt)
		} else {
			err = x.eventDB.StoreReceipt(groupCtx, x.indexerConfig.ChainID, tx.receipt)
		}
		if err != nil {
			// LogEvent(ErrorLevel, "Could not store receipt, retrying", LogData{"cid": x.indexerConfig.ChainID, "bn": log.BlockNumber, "tx": log.TxHash.Hex(), "la": log.Address.String(), "ca": x.addressesToString(x.indexerConfig.Addresses), "e": err.Error()})

			return fmt.Errorf("could not store receipt: %w", err)
		}
		return nil
	})

	if hasTX {
		g.Go(func() error {
			if x.toTip {
				err = x.eventDB.StoreEthTxAtHead(groupCtx, &tx.transaction, x.indexerConfig.ChainID, log.BlockHash, log.BlockNumber, uint64(log.TxIndex))
			} else {
				err = x.eventDB.StoreEthTx(groupCtx, &tx.transaction, x.indexerConfig.ChainID, log.BlockHash, log.BlockNumber, uint64(log.TxIndex))
			}
			if err != nil {
				return fmt.Errorf("could not store tx: %w", err)
			}
			return nil
		})
	}

	g.Go(func() error {
		logs, err := x.prunedReceiptLogs(tx.receipt)
		if err != nil {
			return err
		}
		if x.toTip {
			err = x.eventDB.StoreLogsAtHead(groupCtx, x.indexerConfig.ChainID, logs...)
		} else {
			err = x.eventDB.StoreLogs(groupCtx, x.indexerConfig.ChainID, logs...)
		}
		if err != nil {
			return fmt.Errorf("could not store receipt logs: %w", err)
		}

		return nil
	})

	g.Go(func() error {
		err := x.eventDB.StoreBlockTime(groupCtx, x.indexerConfig.ChainID, tx.blockHeader.Number.Uint64(), tx.blockHeader.Time)
		if err != nil {
			return fmt.Errorf("could not store receipt logs: %w", err)
		}
		return nil
	})

	err = g.Wait()
	if err != nil {
		// LogEvent(ErrorLevel, "Could not store data", LogData{"cid": x.indexerConfig.ChainID, "bn": log.BlockNumber, "tx": log.TxHash.Hex(), "la": log.Address.String(), "ca": x.addressesToString(x.indexerConfig.Addresses), "e": err.Error()})

		return fmt.Errorf("could not store data: %w\n%s on chain %d from %d to %s", err, x.addressesToString(x.indexerConfig.Addresses), x.indexerConfig.ChainID, log.BlockNumber, log.TxHash.String())
	}

	x.cache.Add(log.TxHash, true)
	// LogEvent(InfoLevel, "Log, Receipt, and Tx stored", LogData{"cid": x.indexerConfig.ChainID, "bn": log.BlockNumber, "tx": log.TxHash.Hex(), "la": log.Address.String(), "ca": x.addressesToString(x.indexerConfig.Addresses), "ts": time.Since(startTime).Seconds()})

	return nil
}
func (x *Indexer) getLogs(parentCtx context.Context, startHeight, endHeight uint64) (<-chan types.Log, <-chan string) {
	ctx, span := x.handler.Tracer().Start(parentCtx, "getLogs")
	defer metrics.EndSpan(span)

	logFetcher := NewLogFetcher(x.client[0], big.NewInt(int64(startHeight)), big.NewInt(int64(endHeight)), &x.indexerConfig)
	logsChan, errChan := make(chan types.Log), make(chan string)

	go x.runFetcher(ctx, logFetcher, errChan)
	go x.processLogs(ctx, logFetcher, logsChan, errChan)

	return logsChan, errChan
}

func (x *Indexer) runFetcher(ctx context.Context, logFetcher *LogFetcher, errChan chan<- string) {
	if err := logFetcher.Start(ctx); err != nil {
		select {
		case <-ctx.Done():
			errChan <- fmt.Sprintf("context canceled while appending log to channel %v", ctx.Err())
			return
		case errChan <- err.Error():
			return
		}
	}
}

func (x *Indexer) processLogs(ctx context.Context, logFetcher *LogFetcher, logsChan chan<- types.Log, errChan chan<- string) {
	for {
		select {
		case <-ctx.Done():
			errChan <- fmt.Sprintf("context canceled %v", ctx.Err())
			return
		case logChunks, ok := <-logFetcher.fetchedLogsChan:
			if !ok {
				close(logsChan)
				return
			}
			for _, log := range logChunks {
				select {
				case <-ctx.Done():
					errChan <- fmt.Sprintf("context canceled while loading log chunks to log %v", ctx.Err())
					return
				case logsChan <- log:
				}
			}
		}
	}
}

// prunedReceiptLogs gets all logs from a receipt and prunes null logs.
func (x *Indexer) prunedReceiptLogs(receipt types.Receipt) (logs []types.Log, err error) {
	for i := range receipt.Logs {
		log := receipt.Logs[i]
		if log == nil {
			// LogEvent(ErrorLevel, "log is nil", LogData{"cid": x.indexerConfig.ChainID, "bn": log.BlockNumber, "tx": log.TxHash.Hex(), "la": log.Address.String(), "ca": x.addressesToString(x.indexerConfig.Addresses)})

			return nil, fmt.Errorf("log is nil\nChain: %d\nTxHash: %s\nLog BlockNumber: %d\nLog 's Contract Address: %s\nContract Address: %s", x.indexerConfig.ChainID, log.TxHash.String(), log.BlockNumber, log.Address.String(), x.addressesToString(x.indexerConfig.Addresses))
		}
		logs = append(logs, *log)
	}
	return logs, nil
}

// fetchEventData tries to fetch a transaction from the cache, if it's not there it tries to fetch it from the database.
// nolint: cyclop
func (x *Indexer) fetchEventData(parentCtx context.Context, txhash common.Hash, blockNumber uint64) (tx *txData, err error) {
	ctx, span := x.handler.Tracer().Start(parentCtx, "fetchEventData", trace.WithAttributes(
		attribute.String("tx", txhash.Hex()),
		attribute.String("block", fmt.Sprintf("%d", blockNumber)),
	))

	defer func() {
		metrics.EndSpanWithErr(span, err)
	}()

OUTER:
	// increasing this across more clients puts too much load on the server, results in failed requests. TODO investigate
	for i := range x.client[0:1] {
		tx = &txData{}

		calls := make([]w3types.Caller, 3)

		// setup referencable indexes so we can access errors from the calls
		const (
			receiptIndex = 0
			txIndex      = 1
			headerIndex  = 2
		)

		// get the transaction receipt
		calls[receiptIndex] = eth.TxReceipt(txhash).Returns(&tx.receipt)

		// get the raw transaction
		calls[txIndex] = eth.Tx(txhash).Returns(&tx.transaction)

		// get the block number
		calls[headerIndex] = eth.HeaderByNumber(new(big.Int).SetUint64(blockNumber)).Returns(&tx.blockHeader)

		//nolint: nestif
		if err := x.client[i].BatchWithContext(ctx, calls...); err != nil {
			//nolint: errorlint
			callErr, ok := err.(w3.CallErrors)
			if !ok {
				return nil, fmt.Errorf("could not parse errors: %w", err)
			}

			if callErr[receiptIndex] != nil {
				if callErr[receiptIndex].Error() == txNotFoundError {
					// LogEvent(InfoLevel, "Could not get tx for txHash, attempting with additional confirmations", LogData{"cid": x.indexerConfig.ChainID, "tx": txhash, "ca": x.addressesToString(x.indexerConfig.Addresses), "e": err.Error()})
					continue OUTER
				}
			}

			if callErr[txIndex] != nil {
				switch callErr[txIndex].Error() {
				case txNotSupportedError:
					// LogEvent(InfoLevel, "Invalid tx", LogData{"cid": x.indexerConfig.ChainID, "tx": txhash, "ca": x.addressesToString(x.indexerConfig.Addresses), "e": err.Error()})
					return tx, errNoTx
				case invalidTxVRSError:
					// LogEvent(InfoLevel, "Could not get tx for txHash, attempting with additional confirmations", LogData{"cid": x.indexerConfig.ChainID, "tx": txhash, "ca": x.addressesToString(x.indexerConfig.Addresses), "e": err.Error()})
					return tx, errNoTx
				case txNotFoundError:
					// LogEvent(InfoLevel, "Could not get tx for txHash, attempting with additional confirmations", LogData{"cid": x.indexerConfig.ChainID, "tx": txhash, "ca": x.addressesToString(x.indexerConfig.Addresses), "e": err.Error()})
					continue OUTER
				}
			}

			return nil, fmt.Errorf("could not get tx receipt: %w", err)
		}

		tx.success = true
	}

	if tx == nil || !tx.success {
		return nil, fmt.Errorf("could not get tx data: %w", err)
	}

	return tx, nil
}

func (x *Indexer) addressesToString(addresses []common.Address) string {
	var output string
	for i := range addresses {
		if i == 0 {
			output = addresses[i].String()
		} else {
			output = output + "," + addresses[i].String()
		}
	}
	return output
}