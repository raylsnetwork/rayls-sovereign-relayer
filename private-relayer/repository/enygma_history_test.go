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

func strPtr(s string) *string {
	return &s
}

func TestEnygmaHistoryRepository_InsertEnygmaHistory(t *testing.T) {
	testtools.SilenceLogger()

	pool, teardown := testtools.SetupPostgres(t, dbConf)
	defer teardown()

	repo := repository.NewEnygmaHistoryRepository(pool)

	t.Run("persists history record in database", func(t *testing.T) {
		err := repo.InsertEnygmaHistory(context.TODO(), testdata.EnygmaHistory1)
		require.Nil(t, err)

		got, err := repo.GetEnygmaHistoryByFilters(
			context.TODO(),
			strPtr(testdata.EnygmaHistory1.ResourceId),
			testdata.EnygmaHistory1.BlockNumberPrivateHub,
			testdata.EnygmaHistory1.FromChainId,
			testdata.EnygmaHistory1.EventType,
		)
		require.Nil(t, err)

		require.Len(t, got, 1)
		require.Equal(t, testdata.EnygmaHistory1, got[0])
	})

	t.Run("returns error on duplicate unique key", func(t *testing.T) {
		err := repo.InsertEnygmaHistory(context.TODO(), testdata.EnygmaHistory1)
		require.NotNil(t, err)
	})
}

func TestEnygmaHistoryRepository_GetEnygmaHistoryByFilters(t *testing.T) {
	testtools.SilenceLogger()

	pool, teardown := testtools.SetupPostgres(t, dbConf)
	defer teardown()

	repo := repository.NewEnygmaHistoryRepository(pool)

	// Seed using repository InsertEnygmaHistory method
	err := repo.InsertEnygmaHistory(context.TODO(), testdata.EnygmaHistory1)
	require.Nil(t, err)
	err = repo.InsertEnygmaHistory(context.TODO(), testdata.EnygmaHistory2)
	require.Nil(t, err)
	err = repo.InsertEnygmaHistory(context.TODO(), testdata.EnygmaHistory3SameResource)
	require.Nil(t, err)

	t.Run("returns empty slice when no records match", func(t *testing.T) {
		got, err := repo.GetEnygmaHistoryByFilters(context.TODO(), strPtr("non-existent"), nil, nil, types.EnygmaMint)
		require.Nil(t, err)

		require.Empty(t, got)
	})

	t.Run("returns records matching resource ID and event type", func(t *testing.T) {
		got, err := repo.GetEnygmaHistoryByFilters(
			context.TODO(),
			strPtr("resource-history-1"),
			nil,
			nil,
			types.EnygmaMint,
		)
		require.Nil(t, err)

		require.Len(t, got, 2)
		require.ElementsMatch(
			t,
			[]types.EnygmaHistory{testdata.EnygmaHistory1, testdata.EnygmaHistory3SameResource},
			got,
		)
	})

	t.Run("returns record matching block number", func(t *testing.T) {
		got, err := repo.GetEnygmaHistoryByFilters(
			context.TODO(),
			nil,
			testdata.EnygmaHistory1.BlockNumberPrivateHub,
			nil,
			types.EnygmaMint,
		)
		require.Nil(t, err)

		require.Equal(t, []types.EnygmaHistory{testdata.EnygmaHistory1}, got)
	})

	t.Run("filters by fromChainId", func(t *testing.T) {
		got, err := repo.GetEnygmaHistoryByFilters(context.TODO(), nil, nil, big.NewInt(20), types.EnygmaBurn)
		require.Nil(t, err)

		require.Equal(t, []types.EnygmaHistory{testdata.EnygmaHistory2}, got)
	})

	t.Run("always filters by event type", func(t *testing.T) {
		got, err := repo.GetEnygmaHistoryByFilters(context.TODO(), nil, nil, nil, types.EnygmaBurn)
		require.Nil(t, err)

		require.Equal(t, []types.EnygmaHistory{testdata.EnygmaHistory2}, got)
	})
}

func TestEnygmaHistoryRepository_GetEnygmaHistoryByUniqueKey(t *testing.T) {
	testtools.SilenceLogger()

	client, teardown := testtools.SetupMongo(t, dbConf)
	defer teardown()

	database := client.Database(dbConf.Database)

	repo, err := repository.NewEnygmaHistoryRepository(database, client)
	require.NoError(t, err)

	seed := []interface{}{
		testdata.ModelEnygmaHistory1,
		testdata.ModelEnygmaHistory2,
	}
	err = testtools.SeedCollection(database, repository.EnygmaHistoryCollectionName, seed)
	require.Nil(t, err)

	t.Run("returns nil when no record matches", func(t *testing.T) {
		got, err := repo.GetEnygmaHistoryByUniqueKey(
			context.TODO(),
			"non-existent",
			big.NewInt(9999),
			big.NewInt(1),
			types.EnygmaMint,
		)
		require.Nil(t, err)
		require.Nil(t, got)
	})

	t.Run("returns single record matching unique key", func(t *testing.T) {
		got, err := repo.GetEnygmaHistoryByUniqueKey(
			context.TODO(),
			testdata.EnygmaHistory1.ResourceId,
			testdata.EnygmaHistory1.BlockNumberPrivateHub,
			testdata.EnygmaHistory1.FromChainId,
			testdata.EnygmaHistory1.EventType,
		)
		require.Nil(t, err)
		require.NotNil(t, got)
		require.Equal(t, testdata.EnygmaHistory1, *got)
	})
}

func TestEnygmaHistoryRepository_GetEnygmaHistoryForCheckpoints(t *testing.T) {
	testtools.SilenceLogger()

	pool, teardown := testtools.SetupPostgres(t, dbConf)
	defer teardown()

	repo := repository.NewEnygmaHistoryRepository(pool)

	// Seed using repository InsertEnygmaHistory method
	err := repo.InsertEnygmaHistory(context.TODO(), testdata.EnygmaHistory1)
	require.Nil(t, err)
	err = repo.InsertEnygmaHistory(context.TODO(), testdata.EnygmaHistory2)
	require.Nil(t, err)
	err = repo.InsertEnygmaHistory(context.TODO(), testdata.EnygmaHistory3SameResource)
	require.Nil(t, err)

	t.Run("returns empty slice for empty input", func(t *testing.T) {
		got, err := repo.GetEnygmaHistoryForCheckpoints(context.TODO(), []string{}, []*big.Int{})
		require.Nil(t, err)

		require.Equal(t, []types.EnygmaHistory{}, got)
	})

	t.Run("returns error when slices have different lengths", func(t *testing.T) {
		_, err := repo.GetEnygmaHistoryForCheckpoints(context.TODO(), []string{"id1"}, []*big.Int{})
		require.NotNil(t, err)
	})

	t.Run("returns empty slice when no records match", func(t *testing.T) {
		got, err := repo.GetEnygmaHistoryForCheckpoints(
			context.TODO(),
			[]string{"non-existent"},
			[]*big.Int{big.NewInt(9999)},
		)
		require.Nil(t, err)

		var want []types.EnygmaHistory
		require.Equal(t, want, got)
	})

	t.Run("returns matching history records for given resource and block pairs", func(t *testing.T) {
		got, err := repo.GetEnygmaHistoryForCheckpoints(
			context.TODO(),
			[]string{"resource-history-1", "resource-history-2"},
			[]*big.Int{big.NewInt(1001), big.NewInt(2002)},
		)
		require.Nil(t, err)

		require.Len(t, got, 2)
		require.ElementsMatch(t, []types.EnygmaHistory{testdata.EnygmaHistory1, testdata.EnygmaHistory2}, got)
	})

	t.Run("returns only records matching the resource and block pair", func(t *testing.T) {
		got, err := repo.GetEnygmaHistoryForCheckpoints(
			context.TODO(),
			[]string{"resource-history-1"},
			[]*big.Int{big.NewInt(1001)},
		)
		require.Nil(t, err)

		require.Len(t, got, 1)
		require.Equal(t, testdata.EnygmaHistory1, got[0])
	})
}
