package service

import (
	"encoding/json"
	"fmt"
)

type EnygmaDvpEventType int

const (
	EnygmaDvp721CreationEvent EnygmaDvpEventType = iota
	EnygmaDvp721MintEvent
	EnygmaDvp721BurnEvent
	EnygmaDvp721DepositIntoEnygmaDvpEvent
	EnygmaDvp721WithdrawFromEnygmaDvpEvent
	EnygmaDvp721SwapForEnygmaEvent

	EnygmaDvp1155CreationEvent
	EnygmaDvp1155MintEvent
	EnygmaDvp1155BurnEvent
	EnygmaDvp1155DepositIntoEnygmaDvpEvent
	EnygmaDvp1155WithdrawFromEnygmaDvpEvent
	EnygmaDvp1155SwapForEnygmaEvent

	// Enygma swap events
	EnygmaDvpEnygmaSwapERC721Event
	EnygmaDvpEnygmaSwapERC1155Event

	EnygmaDvpSwapCancelledEvent

	// sentinel element equal to the count of events
	EnygmaDvpEventTypeCount
)

type EnygmaDvpEventBatchBuilder map[EnygmaDvpEventType]EnygmaDvpEventBatch

type EnygmaDvpEventBatch interface {
	Serialize() (EnygmaDvpSerializedEventBatch, error)
	PushEvent(any) error
}
type EnygmaDvpSerializedEventBatch struct {
	BlockNumber      uint64
	Type             EnygmaDvpEventType
	SerializedEvents []byte
}

func (s EnygmaDvpSerializedEventBatch) GetID() string {
	return fmt.Sprintf("%d-%d", s.BlockNumber, s.Type)
}

type EnygmaDvpTypedEventBatch[T any] struct {
	BlockNumber uint64
	Type        EnygmaDvpEventType
	Events      []T
}

func (s EnygmaDvpTypedEventBatch[T]) GetID() string {
	return fmt.Sprintf("%d-%d", s.BlockNumber, s.Type)
}

func (s EnygmaDvpTypedEventBatch[T]) Serialize() (EnygmaDvpSerializedEventBatch, error) {
	serializedEvents, err := json.Marshal(s.Events)
	if err != nil {
		return EnygmaDvpSerializedEventBatch{}, fmt.Errorf("failed to serialzie events")
	}
	return EnygmaDvpSerializedEventBatch{
		BlockNumber:      s.BlockNumber,
		Type:             s.Type,
		SerializedEvents: serializedEvents,
	}, nil
}

func (s *EnygmaDvpTypedEventBatch[T]) PushEvent(ev any) error {
	typedEv, ok := ev.(T)
	if !ok {
		return fmt.Errorf("wrong event type")
	}
	s.Events = append(s.Events, typedEv)
	return nil
}
