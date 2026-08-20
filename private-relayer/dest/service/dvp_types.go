package service

import (
	"github.com/raylsnetwork/rayls-sovereign-relayer/types"
)

type DvpDestMessageType int

const (
	DvpCommitmentsMessage DvpDestMessageType = iota
	DvpNullifierMessage
	DvpSwapInitiatedMessage
	DvpSwapCompletedMessage
	DvpSwapCancelledMessage
	DvpSwapTimedOutMessage
)

type DvpDestMessage struct {
	ID          string
	Type        DvpDestMessageType
	BlockNumber uint64

	SharedID      string
	Commitments   *types.DvpCommitmentsData
	Nullifiers    *types.DvpNullifierData
	SwapInitiated *types.DvpSwapInitiatedData
}

func (m DvpDestMessage) GetID() string {
	return m.ID
}
