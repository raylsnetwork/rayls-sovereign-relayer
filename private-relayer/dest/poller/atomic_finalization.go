// Decommissioning Teleport (vanilla, atomic).

package poller

import (
	"context"
	"log/slog"
	"time"

	"github.com/raylsnetwork/rayls-privacy-relayer-api/private-relayer/repository"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/types"
)

// Deprecated: Decommissioning Teleport (vanilla, atomic).
type FinalizationSignatureService interface {
	HandleDestinationExecuted(context.Context, []string) error
	HandleDestinationReverted(context.Context, []string) error
}

// Deprecated: Decommissioning Teleport (vanilla, atomic).
type FinalizationTransactionRepository interface {
	GetByStateAndAtomicity(
		ctx context.Context,
		state types.TransactionState,
		isAtomic bool,
		opts ...repository.Option,
	) ([]types.Transaction, error)
}

// Deprecated: Decommissioning Teleport (vanilla, atomic).
type FinalizationAtomicStatusRepository interface {
	GetUnprocessedBySharedIDs(
		context.Context,
		[]string,
		...repository.Option,
	) ([]types.AtomicStatusUpdateMessage, error)
}

// Deprecated: Decommissioning Teleport (vanilla, atomic).
type FinalizationPoller struct {
	batchSize int

	signatureService       FinalizationSignatureService
	txRepo                 FinalizationTransactionRepository
	atomicStatusRepository FinalizationAtomicStatusRepository

	ticker     *time.Ticker
	initialRun chan struct{}
}

// Deprecated: Decommissioning Teleport (vanilla, atomic).
func NewFinalizationPoller(
	batchSize int,
	signatureService FinalizationSignatureService,
	txRepo FinalizationTransactionRepository,
	atomicStatusRepository FinalizationAtomicStatusRepository,
) *FinalizationPoller {
	initialRun := make(chan struct{}, 1)
	initialRun <- struct{}{}

	return &FinalizationPoller{
		batchSize: batchSize,

		signatureService:       signatureService,
		txRepo:                 txRepo,
		atomicStatusRepository: atomicStatusRepository,

		ticker:     time.NewTicker(time.Second),
		initialRun: initialRun,
	}
}

// Deprecated: Decommissioning Teleport (vanilla, atomic).
func (p *FinalizationPoller) Run(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		case <-p.ticker.C:
		case <-p.initialRun:
		}

		// Pick up every HubNotifiedExec row regardless of outcome — branching
		// happens in code below. HubNotifiedRevert is intentionally NOT polled:
		// dispatch failed there, no mint exists, nothing for dest to undo (the
		// source side independently handles its escrow release).
		txs, err := p.txRepo.GetByStateAndAtomicity(ctx, types.HubNotifiedExec, true)
		if err != nil {
			slog.Error(
				"finalization: error while trying to get transactions by state and atomicity",
				slog.Any("error", err),
			)
			continue
		}
		if len(txs) == 0 {
			continue
		}

		// Partition by outcome:
		//   - success: PNH-notification call returned OK; consult atomic_status to
		//     learn PNH's verdict (Executed → unlock sigs, Reverted → race lost,
		//     push dest-revert sigs to undo our orphan mint).
		//   - reverted: PNH-notification call returned ErrAlreadyReverted (early
		//     race detection — PNH already considers this batch reverted). No need
		//     to consult atomic_status; route directly to dest-revert sigs.
		var (
			pnhSuccessIDs   []string
			earlyRaceLostIDs []string
		)
		for _, tx := range txs {
			switch tx.Outcome {
			case types.OutcomeSuccess:
				pnhSuccessIDs = append(pnhSuccessIDs, tx.SharedID)
			case types.OutcomeReverted:
				earlyRaceLostIDs = append(earlyRaceLostIDs, tx.SharedID)
			}
		}

		var (
			executedSharedIDs []string
			revertedSharedIDs []string
		)

		if len(pnhSuccessIDs) != 0 {
			atomicStatusUpdateMessages, err := p.atomicStatusRepository.GetUnprocessedBySharedIDs(
				ctx,
				pnhSuccessIDs,
				repository.WithLimit(p.batchSize),
			)
			if err != nil {
				slog.Error("finalization: error while trying to get unprocessed shared IDs", slog.Any("error", err))
				continue
			}

			for _, sum := range atomicStatusUpdateMessages {
				switch sum.Status {
				case types.AtomicExecutedStatus:
					executedSharedIDs = append(executedSharedIDs, sum.SharedID)
				case types.AtomicRevertedStatus:
					// Late race-loss: we minted, told PNH "executed", but PNH ended
					// up reverted. Same recovery as the early-race case below.
					revertedSharedIDs = append(revertedSharedIDs, sum.SharedID)
				}
			}
		}

		// Early-race rows already know PNH's verdict (reverted) from the outcome
		// itself; no atomic_status round-trip needed.
		revertedSharedIDs = append(revertedSharedIDs, earlyRaceLostIDs...)

		if len(executedSharedIDs) != 0 {
			slog.Info("Finalizing successful destination transactions", slog.Int("COUNT", len(executedSharedIDs)))
			err = p.signatureService.HandleDestinationExecuted(ctx, executedSharedIDs)
			if err != nil {
				slog.Error(
					"finalization: error while trying to finalize destination executed transactions",
					slog.Any("error", err),
				)
			}
		}

		if len(revertedSharedIDs) != 0 {
			slog.Info("Finalizing reverted destination transactions", slog.Int("COUNT", len(revertedSharedIDs)))
			err = p.signatureService.HandleDestinationReverted(ctx, revertedSharedIDs)
			if err != nil {
				slog.Error(
					"finalization: error while trying to finalize destination reverted transactions",
					slog.Any("error", err),
				)
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
