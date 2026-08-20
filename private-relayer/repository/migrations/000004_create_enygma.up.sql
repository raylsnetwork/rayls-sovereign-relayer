CREATE TABLE IF NOT EXISTS enygma (
    resource_id TEXT PRIMARY KEY,
    finalized_r TEXT NOT NULL,
    finalized_balance TEXT NOT NULL,
    finalized_block_number BIGINT NOT NULL,
    pending_block_number BIGINT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
