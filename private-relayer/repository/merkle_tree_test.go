//go:build integration

package repository_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/raylsnetwork/rayls-sovereign-relayer/private-relayer/repository"
	"github.com/raylsnetwork/rayls-sovereign-relayer/private-relayer/repository/testdata"
	"github.com/raylsnetwork/rayls-sovereign-relayer/testtools"
)

func TestMerkleTreeRepository_CreateMerkleTree(t *testing.T) {
	testtools.SilenceLogger()

	pool, teardown := testtools.SetupPostgres(t, dbConf)
	defer teardown()

	repo := repository.NewMerkleTreeRepository(pool)

	t.Run("creates and retrieves a merkle tree", func(t *testing.T) {
		tree := testdata.MerkleTree1
		err := repo.CreateMerkleTree(context.Background(), &tree)
		require.Nil(t, err)

		got, err := repo.GetByNumberAndTokenAddress(context.Background(), testdata.MerkleTree1.Number, testdata.MerkleTree1.TokenAddress)
		require.Nil(t, err)
		require.NotNil(t, got)
		require.Equal(t, testdata.MerkleTree1.Type, got.Type)
		require.Equal(t, testdata.MerkleTree1.TokenAddress, got.TokenAddress)
		require.Equal(t, testdata.MerkleTree1.Number, got.Number)
		require.Equal(t, testdata.MerkleTree1.Depth, got.Depth)
		require.Equal(t, testdata.MerkleTree1.Leaves, got.Leaves)
		require.WithinDuration(t, time.Now(), got.CreatedAt, 5*time.Second)
	})

	t.Run("returns error on duplicate type number and token_address", func(t *testing.T) {
		tree := testdata.MerkleTree1
		err := repo.CreateMerkleTree(context.Background(), &tree)
		require.NotNil(t, err)
	})
}

func TestMerkleTreeRepository_GetLatestByTokenAddress(t *testing.T) {
	testtools.SilenceLogger()

	pool, teardown := testtools.SetupPostgres(t, dbConf)
	defer teardown()

	repo := repository.NewMerkleTreeRepository(pool)

	// Seed using repository CreateMerkleTree method
	tree1 := testdata.MerkleTree1
	tree2 := testdata.MerkleTree2
	tree3 := testdata.MerkleTree3
	err := repo.CreateMerkleTree(context.Background(), &tree1)
	require.Nil(t, err)
	err = repo.CreateMerkleTree(context.Background(), &tree2)
	require.Nil(t, err)
	err = repo.CreateMerkleTree(context.Background(), &tree3)
	require.Nil(t, err)

	t.Run("returns nil when no tree exists for token address", func(t *testing.T) {
		got, err := repo.GetLatestByTokenAddress(context.Background(), "0xNonExistent")
		require.Nil(t, err)
		require.Nil(t, got)
	})

	t.Run("returns the tree with the highest number for the token address", func(t *testing.T) {
		got, err := repo.GetLatestByTokenAddress(context.Background(), "0xTokenAddressAAA")
		require.Nil(t, err)
		require.NotNil(t, got)
		require.Equal(t, testdata.MerkleTree2.Number, got.Number)
		require.Equal(t, testdata.MerkleTree2.Leaves, got.Leaves)
	})

	t.Run("returns the only tree when only one exists for token address", func(t *testing.T) {
		got, err := repo.GetLatestByTokenAddress(context.Background(), "0xTokenAddressBBB")
		require.Nil(t, err)
		require.NotNil(t, got)
		require.Equal(t, testdata.MerkleTree3.Number, got.Number)
		require.Equal(t, testdata.MerkleTree3.TokenAddress, got.TokenAddress)
	})
}

func TestMerkleTreeRepository_GetByNumberAndTokenAddress(t *testing.T) {
	testtools.SilenceLogger()

	pool, teardown := testtools.SetupPostgres(t, dbConf)
	defer teardown()

	repo := repository.NewMerkleTreeRepository(pool)

	// Seed using repository CreateMerkleTree method
	tree1 := testdata.MerkleTree1
	tree2 := testdata.MerkleTree2
	tree3 := testdata.MerkleTree3
	err := repo.CreateMerkleTree(context.Background(), &tree1)
	require.Nil(t, err)
	err = repo.CreateMerkleTree(context.Background(), &tree2)
	require.Nil(t, err)
	err = repo.CreateMerkleTree(context.Background(), &tree3)
	require.Nil(t, err)

	t.Run("returns nil when no tree matches", func(t *testing.T) {
		got, err := repo.GetByNumberAndTokenAddress(context.Background(), 999, "0xTokenAddressAAA")
		require.Nil(t, err)
		require.Nil(t, got)
	})

	t.Run("returns correct tree by number and token address", func(t *testing.T) {
		got, err := repo.GetByNumberAndTokenAddress(context.Background(), 1, "0xTokenAddressAAA")
		require.Nil(t, err)
		require.NotNil(t, got)
		require.Equal(t, testdata.MerkleTree1.Number, got.Number)
		require.Equal(t, testdata.MerkleTree1.TokenAddress, got.TokenAddress)
	})

	t.Run("does not cross-match token addresses", func(t *testing.T) {
		got, err := repo.GetByNumberAndTokenAddress(context.Background(), 1, "0xTokenAddressBBB")
		require.Nil(t, err)
		require.NotNil(t, got)
		require.Equal(t, testdata.MerkleTree3.Number, got.Number)
		require.Equal(t, testdata.MerkleTree3.TokenAddress, got.TokenAddress)
	})
}

func TestMerkleTreeRepository_InsertLeaves(t *testing.T) {
	testtools.SilenceLogger()

	pool, teardown := testtools.SetupPostgres(t, dbConf)
	defer teardown()

	repo := repository.NewMerkleTreeRepository(pool)

	// Seed using repository CreateMerkleTree method
	tree1 := testdata.MerkleTree1
	err := repo.CreateMerkleTree(context.Background(), &tree1)
	require.Nil(t, err)

	t.Run("returns error when no tree matches filter", func(t *testing.T) {
		err := repo.InsertLeaves(context.Background(), "0xNonExistent", 1, []string{"new-leaf"})
		require.NotNil(t, err)
		require.Contains(t, err.Error(), "no merkle tree found")
	})

	t.Run("appends new leaves to existing tree", func(t *testing.T) {
		err := repo.InsertLeaves(context.Background(), "0xTokenAddressAAA", 1, []string{"leaf-c", "leaf-d"})
		require.Nil(t, err)

		got, err := repo.GetByNumberAndTokenAddress(context.Background(), 1, "0xTokenAddressAAA")
		require.Nil(t, err)
		require.NotNil(t, got)
		require.Equal(t, []string{"leaf-a", "leaf-b", "leaf-c", "leaf-d"}, got.Leaves)
	})

	t.Run("does not insert duplicate leaves", func(t *testing.T) {
		err := repo.InsertLeaves(context.Background(), "0xTokenAddressAAA", 1, []string{"leaf-a"})
		require.Nil(t, err)

		got, err := repo.GetByNumberAndTokenAddress(context.Background(), 1, "0xTokenAddressAAA")
		require.Nil(t, err)
		require.NotNil(t, got)
		require.ElementsMatch(t, []string{"leaf-a", "leaf-b", "leaf-c", "leaf-d"}, got.Leaves)
	})
}

func TestMerkleTreeRepository_DeleteMerkleTree(t *testing.T) {
	testtools.SilenceLogger()

	pool, teardown := testtools.SetupPostgres(t, dbConf)
	defer teardown()

	repo := repository.NewMerkleTreeRepository(pool)

	t.Run("succeeds on empty collection", func(t *testing.T) {
		err := repo.DeleteMerkleTree(context.Background())
		require.Nil(t, err)
	})

	// Seed using repository CreateMerkleTree method
	tree1 := testdata.MerkleTree1
	tree2 := testdata.MerkleTree2
	tree3 := testdata.MerkleTree3
	err := repo.CreateMerkleTree(context.Background(), &tree1)
	require.Nil(t, err)
	err = repo.CreateMerkleTree(context.Background(), &tree2)
	require.Nil(t, err)
	err = repo.CreateMerkleTree(context.Background(), &tree3)
	require.Nil(t, err)

	t.Run("deletes all trees from the collection", func(t *testing.T) {
		err := repo.DeleteMerkleTree(context.Background())
		require.Nil(t, err)

		got, err := repo.GetLatestByTokenAddress(context.Background(), "0xTokenAddressAAA")
		require.Nil(t, err)
		require.Nil(t, got)
	})
}
