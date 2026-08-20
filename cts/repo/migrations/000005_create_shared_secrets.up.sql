-- Shared Secrets tables for all encryptor prefixes

CREATE TABLE IF NOT EXISTS plaintext_shared_secrets (
    chain_id TEXT NOT NULL,
    initial_block INTEGER NOT NULL,
    encrypted_secret BYTEA NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (chain_id, initial_block)
);
CREATE INDEX IF NOT EXISTS idx_plaintext_shared_secrets_chain_id ON plaintext_shared_secrets(chain_id);

CREATE TABLE IF NOT EXISTS aws_shared_secrets (
    chain_id TEXT NOT NULL,
    initial_block INTEGER NOT NULL,
    encrypted_secret BYTEA NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (chain_id, initial_block)
);
CREATE INDEX IF NOT EXISTS idx_aws_shared_secrets_chain_id ON aws_shared_secrets(chain_id);

CREATE TABLE IF NOT EXISTS gcp_shared_secrets (
    chain_id TEXT NOT NULL,
    initial_block INTEGER NOT NULL,
    encrypted_secret BYTEA NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (chain_id, initial_block)
);
CREATE INDEX IF NOT EXISTS idx_gcp_shared_secrets_chain_id ON gcp_shared_secrets(chain_id);
