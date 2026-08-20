package contractclient_test

import (
	"context"
	"errors"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/contractclient"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/contracts/AuditManagerV1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParticipantStorageClient_GetVenOperatorChainInfo(t *testing.T) {
	t.Run("wraps executor call errors in ParticipantStorageClientError", func(t *testing.T) {
		wantError := errors.New("call error")

		address := common.HexToAddress("0x1")
		executor := &stubExecutor{callErr: wantError}
		client := contractclient.NewParticipantStorageClient(address, big.NewInt(1), big.NewInt(2), executor)

		_, err := client.GetVenOperatorChainInfo(context.Background(), big.NewInt(100))

		require.Error(t, err)
		var clientErr *contractclient.ParticipantStorageClientError
		require.True(t, errors.As(err, &clientErr))
	})
}

func TestParticipantStorageClient_GetMyChainInfo(t *testing.T) {
	t.Run("wraps executor call errors in ParticipantStorageClientError", func(t *testing.T) {
		wantError := errors.New("call error")

		address := common.HexToAddress("0x1")
		executor := &stubExecutor{callErr: wantError}
		client := contractclient.NewParticipantStorageClient(address, big.NewInt(1), big.NewInt(2), executor)

		_, err := client.GetMyChainInfo(context.Background(), big.NewInt(100))

		require.Error(t, err)
		var clientErr *contractclient.ParticipantStorageClientError
		require.True(t, errors.As(err, &clientErr))
	})
}

func TestParticipantStorageClient_GetChainViewData(t *testing.T) {
	t.Run("wraps executor call errors in ParticipantStorageClientError", func(t *testing.T) {
		wantError := errors.New("call error")

		address := common.HexToAddress("0x1")
		executor := &stubExecutor{callErr: wantError}
		client := contractclient.NewParticipantStorageClient(address, big.NewInt(1), big.NewInt(2), executor)

		_, err := client.GetChainViewData(context.Background(), big.NewInt(100), big.NewInt(100))

		require.Error(t, err)
		var clientErr *contractclient.ParticipantStorageClientError
		require.True(t, errors.As(err, &clientErr))
	})
}

func TestParticipantStorageClient_GetEnygmaParticipants(t *testing.T) {
	t.Run("wraps executor call errors in ParticipantStorageClientError", func(t *testing.T) {
		wantError := errors.New("call error")

		address := common.HexToAddress("0x1")
		executor := &stubExecutor{callErr: wantError}
		client := contractclient.NewParticipantStorageClient(address, big.NewInt(1), big.NewInt(2), executor)

		_, err := client.GetEnygmaParticipants(context.Background())

		require.Error(t, err)
		var clientErr *contractclient.ParticipantStorageClientError
		require.True(t, errors.As(err, &clientErr))
	})
}

func TestParticipantStorageClient_SetChainViewData(t *testing.T) {
	t.Run("successfully sets chain view data via executor", func(t *testing.T) {
		address := common.HexToAddress("0x1")
		executor := &stubExecutor{}
		client := contractclient.NewParticipantStorageClient(address, big.NewInt(1), big.NewInt(2), executor)

		err := client.SetChainViewData(context.Background(), big.NewInt(100), "test-view-key", big.NewInt(42))

		require.NoError(t, err)
		assert.Equal(t, address, executor.spyExecuteAddress)
		assert.NotNil(t, executor.spyExecuteCalldata)
	})

	t.Run("wraps executor errors in ParticipantStorageClientError", func(t *testing.T) {
		wantError := errors.New("execute failed")

		address := common.HexToAddress("0x1")
		executor := &stubExecutor{executeErr: wantError}
		client := contractclient.NewParticipantStorageClient(address, big.NewInt(1), big.NewInt(2), executor)

		err := client.SetChainViewData(context.Background(), big.NewInt(100), "test-view-key", big.NewInt(42))

		require.Error(t, err)
		var clientErr *contractclient.ParticipantStorageClientError
		require.True(t, errors.As(err, &clientErr))
		require.ErrorIs(t, err, wantError)
	})
}

func TestParticipantStorageClient_SetAuditInfo(t *testing.T) {
	t.Run("successfully sets audit info via executor", func(t *testing.T) {
		address := common.HexToAddress("0x1")
		executor := &stubExecutor{}
		client := contractclient.NewParticipantStorageClient(address, big.NewInt(1), big.NewInt(2), executor)

		err := client.SetAuditInfo(context.Background(), big.NewInt(100), big.NewInt(42), "audit-key", []byte("encr"), []byte("mac"))

		require.NoError(t, err)
		assert.Equal(t, address, executor.spyExecuteAddress)
	})

	t.Run("wraps executor errors in ParticipantStorageClientError", func(t *testing.T) {
		wantError := errors.New("execute failed")

		address := common.HexToAddress("0x1")
		executor := &stubExecutor{executeErr: wantError}
		client := contractclient.NewParticipantStorageClient(address, big.NewInt(1), big.NewInt(2), executor)

		err := client.SetAuditInfo(context.Background(), big.NewInt(100), big.NewInt(42), "audit-key", []byte("encr"), []byte("mac"))

		require.Error(t, err)
		var clientErr *contractclient.ParticipantStorageClientError
		require.True(t, errors.As(err, &clientErr))
		require.ErrorIs(t, err, wantError)
	})
}

func TestParticipantStorageClient_SetPaymentSpendPublicKey(t *testing.T) {
	t.Run("successfully sets payment spend public key via executor", func(t *testing.T) {
		address := common.HexToAddress("0x1")
		executor := &stubExecutor{}
		client := contractclient.NewParticipantStorageClient(address, big.NewInt(1), big.NewInt(2), executor)

		err := client.SetPaymentSpendPublicKey(context.Background(), big.NewInt(100), big.NewInt(12345), []common.Address{})

		require.NoError(t, err)
		assert.Equal(t, address, executor.spyExecuteAddress)
	})

	t.Run("wraps executor errors in ParticipantStorageClientError", func(t *testing.T) {
		wantError := errors.New("execute failed")

		address := common.HexToAddress("0x1")
		executor := &stubExecutor{executeErr: wantError}
		client := contractclient.NewParticipantStorageClient(address, big.NewInt(1), big.NewInt(2), executor)

		err := client.SetPaymentSpendPublicKey(context.Background(), big.NewInt(100), big.NewInt(12345), []common.Address{})

		require.Error(t, err)
		var clientErr *contractclient.ParticipantStorageClientError
		require.True(t, errors.As(err, &clientErr))
		require.ErrorIs(t, err, wantError)
	})
}

func TestParticipantStorageClient_InitiateKeyAgreement(t *testing.T) {
	t.Run("successfully initiates key agreement via executor", func(t *testing.T) {
		address := common.HexToAddress("0x1")
		executor := &stubExecutor{}
		client := contractclient.NewParticipantStorageClient(address, big.NewInt(1), big.NewInt(2), executor)

		err := client.InitiateKeyAgreement(context.Background(), big.NewInt(100), []byte("ciphertext"), []byte("digest"), big.NewInt(200))

		require.NoError(t, err)
		assert.Equal(t, address, executor.spyExecuteAddress)
	})

	t.Run("wraps executor errors in ParticipantStorageClientError", func(t *testing.T) {
		wantError := errors.New("execute failed")

		address := common.HexToAddress("0x1")
		executor := &stubExecutor{executeErr: wantError}
		client := contractclient.NewParticipantStorageClient(address, big.NewInt(1), big.NewInt(2), executor)

		err := client.InitiateKeyAgreement(context.Background(), big.NewInt(100), []byte("ciphertext"), []byte("digest"), big.NewInt(200))

		require.Error(t, err)
		var clientErr *contractclient.ParticipantStorageClientError
		require.True(t, errors.As(err, &clientErr))
		require.ErrorIs(t, err, wantError)
	})

	t.Run("returns ErrOutdatedKeyAgreement when revert matches BlockNumberLowerThanLatestKeyAgreement selector", func(t *testing.T) {
		selector := AuditManagerV1.AuditManagerV1AuditManagerV1BlockNumberLowerThanLatestKeyAgreementErrorID().Bytes()[:4]
		executor := &stubExecutor{executeErr: contractclient.NewErrorWithRevertData(selector)}

		client := contractclient.NewParticipantStorageClient(common.HexToAddress("0x1"), big.NewInt(1), big.NewInt(2), executor)

		err := client.InitiateKeyAgreement(context.Background(), big.NewInt(100), []byte("ciphertext"), []byte("digest"), big.NewInt(200))
		require.Error(t, err)
		require.ErrorIs(t, err, contractclient.ErrOutdatedKeyAgreement)
	})

	t.Run("does not return ErrOutdatedKeyAgreement for a different revert selector", func(t *testing.T) {
		// Use the unauthorized-caller selector to prove selector compare isn't a false positive.
		otherSelector := AuditManagerV1.AuditManagerV1AuditManagerV1UnauthorizedCallerErrorID().Bytes()[:4]
		executor := &stubExecutor{executeErr: contractclient.NewErrorWithRevertData(otherSelector)}

		client := contractclient.NewParticipantStorageClient(common.HexToAddress("0x1"), big.NewInt(1), big.NewInt(2), executor)

		err := client.InitiateKeyAgreement(context.Background(), big.NewInt(100), []byte("ciphertext"), []byte("digest"), big.NewInt(200))
		require.Error(t, err)
		assert.NotErrorIs(t, err, contractclient.ErrOutdatedKeyAgreement)
	})
}
