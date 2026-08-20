CREATE TABLE IF NOT EXISTS public_revert_signature (
    id TEXT PRIMARY KEY,
    data BYTEA NOT NULL
);

CREATE TABLE IF NOT EXISTS private_revert_signature (
    id TEXT PRIMARY KEY,
    data BYTEA NOT NULL
);
