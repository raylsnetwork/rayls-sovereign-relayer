CREATE TABLE IF NOT EXISTS tx_recovery_data (
    private_hub_tx_hash       TEXT PRIMARY KEY,
    resource_id               TEXT NOT NULL,
    private_hub_block_number  BIGINT NOT NULL,
    from_chain_id             TEXT NOT NULL,
    tx_bytes                  BYTEA,
    event_type                SMALLINT NOT NULL,
    tx_nature                 TEXT NOT NULL,
    status                    SMALLINT NOT NULL,
    created_at                TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_tx_recovery_data_nature_status ON tx_recovery_data(tx_nature, status);

-- Drop recovery columns from enygma_history that were moved to tx_recovery_data.
-- These columns exist in version/2.6.4 and must be removed on upgrade so that
-- InsertEnygmaHistory (which no longer includes them) does not get a NOT NULL
-- constraint violation on "status".
DROP INDEX IF EXISTS idx_enygma_history_status;
ALTER TABLE enygma_history DROP COLUMN IF EXISTS status;
ALTER TABLE enygma_history DROP COLUMN IF EXISTS signed_sender_tx_bytes;
ALTER TABLE enygma_history DROP COLUMN IF EXISTS receiver_mint_tx_map;
