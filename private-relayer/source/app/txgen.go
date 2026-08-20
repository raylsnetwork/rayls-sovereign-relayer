package app

import (
	"fmt"

	"github.com/raylsnetwork/rayls-sovereign-relayer/private-relayer/source/txgen"
)

type TransactionGenerator struct {
	privateHub *txgen.PrivateHubGenerator
}

func (r *SourcePrivateRelayer) initializeTransactionGenerator() error {
	privateHubGen, err := txgen.NewPrivateHubGenerator(r.hubClient)
	if err != nil {
		return fmt.Errorf("failed to create private hub transaction generator: %w", err)
	}

	r.txGen = &TransactionGenerator{
		privateHub: privateHubGen,
	}

	return nil
}
