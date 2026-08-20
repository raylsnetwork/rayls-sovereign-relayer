package service

import (
	"context"
	"fmt"
	"log/slog"
	"math/big"

	"github.com/raylsnetwork/rayls-privacy-relayer-api/cryptography"
	telemetry "github.com/raylsnetwork/rayls-privacy-relayer-api/otel"
	repository2 "github.com/raylsnetwork/rayls-privacy-relayer-api/private-relayer/repository"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/types"
	"go.opentelemetry.io/otel/trace"
)

type syncTracer interface {
	Start(ctx context.Context, spanName string, opts ...trace.SpanStartOption) (context.Context, trace.Span)
}

// var _ syncTracer = (*adapters.OTelTracer)(nil)

type syncTransactionManager interface {
	WithTransaction(ctx context.Context, fn func(ctx context.Context) error) error
}

type syncEnygmaRepository interface {
	GetEnygmaByResourceIds(ctx context.Context, resourceIds []string) ([]types.Enygma, error)
	UpdateEnygma(
		ctx context.Context,
		resourceId string,
		finalizedBalance, finalizedR, finalizedBlockNumber, pendingBlockNumber *big.Int,
	) error
}

type syncEnygmaHistoryRepository interface {
	GetEnygmaHistoryForCheckpoints(
		ctx context.Context,
		resourceIds []string,
		blockNumbers []*big.Int,
	) ([]types.EnygmaHistory, error)
}

type syncEnygmaCheckpointRepository interface {
	GetValidationCandidates(ctx context.Context) ([]types.EnygmaCheckpoint, error)
	MarkAsFinalized(ctx context.Context, resourceId string, finalizedBlockNumber *big.Int) error
	GetLatestCheckpointByFilters(
		ctx context.Context,
		resourceId string,
		status *types.EnygmaCheckpointStatus,
		finalizedBlockNumber *big.Int,
		pendingBlockNumber *big.Int,
	) (*types.EnygmaCheckpoint, error)
	CreateEnygmaCheckpoint(ctx context.Context, checkpoint types.EnygmaCheckpoint) error
}

type syncResync interface {
	ResyncEnygma(ctx context.Context, resourceId string) error
}

// Compile-time interface implementation checks
var (
	_ syncTransactionManager         = (*repository2.TransactionManager)(nil)
	_ syncEnygmaRepository           = (*repository2.EnygmaRepository)(nil)
	_ syncEnygmaHistoryRepository    = (*repository2.EnygmaHistoryRepository)(nil)
	_ syncEnygmaCheckpointRepository = (*repository2.EnygmaCheckpointRepository)(nil)
)

type SyncConfig struct {
	MaxRetries int // max retry attempts before forcing resync
}

type EnygmaSyncService struct {
	conf                        SyncConfig
	txManager                   syncTransactionManager
	enygmaRepository            syncEnygmaRepository
	enygmaHistoryRepository     syncEnygmaHistoryRepository
	enygmaCheckpointsRepository syncEnygmaCheckpointRepository
	enygmaResyncService         syncResync
	retryAttempts               map[string]int // key: checkpointId, value: attemptCount
}

func NewEnygmaSyncService(
	conf SyncConfig,
	txManager syncTransactionManager,
	enygmaRepository syncEnygmaRepository,
	enygmaHistoryRepository syncEnygmaHistoryRepository,
	enygmaCheckpointsRepository syncEnygmaCheckpointRepository,
	enygmaResyncService syncResync,
) *EnygmaSyncService {
	return &EnygmaSyncService{
		conf:                        conf,
		txManager:                   txManager,
		enygmaRepository:            enygmaRepository,
		enygmaHistoryRepository:     enygmaHistoryRepository,
		enygmaCheckpointsRepository: enygmaCheckpointsRepository,
		enygmaResyncService:         enygmaResyncService,
		retryAttempts:               make(map[string]int),
	}
}

func (s *EnygmaSyncService) Run(ctx context.Context) error {
	slog.Debug("Trying to finalize checkpoints...")

	checkpoints, err := s.enygmaCheckpointsRepository.GetValidationCandidates(ctx)
	if err != nil {
		return fmt.Errorf("get validation candidates: %w", err)
	}

	if len(checkpoints) == 0 {
		slog.Debug("No checkpoints to finalize")
		return nil
	}

	resourceIds := make([]string, 0, len(checkpoints))
	blockNumbers := make([]*big.Int, 0, len(checkpoints))

	for _, checkpoint := range checkpoints {
		resourceIds = append(resourceIds, checkpoint.ResourceId)
		blockNumbers = append(blockNumbers, checkpoint.FinalizedBlockNumber)
	}

	// Bulk get history records for all checkpoints
	historyRecords, err := s.enygmaHistoryRepository.GetEnygmaHistoryForCheckpoints(ctx, resourceIds, blockNumbers)
	if err != nil {
		return fmt.Errorf("get enygma history for checkpoints: %w", err)
	}
	// Bulk get enygma records for all unique resource IDs
	enygmaRecords, err := s.enygmaRepository.GetEnygmaByResourceIds(ctx, resourceIds)
	if err != nil {
		return fmt.Errorf("get enygma by resource IDs: %w", err)
	}

	// Create maps for quick lookup
	stateMap := make(map[string]types.Enygma)
	historyMap := make(map[string][]types.EnygmaHistory)

	for _, enygma := range enygmaRecords {
		stateMap[enygma.ResourceId] = enygma
	}

	for _, history := range historyRecords {
		historyMap[history.ResourceId] = append(historyMap[history.ResourceId], history)
	}

	for _, checkpoint := range checkpoints {
		slog.Debug(
			"Finalizing checkpoint",
			"resourceId",
			checkpoint.ResourceId,
			"blockNumberPrivateHub",
			checkpoint.FinalizedBlockNumber,
		)

		checkpointHistory, exists := historyMap[checkpoint.ResourceId]
		if !exists || len(checkpointHistory) == 0 {
			if retryErr := s.checkRetryAndResync(ctx, checkpoint, "No history found"); retryErr != nil {
				return fmt.Errorf("retry failed: %w", retryErr)
			}
			continue
		}

		lastState, exists := stateMap[checkpoint.ResourceId]
		if !exists {
			return fmt.Errorf("no enygma found with resourceId: %s", checkpoint.ResourceId)
		}

		newBalance, newRFactor := computeHistoryChanges(checkpointHistory)
		newFinalizedBalance, newFinalizedR := computeFinalizedValues(
			lastState.FinalizedBalance,
			lastState.FinalizedR,
			newBalance,
			newRFactor,
		)

		isCheckpointValid := validateCheckpointBalance(
			checkpoint.FinalizedPublicBalanceX,
			checkpoint.FinalizedPublicBalanceY,
			newFinalizedBalance,
			newFinalizedR,
		)

		if !isCheckpointValid {
			retryErr := s.checkRetryAndResync(ctx, checkpoint, "Checkpoint validation failed")
			if retryErr != nil {
				return fmt.Errorf("retry failed: %w", retryErr)
			}
			continue
		}

		err = s.finalizeCheckpoint(
			ctx,
			checkpoint.ResourceId,
			checkpoint.PendingBlockNumber,
			checkpoint.FinalizedBlockNumber,
			newFinalizedBalance,
			newFinalizedR,
		)
		if err != nil {
			return fmt.Errorf("finalize checkpoint: %w", err)
		}

		// Cleanup retry attempts for this checkpoint after successful finalization
		if s.retryAttempts[checkpoint.ID] > 0 {
			delete(s.retryAttempts, checkpoint.ID)
			slog.Debug(
				"Cleaned up retry attempts for finalized checkpoint",
				"checkpointId",
				checkpoint.ID,
				"resourceId",
				checkpoint.ResourceId,
			)
		}
	}

	return nil
}

func (s *EnygmaSyncService) finalizeCheckpoint(
	ctx context.Context,
	resourceId string,
	newPendingBlockNumber *big.Int,
	newFinalizedBlockNumber *big.Int,
	newFinalizedBalance, newFinalizedR *big.Int,
) error {
	err := s.txManager.WithTransaction(ctx, func(txCtx context.Context) error {
		err := s.enygmaCheckpointsRepository.MarkAsFinalized(
			txCtx,
			resourceId,
			newFinalizedBlockNumber,
		)
		if err != nil {
			return fmt.Errorf("mark checkpoint as finalized: %w", err)
		}

		err = s.enygmaRepository.UpdateEnygma(
			txCtx,
			resourceId,
			newFinalizedBalance,
			newFinalizedR,
			newFinalizedBlockNumber,
			newPendingBlockNumber,
		)
		if err != nil {
			return fmt.Errorf("update enygma: %w", err)
		}

		return nil
	})
	if err != nil {
		return fmt.Errorf("finalize checkpoint transaction: %w", err)
	}

	return nil
}

func computeHistoryChanges(history []types.EnygmaHistory) (newBalance *big.Int, newRFactor *big.Int) {
	newBalance = big.NewInt(0)
	newRFactor = big.NewInt(0)

	for _, h := range history {
		newBalance = newBalance.Add(newBalance, h.BalanceChange)
		newRFactor = cryptography.AddMod(newRFactor, h.RFactor, cryptography.JubJubPrimeSubGroup)
	}

	return newBalance, newRFactor
}

func computeFinalizedValues(
	lastFinalizedBalance *big.Int,
	lastFinalizedR *big.Int,
	newBalance *big.Int,
	newRFactor *big.Int,
) (newFinalizedBalance *big.Int, newFinalizedR *big.Int) {
	newFinalizedBalance = new(big.Int).Add(lastFinalizedBalance, newBalance)
	newFinalizedR = cryptography.AddMod(lastFinalizedR, newRFactor, cryptography.JubJubPrimeSubGroup)

	return newFinalizedBalance, newFinalizedR
}

// checkRetryAndResync checks if we should retry or force resync based on attempt count
// Always continues to next checkpoint after processing
func (s *EnygmaSyncService) checkRetryAndResync(
	ctx context.Context,
	checkpoint types.EnygmaCheckpoint,
	reason string,
) error {
	attempts := s.retryAttempts[checkpoint.ID]
	if attempts < s.conf.MaxRetries {
		// Increment attempt count and skip resync for now
		s.retryAttempts[checkpoint.ID] = attempts + 1
		slog.Debug(
			reason+", incrementing retry count",
			"checkpointId",
			checkpoint.ID,
			"resourceId",
			checkpoint.ResourceId,
			"attempt",
			attempts+1,
			"maxRetries",
			s.conf.MaxRetries,
		)
		return nil
	}

	// Max retries reached, force resync
	slog.Debug(
		"Max retries reached, forcing resync",
		"checkpointId",
		checkpoint.ID,
		"resourceId",
		checkpoint.ResourceId,
		"attempts",
		attempts,
		"reason",
		reason,
	)
	err := s.enygmaResyncService.ResyncEnygma(ctx, checkpoint.ResourceId)
	if err != nil {
		slog.Error(
			"Failed to resync enygma",
			slog.String("resourceId", checkpoint.ResourceId),
			slog.String("err", err.Error()),
		)
		return fmt.Errorf("resync enygma for resource %s: %w", checkpoint.ResourceId, err)
	}
	slog.Debug(
		"Checkpoint resync was finalized successfully",
		"resourceId",
		checkpoint.ResourceId,
		"blockNumberPrivateHub",
		checkpoint.FinalizedBlockNumber,
	)
	return nil
}

func validateCheckpointBalance(
	publicBalanceX *big.Int,
	publicBalanceY *big.Int,
	newFinalizedBalance *big.Int,
	newFinalizedR *big.Int,
) bool {
	expectedPublicBalance := cryptography.PedersenCommitmentEnygma(newFinalizedBalance, newFinalizedR)

	matchBalanceX := expectedPublicBalance.X.Cmp(publicBalanceX) == 0
	matchBalanceY := expectedPublicBalance.Y.Cmp(publicBalanceY) == 0

	return matchBalanceX && matchBalanceY
}

func ProcessEnygmaFinalizedBalances(
	ctx context.Context,
	enygmaFinalizedBalanceMsgs []*types.EnygmaFinalizedBalance,
	enygmaCheckpointRepository syncEnygmaCheckpointRepository,
	tracer syncTracer,
) error {
	ctx, span := tracer.Start(ctx, telemetry.SPAN_PROCESS_ENYGMA_FINALIZED_BALANCE)
	defer span.End()

	for _, msg := range enygmaFinalizedBalanceMsgs {

		existingCheckPoint, err := enygmaCheckpointRepository.GetLatestCheckpointByFilters(
			ctx,
			msg.ResourceId,
			nil,
			msg.FinalizedBlockNumber,
			msg.PendingBlockNumber,
		)
		if err != nil {
			slog.Error(
				"Failed to get the latest checkpoint for resourceId to check if it was already processed",
				slog.String("resourceId", msg.ResourceId),
				slog.String("finalizedBlockNumber", msg.FinalizedBlockNumber.String()),
				slog.String("pendingBlockNumber", msg.PendingBlockNumber.String()),
				slog.String("err", err.Error()),
			)
			continue
		}

		if existingCheckPoint != nil {
			slog.Info(
				"Checkpoint already processed",
				slog.String("resourceId", msg.ResourceId),
				slog.String("finalizedBlockNumber", msg.FinalizedBlockNumber.String()),
				slog.String("pendingBlockNumber", msg.PendingBlockNumber.String()),
			)
			continue
		}

		zeroX := big.NewInt(0)
		zeroY := big.NewInt(1)
		if msg.BalanceX.Cmp(zeroX) == 0 && msg.BalanceY.Cmp(zeroY) == 0 {
			slog.Debug(
				"Skipping checkpoint creation for zero balance",
				slog.String("resourceId", msg.ResourceId),
				slog.String("blockNumber", msg.FinalizedBlockNumber.String()),
			)
			continue
		}

		// TODO: can we do all of that in a single Upsert query?
		// Use GetLatestCheckpointByFilters with finalizedBlockNumber as upper bound to find
		// the immediately preceding checkpoint. Using (nil,nil,nil) would return the highest
		// block overall, which may be a checkpoint inserted earlier in this same batch loop.
		latestCheckpoint, err := enygmaCheckpointRepository.GetLatestCheckpointByFilters(
			ctx,
			msg.ResourceId,
			nil,
			msg.FinalizedBlockNumber,
			nil,
		)
		if err != nil {
			slog.Error(
				"Failed to get the latest checkpoint for resourceId",
				slog.String("resourceId", msg.ResourceId),
				slog.String("err", err.Error()),
			)
			continue
		}

		// The public balance of our PL was not changed since the last checkpoint(block).
		// This means we were not involved in any transfer since then, not even with 0 amount.
		// So, we don't need to create a new checkpoint
		if latestCheckpoint != nil && latestCheckpoint.FinalizedPublicBalanceX.Cmp(msg.BalanceX) == 0 &&
			latestCheckpoint.FinalizedPublicBalanceY.Cmp(msg.BalanceY) == 0 {
			slog.Info(
				"Public balance of our PL was not changed since the last checkpoint, so we don't need to create a new checkpoint",
				slog.String("resourceId", msg.ResourceId),
				slog.String("blockNumber", msg.FinalizedBlockNumber.String()),
			)
			continue
		}

		err = enygmaCheckpointRepository.CreateEnygmaCheckpoint(ctx, types.EnygmaCheckpoint{
			ResourceId:              msg.ResourceId,
			FinalizedPublicBalanceX: msg.BalanceX,
			FinalizedPublicBalanceY: msg.BalanceY,
			FinalizedBlockNumber:    msg.FinalizedBlockNumber,
			PendingBlockNumber:      msg.PendingBlockNumber,
			Status:                  types.EnygmaCheckpointStatusTentative,
		})
		if err != nil {
			slog.Error(
				"Failed to create Enygma checkpoint",
				slog.String("resourceId", msg.ResourceId),
				slog.String("blockNumber", msg.FinalizedBlockNumber.String()),
				slog.String("err", err.Error()),
			)
			continue
		}
	}

	return nil
}
