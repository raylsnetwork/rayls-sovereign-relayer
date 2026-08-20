package testdata

import (
	"math/big"

	"github.com/raylsnetwork/rayls-sovereign-relayer/private-relayer/repository"
	"github.com/raylsnetwork/rayls-sovereign-relayer/types"
)

var (
	EnygmaHistory1 = types.EnygmaHistory{
		ResourceId:            "resource-history-1",
		FromChainId:           big.NewInt(10),
		BlockNumberPrivateHub: big.NewInt(1001),
		RFactor:               big.NewInt(777),
		BalanceChange:         big.NewInt(500),
		EventType:             types.EnygmaMint,
	}
	EnygmaHistory2 = types.EnygmaHistory{
		ResourceId:            "resource-history-2",
		FromChainId:           big.NewInt(20),
		BlockNumberPrivateHub: big.NewInt(2002),
		RFactor:               big.NewInt(888),
		BalanceChange:         big.NewInt(300),
		EventType:             types.EnygmaBurn,
	}
	EnygmaHistory3SameResource = types.EnygmaHistory{
		ResourceId:            "resource-history-1",
		FromChainId:           big.NewInt(10),
		BlockNumberPrivateHub: big.NewInt(3003),
		RFactor:               big.NewInt(999),
		BalanceChange:         big.NewInt(200),
		EventType:             types.EnygmaMint,
	}
)

var (
	ModelEnygmaHistory1 = repository.EnygmaHistory{
		ResourceId:            EnygmaHistory1.ResourceId,
		FromChainId:           EnygmaHistory1.FromChainId.String(),
		BalanceChange:         EnygmaHistory1.BalanceChange.String(),
		RFactor:               EnygmaHistory1.RFactor.String(),
		BlockNumberPrivateHub: EnygmaHistory1.BlockNumberPrivateHub.Uint64(),
		EventType:             uint8(EnygmaHistory1.EventType),
	}
	ModelEnygmaHistory2 = repository.EnygmaHistory{
		ResourceId:            EnygmaHistory2.ResourceId,
		FromChainId:           EnygmaHistory2.FromChainId.String(),
		BalanceChange:         EnygmaHistory2.BalanceChange.String(),
		RFactor:               EnygmaHistory2.RFactor.String(),
		BlockNumberPrivateHub: EnygmaHistory2.BlockNumberPrivateHub.Uint64(),
		EventType:             uint8(EnygmaHistory2.EventType),
	}
	ModelEnygmaHistory3SameResource = repository.EnygmaHistory{
		ResourceId:            EnygmaHistory3SameResource.ResourceId,
		FromChainId:           EnygmaHistory3SameResource.FromChainId.String(),
		BalanceChange:         EnygmaHistory3SameResource.BalanceChange.String(),
		RFactor:               EnygmaHistory3SameResource.RFactor.String(),
		BlockNumberPrivateHub: EnygmaHistory3SameResource.BlockNumberPrivateHub.Uint64(),
		EventType:             uint8(EnygmaHistory3SameResource.EventType),
	}
)
