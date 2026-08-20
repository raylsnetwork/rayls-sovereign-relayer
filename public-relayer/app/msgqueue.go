// Decommissioning Teleport (vanilla, atomic).

package app

import (
	"context"
	"fmt"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/msgqueue"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/public-relayer/service"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/types"
)

const msgQueueInitTimeout = 30 * time.Second

type MessageQueueConfig struct {
	NATSConn *nats.Conn
	ChainId  string
}

// MessageQueues owns all NATS publishers and consumers used by the public
// relayer.
//
// CTS bridge subjects: TxRequests are addressed to a CTS signing identity
// (publicchain | privatechain) via `cts.send.<identity>`. CTS publishes
// terminal TxResults back on `cts.result.<identity>`. The public relayer
// uses both identities — publicGenerator's forward goes to private chain
// (privatechain identity), its revert goes to public chain (publicchain
// identity); privateGenerator is the mirror image. Publishers and result
// consumers are therefore keyed by identity, not by leg.
type MessageQueues struct {
	publicPublisher *msgqueue.Publisher[service.Message]
	publicConsumer  *msgqueue.Consumer[service.Message]

	privatePublisher *msgqueue.Publisher[service.Message]
	privateConsumer  *msgqueue.Consumer[service.Message]

	deploymentPublisher *msgqueue.Publisher[service.Deployment]
	deploymentConsumer  *msgqueue.Consumer[service.Deployment]

	// CTS bridge — TxRequest publishers, one per signing identity.
	publicChainSendPublisher  *msgqueue.Publisher[types.TxRequest]
	privateChainSendPublisher *msgqueue.Publisher[types.TxRequest]

	// CTS bridge — TxResult consumers, one per signing identity.
	publicChainResultConsumer  *msgqueue.Consumer[types.TxResult]
	privateChainResultConsumer *msgqueue.Consumer[types.TxResult]
}

func (p *PublicRelayer) initializeMessageQueues(config MessageQueueConfig) error {
	js, err := jetstream.New(p.natsConn)
	if err != nil {
		return fmt.Errorf("failed to create jetstream: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), msgQueueInitTimeout)
	defer cancel()

	manager, err := msgqueue.NewManager(ctx, js, config.ChainId)
	if err != nil {
		return fmt.Errorf("failed to create message queue manager: %w", err)
	}

	publicPublisher := msgqueue.NewPublisher[service.Message](manager, "public_messages")
	publicConsumer, err := msgqueue.NewConsumer[service.Message](ctx, manager, "public_generator", "public_messages")
	if err != nil {
		return fmt.Errorf("failed to create public message consumer: %w", err)
	}

	privatePublisher := msgqueue.NewPublisher[service.Message](manager, "private_messages")
	privateConsumer, err := msgqueue.NewConsumer[service.Message](ctx, manager, "private_generator", "private_messages")
	if err != nil {
		return fmt.Errorf("failed to create private message consumer: %w", err)
	}

	deploymentPublisher := msgqueue.NewPublisher[service.Deployment](manager, "deployments")
	deploymentConsumer, err := msgqueue.NewConsumer[service.Deployment](ctx, manager, "deployer", "deployments")
	if err != nil {
		return fmt.Errorf("failed to create deployment consumer: %w", err)
	}

	publicChainSendPublisher := msgqueue.NewPublisher[types.TxRequest](manager, "cts.send.publicchain")
	privateChainSendPublisher := msgqueue.NewPublisher[types.TxRequest](manager, "cts.send.privatechain")

	publicChainResultConsumer, err := msgqueue.NewConsumer[types.TxResult](
		ctx, manager, "public_relayer_publicchain_results", "cts.result.publicchain",
	)
	if err != nil {
		return fmt.Errorf("failed to create publicchain result consumer: %w", err)
	}
	privateChainResultConsumer, err := msgqueue.NewConsumer[types.TxResult](
		ctx, manager, "public_relayer_privatechain_results", "cts.result.privatechain",
	)
	if err != nil {
		return fmt.Errorf("failed to create privatechain result consumer: %w", err)
	}

	p.messageQueues = &MessageQueues{
		publicPublisher: publicPublisher,
		publicConsumer:  publicConsumer,

		privatePublisher: privatePublisher,
		privateConsumer:  privateConsumer,

		deploymentPublisher: deploymentPublisher,
		deploymentConsumer:  deploymentConsumer,

		publicChainSendPublisher:  publicChainSendPublisher,
		privateChainSendPublisher: privateChainSendPublisher,

		publicChainResultConsumer:  publicChainResultConsumer,
		privateChainResultConsumer: privateChainResultConsumer,
	}
	return nil
}

func (m *MessageQueues) Close() {
	// Note: We don't drain the NATS connection here because it's owned by the caller.
	// Consumers and publishers are cleaned up when the NATS connection is drained by the caller.
}
