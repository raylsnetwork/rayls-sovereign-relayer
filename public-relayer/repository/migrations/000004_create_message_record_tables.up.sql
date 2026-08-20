CREATE TABLE IF NOT EXISTS public_message_record (
    id           TEXT PRIMARY KEY,
    status       INTEGER NOT NULL,
    forward_hash TEXT NOT NULL DEFAULT '',
    revert_hash  TEXT NOT NULL DEFAULT '',
    error        TEXT NOT NULL DEFAULT '',
    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at   TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_public_message_record_status ON public_message_record(status);

CREATE TABLE IF NOT EXISTS private_message_record (
    id           TEXT PRIMARY KEY,
    status       INTEGER NOT NULL,
    forward_hash TEXT NOT NULL DEFAULT '',
    revert_hash  TEXT NOT NULL DEFAULT '',
    error        TEXT NOT NULL DEFAULT '',
    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at   TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_private_message_record_status ON private_message_record(status);
