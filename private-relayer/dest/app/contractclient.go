package app

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/contractclient"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/txbatchclient"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/txsim"
)

type ContractClientConfig struct {
	PrivateHubChainID  *big.Int
	PrivateNodeChainID *big.Int
	VENChainID         *big.Int

	PrivateHubDeploymentProxyRegistry  common.Address
	PrivateNodeDeploymentProxyRegistry common.Address
}

type ContractClients struct {
	// Private Hub clients (for listening to source messages)
	hubParticipantStorage *contractclient.ParticipantStorageClient
	hubEncryptor          *contractclient.Encryptor
	hubReceipter          *txbatchclient.TxReceipter
	hubTxSimulator        *txsim.TransactionSimulator
	hubRegistry           *contractclient.DeploymentProxyRegistryClient

	// Private Node clients (for executing on destination)
	nodeEndpoint    *contractclient.EndpointClient
	nodeAuthGen     *contractclient.AuthGen
	nodeDeployer    *contractclient.DeployerClient
	nodeReceipter   *txbatchclient.TxReceipter
	nodeTxSimulator *txsim.TransactionSimulator
	nodeBatchSender *txbatchclient.TxSender
	nodeRegistry    *contractclient.DeploymentProxyRegistryClient

	// Enygma and DVP contract clients
	enygmaIntegrationClient *contractclient.DvpIntegrationClient
	enygmaClient            *contractclient.EnygmaClient
	enygmaHandlerClient     *contractclient.EnygmaHandlerClient
	erc721Client            *contractclient.DvpERC721Client
	erc721HandlerClient     *contractclient.DvpERC721HandlerClient
	erc1155Client           *contractclient.DvpERC1155Client
	erc1155HandlerClient    *contractclient.DvpERC1155HandlerClient

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

func (r *DestPrivateRelayer) initializeContractClients(
	conf ContractClientConfig,
) error { //nolint:contextcheck // initialization code, no parent context available

	// Create auth generators
	//nolint:contextcheck // initialization code, no parent context
	hubSimulator := txsim.NewTransactionSimulator(r.hubClient)

	hubExecutor := contractclient.NewCTSExecutor(contractclient.NewDefaultRetryingTxOpsClient(r.ctsClient.PrivateHubTxOpsServiceClient))

	hubRegistry, err := contractclient.NewDeploymentProxyRegistryClient(conf.PrivateHubDeploymentProxyRegistry, r.hubClient)
	if err != nil {
		return fmt.Errorf("failed to create private hub registry: %w", err)
	}

	nodeSimulator := txsim.NewTransactionSimulator(r.nodeClient)
	nodeAuthGen, err := contractclient.NewAuthGen(context.Background(), r.nodeClient)
	if err != nil {
		return fmt.Errorf("failed to create node auth generator: %w", err)
	}

	nodeExecutor := contractclient.NewCTSExecutor(contractclient.NewDefaultRetryingTxOpsClient(r.ctsClient.PrivateNodeTxOpsServiceClient))
	operatorExecutor := contractclient.NewCTSExecutor(contractclient.NewDefaultRetryingTxOpsClient(r.ctsClient.DVPOperatorTxOpsServiceClient))

	nodeRegistry, err := contractclient.NewDeploymentProxyRegistryClient(conf.PrivateNodeDeploymentProxyRegistry, r.nodeClient)
	if err != nil {
		return fmt.Errorf("failed to create private node registry client factory: %w", err)
	}

	// Populate error map for transaction simulation
	err = txsim.PopulateErrorMap("./contracts/")
	if err != nil {
		return fmt.Errorf("failed to populate error map: %w", err)
	}

	// Create participant storage client (needed for encryptor)
	//nolint:contextcheck // initialization code, no parent context available
	hubParticipantStorageAddress, err := hubRegistry.GetContractAddress("ParticipantStorage")
	if err != nil {
		return fmt.Errorf("failed to get hub participant storage address: %w", err)
	}
	hubParticipantStorage := contractclient.NewParticipantStorageClient(
		hubParticipantStorageAddress,
		conf.PrivateNodeChainID,
		conf.VENChainID,
		hubExecutor,
	)

	// Create encryptor (for decrypting messages from hub)
	hubEncryptor := contractclient.NewEncryptor(
		r.ctsClient,
		hubParticipantStorage,
		r.hubClient,
		conf.PrivateHubChainID,
	)

	// Create private node endpoint client (for executing on destination)
	//nolint:contextcheck // initialization code, no parent context available
	nodeEndpointAddress, err := nodeRegistry.GetContractAddress("Endpoint")
	if err != nil {
		return fmt.Errorf("failed to get node endpoint address: %w", err)
	}
	nodeEndpoint := contractclient.NewEndpointClient(nodeEndpointAddress, nodeExecutor)

	// Get ResourceRegistry address for nodeDeployer
	resourceRegistryAddress, err := hubRegistry.GetContractAddress("ResourceRegistry")
	if err != nil {
		return fmt.Errorf("failed to get ResourceRegistry address: %w", err)
	}
	resourceRegistryClient := contractclient.NewResourceRegistryClient(resourceRegistryAddress, r.hubClient)

	nodeDeployer := contractclient.NewDeployerClient(
		nodeEndpoint,
		resourceRegistryClient,
		60*time.Second,
	)

	hubReceipter := txbatchclient.NewTxReceipter(r.hubClient.Client())
	nodeReceipter := txbatchclient.NewTxReceipter(r.nodeClient.Client())

	// Create batch sender for node
	nodeBatchSender := txbatchclient.NewTxSender(r.nodeClient.Client(), nodeReceipter)

	// Enygma and DVP contract clients
	enygmaIntegrationClient := contractclient.NewDvpIntegrationClient(hubExecutor, hubEncryptor, r.hubClient)

	// Resolve the per-PN ProgrammabilityExecutor — every inbound cross-chain
	// transfer is dispatched through it as an executeProgramData tx.
	programmabilityExecutorAddress, err := nodeRegistry.GetContractAddress("ProgrammabilityExecutor")
	if err != nil {
		return fmt.Errorf("failed to get ProgrammabilityExecutor address: %w", err)
	}

	enygmaClient := contractclient.NewEnygmaClient(hubExecutor, hubEncryptor)
	enygmaHandlerClient := contractclient.NewEnygmaHandlerClient(nodeExecutor, programmabilityExecutorAddress)

	erc721Client := contractclient.NewDvpERC721Client(hubExecutor, hubEncryptor)
	erc721HandlerClient := contractclient.NewDvpERC721HandlerClient(nodeExecutor)

	erc1155Client := contractclient.NewDvpERC1155Client(hubExecutor, hubEncryptor)
	erc1155HandlerClient := contractclient.NewDvpERC1155HandlerClient(nodeExecutor)

	r.contractClients = &ContractClients{
		hubParticipantStorage: hubParticipantStorage,
		hubEncryptor:          hubEncryptor,
		hubReceipter:          hubReceipter,
		hubTxSimulator:        hubSimulator,
		hubRegistry:           hubRegistry,

		nodeEndpoint:    nodeEndpoint,
		nodeAuthGen:     nodeAuthGen,
		nodeDeployer:    nodeDeployer,
		nodeReceipter:   nodeReceipter,
		nodeTxSimulator: nodeSimulator,
		nodeBatchSender: nodeBatchSender,
		nodeRegistry:    nodeRegistry,

		enygmaIntegrationClient: enygmaIntegrationClient,
		enygmaClient:            enygmaClient,
		enygmaHandlerClient:     enygmaHandlerClient,
		erc721Client:            erc721Client,
		erc721HandlerClient:     erc721HandlerClient,
		erc1155Client:           erc1155Client,
		erc1155HandlerClient:    erc1155HandlerClient,

		hubExecutor:      hubExecutor,
		nodeExecutor:     nodeExecutor,
		operatorExecutor: operatorExecutor,
	}

	return nil
}
