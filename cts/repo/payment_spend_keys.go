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
	PaymentSpendKeysCollectionName = "enygma_spend_keys"
)

type PaymentSpendKeysRepository struct {
	pool      *pgxpool.Pool
	tableName string
}

var (
	_ service.PaymentSpendKeysRepository                            = &PaymentSpendKeysRepository{}
	_ service.MigrationRepository[domain.EncryptedPaymentSpendKeys] = &PaymentSpendKeysRepository{}
)

func NewPaymentSpendKeysRepository(
	tableName string,
	pool *pgxpool.Pool,
) *PaymentSpendKeysRepository {
	return &PaymentSpendKeysRepository{
		pool:      pool,
		tableName: tableName,
	}
}

func (r *PaymentSpendKeysRepository) Create(ctx context.Context, keys domain.EncryptedPaymentSpendKeys) error {
	paymentKeysModel := PaymentSpendKeysModel{
		SecretKey: keys.EncryptedSecretKey,
		PublicKey: keys.PublicKey,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	query := fmt.Sprintf(`
		INSERT INTO %s (secret_key, public_key, created_at, updated_at)
		VALUES ($1, $2, $3, $4)
	`, r.tableName)

	_, err := r.pool.Exec(ctx, query,
		paymentKeysModel.SecretKey, paymentKeysModel.PublicKey,
		paymentKeysModel.CreatedAt, paymentKeysModel.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("error creating enygma keys: %w", err)
	}

	return nil
}

func (r *PaymentSpendKeysRepository) Get(ctx context.Context) (domain.EncryptedPaymentSpendKeys, error) {
	query := fmt.Sprintf(`
		SELECT id, secret_key, public_key, created_at, updated_at
		FROM %s
		LIMIT 1
	`, r.tableName)

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return domain.EncryptedPaymentSpendKeys{}, fmt.Errorf("error getting enygma keys: %w", err)
	}
	defer rows.Close()

	paymentKeysModel, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[PaymentSpendKeysModel])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.EncryptedPaymentSpendKeys{}, service.ErrNoPaymentSpendKeySet
		}
		return domain.EncryptedPaymentSpendKeys{}, fmt.Errorf("error decoding enygma keys: %w", err)
	}

	return domain.EncryptedPaymentSpendKeys{
		EncryptedSecretKey: paymentKeysModel.SecretKey,
		PublicKey:          paymentKeysModel.PublicKey,
	}, nil
}

func (r *PaymentSpendKeysRepository) GetAll(ctx context.Context) ([]domain.EncryptedPaymentSpendKeys, error) {
	query := fmt.Sprintf(`
		SELECT id, secret_key, public_key, created_at, updated_at
		FROM %s
	`, r.tableName)

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, withstack.Wrap(fmt.Errorf("querying all payment spend keys: %w", err))
	}
	defer rows.Close()

	keysModel, err := pgx.CollectRows(rows, pgx.RowToStructByName[PaymentSpendKeysModel])
	if err != nil {
		return nil, withstack.Wrap(fmt.Errorf("collecting all payment spend keys rows: %w", err))
	}

	keys := []domain.EncryptedPaymentSpendKeys{}
	for _, keyModel := range keysModel {
		keys = append(keys, toPaymentSpendKeysType(keyModel))
	}

	return keys, nil
}

func (r *PaymentSpendKeysRepository) CreateAll(ctx context.Context, keys []domain.EncryptedPaymentSpendKeys) error {
	if len(keys) == 0 {
		return nil
	}

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return withstack.Wrap(fmt.Errorf("beginning payment spend keys batch transaction: %w", err))
	}
	defer func() { _ = tx.Rollback(ctx) }()

	query := fmt.Sprintf(`
		INSERT INTO %s (secret_key, public_key, created_at, updated_at)
		VALUES ($1, $2, $3, $4)
	`, r.tableName)

	batch := &pgx.Batch{}
	for _, key := range keys {
		model := toPaymentSpendKeysModel(key)
		batch.Queue(query,
			model.SecretKey, model.PublicKey, model.CreatedAt, model.UpdatedAt,
		)
	}

	br := tx.SendBatch(ctx, batch)
	if err := br.Close(); err != nil {
		return withstack.Wrap(fmt.Errorf("closing payment spend keys batch: %w", err))
	}

	if err := tx.Commit(ctx); err != nil {
		return withstack.Wrap(fmt.Errorf("committing payment spend keys batch transaction: %w", err))
	}
	return nil
}

func toPaymentSpendKeysModel(keys domain.EncryptedPaymentSpendKeys) PaymentSpendKeysModel {
	return PaymentSpendKeysModel{
		SecretKey: keys.EncryptedSecretKey,
		PublicKey: keys.PublicKey,
		CreatedAt: time.Now(),
		UpdatedAt: time.Time{},
	}
}

func toPaymentSpendKeysType(model PaymentSpendKeysModel) domain.EncryptedPaymentSpendKeys {
	return domain.EncryptedPaymentSpendKeys{
		EncryptedSecretKey: model.SecretKey,
		PublicKey:          model.PublicKey,
	}
}
