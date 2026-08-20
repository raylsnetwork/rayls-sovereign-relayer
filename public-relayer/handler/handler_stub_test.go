// Decommissioning Teleport (vanilla, atomic).

package handler_test

import (
	"context"
	"errors"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/contracts/RNMessageDispatcherV1"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/contracts/PNTokenCoreV1"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/contracts/PNTokenRegistryV1"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/public-relayer/handler"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/public-relayer/service"
)

// errNotMatched is returned by parser stubs when a log doesn't match the expected event.
var errNotMatched = errors.New("log not matched")

// defaultTokenAddress is the token address shared by the default event and
// token-struct builders so that a triggered deployment resolves consistently.
var defaultTokenAddress = common.HexToAddress("0x3333333333333333333333333333333333333333")

// handlerStub provides default failing implementations for all PublicRelayerHandler dependencies.
// This allows tests to override only the dependencies they care about, with all others
// failing if accidentally called.
type handlerStub struct {
	t *testing.T

	MessagePublisher    func(context.Context, service.Message) error
	DeploymentPublisher func(context.Context, service.Deployment) error
	Dispatcher          func(*ethTypes.Log) (*RNMessageDispatcherV1.RNMessageDispatcherV1MessageDispatched, error)
	TokenCore           func(*ethTypes.Log) (*PNTokenCoreV1.PNTokenCoreV1PublicChainStatusUpdated, error)
	TokenRegistry       func(context.Context, common.Address) (PNTokenRegistryV1.TokenStructsToken, error)
}

// handlerStubOption represents an option for configuring the handlerStub.
type handlerStubOption func(*handlerStub)

// newHandlerStub creates a new stub with all dependencies set to fail by default.
func newHandlerStub(t *testing.T, opts ...handlerStubOption) *handlerStub {
	t.Helper()

	stub := &handlerStub{t: t}

	// Default: all parsers return error (log not recognized) which is normal behavior.
	// The handler uses err == nil checks, so returning an error means "this log didn't match".
	stub.Dispatcher = func(log *ethTypes.Log) (*RNMessageDispatcherV1.RNMessageDispatcherV1MessageDispatched, error) {
		return nil, errNotMatched
	}

	stub.TokenCore = func(log *ethTypes.Log) (*PNTokenCoreV1.PNTokenCoreV1PublicChainStatusUpdated, error) {
		return nil, errNotMatched
	}

	// Default: registry reader returns a well-formed token. Only reached after a
	// PublicChainStatusUpdated event matches, so tests that don't trigger one never hit it.
	stub.TokenRegistry = func(ctx context.Context, tokenAddress common.Address) (PNTokenRegistryV1.TokenStructsToken, error) {
		return newTokenStruct(), nil
	}

	// Default: publishers fail if called unexpectedly.
	stub.MessagePublisher = func(ctx context.Context, msg service.Message) error {
		t.Fatal("shouldn't call MessagePublisher.Push")
		return nil
	}

	stub.DeploymentPublisher = func(ctx context.Context, dpl service.Deployment) error {
		t.Fatal("shouldn't call DeploymentPublisher.Push")
		return nil
	}

	for _, opt := range opts {
		opt(stub)
	}

	return stub
}

// newPublicHandler creates a PublicRelayerHandler in public mode using the stub's configured dependencies.
func (s *handlerStub) newPublicHandler() *handler.PublicRelayerHandler {
	s.t.Helper()

	messageQueue := &MessagePublisherMock{PushFunc: s.MessagePublisher}
	dispatcher := &MessageDispatcherContractMock{UnpackMessageDispatchedEventFunc: s.Dispatcher}

	return handler.NewPublicRelayerPublicHandler(messageQueue, dispatcher)
}

// newPrivateHandler creates a PublicRelayerHandler in private mode using the stub's configured dependencies.
func (s *handlerStub) newPrivateHandler() *handler.PublicRelayerHandler {
	s.t.Helper()

	messageQueue := &MessagePublisherMock{PushFunc: s.MessagePublisher}
	dispatcher := &MessageDispatcherContractMock{UnpackMessageDispatchedEventFunc: s.Dispatcher}
	deploymentQueue := &DeploymentPublisherMock{PushFunc: s.DeploymentPublisher}
	tokenCore := &TokenCoreContractMock{UnpackPublicChainStatusUpdatedEventFunc: s.TokenCore}
	tokenRegistry := &TokenRegistryReaderMock{GetTokenByAddressFunc: s.TokenRegistry}

	return handler.NewPublicRelayerPrivateHandler(messageQueue, dispatcher, deploymentQueue, tokenCore, tokenRegistry)
}

// withDispatcher configures the MessageDispatcherContract dependency.
func withDispatcher(
	fn func(*ethTypes.Log) (*RNMessageDispatcherV1.RNMessageDispatcherV1MessageDispatched, error),
) handlerStubOption {
	return func(s *handlerStub) {
		s.Dispatcher = fn
	}
}

// withTokenCore configures the TokenCoreContract dependency (PublicChainStatusUpdated unpacker).
func withTokenCore(
	fn func(*ethTypes.Log) (*PNTokenCoreV1.PNTokenCoreV1PublicChainStatusUpdated, error),
) handlerStubOption {
	return func(s *handlerStub) {
		s.TokenCore = fn
	}
}

// withTokenRegistry configures the TokenRegistryReader dependency (GetTokenByAddress).
func withTokenRegistry(
	fn func(context.Context, common.Address) (PNTokenRegistryV1.TokenStructsToken, error),
) handlerStubOption {
	return func(s *handlerStub) {
		s.TokenRegistry = fn
	}
}

// withMessagePublisher configures the MessagePublisher dependency.
func withMessagePublisher(fn func(context.Context, service.Message) error) handlerStubOption {
	return func(s *handlerStub) {
		s.MessagePublisher = fn
	}
}

// withDeploymentPublisher configures the DeploymentPublisher dependency.
func withDeploymentPublisher(fn func(context.Context, service.Deployment) error) handlerStubOption {
	return func(s *handlerStub) {
		s.DeploymentPublisher = fn
	}
}

// newMessageDispatchedEvent creates a test MessageDispatched event with sensible defaults.
func newMessageDispatchedEvent(
	opts ...messageDispatchedOption,
) *RNMessageDispatcherV1.RNMessageDispatcherV1MessageDispatched {
	event := &RNMessageDispatcherV1.RNMessageDispatcherV1MessageDispatched{
		MessageId: [32]byte{0xaa, 0xbb, 0xcc},
		From:      common.HexToAddress("0x1111111111111111111111111111111111111111"),
		ToChainId: big.NewInt(1888),
		To:        common.HexToAddress("0x2222222222222222222222222222222222222222"),
		Data:      RNMessageDispatcherV1.RaylsNodeMessage{},
	}

	for _, opt := range opts {
		opt(event)
	}

	return event
}

type messageDispatchedOption func(*RNMessageDispatcherV1.RNMessageDispatcherV1MessageDispatched)

func withMessageID(id [32]byte) messageDispatchedOption {
	return func(e *RNMessageDispatcherV1.RNMessageDispatcherV1MessageDispatched) {
		e.MessageId = id
	}
}

func withFrom(from common.Address) messageDispatchedOption {
	return func(e *RNMessageDispatcherV1.RNMessageDispatcherV1MessageDispatched) {
		e.From = from
	}
}

func withToChainID(chainID *big.Int) messageDispatchedOption {
	return func(e *RNMessageDispatcherV1.RNMessageDispatcherV1MessageDispatched) {
		e.ToChainId = chainID
	}
}

func withTo(to common.Address) messageDispatchedOption {
	return func(e *RNMessageDispatcherV1.RNMessageDispatcherV1MessageDispatched) {
		e.To = to
	}
}

// publicChainStatusPendingDeployment is TokenStructs.PublicChainStatus.PENDING_DEPLOYMENT (1).
const publicChainStatusPendingDeployment uint8 = 1

// newPublicChainStatusUpdatedEvent creates a test PublicChainStatusUpdated event.
// It defaults to newStatus == PENDING_DEPLOYMENT (the relayer's deploy cue).
func newPublicChainStatusUpdatedEvent(
	opts ...publicChainStatusUpdatedOption,
) *PNTokenCoreV1.PNTokenCoreV1PublicChainStatusUpdated {
	event := &PNTokenCoreV1.PNTokenCoreV1PublicChainStatusUpdated{
		TokenAddress:   defaultTokenAddress,
		PreviousStatus: 0,
		NewStatus:      publicChainStatusPendingDeployment,
	}

	for _, opt := range opts {
		opt(event)
	}

	return event
}

type publicChainStatusUpdatedOption func(*PNTokenCoreV1.PNTokenCoreV1PublicChainStatusUpdated)

func withEventTokenAddress(addr common.Address) publicChainStatusUpdatedOption {
	return func(e *PNTokenCoreV1.PNTokenCoreV1PublicChainStatusUpdated) {
		e.TokenAddress = addr
	}
}

func withNewStatus(status uint8) publicChainStatusUpdatedOption {
	return func(e *PNTokenCoreV1.PNTokenCoreV1PublicChainStatusUpdated) {
		e.NewStatus = status
	}
}

// newTokenStruct creates a test registry Token entry with sensible defaults (ERC20).
func newTokenStruct(opts ...tokenStructOption) PNTokenRegistryV1.TokenStructsToken {
	token := PNTokenRegistryV1.TokenStructsToken{
		TokenAddress: defaultTokenAddress,
		ErcStandard:  1, // ERC20
		Uri:          "https://example.com/token",
		Name:         "TestToken",
		Symbol:       "TT",
	}

	for _, opt := range opts {
		opt(&token)
	}

	return token
}

type tokenStructOption func(*PNTokenRegistryV1.TokenStructsToken)

func withTokenAddress(addr common.Address) tokenStructOption {
	return func(t *PNTokenRegistryV1.TokenStructsToken) {
		t.TokenAddress = addr
	}
}

func withStandard(std uint8) tokenStructOption {
	return func(t *PNTokenRegistryV1.TokenStructsToken) {
		t.ErcStandard = std
	}
}

func withURI(uri string) tokenStructOption {
	return func(t *PNTokenRegistryV1.TokenStructsToken) {
		t.Uri = uri
	}
}

func withName(name string) tokenStructOption {
	return func(t *PNTokenRegistryV1.TokenStructsToken) {
		t.Name = name
	}
}

func withSymbol(symbol string) tokenStructOption {
	return func(t *PNTokenRegistryV1.TokenStructsToken) {
		t.Symbol = symbol
	}
}
