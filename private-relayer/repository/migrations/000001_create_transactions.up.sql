CREATE TABLE IF NOT EXISTS transactions (
    id TEXT PRIMARY KEY DEFAULT gen_random_uuid()::TEXT,
    tx_hash TEXT NOT NULL,
    tx_hash_destination TEXT,
    log_index INTEGER NOT NULL,
    shared_id TEXT NOT NULL UNIQUE,
    state SMALLINT NOT NULL,
    outcome TEXT NOT NULL DEFAULT 'pending'
        CHECK (outcome IN ('pending', 'success', 'reverted', 'failed')),
    proof_invalid BOOLEAN NOT NULL DEFAULT FALSE,
    originator_chain_id TEXT NOT NULL,
    destination_chain_id TEXT NOT NULL,
    msg_id BYTEA NOT NULL UNIQUE,
    is_atomic BOOLEAN NOT NULL DEFAULT FALSE,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    batch_id TEXT,
    batch_tx_hash_on_private_hub TEXT,
    resource_id TEXT,
    from_contract_address TEXT,
    from_user_address TEXT,
    transfer_metadata_id TEXT,
    transfer_metadata_amount TEXT,
    block_number BIGINT,
    parent_hash TEXT
);

CREATE INDEX IF NOT EXISTS idx_transactions_state_outcome ON transactions(state, outcome);
CREATE INDEX IF NOT EXISTS idx_transactions_tx_hash_log_index ON transactions(tx_hash, log_index);
