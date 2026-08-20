package contractclient_test

import (
	"context"
	"errors"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/raylsnetwork/rayls-sovereign-relayer/contractclient"
	"github.com/raylsnetwork/rayls-sovereign-relayer/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDvpERC721HandlerClient_GetTotalSupply(t *testing.T) {
	t.Run("wraps executor call errors in DvpERC721HandlerClientError", func(t *testing.T) {
		wantError := errors.New("call error")

		address := common.HexToAddress("0x1")
		executor := &stubExecutor{callErr: wantError}

		client := contractclient.NewDvpERC721HandlerClient(executor)

		_, err := client.GetTotalSupply(context.Background(), address)

		require.Error(t, err)
		var clientErr *contractclient.DvpERC721HandlerClientError
		require.True(t, errors.As(err, &clientErr))
		require.Contains(t, err.Error(), "failed to get total supply")
	})
}

func TestDvpERC721HandlerClient_GetExtraData(t *testing.T) {
	t.Run("wraps executor call errors in DvpERC721HandlerClientError", func(t *testing.T) {
		wantError := errors.New("call error")

		address := common.HexToAddress("0x1")
		executor := &stubExecutor{callErr: wantError}

		client := contractclient.NewDvpERC721HandlerClient(executor)

		_, err := client.GetExtraData(context.Background(), address, big.NewInt(12345))

		require.Error(t, err)
		var clientErr *contractclient.DvpERC721HandlerClientError
		require.True(t, errors.As(err, &clientErr))
		require.Contains(t, err.Error(), "failed to get extra data")
	})
}

func TestDvpERC721HandlerClient_Unlock(t *testing.T) {
	t.Run("successfully unlocks ERC721 via executor", func(t *testing.T) {
		address := common.HexToAddress("0x1")
		executor := &stubExecutor{}

		client := contractclient.NewDvpERC721HandlerClient(executor)

		err := client.Unlock(context.Background(), "test-event-id", address, big.NewInt(12345))

		require.NoError(t, err)
		assert.Equal(t, address, executor.spyExecuteAddress)
		assert.NotNil(t, executor.spyExecuteCalldata)
	})

	t.Run("wraps executor errors in DvpERC721HandlerClientError", func(t *testing.T) {
		wantError := errors.New("execute failed")

		address := common.HexToAddress("0x1")
		executor := &stubExecutor{executeErr: wantError}

		client := contractclient.NewDvpERC721HandlerClient(executor)

		err := client.Unlock(context.Background(), "test-event-id", address, big.NewInt(12345))

		require.Error(t, err)
		var clientErr *contractclient.DvpERC721HandlerClientError
		require.True(t, errors.As(err, &clientErr))
		require.Contains(t, err.Error(), "failed to unlock ERC721 from Enygma DvP")
	})
}

func TestDvpERC721HandlerClient_MarkSwapCompleted(t *testing.T) {
	t.Run("wraps executor errors in DvpERC721HandlerClientError", func(t *testing.T) {
		wantError := errors.New("execute failed")

		address := common.HexToAddress("0x1")
		executor := &stubExecutor{executeErr: wantError}

		client := contractclient.NewDvpERC721HandlerClient(executor)

		err := client.MarkSwapCompleted(
			context.Background(),
			address,
			common.Address{},
			big.NewInt(100),
			big.NewInt(200),
			"resource-id",
			"1",
			big.NewInt(100),
			"abcd1234567890abcd1234567890abcd1234567890abcd1234567890abcd1234",
		)

		require.Error(t, err)
		var clientErr *contractclient.DvpERC721HandlerClientError
		require.True(t, errors.As(err, &clientErr))
	})
}

func TestDvpERC721HandlerClient_NotifySenderWithPNCommunicator(t *testing.T) {
	t.Run("successfully notifies sender via executor", func(t *testing.T) {
		address := common.HexToAddress("0x1")
		executor := &stubExecutor{}

		client := contractclient.NewDvpERC721HandlerClient(executor)

		err := client.NotifySenderWithPNCommunicator(
			context.Background(),
			address,
			"abcd1234567890abcd1234567890abcd1234567890abcd1234567890abcd1234",
			types.SwapError,
			"test message",
		)

		require.NoError(t, err)
		assert.Equal(t, address, executor.spyExecuteAddress)
	})

	t.Run("wraps executor errors in DvpERC721HandlerClientError", func(t *testing.T) {
		wantError := errors.New("execute failed")

		address := common.HexToAddress("0x1")
		executor := &stubExecutor{executeErr: wantError}

		client := contractclient.NewDvpERC721HandlerClient(executor)

		err := client.NotifySenderWithPNCommunicator(
			context.Background(),
			address,
			"abcd1234567890abcd1234567890abcd1234567890abcd1234567890abcd1234",
			types.SwapError,
			"test message",
		)

		require.Error(t, err)
		var clientErr *contractclient.DvpERC721HandlerClientError
		require.True(t, errors.As(err, &clientErr))
	})
}
