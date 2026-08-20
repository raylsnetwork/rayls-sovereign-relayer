package app

import (
	"net/http"
	"time"

	"github.com/raylsnetwork/rayls-sovereign-relayer/client"
	dvpService "github.com/raylsnetwork/rayls-sovereign-relayer/dvp/service"
	merkle "github.com/raylsnetwork/rayls-sovereign-relayer/merkle-service"
	"github.com/raylsnetwork/rayls-sovereign-relayer/private-relayer/repository"
)

// Infrastructure holds additional infrastructure components needed by enygma services
type Infrastructure struct {
	proofAPIClient       *client.ProofAPIClient
	merkleService        *merkle.MerkleService
	transactionManager   *repository.TransactionManager
	dvpMerkleTreeDepth   int
	commitmentCalculator *dvpService.CommitmentCalculator
}

const defaultHTTPClientTimeout = 30 * time.Second

type InfrastructureConfig struct {
	ProofAPIURL        string
	DvpMerkleTreeDepth int
}

func (r *SourcePrivateRelayer) initializeInfrastructure(conf InfrastructureConfig) {
	// Initialize Proof API Client
	proofAPIClient := client.NewProofAPIClient(
		conf.ProofAPIURL,
		client.NewSimpleClient(&http.Client{Timeout: defaultHTTPClientTimeout}),
	)

	// Initialize Merkle Service
	merkleService := merkle.NewMerkleService(
		conf.DvpMerkleTreeDepth,
		r.repositories.merkleTree,
	)

	// Initialize Transaction Manager
	transactionManager := repository.NewTransactionManager(
		r.pool,
	)

	// Initialize Commitment Calculator (shared across enygma and dvp)
	commitmentCalculator := dvpService.NewCommitmentCalculator()

	r.infrastructure = &Infrastructure{
		proofAPIClient:       proofAPIClient,
		merkleService:        merkleService,
		transactionManager:   transactionManager,
		dvpMerkleTreeDepth:   conf.DvpMerkleTreeDepth,
		commitmentCalculator: commitmentCalculator,
	}
}
