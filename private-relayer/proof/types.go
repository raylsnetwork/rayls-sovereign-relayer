package proof

import (
	"math/big"
	"time"
)

// HeaderProofConfig holds configuration for HeaderProofService
type HeaderProofConfig struct {
	PLChainID    *big.Int
	PollInterval time.Duration
	BatchSize    int64
	Timeout      time.Duration
}

const (
	defaultProofPollInterval = 20 * time.Second
	defaultProofBatchSize    = 50
	defaultProofTimeout      = 2 * time.Minute
)

// DefaultHeaderProofConfig returns config with sensible defaults
func DefaultHeaderProofConfig(plChainID *big.Int) HeaderProofConfig {
	return HeaderProofConfig{
		PLChainID:    plChainID,
		PollInterval: defaultProofPollInterval,
		BatchSize:    defaultProofBatchSize,
		Timeout:      defaultProofTimeout,
	}
}
