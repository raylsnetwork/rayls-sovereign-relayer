package contractclient_test

import (
	"context"
	"errors"
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/contractclient"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/contracts/EndpointV1"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type StubDeployerEndpoint struct {
	spyCtx        context.Context
	spyFromChainId *big.Int
	spyFrom        common.Address
	spyTo          common.Address
	spyData        EndpointV1.RaylsMessage
	spyMessageId   [32]byte

	txHash common.Hash
	err    error
}

func (c *StubDeployerEndpoint) ReceivePayload(
	ctx context.Context,
	fromChainId *big.Int,
	from, to common.Address,
	data EndpointV1.RaylsMessage,
	messageId [32]byte,
) (common.Hash, error) {
	c.spyCtx = ctx
	c.spyFromChainId = fromChainId
	c.spyFrom = from
	c.spyTo = to
	c.spyData = data
	c.spyMessageId = messageId
	return c.txHash, c.err
}

type StubResourceRegistryClient struct {
	standard          uint8
	bytecode          []byte
	initializerParams []byte
	err               error

	spyResourceId [32]byte
}

func (r *StubResourceRegistryClient) GetResourceById(resourceId [32]byte) (uint8, []byte, []byte, error) {
	r.spyResourceId = resourceId
	return r.standard, r.bytecode, r.initializerParams, r.err
}

func TestDeployerClient_DeployResourceAndExecute(t *testing.T) {
	t.Run("successfully deploys resource and executes payload", func(t *testing.T) {
		// Setup test data
		wantResourceId := [32]byte{0x01, 0x02, 0x03}
		wantBytecode := []byte{0x60, 0x80, 0x60, 0x40}
		wantInitParams := []byte{0x00, 0x00, 0x00, 0x01}

		message := &types.DispatchedMessageToPrivateHub{
			MessageId:   [32]byte{0xAA, 0xBB, 0xCC, 0xDD},
			From:        common.HexToAddress("0x1111111111111111111111111111111111111111"),
			To:          common.HexToAddress("0x2222222222222222222222222222222222222222"),
			FromChainId: big.NewInt(100),
			ToChainId:   big.NewInt(200),
			Data: EndpointV1.RaylsMessage{
				MessageMetadata: EndpointV1.RaylsMessageMetadata{
					Valid: true,
				},
				Payload: []byte("test payload"),
			},
		}

		txHash := common.HexToHash("0xabcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890")

		// Setup stubs
		resourceRegistry := &StubResourceRegistryClient{
			bytecode:          wantBytecode,
			initializerParams: wantInitParams,
		}

		endpoint := &StubDeployerEndpoint{
			txHash: txHash,
		}

		client := contractclient.NewDeployerClient(
			endpoint,
			resourceRegistry,
			30*time.Second,
		)

		// Execute
		gotTxHash, err := client.DeployResourceAndExecute(
			context.Background(),
			wantResourceId,
			message,
		)

		// Assert
		require.Nil(t, err)
		assert.Equal(t, txHash, gotTxHash)

		// Verify resource registry was called with correct resource ID
		assert.Equal(t, wantResourceId, resourceRegistry.spyResourceId)

		// Verify endpoint was called with correct parameters
		assert.Equal(t, message.FromChainId, endpoint.spyFromChainId)
		assert.Equal(t, message.From, endpoint.spyFrom)
		assert.Equal(t, message.To, endpoint.spyTo)
		assert.Equal(t, message.MessageId, endpoint.spyMessageId)

		// CUSTOM standard (0) → legacy BYTECODE deploy with the issuer-supplied bytecode.
		assert.True(t, endpoint.spyData.MessageMetadata.NewResourceMetadata.Valid)
		assert.Equal(t, uint8(0), endpoint.spyData.MessageMetadata.NewResourceMetadata.ResourceDeployType)
		assert.Equal(t, wantBytecode, endpoint.spyData.MessageMetadata.NewResourceMetadata.Bytecode)
		assert.Equal(t, wantInitParams, endpoint.spyData.MessageMetadata.NewResourceMetadata.InitializerParams)
	})

	t.Run("routes a non-Custom standard through FACTORY mode", func(t *testing.T) {
		message := &types.DispatchedMessageToPrivateHub{
			FromChainId: big.NewInt(100),
			Data:        EndpointV1.RaylsMessage{Payload: []byte("p")},
		}

		// A FACTORY-mode standard token: the registry reports the standard but holds no
		// raw bytecode. The deployer must select FACTORY and pass the template, not empty
		// bytecode (which would revert FactoryV1__EmptyBytecode on the receiver factory).
		resourceRegistry := &StubResourceRegistryClient{
			standard:          uint8(types.ERC20),
			bytecode:          nil,
			initializerParams: []byte{0x00, 0x00, 0x00, 0x01},
		}
		endpoint := &StubDeployerEndpoint{txHash: common.HexToHash("0x01")}

		client := contractclient.NewDeployerClient(endpoint, resourceRegistry, 30*time.Second)

		_, err := client.DeployResourceAndExecute(context.Background(), [32]byte{0x09}, message)
		require.Nil(t, err)

		meta := endpoint.spyData.MessageMetadata.NewResourceMetadata
		assert.True(t, meta.Valid)
		assert.Equal(t, uint8(types.ResourceDeployTypeFactory), meta.ResourceDeployType)
		assert.Equal(t, uint8(types.ERC20), meta.FactoryTemplate)
		assert.Empty(t, meta.Bytecode)
		assert.Equal(t, []byte{0x00, 0x00, 0x00, 0x01}, meta.InitializerParams)
	})

	t.Run("wraps resource registry errors in DeployerClientError", func(t *testing.T) {
		wantError := errors.New("failed to get resource")

		message := &types.DispatchedMessageToPrivateHub{
			FromChainId: big.NewInt(100),
			Data:        EndpointV1.RaylsMessage{},
		}

		resourceRegistry := &StubResourceRegistryClient{
			err: wantError,
		}
		endpoint := &StubDeployerEndpoint{}

		client := contractclient.NewDeployerClient(
			endpoint,
			resourceRegistry,
			30*time.Second,
		)

		_, err := client.DeployResourceAndExecute(
			context.Background(),
			[32]byte{},
			message,
		)

		var gotErr *contractclient.DeployerClientError
		require.ErrorAs(t, err, &gotErr, "error should be wrapped in DeployerClientError")
		require.ErrorIs(t, err, wantError, "underlying error should be preserved")
	})

	t.Run("wraps contract ReceivePayload errors in DeployerClientError", func(t *testing.T) {
		wantError := errors.New("contract error")

		message := &types.DispatchedMessageToPrivateHub{
			FromChainId: big.NewInt(100),
			Data:        EndpointV1.RaylsMessage{},
		}

		resourceRegistry := &StubResourceRegistryClient{
			bytecode:          []byte{0x60, 0x80},
			initializerParams: []byte{0x00},
		}
		endpoint := &StubDeployerEndpoint{
			err: wantError,
		}

		client := contractclient.NewDeployerClient(
			endpoint,
			resourceRegistry,
			30*time.Second,
		)

		_, err := client.DeployResourceAndExecute(
			context.Background(),
			[32]byte{},
			message,
		)

		var gotErr *contractclient.DeployerClientError
		require.ErrorAs(t, err, &gotErr, "error should be wrapped in DeployerClientError")
		require.ErrorIs(t, err, wantError, "underlying error should be preserved")
	})
}
