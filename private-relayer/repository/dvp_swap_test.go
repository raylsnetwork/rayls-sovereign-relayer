//go:build integration

package repository_test

import (
	"context"
	"math/big"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/raylsnetwork/rayls-sovereign-relayer/private-relayer/repository"
	"github.com/raylsnetwork/rayls-sovereign-relayer/private-relayer/repository/testdata"
	"github.com/raylsnetwork/rayls-sovereign-relayer/testtools"
	"github.com/raylsnetwork/rayls-sovereign-relayer/types"
)

func TestDvpSwapRepository_CreateSwap(t *testing.T) {
	testtools.SilenceLogger()

	pool, teardown := testtools.SetupPostgres(t, dbConf)
	defer teardown()

	repo := repository.NewDvpSwapRepository(pool)

	t.Run("persists swap in database", func(t *testing.T) {
		swap := testdata.DvpSwap1
		err := repo.CreateSwap(context.Background(), &swap)
		require.Nil(t, err)

		got, err := repo.GetSwapBySharedID(context.Background(),testdata.DvpSwap1.SharedID)
		require.Nil(t, err)
		require.NotNil(t, got)
		require.Equal(t, testdata.DvpSwap1.SharedID, got.SharedID)
		require.Equal(t, testdata.DvpSwap1.From, got.From)
		require.Equal(t, testdata.DvpSwap1.To, got.To)
		require.Equal(t, testdata.DvpSwap1.SourceChainID, got.SourceChainID)
		require.Equal(t, testdata.DvpSwap1.DestChainID, got.DestChainID)
		require.Equal(t, testdata.DvpSwap1.TokenInAmount, got.TokenInAmount)
		require.Equal(t, testdata.DvpSwap1.TokenOutAmount, got.TokenOutAmount)
		require.Equal(t, testdata.DvpSwap1.Status, got.Status)
		require.Equal(t, testdata.DvpSwap1.ExpiresAt, got.ExpiresAt)
		require.WithinDuration(t, time.Now(), got.CreatedAt, 5*time.Second)
	})

	t.Run("returns error on duplicate shared ID", func(t *testing.T) {
		swap := testdata.DvpSwap1
		err := repo.CreateSwap(context.Background(), &swap)
		require.NotNil(t, err)
	})
}

func TestDvpSwapRepository_GetSwapBySharedID(t *testing.T) {
	testtools.SilenceLogger()

	pool, teardown := testtools.SetupPostgres(t, dbConf)
	defer teardown()

	repo := repository.NewDvpSwapRepository(pool)

	// Seed using repository CreateSwap method
	swap1 := testdata.DvpSwap1
	err := repo.CreateSwap(context.Background(), &swap1)
	require.Nil(t, err)

	t.Run("returns nil when shared ID not found", func(t *testing.T) {
		got, err := repo.GetSwapBySharedID(context.Background(),"non-existent")
		require.Nil(t, err)
		require.Nil(t, got)
	})

	t.Run("returns swap for matching shared ID", func(t *testing.T) {
		got, err := repo.GetSwapBySharedID(context.Background(),testdata.DvpSwap1.SharedID)
		require.Nil(t, err)
		require.NotNil(t, got)
		require.Equal(t, testdata.DvpSwap1.SharedID, got.SharedID)
		require.Equal(t, testdata.DvpSwap1.TokenInAmount, got.TokenInAmount)
	})
}

func TestDvpSwapRepository_GetExpiredSwaps(t *testing.T) {
	testtools.SilenceLogger()

	pool, teardown := testtools.SetupPostgres(t, dbConf)
	defer teardown()

	repo := repository.NewDvpSwapRepository(pool)

	t.Run("returns empty slice when no expired swaps exist", func(t *testing.T) {
		swap1 := testdata.DvpSwap1
		err := repo.CreateSwap(context.Background(), &swap1)
		require.Nil(t, err)

		got, err := repo.GetExpiredSwaps(context.TODO())
		require.Nil(t, err)

		var want []*types.DvpSwap
		require.Equal(t, want, got)
	})

	t.Run("returns only expired swaps with correct status", func(t *testing.T) {
		swapExpired := testdata.DvpSwapExpired
		err := repo.CreateSwap(context.Background(), &swapExpired)
		require.Nil(t, err)

		got, err := repo.GetExpiredSwaps(context.TODO())
		require.Nil(t, err)

		require.Len(t, got, 1)
		require.Equal(t, testdata.DvpSwapExpired.SharedID, got[0].SharedID)
	})
}

func TestDvpSwapRepository_UpdateSwapStatus(t *testing.T) {
	testtools.SilenceLogger()

	pool, teardown := testtools.SetupPostgres(t, dbConf)
	defer teardown()

	repo := repository.NewDvpSwapRepository(pool)

	// Seed using repository CreateSwap method
	swap1 := testdata.DvpSwap1
	err := repo.CreateSwap(context.Background(), &swap1)
	require.Nil(t, err)

	t.Run("updates swap status", func(t *testing.T) {
		err := repo.UpdateSwapStatus(context.Background(),testdata.DvpSwap1.SharedID, types.DvpSwapCancelled)
		require.Nil(t, err)

		got, err := repo.GetSwapBySharedID(context.Background(),testdata.DvpSwap1.SharedID)
		require.Nil(t, err)
		require.NotNil(t, got)
		require.Equal(t, types.DvpSwapCancelled, got.Status)
	})

	t.Run("returns no error when shared ID not found", func(t *testing.T) {
		err := repo.UpdateSwapStatus(context.Background(),"non-existent", types.DvpSwapCancelled)
		require.Nil(t, err)
	})
}

func TestDvpSwapRepository_CancelSwap(t *testing.T) {
	testtools.SilenceLogger()

	pool, teardown := testtools.SetupPostgres(t, dbConf)
	defer teardown()

	repo := repository.NewDvpSwapRepository(pool)

	// Seed using repository CreateSwap method
	swap1 := testdata.DvpSwap1
	err := repo.CreateSwap(context.Background(), &swap1)
	require.Nil(t, err)

	t.Run("cancels swap and sets cancelled_at", func(t *testing.T) {
		err := repo.CancelSwap(context.Background(),testdata.DvpSwap1.SharedID, types.DvpSwapCancelled)
		require.Nil(t, err)

		got, err := repo.GetSwapBySharedID(context.Background(),testdata.DvpSwap1.SharedID)
		require.Nil(t, err)
		require.NotNil(t, got)
		require.Equal(t, types.DvpSwapCancelled, got.Status)
		require.NotNil(t, got.CancelledAt)
		require.WithinDuration(t, time.Now(), *got.CancelledAt, 5*time.Second)
	})
}

func TestDvpSwapRepository_SetSourcePublicKeyAndAddress(t *testing.T) {
	testtools.SilenceLogger()

	pool, teardown := testtools.SetupPostgres(t, dbConf)
	defer teardown()

	repo := repository.NewDvpSwapRepository(pool)

	// Seed using repository CreateSwap method
	swapExpired := testdata.DvpSwapExpired
	err := repo.CreateSwap(context.Background(), &swapExpired)
	require.Nil(t, err)

	t.Run("sets source public key and from address", func(t *testing.T) {
		err := repo.SetSourcePublicKeyAndAddress(context.Background(),testdata.DvpSwapExpired.SharedID, big.NewInt(9999), "0xNewFrom")
		require.Nil(t, err)

		got, err := repo.GetSwapBySharedID(context.Background(),testdata.DvpSwapExpired.SharedID)
		require.Nil(t, err)
		require.NotNil(t, got)
		require.Equal(t, big.NewInt(9999), got.SourcePaymentPublicKey)
		require.Equal(t, "0xNewFrom", got.From)
	})
}

func TestDvpSwapRepository_SetDestinationPublicKeyAndAddress(t *testing.T) {
	testtools.SilenceLogger()

	pool, teardown := testtools.SetupPostgres(t, dbConf)
	defer teardown()

	repo := repository.NewDvpSwapRepository(pool)

	// Seed using repository CreateSwap method
	swapExpired := testdata.DvpSwapExpired
	err := repo.CreateSwap(context.Background(), &swapExpired)
	require.Nil(t, err)

	t.Run("sets destination public key and to address", func(t *testing.T) {
		err := repo.SetDestinationPublicKeyAndAddress(context.Background(),testdata.DvpSwapExpired.SharedID, big.NewInt(8888), "0xNewTo")
		require.Nil(t, err)

		got, err := repo.GetSwapBySharedID(context.Background(),testdata.DvpSwapExpired.SharedID)
		require.Nil(t, err)
		require.NotNil(t, got)
		require.Equal(t, big.NewInt(8888), got.DestPaymentPublicKey)
		require.Equal(t, "0xNewTo", got.To)
	})
}
