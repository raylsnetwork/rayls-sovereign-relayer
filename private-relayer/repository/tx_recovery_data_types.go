package repository

import "time"

const TxRecoveryDataCollectionName = "tx_recovery_data"

type TxRecoveryData struct {
	PrivateHubTxHash      string    `db:"private_hub_tx_hash"`
	ResourceID            string    `db:"resource_id"`
	PrivateHubBlockNumber uint64    `db:"private_hub_block_number"`
	FromChainID           string    `db:"from_chain_id"`
	TxBytes               []byte    `db:"tx_bytes"`
	EventType             uint8     `db:"event_type"`
	TxNature              string    `db:"tx_nature"`
	Status                int       `db:"status"`
	CreatedAt             time.Time `db:"created_at"`
}
