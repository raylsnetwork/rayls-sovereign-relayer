package app

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/listener"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/private-relayer/source/logrouter"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/types"
)

type ListenersConfig struct {
	BatchSize int

	PrivateNodeStartingBlock *big.Int
}

type Listeners struct {
	privateNode *listener.LogListener
}

func (n *SourcePrivateRelayer) initializeListeners(conf ListenersConfig) error {
	endpointAddress, err := n.contractClients.nodeRegistry.GetContractAddress("Endpoint")
	if err != nil {
		return fmt.Errorf("failed to get private node endpoint address: %w", err)
	}

	enygmaEventsAddress, err := n.contractClients.nodeRegistry.GetContractAddress("EnygmaPNEvents")
	if err != nil {
		return fmt.Errorf("failed to get private node enygma events address: %w", err)
	}

	privateNodeListenerConfig := listener.LogListenerConfig{
		Component: types.DocumentIdLastProcessedBlockPrivacyNode,

		StartingBlock: conf.PrivateNodeStartingBlock,
		BatchSize:     conf.BatchSize,

		Addresses: []common.Address{
			endpointAddress,
			enygmaEventsAddress,
		},
	}
	logRouter := logrouter.New(
		endpointAddress,
		enygmaEventsAddress,
		n.msgqueues.endpointBlockPublisher,
		n.msgqueues.enygmaBlockPublisher,
	)
	privateNodeListener, err := listener.NewLogListener(
		context.Background(),
		privateNodeListenerConfig,
		logRouter,
		n.nodeClient,
		n.repositories.lastProcessedBlockRepository,
	)
	if err != nil {
		return fmt.Errorf("failed to create private node listener: %w", err)
	}

	n.listeners = &Listeners{
		privateNode: privateNodeListener,
	}

	return nil
}
