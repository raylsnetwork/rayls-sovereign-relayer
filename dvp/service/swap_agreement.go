package service

import (
	"context"
	"errors"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/raylsnetwork/rayls-sovereign-relayer/types"
)

type swapAgreementEndpointClient interface {
	GetResourceAddress(ctx context.Context, resourceId string) (common.Address, error)
}

type SwapAgreementEnygmaHandlerClient interface {
	NotifySenderWithPNCommunicator(
		ctx context.Context,
		tokenAddress common.Address,
		sharedId string,
		status types.DvpCommunicatorStatus,
		message string,
	) error
}

type SwapAgreementDvpERC721HandlerClient interface {
	NotifySenderWithPNCommunicator(
		ctx context.Context,
		tokenAddress common.Address,
		sharedId string,
		status types.DvpCommunicatorStatus,
		message string,
	) error
}

type SwapAgreementDvpERC1155HandlerClient interface {
	NotifySenderWithPNCommunicator(
		ctx context.Context,
		tokenAddress common.Address,
		sharedId string,
		status types.DvpCommunicatorStatus,
		message string,
	) error
}

type SwapAgreement struct {
	plEndpointClient     swapAgreementEndpointClient
	enygmaHandlerClient  SwapAgreementEnygmaHandlerClient
	erc721HandlerClient  SwapAgreementDvpERC721HandlerClient
	erc1155HandlerClient SwapAgreementDvpERC1155HandlerClient
}

func NewSwapAgreement(
	plEndpointClient swapAgreementEndpointClient,
	enygmaHandlerClient SwapAgreementEnygmaHandlerClient,
	erc721HandlerClient SwapAgreementDvpERC721HandlerClient,
	erc1155HandlerClient SwapAgreementDvpERC1155HandlerClient,
) *SwapAgreement {
	return &SwapAgreement{
		plEndpointClient:     plEndpointClient,
		enygmaHandlerClient:  enygmaHandlerClient,
		erc721HandlerClient:  erc721HandlerClient,
		erc1155HandlerClient: erc1155HandlerClient,
	}
}

func (s *SwapAgreement) Verify(
	ctx context.Context,
	swap *types.DvpSwap,
	destChainId *big.Int,
	tokenInResourceId string,
	tokenInAmount *big.Int,
	tokenInId string,
	tokenInType types.DvpTokenType,
	tokenOutResourceId string,
	tokenOutAmount *big.Int,
	tokenOutId string,
	tokenOutType types.DvpTokenType,
) (string, error) {
	if swap.DestChainID.Cmp(destChainId) != 0 {
		return fmt.Sprintf(
				"Destination chain ID doesn't match. Expected: %s, Got: %s",
				swap.DestChainID.String(),
				destChainId.String(),
			), errors.New(
				"ChainID mismatch",
			)
	}
	if swap.TokenInResourceID != tokenInResourceId {
		return fmt.Sprintf(
				"Token in resourceID doesn't match. Expected: %s, Got: %s",
				tokenInResourceId,
				swap.TokenInResourceID,
			), errors.New(
				"TokenInResourceID mismatch",
			)
	}
	if swap.TokenInAmount.Cmp(tokenInAmount) != 0 {
		return fmt.Sprintf(
				"Token in amount doesn't match. Expected: %s, Got: %s",
				tokenInAmount.String(),
				swap.TokenInAmount.String(),
			), errors.New(
				"TokenInAmount mismatch",
			)
	}
	if swap.TokenInID != tokenInId {
		return fmt.Sprintf(
				"Token in ID doesn't match. Expected: %s, Got: %s",
				tokenInId,
				swap.TokenInID,
			), errors.New(
				"TokenInID mismatch",
			)
	}
	if swap.TokenInType != tokenInType {
		return fmt.Sprintf(
				"Token in type doesn't match. Expected: %d, Got: %d",
				tokenInType,
				swap.TokenInType,
			), errors.New(
				"TokenInType mismatch",
			)
	}
	if swap.TokenOutResourceID != tokenOutResourceId {
		return fmt.Sprintf(
				"Token out resourceID doesn't match. Expected: %s, Got: %s",
				tokenOutResourceId,
				swap.TokenOutResourceID,
			), errors.New(
				"TokenOutResourceID mismatch",
			)
	}
	if swap.TokenOutAmount.Cmp(tokenOutAmount) != 0 {
		return fmt.Sprintf(
				"Token out amount doesn't match. Expected: %s, Got: %s",
				tokenOutAmount.String(),
				swap.TokenOutAmount.String(),
			), errors.New(
				"TokenOutAmount mismatch",
			)
	}
	if swap.TokenOutID != tokenOutId {
		return fmt.Sprintf(
				"Token out ID doesn't match. Expected: %s, Got: %s",
				tokenOutId,
				swap.TokenOutID,
			), errors.New(
				"TokenOutID mismatch",
			)
	}
	if swap.TokenOutType != tokenOutType {
		return fmt.Sprintf(
				"Token out type doesn't match. Expected: %d, Got: %d",
				tokenOutType,
				swap.TokenOutType,
			), errors.New(
				"TokenOutType mismatch",
			)
	}

	return "", nil
}

func (s *SwapAgreement) HandleSwapDisagreement(
	ctx context.Context,
	sharedId string,
	tokenInResourceId string,
	tokenInType types.DvpTokenType,
	message string,
) error {
	tokenInAddress, err := s.plEndpointClient.GetResourceAddress(ctx, tokenInResourceId)
	if err != nil {
		return fmt.Errorf("getting resource address for error notification: %w", err)
	}

	switch tokenInType {
	case types.DvpEnygma:
		return s.enygmaHandlerClient.NotifySenderWithPNCommunicator(ctx, tokenInAddress, sharedId, types.SwapError, message)
	case types.DvpERC721:
		return s.erc721HandlerClient.NotifySenderWithPNCommunicator(ctx, tokenInAddress, sharedId, types.SwapError, message)
	case types.DvpERC1155:
		return s.erc1155HandlerClient.NotifySenderWithPNCommunicator(ctx, tokenInAddress, sharedId, types.SwapError, message)
	default:
		return fmt.Errorf("unsupported token type: %d", tokenInType)
	}
}
