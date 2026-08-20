CREATE TABLE IF NOT EXISTS dvp_deposits (
    commitment TEXT PRIMARY KEY,
    user_address TEXT NOT NULL,
    public_key TEXT NOT NULL,
    token_amount TEXT NOT NULL,
    token_address TEXT NOT NULL,
    token_type SMALLINT NOT NULL,
    token_id TEXT,
    tree_number INTEGER,
    nullifier TEXT,
    status SMALLINT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_dvp_deposits_nullifier ON dvp_deposits(nullifier) WHERE nullifier IS NOT NULL;
