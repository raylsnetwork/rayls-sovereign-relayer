package service_test

import (
	"context"
	"errors"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/raylsnetwork/rayls-sovereign-relayer/contractclient"
	"github.com/raylsnetwork/rayls-sovereign-relayer/enygma/service"
	"github.com/raylsnetwork/rayls-sovereign-relayer/enygma/testutils"
	privatehubservice "github.com/raylsnetwork/rayls-sovereign-relayer/private-relayer/source/service"
	"github.com/raylsnetwork/rayls-sovereign-relayer/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type revertEndpointClientMock struct {
	GetResourceAddressFunc  func(ctx context.Context, resourceId string) (common.Address, error)
	getResourceAddressCalls int
}

func (m *revertEndpointClientMock) GetResourceAddress(ctx context.Context, resourceId string) (common.Address, error) {
	m.getResourceAddressCalls++
	if m.GetResourceAddressFunc == nil {
		return common.Address{}, nil
	}
	return m.GetResourceAddressFunc(ctx, resourceId)
}

type revertEnygmaHandlerClientMock struct {
	RevertSrcTransferBatchFunc func(
		ctx context.Context,
		revertTxs []*types.EnygmaTransferFailed,
	) (map[string]contractclient.BatchResult, error)
	RevertSrcSupplyBatchFunc func(
		ctx context.Context,
		revertTxs []*types.EnygmaSupplyUpdateFailed,
	) (map[string]contractclient.BatchResult, error)

	revertSrcTransferBatchCalls int
	revertSrcSupplyBatchCalls   int
}

func (m *revertEnygmaHandlerClientMock) RevertSrcTransferBatch(
	ctx context.Context,
	revertTxs []*types.EnygmaTransferFailed,
) (map[string]contractclient.BatchResult, error) {
	m.revertSrcTransferBatchCalls++
	if m.RevertSrcTransferBatchFunc == nil {
		return make(map[string]contractclient.BatchResult), nil
	}
	return m.RevertSrcTransferBatchFunc(ctx, revertTxs)
}

func (m *revertEnygmaHandlerClientMock) RevertSrcSupplyBatch(
	ctx context.Context,
	revertTxs []*types.EnygmaSupplyUpdateFailed,
) (map[string]contractclient.BatchResult, error) {
	m.revertSrcSupplyBatchCalls++
	if m.RevertSrcSupplyBatchFunc == nil {
		return make(map[string]contractclient.BatchResult), nil
	}
	return m.RevertSrcSupplyBatchFunc(ctx, revertTxs)
}

// Test data helpers

func revertTestResourceId() string {
	return "test-resource-revert"
}

func revertTestEnygmaAddress() common.Address {
	return common.HexToAddress("0xaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa")
}

func revertTestSenderAddress() common.Address {
	return common.HexToAddress("0xbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb")
}

func revertTestReceiverAddress() common.Address {
	return common.HexToAddress("0xcccccccccccccccccccccccccccccccccccccccc")
}

func createTestSupplyUpdateEvent(
	to common.Address,
	amount *big.Int,
	txHash common.Hash,
) privatehubservice.EnygmaSupplyUpdate {
	return privatehubservice.EnygmaSupplyUpdate{
		TxHash: txHash,
		To:     to,
		Amount: amount,
	}
}

func createTestTransferBatchTx(fromAddr, toAddr common.Address, amount *big.Int) *types.EnygmaTransferBatchTx {
	return &types.EnygmaTransferBatchTx{
		MessageId:     "msg-1",
		ReferenceId:   [32]byte{1, 2, 3, 4},
		FromAddress:   fromAddr,
		ToAmount:      amount,
		ToAddress:     toAddr,
		SendTimestamp: 1234567890,
	}
}

func createTestBatchResultSuccess() contractclient.BatchResult {
	return contractclient.BatchResult{}
}

func createTestBatchResultError(err error) contractclient.BatchResult {
	return contractclient.BatchResult{Err: err}
}

func TestEnygmaRevertService_RevertEnygmaSupplyUpdate(t *testing.T) {
	t.Run("successfully reverts supply updates with mints and burns", func(t *testing.T) {
		ctx := context.Background()
		resourceId := revertTestResourceId()
		enygmaAddress := revertTestEnygmaAddress()

		supplyUpdates := []privatehubservice.EnygmaSupplyUpdate{
			createTestSupplyUpdateEvent(
				revertTestReceiverAddress(),
				big.NewInt(1000),
				common.HexToHash("0x1111111111111111111111111111111111111111111111111111111111111111"),
			), // mint
			createTestSupplyUpdateEvent(
				revertTestSenderAddress(),
				big.NewInt(-500),
				common.HexToHash("0x2222222222222222222222222222222222222222222222222222222222222222"),
			), // burn (negative)
		}

		tracer := &testutils.MockTracer{}
		endpointClient := &revertEndpointClientMock{
			GetResourceAddressFunc: func(_ context.Context, resourceId string) (common.Address, error) {
				return enygmaAddress, nil
			},
		}
		handlerClient := &revertEnygmaHandlerClientMock{
			RevertSrcSupplyBatchFunc: func(ctx context.Context, data []*types.EnygmaSupplyUpdateFailed) (map[string]contractclient.BatchResult, error) {
				assert.Equal(t, 2, len(data))
				assert.Equal(t, big.NewInt(1000), data[0].Amount) // mint amount
				assert.Equal(t, big.NewInt(500), data[1].Amount)  // burn amount (absolute value)
				assert.Equal(t, revertTestReceiverAddress(), data[0].To)
				assert.Equal(t, revertTestSenderAddress(), data[1].To)
				assert.Equal(t, types.EnygmaMint, data[0].Type)
				assert.Equal(t, types.EnygmaBurn, data[1].Type)
				return map[string]contractclient.BatchResult{
					"0": createTestBatchResultSuccess(),
					"1": createTestBatchResultSuccess(),
				}, nil
			},
		}

		svc := service.NewEnygmaRevertService(tracer, endpointClient, handlerClient)

		err := svc.RevertEnygmaSupplyUpdate(ctx, resourceId, supplyUpdates)

		require.NoError(t, err)
		assert.Equal(t, 1, endpointClient.getResourceAddressCalls)
		assert.Equal(t, 1, handlerClient.revertSrcSupplyBatchCalls)
	})

	t.Run("returns error when GetResourceAddress fails", func(t *testing.T) {
		ctx := context.Background()
		resourceId := revertTestResourceId()

		tracer := &testutils.MockTracer{}
		endpointClient := &revertEndpointClientMock{
			GetResourceAddressFunc: func(_ context.Context, resourceId string) (common.Address, error) {
				return common.Address{}, errors.New("address lookup failed")
			},
		}
		handlerClient := &revertEnygmaHandlerClientMock{}

		svc := service.NewEnygmaRevertService(tracer, endpointClient, handlerClient)

		err := svc.RevertEnygmaSupplyUpdate(ctx, resourceId, []privatehubservice.EnygmaSupplyUpdate{})

		require.Error(t, err)
		assert.Contains(t, err.Error(), "address lookup failed")
	})

	t.Run("returns error when RevertSrcSupplyBatch fails", func(t *testing.T) {
		ctx := context.Background()
		resourceId := revertTestResourceId()
		enygmaAddress := revertTestEnygmaAddress()

		supplyUpdates := []privatehubservice.EnygmaSupplyUpdate{
			createTestSupplyUpdateEvent(
				revertTestReceiverAddress(),
				big.NewInt(1000),
				common.HexToHash("0x1111111111111111111111111111111111111111111111111111111111111111"),
			),
		}

		tracer := &testutils.MockTracer{}
		endpointClient := &revertEndpointClientMock{
			GetResourceAddressFunc: func(_ context.Context, resourceId string) (common.Address, error) {
				return enygmaAddress, nil
			},
		}
		handlerClient := &revertEnygmaHandlerClientMock{
			RevertSrcSupplyBatchFunc: func(ctx context.Context, data []*types.EnygmaSupplyUpdateFailed) (map[string]contractclient.BatchResult, error) {
				return nil, errors.New("batcher error")
			},
		}

		svc := service.NewEnygmaRevertService(tracer, endpointClient, handlerClient)

		err := svc.RevertEnygmaSupplyUpdate(ctx, resourceId, supplyUpdates)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "batcher error")
	})

	t.Run("does not return error when individual batch results contain errors", func(t *testing.T) {
		ctx := context.Background()
		resourceId := revertTestResourceId()
		enygmaAddress := revertTestEnygmaAddress()

		supplyUpdates := []privatehubservice.EnygmaSupplyUpdate{
			createTestSupplyUpdateEvent(
				revertTestReceiverAddress(),
				big.NewInt(1000),
				common.HexToHash("0x3333333333333333333333333333333333333333333333333333333333333333"),
			),
			createTestSupplyUpdateEvent(
				revertTestReceiverAddress(),
				big.NewInt(2000),
				common.HexToHash("0x4444444444444444444444444444444444444444444444444444444444444444"),
			),
		}

		tracer := &testutils.MockTracer{}
		endpointClient := &revertEndpointClientMock{
			GetResourceAddressFunc: func(_ context.Context, resourceId string) (common.Address, error) {
				return enygmaAddress, nil
			},
		}
		handlerClient := &revertEnygmaHandlerClientMock{
			RevertSrcSupplyBatchFunc: func(ctx context.Context, data []*types.EnygmaSupplyUpdateFailed) (map[string]contractclient.BatchResult, error) {
				assert.Equal(t, 2, len(data))
				return map[string]contractclient.BatchResult{
					"0": createTestBatchResultSuccess(),
					"1": createTestBatchResultError(errors.New("send failed")),
				}, nil
			},
		}

		svc := service.NewEnygmaRevertService(tracer, endpointClient, handlerClient)

		err := svc.RevertEnygmaSupplyUpdate(ctx, resourceId, supplyUpdates)

		require.NoError(t, err)
		assert.Equal(t, 1, handlerClient.revertSrcSupplyBatchCalls)
	})

	t.Run("handles empty supply update events", func(t *testing.T) {
		ctx := context.Background()
		resourceId := revertTestResourceId()
		enygmaAddress := revertTestEnygmaAddress()

		tracer := &testutils.MockTracer{}
		endpointClient := &revertEndpointClientMock{
			GetResourceAddressFunc: func(_ context.Context, resourceId string) (common.Address, error) {
				return enygmaAddress, nil
			},
		}
		handlerClient := &revertEnygmaHandlerClientMock{
			RevertSrcSupplyBatchFunc: func(ctx context.Context, data []*types.EnygmaSupplyUpdateFailed) (map[string]contractclient.BatchResult, error) {
				// Should send empty list
				assert.Equal(t, 0, len(data))
				return make(map[string]contractclient.BatchResult), nil
			},
		}

		svc := service.NewEnygmaRevertService(tracer, endpointClient, handlerClient)

		err := svc.RevertEnygmaSupplyUpdate(ctx, resourceId, []privatehubservice.EnygmaSupplyUpdate{})

		require.NoError(t, err)
		assert.Equal(t, 1, handlerClient.revertSrcSupplyBatchCalls)
	})
}

func TestEnygmaRevertService_RevertEnygmaTransfer(t *testing.T) {
	t.Run("successfully reverts transfers from multiple senders", func(t *testing.T) {
		ctx := context.Background()
		resourceId := revertTestResourceId()
		enygmaAddress := revertTestEnygmaAddress()

		sender1 := revertTestSenderAddress()
		sender2 := common.HexToAddress("0xdddddddddddddddddddddddddddddddddddddddd")

		transfers := map[string][]*types.EnygmaTransferBatchTx{
			"1": {
				createTestTransferBatchTx(sender1, revertTestReceiverAddress(), big.NewInt(100)),
				createTestTransferBatchTx(sender1, revertTestReceiverAddress(), big.NewInt(200)),
			},
			"2": {
				createTestTransferBatchTx(sender2, revertTestReceiverAddress(), big.NewInt(300)),
			},
		}

		tracer := &testutils.MockTracer{}
		endpointClient := &revertEndpointClientMock{
			GetResourceAddressFunc: func(_ context.Context, resourceId string) (common.Address, error) {
				return enygmaAddress, nil
			},
		}
		handlerClient := &revertEnygmaHandlerClientMock{
			RevertSrcTransferBatchFunc: func(ctx context.Context, data []*types.EnygmaTransferFailed) (map[string]contractclient.BatchResult, error) {
				assert.Equal(t, 2, len(data))

				// Build sender -> amount map for order-independent assertions
				senderAmounts := make(map[common.Address]*big.Int)
				for _, d := range data {
					senderAmounts[d.Sender] = d.Amount
				}

				assert.Equal(t, big.NewInt(300), senderAmounts[sender1])
				assert.Equal(t, big.NewInt(300), senderAmounts[sender2])

				return map[string]contractclient.BatchResult{
					"0": createTestBatchResultSuccess(),
					"1": createTestBatchResultSuccess(),
				}, nil
			},
		}

		svc := service.NewEnygmaRevertService(tracer, endpointClient, handlerClient)

		err := svc.RevertEnygmaTransfer(ctx, resourceId, transfers)

		require.NoError(t, err)
		assert.Equal(t, 1, handlerClient.revertSrcTransferBatchCalls)
	})

	t.Run("returns error when GetResourceAddress fails", func(t *testing.T) {
		ctx := context.Background()
		resourceId := revertTestResourceId()

		transfers := map[string][]*types.EnygmaTransferBatchTx{
			"1": {createTestTransferBatchTx(revertTestSenderAddress(), revertTestReceiverAddress(), big.NewInt(100))},
		}

		tracer := &testutils.MockTracer{}
		endpointClient := &revertEndpointClientMock{
			GetResourceAddressFunc: func(_ context.Context, resourceId string) (common.Address, error) {
				return common.Address{}, errors.New("endpoint error")
			},
		}
		handlerClient := &revertEnygmaHandlerClientMock{}

		svc := service.NewEnygmaRevertService(tracer, endpointClient, handlerClient)

		err := svc.RevertEnygmaTransfer(ctx, resourceId, transfers)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "endpoint error")
	})

	t.Run("handles empty transfers map", func(t *testing.T) {
		ctx := context.Background()
		resourceId := revertTestResourceId()
		enygmaAddress := revertTestEnygmaAddress()

		tracer := &testutils.MockTracer{}
		endpointClient := &revertEndpointClientMock{
			GetResourceAddressFunc: func(_ context.Context, resourceId string) (common.Address, error) {
				return enygmaAddress, nil
			},
		}
		handlerClient := &revertEnygmaHandlerClientMock{
			RevertSrcTransferBatchFunc: func(ctx context.Context, data []*types.EnygmaTransferFailed) (map[string]contractclient.BatchResult, error) {
				assert.Equal(t, 0, len(data))
				return make(map[string]contractclient.BatchResult), nil
			},
		}

		svc := service.NewEnygmaRevertService(tracer, endpointClient, handlerClient)

		err := svc.RevertEnygmaTransfer(ctx, resourceId, make(map[string][]*types.EnygmaTransferBatchTx))

		require.NoError(t, err)
	})
}

func TestEnygmaRevertService_RevertEnygmaDeposit(t *testing.T) {
	t.Run("successfully reverts single deposit", func(t *testing.T) {
		ctx := context.Background()
		resourceId := revertTestResourceId()
		enygmaAddress := revertTestEnygmaAddress()
		sender := revertTestSenderAddress()
		amount := big.NewInt(1000)
		referenceId := [32]byte{1, 2, 3, 4}

		tracer := &testutils.MockTracer{}
		endpointClient := &revertEndpointClientMock{
			GetResourceAddressFunc: func(_ context.Context, resourceId string) (common.Address, error) {
				return enygmaAddress, nil
			},
		}
		handlerClient := &revertEnygmaHandlerClientMock{
			RevertSrcTransferBatchFunc: func(ctx context.Context, data []*types.EnygmaTransferFailed) (map[string]contractclient.BatchResult, error) {
				assert.Equal(t, 1, len(data))
				assert.Equal(t, amount, data[0].Amount)
				assert.Equal(t, sender, data[0].Sender)
				assert.Equal(t, enygmaAddress, data[0].EnygmaAddress)
				assert.Equal(t, "failed to execute deposit on CC", data[0].Reason)
				return map[string]contractclient.BatchResult{
					"0": createTestBatchResultSuccess(),
				}, nil
			},
		}

		svc := service.NewEnygmaRevertService(tracer, endpointClient, handlerClient)

		err := svc.RevertEnygmaDeposit(ctx, resourceId, referenceId, sender, amount)

		require.NoError(t, err)
		assert.Equal(t, 1, handlerClient.revertSrcTransferBatchCalls)
	})

	t.Run("returns error when GetResourceAddressFunc fails", func(t *testing.T) {
		ctx := context.Background()
		resourceId := revertTestResourceId()
		sender := revertTestSenderAddress()
		amount := big.NewInt(1000)
		referenceId := [32]byte{1, 2, 3, 4}

		tracer := &testutils.MockTracer{}
		endpointClient := &revertEndpointClientMock{
			GetResourceAddressFunc: func(_ context.Context, resourceId string) (common.Address, error) {
				return common.Address{}, errors.New("lookup failed")
			},
		}
		handlerClient := &revertEnygmaHandlerClientMock{}

		svc := service.NewEnygmaRevertService(tracer, endpointClient, handlerClient)

		err := svc.RevertEnygmaDeposit(ctx, resourceId, referenceId, sender, amount)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "lookup failed")
	})

	t.Run("returns error when RevertSrcTransferBatch fails", func(t *testing.T) {
		ctx := context.Background()
		resourceId := revertTestResourceId()
		enygmaAddress := revertTestEnygmaAddress()

		tracer := &testutils.MockTracer{}
		endpointClient := &revertEndpointClientMock{
			GetResourceAddressFunc: func(_ context.Context, resourceId string) (common.Address, error) {
				return enygmaAddress, nil
			},
		}
		handlerClient := &revertEnygmaHandlerClientMock{
			RevertSrcTransferBatchFunc: func(ctx context.Context, data []*types.EnygmaTransferFailed) (map[string]contractclient.BatchResult, error) {
				return nil, errors.New("transfer batcher error")
			},
		}

		svc := service.NewEnygmaRevertService(tracer, endpointClient, handlerClient)

		err := svc.RevertEnygmaDeposit(ctx, resourceId, [32]byte{}, revertTestSenderAddress(), big.NewInt(100))

		require.Error(t, err)
		assert.Contains(t, err.Error(), "transfer batcher error")
	})

	t.Run("does not return error when individual batch results contain errors", func(t *testing.T) {
		ctx := context.Background()
		resourceId := revertTestResourceId()
		enygmaAddress := revertTestEnygmaAddress()

		tracer := &testutils.MockTracer{}
		endpointClient := &revertEndpointClientMock{
			GetResourceAddressFunc: func(_ context.Context, resourceId string) (common.Address, error) {
				return enygmaAddress, nil
			},
		}
		handlerClient := &revertEnygmaHandlerClientMock{
			RevertSrcTransferBatchFunc: func(ctx context.Context, data []*types.EnygmaTransferFailed) (map[string]contractclient.BatchResult, error) {
				return map[string]contractclient.BatchResult{
					"0": createTestBatchResultSuccess(),
					"1": createTestBatchResultError(errors.New("send failed")),
				}, nil
			},
		}

		svc := service.NewEnygmaRevertService(tracer, endpointClient, handlerClient)

		err := svc.RevertEnygmaDeposit(ctx, resourceId, [32]byte{}, revertTestSenderAddress(), big.NewInt(100))

		require.NoError(t, err)
		assert.Equal(t, 1, handlerClient.revertSrcTransferBatchCalls)
	})
}
