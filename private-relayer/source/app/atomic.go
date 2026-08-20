// Decommissioning Teleport (vanilla, atomic).

package app

import (
	"fmt"
	"time"

	sharedservice "github.com/raylsnetwork/rayls-sovereign-relayer/private-relayer/shared/service"
	"github.com/raylsnetwork/rayls-sovereign-relayer/private-relayer/source/poller"
	"github.com/raylsnetwork/rayls-sovereign-relayer/private-relayer/source/service"
	"github.com/raylsnetwork/rayls-sovereign-relayer/private-relayer/source/txgen"
	"github.com/raylsnetwork/rayls-sovereign-relayer/txsim"
)

// Deprecated: Decommissioning Teleport (vanilla, atomic).
type AtomicConfig struct {
	BatchSize      int
	ExpirationTime time.Duration
}

// Deprecated: Decommissioning Teleport (vanilla, atomic).
type AtomicServices struct {
	earlyRevert *service.EarlyRevertService
	expired     *service.ExpiredService
	signature   *sharedservice.SignatureService
}

// Deprecated: Decommissioning Teleport (vanilla, atomic).
type AtomicPollers struct {
	earlyRevert     *poller.EarlyRevertPoller
	expired         *poller.ExpiredPoller
	srcFinalization *poller.FinalizationPoller
}

func (r *SourcePrivateRelayer) initializeAtomic(conf AtomicConfig) error {
	privateNodeEndpointAddress, err := r.contractClients.nodeRegistry.GetContractAddress("Endpoint")
	if err != nil {
		return fmt.Errorf("failed to get private node endpoint address: %w", err)
	}

	// Calldata-only generator. Both EarlyRevertService and SignatureService
	// publish via the privatenode publisher — source-revert / early-revert
	// txs execute on the source PN (not the hub), signed by the PN key.
	// CTS owns nonce/gas/sign/broadcast.
	sigCalldataGen, err := txgen.NewSignatureCalldataGenerator(r.nodeClient, privateNodeEndpointAddress)
	if err != nil {
		return fmt.Errorf("failed to create signature calldata generator: %w", err)
	}

	plSim := txsim.NewTransactionSimulator(r.nodeClient)

	// Construct atomic services
	earlyRevertSvc := service.NewEarlyRevertService(
		r.msgqueues.privateNodeSendPublisher,
		sigCalldataGen,
		privateNodeEndpointAddress,
		r.repositories.signatureRepository,
		r.repositories.transactionRepository,
	)
	expiredSvc := service.NewExpiredService(r.contractClients.teleport, r.repositories.transactionRepository)
	signatureSvc := sharedservice.NewSignatureService(
		r.msgqueues.privateNodeSendPublisher,
		r.contractClients.teleport,
		r.nodeClient,
		r.repositories.signatureRepository,
		r.repositories.transactionRepository,
		plSim,
		sigCalldataGen,
		privateNodeEndpointAddress,
	)

	r.atomicServices = &AtomicServices{
		earlyRevert: earlyRevertSvc,
		expired:     expiredSvc,
		signature:   signatureSvc,
	}

	// Construct atomic pollers
	r.atomicPollers = &AtomicPollers{
		earlyRevert: poller.NewEarlyRevertPoller(
			conf.BatchSize,
			earlyRevertSvc,
			r.repositories.transactionRepository,
		),
		expired: poller.NewExpiredPoller(
			conf.BatchSize,
			conf.ExpirationTime,
			expiredSvc,
			r.repositories.transactionRepository,
		),
		srcFinalization: poller.NewFinalizationPoller(
			conf.BatchSize,
			signatureSvc,
			r.repositories.transactionRepository,
			r.repositories.atomicStatusRepository,
		),
	}

	return nil
}
