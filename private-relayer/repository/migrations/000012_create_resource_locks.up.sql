CREATE TABLE IF NOT EXISTS resource_locks (
    resource_id TEXT PRIMARY KEY,
    expires_at TIMESTAMPTZ NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_resource_locks_expires ON resource_locks(expires_at);
