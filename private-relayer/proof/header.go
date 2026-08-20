package proof

import (
	"context"
	"fmt"
	"log/slog"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/raylsnetwork/rayls-sovereign-relayer/contractclient"
	"github.com/raylsnetwork/rayls-sovereign-relayer/contracts/Proofs"
)

//go:generate moq --pkg proof_test -out header_mock_test.go . HeaderBlockchainClient HeaderProofsSubmitter

// HeaderBlockchainClient defines methods for fetching headers from Private Ledger
type HeaderBlockchainClient interface {
	BlockNumber(ctx context.Context) (uint64, error)
	HeaderByNumber(ctx context.Context, number *big.Int) (*ethTypes.Header, error)
}

// HeaderProofsSubmitter defines methods for submitting proofs to Private Hub
type HeaderProofsSubmitter interface {
	GetNextBlockNumber(ctx context.Context, chainID *big.Int) (*big.Int, error)
	SubmitBatchHeaders(
		ctx context.Context,
		chainID *big.Int,
		headers []Proofs.ProofsHeader,
	) (*contractclient.SubmitResult, error)
}

// HeaderProofService polls for new block headers and submits them as proofs
type HeaderProofService struct {
	config    HeaderProofConfig
	pnClient  HeaderBlockchainClient
	submitter HeaderProofsSubmitter

	nextBlockNumber *big.Int
	ticker          *time.Ticker
	initialRun      chan struct{}
}

func NewHeaderProofService(
	config HeaderProofConfig,
	pnClient HeaderBlockchainClient,
	submitter HeaderProofsSubmitter,
) (*HeaderProofService, error) {
	// Get initial block number from contract
	nextBlock, err := submitter.GetNextBlockNumber(context.Background(), config.PLChainID)
	if err != nil {
		return nil, fmt.Errorf("failed to get initial block number: %w", err)
	}

	initialRun := make(chan struct{}, 1)
	initialRun <- struct{}{}

	return &HeaderProofService{
		config:          config,
		pnClient:        pnClient,
		submitter:       submitter,
		nextBlockNumber: nextBlock,

		ticker:     time.NewTicker(config.PollInterval),
		initialRun: initialRun,
	}, nil
}

func (s *HeaderProofService) Run(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			s.ticker.Stop()
			return nil
		case <-s.ticker.C:
		case <-s.initialRun:
		}
		latestBlockUint64, err := s.pnClient.BlockNumber(ctx)
		if err != nil {
			slog.Error("Failed to get latest block number", slog.Any("error", err))
			continue
		}
		latestBlock := new(big.Int).SetUint64(latestBlockUint64)

		lastDelivered := new(big.Int).Sub(s.nextBlockNumber, big.NewInt(1))
		if latestBlock.Cmp(lastDelivered) <= 0 {
			continue
		}

		timeoutCtx, cancel := context.WithTimeout(ctx, s.config.Timeout)

		if err := s.sendHeaderProofs(timeoutCtx, latestBlock); err != nil {
			slog.Error("Error processing headers", slog.Any("error", err))
		}

		cancel()
	}
}

func (s *HeaderProofService) sendHeaderProofs(ctx context.Context, latestBlock *big.Int) error {
	startingBlock := new(big.Int).Set(s.nextBlockNumber)
	endingBlock := new(big.Int).Add(startingBlock, big.NewInt(s.config.BatchSize))

	if endingBlock.Cmp(latestBlock) > 0 {
		endingBlock = latestBlock
	}

	// Fetch headers
	headers, err := s.fetchHeaders(ctx, startingBlock, endingBlock)
	if err != nil {
		return fmt.Errorf("fetching headers: %w", err)
	}

	slog.Info("Sending Header Proof for blocks",
		slog.Int("Beginning", int(startingBlock.Int64())),
		slog.Int("Ending", int(endingBlock.Int64())),
	)

	// Submit batch
	result, err := s.submitter.SubmitBatchHeaders(ctx, s.config.PLChainID, headers)
	if err != nil {
		return fmt.Errorf("submitting batch headers: %w", err)
	}

	// Handle incorrect parent hash event
	if result.IncorrectHashEvent != nil {
		event := result.IncorrectHashEvent
		slog.Error("incorrect header parent hash",
			slog.Int64("block_number", event.BlockNumber.Int64()),
			slog.String("parent_hash", common.Hash(event.ParentHash).String()),
			slog.String("calculated_parent_hash", common.Hash(event.CalculatedParentHash).String()),
		)
	}

	// Update next block number
	s.nextBlockNumber = result.NextExpectedBlock

	return nil
}

func (s *HeaderProofService) fetchHeaders(ctx context.Context, start, end *big.Int) ([]Proofs.ProofsHeader, error) {
	var headers []Proofs.ProofsHeader

	for i := new(big.Int).Set(start); i.Cmp(end) <= 0; i.Add(i, big.NewInt(1)) {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
			header, err := s.pnClient.HeaderByNumber(ctx, new(big.Int).Set(i))
			if err != nil {
				return nil, fmt.Errorf("failed to fetch header %d: %w", i, err)
			}
			headers = append(headers, ConvertEthHeader(header))
		}
	}

	return headers, nil
}
