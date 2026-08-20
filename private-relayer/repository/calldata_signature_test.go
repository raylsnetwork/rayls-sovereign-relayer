//go:build integration

package repository_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/raylsnetwork/rayls-sovereign-relayer/private-relayer/repository"
	"github.com/raylsnetwork/rayls-sovereign-relayer/private-relayer/repository/testdata"
	"github.com/raylsnetwork/rayls-sovereign-relayer/testtools"
	"github.com/raylsnetwork/rayls-sovereign-relayer/types"
)

func TestCalldataSignatureRepository_Get(t *testing.T) {
	testtools.SilenceLogger()

	pool, teardown := testtools.SetupPostgres(t, dbConf)
	defer teardown()

	repo := repository.NewCalldataSignatureRepository(pool)

	t.Run("returns an empty list on no signatures in database", func(t *testing.T) {
		want := []types.CalldataSignature{}

		sharedIDs := []string{"exmaple-shared-id"}

		got, err := repo.GetBySharedIDs(context.TODO(), sharedIDs)
		require.Nil(t, err)

		require.Equal(t, want, got)
	})

	// Seed using repository Create method
	err := repo.Create(context.TODO(), testdata.CalldataSignature1)
	require.Nil(t, err)
	err = repo.Create(context.TODO(), testdata.CalldataSignature2)
	require.Nil(t, err)

	t.Run("returns an empty list on no signatures for shared ID", func(t *testing.T) {
		sharedIDs := []string{"exmaple-shared-id"}

		got, err := repo.GetBySharedIDs(context.TODO(), sharedIDs)
		require.Nil(t, err)

		require.Empty(t, got)
	})

	t.Run("retuns calldata signatures for shared ID", func(t *testing.T) {
		want := []types.CalldataSignature{testdata.CalldataSignature1}

		sharedIDs := []string{testdata.CalldataSignature1.SharedId}

		got, err := repo.GetBySharedIDs(context.TODO(), sharedIDs)
		require.Nil(t, err)

		require.Equal(t, want, got)
	})
}

func TestCalldataSignatureRepository_Create(t *testing.T) {
	testtools.SilenceLogger()

	pool, teardown := testtools.SetupPostgres(t, dbConf)
	defer teardown()

	repo := repository.NewCalldataSignatureRepository(pool)

	t.Run("persists calldata signature in repository", func(t *testing.T) {
		want := testdata.CalldataSignature1

		err := repo.Create(context.TODO(), want)
		require.Nil(t, err)

		sharedIDs := []string{want.SharedId}
		gotSlice, err := repo.GetBySharedIDs(context.TODO(), sharedIDs)
		require.Nil(t, err)

		require.Len(t, gotSlice, 1)
		require.Equal(t, want, gotSlice[0])
	})
}

func TestCalldataSignatureRepository_Delete(t *testing.T) {
	testtools.SilenceLogger()

	pool, teardown := testtools.SetupPostgres(t, dbConf)
	defer teardown()

	repo := repository.NewCalldataSignatureRepository(pool)

	// Seed using repository Create method
	err := repo.Create(context.TODO(), testdata.CalldataSignature1)
	require.Nil(t, err)
	err = repo.Create(context.TODO(), testdata.CalldataSignature2)
	require.Nil(t, err)

	t.Run("returns on nothing to delete", func(t *testing.T) {
		sharedIDs := []string{"exmaple-shared-id"}

		gotErr := repo.DeleteBySharedIDs(context.TODO(), sharedIDs)

		require.Nil(t, gotErr)
	})

	t.Run("deletes signatures with shared ID", func(t *testing.T) {
		sharedIDs := []string{testdata.CalldataSignature1.SharedId}

		err := repo.DeleteBySharedIDs(context.TODO(), sharedIDs)
		require.Nil(t, err)

		got, err := repo.GetBySharedIDs(context.TODO(), sharedIDs)
		require.Nil(t, err)

		require.Empty(t, got)
	})
}

func TestCalldataSignatureRepository_BatchCreate(t *testing.T) {
	testtools.SilenceLogger()

	pool, teardown := testtools.SetupPostgres(t, dbConf)
	defer teardown()

	repo := repository.NewCalldataSignatureRepository(pool)

	t.Run("creates signatures in database", func(t *testing.T) {
		want := []types.CalldataSignature{testdata.CalldataSignature1, testdata.CalldataSignature2}

		err := repo.BatchCreate(context.TODO(), want)
		require.Nil(t, err)

		sharedIDs := []string{testdata.CalldataSignature1.SharedId, testdata.CalldataSignature2.SharedId}
		got, err := repo.GetBySharedIDs(context.TODO(), sharedIDs)
		require.Nil(t, err)

		require.Equal(t, want, got)
	})
}
