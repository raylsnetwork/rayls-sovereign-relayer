package repo

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/cts/domain"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/cts/service"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/withstack"
)

const (
	RaylsViewKeysCollectionName = "rayls_view_keys"
)

type RaylsViewKeysRepository struct {
	pool      *pgxpool.Pool
	tableName string
}

var (
	_ service.RaylsViewKeysRepository                               = &RaylsViewKeysRepository{}
	_ service.MigrationRepository[domain.EncryptedRaylsViewKeyPair] = &RaylsViewKeysRepository{}
)

func NewRaylsViewKeysRepository(
	tableName string,
	pool *pgxpool.Pool,
) *RaylsViewKeysRepository {
	return &RaylsViewKeysRepository{
		pool:      pool,
		tableName: tableName,
	}
}

func (ks *RaylsViewKeysRepository) GetForBlockNumber(
	ctx context.Context,
	blockNumber uint64,
) (domain.EncryptedRaylsViewKeyPair, error) {
	// Check if any records exist
	countQuery := fmt.Sprintf(`SELECT COUNT(*) FROM %s`, ks.tableName)
	var count int64
	err := ks.pool.QueryRow(ctx, countQuery).Scan(&count)
	if err != nil {
		return domain.EncryptedRaylsViewKeyPair{}, withstack.Wrap(fmt.Errorf("counting rayls view keys: %w", err))
	}
	if count == 0 {
		return domain.EncryptedRaylsViewKeyPair{}, service.ErrNoRaylsViewKeysSet
	}

	query := fmt.Sprintf(`
		SELECT initial_block, encrypted_secret_key, public_key, created_at, updated_at
		FROM %s
		WHERE initial_block <= $1
		ORDER BY initial_block DESC
		LIMIT 1
	`, ks.tableName)

	rows, err := ks.pool.Query(ctx, query, blockNumber)
	if err != nil {
		return domain.EncryptedRaylsViewKeyPair{}, withstack.Wrap(fmt.Errorf("querying rayls view keys for block number: %w", err))
	}
	defer rows.Close()

	model, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[RaylsViewKeysModel])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.EncryptedRaylsViewKeyPair{}, service.ErrNoApplicableRaylsViewKeys
		}
		return domain.EncryptedRaylsViewKeyPair{}, withstack.Wrap(fmt.Errorf("collecting rayls view keys row: %w", err))
	}

	return toRaylsViewKeyPairType(model), nil
}

func (ks *RaylsViewKeysRepository) Create(ctx context.Context, pair domain.EncryptedRaylsViewKeyPair) error {
	raylsViewKeysModel := toRaylsViewKeyPairModel(pair)

	query := fmt.Sprintf(`
		INSERT INTO %s (initial_block, encrypted_secret_key, public_key, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
	`, ks.tableName)

	_, err := ks.pool.Exec(ctx, query,
		raylsViewKeysModel.InitialBlock,
		raylsViewKeysModel.EncryptedSecretKey,
		raylsViewKeysModel.RaylsViewPublicKey,
		raylsViewKeysModel.CreatedAt,
		raylsViewKeysModel.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("error creating key: %w", err)
	}

	return nil
}

func (ks *RaylsViewKeysRepository) DeleteByPublicKey(ctx context.Context, pubKey string) error {
	query := fmt.Sprintf(`DELETE FROM %s WHERE public_key = $1`, ks.tableName)

	_, err := ks.pool.Exec(ctx, query, []byte(pubKey))
	if err != nil {
		return fmt.Errorf("error deleting key: %w", err)
	}

	return nil
}

func (ks *RaylsViewKeysRepository) GetAll(ctx context.Context) ([]domain.EncryptedRaylsViewKeyPair, error) {
	query := fmt.Sprintf(`
		SELECT initial_block, encrypted_secret_key, public_key, created_at, updated_at
		FROM %s
	`, ks.tableName)

	rows, err := ks.pool.Query(ctx, query)
	if err != nil {
		return nil, withstack.Wrap(fmt.Errorf("querying all rayls view keys: %w", err))
	}
	defer rows.Close()

	keysModel, err := pgx.CollectRows(rows, pgx.RowToStructByName[RaylsViewKeysModel])
	if err != nil {
		return nil, withstack.Wrap(fmt.Errorf("collecting all rayls view keys rows: %w", err))
	}

	keys := []domain.EncryptedRaylsViewKeyPair{}
	for _, keyModel := range keysModel {
		keys = append(keys, toRaylsViewKeyPairType(keyModel))
	}

	return keys, nil
}

func (ks *RaylsViewKeysRepository) CreateAll(ctx context.Context, keys []domain.EncryptedRaylsViewKeyPair) error {
	if len(keys) == 0 {
		return nil
	}

	tx, err := ks.pool.Begin(ctx)
	if err != nil {
		return withstack.Wrap(fmt.Errorf("beginning rayls view keys batch transaction: %w", err))
	}
	defer func() { _ = tx.Rollback(ctx) }()

	query := fmt.Sprintf(`
		INSERT INTO %s (initial_block, encrypted_secret_key, public_key, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
	`, ks.tableName)

	batch := &pgx.Batch{}
	for _, key := range keys {
		model := toRaylsViewKeyPairModel(key)
		batch.Queue(query,
			model.InitialBlock, model.EncryptedSecretKey,
			model.RaylsViewPublicKey, model.CreatedAt, model.UpdatedAt,
		)
	}

	br := tx.SendBatch(ctx, batch)
	if err := br.Close(); err != nil {
		return withstack.Wrap(fmt.Errorf("closing rayls view keys batch: %w", err))
	}

	if err := tx.Commit(ctx); err != nil {
		return withstack.Wrap(fmt.Errorf("committing rayls view keys batch transaction: %w", err))
	}
	return nil
}

func toRaylsViewKeyPairModel(keys domain.EncryptedRaylsViewKeyPair) RaylsViewKeysModel {
	return RaylsViewKeysModel{
		InitialBlock:       keys.InitialBlock,
		EncryptedSecretKey: keys.EncryptedRaylsViewPrivateKey,
		RaylsViewPublicKey: keys.RaylsViewPublicKey,
		CreatedAt:          time.Now(),
		UpdatedAt:          time.Time{},
	}
}

func toRaylsViewKeyPairType(model RaylsViewKeysModel) domain.EncryptedRaylsViewKeyPair {
	return domain.EncryptedRaylsViewKeyPair{
		InitialBlock:                 model.InitialBlock,
		EncryptedRaylsViewPrivateKey: model.EncryptedSecretKey,
		RaylsViewPublicKey:           model.RaylsViewPublicKey,
	}
}
