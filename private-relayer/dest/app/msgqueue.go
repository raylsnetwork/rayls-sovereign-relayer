package app

import (
	"context"
	"fmt"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/msgqueue"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/private-relayer/dest/logrouter"
	destservice "github.com/raylsnetwork/rayls-privacy-relayer-api/private-relayer/dest/service"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/private-relayer/shared/service"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/types"
)

type MessageQueue struct {
	nc *nats.Conn

	// Routed log queues (by contract)
	endpointBlockPublisher     *msgqueue.Publisher[logrouter.Block]
	teleportBlockPublisher     *msgqueue.Publisher[logrouter.Block]
	enygmaBlockPublisher       *msgqueue.Publisher[logrouter.Block]
	dvpBlockPublisher          *msgqueue.Publisher[logrouter.Block]
	auditManagerBlockPublisher *msgqueue.Publisher[logrouter.Block]

	endpointBlockConsumer     *msgqueue.Consumer[logrouter.Block]
	teleportBlockConsumer     *msgqueue.Consumer[logrouter.Block]
	enygmaBlockConsumer       *msgqueue.Consumer[logrouter.Block]
	dvpBlockConsumer          *msgqueue.Consumer[logrouter.Block]
	auditManagerBlockConsumer *msgqueue.Consumer[logrouter.Block]

	// Service message queues (parsed events)
	crossChainPublisher *msgqueue.Publisher[types.DispatchedMessageToPrivateHub]
	crossChainConsumer  *msgqueue.Consumer[types.DispatchedMessageToPrivateHub]

	privateHubMessagePublisher *msgqueue.Publisher[service.PrivateHubMessage]
	privateHubMessageConsumer  *msgqueue.Consumer[service.PrivateHubMessage]

	enygmaDestPublisher *msgqueue.Publisher[destservice.EnygmaDestMessage]
	enygmaDestConsumer  *msgqueue.Consumer[destservice.EnygmaDestMessage]

	dvpDestPublisher *msgqueue.Publisher[destservice.DvpDestMessage]
	dvpDestConsumer  *msgqueue.Consumer[destservice.DvpDestMessage]

	privateNodeSendPublisher *msgqueue.Publisher[types.TxRequest]
}

func (r *DestPrivateRelayer) initializeMessageQueues(chainId string) error {
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

	// Routed log queues
	endpointBlockPublisher := msgqueue.NewPublisher[logrouter.Block](manager, "dest_endpoint_blocks")
	teleportBlockPublisher := msgqueue.NewPublisher[logrouter.Block](manager, "dest_teleport_blocks")
	enygmaBlockPublisher := msgqueue.NewPublisher[logrouter.Block](manager, "dest_enygma_blocks")
	dvpBlockPublisher := msgqueue.NewPublisher[logrouter.Block](manager, "dest_dvp_blocks")
	auditManagerBlockPublisher := msgqueue.NewPublisher[logrouter.Block](manager, "dest_audit_manager_blocks")

	endpointBlockConsumer, err := msgqueue.NewConsumer[logrouter.Block](
		ctx,
		manager,
		"dest_endpoint_parser",
		"dest_endpoint_blocks",
	)
	if err != nil {
		return fmt.Errorf("failed to initialize endpoint block consumer: %w", err)
	}

	teleportBlockConsumer, err := msgqueue.NewConsumer[logrouter.Block](
		ctx,
		manager,
		"dest_teleport_parser",
		"dest_teleport_blocks",
	)
	if err != nil {
		return fmt.Errorf("failed to initialize teleport block consumer: %w", err)
	}

	enygmaBlockConsumer, err := msgqueue.NewConsumer[logrouter.Block](
		ctx,
		manager,
		"dest_enygma_parser",
		"dest_enygma_blocks",
	)
	if err != nil {
		return fmt.Errorf("failed to initialize enygma block consumer: %w", err)
	}

	dvpBlockConsumer, err := msgqueue.NewConsumer[logrouter.Block](ctx, manager, "dest_dvp_parser", "dest_dvp_blocks")
	if err != nil {
		return fmt.Errorf("failed to initialize dvp block consumer: %w", err)
	}

	auditManagerBlockConsumer, err := msgqueue.NewConsumer[logrouter.Block](
		ctx,
		manager,
		"dest_audit_manager_parser",
		"dest_audit_manager_blocks",
	)
	if err != nil {
		return fmt.Errorf("failed to initialize audit manager block consumer: %w", err)
	}

	// Service message queues
	crossChainPublisher := msgqueue.NewPublisher[types.DispatchedMessageToPrivateHub](
		manager,
		"dest_crosschain_messages",
	)
	crossChainConsumer, err := msgqueue.NewConsumer[types.DispatchedMessageToPrivateHub](
		ctx,
		manager,
		"dest_crosschain_service",
		"dest_crosschain_messages",
	)
	if err != nil {
		return fmt.Errorf("failed to initialize cross chain consumer: %w", err)
	}

	privateHubMessagePublisher := msgqueue.NewPublisher[service.PrivateHubMessage](manager, "dest_private_hub_messages")
	privateHubMessageConsumer, err := msgqueue.NewConsumer[service.PrivateHubMessage](
		ctx,
		manager,
		"dest_crosschain_service_hub",
		"dest_private_hub_messages",
	)
	if err != nil {
		return fmt.Errorf("failed to initialize private hub message consumer: %w", err)
	}

	enygmaDestPublisher := msgqueue.NewPublisher[destservice.EnygmaDestMessage](manager, "dest_enygma_messages")
	enygmaDestConsumer, err := msgqueue.NewConsumer[destservice.EnygmaDestMessage](
		ctx,
		manager,
		"dest_enygma_orchestrator",
		"dest_enygma_messages",
	)
	if err != nil {
		return fmt.Errorf("failed to initialize enygma dest consumer: %w", err)
	}

	dvpDestPublisher := msgqueue.NewPublisher[destservice.DvpDestMessage](manager, "dest_dvp_messages")
	dvpDestConsumer, err := msgqueue.NewConsumer[destservice.DvpDestMessage](
		ctx,
		manager,
		"dest_dvp_orchestrator",
		"dest_dvp_messages",
	)
	if err != nil {
		return fmt.Errorf("failed to initialize dvp dest consumer: %w", err)
	}

	privateNodeSendPublisher := msgqueue.NewPublisher[types.TxRequest](manager, "cts.send.privatenode")

	r.msgqueues = &MessageQueue{
		nc: r.natsConn,

		endpointBlockPublisher:     endpointBlockPublisher,
		teleportBlockPublisher:     teleportBlockPublisher,
		enygmaBlockPublisher:       enygmaBlockPublisher,
		dvpBlockPublisher:          dvpBlockPublisher,
		auditManagerBlockPublisher: auditManagerBlockPublisher,

		endpointBlockConsumer:     endpointBlockConsumer,
		teleportBlockConsumer:     teleportBlockConsumer,
		enygmaBlockConsumer:       enygmaBlockConsumer,
		dvpBlockConsumer:          dvpBlockConsumer,
		auditManagerBlockConsumer: auditManagerBlockConsumer,

		crossChainPublisher: crossChainPublisher,
		crossChainConsumer:  crossChainConsumer,

		privateHubMessagePublisher: privateHubMessagePublisher,
		privateHubMessageConsumer:  privateHubMessageConsumer,

		enygmaDestPublisher: enygmaDestPublisher,
		enygmaDestConsumer:  enygmaDestConsumer,

		dvpDestPublisher: dvpDestPublisher,
		dvpDestConsumer:  dvpDestConsumer,

		privateNodeSendPublisher: privateNodeSendPublisher,
	}

	return nil
}

func (q *MessageQueue) Close() {
	_ = q.nc.Drain()
}
