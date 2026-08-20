package app

import (
	"math/big"
	"time"

	enygmaAdapters "github.com/raylsnetwork/rayls-sovereign-relayer/enygma/adapters"
	enygmaService "github.com/raylsnetwork/rayls-sovereign-relayer/enygma/service"
)

// defaultMaxTxsPerBatch is the default maximum number of transactions per enygma batch.
const defaultMaxTxsPerBatch = 50

type BatcherConfig struct {
	DefaultContextTimeout time.Duration
	PrivateNodeChainID    *big.Int
}

type Batchers struct {
	enygmaBatcher *enygmaService.EnygmaBatcher
}

func (r *DestPrivateRelayer) initializeBatchers(conf BatcherConfig) {
	// Initialize OTel tracer for batcher
	tracer := enygmaAdapters.NewOTelTracer("dest-enygma-batcher")

	// Initialize Enygma Batcher
	enygmaBatcherConfig := &enygmaService.EnygmaBatcherConfig{
		ChainID:        conf.PrivateNodeChainID,
		MaxTxsPerBatch: defaultMaxTxsPerBatch, // Default batch size, could be configurable
	}

	enygmaBatcher := enygmaService.NewEnygmaBatcher(
		enygmaBatcherConfig,
		r.contractClients.hubParticipantStorage,
		tracer,
	)

	r.batchers = &Batchers{
		enygmaBatcher: enygmaBatcher,
	}
}
