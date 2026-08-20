package app

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/raylsnetwork/rayls-sovereign-relayer/contractclient"
	dvpHandler "github.com/raylsnetwork/rayls-sovereign-relayer/dvp/handler"
	dvpService "github.com/raylsnetwork/rayls-sovereign-relayer/dvp/service"

	enygmaAdapters "github.com/raylsnetwork/rayls-sovereign-relayer/enygma/adapters"
	enygmaHandler "github.com/raylsnetwork/rayls-sovereign-relayer/enygma/handler"
	enygmaService "github.com/raylsnetwork/rayls-sovereign-relayer/enygma/service"
	merkle "github.com/raylsnetwork/rayls-sovereign-relayer/merkle-service"
	"github.com/raylsnetwork/rayls-sovereign-relayer/private-relayer/repository"
	"github.com/raylsnetwork/rayls-sovereign-relayer/types"
)

// depositWaiterRetryInterval is the interval between retries when waiting for a DVP deposit.
const depositWaiterRetryInterval = 3 * time.Second

// depositWaiterMaxRetries is the maximum number of retries when waiting for a DVP deposit.
const depositWaiterMaxRetries = 30

// dvpMerkleAdapter composes MerkleService and DvpDepositRepository to satisfy
// the receiverMerkleClient interface which requires both PopulateMerkleDbTree
// and GetNonFungibleDeposit.
type dvpMerkleAdapter struct {
	*merkle.MerkleService
	depositRepo *repository.DvpDepositRepository
}

func (a *dvpMerkleAdapter) GetNonFungibleDeposit(
	ctx context.Context,
	tokenID string,
	tokenAddress string,
	userAddress string,
	tokenType types.DvpTokenType,
	status types.DvpDepositStatus,
) (*types.DvpDeposit, error) {
	return a.depositRepo.GetNonFungibleDeposit(ctx, tokenID, tokenAddress, userAddress, tokenType, status)
}

type ReceiverConfig struct {
	DefaultContextTimeout        time.Duration
	NumberOfJSParamsIn           int
	PrivateNodeChainID           *big.Int
	PrivateHubChainID            *big.Int
	PrivateHubBlockTimeInSeconds float64
}

type Receivers struct {
	enygmaReceiver *enygmaHandler.Receiver
	enygmaDeployer *enygmaService.EnygmaDeployer
	dvpReceiver    *dvpHandler.DvpReceiver
}

func (r *DestPrivateRelayer) initializeReceivers(
	conf ReceiverConfig,
) error { //nolint:contextcheck // initialization code, no parent context available
	// Initialize OTel tracer for receivers
	enygmaTracer := enygmaAdapters.NewOTelTracer("dest-enygma-receiver")
	dvpTracer := enygmaAdapters.NewOTelTracer("dest-dvp-receiver")

	// --- Initialize Enygma Dependencies ---

	// 1. Create Resource Registry Client
	// Get address from hub deployment proxy registry
	resourceRegistryAddress, err := r.contractClients.hubRegistry.GetContractAddress("ResourceRegistry")
	if err != nil {
		return fmt.Errorf("failed to get ResourceRegistry address: %w", err)
	}

	// Create client using abigen v2-style pack/unpack + CallContract
	resourceRegistryClient := contractclient.NewResourceRegistryClient(resourceRegistryAddress, r.hubClient)

	// 2. Initialize Enygma Deployer
	enygmaDeployer := enygmaService.NewEnygmaDeployer(
		r.contractClients.nodeEndpoint,
		resourceRegistryClient,
		enygmaTracer,
	)

	// Initialize Enygma Creation Service
	enygmaCreationService := enygmaService.NewEnygmaCreationService(
		enygmaAdapters.NewOTelTracer("enygma-creation"),
		r.infrastructure.transactionManager,
		r.repositories.enygma,
		r.repositories.enygmaHistory,
	)

	// Create Enygma Teleport Client (for Hub/CC)
	teleportAddress, err := r.contractClients.hubRegistry.GetContractAddress("EnygmaTeleport")
	if err != nil {
		return fmt.Errorf("failed to get EnygmaTeleport contract address: %w", err)
	}

	enygmaTeleportClient := contractclient.NewEnygmaTeleportClient(
		teleportAddress,
		r.contractClients.hubExecutor,
		r.contractClients.hubEncryptor,
	)

	// Initialize Enygma Receiver
	enygmaReceiver := enygmaHandler.NewReceiver(
		conf.PrivateNodeChainID,
		r.contractClients.nodeEndpoint,
		r.contractClients.nodeTxSimulator,
		enygmaDeployer,
		r.repositories.enygmaHistory,
		r.repositories.txRecoveryData,
		enygmaTeleportClient,
		r.contractClients.enygmaHandlerClient,
		enygmaTracer,
		enygmaCreationService,
		conf.DefaultContextTimeout,
	)

	// Initialize Dvp Proof Service
	dvpProofService := dvpService.NewProofService(
		dvpService.ProofServiceConfig{
			ChainID:            conf.PrivateNodeChainID,
			MerkleTreeDepth:    r.infrastructure.dvpMerkleTreeDepth,
			NumberOfJSParamsIn: conf.NumberOfJSParamsIn,
		},
		r.infrastructure.merkleService,
		r.ctsClient,
		r.infrastructure.proofAPIClient,
		r.infrastructure.commitmentCalculator,
		r.repositories.dvpDeposit,
		r.infrastructure.transactionManager,
	)

	// Initialize Dvp Deposit Finder
	dvpDepositFinder := dvpService.NewDepositFinder(r.repositories.dvpDeposit)

	// Initialize Dvp Deposit Waiter (needed by ConsolidationService)
	depositWaiter := dvpService.NewDepositWaiter(
		dvpService.WaitConfig{
			MaxRetries:    depositWaiterMaxRetries,
			RetryInterval: depositWaiterRetryInterval,
		},
		r.repositories.dvpDeposit,
	)

	// Create Dvp Client
	dvpAddress, err := r.contractClients.hubRegistry.GetContractAddress("Dvp")
	if err != nil {
		return fmt.Errorf("failed to get Dvp contract address: %w", err)
	}
	dvpClient := contractclient.NewDvpClient(
		dvpAddress,
		r.contractClients.hubExecutor,
		r.contractClients.operatorExecutor,
		r.contractClients.hubEncryptor,
	)

	// Initialize Dvp Consolidation Service (provides PrepareDepositsForJSProof)
	dvpConsolidationService := dvpService.NewConsolidationService(
		dvpService.ConsolidationConfig{
			ChainID:               conf.PrivateNodeChainID,
			MaxNumberOfJSDeposits: conf.NumberOfJSParamsIn,
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

	// Initialize Dvp Retry Service
	dvpRetryService := dvpService.NewRetryService(dvpTracer)

	// Initialize Dvp Receiver
	dvpReceiver := dvpHandler.NewDvpReceiver(
		dvpHandler.ReceiverConfig{
			ChainID: conf.PrivateNodeChainID,
		},
		r.repositories.dvpSwap,
		r.repositories.dvpDeposit,
		r.contractClients.nodeEndpoint,
		r.contractClients.enygmaHandlerClient,
		r.contractClients.erc721HandlerClient,
		r.contractClients.erc1155HandlerClient,
		dvpDepositFinder,
		dvpConsolidationService, // ConsolidationService implements receiverDepositConsolidator
		r.infrastructure.commitmentCalculator,
		dvpProofService,
		&dvpMerkleAdapter{
			MerkleService: r.infrastructure.merkleService,
			depositRepo:   r.repositories.dvpDeposit,
		},
		dvpRetryService,
		r.infrastructure.transactionManager,
	)

	r.receivers = &Receivers{
		enygmaReceiver: enygmaReceiver,
		enygmaDeployer: enygmaDeployer,
		dvpReceiver:    dvpReceiver,
	}

	return nil
}
