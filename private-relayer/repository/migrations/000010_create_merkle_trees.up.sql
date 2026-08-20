CREATE TABLE IF NOT EXISTS merkle_trees (
    type SMALLINT NOT NULL,
    token_address TEXT NOT NULL,
    number INTEGER NOT NULL,
    depth INTEGER NOT NULL,
    leaves TEXT[] NOT NULL DEFAULT '{}',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (type, number, token_address)
);
