package service

import (
	"encoding/json"
	"fmt"
)

type DvpEventType int

const (
	Dvp721CreationEvent DvpEventType = iota
	Dvp721MintEvent
	Dvp721BurnEvent
	Dvp721DepositIntoDvpEvent
	Dvp721WithdrawFromDvpEvent
	Dvp721SwapForEnygmaEvent

	Dvp1155CreationEvent
	Dvp1155MintEvent
	Dvp1155BurnEvent
	Dvp1155DepositIntoDvpEvent
	Dvp1155WithdrawFromDvpEvent
	Dvp1155SwapForEnygmaEvent

	DvpSwapCancelledEvent

	// Enygma swap events
	DvpEnygmaSwapERC721Event
	DvpEnygmaSwapERC1155Event

	// sentinel element equal to the count of events
	DvpEventTypeCount
)

type DvpEventBatchBuilder map[DvpEventType]DvpEventBatch

type DvpEventBatch interface {
	Serialize() (DvpSerializedEventBatch, error)
	PushEvent(any) error
}
type DvpSerializedEventBatch struct {
	BlockNumber      uint64
	Type             DvpEventType
	SerializedEvents []byte
}

func (s DvpSerializedEventBatch) GetID() string {
	return fmt.Sprintf("%d-%d", s.BlockNumber, s.Type)
}

type DvpTypedEventBatch[T any] struct {
	BlockNumber uint64
	Type        DvpEventType
	Events      []T
}

func (s DvpTypedEventBatch[T]) GetID() string {
	return fmt.Sprintf("%d-%d", s.BlockNumber, s.Type)
}

func (s DvpTypedEventBatch[T]) Serialize() (DvpSerializedEventBatch, error) {
	serializedEvents, err := json.Marshal(s.Events)
	if err != nil {
		return DvpSerializedEventBatch{}, fmt.Errorf("failed to serialzie events")
	}
	return DvpSerializedEventBatch{
		BlockNumber:      s.BlockNumber,
		Type:             s.Type,
		SerializedEvents: serializedEvents,
	}, nil
}

func (s *DvpTypedEventBatch[T]) PushEvent(ev any) error {
	typedEv, ok := ev.(T)
	if !ok {
		return fmt.Errorf("wrong event type")
	}
	s.Events = append(s.Events, typedEv)
	return nil
}
