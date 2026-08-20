package contractclient_test

import (
	"context"
	"encoding/binary"
	"errors"
	"io"
	"log/slog"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/contractclient"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/contracts/RaylsAccessManagerV1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// hasRoleResponse encodes a single boolean as the ABI-encoded payload returned
// by RaylsAccessManagerV1.hasRoleByName: a 32-byte word, 1 if true, 0 otherwise.
func hasRoleResponse(hasRole bool) []byte {
	out := make([]byte, 32)
	if hasRole {
		out[31] = 1
	}
	return out
}

func TestAccessManager_WaitForAuthorization(t *testing.T) {
	// Silence slog output from the polling loop.
	slog.SetDefault(slog.New(slog.NewTextHandler(io.Discard, nil)))

	t.Run("returns error for invalid private key hex strings", func(t *testing.T) {
		wantErrorType := &contractclient.AccessManagerError{}

		address := common.HexToAddress("0x1")
		executor := &stubExecutor{}

		cli := contractclient.NewAccessManager(address, executor)

		err := cli.WaitForAuthorization(context.Background(), "test", []string{"not-a-valid-hex-key"})

		require.ErrorAs(t, err, &wantErrorType, "didn't wrap error in AccessManagerError")
		assert.Contains(t, err.Error(), "failed to convert private keys to addresses")
	})

	t.Run("returns context error on context cancellation before authorization", func(t *testing.T) {
		address := common.HexToAddress("0x1")
		// Always returns "no role yet" so the loop has to wait.
		executor := &stubExecutor{callResult: hasRoleResponse(false)}

		cli := contractclient.NewAccessManager(address, executor)
		cli.PollInterval = 10 * time.Millisecond

		testPrivateKey, err := crypto.GenerateKey()
		require.NoError(t, err)
		testKeyHex := common.Bytes2Hex(crypto.FromECDSA(testPrivateKey))

		// Use an already-cancelled context.
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		err = cli.WaitForAuthorization(ctx, "test", []string{testKeyHex})

		require.ErrorIs(t, err, context.Canceled)
	})

	t.Run("wraps executor Call errors in AccessManagerError", func(t *testing.T) {
		wantError := errors.New("call error")
		wantErrorType := &contractclient.AccessManagerError{}

		address := common.HexToAddress("0x1")
		executor := &stubExecutor{callErr: wantError}

		cli := contractclient.NewAccessManager(address, executor)
		cli.PollInterval = 1 * time.Millisecond

		testPrivateKey, err := crypto.GenerateKey()
		require.NoError(t, err)
		testKeyHex := common.Bytes2Hex(crypto.FromECDSA(testPrivateKey))

		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		defer cancel()

		err = cli.WaitForAuthorization(ctx, "test", []string{testKeyHex})

		require.ErrorAs(t, err, &wantErrorType, "didn't wrap error in AccessManagerError")
		require.ErrorIs(t, err, wantError, "didn't wrap underlying error")
	})

	t.Run("returns nil when all keys have RELAYER role", func(t *testing.T) {
		address := common.HexToAddress("0x1")
		executor := &stubExecutor{callResult: hasRoleResponse(true)}

		cli := contractclient.NewAccessManager(address, executor)
		cli.PollInterval = 1 * time.Millisecond

		testPrivateKey, err := crypto.GenerateKey()
		require.NoError(t, err)
		testKeyHex := common.Bytes2Hex(crypto.FromECDSA(testPrivateKey))
		testAddress := crypto.PubkeyToAddress(testPrivateKey.PublicKey)

		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		defer cancel()

		err = cli.WaitForAuthorization(ctx, "test", []string{testKeyHex})

		require.NoError(t, err)
		assert.Equal(t, address, executor.spyCallAddress, "should call against the registry address")

		// Verify the calldata matches PackHasRoleByName(RoleRelayer, testAddress).
		binding := RaylsAccessManagerV1.NewRaylsAccessManagerV1()
		wantCalldata := binding.PackHasRoleByName(contractclient.RoleRelayer, testAddress)
		assert.Equal(t, wantCalldata, executor.spyCallCalldata, "should query has-role for derived address")
	})
}

func TestAccessManager_GrantAuthorizedSenderRole(t *testing.T) {
	slog.SetDefault(slog.New(slog.NewTextHandler(io.Discard, nil)))

	binding := RaylsAccessManagerV1.NewRaylsAccessManagerV1()

	t.Run("successfully grants AUTHORIZED_SENDER", func(t *testing.T) {
		const wantRoleID uint64 = 42
		wantTokenAddress := common.HexToAddress("0xdeadbeef")
		address := common.HexToAddress("0xABCD")

		executor := &stubExecutor{callResult: encodeUint64Word(wantRoleID)}

		cli := contractclient.NewAccessManager(address, executor)

		err := cli.GrantAuthorizedSenderRole(context.Background(), wantTokenAddress)

		require.NoError(t, err)

		assert.Equal(t, address, executor.spyCallAddress)
		assert.Equal(t, binding.PackGetRoleIdByName(contractclient.RoleAuthorizedSender), executor.spyCallCalldata)

		assert.Equal(t, address, executor.spyExecuteAddress)
		assert.Equal(t, binding.PackGrantRole(wantRoleID, wantTokenAddress, 0), executor.spyExecuteCalldata)
	})

	t.Run("wraps role lookup error", func(t *testing.T) {
		wantErr := errors.New("call failed")
		executor := &stubExecutor{callErr: wantErr}
		cli := contractclient.NewAccessManager(common.Address{}, executor)

		err := cli.GrantAuthorizedSenderRole(context.Background(), common.Address{})

		require.ErrorIs(t, err, wantErr)
		assert.Contains(t, err.Error(), "failed to get AUTHORIZED_SENDER id")
	})

	t.Run("returns error on undecodable role ID payload", func(t *testing.T) {
		executor := &stubExecutor{callResult: []byte{0x01}}
		cli := contractclient.NewAccessManager(common.Address{}, executor)

		err := cli.GrantAuthorizedSenderRole(context.Background(), common.Address{})

		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to unpack role ID")
	})

	t.Run("wraps grant error", func(t *testing.T) {
		wantErr := errors.New("execute failed")
		executor := &stubExecutor{
			callResult: encodeUint64Word(1),
			executeErr: wantErr,
		}
		cli := contractclient.NewAccessManager(common.Address{}, executor)

		err := cli.GrantAuthorizedSenderRole(context.Background(), common.Address{})

		require.ErrorIs(t, err, wantErr)
		assert.Contains(t, err.Error(), "failed to grant AUTHORIZED_SENDER")
	})
}

// encodeUint64Word produces the ABI-encoded form of a uint64 return value:
// a single 32-byte word, big-endian, right-aligned.
func encodeUint64Word(v uint64) []byte {
	out := make([]byte, 32)
	binary.BigEndian.PutUint64(out[24:], v)
	return out
}
