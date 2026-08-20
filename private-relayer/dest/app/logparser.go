package app

import (
	"math/big"

	"github.com/raylsnetwork/rayls-sovereign-relayer/private-relayer/dest/logparser"
	destservice "github.com/raylsnetwork/rayls-sovereign-relayer/private-relayer/dest/service"
)

type LogParsersConfig struct {
	PrivateNodeChainID *big.Int
}

type LogParsers struct {
	endpointParser       *logparser.EndpointLogParser
	teleportParser       *logparser.TeleportParser
	enygmaTeleportParser *logparser.EnygmaTeleportParser
	dvpTeleportParser    *logparser.DvpTeleportParser
	auditManagerParser   *logparser.AuditManagerParser
}

func (r *DestPrivateRelayer) initializeLogParsers(conf LogParsersConfig) error {
	// Initialize endpoint parser (listens to endpoint logs)
	endpointParser := logparser.NewEndpointLogParser(
		conf.PrivateNodeChainID,
		r.msgqueues.endpointBlockConsumer,
		r.msgqueues.privateHubMessagePublisher,
	)

	// Initialize enygma teleport parser (listens to enygma logs)
	enygmaTeleportParser := logparser.NewEnygmaTeleportParser(
		r.msgqueues.enygmaBlockConsumer,
		r.msgqueues.enygmaDestPublisher,
		r.ctsClient,
		conf.PrivateNodeChainID,
	)

	// Initialize dvp teleport parser (listens to dvp logs)
	dvpTeleportParser := logparser.NewDvpTeleportParser(
		r.msgqueues.dvpBlockConsumer,
		r.msgqueues.dvpDestPublisher,
		r.ctsClient,
		r.hubClient,
		r.repositories.dvpSwap,
		conf.PrivateNodeChainID,
	)

	// Initialize SUMService for atomic status update messages
	sumService := destservice.NewSUMService(r.repositories.atomicStatusRepository)

	// Initialize teleport parser (listens to teleport logs for cross chain messages)
	teleportParser := logparser.NewTeleportParser(
		r.msgqueues.teleportBlockConsumer,
		r.msgqueues.crossChainPublisher,
		r.ctsClient,
		r.hubClient, // Provide hub client as EthereumClient
		conf.PrivateNodeChainID,
		sumService,
	)

	// Initialize audit manager parser (listens to audit manager logs for key agreement)
	auditManagerParser := logparser.NewAuditManagerParser(
		conf.PrivateNodeChainID,
		r.msgqueues.auditManagerBlockConsumer,
		r.ctsClient,
	)

	r.logParsers = &LogParsers{
		endpointParser:       endpointParser,
		teleportParser:       teleportParser,
		enygmaTeleportParser: enygmaTeleportParser,
		dvpTeleportParser:    dvpTeleportParser,
		auditManagerParser:   auditManagerParser,
	}

	return nil
}
