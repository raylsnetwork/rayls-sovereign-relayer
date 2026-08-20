-- Enygma Self Secrets tables for all encryptor prefixes

CREATE TABLE IF NOT EXISTS plaintext_enygma_self_secrets (
    initial_block INTEGER PRIMARY KEY,
    encrypted_secret BYTEA NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS aws_enygma_self_secrets (
    initial_block INTEGER PRIMARY KEY,
    encrypted_secret BYTEA NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS gcp_enygma_self_secrets (
    initial_block INTEGER PRIMARY KEY,
    encrypted_secret BYTEA NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
