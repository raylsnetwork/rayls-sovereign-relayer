package proofgen

import (
	"context"
	"log/slog"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/raylsnetwork/rayls-sovereign-relayer/cache"
	"github.com/raylsnetwork/rayls-sovereign-relayer/cryptography/proofs"
	"github.com/raylsnetwork/rayls-sovereign-relayer/ethretry"
	"github.com/raylsnetwork/rayls-sovereign-relayer/faultinjector"
	"github.com/raylsnetwork/rayls-sovereign-relayer/withstack"
)

type (
	generateProofsFunc func(*proofs.Trie, uint, common.Hash) (*proofs.ProofDB, error)
	generateTrieFunc   func([]*types.Transaction) (*proofs.Trie, error)
)

// rpcCallTimeout bounds one receipt/block RPC attempt; ethretry retries per attempt.
const rpcCallTimeout = 30 * time.Second

//go:generate moq --pkg proofgen_test -out proofgen_mock_test.go . EthereumClient

// EthereumClient defines the blockchain methods needed for proof generation.
type EthereumClient interface {
	BlockByHash(context.Context, common.Hash) (*types.Block, error)
	TransactionReceipt(context.Context, common.Hash) (*types.Receipt, error)
}

type ProofGenerator struct {
	generateProofs generateProofsFunc
	generateTrie   generateTrieFunc

	trieCache *cache.LRU[string, *proofs.Trie]

	ethClient EthereumClient
}

func New(ethClient EthereumClient) *ProofGenerator {
	return &ProofGenerator{
		generateProofs: proofs.GenerateProofs,
		generateTrie:   proofs.GenerateTrie,

		trieCache: cache.NewLRU[string, *proofs.Trie](10),

		ethClient: ethClient,
	}
}

func NewCustom(
	ethClient EthereumClient,
	generateProofs generateProofsFunc,
	generateTrie generateTrieFunc,
) *ProofGenerator {
	return &ProofGenerator{
		generateProofs: generateProofs,
		generateTrie:   generateTrie,

		// TODO: Make bigger
		trieCache: cache.NewLRU[string, *proofs.Trie](10),

		ethClient: ethClient,
	}
}

func divCeil(a, b int) int {
	if b == 0 {
		panic("division by zero")
	}
	return (a + b - 1) / b
}

func (g *ProofGenerator) Generate(ctx context.Context, txHash common.Hash) ([]byte, error) {
	if faultErr := faultinjector.Check("private_relayer.proofgen.ProofGenerator.Generate.start"); faultErr != nil {
		return nil, WrapInProofGeneratorError("fault-injected proof generation failure", withstack.Wrap(faultErr))
	}

	// Retry the receipt/block fetches on transient RPC errors (connection, decode,
	// per-attempt timeout) until success or ctx cancellation. ethereum.NotFound and
	// RPC application errors are non-transient and fall through to proof_invalid.
	var receipt *types.Receipt
	if err := ethretry.WithRetry(ctx, func() error {
		ctxReceipt, cancel := context.WithTimeout(ctx, rpcCallTimeout)
		defer cancel()
		var err error
		receipt, err = g.ethClient.TransactionReceipt(ctxReceipt, txHash)
		return err
	}); err != nil {
		return nil, WrapInProofGeneratorError("got error while trying to get transaction receipt", withstack.Wrap(err))
	}

	var block *types.Block
	if err := ethretry.WithRetry(ctx, func() error {
		ctxBlock, cancel := context.WithTimeout(ctx, rpcCallTimeout)
		defer cancel()
		var err error
		block, err = g.ethClient.BlockByHash(ctxBlock, receipt.BlockHash)
		return err
	}); err != nil {
		return nil, WrapInProofGeneratorError("got error while trying to get block by hash", withstack.Wrap(err))
	}

	trie, ok := g.trieCache.Get(block.Hash().String())
	if !ok {
		var err error
		trie, err = g.generateTrie(block.Transactions())
		if err != nil {
			return nil, WrapInProofGeneratorError("got error while generating trie", err)
		}
		g.trieCache.Put(block.Hash().String(), trie)
	}

	proof, err := g.generateProofs(trie, receipt.TransactionIndex, block.TxHash())
	if err != nil {
		return nil, WrapInProofGeneratorError("got error while generating proof", err)
	}

	exported, err := proof.Export()
	if err != nil {
		return nil, WrapInProofGeneratorError("got error while exporting proof", err)
	}

	return exported, nil
}

func (g *ProofGenerator) BatchGenerate(ctx context.Context, txHashes []common.Hash, routineCount int) [][]byte {
	proofSlice := make([][]byte, len(txHashes))

	if len(txHashes) == 0 {
		return proofSlice
	}

	// Limit routines to the number of hashes
	actualRoutineCount := min(routineCount, len(txHashes))
	batchSize := divCeil(len(txHashes), actualRoutineCount)

	var wg sync.WaitGroup
	wg.Add(actualRoutineCount)

	for i := range actualRoutineCount {
		go func() {
			defer wg.Done()

			start := i * batchSize
			if start >= len(txHashes) {
				return // Nothing to process for this routine
			}

			end := (i + 1) * batchSize
			if end > len(txHashes) {
				end = len(txHashes)
			}

			g.batchGenerateRoutine(ctx, txHashes[start:end], proofSlice[start:end])
		}()
	}

	wg.Wait()
	return proofSlice
}

func (g *ProofGenerator) batchGenerateRoutine(ctx context.Context, txHashes []common.Hash, proofSlice [][]byte) {
	for i, txHash := range txHashes {
		proof, err := g.Generate(ctx, txHash)
		if err != nil {
			// Log the cause and leave the slot nil; the dispatch loop records it
			// as proof_invalid rather than dispatching a nil proof.
			slog.Error(
				"Failed to generate proof for transaction",
				slog.Any("tx_hash", txHash),
				slog.Any("error", err),
			)
			continue
		}
		proofSlice[i] = proof
	}
}
