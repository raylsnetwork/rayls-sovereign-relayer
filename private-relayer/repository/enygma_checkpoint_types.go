package repository

import "time"

const EnygmaCheckpointCollectionName = "enygma_checkpoints"

type EnygmaCheckpoint struct {
	ID                      string    `db:"id"`
	ResourceId              string    `db:"resource_id"`
	FinalizedPublicBalanceX string    `db:"finalized_public_balance_x"`
	FinalizedPublicBalanceY string    `db:"finalized_public_balance_y"`
	FinalizedBlockNumber    uint64    `db:"finalized_block_number"`
	PendingBlockNumber      uint64    `db:"pending_block_number"`
	Status                  uint8     `db:"status"`
	CreatedAt               time.Time `db:"created_at"`
	UpdatedAt               time.Time `db:"updated_at"`
}
