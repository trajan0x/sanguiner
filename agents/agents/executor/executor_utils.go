package executor

import (
	"context"
	"fmt"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/jpillora/backoff"
	"github.com/synapsecns/sanguine/agents/contracts/inbox"
	"github.com/synapsecns/sanguine/agents/contracts/lightinbox"
	"github.com/synapsecns/sanguine/agents/contracts/origin"
	"github.com/synapsecns/sanguine/agents/types"
	"math/big"
	"time"
)

// logToMessage converts the log to a leaf data.
func (e Executor) logToMessage(log ethTypes.Log, chainID uint32) (types.Message, error) {
	message, ok := e.chainExecutors[chainID].originParser.ParseSent(log)
	if !ok {
		return nil, fmt.Errorf("could not parse committed message")
	}

	if message == nil {
		//nolint:nilnil
		return nil, nil
	}

	return message, nil
}

func (e Executor) logToInterface(log ethTypes.Log, chainID uint32) (any, error) {
	switch {
	case e.isSnapshotAcceptedEvent(log, chainID):
		wrappedSnapshot, err := e.chainExecutors[chainID].inboxParser.ParseSnapshotAccepted(log)
		if err != nil {
			return nil, fmt.Errorf("could not parse snapshot: %w", err)
		}

		return wrappedSnapshot.Snapshot, nil
	case e.isSentEvent(log, chainID):
		return e.logToMessage(log, chainID)
	case e.isAttestationAcceptedEvent(log, chainID):
		wrappedAttestation, err := e.chainExecutors[chainID].lightInboxParser.ParseAttestationAccepted(log)
		if err != nil {
			return nil, fmt.Errorf("could not parse attestation: %w", err)
		}

		return wrappedAttestation.Attestation, nil
	default:
		//nolint:nilnil
		return nil, nil
	}
}

func (e Executor) isSnapshotAcceptedEvent(log ethTypes.Log, chainID uint32) bool {
	if e.chainExecutors[chainID].inboxParser == nil {
		return false
	}

	inboxEvent, ok := e.chainExecutors[chainID].inboxParser.EventType(log)
	return ok && inboxEvent == inbox.SnapshotAcceptedEvent
}

func (e Executor) isSentEvent(log ethTypes.Log, chainID uint32) bool {
	if e.chainExecutors[chainID].originParser == nil {
		return false
	}

	originEvent, ok := e.chainExecutors[chainID].originParser.EventType(log)
	return ok && originEvent == origin.SentEvent
}

func (e Executor) isAttestationAcceptedEvent(log ethTypes.Log, chainID uint32) bool {
	if e.chainExecutors[chainID].lightInboxParser == nil {
		return false
	}

	lightManagerEvent, ok := e.chainExecutors[chainID].lightInboxParser.EventType(log)
	return ok && lightManagerEvent == lightinbox.AttestationAcceptedEvent
}

// processMessage processes and stores a message.
func (e Executor) processMessage(ctx context.Context, message types.Message, logBlockNumber uint64) error {
	merkleIndex := e.chainExecutors[message.OriginDomain()].merkleTree.NumOfItems()
	leaf, err := message.ToLeaf()
	if err != nil {
		return fmt.Errorf("could not convert message to leaf: %w", err)
	}

	// Make sure the nonce of the message is being inserted at the right index.
	switch {
	case merkleIndex+1 > message.Nonce():
		return nil
	case merkleIndex+1 < message.Nonce():
		return fmt.Errorf("nonce is not correct. expected: %d, got: %d", merkleIndex+1, message.Nonce())
	default:
	}

	e.chainExecutors[message.OriginDomain()].merkleTree.Insert(leaf[:])

	err = e.executorDB.StoreMessage(ctx, message, logBlockNumber, false, 0)
	if err != nil {
		return fmt.Errorf("could not store message: %w", err)
	}

	return nil
}

// processAttestation processes and stores an attestation.
func (e Executor) processSnapshot(ctx context.Context, snapshot types.Snapshot, logBlockNumber uint64) error {
	snapshotRoot, proofs, err := snapshot.SnapshotRootAndProofs()
	if err != nil {
		return fmt.Errorf("could not get snapshot root and proofs: %w", err)
	}

	err = e.executorDB.StoreStates(ctx, snapshot.States(), snapshotRoot, proofs, logBlockNumber)
	if err != nil {
		return fmt.Errorf("could not store states: %w", err)
	}

	return nil
}

// processAttestation processes and stores an attestation.
func (e Executor) processAttestation(ctx context.Context, attestation types.Attestation, chainID uint32, logBlockNumber uint64) error {
	b := &backoff.Backoff{
		Factor: 2,
		Jitter: true,
		Min:    30 * time.Millisecond,
		Max:    3 * time.Second,
	}

	timeout := time.Duration(0)

	var logHeader *ethTypes.Header
	var err error

retryLoop:
	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("context canceled: %w", ctx.Err())
		case <-time.After(timeout):
			if b.Attempt() >= rpcRetry {
				return fmt.Errorf("could not get log header: %w", err)
			}
			logHeader, err = e.chainExecutors[chainID].rpcClient.HeaderByNumber(ctx, big.NewInt(int64(logBlockNumber)))
			if err != nil {
				timeout = b.Duration()

				continue
			}

			break retryLoop
		}
	}

	err = e.executorDB.StoreAttestation(ctx, attestation, chainID, logBlockNumber, logHeader.Time)
	if err != nil {
		return fmt.Errorf("could not store attestation: %w", err)
	}

	return nil
}
