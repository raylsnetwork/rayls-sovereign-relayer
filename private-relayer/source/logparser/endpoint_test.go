package logparser_test

import (
	"context"
	"errors"
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/contracts/EndpointV1"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/msgqueue"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/private-relayer/shared/service"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/private-relayer/source/logparser"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/private-relayer/source/logrouter"
	srcservice "github.com/raylsnetwork/rayls-privacy-relayer-api/private-relayer/source/service"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/testtools"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/testtools/fake"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/testtools/spy"
	"github.com/stretchr/testify/assert"
)

// newMessageDispatchedLog creates a valid ABI-encoded MessageDispatched event log.
func newMessageDispatchedLog(
	messageID [32]byte,
	from common.Address,
	toChainID *big.Int,
	to common.Address,
) ethTypes.Log {
	abi, _ := EndpointV1.EndpointV1MetaData.ParseABI()

	data := EndpointV1.RaylsMessage{
		MessageMetadata: EndpointV1.RaylsMessageMetadata{
			Nonce: big.NewInt(0),
			NewResourceMetadata: EndpointV1.NewResourceMetadata{
				Bytecode:          []byte{},
				InitializerParams: []byte{},
			},
			LockData:                  []byte{},
			RevertPayloadDataSender:   []byte{},
			RevertPayloadDataReceiver: []byte{},
			TransferMetadata: EndpointV1.BridgedTransferMetadata{
				Id:     big.NewInt(0),
				Amount: big.NewInt(0),
			},
		},
		Payload: []byte{},
	}

	packedData, _ := abi.Events["MessageDispatched"].Inputs.NonIndexed().Pack(to, data)

	return ethTypes.Log{
		Address:     common.HexToAddress("0x1234567890123456789012345678901234567890"),
		BlockNumber: 100,
		TxHash:      common.HexToHash("0xabcdef"),
		TxIndex:     1,
		BlockHash:   common.HexToHash("0x123456"),
		Index:       0,
		Topics: []common.Hash{
			crypto.Keccak256Hash(
				[]byte(
					"MessageDispatched(bytes32,address,uint256,address,((bool,uint256,(bool,uint8,bytes,uint8,bytes),bytes32,bytes,bytes,bytes,(uint8,uint256,address,address,address,uint256),bool),bytes))",
				),
			),
			common.BytesToHash(messageID[:]),
			common.BytesToHash(from.Bytes()),
			common.BytesToHash(toChainID.Bytes()),
		},
		Data: packedData,
	}
}

// newUpdateViewKeysRequestLog creates a valid ABI-encoded UpdateViewKeysRequest event log.
func newUpdateViewKeysRequestLog(blockNumber *big.Int) ethTypes.Log {
	abi, _ := EndpointV1.EndpointV1MetaData.ParseABI()
	packedData, _ := abi.Events["UpdateRaylsViewKeysRequest"].Inputs.NonIndexed().Pack(blockNumber)

	return ethTypes.Log{
		Address:     common.HexToAddress("0x1234567890123456789012345678901234567890"),
		BlockNumber: 100,
		TxHash:      common.HexToHash("0xabcdef"),
		TxIndex:     1,
		BlockHash:   common.HexToHash("0x123456"),
		Index:       0,
		Topics: []common.Hash{
			crypto.Keccak256Hash([]byte("UpdateRaylsViewKeysRequest(uint256)")),
		},
		Data: packedData,
	}
}

func TestEndpointLogParser(t *testing.T) {
	testtools.SilenceLogger()

	privateHubChainID := big.NewInt(999)
	crossChainID := big.NewInt(42)

	t.Run("routes message to private hub MQ when destination is private hub", func(t *testing.T) {
		wantMessageID := [32]byte{0x01}
		wantFrom := common.HexToAddress("0x1111111111111111111111111111111111111111")
		wantTo := common.HexToAddress("0x2222222222222222222222222222222222222222")

		log := newMessageDispatchedLog(wantMessageID, wantFrom, privateHubChainID, wantTo)

		ackSpy := spy.NewAck()
		blockMQ := &EndpointMQMock{
			NextFunc: fake.NextMQ(msgqueue.Message[logrouter.Block]{
				V: logrouter.Block{
					Number: 100,
					Logs:   []ethTypes.Log{log},
				},
				Ack: ackSpy.Fn(),
			}),
		}

		calledPrivateHub := false
		privateHubMQ := &PrivateHubMQMock{
			PushFunc: func(ctx context.Context, msg service.PrivateHubMessage) error {
				calledPrivateHub = true
				assert.Equal(t, wantTo, msg.To)
				return nil
			},
		}

		crossChainMQ := &CrossChainMQMock{
			PushFunc: func(ctx context.Context, msg srcservice.CrossChainMessage) error {
				assert.Fail(t, "should not push to cross chain MQ for private hub destination")
				return nil
			},
		}

		keysService := &KeysServiceMock{}

		parser := logparser.NewEndpointLogParser(
			privateHubChainID, blockMQ, crossChainMQ, privateHubMQ, keysService,
		)

		ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond)
		defer cancel()

		_ = parser.Fetch(ctx)

		assert.True(t, calledPrivateHub, "should push to private hub MQ")
		ackSpy.AssertCalled(t)
	})

	t.Run("routes message to cross chain MQ when destination is not private hub", func(t *testing.T) {
		wantMessageID := [32]byte{0x02}
		wantFrom := common.HexToAddress("0x1111111111111111111111111111111111111111")
		wantTo := common.HexToAddress("0x2222222222222222222222222222222222222222")

		log := newMessageDispatchedLog(wantMessageID, wantFrom, crossChainID, wantTo)

		ackSpy := spy.NewAck()
		blockMQ := &EndpointMQMock{
			NextFunc: fake.NextMQ(msgqueue.Message[logrouter.Block]{
				V: logrouter.Block{
					Number: 100,
					Logs:   []ethTypes.Log{log},
				},
				Ack: ackSpy.Fn(),
			}),
		}

		privateHubMQ := &PrivateHubMQMock{
			PushFunc: func(ctx context.Context, msg service.PrivateHubMessage) error {
				assert.Fail(t, "should not push to private hub MQ for cross chain destination")
				return nil
			},
		}

		calledCrossChain := false
		crossChainMQ := &CrossChainMQMock{
			PushFunc: func(ctx context.Context, msg srcservice.CrossChainMessage) error {
				calledCrossChain = true
				assert.Equal(t, crossChainID, msg.ToChainID)
				return nil
			},
		}

		keysService := &KeysServiceMock{}

		parser := logparser.NewEndpointLogParser(
			privateHubChainID, blockMQ, crossChainMQ, privateHubMQ, keysService,
		)

		ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond)
		defer cancel()

		_ = parser.Fetch(ctx)

		assert.True(t, calledCrossChain, "should push to cross chain MQ")
		ackSpy.AssertCalled(t)
	})

	t.Run("does not ack block when private hub push fails", func(t *testing.T) {
		log := newMessageDispatchedLog([32]byte{0x03}, common.Address{}, privateHubChainID, common.Address{})

		ackSpy := spy.NewAck()
		blockMQ := &EndpointMQMock{
			NextFunc: fake.NextMQ(msgqueue.Message[logrouter.Block]{
				V: logrouter.Block{
					Number: 100,
					Logs:   []ethTypes.Log{log},
				},
				Ack: ackSpy.Fn(),
			}),
		}

		privateHubMQ := &PrivateHubMQMock{
			PushFunc: func(ctx context.Context, msg service.PrivateHubMessage) error {
				return errors.New("private hub MQ unavailable")
			},
		}
		crossChainMQ := &CrossChainMQMock{}
		keysService := &KeysServiceMock{}

		parser := logparser.NewEndpointLogParser(
			privateHubChainID, blockMQ, crossChainMQ, privateHubMQ, keysService,
		)

		ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond)
		defer cancel()

		_ = parser.Fetch(ctx)

		ackSpy.AssertNotCalled(t, "should not ack block when private hub push fails, allowing MQ redelivery")
	})

	t.Run("does not ack block when cross chain push fails", func(t *testing.T) {
		log := newMessageDispatchedLog([32]byte{0x04}, common.Address{}, crossChainID, common.Address{})

		ackSpy := spy.NewAck()
		blockMQ := &EndpointMQMock{
			NextFunc: fake.NextMQ(msgqueue.Message[logrouter.Block]{
				V: logrouter.Block{
					Number: 100,
					Logs:   []ethTypes.Log{log},
				},
				Ack: ackSpy.Fn(),
			}),
		}

		privateHubMQ := &PrivateHubMQMock{}
		crossChainMQ := &CrossChainMQMock{
			PushFunc: func(ctx context.Context, msg srcservice.CrossChainMessage) error {
				return errors.New("cross chain MQ unavailable")
			},
		}
		keysService := &KeysServiceMock{}

		parser := logparser.NewEndpointLogParser(
			privateHubChainID, blockMQ, crossChainMQ, privateHubMQ, keysService,
		)

		ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond)
		defer cancel()

		_ = parser.Fetch(ctx)

		ackSpy.AssertNotCalled(t, "should not ack block when cross chain push fails, allowing MQ redelivery")
	})

	t.Run("does not ack block when UpdateViewKeys fails", func(t *testing.T) {
		log := newUpdateViewKeysRequestLog(big.NewInt(42))

		ackSpy := spy.NewAck()
		blockMQ := &EndpointMQMock{
			NextFunc: fake.NextMQ(msgqueue.Message[logrouter.Block]{
				V: logrouter.Block{
					Number: 100,
					Logs:   []ethTypes.Log{log},
				},
				Ack: ackSpy.Fn(),
			}),
		}

		privateHubMQ := &PrivateHubMQMock{}
		crossChainMQ := &CrossChainMQMock{}
		keysService := &KeysServiceMock{
			UpdateRaylsViewKeysFunc: func(ctx context.Context, blockNumber *big.Int) error {
				return errors.New("DH keys service unavailable")
			},
		}

		parser := logparser.NewEndpointLogParser(
			privateHubChainID, blockMQ, crossChainMQ, privateHubMQ, keysService,
		)

		ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond)
		defer cancel()

		_ = parser.Fetch(ctx)

		ackSpy.AssertNotCalled(t, "should not ack block when UpdateViewKeys fails, allowing MQ redelivery")
	})

	t.Run("supports graceful shutdown", func(t *testing.T) {
		blockMQ := &EndpointMQMock{
			NextFunc: func(ctx context.Context) (msgqueue.Message[logrouter.Block], error) {
				<-ctx.Done()
				return msgqueue.Message[logrouter.Block]{}, ctx.Err()
			},
		}
		privateHubMQ := &PrivateHubMQMock{}
		crossChainMQ := &CrossChainMQMock{}
		keysService := &KeysServiceMock{}

		parser := logparser.NewEndpointLogParser(
			privateHubChainID, blockMQ, crossChainMQ, privateHubMQ, keysService,
		)

		hasGracefulShutdown := testtools.ShutdownFixture(t, parser.Fetch, time.Millisecond)
		assert.True(t, hasGracefulShutdown, "service should shutdown gracefully when context is cancelled")
	})
}
