package service_test

import (
	"context"
	"fmt"
	"math/big"
	"testing"

	"github.com/raylsnetwork/rayls-privacy-relayer-api/enygma/service"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/enygma/testutils"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProcessEnygmaFinalizedBalances(t *testing.T) {
	ctx := context.Background()

	t.Run("successfully creates checkpoint", func(t *testing.T) {
		repo := &MockEnygmaCheckpointRepository{
			GetLatestCheckpointByFiltersFunc: func(ctx context.Context, resourceId string, status *types.EnygmaCheckpointStatus, finalizedBlockNumber *big.Int, pendingBlockNumber *big.Int) (*types.EnygmaCheckpoint, error) {
				//nolint:nilnil // intentional nil return in test mock - no existing checkpoint
				return nil, nil
			},
			CreateEnygmaCheckpointFunc: func(ctx context.Context, checkpoint types.EnygmaCheckpoint) error {
				return nil // Success
			},
		}
		tracer := &testutils.MockTracer{}

		msg := createEnygmaFinalizedBalance("resource1")
		messages := []*types.EnygmaFinalizedBalance{msg}

		err := service.ProcessEnygmaFinalizedBalances(ctx, messages, repo, tracer)

		require.NoError(t, err)
		assert.Len(t, repo.CreateEnygmaCheckpointCalls.Calls, 1)
		createdCheckpoint := repo.CreateEnygmaCheckpointCalls.Calls[0].Checkpoint
		assert.Equal(t, "resource1", createdCheckpoint.ResourceId)
		assert.Equal(t, msg.FinalizedBlockNumber, createdCheckpoint.FinalizedBlockNumber)
		assert.Equal(t, msg.PendingBlockNumber, createdCheckpoint.PendingBlockNumber)
		assert.Equal(t, msg.BalanceX, createdCheckpoint.FinalizedPublicBalanceX)
		assert.Equal(t, msg.BalanceY, createdCheckpoint.FinalizedPublicBalanceY)
		assert.Equal(t, types.EnygmaCheckpointStatusTentative, createdCheckpoint.Status)
	})

	t.Run("processes multiple messages sequentially", func(t *testing.T) {
		repo := &MockEnygmaCheckpointRepository{
			GetLatestCheckpointByFiltersFunc: func(ctx context.Context, resourceId string, status *types.EnygmaCheckpointStatus, finalizedBlockNumber *big.Int, pendingBlockNumber *big.Int) (*types.EnygmaCheckpoint, error) {
				return nil, nil //nolint:nilnil // intentional nil return in test mock
			},
			CreateEnygmaCheckpointFunc: func(ctx context.Context, checkpoint types.EnygmaCheckpoint) error {
				return nil
			},
		}
		tracer := &testutils.MockTracer{}

		msg1 := createEnygmaFinalizedBalance("resource1")
		msg2 := createEnygmaFinalizedBalance("resource2")
		messages := []*types.EnygmaFinalizedBalance{msg1, msg2}

		err := service.ProcessEnygmaFinalizedBalances(ctx, messages, repo, tracer)

		require.NoError(t, err)
		assert.Len(t, repo.CreateEnygmaCheckpointCalls.Calls, 2)
		assert.Equal(t, "resource1", repo.CreateEnygmaCheckpointCalls.Calls[0].Checkpoint.ResourceId)
		assert.Equal(t, "resource2", repo.CreateEnygmaCheckpointCalls.Calls[1].Checkpoint.ResourceId)
	})

	t.Run("skips checkpoint already processed", func(t *testing.T) {
		existingCheckpoint := createEnygmaCheckpoint(
			"id1",
			"resource1",
			big.NewInt(12345),
			big.NewInt(67890),
			big.NewInt(100),
			big.NewInt(101),
		)

		repo := &MockEnygmaCheckpointRepository{
			GetLatestCheckpointByFiltersFunc: func(ctx context.Context, resourceId string, status *types.EnygmaCheckpointStatus, finalizedBlockNumber *big.Int, pendingBlockNumber *big.Int) (*types.EnygmaCheckpoint, error) {
				return existingCheckpoint, nil
			},
		}
		tracer := &testutils.MockTracer{}

		msg := createEnygmaFinalizedBalance("resource1")
		messages := []*types.EnygmaFinalizedBalance{msg}

		err := service.ProcessEnygmaFinalizedBalances(ctx, messages, repo, tracer)

		require.NoError(t, err)
		assert.Len(t, repo.CreateEnygmaCheckpointCalls.Calls, 0)
	})

	t.Run("skips zero balance checkpoint", func(t *testing.T) {
		repo := &MockEnygmaCheckpointRepository{
			GetLatestCheckpointByFiltersFunc: func(ctx context.Context, resourceId string, status *types.EnygmaCheckpointStatus, finalizedBlockNumber *big.Int, pendingBlockNumber *big.Int) (*types.EnygmaCheckpoint, error) {
				return nil, nil //nolint:nilnil // intentional nil return in test mock
			},
		}
		tracer := &testutils.MockTracer{}

		// Zero point on the curve: (0, 1)
		msg := createEnygmaFinalizedBalanceWithValues(
			"resource1",
			big.NewInt(0),
			big.NewInt(1),
			big.NewInt(100),
			big.NewInt(101),
		)
		messages := []*types.EnygmaFinalizedBalance{msg}

		err := service.ProcessEnygmaFinalizedBalances(ctx, messages, repo, tracer)

		require.NoError(t, err)
		assert.Len(t, repo.CreateEnygmaCheckpointCalls.Calls, 0)
	})

	t.Run("skips unchanged public balance", func(t *testing.T) {
		balanceX := big.NewInt(12345)
		balanceY := big.NewInt(67890)
		latestCheckpoint := createEnygmaCheckpoint(
			"id1",
			"resource1",
			balanceX,
			balanceY,
			big.NewInt(99),
			big.NewInt(100),
		)

		repo := &MockEnygmaCheckpointRepository{
			GetLatestCheckpointByFiltersFunc: func(ctx context.Context, resourceId string, status *types.EnygmaCheckpointStatus, finalizedBlockNumber *big.Int, pendingBlockNumber *big.Int) (*types.EnygmaCheckpoint, error) {
				// First call checks exact checkpoint existence
				if finalizedBlockNumber != nil && pendingBlockNumber != nil {
					return nil, nil //nolint:nilnil // intentional nil return in test mock
				}
				// Second call gets latest checkpoint
				return latestCheckpoint, nil
			},
		}
		tracer := &testutils.MockTracer{}

		// Same balance as latest checkpoint
		msg := createEnygmaFinalizedBalanceWithValues("resource1", balanceX, balanceY, big.NewInt(100), big.NewInt(101))
		messages := []*types.EnygmaFinalizedBalance{msg}

		err := service.ProcessEnygmaFinalizedBalances(ctx, messages, repo, tracer)

		require.NoError(t, err)
		assert.Len(t, repo.CreateEnygmaCheckpointCalls.Calls, 0)
	})

	t.Run("continues on checkpoint existence query error", func(t *testing.T) {
		repo := &MockEnygmaCheckpointRepository{
			GetLatestCheckpointByFiltersFunc: func(ctx context.Context, resourceId string, status *types.EnygmaCheckpointStatus, finalizedBlockNumber *big.Int, pendingBlockNumber *big.Int) (*types.EnygmaCheckpoint, error) {
				return nil, fmt.Errorf("database error")
			},
		}
		tracer := &testutils.MockTracer{}

		msg := createEnygmaFinalizedBalance("resource1")
		messages := []*types.EnygmaFinalizedBalance{msg}

		err := service.ProcessEnygmaFinalizedBalances(ctx, messages, repo, tracer)

		require.NoError(t, err)
		assert.Len(t, repo.CreateEnygmaCheckpointCalls.Calls, 0)
	})

	t.Run("continues on checkpoint creation error", func(t *testing.T) {
		repo := &MockEnygmaCheckpointRepository{
			GetLatestCheckpointByFiltersFunc: func(ctx context.Context, resourceId string, status *types.EnygmaCheckpointStatus, finalizedBlockNumber *big.Int, pendingBlockNumber *big.Int) (*types.EnygmaCheckpoint, error) {
				return nil, nil //nolint:nilnil // intentional nil return in test mock
			},
			CreateEnygmaCheckpointFunc: func(ctx context.Context, checkpoint types.EnygmaCheckpoint) error {
				return fmt.Errorf("creation failed")
			},
		}
		tracer := &testutils.MockTracer{}

		msg := createEnygmaFinalizedBalance("resource1")
		messages := []*types.EnygmaFinalizedBalance{msg}

		err := service.ProcessEnygmaFinalizedBalances(ctx, messages, repo, tracer)

		require.NoError(t, err)
	})

	t.Run("handles empty input array", func(t *testing.T) {
		repo := &MockEnygmaCheckpointRepository{}
		tracer := &testutils.MockTracer{}

		err := service.ProcessEnygmaFinalizedBalances(ctx, []*types.EnygmaFinalizedBalance{}, repo, tracer)

		require.NoError(t, err)
		assert.Len(t, repo.CreateEnygmaCheckpointCalls.Calls, 0)
	})

	t.Run("processes multiple messages for same resource", func(t *testing.T) {
		repo := &MockEnygmaCheckpointRepository{
			GetLatestCheckpointByFiltersFunc: func(ctx context.Context, resourceId string, status *types.EnygmaCheckpointStatus, finalizedBlockNumber *big.Int, pendingBlockNumber *big.Int) (*types.EnygmaCheckpoint, error) {
				return nil, nil //nolint:nilnil // intentional nil return in test mock
			},
			CreateEnygmaCheckpointFunc: func(ctx context.Context, checkpoint types.EnygmaCheckpoint) error {
				return nil
			},
		}
		tracer := &testutils.MockTracer{}

		msg1 := createEnygmaFinalizedBalanceWithValues(
			"resource1",
			big.NewInt(100),
			big.NewInt(200),
			big.NewInt(100),
			big.NewInt(101),
		)
		msg2 := createEnygmaFinalizedBalanceWithValues(
			"resource1",
			big.NewInt(300),
			big.NewInt(400),
			big.NewInt(102),
			big.NewInt(103),
		)
		messages := []*types.EnygmaFinalizedBalance{msg1, msg2}

		err := service.ProcessEnygmaFinalizedBalances(ctx, messages, repo, tracer)

		require.NoError(t, err)
		assert.Len(t, repo.CreateEnygmaCheckpointCalls.Calls, 2)
		assert.Equal(t, big.NewInt(100), repo.CreateEnygmaCheckpointCalls.Calls[0].Checkpoint.FinalizedPublicBalanceX)
		assert.Equal(t, big.NewInt(300), repo.CreateEnygmaCheckpointCalls.Calls[1].Checkpoint.FinalizedPublicBalanceX)
	})
}

func TestEnygmaSyncService_Run(t *testing.T) {
	ctx := context.Background()
	syncConfig := service.SyncConfig{MaxRetries: 2}

	t.Run("returns early when no validation candidates exist", func(t *testing.T) {
		checkpointRepo := &MockEnygmaCheckpointRepository{
			GetValidationCandidatesFunc: func(ctx context.Context) ([]types.EnygmaCheckpoint, error) {
				return []types.EnygmaCheckpoint{}, nil
			},
		}
		historyRepo := &MockEnygmaHistoryRepository{}
		enygmaRepo := &MockEnygmaRepository{}
		txManager := &MockTransactionManager{}
		resyncService := &MockResyncService{}

		svc := service.NewEnygmaSyncService(
			syncConfig,
			txManager,
			enygmaRepo,
			historyRepo,
			checkpointRepo,
			resyncService,
		)

		err := svc.Run(ctx)

		require.NoError(t, err)
		assert.Len(t, checkpointRepo.GetValidationCandidatesCalls.Calls, 1)
	})

	t.Run("returns error when GetValidationCandidates fails", func(t *testing.T) {
		checkpointRepo := &MockEnygmaCheckpointRepository{
			GetValidationCandidatesFunc: func(ctx context.Context) ([]types.EnygmaCheckpoint, error) {
				return nil, fmt.Errorf("database error")
			},
		}
		historyRepo := &MockEnygmaHistoryRepository{}
		enygmaRepo := &MockEnygmaRepository{}
		txManager := &MockTransactionManager{}
		resyncService := &MockResyncService{}

		svc := service.NewEnygmaSyncService(
			syncConfig,
			txManager,
			enygmaRepo,
			historyRepo,
			checkpointRepo,
			resyncService,
		)

		err := svc.Run(ctx)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "database error")
	})

	t.Run("returns error when GetEnygmaHistoryForCheckpoints fails", func(t *testing.T) {
		checkpoint := createEnygmaCheckpoint(
			"id1",
			"resource1",
			big.NewInt(100),
			big.NewInt(200),
			big.NewInt(100),
			big.NewInt(101),
		)

		checkpointRepo := &MockEnygmaCheckpointRepository{
			GetValidationCandidatesFunc: func(ctx context.Context) ([]types.EnygmaCheckpoint, error) {
				return []types.EnygmaCheckpoint{*checkpoint}, nil
			},
		}
		historyRepo := &MockEnygmaHistoryRepository{
			GetEnygmaHistoryForCheckpointsFunc: func(ctx context.Context, resourceIds []string, blockNumbers []*big.Int) ([]types.EnygmaHistory, error) {
				return nil, fmt.Errorf("history fetch error")
			},
		}
		enygmaRepo := &MockEnygmaRepository{}
		txManager := &MockTransactionManager{}
		resyncService := &MockResyncService{}

		svc := service.NewEnygmaSyncService(
			syncConfig,
			txManager,
			enygmaRepo,
			historyRepo,
			checkpointRepo,
			resyncService,
		)

		err := svc.Run(ctx)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "history fetch error")
	})

	t.Run("returns error when GetEnygmaByResourceIds fails", func(t *testing.T) {
		checkpoint := createEnygmaCheckpoint(
			"id1",
			"resource1",
			big.NewInt(100),
			big.NewInt(200),
			big.NewInt(100),
			big.NewInt(101),
		)

		checkpointRepo := &MockEnygmaCheckpointRepository{
			GetValidationCandidatesFunc: func(ctx context.Context) ([]types.EnygmaCheckpoint, error) {
				return []types.EnygmaCheckpoint{*checkpoint}, nil
			},
		}
		history := createEnygmaHistory("resource1", big.NewInt(50), big.NewInt(10))
		historyRepo := &MockEnygmaHistoryRepository{
			GetEnygmaHistoryForCheckpointsFunc: func(ctx context.Context, resourceIds []string, blockNumbers []*big.Int) ([]types.EnygmaHistory, error) {
				return []types.EnygmaHistory{*history}, nil
			},
		}
		enygmaRepo := &MockEnygmaRepository{
			GetEnygmaByResourceIdsFunc: func(ctx context.Context, resourceIds []string) ([]types.Enygma, error) {
				return nil, fmt.Errorf("enygma fetch error")
			},
		}
		txManager := &MockTransactionManager{}
		resyncService := &MockResyncService{}

		svc := service.NewEnygmaSyncService(
			syncConfig,
			txManager,
			enygmaRepo,
			historyRepo,
			checkpointRepo,
			resyncService,
		)

		err := svc.Run(ctx)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "enygma fetch error")
	})

	t.Run("skips finalization when no history found and resyncs not needed yet", func(t *testing.T) {
		checkpoint := createEnygmaCheckpoint(
			"id1",
			"resource1",
			big.NewInt(100),
			big.NewInt(200),
			big.NewInt(100),
			big.NewInt(101),
		)

		checkpointRepo := &MockEnygmaCheckpointRepository{
			MarkAsFinalizedFunc: func(ctx context.Context, resourceId string, finalizedBlockNumber *big.Int) error {
				assert.Fail(t, "MarkAsFinalized should not be called")
				return nil
			},
			GetValidationCandidatesFunc: func(ctx context.Context) ([]types.EnygmaCheckpoint, error) {
				return []types.EnygmaCheckpoint{*checkpoint}, nil
			},
		}
		historyRepo := &MockEnygmaHistoryRepository{
			GetEnygmaHistoryForCheckpointsFunc: func(ctx context.Context, resourceIds []string, blockNumbers []*big.Int) ([]types.EnygmaHistory, error) {
				return []types.EnygmaHistory{}, nil // No history
			},
		}
		enygmaRepo := &MockEnygmaRepository{
			UpdateEnygmaFunc: func(ctx context.Context, resourceId string, finalizedBalance, finalizedR, finalizedBlockNumber, pendingBlockNumber *big.Int) error {
				assert.Fail(t, "UpdateEnygma should not be called")
				return nil
			},
			GetEnygmaByResourceIdsFunc: func(ctx context.Context, resourceIds []string) ([]types.Enygma, error) {
				return []types.Enygma{}, nil
			},
		}
		txManager := &MockTransactionManager{}
		resyncService := &MockResyncService{}

		svc := service.NewEnygmaSyncService(
			syncConfig,
			txManager,
			enygmaRepo,
			historyRepo,
			checkpointRepo,
			resyncService,
		)

		err := svc.Run(ctx)

		require.NoError(t, err)
		// Should not finalize since no history was found
		assert.Len(t, txManager.WithTransactionCalls.Calls, 0)
		assert.Len(t, checkpointRepo.MarkAsFinalizedCalls.Calls, 0)
		assert.Len(t, enygmaRepo.UpdateEnygmaCalls.Calls, 0)
	})

	t.Run("returns error when enygma state is not found", func(t *testing.T) {
		checkpoint := createEnygmaCheckpoint(
			"id1",
			"resource1",
			big.NewInt(100),
			big.NewInt(200),
			big.NewInt(100),
			big.NewInt(101),
		)

		checkpointRepo := &MockEnygmaCheckpointRepository{
			GetValidationCandidatesFunc: func(ctx context.Context) ([]types.EnygmaCheckpoint, error) {
				return []types.EnygmaCheckpoint{*checkpoint}, nil
			},
		}
		history := createEnygmaHistory("resource1", big.NewInt(50), big.NewInt(10))
		historyRepo := &MockEnygmaHistoryRepository{
			GetEnygmaHistoryForCheckpointsFunc: func(ctx context.Context, resourceIds []string, blockNumbers []*big.Int) ([]types.EnygmaHistory, error) {
				return []types.EnygmaHistory{*history}, nil
			},
		}
		enygmaRepo := &MockEnygmaRepository{
			GetEnygmaByResourceIdsFunc: func(ctx context.Context, resourceIds []string) ([]types.Enygma, error) {
				return []types.Enygma{}, nil // Empty list
			},
		}
		txManager := &MockTransactionManager{}
		resyncService := &MockResyncService{}

		svc := service.NewEnygmaSyncService(
			syncConfig,
			txManager,
			enygmaRepo,
			historyRepo,
			checkpointRepo,
			resyncService,
		)

		err := svc.Run(ctx)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "no enygma found")
	})

	t.Run("skips finalization when checkpoint is not valid", func(t *testing.T) {
		checkpoint := createEnygmaCheckpoint(
			"id1",
			"resource1",
			big.NewInt(100),
			big.NewInt(200),
			big.NewInt(100),
			big.NewInt(101),
		)
		enygma := createEnygma("resource1", big.NewInt(50), big.NewInt(10), big.NewInt(99), big.NewInt(100))
		history := createEnygmaHistory("resource1", big.NewInt(50), big.NewInt(90))

		checkpointRepo := &MockEnygmaCheckpointRepository{
			MarkAsFinalizedFunc: func(ctx context.Context, resourceId string, finalizedBlockNumber *big.Int) error {
				assert.Fail(t, "MarkAsFinalized should not be called")
				return nil
			},
			GetValidationCandidatesFunc: func(ctx context.Context) ([]types.EnygmaCheckpoint, error) {
				return []types.EnygmaCheckpoint{*checkpoint}, nil
			},
		}
		historyRepo := &MockEnygmaHistoryRepository{
			GetEnygmaHistoryForCheckpointsFunc: func(ctx context.Context, resourceIds []string, blockNumbers []*big.Int) ([]types.EnygmaHistory, error) {
				return []types.EnygmaHistory{*history}, nil
			},
		}
		enygmaRepo := &MockEnygmaRepository{
			UpdateEnygmaFunc: func(ctx context.Context, resourceId string, finalizedBalance, finalizedR, finalizedBlockNumber, pendingBlockNumber *big.Int) error {
				assert.Fail(t, "UpdateEnygma should not be called")
				return nil
			},
			GetEnygmaByResourceIdsFunc: func(ctx context.Context, resourceIds []string) ([]types.Enygma, error) {
				return []types.Enygma{*enygma}, nil
			},
		}
		txManager := &MockTransactionManager{}
		resyncService := &MockResyncService{}

		svc := service.NewEnygmaSyncService(
			syncConfig,
			txManager,
			enygmaRepo,
			historyRepo,
			checkpointRepo,
			resyncService,
		)

		err := svc.Run(ctx)

		require.NoError(t, err)
		assert.Len(t, txManager.WithTransactionCalls.Calls, 0)
		assert.Len(t, checkpointRepo.MarkAsFinalizedCalls.Calls, 0)
		assert.Len(t, enygmaRepo.UpdateEnygmaCalls.Calls, 0)
	})

	t.Run("successfully finalizes checkpoint with valid cryptographic values", func(t *testing.T) {
		// This test demonstrates successful checkpoint finalization using real Pedersen commitment values
		checkpoint, prevEnygma, history, newBalance, newR := createValidCheckpointWithCryptography()

		checkpointRepo := &MockEnygmaCheckpointRepository{
			GetValidationCandidatesFunc: func(ctx context.Context) ([]types.EnygmaCheckpoint, error) {
				return []types.EnygmaCheckpoint{*checkpoint}, nil
			},
			MarkAsFinalizedFunc: func(ctx context.Context, resourceId string, finalizedBlockNumber *big.Int) error {
				return nil
			},
		}
		historyRepo := &MockEnygmaHistoryRepository{
			GetEnygmaHistoryForCheckpointsFunc: func(ctx context.Context, resourceIds []string, blockNumbers []*big.Int) ([]types.EnygmaHistory, error) {
				return []types.EnygmaHistory{*history}, nil
			},
		}
		enygmaRepo := &MockEnygmaRepository{
			GetEnygmaByResourceIdsFunc: func(ctx context.Context, resourceIds []string) ([]types.Enygma, error) {
				return []types.Enygma{*prevEnygma}, nil
			},
			UpdateEnygmaFunc: func(ctx context.Context, resourceId string, finalizedBalance, finalizedR, finalizedBlockNumber, pendingBlockNumber *big.Int) error {
				// Verify the correct values are passed
				assert.Equal(t, "resource1", resourceId)
				assert.Equal(t, newBalance, finalizedBalance, "finalized balance should match computed value")
				assert.Equal(t, newR, finalizedR, "finalized R should match computed value")
				assert.Equal(t, checkpoint.FinalizedBlockNumber, finalizedBlockNumber)
				assert.Equal(t, checkpoint.PendingBlockNumber, pendingBlockNumber)
				return nil
			},
		}
		txManager := &MockTransactionManager{
			WithTransactionFunc: func(ctx context.Context, fn func(ctx context.Context) error) error {
				return fn(ctx)
			},
		}
		resyncService := &MockResyncService{}

		svc := service.NewEnygmaSyncService(
			syncConfig,
			txManager,
			enygmaRepo,
			historyRepo,
			checkpointRepo,
			resyncService,
		)

		err := svc.Run(ctx)

		// Should succeed with valid cryptographic values
		require.NoError(t, err)
		assert.Len(
			t,
			txManager.WithTransactionCalls.Calls,
			1,
			"transaction should be called once for successful finalization",
		)
		assert.Len(t, checkpointRepo.MarkAsFinalizedCalls.Calls, 1, "checkpoint should be marked as finalized")
		assert.Len(t, enygmaRepo.UpdateEnygmaCalls.Calls, 1, "enygma state should be updated")
	})

	t.Run("finalizeCheckpoint returns error when MarkAsFinalized fails", func(t *testing.T) {
		// This test verifies that MarkAsFinalized errors are propagated and transaction is rolled back
		checkpoint, prevEnygma, history, _, _ := createValidCheckpointWithCryptography()

		checkpointRepo := &MockEnygmaCheckpointRepository{
			GetValidationCandidatesFunc: func(ctx context.Context) ([]types.EnygmaCheckpoint, error) {
				return []types.EnygmaCheckpoint{*checkpoint}, nil
			},
			MarkAsFinalizedFunc: func(ctx context.Context, resourceId string, finalizedBlockNumber *big.Int) error {
				return fmt.Errorf("failed to mark checkpoint as finalized")
			},
		}
		historyRepo := &MockEnygmaHistoryRepository{
			GetEnygmaHistoryForCheckpointsFunc: func(ctx context.Context, resourceIds []string, blockNumbers []*big.Int) ([]types.EnygmaHistory, error) {
				return []types.EnygmaHistory{*history}, nil
			},
		}
		enygmaRepo := &MockEnygmaRepository{
			GetEnygmaByResourceIdsFunc: func(ctx context.Context, resourceIds []string) ([]types.Enygma, error) {
				return []types.Enygma{*prevEnygma}, nil
			},
			UpdateEnygmaFunc: func(ctx context.Context, resourceId string, finalizedBalance, finalizedR, finalizedBlockNumber, pendingBlockNumber *big.Int) error {
				// Should NOT be called since MarkAsFinalized fails first
				t.Errorf("UpdateEnygma should not be called when MarkAsFinalized fails")
				return nil
			},
		}
		txManager := &MockTransactionManager{
			WithTransactionFunc: func(ctx context.Context, fn func(ctx context.Context) error) error {
				return fn(ctx)
			},
		}
		resyncService := &MockResyncService{}

		svc := service.NewEnygmaSyncService(
			syncConfig,
			txManager,
			enygmaRepo,
			historyRepo,
			checkpointRepo,
			resyncService,
		)

		err := svc.Run(ctx)

		// Should propagate the MarkAsFinalized error
		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to mark checkpoint as finalized")
		// UpdateEnygma should not be called due to transaction failure
		assert.Len(
			t,
			enygmaRepo.UpdateEnygmaCalls.Calls,
			0,
			"UpdateEnygma should not be called when MarkAsFinalized fails",
		)
	})

	t.Run("finalizeCheckpoint returns error when UpdateEnygma fails", func(t *testing.T) {
		// This test verifies that UpdateEnygma errors are propagated and MarkAsFinalized happens first
		checkpoint, prevEnygma, history, _, _ := createValidCheckpointWithCryptography()

		checkpointRepo := &MockEnygmaCheckpointRepository{
			GetValidationCandidatesFunc: func(ctx context.Context) ([]types.EnygmaCheckpoint, error) {
				return []types.EnygmaCheckpoint{*checkpoint}, nil
			},
			MarkAsFinalizedFunc: func(ctx context.Context, resourceId string, finalizedBlockNumber *big.Int) error {
				return nil
			},
		}
		historyRepo := &MockEnygmaHistoryRepository{
			GetEnygmaHistoryForCheckpointsFunc: func(ctx context.Context, resourceIds []string, blockNumbers []*big.Int) ([]types.EnygmaHistory, error) {
				return []types.EnygmaHistory{*history}, nil
			},
		}
		enygmaRepo := &MockEnygmaRepository{
			GetEnygmaByResourceIdsFunc: func(ctx context.Context, resourceIds []string) ([]types.Enygma, error) {
				return []types.Enygma{*prevEnygma}, nil
			},
			UpdateEnygmaFunc: func(ctx context.Context, resourceId string, finalizedBalance, finalizedR, finalizedBlockNumber, pendingBlockNumber *big.Int) error {
				return fmt.Errorf("failed to update enygma state")
			},
		}
		txManager := &MockTransactionManager{
			WithTransactionFunc: func(ctx context.Context, fn func(ctx context.Context) error) error {
				return fn(ctx)
			},
		}
		resyncService := &MockResyncService{}

		svc := service.NewEnygmaSyncService(
			syncConfig,
			txManager,
			enygmaRepo,
			historyRepo,
			checkpointRepo,
			resyncService,
		)

		err := svc.Run(ctx)

		// Should propagate the UpdateEnygma error
		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to update enygma state")
		// MarkAsFinalized should have been called despite UpdateEnygma failure
		// (it's part of the transaction, so in real execution it would be rolled back)
		assert.Len(
			t,
			checkpointRepo.MarkAsFinalizedCalls.Calls,
			1,
			"MarkAsFinalized should be called before UpdateEnygma",
		)
	})

	t.Run("successfully finalizes checkpoint when ResyncEnygma is called", func(t *testing.T) {
		// This test verifies that after 10 failed validations, ResyncEnygma is called
		checkpoint, prevEnygma, history := createCheckpointWithValidationFailure()

		checkpointRepo := &MockEnygmaCheckpointRepository{
			GetValidationCandidatesFunc: func(ctx context.Context) ([]types.EnygmaCheckpoint, error) {
				return []types.EnygmaCheckpoint{*checkpoint}, nil
			},
		}
		historyRepo := &MockEnygmaHistoryRepository{
			GetEnygmaHistoryForCheckpointsFunc: func(ctx context.Context, resourceIds []string, blockNumbers []*big.Int) ([]types.EnygmaHistory, error) {
				return []types.EnygmaHistory{*history}, nil
			},
		}
		enygmaRepo := &MockEnygmaRepository{
			GetEnygmaByResourceIdsFunc: func(ctx context.Context, resourceIds []string) ([]types.Enygma, error) {
				return []types.Enygma{*prevEnygma}, nil
			},
		}
		txManager := &MockTransactionManager{
			WithTransactionFunc: func(ctx context.Context, fn func(ctx context.Context) error) error {
				return fn(ctx)
			},
		}
		resyncService := &MockResyncService{
			ResyncEnygmaFunc: func(ctx context.Context, resourceId string) error {
				return nil // Success
			},
		}

		svc := service.NewEnygmaSyncService(
			syncConfig,
			txManager,
			enygmaRepo,
			historyRepo,
			checkpointRepo,
			resyncService,
		)

		// Call Run 3 times (i=0 to 2) to trigger maxRetries (maxRetries=2: attempts 0,1 increment, then 2 resyncs)
		for i := 0; i <= syncConfig.MaxRetries; i++ {
			err := svc.Run(ctx)
			require.NoError(t, err, "Run should not return error on attempt %d", i+1)
		}

		// Verify ResyncEnygma was called exactly once on the 3rd attempt
		resyncService.ResyncEnygmaCalls.RLock()
		calls := len(resyncService.ResyncEnygmaCalls.Calls)
		resourceIds := make([]string, 0, calls)
		for _, call := range resyncService.ResyncEnygmaCalls.Calls {
			resourceIds = append(resourceIds, call.ResourceId)
		}
		resyncService.ResyncEnygmaCalls.RUnlock()

		assert.Equal(t, 1, calls, "ResyncEnygma should be called exactly once")
		assert.Equal(t, []string{"resource1"}, resourceIds, "ResyncEnygma should be called with correct resourceId")
	})

	t.Run("returns error when ResyncEnygma fails after maxRetries", func(t *testing.T) {
		// This test verifies that ResyncEnygma errors are propagated
		checkpoint, prevEnygma, history := createCheckpointWithValidationFailure()

		checkpointRepo := &MockEnygmaCheckpointRepository{
			GetValidationCandidatesFunc: func(ctx context.Context) ([]types.EnygmaCheckpoint, error) {
				return []types.EnygmaCheckpoint{*checkpoint}, nil
			},
		}
		historyRepo := &MockEnygmaHistoryRepository{
			GetEnygmaHistoryForCheckpointsFunc: func(ctx context.Context, resourceIds []string, blockNumbers []*big.Int) ([]types.EnygmaHistory, error) {
				return []types.EnygmaHistory{*history}, nil
			},
		}
		enygmaRepo := &MockEnygmaRepository{
			GetEnygmaByResourceIdsFunc: func(ctx context.Context, resourceIds []string) ([]types.Enygma, error) {
				return []types.Enygma{*prevEnygma}, nil
			},
		}
		txManager := &MockTransactionManager{
			WithTransactionFunc: func(ctx context.Context, fn func(ctx context.Context) error) error {
				return fn(ctx)
			},
		}
		resyncService := &MockResyncService{
			ResyncEnygmaFunc: func(ctx context.Context, resourceId string) error {
				return fmt.Errorf("resync failed: network error")
			},
		}

		svc := service.NewEnygmaSyncService(
			syncConfig,
			txManager,
			enygmaRepo,
			historyRepo,
			checkpointRepo,
			resyncService,
		)

		// Call Run 2 times (i=0 to 1) to increment retry counter (no resync yet)
		for i := 0; i < syncConfig.MaxRetries; i++ {
			err := svc.Run(ctx)
			require.NoError(t, err)
		}

		// 3rd call should trigger ResyncEnygma and propagate its error
		err := svc.Run(ctx)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "resync failed: network error")

		// Verify ResyncEnygma was called
		resyncService.ResyncEnygmaCalls.RLock()
		assert.Equal(t, 1, len(resyncService.ResyncEnygmaCalls.Calls), "ResyncEnygma should be called once")
		resyncService.ResyncEnygmaCalls.RUnlock()
	})
}
