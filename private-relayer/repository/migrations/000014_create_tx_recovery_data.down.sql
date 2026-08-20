DROP TABLE IF EXISTS tx_recovery_data;

-- Restore enygma_history columns that were dropped in the up migration.
-- DEFAULT 0 is required because status was NOT NULL; operators rolling back
-- must treat the restored schema as equivalent to the pre-upgrade state.
ALTER TABLE enygma_history ADD COLUMN IF NOT EXISTS signed_sender_tx_bytes BYTEA;
ALTER TABLE enygma_history ADD COLUMN IF NOT EXISTS receiver_mint_tx_map BYTEA;
ALTER TABLE enygma_history ADD COLUMN IF NOT EXISTS status SMALLINT NOT NULL DEFAULT 0;
CREATE INDEX IF NOT EXISTS idx_enygma_history_status ON enygma_history(status);
