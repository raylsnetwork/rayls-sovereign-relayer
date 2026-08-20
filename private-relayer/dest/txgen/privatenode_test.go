package txgen_test

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/raylsnetwork/rayls-sovereign-relayer/contracts/EndpointV1"
	"github.com/raylsnetwork/rayls-sovereign-relayer/private-relayer/dest/txgen"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func validEndpointABI(t *testing.T) *abi.ABI {
	t.Helper()
	endpointABI, err := EndpointV1.EndpointV1MetaData.ParseABI()
	require.NoError(t, err, "failed to parse EndpointV1 ABI")
	return endpointABI
}

func validRaylsMessage() EndpointV1.RaylsMessage {
	return EndpointV1.RaylsMessage{
		Payload: []byte{0xFF, 0x00, 0xFF},
		MessageMetadata: EndpointV1.RaylsMessageMetadata{
			Valid:        true,
			Nonce:        big.NewInt(1),
			ResourceId:   [32]byte{0x01},
			IgnoresNonce: false,
			NewResourceMetadata: EndpointV1.NewResourceMetadata{
				Valid:              false,
				ResourceDeployType: 0,
				Bytecode:           []byte{},
				FactoryTemplate:    0,
				InitializerParams:  []byte{},
			},
			LockData:                  []byte{},
			RevertPayloadDataSender:   []byte{},
			RevertPayloadDataReceiver: []byte{},
			TransferMetadata: EndpointV1.BridgedTransferMetadata{
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

func TestNewPrivateNodeGenerator(t *testing.T) {
	t.Run("creates generator with valid inputs", func(t *testing.T) {
		client := &EthereumClientMock{}

		gen, err := txgen.NewPrivateNodeGenerator(client)
		require.NoError(t, err)
		assert.NotNil(t, gen)
	})
}

func TestPrivateNodeGenerator_Generate(t *testing.T) {
	t.Run("packs valid receivePayload calldata", func(t *testing.T) {
		client := &EthereumClientMock{}
		gen, err := txgen.NewPrivateNodeGenerator(client)
		require.NoError(t, err)

		data, err := gen.Generate(
			big.NewInt(1337),
			common.HexToAddress("0x1111111111111111111111111111111111111111"),
			common.HexToAddress("0x2222222222222222222222222222222222222222"),
			validRaylsMessage(),
			common.HexToHash("0xabcdef"),
		)
		require.NoError(t, err)
		assert.NotEmpty(t, data)
		assert.Greater(t, len(data), 4, "calldata should contain selector + encoded args")
	})

	t.Run("returns data with correct function selector", func(t *testing.T) {
		client := &EthereumClientMock{}
		gen, err := txgen.NewPrivateNodeGenerator(client)
		require.NoError(t, err)

		data, err := gen.Generate(
			big.NewInt(1),
			common.Address{},
			common.Address{},
			validRaylsMessage(),
			common.Hash{},
		)
		require.NoError(t, err)
		assert.Equal(t, receivePayloadSelector(t), data[:4])
	})

	t.Run("produces different output for different message IDs", func(t *testing.T) {
		client := &EthereumClientMock{}
		gen, err := txgen.NewPrivateNodeGenerator(client)
		require.NoError(t, err)

		msg := validRaylsMessage()
		data1, err := gen.Generate(big.NewInt(1), common.Address{}, common.Address{}, msg, common.HexToHash("0x01"))
		require.NoError(t, err)

		data2, err := gen.Generate(big.NewInt(1), common.Address{}, common.Address{}, msg, common.HexToHash("0x02"))
		require.NoError(t, err)

		assert.NotEqual(t, data1, data2)
	})

	t.Run("produces different output for different chain IDs", func(t *testing.T) {
		client := &EthereumClientMock{}
		gen, err := txgen.NewPrivateNodeGenerator(client)
		require.NoError(t, err)

		msg := validRaylsMessage()
		id := common.HexToHash("0xAA")

		data1, err := gen.Generate(big.NewInt(1), common.Address{}, common.Address{}, msg, id)
		require.NoError(t, err)

		data2, err := gen.Generate(big.NewInt(9999), common.Address{}, common.Address{}, msg, id)
		require.NoError(t, err)

		assert.NotEqual(t, data1, data2)
	})
}
