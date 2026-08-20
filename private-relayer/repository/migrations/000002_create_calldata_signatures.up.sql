CREATE TABLE IF NOT EXISTS calldata_signatures (
    shared_id TEXT NOT NULL,
    status SMALLINT NOT NULL,
    signature BYTEA NOT NULL,
    resource_id BYTEA NOT NULL,
    signature_execute_chain_id TEXT NOT NULL,
    destination_chain_id TEXT NOT NULL,
    signature_type SMALLINT NOT NULL,
    PRIMARY KEY (shared_id, signature_type)
);

CREATE INDEX IF NOT EXISTS idx_calldata_signatures_shared_id ON calldata_signatures(shared_id);
