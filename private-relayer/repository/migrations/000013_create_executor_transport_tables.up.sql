-- Executor tables for private-relayer service

CREATE TABLE IF NOT EXISTS to_private_hub_executor_transactions (
    id TEXT PRIMARY KEY,
    type INTEGER NOT NULL,
    data BYTEA,
    address TEXT NOT NULL,
    hash TEXT,
    result BIGINT,
    error TEXT,
    status INTEGER NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_to_private_hub_hash ON to_private_hub_executor_transactions(hash);
CREATE INDEX IF NOT EXISTS idx_to_private_hub_status ON to_private_hub_executor_transactions(status);

CREATE TABLE IF NOT EXISTS to_private_node_executor_transactions (
    id TEXT PRIMARY KEY,
    type INTEGER NOT NULL,
    data BYTEA,
    address TEXT NOT NULL,
    hash TEXT,
    result BIGINT,
    error TEXT,
    status INTEGER NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_to_private_node_hash ON to_private_node_executor_transactions(hash);
CREATE INDEX IF NOT EXISTS idx_to_private_node_status ON to_private_node_executor_transactions(status);
