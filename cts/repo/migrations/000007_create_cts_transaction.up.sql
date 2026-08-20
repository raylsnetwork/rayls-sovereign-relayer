CREATE TABLE IF NOT EXISTS cts_transaction (
    correlation_id    TEXT        NOT NULL,
    identity          TEXT        NOT NULL,
    message_type      TEXT        NOT NULL,
    address           BYTEA       NOT NULL,
    calldata          BYTEA       NOT NULL,
    status            TEXT        NOT NULL
                      CHECK (status IN ('pending','sent','finished','failed')),
    tx_hash           BYTEA,
    receipt_status    SMALLINT,
    revert_data       BYTEA,
    error_reason      TEXT,
    send_attempts     INT         NOT NULL DEFAULT 0,
    receipt_attempts  INT         NOT NULL DEFAULT 0,
    sent_at           TIMESTAMPTZ,
    created_at        TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at        TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    -- Composite PK: a single shared id can carry multiple message
    -- types (e.g. crosschain.atomic + atomic.destination-unlock +
    -- atomic.destination-revert). Each is a distinct row.
    PRIMARY KEY (correlation_id, message_type)
);

CREATE INDEX IF NOT EXISTS idx_cts_tx_pending ON cts_transaction (identity, created_at)
    WHERE status = 'pending';
CREATE INDEX IF NOT EXISTS idx_cts_tx_sent    ON cts_transaction (identity, sent_at)
    WHERE status = 'sent';
