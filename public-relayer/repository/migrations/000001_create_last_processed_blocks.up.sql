CREATE TABLE IF NOT EXISTS last_processed_block_numbers (
    chain TEXT PRIMARY KEY,
    last_block TEXT NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
