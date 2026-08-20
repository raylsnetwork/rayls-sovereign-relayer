CREATE TABLE IF NOT EXISTS enygma_checkpoints (
    id TEXT PRIMARY KEY DEFAULT gen_random_uuid()::TEXT,
    resource_id TEXT NOT NULL,
    finalized_public_balance_x TEXT NOT NULL,
    finalized_public_balance_y TEXT NOT NULL,
    finalized_block_number BIGINT NOT NULL,
    pending_block_number BIGINT NOT NULL,
    status SMALLINT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (resource_id, finalized_block_number)
);

CREATE INDEX IF NOT EXISTS idx_enygma_checkpoints_status ON enygma_checkpoints(status);
