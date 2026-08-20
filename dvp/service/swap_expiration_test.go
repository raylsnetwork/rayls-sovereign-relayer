package service_test

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"
	"time"

	"github.com/raylsnetwork/rayls-privacy-relayer-api/dvp/service"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/dvp/testdata"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/testtools"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/types"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// pendingSwap thinly wraps testdata.NewDvpSwap so this suite's fixtures
// stay easy to read at the call sites. The expiration service only
// inspects SharedID; the rest of the row is realistic but irrelevant.
func pendingSwap(sharedID string, status types.DvpSwapStatus) *types.DvpSwap {
	return testdata.NewDvpSwap(
		testdata.WithSharedID(sharedID),
		testdata.WithStatus(status),
		testdata.WithCreatedAt(time.Now()),
	)
}

// runOnceAndCancel starts the service in a goroutine and cancels the context
// once `done` fires. It returns when the goroutine exits (or after a timeout).
// The expiration service runs an immediate first iteration before waiting on
// the ticker, so a single signal on `done` after the first poll is enough to
// observe one full pass.
func runOnceAndCancel(t *testing.T, svc *service.ExpirationService, done <-chan struct{}) {
	t.Helper()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	exitCh := make(chan struct{})
	go func() {
		_ = svc.Run(ctx)
		close(exitCh)
	}()

	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("test timed out waiting for first poll")
	}
	cancel()
	select {
	case <-exitCh:
	case <-time.After(time.Second):
		t.Fatal("service did not exit after context cancellation")
	}
}

func TestExpirationService_Run(t *testing.T) {
	testtools.SilenceLogger()

	t.Run("polls pending swaps and expires only those the contract reports as expired", func(t *testing.T) {
		swaps := []*types.DvpSwap{
			pendingSwap("expired-1", types.DvpSwapInitiated),
			pendingSwap("still-pending", types.DvpSwapInitiated),
			pendingSwap("expired-2", types.DvpSwapWaitingConfirmation),
		}
		expirationByID := map[string]bool{
			"expired-1":     true,
			"still-pending": false,
			"expired-2":     true,
		}

		var pollCount atomic.Int32
		done := make(chan struct{})
		repo := &ExpirationSwapRepositoryMock{
			GetPendingSwapsFunc: func(ctx context.Context) ([]*types.DvpSwap, error) {
				if pollCount.Add(1) == 1 {
					defer close(done)
				}
				return swaps, nil
			},
		}

		var expiredIDs []string
		client := &ExpirationDvpClientMock{
			IsSwapExpiredFunc: func(ctx context.Context, sharedId string) (bool, error) {
				return expirationByID[sharedId], nil
			},
			ExpireSwapFunc: func(ctx context.Context, sharedId string) error {
				expiredIDs = append(expiredIDs, sharedId)
				return nil
			},
		}

		svc := service.NewExpirationService(
			service.ExpirationConfig{Interval: time.Hour},
			repo,
			client,
		)
		runOnceAndCancel(t, svc, done)

		// First poll happened before context cancellation; subsequent ticker
		// iterations could have run if the test was slow, but the deterministic
		// minimum is one iteration.
		assert.GreaterOrEqual(t, pollCount.Load(), int32(1))

		// IsSwapExpired is invoked once per pending row in the first iteration.
		isExpiredCalls := client.IsSwapExpiredCalls()
		require.GreaterOrEqual(t, len(isExpiredCalls), 3)
		seen := map[string]bool{}
		for _, c := range isExpiredCalls[:3] {
			seen[c.SharedId] = true
		}
		assert.Equal(t, map[string]bool{"expired-1": true, "still-pending": true, "expired-2": true}, seen)

		// ExpireSwap is called only for the two expired rows.
		require.GreaterOrEqual(t, len(expiredIDs), 2)
		assert.ElementsMatch(t, []string{"expired-1", "expired-2"}, expiredIDs[:2])
	})

	t.Run("logs and continues when GetPendingSwaps fails", func(t *testing.T) {
		done := make(chan struct{})
		var first atomic.Bool
		repo := &ExpirationSwapRepositoryMock{
			GetPendingSwapsFunc: func(ctx context.Context) ([]*types.DvpSwap, error) {
				if !first.Swap(true) {
					defer close(done)
				}
				return nil, errors.New("db down")
			},
		}
		client := &ExpirationDvpClientMock{}

		svc := service.NewExpirationService(
			service.ExpirationConfig{Interval: time.Hour},
			repo,
			client,
		)
		runOnceAndCancel(t, svc, done)

		// IsSwapExpired and ExpireSwap must not be invoked when GetPendingSwaps errors.
		assert.Empty(t, client.IsSwapExpiredCalls())
		assert.Empty(t, client.ExpireSwapCalls())
	})

	t.Run("logs and continues when IsSwapExpired errors for a row", func(t *testing.T) {
		swaps := []*types.DvpSwap{
			pendingSwap("good", types.DvpSwapInitiated),
			pendingSwap("checker-fails", types.DvpSwapInitiated),
		}

		done := make(chan struct{})
		var first atomic.Bool
		repo := &ExpirationSwapRepositoryMock{
			GetPendingSwapsFunc: func(ctx context.Context) ([]*types.DvpSwap, error) {
				if !first.Swap(true) {
					defer close(done)
				}
				return swaps, nil
			},
		}
		client := &ExpirationDvpClientMock{
			IsSwapExpiredFunc: func(ctx context.Context, sharedId string) (bool, error) {
				if sharedId == "checker-fails" {
					return false, errors.New("rpc timeout")
				}
				return true, nil
			},
			ExpireSwapFunc: func(ctx context.Context, sharedId string) error {
				return nil
			},
		}

		svc := service.NewExpirationService(
			service.ExpirationConfig{Interval: time.Hour},
			repo,
			client,
		)
		runOnceAndCancel(t, svc, done)

		// "good" must be expired; "checker-fails" must NOT trigger an ExpireSwap.
		expireCalls := client.ExpireSwapCalls()
		require.NotEmpty(t, expireCalls)
		var sawCheckerFails bool
		for _, c := range expireCalls {
			if c.SharedId == "checker-fails" {
				sawCheckerFails = true
			}
		}
		assert.False(t, sawCheckerFails, "checker-fails must not be expired")
	})

	t.Run("continues to next swap when ExpireSwap errors", func(t *testing.T) {
		swaps := []*types.DvpSwap{
			pendingSwap("expire-fails", types.DvpSwapInitiated),
			pendingSwap("expire-ok", types.DvpSwapInitiated),
		}

		done := make(chan struct{})
		var first atomic.Bool
		repo := &ExpirationSwapRepositoryMock{
			GetPendingSwapsFunc: func(ctx context.Context) ([]*types.DvpSwap, error) {
				if !first.Swap(true) {
					defer close(done)
				}
				return swaps, nil
			},
		}
		client := &ExpirationDvpClientMock{
			IsSwapExpiredFunc: func(ctx context.Context, sharedId string) (bool, error) {
				return true, nil
			},
			ExpireSwapFunc: func(ctx context.Context, sharedId string) error {
				if sharedId == "expire-fails" {
					return errors.New("expire failed")
				}
				return nil
			},
		}

		svc := service.NewExpirationService(
			service.ExpirationConfig{Interval: time.Hour},
			repo,
			client,
		)
		runOnceAndCancel(t, svc, done)

		expireCalls := client.ExpireSwapCalls()
		require.GreaterOrEqual(t, len(expireCalls), 2)
		ids := []string{expireCalls[0].SharedId, expireCalls[1].SharedId}
		assert.ElementsMatch(t, []string{"expire-fails", "expire-ok"}, ids)
	})

	t.Run("returns nil when context is cancelled before first iteration", func(t *testing.T) {
		repo := &ExpirationSwapRepositoryMock{
			GetPendingSwapsFunc: func(ctx context.Context) ([]*types.DvpSwap, error) {
				return nil, nil
			},
		}
		client := &ExpirationDvpClientMock{}

		svc := service.NewExpirationService(
			service.ExpirationConfig{Interval: time.Hour},
			repo,
			client,
		)

		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		err := svc.Run(ctx)
		require.NoError(t, err)
	})
}
