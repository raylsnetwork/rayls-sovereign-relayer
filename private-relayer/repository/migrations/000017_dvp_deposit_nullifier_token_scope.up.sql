DROP INDEX IF EXISTS idx_dvp_deposits_nullifier;

CREATE UNIQUE INDEX idx_dvp_deposits_token_nullifier
    ON dvp_deposits (token_address, nullifier)
    WHERE nullifier <> '';
