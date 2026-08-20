//go:build integration

package repository_test

import (
	"context"
	"math/big"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/raylsnetwork/rayls-privacy-relayer-api/listener"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/private-relayer/repository"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/private-relayer/repository/testdata"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/testtools"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/types"
)

func TestLastProcessedBlockRepository_Get(t *testing.T) {
	testtools.SilenceLogger()

	pool, teardown := testtools.SetupPostgres(t, dbConf)
	defer teardown()

	repo := repository.NewLastProcessedBlockRepository(pool)

	t.Run("returns ErrLastProcessedBlockNotFound on no record for chain id", func(t *testing.T) {
		wantErr := listener.ErrLastProcessedBlockNotFound

		chainID := testdata.ChainID1

		_, gotErr := repo.Get(context.TODO(), chainID)

		require.ErrorIs(t, gotErr, wantErr)
	})

	// Seed using repository Create method
	err := repo.Create(context.TODO(), testdata.ChainID1, testdata.LastProcessedBlock1)
	require.Nil(t, err)
	err = repo.Create(context.TODO(), testdata.ChainID2, testdata.LastProcessedBlock2)
	require.Nil(t, err)

	t.Run("returns block number on existing record for chain id", func(t *testing.T) {
		want := testdata.LastProcessedBlock1

		chainID := testdata.ChainID1

		got, err := repo.Get(context.TODO(), chainID)
		require.Nil(t, err)

		require.Equal(t, want, got)
	})
}

func TestLastProcessedBlockReposiory_Create(t *testing.T) {
	testtools.SilenceLogger()

	pool, teardown := testtools.SetupPostgres(t, dbConf)
	defer teardown()

	repo := repository.NewLastProcessedBlockRepository(pool)

	t.Run("persists last processed block for chain ID in database", func(t *testing.T) {
		want := new(big.Int).SetUint64(7007)

		chainID := types.DocumentIdLastProcessedBlockPrivateHub

		err := repo.Create(context.TODO(), chainID, want)
		require.Nil(t, err)

		got, err := repo.Get(context.TODO(), chainID)
		require.Nil(t, err)

		require.Equal(t, want, got)
	})

	t.Run("returns ErrLastProcessedBlockAlreadyExists if a record for chain ID aready exists", func(t *testing.T) {
		wantErr := listener.ErrLastProcessedBlockAlreadyExists

		chainID := types.DocumentIdLastProcessedBlockPrivateHub
		lastBlock := new(big.Int).SetUint64(7007)

		gotErr := repo.Create(context.TODO(), chainID, lastBlock)

		require.ErrorIs(t, gotErr, wantErr)
	})
}

func TestLastProcessedBlockRepository_Update(t *testing.T) {
	testtools.SilenceLogger()

	pool, teardown := testtools.SetupPostgres(t, dbConf)
	defer teardown()

	repo := repository.NewLastProcessedBlockRepository(pool)

	t.Run("returns ErrLastProcessedBlockNotFound if a record for chain ID doesn't exist", func(t *testing.T) {
		wantErr := listener.ErrLastProcessedBlockNotFound

		chainID := testdata.ChainID1
		lastBlock := new(big.Int).SetUint64(1234)

		gotErr := repo.Update(context.TODO(), chainID, lastBlock)

		require.ErrorIs(t, gotErr, wantErr)
	})

	// Seed using repository Create method
	err := repo.Create(context.TODO(), testdata.ChainID1, testdata.LastProcessedBlock1)
	require.Nil(t, err)

	t.Run("updates last processed block for chain ID", func(t *testing.T) {
		want := new(big.Int).SetUint64(1234)

		chainID := testdata.ChainID1

		err := repo.Update(context.TODO(), chainID, want)
		require.Nil(t, err)

		got, err := repo.Get(context.TODO(), chainID)
		require.Nil(t, err)

		require.Equal(t, want, got)
	})
}

func TestLastProcessedBlockRepository_UpdateWithLock(t *testing.T) {
	testtools.SilenceLogger()

	pool, teardown := testtools.SetupPostgres(t, dbConf)
	defer teardown()

	repo := repository.NewLastProcessedBlockRepository(pool)

	// Seed using repository Create method
	err := repo.Create(context.TODO(), testdata.ChainID1, testdata.LastProcessedBlock1)
	require.Nil(t, err)

	t.Run("updates last processed block with lock", func(t *testing.T) {
		want := new(big.Int).SetUint64(9999)

		chainID := testdata.ChainID1

		err := repo.UpdateWithLock(context.TODO(), chainID, want)
		require.Nil(t, err)

		got, err := repo.Get(context.TODO(), chainID)
		require.Nil(t, err)

		require.Equal(t, want, got)
	})

	t.Run("returns ErrLastProcessedBlockNotFound if a record for chain ID doesn't exist", func(t *testing.T) {
		wantErr := listener.ErrLastProcessedBlockNotFound

		chainID := types.DocumentIdLastProcessedBlockPrivateNode
		lastBlock := new(big.Int).SetUint64(1234)

		gotErr := repo.UpdateWithLock(context.TODO(), chainID, lastBlock)

		require.ErrorIs(t, gotErr, wantErr)
	})
}
