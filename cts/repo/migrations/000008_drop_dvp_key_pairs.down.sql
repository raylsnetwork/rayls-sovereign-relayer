CREATE TABLE IF NOT EXISTS plaintext_dvp_key_pairs (
    public_key BYTEA PRIMARY KEY,
    private_key BYTEA NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS aws_dvp_key_pairs (
    public_key BYTEA PRIMARY KEY,
    private_key BYTEA NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS gcp_dvp_key_pairs (
    public_key BYTEA PRIMARY KEY,
    private_key BYTEA NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);
