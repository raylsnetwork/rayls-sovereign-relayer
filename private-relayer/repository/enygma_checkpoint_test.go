//go:build integration

package repository_test

import (
	"context"
	"math/big"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/raylsnetwork/rayls-privacy-relayer-api/private-relayer/repository"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/private-relayer/repository/testdata"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/testtools"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/types"
)

func TestEnygmaCheckpointRepository_CreateEnygmaCheckpoint(t *testing.T) {
	testtools.SilenceLogger()

	pool, teardown := testtools.SetupPostgres(t, dbConf)
	defer teardown()

	repo := repository.NewEnygmaCheckpointRepository(pool)

	t.Run("persists checkpoint in database", func(t *testing.T) {
		err := repo.CreateEnygmaCheckpoint(context.TODO(), testdata.EnygmaCheckpoint1)
		require.Nil(t, err)

		got, err := repo.GetLatestCheckpointByFilters(
			context.TODO(),
			testdata.EnygmaCheckpoint1.ResourceId,
			nil,
			nil,
			nil,
		)
		require.Nil(t, err)
		require.NotNil(t, got)
		require.NotEmpty(t, got.ID)
		require.Equal(t, testdata.EnygmaCheckpoint1.ResourceId, got.ResourceId)
		require.Equal(t, testdata.EnygmaCheckpoint1.FinalizedPublicBalanceX, got.FinalizedPublicBalanceX)
		require.Equal(t, testdata.EnygmaCheckpoint1.FinalizedPublicBalanceY, got.FinalizedPublicBalanceY)
		require.Equal(t, testdata.EnygmaCheckpoint1.FinalizedBlockNumber, got.FinalizedBlockNumber)
		require.Equal(t, testdata.EnygmaCheckpoint1.PendingBlockNumber, got.PendingBlockNumber)
		require.Equal(t, testdata.EnygmaCheckpoint1.Status, got.Status)
	})

	t.Run("returns error on duplicate resource_id and finalized_block_number", func(t *testing.T) {
		err := repo.CreateEnygmaCheckpoint(context.TODO(), testdata.EnygmaCheckpoint1)
		require.NotNil(t, err)
	})
}

func TestEnygmaCheckpointRepository_GetLatestCheckpointByFilters(t *testing.T) {
	testtools.SilenceLogger()

	pool, teardown := testtools.SetupPostgres(t, dbConf)
	defer teardown()

	repo := repository.NewEnygmaCheckpointRepository(pool)

	// Seed using repository CreateEnygmaCheckpoint method
	err := repo.CreateEnygmaCheckpoint(context.TODO(), testdata.EnygmaCheckpoint1)
	require.Nil(t, err)
	err = repo.CreateEnygmaCheckpoint(context.TODO(), testdata.EnygmaCheckpoint1HigherBlock)
	require.Nil(t, err)
	err = repo.CreateEnygmaCheckpoint(context.TODO(), testdata.EnygmaCheckpoint2)
	require.Nil(t, err)
	err = repo.CreateEnygmaCheckpoint(context.TODO(), testdata.EnygmaCheckpoint3Finalized)
	require.Nil(t, err)

	t.Run("returns nil when no checkpoints exist for resource ID", func(t *testing.T) {
		got, err := repo.GetLatestCheckpointByFilters(context.TODO(), "non-existent", nil, nil, nil)
		require.Nil(t, err)
		require.Nil(t, got)
	})

	t.Run("returns latest checkpoint by finalized_block_number for resource ID", func(t *testing.T) {
		got, err := repo.GetLatestCheckpointByFilters(context.TODO(), "resource-checkpoint-1", nil, nil, nil)
		require.Nil(t, err)
		require.NotNil(t, got)
		require.Equal(t, testdata.EnygmaCheckpoint1HigherBlock.FinalizedBlockNumber, got.FinalizedBlockNumber)
		require.Equal(t, testdata.EnygmaCheckpoint1HigherBlock.FinalizedPublicBalanceX, got.FinalizedPublicBalanceX)
	})

	t.Run("filters by status", func(t *testing.T) {
		status := types.EnygmaCheckpointStatusFinal
		got, err := repo.GetLatestCheckpointByFilters(context.TODO(), "", &status, nil, nil)
		require.Nil(t, err)
		require.NotNil(t, got)
		require.Equal(t, testdata.EnygmaCheckpoint3Finalized.ResourceId, got.ResourceId)
	})

	t.Run("filters by finalized block number", func(t *testing.T) {
		got, err := repo.GetLatestCheckpointByFilters(
			context.TODO(),
			"resource-checkpoint-1",
			nil,
			big.NewInt(500),
			nil,
		)
		require.Nil(t, err)
		require.NotNil(t, got)
		require.Equal(t, testdata.EnygmaCheckpoint1.FinalizedBlockNumber, got.FinalizedBlockNumber)
	})

	t.Run("returns nil when no checkpoints match filters", func(t *testing.T) {
		got, err := repo.GetLatestCheckpointByFilters(
			context.TODO(),
			"resource-checkpoint-1",
			nil,
			big.NewInt(9999),
			nil,
		)
		require.Nil(t, err)
		require.Nil(t, got)
	})
}

func TestEnygmaCheckpointRepository_GetValidationCandidates(t *testing.T) {
	testtools.SilenceLogger()

	pool, teardown := testtools.SetupPostgres(t, dbConf)
	defer teardown()

	repo := repository.NewEnygmaCheckpointRepository(pool)

	t.Run("returns empty slice when no tentative checkpoints exist", func(t *testing.T) {
		got, err := repo.GetValidationCandidates(context.TODO())
		require.Nil(t, err)
		require.Empty(t, got)
	})

	// Seed using repository CreateEnygmaCheckpoint method
	err := repo.CreateEnygmaCheckpoint(context.TODO(), testdata.EnygmaCheckpoint1)
	require.Nil(t, err)
	err = repo.CreateEnygmaCheckpoint(context.TODO(), testdata.EnygmaCheckpoint1HigherBlock)
	require.Nil(t, err)
	err = repo.CreateEnygmaCheckpoint(context.TODO(), testdata.EnygmaCheckpoint2)
	require.Nil(t, err)
	err = repo.CreateEnygmaCheckpoint(context.TODO(), testdata.EnygmaCheckpoint3Finalized)
	require.Nil(t, err)

	t.Run("returns earliest tentative checkpoint per resource ID", func(t *testing.T) {
		got, err := repo.GetValidationCandidates(context.TODO())
		require.Nil(t, err)
		require.Len(t, got, 2)

		resourceBlockMap := make(map[string]*big.Int)
		for _, cp := range got {
			resourceBlockMap[cp.ResourceId] = cp.FinalizedBlockNumber
		}

		require.Equal(t, big.NewInt(500), resourceBlockMap["resource-checkpoint-1"])
		require.Equal(t, big.NewInt(700), resourceBlockMap["resource-checkpoint-2"])
	})
}

func TestEnygmaCheckpointRepository_MarkAsFinalized(t *testing.T) {
	testtools.SilenceLogger()

	pool, teardown := testtools.SetupPostgres(t, dbConf)
	defer teardown()

	repo := repository.NewEnygmaCheckpointRepository(pool)

	// Seed using repository CreateEnygmaCheckpoint method
	err := repo.CreateEnygmaCheckpoint(context.TODO(), testdata.EnygmaCheckpoint1)
	require.Nil(t, err)

	t.Run("marks checkpoint as finalized", func(t *testing.T) {
		err := repo.MarkAsFinalized(
			context.TODO(),
			testdata.EnygmaCheckpoint1.ResourceId,
			testdata.EnygmaCheckpoint1.FinalizedBlockNumber,
		)
		require.Nil(t, err)

		status := types.EnygmaCheckpointStatusFinal
		got, err := repo.GetLatestCheckpointByFilters(
			context.TODO(),
			testdata.EnygmaCheckpoint1.ResourceId,
			&status,
			nil,
			nil,
		)
		require.Nil(t, err)
		require.NotNil(t, got)
		require.Equal(t, types.EnygmaCheckpointStatusFinal, got.Status)
	})

	t.Run("returns error when no checkpoint matches", func(t *testing.T) {
		err := repo.MarkAsFinalized(context.TODO(), "non-existent", big.NewInt(9999))
		require.NotNil(t, err)
		require.Contains(t, err.Error(), "no checkpoint found")
	})
}
