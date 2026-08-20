package service

import (
	"context"
	"fmt"
	"math/big"

	"github.com/raylsnetwork/rayls-privacy-relayer-api/private-relayer/repository"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/types"
)

type finderDepositRepository interface {
	GetDepositsByToken(
		ctx context.Context,
		tokenAddress string,
		tokenId string,
		tokenType types.DvpTokenType,
		userAddress string,
		status types.DvpDepositStatus,
	) ([]types.DvpDeposit, error)
	GetFungibleDeposits(
		ctx context.Context,
		tokenAddress string,
		userAddress string,
		tokenType types.DvpTokenType,
		status types.DvpDepositStatus,
	) ([]types.DvpDeposit, error)
	GetNonFungibleDeposit(
		ctx context.Context,
		tokenId string,
		tokenAddress string,
		userAddress string,
		tokenType types.DvpTokenType,
		status types.DvpDepositStatus,
	) (*types.DvpDeposit, error)
}

var _ finderDepositRepository = (*repository.DvpDepositRepository)(nil)

type DepositFinder struct {
	depositRepository finderDepositRepository
}

func NewDepositFinder(depositRepository finderDepositRepository) *DepositFinder {
	return &DepositFinder{
		depositRepository: depositRepository,
	}
}

// TODO: can we do all of this with a single DB query, instead of doing it in memory?
func (s *DepositFinder) FindERC1155DepositsForJSProof(
	ctx context.Context,
	userAddress string,
	tokenAddress string,
	tokenId string,
	paymentAmount *big.Int,
) ([]*types.DvpDeposit, error) {
	deposits, err := s.depositRepository.GetDepositsByToken(
		ctx,
		tokenAddress,
		tokenId,
		types.DvpERC1155,
		userAddress,
		types.DvpDepositUnspent,
	)
	if err != nil {
		return nil, fmt.Errorf("getting ERC1155 deposits by token: %w", err)
	}

	jsDeposits := make([]*types.DvpDeposit, 0)
	jsTotalAmount := big.NewInt(0)

	for _, deposit := range deposits {
		if jsTotalAmount.Cmp(paymentAmount) >= 0 {
			break
		}

		jsDeposits = append(jsDeposits, &deposit)
		jsTotalAmount = jsTotalAmount.Add(jsTotalAmount, deposit.TokenAmount)
	}

	if jsTotalAmount.Cmp(paymentAmount) < 0 {
		return nil, fmt.Errorf(
			"not enough deposits. Found %d deposits with total amount %s | required %s",
			len(jsDeposits),
			jsTotalAmount.String(),
			paymentAmount.String(),
		)
	}

	return jsDeposits, nil
}

// TODO: can we do all of this with a single DB query, instead of doing it in memory?
func (s *DepositFinder) FindEnygmaDeposits(
	ctx context.Context,
	userAddress string,
	tokenAddress string,
	paymentAmount *big.Int,
) ([]*types.DvpDeposit, error) {
	deposits, err := s.depositRepository.GetFungibleDeposits(
		ctx,
		tokenAddress,
		userAddress,
		types.DvpEnygma,
		types.DvpDepositUnspent,
	)
	if err != nil {
		return nil, fmt.Errorf("getting fungible enygma deposits: %w", err)
	}

	jsDeposits := make([]*types.DvpDeposit, 0)
	jsTotalAmount := big.NewInt(0)

	for _, deposit := range deposits {
		if jsTotalAmount.Cmp(paymentAmount) >= 0 {
			break
		}

		jsDeposits = append(jsDeposits, &deposit)
		jsTotalAmount = jsTotalAmount.Add(jsTotalAmount, deposit.TokenAmount)
	}

	if jsTotalAmount.Cmp(paymentAmount) < 0 {
		return nil, fmt.Errorf(
			"not enough deposits. Found %d deposits with total amount %s | required %s",
			len(jsDeposits),
			jsTotalAmount.String(),
			paymentAmount.String(),
		)
	}

	return jsDeposits, nil
}

func (s *DepositFinder) FindERC721Deposit(
	ctx context.Context,
	userAddress string,
	tokenAddress string,
	tokenId string,
) (*types.DvpDeposit, error) {
	deposit, err := s.depositRepository.GetNonFungibleDeposit(
		ctx,
		tokenId,
		tokenAddress,
		userAddress,
		types.DvpERC721,
		types.DvpDepositUnspent,
	)
	if err != nil {
		return nil, fmt.Errorf("getting non-fungible ERC721 deposit: %w", err)
	}

	if deposit == nil {
		return nil, fmt.Errorf("ERC721 deposit not found")
	}

	return deposit, nil
}
