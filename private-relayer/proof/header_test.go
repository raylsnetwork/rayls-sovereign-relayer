package proof_test

import (
	"context"
	"errors"
	"math/big"
	"testing"
	"time"

	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/contractclient"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/contracts/Proofs"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/private-relayer/proof"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/testtools"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func defaultTestConfig() proof.HeaderProofConfig {
	return proof.HeaderProofConfig{
		PLChainID:    big.NewInt(99999),
		PollInterval: 1 * time.Hour,
		BatchSize:    100,
		Timeout:      2 * time.Minute,
	}
}

func newDefaultBlockchainClientMock() *HeaderBlockchainClientMock {
	return &HeaderBlockchainClientMock{
		BlockNumberFunc: func(ctx context.Context) (uint64, error) {
			return 10, nil
		},
		HeaderByNumberFunc: func(ctx context.Context, number *big.Int) (*ethTypes.Header, error) {
			return &ethTypes.Header{
				Number:     new(big.Int).Set(number),
				Difficulty: big.NewInt(0),
			}, nil
		},
	}
}

func newDefaultSubmitterMock(initialBlock int64) *HeaderProofsSubmitterMock {
	return &HeaderProofsSubmitterMock{
		GetNextBlockNumberFunc: func(ctx context.Context, chainID *big.Int) (*big.Int, error) {
			return big.NewInt(initialBlock), nil
		},
		SubmitBatchHeadersFunc: func(ctx context.Context, chainID *big.Int, headers []Proofs.ProofsHeader) (*contractclient.SubmitResult, error) {
			lastHeader := headers[len(headers)-1]
			return &contractclient.SubmitResult{
				StartBlock:        headers[0].Number,
				EndBlock:          lastHeader.Number,
				NextExpectedBlock: new(big.Int).Add(lastHeader.Number, big.NewInt(1)),
			}, nil
		},
	}
}

func TestNewHeaderProofService(t *testing.T) {
	testtools.SilenceLogger()

	t.Run("returns error when GetNextBlockNumber fails", func(t *testing.T) {
		submitter := &HeaderProofsSubmitterMock{
			GetNextBlockNumberFunc: func(ctx context.Context, chainID *big.Int) (*big.Int, error) {
				return nil, errors.New("contract call failed")
			},
		}
		client := newDefaultBlockchainClientMock()

		svc, err := proof.NewHeaderProofService(defaultTestConfig(), client, submitter)

		require.Error(t, err)
		assert.Nil(t, svc)
		assert.Contains(t, err.Error(), "failed to get initial block number")
	})

	t.Run("creates service with initial block number from submitter", func(t *testing.T) {
		submitter := newDefaultSubmitterMock(5)
		client := newDefaultBlockchainClientMock()

		svc, err := proof.NewHeaderProofService(defaultTestConfig(), client, submitter)

		require.NoError(t, err)
		assert.NotNil(t, svc)
		require.Len(t, submitter.GetNextBlockNumberCalls(), 1)
		assert.Equal(t, defaultTestConfig().PLChainID, submitter.GetNextBlockNumberCalls()[0].ChainID)
	})
}

func TestRun(t *testing.T) {
	testtools.SilenceLogger()

	t.Run("returns nil on context cancellation", func(t *testing.T) {
		submitter := newDefaultSubmitterMock(5)
		client := newDefaultBlockchainClientMock()

		svc, err := proof.NewHeaderProofService(defaultTestConfig(), client, submitter)
		require.NoError(t, err)

		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		err = svc.Run(ctx)
		assert.NoError(t, err)
	})

	t.Run("skips processing when no new blocks exist", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		// nextBlock = 10, latestBlock = 9 → no gap
		client := &HeaderBlockchainClientMock{
			BlockNumberFunc: func(ctx context.Context) (uint64, error) {
				cancel()
				return 9, nil
			},
		}
		submitter := newDefaultSubmitterMock(10)

		svc, err := proof.NewHeaderProofService(defaultTestConfig(), client, submitter)
		require.NoError(t, err)

		err = svc.Run(ctx)
		assert.NoError(t, err)

		assert.Len(t, client.BlockNumberCalls(), 1)
		assert.Len(t, submitter.SubmitBatchHeadersCalls(), 0)
	})

	t.Run("fetches and submits headers for block gap", func(t *testing.T) {
		// nextBlock=5, latestBlock=7 → fetch blocks 5, 6, 7
		ctx, cancel := context.WithCancel(context.Background())
		submitter := &HeaderProofsSubmitterMock{
			GetNextBlockNumberFunc: func(ctx context.Context, chainID *big.Int) (*big.Int, error) {
				return big.NewInt(5), nil
			},
			SubmitBatchHeadersFunc: func(ctx context.Context, chainID *big.Int, headers []Proofs.ProofsHeader) (*contractclient.SubmitResult, error) {
				cancel()
				return &contractclient.SubmitResult{
					StartBlock:        big.NewInt(5),
					EndBlock:          big.NewInt(7),
					NextExpectedBlock: big.NewInt(8),
				}, nil
			},
		}
		client := &HeaderBlockchainClientMock{
			BlockNumberFunc: func(ctx context.Context) (uint64, error) {
				return 7, nil
			},
			HeaderByNumberFunc: func(ctx context.Context, number *big.Int) (*ethTypes.Header, error) {
				return &ethTypes.Header{
					Number:     new(big.Int).Set(number),
					Difficulty: big.NewInt(0),
				}, nil
			},
		}

		svc, err := proof.NewHeaderProofService(defaultTestConfig(), client, submitter)
		require.NoError(t, err)

		err = svc.Run(ctx)
		assert.NoError(t, err)

		headerCalls := client.HeaderByNumberCalls()
		require.Len(t, headerCalls, 3)
		assert.Equal(t, int64(5), headerCalls[0].Number.Int64())
		assert.Equal(t, int64(6), headerCalls[1].Number.Int64())
		assert.Equal(t, int64(7), headerCalls[2].Number.Int64())

		submitCalls := submitter.SubmitBatchHeadersCalls()
		require.Len(t, submitCalls, 1)
		assert.Len(t, submitCalls[0].Headers, 3)
		assert.Equal(t, defaultTestConfig().PLChainID, submitCalls[0].ChainID)
	})

	t.Run("clamps batch to latest block when gap exceeds batch size", func(t *testing.T) {
		// batchSize=2, nextBlock=5, latestBlock=100
		// endingBlock = 5 + 2 = 7, clamped to min(7, 100) = 7
		ctx, cancel := context.WithCancel(context.Background())
		config := defaultTestConfig()
		config.BatchSize = 2

		submitter := &HeaderProofsSubmitterMock{
			GetNextBlockNumberFunc: func(ctx context.Context, chainID *big.Int) (*big.Int, error) {
				return big.NewInt(5), nil
			},
			SubmitBatchHeadersFunc: func(ctx context.Context, chainID *big.Int, headers []Proofs.ProofsHeader) (*contractclient.SubmitResult, error) {
				cancel()
				return &contractclient.SubmitResult{
					NextExpectedBlock: big.NewInt(8),
				}, nil
			},
		}
		client := &HeaderBlockchainClientMock{
			BlockNumberFunc: func(ctx context.Context) (uint64, error) {
				return 100, nil
			},
			HeaderByNumberFunc: func(ctx context.Context, number *big.Int) (*ethTypes.Header, error) {
				return &ethTypes.Header{
					Number:     new(big.Int).Set(number),
					Difficulty: big.NewInt(0),
				}, nil
			},
		}

		svc, err := proof.NewHeaderProofService(config, client, submitter)
		require.NoError(t, err)

		err = svc.Run(ctx)
		assert.NoError(t, err)

		headerCalls := client.HeaderByNumberCalls()
		require.Len(t, headerCalls, 3)
		assert.Equal(t, int64(5), headerCalls[0].Number.Int64())
		assert.Equal(t, int64(7), headerCalls[len(headerCalls)-1].Number.Int64())
	})

	t.Run("updates next block number from submit result", func(t *testing.T) {
		// First cycle: nextBlock=5, latestBlock=7, result.NextExpectedBlock=10
		// Second cycle: latestBlock=9, lastDelivered=10-1=9, 9<=9 → skip
		ctx, cancel := context.WithCancel(context.Background())
		config := defaultTestConfig()
		config.PollInterval = 1 * time.Millisecond

		blockNumberCallCount := 0
		client := &HeaderBlockchainClientMock{
			BlockNumberFunc: func(ctx context.Context) (uint64, error) {
				blockNumberCallCount++
				if blockNumberCallCount == 1 {
					return 7, nil
				}
				// Second call: latestBlock=9, lastDelivered=10-1=9, no gap → skip
				cancel()
				return 9, nil
			},
			HeaderByNumberFunc: func(ctx context.Context, number *big.Int) (*ethTypes.Header, error) {
				return &ethTypes.Header{
					Number:     new(big.Int).Set(number),
					Difficulty: big.NewInt(0),
				}, nil
			},
		}
		submitter := &HeaderProofsSubmitterMock{
			GetNextBlockNumberFunc: func(ctx context.Context, chainID *big.Int) (*big.Int, error) {
				return big.NewInt(5), nil
			},
			SubmitBatchHeadersFunc: func(ctx context.Context, chainID *big.Int, headers []Proofs.ProofsHeader) (*contractclient.SubmitResult, error) {
				return &contractclient.SubmitResult{
					NextExpectedBlock: big.NewInt(10),
				}, nil
			},
		}

		svc, err := proof.NewHeaderProofService(config, client, submitter)
		require.NoError(t, err)

		err = svc.Run(ctx)
		assert.NoError(t, err)

		assert.GreaterOrEqual(t, blockNumberCallCount, 2)
		assert.Len(t, submitter.SubmitBatchHeadersCalls(), 1)
	})

	t.Run("continues processing when BlockNumber returns error", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		config := defaultTestConfig()
		config.PollInterval = 1 * time.Millisecond

		callCount := 0
		client := &HeaderBlockchainClientMock{
			BlockNumberFunc: func(ctx context.Context) (uint64, error) {
				callCount++
				if callCount >= 2 {
					cancel()
				}
				return 0, errors.New("rpc error")
			},
		}
		submitter := newDefaultSubmitterMock(5)

		svc, err := proof.NewHeaderProofService(config, client, submitter)
		require.NoError(t, err)

		err = svc.Run(ctx)
		assert.NoError(t, err)

		assert.GreaterOrEqual(t, callCount, 2)
		assert.Len(t, submitter.SubmitBatchHeadersCalls(), 0)
	})

	t.Run("continues processing when HeaderByNumber returns error", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		config := defaultTestConfig()
		config.PollInterval = 1 * time.Millisecond

		blockNumberCallCount := 0
		client := &HeaderBlockchainClientMock{
			BlockNumberFunc: func(ctx context.Context) (uint64, error) {
				blockNumberCallCount++
				if blockNumberCallCount >= 2 {
					cancel()
				}
				return 10, nil
			},
			HeaderByNumberFunc: func(ctx context.Context, number *big.Int) (*ethTypes.Header, error) {
				return nil, errors.New("header fetch failed")
			},
		}
		submitter := newDefaultSubmitterMock(5)

		svc, err := proof.NewHeaderProofService(config, client, submitter)
		require.NoError(t, err)

		err = svc.Run(ctx)
		assert.NoError(t, err)

		assert.GreaterOrEqual(t, blockNumberCallCount, 2)
		assert.Len(t, submitter.SubmitBatchHeadersCalls(), 0)
	})

	t.Run("continues processing when SubmitBatchHeaders returns error", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		config := defaultTestConfig()
		config.PollInterval = 1 * time.Millisecond

		blockNumberCallCount := 0
		client := &HeaderBlockchainClientMock{
			BlockNumberFunc: func(ctx context.Context) (uint64, error) {
				blockNumberCallCount++
				if blockNumberCallCount >= 2 {
					cancel()
				}
				return 10, nil
			},
			HeaderByNumberFunc: func(ctx context.Context, number *big.Int) (*ethTypes.Header, error) {
				return &ethTypes.Header{
					Number:     new(big.Int).Set(number),
					Difficulty: big.NewInt(0),
				}, nil
			},
		}
		submitter := &HeaderProofsSubmitterMock{
			GetNextBlockNumberFunc: func(ctx context.Context, chainID *big.Int) (*big.Int, error) {
				return big.NewInt(5), nil
			},
			SubmitBatchHeadersFunc: func(ctx context.Context, chainID *big.Int, headers []Proofs.ProofsHeader) (*contractclient.SubmitResult, error) {
				return nil, errors.New("submission failed")
			},
		}

		svc, err := proof.NewHeaderProofService(config, client, submitter)
		require.NoError(t, err)

		err = svc.Run(ctx)
		assert.NoError(t, err)

		assert.GreaterOrEqual(t, blockNumberCallCount, 2)
	})

	t.Run("handles incorrect parent hash event without panic", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		submitter := &HeaderProofsSubmitterMock{
			GetNextBlockNumberFunc: func(ctx context.Context, chainID *big.Int) (*big.Int, error) {
				return big.NewInt(5), nil
			},
			SubmitBatchHeadersFunc: func(ctx context.Context, chainID *big.Int, headers []Proofs.ProofsHeader) (*contractclient.SubmitResult, error) {
				cancel()
				return &contractclient.SubmitResult{
					NextExpectedBlock: big.NewInt(6),
					IncorrectHashEvent: &Proofs.ProofsIncorrectParentHashEvent{
						ChainId:              big.NewInt(99999),
						BlockNumber:          big.NewInt(6),
						ParentHash:           [32]byte{0xaa},
						CalculatedParentHash: [32]byte{0xbb},
					},
				}, nil
			},
		}
		client := &HeaderBlockchainClientMock{
			BlockNumberFunc: func(ctx context.Context) (uint64, error) {
				return 7, nil
			},
			HeaderByNumberFunc: func(ctx context.Context, number *big.Int) (*ethTypes.Header, error) {
				return &ethTypes.Header{
					Number:     new(big.Int).Set(number),
					Difficulty: big.NewInt(0),
				}, nil
			},
		}

		svc, err := proof.NewHeaderProofService(defaultTestConfig(), client, submitter)
		require.NoError(t, err)

		err = svc.Run(ctx)
		assert.NoError(t, err)

		assert.Len(t, submitter.SubmitBatchHeadersCalls(), 1)
	})

	t.Run("supports graceful shutdown", func(t *testing.T) {
		submitter := newDefaultSubmitterMock(5)
		client := &HeaderBlockchainClientMock{
			BlockNumberFunc: func(ctx context.Context) (uint64, error) {
				<-ctx.Done()
				return 0, ctx.Err()
			},
		}

		svc, err := proof.NewHeaderProofService(defaultTestConfig(), client, submitter)
		require.NoError(t, err)

		hasGracefulShutdown := testtools.ShutdownFixture(t, svc.Run, time.Second)
		assert.True(t, hasGracefulShutdown)
	})
}
