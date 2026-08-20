package repo

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/cts/domain"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/cts/service"
)

const (
	SharedSecretsCollectionName = "shared_secrets"
)

type SharedSecretsRepository struct {
	pool      *pgxpool.Pool
	tableName string
}

var _ service.SharedSecretsRepository = &SharedSecretsRepository{}

func NewSharedSecretsRepository(
	tableName string,
	pool *pgxpool.Pool,
) *SharedSecretsRepository {
	return &SharedSecretsRepository{
		pool:      pool,
		tableName: tableName,
	}
}

func (r *SharedSecretsRepository) Create(ctx context.Context, secret domain.EncryptedSharedSecret) error {
	model := toSharedSecretModel(secret)

	query := fmt.Sprintf(`
		INSERT INTO %s (chain_id, encrypted_secret, initial_block, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
	`, r.tableName)

	_, err := r.pool.Exec(ctx, query,
		model.ChainId, model.EncryptedSecret, model.InitialBlock,
		model.CreatedAt, model.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("error creating shared secret: %w", err)
	}
	return nil
}

func (r *SharedSecretsRepository) GetByChainId(
	ctx context.Context,
	chainId string,
	blockNumber uint64,
) (domain.EncryptedSharedSecret, error) {
	// Filter: chain_id matches AND initial_block <= blockNumber
	// Sort by initial_block descending to get the most recent valid secret
	query := fmt.Sprintf(`
		SELECT chain_id, encrypted_secret, initial_block, created_at, updated_at
		FROM %s
		WHERE chain_id = $1 AND initial_block <= $2
		ORDER BY initial_block DESC
		LIMIT 1
	`, r.tableName)

	rows, err := r.pool.Query(ctx, query, chainId, blockNumber)
	if err != nil {
		return domain.EncryptedSharedSecret{}, fmt.Errorf("error getting shared secret: %w", err)
	}
	defer rows.Close()

	model, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[SharedSecretModel])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.EncryptedSharedSecret{}, service.ErrNoApplicableSharedSecret
		}
		return domain.EncryptedSharedSecret{}, fmt.Errorf("error decoding shared secret: %w", err)
	}

	return toSharedSecretType(model), nil
}

func (r *SharedSecretsRepository) GetByChainIds(
	ctx context.Context,
	chainIds []*big.Int,
	blockNumber uint64,
) ([]domain.EncryptedSharedSecret, error) {
	if len(chainIds) == 0 {
		return []domain.EncryptedSharedSecret{}, nil
	}

	// Convert chain IDs to strings
	chainIdsStrings := make([]string, len(chainIds))
	for i, chainId := range chainIds {
		chainIdsStrings[i] = chainId.String()
	}

	// Use DISTINCT ON to get the most recent valid secret for each chain_id
	query := fmt.Sprintf(`
		SELECT DISTINCT ON (chain_id)
		       chain_id, encrypted_secret, initial_block, created_at, updated_at
		FROM %s
		WHERE chain_id = ANY($1) AND initial_block <= $2
		ORDER BY chain_id, initial_block DESC
	`, r.tableName)

	rows, err := r.pool.Query(ctx, query, chainIdsStrings, blockNumber)
	if err != nil {
		return nil, fmt.Errorf("error aggregating shared secrets: %w", err)
	}
	defer rows.Close()

	models, err := pgx.CollectRows(rows, pgx.RowToStructByName[SharedSecretModel])
	if err != nil {
		return nil, fmt.Errorf("error decoding shared secrets: %w", err)
	}

	secrets := make([]domain.EncryptedSharedSecret, len(models))
	for i, model := range models {
		secrets[i] = toSharedSecretType(model)
	}

	return secrets, nil
}

func (r *SharedSecretsRepository) GetAll(ctx context.Context, blockNumber uint64) ([]domain.EncryptedSharedSecret, error) {
	// Use DISTINCT ON to get the most recent valid secret for each chain_id
	query := fmt.Sprintf(`
		SELECT DISTINCT ON (chain_id)
		       chain_id, encrypted_secret, initial_block, created_at, updated_at
		FROM %s
		WHERE initial_block <= $1
		ORDER BY chain_id, initial_block DESC
	`, r.tableName)

	rows, err := r.pool.Query(ctx, query, blockNumber)
	if err != nil {
		return nil, fmt.Errorf("error aggregating all shared secrets: %w", err)
	}
	defer rows.Close()

	models, err := pgx.CollectRows(rows, pgx.RowToStructByName[SharedSecretModel])
	if err != nil {
		return nil, fmt.Errorf("error decoding shared secrets: %w", err)
	}

	secrets := make([]domain.EncryptedSharedSecret, len(models))
	for i, model := range models {
		secrets[i] = toSharedSecretType(model)
	}

	return secrets, nil
}

func toSharedSecretModel(secret domain.EncryptedSharedSecret) SharedSecretModel {
	return SharedSecretModel{
		ChainId:         secret.ChainId,
		EncryptedSecret: secret.EncryptedSecret,
		InitialBlock:    secret.InitialBlock,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}
}

func toSharedSecretType(model SharedSecretModel) domain.EncryptedSharedSecret {
	return domain.EncryptedSharedSecret{
		ChainId:         model.ChainId,
		EncryptedSecret: model.EncryptedSecret,
		InitialBlock:    model.InitialBlock,
	}
}
