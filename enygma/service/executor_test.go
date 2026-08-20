package service_test

import (
	"context"
	"errors"
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"

	keyspb "github.com/raylsnetwork/rayls-privacy-relayer-api/cts/gen/proto/keys"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/dvp"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/enygma"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/enygma/service"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/enygma/testutils"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/types"
)

// NOTE on this rewrite:
// The previous test file in `issue/173` exercised an older richer executor
// that owned a transaction-recovery repository, a TxOps sender, a contract
// clients factory, a broadcaster, a receipter and a simulator. The CTS
// migration on `feat/cts` collapsed those responsibilities — the
// `EnygmaExecutor` now talks directly to a single `executorEnygmaClient`
// and `ExecutorDvpIntegrationClient` plus a `executorKeysClient` (gRPC).
// Test cases that exercised the removed pipeline (crash recovery resume,
// contract factory failures, broadcaster/receipter assertions, TxOps-send
// counters) have been dropped because the underlying behaviour no longer
// exists in `executor.go`. The remaining cases still cover happy-path
// success and every error branch the executor produces today.

func executorTestResourceId() string {
	return "test-resource-123"
}

func executorTestBlockNumber() uint64 {
	return 100
}

func executorTestTxHash() common.Hash {
	return common.HexToHash("0xabcd")
}

func executorTestFromAddress() common.Address {
	return common.HexToAddress("0xfrom1234")
}

func executorTestEnygmaAddress() common.Address {
	return common.HexToAddress("0x1234567890123456789012345678901234567890")
}

func executorTestAmount() *big.Int {
	return big.NewInt(1000)
}

func executorTestPnChainId() *big.Int {
	return big.NewInt(1)
}

func executorTestConfig() service.ExecutorConfig {
	return service.ExecutorConfig{
		DefaultContextTimeout: 5 * time.Second,
		MaxNumberOfJSDeposits: 10,
	}
}

// successKeysClient returns a deterministic payment-spend key tuple.
func successKeysClient() *executorKeysClientMock {
	return &executorKeysClientMock{
		GetPaymentSpendKeyFunc: func(
			ctx context.Context,
			in *keyspb.GetPaymentSpendKeyRequest,
			opts ...grpc.CallOption,
		) (*keyspb.PaymentSpendKeyResponse, error) {
			return &keyspb.PaymentSpendKeyResponse{
				PublicKey: big.NewInt(1).Bytes(),
				SecretKey: big.NewInt(999).Bytes(),
			}, nil
		},
	}
}

func TestEnygmaExecutor_ExecuteEnygmaSupplyUpdate(t *testing.T) {
	t.Run("successfully executes mint supply update", func(t *testing.T) {
		ctx := context.Background()
		conf := executorTestConfig()
		resourceId := executorTestResourceId()
		blockNumber := executorTestBlockNumber()
		enygmaAddress := executorTestEnygmaAddress()
		amount := executorTestAmount()
		plChainId := executorTestPnChainId()

		tracer := &testutils.MockTracer{}
		batcher := &executorEnygmaBatcherMock{}
		proofGen := &executorProofGeneratorMock{}
		repository := &executorEnygmaHistoryRepositoryMock{
			InsertEnygmaHistoryFunc: func(ctx context.Context, history types.EnygmaHistory) error {
				return nil
			},
		}
		keysClient := successKeysClient()
		enygmaClient := &executorEnygmaClientMock{
			SupplyUpdateFunc: func(
				ctx context.Context,
				_ string,
				tokenAddress common.Address,
				senderChainId *big.Int,
				blockNumber *big.Int,
				update types.EnygmaSupplyUpdate,
			) error {
				return nil
			},
		}
		integrationClient := &ExecutorDvpIntegrationClientMock{}
		commitmentCalc := &executorCommitmentCalculatorMock{}

		executor := service.NewEnygmaExecutor(
			conf,
			tracer,
			batcher,
			proofGen,
			repository,
			keysClient,
			enygmaClient,
			integrationClient,
			commitmentCalc,
			plChainId,
		)

		supplyUpdate := types.EnygmaSupplyUpdate{
			Type:   types.EnygmaMint,
			Amount: amount,
		}
		err := executor.ExecuteEnygmaSupplyUpdate(ctx, "test-batch-id", resourceId, blockNumber, supplyUpdate, enygmaAddress)

		require.NoError(t, err)
		require.Len(t, repository.InsertEnygmaHistoryCalls(), 1)
		history := repository.InsertEnygmaHistoryCalls()[0].History
		assert.Equal(t, types.EnygmaMint, history.EventType)
		assert.Equal(t, amount, history.BalanceChange)
		assert.Equal(t, plChainId, history.FromChainId)
		assert.Equal(t, big.NewInt(0), history.RFactor)
		assert.Equal(t, new(big.Int).SetUint64(blockNumber), history.BlockNumberPrivateHub)
		assert.Len(t, enygmaClient.SupplyUpdateCalls(), 1)
	})

	t.Run("successfully executes burn supply update", func(t *testing.T) {
		ctx := context.Background()
		conf := executorTestConfig()
		resourceId := executorTestResourceId()
		blockNumber := executorTestBlockNumber()
		enygmaAddress := executorTestEnygmaAddress()
		amount := executorTestAmount()
		plChainId := executorTestPnChainId()

		tracer := &testutils.MockTracer{}
		batcher := &executorEnygmaBatcherMock{}
		proofGen := &executorProofGeneratorMock{}
		repository := &executorEnygmaHistoryRepositoryMock{
			InsertEnygmaHistoryFunc: func(ctx context.Context, history types.EnygmaHistory) error {
				return nil
			},
		}
		keysClient := successKeysClient()
		enygmaClient := &executorEnygmaClientMock{
			SupplyUpdateFunc: func(
				ctx context.Context,
				_ string,
				tokenAddress common.Address,
				senderChainId *big.Int,
				blockNumber *big.Int,
				update types.EnygmaSupplyUpdate,
			) error {
				return nil
			},
		}
		integrationClient := &ExecutorDvpIntegrationClientMock{}
		commitmentCalc := &executorCommitmentCalculatorMock{}

		executor := service.NewEnygmaExecutor(
			conf,
			tracer,
			batcher,
			proofGen,
			repository,
			keysClient,
			enygmaClient,
			integrationClient,
			commitmentCalc,
			plChainId,
		)

		supplyUpdate := types.EnygmaSupplyUpdate{
			Type:   types.EnygmaBurn,
			Amount: amount,
		}
		err := executor.ExecuteEnygmaSupplyUpdate(ctx, "test-batch-id", resourceId, blockNumber, supplyUpdate, enygmaAddress)

		require.NoError(t, err)
		require.Len(t, repository.InsertEnygmaHistoryCalls(), 1)
		history := repository.InsertEnygmaHistoryCalls()[0].History
		assert.Equal(t, types.EnygmaBurn, history.EventType)
		assert.Equal(t, plChainId, history.FromChainId)
		assert.Equal(t, big.NewInt(0), history.RFactor)
		assert.Equal(t, new(big.Int).SetUint64(blockNumber), history.BlockNumberPrivateHub)
		// For burn, balance change should be negative.
		assert.Equal(t, big.NewInt(-1000), history.BalanceChange)
	})

	t.Run("returns error for invalid supply update type", func(t *testing.T) {
		ctx := context.Background()
		conf := executorTestConfig()
		resourceId := executorTestResourceId()
		blockNumber := executorTestBlockNumber()
		enygmaAddress := executorTestEnygmaAddress()
		amount := executorTestAmount()
		plChainId := executorTestPnChainId()

		tracer := &testutils.MockTracer{}
		batcher := &executorEnygmaBatcherMock{}
		proofGen := &executorProofGeneratorMock{}
		repository := &executorEnygmaHistoryRepositoryMock{}
		keysClient := successKeysClient()
		enygmaClient := &executorEnygmaClientMock{}
		integrationClient := &ExecutorDvpIntegrationClientMock{}
		commitmentCalc := &executorCommitmentCalculatorMock{}

		executor := service.NewEnygmaExecutor(
			conf,
			tracer,
			batcher,
			proofGen,
			repository,
			keysClient,
			enygmaClient,
			integrationClient,
			commitmentCalc,
			plChainId,
		)

		supplyUpdate := types.EnygmaSupplyUpdate{
			Type:   types.EnygmaTransfer,
			Amount: amount,
		}
		err := executor.ExecuteEnygmaSupplyUpdate(ctx, "test-batch-id", resourceId, blockNumber, supplyUpdate, enygmaAddress)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "invalid supply update type")
		assert.Empty(t, repository.InsertEnygmaHistoryCalls())
	})

	t.Run("returns error when SupplyUpdate fails", func(t *testing.T) {
		ctx := context.Background()
		conf := executorTestConfig()
		resourceId := executorTestResourceId()
		blockNumber := executorTestBlockNumber()
		enygmaAddress := executorTestEnygmaAddress()
		amount := executorTestAmount()
		plChainId := executorTestPnChainId()

		tracer := &testutils.MockTracer{}
		batcher := &executorEnygmaBatcherMock{}
		proofGen := &executorProofGeneratorMock{}
		repository := &executorEnygmaHistoryRepositoryMock{}
		keysClient := successKeysClient()
		enygmaClient := &executorEnygmaClientMock{
			SupplyUpdateFunc: func(
				ctx context.Context,
				_ string,
				tokenAddress common.Address,
				senderChainId *big.Int,
				blockNumber *big.Int,
				update types.EnygmaSupplyUpdate,
			) error {
				return errors.New("supply update failed")
			},
		}
		integrationClient := &ExecutorDvpIntegrationClientMock{}
		commitmentCalc := &executorCommitmentCalculatorMock{}

		executor := service.NewEnygmaExecutor(
			conf,
			tracer,
			batcher,
			proofGen,
			repository,
			keysClient,
			enygmaClient,
			integrationClient,
			commitmentCalc,
			plChainId,
		)

		supplyUpdate := types.EnygmaSupplyUpdate{
			Type:   types.EnygmaMint,
			Amount: amount,
		}
		err := executor.ExecuteEnygmaSupplyUpdate(ctx, "test-batch-id", resourceId, blockNumber, supplyUpdate, enygmaAddress)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "supply update failed")
		assert.Empty(t, repository.InsertEnygmaHistoryCalls())
	})

	t.Run("returns error when repository insertion fails", func(t *testing.T) {
		ctx := context.Background()
		conf := executorTestConfig()
		resourceId := executorTestResourceId()
		blockNumber := executorTestBlockNumber()
		enygmaAddress := executorTestEnygmaAddress()
		amount := executorTestAmount()
		plChainId := executorTestPnChainId()

		tracer := &testutils.MockTracer{}
		batcher := &executorEnygmaBatcherMock{}
		proofGen := &executorProofGeneratorMock{}
		repository := &executorEnygmaHistoryRepositoryMock{
			InsertEnygmaHistoryFunc: func(ctx context.Context, history types.EnygmaHistory) error {
				return errors.New("database error")
			},
		}
		keysClient := successKeysClient()
		enygmaClient := &executorEnygmaClientMock{
			SupplyUpdateFunc: func(
				ctx context.Context,
				_ string,
				tokenAddress common.Address,
				senderChainId *big.Int,
				blockNumber *big.Int,
				update types.EnygmaSupplyUpdate,
			) error {
				return nil
			},
		}
		integrationClient := &ExecutorDvpIntegrationClientMock{}
		commitmentCalc := &executorCommitmentCalculatorMock{}

		executor := service.NewEnygmaExecutor(
			conf,
			tracer,
			batcher,
			proofGen,
			repository,
			keysClient,
			enygmaClient,
			integrationClient,
			commitmentCalc,
			plChainId,
		)

		supplyUpdate := types.EnygmaSupplyUpdate{
			Type:   types.EnygmaMint,
			Amount: amount,
		}
		err := executor.ExecuteEnygmaSupplyUpdate(ctx, "test-batch-id", resourceId, blockNumber, supplyUpdate, enygmaAddress)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "database error")
	})
}

func TestEnygmaExecutor_ExecuteEnygmaCrossTransfer(t *testing.T) {
	t.Run("successfully executes cross transfer", func(t *testing.T) {
		ctx := context.Background()
		conf := executorTestConfig()
		resourceId := executorTestResourceId()
		blockNumber := executorTestBlockNumber()
		enygmaAddress := executorTestEnygmaAddress()
		plChainId := executorTestPnChainId()
		toChainId := big.NewInt(2)
		senderRFactor := big.NewInt(100)

		txsByChainID := make(map[string][]*types.EnygmaTransferBatchTx)

		tracer := &testutils.MockTracer{}
		batcher := &executorEnygmaBatcherMock{
			CreateBatchesWithAnonimityFunc: func(
				ctx context.Context,
				resourceId string,
				blockNumber *big.Int,
				txsByChainID map[string][]*types.EnygmaTransferBatchTx,
			) ([]*types.EnygmaTransferBatch, error) {
				return []*types.EnygmaTransferBatch{
					{
						ResourceId:   resourceId,
						ToChainID:    plChainId,
						FromChainID:  plChainId,
						Transactions: []*types.EnygmaTransferBatchTx{},
					},
					{
						ResourceId:  resourceId,
						ToChainID:   toChainId,
						FromChainID: plChainId,
						Transactions: []*types.EnygmaTransferBatchTx{
							{
								MessageId:   "msg-1",
								ReferenceId: [32]byte{1, 2, 3, 4},
								FromAddress: common.HexToAddress("0x1111111111111111111111111111111111111111"),
								ToAmount:    big.NewInt(100),
								ToAddress:   common.HexToAddress("0x2222222222222222222222222222222222222222"),
							},
						},
					},
				}, nil
			},
		}
		proofGen := &executorProofGeneratorMock{
			GenerateTransferProofFunc: func(
				ctx context.Context,
				params enygma.TransferProofParams,
			) (*types.EnygmaProofResponse, []*types.Point, []*big.Int, []*big.Int, error) {
				return &types.EnygmaProofResponse{}, []*types.Point{}, []*big.Int{
					senderRFactor,
					big.NewInt(200),
				}, []*big.Int{}, nil
			},
		}
		repository := &executorEnygmaHistoryRepositoryMock{
			InsertEnygmaHistoryFunc: func(ctx context.Context, history types.EnygmaHistory) error {
				return nil
			},
		}
		keysClient := successKeysClient()
		enygmaClient := &executorEnygmaClientMock{
			TransferBatchFunc: func(
				ctx context.Context,
				_ string,
				tokenAddress common.Address,
				batches []*types.EnygmaTransferBatch,
				proof *types.EnygmaProofResponse,
				blockNumber *big.Int,
			) error {
				return nil
			},
		}
		integrationClient := &ExecutorDvpIntegrationClientMock{}
		commitmentCalc := &executorCommitmentCalculatorMock{}

		executor := service.NewEnygmaExecutor(
			conf,
			tracer,
			batcher,
			proofGen,
			repository,
			keysClient,
			enygmaClient,
			integrationClient,
			commitmentCalc,
			plChainId,
		)

		err := executor.ExecuteEnygmaCrossTransfer(ctx, "test-batch-id", blockNumber, resourceId, txsByChainID, enygmaAddress)

		require.NoError(t, err)
		assert.Len(t, batcher.CreateBatchesWithAnonimityCalls(), 1)
		assert.Len(t, proofGen.GenerateTransferProofCalls(), 1)
		assert.Len(t, enygmaClient.TransferBatchCalls(), 1)
		require.Len(t, repository.InsertEnygmaHistoryCalls(), 1)
		history := repository.InsertEnygmaHistoryCalls()[0].History
		assert.Equal(t, types.EnygmaTransfer, history.EventType)
		assert.Equal(t, senderRFactor, history.RFactor)
		assert.Equal(t, big.NewInt(-100), history.BalanceChange)
		assert.Equal(t, plChainId, history.FromChainId)
		assert.Equal(t, new(big.Int).SetUint64(blockNumber), history.BlockNumberPrivateHub)
	})

	t.Run("returns error when batcher fails", func(t *testing.T) {
		ctx := context.Background()
		conf := executorTestConfig()
		resourceId := executorTestResourceId()
		blockNumber := executorTestBlockNumber()
		enygmaAddress := executorTestEnygmaAddress()
		plChainId := executorTestPnChainId()

		txsByChainID := make(map[string][]*types.EnygmaTransferBatchTx)

		tracer := &testutils.MockTracer{}
		batcher := &executorEnygmaBatcherMock{
			CreateBatchesWithAnonimityFunc: func(
				ctx context.Context,
				resourceId string,
				blockNumber *big.Int,
				txsByChainID map[string][]*types.EnygmaTransferBatchTx,
			) ([]*types.EnygmaTransferBatch, error) {
				return nil, errors.New("batching failed")
			},
		}
		proofGen := &executorProofGeneratorMock{}
		repository := &executorEnygmaHistoryRepositoryMock{}
		keysClient := successKeysClient()
		enygmaClient := &executorEnygmaClientMock{}
		integrationClient := &ExecutorDvpIntegrationClientMock{}
		commitmentCalc := &executorCommitmentCalculatorMock{}

		executor := service.NewEnygmaExecutor(
			conf,
			tracer,
			batcher,
			proofGen,
			repository,
			keysClient,
			enygmaClient,
			integrationClient,
			commitmentCalc,
			plChainId,
		)

		err := executor.ExecuteEnygmaCrossTransfer(ctx, "test-batch-id", blockNumber, resourceId, txsByChainID, enygmaAddress)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "batching failed")
		assert.Empty(t, proofGen.GenerateTransferProofCalls())
	})

	t.Run("returns error when proof generation fails", func(t *testing.T) {
		ctx := context.Background()
		conf := executorTestConfig()
		resourceId := executorTestResourceId()
		blockNumber := executorTestBlockNumber()
		enygmaAddress := executorTestEnygmaAddress()
		plChainId := executorTestPnChainId()

		txsByChainID := make(map[string][]*types.EnygmaTransferBatchTx)

		tracer := &testutils.MockTracer{}
		batcher := &executorEnygmaBatcherMock{
			CreateBatchesWithAnonimityFunc: func(
				ctx context.Context,
				resourceId string,
				blockNumber *big.Int,
				txsByChainID map[string][]*types.EnygmaTransferBatchTx,
			) ([]*types.EnygmaTransferBatch, error) {
				return []*types.EnygmaTransferBatch{}, nil
			},
		}
		proofGen := &executorProofGeneratorMock{
			GenerateTransferProofFunc: func(
				ctx context.Context,
				params enygma.TransferProofParams,
			) (*types.EnygmaProofResponse, []*types.Point, []*big.Int, []*big.Int, error) {
				return nil, nil, nil, nil, errors.New("proof generation failed")
			},
		}
		repository := &executorEnygmaHistoryRepositoryMock{}
		keysClient := successKeysClient()
		enygmaClient := &executorEnygmaClientMock{}
		integrationClient := &ExecutorDvpIntegrationClientMock{}
		commitmentCalc := &executorCommitmentCalculatorMock{}

		executor := service.NewEnygmaExecutor(
			conf,
			tracer,
			batcher,
			proofGen,
			repository,
			keysClient,
			enygmaClient,
			integrationClient,
			commitmentCalc,
			plChainId,
		)

		err := executor.ExecuteEnygmaCrossTransfer(ctx, "test-batch-id", blockNumber, resourceId, txsByChainID, enygmaAddress)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "proof generation failed")
		assert.Empty(t, enygmaClient.TransferBatchCalls())
	})

	t.Run("returns error when TransferBatch fails", func(t *testing.T) {
		ctx := context.Background()
		conf := executorTestConfig()
		resourceId := executorTestResourceId()
		blockNumber := executorTestBlockNumber()
		enygmaAddress := executorTestEnygmaAddress()
		plChainId := executorTestPnChainId()

		txsByChainID := make(map[string][]*types.EnygmaTransferBatchTx)

		tracer := &testutils.MockTracer{}
		batcher := &executorEnygmaBatcherMock{
			CreateBatchesWithAnonimityFunc: func(
				ctx context.Context,
				resourceId string,
				blockNumber *big.Int,
				txsByChainID map[string][]*types.EnygmaTransferBatchTx,
			) ([]*types.EnygmaTransferBatch, error) {
				return []*types.EnygmaTransferBatch{}, nil
			},
		}
		proofGen := &executorProofGeneratorMock{
			GenerateTransferProofFunc: func(
				ctx context.Context,
				params enygma.TransferProofParams,
			) (*types.EnygmaProofResponse, []*types.Point, []*big.Int, []*big.Int, error) {
				return &types.EnygmaProofResponse{}, []*types.Point{}, []*big.Int{}, []*big.Int{}, nil
			},
		}
		repository := &executorEnygmaHistoryRepositoryMock{}
		keysClient := successKeysClient()
		enygmaClient := &executorEnygmaClientMock{
			TransferBatchFunc: func(
				ctx context.Context,
				_ string,
				tokenAddress common.Address,
				batches []*types.EnygmaTransferBatch,
				proof *types.EnygmaProofResponse,
				blockNumber *big.Int,
			) error {
				return errors.New("send transfer batch failed")
			},
		}
		integrationClient := &ExecutorDvpIntegrationClientMock{}
		commitmentCalc := &executorCommitmentCalculatorMock{}

		executor := service.NewEnygmaExecutor(
			conf,
			tracer,
			batcher,
			proofGen,
			repository,
			keysClient,
			enygmaClient,
			integrationClient,
			commitmentCalc,
			plChainId,
		)

		err := executor.ExecuteEnygmaCrossTransfer(ctx, "test-batch-id", blockNumber, resourceId, txsByChainID, enygmaAddress)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "send transfer batch failed")
		assert.Empty(t, repository.InsertEnygmaHistoryCalls())
	})

	t.Run("returns error when repository insertion fails", func(t *testing.T) {
		ctx := context.Background()
		conf := executorTestConfig()
		resourceId := executorTestResourceId()
		blockNumber := executorTestBlockNumber()
		enygmaAddress := executorTestEnygmaAddress()
		plChainId := executorTestPnChainId()

		txsByChainID := make(map[string][]*types.EnygmaTransferBatchTx)

		tracer := &testutils.MockTracer{}
		batcher := &executorEnygmaBatcherMock{
			CreateBatchesWithAnonimityFunc: func(
				ctx context.Context,
				resourceId string,
				blockNumber *big.Int,
				txsByChainID map[string][]*types.EnygmaTransferBatchTx,
			) ([]*types.EnygmaTransferBatch, error) {
				return []*types.EnygmaTransferBatch{}, nil
			},
		}
		proofGen := &executorProofGeneratorMock{
			GenerateTransferProofFunc: func(
				ctx context.Context,
				params enygma.TransferProofParams,
			) (*types.EnygmaProofResponse, []*types.Point, []*big.Int, []*big.Int, error) {
				return &types.EnygmaProofResponse{}, []*types.Point{}, []*big.Int{}, []*big.Int{}, nil
			},
		}
		repository := &executorEnygmaHistoryRepositoryMock{
			InsertEnygmaHistoryFunc: func(ctx context.Context, history types.EnygmaHistory) error {
				return errors.New("database error")
			},
		}
		keysClient := successKeysClient()
		enygmaClient := &executorEnygmaClientMock{
			TransferBatchFunc: func(
				ctx context.Context,
				_ string,
				tokenAddress common.Address,
				batches []*types.EnygmaTransferBatch,
				proof *types.EnygmaProofResponse,
				blockNumber *big.Int,
			) error {
				return nil
			},
		}
		integrationClient := &ExecutorDvpIntegrationClientMock{}
		commitmentCalc := &executorCommitmentCalculatorMock{}

		executor := service.NewEnygmaExecutor(
			conf,
			tracer,
			batcher,
			proofGen,
			repository,
			keysClient,
			enygmaClient,
			integrationClient,
			commitmentCalc,
			plChainId,
		)

		err := executor.ExecuteEnygmaCrossTransfer(ctx, "test-batch-id", blockNumber, resourceId, txsByChainID, enygmaAddress)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "database error")
	})
}

func TestEnygmaExecutor_ExecuteEnygmaDeposit(t *testing.T) {
	t.Run("successfully executes deposit", func(t *testing.T) {
		ctx := context.Background()
		conf := executorTestConfig()
		resourceId := executorTestResourceId()
		blockNumber := executorTestBlockNumber()
		enygmaAddress := executorTestEnygmaAddress()
		amount := executorTestAmount()
		txHash := executorTestTxHash()
		fromAddress := executorTestFromAddress()
		depositCommitment := big.NewInt(100)
		depositSalt := big.NewInt(200)
		plChainId := executorTestPnChainId()
		toChainId := big.NewInt(2)
		integrationAddress := common.HexToAddress("0xintegration1234")

		senderRFactor := big.NewInt(456)

		tracer := &testutils.MockTracer{}
		batcher := &executorEnygmaBatcherMock{
			CreateBatchesWithAnonimityFunc: func(
				ctx context.Context,
				resourceId string,
				blockNumber *big.Int,
				txsByChainID map[string][]*types.EnygmaTransferBatchTx,
			) ([]*types.EnygmaTransferBatch, error) {
				return []*types.EnygmaTransferBatch{
					{
						ResourceId:   resourceId,
						ToChainID:    plChainId,
						FromChainID:  plChainId,
						Transactions: []*types.EnygmaTransferBatchTx{},
					},
					{
						ResourceId:   resourceId,
						ToChainID:    toChainId,
						FromChainID:  plChainId,
						Transactions: []*types.EnygmaTransferBatchTx{},
					},
				}, nil
			},
		}
		proofGen := &executorProofGeneratorMock{
			GenerateDepositProofFunc: func(
				ctx context.Context,
				params enygma.DepositProofParams,
			) (*types.EnygmaProofResponse, []*types.Point, []*big.Int, []*big.Int, error) {
				return &types.EnygmaProofResponse{}, []*types.Point{}, []*big.Int{
					senderRFactor,
					big.NewInt(789),
				}, []*big.Int{}, nil
			},
		}
		repository := &executorEnygmaHistoryRepositoryMock{
			InsertEnygmaHistoryFunc: func(ctx context.Context, history types.EnygmaHistory) error {
				return nil
			},
		}
		keysClient := successKeysClient()
		enygmaClient := &executorEnygmaClientMock{
			GetDvpIntegrationContractAddressFunc: func(_ context.Context, tokenAddress common.Address) (common.Address, error) {
				return integrationAddress, nil
			},
		}
		integrationClient := &ExecutorDvpIntegrationClientMock{
			DepositFunc: func(
				ctx context.Context,
				_ string,
				batches []*types.EnygmaTransferBatch,
				proof *types.EnygmaProofResponse,
				blockNumber *big.Int,
				chainId *big.Int,
				resourceId string,
				amount *big.Int,
				from common.Address,
				sourceTxHash common.Hash,
				dvpIntegrationAddress common.Address,
			) error {
				return nil
			},
		}
		commitmentCalc := &executorCommitmentCalculatorMock{}

		executor := service.NewEnygmaExecutor(
			conf,
			tracer,
			batcher,
			proofGen,
			repository,
			keysClient,
			enygmaClient,
			integrationClient,
			commitmentCalc,
			plChainId,
		)

		err := executor.ExecuteEnygmaDeposit(
			ctx,
			"test-batch-id",
			resourceId,
			amount,
			blockNumber,
			depositCommitment,
			depositSalt,
			fromAddress,
			txHash,
			enygmaAddress,
		)

		require.NoError(t, err)
		assert.Len(t, batcher.CreateBatchesWithAnonimityCalls(), 1)
		assert.Len(t, proofGen.GenerateDepositProofCalls(), 1)
		assert.Len(t, enygmaClient.GetDvpIntegrationContractAddressCalls(), 1)
		assert.Len(t, integrationClient.DepositCalls(), 1)
		require.Len(t, repository.InsertEnygmaHistoryCalls(), 1)
		history := repository.InsertEnygmaHistoryCalls()[0].History
		assert.Equal(t, types.EnygmaDepositToDvp, history.EventType)
		assert.Equal(t, plChainId, history.FromChainId)
		assert.Equal(t, senderRFactor, history.RFactor)
		assert.Equal(t, new(big.Int).Neg(amount), history.BalanceChange)
		assert.Equal(t, new(big.Int).SetUint64(blockNumber), history.BlockNumberPrivateHub)
		// The integration client should receive the address resolved by the enygma client.
		assert.Equal(t, integrationAddress, integrationClient.DepositCalls()[0].DvpIntegrationAddress)
	})

	t.Run("returns error when batcher fails", func(t *testing.T) {
		ctx := context.Background()
		conf := executorTestConfig()
		resourceId := executorTestResourceId()
		blockNumber := executorTestBlockNumber()
		enygmaAddress := executorTestEnygmaAddress()
		amount := executorTestAmount()
		txHash := executorTestTxHash()
		fromAddress := executorTestFromAddress()
		depositCommitment := big.NewInt(100)
		depositSalt := big.NewInt(200)
		plChainId := executorTestPnChainId()

		tracer := &testutils.MockTracer{}
		batcher := &executorEnygmaBatcherMock{
			CreateBatchesWithAnonimityFunc: func(
				ctx context.Context,
				resourceId string,
				blockNumber *big.Int,
				txsByChainID map[string][]*types.EnygmaTransferBatchTx,
			) ([]*types.EnygmaTransferBatch, error) {
				return nil, errors.New("batching failed")
			},
		}
		proofGen := &executorProofGeneratorMock{}
		repository := &executorEnygmaHistoryRepositoryMock{}
		keysClient := successKeysClient()
		enygmaClient := &executorEnygmaClientMock{}
		integrationClient := &ExecutorDvpIntegrationClientMock{}
		commitmentCalc := &executorCommitmentCalculatorMock{}

		executor := service.NewEnygmaExecutor(
			conf,
			tracer,
			batcher,
			proofGen,
			repository,
			keysClient,
			enygmaClient,
			integrationClient,
			commitmentCalc,
			plChainId,
		)

		err := executor.ExecuteEnygmaDeposit(
			ctx,
			"test-batch-id",
			resourceId,
			amount,
			blockNumber,
			depositCommitment,
			depositSalt,
			fromAddress,
			txHash,
			enygmaAddress,
		)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "batching failed")
	})

	t.Run("returns error when keys client fails", func(t *testing.T) {
		ctx := context.Background()
		conf := executorTestConfig()
		resourceId := executorTestResourceId()
		blockNumber := executorTestBlockNumber()
		enygmaAddress := executorTestEnygmaAddress()
		amount := executorTestAmount()
		txHash := executorTestTxHash()
		fromAddress := executorTestFromAddress()
		depositCommitment := big.NewInt(100)
		depositSalt := big.NewInt(200)
		plChainId := executorTestPnChainId()

		tracer := &testutils.MockTracer{}
		batcher := &executorEnygmaBatcherMock{
			CreateBatchesWithAnonimityFunc: func(
				ctx context.Context,
				resourceId string,
				blockNumber *big.Int,
				txsByChainID map[string][]*types.EnygmaTransferBatchTx,
			) ([]*types.EnygmaTransferBatch, error) {
				return []*types.EnygmaTransferBatch{}, nil
			},
		}
		proofGen := &executorProofGeneratorMock{}
		repository := &executorEnygmaHistoryRepositoryMock{}
		keysClient := &executorKeysClientMock{
			GetPaymentSpendKeyFunc: func(
				ctx context.Context,
				in *keyspb.GetPaymentSpendKeyRequest,
				opts ...grpc.CallOption,
			) (*keyspb.PaymentSpendKeyResponse, error) {
				return nil, errors.New("keys client failed")
			},
		}
		enygmaClient := &executorEnygmaClientMock{}
		integrationClient := &ExecutorDvpIntegrationClientMock{}
		commitmentCalc := &executorCommitmentCalculatorMock{}

		executor := service.NewEnygmaExecutor(
			conf,
			tracer,
			batcher,
			proofGen,
			repository,
			keysClient,
			enygmaClient,
			integrationClient,
			commitmentCalc,
			plChainId,
		)

		err := executor.ExecuteEnygmaDeposit(
			ctx,
			"test-batch-id",
			resourceId,
			amount,
			blockNumber,
			depositCommitment,
			depositSalt,
			fromAddress,
			txHash,
			enygmaAddress,
		)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "keys client failed")
	})

	t.Run("returns error when proof generation fails", func(t *testing.T) {
		ctx := context.Background()
		conf := executorTestConfig()
		resourceId := executorTestResourceId()
		blockNumber := executorTestBlockNumber()
		enygmaAddress := executorTestEnygmaAddress()
		amount := executorTestAmount()
		txHash := executorTestTxHash()
		fromAddress := executorTestFromAddress()
		depositCommitment := big.NewInt(100)
		depositSalt := big.NewInt(200)
		plChainId := executorTestPnChainId()

		tracer := &testutils.MockTracer{}
		batcher := &executorEnygmaBatcherMock{
			CreateBatchesWithAnonimityFunc: func(
				ctx context.Context,
				resourceId string,
				blockNumber *big.Int,
				txsByChainID map[string][]*types.EnygmaTransferBatchTx,
			) ([]*types.EnygmaTransferBatch, error) {
				return []*types.EnygmaTransferBatch{}, nil
			},
		}
		proofGen := &executorProofGeneratorMock{
			GenerateDepositProofFunc: func(
				ctx context.Context,
				params enygma.DepositProofParams,
			) (*types.EnygmaProofResponse, []*types.Point, []*big.Int, []*big.Int, error) {
				return nil, nil, nil, nil, errors.New("proof generation failed")
			},
		}
		repository := &executorEnygmaHistoryRepositoryMock{}
		keysClient := successKeysClient()
		enygmaClient := &executorEnygmaClientMock{}
		integrationClient := &ExecutorDvpIntegrationClientMock{}
		commitmentCalc := &executorCommitmentCalculatorMock{}

		executor := service.NewEnygmaExecutor(
			conf,
			tracer,
			batcher,
			proofGen,
			repository,
			keysClient,
			enygmaClient,
			integrationClient,
			commitmentCalc,
			plChainId,
		)

		err := executor.ExecuteEnygmaDeposit(
			ctx,
			"test-batch-id",
			resourceId,
			amount,
			blockNumber,
			depositCommitment,
			depositSalt,
			fromAddress,
			txHash,
			enygmaAddress,
		)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "proof generation failed")
	})

	t.Run("returns error when getting integration address fails", func(t *testing.T) {
		ctx := context.Background()
		conf := executorTestConfig()
		resourceId := executorTestResourceId()
		blockNumber := executorTestBlockNumber()
		enygmaAddress := executorTestEnygmaAddress()
		amount := executorTestAmount()
		txHash := executorTestTxHash()
		fromAddress := executorTestFromAddress()
		depositCommitment := big.NewInt(100)
		depositSalt := big.NewInt(200)
		plChainId := executorTestPnChainId()

		tracer := &testutils.MockTracer{}
		batcher := &executorEnygmaBatcherMock{
			CreateBatchesWithAnonimityFunc: func(
				ctx context.Context,
				resourceId string,
				blockNumber *big.Int,
				txsByChainID map[string][]*types.EnygmaTransferBatchTx,
			) ([]*types.EnygmaTransferBatch, error) {
				return []*types.EnygmaTransferBatch{}, nil
			},
		}
		proofGen := &executorProofGeneratorMock{
			GenerateDepositProofFunc: func(
				ctx context.Context,
				params enygma.DepositProofParams,
			) (*types.EnygmaProofResponse, []*types.Point, []*big.Int, []*big.Int, error) {
				return &types.EnygmaProofResponse{}, []*types.Point{}, []*big.Int{}, []*big.Int{}, nil
			},
		}
		repository := &executorEnygmaHistoryRepositoryMock{}
		keysClient := successKeysClient()
		enygmaClient := &executorEnygmaClientMock{
			GetDvpIntegrationContractAddressFunc: func(_ context.Context, tokenAddress common.Address) (common.Address, error) {
				return common.Address{}, errors.New("failed to get integration address")
			},
		}
		integrationClient := &ExecutorDvpIntegrationClientMock{}
		commitmentCalc := &executorCommitmentCalculatorMock{}

		executor := service.NewEnygmaExecutor(
			conf,
			tracer,
			batcher,
			proofGen,
			repository,
			keysClient,
			enygmaClient,
			integrationClient,
			commitmentCalc,
			plChainId,
		)

		err := executor.ExecuteEnygmaDeposit(
			ctx,
			"test-batch-id",
			resourceId,
			amount,
			blockNumber,
			depositCommitment,
			depositSalt,
			fromAddress,
			txHash,
			enygmaAddress,
		)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to get integration address")
	})

	t.Run("returns error when Deposit fails", func(t *testing.T) {
		ctx := context.Background()
		conf := executorTestConfig()
		resourceId := executorTestResourceId()
		blockNumber := executorTestBlockNumber()
		enygmaAddress := executorTestEnygmaAddress()
		amount := executorTestAmount()
		txHash := executorTestTxHash()
		fromAddress := executorTestFromAddress()
		depositCommitment := big.NewInt(100)
		depositSalt := big.NewInt(200)
		plChainId := executorTestPnChainId()

		tracer := &testutils.MockTracer{}
		batcher := &executorEnygmaBatcherMock{
			CreateBatchesWithAnonimityFunc: func(
				ctx context.Context,
				resourceId string,
				blockNumber *big.Int,
				txsByChainID map[string][]*types.EnygmaTransferBatchTx,
			) ([]*types.EnygmaTransferBatch, error) {
				return []*types.EnygmaTransferBatch{}, nil
			},
		}
		proofGen := &executorProofGeneratorMock{
			GenerateDepositProofFunc: func(
				ctx context.Context,
				params enygma.DepositProofParams,
			) (*types.EnygmaProofResponse, []*types.Point, []*big.Int, []*big.Int, error) {
				return &types.EnygmaProofResponse{}, []*types.Point{}, []*big.Int{}, []*big.Int{}, nil
			},
		}
		repository := &executorEnygmaHistoryRepositoryMock{}
		keysClient := successKeysClient()
		enygmaClient := &executorEnygmaClientMock{
			GetDvpIntegrationContractAddressFunc: func(_ context.Context, tokenAddress common.Address) (common.Address, error) {
				return common.HexToAddress("0xintegration1234"), nil
			},
		}
		integrationClient := &ExecutorDvpIntegrationClientMock{
			DepositFunc: func(
				ctx context.Context,
				_ string,
				batches []*types.EnygmaTransferBatch,
				proof *types.EnygmaProofResponse,
				blockNumber *big.Int,
				chainId *big.Int,
				resourceId string,
				amount *big.Int,
				from common.Address,
				sourceTxHash common.Hash,
				dvpIntegrationAddress common.Address,
			) error {
				return errors.New("deposit failed")
			},
		}
		commitmentCalc := &executorCommitmentCalculatorMock{}

		executor := service.NewEnygmaExecutor(
			conf,
			tracer,
			batcher,
			proofGen,
			repository,
			keysClient,
			enygmaClient,
			integrationClient,
			commitmentCalc,
			plChainId,
		)

		err := executor.ExecuteEnygmaDeposit(
			ctx,
			"test-batch-id",
			resourceId,
			amount,
			blockNumber,
			depositCommitment,
			depositSalt,
			fromAddress,
			txHash,
			enygmaAddress,
		)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "deposit failed")
		assert.Empty(t, repository.InsertEnygmaHistoryCalls())
	})

	t.Run("returns error when repository insertion fails", func(t *testing.T) {
		ctx := context.Background()
		conf := executorTestConfig()
		resourceId := executorTestResourceId()
		blockNumber := executorTestBlockNumber()
		enygmaAddress := executorTestEnygmaAddress()
		amount := executorTestAmount()
		txHash := executorTestTxHash()
		fromAddress := executorTestFromAddress()
		depositCommitment := big.NewInt(100)
		depositSalt := big.NewInt(200)
		plChainId := executorTestPnChainId()

		tracer := &testutils.MockTracer{}
		batcher := &executorEnygmaBatcherMock{
			CreateBatchesWithAnonimityFunc: func(
				ctx context.Context,
				resourceId string,
				blockNumber *big.Int,
				txsByChainID map[string][]*types.EnygmaTransferBatchTx,
			) ([]*types.EnygmaTransferBatch, error) {
				return []*types.EnygmaTransferBatch{}, nil
			},
		}
		proofGen := &executorProofGeneratorMock{
			GenerateDepositProofFunc: func(
				ctx context.Context,
				params enygma.DepositProofParams,
			) (*types.EnygmaProofResponse, []*types.Point, []*big.Int, []*big.Int, error) {
				return &types.EnygmaProofResponse{}, []*types.Point{}, []*big.Int{}, []*big.Int{}, nil
			},
		}
		repository := &executorEnygmaHistoryRepositoryMock{
			InsertEnygmaHistoryFunc: func(ctx context.Context, history types.EnygmaHistory) error {
				return errors.New("database error")
			},
		}
		keysClient := successKeysClient()
		enygmaClient := &executorEnygmaClientMock{
			GetDvpIntegrationContractAddressFunc: func(_ context.Context, tokenAddress common.Address) (common.Address, error) {
				return common.HexToAddress("0xintegration1234"), nil
			},
		}
		integrationClient := &ExecutorDvpIntegrationClientMock{
			DepositFunc: func(
				ctx context.Context,
				_ string,
				batches []*types.EnygmaTransferBatch,
				proof *types.EnygmaProofResponse,
				blockNumber *big.Int,
				chainId *big.Int,
				resourceId string,
				amount *big.Int,
				from common.Address,
				sourceTxHash common.Hash,
				dvpIntegrationAddress common.Address,
			) error {
				return nil
			},
		}
		commitmentCalc := &executorCommitmentCalculatorMock{}

		executor := service.NewEnygmaExecutor(
			conf,
			tracer,
			batcher,
			proofGen,
			repository,
			keysClient,
			enygmaClient,
			integrationClient,
			commitmentCalc,
			plChainId,
		)

		err := executor.ExecuteEnygmaDeposit(
			ctx,
			"test-batch-id",
			resourceId,
			amount,
			blockNumber,
			depositCommitment,
			depositSalt,
			fromAddress,
			txHash,
			enygmaAddress,
		)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "database error")
	})
}

func TestEnygmaExecutor_ExecuteEnygmaWithdrawal(t *testing.T) {
	t.Run("successfully executes withdrawal", func(t *testing.T) {
		ctx := context.Background()
		conf := executorTestConfig()
		conf.MaxNumberOfJSDeposits = 2
		resourceId := executorTestResourceId()
		blockNumber := executorTestBlockNumber()
		enygmaAddress := executorTestEnygmaAddress()
		amount := executorTestAmount()
		txHash := executorTestTxHash()
		fromAddress := executorTestFromAddress()
		plChainId := executorTestPnChainId()
		toChainId := big.NewInt(2)
		senderRFactor := big.NewInt(456)
		integrationAddress := common.HexToAddress("0xintegration5678")

		deposits := []*types.DvpDeposit{
			{
				Salt:         big.NewInt(1),
				TokenAmount:  big.NewInt(500),
				TokenAddress: "0xtoken1",
			},
		}

		jsProof := &dvp.ProofReceipt{}

		tracer := &testutils.MockTracer{}
		batcher := &executorEnygmaBatcherMock{
			CreateBatchesWithAnonimityFunc: func(
				ctx context.Context,
				resourceId string,
				blockNumber *big.Int,
				txsByChainID map[string][]*types.EnygmaTransferBatchTx,
			) ([]*types.EnygmaTransferBatch, error) {
				return []*types.EnygmaTransferBatch{
					{ToChainID: plChainId},
					{ToChainID: toChainId},
				}, nil
			},
		}
		proofGen := &executorProofGeneratorMock{
			GenerateWithdrawProofFunc: func(
				ctx context.Context,
				params enygma.WithdrawProofParams,
			) (*types.EnygmaProofResponse, []*types.Point, []*big.Int, []*big.Int, error) {
				assert.Equal(t, conf.MaxNumberOfJSDeposits, len(params.DepositCommitments))
				assert.Equal(t, conf.MaxNumberOfJSDeposits, len(params.DepositSecretKeys))
				assert.Equal(t, conf.MaxNumberOfJSDeposits, len(params.DepositAmounts))
				assert.Equal(t, 2, len(params.Batches))
				assert.Equal(t, resourceId, params.ResourceId)
				assert.Equal(t, amount, params.SenderAmount)
				assert.Equal(t, new(big.Int).SetUint64(blockNumber), params.BlockNumber)
				assert.Equal(t, enygmaAddress, params.TokenAddress)
				return &types.EnygmaProofResponse{}, []*types.Point{}, []*big.Int{
					senderRFactor,
					big.NewInt(789),
				}, []*big.Int{}, nil
			},
		}
		repository := &executorEnygmaHistoryRepositoryMock{
			InsertEnygmaHistoryFunc: func(ctx context.Context, history types.EnygmaHistory) error {
				return nil
			},
		}
		keysClient := successKeysClient()
		enygmaClient := &executorEnygmaClientMock{
			GetDvpIntegrationContractAddressFunc: func(_ context.Context, tokenAddress common.Address) (common.Address, error) {
				return integrationAddress, nil
			},
		}
		integrationClient := &ExecutorDvpIntegrationClientMock{
			WithdrawFunc: func(
				ctx context.Context,
				_ string,
				batches []*types.EnygmaTransferBatch,
				proof *types.EnygmaProofResponse,
				blockNumber *big.Int,
				jsProof *dvp.ProofReceipt,
				chainId *big.Int,
				resourceId string,
				amount *big.Int,
				from common.Address,
				sourceTxHash common.Hash,
				dvpIntegrationAddress common.Address,
			) error {
				return nil
			},
		}
		commitmentCalc := &executorCommitmentCalculatorMock{
			CalculatePaymentCommitmentFunc: func(spendPK, salt, paymentAmount *big.Int, tokenAddress string) (*big.Int, error) {
				return big.NewInt(250), nil
			},
		}

		executor := service.NewEnygmaExecutor(
			conf,
			tracer,
			batcher,
			proofGen,
			repository,
			keysClient,
			enygmaClient,
			integrationClient,
			commitmentCalc,
			plChainId,
		)

		err := executor.ExecuteEnygmaWithdrawal(
			ctx,
			"test-batch-id",
			resourceId,
			amount,
			deposits,
			blockNumber,
			enygmaAddress,
			jsProof,
			fromAddress,
			txHash,
		)

		require.NoError(t, err)
		assert.Len(t, batcher.CreateBatchesWithAnonimityCalls(), 1)
		assert.Len(t, proofGen.GenerateWithdrawProofCalls(), 1)
		assert.Len(t, keysClient.GetPaymentSpendKeyCalls(), 1)
		assert.Len(t, commitmentCalc.CalculatePaymentCommitmentCalls(), 1)
		assert.Len(t, integrationClient.WithdrawCalls(), 1)
		require.Len(t, repository.InsertEnygmaHistoryCalls(), 1)
		history := repository.InsertEnygmaHistoryCalls()[0].History
		assert.Equal(t, types.EnygmaWithdrawFromDvp, history.EventType)
		assert.Equal(t, plChainId, history.FromChainId)
		assert.Equal(t, senderRFactor, history.RFactor)
		assert.Equal(t, amount, history.BalanceChange)
		assert.Equal(t, new(big.Int).SetUint64(blockNumber), history.BlockNumberPrivateHub)
		assert.Equal(t, integrationAddress, integrationClient.WithdrawCalls()[0].DvpIntegrationAddress)
	})

	t.Run("returns error when keys client fails", func(t *testing.T) {
		ctx := context.Background()
		conf := executorTestConfig()
		resourceId := executorTestResourceId()
		blockNumber := executorTestBlockNumber()
		enygmaAddress := executorTestEnygmaAddress()
		amount := executorTestAmount()
		txHash := executorTestTxHash()
		fromAddress := executorTestFromAddress()
		plChainId := executorTestPnChainId()

		deposits := []*types.DvpDeposit{
			{
				Salt:         big.NewInt(1),
				TokenAmount:  big.NewInt(500),
				TokenAddress: "0xtoken1",
			},
		}

		jsProof := &dvp.ProofReceipt{}

		tracer := &testutils.MockTracer{}
		batcher := &executorEnygmaBatcherMock{}
		proofGen := &executorProofGeneratorMock{}
		repository := &executorEnygmaHistoryRepositoryMock{}
		keysClient := &executorKeysClientMock{
			GetPaymentSpendKeyFunc: func(
				ctx context.Context,
				in *keyspb.GetPaymentSpendKeyRequest,
				opts ...grpc.CallOption,
			) (*keyspb.PaymentSpendKeyResponse, error) {
				return nil, errors.New("failed to get key pair")
			},
		}
		enygmaClient := &executorEnygmaClientMock{}
		integrationClient := &ExecutorDvpIntegrationClientMock{}
		commitmentCalc := &executorCommitmentCalculatorMock{}

		executor := service.NewEnygmaExecutor(
			conf,
			tracer,
			batcher,
			proofGen,
			repository,
			keysClient,
			enygmaClient,
			integrationClient,
			commitmentCalc,
			plChainId,
		)

		err := executor.ExecuteEnygmaWithdrawal(
			ctx,
			"test-batch-id",
			resourceId,
			amount,
			deposits,
			blockNumber,
			enygmaAddress,
			jsProof,
			fromAddress,
			txHash,
		)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to get key pair")
	})

	t.Run("returns error when commitment calculation fails", func(t *testing.T) {
		ctx := context.Background()
		conf := executorTestConfig()
		resourceId := executorTestResourceId()
		blockNumber := executorTestBlockNumber()
		enygmaAddress := executorTestEnygmaAddress()
		amount := executorTestAmount()
		txHash := executorTestTxHash()
		fromAddress := executorTestFromAddress()
		plChainId := executorTestPnChainId()

		deposits := []*types.DvpDeposit{
			{
				Salt:         big.NewInt(1),
				TokenAmount:  big.NewInt(500),
				TokenAddress: "0xtoken1",
			},
		}

		jsProof := &dvp.ProofReceipt{}

		tracer := &testutils.MockTracer{}
		batcher := &executorEnygmaBatcherMock{}
		proofGen := &executorProofGeneratorMock{}
		repository := &executorEnygmaHistoryRepositoryMock{}
		keysClient := successKeysClient()
		enygmaClient := &executorEnygmaClientMock{}
		integrationClient := &ExecutorDvpIntegrationClientMock{}
		commitmentCalc := &executorCommitmentCalculatorMock{
			CalculatePaymentCommitmentFunc: func(spendPK, salt, paymentAmount *big.Int, tokenAddress string) (*big.Int, error) {
				return nil, errors.New("commitment calculation failed")
			},
		}

		executor := service.NewEnygmaExecutor(
			conf,
			tracer,
			batcher,
			proofGen,
			repository,
			keysClient,
			enygmaClient,
			integrationClient,
			commitmentCalc,
			plChainId,
		)

		err := executor.ExecuteEnygmaWithdrawal(
			ctx,
			"test-batch-id",
			resourceId,
			amount,
			deposits,
			blockNumber,
			enygmaAddress,
			jsProof,
			fromAddress,
			txHash,
		)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "commitment calculation failed")
	})

	t.Run("returns error when batcher fails", func(t *testing.T) {
		ctx := context.Background()
		conf := executorTestConfig()
		resourceId := executorTestResourceId()
		blockNumber := executorTestBlockNumber()
		enygmaAddress := executorTestEnygmaAddress()
		amount := executorTestAmount()
		txHash := executorTestTxHash()
		fromAddress := executorTestFromAddress()
		plChainId := executorTestPnChainId()

		deposits := []*types.DvpDeposit{
			{
				Salt:         big.NewInt(1),
				TokenAmount:  big.NewInt(500),
				TokenAddress: "0xtoken1",
			},
		}

		jsProof := &dvp.ProofReceipt{}

		tracer := &testutils.MockTracer{}
		batcher := &executorEnygmaBatcherMock{
			CreateBatchesWithAnonimityFunc: func(
				ctx context.Context,
				resourceId string,
				blockNumber *big.Int,
				txsByChainID map[string][]*types.EnygmaTransferBatchTx,
			) ([]*types.EnygmaTransferBatch, error) {
				return nil, errors.New("batching failed")
			},
		}
		proofGen := &executorProofGeneratorMock{}
		repository := &executorEnygmaHistoryRepositoryMock{}
		keysClient := successKeysClient()
		enygmaClient := &executorEnygmaClientMock{}
		integrationClient := &ExecutorDvpIntegrationClientMock{}
		commitmentCalc := &executorCommitmentCalculatorMock{
			CalculatePaymentCommitmentFunc: func(spendPK, salt, paymentAmount *big.Int, tokenAddress string) (*big.Int, error) {
				return big.NewInt(250), nil
			},
		}

		executor := service.NewEnygmaExecutor(
			conf,
			tracer,
			batcher,
			proofGen,
			repository,
			keysClient,
			enygmaClient,
			integrationClient,
			commitmentCalc,
			plChainId,
		)

		err := executor.ExecuteEnygmaWithdrawal(
			ctx,
			"test-batch-id",
			resourceId,
			amount,
			deposits,
			blockNumber,
			enygmaAddress,
			jsProof,
			fromAddress,
			txHash,
		)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "batching failed")
	})

	t.Run("returns error when proof generation fails", func(t *testing.T) {
		ctx := context.Background()
		conf := executorTestConfig()
		resourceId := executorTestResourceId()
		blockNumber := executorTestBlockNumber()
		enygmaAddress := executorTestEnygmaAddress()
		amount := executorTestAmount()
		txHash := executorTestTxHash()
		fromAddress := executorTestFromAddress()
		plChainId := executorTestPnChainId()

		deposits := []*types.DvpDeposit{
			{
				Salt:         big.NewInt(1),
				TokenAmount:  big.NewInt(500),
				TokenAddress: "0xtoken1",
			},
		}

		jsProof := &dvp.ProofReceipt{}

		tracer := &testutils.MockTracer{}
		batcher := &executorEnygmaBatcherMock{
			CreateBatchesWithAnonimityFunc: func(
				ctx context.Context,
				resourceId string,
				blockNumber *big.Int,
				txsByChainID map[string][]*types.EnygmaTransferBatchTx,
			) ([]*types.EnygmaTransferBatch, error) {
				return []*types.EnygmaTransferBatch{}, nil
			},
		}
		proofGen := &executorProofGeneratorMock{
			GenerateWithdrawProofFunc: func(
				ctx context.Context,
				params enygma.WithdrawProofParams,
			) (*types.EnygmaProofResponse, []*types.Point, []*big.Int, []*big.Int, error) {
				return nil, nil, nil, nil, errors.New("proof generation failed")
			},
		}
		repository := &executorEnygmaHistoryRepositoryMock{}
		keysClient := successKeysClient()
		enygmaClient := &executorEnygmaClientMock{}
		integrationClient := &ExecutorDvpIntegrationClientMock{}
		commitmentCalc := &executorCommitmentCalculatorMock{
			CalculatePaymentCommitmentFunc: func(spendPK, salt, paymentAmount *big.Int, tokenAddress string) (*big.Int, error) {
				return big.NewInt(250), nil
			},
		}

		executor := service.NewEnygmaExecutor(
			conf,
			tracer,
			batcher,
			proofGen,
			repository,
			keysClient,
			enygmaClient,
			integrationClient,
			commitmentCalc,
			plChainId,
		)

		err := executor.ExecuteEnygmaWithdrawal(
			ctx,
			"test-batch-id",
			resourceId,
			amount,
			deposits,
			blockNumber,
			enygmaAddress,
			jsProof,
			fromAddress,
			txHash,
		)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "proof generation failed")
	})

	t.Run("returns error when getting integration address fails", func(t *testing.T) {
		ctx := context.Background()
		conf := executorTestConfig()
		resourceId := executorTestResourceId()
		blockNumber := executorTestBlockNumber()
		enygmaAddress := executorTestEnygmaAddress()
		amount := executorTestAmount()
		txHash := executorTestTxHash()
		fromAddress := executorTestFromAddress()
		plChainId := executorTestPnChainId()

		deposits := []*types.DvpDeposit{
			{
				Salt:         big.NewInt(1),
				TokenAmount:  big.NewInt(500),
				TokenAddress: "0xtoken1",
			},
		}
		jsProof := &dvp.ProofReceipt{}

		tracer := &testutils.MockTracer{}
		batcher := &executorEnygmaBatcherMock{
			CreateBatchesWithAnonimityFunc: func(
				ctx context.Context,
				resourceId string,
				blockNumber *big.Int,
				txsByChainID map[string][]*types.EnygmaTransferBatchTx,
			) ([]*types.EnygmaTransferBatch, error) {
				return []*types.EnygmaTransferBatch{}, nil
			},
		}
		proofGen := &executorProofGeneratorMock{
			GenerateWithdrawProofFunc: func(
				ctx context.Context,
				params enygma.WithdrawProofParams,
			) (*types.EnygmaProofResponse, []*types.Point, []*big.Int, []*big.Int, error) {
				return &types.EnygmaProofResponse{}, []*types.Point{}, []*big.Int{}, []*big.Int{}, nil
			},
		}
		repository := &executorEnygmaHistoryRepositoryMock{}
		keysClient := successKeysClient()
		enygmaClient := &executorEnygmaClientMock{
			GetDvpIntegrationContractAddressFunc: func(_ context.Context, tokenAddress common.Address) (common.Address, error) {
				return common.Address{}, errors.New("failed to get integration address")
			},
		}
		integrationClient := &ExecutorDvpIntegrationClientMock{}
		commitmentCalc := &executorCommitmentCalculatorMock{
			CalculatePaymentCommitmentFunc: func(spendPK, salt, paymentAmount *big.Int, tokenAddress string) (*big.Int, error) {
				return big.NewInt(250), nil
			},
		}

		executor := service.NewEnygmaExecutor(
			conf,
			tracer,
			batcher,
			proofGen,
			repository,
			keysClient,
			enygmaClient,
			integrationClient,
			commitmentCalc,
			plChainId,
		)

		err := executor.ExecuteEnygmaWithdrawal(
			ctx,
			"test-batch-id",
			resourceId,
			amount,
			deposits,
			blockNumber,
			enygmaAddress,
			jsProof,
			fromAddress,
			txHash,
		)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to get integration address")
	})

	t.Run("returns error when Withdraw fails", func(t *testing.T) {
		ctx := context.Background()
		conf := executorTestConfig()
		resourceId := executorTestResourceId()
		blockNumber := executorTestBlockNumber()
		enygmaAddress := executorTestEnygmaAddress()
		amount := executorTestAmount()
		txHash := executorTestTxHash()
		fromAddress := executorTestFromAddress()
		plChainId := executorTestPnChainId()

		deposits := []*types.DvpDeposit{
			{
				Salt:         big.NewInt(1),
				TokenAmount:  big.NewInt(500),
				TokenAddress: "0xtoken1",
			},
		}

		jsProof := &dvp.ProofReceipt{}

		tracer := &testutils.MockTracer{}
		batcher := &executorEnygmaBatcherMock{
			CreateBatchesWithAnonimityFunc: func(
				ctx context.Context,
				resourceId string,
				blockNumber *big.Int,
				txsByChainID map[string][]*types.EnygmaTransferBatchTx,
			) ([]*types.EnygmaTransferBatch, error) {
				return []*types.EnygmaTransferBatch{}, nil
			},
		}
		proofGen := &executorProofGeneratorMock{
			GenerateWithdrawProofFunc: func(
				ctx context.Context,
				params enygma.WithdrawProofParams,
			) (*types.EnygmaProofResponse, []*types.Point, []*big.Int, []*big.Int, error) {
				return &types.EnygmaProofResponse{}, []*types.Point{}, []*big.Int{}, []*big.Int{}, nil
			},
		}
		repository := &executorEnygmaHistoryRepositoryMock{}
		keysClient := successKeysClient()
		enygmaClient := &executorEnygmaClientMock{
			GetDvpIntegrationContractAddressFunc: func(_ context.Context, tokenAddress common.Address) (common.Address, error) {
				return common.HexToAddress("0xintegration5678"), nil
			},
		}
		integrationClient := &ExecutorDvpIntegrationClientMock{
			WithdrawFunc: func(
				ctx context.Context,
				_ string,
				batches []*types.EnygmaTransferBatch,
				proof *types.EnygmaProofResponse,
				blockNumber *big.Int,
				jsProof *dvp.ProofReceipt,
				chainId *big.Int,
				resourceId string,
				amount *big.Int,
				from common.Address,
				sourceTxHash common.Hash,
				dvpIntegrationAddress common.Address,
			) error {
				return errors.New("withdrawal failed")
			},
		}
		commitmentCalc := &executorCommitmentCalculatorMock{
			CalculatePaymentCommitmentFunc: func(spendPK, salt, paymentAmount *big.Int, tokenAddress string) (*big.Int, error) {
				return big.NewInt(250), nil
			},
		}

		executor := service.NewEnygmaExecutor(
			conf,
			tracer,
			batcher,
			proofGen,
			repository,
			keysClient,
			enygmaClient,
			integrationClient,
			commitmentCalc,
			plChainId,
		)

		err := executor.ExecuteEnygmaWithdrawal(
			ctx,
			"test-batch-id",
			resourceId,
			amount,
			deposits,
			blockNumber,
			enygmaAddress,
			jsProof,
			fromAddress,
			txHash,
		)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "withdrawal failed")
		assert.Empty(t, repository.InsertEnygmaHistoryCalls())
	})

	t.Run("returns error when repository insertion fails", func(t *testing.T) {
		ctx := context.Background()
		conf := executorTestConfig()
		resourceId := executorTestResourceId()
		blockNumber := executorTestBlockNumber()
		enygmaAddress := executorTestEnygmaAddress()
		amount := executorTestAmount()
		txHash := executorTestTxHash()
		fromAddress := executorTestFromAddress()
		plChainId := executorTestPnChainId()

		deposits := []*types.DvpDeposit{
			{
				Salt:         big.NewInt(1),
				TokenAmount:  big.NewInt(500),
				TokenAddress: "0xtoken1",
			},
		}

		jsProof := &dvp.ProofReceipt{}

		tracer := &testutils.MockTracer{}
		batcher := &executorEnygmaBatcherMock{
			CreateBatchesWithAnonimityFunc: func(
				ctx context.Context,
				resourceId string,
				blockNumber *big.Int,
				txsByChainID map[string][]*types.EnygmaTransferBatchTx,
			) ([]*types.EnygmaTransferBatch, error) {
				return []*types.EnygmaTransferBatch{}, nil
			},
		}
		proofGen := &executorProofGeneratorMock{
			GenerateWithdrawProofFunc: func(
				ctx context.Context,
				params enygma.WithdrawProofParams,
			) (*types.EnygmaProofResponse, []*types.Point, []*big.Int, []*big.Int, error) {
				return &types.EnygmaProofResponse{}, []*types.Point{}, []*big.Int{}, []*big.Int{}, nil
			},
		}
		repository := &executorEnygmaHistoryRepositoryMock{
			InsertEnygmaHistoryFunc: func(ctx context.Context, history types.EnygmaHistory) error {
				return errors.New("database error")
			},
		}
		keysClient := successKeysClient()
		enygmaClient := &executorEnygmaClientMock{
			GetDvpIntegrationContractAddressFunc: func(_ context.Context, tokenAddress common.Address) (common.Address, error) {
				return common.HexToAddress("0xintegration5678"), nil
			},
		}
		integrationClient := &ExecutorDvpIntegrationClientMock{
			WithdrawFunc: func(
				ctx context.Context,
				_ string,
				batches []*types.EnygmaTransferBatch,
				proof *types.EnygmaProofResponse,
				blockNumber *big.Int,
				jsProof *dvp.ProofReceipt,
				chainId *big.Int,
				resourceId string,
				amount *big.Int,
				from common.Address,
				sourceTxHash common.Hash,
				dvpIntegrationAddress common.Address,
			) error {
				return nil
			},
		}
		commitmentCalc := &executorCommitmentCalculatorMock{
			CalculatePaymentCommitmentFunc: func(spendPK, salt, paymentAmount *big.Int, tokenAddress string) (*big.Int, error) {
				return big.NewInt(250), nil
			},
		}

		executor := service.NewEnygmaExecutor(
			conf,
			tracer,
			batcher,
			proofGen,
			repository,
			keysClient,
			enygmaClient,
			integrationClient,
			commitmentCalc,
			plChainId,
		)

		err := executor.ExecuteEnygmaWithdrawal(
			ctx,
			"test-batch-id",
			resourceId,
			amount,
			deposits,
			blockNumber,
			enygmaAddress,
			jsProof,
			fromAddress,
			txHash,
		)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "database error")
	})
}
