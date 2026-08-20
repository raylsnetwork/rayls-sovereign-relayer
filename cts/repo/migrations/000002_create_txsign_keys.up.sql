-- Transaction Signing Keys tables for all encryptor prefixes

CREATE TABLE IF NOT EXISTS plaintext_txsign_keys (
    kind TEXT PRIMARY KEY,
    public_chain_keys BYTEA[],
    private_chain_keys BYTEA[],
    private_hub_keys BYTEA[],
    private_hub_dvp_operator_keys BYTEA[],
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS aws_txsign_keys (
    kind TEXT PRIMARY KEY,
    public_chain_keys BYTEA[],
    private_chain_keys BYTEA[],
    private_hub_keys BYTEA[],
    private_hub_dvp_operator_keys BYTEA[],
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS gcp_txsign_keys (
    kind TEXT PRIMARY KEY,
    public_chain_keys BYTEA[],
    private_chain_keys BYTEA[],
    private_hub_keys BYTEA[],
    private_hub_dvp_operator_keys BYTEA[],
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
