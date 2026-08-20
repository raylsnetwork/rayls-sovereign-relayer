// Decommissioning Teleport (vanilla, atomic).

// Package handler implements the legacy public-chain (RN) Teleport bridge relayer.
//
// Deprecated: Decommissioning Teleport (vanilla, atomic).
package handler

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/ethereum/go-ethereum/common"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/contracts/RNMessageDispatcherV1"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/contracts/PNTokenCoreV1"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/contracts/PNTokenRegistryV1"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/faultinjector"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/public-relayer/service"
)

//go:generate moq --pkg handler_test -out handler_mock_test.go . MessageDispatcherContract TokenCoreContract TokenRegistryReader MessagePublisher DeploymentPublisher

// publicChainStatusPendingDeployment is TokenStructs.PublicChainStatus.PENDING_DEPLOYMENT (1):
// the status the PN TokenCore transitions to on submitToPublicChain(), which is
// the relayer's cue to deploy the public-chain token mirror.
const publicChainStatusPendingDeployment uint8 = 1

type MessageDispatcherContract interface {
	UnpackMessageDispatchedEvent(log *ethTypes.Log) (*RNMessageDispatcherV1.RNMessageDispatcherV1MessageDispatched, error)
}

// TokenCoreContract unpacks the PNTokenCoreV1 module's PublicChainStatusUpdated
// event, emitted when a token's public-chain lifecycle status changes.
type TokenCoreContract interface {
	UnpackPublicChainStatusUpdatedEvent(log *ethTypes.Log) (*PNTokenCoreV1.PNTokenCoreV1PublicChainStatusUpdated, error)
}

// TokenRegistryReader reads a token's full registry entry from the
// PNTokenRegistryV1 facade. The PublicChainStatusUpdated event carries no token
// metadata, so name/symbol/uri/standard are fetched here.
type TokenRegistryReader interface {
	GetTokenByAddress(ctx context.Context, tokenAddress common.Address) (PNTokenRegistryV1.TokenStructsToken, error)
}

type MessagePublisher interface {
	Push(context.Context, service.Message) error
}

type DeploymentPublisher interface {
	Push(context.Context, service.Deployment) error
}

type PublicRelayerHandler struct {
	isPrivate bool

	deploymentQueue DeploymentPublisher
	tokenCore       TokenCoreContract
	tokenRegistry   TokenRegistryReader

	messageQueue MessagePublisher
	dispatcher   MessageDispatcherContract
}

func NewPublicRelayerPublicHandler(
	messageQueue MessagePublisher,
	dispatcher MessageDispatcherContract,
) *PublicRelayerHandler {
	return &PublicRelayerHandler{
		isPrivate: false,

		messageQueue: messageQueue,
		dispatcher:   dispatcher,
	}
}

func NewPublicRelayerPrivateHandler(
	messageQueue MessagePublisher,
	dispatcher MessageDispatcherContract,
	deploymentQueue DeploymentPublisher,
	tokenCore TokenCoreContract,
	tokenRegistry TokenRegistryReader,
) *PublicRelayerHandler {
	return &PublicRelayerHandler{
		isPrivate: true,

		messageQueue:    messageQueue,
		dispatcher:      dispatcher,
		deploymentQueue: deploymentQueue,
		tokenCore:       tokenCore,
		tokenRegistry:   tokenRegistry,
	}
}

func (h *PublicRelayerHandler) Handle(ctx context.Context, logs []ethTypes.Log) error {
	slog.Debug("Handling logs...", slog.Int("count", len(logs)))

	var count int

	for _, log := range logs {
		if event, err := h.dispatcher.UnpackMessageDispatchedEvent(&log); err == nil {
			h.processMessageDispatched(ctx, event)
			count++
		}

		if h.isPrivate {
			if event, err := h.tokenCore.UnpackPublicChainStatusUpdatedEvent(&log); err == nil {
				h.processPublicChainStatusUpdated(ctx, event)
			}
		}
	}

	if count != 0 {
		slog.Info("Pushed messages to queue", slog.Int("COUNT", count))
	}

	// Block-level events are now durable in NATS but the listener's
	// last_processed_block checkpoint hasn't moved yet — returning an error
	// here prevents the checkpoint from advancing so redelivery is exercised
	// on next listener iteration.
	if err := faultinjector.Check("public_relayer.handler.PublicRelayerHandler.Handle.after_push_messages"); err != nil {
		return fmt.Errorf("fault injection at after_push_messages: %w", err)
	}

	return nil
}

func (h *PublicRelayerHandler) processMessageDispatched(
	ctx context.Context,
	event *RNMessageDispatcherV1.RNMessageDispatcherV1MessageDispatched,
) {
	message := service.Message{
		ID:          common.BytesToHash(event.MessageId[:]),
		FromAddress: event.From,
		ToChainID:   event.ToChainId,
		ToAddress:   event.To,
		Data:        event.Data,
	}
	err := h.messageQueue.Push(ctx, message)
	if err != nil {
		slog.Error("Failed to push endpoint message", slog.Any("error", err))
	}
}

// processPublicChainStatusUpdated reacts to a token's public-chain status
// transition. Only PENDING_DEPLOYMENT (set by submitToPublicChain on the PN) is
// the relayer's cue to deploy the public-chain mirror. The event carries no
// token metadata, so it is fetched from the registry facade.
func (h *PublicRelayerHandler) processPublicChainStatusUpdated(
	ctx context.Context,
	event *PNTokenCoreV1.PNTokenCoreV1PublicChainStatusUpdated,
) {
	if event.NewStatus != publicChainStatusPendingDeployment {
		return
	}

	token, err := h.tokenRegistry.GetTokenByAddress(ctx, event.TokenAddress)
	if err != nil {
		slog.Error("Failed to get token metadata", slog.Any("tokenAddress", event.TokenAddress), slog.Any("error", err))
		return
	}

	standard, err := getTokenStandard(token.ErcStandard)
	if err != nil {
		slog.Error("Failed to get token standard", slog.Int("standard", int(token.ErcStandard)), slog.Any("error", err))
		return
	}

	deployment := service.Deployment{
		PrivateAddress: token.TokenAddress,

		Standard: standard,

		URI:    token.Uri,
		Name:   token.Name,
		Symbol: token.Symbol,
	}

	err = h.deploymentQueue.Push(ctx, deployment)
	if err != nil {
		slog.Error("Failed to push deployment message", slog.Any("error", err))
	}
}

// getTokenStandard maps RaylsNodeBridgeableERC enum values to TokenStandard.
// Only ERC20, ERC721, and ERC1155 are deployable on the public chain.
// CUSTOM, ENYGMA, and ZKDVP* types do not need public chain mirrors.
func getTokenStandard(standard uint8) (service.TokenStandard, error) {
	ts := service.TokenStandard(standard)
	switch ts {
	case service.ERC20:
		return service.ERC20, nil
	case service.ERC721:
		return service.ERC721, nil
	case service.ERC1155:
		return service.ERC1155, nil
	case service.CUSTOM:
		return ts, fmt.Errorf("CUSTOM tokens do not need public chain mirrors")
	case service.ENYGMA:
		return ts, fmt.Errorf("ENYGMA tokens use Path C (private relayer), not Path B")
	case service.ZKDVPERC721:
		return ts, fmt.Errorf("ZKDVPERC721 tokens use the DvP flow, not a public chain mirror")
	case service.ZKDVPERC1155:
		return ts, fmt.Errorf("ZKDVPERC1155 tokens use the DvP flow, not a public chain mirror")
	default:
		return service.TokenStandard(-1), fmt.Errorf("unknown token standard: %d", standard)
	}
}
