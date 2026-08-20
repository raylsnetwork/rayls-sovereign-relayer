package testdata

import (
	"math/big"

	"github.com/raylsnetwork/rayls-sovereign-relayer/private-relayer/repository"
	"github.com/raylsnetwork/rayls-sovereign-relayer/types"
)

var (
	ChainID1 = types.LastProcessedBlockDocumentPublicChain
	ChainID2 = types.DocumentIdLastProcessedBlockPrivacyNode
)

var (
	LastProcessedBlock1 = new(big.Int).SetInt64(7007)
	LastProcessedBlock2 = new(big.Int).SetInt64(10)
)

var (
	LastProcessedBlockModel1 = repository.LastProcessedBlockNumber{
		Chain:     string(ChainID1),
		LastBlock: LastProcessedBlock1.String(),
	}
	LastProcessedBlockModel2 = repository.LastProcessedBlockNumber{
		Chain:     string(ChainID2),
		LastBlock: LastProcessedBlock2.String(),
	}
)
