ALTER TABLE dvp_swaps
    DROP COLUMN IF EXISTS self_salt,
    DROP COLUMN IF EXISTS dest_salt;

ALTER TABLE dvp_swaps ADD COLUMN source_payment_public_key TEXT;
ALTER TABLE dvp_swaps ADD COLUMN dest_payment_public_key TEXT;

UPDATE dvp_swaps SET expires_at = NOW() WHERE expires_at IS NULL;
ALTER TABLE dvp_swaps ALTER COLUMN expires_at SET NOT NULL;

CREATE TABLE IF NOT EXISTS dvp_swap_metadata (
    shared_id TEXT PRIMARY KEY,
    source_proof BYTEA,
    source_proof_type SMALLINT,
    source_calldata BYTEA,
    dest_proof BYTEA,
    dest_proof_type SMALLINT,
    dest_calldata BYTEA
);
