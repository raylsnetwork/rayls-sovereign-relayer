package service

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/contracts/EndpointV1"
)

type CrossChainMessage struct {
	ID string

	MessageID [32]byte
	From      common.Address
	ToChainID *big.Int
	To        common.Address
	Data      EndpointV1.RaylsMessage

	BlockHash common.Hash
	TxHash    common.Hash

	BlockNumber uint64
	TxIdx       uint
	LogIdx      uint

	TokenAddress common.Address
}

func (m CrossChainMessage) GetID() string {
	return m.ID
}
