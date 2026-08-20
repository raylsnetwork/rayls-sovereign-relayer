// Decommissioning Teleport (vanilla, atomic).

package service

import (
	"context"
	"log/slog"

	"github.com/ethereum/go-ethereum/common"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/faultinjector"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/msgqueue"
)

//go:generate moq --pkg service_test -out deployer_mock_test.go . DeploymentConsumer DeployerClient TokenGovernanceClient AccessManagerClient
type DeploymentConsumer interface {
	Next(ctx context.Context) (msgqueue.Message[Deployment], error)
}

type DeployerClient interface {
	DeployPublicChainERC20(
		ctx context.Context,
		name, symbol string,
		privateAddress, raylsPublicRelayerEndpoint common.Address,
	) (common.Address, error)
	DeployPublicChainERC721(
		ctx context.Context,
		uri, name, symbol string,
		privateAddress, raylsPublicRelayerEndpoint common.Address,
	) (common.Address, error)
	DeployPublicChainERC1155(
		ctx context.Context,
		uri, name string,
		privateAddress, raylsPublicRelayerEndpoint common.Address,
	) (common.Address, error)
}

type TokenGovernanceClient interface {
	UpdatePublicTokenAddress(ctx context.Context, privateAddress, publicAddress common.Address) error
	GetPublicAddressByPrivateAddress(ctx context.Context, privateAddress common.Address) (common.Address, error)
}

type AccessManagerClient interface {
	GrantAuthorizedSenderRole(ctx context.Context, tokenAddress common.Address) error
}

type DeployerService struct {
	endpointAddress common.Address

	deployerClient        DeployerClient
	tokenGovernanceClient TokenGovernanceClient
	accessManagerClient   AccessManagerClient

	consumer DeploymentConsumer
}

func NewDeployerService(
	consumer DeploymentConsumer,
	endpointAddress common.Address,
	deployerClient DeployerClient,
	tokenGovernanceClient TokenGovernanceClient,
	accessManagerClient AccessManagerClient,
) *DeployerService {
	return &DeployerService{
		consumer: consumer,

		endpointAddress: endpointAddress,

		deployerClient:        deployerClient,
		tokenGovernanceClient: tokenGovernanceClient,
		accessManagerClient:   accessManagerClient,
	}
}

func (s *DeployerService) Deploy(ctx context.Context) error {
	for {
		msg, err := s.consumer.Next(ctx)
		if err != nil {
			if ctx.Err() != nil {
				return nil //nolint:nilerr // context cancellation is a clean shutdown signal
			}
			slog.Error("Failed to read from deployment queue", slog.Any("error", err))
			continue
		}

		dpl := msg.V

		onChainPublicAddress, err := s.tokenGovernanceClient.GetPublicAddressByPrivateAddress(ctx, dpl.PrivateAddress)
		if err != nil {
			slog.Error("Failed to get public address for token", slog.String("name", dpl.Name), slog.Any("error", err))
			continue // un-acked, NATS will redeliver
		}

		var emptyAddress common.Address
		if onChainPublicAddress != emptyAddress {
			if ackErr := msg.Ack(ctx); ackErr != nil {
				slog.Error("Failed to acknowledge deployment message", slog.Any("error", ackErr))
			}
			continue
		}

		var publicAddress common.Address
		switch dpl.Standard {
		case ERC20:
			publicAddress, err = s.deployerClient.DeployPublicChainERC20(
				ctx,
				dpl.Name,
				dpl.Symbol,
				dpl.PrivateAddress,
				s.endpointAddress,
			)
		case ERC721:
			publicAddress, err = s.deployerClient.DeployPublicChainERC721(
				ctx,
				dpl.URI,
				dpl.Name,
				dpl.Symbol,
				dpl.PrivateAddress,
				s.endpointAddress,
			)
		case ERC1155:
			publicAddress, err = s.deployerClient.DeployPublicChainERC1155(
				ctx,
				dpl.URI,
				dpl.Name,
				dpl.PrivateAddress,
				s.endpointAddress,
			)
		default:
			slog.Error(
				"Unsupported token standard for public chain deployment",
				slog.Int("standard", int(dpl.Standard)),
				slog.String("name", dpl.Name),
			)
			continue
		}
		if err != nil {
			slog.Error("Failed to deploy token", slog.Any("name", dpl.Name), slog.Any("error", err))
			continue // un-acked, NATS will redeliver
		}
		slog.Info("Successfully deployed token", slog.Any("name", dpl.Name))

		err = s.tokenGovernanceClient.UpdatePublicTokenAddress(ctx, dpl.PrivateAddress, publicAddress)
		if err != nil {
			slog.Error("Failed to set public address for token", slog.String("name", dpl.Name), slog.Any("error", err))
		}

		// Fault injection: the public-chain mapping for the token is now
		// live on-chain (UpdatePublicTokenAddress committed) but NATS hasn't
		// been acked. A crash here triggers redelivery; the idempotency
		// guard at GetPublicAddressByPrivateAddress (line 84 above) is what
		// keeps the next attempt from re-deploying the same token.
		if faultErr := faultinjector.Check("public_relayer.service.DeployerService.Deploy.after_governance_update"); faultErr != nil {
			slog.Error("Fault injection at after_governance_update", slog.Any("error", faultErr))
			continue
		}

		// Grant the AUTHORIZED_SENDER to the newly deployed token via the access manager.
		// If granting fails the message is still acked because the token is already deployed
		// and its governance address recorded. Re-delivering would re-attempt the full deploy
		// flow. Role grant failures should be resolved via manual intervention or a separate
		// reconciliation process.
		err = s.accessManagerClient.GrantAuthorizedSenderRole(ctx, publicAddress)
		if err != nil {
			slog.Error(
				"Failed to grant AUTHORIZED_SENDER to token",
				slog.String("name", dpl.Name),
				slog.String("publicAddress", publicAddress.Hex()),
				slog.Any("error", err),
			)
		} else {
			slog.Info(
				"Successfully granted AUTHORIZED_SENDER to token",
				slog.String("name", dpl.Name),
				slog.String("publicAddress", publicAddress.Hex()),
			)
		}

		// Fault injection: deployment side-effects (chain deploy, governance
		// update, role grant) have all run but NATS hasn't been acked.
		// A crash here triggers redelivery; same idempotency guard as above.
		if faultErr := faultinjector.Check("public_relayer.service.DeployerService.Deploy.before_ack"); faultErr != nil {
			slog.Error("Fault injection at before_ack", slog.Any("error", faultErr))
			continue
		}
		if err := msg.Ack(ctx); err != nil {
			slog.Error("Failed to acknowledge deployment message", slog.Any("error", err))
		}
	}
}
