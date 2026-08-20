CREATE TABLE IF NOT EXISTS dvp_swaps (
    shared_id TEXT PRIMARY KEY,
    from_address TEXT NOT NULL,
    to_address TEXT NOT NULL,
    source_payment_public_key TEXT,
    source_chain_id TEXT NOT NULL,
    dest_payment_public_key TEXT,
    dest_chain_id TEXT NOT NULL,
    token_in_amount TEXT NOT NULL,
    token_in_address TEXT NOT NULL,
    token_in_resource_id TEXT,
    token_in_type SMALLINT NOT NULL,
    token_in_id TEXT,
    token_out_amount TEXT NOT NULL,
    token_out_address TEXT NOT NULL,
    token_out_resource_id TEXT,
    token_out_type SMALLINT NOT NULL,
    token_out_id TEXT,
    status SMALLINT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    expires_at TIMESTAMPTZ NOT NULL,
    cancelled_at TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_dvp_swaps_expires_status ON dvp_swaps(expires_at, status);
