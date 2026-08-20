// Decommissioning Teleport (vanilla, atomic).

package poller

import (
	"context"
	"log/slog"
	"slices"
	"time"

	"github.com/raylsnetwork/rayls-sovereign-relayer/private-relayer/repository"
	"github.com/raylsnetwork/rayls-sovereign-relayer/types"
)

type finalizationSignatureService interface {
	HandleSourceExecuted(context.Context, []string) error
	HandleSourceReverted(context.Context, []string) error
}

type finalizationTransactionRepository interface {
	GetByStateOutcomeAndAtomicity(
		ctx context.Context,
		state types.TransactionState,
		outcome types.TransactionOutcome,
		isAtomic bool,
		opts ...repository.Option,
	) ([]types.Transaction, error)
	GetByStateOutcomesAndAtomicity(
		ctx context.Context,
		state types.TransactionState,
		outcomes []types.TransactionOutcome,
		isAtomic bool,
		opts ...repository.Option,
	) ([]types.Transaction, error)
}

type finalizationAtomicStatusRepository interface {
	GetUnprocessedBySharedIDs(
		context.Context,
		[]string,
		...repository.Option,
	) ([]types.AtomicStatusUpdateMessage, error)
}

// Deprecated: Decommissioning Teleport (vanilla, atomic).
type FinalizationPoller struct {
	batchSize int

	svc      finalizationSignatureService
	txRepo   finalizationTransactionRepository
	sumsRepo finalizationAtomicStatusRepository

	ticker     *time.Ticker
	initialRun chan struct{}
}

// Deprecated: Decommissioning Teleport (vanilla, atomic).
func NewFinalizationPoller(
	batchSize int,
	svc finalizationSignatureService,
	txRepo finalizationTransactionRepository,
	sumsRepo finalizationAtomicStatusRepository,
) *FinalizationPoller {
	initialRun := make(chan struct{}, 1)
	initialRun <- struct{}{}

	return &FinalizationPoller{
		batchSize: batchSize,

		svc:      svc,
		txRepo:   txRepo,
		sumsRepo: sumsRepo,

		ticker:     time.NewTicker(time.Second),
		initialRun: initialRun,
	}
}

// Deprecated: Decommissioning Teleport (vanilla, atomic).
func (p *FinalizationPoller) Run(ctx context.Context) error { //nolint:gocognit // finalization polling with retry logic
	defer p.ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return nil
		case <-p.ticker.C:
		case <-p.initialRun:
		}

		var (
			err       error
			sharedIDs []string
			sums      []types.AtomicStatusUpdateMessage

			executedSharedIDs []string
			revertedSharedIDs []string
		)

		txs, err := p.txRepo.GetByStateOutcomeAndAtomicity(ctx, types.SourcePublish, types.OutcomeSuccess, true)
		if err != nil {
			slog.Error(
				"finalization: error while trying to get transactions by state SourcePublish+success",
				slog.Any("error", err),
			)
			continue
		}
		failedTxs, err := p.txRepo.GetByStateOutcomesAndAtomicity(
			ctx,
			types.SourcePublish,
			[]types.TransactionOutcome{types.OutcomeReverted, types.OutcomeFailed},
			true,
		)
		if err != nil {
			slog.Error(
				"finalization: error while getting transactions in state SourcePublish+failed/reverted",
				slog.Any("error", err),
			)
			continue
		}
		revertedTxs, err := p.txRepo.GetByStateOutcomeAndAtomicity(ctx, types.SourceTimeoutRevert, types.OutcomeSuccess, true)
		if err != nil {
			slog.Error(
				"finalization: error while trying to get transactions in state SourceTimeoutRevert+success",
				slog.Any("error", err),
			)
			continue
		}

		if len(txs) == 0 && len(failedTxs) == 0 && len(revertedTxs) == 0 {
			continue
		}

		sharedIDs = getSharedIDsForTransactions(slices.Concat(txs, failedTxs, revertedTxs))
		sums, err = p.sumsRepo.GetUnprocessedBySharedIDs(ctx, sharedIDs, repository.WithLimit(p.batchSize))
		if err != nil {
			slog.Error("finalization: error while trying to get unprocessed shared IDs", slog.Any("error", err))
			continue
		}
		if len(sums) == 0 {
			continue
		}

		for _, sum := range sums {
			switch sum.Status {
			case types.AtomicExecutedStatus:
				executedSharedIDs = append(executedSharedIDs, sum.SharedID)
			case types.AtomicRevertedStatus:
				revertedSharedIDs = append(revertedSharedIDs, sum.SharedID)
			}
		}

		if len(executedSharedIDs) != 0 {
			slog.Info("Finalizing executed source transactions", slog.Int("COUNT", len(executedSharedIDs)))
			err = p.svc.HandleSourceExecuted(ctx, executedSharedIDs)
			if err != nil {
				slog.Error(
					"AtomicFinalization: error while trying to finalize source executed transactions",
					slog.Any("error", err),
				)
			} else {
				slog.Info("successfully executed source transactions", slog.Int("COUNT", len(executedSharedIDs)))
			}
		}

		if len(revertedSharedIDs) != 0 {
			slog.Info("Finalizing reverted source transactions", slog.Int("COUNT", len(revertedSharedIDs)))
			err = p.svc.HandleSourceReverted(ctx, revertedSharedIDs)
			if err != nil {
				slog.Error(
					"AtomicFinalization: error while trying to finalize source reverted transactions",
					slog.Any("error", err),
				)
			} else {
				slog.Info("Successfully reverted source transactions", slog.Int("COUNT", len(revertedSharedIDs)))
			}
		}
	}
}

func getSharedIDsForTransactions(txs []types.Transaction) []string {
	sharedIDs := make([]string, 0, len(txs))
	for _, tx := range txs {
		sharedIDs = append(sharedIDs, tx.SharedID)
	}

	return sharedIDs
}
