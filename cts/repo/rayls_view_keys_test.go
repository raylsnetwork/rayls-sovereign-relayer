package repo_test

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/raylsnetwork/rayls-sovereign-relayer/cts/domain"
	"github.com/raylsnetwork/rayls-sovereign-relayer/cts/keygen"
	"github.com/raylsnetwork/rayls-sovereign-relayer/cts/repo"
	"github.com/raylsnetwork/rayls-sovereign-relayer/cts/repo/testutil"
	"github.com/raylsnetwork/rayls-sovereign-relayer/cts/service"
	"github.com/stretchr/testify/require"
)

const testTableName = "test_rayls_view_keys"

var dbConf = testutil.DBConfig{
	User:     "user",
	Pass:     "pass",
	Database: "testdb",
}

func createTestTable(t testing.TB, pool *pgxpool.Pool) {
	t.Helper()
	_, err := pool.Exec(context.Background(), `
		CREATE TABLE IF NOT EXISTS test_rayls_view_keys (
			initial_block INTEGER PRIMARY KEY,
			encrypted_secret_key BYTEA NOT NULL,
			public_key BYTEA NOT NULL,
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)
	`)
	require.NoError(t, err)
}

func clearTestTable(t testing.TB, pool *pgxpool.Pool) {
	t.Helper()
	_, err := pool.Exec(context.Background(), `DELETE FROM test_rayls_view_keys`)
	require.NoError(t, err)
}

func TestRaylsViewKeysRepository(t *testing.T) {
	testutil.SetupLogger()

	pool, teardown := testutil.SetupPostgres(t, dbConf)
	defer teardown()

	createTestTable(t, pool)

	repository := repo.NewRaylsViewKeysRepository(testTableName, pool)
	require.NotNil(t, repository)

	t.Run("returns service.ErrNoRaylsViewKeysSet on no keys in repo", func(t *testing.T) {
		clearTestTable(t, pool)

		wantErr := service.ErrNoRaylsViewKeysSet

		blockNumber := uint64(110)
		_, gotErr := repository.GetForBlockNumber(context.TODO(), blockNumber)

		require.Equal(t, wantErr, gotErr)
	})

	t.Run("creates new key", func(t *testing.T) {
		clearTestTable(t, pool)

		keyPair, err := keygen.GenerateRaylsViewKeys()
		require.NoError(t, err)

		keyPair.InitialBlock = 1337

		enc := testutil.NewCeaserEncryptor(10)

		encKeyPair, err := keyPair.Encrypt(enc)
		require.NoError(t, err)

		err = repository.Create(context.TODO(), encKeyPair)
		require.Nil(t, err)

		got, err := repository.GetForBlockNumber(context.TODO(), keyPair.InitialBlock)
		require.Nil(t, err)

		require.Equal(t, encKeyPair, got)
	})
}

func TestRaylsViewKeysRepositoryBatchFunctions(t *testing.T) {
	testutil.SetupLogger()

	pool, teardown := testutil.SetupPostgres(t, dbConf)
	defer teardown()

	createTestTable(t, pool)

	raylsViewKeysRepo := repo.NewRaylsViewKeysRepository(testTableName, pool)
	require.NotNil(t, raylsViewKeysRepo)

	keys := generateEncryptedKeyPairSlice(t)

	t.Run("returns an empty list on no records in database", func(t *testing.T) {
		clearTestTable(t, pool)

		want := []domain.EncryptedRaylsViewKeyPair{}

		got, err := raylsViewKeysRepo.GetAll(context.Background())
		require.Nil(t, err)

		require.Equal(t, want, got)
	})

	t.Run("inserts records into database and gets them", func(t *testing.T) {
		clearTestTable(t, pool)

		err := raylsViewKeysRepo.CreateAll(context.Background(), keys)
		require.Nil(t, err)

		got, err := raylsViewKeysRepo.GetAll(context.Background())
		require.Nil(t, err)

		for _, elem := range got {
			require.Contains(t, keys, elem)
		}
	})
}

func generateEncryptedKeyPairSlice(t testing.TB) []domain.EncryptedRaylsViewKeyPair {
	t.Helper()

	enc := testutil.NewCeaserEncryptor(10)

	firstKeyPair, err := keygen.GenerateRaylsViewKeys()
	require.NoError(t, err)

	firstKeyPair.InitialBlock = 100

	encrFirstKeyPair, err := firstKeyPair.Encrypt(enc)
	require.NoError(t, err)

	secondKeyPair, err := keygen.GenerateRaylsViewKeys()
	require.NoError(t, err)

	secondKeyPair.InitialBlock = 200

	encSecondPair, err := secondKeyPair.Encrypt(enc)
	require.NoError(t, err)

	return []domain.EncryptedRaylsViewKeyPair{encrFirstKeyPair, encSecondPair}
}
