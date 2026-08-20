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

type StubEnygmaEncryptor struct {
	encryptedBatches [][]byte
	err              error

	spyBatches     []*types.EnygmaTransferBatch
	spyBlockNumber *big.Int
}

func (s *StubEnygmaEncryptor) EncryptEnygmaTransferBatches(
	_ context.Context,
	batches []*types.EnygmaTransferBatch,
	blockNumber *big.Int,
) ([][]byte, error) {
	s.spyBatches = batches
	s.spyBlockNumber = blockNumber
	return s.encryptedBatches, s.err
}

func TestEnygmaClient_GetDvpIntegrationContractAddress(t *testing.T) {
	t.Run("wraps executor call errors in EnygmaClientError", func(t *testing.T) {
		wantError := errors.New("contract call failed")

		address := common.HexToAddress("0x1")
		executor := &stubExecutor{
			callErr: wantError,
		}

		client := contractclient.NewEnygmaClient(executor, nil)

		_, err := client.GetDvpIntegrationContractAddress(context.Background(), address)

		require.NotNil(t, err)
		var wrappedErr *contractclient.EnygmaClientError
		require.ErrorAs(t, err, &wrappedErr, "error should be wrapped in EnygmaClientError")
		require.ErrorIs(t, err, wantError, "underlying error should be preserved")
	})
}

func TestEnygmaClient_GetPublicValuesFinalised(t *testing.T) {
	t.Run("wraps executor call errors in EnygmaClientError", func(t *testing.T) {
		wantError := errors.New("contract call failed")

		address := common.HexToAddress("0x1")
		executor := &stubExecutor{
			callErr: wantError,
		}

		client := contractclient.NewEnygmaClient(executor, nil)

		_, err := client.GetPublicValuesFinalised(context.Background(), address)

		require.NotNil(t, err)
		var wrappedErr *contractclient.EnygmaClientError
		require.ErrorAs(t, err, &wrappedErr, "error should be wrapped in EnygmaClientError")
		require.ErrorIs(t, err, wantError, "underlying error should be preserved")
	})
}

func TestEnygmaClient_SignSupplyUpdate(t *testing.T) {
	t.Run("successfully signs supply update via executor", func(t *testing.T) {
		wantSenderChainId := big.NewInt(1)
		wantBlockNumber := big.NewInt(100)
		wantAmount := big.NewInt(5000)
		wantUpdate := types.EnygmaSupplyUpdate{
			Amount: wantAmount,
			Type:   types.EnygmaDepositToDvp,
		}

		address := common.HexToAddress("0x1")
		executor := &stubExecutor{}

		client := contractclient.NewEnygmaClient(executor, nil)

		err := client.SupplyUpdate(context.Background(), "test-event-id", address, wantSenderChainId, wantBlockNumber, wantUpdate)

		require.Nil(t, err)
		assert.Equal(t, address, executor.spyExecuteAddress)
		assert.NotNil(t, executor.spyExecuteCalldata)
	})

	t.Run("wraps executor sign errors in EnygmaClientError", func(t *testing.T) {
		wantError := errors.New("sign failed")

		address := common.HexToAddress("0x1")
		executor := &stubExecutor{
			executeErr: wantError,
		}

		client := contractclient.NewEnygmaClient(executor, nil)

		err := client.SupplyUpdate(
			context.Background(), "test-event-id",
			address,
			big.NewInt(1),
			big.NewInt(100),
			types.EnygmaSupplyUpdate{Amount: big.NewInt(0)},
		)

		require.NotNil(t, err)
		var wrappedErr *contractclient.EnygmaClientError
		require.ErrorAs(t, err, &wrappedErr, "error should be wrapped")
		require.ErrorIs(t, err, wantError, "underlying error should be preserved")
	})
}

func TestEnygmaClient_SignTransferBatch(t *testing.T) {
	createTestBatchesEnygma := func() []*types.EnygmaTransferBatch {
		return []*types.EnygmaTransferBatch{
			{
				ToChainID: big.NewInt(100),
			},
		}
	}

	createTestProofEnygma := func() *types.EnygmaProofResponse {
		return &types.EnygmaProofResponse{
			PiA: [2]*big.Int{big.NewInt(1), big.NewInt(2)},
			PiB: [2][2]*big.Int{
				{big.NewInt(3), big.NewInt(4)},
				{big.NewInt(5), big.NewInt(6)},
			},
			PiC:          [2]*big.Int{big.NewInt(7), big.NewInt(8)},
			PublicSignal: []*big.Int{big.NewInt(9), big.NewInt(10)},
		}
	}

	t.Run("successfully signs transfer batch via executor", func(t *testing.T) {
		wantBlockNumber := big.NewInt(100)
		wantBatches := createTestBatchesEnygma()
		wantProof := createTestProofEnygma()

		address := common.HexToAddress("0x1")
		executor := &stubExecutor{}
		encryptor := &StubEnygmaEncryptor{
			encryptedBatches: [][]byte{{0xaa, 0xbb}, {0xcc, 0xdd}},
		}

		client := contractclient.NewEnygmaClient(executor, encryptor)

		err := client.TransferBatch(context.Background(), "test-event-id", address, wantBatches, wantProof, wantBlockNumber)

		require.Nil(t, err)
		assert.Equal(t, wantBatches, encryptor.spyBatches)
		assert.Equal(t, wantBlockNumber, encryptor.spyBlockNumber)
		assert.Equal(t, address, executor.spyExecuteAddress)
		assert.NotNil(t, executor.spyExecuteCalldata)
	})

	t.Run("wraps encryption errors in EnygmaClientError", func(t *testing.T) {
		wantError := errors.New("encryption failed")

		address := common.HexToAddress("0x1")
		executor := &stubExecutor{}
		encryptor := &StubEnygmaEncryptor{
			err: wantError,
		}

		client := contractclient.NewEnygmaClient(executor, encryptor)

		err := client.TransferBatch(
			context.Background(), "test-event-id",
			address,
			createTestBatchesEnygma(),
			createTestProofEnygma(),
			big.NewInt(100),
		)

		require.NotNil(t, err)
		var wrappedErr *contractclient.EnygmaClientError
		require.ErrorAs(t, err, &wrappedErr, "error should be wrapped")
		require.ErrorIs(t, err, wantError, "underlying error should be preserved")
	})

	t.Run("wraps executor sign errors in EnygmaClientError", func(t *testing.T) {
		wantError := errors.New("sign failed")

		address := common.HexToAddress("0x1")
		executor := &stubExecutor{
			executeErr: wantError,
		}
		encryptor := &StubEnygmaEncryptor{
			encryptedBatches: [][]byte{{0xaa}},
		}

		client := contractclient.NewEnygmaClient(executor, encryptor)

		err := client.TransferBatch(
			context.Background(), "test-event-id",
			address,
			createTestBatchesEnygma(),
			createTestProofEnygma(),
			big.NewInt(100),
		)

		require.NotNil(t, err)
		var wrappedErr *contractclient.EnygmaClientError
		require.ErrorAs(t, err, &wrappedErr, "error should be wrapped")
		require.ErrorIs(t, err, wantError, "underlying error should be preserved")
	})
}
