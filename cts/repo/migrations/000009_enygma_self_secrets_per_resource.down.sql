TRUNCATE TABLE plaintext_enygma_self_secrets;
TRUNCATE TABLE aws_enygma_self_secrets;
TRUNCATE TABLE gcp_enygma_self_secrets;

ALTER TABLE plaintext_enygma_self_secrets
    DROP CONSTRAINT plaintext_enygma_self_secrets_pkey,
    DROP COLUMN resource_id,
    ADD PRIMARY KEY (initial_block);

ALTER TABLE aws_enygma_self_secrets
    DROP CONSTRAINT aws_enygma_self_secrets_pkey,
    DROP COLUMN resource_id,
    ADD PRIMARY KEY (initial_block);

ALTER TABLE gcp_enygma_self_secrets
    DROP CONSTRAINT gcp_enygma_self_secrets_pkey,
    DROP COLUMN resource_id,
    ADD PRIMARY KEY (initial_block);
