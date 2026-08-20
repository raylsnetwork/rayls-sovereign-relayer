ALTER TABLE dvp_deposits ADD COLUMN public_key TEXT;
ALTER TABLE dvp_deposits DROP COLUMN IF EXISTS salt;
