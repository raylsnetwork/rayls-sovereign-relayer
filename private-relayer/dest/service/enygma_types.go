package service

import (
	"math/big"

	"github.com/raylsnetwork/rayls-sovereign-relayer/types"
)

type EnygmaDestMessageType int

const (
	EnygmaTransferBatchMessage EnygmaDestMessageType = iota
	EnygmaFinalizedBalanceMessage
)

type EnygmaDestMessage struct {
	ID          string
	Type        EnygmaDestMessageType
	BlockNumber uint64

	TransferBatch    *types.EnygmaTransferBatch
	FinalizedBalance *types.EnygmaFinalizedBalance
}

func (m EnygmaDestMessage) GetID() string {
	return m.ID
}

type EnygmaFinalizedBalanceEvent struct {
	ResourceId           [32]byte
	FinalizedBlockNumber *big.Int
	PendingBlockNumber   *big.Int
	Balances             []EnygmaPointWithChainId
}

type EnygmaPointWithChainId struct {
	C1      *big.Int
	C2      *big.Int
	ChainId *big.Int
}
