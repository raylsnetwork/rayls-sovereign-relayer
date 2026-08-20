package contractclient_test

import (
	"context"
	"errors"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/raylsnetwork/rayls-sovereign-relayer/contractclient"
	"github.com/stretchr/testify/require"
)

type stubResourceRegistryBackend struct {
	result []byte
	err    error
}

func (s *stubResourceRegistryBackend) CallContract(
	_ context.Context,
	_ ethereum.CallMsg,
	_ *big.Int,
) ([]byte, error) {
	return s.result, s.err
}

func TestResourceRegistryClient_GetResourceById(t *testing.T) {
	t.Run("wraps backend errors in ResourceRegistryClientError", func(t *testing.T) {
		wantResourceId := [32]byte{0x01, 0x02, 0x03, 0x04}
		wantError := errors.New("contract error")

		backend := &stubResourceRegistryBackend{
			err: wantError,
		}

		client := contractclient.NewResourceRegistryClient(common.Address{}, backend)

		_, _, _, gotErr := client.GetResourceById(wantResourceId)

		var wrappedErr *contractclient.ResourceRegistryClientError
		require.NotNil(t, gotErr, "should return error when contract call fails")
		require.ErrorAs(t, gotErr, &wrappedErr, "error should be wrapped in ResourceRegistryClientError")
		require.ErrorIs(t, gotErr, wantError, "underlying error should be preserved")
	})
}
