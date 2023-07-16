package scribe_test

import (
	"context"
	"fmt"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/params"
	. "github.com/stretchr/testify/assert"
	"github.com/synapsecns/sanguine/ethergo/backends/geth"
	"github.com/synapsecns/sanguine/services/scribe/backend"
	"github.com/synapsecns/sanguine/services/scribe/config"
	"github.com/synapsecns/sanguine/services/scribe/db"
	"github.com/synapsecns/sanguine/services/scribe/scribe"
	"github.com/synapsecns/sanguine/services/scribe/scribe/indexer"
	"github.com/synapsecns/sanguine/services/scribe/testutil"
	"math/big"
	"time"
)

// TestIndexToBlock tests using a contractBackfiller for recording receipts and logs in a database.
func (s *ScribeSuite) TestIndexToBlock() {
	// Get simulated blockchain, deploy the test contract, and set up test variables.
	simulatedChain := geth.NewEmbeddedBackendForChainID(s.GetSuiteContext(), s.T(), big.NewInt(142))
	simulatedClient, err := backend.DialBackend(s.GetTestContext(), simulatedChain.RPCAddress(), s.nullMetrics)
	Nil(s.T(), err)

	simulatedChain.FundAccount(s.GetTestContext(), s.wallet.Address(), *big.NewInt(params.Ether))
	testContract, testRef := s.manager.GetTestContract(s.GetTestContext(), simulatedChain)
	transactOpts := simulatedChain.GetTxContext(s.GetTestContext(), nil)

	// Set config.
	contractConfig := config.ContractConfig{
		Address:    testContract.Address().String(),
		StartBlock: 0,
	}

	simulatedChainArr := []backend.ScribeBackend{simulatedClient, simulatedClient}
	chainConfig := config.ChainConfig{
		ChainID:              142,
		GetLogsBatchAmount:   1,
		Confirmations:        0,
		StoreConcurrency:     1,
		GetLogsRange:         1,
		ConcurrencyThreshold: 100,
		Contracts:            []config.ContractConfig{contractConfig},
	}

	chainIndexer, err := scribe.NewChainIndexer(s.testDB, simulatedChainArr, chainConfig, s.nullMetrics)
	Nil(s.T(), err)

	// Emit events for the backfiller to read.
	tx, err := testRef.EmitEventA(transactOpts.TransactOpts, big.NewInt(1), big.NewInt(2), big.NewInt(3))
	Nil(s.T(), err)
	simulatedChain.WaitForConfirmation(s.GetTestContext(), tx)

	tx, err = testRef.EmitEventA(transactOpts.TransactOpts, big.NewInt(1), big.NewInt(2), big.NewInt(3))
	Nil(s.T(), err)

	simulatedChain.WaitForConfirmation(s.GetTestContext(), tx)
	tx, err = testRef.EmitEventB(transactOpts.TransactOpts, []byte{4}, big.NewInt(5), big.NewInt(6))
	Nil(s.T(), err)
	simulatedChain.WaitForConfirmation(s.GetTestContext(), tx)

	// Emit two logs in one receipt.
	tx, err = testRef.EmitEventAandB(transactOpts.TransactOpts, big.NewInt(7), big.NewInt(8), big.NewInt(9))
	Nil(s.T(), err)

	simulatedChain.WaitForConfirmation(s.GetTestContext(), tx)

	// Get the block that the last transaction was executed in.
	txBlockNumber, err := testutil.GetTxBlockNumber(s.GetTestContext(), simulatedChain, tx)
	Nil(s.T(), err)

	// TODO use no-op meter
	blockHeightMeter, err := s.nullMetrics.Meter().NewHistogram(fmt.Sprint("scribe_block_meter", chainConfig.ChainID), "block_histogram", "a block height meter", "blocks")
	Nil(s.T(), err)

	contracts := []common.Address{common.HexToAddress(contractConfig.Address)}
	indexer, err := indexer.NewIndexer(chainConfig, contracts, s.testDB, simulatedChainArr, s.nullMetrics, blockHeightMeter)
	Nil(s.T(), err)

	err = chainIndexer.IndexToBlock(s.GetTestContext(), nil, uint64(0), indexer)
	Nil(s.T(), err)

	// Get all receipts.
	receipts, err := s.testDB.RetrieveReceiptsWithFilter(s.GetTestContext(), db.ReceiptFilter{}, 1)
	Nil(s.T(), err)

	// Check to see if 3 receipts were collected.
	Equal(s.T(), 4, len(receipts))

	// Get all logs.
	logs, err := s.testDB.RetrieveLogsWithFilter(s.GetTestContext(), db.LogFilter{}, 1)
	Nil(s.T(), err)

	// Check to see if 4 logs were collected.
	Equal(s.T(), 5, len(logs))

	// Check to see if the last receipt has two logs.
	Equal(s.T(), 2, len(receipts[0].Logs))

	// Ensure last indexed block is correct.
	lastIndexed, err := s.testDB.RetrieveLastIndexed(s.GetTestContext(), testContract.Address(), uint32(testContract.ChainID().Uint64()))
	Nil(s.T(), err)
	Equal(s.T(), txBlockNumber, lastIndexed)
}

// TestChainIndexer tests that the ChainIndexer can backfill events from a chain.
func (s *ScribeSuite) TestChainIndexer() {
	const numberOfContracts = 3
	const desiredBlockHeight = 20
	chainID := gofakeit.Uint32()
	chainBackends := make(map[uint32]geth.Backend)
	newBackend := geth.NewEmbeddedBackendForChainID(s.GetTestContext(), s.T(), big.NewInt(int64(chainID)))
	chainBackends[chainID] = *newBackend

	// Create contract managers
	managers := []*testutil.DeployManager{s.manager}
	if numberOfContracts > 1 {
		for i := 1; i < numberOfContracts; i++ {
			managers = append(managers, testutil.NewDeployManager(s.T()))
		}
	}

	testChainHandlerMap, chainBackendMap, err := testutil.PopulateChainsWithLogs(s.GetTestContext(), s.T(), chainBackends, desiredBlockHeight, managers, s.nullMetrics)
	Nil(s.T(), err)

	var contractConfigs []config.ContractConfig
	addresses := testChainHandlerMap[chainID].Addresses
	for i := range addresses {
		contractConfig := config.ContractConfig{
			Address: addresses[i].String(),
		}
		contractConfigs = append(contractConfigs, contractConfig)
	}
	chainConfig := config.ChainConfig{
		ChainID:            chainID,
		Confirmations:      0,
		GetLogsBatchAmount: 1,
		StoreConcurrency:   1,
		GetLogsRange:       1,
		Contracts:          contractConfigs,
	}
	killableContext, cancel := context.WithTimeout(s.GetTestContext(), 20*time.Second)
	defer cancel()
	chainIndexer, err := scribe.NewChainIndexer(s.testDB, chainBackendMap[chainID], chainConfig, s.nullMetrics)
	Nil(s.T(), err)
	_ = chainIndexer.Index(killableContext, nil)
	sum := uint64(0)
	for _, value := range testChainHandlerMap[chainID].EventsEmitted {
		sum += value
	}
	logs, err := s.testDB.RetrieveLogsWithFilter(s.GetTestContext(), db.LogFilter{}, 1)
	Nil(s.T(), err)
	Equal(s.T(), sum, uint64(len(logs)))
	receipts, err := s.testDB.RetrieveReceiptsWithFilter(s.GetTestContext(), db.ReceiptFilter{}, 1)
	Nil(s.T(), err)
	Equal(s.T(), sum, uint64(len(receipts)))
}

// TestChainIndexerLivefill tests a ChainIndexer's ablity to livefill and handle passing events from backfill to livefill.
//
// nolint:cyclop
func (s *ScribeSuite) TestChainIndexerLivefill() {
	const numberOfContracts = 5
	currentBlockHeight := uint64(0) // starting with zero to emit events while indexing.
	chainID := gofakeit.Uint32()
	chainBackends := make(map[uint32]geth.Backend)
	newBackend := geth.NewEmbeddedBackendForChainID(s.GetTestContext(), s.T(), big.NewInt(int64(chainID)))
	chainBackends[chainID] = *newBackend

	// Create contract managers
	deployManagers := []*testutil.DeployManager{s.manager}
	if numberOfContracts > 1 {
		for i := 1; i < numberOfContracts; i++ {
			deployManagers = append(deployManagers, testutil.NewDeployManager(s.T()))
		}
	}

	testChainHandlerMap, chainBackendMap, err := testutil.PopulateChainsWithLogs(s.GetTestContext(), s.T(), chainBackends, currentBlockHeight, deployManagers, s.nullMetrics)
	Nil(s.T(), err)
	addresses := testChainHandlerMap[chainID].Addresses
	// Differing start blocks and refresh rates to test contracts reaching livefill at different times.
	contractConfig1 := config.ContractConfig{
		Address:     addresses[0].String(),
		StartBlock:  0,
		RefreshRate: 4,
	}
	contractConfig2 := config.ContractConfig{
		Address:     addresses[1].String(),
		StartBlock:  25,
		RefreshRate: 1,
	}
	contractConfig3 := config.ContractConfig{
		Address:     addresses[2].String(),
		StartBlock:  30,
		RefreshRate: 3,
	}
	contractConfig4 := config.ContractConfig{
		Address:     addresses[3].String(),
		StartBlock:  30,
		RefreshRate: 1,
	}
	contractConfig5 := config.ContractConfig{
		Address:     addresses[4].String(),
		StartBlock:  0,
		RefreshRate: 3,
	}

	contractConfigs := []config.ContractConfig{contractConfig1, contractConfig2, contractConfig3, contractConfig4, contractConfig5}
	chainConfig := config.ChainConfig{
		ChainID:            chainID,
		Confirmations:      0,
		GetLogsBatchAmount: 1,
		StoreConcurrency:   1,
		GetLogsRange:       1,
		// livefill threshold kept small to ensure that the indexer does not reach the head before the continuous event emitting starts
		LivefillThreshold: 0,
		Contracts:         contractConfigs,
	}

	// Update start blocks
	for i := range contractConfigs {
		contract := contractConfigs[i]
		contractAddress := common.HexToAddress(contract.Address)
		testChainHandlerMap[chainID].ContractStartBlocks[contractAddress] = contract.StartBlock
	}

	chainIndexer, err := scribe.NewChainIndexer(s.testDB, chainBackendMap[chainID], chainConfig, s.nullMetrics)
	Nil(s.T(), err)
	Equal(s.T(), 0, len(chainIndexer.GetLivefillContracts()))
	currentBlockHeight = 30

	emittingContext, cancelEmitting := context.WithTimeout(s.GetTestContext(), 60*time.Second)
	defer cancelEmitting()

	// Emit an event for every contract every second. This will terminate 10 seconds before indexing terminates.
	go func() {
		for {
			select {
			case <-emittingContext.Done():
				return
			case <-time.After(1 * time.Second):
				currentBlockHeight += 2
				emitErr := testutil.EmitEvents(s.GetTestContext(), s.T(), newBackend, currentBlockHeight, testChainHandlerMap[chainID])
				Nil(s.T(), emitErr)
			}
		}
	}()

	<-time.After(40 * time.Second) // wait for 200 seconds before indexing to get some events on chain before indexing.

	// Cap indexing for 60 seconds.
	indexingContext, cancelIndexing := context.WithTimeout(s.GetTestContext(), 30*time.Second)
	defer cancelIndexing()

	// Check that the number of livefill contracts is correct.
	numberLivefillContracts := 0
	go func() {
		currentLength := 0
		for {
			select {
			case <-indexingContext.Done():
				return
			default:
				contracts := chainIndexer.GetLivefillContracts()
				if currentLength != len(contracts) {
					currentLength = len(contracts)
					newContract := contracts[currentLength-1]

					lastIndexed, indexErr := s.testDB.RetrieveLastIndexed(s.GetTestContext(), common.HexToAddress(newContract.Address), chainID)
					Nil(s.T(), indexErr)
					numberLivefillContracts = len(contracts)
					currentBlock, indexErr := newBackend.BlockNumber(s.GetTestContext())
					Nil(s.T(), indexErr)
					GreaterOrEqual(s.T(), int(lastIndexed), int(currentBlock)-int(chainConfig.LivefillThreshold))
				}
			}
		}
	}()

	// Index events
	_ = chainIndexer.Index(indexingContext, nil)

	<-indexingContext.Done()
	sum := uint64(0)
	for _, value := range testChainHandlerMap[chainID].EventsEmitted {
		sum += value
	}

	logs, err := testutil.GetLogsUntilNoneLeft(s.GetTestContext(), s.testDB, db.LogFilter{})
	Nil(s.T(), err)
	Equal(s.T(), sum, uint64(len(logs)))
	receipts, err := testutil.GetReceiptsUntilNoneLeft(s.GetTestContext(), s.testDB, db.ReceiptFilter{})
	Nil(s.T(), err)
	Equal(s.T(), sum, uint64(len(receipts)))
	Equal(s.T(), numberOfContracts, numberLivefillContracts)
}
