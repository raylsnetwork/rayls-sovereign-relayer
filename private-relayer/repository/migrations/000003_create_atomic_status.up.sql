CREATE TABLE IF NOT EXISTS atomic_status (
    shared_id TEXT PRIMARY KEY,
    status SMALLINT NOT NULL,
    is_processed BOOLEAN NOT NULL DEFAULT FALSE
);
