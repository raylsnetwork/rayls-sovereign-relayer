package app

import (
	"math/big"

	"github.com/raylsnetwork/rayls-sovereign-relayer/private-relayer/source/logparser"
)

type LogpParsersConfig struct {
	PrivateHubChainID *big.Int
}
type LogParsers struct {
	endpointParser *logparser.EndpointLogParser
	enygmaParser   *logparser.EnygmaLogParser
}

func (r *SourcePrivateRelayer) initializeLogParsers(config LogpParsersConfig) error {
	endpointParser := logparser.NewEndpointLogParser(
		config.PrivateHubChainID,
		r.msgqueues.endpointBlockConsumer,
		r.msgqueues.crossChainPublisher,
		r.msgqueues.privateHubPublisher,
		r.services.keysService,
	)

	enygmaParser := logparser.NewEnygmaLogParser(
		r.msgqueues.enygmaBlockConsumer,
		r.msgqueues.enygmaBatchPublisher,
		r.msgqueues.dvpBatchPublisher,
	)

	r.logParsers = &LogParsers{
		endpointParser: endpointParser,
		enygmaParser:   enygmaParser,
	}

	return nil
}
