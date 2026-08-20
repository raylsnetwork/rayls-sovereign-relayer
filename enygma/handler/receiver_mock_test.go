package handler_test

import (
	"context"
	"math/big"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"

	"github.com/raylsnetwork/rayls-privacy-relayer-api/contractclient"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/enygma/handler"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/enygma/testutils"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/txsim"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/types"
)

type MockReceiverEndpointClient struct {
	GetResourceAddressFunc func(ctx context.Context, resourceId string) (common.Address, error)
	calls                  struct {
		GetResourceAddress []struct {
			Ctx        context.Context
			ResourceId string
		}
	}
	lock sync.RWMutex
}

func (m *MockReceiverEndpointClient) GetResourceAddress(ctx context.Context, resourceId string) (common.Address, error) {
	if m.GetResourceAddressFunc == nil {
		panic("MockReceiverEndpointClient.GetResourceAddressFunc: method is nil but GetResourceAddress was just called")
	}
	m.lock.Lock()
	m.calls.GetResourceAddress = append(m.calls.GetResourceAddress, struct {
		Ctx        context.Context
		ResourceId string
	}{Ctx: ctx, ResourceId: resourceId})
	m.lock.Unlock()
	return m.GetResourceAddressFunc(ctx, resourceId)
}

func (m *MockReceiverEndpointClient) GetResourceAddressCalls() []struct {
	Ctx        context.Context
	ResourceId string
} {
	m.lock.RLock()
	defer m.lock.RUnlock()
	return m.calls.GetResourceAddress
}

// MockReceiverDeployer mocks the receiverDeployer interface
type MockReceiverDeployer struct {
	DeployFunc func(ctx context.Context, resourceId [32]byte, initiatorChainId *big.Int) (common.Address, error)
	calls      struct {
		Deploy []struct {
			Ctx              context.Context
			ResourceId       [32]byte
			InitiatorChainId *big.Int
		}
	}
	lock sync.RWMutex
}

func (m *MockReceiverDeployer) Deploy(ctx context.Context, resourceId [32]byte, initiatorChainId *big.Int) (common.Address, error) {
	if m.DeployFunc == nil {
		panic("MockReceiverDeployer.DeployFunc: method is nil but Deploy was just called")
	}
	m.lock.Lock()
	m.calls.Deploy = append(m.calls.Deploy, struct {
		Ctx              context.Context
		ResourceId       [32]byte
		InitiatorChainId *big.Int
	}{Ctx: ctx, ResourceId: resourceId, InitiatorChainId: initiatorChainId})
	m.lock.Unlock()
	return m.DeployFunc(ctx, resourceId, initiatorChainId)
}

func (m *MockReceiverDeployer) DeployCalls() []struct {
	Ctx              context.Context
	ResourceId       [32]byte
	InitiatorChainId *big.Int
} {
	m.lock.RLock()
	defer m.lock.RUnlock()
	return m.calls.Deploy
}

type MockReceiverEnygmaHistoryRepository struct {
	InsertEnygmaHistoryFunc         func(ctx context.Context, history types.EnygmaHistory) error
	GetEnygmaHistoryByUniqueKeyFunc func(ctx context.Context, resourceId string, blockNumberPrivateHub *big.Int, fromChainId *big.Int, eventType types.EnygmaEventType) (*types.EnygmaHistory, error)
	calls                           struct {
		InsertEnygmaHistory []struct {
			Ctx     context.Context
			History types.EnygmaHistory
		}
		GetEnygmaHistoryByUniqueKey []struct {
			Ctx                   context.Context
			ResourceId            string
			BlockNumberPrivateHub *big.Int
			FromChainId           *big.Int
			EventType             types.EnygmaEventType
		}
	}
	lock sync.RWMutex
}

func (m *MockReceiverEnygmaHistoryRepository) InsertEnygmaHistory(ctx context.Context, history types.EnygmaHistory) error {
	if m.InsertEnygmaHistoryFunc == nil {
		panic("MockReceiverEnygmaHistoryRepository.InsertEnygmaHistoryFunc: method is nil but InsertEnygmaHistory was just called")
	}
	m.lock.Lock()
	m.calls.InsertEnygmaHistory = append(m.calls.InsertEnygmaHistory, struct {
		Ctx     context.Context
		History types.EnygmaHistory
	}{Ctx: ctx, History: history})
	m.lock.Unlock()
	return m.InsertEnygmaHistoryFunc(ctx, history)
}

func (m *MockReceiverEnygmaHistoryRepository) InsertEnygmaHistoryCalls() []struct {
	Ctx     context.Context
	History types.EnygmaHistory
} {
	m.lock.RLock()
	defer m.lock.RUnlock()
	return m.calls.InsertEnygmaHistory
}

func (m *MockReceiverEnygmaHistoryRepository) GetEnygmaHistoryByUniqueKey(ctx context.Context, resourceId string, blockNumberPrivateHub *big.Int, fromChainId *big.Int, eventType types.EnygmaEventType) (*types.EnygmaHistory, error) {
	m.lock.Lock()
	m.calls.GetEnygmaHistoryByUniqueKey = append(m.calls.GetEnygmaHistoryByUniqueKey, struct {
		Ctx                   context.Context
		ResourceId            string
		BlockNumberPrivateHub *big.Int
		FromChainId           *big.Int
		EventType             types.EnygmaEventType
	}{Ctx: ctx, ResourceId: resourceId, BlockNumberPrivateHub: blockNumberPrivateHub, FromChainId: fromChainId, EventType: eventType})
	m.lock.Unlock()
	if m.GetEnygmaHistoryByUniqueKeyFunc == nil {
		return nil, nil
	}
	return m.GetEnygmaHistoryByUniqueKeyFunc(ctx, resourceId, blockNumberPrivateHub, fromChainId, eventType)
}

func (m *MockReceiverEnygmaHistoryRepository) GetEnygmaHistoryByUniqueKeyCalls() []struct {
	Ctx                   context.Context
	ResourceId            string
	BlockNumberPrivateHub *big.Int
	FromChainId           *big.Int
	EventType             types.EnygmaEventType
} {
	m.lock.RLock()
	defer m.lock.RUnlock()
	return m.calls.GetEnygmaHistoryByUniqueKey
}

// MockEnygmaHandlerClient mocks the enygmaHandlerClient interface (mint + revert dest batches).
type MockEnygmaHandlerClient struct {
	ReceiveDestTransferBatchFunc func(ctx context.Context, transfers []*types.EnygmaCrossTransferData) (map[string]contractclient.BatchResult, error)
	RevertDestTransferBatchFunc  func(ctx context.Context, tokenAddress common.Address, reverts []*types.EnygmaCrossTransferData) (map[string]contractclient.BatchResult, error)
	calls                        struct {
		ReceiveDestTransferBatch []struct {
			Ctx       context.Context
			Transfers []*types.EnygmaCrossTransferData
		}
		RevertDestTransferBatch []struct {
			Ctx          context.Context
			TokenAddress common.Address
			Reverts      []*types.EnygmaCrossTransferData
		}
	}
	lock sync.RWMutex
}

func (m *MockEnygmaHandlerClient) ReceiveDestTransferBatch(ctx context.Context, transfers []*types.EnygmaCrossTransferData) (map[string]contractclient.BatchResult, error) {
	if m.ReceiveDestTransferBatchFunc == nil {
		panic("MockEnygmaHandlerClient.ReceiveDestTransferBatchFunc: method is nil but ReceiveDestTransferBatch was just called")
	}
	m.lock.Lock()
	m.calls.ReceiveDestTransferBatch = append(m.calls.ReceiveDestTransferBatch, struct {
		Ctx       context.Context
		Transfers []*types.EnygmaCrossTransferData
	}{Ctx: ctx, Transfers: transfers})
	m.lock.Unlock()
	return m.ReceiveDestTransferBatchFunc(ctx, transfers)
}

func (m *MockEnygmaHandlerClient) ReceiveDestTransferBatchCalls() []struct {
	Ctx       context.Context
	Transfers []*types.EnygmaCrossTransferData
} {
	m.lock.RLock()
	defer m.lock.RUnlock()
	return m.calls.ReceiveDestTransferBatch
}

func (m *MockEnygmaHandlerClient) RevertDestTransferBatch(ctx context.Context, tokenAddress common.Address, reverts []*types.EnygmaCrossTransferData) (map[string]contractclient.BatchResult, error) {
	if m.RevertDestTransferBatchFunc == nil {
		panic("MockEnygmaHandlerClient.RevertDestTransferBatchFunc: method is nil but RevertDestTransferBatch was just called")
	}
	m.lock.Lock()
	m.calls.RevertDestTransferBatch = append(m.calls.RevertDestTransferBatch, struct {
		Ctx          context.Context
		TokenAddress common.Address
		Reverts      []*types.EnygmaCrossTransferData
	}{Ctx: ctx, TokenAddress: tokenAddress, Reverts: reverts})
	m.lock.Unlock()
	return m.RevertDestTransferBatchFunc(ctx, tokenAddress, reverts)
}

func (m *MockEnygmaHandlerClient) RevertDestTransferBatchCalls() []struct {
	Ctx          context.Context
	TokenAddress common.Address
	Reverts      []*types.EnygmaCrossTransferData
} {
	m.lock.RLock()
	defer m.lock.RUnlock()
	return m.calls.RevertDestTransferBatch
}

type MockReceiverTransactionSimulator struct {
	GetRevertReasonFunc   func(ctx context.Context, txHash common.Hash) (txsim.ContractError, error)
	DecodeRevertBytesFunc func(data []byte) (txsim.ContractError, error)
	calls                 struct {
		GetRevertReason []struct {
			Ctx    context.Context
			TxHash common.Hash
		}
		DecodeRevertBytes []struct {
			Data []byte
		}
	}
	lock sync.RWMutex
}

func (m *MockReceiverTransactionSimulator) GetRevertReason(ctx context.Context, txHash common.Hash) (txsim.ContractError, error) {
	if m.GetRevertReasonFunc == nil {
		panic("MockReceiverTransactionSimulator.GetRevertReasonFunc: method is nil but GetRevertReason was just called")
	}
	m.lock.Lock()
	m.calls.GetRevertReason = append(m.calls.GetRevertReason, struct {
		Ctx    context.Context
		TxHash common.Hash
	}{Ctx: ctx, TxHash: txHash})
	m.lock.Unlock()
	return m.GetRevertReasonFunc(ctx, txHash)
}

func (m *MockReceiverTransactionSimulator) GetRevertReasonCalls() []struct {
	Ctx    context.Context
	TxHash common.Hash
} {
	m.lock.RLock()
	defer m.lock.RUnlock()
	return m.calls.GetRevertReason
}

func (m *MockReceiverTransactionSimulator) DecodeRevertBytes(data []byte) (txsim.ContractError, error) {
	if m.DecodeRevertBytesFunc == nil {
		panic("MockReceiverTransactionSimulator.DecodeRevertBytesFunc: method is nil but DecodeRevertBytes was just called")
	}
	m.lock.Lock()
	m.calls.DecodeRevertBytes = append(m.calls.DecodeRevertBytes, struct {
		Data []byte
	}{Data: data})
	m.lock.Unlock()
	return m.DecodeRevertBytesFunc(data)
}

func (m *MockReceiverTransactionSimulator) DecodeRevertBytesCalls() []struct {
	Data []byte
} {
	m.lock.RLock()
	defer m.lock.RUnlock()
	return m.calls.DecodeRevertBytes
}

type MockReceiverTeleportClient struct {
	SendTransferCompletedFunc func(ctx context.Context, messages []types.EnygmaTransferCompleted) error
	calls                     struct {
		SendTransferCompleted []struct {
			Ctx      context.Context
			Messages []types.EnygmaTransferCompleted
		}
	}
	lock sync.RWMutex
}

func (m *MockReceiverTeleportClient) SendTransferCompleted(ctx context.Context, messages []types.EnygmaTransferCompleted) error {
	if m.SendTransferCompletedFunc == nil {
		panic("MockReceiverTeleportClient.SendTransferCompletedFunc: method is nil but SendTransferCompleted was just called")
	}
	m.lock.Lock()
	m.calls.SendTransferCompleted = append(m.calls.SendTransferCompleted, struct {
		Ctx      context.Context
		Messages []types.EnygmaTransferCompleted
	}{Ctx: ctx, Messages: messages})
	m.lock.Unlock()
	return m.SendTransferCompletedFunc(ctx, messages)
}

func (m *MockReceiverTeleportClient) SendTransferCompletedCalls() []struct {
	Ctx      context.Context
	Messages []types.EnygmaTransferCompleted
} {
	m.lock.RLock()
	defer m.lock.RUnlock()
	return m.calls.SendTransferCompleted
}

type MockEnygmaCreationService struct {
	CreateEnygmaFunc func(ctx context.Context, resourceId string, fromChainId *big.Int, blockNumberPrivateHub *big.Int) error
	calls            struct {
		CreateEnygma []struct {
			Ctx                   context.Context
			ResourceId            string
			FromChainId           *big.Int
			BlockNumberPrivateHub *big.Int
		}
	}
	lock sync.RWMutex
}

func (m *MockEnygmaCreationService) CreateEnygma(ctx context.Context, resourceId string, fromChainId *big.Int, blockNumberPrivateHub *big.Int) error {
	if m.CreateEnygmaFunc == nil {
		panic("MockEnygmaCreationService.CreateEnygmaFunc: method is nil but CreateEnygma was just called")
	}
	m.lock.Lock()
	m.calls.CreateEnygma = append(m.calls.CreateEnygma, struct {
		Ctx                   context.Context
		ResourceId            string
		FromChainId           *big.Int
		BlockNumberPrivateHub *big.Int
	}{Ctx: ctx, ResourceId: resourceId, FromChainId: fromChainId, BlockNumberPrivateHub: blockNumberPrivateHub})
	m.lock.Unlock()
	return m.CreateEnygmaFunc(ctx, resourceId, fromChainId, blockNumberPrivateHub)
}

func (m *MockEnygmaCreationService) CreateEnygmaCalls() []struct {
	Ctx                   context.Context
	ResourceId            string
	FromChainId           *big.Int
	BlockNumberPrivateHub *big.Int
} {
	m.lock.RLock()
	defer m.lock.RUnlock()
	return m.calls.CreateEnygma
}

// MockReceiverTxRecoveryRepository mocks the receiverTxRecoveryRepository interface.
type MockReceiverTxRecoveryRepository struct {
	InsertFunc                func(ctx context.Context, data types.TxRecoveryData) error
	GetByPrivateHubTxHashFunc func(ctx context.Context, privateHubTxHash string) (*types.TxRecoveryData, error)
	MarkConfirmedFunc         func(ctx context.Context, privateHubTxHash string) error
}

func (m *MockReceiverTxRecoveryRepository) Insert(ctx context.Context, data types.TxRecoveryData) error {
	if m.InsertFunc == nil {
		return nil
	}
	return m.InsertFunc(ctx, data)
}

func (m *MockReceiverTxRecoveryRepository) GetByPrivateHubTxHash(ctx context.Context, privateHubTxHash string) (*types.TxRecoveryData, error) {
	if m.GetByPrivateHubTxHashFunc == nil {
		return nil, nil
	}
	return m.GetByPrivateHubTxHashFunc(ctx, privateHubTxHash)
}

func (m *MockReceiverTxRecoveryRepository) MarkConfirmed(ctx context.Context, privateHubTxHash string) error {
	if m.MarkConfirmedFunc == nil {
		return nil
	}
	return m.MarkConfirmedFunc(ctx, privateHubTxHash)
}

type receiverDeps struct {
	Tracer               *testutils.MockTracer
	EndpointClient       *MockReceiverEndpointClient
	Deployer             *MockReceiverDeployer
	HistoryRepository    *MockReceiverEnygmaHistoryRepository
	TxRecoveryRepository *MockReceiverTxRecoveryRepository
	HandlerClient        *MockEnygmaHandlerClient
	TeleportClient       *MockReceiverTeleportClient
	Simulator            *MockReceiverTransactionSimulator
	CreationServiceMock  *MockEnygmaCreationService
}

func setupReceiverDeps() *receiverDeps {
	return &receiverDeps{
		Tracer:               &testutils.MockTracer{},
		EndpointClient:       &MockReceiverEndpointClient{},
		Deployer:             &MockReceiverDeployer{},
		HistoryRepository:    &MockReceiverEnygmaHistoryRepository{},
		TxRecoveryRepository: &MockReceiverTxRecoveryRepository{},
		HandlerClient:        &MockEnygmaHandlerClient{},
		TeleportClient:       &MockReceiverTeleportClient{},
		Simulator:            &MockReceiverTransactionSimulator{},
		CreationServiceMock: &MockEnygmaCreationService{
			CreateEnygmaFunc: func(ctx context.Context, resourceId string, fromChainId *big.Int, blockNumberPrivateHub *big.Int) error {
				return nil
			},
		},
	}
}

func setupReceiver(deps *receiverDeps) *handler.Receiver {
	return handler.NewReceiver(
		big.NewInt(1),
		deps.EndpointClient,
		deps.Simulator,
		deps.Deployer,
		deps.HistoryRepository,
		deps.TxRecoveryRepository,
		deps.TeleportClient,
		deps.HandlerClient,
		deps.Tracer,
		deps.CreationServiceMock,
		30*time.Second,
	)
}
