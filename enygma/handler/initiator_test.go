package handler_test

import (
	"context"
	"crypto/mlkem"
	"encoding/hex"
	"errors"
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/raylsnetwork/rayls-sovereign-relayer/cts/gen/proto/keys"
	"github.com/raylsnetwork/rayls-sovereign-relayer/dvp"
	"github.com/raylsnetwork/rayls-sovereign-relayer/enygma/handler"
	"github.com/raylsnetwork/rayls-sovereign-relayer/types"
	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
)

func testViewPublicKey() []byte {
	dk, err := mlkem.GenerateKey768()
	if err != nil {
		panic(err)
	}
	return dk.EncapsulationKey().Bytes()
}

func setupInitiator(deps *initiatorDeps) *handler.Initiator {
	return handler.NewInitiator(
		handler.InitiatorConfig{
			DefaultContextTimeout: 5 * time.Second,
			ViewPublicKey:         new(big.Int),
		},
		deps.KeysClient,
		deps.EnygmaHandlerClient,
		deps.CCEndpointClient,
		deps.PLEndpointClient,
		deps.EnygmaHistoryRepository,
		deps.DvpDepositRepository,
		big.NewInt(1),
		deps.DepositFinder,
		deps.CommitmentCalculator,
		deps.DepositConsolidator,
		deps.DvpProofGenerator,
		deps.Tracer,
		deps.RetryService,
		deps.Executor,
		deps.Finalization,
		deps.CreationService,
		deps.TxManager,
	)
}

type initiatorDeps struct {
	KeysClient              *initiatorKeysClientMock
	EnygmaHandlerClient     *initiatorEnygmaHandlerClientMock
	CCEndpointClient        *initiatorEndpointClientMock
	PLEndpointClient        *initiatorEndpointClientMock
	EnygmaHistoryRepository *initiatorEnygmaHistoryRepositoryMock
	DvpDepositRepository    *initiatorDvpDepositRepositoryMock
	DepositFinder           *initiatorDepositFinderMock
	CommitmentCalculator    *initiatorCommitmentCalculatorMock
	DepositConsolidator     *initiatorDepositConsolidatorMock
	DvpProofGenerator       *initiatorDvpProofGeneratorMock
	Tracer                  *initiatorTracerMock
	RetryService            *initiatorRetryServiceMock
	Executor                *initiatorEnygmaExecutorMock
	Finalization            *initiatorEnygmaFinalizationServiceMock
	CreationService         *initiatorEnygmaCreationServiceMock
	TxManager               *initiatorTxManagerMock
}

func setupInitiatorDeps() *initiatorDeps {
	noopSpan := trace.SpanFromContext(context.Background())
	defaultViewPubKeyHex := hex.EncodeToString(testViewPublicKey())
	return &initiatorDeps{
		KeysClient: &initiatorKeysClientMock{
			GetViewPublicKeyFunc: func(ctx context.Context, in *keys.GetViewPublicKeyRequest, opts ...grpc.CallOption) (*keys.GetViewPublicKeyResponse, error) {
				return &keys.GetViewPublicKeyResponse{
					PublicKey: defaultViewPubKeyHex,
				}, nil
			},
		},
		EnygmaHandlerClient:     &initiatorEnygmaHandlerClientMock{},
		CCEndpointClient:        &initiatorEndpointClientMock{},
		PLEndpointClient:        &initiatorEndpointClientMock{},
		EnygmaHistoryRepository: &initiatorEnygmaHistoryRepositoryMock{},
		DvpDepositRepository:    &initiatorDvpDepositRepositoryMock{},
		DepositFinder:           &initiatorDepositFinderMock{},
		CommitmentCalculator:    &initiatorCommitmentCalculatorMock{},
		DepositConsolidator:     &initiatorDepositConsolidatorMock{},
		DvpProofGenerator:       &initiatorDvpProofGeneratorMock{},
		Tracer: &initiatorTracerMock{
			StartFunc: func(ctx context.Context, spanName string, opts ...trace.SpanStartOption) (context.Context, trace.Span) {
				return ctx, noopSpan
			},
		},
		RetryService:    &initiatorRetryServiceMock{},
		Executor:        &initiatorEnygmaExecutorMock{},
		Finalization:    &initiatorEnygmaFinalizationServiceMock{},
		CreationService: &initiatorEnygmaCreationServiceMock{},
		TxManager: &initiatorTxManagerMock{
			WithTransactionFunc: func(ctx context.Context, fn func(ctx context.Context) error) error {
				return fn(ctx)
			},
		},
	}
}

func createTestResourceId() string {
	return "0000000000000000000000000000000000000000000000000000000000000001"
}

func createTestPaymentSpendKey() *keys.PaymentSpendKeyResponse {
	return &keys.PaymentSpendKeyResponse{
		PublicKey: big.NewInt(12345).Bytes(),
		SecretKey: big.NewInt(67890).Bytes(),
	}
}

func createTestDvpDeposit(tokenAddress string, amount *big.Int) *types.DvpDeposit {
	return &types.DvpDeposit{
		UserAddress:  "0x1111111111111111111111111111111111111111",
		Salt:         big.NewInt(12345),
		TokenAmount:  amount,
		Commitment:   big.NewInt(99999),
		TokenAddress: tokenAddress,
		TokenID:      "",
		TokenType:    types.DvpEnygma,
		Status:       types.DvpDepositPending,
	}
}

func createTestProofReceipt() *dvp.ProofReceipt {
	return &dvp.ProofReceipt{
		Nullifiers: []*big.Int{big.NewInt(111), big.NewInt(222)},
	}
}

func TestHandleEnygmaCreation(t *testing.T) {
	t.Run("Successful creation with initial supply", func(t *testing.T) {
		deps := setupInitiatorDeps()

		deps.PLEndpointClient.GetResourceAddressFunc = func(_ context.Context, resourceId string) (common.Address, error) {
			return common.HexToAddress("0x1234567890123456789012345678901234567890"), nil
		}
		deps.CreationService.CreateEnygmaFunc = func(ctx context.Context, resourceId string, fromChainId *big.Int, blockNumberPrivateHub *big.Int) error {
			return nil
		}
		deps.CCEndpointClient.GetResourceAddressFunc = func(_ context.Context, resourceId string) (common.Address, error) {
			return common.HexToAddress("0x1234567890123456789012345678901234567890"), nil
		}
		deps.RetryService.RetryOperationFunc = func(ctx context.Context, operationName string, maxRetries int, blockNumber uint64, executeOperation func(ctx context.Context, nextBlockNumber uint64) error) (uint64, error) {
			err := executeOperation(ctx, blockNumber)
			return blockNumber, err
		}
		deps.Executor.ExecuteEnygmaSupplyUpdateFunc = func(ctx context.Context, _ string, resourceId string, blockNumber uint64, batch types.EnygmaSupplyUpdate, enygmaAddress common.Address) error {
			return nil
		}

		supply := big.NewInt(1000)
		init := setupInitiator(deps)

		blockNum, err := init.HandleEnygmaCreation(context.Background(), "test-event-id", createTestResourceId(), 100, supply)

		assert.NoError(t, err)
		assert.Equal(t, uint64(100), blockNum)
		assert.Equal(t, 1, len(deps.CreationService.CreateEnygmaCalls()))
		assert.Equal(t, 1, len(deps.CCEndpointClient.GetResourceAddressCalls()))
		assert.Equal(t, 1, len(deps.RetryService.RetryOperationCalls()))
		assert.Equal(t, 1, len(deps.Executor.ExecuteEnygmaSupplyUpdateCalls()))
		// Finalization is no longer called inside HandleEnygmaCreation
		assert.Equal(t, 0, len(deps.Finalization.ExecuteFinalizationCalls()))
	})

	t.Run("Zero initial supply - early exit", func(t *testing.T) {
		deps := setupInitiatorDeps()
		deps.PLEndpointClient.GetResourceAddressFunc = func(_ context.Context, resourceId string) (common.Address, error) {
			return common.HexToAddress("0x1234567890123456789012345678901234567890"), nil
		}
		deps.CreationService.CreateEnygmaFunc = func(ctx context.Context, resourceId string, fromChainId *big.Int, blockNumberPrivateHub *big.Int) error {
			return nil
		}

		supply := big.NewInt(0)
		init := setupInitiator(deps)

		blockNum, err := init.HandleEnygmaCreation(context.Background(), "test-event-id", createTestResourceId(), 100, supply)

		assert.NoError(t, err)
		assert.Equal(t, uint64(0), blockNum)
		assert.Equal(t, 1, len(deps.CreationService.CreateEnygmaCalls()))
		// Verify that PNH resource address was not called for zero supply
		assert.Equal(t, 0, len(deps.CCEndpointClient.GetResourceAddressCalls()))
	})

	t.Run("returns error if failed to retrieve PNH resource address", func(t *testing.T) {
		deps := setupInitiatorDeps()
		deps.PLEndpointClient.GetResourceAddressFunc = func(_ context.Context, resourceId string) (common.Address, error) {
			return common.HexToAddress("0x1234567890123456789012345678901234567890"), nil
		}
		deps.CreationService.CreateEnygmaFunc = func(ctx context.Context, resourceId string, fromChainId *big.Int, blockNumberPrivateHub *big.Int) error {
			return nil
		}
		deps.CCEndpointClient.GetResourceAddressFunc = func(_ context.Context, resourceId string) (common.Address, error) {
			return common.Address{}, errors.New("endpoint error")
		}

		init := setupInitiator(deps)

		blockNum, err := init.HandleEnygmaCreation(context.Background(), "test-event-id", createTestResourceId(), 100, big.NewInt(1000))

		assert.Error(t, err)
		assert.Equal(t, uint64(0), blockNum)
		assert.Equal(t, 1, len(deps.CCEndpointClient.GetResourceAddressCalls()))
		// Verify unused mocks were not called
		assert.Equal(t, 0, len(deps.RetryService.RetryOperationCalls()))
	})

	t.Run("returns error if creation service fails", func(t *testing.T) {
		deps := setupInitiatorDeps()
		deps.PLEndpointClient.GetResourceAddressFunc = func(_ context.Context, resourceId string) (common.Address, error) {
			return common.HexToAddress("0x1234567890123456789012345678901234567890"), nil
		}
		deps.CreationService.CreateEnygmaFunc = func(ctx context.Context, resourceId string, fromChainId *big.Int, blockNumberPrivateHub *big.Int) error {
			return errors.New("creation failed")
		}

		init := setupInitiator(deps)

		blockNum, err := init.HandleEnygmaCreation(context.Background(), "test-event-id", createTestResourceId(), 100, big.NewInt(1000))

		assert.Error(t, err)
		assert.Equal(t, uint64(0), blockNum)
		assert.Equal(t, 1, len(deps.CreationService.CreateEnygmaCalls()))
	})

	t.Run("returns error if failed to retry operation", func(t *testing.T) {
		deps := setupInitiatorDeps()

		deps.PLEndpointClient.GetResourceAddressFunc = func(_ context.Context, resourceId string) (common.Address, error) {
			return common.HexToAddress("0x1234567890123456789012345678901234567890"), nil
		}
		deps.CreationService.CreateEnygmaFunc = func(ctx context.Context, resourceId string, fromChainId *big.Int, blockNumberPrivateHub *big.Int) error {
			return nil
		}
		deps.CCEndpointClient.GetResourceAddressFunc = func(_ context.Context, resourceId string) (common.Address, error) {
			return common.HexToAddress("0x1234567890123456789012345678901234567890"), nil
		}
		deps.RetryService.RetryOperationFunc = func(ctx context.Context, operationName string, maxRetries int, blockNumber uint64, executeOperation func(ctx context.Context, nextBlockNumber uint64) error) (uint64, error) {
			return 0, errors.New("retry failed")
		}

		init := setupInitiator(deps)

		blockNum, err := init.HandleEnygmaCreation(context.Background(), "test-event-id", createTestResourceId(), 100, big.NewInt(1000))

		assert.Error(t, err)
		assert.Equal(t, uint64(0), blockNum)
		assert.Equal(t, 1, len(deps.CreationService.CreateEnygmaCalls()))
		assert.Equal(t, 1, len(deps.CCEndpointClient.GetResourceAddressCalls()))
		assert.Equal(t, 1, len(deps.RetryService.RetryOperationCalls()))
		// Verify that executor was not called
		assert.Equal(t, 0, len(deps.Executor.ExecuteEnygmaSupplyUpdateCalls()))
	})

	t.Run("returns block number for caller to handle finalization", func(t *testing.T) {
		deps := setupInitiatorDeps()

		deps.PLEndpointClient.GetResourceAddressFunc = func(_ context.Context, resourceId string) (common.Address, error) {
			return common.HexToAddress("0x1234567890123456789012345678901234567890"), nil
		}
		deps.CreationService.CreateEnygmaFunc = func(ctx context.Context, resourceId string, fromChainId *big.Int, blockNumberPrivateHub *big.Int) error {
			return nil
		}
		deps.CCEndpointClient.GetResourceAddressFunc = func(_ context.Context, resourceId string) (common.Address, error) {
			return common.HexToAddress("0x1234567890123456789012345678901234567890"), nil
		}
		wantBlockNumber := uint64(100)
		deps.RetryService.RetryOperationFunc = func(ctx context.Context, operationName string, maxRetries int, blockNumber uint64, executeOperation func(ctx context.Context, nextBlockNumber uint64) error) (uint64, error) {
			err := executeOperation(ctx, blockNumber)
			return blockNumber, err
		}
		deps.Executor.ExecuteEnygmaSupplyUpdateFunc = func(ctx context.Context, _ string, resourceId string, blockNumber uint64, batch types.EnygmaSupplyUpdate, enygmaAddress common.Address) error {
			return nil
		}

		init := setupInitiator(deps)

		blockNum, err := init.HandleEnygmaCreation(
			context.Background(),
			"test-event-id",
			createTestResourceId(),
			wantBlockNumber,
			big.NewInt(1000),
		)

		assert.NoError(t, err)
		assert.Equal(t, wantBlockNumber, blockNum)
		assert.Equal(t, 1, len(deps.CreationService.CreateEnygmaCalls()))
		assert.Equal(t, 1, len(deps.CCEndpointClient.GetResourceAddressCalls()))
		assert.Equal(t, 1, len(deps.RetryService.RetryOperationCalls()))
		assert.Equal(t, 1, len(deps.Executor.ExecuteEnygmaSupplyUpdateCalls()))
		// Finalization is no longer called inside HandleEnygmaCreation — orchestrator handles it
		assert.Equal(t, 0, len(deps.Finalization.ExecuteFinalizationCalls()))
	})
}

func TestHandleEnygmaSupplyUpdates(t *testing.T) {
	t.Run("Successful supply update", func(t *testing.T) {
		deps := setupInitiatorDeps()
		deps.CCEndpointClient.GetResourceAddressFunc = func(_ context.Context, resourceId string) (common.Address, error) {
			return common.HexToAddress("0x1234567890123456789012345678901234567890"), nil
		}
		deps.RetryService.RetryOperationFunc = func(ctx context.Context, operationName string, maxRetries int, blockNumber uint64, executeOperation func(ctx context.Context, nextBlockNumber uint64) error) (uint64, error) {
			err := executeOperation(ctx, blockNumber)
			return blockNumber, err
		}
		deps.Executor.ExecuteEnygmaSupplyUpdateFunc = func(ctx context.Context, _ string, resourceId string, blockNumber uint64, batch types.EnygmaSupplyUpdate, enygmaAddress common.Address) error {
			return nil
		}

		init := setupInitiator(deps)

		batch := types.EnygmaSupplyUpdate{
			Amount: big.NewInt(100),
			Type:   types.EnygmaMint,
		}
		newBlockNumber, err := init.HandleEnygmaSupplyUpdates(context.Background(), "test-event-id", createTestResourceId(), 100, batch)

		assert.NoError(t, err)
		assert.Equal(t, uint64(100), newBlockNumber)
		assert.Equal(t, 1, len(deps.CCEndpointClient.GetResourceAddressCalls()))
		assert.Equal(t, 1, len(deps.RetryService.RetryOperationCalls()))
	})

	t.Run("returns error if zero amount supply update", func(t *testing.T) {
		deps := setupInitiatorDeps()
		init := setupInitiator(deps)

		batch := types.EnygmaSupplyUpdate{
			Amount: big.NewInt(0),
			Type:   types.EnygmaMint,
		}
		newBlockNumber, err := init.HandleEnygmaSupplyUpdates(context.Background(), "test-event-id", createTestResourceId(), 100, batch)

		assert.Error(t, err)
		assert.Equal(t, handler.ErrNoSupplyChanges, err)
		assert.Equal(t, uint64(100), newBlockNumber)
		assert.Equal(t, 0, len(deps.CCEndpointClient.GetResourceAddressCalls()))
	})

	t.Run("returns error if failed to retrieve resource address", func(t *testing.T) {
		deps := setupInitiatorDeps()
		deps.CCEndpointClient.GetResourceAddressFunc = func(_ context.Context, resourceId string) (common.Address, error) {
			return common.Address{}, errors.New("endpoint error")
		}

		init := setupInitiator(deps)

		batch := types.EnygmaSupplyUpdate{
			Amount: big.NewInt(100),
			Type:   types.EnygmaMint,
		}
		_, err := init.HandleEnygmaSupplyUpdates(context.Background(), "test-event-id", createTestResourceId(), 100, batch)

		assert.Error(t, err)
		assert.Equal(t, 1, len(deps.CCEndpointClient.GetResourceAddressCalls()))
	})

	t.Run("returns error if failed to retry operation", func(t *testing.T) {
		deps := setupInitiatorDeps()

		deps.CreationService.CreateEnygmaFunc = func(ctx context.Context, resourceId string, fromChainId *big.Int, blockNumberPrivateHub *big.Int) error {
			return nil
		}
		deps.CCEndpointClient.GetResourceAddressFunc = func(_ context.Context, resourceId string) (common.Address, error) {
			return common.HexToAddress("0x1234567890123456789012345678901234567890"), nil
		}
		deps.RetryService.RetryOperationFunc = func(ctx context.Context, operationName string, maxRetries int, blockNumber uint64, executeOperation func(ctx context.Context, nextBlockNumber uint64) error) (uint64, error) {
			assert.Equal(t, "supply update", operationName)
			return 0, errors.New("retry failed")
		}

		init := setupInitiator(deps)

		batch := types.EnygmaSupplyUpdate{
			Amount: big.NewInt(100),
			Type:   types.EnygmaMint,
		}
		_, err := init.HandleEnygmaSupplyUpdates(context.Background(), "test-event-id", createTestResourceId(), 100, batch)

		assert.Error(t, err)
		assert.Equal(t, 1, len(deps.CCEndpointClient.GetResourceAddressCalls()))
		assert.Equal(t, 1, len(deps.RetryService.RetryOperationCalls()))
	})
}

func TestHandleSendEnygmaCrossTransfer(t *testing.T) {
	t.Run("Successful cross-transfer", func(t *testing.T) {
		deps := setupInitiatorDeps()
		deps.CCEndpointClient.GetResourceAddressFunc = func(_ context.Context, resourceId string) (common.Address, error) {
			return common.HexToAddress("0x1234567890123456789012345678901234567890"), nil
		}
		deps.RetryService.RetryOperationFunc = func(ctx context.Context, operationName string, maxRetries int, blockNumber uint64, executeOperation func(ctx context.Context, nextBlockNumber uint64) error) (uint64, error) {
			assert.Equal(t, "cross transfer", operationName)
			err := executeOperation(ctx, blockNumber)
			return blockNumber, err
		}
		deps.Executor.ExecuteEnygmaCrossTransferFunc = func(ctx context.Context, _ string, blockNumber uint64, resourceId string, batch map[string][]*types.EnygmaTransferBatchTx, enygmaAddress common.Address) error {
			return nil
		}

		init := setupInitiator(deps)
		batch := make(map[string][]*types.EnygmaTransferBatchTx)

		newBlockNumber, err := init.HandleEnygmaCrossTransfer(context.Background(), "test-event-id", createTestResourceId(), 100, batch)

		assert.NoError(t, err)
		assert.Equal(t, uint64(100), newBlockNumber)
		assert.Equal(t, 1, len(deps.CCEndpointClient.GetResourceAddressCalls()))
		assert.Equal(t, 1, len(deps.RetryService.RetryOperationCalls()))
	})

	t.Run("returns error if failed to retrieve resource address", func(t *testing.T) {
		deps := setupInitiatorDeps()
		deps.CCEndpointClient.GetResourceAddressFunc = func(_ context.Context, resourceId string) (common.Address, error) {
			return common.Address{}, errors.New("endpoint error")
		}

		init := setupInitiator(deps)
		batch := make(map[string][]*types.EnygmaTransferBatchTx)

		_, err := init.HandleEnygmaCrossTransfer(context.Background(), "test-event-id", createTestResourceId(), 100, batch)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "endpoint error")
		assert.Equal(t, 0, len(deps.RetryService.RetryOperationCalls()))
	})

	t.Run("returns error if failed to retry operation", func(t *testing.T) {
		deps := setupInitiatorDeps()

		deps.CreationService.CreateEnygmaFunc = func(ctx context.Context, resourceId string, fromChainId *big.Int, blockNumberPrivateHub *big.Int) error {
			return nil
		}
		deps.CCEndpointClient.GetResourceAddressFunc = func(_ context.Context, resourceId string) (common.Address, error) {
			return common.HexToAddress("0x1234567890123456789012345678901234567890"), nil
		}
		deps.RetryService.RetryOperationFunc = func(ctx context.Context, operationName string, maxRetries int, blockNumber uint64, executeOperation func(ctx context.Context, nextBlockNumber uint64) error) (uint64, error) {
			return 0, errors.New("retry failed")
		}

		init := setupInitiator(deps)

		batch := make(map[string][]*types.EnygmaTransferBatchTx)
		_, err := init.HandleEnygmaCrossTransfer(context.Background(), "test-event-id", createTestResourceId(), 100, batch)

		assert.Error(t, err)
		assert.Equal(t, 1, len(deps.CCEndpointClient.GetResourceAddressCalls()))
		assert.Equal(t, 1, len(deps.RetryService.RetryOperationCalls()))
	})
}

func TestHandleEnygmaDeposit(t *testing.T) {
	t.Run("Successful deposit", func(t *testing.T) {
		deps := setupInitiatorDeps()
		keypair := createTestPaymentSpendKey()

		deps.CCEndpointClient.GetResourceAddressFunc = func(_ context.Context, resourceId string) (common.Address, error) {
			return common.HexToAddress("0x1234567890123456789012345678901234567890"), nil
		}
		deps.KeysClient.GetPaymentSpendKeyFunc = func(ctx context.Context, in *keys.GetPaymentSpendKeyRequest, opts ...grpc.CallOption) (*keys.PaymentSpendKeyResponse, error) {
			return keypair, nil
		}
		deps.CommitmentCalculator.CalculatePaymentCommitmentFunc = func(spendPK, salt, paymentAmount *big.Int, tokenAddress string) (*big.Int, error) {
			return big.NewInt(88888), nil
		}
		deps.DvpDepositRepository.CreateDepositFunc = func(ctx context.Context, deposit *types.DvpDeposit) error {
			return nil
		}
		deps.RetryService.RetryOperationFunc = func(ctx context.Context, operationName string, maxRetries int, blockNumber uint64, executeOperation func(ctx context.Context, nextBlockNumber uint64) error) (uint64, error) {
			err := executeOperation(ctx, blockNumber)
			return blockNumber, err
		}
		deps.Executor.ExecuteEnygmaDepositFunc = func(ctx context.Context, _ string, resourceId string, amount *big.Int, blockNumber uint64, commitment *big.Int, keypairCoin *big.Int, from common.Address, txHash common.Hash, enygmaAddress common.Address) error {
			return nil
		}

		init := setupInitiator(deps)

		referenceId := [32]byte{}

		newBlockNumber, err := init.HandleEnygmaDeposit(
			context.Background(),
			"test-event-id",
			100,
			createTestResourceId(),
			referenceId,
			common.HexToAddress("0xAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA"),
			big.NewInt(1000),
			common.HexToHash("0xabcd"),
		)

		assert.NoError(t, err)
		assert.Equal(t, uint64(100), newBlockNumber)
		assert.Equal(t, 1, len(deps.CCEndpointClient.GetResourceAddressCalls()))
		assert.Equal(t, 1, len(deps.KeysClient.GetPaymentSpendKeyCalls()))
		assert.Equal(t, 1, len(deps.CommitmentCalculator.CalculatePaymentCommitmentCalls()))
		assert.Equal(t, 1, len(deps.DvpDepositRepository.CreateDepositCalls()))
		assert.Equal(t, 1, len(deps.Executor.ExecuteEnygmaDepositCalls()))
	})

	t.Run("returns error if failed to retrieve resource address", func(t *testing.T) {
		deps := setupInitiatorDeps()
		deps.CCEndpointClient.GetResourceAddressFunc = func(_ context.Context, resourceId string) (common.Address, error) {
			return common.Address{}, errors.New("endpoint error")
		}

		init := setupInitiator(deps)

		referenceId := [32]byte{}

		_, err := init.HandleEnygmaDeposit(
			context.Background(),
			"test-event-id",
			100,
			createTestResourceId(),
			referenceId,
			common.HexToAddress("0xAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA"),
			big.NewInt(1000),
			common.HexToHash("0xabcd"),
		)

		assert.Error(t, err)
		assert.Equal(t, 0, len(deps.KeysClient.GetPaymentSpendKeyCalls()))
	})

	t.Run("returns error if failed to create keypair", func(t *testing.T) {
		deps := setupInitiatorDeps()
		deps.CCEndpointClient.GetResourceAddressFunc = func(_ context.Context, resourceId string) (common.Address, error) {
			return common.HexToAddress("0x1234567890123456789012345678901234567890"), nil
		}
		deps.KeysClient.GetPaymentSpendKeyFunc = func(ctx context.Context, in *keys.GetPaymentSpendKeyRequest, opts ...grpc.CallOption) (*keys.PaymentSpendKeyResponse, error) {
			return &keys.PaymentSpendKeyResponse{}, errors.New("keypair creation failed")
		}

		init := setupInitiator(deps)

		referenceId := [32]byte{}

		_, err := init.HandleEnygmaDeposit(
			context.Background(),
			"test-event-id",
			100,
			createTestResourceId(),
			referenceId,
			common.HexToAddress("0xAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA"),
			big.NewInt(1000),
			common.HexToHash("0xabcd"),
		)

		assert.Error(t, err)
		assert.Equal(t, 1, len(deps.KeysClient.GetPaymentSpendKeyCalls()))
		assert.Equal(t, 0, len(deps.CommitmentCalculator.CalculatePaymentCommitmentCalls()))
	})

	t.Run("returns error if failed to calculate commitment", func(t *testing.T) {
		deps := setupInitiatorDeps()
		keypair := createTestPaymentSpendKey()

		deps.CCEndpointClient.GetResourceAddressFunc = func(_ context.Context, resourceId string) (common.Address, error) {
			return common.HexToAddress("0x1234567890123456789012345678901234567890"), nil
		}
		deps.KeysClient.GetPaymentSpendKeyFunc = func(ctx context.Context, in *keys.GetPaymentSpendKeyRequest, opts ...grpc.CallOption) (*keys.PaymentSpendKeyResponse, error) {
			return keypair, nil
		}
		deps.DvpDepositRepository.CreateDepositFunc = func(ctx context.Context, deposit *types.DvpDeposit) error {
			return nil
		}
		deps.CommitmentCalculator.CalculatePaymentCommitmentFunc = func(spendPK, salt, paymentAmount *big.Int, tokenAddress string) (*big.Int, error) {
			return nil, errors.New("commitment calculation failed")
		}

		init := setupInitiator(deps)

		referenceId := [32]byte{}

		_, err := init.HandleEnygmaDeposit(
			context.Background(),
			"test-event-id",
			100,
			createTestResourceId(),
			referenceId,
			common.HexToAddress("0xAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA"),
			big.NewInt(1000),
			common.HexToHash("0xabcd"),
		)

		assert.Error(t, err)
		assert.Equal(t, 0, len(deps.DvpDepositRepository.CreateDepositCalls()))
	})

	t.Run("returns error if failed to create deposit in database", func(t *testing.T) {
		deps := setupInitiatorDeps()
		keypair := createTestPaymentSpendKey()

		deps.CCEndpointClient.GetResourceAddressFunc = func(_ context.Context, resourceId string) (common.Address, error) {
			return common.HexToAddress("0x1234567890123456789012345678901234567890"), nil
		}
		deps.KeysClient.GetPaymentSpendKeyFunc = func(ctx context.Context, in *keys.GetPaymentSpendKeyRequest, opts ...grpc.CallOption) (*keys.PaymentSpendKeyResponse, error) {
			return keypair, nil
		}
		deps.CommitmentCalculator.CalculatePaymentCommitmentFunc = func(spendPK, salt, paymentAmount *big.Int, tokenAddress string) (*big.Int, error) {
			return big.NewInt(88888), nil
		}
		deps.DvpDepositRepository.CreateDepositFunc = func(ctx context.Context, deposit *types.DvpDeposit) error {
			return errors.New("database error")
		}

		init := setupInitiator(deps)

		referenceId := [32]byte{}

		_, err := init.HandleEnygmaDeposit(
			context.Background(),
			"test-event-id",
			100,
			createTestResourceId(),
			referenceId,
			common.HexToAddress("0xAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA"),
			big.NewInt(1000),
			common.HexToHash("0xabcd"),
		)

		assert.Error(t, err)
		assert.Equal(t, 0, len(deps.RetryService.RetryOperationCalls()))
	})

	t.Run("returns error if failed to retry operation", func(t *testing.T) {
		deps := setupInitiatorDeps()
		keypair := createTestPaymentSpendKey()

		deps.CCEndpointClient.GetResourceAddressFunc = func(_ context.Context, resourceId string) (common.Address, error) {
			return common.HexToAddress("0x1234567890123456789012345678901234567890"), nil
		}
		deps.KeysClient.GetPaymentSpendKeyFunc = func(ctx context.Context, in *keys.GetPaymentSpendKeyRequest, opts ...grpc.CallOption) (*keys.PaymentSpendKeyResponse, error) {
			return keypair, nil
		}
		deps.CommitmentCalculator.CalculatePaymentCommitmentFunc = func(spendPK, salt, paymentAmount *big.Int, tokenAddress string) (*big.Int, error) {
			return big.NewInt(88888), nil
		}
		deps.DvpDepositRepository.CreateDepositFunc = func(ctx context.Context, deposit *types.DvpDeposit) error {
			return nil
		}
		deps.RetryService.RetryOperationFunc = func(ctx context.Context, operationName string, maxRetries int, blockNumber uint64, executeOperation func(ctx context.Context, nextBlockNumber uint64) error) (uint64, error) {
			return 0, errors.New("retry failed")
		}

		init := setupInitiator(deps)

		referenceId := [32]byte{}
		_, err := init.HandleEnygmaDeposit(
			context.Background(),
			"test-event-id",
			100,
			createTestResourceId(),
			referenceId,
			common.HexToAddress("0xAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA"),
			big.NewInt(1000),
			common.HexToHash("0xabcd"),
		)

		assert.Error(t, err)
		assert.Equal(t, 1, len(deps.CCEndpointClient.GetResourceAddressCalls()))
		assert.Equal(t, 1, len(deps.KeysClient.GetPaymentSpendKeyCalls()))
		assert.Equal(t, 1, len(deps.CommitmentCalculator.CalculatePaymentCommitmentCalls()))
		assert.Equal(t, 1, len(deps.DvpDepositRepository.CreateDepositCalls()))
		// Executor is not called since retry fails before executing the operation
		assert.Equal(t, 0, len(deps.Executor.ExecuteEnygmaDepositCalls()))
		assert.Equal(t, 1, len(deps.RetryService.RetryOperationCalls()))
	})
}

func TestHandleEnygmaWithdrawal(t *testing.T) {
	t.Run("Successful withdrawal", func(t *testing.T) {
		deps := setupInitiatorDeps()
		deposit := createTestDvpDeposit("0x1234567890123456789012345678901234567890", big.NewInt(1000))
		proof := createTestProofReceipt()
		keypair := createTestPaymentSpendKey()
		to := common.HexToAddress("0xBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBB")
		amount := big.NewInt(1000)

		deps.CCEndpointClient.GetResourceAddressFunc = func(_ context.Context, resourceId string) (common.Address, error) {
			return common.HexToAddress("0x1234567890123456789012345678901234567890"), nil
		}
		deps.PLEndpointClient.GetResourceAddressFunc = func(_ context.Context, resourceId string) (common.Address, error) {
			return common.HexToAddress("0x1234567890123456789012345678901234567890"), nil
		}
		deps.DepositFinder.FindEnygmaDepositsFunc = func(ctx context.Context, userAddress string, tokenAddress string, paymentAmount *big.Int) ([]*types.DvpDeposit, error) {
			return []*types.DvpDeposit{deposit}, nil
		}
		deps.DepositConsolidator.PrepareDepositsForJSProofFunc = func(ctx context.Context, _ string, sourceViewPublicKey []byte, deposits []*types.DvpDeposit) ([]*types.DvpDeposit, error) {
			return deposits, nil
		}
		deps.KeysClient.GetPaymentSpendKeyFunc = func(ctx context.Context, in *keys.GetPaymentSpendKeyRequest, opts ...grpc.CallOption) (*keys.PaymentSpendKeyResponse, error) {
			return keypair, nil
		}
		deps.DvpProofGenerator.GenerateEnygmaJSProofFunc = func(ctx context.Context, sourceViewPublicKey []byte, nftCommitment *big.Int, destinationPaymentPublicKey *big.Int, destinationSalt *big.Int, paymentAmount *big.Int, tokenAddress string, deposits []*types.DvpDeposit) (*dvp.ProofReceipt, error) {
			return proof, nil
		}
		deps.DvpDepositRepository.BatchUpdateStatusForCommitmentsFunc = func(ctx context.Context, commitments []string, status types.DvpDepositStatus) error {
			return nil
		}
		deps.DvpDepositRepository.BatchUpsertNullifiersFunc = func(ctx context.Context, commitmentNullifierMap map[string]string) error {
			return nil
		}
		deps.RetryService.RetryOperationFunc = func(ctx context.Context, operationName string, maxRetries int, blockNumber uint64, executeOperation func(ctx context.Context, nextBlockNumber uint64) error) (uint64, error) {
			err := executeOperation(ctx, blockNumber)
			return blockNumber, err
		}
		deps.Executor.ExecuteEnygmaWithdrawalFunc = func(ctx context.Context, _ string, resourceId string, amount *big.Int, deposits []*types.DvpDeposit, blockNumber uint64, enygmaAddress common.Address, proof *dvp.ProofReceipt, to common.Address, txHash common.Hash) error {
			return nil
		}
		deps.EnygmaHandlerClient.ReceiveWithdrawFunc = func(ctx context.Context, _ string, tokenAddress common.Address, toAddr common.Address, value *big.Int, referenceId [32]byte) error {
			assert.Equal(t, amount, value)
			assert.Equal(t, to, toAddr)
			return nil
		}

		init := setupInitiator(deps)

		referenceId := [32]byte{}
		newBlockNumber, err := init.HandleEnygmaWithdrawal(
			context.Background(),
			"test-event-id",
			100,
			createTestResourceId(),
			referenceId,
			to,
			amount,
			common.HexToHash("0xabcd"),
		)

		assert.NoError(t, err)
		assert.Equal(t, uint64(100), newBlockNumber)
		assert.Equal(t, 1, len(deps.DepositFinder.FindEnygmaDepositsCalls()))
		assert.Equal(t, 1, len(deps.DepositConsolidator.PrepareDepositsForJSProofCalls()))
		assert.Equal(t, 1, len(deps.DvpProofGenerator.GenerateEnygmaJSProofCalls()))
		assert.Equal(t, 1, len(deps.Executor.ExecuteEnygmaWithdrawalCalls()))
	})

	t.Run("returns error if failed to retrieve resource address", func(t *testing.T) {
		deps := setupInitiatorDeps()
		deps.CCEndpointClient.GetResourceAddressFunc = func(_ context.Context, resourceId string) (common.Address, error) {
			return common.Address{}, errors.New("endpoint error")
		}

		init := setupInitiator(deps)

		referenceId := [32]byte{}
		_, err := init.HandleEnygmaWithdrawal(
			context.Background(),
			"test-event-id",
			100,
			createTestResourceId(),
			referenceId,
			common.HexToAddress("0xBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBB"),
			big.NewInt(1000),
			common.HexToHash("0xabcd"),
		)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "endpoint error")
	})

	t.Run("returns error if no deposits found", func(t *testing.T) {
		deps := setupInitiatorDeps()
		deps.CCEndpointClient.GetResourceAddressFunc = func(_ context.Context, resourceId string) (common.Address, error) {
			return common.HexToAddress("0x1234567890123456789012345678901234567890"), nil
		}
		deps.DepositFinder.FindEnygmaDepositsFunc = func(ctx context.Context, userAddress string, tokenAddress string, paymentAmount *big.Int) ([]*types.DvpDeposit, error) {
			return []*types.DvpDeposit{}, nil
		}
		deps.DepositConsolidator.PrepareDepositsForJSProofFunc = func(ctx context.Context, _ string, sourceViewPublicKey []byte, deposits []*types.DvpDeposit) ([]*types.DvpDeposit, error) {
			return deposits, nil
		}

		init := setupInitiator(deps)

		referenceId := [32]byte{}
		_, err := init.HandleEnygmaWithdrawal(
			context.Background(),
			"test-event-id",
			100,
			createTestResourceId(),
			referenceId,
			common.HexToAddress("0xBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBB"),
			big.NewInt(1000),
			common.HexToHash("0xabcd"),
		)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "no deposits found to withdraw from")
		assert.Equal(t, 0, len(deps.DvpProofGenerator.GenerateEnygmaJSProofCalls()))
	})

	t.Run("returns error if failed to find deposits", func(t *testing.T) {
		deps := setupInitiatorDeps()
		deps.CCEndpointClient.GetResourceAddressFunc = func(_ context.Context, resourceId string) (common.Address, error) {
			return common.HexToAddress("0x1234567890123456789012345678901234567890"), nil
		}
		deps.DepositFinder.FindEnygmaDepositsFunc = func(ctx context.Context, userAddress string, tokenAddress string, paymentAmount *big.Int) ([]*types.DvpDeposit, error) {
			return nil, errors.New("deposit finder error")
		}

		init := setupInitiator(deps)

		referenceId := [32]byte{}
		_, err := init.HandleEnygmaWithdrawal(
			context.Background(),
			"test-event-id",
			100,
			createTestResourceId(),
			referenceId,
			common.HexToAddress("0xBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBB"),
			big.NewInt(1000),
			common.HexToHash("0xabcd"),
		)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "deposit finder error")
	})

	t.Run("returns error if failed to prepare deposits for proof", func(t *testing.T) {
		deps := setupInitiatorDeps()
		deposit := createTestDvpDeposit("0x1234567890123456789012345678901234567890", big.NewInt(1000))

		deps.CCEndpointClient.GetResourceAddressFunc = func(_ context.Context, resourceId string) (common.Address, error) {
			return common.HexToAddress("0x1234567890123456789012345678901234567890"), nil
		}
		deps.DepositFinder.FindEnygmaDepositsFunc = func(ctx context.Context, userAddress string, tokenAddress string, paymentAmount *big.Int) ([]*types.DvpDeposit, error) {
			return []*types.DvpDeposit{deposit}, nil
		}
		deps.DepositConsolidator.PrepareDepositsForJSProofFunc = func(ctx context.Context, _ string, sourceViewPublicKey []byte, deposits []*types.DvpDeposit) ([]*types.DvpDeposit, error) {
			return nil, errors.New("prepare deposits failed")
		}

		init := setupInitiator(deps)

		referenceId := [32]byte{}
		_, err := init.HandleEnygmaWithdrawal(
			context.Background(),
			"test-event-id",
			100,
			createTestResourceId(),
			referenceId,
			common.HexToAddress("0xBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBB"),
			big.NewInt(1000),
			common.HexToHash("0xabcd"),
		)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "prepare deposits failed")
		assert.Equal(t, 0, len(deps.DvpProofGenerator.GenerateEnygmaJSProofCalls()))
	})

	t.Run("returns error if failed to generate proof", func(t *testing.T) {
		deps := setupInitiatorDeps()
		deposit := createTestDvpDeposit("0x1234567890123456789012345678901234567890", big.NewInt(1000))
		keypair := createTestPaymentSpendKey()

		deps.CCEndpointClient.GetResourceAddressFunc = func(_ context.Context, resourceId string) (common.Address, error) {
			return common.HexToAddress("0x1234567890123456789012345678901234567890"), nil
		}
		deps.DepositFinder.FindEnygmaDepositsFunc = func(ctx context.Context, userAddress string, tokenAddress string, paymentAmount *big.Int) ([]*types.DvpDeposit, error) {
			return []*types.DvpDeposit{deposit}, nil
		}
		deps.DepositConsolidator.PrepareDepositsForJSProofFunc = func(ctx context.Context, _ string, sourceViewPublicKey []byte, deposits []*types.DvpDeposit) ([]*types.DvpDeposit, error) {
			return deposits, nil
		}
		deps.KeysClient.GetPaymentSpendKeyFunc = func(ctx context.Context, in *keys.GetPaymentSpendKeyRequest, opts ...grpc.CallOption) (*keys.PaymentSpendKeyResponse, error) {
			return keypair, nil
		}
		deps.DvpProofGenerator.GenerateEnygmaJSProofFunc = func(ctx context.Context, sourceViewPublicKey []byte, nftCommitment *big.Int, destinationPaymentPublicKey *big.Int, destinationSalt *big.Int, paymentAmount *big.Int, tokenAddress string, deposits []*types.DvpDeposit) (*dvp.ProofReceipt, error) {
			return nil, errors.New("proof generation failed")
		}
		deps.DvpDepositRepository.BatchUpdateStatusForCommitmentsFunc = func(ctx context.Context, commitments []string, status types.DvpDepositStatus) error {
			return nil
		}

		init := setupInitiator(deps)

		referenceId := [32]byte{}
		_, err := init.HandleEnygmaWithdrawal(
			context.Background(),
			"test-event-id",
			100,
			createTestResourceId(),
			referenceId,
			common.HexToAddress("0xBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBB"),
			big.NewInt(1000),
			common.HexToHash("0xabcd"),
		)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "proof generation failed")
		assert.Equal(t, 0, len(deps.Executor.ExecuteEnygmaWithdrawalCalls()))
	})

	t.Run("returns error if failed to update deposit status", func(t *testing.T) {
		deps := setupInitiatorDeps()
		deposit := createTestDvpDeposit("0x1234567890123456789012345678901234567890", big.NewInt(1000))
		keypair := createTestPaymentSpendKey()

		deps.CCEndpointClient.GetResourceAddressFunc = func(_ context.Context, resourceId string) (common.Address, error) {
			return common.HexToAddress("0x1234567890123456789012345678901234567890"), nil
		}
		deps.DepositFinder.FindEnygmaDepositsFunc = func(ctx context.Context, userAddress string, tokenAddress string, paymentAmount *big.Int) ([]*types.DvpDeposit, error) {
			return []*types.DvpDeposit{deposit}, nil
		}
		deps.DepositConsolidator.PrepareDepositsForJSProofFunc = func(ctx context.Context, _ string, sourceViewPublicKey []byte, deposits []*types.DvpDeposit) ([]*types.DvpDeposit, error) {
			return deposits, nil
		}
		deps.KeysClient.GetPaymentSpendKeyFunc = func(ctx context.Context, in *keys.GetPaymentSpendKeyRequest, opts ...grpc.CallOption) (*keys.PaymentSpendKeyResponse, error) {
			return keypair, nil
		}
		deps.DvpProofGenerator.GenerateEnygmaJSProofFunc = func(ctx context.Context, sourceViewPublicKey []byte, nftCommitment *big.Int, publicKey *big.Int, destinationSalt *big.Int, amount *big.Int, tokenAddress string, deposits []*types.DvpDeposit) (*dvp.ProofReceipt, error) {
			return createTestProofReceipt(), nil
		}
		deps.DvpDepositRepository.BatchUpdateStatusForCommitmentsFunc = func(ctx context.Context, commitments []string, status types.DvpDepositStatus) error {
			return errors.New("update status failed")
		}

		init := setupInitiator(deps)

		referenceId := [32]byte{}
		_, err := init.HandleEnygmaWithdrawal(
			context.Background(),
			"test-event-id",
			100,
			createTestResourceId(),
			referenceId,
			common.HexToAddress("0xBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBB"),
			big.NewInt(1000),
			common.HexToHash("0xabcd"),
		)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "update status failed")
	})

	t.Run("returns error if failed to get dvp key pair", func(t *testing.T) {
		deps := setupInitiatorDeps()
		deposit := createTestDvpDeposit("0x1234567890123456789012345678901234567890", big.NewInt(1000))

		deps.CCEndpointClient.GetResourceAddressFunc = func(_ context.Context, resourceId string) (common.Address, error) {
			return common.HexToAddress("0x1234567890123456789012345678901234567890"), nil
		}
		deps.DepositFinder.FindEnygmaDepositsFunc = func(ctx context.Context, userAddress string, tokenAddress string, paymentAmount *big.Int) ([]*types.DvpDeposit, error) {
			return []*types.DvpDeposit{deposit}, nil
		}
		deps.DepositConsolidator.PrepareDepositsForJSProofFunc = func(ctx context.Context, _ string, sourceViewPublicKey []byte, deposits []*types.DvpDeposit) ([]*types.DvpDeposit, error) {
			return deposits, nil
		}
		deps.DvpDepositRepository.BatchUpdateStatusForCommitmentsFunc = func(ctx context.Context, commitments []string, status types.DvpDepositStatus) error {
			return nil
		}
		deps.KeysClient.GetPaymentSpendKeyFunc = func(ctx context.Context, in *keys.GetPaymentSpendKeyRequest, opts ...grpc.CallOption) (*keys.PaymentSpendKeyResponse, error) {
			return &keys.PaymentSpendKeyResponse{}, errors.New("keypair retrieval failed")
		}

		init := setupInitiator(deps)

		referenceId := [32]byte{}
		_, err := init.HandleEnygmaWithdrawal(
			context.Background(),
			"test-event-id",
			100,
			createTestResourceId(),
			referenceId,
			common.HexToAddress("0xBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBB"),
			big.NewInt(1000),
			common.HexToHash("0xabcd"),
		)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "keypair retrieval failed")
		// Verify these were NOT called
		assert.Equal(t, 0, len(deps.DvpProofGenerator.GenerateEnygmaJSProofCalls()))
		assert.Equal(t, 0, len(deps.DvpDepositRepository.BatchUpsertNullifiersCalls()))
		assert.Equal(t, 0, len(deps.RetryService.RetryOperationCalls()))
	})

	t.Run("returns error if failed to upsert nullifiers", func(t *testing.T) {
		deps := setupInitiatorDeps()
		deposit := createTestDvpDeposit("0x1234567890123456789012345678901234567890", big.NewInt(1000))
		proof := createTestProofReceipt()
		keypair := createTestPaymentSpendKey()

		deps.CCEndpointClient.GetResourceAddressFunc = func(_ context.Context, resourceId string) (common.Address, error) {
			return common.HexToAddress("0x1234567890123456789012345678901234567890"), nil
		}
		deps.DepositFinder.FindEnygmaDepositsFunc = func(ctx context.Context, userAddress string, tokenAddress string, paymentAmount *big.Int) ([]*types.DvpDeposit, error) {
			return []*types.DvpDeposit{deposit}, nil
		}
		deps.DepositConsolidator.PrepareDepositsForJSProofFunc = func(ctx context.Context, _ string, sourceViewPublicKey []byte, deposits []*types.DvpDeposit) ([]*types.DvpDeposit, error) {
			return deposits, nil
		}
		deps.DvpDepositRepository.BatchUpdateStatusForCommitmentsFunc = func(ctx context.Context, commitments []string, status types.DvpDepositStatus) error {
			return nil
		}
		deps.KeysClient.GetPaymentSpendKeyFunc = func(ctx context.Context, in *keys.GetPaymentSpendKeyRequest, opts ...grpc.CallOption) (*keys.PaymentSpendKeyResponse, error) {
			return keypair, nil
		}
		deps.DvpProofGenerator.GenerateEnygmaJSProofFunc = func(ctx context.Context, sourceViewPublicKey []byte, nftCommitment *big.Int, destinationPaymentPublicKey *big.Int, destinationSalt *big.Int, paymentAmount *big.Int, tokenAddress string, deposits []*types.DvpDeposit) (*dvp.ProofReceipt, error) {
			return proof, nil
		}
		deps.DvpDepositRepository.BatchUpsertNullifiersFunc = func(ctx context.Context, commitmentNullifierMap map[string]string) error {
			return errors.New("nullifier upsert failed")
		}

		init := setupInitiator(deps)

		referenceId := [32]byte{}
		_, err := init.HandleEnygmaWithdrawal(
			context.Background(),
			"test-event-id",
			100,
			createTestResourceId(),
			referenceId,
			common.HexToAddress("0xBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBB"),
			big.NewInt(1000),
			common.HexToHash("0xabcd"),
		)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "nullifier upsert failed")
		// Verify these were called
		assert.Equal(t, 1, len(deps.DvpProofGenerator.GenerateEnygmaJSProofCalls()))
		assert.Equal(t, 1, len(deps.DvpDepositRepository.BatchUpsertNullifiersCalls()))
		// Verify this was NOT called
		assert.Equal(t, 0, len(deps.RetryService.RetryOperationCalls()))
	})

	t.Run("returns error if failed to retry operation", func(t *testing.T) {
		deps := setupInitiatorDeps()
		deposit := createTestDvpDeposit("0x1234567890123456789012345678901234567890", big.NewInt(1000))
		proof := createTestProofReceipt()
		keypair := createTestPaymentSpendKey()

		deps.CCEndpointClient.GetResourceAddressFunc = func(_ context.Context, resourceId string) (common.Address, error) {
			return common.HexToAddress("0x1234567890123456789012345678901234567890"), nil
		}
		deps.DepositFinder.FindEnygmaDepositsFunc = func(ctx context.Context, userAddress string, tokenAddress string, paymentAmount *big.Int) ([]*types.DvpDeposit, error) {
			return []*types.DvpDeposit{deposit}, nil
		}
		deps.DepositConsolidator.PrepareDepositsForJSProofFunc = func(ctx context.Context, _ string, sourceViewPublicKey []byte, deposits []*types.DvpDeposit) ([]*types.DvpDeposit, error) {
			return deposits, nil
		}
		deps.DvpDepositRepository.BatchUpdateStatusForCommitmentsFunc = func(ctx context.Context, commitments []string, status types.DvpDepositStatus) error {
			return nil
		}
		deps.KeysClient.GetPaymentSpendKeyFunc = func(ctx context.Context, in *keys.GetPaymentSpendKeyRequest, opts ...grpc.CallOption) (*keys.PaymentSpendKeyResponse, error) {
			return keypair, nil
		}
		deps.DvpProofGenerator.GenerateEnygmaJSProofFunc = func(ctx context.Context, sourceViewPublicKey []byte, nftCommitment *big.Int, destinationPaymentPublicKey *big.Int, destinationSalt *big.Int, paymentAmount *big.Int, tokenAddress string, deposits []*types.DvpDeposit) (*dvp.ProofReceipt, error) {
			return proof, nil
		}
		deps.DvpDepositRepository.BatchUpsertNullifiersFunc = func(ctx context.Context, commitmentNullifierMap map[string]string) error {
			return nil
		}
		deps.RetryService.RetryOperationFunc = func(ctx context.Context, operationName string, maxRetries int, blockNumber uint64, executeOperation func(ctx context.Context, nextBlockNumber uint64) error) (uint64, error) {
			assert.Equal(t, "withdrawal", operationName)
			return 0, errors.New("retry failed")
		}

		init := setupInitiator(deps)

		referenceId := [32]byte{}
		_, err := init.HandleEnygmaWithdrawal(
			context.Background(),
			"test-event-id",
			100,
			createTestResourceId(),
			referenceId,
			common.HexToAddress("0xBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBB"),
			big.NewInt(1000),
			common.HexToHash("0xabcd"),
		)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "retry failed")
		// Verify all steps up to retry were called
		assert.Equal(
			t,
			2,
			len(deps.DvpDepositRepository.BatchUpdateStatusForCommitmentsCalls()),
		) // Lock and then Unspent cleanup
		assert.Equal(t, 1, len(deps.DvpProofGenerator.GenerateEnygmaJSProofCalls()))
		assert.Equal(t, 1, len(deps.DvpDepositRepository.BatchUpsertNullifiersCalls()))
		assert.Equal(t, 1, len(deps.RetryService.RetryOperationCalls()))
		// Verify cleanup: deposits were reverted to Unspent status
		calls := deps.DvpDepositRepository.BatchUpdateStatusForCommitmentsCalls()
		assert.Equal(t, types.DvpDepositUnspent, calls[1].Status) // Second call should revert to Unspent
		// Verify withdrawal was NOT executed
		assert.Equal(t, 0, len(deps.EnygmaHandlerClient.ReceiveWithdrawCalls()))
	})

	t.Run("returns error if failed to revert deposits to unspent after retry operation fails", func(t *testing.T) {
		deps := setupInitiatorDeps()
		deposit := createTestDvpDeposit("0x1234567890123456789012345678901234567890", big.NewInt(1000))
		proof := createTestProofReceipt()
		keypair := createTestPaymentSpendKey()

		deps.CCEndpointClient.GetResourceAddressFunc = func(_ context.Context, resourceId string) (common.Address, error) {
			return common.HexToAddress("0x1234567890123456789012345678901234567890"), nil
		}
		deps.DepositFinder.FindEnygmaDepositsFunc = func(ctx context.Context, userAddress string, tokenAddress string, paymentAmount *big.Int) ([]*types.DvpDeposit, error) {
			return []*types.DvpDeposit{deposit}, nil
		}
		deps.DepositConsolidator.PrepareDepositsForJSProofFunc = func(ctx context.Context, _ string, sourceViewPublicKey []byte, deposits []*types.DvpDeposit) ([]*types.DvpDeposit, error) {
			return deposits, nil
		}

		// Track calls to BatchUpdateStatusForCommitments
		callCount := 0
		deps.DvpDepositRepository.BatchUpdateStatusForCommitmentsFunc = func(ctx context.Context, commitments []string, status types.DvpDepositStatus) error {
			callCount++
			// First call (locking) succeeds, second call (cleanup revert) fails
			if callCount == 2 && status == types.DvpDepositUnspent {
				return errors.New("revert failed")
			}
			return nil
		}

		deps.KeysClient.GetPaymentSpendKeyFunc = func(ctx context.Context, in *keys.GetPaymentSpendKeyRequest, opts ...grpc.CallOption) (*keys.PaymentSpendKeyResponse, error) {
			return keypair, nil
		}
		deps.DvpProofGenerator.GenerateEnygmaJSProofFunc = func(ctx context.Context, sourceViewPublicKey []byte, nftCommitment *big.Int, destinationPaymentPublicKey *big.Int, destinationSalt *big.Int, paymentAmount *big.Int, tokenAddress string, deposits []*types.DvpDeposit) (*dvp.ProofReceipt, error) {
			return proof, nil
		}
		deps.DvpDepositRepository.BatchUpsertNullifiersFunc = func(ctx context.Context, commitmentNullifierMap map[string]string) error {
			return nil
		}
		deps.RetryService.RetryOperationFunc = func(ctx context.Context, operationName string, maxRetries int, blockNumber uint64, executeOperation func(ctx context.Context, nextBlockNumber uint64) error) (uint64, error) {
			return 0, errors.New("retry failed")
		}

		init := setupInitiator(deps)

		referenceId := [32]byte{}
		_, err := init.HandleEnygmaWithdrawal(
			context.Background(),
			"test-event-id",
			100,
			createTestResourceId(),
			referenceId,
			common.HexToAddress("0xBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBB"),
			big.NewInt(1000),
			common.HexToHash("0xabcd"),
		)

		// Should return the cleanup error, not the retry error
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "revert failed")
		// Verify BatchUpdateStatusForCommitments was called twice (lock + cleanup attempt)
		assert.Equal(t, 2, len(deps.DvpDepositRepository.BatchUpdateStatusForCommitmentsCalls()))
		// Verify PL operations were NOT attempted
		assert.Equal(t, 0, len(deps.PLEndpointClient.GetResourceAddressCalls()))
		assert.Equal(t, 0, len(deps.EnygmaHandlerClient.ReceiveWithdrawCalls()))
	})

	t.Run("returns error if failed to get resource address on PL", func(t *testing.T) {
		deps := setupInitiatorDeps()
		deposit := createTestDvpDeposit("0x1234567890123456789012345678901234567890", big.NewInt(1000))
		proof := createTestProofReceipt()
		keypair := createTestPaymentSpendKey()
		to := common.HexToAddress("0xBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBB")
		amount := big.NewInt(1000)

		deps.CCEndpointClient.GetResourceAddressFunc = func(_ context.Context, resourceId string) (common.Address, error) {
			return common.HexToAddress("0x1234567890123456789012345678901234567890"), nil
		}
		deps.DepositFinder.FindEnygmaDepositsFunc = func(ctx context.Context, userAddress string, tokenAddress string, paymentAmount *big.Int) ([]*types.DvpDeposit, error) {
			return []*types.DvpDeposit{deposit}, nil
		}
		deps.DepositConsolidator.PrepareDepositsForJSProofFunc = func(ctx context.Context, _ string, sourceViewPublicKey []byte, deposits []*types.DvpDeposit) ([]*types.DvpDeposit, error) {
			return deposits, nil
		}
		deps.DvpDepositRepository.BatchUpdateStatusForCommitmentsFunc = func(ctx context.Context, commitments []string, status types.DvpDepositStatus) error {
			return nil
		}
		deps.KeysClient.GetPaymentSpendKeyFunc = func(ctx context.Context, in *keys.GetPaymentSpendKeyRequest, opts ...grpc.CallOption) (*keys.PaymentSpendKeyResponse, error) {
			return keypair, nil
		}
		deps.DvpProofGenerator.GenerateEnygmaJSProofFunc = func(ctx context.Context, sourceViewPublicKey []byte, nftCommitment *big.Int, destinationPaymentPublicKey *big.Int, destinationSalt *big.Int, paymentAmount *big.Int, tokenAddress string, deposits []*types.DvpDeposit) (*dvp.ProofReceipt, error) {
			return proof, nil
		}
		deps.DvpDepositRepository.BatchUpsertNullifiersFunc = func(ctx context.Context, commitmentNullifierMap map[string]string) error {
			return nil
		}
		deps.RetryService.RetryOperationFunc = func(ctx context.Context, operationName string, maxRetries int, blockNumber uint64, executeOperation func(ctx context.Context, nextBlockNumber uint64) error) (uint64, error) {
			err := executeOperation(ctx, blockNumber)
			return blockNumber, err
		}
		deps.Executor.ExecuteEnygmaWithdrawalFunc = func(ctx context.Context, _ string, resourceId string, amount *big.Int, deposits []*types.DvpDeposit, blockNumber uint64, enygmaAddress common.Address, proof *dvp.ProofReceipt, to common.Address, txHash common.Hash) error {
			return nil
		}
		deps.PLEndpointClient.GetResourceAddressFunc = func(_ context.Context, resourceId string) (common.Address, error) {
			return common.Address{}, errors.New("PL endpoint error")
		}

		init := setupInitiator(deps)

		referenceId := [32]byte{}
		_, err := init.HandleEnygmaWithdrawal(
			context.Background(),
			"test-event-id",
			100,
			createTestResourceId(),
			referenceId,
			to,
			amount,
			common.HexToHash("0xabcd"),
		)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "PL endpoint error")
		// Verify execution succeeded
		assert.Equal(t, 1, len(deps.Executor.ExecuteEnygmaWithdrawalCalls()))
		// Verify PL GetResourceAddress was called
		assert.Equal(t, 1, len(deps.PLEndpointClient.GetResourceAddressCalls()))
		// Verify subsequent operations were NOT called
		assert.Equal(t, 0, len(deps.EnygmaHandlerClient.ReceiveWithdrawCalls()))
	})

	t.Run("returns error if failed to receive withdrawal", func(t *testing.T) {
		deps := setupInitiatorDeps()
		deposit := createTestDvpDeposit("0x1234567890123456789012345678901234567890", big.NewInt(1000))
		proof := createTestProofReceipt()
		keypair := createTestPaymentSpendKey()
		to := common.HexToAddress("0xBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBB")
		amount := big.NewInt(1000)

		deps.CCEndpointClient.GetResourceAddressFunc = func(_ context.Context, resourceId string) (common.Address, error) {
			return common.HexToAddress("0x1234567890123456789012345678901234567890"), nil
		}
		deps.DepositFinder.FindEnygmaDepositsFunc = func(ctx context.Context, userAddress string, tokenAddress string, paymentAmount *big.Int) ([]*types.DvpDeposit, error) {
			return []*types.DvpDeposit{deposit}, nil
		}
		deps.DepositConsolidator.PrepareDepositsForJSProofFunc = func(ctx context.Context, _ string, sourceViewPublicKey []byte, deposits []*types.DvpDeposit) ([]*types.DvpDeposit, error) {
			return deposits, nil
		}
		deps.DvpDepositRepository.BatchUpdateStatusForCommitmentsFunc = func(ctx context.Context, commitments []string, status types.DvpDepositStatus) error {
			return nil
		}
		deps.KeysClient.GetPaymentSpendKeyFunc = func(ctx context.Context, in *keys.GetPaymentSpendKeyRequest, opts ...grpc.CallOption) (*keys.PaymentSpendKeyResponse, error) {
			return keypair, nil
		}
		deps.DvpProofGenerator.GenerateEnygmaJSProofFunc = func(ctx context.Context, sourceViewPublicKey []byte, nftCommitment *big.Int, destinationPaymentPublicKey *big.Int, destinationSalt *big.Int, paymentAmount *big.Int, tokenAddress string, deposits []*types.DvpDeposit) (*dvp.ProofReceipt, error) {
			return proof, nil
		}
		deps.DvpDepositRepository.BatchUpsertNullifiersFunc = func(ctx context.Context, commitmentNullifierMap map[string]string) error {
			return nil
		}
		deps.RetryService.RetryOperationFunc = func(ctx context.Context, operationName string, maxRetries int, blockNumber uint64, executeOperation func(ctx context.Context, nextBlockNumber uint64) error) (uint64, error) {
			err := executeOperation(ctx, blockNumber)
			return blockNumber, err
		}
		deps.Executor.ExecuteEnygmaWithdrawalFunc = func(ctx context.Context, _ string, resourceId string, amount *big.Int, deposits []*types.DvpDeposit, blockNumber uint64, enygmaAddress common.Address, proof *dvp.ProofReceipt, to common.Address, txHash common.Hash) error {
			return nil
		}
		deps.PLEndpointClient.GetResourceAddressFunc = func(_ context.Context, resourceId string) (common.Address, error) {
			return common.HexToAddress("0x9999999999999999999999999999999999999999"), nil
		}
		deps.EnygmaHandlerClient.ReceiveWithdrawFunc = func(ctx context.Context, _ string, tokenAddress common.Address, toAddr common.Address, value *big.Int, referenceId [32]byte) error {
			return errors.New("withdrawal failed")
		}

		init := setupInitiator(deps)

		referenceId := [32]byte{}
		_, err := init.HandleEnygmaWithdrawal(
			context.Background(),
			"test-event-id",
			100,
			createTestResourceId(),
			referenceId,
			to,
			amount,
			common.HexToHash("0xabcd"),
		)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "withdrawal failed")
		// Verify all operations were called
		assert.Equal(t, 1, len(deps.Executor.ExecuteEnygmaWithdrawalCalls()))
		assert.Equal(t, 1, len(deps.PLEndpointClient.GetResourceAddressCalls()))
		assert.Equal(t, 1, len(deps.EnygmaHandlerClient.ReceiveWithdrawCalls()))
	})

	t.Run("no consolidation needed - withdrawal amount equals total deposits", func(t *testing.T) {
		deps := setupInitiatorDeps()
		deposit := createTestDvpDeposit("0x1234567890123456789012345678901234567890", big.NewInt(1000))
		deposit.Commitment = big.NewInt(111)
		proof := createTestProofReceipt()
		keypair := createTestPaymentSpendKey()
		to := common.HexToAddress("0xBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBB")
		amount := big.NewInt(1000)

		deps.CCEndpointClient.GetResourceAddressFunc = func(_ context.Context, resourceId string) (common.Address, error) {
			return common.HexToAddress("0x1234567890123456789012345678901234567890"), nil
		}
		deps.DepositFinder.FindEnygmaDepositsFunc = func(ctx context.Context, userAddress string, tokenAddress string, paymentAmount *big.Int) ([]*types.DvpDeposit, error) {
			return []*types.DvpDeposit{deposit}, nil
		}
		deps.DepositConsolidator.PrepareDepositsForJSProofFunc = func(ctx context.Context, _ string, sourceViewPublicKey []byte, deposits []*types.DvpDeposit) ([]*types.DvpDeposit, error) {
			return deposits, nil
		}
		deps.DvpDepositRepository.BatchUpdateStatusForCommitmentsFunc = func(ctx context.Context, commitments []string, status types.DvpDepositStatus) error {
			return nil
		}
		deps.KeysClient.GetPaymentSpendKeyFunc = func(ctx context.Context, in *keys.GetPaymentSpendKeyRequest, opts ...grpc.CallOption) (*keys.PaymentSpendKeyResponse, error) {
			return keypair, nil
		}
		deps.DvpProofGenerator.GenerateEnygmaJSProofFunc = func(ctx context.Context, sourceViewPublicKey []byte, nftCommitment *big.Int, destinationPaymentPublicKey *big.Int, destinationSalt *big.Int, paymentAmount *big.Int, tokenAddress string, deposits []*types.DvpDeposit) (*dvp.ProofReceipt, error) {
			// Verify original deposit is used (1000 tokens)
			assert.Equal(t, 1, len(deposits))
			assert.Equal(t, big.NewInt(1000), deposits[0].TokenAmount)
			return proof, nil
		}
		deps.DvpDepositRepository.BatchUpsertNullifiersFunc = func(ctx context.Context, commitmentNullifierMap map[string]string) error {
			return nil
		}
		deps.RetryService.RetryOperationFunc = func(ctx context.Context, operationName string, maxRetries int, blockNumber uint64, executeOperation func(ctx context.Context, nextBlockNumber uint64) error) (uint64, error) {
			err := executeOperation(ctx, blockNumber)
			return blockNumber, err
		}
		deps.Executor.ExecuteEnygmaWithdrawalFunc = func(ctx context.Context, _ string, resourceId string, amount *big.Int, deposits []*types.DvpDeposit, blockNumber uint64, enygmaAddress common.Address, proof *dvp.ProofReceipt, to common.Address, txHash common.Hash) error {
			return nil
		}
		deps.PLEndpointClient.GetResourceAddressFunc = func(_ context.Context, resourceId string) (common.Address, error) {
			return common.HexToAddress("0x9999999999999999999999999999999999999999"), nil
		}
		deps.EnygmaHandlerClient.ReceiveWithdrawFunc = func(ctx context.Context, _ string, tokenAddress common.Address, toAddr common.Address, value *big.Int, referenceId [32]byte) error {
			return nil
		}

		init := setupInitiator(deps)

		referenceId := [32]byte{}
		_, err := init.HandleEnygmaWithdrawal(
			context.Background(),
			"test-event-id",
			100,
			createTestResourceId(),
			referenceId,
			to,
			amount,
			common.HexToHash("0xabcd"),
		)

		assert.NoError(t, err)
		// Verify consolidation was NOT called (amount equals total)
		assert.Equal(t, 0, len(deps.DepositConsolidator.ConsolidateEnygmaDepositsCalls()))
		// Verify original deposit was used in proof generation
		assert.Equal(t, 1, len(deps.DvpProofGenerator.GenerateEnygmaJSProofCalls()))
	})

	t.Run("consolidation triggered - exact match found", func(t *testing.T) {
		deps := setupInitiatorDeps()

		// User has more than withdrawal amount: 600 + 400 = 1000
		deposit1 := createTestDvpDeposit("0x1234567890123456789012345678901234567890", big.NewInt(600))
		deposit1.Commitment = big.NewInt(111)

		deposit2 := createTestDvpDeposit("0x1234567890123456789012345678901234567890", big.NewInt(400))
		deposit2.Commitment = big.NewInt(222)

		// After consolidation, returns exact match plus remainder
		consolidatedMain := createTestDvpDeposit("0x1234567890123456789012345678901234567890", big.NewInt(750))
		consolidatedMain.Commitment = big.NewInt(333)
		consolidatedMain.Salt = big.NewInt(99999)

		consolidatedRemainder := createTestDvpDeposit("0x1234567890123456789012345678901234567890", big.NewInt(250))
		consolidatedRemainder.Commitment = big.NewInt(334)
		consolidatedRemainder.Salt = big.NewInt(99998)

		proof := createTestProofReceipt()
		keypair := createTestPaymentSpendKey()
		to := common.HexToAddress("0xBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBB")
		amount := big.NewInt(750)

		deps.CCEndpointClient.GetResourceAddressFunc = func(_ context.Context, resourceId string) (common.Address, error) {
			return common.HexToAddress("0x1234567890123456789012345678901234567890"), nil
		}
		deps.DepositFinder.FindEnygmaDepositsFunc = func(ctx context.Context, userAddress string, tokenAddress string, paymentAmount *big.Int) ([]*types.DvpDeposit, error) {
			return []*types.DvpDeposit{deposit1, deposit2}, nil
		}
		deps.DepositConsolidator.PrepareDepositsForJSProofFunc = func(ctx context.Context, _ string, sourceViewPublicKey []byte, deposits []*types.DvpDeposit) ([]*types.DvpDeposit, error) {
			return deposits, nil
		}
		deps.DepositConsolidator.ConsolidateEnygmaDepositsFunc = func(ctx context.Context, _ string, sourceViewPublicKey []byte, deposits []*types.DvpDeposit, consolidationAmount *big.Int) ([]*types.DvpDeposit, error) {
			// Verify correct amount passed to consolidation
			assert.Equal(t, big.NewInt(750), consolidationAmount)
			// Return consolidated deposits with exact match
			return []*types.DvpDeposit{consolidatedMain, consolidatedRemainder}, nil
		}
		deps.DvpDepositRepository.BatchUpdateStatusForCommitmentsFunc = func(ctx context.Context, commitments []string, status types.DvpDepositStatus) error {
			return nil
		}
		deps.KeysClient.GetPaymentSpendKeyFunc = func(ctx context.Context, in *keys.GetPaymentSpendKeyRequest, opts ...grpc.CallOption) (*keys.PaymentSpendKeyResponse, error) {
			return keypair, nil
		}
		deps.DvpProofGenerator.GenerateEnygmaJSProofFunc = func(ctx context.Context, sourceViewPublicKey []byte, nftCommitment *big.Int, destinationPaymentPublicKey *big.Int, destinationSalt *big.Int, paymentAmount *big.Int, tokenAddress string, deposits []*types.DvpDeposit) (*dvp.ProofReceipt, error) {
			// Verify only the exact match deposit is passed (750 tokens)
			assert.Equal(t, 1, len(deposits))
			assert.Equal(t, big.NewInt(750), deposits[0].TokenAmount)
			return proof, nil
		}
		deps.DvpDepositRepository.BatchUpsertNullifiersFunc = func(ctx context.Context, commitmentNullifierMap map[string]string) error {
			return nil
		}
		deps.RetryService.RetryOperationFunc = func(ctx context.Context, operationName string, maxRetries int, blockNumber uint64, executeOperation func(ctx context.Context, nextBlockNumber uint64) error) (uint64, error) {
			err := executeOperation(ctx, blockNumber)
			return blockNumber, err
		}
		deps.Executor.ExecuteEnygmaWithdrawalFunc = func(ctx context.Context, _ string, resourceId string, amount *big.Int, deposits []*types.DvpDeposit, blockNumber uint64, enygmaAddress common.Address, proof *dvp.ProofReceipt, to common.Address, txHash common.Hash) error {
			return nil
		}
		deps.PLEndpointClient.GetResourceAddressFunc = func(_ context.Context, resourceId string) (common.Address, error) {
			return common.HexToAddress("0x9999999999999999999999999999999999999999"), nil
		}
		deps.EnygmaHandlerClient.ReceiveWithdrawFunc = func(ctx context.Context, _ string, tokenAddress common.Address, toAddr common.Address, value *big.Int, referenceId [32]byte) error {
			return nil
		}

		init := setupInitiator(deps)

		referenceId := [32]byte{}
		_, err := init.HandleEnygmaWithdrawal(
			context.Background(),
			"test-event-id",
			100,
			createTestResourceId(),
			referenceId,
			to,
			amount,
			common.HexToHash("0xabcd"),
		)

		assert.NoError(t, err)
		// Verify consolidation was called with correct amount
		consolidateCalls := deps.DepositConsolidator.ConsolidateEnygmaDepositsCalls()
		assert.Equal(t, 1, len(consolidateCalls))
		assert.Equal(t, big.NewInt(750), consolidateCalls[0].ConsolidationAmount)
		assert.Equal(t, 2, len(consolidateCalls[0].DepositsToConsolidate))
		// Verify proof was generated with filtered deposit (750 only)
		assert.Equal(t, 1, len(deps.DvpProofGenerator.GenerateEnygmaJSProofCalls()))
	})

	t.Run("consolidation triggered - no exact match found", func(t *testing.T) {
		deps := setupInitiatorDeps()

		// User has more than needed: 600 + 500 = 1100
		deposit1 := createTestDvpDeposit("0x1234567890123456789012345678901234567890", big.NewInt(600))
		deposit1.Commitment = big.NewInt(111)

		deposit2 := createTestDvpDeposit("0x1234567890123456789012345678901234567890", big.NewInt(500))
		deposit2.Commitment = big.NewInt(222)

		// After consolidation, no deposit matches exact 750
		consolidatedDeposit1 := createTestDvpDeposit("0x1234567890123456789012345678901234567890", big.NewInt(800))
		consolidatedDeposit1.Commitment = big.NewInt(333)

		consolidatedDeposit2 := createTestDvpDeposit("0x1234567890123456789012345678901234567890", big.NewInt(300))
		consolidatedDeposit2.Commitment = big.NewInt(334)

		proof := createTestProofReceipt()
		keypair := createTestPaymentSpendKey()
		to := common.HexToAddress("0xBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBB")
		amount := big.NewInt(750)

		deps.CCEndpointClient.GetResourceAddressFunc = func(_ context.Context, resourceId string) (common.Address, error) {
			return common.HexToAddress("0x1234567890123456789012345678901234567890"), nil
		}
		deps.DepositFinder.FindEnygmaDepositsFunc = func(ctx context.Context, userAddress string, tokenAddress string, paymentAmount *big.Int) ([]*types.DvpDeposit, error) {
			return []*types.DvpDeposit{deposit1, deposit2}, nil
		}
		deps.DepositConsolidator.PrepareDepositsForJSProofFunc = func(ctx context.Context, _ string, sourceViewPublicKey []byte, deposits []*types.DvpDeposit) ([]*types.DvpDeposit, error) {
			return deposits, nil
		}
		deps.DepositConsolidator.ConsolidateEnygmaDepositsFunc = func(ctx context.Context, _ string, sourceViewPublicKey []byte, deposits []*types.DvpDeposit, consolidationAmount *big.Int) ([]*types.DvpDeposit, error) {
			// Return deposits without exact 750 match
			return []*types.DvpDeposit{consolidatedDeposit1, consolidatedDeposit2}, nil
		}
		deps.DvpDepositRepository.BatchUpdateStatusForCommitmentsFunc = func(ctx context.Context, commitments []string, status types.DvpDepositStatus) error {
			return nil
		}
		deps.KeysClient.GetPaymentSpendKeyFunc = func(ctx context.Context, in *keys.GetPaymentSpendKeyRequest, opts ...grpc.CallOption) (*keys.PaymentSpendKeyResponse, error) {
			return keypair, nil
		}
		deps.DvpProofGenerator.GenerateEnygmaJSProofFunc = func(ctx context.Context, sourceViewPublicKey []byte, nftCommitment *big.Int, destinationPaymentPublicKey *big.Int, destinationSalt *big.Int, paymentAmount *big.Int, tokenAddress string, deposits []*types.DvpDeposit) (*dvp.ProofReceipt, error) {
			// Verify both deposits passed (no filtering)
			assert.Equal(t, 2, len(deposits))
			return proof, nil
		}
		deps.DvpDepositRepository.BatchUpsertNullifiersFunc = func(ctx context.Context, commitmentNullifierMap map[string]string) error {
			return nil
		}
		deps.RetryService.RetryOperationFunc = func(ctx context.Context, operationName string, maxRetries int, blockNumber uint64, executeOperation func(ctx context.Context, nextBlockNumber uint64) error) (uint64, error) {
			err := executeOperation(ctx, blockNumber)
			return blockNumber, err
		}
		deps.Executor.ExecuteEnygmaWithdrawalFunc = func(ctx context.Context, _ string, resourceId string, amount *big.Int, deposits []*types.DvpDeposit, blockNumber uint64, enygmaAddress common.Address, proof *dvp.ProofReceipt, to common.Address, txHash common.Hash) error {
			return nil
		}
		deps.PLEndpointClient.GetResourceAddressFunc = func(_ context.Context, resourceId string) (common.Address, error) {
			return common.HexToAddress("0x9999999999999999999999999999999999999999"), nil
		}
		deps.EnygmaHandlerClient.ReceiveWithdrawFunc = func(ctx context.Context, _ string, tokenAddress common.Address, toAddr common.Address, value *big.Int, referenceId [32]byte) error {
			return nil
		}

		init := setupInitiator(deps)

		referenceId := [32]byte{}
		_, err := init.HandleEnygmaWithdrawal(
			context.Background(),
			"test-event-id",
			100,
			createTestResourceId(),
			referenceId,
			to,
			amount,
			common.HexToHash("0xabcd"),
		)

		assert.NoError(t, err)
		// Verify consolidation was called
		consolidateCalls := deps.DepositConsolidator.ConsolidateEnygmaDepositsCalls()
		assert.Equal(t, 1, len(consolidateCalls))
		// Verify proof was generated with both deposits (no exact match to filter)
		proofCalls := deps.DvpProofGenerator.GenerateEnygmaJSProofCalls()
		assert.Equal(t, 1, len(proofCalls))
		assert.Equal(t, 2, len(proofCalls[0].Deposits))
	})
}
