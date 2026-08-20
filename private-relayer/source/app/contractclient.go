package app

import (
	"fmt"
	"log/slog"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/raylsnetwork/rayls-sovereign-relayer/contractclient"
	"github.com/raylsnetwork/rayls-sovereign-relayer/txsim"
)

type ContractClientConfig struct {
	PrivateHubChainID  *big.Int
	PrivateNodeChainID *big.Int
	VENChainID         *big.Int

	// Registry addresses for contract discovery
	PrivateHubDeploymentProxyRegistry  common.Address
	PrivateNodeDeploymentProxyRegistry common.Address
}

type ContractClients struct {
	teleport           *contractclient.TeleportClient
	participantStorage *contractclient.ParticipantStorageClient

	// Endpoint clients for both chains
	privateHubEndpoint  *contractclient.EndpointClient
	privateNodeEndpoint *contractclient.EndpointClient

	// Dvp clients
	dvpClient *contractclient.DvpClient

	// Enygma and DVP contract clients
	enygmaIntegrationClient *contractclient.DvpIntegrationClient
	enygmaClient            *contractclient.EnygmaClient
	enygmaHandlerClient     *contractclient.EnygmaHandlerClient
	erc721Client            *contractclient.DvpERC721Client
	erc721HandlerClient     *contractclient.DvpERC721HandlerClient
	erc1155Client           *contractclient.DvpERC1155Client
	erc1155HandlerClient    *contractclient.DvpERC1155HandlerClient

	hubRegistry  *contractclient.DeploymentProxyRegistryClient
	nodeRegistry *contractclient.DeploymentProxyRegistryClient

	// Executors implement Executor. hub and node route through CTS TxOps
	// (no local key queue, no local authGen). operator still uses the
	// local contractclient.Executor because the DvP operator signer
	// identity hasn't been migrated into CTS yet — it has its own
	// single-key queue with narrow scope.
	hubExecutor      contractclient.Executor
	hubSigner        contractclient.Signer
	nodeExecutor     contractclient.Executor
	operatorExecutor contractclient.Executor
}

func (p *SourcePrivateRelayer) initializeContractClients(
	config ContractClientConfig,
) error {
	hubExecutor := contractclient.NewCTSExecutor(contractclient.NewDefaultRetryingTxOpsClient(p.ctsClient.PrivateHubTxOpsServiceClient))

	hubRegistry, err := contractclient.NewDeploymentProxyRegistryClient(config.PrivateHubDeploymentProxyRegistry, p.hubClient)
	if err != nil {
		return fmt.Errorf("failed to create private hub registry: %w", err)
	}

	nodeExecutor := contractclient.NewCTSExecutor(contractclient.NewDefaultRetryingTxOpsClient(p.ctsClient.PrivateNodeTxOpsServiceClient))
	operatorExecutor := contractclient.NewCTSExecutor(contractclient.NewDefaultRetryingTxOpsClient(p.ctsClient.DVPOperatorTxOpsServiceClient))

	nodeRegistry, err := contractclient.NewDeploymentProxyRegistryClient(config.PrivateNodeDeploymentProxyRegistry, p.nodeClient)
	if err != nil {
		return fmt.Errorf("failed to create private node registry client factory: %w", err)
	}

	// Populate error map for transaction simulation
	err = txsim.PopulateErrorMap("./contracts/")
	if err != nil {
		return fmt.Errorf("failed to populate error map: %w", err)
	}
	// Create participant storage first (needed by encryptor)

	participantStorageAddress, err := hubRegistry.GetContractAddress("ParticipantStorage")
	if err != nil {
		return fmt.Errorf("failed to get ParticipantStorage contract address: %w", err)
	}

	participantStorageClient := contractclient.NewParticipantStorageClient(
		participantStorageAddress,
		config.PrivateNodeChainID,
		config.VENChainID,
		hubExecutor,
	)
	if err != nil {
		return fmt.Errorf("failed to create participant storage client: %w", err)
	}

	// Create encryptor (needs participant storage)
	encryptor := contractclient.NewEncryptor(
		p.ctsClient,
		participantStorageClient,
		p.hubClient,
		config.PrivateHubChainID,
	)

	// Create private hub teleport client
	teleportAddress, err := hubRegistry.GetContractAddress("Teleport")
	if err != nil {
		return fmt.Errorf("failed to get Teleport contract address: %w", err)
	}
	teleportClient := contractclient.NewTeleportClient(
		teleportAddress,
		hubExecutor,
		encryptor,
	)

	privateHubEndpointAddress, err := hubRegistry.GetContractAddress("Endpoint")
	if err != nil {
		return fmt.Errorf("failed to get private hub endpoint address: %w", err)
	}
	privateHubEndpoint := contractclient.NewEndpointClient(
		privateHubEndpointAddress,
		hubExecutor,
	)

	// Create Dvp client
	dvpAddress, err := hubRegistry.GetContractAddress("Dvp")
	if err != nil {
		return fmt.Errorf("failed to get Dvp contract address: %w", err)
	}
	dvpClient := contractclient.NewDvpClient(
		dvpAddress,
		hubExecutor,
		operatorExecutor,
		encryptor,
	)

	// Create private node endpoint client
	privateNodeEndpointAddress, err := nodeRegistry.GetContractAddress("Endpoint")
	if err != nil {
		return fmt.Errorf("failed to get private node endpoint address: %w", err)
	}
	privateNodeEndpoint := contractclient.NewEndpointClient(
		privateNodeEndpointAddress,
		nodeExecutor,
	)

	// Enygma and DVP contract clients
	enygmaIntegrationClient := contractclient.NewDvpIntegrationClient(hubExecutor, encryptor, p.hubClient)

	// The source relayer's EnygmaHandlerClient never calls ReceiveDestTransferBatch — the only
	// method that uses programmabilityExecutorAddress. The source side only does ReceiveWithdraw,
	// RevertSrc*, and DVP notifications. So a missing ProgrammabilityExecutor registry entry must
	// not block startup here; we fall back to the zero address with a warning. The dest relayer,
	// which dispatches every inbound transfer through the executor, keeps its hard-fail.
	programmabilityExecutorAddress, err := nodeRegistry.GetContractAddress("ProgrammabilityExecutor")
	if err != nil {
		slog.Warn("ProgrammabilityExecutor not found in registry; source relayer continues with zero address (it never dispatches through the executor)",
			slog.Any("error", err))
		programmabilityExecutorAddress = common.Address{}
	}

	enygmaClient := contractclient.NewEnygmaClient(hubExecutor, encryptor)
	enygmaHandlerClient := contractclient.NewEnygmaHandlerClient(nodeExecutor, programmabilityExecutorAddress)

	erc721Client := contractclient.NewDvpERC721Client(operatorExecutor, encryptor)
	erc721HandlerClient := contractclient.NewDvpERC721HandlerClient(nodeExecutor)

	erc1155Client := contractclient.NewDvpERC1155Client(operatorExecutor, encryptor)
	erc1155HandlerClient := contractclient.NewDvpERC1155HandlerClient(nodeExecutor)

	p.contractClients = &ContractClients{
		teleport:            teleportClient,
		participantStorage:  participantStorageClient,
		privateHubEndpoint:  privateHubEndpoint,
		privateNodeEndpoint: privateNodeEndpoint,
		dvpClient:           dvpClient,

		enygmaIntegrationClient: enygmaIntegrationClient,
		enygmaClient:            enygmaClient,
		enygmaHandlerClient:     enygmaHandlerClient,
		erc721Client:            erc721Client,
		erc721HandlerClient:     erc721HandlerClient,
		erc1155Client:           erc1155Client,
		erc1155HandlerClient:    erc1155HandlerClient,

		hubRegistry:  hubRegistry,
		nodeRegistry: nodeRegistry,

		hubExecutor:      hubExecutor,
		nodeExecutor:     nodeExecutor,
		operatorExecutor: operatorExecutor,
	}
	return nil
}
