package proofgen_test

import (
	"context"
	"errors"
	"math/big"
	"net"
	"sync/atomic"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/raylsnetwork/rayls-sovereign-relayer/cryptography/proofs"
	"github.com/raylsnetwork/rayls-sovereign-relayer/private-relayer/proofgen"
	"github.com/raylsnetwork/rayls-sovereign-relayer/testtools"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// SpyProofsGeneratorFunc captures arguments passed to generateProofs.
type SpyProofsGeneratorFunc struct {
	proof *proofs.ProofDB
	err   error

	spyTrie       *proofs.Trie
	spyIndex      uint
	spyTxRootHash common.Hash
}

func (s *SpyProofsGeneratorFunc) GenerateProofs(
	trie *proofs.Trie,
	txToProve uint,
	transactionRoot common.Hash,
) (*proofs.ProofDB, error) {
	s.spyTrie = trie
	s.spyIndex = txToProve
	s.spyTxRootHash = transactionRoot

	return s.proof, s.err
}

// SpyTrieGeneratorFunc captures arguments passed to generateTrie.
type SpyTrieGeneratorFunc struct {
	trie *proofs.Trie
	err  error

	spyTransactions []*types.Transaction
	callCount       int
}

func (s *SpyTrieGeneratorFunc) GenerateTrie(txs []*types.Transaction) (*proofs.Trie, error) {
	s.spyTransactions = txs
	s.callCount++
	return s.trie, s.err
}

// newTestBlock creates a block with the given header TxHash and transactions.
func newTestBlock(txRootHash common.Hash, txs []*types.Transaction) *types.Block {
	header := &types.Header{TxHash: txRootHash}
	return types.NewBlockWithHeader(header).WithBody(types.Body{Transactions: txs})
}

func TestProofsGenerator_Generate(t *testing.T) {
	t.Run("returns proof for transaction hash", func(t *testing.T) {
		wantProof := proofs.NewProofDB()
		require.NoError(t, wantProof.Put([]byte("key"), []byte("value")))

		wantTrie := proofs.NewTrie()
		require.NoError(t, wantTrie.Put([]byte("another-key"), []byte("another-value")))

		wantTransactions := []*types.Transaction{
			types.NewTx(&types.LegacyTx{Value: new(big.Int).SetUint64(1337)}),
		}

		wantBlockHash := common.Hash([32]byte{0xDE, 0xAD, 0xBE, 0xEF})
		wantTxHash := common.Hash([32]byte{0xC0, 0x01, 0xBA, 0xBE})
		wantTxRootHash := common.Hash([32]byte{0xCA, 0xFE, 0xBA, 0xBE})
		wantIndex := uint(42)

		block := newTestBlock(wantTxRootHash, wantTransactions)
		receipt := &types.Receipt{BlockHash: wantBlockHash, TransactionIndex: wantIndex}

		ethClient := &EthereumClientMock{
			TransactionReceiptFunc: func(_ context.Context, hash common.Hash) (*types.Receipt, error) {
				return receipt, nil
			},
			BlockByHashFunc: func(_ context.Context, hash common.Hash) (*types.Block, error) {
				return block, nil
			},
		}
		spyProofsGen := SpyProofsGeneratorFunc{proof: wantProof}
		spyTrieGen := SpyTrieGeneratorFunc{trie: wantTrie}
		gen := proofgen.NewCustom(ethClient, spyProofsGen.GenerateProofs, spyTrieGen.GenerateTrie)

		gotProof, err := gen.Generate(context.Background(), wantTxHash)
		require.NoError(t, err)

		// Verify eth client received correct arguments
		require.Len(t, ethClient.TransactionReceiptCalls(), 1)
		assert.Equal(t, wantTxHash, ethClient.TransactionReceiptCalls()[0].Hash)

		require.Len(t, ethClient.BlockByHashCalls(), 1)
		assert.Equal(t, wantBlockHash, ethClient.BlockByHashCalls()[0].Hash)

		// Verify trie generation received correct transactions
		assert.Equal(t, wantTransactions, spyTrieGen.spyTransactions)

		// Verify proof generation received correct arguments
		assert.Equal(t, wantTrie, spyProofsGen.spyTrie)
		assert.Equal(t, wantIndex, spyProofsGen.spyIndex)
		assert.Equal(t, wantTxRootHash, spyProofsGen.spyTxRootHash)

		// Verify proof bytes
		wantProofBytes, err := wantProof.Export()
		require.NoError(t, err)
		assert.Equal(t, wantProofBytes, gotProof)
	})

	t.Run("wraps error in ProofsGeneratorError", func(t *testing.T) {
		wantErr := errors.New("example error")

		block := newTestBlock(
			common.Hash([32]byte{0xCA, 0xFE, 0xBA, 0xBE}),
			[]*types.Transaction{types.NewTx(&types.LegacyTx{Value: new(big.Int).SetUint64(1337)})},
		)
		receipt := &types.Receipt{
			BlockHash:        common.Hash([32]byte{0xDE, 0xAD, 0xBE, 0xEF}),
			TransactionIndex: 42,
		}

		cases := []struct {
			name                  string
			transactionReceiptErr error
			blockByHashErr        error
			generateTrieErr       error
			generateProofsErr     error
		}{
			{name: "on error from TransactionReceipt", transactionReceiptErr: wantErr},
			{name: "on error from BlockByHash", blockByHashErr: wantErr},
			{name: "on error from GenerateTrie", generateTrieErr: wantErr},
			{name: "on error from GenerateProofs", generateProofsErr: wantErr},
		}

		for _, tc := range cases {
			t.Run(tc.name, func(t *testing.T) {
				ethClient := &EthereumClientMock{
					TransactionReceiptFunc: func(_ context.Context, _ common.Hash) (*types.Receipt, error) {
						return receipt, tc.transactionReceiptErr
					},
					BlockByHashFunc: func(_ context.Context, _ common.Hash) (*types.Block, error) {
						return block, tc.blockByHashErr
					},
				}
				spyProofsGen := SpyProofsGeneratorFunc{proof: proofs.NewProofDB(), err: tc.generateProofsErr}
				spyTrieGen := SpyTrieGeneratorFunc{trie: proofs.NewTrie(), err: tc.generateTrieErr}

				gen := proofgen.NewCustom(ethClient, spyProofsGen.GenerateProofs, spyTrieGen.GenerateTrie)
				_, err := gen.Generate(context.Background(), common.Hash([32]byte{0xC0, 0x01, 0xBA, 0xBE}))

				assert.ErrorAs(t, err, new(*proofgen.ProofGeneratorError))
				assert.ErrorIs(t, err, wantErr)
			})
		}
	})

	t.Run("uses cached trie on second call for same block", func(t *testing.T) {
		wantTrie := proofs.NewTrie()
		require.NoError(t, wantTrie.Put([]byte("key"), []byte("value")))

		tx := types.NewTx(&types.LegacyTx{Value: new(big.Int).SetUint64(1)})
		block := newTestBlock(common.Hash([32]byte{0xAA}), []*types.Transaction{tx})
		receipt := &types.Receipt{BlockHash: common.Hash([32]byte{0xBB}), TransactionIndex: 0}

		ethClient := &EthereumClientMock{
			TransactionReceiptFunc: func(_ context.Context, _ common.Hash) (*types.Receipt, error) {
				return receipt, nil
			},
			BlockByHashFunc: func(_ context.Context, _ common.Hash) (*types.Block, error) {
				return block, nil
			},
		}
		spyProofsGen := SpyProofsGeneratorFunc{proof: proofs.NewProofDB()}
		spyTrieGen := SpyTrieGeneratorFunc{trie: wantTrie}

		gen := proofgen.NewCustom(ethClient, spyProofsGen.GenerateProofs, spyTrieGen.GenerateTrie)

		_, err := gen.Generate(context.Background(), common.Hash([32]byte{0x01}))
		require.NoError(t, err)

		_, err = gen.Generate(context.Background(), common.Hash([32]byte{0x02}))
		require.NoError(t, err)

		assert.Equal(t, 1, spyTrieGen.callCount, "trie should be generated only once; second call should use cache")
	})
}

func TestProofsGenerator_BatchGenerate(t *testing.T) {
	testtools.SilenceLogger()

	t.Run("returns empty slice for empty input", func(t *testing.T) {
		ethClient := &EthereumClientMock{}
		gen := proofgen.NewCustom(ethClient, nil, nil)

		result := gen.BatchGenerate(context.Background(), nil, 4)
		assert.Empty(t, result)
	})

	t.Run("generates proofs for multiple transactions", func(t *testing.T) {
		wantProof := proofs.NewProofDB()
		require.NoError(t, wantProof.Put([]byte("k"), []byte("v")))

		tx := types.NewTx(&types.LegacyTx{Value: new(big.Int).SetUint64(1)})
		block := newTestBlock(common.Hash([32]byte{0xAA}), []*types.Transaction{tx})
		receipt := &types.Receipt{BlockHash: common.Hash([32]byte{0xBB}), TransactionIndex: 0}

		ethClient := &EthereumClientMock{
			TransactionReceiptFunc: func(_ context.Context, _ common.Hash) (*types.Receipt, error) {
				return receipt, nil
			},
			BlockByHashFunc: func(_ context.Context, _ common.Hash) (*types.Block, error) {
				return block, nil
			},
		}

		var proofCalls atomic.Int64
		generateProofs := func(_ *proofs.Trie, _ uint, _ common.Hash) (*proofs.ProofDB, error) {
			proofCalls.Add(1)
			return wantProof, nil
		}
		generateTrie := func(_ []*types.Transaction) (*proofs.Trie, error) {
			return proofs.NewTrie(), nil
		}

		gen := proofgen.NewCustom(ethClient, generateProofs, generateTrie)

		hashes := []common.Hash{{0x01}, {0x02}, {0x03}}
		result := gen.BatchGenerate(context.Background(), hashes, 2)

		require.Len(t, result, 3)
		wantProofBytes, err := wantProof.Export()
		require.NoError(t, err)
		for i, proof := range result {
			assert.NotNil(t, proof, "proof at index %d should not be nil", i)
			assert.Equal(t, wantProofBytes, proof)
		}
		assert.Equal(t, int64(3), proofCalls.Load())
	})

	t.Run("limits routines to number of hashes", func(t *testing.T) {
		tx := types.NewTx(&types.LegacyTx{Value: new(big.Int).SetUint64(1)})
		block := newTestBlock(common.Hash([32]byte{0xAA}), []*types.Transaction{tx})
		receipt := &types.Receipt{BlockHash: common.Hash([32]byte{0xBB}), TransactionIndex: 0}

		ethClient := &EthereumClientMock{
			TransactionReceiptFunc: func(_ context.Context, _ common.Hash) (*types.Receipt, error) {
				return receipt, nil
			},
			BlockByHashFunc: func(_ context.Context, _ common.Hash) (*types.Block, error) {
				return block, nil
			},
		}
		generateProofs := func(_ *proofs.Trie, _ uint, _ common.Hash) (*proofs.ProofDB, error) {
			return proofs.NewProofDB(), nil
		}
		generateTrie := func(_ []*types.Transaction) (*proofs.Trie, error) {
			return proofs.NewTrie(), nil
		}

		gen := proofgen.NewCustom(ethClient, generateProofs, generateTrie)

		// 2 hashes with 10 routines — should not panic
		hashes := []common.Hash{{0x01}, {0x02}}
		result := gen.BatchGenerate(context.Background(), hashes, 10)

		require.Len(t, result, 2)
		for i, proof := range result {
			assert.NotNil(t, proof, "proof at index %d should not be nil", i)
		}
	})

	t.Run("returns nil when proof generation fails with a non-transient error", func(t *testing.T) {
		var receiptCalls atomic.Int64
		ethClient := &EthereumClientMock{
			TransactionReceiptFunc: func(_ context.Context, _ common.Hash) (*types.Receipt, error) {
				receiptCalls.Add(1)
				return nil, errors.New("execution reverted")
			},
			BlockByHashFunc: func(_ context.Context, _ common.Hash) (*types.Block, error) {
				return nil, errors.New("should not be called")
			},
		}
		generateProofs := func(_ *proofs.Trie, _ uint, _ common.Hash) (*proofs.ProofDB, error) {
			return nil, errors.New("should not be called")
		}
		generateTrie := func(_ []*types.Transaction) (*proofs.Trie, error) {
			return nil, errors.New("should not be called")
		}

		gen := proofgen.NewCustom(ethClient, generateProofs, generateTrie)

		hashes := []common.Hash{{0x01}, {0x02}, {0x03}}
		result := gen.BatchGenerate(context.Background(), hashes, 2)

		require.Len(t, result, 3)
		for i, proof := range result {
			assert.Nil(t, proof, "proof at index %d should be nil on failure", i)
		}
		// A non-transient error is not retried by ethretry: exactly one attempt per hash.
		assert.Equal(t, int64(len(hashes)), receiptCalls.Load(),
			"a non-transient error must not be retried")
	})

	t.Run("routes RPC through ethretry: retries transient errors and stops on ctx cancellation", func(t *testing.T) {
		var receiptCalls atomic.Int64
		ethClient := &EthereumClientMock{
			TransactionReceiptFunc: func(_ context.Context, _ common.Hash) (*types.Receipt, error) {
				receiptCalls.Add(1)
				// A connection error is transient, so ethretry enters its retry loop
				// (rather than returning immediately as it would for a plain error).
				return nil, &net.OpError{Op: "dial", Net: "tcp", Err: errors.New("connection refused")}
			},
			BlockByHashFunc: func(_ context.Context, _ common.Hash) (*types.Block, error) {
				return nil, errors.New("should not be called")
			},
		}
		generateProofs := func(_ *proofs.Trie, _ uint, _ common.Hash) (*proofs.ProofDB, error) {
			return nil, errors.New("should not be called")
		}
		generateTrie := func(_ []*types.Transaction) (*proofs.Trie, error) {
			return nil, errors.New("should not be called")
		}

		gen := proofgen.NewCustom(ethClient, generateProofs, generateTrie)

		// Already-cancelled context: the transient error routes into ethretry's
		// retry loop, which then observes the cancelled context and gives up — so
		// the routine surfaces a nil proof instead of blocking forever on a
		// persistent connection error. This proves the RPC path is wrapped in
		// ethretry and honours context cancellation.
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		result := gen.BatchGenerate(ctx, []common.Hash{{0x01}}, 1)

		require.Len(t, result, 1)
		assert.Nil(t, result[0], "proof should be nil when the retry loop is cancelled")
		assert.Equal(t, int64(1), receiptCalls.Load(),
			"the RPC is attempted once, then the cancelled context stops the retry loop")
	})
}
