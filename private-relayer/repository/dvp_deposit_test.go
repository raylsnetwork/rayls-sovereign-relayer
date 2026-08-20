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

func TestDvpDepositRepository_CreateDeposit(t *testing.T) {
	testtools.SilenceLogger()

	pool, teardown := testtools.SetupPostgres(t, dbConf)
	defer teardown()

	repo := repository.NewDvpDepositRepository(pool)

	t.Run("persists deposit in database", func(t *testing.T) {
		deposit := testdata.DvpDeposit1
		err := repo.CreateDeposit(context.Background(), &deposit)
		require.Nil(t, err)

		got, err := repo.GetDepositByCommitment(context.Background(),testdata.DvpDeposit1.Commitment)
		require.Nil(t, err)
		require.NotNil(t, got)
		require.Equal(t, testdata.DvpDeposit1.UserAddress, got.UserAddress)
		require.Equal(t, testdata.DvpDeposit1.PublicKey, got.PublicKey)
		require.Equal(t, testdata.DvpDeposit1.TokenAmount, got.TokenAmount)
		require.Equal(t, testdata.DvpDeposit1.Commitment, got.Commitment)
		require.Equal(t, testdata.DvpDeposit1.Status, got.Status)
		require.WithinDuration(t, time.Now(), got.CreatedAt, 5*time.Second)
	})

	t.Run("returns error on duplicate commitment", func(t *testing.T) {
		deposit := testdata.DvpDeposit1
		err := repo.CreateDeposit(context.Background(), &deposit)
		require.NotNil(t, err)
	})
}

func TestDvpDepositRepository_GetDepositByCommitment(t *testing.T) {
	testtools.SilenceLogger()

	pool, teardown := testtools.SetupPostgres(t, dbConf)
	defer teardown()

	repo := repository.NewDvpDepositRepository(pool)

	// Seed using repository CreateDeposit method
	deposit1 := testdata.DvpDeposit1
	err := repo.CreateDeposit(context.Background(), &deposit1)
	require.Nil(t, err)
	deposit2 := testdata.DvpDeposit2
	err = repo.CreateDeposit(context.Background(), &deposit2)
	require.Nil(t, err)

	t.Run("returns nil when commitment not found", func(t *testing.T) {
		got, err := repo.GetDepositByCommitment(context.Background(),big.NewInt(9999))
		require.Nil(t, err)
		require.Nil(t, got)
	})

	t.Run("returns deposit for matching commitment", func(t *testing.T) {
		got, err := repo.GetDepositByCommitment(context.Background(),testdata.DvpDeposit1.Commitment)
		require.Nil(t, err)
		require.NotNil(t, got)
		require.Equal(t, testdata.DvpDeposit1.Commitment, got.Commitment)
		require.Equal(t, testdata.DvpDeposit1.UserAddress, got.UserAddress)
	})
}

func TestDvpDepositRepository_GetDepositByCommitmentAndStatus(t *testing.T) {
	testtools.SilenceLogger()

	pool, teardown := testtools.SetupPostgres(t, dbConf)
	defer teardown()

	repo := repository.NewDvpDepositRepository(pool)

	// Seed using repository CreateDeposit method
	deposit1 := testdata.DvpDeposit1
	err := repo.CreateDeposit(context.Background(), &deposit1)
	require.Nil(t, err)
	deposit2 := testdata.DvpDeposit2
	err = repo.CreateDeposit(context.Background(), &deposit2)
	require.Nil(t, err)

	t.Run("returns nil when no deposit matches commitment", func(t *testing.T) {
		got, err := repo.GetDepositByCommitmentAndStatus(context.Background(),big.NewInt(9999), types.DvpDepositPending)
		require.Nil(t, err)
		require.Nil(t, got)
	})

	t.Run("returns nil when commitment matches but status does not", func(t *testing.T) {
		got, err := repo.GetDepositByCommitmentAndStatus(context.Background(),testdata.DvpDeposit1.Commitment, types.DvpDepositUnspent)
		require.Nil(t, err)
		require.Nil(t, got)
	})

	t.Run("returns deposit when commitment and status both match", func(t *testing.T) {
		got, err := repo.GetDepositByCommitmentAndStatus(context.Background(),testdata.DvpDeposit1.Commitment, types.DvpDepositPending)
		require.Nil(t, err)
		require.NotNil(t, got)
		require.Equal(t, testdata.DvpDeposit1.Commitment, got.Commitment)
	})
}

func TestDvpDepositRepository_ConfirmDeposit(t *testing.T) {
	testtools.SilenceLogger()

	pool, teardown := testtools.SetupPostgres(t, dbConf)
	defer teardown()

	repo := repository.NewDvpDepositRepository(pool)

	// Seed using repository CreateDeposit method
	deposit1 := testdata.DvpDeposit1
	err := repo.CreateDeposit(context.Background(), &deposit1)
	require.Nil(t, err)

	t.Run("updates tree number and status to unspent", func(t *testing.T) {
		err := repo.ConfirmDeposit(context.Background(),testdata.DvpDeposit1.Commitment, big.NewInt(5))
		require.Nil(t, err)

		got, err := repo.GetDepositByCommitment(context.Background(),testdata.DvpDeposit1.Commitment)
		require.Nil(t, err)
		require.NotNil(t, got)
		require.Equal(t, types.DvpDepositUnspent, got.Status)
		require.Equal(t, 5, got.TreeNumber)
	})
}

func TestDvpDepositRepository_UpdateDepositStatus(t *testing.T) {
	testtools.SilenceLogger()

	pool, teardown := testtools.SetupPostgres(t, dbConf)
	defer teardown()

	repo := repository.NewDvpDepositRepository(pool)

	// Seed using repository CreateDeposit method
	deposit1 := testdata.DvpDeposit1
	err := repo.CreateDeposit(context.Background(), &deposit1)
	require.Nil(t, err)

	t.Run("updates deposit status", func(t *testing.T) {
		err := repo.UpdateDepositStatus(context.Background(),testdata.DvpDeposit1.Commitment, types.DvpDepositSpent)
		require.Nil(t, err)

		got, err := repo.GetDepositByCommitment(context.Background(),testdata.DvpDeposit1.Commitment)
		require.Nil(t, err)
		require.NotNil(t, got)
		require.Equal(t, types.DvpDepositSpent, got.Status)
	})
}

func TestDvpDepositRepository_BatchUpdateStatusForCommitments(t *testing.T) {
	testtools.SilenceLogger()

	pool, teardown := testtools.SetupPostgres(t, dbConf)
	defer teardown()

	repo := repository.NewDvpDepositRepository(pool)

	// Seed using repository CreateDeposit method
	deposit1 := testdata.DvpDeposit1
	err := repo.CreateDeposit(context.Background(), &deposit1)
	require.Nil(t, err)
	deposit2 := testdata.DvpDeposit2
	err = repo.CreateDeposit(context.Background(), &deposit2)
	require.Nil(t, err)

	t.Run("updates status for multiple commitments in a transaction", func(t *testing.T) {
		commitments := []string{testdata.DvpDeposit1.Commitment.String(), testdata.DvpDeposit2.Commitment.String()}
		err := repo.BatchUpdateStatusForCommitments(context.TODO(), commitments, types.DvpDepositLocked)
		require.Nil(t, err)

		got1, err := repo.GetDepositByCommitment(context.Background(), testdata.DvpDeposit1.Commitment)
		require.Nil(t, err)
		require.Equal(t, types.DvpDepositLocked, got1.Status)

		got2, err := repo.GetDepositByCommitment(context.Background(), testdata.DvpDeposit2.Commitment)
		require.Nil(t, err)
		require.Equal(t, types.DvpDepositLocked, got2.Status)
	})

	t.Run("returns no error when no commitments match", func(t *testing.T) {
		err := repo.BatchUpdateStatusForCommitments(context.TODO(), []string{"non-existent"}, types.DvpDepositLocked)
		require.Nil(t, err)
	})
}

func TestDvpDepositRepository_GetNonFungibleDeposit(t *testing.T) {
	testtools.SilenceLogger()

	pool, teardown := testtools.SetupPostgres(t, dbConf)
	defer teardown()

	repo := repository.NewDvpDepositRepository(pool)

	// Seed using repository CreateDeposit method
	deposit3NFT := testdata.DvpDeposit3NFT
	err := repo.CreateDeposit(context.Background(), &deposit3NFT)
	require.Nil(t, err)

	t.Run("returns nil when no NFT deposit matches", func(t *testing.T) {
		got, err := repo.GetNonFungibleDeposit(context.Background(),
			"wrong-id",
			"0xTokenB",
			"0xUser2",
			types.DvpERC721,
			types.DvpDepositUnspent,
		)
		require.Nil(t, err)
		require.Nil(t, got)
	})

	t.Run("returns NFT deposit matching token ID and address", func(t *testing.T) {
		got, err := repo.GetNonFungibleDeposit(context.Background(),
			"nft-token-id-1",
			"0xTokenB",
			"0xUser2",
			types.DvpERC721,
			types.DvpDepositUnspent,
		)
		require.Nil(t, err)
		require.NotNil(t, got)
		require.Equal(t, testdata.DvpDeposit3NFT.Commitment, got.Commitment)
	})
}

func TestDvpDepositRepository_GetFungibleDeposits(t *testing.T) {
	testtools.SilenceLogger()

	pool, teardown := testtools.SetupPostgres(t, dbConf)
	defer teardown()

	repo := repository.NewDvpDepositRepository(pool)

	// Seed using repository CreateDeposit method
	deposit1 := testdata.DvpDeposit1
	err := repo.CreateDeposit(context.Background(), &deposit1)
	require.Nil(t, err)
	deposit2 := testdata.DvpDeposit2
	err = repo.CreateDeposit(context.Background(), &deposit2)
	require.Nil(t, err)

	t.Run("returns empty when no fungible deposits match", func(t *testing.T) {
		got, err := repo.GetFungibleDeposits(context.Background(),"0xOtherToken", "0xUser1", types.DvpERC20, types.DvpDepositPending)
		require.Nil(t, err)
		require.Empty(t, got)
	})

	t.Run("returns matching fungible deposits", func(t *testing.T) {
		got, err := repo.GetFungibleDeposits(context.Background(),"0xTokenA", "0xUser1", types.DvpERC20, types.DvpDepositPending)
		require.Nil(t, err)
		require.Len(t, got, 1)
		require.Equal(t, testdata.DvpDeposit1.Commitment, got[0].Commitment)
	})
}

func TestDvpDepositRepository_GetDepositsByToken(t *testing.T) {
	testtools.SilenceLogger()

	pool, teardown := testtools.SetupPostgres(t, dbConf)
	defer teardown()

	repo := repository.NewDvpDepositRepository(pool)

	// Seed using repository CreateDeposit method
	deposit3NFT := testdata.DvpDeposit3NFT
	err := repo.CreateDeposit(context.Background(), &deposit3NFT)
	require.Nil(t, err)

	t.Run("returns deposits matching all filters", func(t *testing.T) {
		got, err := repo.GetDepositsByToken(context.Background(),
			"0xTokenB",
			"nft-token-id-1",
			types.DvpERC721,
			"0xUser2",
			types.DvpDepositUnspent,
		)
		require.Nil(t, err)
		require.Len(t, got, 1)
		require.Equal(t, testdata.DvpDeposit3NFT.Commitment, got[0].Commitment)
	})

	t.Run("returns empty when user address does not match", func(t *testing.T) {
		got, err := repo.GetDepositsByToken(context.Background(),
			"0xTokenB",
			"nft-token-id-1",
			types.DvpERC721,
			"0xWrongUser",
			types.DvpDepositUnspent,
		)
		require.Nil(t, err)
		require.Empty(t, got)
	})
}

func TestDvpDepositRepository_UpsertDepositNullifier(t *testing.T) {
	testtools.SilenceLogger()

	pool, teardown := testtools.SetupPostgres(t, dbConf)
	defer teardown()

	repo := repository.NewDvpDepositRepository(pool)

	// Seed using repository CreateDeposit method
	deposit1 := testdata.DvpDeposit1
	err := repo.CreateDeposit(context.Background(), &deposit1)
	require.Nil(t, err)

	t.Run("updates nullifier for matching commitment", func(t *testing.T) {
		err := repo.UpsertDepositNullifier(context.Background(),testdata.DvpDeposit1.Commitment, big.NewInt(7777))
		require.Nil(t, err)

		got, err := repo.GetDepositByCommitment(context.Background(),testdata.DvpDeposit1.Commitment)
		require.Nil(t, err)
		require.NotNil(t, got)
		require.Equal(t, big.NewInt(7777), got.Nullifier)
	})
}

func TestDvpDepositRepository_BatchUpsertNullifiers(t *testing.T) {
	testtools.SilenceLogger()

	pool, teardown := testtools.SetupPostgres(t, dbConf)
	defer teardown()

	repo := repository.NewDvpDepositRepository(pool)

	// Seed using repository CreateDeposit method
	deposit1 := testdata.DvpDeposit1
	err := repo.CreateDeposit(context.Background(), &deposit1)
	require.Nil(t, err)
	deposit2 := testdata.DvpDeposit2
	err = repo.CreateDeposit(context.Background(), &deposit2)
	require.Nil(t, err)

	t.Run("returns no error for empty map", func(t *testing.T) {
		err := repo.BatchUpsertNullifiers(context.TODO(), map[string]string{})
		require.Nil(t, err)
	})

	t.Run("updates nullifiers for multiple deposits in bulk", func(t *testing.T) {
		nullifierMap := map[string]string{
			testdata.DvpDeposit1.Commitment.String(): "1111",
			testdata.DvpDeposit2.Commitment.String(): "2222",
		}
		err := repo.BatchUpsertNullifiers(context.TODO(), nullifierMap)
		require.Nil(t, err)

		got1, err := repo.GetDepositByCommitment(context.Background(), testdata.DvpDeposit1.Commitment)
		require.Nil(t, err)
		require.Equal(t, big.NewInt(1111), got1.Nullifier)

		got2, err := repo.GetDepositByCommitment(context.Background(), testdata.DvpDeposit2.Commitment)
		require.Nil(t, err)
		require.Equal(t, big.NewInt(2222), got2.Nullifier)
	})
}

func TestDvpDepositRepository_UpdateDepositStatusByNullifier(t *testing.T) {
	testtools.SilenceLogger()

	pool, teardown := testtools.SetupPostgres(t, dbConf)
	defer teardown()

	repo := repository.NewDvpDepositRepository(pool)

	// Seed two deposits that share the same nullifier value across two
	// different token addresses — this is the cross-vault collision the
	// fix is meant to disambiguate.
	depositTokenA := testdata.DvpDeposit2
	require.Nil(t, repo.CreateDeposit(context.Background(), &depositTokenA))

	depositTokenB := testdata.DvpDeposit2
	depositTokenB.Commitment = big.NewInt(8002)
	depositTokenB.TokenAddress = "0xTokenB"
	require.Nil(t, repo.CreateDeposit(context.Background(), &depositTokenB))

	t.Run("flips only the deposit matching the (token_address, nullifier) pair", func(t *testing.T) {
		err := repo.UpdateDepositStatusByNullifier(
			context.Background(),
			depositTokenA.TokenAddress,
			big.NewInt(5555),
			types.DvpDepositSpent,
		)
		require.Nil(t, err)

		gotA, err := repo.GetDepositByCommitment(context.Background(), depositTokenA.Commitment)
		require.Nil(t, err)
		require.NotNil(t, gotA)
		require.Equal(t, types.DvpDepositSpent, gotA.Status)

		gotB, err := repo.GetDepositByCommitment(context.Background(), depositTokenB.Commitment)
		require.Nil(t, err)
		require.NotNil(t, gotB)
		require.Equal(t, depositTokenB.Status, gotB.Status, "colliding deposit on a different token must not be touched")
	})

	t.Run("no-op when no row matches the (token_address, nullifier) pair", func(t *testing.T) {
		err := repo.UpdateDepositStatusByNullifier(
			context.Background(),
			"0xUnknownToken",
			big.NewInt(5555),
			types.DvpDepositSpent,
		)
		require.Nil(t, err)
	})
}
