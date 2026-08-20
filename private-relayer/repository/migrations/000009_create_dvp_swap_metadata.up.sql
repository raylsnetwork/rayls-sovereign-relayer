CREATE TABLE IF NOT EXISTS dvp_swap_metadata (
    shared_id TEXT PRIMARY KEY,
    source_proof BYTEA,
    source_proof_type SMALLINT,
    source_calldata BYTEA,
    dest_proof BYTEA,
    dest_proof_type SMALLINT,
    dest_calldata BYTEA
);
