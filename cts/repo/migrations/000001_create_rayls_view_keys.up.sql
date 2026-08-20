-- Rayls View Keys tables for all encryptor prefixes

CREATE TABLE IF NOT EXISTS plaintext_rayls_view_keys (
    initial_block INTEGER PRIMARY KEY,
    encrypted_secret_key BYTEA NOT NULL,
    public_key BYTEA NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS aws_rayls_view_keys (
    initial_block INTEGER PRIMARY KEY,
    encrypted_secret_key BYTEA NOT NULL,
    public_key BYTEA NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS gcp_rayls_view_keys (
    initial_block INTEGER PRIMARY KEY,
    encrypted_secret_key BYTEA NOT NULL,
    public_key BYTEA NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
