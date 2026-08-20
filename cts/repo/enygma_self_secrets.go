package repo

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/raylsnetwork/rayls-sovereign-relayer/cts/domain"
	"github.com/raylsnetwork/rayls-sovereign-relayer/cts/service"
)

const (
	//nolint:gosec // not a credential, just a collection name
	EnygmaSelfSecretsCollectionName = "enygma_self_secrets"
)

type EnygmaSelfSecretsRepository struct {
	pool      *pgxpool.Pool
	tableName string
}

var _ service.EnygmaSelfSecretsRepository = &EnygmaSelfSecretsRepository{}

func NewEnygmaSelfSecretsRepository(
	tableName string,
	pool *pgxpool.Pool,
) *EnygmaSelfSecretsRepository {
	return &EnygmaSelfSecretsRepository{
		pool:      pool,
		tableName: tableName,
	}
}

func (r *EnygmaSelfSecretsRepository) Create(ctx context.Context, secret domain.EncryptedEnygmaSelfSecret) error {
	model := toSelfSecretModel(secret)

	query := fmt.Sprintf(`
		INSERT INTO %s (encrypted_secret, initial_block, resource_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (initial_block, resource_id) DO NOTHING
	`, r.tableName)

	_, err := r.pool.Exec(ctx, query,
		model.EncryptedSecret, model.InitialBlock, model.ResourceID, model.CreatedAt, model.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("error creating self secret: %w", err)
	}
	return nil
}

func (r *EnygmaSelfSecretsRepository) GetByBlockNumberAndResource(
	ctx context.Context,
	blockNumber uint64,
	resourceID []byte,
) (domain.EncryptedEnygmaSelfSecret, error) {
	// Filter: resource_id = $1 AND initial_block <= $2
	// Sort by initial_block descending to get the most recent valid secret for this resource.
	query := fmt.Sprintf(`
		SELECT encrypted_secret, initial_block, resource_id, created_at, updated_at
		FROM %s
		WHERE resource_id = $1 AND initial_block <= $2
		ORDER BY initial_block DESC
		LIMIT 1
	`, r.tableName)

	rows, err := r.pool.Query(ctx, query, resourceID, blockNumber)
	if err != nil {
		return domain.EncryptedEnygmaSelfSecret{}, fmt.Errorf("error getting self secret: %w", err)
	}
	defer rows.Close()

	model, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[EnygmaSelfSecretModel])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.EncryptedEnygmaSelfSecret{}, service.ErrNoApplicableEnygmaSelfSecret
		}
		return domain.EncryptedEnygmaSelfSecret{}, fmt.Errorf("error decoding self secret: %w", err)
	}

	return toSelfSecretType(model), nil
}

func toSelfSecretModel(secret domain.EncryptedEnygmaSelfSecret) EnygmaSelfSecretModel {
	return EnygmaSelfSecretModel{
		EncryptedSecret: secret.EncryptedSecret,
		InitialBlock:    secret.InitialBlock,
		ResourceID:      secret.ResourceID,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}
}

func toSelfSecretType(model EnygmaSelfSecretModel) domain.EncryptedEnygmaSelfSecret {
	return domain.EncryptedEnygmaSelfSecret{
		EncryptedSecret: model.EncryptedSecret,
		InitialBlock:    model.InitialBlock,
		ResourceID:      model.ResourceID,
	}
}
