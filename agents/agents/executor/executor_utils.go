package executor

import (
	"fmt"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	execTypes "github.com/synapsecns/sanguine/agents/agents/executor/types"
	"github.com/synapsecns/sanguine/agents/contracts/inbox"
	"github.com/synapsecns/sanguine/agents/contracts/lightinbox"
	"github.com/synapsecns/sanguine/agents/contracts/origin"
	"github.com/synapsecns/sanguine/agents/types"
)

// logToMessage converts the log to a leaf data.
func (e Executor) logToMessage(log ethTypes.Log, chainID uint32) (*types.Message, error) {
	message, ok := e.chainExecutors[chainID].originParser.ParseSent(log)
	if !ok {
		return nil, fmt.Errorf("could not parse committed message")
	}

	return &message, nil
}

// logToAttestation converts the log to an attestation.
func (e Executor) logToAttestation(log ethTypes.Log, chainID uint32) (*types.Attestation, error) {
	attestation, ok := (*e.chainExecutors[chainID].lightInboxParser).ParseAttestationAccepted(log)
	if !ok {
		return nil, fmt.Errorf("could not parse attestation")
	}

	return &attestation, nil
}

// logToSnapshot converts the log to a snapshot.
func (e Executor) logToSnapshot(log ethTypes.Log, chainID uint32) (*types.Snapshot, error) {
	snapshot, domain, ok := (*e.chainExecutors[chainID].inboxParser).ParseSnapshotAccepted(log)
	if !ok {
		return nil, fmt.Errorf("could not parse snapshot")
	}

	if domain == 0 {
		//nolint:nilnil
		return nil, nil
	}

	return &snapshot, nil
}

// logType determines whether a log is a `Sent` from Origin.sol or `AttestationAccepted` from Destination.sol.
// TODO: Clean with switch case.
func (e Executor) logType(log ethTypes.Log, chainID uint32) execTypes.ContractType {
	contractType := execTypes.Other

	//nolint:nestif
	if e.chainExecutors[chainID].inboxParser != nil {
		if summitEvent, ok := (*e.chainExecutors[chainID].inboxParser).EventType(log); ok && summitEvent == inbox.SnapshotAcceptedEvent {
			contractType = execTypes.InboxContract
		}

		return contractType
	}

	//nolint:nestif
	if originEvent, ok := e.chainExecutors[chainID].originParser.EventType(log); ok && originEvent == origin.SentEvent {
		contractType = execTypes.OriginContract
	} else if lightManagerEvent, ok := (*e.chainExecutors[chainID].lightInboxParser).EventType(log); ok && lightManagerEvent == lightinbox.AttestationAcceptedEvent {
		contractType = execTypes.LightInboxContract
	}

	return contractType
}
