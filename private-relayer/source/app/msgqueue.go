package app

import (
	"context"
	"fmt"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/msgqueue"
	sharedservice "github.com/raylsnetwork/rayls-privacy-relayer-api/private-relayer/shared/service"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/private-relayer/source/logrouter"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/private-relayer/source/service"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/types"
)

type MessageQueueConfig struct {
	NATSUrl string
}

type MessageQueue struct {
	nc *nats.Conn

	endpointBlockPublisher *msgqueue.Publisher[logrouter.Block]
	enygmaBlockPublisher   *msgqueue.Publisher[logrouter.Block]

	endpointBlockConsumer *msgqueue.Consumer[logrouter.Block]
	enygmaBlockConsumer   *msgqueue.Consumer[logrouter.Block]

	crossChainPublisher *msgqueue.Publisher[service.CrossChainMessage]
	crossChainConsumer  *msgqueue.Consumer[service.CrossChainMessage]

	privateHubPublisher *msgqueue.Publisher[sharedservice.PrivateHubMessage]
	privateHubConsumer  *msgqueue.Consumer[sharedservice.PrivateHubMessage]

	enygmaBatchPublisher *msgqueue.Publisher[service.EnygmaSerializedEvent]
	enygmaBatchConsumer  *msgqueue.Consumer[service.EnygmaSerializedEvent]

	dvpBatchPublisher *msgqueue.Publisher[service.DvpSerializedEventBatch]
	dvpBatchConsumer  *msgqueue.Consumer[service.DvpSerializedEventBatch]

	privateHubSendPublisher  *msgqueue.Publisher[types.TxRequest]
	privateNodeSendPublisher *msgqueue.Publisher[types.TxRequest]
}

func (r *SourcePrivateRelayer) initializeMessageQueues(chainId string) error {
	js, err := jetstream.New(r.natsConn)
	if err != nil {
		return fmt.Errorf("failed to create jetstream connection: %w", err)
	}

	// Create context with timeout for message queue initialization
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	manager, err := msgqueue.NewManager(ctx, js, chainId)
	if err != nil {
		return fmt.Errorf("failed to create message queue manager: %w", err)
	}

	endpointBlockPublisher := msgqueue.NewPublisher[logrouter.Block](manager, "endpoint_blocks")
	enygmaBlockPublisher := msgqueue.NewPublisher[logrouter.Block](manager, "enygma_blocks")

	endpointBlockConsumer, err := msgqueue.NewConsumer[logrouter.Block](
		ctx,
		manager,
		"endpoint_parser",
		"endpoint_blocks",
	)
	if err != nil {
		return fmt.Errorf("failed to initialize endpoint block consumer: %w", err)
	}

	enygmaBlockConsumer, err := msgqueue.NewConsumer[logrouter.Block](ctx, manager, "enygma_parser", "enygma_blocks")
	if err != nil {
		return fmt.Errorf("failed to initialize enygma block consumer: %w", err)
	}

	crossChainPublisher := msgqueue.NewPublisher[service.CrossChainMessage](manager, "cross_chain_events")
	crossChainConsumer, err := msgqueue.NewConsumer[service.CrossChainMessage](
		ctx,
		manager,
		"cross_chain_service",
		"cross_chain_events",
	)
	if err != nil {
		return fmt.Errorf("failed to initialize cross chain consumer: %w", err)
	}

	privateHubPublisher := msgqueue.NewPublisher[sharedservice.PrivateHubMessage](manager, "private_hub_events")
	privateHubConsumer, err := msgqueue.NewConsumer[sharedservice.PrivateHubMessage](
		ctx,
		manager,
		"private_hub_service",
		"private_hub_events",
	)
	if err != nil {
		return fmt.Errorf("failed to initialize private hub consumer: %w", err)
	}

	enygmaBatchPublisher := msgqueue.NewPublisher[service.EnygmaSerializedEvent](manager, "enygma_batches")
	enygmaBatchConsumer, err := msgqueue.NewConsumer[service.EnygmaSerializedEvent](ctx, manager, "enygma_orchestrator", "enygma_batches")
	if err != nil {
		return fmt.Errorf("failed to initialize enygma batch consumer: %w", err)
	}

	dvpBatchPublisher := msgqueue.NewPublisher[service.DvpSerializedEventBatch](manager, "dvp_batches")
	dvpBatchConsumer, err := msgqueue.NewConsumer[service.DvpSerializedEventBatch](
		ctx,
		manager,
		"dvp_orchestrator",
		"dvp_batches",
	)
	if err != nil {
		return fmt.Errorf("failed to initialize dvp batch consumer: %w", err)
	}

	privateHubSendPublisher := msgqueue.NewPublisher[types.TxRequest](manager, "cts.send.privatehub")
	privateNodeSendPublisher := msgqueue.NewPublisher[types.TxRequest](manager, "cts.send.privatenode")

	r.msgqueues = &MessageQueue{
		nc: r.natsConn,

		endpointBlockPublisher: endpointBlockPublisher,
		enygmaBlockPublisher:   enygmaBlockPublisher,

		endpointBlockConsumer: endpointBlockConsumer,
		enygmaBlockConsumer:   enygmaBlockConsumer,

		crossChainPublisher: crossChainPublisher,
		crossChainConsumer:  crossChainConsumer,

		privateHubPublisher: privateHubPublisher,
		privateHubConsumer:  privateHubConsumer,

		enygmaBatchPublisher: enygmaBatchPublisher,
		enygmaBatchConsumer:  enygmaBatchConsumer,

		dvpBatchPublisher: dvpBatchPublisher,
		dvpBatchConsumer:  dvpBatchConsumer,

		privateHubSendPublisher:  privateHubSendPublisher,
		privateNodeSendPublisher: privateNodeSendPublisher,
	}

	return nil
}

func (q *MessageQueue) Close() {
	_ = q.nc.Drain()
}
