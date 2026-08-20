DROP INDEX IF EXISTS idx_dvp_deposits_token_nullifier;

CREATE INDEX IF NOT EXISTS idx_dvp_deposits_nullifier
    ON dvp_deposits (nullifier)
    WHERE nullifier <> '';
