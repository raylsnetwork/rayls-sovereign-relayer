package contractclient_test

import (
	"context"
	"errors"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/contractclient"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/contracts/EndpointV1"
	"github.com/stretchr/testify/require"
)

// Tests for GetResourceAddress
func TestEndpointClient_GetResourceAddress(t *testing.T) {
	t.Run("wraps executor call errors in EndpointClientError", func(t *testing.T) {
		// Valid 64-character hex string (32 bytes) without "0x" prefix
		wantResourceId := "1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef"
		wantError := errors.New("call error")

		address := common.HexToAddress("0x1")
		executor := &stubExecutor{
			callErr: wantError,
		}

		client := contractclient.NewEndpointClient(address, executor)

		_, gotErr := client.GetResourceAddress(context.Background(), wantResourceId)

		require.NotNil(t, gotErr, "should return error when call fails")
		require.ErrorIs(t, gotErr, wantError, "underlying error should be preserved")
	})

	t.Run("returns error for invalid resource ID format", func(t *testing.T) {
		invalidResourceId := "invalid-hex-string"

		address := common.HexToAddress("0x1")
		executor := &stubExecutor{}

		client := contractclient.NewEndpointClient(address, executor)

		_, gotErr := client.GetResourceAddress(context.Background(), invalidResourceId)

		require.NotNil(t, gotErr, "should return error for invalid resource ID")
	})
}

// Tests for ReceivePayload
func TestEndpointClient_ReceivePayload(t *testing.T) {
	t.Run("wraps executor errors in EndpointClientError", func(t *testing.T) {
		wantError := errors.New("execute failed")

		message := EndpointV1.RaylsMessage{
			MessageMetadata: EndpointV1.RaylsMessageMetadata{
				Nonce: big.NewInt(0),
				TransferMetadata: EndpointV1.BridgedTransferMetadata{
					Id:     big.NewInt(0),
					Amount: big.NewInt(0),
				},
			},
			Payload: []byte("test"),
		}

		address := common.HexToAddress("0x1")
		executor := &stubExecutor{
			executeErr: wantError,
		}

		client := contractclient.NewEndpointClient(address, executor)

		_, err := client.ReceivePayload(
			context.Background(),
			big.NewInt(100),
			common.Address{},
			common.Address{},
			message,
			[32]byte{},
		)

		var gotErr *contractclient.EndpointClientError
		require.ErrorAs(t, err, &gotErr, "error should be wrapped in EndpointClientError")
		require.ErrorIs(t, err, wantError, "underlying error should be preserved")
	})
}
