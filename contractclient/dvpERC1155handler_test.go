package contractclient_test

import (
	"context"
	"errors"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/raylsnetwork/rayls-sovereign-relayer/contractclient"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDvpERC1155HandlerClient_GetAllTokenIdsWithSupply(t *testing.T) {
	t.Run("wraps executor call errors in DvpERC1155HandlerClientError", func(t *testing.T) {
		wantError := errors.New("call error")

		address := common.HexToAddress("0x1")
		executor := &stubExecutor{callErr: wantError}

		client := contractclient.NewDvpERC1155HandlerClient(executor)

		_, err := client.GetAllTokenIdsWithSupply(context.Background(), address)

		require.Error(t, err)
		var clientErr *contractclient.DvpERC1155HandlerClientError
		require.True(t, errors.As(err, &clientErr))
		require.Contains(t, err.Error(), "failed to get all token ids with supply")
	})
}

func TestDvpERC1155HandlerClient_GetTokenExtraData(t *testing.T) {
	t.Run("wraps executor call errors in DvpERC1155HandlerClientError", func(t *testing.T) {
		wantError := errors.New("call error")

		address := common.HexToAddress("0x1")
		executor := &stubExecutor{callErr: wantError}

		client := contractclient.NewDvpERC1155HandlerClient(executor)

		_, err := client.GetTokenExtraData(context.Background(), address, big.NewInt(12345))

		require.Error(t, err)
		var clientErr *contractclient.DvpERC1155HandlerClientError
		require.True(t, errors.As(err, &clientErr))
		require.Contains(t, err.Error(), "failed to get token extra data")
	})
}

func TestDvpERC1155HandlerClient_Unlock(t *testing.T) {
	t.Run("successfully unlocks ERC1155 tokens via executor", func(t *testing.T) {
		address := common.HexToAddress("0x1")
		executor := &stubExecutor{}

		client := contractclient.NewDvpERC1155HandlerClient(executor)

		tokenId := big.NewInt(12345)
		tokenAmount := big.NewInt(100)
		to := common.HexToAddress("0x1234567890123456789012345678901234567890")
		err := client.Unlock(context.Background(), "test-event-id", address, tokenId, tokenAmount, to)

		require.NoError(t, err)
		assert.Equal(t, address, executor.spyExecuteAddress)
		assert.NotNil(t, executor.spyExecuteCalldata)
	})

	t.Run("wraps executor errors in DvpERC1155HandlerClientError", func(t *testing.T) {
		wantError := errors.New("execute failed")

		address := common.HexToAddress("0x1")
		executor := &stubExecutor{executeErr: wantError}

		client := contractclient.NewDvpERC1155HandlerClient(executor)

		err := client.Unlock(context.Background(), "test-event-id", address, big.NewInt(12345), big.NewInt(100), common.Address{})

		require.Error(t, err)
		var clientErr *contractclient.DvpERC1155HandlerClientError
		require.True(t, errors.As(err, &clientErr))
		require.Contains(t, err.Error(), "failed to unlock ERC1155 from Enygma DvP")
	})
}

func TestDvpERC1155HandlerClient_MarkSwapCompleted(t *testing.T) {
	t.Run("wraps executor errors in DvpERC1155HandlerClientError", func(t *testing.T) {
		wantError := errors.New("execute failed")

		address := common.HexToAddress("0x1")
		executor := &stubExecutor{executeErr: wantError}

		client := contractclient.NewDvpERC1155HandlerClient(executor)

		err := client.MarkSwapCompleted(
			context.Background(),
			address,
			common.Address{},
			common.Address{},
			big.NewInt(100),
			big.NewInt(200),
			"resource-id",
			"1",
			big.NewInt(100),
			[]byte("data"),
			"abcd1234567890abcd1234567890abcd1234567890abcd1234567890abcd1234",
		)

		require.Error(t, err)
		var clientErr *contractclient.DvpERC1155HandlerClientError
		require.True(t, errors.As(err, &clientErr))
	})
}
