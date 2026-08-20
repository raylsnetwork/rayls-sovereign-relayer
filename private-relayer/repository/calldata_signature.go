// Decommissioning Teleport (vanilla, atomic).

package repository

import (
	"context"
	"fmt"
	"log/slog"
	"math/big"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/raylsnetwork/rayls-sovereign-relayer/types"
	"github.com/raylsnetwork/rayls-sovereign-relayer/withstack"
)

// Deprecated: Decommissioning Teleport (vanilla, atomic).
type CalldataSignatureRepository struct {
	pool      *pgxpool.Pool
	tableName string
}

// Deprecated: Decommissioning Teleport (vanilla, atomic).
func NewCalldataSignatureRepository(pool *pgxpool.Pool) *CalldataSignatureRepository {
	return &CalldataSignatureRepository{
		pool:      pool,
		tableName: CalldataSignatureCollectionName,
	}
}

// Deprecated: Decommissioning Teleport (vanilla, atomic).
func (sigRepo *CalldataSignatureRepository) GetBySharedIDs(
	ctx context.Context,
	sharedIds []string,
) ([]types.CalldataSignature, error) {
	query := fmt.Sprintf(`
		SELECT shared_id, status, signature, resource_id, signature_execute_chain_id,
		       destination_chain_id, signature_type
		FROM %s
		WHERE shared_id = ANY($1)
	`, sigRepo.tableName)

	rows, err := sigRepo.pool.Query(ctx, query, sharedIds)
	if err != nil {
		return nil, withstack.Wrap(fmt.Errorf("error getting signatures for sharedIds %v: %w", sharedIds, err))
	}
	defer rows.Close()

	dbResults, err := pgx.CollectRows(rows, pgx.RowToStructByName[CalldataSignature])
	if err != nil {
		return nil, withstack.Wrap(fmt.Errorf("error decoding signatures for sharedIds %v: %w", sharedIds, err))
	}

	results := []types.CalldataSignature{}
	for _, result := range dbResults {
		signature, err := calldataModelToCalldataType(result)
		if err != nil {
			return nil, withstack.Wrap(err)
		}

		results = append(results, signature)
	}

	return results, nil
}

// Deprecated: Decommissioning Teleport (vanilla, atomic).
func (sigRepo *CalldataSignatureRepository) GetDestinationUnlocksForSharedIDs(
	ctx context.Context,
	sharedIDs []string,
) ([]types.CalldataSignature, error) {
	return sigRepo.getBySignatureType(ctx, sharedIDs, types.UnlockOnDestinationSide)
}

// Deprecated: Decommissioning Teleport (vanilla, atomic).
func (sigRepo *CalldataSignatureRepository) GetDestinationRevertsForSharedIDs(
	ctx context.Context,
	sharedIDs []string,
) ([]types.CalldataSignature, error) {
	return sigRepo.getBySignatureType(ctx, sharedIDs, types.RevertOnDestinationSide)
}

// Deprecated: Decommissioning Teleport (vanilla, atomic).
func (sigRepo *CalldataSignatureRepository) GetSourceRevertForSharedIDs(
	ctx context.Context,
	sharedIDs []string,
) ([]types.CalldataSignature, error) {
	return sigRepo.getBySignatureType(ctx, sharedIDs, types.RevertOnSenderSide)
}

func (sigRepo *CalldataSignatureRepository) getBySignatureType(
	ctx context.Context,
	sharedIDs []string,
	sigType types.CallDataSignatureType,
) ([]types.CalldataSignature, error) {
	query := fmt.Sprintf(`
		SELECT shared_id, status, signature, resource_id, signature_execute_chain_id,
		       destination_chain_id, signature_type
		FROM %s
		WHERE shared_id = ANY($1) AND signature_type = $2
	`, sigRepo.tableName)

	rows, err := sigRepo.pool.Query(ctx, query, sharedIDs, sigType)
	if err != nil {
		return nil, withstack.Wrap(fmt.Errorf("error getting signatures for sharedIds %v: %w", sharedIDs, err))
	}
	defer rows.Close()

	dbResults, err := pgx.CollectRows(rows, pgx.RowToStructByName[CalldataSignature])
	if err != nil {
		return nil, withstack.Wrap(fmt.Errorf("error decoding signatures for sharedIds %v: %w", sharedIDs, err))
	}

	results := []types.CalldataSignature{}
	for _, result := range dbResults {
		signature, err := calldataModelToCalldataType(result)
		if err != nil {
			return nil, withstack.Wrap(err)
		}

		results = append(results, signature)
	}

	return results, nil
}

// Deprecated: Decommissioning Teleport (vanilla, atomic).
func (sigRepo *CalldataSignatureRepository) Create(ctx context.Context, sig types.CalldataSignature) error {
	model := calldataTypeToCalldataModel(sig)

	query := fmt.Sprintf(`
		INSERT INTO %s (
			shared_id, status, signature, resource_id, signature_execute_chain_id,
			destination_chain_id, signature_type
		) VALUES ($1, $2, $3, $4, $5, $6, $7)
	`, sigRepo.tableName)

	_, err := sigRepo.pool.Exec(ctx, query,
		model.SharedId, model.Status, model.Signature, model.ResourceId[:],
		model.SignatureExecuteChainId, model.DestinationChainId, model.SignatureType,
	)
	if err != nil {
		return withstack.Wrap(fmt.Errorf("error adding signature for sharedId %s: %w", sig.SharedId, err))
	}

	return nil
}

// DeleteBySharedIDs deletes all rows from the signatures table that match the sharedId
//
// Deprecated: Decommissioning Teleport (vanilla, atomic).
func (sigRepo *CalldataSignatureRepository) DeleteBySharedIDs(ctx context.Context, sharedIds []string) error {
	query := fmt.Sprintf(`DELETE FROM %s WHERE shared_id = ANY($1)`, sigRepo.tableName)

	res, err := sigRepo.pool.Exec(ctx, query, sharedIds)
	if err != nil {
		return withstack.Wrap(err)
	}

	if res.RowsAffected() == 0 {
		return nil
	}

	// Successfully deleted
	slog.Debug("Successfully deleted signatures", slog.Int("deleted_count", int(res.RowsAffected())))
	return nil
}

// Deprecated: Decommissioning Teleport (vanilla, atomic).
func (sigRepo *CalldataSignatureRepository) BatchCreate(
	ctx context.Context,
	signatures []types.CalldataSignature,
) error {
	tx, err := sigRepo.pool.Begin(ctx)
	if err != nil {
		return withstack.Wrap(fmt.Errorf("error creating batch signatures: %w", err))
	}
	defer func() { _ = tx.Rollback(ctx) }()

	query := fmt.Sprintf(`
		INSERT INTO %s (
			shared_id, status, signature, resource_id, signature_execute_chain_id,
			destination_chain_id, signature_type
		) VALUES ($1, $2, $3, $4, $5, $6, $7)
	`, sigRepo.tableName)

	batch := &pgx.Batch{}
	for _, sig := range signatures {
		model := calldataTypeToCalldataModel(sig)
		batch.Queue(query,
			model.SharedId, model.Status, model.Signature, model.ResourceId[:],
			model.SignatureExecuteChainId, model.DestinationChainId, model.SignatureType,
		)
	}

	br := tx.SendBatch(ctx, batch)
	if err := br.Close(); err != nil {
		return withstack.Wrap(fmt.Errorf("failed to create batch signatures: %w", err))
	}

	if err := tx.Commit(ctx); err != nil {
		return withstack.Wrap(fmt.Errorf("transaction failed: %w", err))
	}

	return nil
}

func calldataTypeToCalldataModel(sig types.CalldataSignature) CalldataSignature {
	return CalldataSignature{
		SharedId:                sig.SharedId,
		Status:                  0,
		Signature:               sig.Signature,
		ResourceId:              sig.ResourceId[:],
		SignatureExecuteChainId: "0",
		DestinationChainId:      "0",
		SignatureType:           uint8(sig.SignatureType),
	}
}

func calldataModelToCalldataType(model CalldataSignature) (types.CalldataSignature, error) {
	if _, ok := new(big.Int).SetString(model.SignatureExecuteChainId, 10); !ok {
		return types.CalldataSignature{}, fmt.Errorf("invalid SignatureExecuteChainId: %s", model.SignatureExecuteChainId)
	}
	if _, ok := new(big.Int).SetString(model.DestinationChainId, 10); !ok {
		return types.CalldataSignature{}, fmt.Errorf("invalid DestinationChainId: %s", model.DestinationChainId)
	}

	var resourceId [32]byte
	copy(resourceId[:], model.ResourceId)

	return types.CalldataSignature{
		SharedId:      model.SharedId,
		Signature:     model.Signature,
		ResourceId:    resourceId,
		SignatureType: types.CallDataSignatureType(model.SignatureType),
	}, nil
}
