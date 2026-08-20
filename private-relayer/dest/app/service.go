// Decommissioning Teleport (vanilla, atomic): atomic members below marked; shared/generic/Enygma/DVP retained.

package app

import (
	"fmt"
	"math/big"
	"time"

	"github.com/raylsnetwork/rayls-privacy-relayer-api/batcher"
	enygmaAdapters "github.com/raylsnetwork/rayls-privacy-relayer-api/enygma/adapters"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/enygma/resync"
	enygmaService "github.com/raylsnetwork/rayls-privacy-relayer-api/enygma/service"
	destservice "github.com/raylsnetwork/rayls-privacy-relayer-api/private-relayer/dest/service"
	sharedservice "github.com/raylsnetwork/rayls-privacy-relayer-api/private-relayer/shared/service"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/resultrouter"
)

// MessageType strings for the four publisher flavours that target the
// privatenode signing identity. They double as the keys the result
// router uses to fan TxResult batches into per-callback handlers, so
// they MUST match what the producer publishes (here) and what the
// router expects (router registration below).
const (
	// Decommissioning Teleport (vanilla, atomic): msgTypeCrosschainAtomic only; msgTypeCrosschainVanilla is retained (generic non-atomic).
	msgTypeCrosschainAtomic  = "crosschain.atomic"
	msgTypeCrosschainVanilla = "crosschain.vanilla"
	// Decommissioning Teleport (vanilla, atomic).
	msgTypeAtomicDestinationUnlock = "atomic.destination-unlock"
	msgTypeAtomicDestinationRevert = "atomic.destination-revert"
	msgTypeAtomicSourceRevert      = "atomic.source-revert"
	msgTypePrivateHubExecute       = "privatehub.execute"
)

const (
	routerBatchSize = 100
	routerInterval  = time.Second
)

type ServiceConfig struct {
	MyChainID              *big.Int
	PrivateHubChainID      *big.Int
	DefaultContextTimeout  time.Duration
	PrivateHubTickerPeriod time.Duration
	EnygmaSyncMaxRetries   int
}

type Services struct {
	privateHub         *sharedservice.PrivateHubService
	crossChain         *destservice.CrossChainService
	enygmaOrchestrator *destservice.EnygmaOrchestrator
	enygmaSyncService  *enygmaService.EnygmaSyncService
	dvpOrchestrator    *destservice.DvpOrchestrator
}

func (r *DestPrivateRelayer) initializeServices(conf ServiceConfig) error {
	// 2. PrivateHub Service (consumes endpoint-parsed messages, generates txs, pushes to executor)
	// Get destination endpoint address (Private Node endpoint where we execute)
	destEndpointAddress, err := r.contractClients.nodeRegistry.GetContractAddress("Endpoint")
	if err != nil {
		return fmt.Errorf("failed to get destination endpoint address: %w", err)
	}

	// All four batchers share the same publisher (cts.send.privatenode);
	// they only differ by messageType, which is what the router uses to
	// fan results back to the right callback.
	privateHubBatcher := batcher.NewBatcher(msgTypePrivateHubExecute, r.msgqueues.privateNodeSendPublisher)
	atomicBatcher := batcher.NewBatcher(msgTypeCrosschainAtomic, r.msgqueues.privateNodeSendPublisher)
	vanillaBatcher := batcher.NewBatcher(msgTypeCrosschainVanilla, r.msgqueues.privateNodeSendPublisher)

	privateHubService := sharedservice.NewPrivateHubService(
		conf.PrivateHubTickerPeriod,
		conf.PrivateHubChainID, // srcChainID - the chain sending messages (Private Hub/CC)
		destEndpointAddress,    // destEndpointAddress - where we execute on Private Node
		r.msgqueues.privateHubMessageConsumer,
		r.txGen.privateNode, // TransactionGenerator
		privateHubBatcher,
	)

	// 3. Cross Chain Service
	crossChainService := destservice.NewCrossChainService(
		destEndpointAddress,
		100,
		r.msgqueues.crossChainConsumer,
		r.repositories.transactionRepository,
		r.repositories.signatureRepository,
		r.txGen.privateNode, // TransactionGenerator
		atomicBatcher,
		vanillaBatcher,
		r.atomicServices.destReceipt,
		r.atomicServices.destVanillaReceipt,
		r.contractClients.nodeEndpoint,
		r.contractClients.nodeDeployer,
	)

	// 4. Enygma Orchestrator
	enygmaOrchestrator := destservice.NewEnygmaOrchestrator(
		r.msgqueues.enygmaDestConsumer,
		r.receivers.enygmaReceiver,
		r.repositories.enygmaCheckpoint,
	)

	// 5. Dvp Orchestrator
	dvpOrchestrator := destservice.NewDvpOrchestrator(
		r.msgqueues.dvpDestConsumer,
		r.receivers.dvpReceiver,
	)

	// 6. Enygma Resync Service (for recovering missing history)
	teleportAddress, err := r.contractClients.hubRegistry.GetContractAddress("EnygmaTeleport")
	if err != nil {
		return fmt.Errorf("failed to get EnygmaTeleport address for resync service: %w", err)
	}

	enygmaResyncTracer := enygmaAdapters.NewOTelTracer("enygma-resync")
	enygmaResyncService, err := resync.NewEnygmaResyncService(
		teleportAddress,
		conf.MyChainID,
		r.repositories.resourceLock,
		r.repositories.enygmaCheckpoint,
		r.hubClient,
		r.ctsClient,
		r.receivers.enygmaReceiver,
		enygmaResyncTracer,
	)
	if err != nil {
		return fmt.Errorf("failed to create enygma resync service: %w", err)
	}

	// 7. Enygma Sync Service (validates checkpoints, triggers resync when needed)
	enygmaSyncService := enygmaService.NewEnygmaSyncService(
		enygmaService.SyncConfig{MaxRetries: conf.EnygmaSyncMaxRetries},
		r.infrastructure.transactionManager,
		r.repositories.enygma,
		r.repositories.enygmaHistory,
		r.repositories.enygmaCheckpoint,
		enygmaResyncService,
	)

	r.pnhResultRouter.Register(msgTypePrivateHubExecute, resultrouter.HandlerFunc(privateHubService.HandleResults))

	signatureSvc := r.atomicServices.signature
	// Decommissioning Teleport (vanilla, atomic): atomic registration below; the vanilla registration is retained.
	r.pnResultRouter.Register(msgTypeCrosschainAtomic, resultrouter.HandlerFunc(crossChainService.HandleAtomicResults))
	r.pnResultRouter.Register(msgTypeCrosschainVanilla, resultrouter.HandlerFunc(crossChainService.HandleVanillaResults))
	// Decommissioning Teleport (vanilla, atomic).
	r.pnResultRouter.Register(msgTypeAtomicDestinationUnlock, resultrouter.HandlerFunc(signatureSvc.HandleDestinationExecutedCallback))
	r.pnResultRouter.Register(msgTypeAtomicDestinationRevert, resultrouter.HandlerFunc(signatureSvc.HandleDestinationRevertedCallback))
	r.pnResultRouter.Register(msgTypeAtomicSourceRevert, resultrouter.HandlerFunc(signatureSvc.HandleSourceRevertedCallback))
	// privatehub.execute results currently have no callback — the
	// router acks them and logs a warning. Add a HandlePrivateHubResults
	// here when observability of those flows is needed.

	r.services = &Services{
		privateHub:         privateHubService,
		crossChain:         crossChainService,
		enygmaOrchestrator: enygmaOrchestrator,
		enygmaSyncService:  enygmaSyncService,
		dvpOrchestrator:    dvpOrchestrator,
	}

	return nil
}
