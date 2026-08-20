package testdata

import (
	"math/big"

	"github.com/raylsnetwork/rayls-sovereign-relayer/private-relayer/repository"
	"github.com/raylsnetwork/rayls-sovereign-relayer/types"
)

var (
	Enygma1 = types.Enygma{
		ResourceId:           "resource-enygma-1",
		FinalizedR:           big.NewInt(111),
		FinalizedBalance:     big.NewInt(999),
		FinalizedBlockNumber: big.NewInt(1000),
		PendingBlockNumber:   big.NewInt(2000),
	}
	Enygma2 = types.Enygma{
		ResourceId:           "resource-enygma-2",
		FinalizedR:           big.NewInt(222),
		FinalizedBalance:     big.NewInt(5000),
		FinalizedBlockNumber: big.NewInt(3000),
		PendingBlockNumber:   big.NewInt(4000),
	}
)

var (
	ModelEnygma1 = repository.Enygma{
		ResourceId:           Enygma1.ResourceId,
		FinalizedR:           Enygma1.FinalizedR.String(),
		FinalizedBalance:     Enygma1.FinalizedBalance.String(),
		FinalizedBlockNumber: Enygma1.FinalizedBlockNumber.Uint64(),
		PendingBlockNumber:   Enygma1.PendingBlockNumber.Uint64(),
	}
	ModelEnygma2 = repository.Enygma{
		ResourceId:           Enygma2.ResourceId,
		FinalizedR:           Enygma2.FinalizedR.String(),
		FinalizedBalance:     Enygma2.FinalizedBalance.String(),
		FinalizedBlockNumber: Enygma2.FinalizedBlockNumber.Uint64(),
		PendingBlockNumber:   Enygma2.PendingBlockNumber.Uint64(),
	}
)
