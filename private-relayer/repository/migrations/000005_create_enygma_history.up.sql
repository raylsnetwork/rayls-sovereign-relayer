CREATE TABLE IF NOT EXISTS enygma_history (
    resource_id TEXT NOT NULL,
    from_chain_id TEXT NOT NULL,
    balance_change TEXT NOT NULL,
    block_number_private_hub BIGINT NOT NULL,
    r_factor TEXT NOT NULL,
    event_type SMALLINT NOT NULL,
    private_hub_tx_hash TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (resource_id, block_number_private_hub, from_chain_id, event_type)
);
