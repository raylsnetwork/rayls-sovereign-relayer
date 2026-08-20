// Decommissioning Teleport (vanilla, atomic).

package poller_test

import (
	"context"
	"testing"
	"time"

	"github.com/raylsnetwork/rayls-sovereign-relayer/private-relayer/repository"
	"github.com/raylsnetwork/rayls-sovereign-relayer/private-relayer/source/poller"
	"github.com/raylsnetwork/rayls-sovereign-relayer/testtools"
	"github.com/raylsnetwork/rayls-sovereign-relayer/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type StubEarlyRevertTransactionRepository struct {
	txs         []types.Transaction
	spyState    types.TransactionState
	spyOutcomes []types.TransactionOutcome
	spyIsAtomic bool
}

func (r *StubEarlyRevertTransactionRepository) GetByStateOutcomesAndAtomicity(
	ctx context.Context,
	state types.TransactionState,
	outcomes []types.TransactionOutcome,
	isAtomic bool,
	opts ...repository.Option,
) ([]types.Transaction, error) {
	r.spyState = state
	r.spyOutcomes = outcomes
	r.spyIsAtomic = isAtomic
	return r.txs, nil
}

type StubEarlyRevertService struct {
	spyCalledHandleEarlyRevert bool
	spySharedIDs               []string
}

func (s *StubEarlyRevertService) HandleEarlyRevert(ctx context.Context, sharedIDs []string) error {
	s.spyCalledHandleEarlyRevert = true
	s.spySharedIDs = sharedIDs
	return nil
}

func TestEarlyRevertPoller(t *testing.T) {
	batchSize := 100

	t.Run("exits gracefully on canceled context", func(t *testing.T) {
		txRepo := &StubEarlyRevertTransactionRepository{}
		svc := &StubEarlyRevertService{}
		pol := poller.NewEarlyRevertPoller(batchSize, svc, txRepo)

		hasGracefulShutdown := testtools.ShutdownFixture(t, pol.Run, time.Millisecond)
		require.True(t, hasGracefulShutdown, "poller doesn't support graceful shutdown")
	})

	t.Run("reverts transactions whose source publish reverted or failed", func(t *testing.T) {
		wantSharedIDs := []string{"first-shared-id", "second-shared-id"}
		wantState := types.SourcePublish
		wantOutcomes := []types.TransactionOutcome{types.OutcomeReverted, types.OutcomeFailed}
		wantIsAtomic := true

		txs := []types.Transaction{
			{SharedID: wantSharedIDs[0]},
			{SharedID: wantSharedIDs[1]},
		}

		txRepo := &StubEarlyRevertTransactionRepository{txs: txs}
		svc := &StubEarlyRevertService{}
		pol := poller.NewEarlyRevertPoller(batchSize, svc, txRepo)
		ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond)
		defer cancel()

		err := pol.Run(ctx)
		require.NoError(t, err)

		assert.Equal(t, wantState, txRepo.spyState, "didn't request transactions with correct state from repository")
		assert.ElementsMatch(t, wantOutcomes, txRepo.spyOutcomes, "didn't request transactions with correct outcomes")
		assert.Equal(
			t,
			wantIsAtomic,
			txRepo.spyIsAtomic,
			"didn't request transactions with correct atomicity from repository",
		)
		assert.ElementsMatch(t, wantSharedIDs, svc.spySharedIDs, "reverted wrong shared IDs")
	})
}
