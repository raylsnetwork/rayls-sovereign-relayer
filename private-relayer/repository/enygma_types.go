package repository

import "time"

const EnygmaCollectionName = "enygma"

type Enygma struct {
	ResourceId           string    `db:"resource_id"`
	FinalizedR           string    `db:"finalized_r"`
	FinalizedBalance     string    `db:"finalized_balance"`
	FinalizedBlockNumber uint64    `db:"finalized_block_number"`
	PendingBlockNumber   uint64    `db:"pending_block_number"`
	CreatedAt            time.Time `db:"created_at"`
	UpdatedAt            time.Time `db:"updated_at"`
}
