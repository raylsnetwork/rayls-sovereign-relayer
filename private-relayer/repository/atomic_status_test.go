// Decommissioning Teleport (vanilla, atomic).

//go:build integration

package repository_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/raylsnetwork/rayls-privacy-relayer-api/private-relayer/repository"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/private-relayer/repository/testdata"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/testtools"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/types"
)

func TestAtomicStatusRepository_BatchCreate(t *testing.T) {
	testtools.SilenceLogger()

	pool, teardown := testtools.SetupPostgres(t, dbConf)
	defer teardown()

	repo := repository.NewAtomicStatusRepository(pool)

	t.Run("creates a single status update message", func(t *testing.T) {
		err := repo.BatchCreate(context.TODO(), []types.AtomicStatusUpdateMessage{testdata.AtomicSUM1})
		require.Nil(t, err)

		got, err := repo.GetUnprocessedBySharedIDs(context.TODO(), []string{testdata.AtomicSUM1.SharedID})
		require.Nil(t, err)

		require.Equal(t, []types.AtomicStatusUpdateMessage{testdata.AtomicSUM1}, got)
	})

	t.Run("creates multiple status update messages in a single transaction", func(t *testing.T) {
		err := repo.BatchCreate(
			context.TODO(),
			[]types.AtomicStatusUpdateMessage{testdata.AtomicSUM2, testdata.AtomicSUM3},
		)
		require.Nil(t, err)

		got, err := repo.GetUnprocessedBySharedIDs(
			context.TODO(),
			[]string{testdata.AtomicSUM2.SharedID, testdata.AtomicSUM3.SharedID},
		)
		require.Nil(t, err)

		require.Len(t, got, 2)
		require.ElementsMatch(t, []types.AtomicStatusUpdateMessage{testdata.AtomicSUM2, testdata.AtomicSUM3}, got)
	})
}

func TestAtomicStatusRepository_GetUnprocessedBySharedIDs(t *testing.T) {
	testtools.SilenceLogger()

	pool, teardown := testtools.SetupPostgres(t, dbConf)
	defer teardown()

	repo := repository.NewAtomicStatusRepository(pool)

	// Seed using repository BatchCreate method
	err := repo.BatchCreate(context.TODO(), []types.AtomicStatusUpdateMessage{testdata.AtomicSUM1, testdata.AtomicSUM2})
	require.Nil(t, err)
	// Mark AtomicSUM3 as processed by using MarkAsProcessed
	err = repo.BatchCreate(context.TODO(), []types.AtomicStatusUpdateMessage{testdata.AtomicSUM3})
	require.Nil(t, err)
	err = repo.MarkAsProcessed(context.TODO(), []string{testdata.AtomicSUM3.SharedID})
	require.Nil(t, err)

	t.Run("returns empty slice when no shared IDs match", func(t *testing.T) {
		got, err := repo.GetUnprocessedBySharedIDs(context.TODO(), []string{"non-existent-id"})
		require.Nil(t, err)

		require.Empty(t, got)
	})

	t.Run("returns only unprocessed items for matching shared IDs", func(t *testing.T) {
		got, err := repo.GetUnprocessedBySharedIDs(
			context.TODO(),
			[]string{testdata.AtomicSUM1.SharedID, testdata.AtomicSUM3.SharedID},
		)
		require.Nil(t, err)

		require.Equal(t, []types.AtomicStatusUpdateMessage{testdata.AtomicSUM1}, got)
	})

	t.Run("returns all unprocessed items for matching shared IDs", func(t *testing.T) {
		got, err := repo.GetUnprocessedBySharedIDs(
			context.TODO(),
			[]string{testdata.AtomicSUM1.SharedID, testdata.AtomicSUM2.SharedID},
		)
		require.Nil(t, err)

		require.Len(t, got, 2)
		require.ElementsMatch(t, []types.AtomicStatusUpdateMessage{testdata.AtomicSUM1, testdata.AtomicSUM2}, got)
	})

	t.Run("respects WithLimit option", func(t *testing.T) {
		got, err := repo.GetUnprocessedBySharedIDs(
			context.TODO(),
			[]string{testdata.AtomicSUM1.SharedID, testdata.AtomicSUM2.SharedID},
			repository.WithLimit(1),
		)
		require.Nil(t, err)

		require.Len(t, got, 1)
	})
}
