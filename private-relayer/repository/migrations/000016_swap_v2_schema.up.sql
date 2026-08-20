ALTER TABLE dvp_swaps ALTER COLUMN expires_at DROP NOT NULL;
ALTER TABLE dvp_swaps ALTER COLUMN expires_at SET DEFAULT NULL;
DROP TABLE IF EXISTS dvp_swap_metadata;

ALTER TABLE dvp_swaps DROP COLUMN IF EXISTS source_payment_public_key;
ALTER TABLE dvp_swaps DROP COLUMN IF EXISTS dest_payment_public_key;

ALTER TABLE dvp_swaps
    ADD COLUMN self_salt TEXT,
    ADD COLUMN dest_salt TEXT,
    ADD COLUMN cancel_preimage TEXT;
