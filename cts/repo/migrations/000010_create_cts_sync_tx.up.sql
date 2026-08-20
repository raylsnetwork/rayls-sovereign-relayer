-- cts_sync_tx persists durable state for the sync SignAndSend idempotency
-- path. Each (id) row is created after signing but before broadcast, so a
-- crash-recovery retry can resume against the same tx_hash instead of
-- re-signing (which would burn a new nonce and risk a duplicate on-chain
-- effect). The full RLP-encoded tx is stored alongside the hash so the
-- recovery path can replay it as an eth_call to extract revert reasons
-- without needing to re-fetch the tx from the chain.
--
-- result_state is the row's lifecycle:
--   pending  — created, no terminal verdict yet (may or may not be broadcast)
--   mined    — hard-terminal, tx mined with status==1 (receipt_json is durable result)
--   reverted — hard-terminal, tx mined with status==0 (revert_data is durable result)
--   failed   — soft-terminal, recovery exhausted; a relayer retry reopens it to pending
--
-- version backs optimistic concurrency on Save: a read-modify-write asserts the
-- version it loaded and bumps it, so a concurrent writer (crash recovery, a
-- parallel same-id caller, the retention job) can't silently clobber a verdict.
--
-- Distinct from cts_transaction (which serves the async BatchPipeline) —
-- the schemas don't overlap meaningfully and keeping them separate avoids
-- nullable-by-default columns in either table.
CREATE TABLE IF NOT EXISTS cts_sync_tx (
    id            TEXT        NOT NULL PRIMARY KEY,
    from_address  BYTEA       NOT NULL,
    tx_hash       BYTEA       NOT NULL,
    tx_rlp        BYTEA       NOT NULL,
    result_state  TEXT        NOT NULL DEFAULT 'pending'
                  CHECK (result_state IN ('pending', 'mined', 'reverted', 'failed')),
    receipt_json  BYTEA,
    revert_data   BYTEA,
    version       BIGINT      NOT NULL DEFAULT 0,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Retention job scans for hard-terminal rows past the retention window.
CREATE INDEX IF NOT EXISTS idx_cts_sync_tx_state_updated
    ON cts_sync_tx (result_state, updated_at);
