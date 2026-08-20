package service_test

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/dvp/service"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSwapAgreement_Verify(t *testing.T) {
	t.Run("successfully verifies when all parameters match", func(t *testing.T) {
		ctx := context.Background()
		swap := createTestDvpSwap()

		svc := service.NewSwapAgreement(nil, nil, nil, nil)

		message, err := svc.Verify(
			ctx,
			swap,
			swap.DestChainID,
			swap.TokenInResourceID,
			swap.TokenInAmount,
			swap.TokenInID,
			swap.TokenInType,
			swap.TokenOutResourceID,
			swap.TokenOutAmount,
			swap.TokenOutID,
			swap.TokenOutType,
		)

		require.NoError(t, err)
		assert.Empty(t, message)
	})

	t.Run("returns error when DestChainID does not match", func(t *testing.T) {
		ctx := context.Background()
		swap := createTestDvpSwap()
		swap.DestChainID = big.NewInt(1)
		differentChainID := big.NewInt(2)

		svc := service.NewSwapAgreement(nil, nil, nil, nil)

		message, err := svc.Verify(
			ctx,
			swap,
			differentChainID,
			swap.TokenInResourceID,
			swap.TokenInAmount,
			swap.TokenInID,
			swap.TokenInType,
			swap.TokenOutResourceID,
			swap.TokenOutAmount,
			swap.TokenOutID,
			swap.TokenOutType,
		)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "ChainID mismatch")
		assert.Contains(t, message, "Destination chain ID doesn't match")
		assert.Contains(t, message, "Expected: 1")
		assert.Contains(t, message, "Got: 2")
	})

	t.Run("returns error when TokenInResourceID does not match", func(t *testing.T) {
		ctx := context.Background()
		swap := createTestDvpSwap()
		swap.TokenInResourceID = "resource1"
		differentResourceID := "resource2"

		svc := service.NewSwapAgreement(nil, nil, nil, nil)

		message, err := svc.Verify(
			ctx,
			swap,
			swap.DestChainID,
			differentResourceID,
			swap.TokenInAmount,
			swap.TokenInID,
			swap.TokenInType,
			swap.TokenOutResourceID,
			swap.TokenOutAmount,
			swap.TokenOutID,
			swap.TokenOutType,
		)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "TokenInResourceID mismatch")
		assert.Contains(t, message, "Token in resourceID doesn't match")
		assert.Contains(t, message, "Expected: resource2")
		assert.Contains(t, message, "Got: resource1")
	})

	t.Run("returns error when TokenInAmount does not match", func(t *testing.T) {
		ctx := context.Background()
		swap := createTestDvpSwap()
		swap.TokenInAmount = big.NewInt(100)
		differentAmount := big.NewInt(200)

		svc := service.NewSwapAgreement(nil, nil, nil, nil)

		message, err := svc.Verify(
			ctx,
			swap,
			swap.DestChainID,
			swap.TokenInResourceID,
			differentAmount,
			swap.TokenInID,
			swap.TokenInType,
			swap.TokenOutResourceID,
			swap.TokenOutAmount,
			swap.TokenOutID,
			swap.TokenOutType,
		)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "TokenInAmount mismatch")
		assert.Contains(t, message, "Token in amount doesn't match")
		assert.Contains(t, message, "Expected: 200")
		assert.Contains(t, message, "Got: 100")
	})

	t.Run("returns error when TokenInID does not match", func(t *testing.T) {
		ctx := context.Background()
		swap := createTestDvpSwap()
		swap.TokenInID = "id1"
		differentID := "id2"

		svc := service.NewSwapAgreement(nil, nil, nil, nil)

		message, err := svc.Verify(
			ctx,
			swap,
			swap.DestChainID,
			swap.TokenInResourceID,
			swap.TokenInAmount,
			differentID,
			swap.TokenInType,
			swap.TokenOutResourceID,
			swap.TokenOutAmount,
			swap.TokenOutID,
			swap.TokenOutType,
		)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "TokenInID mismatch")
		assert.Contains(t, message, "Token in ID doesn't match")
		assert.Contains(t, message, "Expected: id2")
		assert.Contains(t, message, "Got: id1")
	})

	t.Run("returns error when TokenInType does not match", func(t *testing.T) {
		ctx := context.Background()
		swap := createTestDvpSwap()
		swap.TokenInType = types.DvpEnygma
		differentType := types.DvpERC721

		svc := service.NewSwapAgreement(nil, nil, nil, nil)

		message, err := svc.Verify(
			ctx,
			swap,
			swap.DestChainID,
			swap.TokenInResourceID,
			swap.TokenInAmount,
			swap.TokenInID,
			differentType,
			swap.TokenOutResourceID,
			swap.TokenOutAmount,
			swap.TokenOutID,
			swap.TokenOutType,
		)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "TokenInType mismatch")
		assert.Contains(t, message, "Token in type doesn't match")
		assert.Contains(t, message, fmt.Sprintf("Expected: %d", differentType))
		assert.Contains(t, message, fmt.Sprintf("Got: %d", swap.TokenInType))
	})

	t.Run("returns error when TokenOutResourceID does not match", func(t *testing.T) {
		ctx := context.Background()
		swap := createTestDvpSwap()
		swap.TokenOutResourceID = "resource1"
		differentResourceID := "resource2"

		svc := service.NewSwapAgreement(nil, nil, nil, nil)

		message, err := svc.Verify(
			ctx,
			swap,
			swap.DestChainID,
			swap.TokenInResourceID,
			swap.TokenInAmount,
			swap.TokenInID,
			swap.TokenInType,
			differentResourceID,
			swap.TokenOutAmount,
			swap.TokenOutID,
			swap.TokenOutType,
		)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "TokenOutResourceID mismatch")
		assert.Contains(t, message, "Token out resourceID doesn't match")
		assert.Contains(t, message, "Expected: resource2")
		assert.Contains(t, message, "Got: resource1")
	})

	t.Run("returns error when TokenOutAmount does not match", func(t *testing.T) {
		ctx := context.Background()
		swap := createTestDvpSwap()
		swap.TokenOutAmount = big.NewInt(100)
		differentAmount := big.NewInt(200)

		svc := service.NewSwapAgreement(nil, nil, nil, nil)

		message, err := svc.Verify(
			ctx,
			swap,
			swap.DestChainID,
			swap.TokenInResourceID,
			swap.TokenInAmount,
			swap.TokenInID,
			swap.TokenInType,
			swap.TokenOutResourceID,
			differentAmount,
			swap.TokenOutID,
			swap.TokenOutType,
		)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "TokenOutAmount mismatch")
		assert.Contains(t, message, "Token out amount doesn't match")
		assert.Contains(t, message, "Expected: 200")
		assert.Contains(t, message, "Got: 100")
	})

	t.Run("returns error when TokenOutID does not match", func(t *testing.T) {
		ctx := context.Background()
		swap := createTestDvpSwap()
		swap.TokenOutID = "id1"
		differentID := "id2"

		svc := service.NewSwapAgreement(nil, nil, nil, nil)

		message, err := svc.Verify(
			ctx,
			swap,
			swap.DestChainID,
			swap.TokenInResourceID,
			swap.TokenInAmount,
			swap.TokenInID,
			swap.TokenInType,
			swap.TokenOutResourceID,
			swap.TokenOutAmount,
			differentID,
			swap.TokenOutType,
		)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "TokenOutID mismatch")
		assert.Contains(t, message, "Token out ID doesn't match")
		assert.Contains(t, message, "Expected: id2")
		assert.Contains(t, message, "Got: id1")
	})

	t.Run("returns error when TokenOutType does not match", func(t *testing.T) {
		ctx := context.Background()
		swap := createTestDvpSwap()
		swap.TokenOutType = types.DvpERC721
		differentType := types.DvpERC1155

		svc := service.NewSwapAgreement(nil, nil, nil, nil)

		message, err := svc.Verify(
			ctx,
			swap,
			swap.DestChainID,
			swap.TokenInResourceID,
			swap.TokenInAmount,
			swap.TokenInID,
			swap.TokenInType,
			swap.TokenOutResourceID,
			swap.TokenOutAmount,
			swap.TokenOutID,
			differentType,
		)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "TokenOutType mismatch")
		assert.Contains(t, message, "Token out type doesn't match")
		assert.Contains(t, message, fmt.Sprintf("Expected: %d", differentType))
		assert.Contains(t, message, fmt.Sprintf("Got: %d", swap.TokenOutType))
	})
}

func TestSwapAgreement_HandleSwapDisagreement(t *testing.T) {
	const (
		sharedID    = "test-shared-id"
		resourceID  = "resource-in-123"
		testMessage = "test disagreement message"
	)
	testAddress := createTestAddress("0x1234567890123456789012345678901234567890")

	endpointOK := func() *swapAgreementEndpointClientMock {
		return &swapAgreementEndpointClientMock{
			GetResourceAddressFunc: func(_ context.Context, resourceId string) (common.Address, error) {
				return testAddress, nil
			},
		}
	}

	t.Run("successfully handles disagreement for Enygma token type", func(t *testing.T) {
		ctx := context.Background()

		enygmaHandlerMock := &swapAgreementEnygmaHandlerClientMock{
			NotifySenderWithPNCommunicatorFunc: func(ctx context.Context, tokenAddress common.Address, sharedId string, status types.DvpCommunicatorStatus, message string) error {
				return nil
			},
		}
		endpointMock := endpointOK()

		svc := service.NewSwapAgreement(endpointMock, enygmaHandlerMock, nil, nil)

		err := svc.HandleSwapDisagreement(ctx, sharedID, resourceID, types.DvpEnygma, testMessage)

		require.NoError(t, err)

		endpointCalls := endpointMock.GetResourceAddressCalls()
		require.Len(t, endpointCalls, 1)
		assert.Equal(t, resourceID, endpointCalls[0].ResourceId)

		notifyCalls := enygmaHandlerMock.NotifySenderWithPNCommunicatorCalls()
		require.Len(t, notifyCalls, 1)
		assert.Equal(t, testAddress, notifyCalls[0].TokenAddress)
		assert.Equal(t, sharedID, notifyCalls[0].SharedId)
		assert.Equal(t, types.SwapError, notifyCalls[0].Status)
		assert.Equal(t, testMessage, notifyCalls[0].Message)
	})

	t.Run("successfully handles disagreement for ERC721 token type", func(t *testing.T) {
		ctx := context.Background()

		erc721HandlerMock := &swapAgreementDvpERC721HandlerClientMock{
			NotifySenderWithPNCommunicatorFunc: func(ctx context.Context, tokenAddress common.Address, sharedId string, status types.DvpCommunicatorStatus, message string) error {
				return nil
			},
		}
		endpointMock := endpointOK()

		svc := service.NewSwapAgreement(endpointMock, nil, erc721HandlerMock, nil)

		err := svc.HandleSwapDisagreement(ctx, sharedID, resourceID, types.DvpERC721, testMessage)

		require.NoError(t, err)

		notifyCalls := erc721HandlerMock.NotifySenderWithPNCommunicatorCalls()
		require.Len(t, notifyCalls, 1)
		assert.Equal(t, testAddress, notifyCalls[0].TokenAddress)
		assert.Equal(t, sharedID, notifyCalls[0].SharedId)
		assert.Equal(t, types.SwapError, notifyCalls[0].Status)
		assert.Equal(t, testMessage, notifyCalls[0].Message)
	})

	t.Run("successfully handles disagreement for ERC1155 token type", func(t *testing.T) {
		ctx := context.Background()

		erc1155HandlerMock := &swapAgreementDvpERC1155HandlerClientMock{
			NotifySenderWithPNCommunicatorFunc: func(ctx context.Context, tokenAddress common.Address, sharedId string, status types.DvpCommunicatorStatus, message string) error {
				return nil
			},
		}
		endpointMock := endpointOK()

		svc := service.NewSwapAgreement(endpointMock, nil, nil, erc1155HandlerMock)

		err := svc.HandleSwapDisagreement(ctx, sharedID, resourceID, types.DvpERC1155, testMessage)

		require.NoError(t, err)

		notifyCalls := erc1155HandlerMock.NotifySenderWithPNCommunicatorCalls()
		require.Len(t, notifyCalls, 1)
		assert.Equal(t, testAddress, notifyCalls[0].TokenAddress)
		assert.Equal(t, sharedID, notifyCalls[0].SharedId)
		assert.Equal(t, types.SwapError, notifyCalls[0].Status)
		assert.Equal(t, testMessage, notifyCalls[0].Message)
	})

	t.Run("returns error when GetResourceAddress fails", func(t *testing.T) {
		ctx := context.Background()

		endpointMock := &swapAgreementEndpointClientMock{
			GetResourceAddressFunc: func(_ context.Context, resourceId string) (common.Address, error) {
				return common.Address{}, errors.New("resource address lookup failed")
			},
		}

		svc := service.NewSwapAgreement(endpointMock, nil, nil, nil)

		err := svc.HandleSwapDisagreement(ctx, sharedID, resourceID, types.DvpEnygma, testMessage)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "resource address lookup failed")
		assert.Len(t, endpointMock.GetResourceAddressCalls(), 1)
	})

	t.Run("returns error when Enygma NotifySenderWithPNCommunicator fails", func(t *testing.T) {
		ctx := context.Background()

		enygmaHandlerMock := &swapAgreementEnygmaHandlerClientMock{
			NotifySenderWithPNCommunicatorFunc: func(ctx context.Context, tokenAddress common.Address, sharedId string, status types.DvpCommunicatorStatus, message string) error {
				return errors.New("notification failed")
			},
		}

		svc := service.NewSwapAgreement(endpointOK(), enygmaHandlerMock, nil, nil)

		err := svc.HandleSwapDisagreement(ctx, sharedID, resourceID, types.DvpEnygma, testMessage)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "notification failed")
		assert.Len(t, enygmaHandlerMock.NotifySenderWithPNCommunicatorCalls(), 1)
	})

	t.Run("returns error when ERC721 NotifySenderWithPNCommunicator fails", func(t *testing.T) {
		ctx := context.Background()

		erc721HandlerMock := &swapAgreementDvpERC721HandlerClientMock{
			NotifySenderWithPNCommunicatorFunc: func(ctx context.Context, tokenAddress common.Address, sharedId string, status types.DvpCommunicatorStatus, message string) error {
				return errors.New("erc721 notification failed")
			},
		}

		svc := service.NewSwapAgreement(endpointOK(), nil, erc721HandlerMock, nil)

		err := svc.HandleSwapDisagreement(ctx, sharedID, resourceID, types.DvpERC721, testMessage)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "erc721 notification failed")
		assert.Len(t, erc721HandlerMock.NotifySenderWithPNCommunicatorCalls(), 1)
	})

	t.Run("returns error when ERC1155 NotifySenderWithPNCommunicator fails", func(t *testing.T) {
		ctx := context.Background()

		erc1155HandlerMock := &swapAgreementDvpERC1155HandlerClientMock{
			NotifySenderWithPNCommunicatorFunc: func(ctx context.Context, tokenAddress common.Address, sharedId string, status types.DvpCommunicatorStatus, message string) error {
				return errors.New("erc1155 notification failed")
			},
		}

		svc := service.NewSwapAgreement(endpointOK(), nil, nil, erc1155HandlerMock)

		err := svc.HandleSwapDisagreement(ctx, sharedID, resourceID, types.DvpERC1155, testMessage)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "erc1155 notification failed")
		assert.Len(t, erc1155HandlerMock.NotifySenderWithPNCommunicatorCalls(), 1)
	})

	t.Run("returns error for unsupported token type", func(t *testing.T) {
		ctx := context.Background()

		endpointMock := endpointOK()
		svc := service.NewSwapAgreement(endpointMock, nil, nil, nil)

		err := svc.HandleSwapDisagreement(ctx, sharedID, resourceID, types.DvpERC20, testMessage)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "unsupported token type")
		assert.Len(t, endpointMock.GetResourceAddressCalls(), 1)
	})
}

// Helper functions

func createTestDvpSwap() *types.DvpSwap {
	return &types.DvpSwap{
		SharedID:               "test-shared-id-123",
		From:                   "0xSenderAddress",
		To:                     "0xReceiverAddress",
		SourceChainID:          big.NewInt(1),
		DestChainID:            big.NewInt(2),
		TokenInAmount:          big.NewInt(1000),
		TokenInResourceID:      "resource-in-123",
		TokenInType:            types.DvpEnygma,
		TokenInID:              "token-in-id-1",
		TokenOutAmount:         big.NewInt(500),
		TokenOutResourceID:     "resource-out-456",
		TokenOutType:           types.DvpERC721,
		TokenOutID:             "token-out-id-2",
		Status:                 types.DvpSwapInitiated,
		CreatedAt:              time.Now(),
	}
}

func createTestAddress(hex string) common.Address {
	return common.HexToAddress(hex)
}
