CREATE TABLE IF NOT EXISTS public_transport (
    id TEXT PRIMARY KEY,
    type INTEGER NOT NULL,
    data BYTEA,
    address TEXT NOT NULL,
    hash TEXT,
    result BIGINT,
    error TEXT,
    status INTEGER NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_public_transport_hash ON public_transport(hash);
CREATE INDEX IF NOT EXISTS idx_public_transport_status ON public_transport(status);

CREATE TABLE IF NOT EXISTS private_transport (
    id TEXT PRIMARY KEY,
    type INTEGER NOT NULL,
    data BYTEA,
    address TEXT NOT NULL,
    hash TEXT,
    result BIGINT,
    error TEXT,
    status INTEGER NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_private_transport_hash ON private_transport(hash);
CREATE INDEX IF NOT EXISTS idx_private_transport_status ON private_transport(status);
