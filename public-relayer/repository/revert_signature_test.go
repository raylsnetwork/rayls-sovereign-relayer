// Decommissioning Teleport (vanilla, atomic).

//go:build integration

package repository_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/raylsnetwork/rayls-privacy-relayer-api/public-relayer/repository"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/public-relayer/repository/testdata"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/public-relayer/service"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/testtools"
)

var dbConf = testtools.DBConfig{
	User:           "test",
	Pass:           "test",
	Database:       "testdb",
	MigrationsPath: "file://./migrations",
}

func setupRepo(t *testing.T, collectionName string) *repository.RevertSignatureRepository {
	t.Helper()
	testtools.SilenceLogger()

	pool, cleanup := testtools.SetupPostgres(t, dbConf)
	t.Cleanup(cleanup)

	repo := repository.NewRevertSignatureRepository(collectionName, pool)

	return repo
}

func TestNewRevertSignatureRepository(t *testing.T) {
	t.Run("creates repository without error", func(t *testing.T) {
		repo := setupRepo(t, "test_revert_new")
		assert.NotNil(t, repo)
	})
}

func TestRevertSignatureRepository_Create(t *testing.T) {
	t.Run("creates a revert signature record", func(t *testing.T) {
		repo := setupRepo(t, "test_revert_create")

		err := repo.Create(context.Background(), testdata.RevertSig1)
		require.NoError(t, err)

		results, err := repo.GetByIDs(context.Background(), []string{testdata.RevertSig1.ID})
		require.NoError(t, err)
		require.Len(t, results, 1)
		assert.Equal(t, testdata.RevertSig1.ID, results[0].ID)
		assert.Equal(t, testdata.RevertSig1.Data, results[0].Data)
	})

	t.Run("returns error on duplicate ID", func(t *testing.T) {
		repo := setupRepo(t, "test_revert_create_dup")

		err := repo.Create(context.Background(), testdata.RevertSig1)
		require.NoError(t, err)

		err = repo.Create(context.Background(), testdata.RevertSig1)
		require.Error(t, err)
	})
}

func TestRevertSignatureRepository_BatchCreate(t *testing.T) {
	t.Run("returns nil on empty slice", func(t *testing.T) {
		repo := setupRepo(t, "test_revert_batch_empty")

		err := repo.BatchCreate(context.Background(), []service.RevertSignature{})
		require.NoError(t, err)
	})

	t.Run("creates multiple revert signature records", func(t *testing.T) {
		repo := setupRepo(t, "test_revert_batch_multi")

		sigs := []service.RevertSignature{testdata.RevertSig1, testdata.RevertSig2, testdata.RevertSig3}
		err := repo.BatchCreate(context.Background(), sigs)
		require.NoError(t, err)

		results, err := repo.GetByIDs(context.Background(), []string{
			testdata.RevertSig1.ID,
			testdata.RevertSig2.ID,
			testdata.RevertSig3.ID,
		})
		require.NoError(t, err)
		assert.Len(t, results, 3)
	})

	t.Run("silently skips duplicate records", func(t *testing.T) {
		repo := setupRepo(t, "test_revert_batch_dup")

		err := repo.Create(context.Background(), testdata.RevertSig1)
		require.NoError(t, err)

		// Batch includes the existing sig1 plus new sig2
		sigs := []service.RevertSignature{testdata.RevertSig1, testdata.RevertSig2}
		err = repo.BatchCreate(context.Background(), sigs)
		require.NoError(t, err)

		results, err := repo.GetByIDs(context.Background(), []string{
			testdata.RevertSig1.ID,
			testdata.RevertSig2.ID,
		})
		require.NoError(t, err)
		assert.Len(t, results, 2)
	})
}

func TestRevertSignatureRepository_GetByIDs(t *testing.T) {
	t.Run("returns empty slice on empty input", func(t *testing.T) {
		repo := setupRepo(t, "test_revert_get_empty")

		results, err := repo.GetByIDs(context.Background(), []string{})
		require.NoError(t, err)
		assert.Empty(t, results)
	})

	t.Run("returns empty slice when no matching IDs", func(t *testing.T) {
		repo := setupRepo(t, "test_revert_get_nomatch")

		results, err := repo.GetByIDs(context.Background(), []string{"nonexistent-1", "nonexistent-2"})
		require.NoError(t, err)
		assert.Empty(t, results)
	})

	t.Run("returns only matching records when some IDs dont exist", func(t *testing.T) {
		repo := setupRepo(t, "test_revert_get_partial")

		err := repo.Create(context.Background(), testdata.RevertSig1)
		require.NoError(t, err)

		results, err := repo.GetByIDs(context.Background(), []string{testdata.RevertSig1.ID, "nonexistent"})
		require.NoError(t, err)
		require.Len(t, results, 1)
		assert.Equal(t, testdata.RevertSig1.ID, results[0].ID)
	})
}
