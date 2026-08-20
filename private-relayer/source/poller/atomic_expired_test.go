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

type StubExpiredExpiredService struct {
	spyCalledHandle bool
	spySharedIDs    []string
}

func (s *StubExpiredExpiredService) Handle(ctx context.Context, sharedIDs []string) error {
	s.spyCalledHandle = true
	s.spySharedIDs = sharedIDs
	return nil
}

type StubExpiredTransactionRepository struct {
	txs []types.Transaction

	spyState    types.TransactionState
	spyOutcome  types.TransactionOutcome
	spyIsAtomic bool
}

func (r *StubExpiredTransactionRepository) GetByStateOutcomeAndAtomicity(
	ctx context.Context,
	state types.TransactionState,
	outcome types.TransactionOutcome,
	isAtomic bool,
	opts ...repository.Option,
) ([]types.Transaction, error) {
	r.spyState = state
	r.spyOutcome = outcome
	r.spyIsAtomic = isAtomic

	return r.txs, nil
}

func TestExpiredPoller(t *testing.T) {
	batchSize := 100

	t.Run("exits gracefully on canceled context", func(t *testing.T) {
		svc := &StubExpiredExpiredService{}
		txRepo := &StubExpiredTransactionRepository{}
		pol := poller.NewExpiredPoller(batchSize, time.Second, svc, txRepo)

		hasGracefulShutdown := testtools.ShutdownFixture(t, pol.Run, time.Millisecond)

		require.True(t, hasGracefulShutdown, "poller doesn't support graceful shutdown")
	})

	t.Run("reverts expired transactions", func(t *testing.T) {
		wantSharedIDs := []string{"first-shared-id", "second-shared-id"}

		wantState := types.SourcePublish
		wantOutcome := types.OutcomeSuccess
		wantIsAtomic := true

		txs := []types.Transaction{
			{
				SharedID:  wantSharedIDs[0],
				UpdatedAt: time.Now().Add(-time.Second),
			},
			{
				SharedID:  wantSharedIDs[1],
				UpdatedAt: time.Now().Add(-time.Second),
			},
			{
				SharedID:  "not-expired-shared-id",
				UpdatedAt: time.Now(),
			},

			{
				SharedID:  "second-not-expired-shared-id",
				UpdatedAt: time.Now(),
			},
		}

		svc := &StubExpiredExpiredService{}
		txRepo := &StubExpiredTransactionRepository{
			txs: txs,
		}
		pol := poller.NewExpiredPoller(batchSize, time.Second/2, svc, txRepo)

		ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond)
		defer cancel()

		err := pol.Run(ctx)
		require.Nil(t, err)

		assert.Equal(t, wantState, txRepo.spyState, "didn't request transactions with correct state from repository")
		assert.Equal(t, wantOutcome, txRepo.spyOutcome, "didn't request transactions with correct outcome from repository")
		assert.Equal(
			t,
			wantIsAtomic,
			txRepo.spyIsAtomic,
			"didn't request transactions with correct atomicity from repository",
		)

		assert.ElementsMatch(t, wantSharedIDs, svc.spySharedIDs, "didn't revert correct shared IDs")
	})

	t.Run("continues on no transactions to process", func(t *testing.T) {
		svc := &StubExpiredExpiredService{}
		txRepo := &StubExpiredTransactionRepository{}
		pol := poller.NewExpiredPoller(batchSize, time.Second/2, svc, txRepo)

		ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond)
		defer cancel()

		err := pol.Run(ctx)
		require.Nil(t, err)

		assert.False(t, svc.spyCalledHandle, "called expired revert with no transactions to process")
	})
}
