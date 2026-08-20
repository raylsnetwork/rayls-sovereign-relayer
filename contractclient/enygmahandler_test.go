package contractclient_test

import (
	"context"
	"errors"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/contractclient"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// progSteps builds a per-recipient programmability step array from raw arg payloads — one step
// per payload, with placeholder resourceId/selector (these tests assert dispatch wiring, not
// step structure).
func progSteps(payloads ...[]byte) []types.EnygmaProgramData {
	steps := make([]types.EnygmaProgramData, len(payloads))
	for i, p := range payloads {
		steps[i] = types.EnygmaProgramData{ResourceId: [32]byte{byte(i + 1)}, Args: p}
	}
	return steps
}

func TestEnygmaHandlerClient_ReceiveWithdraw(t *testing.T) {
	t.Run("successfully receives withdrawal from dvp", func(t *testing.T) {
		wantTo := common.HexToAddress("0x1111111111111111111111111111111111111111")
		wantValue := big.NewInt(100)
		wantReferenceId := [32]byte{0xAA, 0xBB, 0xCC, 0xDD}

		address := common.HexToAddress("0x1")
		executor := &stubExecutor{}

		client := contractclient.NewEnygmaHandlerClient(executor, common.HexToAddress("0x000000000000000000000000000000000000E0EC"))

		err := client.ReceiveWithdraw(context.Background(), "test-event-id", address, wantTo, wantValue, wantReferenceId)

		require.Nil(t, err)
		assert.Equal(t, address, executor.spyExecuteAddress)
		assert.NotNil(t, executor.spyExecuteCalldata)
	})

	t.Run("wraps executor errors in EnygmaHandlerClientError", func(t *testing.T) {
		wantError := errors.New("execute failed")

		address := common.HexToAddress("0x1")
		executor := &stubExecutor{
			executeErr: wantError,
		}

		client := contractclient.NewEnygmaHandlerClient(executor, common.HexToAddress("0x000000000000000000000000000000000000E0EC"))

		err := client.ReceiveWithdraw(context.Background(), "test-event-id", address, common.Address{}, big.NewInt(0), [32]byte{})

		require.NotNil(t, err)
		var wrappedErr *contractclient.EnygmaHandlerClientError
		require.ErrorAs(t, err, &wrappedErr, "error should be wrapped")
		require.ErrorIs(t, err, wantError, "underlying error should be preserved")
	})
}

func TestEnygmaHandlerClient_ReceiveDestTransferBatch(t *testing.T) {
	executorAddr := common.HexToAddress("0x000000000000000000000000000000000000E0EC")

	newTransfer := func(msgID string, amount *big.Int, from common.Address, programData []types.EnygmaProgramData) *types.EnygmaCrossTransferData {
		return &types.EnygmaCrossTransferData{
			EnygmaTransferBatchTx: &types.EnygmaTransferBatchTx{
				MessageId:   msgID,
				ToAmount:    amount,
				FromAddress: from,
				ProgramData: programData,
			},
			EnygmaAddress: common.HexToAddress("0xToKeN"),
		}
	}

	t.Run("dispatches each transfer to the ProgrammabilityExecutor, not the token", func(t *testing.T) {
		// ProgramData blobs here are illustrative raw bytes — they would not ABI-decode as
		// valid (bytes32, bytes4, bytes) blobs on-chain. This test only verifies the dispatch
		// wiring (target address + non-empty packed calldata), not blob structure.
		transfers := []*types.EnygmaCrossTransferData{
			newTransfer("msg-1", big.NewInt(100), common.HexToAddress("0xa11ce"), progSteps([]byte{0x01, 0x02})),
			newTransfer("msg-2", big.NewInt(200), common.HexToAddress("0xb0b"), progSteps([]byte{0x03}, []byte{0x04})),
		}

		executor := &stubExecutor{batchResults: map[string]contractclient.BatchResult{}}
		client := contractclient.NewEnygmaHandlerClient(executor, executorAddr)

		_, err := client.ReceiveDestTransferBatch(context.Background(), transfers)

		require.NoError(t, err)
		require.Len(t, executor.spyBatchItems, 2)

		byMsgID := make(map[string]contractclient.BatchInput, len(executor.spyBatchItems))
		for _, item := range executor.spyBatchItems {
			byMsgID[item.MsgID] = item
		}
		for _, msgID := range []string{"msg-1", "msg-2"} {
			item, ok := byMsgID[msgID]
			require.True(t, ok, "expected an item for %s", msgID)
			// Every dispatch targets the executor; the token address is never the tx target.
			assert.Equal(t, executorAddr, item.Address)
			assert.NotEmpty(t, item.Data, "packed executeProgramData calldata must not be empty")
		}
	})

	t.Run("wraps executor errors in EnygmaHandlerClientError", func(t *testing.T) {
		wantError := errors.New("batch failed")
		executor := &stubExecutor{batchExecuteErr: wantError}
		client := contractclient.NewEnygmaHandlerClient(executor, executorAddr)

		_, err := client.ReceiveDestTransferBatch(
			context.Background(),
			[]*types.EnygmaCrossTransferData{newTransfer("msg-1", big.NewInt(1), common.Address{}, progSteps([]byte{0x01}))},
		)

		require.NotNil(t, err)
		var wrappedErr *contractclient.EnygmaHandlerClientError
		require.ErrorAs(t, err, &wrappedErr, "error should be wrapped")
		require.ErrorIs(t, err, wantError, "underlying error should be preserved")
	})
}

func TestEnygmaHandlerClient_MarkSwapCompleted(t *testing.T) {
	t.Run("successfully marks dvp swap as completed", func(t *testing.T) {
		wantDestChainId := big.NewInt(100)
		wantSharedId := "abcd1234567890abcd1234567890abcd1234567890abcd1234567890abcd1234"

		address := common.HexToAddress("0x1")
		executor := &stubExecutor{}

		client := contractclient.NewEnygmaHandlerClient(executor, common.HexToAddress("0x000000000000000000000000000000000000E0EC"))

		err := client.MarkSwapCompleted(context.Background(), address, wantDestChainId, wantSharedId)

		require.Nil(t, err)
		assert.Equal(t, address, executor.spyExecuteAddress)
		assert.NotNil(t, executor.spyExecuteCalldata)
	})

	t.Run("wraps string to bytes32 conversion errors in EnygmaHandlerClientError", func(t *testing.T) {
		address := common.HexToAddress("0x1")
		executor := &stubExecutor{}

		client := contractclient.NewEnygmaHandlerClient(executor, common.HexToAddress("0x000000000000000000000000000000000000E0EC"))

		err := client.MarkSwapCompleted(context.Background(), address, big.NewInt(100), "invalid-hex-string")

		require.NotNil(t, err)
		var wrappedErr *contractclient.EnygmaHandlerClientError
		require.ErrorAs(t, err, &wrappedErr, "error should be wrapped")
	})

	t.Run("wraps executor errors in EnygmaHandlerClientError", func(t *testing.T) {
		wantError := errors.New("execute failed")

		address := common.HexToAddress("0x1")
		executor := &stubExecutor{
			executeErr: wantError,
		}

		client := contractclient.NewEnygmaHandlerClient(executor, common.HexToAddress("0x000000000000000000000000000000000000E0EC"))

		err := client.MarkSwapCompleted(
			context.Background(),
			address,
			big.NewInt(100),
			"abcd1234567890abcd1234567890abcd1234567890abcd1234567890abcd1234",
		)

		require.NotNil(t, err)
		var wrappedErr *contractclient.EnygmaHandlerClientError
		require.ErrorAs(t, err, &wrappedErr, "error should be wrapped")
		require.ErrorIs(t, err, wantError, "underlying error should be preserved")
	})
}

func TestEnygmaHandlerClient_NotifySenderWithPNCommunicator(t *testing.T) {
	t.Run("successfully notifies sender with pl communicator", func(t *testing.T) {
		wantSharedId := "abcd1234567890abcd1234567890abcd1234567890abcd1234567890abcd1234"

		address := common.HexToAddress("0x1")
		executor := &stubExecutor{}

		client := contractclient.NewEnygmaHandlerClient(executor, common.HexToAddress("0x000000000000000000000000000000000000E0EC"))

		err := client.NotifySenderWithPNCommunicator(context.Background(), address, wantSharedId, types.SwapError, "test message")

		require.Nil(t, err)
		assert.Equal(t, address, executor.spyExecuteAddress)
		assert.NotNil(t, executor.spyExecuteCalldata)
	})

	t.Run("wraps string to bytes32 conversion errors in EnygmaHandlerClientError", func(t *testing.T) {
		address := common.HexToAddress("0x1")
		executor := &stubExecutor{}

		client := contractclient.NewEnygmaHandlerClient(executor, common.HexToAddress("0x000000000000000000000000000000000000E0EC"))

		err := client.NotifySenderWithPNCommunicator(
			context.Background(),
			address,
			"invalid-hex",
			types.SwapError,
			"test message",
		)

		require.NotNil(t, err)
		var wrappedErr *contractclient.EnygmaHandlerClientError
		require.ErrorAs(t, err, &wrappedErr, "error should be wrapped")
	})

	t.Run("wraps executor errors in EnygmaHandlerClientError", func(t *testing.T) {
		wantError := errors.New("execute failed")

		address := common.HexToAddress("0x1")
		executor := &stubExecutor{
			executeErr: wantError,
		}

		client := contractclient.NewEnygmaHandlerClient(executor, common.HexToAddress("0x000000000000000000000000000000000000E0EC"))

		err := client.NotifySenderWithPNCommunicator(
			context.Background(),
			address,
			"abcd1234567890abcd1234567890abcd1234567890abcd1234567890abcd1234",
			types.SwapError,
			"test message",
		)

		require.NotNil(t, err)
		var wrappedErr *contractclient.EnygmaHandlerClientError
		require.ErrorAs(t, err, &wrappedErr, "error should be wrapped")
		require.ErrorIs(t, err, wantError, "underlying error should be preserved")
	})
}
