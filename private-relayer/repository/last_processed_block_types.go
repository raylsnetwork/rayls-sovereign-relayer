package repository

import "time"

const LastProcessedBlockNumberCollectionName = "last_processed_block_numbers"

type LastProcessedBlockNumber struct {
	Chain     string    `db:"chain"`
	LastBlock string    `db:"last_block"`
	UpdatedAt time.Time `db:"updated_at"`
}
