package app

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/raylsnetwork/rayls-sovereign-relayer/contractclient"
	keyspb "github.com/raylsnetwork/rayls-sovereign-relayer/cts/gen/proto/keys"
	dvpHandler "github.com/raylsnetwork/rayls-sovereign-relayer/dvp/handler"
	dvpService "github.com/raylsnetwork/rayls-sovereign-relayer/dvp/service"
	"github.com/raylsnetwork/rayls-sovereign-relayer/private-relayer/source/service"
)

const (
	dvpDepositWaiterMaxRetries    = 30
	dvpDepositWaiterRetryInterval = 3 * time.Second
	dvpSwapWaiterRetryInterval    = 6 * time.Second
)

// DvpServiceConfig holds configuration for dvp services
type DvpServiceConfig struct {
	PrivateHubChainID            *big.Int
	MyChainID                    *big.Int
	NumberOfJSParamsIn           int
	DefaultContextTimeout        time.Duration
	PrivateHubBlockTimeInSeconds float64
}

// initializeDvpServices initializes all dvp-related services
func (r *SourcePrivateRelayer) initializeDvpServices(
	ctx context.Context,
	config DvpServiceConfig,
) (*service.DvpOrchestrator, error) {
	// Use pre-created clients
	dvpClient := r.contractClients.dvpClient
	privateHubEndpointClient := r.contractClients.privateHubEndpoint
	privateNodeEndpointClient := r.contractClients.privateNodeEndpoint

	// Get dvp contract address from hub registry (needed for initiator config)
	dvpContractAddress, err := r.contractClients.hubRegistry.GetContractAddress("Dvp")
	if err != nil {
		return nil, fmt.Errorf("getting dvp contract address: %w", err)
	}

	// Create encryptor (needed for contract client factory)
	encryptor := contractclient.NewEncryptor(
		r.ctsClient,
		r.contractClients.participantStorage,
		r.hubClient,
		config.PrivateHubChainID,
	)

	// Initialize dvp services (use shared commitment calculator from infrastructure)
	dvpDepositFinder := dvpService.NewDepositFinder(r.repositories.dvpDeposit)

	// Initialize dvp proof service (reuse from infrastructure)
	dvpProofService := dvpService.NewProofService(
		dvpService.ProofServiceConfig{
			ChainID:            config.MyChainID,
			MerkleTreeDepth:    r.infrastructure.dvpMerkleTreeDepth,
			NumberOfJSParamsIn: config.NumberOfJSParamsIn,
		},
		r.infrastructure.merkleService,
		r.ctsClient,
		r.infrastructure.proofAPIClient,
		r.infrastructure.commitmentCalculator,
		r.repositories.dvpDeposit,
		r.infrastructure.transactionManager,
	)

	// Initialize deposit waiter
	depositWaiter := dvpService.NewDepositWaiter(
		dvpService.WaitConfig{
			MaxRetries:    dvpDepositWaiterMaxRetries,
			RetryInterval: dvpDepositWaiterRetryInterval,
		},
		r.repositories.dvpDeposit,
	)

	// Initialize swap waiter
	swapWaiter := dvpService.NewSwapWaiter(
		dvpService.WaitConfig{
			MaxRetries:    10,
			RetryInterval: dvpSwapWaiterRetryInterval,
		},
		r.repositories.dvpSwap,
	)

	// Initialize dvp consolidation service
	dvpConsolidationService := dvpService.NewConsolidationService(
		dvpService.ConsolidationConfig{
			ChainID:               config.MyChainID,
			MaxNumberOfJSDeposits: config.NumberOfJSParamsIn,
		},
		r.repositories.dvpDeposit,
		r.ctsClient,
		r.infrastructure.commitmentCalculator,
		dvpProofService,
		dvpClient,
		r.contractClients.enygmaClient,
		r.contractClients.enygmaIntegrationClient,
		depositWaiter,
		r.infrastructure.transactionManager,
	)

	dvpOperatorAddress, err := getDvpOperatorAddress(r.ctsClient)
	if err != nil {
		return nil, fmt.Errorf("getting dvp operator address: %w", err)
	}

	// Initialzie swap agreement
	swapAgreement := dvpService.NewSwapAgreement(
		r.contractClients.privateNodeEndpoint,
		r.contractClients.enygmaHandlerClient,
		r.contractClients.erc721HandlerClient,
		r.contractClients.erc1155HandlerClient,
	)

	// Initialize dvp initiator
	dvpInitiator := dvpHandler.NewDvpInitiator(
		dvpHandler.InitiatorConfig{
			ChainID:            config.MyChainID,
			DvpContractAddress: dvpContractAddress,
			DvpOperatorAddress: dvpOperatorAddress,
		},
		r.repositories.dvpSwap,
		r.repositories.dvpDeposit,
		r.contractClients.participantStorage,
		r.ctsClient,
		privateHubEndpointClient,
		privateNodeEndpointClient,

		r.contractClients.erc721Client,
		r.contractClients.erc721HandlerClient,
		r.contractClients.erc1155Client,
		r.contractClients.erc1155HandlerClient,

		dvpClient,
		encryptor,
		dvpDepositFinder,
		dvpConsolidationService,
		r.infrastructure.commitmentCalculator,
		dvpProofService,
		swapWaiter,
		swapAgreement,
		r.nodeClient,
		r.hubClient,
		r.infrastructure.transactionManager,
	)

	// Initialize dvp orchestrator
	dvpOrchestrator := service.NewDvpOrchestratorWithDeps(
		r.msgqueues.dvpBatchConsumer,
		dvpInitiator,
	)

	return dvpOrchestrator, nil
}

func getDvpOperatorAddress(ctsClient keyspb.KeysServiceClient) (common.Address, error) {
	resp, err := ctsClient.GetSignAddresses(context.Background(), &keyspb.GetSignAddressesRequest{
		ServiceType: keyspb.ServiceType_SERVICE_TYPE_PRIVATE_RELAYER,
	})
	if err != nil {
		return common.Address{}, fmt.Errorf("calling GetSignAddresses: %w", err)
	}

	dvpAddresses := resp.GetAddressSets()[int32(keyspb.KeySetName_KEY_SET_NAME_PRIVATE_HUB_DVP_OPERATOR)]
	if dvpAddresses == nil || len(dvpAddresses.GetValues()) == 0 {
		return common.Address{}, fmt.Errorf("no DVP operator address returned from CTS")
	}

	return common.HexToAddress(dvpAddresses.GetValues()[0]), nil
}
