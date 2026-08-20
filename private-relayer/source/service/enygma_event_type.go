package service

import (
	"encoding/json"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
)

type EnygmaEventType int

// Currently the order of execution of enygma events
// is set by the order of events in this enumeration.
// !!! DO NOT CHANGE UNLESS !!!
// INTENTIONALLY REORDERING EVENT EXECUTION
const (
	EnygmaCreationEvent EnygmaEventType = iota
	EnygmaSupplyUpdateEvent
	EnygmaTransferEvent
	EnygmaDepositEvent
	EnygmaWithdrawEvent

	// sentinel element equal to the count of events
	EnygmaEventTypeCount
)

// EnygmaSerializedEvent represents a single Enygma event for individual message queue publishing
type EnygmaSerializedEvent struct {
	Id              string
	BlockNumber     uint64
	LogIndex        uint
	TxHash          common.Hash
	ResourceID      string
	Type            EnygmaEventType
	SerializedEvent []byte
}

func (e EnygmaSerializedEvent) GetID() string {
	return e.Id
}

// EnygmaTypedEvent wraps an event with log metadata for ID generation
type EnygmaTypedEvent[T any] struct {
	Id          string
	BlockNumber uint64
	LogIndex    uint
	TxHash      common.Hash
	ResourceID  string
	Type        EnygmaEventType
	Event       T
}

func NewEnygmaTypedEvent[T any](
	blockNumber uint64,
	logIndex uint,
	txHash common.Hash,
	resourceID string,
	eventType EnygmaEventType,
	event T,
) EnygmaTypedEvent[T] {
	var id string

	switch e := any(event).(type) {
	case *EnygmaTransferTx:
		// We can have many transfers in one event, so we need to make sure each one of them has unique ID.
		id = fmt.Sprintf("%s-%d-%d", txHash.Hex(), logIndex, e.DestIdx)
	default:
		// Other event types use base ID (txHash-logIndex)
		id = fmt.Sprintf("%s-%d", txHash.Hex(), logIndex)
	}

	return EnygmaTypedEvent[T]{
		Id:          id,
		BlockNumber: blockNumber,
		LogIndex:    logIndex,
		TxHash:      txHash,
		ResourceID:  resourceID,
		Type:        eventType,
		Event:       event,
	}
}

func (e EnygmaTypedEvent[T]) GetID() string {
	return e.Id
}

func (e EnygmaTypedEvent[T]) Serialize() (EnygmaSerializedEvent, error) {
	serialized, err := json.Marshal(e.Event)
	if err != nil {
		return EnygmaSerializedEvent{}, fmt.Errorf("failed to serialize event: %w", err)
	}

	return EnygmaSerializedEvent{
		Id:              e.GetID(),
		BlockNumber:     e.BlockNumber,
		LogIndex:        e.LogIndex,
		TxHash:          e.TxHash,
		ResourceID:      e.ResourceID,
		Type:            e.Type,
		SerializedEvent: serialized,
	}, nil
}
