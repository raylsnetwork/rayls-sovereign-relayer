// Decommissioning Teleport (vanilla, atomic): atomic members below marked; shared/generic/Enygma/DVP retained.

package app

import (
	"fmt"
	"math/big"

	"github.com/raylsnetwork/rayls-sovereign-relayer/contractclient"
	destpoller "github.com/raylsnetwork/rayls-sovereign-relayer/private-relayer/dest/poller"
	destservice "github.com/raylsnetwork/rayls-sovereign-relayer/private-relayer/dest/service"
	desttxgen "github.com/raylsnetwork/rayls-sovereign-relayer/private-relayer/dest/txgen"
	"github.com/raylsnetwork/rayls-sovereign-relayer/private-relayer/proofgen"
	sharedservice "github.com/raylsnetwork/rayls-sovereign-relayer/private-relayer/shared/service"
	"github.com/raylsnetwork/rayls-sovereign-relayer/txsim"
)

type AtomicConfig struct {
	// Decommissioning Teleport (vanilla, atomic).
	BatchSize int
	MyChainID *big.Int
}

// AtomicServices holds the destination-side atomic flow services.
//
// The receipt + vanilla-receipt services are no longer driven by a
// poller — `CrossChainService.HandleAtomicResults` /
// `HandleVanillaResults` invoke them directly when terminal TxResults
// arrive on `cts.result.privatenode`. The pollers are gone with the
// executor DB.
type AtomicServices struct {
	// Decommissioning Teleport (vanilla, atomic).
	destReceipt        *destservice.ReceiptService
	destVanillaReceipt *destservice.VanillaReceiptService
	// Decommissioning Teleport (vanilla, atomic).
	signature *sharedservice.SignatureService
}

// AtomicPollers retains only the FinalizationPoller. The old
// ReceiptPoller polled the executor DB to dispatch receipt callbacks;
// in the async flow that path is owned by the result router. The
// FinalizationPoller is unrelated — it runs against the atomic_status
// table and the SignatureService.
//
// Deprecated: Decommissioning Teleport (vanilla, atomic).
type AtomicPollers struct {
	destFinalization *destpoller.FinalizationPoller
}

func (r *DestPrivateRelayer) initializeAtomic(conf AtomicConfig) error {
	// Create Hub-side teleport client (for atomic message batch operations on CC)
	teleportAddress, err := r.contractClients.hubRegistry.GetContractAddress("Teleport")
	if err != nil {
		return fmt.Errorf("failed to get hub teleport address: %w", err)
	}
	hubTeleport := contractclient.NewTeleportClient(
		teleportAddress,
		r.contractClients.hubExecutor,
		r.contractClients.hubEncryptor,
	)

	// Create Hub-side (CC) transaction simulator (for revert reason analysis)
	hubTxSim := txsim.NewTransactionSimulator(r.hubClient)

	// Create PL-side proof generator (for vanilla receipt proofs)
	proofGen := proofgen.New(r.nodeClient)

	// Get PL-side endpoint address — the SignatureService publishes
	// against this address (it is the calldata target for unlock /
	// destination-revert / source-revert signatures).
	nodeEndpointAddress, err := r.contractClients.nodeRegistry.GetContractAddress("Endpoint")
	if err != nil {
		return fmt.Errorf("failed to get node endpoint address: %w", err)
	}

	// PL-side signature calldata generator (Generate(types.CalldataSignature) ([]byte, error)).
	sigGen, err := desttxgen.NewSignatureGenerator(r.nodeClient, nodeEndpointAddress)
	if err != nil {
		return fmt.Errorf("failed to create signature generator: %w", err)
	}

	// Create PL-side transaction simulator (for signature revert analysis)
	plSim := txsim.NewTransactionSimulator(r.nodeClient)

	// Construct atomic services
	// Decommissioning Teleport (vanilla, atomic).
	destReceiptSvc := destservice.NewReceiptService(
		hubTeleport,
		r.nodeClient,
		r.repositories.transactionRepository,
		hubTxSim,
	)
	destVanillaReceiptSvc := destservice.NewVanillaReceiptService(
		conf.MyChainID,
		hubTeleport,
		r.nodeClient,
		proofGen,
		r.repositories.transactionRepository,
	)
	// SignatureService now publishes TxRequest batches directly to
	// `cts.send.privatenode`; CTS signs/broadcasts and the result
	// router fans the resulting TxResults back into the per-flavour
	// callbacks. The previous local key-queue + txbatchclient path is
	// gone.
	// Decommissioning Teleport (vanilla, atomic).
	signatureSvc := sharedservice.NewSignatureService(
		r.msgqueues.privateNodeSendPublisher,
		hubTeleport,
		r.nodeClient,
		r.repositories.signatureRepository,
		r.repositories.transactionRepository,
		plSim,
		sigGen,
		nodeEndpointAddress,
	)

	r.atomicServices = &AtomicServices{
		destReceipt:        destReceiptSvc,
		destVanillaReceipt: destVanillaReceiptSvc,
		signature:          signatureSvc,
	}

	r.atomicPollers = &AtomicPollers{
		// Decommissioning Teleport (vanilla, atomic).
		destFinalization: destpoller.NewFinalizationPoller(
			conf.BatchSize,
			signatureSvc,
			r.repositories.transactionRepository,
			r.repositories.atomicStatusRepository,
		),
	}

	return nil
}
