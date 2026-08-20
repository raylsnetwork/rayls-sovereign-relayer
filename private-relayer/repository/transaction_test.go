//go:build integration

package repository_test

import (
	"context"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	"github.com/raylsnetwork/rayls-privacy-relayer-api/private-relayer/repository"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/private-relayer/repository/testdata"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/testtools"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/types"
)

var dbConf = testtools.DBConfig{
	User:           "test",
	Pass:           "test",
	Database:       "testdb",
	MigrationsPath: "file://./migrations",
}

func TestTransactionsRepository_GetTransactions(t *testing.T) {
	testtools.SilenceLogger()

	pool, teardown := testtools.SetupPostgres(t, dbConf)
	defer teardown()

	repo := repository.NewTransactionRepository(pool)

	t.Run("returns an empty slice on no transactions in database", func(t *testing.T) {
		var want []types.Transaction

		hashes := []string{"example-hash"}

		got, err := repo.GetByHashes(context.TODO(), hashes)
		require.Nil(t, err)

		require.Equal(t, want, got)
	})

	// Seed using repository Create method
	err := repo.Create(context.TODO(), testdata.Transaction1)
	require.NoError(t, err)
	err = repo.Create(context.TODO(), testdata.Transaction2)
	require.NoError(t, err)
	err = repo.Create(context.TODO(), testdata.Transaction3)
	require.NoError(t, err)

	t.Run("returns an empty slice on no transacitons with matching hashes", func(t *testing.T) {
		var want []types.Transaction

		hashes := []string{"example-hash"}

		got, err := repo.GetByHashes(context.TODO(), hashes)
		require.Nil(t, err)

		require.Equal(t, want, got)
	})

	t.Run("returns transaction with matching hash", func(t *testing.T) {
		want := []types.Transaction{testdata.Transaction1}

		hashes := []string{testdata.Transaction1.TxHash}

		got, err := repo.GetByHashes(context.TODO(), hashes)
		require.Nil(t, err)

		require.Equal(t, want, got)
	})
}

func TestTransactionRepository_GetByState(t *testing.T) {
	testtools.SilenceLogger()

	pool, teardown := testtools.SetupPostgres(t, dbConf)
	defer teardown()

	repo := repository.NewTransactionRepository(pool)

	t.Run("returns an empty slice on no transactions in database", func(t *testing.T) {
		var want []types.Transaction

		state := types.SourcePublish

		got, err := repo.GetByState(context.Background(), state)
		require.Nil(t, err)

		require.Equal(t, want, got)
	})

	// Seed using repository Create method
	err := repo.Create(context.TODO(), testdata.Transaction1)
	require.NoError(t, err)
	err = repo.Create(context.TODO(), testdata.Transaction2)
	require.NoError(t, err)
	err = repo.Create(context.TODO(), testdata.Transaction3)
	require.NoError(t, err)

	t.Run("returns an empty slice on no transactions with matching state", func(t *testing.T) {
		var want []types.Transaction

		state := types.SourceFinalized

		got, err := repo.GetByState(context.Background(), state)
		require.Nil(t, err)

		require.Equal(t, want, got)
	})

	t.Run("returns transaction with matching state", func(t *testing.T) {
		want := []types.Transaction{testdata.Transaction1}

		state := testdata.Transaction1.State

		got, err := repo.GetByState(context.Background(), state)
		require.Nil(t, err)

		require.Equal(t, want, got)
	})
}

func TestTransactionRepository_GetByState_WithLimit(t *testing.T) {
	testtools.SilenceLogger()

	pool, teardown := testtools.SetupPostgres(t, dbConf)
	defer teardown()

	repo := repository.NewTransactionRepository(pool)

	state := types.HubNotifiedExec

	seedTransactions := []types.Transaction{
		{
			SharedID: "first transaction",
			MsgID:    common.HexToHash("0xc0febabe"),
			State:    state,
		},
		{
			SharedID: "second transaction",
			MsgID:    common.HexToHash("0xdeadc0de"),
			State:    state,
		},
		{
			SharedID: "third transaction",
			MsgID:    common.HexToHash("0xc0cac01a"),
			State:    state,
		},
		{
			SharedID: "fourth transaction",
			MsgID:    common.HexToHash("0xaaaaaaaa"),
			State:    state,
		},
		{
			SharedID: "fifth transaction",
			MsgID:    common.HexToHash("0x77777777"),
			State:    state,
		},
	}

	for _, tx := range seedTransactions {
		err := repo.Create(context.TODO(), tx)
		require.NoError(t, err)
	}

	t.Run("limits transactions to one", func(t *testing.T) {
		limit := 1

		got, err := repo.GetByState(context.TODO(), state, repository.WithLimit(limit))
		require.Nil(t, err)

		require.Equal(t, limit, len(got))
	})

	t.Run("limits transactions to two", func(t *testing.T) {
		limit := 2

		got, err := repo.GetByState(context.TODO(), state, repository.WithLimit(limit))
		require.Nil(t, err)

		require.Equal(t, limit, len(got))
	})

	// PostgreSQL LIMIT 0 returns no rows, but we maintain backward compatibility
	// by returning all rows when limit is 0
	t.Run("returns all transactions when limit is set to zero", func(t *testing.T) {
		limit := 0

		got, err := repo.GetByState(context.TODO(), state, repository.WithLimit(limit))
		require.Nil(t, err)

		require.Equal(t, len(seedTransactions), len(got))
	})

	t.Run("returns all transactions when less elements than limit", func(t *testing.T) {
		limit := 10

		got, err := repo.GetByState(context.TODO(), state, repository.WithLimit(limit))
		require.Nil(t, err)

		require.Equal(t, len(seedTransactions), len(got))
	})
}

func TestTransactionRepository_GetByStateOutcomeAndAtomicity(t *testing.T) {
	testtools.SilenceLogger()

	pool, teardown := testtools.SetupPostgres(t, dbConf)
	defer teardown()

	repo := repository.NewTransactionRepository(pool)

	t.Run("returns an empty slice on no transactions in database", func(t *testing.T) {
		var want []types.Transaction

		got, err := repo.GetByStateOutcomeAndAtomicity(
			context.TODO(), types.SourcePublish, types.OutcomePending, false,
		)
		require.Nil(t, err)

		require.Equal(t, want, got)
	})

	// Seed using repository Create method
	err := repo.Create(context.TODO(), testdata.Transaction1)
	require.NoError(t, err)
	err = repo.Create(context.TODO(), testdata.Transaction2)
	require.NoError(t, err)
	err = repo.Create(context.TODO(), testdata.Transaction3)
	require.NoError(t, err)

	t.Run("returns an empty slice on no transactions with matching state, outcome, and atomicity", func(t *testing.T) {
		var want []types.Transaction

		got, err := repo.GetByStateOutcomeAndAtomicity(
			context.TODO(), types.SourceFinalized, types.OutcomePending, false,
		)
		require.Nil(t, err)

		require.Equal(t, want, got)
	})

	t.Run("returns an empty slice on matching state but not atomicity", func(t *testing.T) {
		var want []types.Transaction

		got, err := repo.GetByStateOutcomeAndAtomicity(
			context.TODO(), testdata.Transaction1.State, testdata.Transaction1.Outcome, false,
		)
		require.Nil(t, err)

		require.Equal(t, want, got)
	})

	t.Run("returns an empty slice on matching atomicity but not state", func(t *testing.T) {
		var want []types.Transaction

		got, err := repo.GetByStateOutcomeAndAtomicity(
			context.TODO(), types.SourceFinalized, types.OutcomePending, true,
		)
		require.Nil(t, err)

		require.Equal(t, want, got)
	})

	t.Run("returns transaction with matching state, outcome, and atomicity", func(t *testing.T) {
		want := []types.Transaction{testdata.Transaction1}

		got, err := repo.GetByStateOutcomeAndAtomicity(
			context.TODO(), testdata.Transaction1.State, testdata.Transaction1.Outcome, true,
		)
		require.Nil(t, err)

		require.Equal(t, want, got)
	})
}

func TestTransactionRepository_GetByStateOutcomeAndAtomicity_WithLimit(t *testing.T) {
	testtools.SilenceLogger()

	pool, teardown := testtools.SetupPostgres(t, dbConf)
	defer teardown()

	repo := repository.NewTransactionRepository(pool)

	state := types.HubNotifiedExec
	outcome := types.OutcomePending
	isAtomic := true

	seedTransactions := []types.Transaction{
		{
			SharedID: "first transaction",
			MsgID:    common.HexToHash("0xc0febabe"),
			State:    state,
			Outcome:  outcome,
			IsAtomic: isAtomic,
		},
		{
			SharedID: "second transaction",
			MsgID:    common.HexToHash("0xdeadc0de"),
			State:    state,
			Outcome:  outcome,
			IsAtomic: isAtomic,
		},
		{
			SharedID: "third transaction",
			MsgID:    common.HexToHash("0xc0cac01a"),
			State:    state,
			Outcome:  outcome,
			IsAtomic: isAtomic,
		},
		{
			SharedID: "fourth transaction",
			MsgID:    common.HexToHash("0xaaaaaaaa"),
			State:    state,
			Outcome:  outcome,
			IsAtomic: isAtomic,
		},
		{
			SharedID: "fifth transaction",
			MsgID:    common.HexToHash("0x77777777"),
			State:    state,
			Outcome:  outcome,
			IsAtomic: isAtomic,
		},
	}

	for _, tx := range seedTransactions {
		err := repo.Create(context.TODO(), tx)
		require.NoError(t, err)
	}

	t.Run("limits transactions to one", func(t *testing.T) {
		limit := 1

		got, err := repo.GetByStateOutcomeAndAtomicity(context.TODO(), state, outcome, isAtomic, repository.WithLimit(limit))
		require.Nil(t, err)

		require.Equal(t, limit, len(got))
	})

	t.Run("limits transactions to two", func(t *testing.T) {
		limit := 2

		got, err := repo.GetByStateOutcomeAndAtomicity(context.TODO(), state, outcome, isAtomic, repository.WithLimit(limit))
		require.Nil(t, err)

		require.Equal(t, limit, len(got))
	})

	// PostgreSQL LIMIT 0 returns no rows, but we maintain backward compatibility
	// by returning all rows when limit is 0
	t.Run("returns all transactions when limit is set to zero", func(t *testing.T) {
		limit := 0

		got, err := repo.GetByStateOutcomeAndAtomicity(context.TODO(), state, outcome, isAtomic, repository.WithLimit(limit))
		require.Nil(t, err)

		require.Equal(t, len(seedTransactions), len(got))
	})

	t.Run("returns all transactions when less elements than limit", func(t *testing.T) {
		limit := 10

		got, err := repo.GetByStateOutcomeAndAtomicity(context.TODO(), state, outcome, isAtomic, repository.WithLimit(limit))
		require.Nil(t, err)

		require.Equal(t, len(seedTransactions), len(got))
	})
}

func TestTransactionRepository_GetBySharedIDs(t *testing.T) {
	testtools.SilenceLogger()

	pool, teardown := testtools.SetupPostgres(t, dbConf)
	defer teardown()

	repo := repository.NewTransactionRepository(pool)

	t.Run("returns an empty slice on no transactions in database", func(t *testing.T) {
		var want []types.Transaction

		sharedIDs := []string{"example-shared-id"}

		got, err := repo.GetBySharedIDs(context.TODO(), sharedIDs)
		require.Nil(t, err)

		require.Equal(t, want, got)
	})

	// Seed using repository Create method
	err := repo.Create(context.TODO(), testdata.Transaction1)
	require.NoError(t, err)
	err = repo.Create(context.TODO(), testdata.Transaction2)
	require.NoError(t, err)
	err = repo.Create(context.TODO(), testdata.Transaction3)
	require.NoError(t, err)

	t.Run("returns an empty slice on no transactions with matching shared id", func(t *testing.T) {
		var want []types.Transaction

		sharedIDs := []string{"example-shared-id"}

		got, err := repo.GetBySharedIDs(context.TODO(), sharedIDs)
		require.Nil(t, err)

		require.Equal(t, want, got)
	})

	t.Run("returns transaction with matching shared id", func(t *testing.T) {
		want := []types.Transaction{testdata.Transaction1}

		sharedIDs := []string{testdata.Transaction1.SharedID}

		got, err := repo.GetBySharedIDs(context.TODO(), sharedIDs)
		require.Nil(t, err)

		require.Equal(t, want, got)
	})

	t.Run("returns multiple transactions with matching shared ids", func(t *testing.T) {
		want := []types.Transaction{testdata.Transaction1, testdata.Transaction2}

		sharedIDs := []string{testdata.Transaction1.SharedID, testdata.Transaction2.SharedID}

		got, err := repo.GetBySharedIDs(context.TODO(), sharedIDs)
		require.Nil(t, err)

		require.ElementsMatch(t, want, got)
	})
}

func TestTransactionRepository_UpdateDestinationHashForSharedID(t *testing.T) {
	testtools.SilenceLogger()

	pool, teardown := testtools.SetupPostgres(t, dbConf)
	defer teardown()

	repo := repository.NewTransactionRepository(pool)

	t.Run("returns on no transactions to update", func(t *testing.T) {
		sharedID := "example-shared-id"
		destinationHash := common.HexToHash("0xDEADBEEF")

		err := repo.UpdateDestinationHashForSharedID(context.TODO(), sharedID, destinationHash)
		require.Nil(t, err)
	})

	// Seed using repository Create method
	err := repo.Create(context.TODO(), testdata.Transaction1)
	require.NoError(t, err)
	err = repo.Create(context.TODO(), testdata.Transaction2)
	require.NoError(t, err)
	err = repo.Create(context.TODO(), testdata.Transaction3)
	require.NoError(t, err)

	t.Run("updates destination hash for shared ID", func(t *testing.T) {
		sharedID := testdata.Transaction1.SharedID
		want := common.HexToHash("0xDEADBEEF")

		err := repo.UpdateDestinationHashForSharedID(context.TODO(), sharedID, want)
		require.Nil(t, err)

		gotTxSlice, err := repo.GetBySharedIDs(context.TODO(), []string{sharedID})
		require.Nil(t, err)

		require.Len(t, gotTxSlice, 1)

		got := gotTxSlice[0].TxHashDestination
		require.Equal(t, want, got)
	})
}

func TestTransactionRepository_BatchSetStateAndOutcome(t *testing.T) {
	testtools.SilenceLogger()

	pool, teardown := testtools.SetupPostgres(t, dbConf)
	defer teardown()

	repo := repository.NewTransactionRepository(pool)

	t.Run("returns on no transactions to update", func(t *testing.T) {
		sharedIDs := []string{"example-shared-id"}

		err := repo.BatchSetStateAndOutcome(context.TODO(), sharedIDs, types.SourceFinalized, types.OutcomeSuccess)
		require.Nil(t, err)
	})

	// Seed using repository Create method
	err := repo.Create(context.TODO(), testdata.Transaction1)
	require.NoError(t, err)
	err = repo.Create(context.TODO(), testdata.Transaction2)
	require.NoError(t, err)
	err = repo.Create(context.TODO(), testdata.Transaction3)
	require.NoError(t, err)

	t.Run("updates transaction state and outcome for shared ids", func(t *testing.T) {
		transactions := []types.Transaction{testdata.Transaction1, testdata.Transaction2}

		wantState := types.SourcePublish
		wantOutcome := types.OutcomeSuccess

		want := make([]types.Transaction, len(transactions))
		copy(want, transactions)
		want[0].State = wantState
		want[0].Outcome = wantOutcome
		want[1].State = wantState
		want[1].Outcome = wantOutcome

		sharedIDs := []string{transactions[0].SharedID, transactions[1].SharedID}

		err := repo.BatchSetStateAndOutcome(context.TODO(), sharedIDs, wantState, wantOutcome)
		require.Nil(t, err)

		got, err := repo.GetBySharedIDs(context.TODO(), sharedIDs)
		require.Nil(t, err)

		require.ElementsMatch(t, want, got)
	})
}

func TestTransactionRepository_BatchCreateWithStateAndOutcome(t *testing.T) {
	testtools.SilenceLogger()

	pool, teardown := testtools.SetupPostgres(t, dbConf)
	defer teardown()

	repo := repository.NewTransactionRepository(pool)

	t.Run("batch creates transactions and sets state and outcome", func(t *testing.T) {
		transactions := []types.Transaction{testdata.Transaction1, testdata.Transaction2}

		wantState := types.SourcePublish
		wantOutcome := types.OutcomeSuccess

		want := make([]types.Transaction, len(transactions))
		copy(want, transactions)
		want[0].State = wantState
		want[0].Outcome = wantOutcome
		want[1].State = wantState
		want[1].Outcome = wantOutcome

		err := repo.BatchCreateWithStateAndOutcome(context.TODO(), transactions, wantState, wantOutcome)
		require.Nil(t, err)

		got, err := repo.GetByStateOutcomeAndAtomicity(context.Background(), wantState, wantOutcome, true)
		require.Nil(t, err)

		require.NotEmpty(t, got)
	})
}

func TestTransactionRepository_Create(t *testing.T) {
	testtools.SilenceLogger()

	pool, teardown := testtools.SetupPostgres(t, dbConf)
	defer teardown()

	repo := repository.NewTransactionRepository(pool)

	t.Run("persists transaction in database", func(t *testing.T) {
		want := testdata.Transaction1

		err := repo.Create(context.TODO(), want)
		require.Nil(t, err)

		sharedIDs := []string{want.SharedID}
		gotSlice, err := repo.GetBySharedIDs(context.Background(), sharedIDs)
		require.Nil(t, err)

		require.Len(t, gotSlice, 1)
		require.Equal(t, want, gotSlice[0])
	})
}
