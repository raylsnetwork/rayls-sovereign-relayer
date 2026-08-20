package service_test

import (
	"context"
	"math/big"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/cts/gen/proto/keys"
	"google.golang.org/grpc"

	"github.com/raylsnetwork/rayls-privacy-relayer-api/types"
)

var _ interface {
	GetPaymentSpendKey(ctx context.Context, in *keys.GetPaymentSpendKeyRequest, opts ...grpc.CallOption) (*keys.PaymentSpendKeyResponse, error)
	GenerateEnygmaSharedSecrets(ctx context.Context, in *keys.GenerateEnygmaSharedSecretsRequest, opts ...grpc.CallOption) (*keys.GenerateEnygmaSharedSecretsResponse, error)
} = (*MockProofKOSClient)(nil)

type MockProofKOSClient struct {
	GetPaymentSpendKeyFunc          func(ctx context.Context, in *keys.GetPaymentSpendKeyRequest, opts ...grpc.CallOption) (*keys.PaymentSpendKeyResponse, error)
	GenerateEnygmaSharedSecretsFunc func(ctx context.Context, in *keys.GenerateEnygmaSharedSecretsRequest, opts ...grpc.CallOption) (*keys.GenerateEnygmaSharedSecretsResponse, error)
	calls                           struct {
		GetPaymentSpendKey []struct {
			Ctx context.Context
			In  *keys.GetPaymentSpendKeyRequest
		}
		GenerateEnygmaSharedSecrets []struct {
			Ctx context.Context
			In  *keys.GenerateEnygmaSharedSecretsRequest
		}
	}
	lock sync.RWMutex
}

func (m *MockProofKOSClient) GetPaymentSpendKey(ctx context.Context, in *keys.GetPaymentSpendKeyRequest, opts ...grpc.CallOption) (*keys.PaymentSpendKeyResponse, error) {
	m.lock.Lock()
	m.calls.GetPaymentSpendKey = append(m.calls.GetPaymentSpendKey, struct {
		Ctx context.Context
		In  *keys.GetPaymentSpendKeyRequest
	}{Ctx: ctx, In: in})
	m.lock.Unlock()

	if m.GetPaymentSpendKeyFunc == nil {
		return &keys.PaymentSpendKeyResponse{}, nil
	}
	return m.GetPaymentSpendKeyFunc(ctx, in, opts...)
}

func (m *MockProofKOSClient) GetPaymentSpendKeyCalls() []struct {
	Ctx context.Context
	In  *keys.GetPaymentSpendKeyRequest
} {
	m.lock.RLock()
	defer m.lock.RUnlock()
	return m.calls.GetPaymentSpendKey
}

func (m *MockProofKOSClient) GenerateEnygmaSharedSecrets(ctx context.Context, in *keys.GenerateEnygmaSharedSecretsRequest, opts ...grpc.CallOption) (*keys.GenerateEnygmaSharedSecretsResponse, error) {
	m.lock.Lock()
	m.calls.GenerateEnygmaSharedSecrets = append(m.calls.GenerateEnygmaSharedSecrets, struct {
		Ctx context.Context
		In  *keys.GenerateEnygmaSharedSecretsRequest
	}{Ctx: ctx, In: in})
	m.lock.Unlock()

	if m.GenerateEnygmaSharedSecretsFunc == nil {
		return &keys.GenerateEnygmaSharedSecretsResponse{}, nil
	}
	return m.GenerateEnygmaSharedSecretsFunc(ctx, in, opts...)
}

func (m *MockProofKOSClient) GenerateEnygmaSharedSecretsCalls() []struct {
	Ctx context.Context
	In  *keys.GenerateEnygmaSharedSecretsRequest
} {
	m.lock.RLock()
	defer m.lock.RUnlock()
	return m.calls.GenerateEnygmaSharedSecrets
}

type MockProofAPIClient struct {
	CreateTransferProofFunc func(anonymityIndex int, request types.TransferProofRequest) (types.EnygmaProofResponse, error)
	CreateDepositProofFunc  func(anonymityIndex int, request types.DepositProofRequest) (types.EnygmaProofResponse, error)
	CreateWithdrawProofFunc func(anonymityIndex int, request types.WithdrawProofRequest) (types.EnygmaProofResponse, error)

	calls struct {
		CreateTransferProof []struct {
			AnonymityIndex int
			Request        types.TransferProofRequest
		}
		CreateDepositProof []struct {
			AnonymityIndex int
			Request        types.DepositProofRequest
		}
		CreateWithdrawProof []struct {
			AnonymityIndex int
			Request        types.WithdrawProofRequest
		}
	}
	lock sync.RWMutex
}

func (m *MockProofAPIClient) CreateTransferProof(anonymityIndex int, request types.TransferProofRequest) (types.EnygmaProofResponse, error) {
	if m.CreateTransferProofFunc == nil {
		panic("MockProofAPIClient.CreateTransferProofFunc: method is nil but CreateTransferProof was just called")
	}
	m.lock.Lock()
	m.calls.CreateTransferProof = append(m.calls.CreateTransferProof, struct {
		AnonymityIndex int
		Request        types.TransferProofRequest
	}{AnonymityIndex: anonymityIndex, Request: request})
	m.lock.Unlock()
	return m.CreateTransferProofFunc(anonymityIndex, request)
}

func (m *MockProofAPIClient) CreateTransferProofCalls() []struct {
	AnonymityIndex int
	Request        types.TransferProofRequest
} {
	m.lock.RLock()
	defer m.lock.RUnlock()
	return m.calls.CreateTransferProof
}

func (m *MockProofAPIClient) CreateDepositProof(anonymityIndex int, request types.DepositProofRequest) (types.EnygmaProofResponse, error) {
	if m.CreateDepositProofFunc == nil {
		panic("MockProofAPIClient.CreateDepositProofFunc: method is nil but CreateDepositProof was just called")
	}
	m.lock.Lock()
	m.calls.CreateDepositProof = append(m.calls.CreateDepositProof, struct {
		AnonymityIndex int
		Request        types.DepositProofRequest
	}{AnonymityIndex: anonymityIndex, Request: request})
	m.lock.Unlock()
	return m.CreateDepositProofFunc(anonymityIndex, request)
}

func (m *MockProofAPIClient) CreateDepositProofCalls() []struct {
	AnonymityIndex int
	Request        types.DepositProofRequest
} {
	m.lock.RLock()
	defer m.lock.RUnlock()
	return m.calls.CreateDepositProof
}

func (m *MockProofAPIClient) CreateWithdrawProof(anonymityIndex int, request types.WithdrawProofRequest) (types.EnygmaProofResponse, error) {
	if m.CreateWithdrawProofFunc == nil {
		panic("MockProofAPIClient.CreateWithdrawProofFunc: method is nil but CreateWithdrawProof was just called")
	}
	m.lock.Lock()
	m.calls.CreateWithdrawProof = append(m.calls.CreateWithdrawProof, struct {
		AnonymityIndex int
		Request        types.WithdrawProofRequest
	}{AnonymityIndex: anonymityIndex, Request: request})
	m.lock.Unlock()
	return m.CreateWithdrawProofFunc(anonymityIndex, request)
}

func (m *MockProofAPIClient) CreateWithdrawProofCalls() []struct {
	AnonymityIndex int
	Request        types.WithdrawProofRequest
} {
	m.lock.RLock()
	defer m.lock.RUnlock()
	return m.calls.CreateWithdrawProof
}

type MockProofEnygmaRepository struct {
	GetEnygmaByResourceIdFunc func(ctx context.Context, resourceId string) (types.Enygma, error)

	calls struct {
		GetEnygmaByResourceId []struct {
			Ctx        context.Context
			ResourceId string
		}
	}
	lock sync.RWMutex
}

func (m *MockProofEnygmaRepository) GetEnygmaByResourceId(ctx context.Context, resourceId string) (types.Enygma, error) {
	if m.GetEnygmaByResourceIdFunc == nil {
		panic("MockProofEnygmaRepository.GetEnygmaByResourceIdFunc: method is nil but GetEnygmaByResourceId was just called")
	}
	m.lock.Lock()
	m.calls.GetEnygmaByResourceId = append(m.calls.GetEnygmaByResourceId, struct {
		Ctx        context.Context
		ResourceId string
	}{Ctx: ctx, ResourceId: resourceId})
	m.lock.Unlock()
	return m.GetEnygmaByResourceIdFunc(ctx, resourceId)
}

func (m *MockProofEnygmaRepository) GetEnygmaByResourceIdCalls() []struct {
	Ctx        context.Context
	ResourceId string
} {
	m.lock.RLock()
	defer m.lock.RUnlock()
	return m.calls.GetEnygmaByResourceId
}

type MockProofEncryptor struct {
	GetEnygmaSharedSecretsFunc func(participantsIDs []*big.Int, chainID *big.Int, blockNumber *big.Int) ([]*big.Int, []*big.Int, error)

	calls struct {
		GetEnygmaSharedSecrets []struct {
			ParticipantsIDs []*big.Int
			ChainID         *big.Int
			BlockNumber     *big.Int
		}
	}
	lock sync.RWMutex
}

func (m *MockProofEncryptor) GetEnygmaSharedSecrets(participantsIDs []*big.Int, chainID *big.Int, blockNumber *big.Int) ([]*big.Int, []*big.Int, error) {
	m.lock.Lock()
	m.calls.GetEnygmaSharedSecrets = append(m.calls.GetEnygmaSharedSecrets, struct {
		ParticipantsIDs []*big.Int
		ChainID         *big.Int
		BlockNumber     *big.Int
	}{ParticipantsIDs: participantsIDs, ChainID: chainID, BlockNumber: blockNumber})
	m.lock.Unlock()

	if m.GetEnygmaSharedSecretsFunc == nil {
		// Return default empty secrets
		emptySecrets := make([]*big.Int, len(participantsIDs))
		return emptySecrets, emptySecrets, nil
	}
	return m.GetEnygmaSharedSecretsFunc(participantsIDs, chainID, blockNumber)
}

func (m *MockProofEncryptor) GetEnygmaSharedSecretsCalls() []struct {
	ParticipantsIDs []*big.Int
	ChainID         *big.Int
	BlockNumber     *big.Int
} {
	m.lock.RLock()
	defer m.lock.RUnlock()
	return m.calls.GetEnygmaSharedSecrets
}

type MockEnygmaClient struct {
	GetPublicValuesFinalisedFunc func(ctx context.Context, tokenAddress common.Address) (*types.EnygmaPublicValues, error)

	calls struct {
		GetPublicValuesFinalised []struct {
			Ctx          context.Context
			TokenAddress common.Address
		}
	}
	lock sync.RWMutex
}

func (m *MockEnygmaClient) GetPublicValuesFinalised(ctx context.Context, tokenAddress common.Address) (*types.EnygmaPublicValues, error) {
	if m.GetPublicValuesFinalisedFunc == nil {
		panic("MockEnygmaClient.GetPublicValuesFinalisedFunc: method is nil but GetPublicValuesFinalised was just called")
	}
	m.lock.Lock()
	m.calls.GetPublicValuesFinalised = append(m.calls.GetPublicValuesFinalised, struct {
		Ctx          context.Context
		TokenAddress common.Address
	}{Ctx: ctx, TokenAddress: tokenAddress})
	m.lock.Unlock()
	return m.GetPublicValuesFinalisedFunc(ctx, tokenAddress)
}

func (m *MockEnygmaClient) GetPublicValuesFinalisedCalls() []struct {
	Ctx          context.Context
	TokenAddress common.Address
} {
	m.lock.RLock()
	defer m.lock.RUnlock()
	return m.calls.GetPublicValuesFinalised
}
