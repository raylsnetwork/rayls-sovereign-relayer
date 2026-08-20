package app

import (
	"fmt"

	"github.com/raylsnetwork/rayls-privacy-relayer-api/private-relayer/dest/txgen"
)

type TransactionGenerator struct {
	privateNode *txgen.PrivateNodeGenerator
}

func (r *DestPrivateRelayer) initializeTransactionGenerator() error {
	privateNodeGen, err := txgen.NewPrivateNodeGenerator(r.nodeClient)
	if err != nil {
		return fmt.Errorf("failed to create private node transaction generator: %w", err)
	}

	r.txGen = &TransactionGenerator{
		privateNode: privateNodeGen,
	}

	return nil
}
