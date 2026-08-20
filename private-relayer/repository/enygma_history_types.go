package repository

import "time"

const EnygmaHistoryCollectionName = "enygma_history"

type EnygmaHistory struct {
	ResourceId            string    `db:"resource_id"`
	FromChainId           string    `db:"from_chain_id"`
	BalanceChange         string    `db:"balance_change"`
	BlockNumberPrivateHub uint64    `db:"block_number_private_hub"`
	RFactor               string    `db:"r_factor"`
	EventType             uint8     `db:"event_type"`
	PrivateHubTxHash      string    `db:"private_hub_tx_hash"`
	CreatedAt             time.Time `db:"created_at"`
}
