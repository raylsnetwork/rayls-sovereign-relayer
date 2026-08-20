// Decommissioning Teleport (vanilla, atomic).

package txgen_test

import (
	"context"
	"errors"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/private-relayer/dest/txgen"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	testChainID  = big.NewInt(1337)
	testGasPrice = big.NewInt(1_000_000_000)
	testEndpoint = common.HexToAddress("0x3333333333333333333333333333333333333333")
)

func validEthereumClientMock() *EthereumClientMock {
	return &EthereumClientMock{
		ChainIDFunc: func(ctx context.Context) (*big.Int, error) {
			return testChainID, nil
		},
		SuggestGasPriceFunc: func(ctx context.Context) (*big.Int, error) {
			return testGasPrice, nil
		},
		EstimateGasFunc: func(ctx context.Context, msg ethereum.CallMsg) (uint64, error) {
			return 21000, nil
		},
	}
}

func validCalldataSignature() types.CalldataSignature {
	return types.CalldataSignature{
		SharedId:   "test-shared-id-123",
		Signature:  []byte{0xAA, 0xBB, 0xCC},
		ResourceId: [32]byte{0x01, 0x02, 0x03},
	}
}

func newSignatureGenerator(t *testing.T) *txgen.SignatureGenerator {
	t.Helper()
	gen, err := txgen.NewSignatureGenerator(validEthereumClientMock(), testEndpoint)
	require.NoError(t, err)
	return gen
}

func TestNewSignatureGenerator(t *testing.T) {
	t.Run("creates generator with valid client", func(t *testing.T) {
		client := validEthereumClientMock()

		gen, err := txgen.NewSignatureGenerator(client, testEndpoint)
		require.NoError(t, err)
		assert.NotNil(t, gen)
		assert.Len(t, client.ChainIDCalls(), 1)
	})

	t.Run("returns error when ChainID fails", func(t *testing.T) {
		client := &EthereumClientMock{
			ChainIDFunc: func(ctx context.Context) (*big.Int, error) {
				return nil, errors.New("chain ID unavailable")
			},
		}

		gen, err := txgen.NewSignatureGenerator(client, testEndpoint)
		require.Error(t, err)
		assert.Nil(t, gen)
		assert.Contains(t, err.Error(), "error getting chain ID")
	})
}

func TestSignatureGenerator_Generate(t *testing.T) {
	t.Run("returns calldata bytes", func(t *testing.T) {
		gen := newSignatureGenerator(t)

		data, err := gen.Generate(validCalldataSignature())
		require.NoError(t, err)
		assert.NotEmpty(t, data)
	})

	t.Run("calldata starts with receivePayload selector", func(t *testing.T) {
		gen := newSignatureGenerator(t)

		data, err := gen.Generate(validCalldataSignature())
		require.NoError(t, err)

		selector := receivePayloadSelector(t)
		assert.Equal(t, selector, data[:4])
	})

	t.Run("builds RaylsMessage with msgId derived from SharedId", func(t *testing.T) {
		gen := newSignatureGenerator(t)

		sig := validCalldataSignature()
		data, err := gen.Generate(sig)
		require.NoError(t, err)

		// Decode args and check the msgId field — it must be Keccak256(SharedId).
		endpointABI := validEndpointABI(t)
		method := endpointABI.Methods["receivePayload"]
		args, err := method.Inputs.Unpack(data[4:])
		require.NoError(t, err)
		require.Len(t, args, 5)

		msgId, ok := args[4].([32]byte)
		require.True(t, ok)
		expectedMsgId := crypto.Keccak256Hash([]byte(sig.SharedId))
		assert.Equal(t, expectedMsgId, common.Hash(msgId))
	})

}
