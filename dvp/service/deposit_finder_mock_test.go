package service_test

import (
	"context"
	"sync"

	"github.com/raylsnetwork/rayls-sovereign-relayer/types"
)

type finderDepositRepositoryMock struct {
	GetDepositsByTokenFunc    func(ctx context.Context, tokenAddress string, tokenId string, tokenType types.DvpTokenType, userAddress string, status types.DvpDepositStatus) ([]types.DvpDeposit, error)
	GetFungibleDepositsFunc   func(ctx context.Context, tokenAddress string, userAddress string, tokenType types.DvpTokenType, status types.DvpDepositStatus) ([]types.DvpDeposit, error)
	GetNonFungibleDepositFunc func(ctx context.Context, tokenId string, tokenAddress string, userAddress string, tokenType types.DvpTokenType, status types.DvpDepositStatus) (*types.DvpDeposit, error)

	calls struct {
		GetDepositsByToken []struct {
			Ctx          context.Context
			TokenAddress string
			TokenId      string
			TokenType    types.DvpTokenType
			UserAddress  string
			Status       types.DvpDepositStatus
		}
		GetFungibleDeposits []struct {
			Ctx          context.Context
			TokenAddress string
			UserAddress  string
			TokenType    types.DvpTokenType
			Status       types.DvpDepositStatus
		}
		GetNonFungibleDeposit []struct {
			Ctx          context.Context
			TokenAddress string
			TokenId      string
			UserAddress  string
			TokenType    types.DvpTokenType
			Status       types.DvpDepositStatus
		}
	}
	lockGetDepositsByToken    sync.RWMutex
	lockGetFungibleDeposits   sync.RWMutex
	lockGetNonFungibleDeposit sync.RWMutex
}

func (m *finderDepositRepositoryMock) GetDepositsByToken(ctx context.Context, tokenAddress string, tokenId string, tokenType types.DvpTokenType, userAddress string, status types.DvpDepositStatus) ([]types.DvpDeposit, error) {
	m.lockGetDepositsByToken.Lock()
	defer m.lockGetDepositsByToken.Unlock()

	m.calls.GetDepositsByToken = append(m.calls.GetDepositsByToken, struct {
		Ctx          context.Context
		TokenAddress string
		TokenId      string
		TokenType    types.DvpTokenType
		UserAddress  string
		Status       types.DvpDepositStatus
	}{Ctx: ctx, TokenAddress: tokenAddress, TokenId: tokenId, TokenType: tokenType, UserAddress: userAddress, Status: status})

	return m.GetDepositsByTokenFunc(ctx, tokenAddress, tokenId, tokenType, userAddress, status)
}

func (m *finderDepositRepositoryMock) GetFungibleDeposits(ctx context.Context, tokenAddress string, userAddress string, tokenType types.DvpTokenType, status types.DvpDepositStatus) ([]types.DvpDeposit, error) {
	m.lockGetFungibleDeposits.Lock()
	defer m.lockGetFungibleDeposits.Unlock()

	m.calls.GetFungibleDeposits = append(m.calls.GetFungibleDeposits, struct {
		Ctx          context.Context
		TokenAddress string
		UserAddress  string
		TokenType    types.DvpTokenType
		Status       types.DvpDepositStatus
	}{Ctx: ctx, TokenAddress: tokenAddress, UserAddress: userAddress, TokenType: tokenType, Status: status})

	return m.GetFungibleDepositsFunc(ctx, tokenAddress, userAddress, tokenType, status)
}

func (m *finderDepositRepositoryMock) GetDepositsByTokenCalls() []struct {
	Ctx          context.Context
	TokenAddress string
	TokenId      string
	TokenType    types.DvpTokenType
	UserAddress  string
	Status       types.DvpDepositStatus
} {
	m.lockGetDepositsByToken.RLock()
	defer m.lockGetDepositsByToken.RUnlock()
	return m.calls.GetDepositsByToken
}

func (m *finderDepositRepositoryMock) GetFungibleDepositsCalls() []struct {
	Ctx          context.Context
	TokenAddress string
	UserAddress  string
	TokenType    types.DvpTokenType
	Status       types.DvpDepositStatus
} {
	m.lockGetFungibleDeposits.RLock()
	defer m.lockGetFungibleDeposits.RUnlock()
	return m.calls.GetFungibleDeposits
}

func (m *finderDepositRepositoryMock) GetNonFungibleDeposit(ctx context.Context, tokenId string, tokenAddress string, userAddress string, tokenType types.DvpTokenType, status types.DvpDepositStatus) (*types.DvpDeposit, error) {
	m.lockGetNonFungibleDeposit.Lock()
	defer m.lockGetNonFungibleDeposit.Unlock()

	m.calls.GetNonFungibleDeposit = append(m.calls.GetNonFungibleDeposit, struct {
		Ctx          context.Context
		TokenAddress string
		TokenId      string
		UserAddress  string
		TokenType    types.DvpTokenType
		Status       types.DvpDepositStatus
	}{Ctx: ctx, TokenAddress: tokenAddress, TokenId: tokenId, UserAddress: userAddress, TokenType: tokenType, Status: status})

	return m.GetNonFungibleDepositFunc(ctx, tokenId, tokenAddress, userAddress, tokenType, status)
}

func (m *finderDepositRepositoryMock) GetNonFungibleDepositCalls() []struct {
	Ctx          context.Context
	TokenAddress string
	TokenId      string
	UserAddress  string
	TokenType    types.DvpTokenType
	Status       types.DvpDepositStatus
} {
	m.lockGetNonFungibleDeposit.RLock()
	defer m.lockGetNonFungibleDeposit.RUnlock()
	return m.calls.GetNonFungibleDeposit
}
