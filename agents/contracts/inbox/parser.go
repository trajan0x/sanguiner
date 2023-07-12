package inbox

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/synapsecns/sanguine/agents/types"
)

// Parser parses events from the inbox contract.
type Parser interface {
	// EventType is the event type.
	EventType(log ethTypes.Log) (_ EventType, ok bool)
	// ParseSnapshotAccepted parses a SnapshotAccepted event.
	ParseSnapshotAccepted(log ethTypes.Log) (_ types.Snapshot, domain uint32, agentSig []byte, ok bool)
}

type parserImpl struct {
	// filterer is the parser filterer we use to parse events
	filterer *InboxFilterer
}

// NewParser creates a new parser for the inbox contract.
func NewParser(inboxAddress common.Address) (Parser, error) {
	parser, err := NewInboxFilterer(inboxAddress, nil)
	if err != nil {
		return nil, fmt.Errorf("could not create %T: %w", InboxFilterer{}, err)
	}

	return &parserImpl{filterer: parser}, nil
}

func (p parserImpl) EventType(log ethTypes.Log) (_ EventType, ok bool) {
	for _, logTopic := range log.Topics {
		eventType := eventTypeFromTopic(logTopic)
		if eventType == nil {
			continue
		}

		return *eventType, true
	}
	// return an unknown event to avoid cases where user failed to check the event type
	return EventType(len(AllEventTypes) + 2), false
}

// ParseSnapshotAccepted parses a SnapshotAccepted event.
func (p parserImpl) ParseSnapshotAccepted(log ethTypes.Log) (_ types.Snapshot, domain uint32, agentSig []byte, ok bool) {
	fmt.Printf("ParseSnapshotAccepted: %v\n", log)
	inboxSnapshot, err := p.filterer.ParseSnapshotAccepted(log)
	if err != nil {
		return nil, 0, nil, false
	}

	snapshot, err := types.DecodeSnapshot(inboxSnapshot.SnapPayload)
	if err != nil {
		return nil, 0, nil, false
	}
	fmt.Printf("inboxSnapshot: %v\n", inboxSnapshot)
	fmt.Printf("inboxSnapshot.Signature: %v\n", inboxSnapshot.SnapSignature)

	return snapshot, inboxSnapshot.Domain, inboxSnapshot.SnapSignature, true
}

// EventType is the type of the summit events
//
//go:generate go run golang.org/x/tools/cmd/stringer -type=EventType
type EventType uint

const (
	// SnapshotAcceptedEvent is a SnapshotAccepted event.
	SnapshotAcceptedEvent EventType = iota
)

// Int gets the int for an event type.
func (i EventType) Int() uint8 {
	return uint8(i)
}

// AllEventTypes contains all event types.
var AllEventTypes = []EventType{SnapshotAcceptedEvent}
