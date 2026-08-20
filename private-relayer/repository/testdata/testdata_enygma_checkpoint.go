package testdata

import (
	"math/big"

	"github.com/raylsnetwork/rayls-privacy-relayer-api/private-relayer/repository"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/types"
)

var (
	EnygmaCheckpoint1 = types.EnygmaCheckpoint{
		ResourceId:              "resource-checkpoint-1",
		FinalizedPublicBalanceX: big.NewInt(100),
		FinalizedPublicBalanceY: big.NewInt(200),
		FinalizedBlockNumber:    big.NewInt(500),
		PendingBlockNumber:      big.NewInt(600),
		Status:                  types.EnygmaCheckpointStatusTentative,
	}
	EnygmaCheckpoint2 = types.EnygmaCheckpoint{
		ResourceId:              "resource-checkpoint-2",
		FinalizedPublicBalanceX: big.NewInt(300),
		FinalizedPublicBalanceY: big.NewInt(400),
		FinalizedBlockNumber:    big.NewInt(700),
		PendingBlockNumber:      big.NewInt(800),
		Status:                  types.EnygmaCheckpointStatusTentative,
	}
	EnygmaCheckpoint1HigherBlock = types.EnygmaCheckpoint{
		ResourceId:              "resource-checkpoint-1",
		FinalizedPublicBalanceX: big.NewInt(101),
		FinalizedPublicBalanceY: big.NewInt(201),
		FinalizedBlockNumber:    big.NewInt(900),
		PendingBlockNumber:      big.NewInt(1000),
		Status:                  types.EnygmaCheckpointStatusTentative,
	}
	EnygmaCheckpoint3Finalized = types.EnygmaCheckpoint{
		ResourceId:              "resource-checkpoint-3",
		FinalizedPublicBalanceX: big.NewInt(500),
		FinalizedPublicBalanceY: big.NewInt(600),
		FinalizedBlockNumber:    big.NewInt(1100),
		PendingBlockNumber:      big.NewInt(1200),
		Status:                  types.EnygmaCheckpointStatusFinal,
	}
)

var (
	checkpointID1          = "checkpoint-id-1"
	checkpointID2          = "checkpoint-id-2"
	checkpointID1Higher    = "checkpoint-id-1-higher"
	checkpointID3Finalized = "checkpoint-id-3-finalized"
)

var (
	ModelEnygmaCheckpoint1 = repository.EnygmaCheckpoint{
		ID:                      checkpointID1,
		ResourceId:              EnygmaCheckpoint1.ResourceId,
		FinalizedPublicBalanceX: EnygmaCheckpoint1.FinalizedPublicBalanceX.String(),
		FinalizedPublicBalanceY: EnygmaCheckpoint1.FinalizedPublicBalanceY.String(),
		FinalizedBlockNumber:    EnygmaCheckpoint1.FinalizedBlockNumber.Uint64(),
		PendingBlockNumber:      EnygmaCheckpoint1.PendingBlockNumber.Uint64(),
		Status:                  uint8(EnygmaCheckpoint1.Status),
	}
	ModelEnygmaCheckpoint2 = repository.EnygmaCheckpoint{
		ID:                      checkpointID2,
		ResourceId:              EnygmaCheckpoint2.ResourceId,
		FinalizedPublicBalanceX: EnygmaCheckpoint2.FinalizedPublicBalanceX.String(),
		FinalizedPublicBalanceY: EnygmaCheckpoint2.FinalizedPublicBalanceY.String(),
		FinalizedBlockNumber:    EnygmaCheckpoint2.FinalizedBlockNumber.Uint64(),
		PendingBlockNumber:      EnygmaCheckpoint2.PendingBlockNumber.Uint64(),
		Status:                  uint8(EnygmaCheckpoint2.Status),
	}
	ModelEnygmaCheckpoint1HigherBlock = repository.EnygmaCheckpoint{
		ID:                      checkpointID1Higher,
		ResourceId:              EnygmaCheckpoint1HigherBlock.ResourceId,
		FinalizedPublicBalanceX: EnygmaCheckpoint1HigherBlock.FinalizedPublicBalanceX.String(),
		FinalizedPublicBalanceY: EnygmaCheckpoint1HigherBlock.FinalizedPublicBalanceY.String(),
		FinalizedBlockNumber:    EnygmaCheckpoint1HigherBlock.FinalizedBlockNumber.Uint64(),
		PendingBlockNumber:      EnygmaCheckpoint1HigherBlock.PendingBlockNumber.Uint64(),
		Status:                  uint8(EnygmaCheckpoint1HigherBlock.Status),
	}
	ModelEnygmaCheckpoint3Finalized = repository.EnygmaCheckpoint{
		ID:                      checkpointID3Finalized,
		ResourceId:              EnygmaCheckpoint3Finalized.ResourceId,
		FinalizedPublicBalanceX: EnygmaCheckpoint3Finalized.FinalizedPublicBalanceX.String(),
		FinalizedPublicBalanceY: EnygmaCheckpoint3Finalized.FinalizedPublicBalanceY.String(),
		FinalizedBlockNumber:    EnygmaCheckpoint3Finalized.FinalizedBlockNumber.Uint64(),
		PendingBlockNumber:      EnygmaCheckpoint3Finalized.PendingBlockNumber.Uint64(),
		Status:                  uint8(EnygmaCheckpoint3Finalized.Status),
	}
)
