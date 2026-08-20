-- Enygma Spend Keys tables for all encryptor prefixes

CREATE TABLE IF NOT EXISTS plaintext_enygma_spend_keys (
    id SERIAL PRIMARY KEY,
    secret_key BYTEA NOT NULL,
    public_key BYTEA NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS aws_enygma_spend_keys (
    id SERIAL PRIMARY KEY,
    secret_key BYTEA NOT NULL,
    public_key BYTEA NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS gcp_enygma_spend_keys (
    id SERIAL PRIMARY KEY,
    secret_key BYTEA NOT NULL,
    public_key BYTEA NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
