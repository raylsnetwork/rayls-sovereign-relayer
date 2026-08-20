package repository

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/types"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/withstack"
)

// ErrAlreadyProcessed is returned when an InsertEnygmaHistory call hits
// a duplicate-key conflict, meaning the event has already been recorded.
var ErrAlreadyProcessed = errors.New("enygma history already processed")

type EnygmaHistoryRepository struct {
	pool      *pgxpool.Pool
	tableName string
}

func NewEnygmaHistoryRepository(pool *pgxpool.Pool) *EnygmaHistoryRepository {
	return &EnygmaHistoryRepository{
		pool:      pool,
		tableName: EnygmaHistoryCollectionName,
	}
}

func (repo *EnygmaHistoryRepository) InsertEnygmaHistory(
	ctx context.Context,
	history types.EnygmaHistory,
) error {
	model := enygmaTypeToEnygmaHistoryModel(history)
	model.CreatedAt = time.Now()

	query := fmt.Sprintf(`
		INSERT INTO %s (
			resource_id, from_chain_id, balance_change, r_factor,
			block_number_private_hub, event_type, private_hub_tx_hash, created_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`, repo.tableName)

	q := GetQuerier(ctx, repo.pool)
	_, err := q.Exec(ctx, query,
		model.ResourceId, model.FromChainId, model.BalanceChange, model.RFactor,
		model.BlockNumberPrivateHub, model.EventType, model.PrivateHubTxHash, model.CreatedAt,
	)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return ErrAlreadyProcessed
		}
		return withstack.Wrap(fmt.Errorf("PostgreSQL InsertEnygmaHistory: %w", err))
	}
	return nil
}

const enygmaHistorySelectColumns = `resource_id, from_chain_id, balance_change, r_factor,
		       block_number_private_hub, event_type, private_hub_tx_hash, created_at`

// GetEnygmaHistoryByUniqueKey returns the single EnygmaHistory record matching the
// compound unique index (resource_id, block_number_private_hub, from_chain_id, event_type).
// Returns (nil, nil) if no record exists.
func (repo *EnygmaHistoryRepository) GetEnygmaHistoryByUniqueKey(
	ctx context.Context,
	resourceId string,
	blockNumberPrivateHub *big.Int,
	fromChainId *big.Int,
	eventType types.EnygmaEventType,
) (*types.EnygmaHistory, error) {
	query := fmt.Sprintf(`
		SELECT %s
		FROM %s
		WHERE resource_id = $1
		  AND block_number_private_hub = $2
		  AND from_chain_id = $3
		  AND event_type = $4
	`, enygmaHistorySelectColumns, repo.tableName)

	q := GetQuerier(ctx, repo.pool)
	rows, err := q.Query(ctx, query,
		resourceId, blockNumberPrivateHub.Uint64(), fromChainId.String(), uint8(eventType),
	)
	if err != nil {
		return nil, withstack.Wrap(
			fmt.Errorf(
				"error getting enygma history for resource_id %s block %d: %w",
				resourceId,
				blockNumberPrivateHub.Uint64(),
				err,
			),
		)
	}
	defer rows.Close()

	model, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[EnygmaHistory])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil //nolint:nilnil // returning nil,nil is intentional for "not found" pattern
		}
		return nil, withstack.Wrap(
			fmt.Errorf(
				"error getting enygma history for resource_id %s block %d: %w",
				resourceId,
				blockNumberPrivateHub.Uint64(),
				err,
			),
		)
	}

	history, err := enygmaHistoryModelToEnygmaHistoryType(model)
	if err != nil {
		return nil, withstack.Wrap(err)
	}

	return &history, nil
}

func (repo *EnygmaHistoryRepository) GetEnygmaHistoryForCheckpoints(
	ctx context.Context,
	resourceIds []string,
	blockNumbers []*big.Int,
) ([]types.EnygmaHistory, error) {
	if len(resourceIds) != len(blockNumbers) {
		return nil, withstack.Wrap(fmt.Errorf("resourceIds and blockNumbers slices must have the same length"))
	}

	if len(resourceIds) == 0 {
		return []types.EnygmaHistory{}, nil
	}

	query := fmt.Sprintf(`
		SELECT %s
		FROM %s
		WHERE (
	`, enygmaHistorySelectColumns, repo.tableName)
	args := []any{}

	for i, resourceId := range resourceIds {
		if i > 0 {
			query += " OR "
		}
		argIdx := len(args) + 1
		query += fmt.Sprintf("(resource_id = $%d AND block_number_private_hub = $%d)", argIdx, argIdx+1)
		args = append(args, resourceId, blockNumbers[i].Uint64())
	}
	query += ")"

	q := GetQuerier(ctx, repo.pool)
	rows, err := q.Query(ctx, query, args...)
	if err != nil {
		return nil, withstack.Wrap(fmt.Errorf("PostgreSQL GetEnygmaHistoryForCheckpoints: %w", err))
	}
	defer rows.Close()

	models, err := pgx.CollectRows(rows, pgx.RowToStructByName[EnygmaHistory])
	if err != nil {
		return nil, withstack.Wrap(fmt.Errorf("PostgreSQL GetEnygmaHistoryForCheckpoints cursor: %w", err))
	}

	var results []types.EnygmaHistory
	for _, model := range models {
		history, err := enygmaHistoryModelToEnygmaHistoryType(model)
		if err != nil {
			return nil, withstack.Wrap(err)
		}

		results = append(results, history)
	}

	return results, nil
}

func enygmaTypeToEnygmaHistoryModel(history types.EnygmaHistory) EnygmaHistory {
	var fromChainId string
	if history.FromChainId != nil {
		fromChainId = history.FromChainId.String()
	}
	var balanceChange string
	if history.BalanceChange != nil {
		balanceChange = history.BalanceChange.String()
	}
	var rFactor string
	if history.RFactor != nil {
		rFactor = history.RFactor.String()
	}
	var blockNumberPrivateHub uint64
	if history.BlockNumberPrivateHub != nil {
		blockNumberPrivateHub = history.BlockNumberPrivateHub.Uint64()
	}

	return EnygmaHistory{
		ResourceId:            history.ResourceId,
		FromChainId:           fromChainId,
		BalanceChange:         balanceChange,
		RFactor:               rFactor,
		BlockNumberPrivateHub: blockNumberPrivateHub,
		EventType:             uint8(history.EventType),
		PrivateHubTxHash:      history.PrivateHubTxHash,
	}
}

func enygmaHistoryModelToEnygmaHistoryType(model EnygmaHistory) (types.EnygmaHistory, error) {
	var fromChainId *big.Int
	if model.FromChainId != "" {
		v, ok := new(big.Int).SetString(model.FromChainId, 10)
		if !ok {
			return types.EnygmaHistory{}, withstack.Wrap(fmt.Errorf("invalid fromChainId: %s", model.FromChainId))
		}
		fromChainId = v
	}

	var balanceChange *big.Int
	if model.BalanceChange != "" {
		v, ok := new(big.Int).SetString(model.BalanceChange, 10)
		if !ok {
			return types.EnygmaHistory{}, withstack.Wrap(fmt.Errorf("invalid balanceChange: %s", model.BalanceChange))
		}
		balanceChange = v
	}

	var rFactor *big.Int
	if model.RFactor != "" {
		v, ok := new(big.Int).SetString(model.RFactor, 10)
		if !ok {
			return types.EnygmaHistory{}, withstack.Wrap(fmt.Errorf("invalid rFactor: %s", model.RFactor))
		}
		rFactor = v
	}

	var blockNumberPrivateHub *big.Int
	if model.BlockNumberPrivateHub != 0 {
		blockNumberPrivateHub = new(big.Int).SetUint64(model.BlockNumberPrivateHub)
	}

	return types.EnygmaHistory{
		ResourceId:            model.ResourceId,
		FromChainId:           fromChainId,
		BalanceChange:         balanceChange,
		BlockNumberPrivateHub: blockNumberPrivateHub,
		RFactor:               rFactor,
		EventType:             types.EnygmaEventType(model.EventType),
		PrivateHubTxHash:      model.PrivateHubTxHash,
	}, nil
}
