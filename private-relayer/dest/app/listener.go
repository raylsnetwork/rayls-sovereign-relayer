package app

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/listener"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/private-relayer/dest/logrouter"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/types"
)

type ListenersConfig struct {
	BatchSize               int
	PrivateHubStartingBlock *big.Int
}

type Listeners struct {
	privateHub *listener.LogListener
}

// LogRouter wraps the actual logrouter.LogRouter and implements listener.LogHandler
type LogRouter struct {
	router *logrouter.LogRouter
}

// Handle implements listener.LogHandler interface
func (l *LogRouter) Handle(ctx context.Context, logs []ethTypes.Log) error {
	return l.router.Handle(ctx, logs)
}

func (r *DestPrivateRelayer) initializeListeners(conf ListenersConfig) error {
	// Resolve contract addresses from hub deployment proxy registry
	endpointAddress, err := r.contractClients.hubRegistry.GetContractAddress("Endpoint")
	if err != nil {
		return fmt.Errorf("failed to get hub endpoint address: %w", err)
	}

	teleportAddress, err := r.contractClients.hubRegistry.GetContractAddress("Teleport")
	if err != nil {
		return fmt.Errorf("failed to get hub teleport address: %w", err)
	}

	enygmaTeleportAddress, err := r.contractClients.hubRegistry.GetContractAddress("EnygmaTeleport")
	if err != nil {
		return fmt.Errorf("failed to get hub enygma teleport address: %w", err)
	}

	dvpTeleportAddress, err := r.contractClients.hubRegistry.GetContractAddress("DvpTeleport")
	if err != nil {
		return fmt.Errorf("failed to get hub dvp teleport address: %w", err)
	}

	auditManagerAddress, err := r.contractClients.hubRegistry.GetContractAddress("AuditManager")
	if err != nil {
		return fmt.Errorf("failed to get hub audit manager address: %w", err)
	}

	// Create LogRouter (routes logs by contract address to appropriate message queues)
	logRouter := logrouter.New(
		logrouter.LogRouterConfig{
			EndpointAddress:     endpointAddress,
			TeleportAddress:     teleportAddress,
			EnygmaAddress:       enygmaTeleportAddress,
			DvpAddress:          dvpTeleportAddress,
			AuditManagerAddress: auditManagerAddress,
		},
		r.msgqueues.endpointBlockPublisher,
		r.msgqueues.teleportBlockPublisher,
		r.msgqueues.enygmaBlockPublisher,
		r.msgqueues.dvpBlockPublisher,
		r.msgqueues.auditManagerBlockPublisher,
	)

	r.logRouter = &LogRouter{
		router: logRouter,
	}

	// Create Private Hub listener - listens to messages coming from source on the Private Hub
	privateHubListenerConfig := listener.LogListenerConfig{
		Component:     types.DocumentIdLastProcessedBlockPrivateHub,
		StartingBlock: conf.PrivateHubStartingBlock,
		BatchSize:     conf.BatchSize,
		Addresses: []common.Address{
			endpointAddress,
			teleportAddress,
			enygmaTeleportAddress,
			dvpTeleportAddress,
			auditManagerAddress,
		},
	}

	privateHubListener, err := listener.NewLogListener(
		context.Background(),
		privateHubListenerConfig,
		r.logRouter, // LogRouter implements LogHandler interface
		r.hubClient, // Private Hub client
		r.repositories.lastProcessedBlockRepository,
	)
	if err != nil {
		return fmt.Errorf("failed to create private hub listener: %w", err)
	}

	r.listeners = &Listeners{
		privateHub: privateHubListener,
	}

	return nil
}
