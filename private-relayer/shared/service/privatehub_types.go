package service

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/contracts/EndpointV1"
)

type PrivateHubMessage struct {
	ID string

	MessageID common.Hash
	From      common.Address
	To        common.Address
	Data      EndpointV1.RaylsMessage
}

func (m PrivateHubMessage) GetID() string {
	return m.ID
}
