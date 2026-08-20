// Decommissioning Teleport (vanilla, atomic).

package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/public-relayer/service"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/withstack"
)

type RevertSignatureRepository struct {
	pool      *pgxpool.Pool
	tableName string
}

func NewRevertSignatureRepository(
	collectionName string,
	pool *pgxpool.Pool,
) *RevertSignatureRepository {
	return &RevertSignatureRepository{
		pool:      pool,
		tableName: collectionName,
	}
}

func (repo *RevertSignatureRepository) Create(ctx context.Context, sig service.RevertSignature) error {
	record := revertSignatureTypeToModel(sig)

	query := fmt.Sprintf(`INSERT INTO %s (id, data) VALUES ($1, $2)`, repo.tableName)

	if _, err := repo.pool.Exec(ctx, query, record.ID, record.Data); err != nil {
		return withstack.Wrap(fmt.Errorf("failed to create revert signature record: %w", err))
	}
	return nil
}

// BatchCreate inserts multiple revert signature records in a single bulk operation.
// Uses ON CONFLICT DO NOTHING so valid records persist even when duplicates exist.
func (repo *RevertSignatureRepository) BatchCreate(ctx context.Context, sigs []service.RevertSignature) error {
	if len(sigs) == 0 {
		return nil
	}

	tx, err := repo.pool.Begin(ctx)
	if err != nil {
		return withstack.Wrap(fmt.Errorf("failed to begin transaction: %w", err))
	}
	defer func() { _ = tx.Rollback(ctx) }()

	query := fmt.Sprintf(`INSERT INTO %s (id, data) VALUES ($1, $2) ON CONFLICT (id) DO NOTHING`, repo.tableName)

	batch := &pgx.Batch{}
	for _, s := range sigs {
		record := revertSignatureTypeToModel(s)
		batch.Queue(query, record.ID, record.Data)
	}

	br := tx.SendBatch(ctx, batch)
	if err := br.Close(); err != nil {
		return withstack.Wrap(fmt.Errorf("failed to create revert signature records: %w", err))
	}

	if err := tx.Commit(ctx); err != nil {
		return withstack.Wrap(fmt.Errorf("failed to commit transaction: %w", err))
	}

	return nil
}

func revertSignatureTypeToModel(rs service.RevertSignature) RevertSignature {
	return RevertSignature{
		ID:   rs.ID,
		Data: rs.Data,
	}
}

// GetByIDs fetches revert signatures for the provided IDs.
func (repo *RevertSignatureRepository) GetByIDs(
	ctx context.Context,
	idSlice []string,
) ([]service.RevertSignature, error) {
	if len(idSlice) == 0 {
		return []service.RevertSignature{}, nil
	}

	query := fmt.Sprintf(`SELECT id, data FROM %s WHERE id = ANY($1)`, repo.tableName)

	rows, err := repo.pool.Query(ctx, query, idSlice)
	if err != nil {
		return nil, withstack.Wrap(fmt.Errorf("failed to get revert signature records: %w", err))
	}
	defer rows.Close()

	dbResults, err := pgx.CollectRows(rows, pgx.RowToStructByName[RevertSignature])
	if err != nil {
		return nil, withstack.Wrap(fmt.Errorf("failed to decode revert signature records: %w", err))
	}

	results := make([]service.RevertSignature, 0, len(dbResults))
	for _, r := range dbResults {
		results = append(results, service.RevertSignature{ID: r.ID, Data: r.Data})
	}
	return results, nil
}
