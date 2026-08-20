package service_test

import (
	"context"
	"errors"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/contracts/EndpointV1"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/enygma/service"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

type deployerEndpointClientMock struct {
	GetResourceAddressFunc  func(ctx context.Context, resourceId string) (common.Address, error)
	ReceivePayloadFunc      func(ctx context.Context, fromChainId *big.Int, from common.Address, to common.Address, data EndpointV1.RaylsMessage, messageId [32]byte) (common.Hash, error)
	getResourceAddressCalls int
	receivePayloadCalls     int
}

func (m *deployerEndpointClientMock) GetResourceAddress(ctx context.Context, resourceId string) (common.Address, error) {
	m.getResourceAddressCalls++
	return m.GetResourceAddressFunc(ctx, resourceId)
}

func (m *deployerEndpointClientMock) ReceivePayload(
	ctx context.Context,
	fromChainId *big.Int,
	from common.Address,
	to common.Address,
	data EndpointV1.RaylsMessage,
	messageId [32]byte,
) (common.Hash, error) {
	m.receivePayloadCalls++
	return m.ReceivePayloadFunc(ctx, fromChainId, from, to, data, messageId)
}

type deployerResourceRegistryClientMock struct {
	GetResourceByIdFunc  func(resourceId [32]byte) (uint8, []byte, []byte, error)
	getResourceByIdCalls int
}

func (m *deployerResourceRegistryClientMock) GetResourceById(resourceId [32]byte) (uint8, []byte, []byte, error) {
	m.getResourceByIdCalls++
	return m.GetResourceByIdFunc(resourceId)
}

type deployerTracerMock struct {
	StartFunc  func(ctx context.Context, spanName string, opts ...trace.SpanStartOption) (context.Context, trace.Span)
	startCalls int
	lastSpan   string
}

func (m *deployerTracerMock) Start(
	ctx context.Context,
	spanName string,
	opts ...trace.SpanStartOption,
) (context.Context, trace.Span) {
	m.startCalls++
	m.lastSpan = spanName
	if m.StartFunc != nil {
		return m.StartFunc(ctx, spanName, opts...)
	}
	// Default: return noop span from OpenTelemetry
	tracer := otel.GetTracerProvider().Tracer("test")
	return tracer.Start(ctx, spanName, opts...)
}

func testResourceId() [32]byte {
	return [32]byte{
		1,
		2,
		3,
		4,
		5,
		6,
		7,
		8,
		9,
		10,
		11,
		12,
		13,
		14,
		15,
		16,
		17,
		18,
		19,
		20,
		21,
		22,
		23,
		24,
		25,
		26,
		27,
		28,
		29,
		30,
		31,
		32,
	}
}

func testChainId() *big.Int {
	return big.NewInt(1)
}

func testDeployedAddress() common.Address {
	return common.HexToAddress("0x1234567890123456789012345678901234567890")
}

func testBytecode() []byte {
	return []byte{0x60, 0x80, 0x60, 0x40, 0x52}
}

func testInitParams() []byte {
	return []byte{0x00, 0x01}
}

func TestNewEnygmaDeployer(t *testing.T) {
	endpointMock := &deployerEndpointClientMock{}
	registryMock := &deployerResourceRegistryClientMock{}
	tracerMock := &deployerTracerMock{}

	deployer := service.NewEnygmaDeployer(endpointMock, registryMock, tracerMock)

	require.NotNil(t, deployer)
}

func TestEnygmaDeployer_Deploy(t *testing.T) {
	t.Run("successfully deploys contract", func(t *testing.T) {
		resourceId := testResourceId()
		chainId := testChainId()
		deployedAddress := testDeployedAddress()
		bytecode := testBytecode()
		initParams := testInitParams()

		endpointMock := &deployerEndpointClientMock{
			ReceivePayloadFunc: func(ctx context.Context, fromChainId *big.Int, from common.Address, to common.Address, data EndpointV1.RaylsMessage, messageId [32]byte) (common.Hash, error) {
				return common.Hash{}, nil
			},
			GetResourceAddressFunc: func(_ context.Context, resourceId string) (common.Address, error) {
				return deployedAddress, nil
			},
		}

		registryMock := &deployerResourceRegistryClientMock{
			GetResourceByIdFunc: func(resourceId [32]byte) (uint8, []byte, []byte, error) {
				return uint8(types.ENYGMA), bytecode, initParams, nil
			},
		}

		tracerMock := &deployerTracerMock{}

		deployer := service.NewEnygmaDeployer(endpointMock, registryMock, tracerMock)
		address, err := deployer.Deploy(context.Background(), resourceId, chainId)

		require.NoError(t, err)
		assert.Equal(t, deployedAddress, address)
		assert.Equal(t, 1, registryMock.getResourceByIdCalls)
		assert.Equal(t, 1, endpointMock.receivePayloadCalls)
		assert.Equal(t, 1, endpointMock.getResourceAddressCalls)
		assert.Equal(t, 1, tracerMock.startCalls)
	})

	t.Run("returns error when resource registry lookup fails", func(t *testing.T) {
		resourceId := testResourceId()
		chainId := testChainId()

		endpointMock := &deployerEndpointClientMock{}

		registryMock := &deployerResourceRegistryClientMock{
			GetResourceByIdFunc: func(resourceId [32]byte) (uint8, []byte, []byte, error) {
				return 0, nil, nil, errors.New("registry lookup failed")
			},
		}

		tracerMock := &deployerTracerMock{}

		deployer := service.NewEnygmaDeployer(endpointMock, registryMock, tracerMock)
		address, err := deployer.Deploy(context.Background(), resourceId, chainId)

		assert.Error(t, err)
		assert.Equal(t, common.Address{}, address)
		assert.Equal(t, 1, registryMock.getResourceByIdCalls)
		assert.Equal(t, 0, endpointMock.receivePayloadCalls)
	})

	t.Run("returns error when endpoint receive payload fails", func(t *testing.T) {
		resourceId := testResourceId()
		chainId := testChainId()
		bytecode := testBytecode()
		initParams := testInitParams()

		endpointMock := &deployerEndpointClientMock{
			ReceivePayloadFunc: func(ctx context.Context, fromChainId *big.Int, from common.Address, to common.Address, data EndpointV1.RaylsMessage, messageId [32]byte) (common.Hash, error) {
				return common.Hash{}, errors.New("endpoint receive failed")
			},
		}

		registryMock := &deployerResourceRegistryClientMock{
			GetResourceByIdFunc: func(resourceId [32]byte) (uint8, []byte, []byte, error) {
				return uint8(types.ENYGMA), bytecode, initParams, nil
			},
		}

		tracerMock := &deployerTracerMock{}

		deployer := service.NewEnygmaDeployer(endpointMock, registryMock, tracerMock)
		address, err := deployer.Deploy(context.Background(), resourceId, chainId)

		assert.Error(t, err)
		assert.Equal(t, common.Address{}, address)
		assert.Equal(t, 1, registryMock.getResourceByIdCalls)
		assert.Equal(t, 1, endpointMock.receivePayloadCalls)
		assert.Equal(t, 0, endpointMock.getResourceAddressCalls)
	})

	t.Run("returns error when resource address resolution fails", func(t *testing.T) {
		resourceId := testResourceId()
		chainId := testChainId()
		bytecode := testBytecode()
		initParams := testInitParams()

		endpointMock := &deployerEndpointClientMock{
			ReceivePayloadFunc: func(ctx context.Context, fromChainId *big.Int, from common.Address, to common.Address, data EndpointV1.RaylsMessage, messageId [32]byte) (common.Hash, error) {
				return common.Hash{}, nil
			},
			GetResourceAddressFunc: func(_ context.Context, resourceId string) (common.Address, error) {
				return common.Address{}, errors.New("address resolution failed")
			},
		}

		registryMock := &deployerResourceRegistryClientMock{
			GetResourceByIdFunc: func(resourceId [32]byte) (uint8, []byte, []byte, error) {
				return uint8(types.ENYGMA), bytecode, initParams, nil
			},
		}

		tracerMock := &deployerTracerMock{}

		deployer := service.NewEnygmaDeployer(endpointMock, registryMock, tracerMock)
		address, err := deployer.Deploy(context.Background(), resourceId, chainId)

		assert.Error(t, err)
		assert.Equal(t, common.Address{}, address)
		assert.Equal(t, 1, registryMock.getResourceByIdCalls)
		assert.Equal(t, 1, endpointMock.receivePayloadCalls)
		assert.Equal(t, 1, endpointMock.getResourceAddressCalls)
	})

	t.Run("tracer span is started during deployment", func(t *testing.T) {
		resourceId := testResourceId()
		chainId := testChainId()
		bytecode := testBytecode()
		initParams := testInitParams()
		deployedAddress := testDeployedAddress()

		endpointMock := &deployerEndpointClientMock{
			ReceivePayloadFunc: func(ctx context.Context, fromChainId *big.Int, from common.Address, to common.Address, data EndpointV1.RaylsMessage, messageId [32]byte) (common.Hash, error) {
				return common.Hash{}, nil
			},
			GetResourceAddressFunc: func(_ context.Context, resourceId string) (common.Address, error) {
				return deployedAddress, nil
			},
		}

		registryMock := &deployerResourceRegistryClientMock{
			GetResourceByIdFunc: func(resourceId [32]byte) (uint8, []byte, []byte, error) {
				return uint8(types.ENYGMA), bytecode, initParams, nil
			},
		}

		tracerMock := &deployerTracerMock{}

		deployer := service.NewEnygmaDeployer(endpointMock, registryMock, tracerMock)
		_, err := deployer.Deploy(context.Background(), resourceId, chainId)

		require.NoError(t, err)
		assert.Equal(t, 1, tracerMock.startCalls)
		assert.NotEmpty(t, tracerMock.lastSpan)
	})
}
