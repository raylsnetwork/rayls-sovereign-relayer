package service_test

import (
	"context"
	"math/big"
	"sync"

	"github.com/raylsnetwork/rayls-privacy-relayer-api/cryptography"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/types"
)

type MockEnygmaCheckpointRepository struct {
	GetLatestCheckpointByFiltersFunc  func(ctx context.Context, resourceId string, status *types.EnygmaCheckpointStatus, finalizedBlockNumber *big.Int, pendingBlockNumber *big.Int) (*types.EnygmaCheckpoint, error)
	GetLatestCheckpointByFiltersCalls struct {
		sync.RWMutex
		Calls []struct {
			Ctx                  context.Context
			ResourceId           string
			Status               *types.EnygmaCheckpointStatus
			FinalizedBlockNumber *big.Int
			PendingBlockNumber   *big.Int
		}
	}

	CreateEnygmaCheckpointFunc  func(ctx context.Context, checkpoint types.EnygmaCheckpoint) error
	CreateEnygmaCheckpointCalls struct {
		sync.RWMutex
		Calls []struct {
			Ctx        context.Context
			Checkpoint types.EnygmaCheckpoint
		}
	}

	GetValidationCandidatesFunc  func(ctx context.Context) ([]types.EnygmaCheckpoint, error)
	GetValidationCandidatesCalls struct {
		sync.RWMutex
		Calls []struct {
			Ctx context.Context
		}
	}

	MarkAsFinalizedFunc  func(ctx context.Context, resourceId string, finalizedBlockNumber *big.Int) error
	MarkAsFinalizedCalls struct {
		sync.RWMutex
		Calls []struct {
			Ctx                  context.Context
			ResourceId           string
			FinalizedBlockNumber *big.Int
		}
	}
}

func (m *MockEnygmaCheckpointRepository) GetLatestCheckpointByFilters(ctx context.Context, resourceId string, status *types.EnygmaCheckpointStatus, finalizedBlockNumber *big.Int, pendingBlockNumber *big.Int) (*types.EnygmaCheckpoint, error) {
	m.GetLatestCheckpointByFiltersCalls.Lock()
	defer m.GetLatestCheckpointByFiltersCalls.Unlock()
	m.GetLatestCheckpointByFiltersCalls.Calls = append(m.GetLatestCheckpointByFiltersCalls.Calls, struct {
		Ctx                  context.Context
		ResourceId           string
		Status               *types.EnygmaCheckpointStatus
		FinalizedBlockNumber *big.Int
		PendingBlockNumber   *big.Int
	}{
		Ctx:                  ctx,
		ResourceId:           resourceId,
		Status:               status,
		FinalizedBlockNumber: finalizedBlockNumber,
		PendingBlockNumber:   pendingBlockNumber,
	})

	if m.GetLatestCheckpointByFiltersFunc != nil {
		return m.GetLatestCheckpointByFiltersFunc(ctx, resourceId, status, finalizedBlockNumber, pendingBlockNumber)
	}
	return nil, nil //nolint:nilnil // intentional nil return in test mock
}

func (m *MockEnygmaCheckpointRepository) CreateEnygmaCheckpoint(ctx context.Context, checkpoint types.EnygmaCheckpoint) error {
	m.CreateEnygmaCheckpointCalls.Lock()
	defer m.CreateEnygmaCheckpointCalls.Unlock()
	m.CreateEnygmaCheckpointCalls.Calls = append(m.CreateEnygmaCheckpointCalls.Calls, struct {
		Ctx        context.Context
		Checkpoint types.EnygmaCheckpoint
	}{
		Ctx:        ctx,
		Checkpoint: checkpoint,
	})

	if m.CreateEnygmaCheckpointFunc != nil {
		return m.CreateEnygmaCheckpointFunc(ctx, checkpoint)
	}
	return nil
}

func (m *MockEnygmaCheckpointRepository) GetValidationCandidates(ctx context.Context) ([]types.EnygmaCheckpoint, error) {
	m.GetValidationCandidatesCalls.Lock()
	defer m.GetValidationCandidatesCalls.Unlock()
	m.GetValidationCandidatesCalls.Calls = append(m.GetValidationCandidatesCalls.Calls, struct {
		Ctx context.Context
	}{
		Ctx: ctx,
	})

	if m.GetValidationCandidatesFunc != nil {
		return m.GetValidationCandidatesFunc(ctx)
	}
	return nil, nil //nolint:nilnil // intentional nil return in test mock
}

func (m *MockEnygmaCheckpointRepository) MarkAsFinalized(ctx context.Context, resourceId string, finalizedBlockNumber *big.Int) error {
	m.MarkAsFinalizedCalls.Lock()
	defer m.MarkAsFinalizedCalls.Unlock()
	m.MarkAsFinalizedCalls.Calls = append(m.MarkAsFinalizedCalls.Calls, struct {
		Ctx                  context.Context
		ResourceId           string
		FinalizedBlockNumber *big.Int
	}{
		Ctx:                  ctx,
		ResourceId:           resourceId,
		FinalizedBlockNumber: finalizedBlockNumber,
	})

	if m.MarkAsFinalizedFunc != nil {
		return m.MarkAsFinalizedFunc(ctx, resourceId, finalizedBlockNumber)
	}
	return nil
}

func createEnygmaFinalizedBalance(resourceId string) *types.EnygmaFinalizedBalance {
	return &types.EnygmaFinalizedBalance{
		ResourceId:           resourceId,
		FinalizedBlockNumber: big.NewInt(100),
		PendingBlockNumber:   big.NewInt(101),
		BalanceX:             big.NewInt(12345),
		BalanceY:             big.NewInt(67890),
	}
}

func createEnygmaFinalizedBalanceWithValues(resourceId string, balanceX, balanceY, finalizedBlockNum, pendingBlockNum *big.Int) *types.EnygmaFinalizedBalance {
	return &types.EnygmaFinalizedBalance{
		ResourceId:           resourceId,
		FinalizedBlockNumber: finalizedBlockNum,
		PendingBlockNumber:   pendingBlockNum,
		BalanceX:             balanceX,
		BalanceY:             balanceY,
	}
}

func createEnygmaCheckpoint(id, resourceId string, balanceX, balanceY, finalizedBlockNum, pendingBlockNum *big.Int) *types.EnygmaCheckpoint {
	return &types.EnygmaCheckpoint{
		ID:                      id,
		ResourceId:              resourceId,
		FinalizedPublicBalanceX: balanceX,
		FinalizedPublicBalanceY: balanceY,
		FinalizedBlockNumber:    finalizedBlockNum,
		PendingBlockNumber:      pendingBlockNum,
		Status:                  types.EnygmaCheckpointStatusTentative,
	}
}

type MockEnygmaHistoryRepository struct {
	GetEnygmaHistoryForCheckpointsFunc  func(ctx context.Context, resourceIds []string, blockNumbers []*big.Int) ([]types.EnygmaHistory, error)
	GetEnygmaHistoryForCheckpointsCalls struct {
		sync.RWMutex
		Calls []struct {
			Ctx          context.Context
			ResourceIds  []string
			BlockNumbers []*big.Int
		}
	}
}

func (m *MockEnygmaHistoryRepository) GetEnygmaHistoryForCheckpoints(ctx context.Context, resourceIds []string, blockNumbers []*big.Int) ([]types.EnygmaHistory, error) {
	m.GetEnygmaHistoryForCheckpointsCalls.Lock()
	defer m.GetEnygmaHistoryForCheckpointsCalls.Unlock()
	m.GetEnygmaHistoryForCheckpointsCalls.Calls = append(m.GetEnygmaHistoryForCheckpointsCalls.Calls, struct {
		Ctx          context.Context
		ResourceIds  []string
		BlockNumbers []*big.Int
	}{
		Ctx:          ctx,
		ResourceIds:  resourceIds,
		BlockNumbers: blockNumbers,
	})

	if m.GetEnygmaHistoryForCheckpointsFunc != nil {
		return m.GetEnygmaHistoryForCheckpointsFunc(ctx, resourceIds, blockNumbers)
	}
	return nil, nil //nolint:nilnil // intentional nil return in test mock
}

type MockEnygmaRepository struct {
	GetEnygmaByResourceIdsFunc  func(ctx context.Context, resourceIds []string) ([]types.Enygma, error)
	GetEnygmaByResourceIdsCalls struct {
		sync.RWMutex
		Calls []struct {
			Ctx         context.Context
			ResourceIds []string
		}
	}

	UpdateEnygmaFunc  func(ctx context.Context, resourceId string, finalizedBalance, finalizedR, finalizedBlockNumber, pendingBlockNumber *big.Int) error
	UpdateEnygmaCalls struct {
		sync.RWMutex
		Calls []struct {
			Ctx                  context.Context
			ResourceId           string
			FinalizedBalance     *big.Int
			FinalizedR           *big.Int
			FinalizedBlockNumber *big.Int
			PendingBlockNumber   *big.Int
		}
	}
}

func (m *MockEnygmaRepository) GetEnygmaByResourceIds(ctx context.Context, resourceIds []string) ([]types.Enygma, error) {
	m.GetEnygmaByResourceIdsCalls.Lock()
	defer m.GetEnygmaByResourceIdsCalls.Unlock()
	m.GetEnygmaByResourceIdsCalls.Calls = append(m.GetEnygmaByResourceIdsCalls.Calls, struct {
		Ctx         context.Context
		ResourceIds []string
	}{
		Ctx:         ctx,
		ResourceIds: resourceIds,
	})

	if m.GetEnygmaByResourceIdsFunc != nil {
		return m.GetEnygmaByResourceIdsFunc(ctx, resourceIds)
	}
	return nil, nil //nolint:nilnil // intentional nil return in test mock
}

func (m *MockEnygmaRepository) UpdateEnygma(ctx context.Context, resourceId string, finalizedBalance, finalizedR, finalizedBlockNumber, pendingBlockNumber *big.Int) error {
	m.UpdateEnygmaCalls.Lock()
	defer m.UpdateEnygmaCalls.Unlock()
	m.UpdateEnygmaCalls.Calls = append(m.UpdateEnygmaCalls.Calls, struct {
		Ctx                  context.Context
		ResourceId           string
		FinalizedBalance     *big.Int
		FinalizedR           *big.Int
		FinalizedBlockNumber *big.Int
		PendingBlockNumber   *big.Int
	}{
		Ctx:                  ctx,
		ResourceId:           resourceId,
		FinalizedBalance:     finalizedBalance,
		FinalizedR:           finalizedR,
		FinalizedBlockNumber: finalizedBlockNumber,
		PendingBlockNumber:   pendingBlockNumber,
	})

	if m.UpdateEnygmaFunc != nil {
		return m.UpdateEnygmaFunc(ctx, resourceId, finalizedBalance, finalizedR, finalizedBlockNumber, pendingBlockNumber)
	}
	return nil
}

type MockTransactionManager struct {
	WithTransactionFunc  func(ctx context.Context, fn func(ctx context.Context) error) error
	WithTransactionCalls struct {
		sync.RWMutex
		Calls []struct {
			Ctx context.Context
		}
	}
}

func (m *MockTransactionManager) WithTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
	m.WithTransactionCalls.Lock()
	defer m.WithTransactionCalls.Unlock()
	m.WithTransactionCalls.Calls = append(m.WithTransactionCalls.Calls, struct {
		Ctx context.Context
	}{
		Ctx: ctx,
	})

	if m.WithTransactionFunc != nil {
		return m.WithTransactionFunc(ctx, fn)
	}
	// Call fn for testing purposes
	return fn(ctx)
}

type MockResyncService struct {
	ResyncEnygmaFunc  func(ctx context.Context, resourceId string) error
	ResyncEnygmaCalls struct {
		sync.RWMutex
		Calls []struct {
			Ctx        context.Context
			ResourceId string
		}
	}
}

func (m *MockResyncService) ResyncEnygma(ctx context.Context, resourceId string) error {
	m.ResyncEnygmaCalls.Lock()
	defer m.ResyncEnygmaCalls.Unlock()
	m.ResyncEnygmaCalls.Calls = append(m.ResyncEnygmaCalls.Calls, struct {
		Ctx        context.Context
		ResourceId string
	}{
		Ctx:        ctx,
		ResourceId: resourceId,
	})

	if m.ResyncEnygmaFunc != nil {
		return m.ResyncEnygmaFunc(ctx, resourceId)
	}
	return nil
}

func createEnygmaHistory(resourceId string, balanceChange, rFactor *big.Int) *types.EnygmaHistory {
	return &types.EnygmaHistory{
		ResourceId:            resourceId,
		FromChainId:           big.NewInt(1),
		BlockNumberPrivateHub: big.NewInt(100),
		BalanceChange:         balanceChange,
		RFactor:               rFactor,
		EventType:             types.EnygmaTransfer,
	}
}

func createEnygma(resourceId string, finalizedBalance, finalizedR, finalizedBlockNum, pendingBlockNum *big.Int) *types.Enygma {
	return &types.Enygma{
		ResourceId:           resourceId,
		FinalizedBalance:     finalizedBalance,
		FinalizedR:           finalizedR,
		FinalizedBlockNumber: finalizedBlockNum,
		PendingBlockNumber:   pendingBlockNum,
	}
}

// createValidCheckpointWithCryptography generates checkpoint data with valid Pedersen commitment
// Returns: checkpoint, enygma (previous state), history, newBalance, newR
func createValidCheckpointWithCryptography() (*types.EnygmaCheckpoint, *types.Enygma, *types.EnygmaHistory, *big.Int, *big.Int) {
	prevBalance := big.NewInt(100)
	prevR := big.NewInt(0)

	balanceChange := big.NewInt(50)
	rChange := big.NewInt(50)

	// New finalized values: balance=150, r=50
	newBalance := new(big.Int).Add(prevBalance, balanceChange)
	newR := cryptography.AddMod(prevR, rChange, cryptography.JubJubPrimeSubGroup)

	commitment := cryptography.PedersenCommitmentEnygma(newBalance, newR)

	checkpoint := &types.EnygmaCheckpoint{
		ID:                      "id1",
		ResourceId:              "resource1",
		FinalizedPublicBalanceX: commitment.X,
		FinalizedPublicBalanceY: commitment.Y,
		FinalizedBlockNumber:    big.NewInt(100),
		PendingBlockNumber:      big.NewInt(101),
		Status:                  types.EnygmaCheckpointStatusTentative,
	}

	enygma := &types.Enygma{
		ResourceId:           "resource1",
		FinalizedBalance:     prevBalance,
		FinalizedR:           prevR,
		FinalizedBlockNumber: big.NewInt(99),
		PendingBlockNumber:   big.NewInt(100),
	}

	history := &types.EnygmaHistory{
		ResourceId:            "resource1",
		FromChainId:           big.NewInt(1),
		BlockNumberPrivateHub: big.NewInt(100),
		BalanceChange:         balanceChange,
		RFactor:               rChange,
		EventType:             types.EnygmaTransfer,
	}

	return checkpoint, enygma, history, newBalance, newR
}

// createCheckpointWithValidationFailure generates checkpoint data with mismatched balance/R values
// to intentionally fail validation. Useful for testing retry and resync logic.
func createCheckpointWithValidationFailure() (*types.EnygmaCheckpoint, *types.Enygma, *types.EnygmaHistory) {
	prevBalance := big.NewInt(100)
	prevR := big.NewInt(0)

	balanceChange := big.NewInt(50)
	rChange := big.NewInt(50)

	// Create checkpoint with INTENTIONALLY WRONG commitment (to fail validation)
	// (We don't use newBalance/newR here because we want validation to fail)
	checkpoint := &types.EnygmaCheckpoint{
		ID:                      "id1",
		ResourceId:              "resource1",
		FinalizedPublicBalanceX: big.NewInt(999), // Wrong value - will fail validation
		FinalizedPublicBalanceY: big.NewInt(999), // Wrong value - will fail validation
		FinalizedBlockNumber:    big.NewInt(100),
		PendingBlockNumber:      big.NewInt(101),
		Status:                  types.EnygmaCheckpointStatusTentative,
	}

	enygma := &types.Enygma{
		ResourceId:           "resource1",
		FinalizedBalance:     prevBalance,
		FinalizedR:           prevR,
		FinalizedBlockNumber: big.NewInt(99),
		PendingBlockNumber:   big.NewInt(100),
	}

	history := &types.EnygmaHistory{
		ResourceId:            "resource1",
		FromChainId:           big.NewInt(1),
		BlockNumberPrivateHub: big.NewInt(100),
		BalanceChange:         balanceChange,
		RFactor:               rChange,
		EventType:             types.EnygmaTransfer,
	}

	return checkpoint, enygma, history
}
