//go:build integration

package repository_test

import (
	"context"
	"math/big"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/raylsnetwork/rayls-sovereign-relayer/private-relayer/repository"
	"github.com/raylsnetwork/rayls-sovereign-relayer/private-relayer/repository/testdata"
	"github.com/raylsnetwork/rayls-sovereign-relayer/testtools"
	"github.com/raylsnetwork/rayls-sovereign-relayer/types"
)

func TestEnygmaRepository_CreateEnygma(t *testing.T) {
	testtools.SilenceLogger()

	pool, teardown := testtools.SetupPostgres(t, dbConf)
	defer teardown()

	repo := repository.NewEnygmaRepository(pool)

	t.Run("persists enygma record in database", func(t *testing.T) {
		err := repo.CreateEnygma(context.TODO(), testdata.Enygma1)
		require.Nil(t, err)

		got, err := repo.GetEnygmaByResourceId(context.TODO(), testdata.Enygma1.ResourceId)
		require.Nil(t, err)

		require.Equal(t, testdata.Enygma1, got)
	})

	t.Run("returns error on duplicate resource ID", func(t *testing.T) {
		err := repo.CreateEnygma(context.TODO(), testdata.Enygma1)
		require.NotNil(t, err)
	})
}

func TestEnygmaRepository_UpdateEnygma(t *testing.T) {
	testtools.SilenceLogger()

	pool, teardown := testtools.SetupPostgres(t, dbConf)
	defer teardown()

	repo := repository.NewEnygmaRepository(pool)

	// Seed using repository CreateEnygma method
	err := repo.CreateEnygma(context.TODO(), testdata.Enygma1)
	require.Nil(t, err)
	err = repo.CreateEnygma(context.TODO(), testdata.Enygma2)
	require.Nil(t, err)

	t.Run("updates finalized balance and block numbers for existing resource", func(t *testing.T) {
		newBalance := big.NewInt(7777)
		newR := big.NewInt(8888)
		newFinalizedBlock := big.NewInt(5000)
		newPendingBlock := big.NewInt(6000)

		err := repo.UpdateEnygma(
			context.TODO(),
			testdata.Enygma1.ResourceId,
			newBalance,
			newR,
			newFinalizedBlock,
			newPendingBlock,
		)
		require.Nil(t, err)

		got, err := repo.GetEnygmaByResourceId(context.TODO(), testdata.Enygma1.ResourceId)
		require.Nil(t, err)

		require.Equal(t, newBalance, got.FinalizedBalance)
		require.Equal(t, newR, got.FinalizedR)
		require.Equal(t, newFinalizedBlock, got.FinalizedBlockNumber)
		require.Equal(t, newPendingBlock, got.PendingBlockNumber)
	})

	t.Run("returns no error when resource does not exist", func(t *testing.T) {
		err := repo.UpdateEnygma(
			context.TODO(),
			"non-existent",
			big.NewInt(1),
			big.NewInt(1),
			big.NewInt(1),
			big.NewInt(1),
		)
		require.Nil(t, err)
	})
}

func TestEnygmaRepository_GetEnygmaByResourceId(t *testing.T) {
	testtools.SilenceLogger()

	pool, teardown := testtools.SetupPostgres(t, dbConf)
	defer teardown()

	repo := repository.NewEnygmaRepository(pool)

	t.Run("returns error when resource does not exist", func(t *testing.T) {
		_, err := repo.GetEnygmaByResourceId(context.TODO(), "non-existent")
		require.NotNil(t, err)
		require.Contains(t, err.Error(), "no enygma found")
	})

	// Seed using repository CreateEnygma method
	err := repo.CreateEnygma(context.TODO(), testdata.Enygma1)
	require.Nil(t, err)
	err = repo.CreateEnygma(context.TODO(), testdata.Enygma2)
	require.Nil(t, err)

	t.Run("returns enygma for existing resource ID", func(t *testing.T) {
		got, err := repo.GetEnygmaByResourceId(context.TODO(), testdata.Enygma1.ResourceId)
		require.Nil(t, err)

		require.Equal(t, testdata.Enygma1, got)
	})
}

func TestEnygmaRepository_GetEnygmaByResourceIds(t *testing.T) {
	testtools.SilenceLogger()

	pool, teardown := testtools.SetupPostgres(t, dbConf)
	defer teardown()

	repo := repository.NewEnygmaRepository(pool)

	t.Run("returns empty slice for empty input", func(t *testing.T) {
		got, err := repo.GetEnygmaByResourceIds(context.TODO(), []string{})
		require.Nil(t, err)

		require.Equal(t, []types.Enygma{}, got)
	})

	t.Run("returns empty slice when no records in database", func(t *testing.T) {
		got, err := repo.GetEnygmaByResourceIds(context.TODO(), []string{"some-id"})
		require.Nil(t, err)

		var want []types.Enygma
		require.Equal(t, want, got)
	})

	// Seed using repository CreateEnygma method
	err := repo.CreateEnygma(context.TODO(), testdata.Enygma1)
	require.Nil(t, err)
	err = repo.CreateEnygma(context.TODO(), testdata.Enygma2)
	require.Nil(t, err)

	t.Run("returns empty slice when no resource IDs match", func(t *testing.T) {
		got, err := repo.GetEnygmaByResourceIds(context.TODO(), []string{"non-existent"})
		require.Nil(t, err)

		var want []types.Enygma
		require.Equal(t, want, got)
	})

	t.Run("returns matching enygma records for given resource IDs", func(t *testing.T) {
		got, err := repo.GetEnygmaByResourceIds(context.TODO(), []string{testdata.Enygma1.ResourceId})
		require.Nil(t, err)

		require.Equal(t, []types.Enygma{testdata.Enygma1}, got)
	})

	t.Run("returns multiple matching records", func(t *testing.T) {
		got, err := repo.GetEnygmaByResourceIds(
			context.TODO(),
			[]string{testdata.Enygma1.ResourceId, testdata.Enygma2.ResourceId},
		)
		require.Nil(t, err)

		require.Len(t, got, 2)
		require.ElementsMatch(t, []types.Enygma{testdata.Enygma1, testdata.Enygma2}, got)
	})
}
