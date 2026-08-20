package repository

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/raylsnetwork/rayls-sovereign-relayer/conv"
	"github.com/raylsnetwork/rayls-sovereign-relayer/types"
	"github.com/raylsnetwork/rayls-sovereign-relayer/withstack"
)

type DvpSwapRepository struct {
	pool      *pgxpool.Pool
	tableName string
}

func NewDvpSwapRepository(pool *pgxpool.Pool) *DvpSwapRepository {
	return &DvpSwapRepository{
		pool:      pool,
		tableName: DvpSwapCollectionName,
	}
}

func (r *DvpSwapRepository) CreateSwap(ctx context.Context, swap *types.DvpSwap) error {
	model := swapDomainToModel(swap)

	query := fmt.Sprintf(`
		INSERT INTO %s (
			shared_id, from_address, to_address,
			source_chain_id, dest_chain_id,
			token_in_amount, token_in_address, token_in_resource_id,
			token_in_type, token_in_id, token_out_amount, token_out_address,
			token_out_resource_id, token_out_type, token_out_id,
			status, created_at, cancelled_at,
			self_salt, dest_salt, cancel_preimage
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21)
	`, r.tableName)

	q := GetQuerier(ctx, r.pool)
	_, err := q.Exec(ctx, query,
		model.SharedID, model.From, model.To,
		model.SourceChainID, model.DestChainID,
		model.TokenInAmount, model.TokenInAddress, model.TokenInResourceID,
		model.TokenInType, model.TokenInID, model.TokenOutAmount, model.TokenOutAddress,
		model.TokenOutResourceID, model.TokenOutType, model.TokenOutID,
		model.Status, model.CreatedAt, model.CancelledAt,
		model.SelfSalt, model.DestSalt, model.CancelPreimage,
	)
	if err != nil {
		return withstack.Wrap(fmt.Errorf("failed to insert swap: %w", err))
	}
	return nil
}

func (r *DvpSwapRepository) GetSwapBySharedID(ctx context.Context, sharedID string) (*types.DvpSwap, error) {
	query := fmt.Sprintf(`
		SELECT shared_id, from_address, to_address,
		       source_chain_id, dest_chain_id,
		       token_in_amount, token_in_address, token_in_resource_id,
		       token_in_type, token_in_id, token_out_amount, token_out_address,
		       token_out_resource_id, token_out_type, token_out_id,
		       status, created_at, cancelled_at,
		       self_salt, dest_salt, cancel_preimage
		FROM %s
		WHERE shared_id = $1
	`, r.tableName)

	q := GetQuerier(ctx, r.pool)
	rows, err := q.Query(ctx, query, sharedID)
	if err != nil {
		return nil, withstack.Wrap(fmt.Errorf("failed to get swap with shared_id %s: %w", sharedID, err))
	}
	defer rows.Close()

	swap, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[DvpSwap])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil //nolint:nilnil // returning nil,nil is intentional for "not found" pattern
		}
		return nil, withstack.Wrap(fmt.Errorf("failed to get swap with shared_id %s: %w", sharedID, err))
	}

	domain, err := swapModelToDomain(swap)
	if err != nil {
		return nil, withstack.Wrap(fmt.Errorf("failed to decode swap with shared_id %s: %w", sharedID, err))
	}
	return domain, nil
}

func (r *DvpSwapRepository) GetPendingSwaps(ctx context.Context) ([]*types.DvpSwap, error) {
	query := fmt.Sprintf(`
		SELECT shared_id, from_address, to_address,
		       source_chain_id, dest_chain_id,
		       token_in_amount, token_in_address, token_in_resource_id,
		       token_in_type, token_in_id, token_out_amount, token_out_address,
		       token_out_resource_id, token_out_type, token_out_id,
		       status, created_at, cancelled_at,
		       self_salt, dest_salt, cancel_preimage
		FROM %s
		WHERE status IN ($1, $2)
	`, r.tableName)

	q := GetQuerier(ctx, r.pool)
	rows, err := q.Query(ctx, query, types.DvpSwapInitiated, types.DvpSwapWaitingConfirmation)
	if err != nil {
		return nil, withstack.Wrap(fmt.Errorf("failed to get pending swaps: %w", err))
	}
	defer rows.Close()

	dbResults, err := pgx.CollectRows(rows, pgx.RowToStructByName[DvpSwap])
	if err != nil {
		return nil, withstack.Wrap(fmt.Errorf("failed to decode pending swaps: %w", err))
	}

	var results []*types.DvpSwap
	for _, dbResult := range dbResults {
		domain, err := swapModelToDomain(dbResult)
		if err != nil {
			return nil, withstack.Wrap(fmt.Errorf("failed to convert pending swap: %w", err))
		}
		results = append(results, domain)
	}
	return results, nil
}

func (r *DvpSwapRepository) CancelSwap(ctx context.Context, sharedID string, status types.DvpSwapStatus) error {
	query := fmt.Sprintf(`UPDATE %s SET status = $1, cancelled_at = $2 WHERE shared_id = $3`, r.tableName)

	q := GetQuerier(ctx, r.pool)
	_, err := q.Exec(ctx, query, status, time.Now(), sharedID)
	if err != nil {
		return withstack.Wrap(fmt.Errorf("failed to cancel swap for shared_id %s: %w", sharedID, err))
	}

	return nil
}

func (r *DvpSwapRepository) UpdateSwapStatus(ctx context.Context, sharedID string, status types.DvpSwapStatus) error {
	query := fmt.Sprintf(`UPDATE %s SET status = $1 WHERE shared_id = $2`, r.tableName)

	q := GetQuerier(ctx, r.pool)
	_, err := q.Exec(ctx, query, status, sharedID)
	if err != nil {
		return withstack.Wrap(
			fmt.Errorf("failed to update swap status for shared_id %s to %v: %w", sharedID, status, err),
		)
	}
	return nil
}

func (r *DvpSwapRepository) UpdateSwapTo(ctx context.Context, sharedID string, to string) error {
	query := fmt.Sprintf(`UPDATE %s SET to_address = $1 WHERE shared_id = $2`, r.tableName)

	q := GetQuerier(ctx, r.pool)
	_, err := q.Exec(ctx, query, to, sharedID)
	if err != nil {
		return withstack.Wrap(
			fmt.Errorf("failed to update to_address for shared_id %s: %w", sharedID, err),
		)
	}
	return nil
}

func (r *DvpSwapRepository) UpdateSwapFrom(ctx context.Context, sharedID string, from string) error {
	query := fmt.Sprintf(`UPDATE %s SET from_address = $1 WHERE shared_id = $2`, r.tableName)

	q := GetQuerier(ctx, r.pool)
	_, err := q.Exec(ctx, query, from, sharedID)
	if err != nil {
		return withstack.Wrap(
			fmt.Errorf("failed to update from_address for shared_id %s: %w", sharedID, err),
		)
	}
	return nil
}

func swapModelToDomain(swap DvpSwap) (*types.DvpSwap, error) {
	sourceChainID, err := conv.StringToBigInt(swap.SourceChainID)
	if err != nil {
		return nil, fmt.Errorf("converting source chain ID: %w", err)
	}

	destChainID, err := conv.StringToBigInt(swap.DestChainID)
	if err != nil {
		return nil, fmt.Errorf("converting dest chain ID: %w", err)
	}

	tokenInAmount, err := conv.StringToBigInt(swap.TokenInAmount)
	if err != nil {
		return nil, fmt.Errorf("converting token in amount: %w", err)
	}

	tokenOutAmount, err := conv.StringToBigInt(swap.TokenOutAmount)
	if err != nil {
		return nil, fmt.Errorf("converting token out amount: %w", err)
	}

	var selfSalt *big.Int
	if swap.SelfSalt != nil {
		selfSalt, err = conv.StringToBigInt(*swap.SelfSalt)
		if err != nil {
			return nil, fmt.Errorf("converting self salt: %w", err)
		}
	}

	var destSalt *big.Int
	if swap.DestSalt != nil {
		destSalt, err = conv.StringToBigInt(*swap.DestSalt)
		if err != nil {
			return nil, fmt.Errorf("converting dest salt: %w", err)
		}
	}

	var cancelPreimage *big.Int
	if swap.CancelPreimage != nil {
		cancelPreimage, err = conv.StringToBigInt(*swap.CancelPreimage)
		if err != nil {
			return nil, fmt.Errorf("converting cancel preimage: %w", err)
		}
	}

	return &types.DvpSwap{
		SharedID:       swap.SharedID,
		From:           swap.From,
		To:             swap.To,
		SourceChainID:  sourceChainID,
		DestChainID:    destChainID,
		TokenInAmount:  tokenInAmount,
		TokenInAddress: swap.TokenInAddress,
		TokenInResourceID: swap.TokenInResourceID,
		TokenInType:    swap.TokenInType,
		TokenInID:      swap.TokenInID,
		TokenOutAmount: tokenOutAmount,
		TokenOutAddress: swap.TokenOutAddress,
		TokenOutResourceID: swap.TokenOutResourceID,
		TokenOutType:   swap.TokenOutType,
		TokenOutID:     swap.TokenOutID,
		Status:         swap.Status,
		CreatedAt:      swap.CreatedAt,
		CancelledAt:    swap.CancelledAt,
		SelfSalt:       selfSalt,
		DestSalt:       destSalt,
		CancelPreimage: cancelPreimage,
	}, nil
}

func swapDomainToModel(swap *types.DvpSwap) DvpSwap {
	var selfSalt *string
	if swap.SelfSalt != nil {
		s := swap.SelfSalt.String()
		selfSalt = &s
	}

	var destSalt *string
	if swap.DestSalt != nil {
		s := swap.DestSalt.String()
		destSalt = &s
	}

	var cancelPreimage *string
	if swap.CancelPreimage != nil {
		s := swap.CancelPreimage.String()
		cancelPreimage = &s
	}

	model := DvpSwap{
		SharedID:       swap.SharedID,
		From:           swap.From,
		To:             swap.To,
		SourceChainID:  swap.SourceChainID.String(),
		DestChainID:    swap.DestChainID.String(),
		TokenInAmount:  swap.TokenInAmount.String(),
		TokenInAddress: swap.TokenInAddress,
		TokenInResourceID: swap.TokenInResourceID,
		TokenInType:    swap.TokenInType,
		TokenInID:      swap.TokenInID,
		TokenOutAmount: swap.TokenOutAmount.String(),
		TokenOutAddress: swap.TokenOutAddress,
		TokenOutResourceID: swap.TokenOutResourceID,
		TokenOutType:   swap.TokenOutType,
		TokenOutID:     swap.TokenOutID,
		Status:         swap.Status,
		CancelledAt:    swap.CancelledAt,
		CreatedAt:      time.Now(),
		SelfSalt:       selfSalt,
		DestSalt:       destSalt,
		CancelPreimage: cancelPreimage,
	}

	return model
}
