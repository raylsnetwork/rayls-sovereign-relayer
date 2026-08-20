// Decommissioning Teleport (vanilla, atomic).

package txgen_test

import (
	"context"
	"fmt"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/raylsnetwork/rayls-sovereign-relayer/contracts/RNEndpointV1"
	"github.com/raylsnetwork/rayls-sovereign-relayer/contracts/RNMessageDispatcherV1"
	"github.com/raylsnetwork/rayls-sovereign-relayer/public-relayer/txgen"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func validEndpointABI(t *testing.T) *abi.ABI {
	t.Helper()
	endpointABI, err := RNEndpointV1.RNEndpointV1MetaData.ParseABI()
	require.NoError(t, err, "failed to parse RNEndpointV1 ABI")
	return endpointABI
}

func validRaylsNodeMessage() RNMessageDispatcherV1.RaylsNodeMessage {
	return RNMessageDispatcherV1.RaylsNodeMessage{
		Payload: []byte{0xFF, 0x00, 0xFF},
		MessageMetadata: RNMessageDispatcherV1.RaylsNodeMessageMetadata{
			Nonce: big.NewInt(1),
			NewResourceMetadata: RNMessageDispatcherV1.RaylsNodeNewResourceMetadata{
				ResourceDeployType: 0,
				Bytecode:           []byte{},
				FactoryTemplate:    0,
				InitializerParams:  []byte{},
			},
			RevertPayloadData: []byte{},
			TransferMetadata: RNMessageDispatcherV1.RaylsNodeBridgedTransferMetadata{
				AssetType: 0,
				Id:        big.NewInt(0),
				Amount:    big.NewInt(0),
			},
		},
	}
}

func receivePayloadSelector(t *testing.T) []byte {
	t.Helper()
	endpointABI := validEndpointABI(t)
	method, ok := endpointABI.Methods["receivePayload"]
	require.True(t, ok, "receivePayload method not found in ABI")
	return method.ID
}

func mockClient(chainID *big.Int) *ERC1564EthereumClientMock {
	return &ERC1564EthereumClientMock{
		ChainIDFunc: func(context.Context) (*big.Int, error) {
			return chainID, nil
		},
	}
}

func TestNewEIP5164Generators(t *testing.T) {
	t.Run("creates generator with valid inputs", func(t *testing.T) {
		client := mockClient(big.NewInt(1337))

		gen, err := txgen.NewEIP5164Generators(client)
		require.NoError(t, err)
		assert.NotNil(t, gen)
		assert.Equal(t, 1, len(client.ChainIDCalls()))
	})

	t.Run("returns error when ChainID call fails", func(t *testing.T) {
		client := &ERC1564EthereumClientMock{
			ChainIDFunc: func(context.Context) (*big.Int, error) {
				return nil, fmt.Errorf("rpc connection refused")
			},
		}

		gen, err := txgen.NewEIP5164Generators(client)
		require.Error(t, err)
		assert.Nil(t, gen)
		assert.Contains(t, err.Error(), "error getting chain ID")
	})
}

func TestERC1564Generator_Generate(t *testing.T) {
	t.Run("packs valid receivePayload calldata", func(t *testing.T) {
		client := mockClient(big.NewInt(1337))
		gen, err := txgen.NewEIP5164Generators(client)
		require.NoError(t, err)

		data, err := gen.Generate(
			common.HexToAddress("0x1111111111111111111111111111111111111111"),
			common.HexToAddress("0x2222222222222222222222222222222222222222"),
			validRaylsNodeMessage(),
			common.HexToHash("0xabcdef"),
		)
		require.NoError(t, err)
		assert.NotEmpty(t, data)
		assert.Greater(t, len(data), 4, "calldata should contain selector + encoded args")
	})

	t.Run("returns data with correct function selector", func(t *testing.T) {
		client := mockClient(big.NewInt(1))
		gen, err := txgen.NewEIP5164Generators(client)
		require.NoError(t, err)

		data, err := gen.Generate(
			common.Address{},
			common.Address{},
			validRaylsNodeMessage(),
			common.Hash{},
		)
		require.NoError(t, err)
		assert.Equal(t, receivePayloadSelector(t), data[:4])
	})

	t.Run("produces different output for different message IDs", func(t *testing.T) {
		client := mockClient(big.NewInt(1))
		gen, err := txgen.NewEIP5164Generators(client)
		require.NoError(t, err)

		msg := validRaylsNodeMessage()
		data1, err := gen.Generate(common.Address{}, common.Address{}, msg, common.HexToHash("0x01"))
		require.NoError(t, err)

		data2, err := gen.Generate(common.Address{}, common.Address{}, msg, common.HexToHash("0x02"))
		require.NoError(t, err)

		assert.NotEqual(t, data1, data2)
	})

	t.Run("produces different output for different from addresses", func(t *testing.T) {
		client := mockClient(big.NewInt(1))
		gen, err := txgen.NewEIP5164Generators(client)
		require.NoError(t, err)

		msg := validRaylsNodeMessage()
		id := common.HexToHash("0xAA")

		data1, err := gen.Generate(common.HexToAddress("0x01"), common.Address{}, msg, id)
		require.NoError(t, err)

		data2, err := gen.Generate(common.HexToAddress("0x02"), common.Address{}, msg, id)
		require.NoError(t, err)

		assert.NotEqual(t, data1, data2)
	})
}
