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
	RaylsSignKeysCollectionName = "txsign_keys"

	kindPublicRelayer  = "public_relayer"
	kindPrivateRelayer = "private_relayer"
	kindAtomicService  = "atomic_service"
)

type RaylsSignKeysRepository struct {
	pool      *pgxpool.Pool
	tableName string
}

var _ service.RaylsSignKeysRepository = &RaylsSignKeysRepository{}

func NewRaylsSignKeysRepository(
	tableName string,
	pool *pgxpool.Pool,
) *RaylsSignKeysRepository {
	return &RaylsSignKeysRepository{
		pool:      pool,
		tableName: tableName,
	}
}

func (r *RaylsSignKeysRepository) CreatePublicRelayerRaylsSignKeys(
	ctx context.Context,
	keys domain.EncryptedPublicRelayerRaylsSignKeys,
) error {
	model := toPublicRelayerRaylsSignKeysModel(keys)
	query := fmt.Sprintf(`
		INSERT INTO %s (kind, public_chain_keys, private_chain_keys, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (kind) DO NOTHING
	`, r.tableName)

	result, err := r.pool.Exec(ctx, query,
		model.Kind, model.PublicChainKeys, model.PrivateChainKeys,
		model.CreatedAt, model.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("error creating public relayer keys: %w", err)
	}
	// RowsAffected() == 0 means ON CONFLICT triggered, so key already exists
	if result.RowsAffected() == 0 {
		return service.ErrRaylsSignKeysAlreadyExists
	}
	return nil
}

func (r *RaylsSignKeysRepository) GetPublicRelayerRaylsSignKeys(
	ctx context.Context,
) (domain.EncryptedPublicRelayerRaylsSignKeys, error) {
	query := fmt.Sprintf(`
		SELECT kind, public_chain_keys, private_chain_keys, created_at, updated_at
		FROM %s
		WHERE kind = $1
	`, r.tableName)

	rows, err := r.pool.Query(ctx, query, kindPublicRelayer)
	if err != nil {
		return domain.EncryptedPublicRelayerRaylsSignKeys{}, withstack.Wrap(fmt.Errorf("querying public relayer sign keys: %w", err))
	}
	defer rows.Close()

	model, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[PublicRelayerRaylsSignKeysModel])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.EncryptedPublicRelayerRaylsSignKeys{}, service.ErrNoRaylsSignKeys
		}
		return domain.EncryptedPublicRelayerRaylsSignKeys{}, withstack.Wrap(fmt.Errorf("collecting public relayer sign keys row: %w", err))
	}
	return fromPublicRelayerRaylsSignKeysModel(model), nil
}

func (r *RaylsSignKeysRepository) CreatePrivateRelayerRaylsSignKeys(
	ctx context.Context,
	keys domain.EncryptedPrivateRelayerRaylsSignKeys,
) error {
	model := toPrivateRelayerRaylsSignKeysModel(keys)
	query := fmt.Sprintf(`
		INSERT INTO %s (kind, private_hub_keys, private_hub_dvp_operator_keys, private_chain_keys, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (kind) DO NOTHING
	`, r.tableName)

	result, err := r.pool.Exec(ctx, query,
		model.Kind, model.PrivateHubKeys, model.PrivateHubDvpOperatorKeys,
		model.PrivateChainKeys, model.CreatedAt, model.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("error creating private relayer keys: %w", err)
	}
	// RowsAffected() == 0 means ON CONFLICT triggered, so key already exists
	if result.RowsAffected() == 0 {
		return service.ErrRaylsSignKeysAlreadyExists
	}
	return nil
}

func (r *RaylsSignKeysRepository) GetPrivateRelayerRaylsSignKeys(
	ctx context.Context,
) (domain.EncryptedPrivateRelayerRaylsSignKeys, error) {
	query := fmt.Sprintf(`
		SELECT kind, private_hub_keys, private_hub_dvp_operator_keys, private_chain_keys, created_at, updated_at
		FROM %s
		WHERE kind = $1
	`, r.tableName)

	rows, err := r.pool.Query(ctx, query, kindPrivateRelayer)
	if err != nil {
		return domain.EncryptedPrivateRelayerRaylsSignKeys{}, withstack.Wrap(fmt.Errorf("querying private relayer sign keys: %w", err))
	}
	defer rows.Close()

	model, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[PrivateRelayerRaylsSignKeysModel])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.EncryptedPrivateRelayerRaylsSignKeys{}, service.ErrNoRaylsSignKeys
		}
		return domain.EncryptedPrivateRelayerRaylsSignKeys{}, withstack.Wrap(fmt.Errorf("collecting private relayer sign keys row: %w", err))
	}
	return fromPrivateRelayerRaylsSignKeysModel(model), nil
}

func (r *RaylsSignKeysRepository) CreateAtomicServiceRaylsSignKeys(
	ctx context.Context,
	keys domain.EncryptedAtomicServiceRaylsSignKeys,
) error {
	model := toAtomicServiceRaylsSignKeysModel(keys)
	query := fmt.Sprintf(`
		INSERT INTO %s (kind, private_hub_keys, private_chain_keys, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (kind) DO NOTHING
	`, r.tableName)

	result, err := r.pool.Exec(ctx, query,
		model.Kind, model.PrivateHubKeys, model.PrivateChainKeys,
		model.CreatedAt, model.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("error creating atomic service keys: %w", err)
	}
	// RowsAffected() == 0 means ON CONFLICT triggered, so key already exists
	if result.RowsAffected() == 0 {
		return service.ErrRaylsSignKeysAlreadyExists
	}
	return nil
}

func (r *RaylsSignKeysRepository) GetAtomicServiceRaylsSignKeys(
	ctx context.Context,
) (domain.EncryptedAtomicServiceRaylsSignKeys, error) {
	query := fmt.Sprintf(`
		SELECT kind, private_hub_keys, private_chain_keys, created_at, updated_at
		FROM %s
		WHERE kind = $1
	`, r.tableName)

	rows, err := r.pool.Query(ctx, query, kindAtomicService)
	if err != nil {
		return domain.EncryptedAtomicServiceRaylsSignKeys{}, withstack.Wrap(fmt.Errorf("querying atomic service sign keys: %w", err))
	}
	defer rows.Close()

	model, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[AtomicServiceRaylsSignKeysModel])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.EncryptedAtomicServiceRaylsSignKeys{}, service.ErrNoRaylsSignKeys
		}
		return domain.EncryptedAtomicServiceRaylsSignKeys{}, withstack.Wrap(fmt.Errorf("collecting atomic service sign keys row: %w", err))
	}
	return fromAtomicServiceRaylsSignKeysModel(model), nil
}

func toPublicRelayerRaylsSignKeysModel(
	keys domain.EncryptedPublicRelayerRaylsSignKeys,
) PublicRelayerRaylsSignKeysModel {
	return PublicRelayerRaylsSignKeysModel{
		Kind:             kindPublicRelayer,
		PublicChainKeys:  keys.PublicChainKeys,
		PrivateChainKeys: keys.PrivateChainKeys,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Time{},
	}
}

func fromPublicRelayerRaylsSignKeysModel(m PublicRelayerRaylsSignKeysModel) domain.EncryptedPublicRelayerRaylsSignKeys {
	return domain.EncryptedPublicRelayerRaylsSignKeys{
		PublicChainKeys:  m.PublicChainKeys,
		PrivateChainKeys: m.PrivateChainKeys,
	}
}

func toPrivateRelayerRaylsSignKeysModel(
	keys domain.EncryptedPrivateRelayerRaylsSignKeys,
) PrivateRelayerRaylsSignKeysModel {
	return PrivateRelayerRaylsSignKeysModel{
		Kind:                      kindPrivateRelayer,
		PrivateHubKeys:            keys.PrivateHubKeys,
		PrivateHubDvpOperatorKeys: keys.PrivateHubDvpOperatorKeys,
		PrivateChainKeys:          keys.PrivateNodeKeys,
		CreatedAt:                 time.Now(),
		UpdatedAt:                 time.Time{},
	}
}

func fromPrivateRelayerRaylsSignKeysModel(
	m PrivateRelayerRaylsSignKeysModel,
) domain.EncryptedPrivateRelayerRaylsSignKeys {
	return domain.EncryptedPrivateRelayerRaylsSignKeys{
		PrivateHubKeys:            m.PrivateHubKeys,
		PrivateHubDvpOperatorKeys: m.PrivateHubDvpOperatorKeys,
		PrivateNodeKeys:           m.PrivateChainKeys,
	}
}

func toAtomicServiceRaylsSignKeysModel(
	keys domain.EncryptedAtomicServiceRaylsSignKeys,
) AtomicServiceRaylsSignKeysModel {
	return AtomicServiceRaylsSignKeysModel{
		Kind:             kindAtomicService,
		PrivateHubKeys:   keys.PrivateHubKeys,
		PrivateChainKeys: keys.PrivateChainKeys,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Time{},
	}
}

func fromAtomicServiceRaylsSignKeysModel(m AtomicServiceRaylsSignKeysModel) domain.EncryptedAtomicServiceRaylsSignKeys {
	return domain.EncryptedAtomicServiceRaylsSignKeys{
		PrivateHubKeys:   m.PrivateHubKeys,
		PrivateChainKeys: m.PrivateChainKeys,
	}
}

// Batch operations for typed key groups
func (r *RaylsSignKeysRepository) GetAll(ctx context.Context) ([]domain.EncryptedPublicRelayerRaylsSignKeys, error) {
	query := fmt.Sprintf(`
		SELECT kind, public_chain_keys, private_chain_keys, created_at, updated_at
		FROM %s
		WHERE kind = $1
	`, r.tableName)

	rows, err := r.pool.Query(ctx, query, kindPublicRelayer)
	if err != nil {
		return nil, withstack.Wrap(fmt.Errorf("querying all sign keys: %w", err))
	}
	defer rows.Close()

	models, err := pgx.CollectRows(rows, pgx.RowToStructByName[PublicRelayerRaylsSignKeysModel])
	if err != nil {
		return nil, withstack.Wrap(fmt.Errorf("collecting all sign keys rows: %w", err))
	}

	out := make([]domain.EncryptedPublicRelayerRaylsSignKeys, 0, len(models))
	for _, m := range models {
		out = append(out, fromPublicRelayerRaylsSignKeysModel(m))
	}
	return out, nil
}

func (r *RaylsSignKeysRepository) CreateAll(
	ctx context.Context,
	keys []domain.EncryptedPublicRelayerRaylsSignKeys,
) error {
	if len(keys) == 0 {
		return nil
	}

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return withstack.Wrap(fmt.Errorf("beginning sign keys batch transaction: %w", err))
	}
	defer func() { _ = tx.Rollback(ctx) }()

	query := fmt.Sprintf(`
		INSERT INTO %s (kind, public_chain_keys, private_chain_keys, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
	`, r.tableName)

	batch := &pgx.Batch{}
	for _, k := range keys {
		model := toPublicRelayerRaylsSignKeysModel(k)
		batch.Queue(query,
			model.Kind, model.PublicChainKeys, model.PrivateChainKeys,
			model.CreatedAt, model.UpdatedAt,
		)
	}

	br := tx.SendBatch(ctx, batch)
	if err := br.Close(); err != nil {
		return withstack.Wrap(fmt.Errorf("closing sign keys batch: %w", err))
	}

	if err := tx.Commit(ctx); err != nil {
		return withstack.Wrap(fmt.Errorf("committing sign keys batch transaction: %w", err))
	}
	return nil
}
