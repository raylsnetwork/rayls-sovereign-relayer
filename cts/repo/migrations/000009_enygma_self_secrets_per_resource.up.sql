-- Scope enygma self secrets per (initial_block, resource_id).
-- The stored value is Poseidon(prevRFactor, paymentSpendKey) where prevRFactor
-- is per-resource (enygmaState.FinalizedR), so two resources at the same block
-- legitimately produce different self secrets. The original schema keyed only
-- on initial_block, which both rejected legitimate writes for distinct
-- resources at the same block and made the column "resource scope" implicit.
--
-- Existing rows have no resource_id and cannot be backfilled meaningfully —
-- self secrets are derivable from on-chain state, so truncate first.

TRUNCATE TABLE plaintext_enygma_self_secrets;
TRUNCATE TABLE aws_enygma_self_secrets;
TRUNCATE TABLE gcp_enygma_self_secrets;

ALTER TABLE plaintext_enygma_self_secrets
    DROP CONSTRAINT plaintext_enygma_self_secrets_pkey,
    ADD COLUMN resource_id BYTEA NOT NULL,
    ADD PRIMARY KEY (initial_block, resource_id);

ALTER TABLE aws_enygma_self_secrets
    DROP CONSTRAINT aws_enygma_self_secrets_pkey,
    ADD COLUMN resource_id BYTEA NOT NULL,
    ADD PRIMARY KEY (initial_block, resource_id);

ALTER TABLE gcp_enygma_self_secrets
    DROP CONSTRAINT gcp_enygma_self_secrets_pkey,
    ADD COLUMN resource_id BYTEA NOT NULL,
    ADD PRIMARY KEY (initial_block, resource_id);
