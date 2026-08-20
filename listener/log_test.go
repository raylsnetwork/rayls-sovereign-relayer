//go:build integration

package listener_test

import (
	"context"
	"errors"
	"math/big"
	"testing"
	"time"

	"raylz-relayer/listener"
	"raylz-relayer/testtools"
	"raylz-relayer/types"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/raylsnetwork/rayls-privacy-relayer-api/listener"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/types"
)

func TestNewLogListener(t *testing.T) {
	testtools.SilenceLogger()

	batchSize := 100
	addresses := []common.Address{}

	t.Run("sets last processed block to starting block when not found in repository", func(t *testing.T) {
		wantBlock := new(big.Int).SetUint64(1337)
		wantComponent := types.LastProcessedBlockDocumentPublicChain

		ethClient := &LogEthereumClientMock{}
		lpbRepo := &LogLastProcessedBlockRepositoryMock{
			GetFunc: func(_ context.Context, component types.LastProcessedBlockDocument) (*big.Int, error) {
				assert.Equal(t, wantComponent, component)
				return nil, listener.ErrLastProcessedBlockNotFound
			},
			CreateFunc: func(_ context.Context, component types.LastProcessedBlockDocument, blockNumber *big.Int) error {
				assert.Equal(t, wantComponent, component)
				assert.Equal(t, wantBlock, blockNumber)
				return nil
			},
		}
		logHandler := &LogHandlerMock{}

		config := listener.LogListenerConfig{
			Component:     wantComponent,
			StartingBlock: wantBlock,
			BatchSize:     batchSize,
			Addresses:     addresses,
		}

		_, err := listener.NewLogListener(config, logHandler, ethClient, lpbRepo)

		require.NoError(t, err)
		require.Len(t, lpbRepo.GetCalls(), 1)
		require.Len(t, lpbRepo.CreateCalls(), 1)
	})

	t.Run("uses existing block from repository when found", func(t *testing.T) {
		existingBlock := new(big.Int).SetUint64(500)

		ethClient := &LogEthereumClientMock{}
		lpbRepo := &LogLastProcessedBlockRepositoryMock{
			GetFunc: func(_ context.Context, _ types.LastProcessedBlockDocument) (*big.Int, error) {
				return existingBlock, nil
			},
		}
		logHandler := &LogHandlerMock{}

		config := listener.LogListenerConfig{
			Component:     types.LastProcessedBlockDocumentPublicChain,
			StartingBlock: big.NewInt(1),
			BatchSize:     batchSize,
			Addresses:     addresses,
		}

		_, err := listener.NewLogListener(config, logHandler, ethClient, lpbRepo)

		require.NoError(t, err)
		require.Len(t, lpbRepo.GetCalls(), 1)
		require.Len(t, lpbRepo.CreateCalls(), 0)
	})

	t.Run("wraps repository Get error in ListenerError", func(t *testing.T) {
		wantErr := errors.New("database connection failed")

		ethClient := &LogEthereumClientMock{}
		lpbRepo := &LogLastProcessedBlockRepositoryMock{
			GetFunc: func(_ context.Context, _ types.LastProcessedBlockDocument) (*big.Int, error) {
				return nil, wantErr
			},
		}
		logHandler := &LogHandlerMock{}

		config := listener.LogListenerConfig{
			Component:     types.DocumentIdLastProcessedBlockPrivacyNode,
			StartingBlock: big.NewInt(1337),
			BatchSize:     batchSize,
			Addresses:     addresses,
		}

		_, gotErr := listener.NewLogListener(config, logHandler, ethClient, lpbRepo)

		require.ErrorAs(t, gotErr, new(*listener.ListenerError))
		require.ErrorIs(t, gotErr, wantErr)
	})

	t.Run("wraps repository Create error in ListenerError", func(t *testing.T) {
		wantErr := errors.New("insert failed")

		ethClient := &LogEthereumClientMock{}
		lpbRepo := &LogLastProcessedBlockRepositoryMock{
			GetFunc: func(_ context.Context, _ types.LastProcessedBlockDocument) (*big.Int, error) {
				return nil, listener.ErrLastProcessedBlockNotFound
			},
			CreateFunc: func(_ context.Context, _ types.LastProcessedBlockDocument, _ *big.Int) error {
				return wantErr
			},
		}
		logHandler := &LogHandlerMock{}

		config := listener.LogListenerConfig{
			Component:     types.DocumentIdLastProcessedBlockPrivacyNode,
			StartingBlock: big.NewInt(1),
			BatchSize:     batchSize,
			Addresses:     addresses,
		}

		_, gotErr := listener.NewLogListener(config, logHandler, ethClient, lpbRepo)

		require.ErrorAs(t, gotErr, new(*listener.ListenerError))
		require.ErrorIs(t, gotErr, wantErr)
	})
}

func TestRun(t *testing.T) {
	testtools.SilenceLogger()

	t.Run("supports graceful shutdown", func(t *testing.T) {
		ethClient := &LogEthereumClientMock{
			BlockNumberFunc: func(ctx context.Context) (uint64, error) {
				<-ctx.Done()
				return 0, ctx.Err()
			},
		}
		lpbRepo := &LogLastProcessedBlockRepositoryMock{
			GetFunc: func(_ context.Context, _ types.LastProcessedBlockDocument) (*big.Int, error) {
				return big.NewInt(1), nil
			},
		}
		logHandler := &LogHandlerMock{}

		config := listener.LogListenerConfig{
			Component:     types.DocumentIdLastProcessedBlockPrivateHub,
			StartingBlock: big.NewInt(1),
			BatchSize:     10,
			Addresses:     []common.Address{common.HexToAddress("0x1")},
		}

		logListener, err := listener.NewLogListener(config, logHandler, ethClient, lpbRepo)
		require.NoError(t, err)

		hasGracefulShutdown := testtools.ShutdownFixture(t, logListener.Run, time.Second)
		assert.True(t, hasGracefulShutdown)
	})

	t.Run("continues polling when BlockNumber returns error", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())

		blockNumberCallCount := 0
		ethClient := &LogEthereumClientMock{
			BlockNumberFunc: func(_ context.Context) (uint64, error) {
				blockNumberCallCount++
				if blockNumberCallCount >= 2 {
					cancel()
				}
				return 0, errors.New("rpc error")
			},
		}
		lpbRepo := &LogLastProcessedBlockRepositoryMock{
			GetFunc: func(_ context.Context, _ types.LastProcessedBlockDocument) (*big.Int, error) {
				return big.NewInt(1), nil
			},
		}
		logHandler := &LogHandlerMock{}

		config := listener.LogListenerConfig{
			Component:     types.DocumentIdLastProcessedBlockPrivateHub,
			StartingBlock: big.NewInt(1),
			BatchSize:     10,
			Addresses:     []common.Address{common.HexToAddress("0x1")},
		}

		logListener, err := listener.NewLogListener(config, logHandler, ethClient, lpbRepo)
		require.NoError(t, err)

		err = logListener.Run(ctx)

		assert.NoError(t, err)
		assert.GreaterOrEqual(t, blockNumberCallCount, 2)
		assert.Len(t, logHandler.HandleCalls(), 0)
	})

	t.Run("skips processing when no new blocks exist", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())

		// currentBlock = 100, latestBlock = 50 → no new blocks
		ethClient := &LogEthereumClientMock{
			BlockNumberFunc: func(_ context.Context) (uint64, error) {
				cancel()
				return 50, nil
			},
		}
		lpbRepo := &LogLastProcessedBlockRepositoryMock{
			GetFunc: func(_ context.Context, _ types.LastProcessedBlockDocument) (*big.Int, error) {
				return big.NewInt(100), nil
			},
		}
		logHandler := &LogHandlerMock{}

		config := listener.LogListenerConfig{
			Component:     types.DocumentIdLastProcessedBlockPrivateHub,
			StartingBlock: big.NewInt(100),
			BatchSize:     10,
			Addresses:     []common.Address{common.HexToAddress("0x1")},
		}

		logListener, err := listener.NewLogListener(config, logHandler, ethClient, lpbRepo)
		require.NoError(t, err)

		err = logListener.Run(ctx)

		assert.NoError(t, err)
		assert.Len(t, logHandler.HandleCalls(), 0)
	})
}

func TestRun_NearTipMode(t *testing.T) {
	testtools.SilenceLogger()

	t.Run("processes single block with one log", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())

		currentBlock := big.NewInt(100)
		latestBlock := uint64(100) // same block = near tip

		wantComponent := types.DocumentIdLastProcessedBlockPrivateHub
		wantAddress := common.HexToAddress("0x1")

		ethClient := &LogEthereumClientMock{
			BlockNumberFunc: func(_ context.Context) (uint64, error) {
				return latestBlock, nil
			},
			FilterLogsFunc: func(_ context.Context, q ethereum.FilterQuery) ([]ethTypes.Log, error) {
				return []ethTypes.Log{
					{BlockNumber: q.FromBlock.Uint64(), Address: wantAddress},
				}, nil
			},
		}
		lpbRepo := &LogLastProcessedBlockRepositoryMock{
			GetFunc: func(_ context.Context, _ types.LastProcessedBlockDocument) (*big.Int, error) {
				return currentBlock, nil
			},
			UpdateFunc: func(_ context.Context, component types.LastProcessedBlockDocument, blockNumber *big.Int) error {
				assert.Equal(t, wantComponent, component)
				assert.Equal(t, big.NewInt(101), blockNumber)
				return nil
			},
		}
		logHandler := &LogHandlerMock{
			HandleFunc: func(_ context.Context, logs []ethTypes.Log) error {
				require.Len(t, logs, 1)
				assert.Equal(t, uint64(100), logs[0].BlockNumber)
				cancel()
				return nil
			},
		}

		config := listener.LogListenerConfig{
			Component:     wantComponent,
			StartingBlock: currentBlock,
			BatchSize:     10,
			Addresses:     []common.Address{wantAddress},
		}

		logListener, err := listener.NewLogListener(config, logHandler, ethClient, lpbRepo)
		require.NoError(t, err)

		err = logListener.Run(ctx)

		assert.NoError(t, err)
		assert.Len(t, logHandler.HandleCalls(), 1)
		assert.Len(t, lpbRepo.UpdateCalls(), 1)
	})

	t.Run("processes single block with multiple logs", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())

		currentBlock := big.NewInt(100)
		latestBlock := uint64(100)

		ethClient := &LogEthereumClientMock{
			BlockNumberFunc: func(_ context.Context) (uint64, error) {
				return latestBlock, nil
			},
			FilterLogsFunc: func(_ context.Context, q ethereum.FilterQuery) ([]ethTypes.Log, error) {
				logs := make([]ethTypes.Log, 10)
				for i := range logs {
					logs[i] = ethTypes.Log{BlockNumber: q.FromBlock.Uint64(), Index: uint(i)}
				}
				return logs, nil
			},
		}
		lpbRepo := &LogLastProcessedBlockRepositoryMock{
			GetFunc: func(_ context.Context, _ types.LastProcessedBlockDocument) (*big.Int, error) {
				return currentBlock, nil
			},
			UpdateFunc: func(_ context.Context, _ types.LastProcessedBlockDocument, blockNumber *big.Int) error {
				assert.Equal(t, big.NewInt(101), blockNumber)
				return nil
			},
		}
		logHandler := &LogHandlerMock{
			HandleFunc: func(_ context.Context, logs []ethTypes.Log) error {
				require.Len(t, logs, 10)
				cancel()
				return nil
			},
		}

		config := listener.LogListenerConfig{
			Component:     types.DocumentIdLastProcessedBlockPrivateHub,
			StartingBlock: currentBlock,
			BatchSize:     10,
			Addresses:     []common.Address{common.HexToAddress("0x1")},
		}

		logListener, err := listener.NewLogListener(config, logHandler, ethClient, lpbRepo)
		require.NoError(t, err)

		err = logListener.Run(ctx)

		assert.NoError(t, err)
		assert.Len(t, logHandler.HandleCalls(), 1)
		assert.Len(t, lpbRepo.UpdateCalls(), 1)
	})

	t.Run("processes multiple blocks one by one near tip", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())

		currentBlock := big.NewInt(98)
		latestBlock := uint64(100) // 3 blocks: 98, 99, 100

		ethClient := &LogEthereumClientMock{
			BlockNumberFunc: func(_ context.Context) (uint64, error) {
				return latestBlock, nil
			},
			FilterLogsFunc: func(_ context.Context, q ethereum.FilterQuery) ([]ethTypes.Log, error) {
				return []ethTypes.Log{
					{BlockNumber: q.FromBlock.Uint64()},
				}, nil
			},
		}

		updateCallCount := 0
		lpbRepo := &LogLastProcessedBlockRepositoryMock{
			GetFunc: func(_ context.Context, _ types.LastProcessedBlockDocument) (*big.Int, error) {
				return currentBlock, nil
			},
			UpdateFunc: func(_ context.Context, _ types.LastProcessedBlockDocument, blockNumber *big.Int) error {
				updateCallCount++
				// Expect updates: 99, 100, 101
				expectedBlock := big.NewInt(int64(98 + updateCallCount))
				assert.Equal(t, expectedBlock, blockNumber)
				return nil
			},
		}
		handleCallCount := 0
		logHandler := &LogHandlerMock{
			HandleFunc: func(_ context.Context, logs []ethTypes.Log) error {
				handleCallCount++
				require.Len(t, logs, 1)
				if handleCallCount == 3 {
					cancel()
				}
				return nil
			},
		}

		config := listener.LogListenerConfig{
			Component:     types.DocumentIdLastProcessedBlockPrivateHub,
			StartingBlock: currentBlock,
			BatchSize:     10,
			Addresses:     []common.Address{common.HexToAddress("0x1")},
		}

		logListener, err := listener.NewLogListener(config, logHandler, ethClient, lpbRepo)
		require.NoError(t, err)

		err = logListener.Run(ctx)

		assert.NoError(t, err)
		assert.Len(t, logHandler.HandleCalls(), 3)
		assert.Len(t, lpbRepo.UpdateCalls(), 3)
	})

	t.Run("skips handler call when block has no logs", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())

		currentBlock := big.NewInt(100)
		latestBlock := uint64(100)

		ethClient := &LogEthereumClientMock{
			BlockNumberFunc: func(_ context.Context) (uint64, error) {
				return latestBlock, nil
			},
			FilterLogsFunc: func(_ context.Context, _ ethereum.FilterQuery) ([]ethTypes.Log, error) {
				return []ethTypes.Log{}, nil
			},
			HeaderByNumberFunc: func(_ context.Context, _ *big.Int) (*ethTypes.Header, error) {
				return &ethTypes.Header{Bloom: ethTypes.Bloom{}}, nil
			},
		}
		lpbRepo := &LogLastProcessedBlockRepositoryMock{
			GetFunc: func(_ context.Context, _ types.LastProcessedBlockDocument) (*big.Int, error) {
				return currentBlock, nil
			},
			UpdateFunc: func(_ context.Context, _ types.LastProcessedBlockDocument, blockNumber *big.Int) error {
				assert.Equal(t, big.NewInt(101), blockNumber)
				cancel()
				return nil
			},
		}
		logHandler := &LogHandlerMock{}

		config := listener.LogListenerConfig{
			Component:     types.DocumentIdLastProcessedBlockPrivateHub,
			StartingBlock: currentBlock,
			BatchSize:     10,
			Addresses:     []common.Address{common.HexToAddress("0x1")},
		}

		logListener, err := listener.NewLogListener(config, logHandler, ethClient, lpbRepo)
		require.NoError(t, err)

		err = logListener.Run(ctx)

		assert.NoError(t, err)
		assert.Len(t, logHandler.HandleCalls(), 0)
		assert.Len(t, lpbRepo.UpdateCalls(), 1)
	})

	t.Run("retries on FilterLogs error", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())

		currentBlock := big.NewInt(100)
		latestBlock := uint64(102) // within 10 blocks = near tip

		filterLogsCallCount := 0
		ethClient := &LogEthereumClientMock{
			BlockNumberFunc: func(_ context.Context) (uint64, error) {
				return latestBlock, nil
			},
			FilterLogsFunc: func(_ context.Context, q ethereum.FilterQuery) ([]ethTypes.Log, error) {
				filterLogsCallCount++
				// Fail on first call for block 100, succeed on retries
				if filterLogsCallCount == 1 && q.FromBlock.Uint64() == 100 {
					return nil, errors.New("transient network error")
				}
				return []ethTypes.Log{{BlockNumber: q.FromBlock.Uint64()}}, nil
			},
		}
		updateCallCount := 0
		lpbRepo := &LogLastProcessedBlockRepositoryMock{
			GetFunc: func(_ context.Context, _ types.LastProcessedBlockDocument) (*big.Int, error) {
				return currentBlock, nil
			},
			UpdateFunc: func(_ context.Context, _ types.LastProcessedBlockDocument, blockNumber *big.Int) error {
				updateCallCount++
				if updateCallCount == 1 {
					assert.Equal(t, big.NewInt(101), blockNumber)
					cancel()
				}
				return nil
			},
		}
		handleCallCount := 0
		logHandler := &LogHandlerMock{
			HandleFunc: func(_ context.Context, logs []ethTypes.Log) error {
				handleCallCount++
				require.Len(t, logs, 1)
				return nil
			},
		}

		config := listener.LogListenerConfig{
			Component:     types.LastProcessedBlockDocumentPublicChain,
			StartingBlock: currentBlock,
			BatchSize:     10,
			Addresses:     []common.Address{common.HexToAddress("0x1")},
		}

		logListener, err := listener.NewLogListener(config, logHandler, ethClient, lpbRepo)
		require.NoError(t, err)

		err = logListener.Run(ctx)

		assert.NoError(t, err)
		assert.GreaterOrEqual(t, filterLogsCallCount, 2)
		assert.GreaterOrEqual(t, handleCallCount, 1)
		assert.GreaterOrEqual(t, updateCallCount, 1)
	})

	t.Run("retries on handler error", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())

		currentBlock := big.NewInt(100)
		latestBlock := uint64(102)

		ethClient := &LogEthereumClientMock{
			BlockNumberFunc: func(_ context.Context) (uint64, error) {
				return latestBlock, nil
			},
			FilterLogsFunc: func(_ context.Context, q ethereum.FilterQuery) ([]ethTypes.Log, error) {
				return []ethTypes.Log{{BlockNumber: q.FromBlock.Uint64()}}, nil
			},
		}
		updateCallCount := 0
		lpbRepo := &LogLastProcessedBlockRepositoryMock{
			GetFunc: func(_ context.Context, _ types.LastProcessedBlockDocument) (*big.Int, error) {
				return currentBlock, nil
			},
			UpdateFunc: func(_ context.Context, _ types.LastProcessedBlockDocument, blockNumber *big.Int) error {
				updateCallCount++
				if updateCallCount == 1 {
					assert.Equal(t, big.NewInt(101), blockNumber)
					cancel()
				}
				return nil
			},
		}
		handleCallCount := 0
		logHandler := &LogHandlerMock{
			HandleFunc: func(_ context.Context, logs []ethTypes.Log) error {
				handleCallCount++
				require.Len(t, logs, 1)
				// Fail on first call for block 100, succeed on retries
				if handleCallCount == 1 && logs[0].BlockNumber == 100 {
					return errors.New("transient handler error")
				}
				return nil
			},
		}

		config := listener.LogListenerConfig{
			Component:     types.LastProcessedBlockDocumentPublicChain,
			StartingBlock: currentBlock,
			BatchSize:     10,
			Addresses:     []common.Address{common.HexToAddress("0x1")},
		}

		logListener, err := listener.NewLogListener(config, logHandler, ethClient, lpbRepo)
		require.NoError(t, err)

		err = logListener.Run(ctx)

		assert.NoError(t, err)
		assert.GreaterOrEqual(t, handleCallCount, 2)
		assert.GreaterOrEqual(t, updateCallCount, 1)
	})
}

func TestRun_BatchMode(t *testing.T) {
	testtools.SilenceLogger()

	t.Run("processes batch of blocks far from tip", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())

		currentBlock := big.NewInt(100)
		latestBlock := uint64(200) // >10 blocks away = batch mode

		ethClient := &LogEthereumClientMock{
			BlockNumberFunc: func(_ context.Context) (uint64, error) {
				return latestBlock, nil
			},
			FilterLogsFunc: func(_ context.Context, q ethereum.FilterQuery) ([]ethTypes.Log, error) {
				// Return one log per FilterLogs call
				return []ethTypes.Log{{BlockNumber: q.FromBlock.Uint64()}}, nil
			},
		}
		lpbRepo := &LogLastProcessedBlockRepositoryMock{
			GetFunc: func(_ context.Context, _ types.LastProcessedBlockDocument) (*big.Int, error) {
				return currentBlock, nil
			},
			UpdateFunc: func(_ context.Context, _ types.LastProcessedBlockDocument, blockNumber *big.Int) error {
				// First batch: from=100, to=109 (batchSize=10, inclusive range),
				// update to next unprocessed = 110
				assert.Equal(t, big.NewInt(110), blockNumber)
				cancel()
				return nil
			},
		}
		logHandler := &LogHandlerMock{
			HandleFunc: func(_ context.Context, logs []ethTypes.Log) error {
				require.GreaterOrEqual(t, len(logs), 1)
				return nil
			},
		}

		config := listener.LogListenerConfig{
			Component:     types.DocumentIdLastProcessedBlockPrivateHub,
			StartingBlock: currentBlock,
			BatchSize:     10,
			Addresses:     []common.Address{common.HexToAddress("0x1")},
		}

		logListener, err := listener.NewLogListener(config, logHandler, ethClient, lpbRepo)
		require.NoError(t, err)

		err = logListener.Run(ctx)

		assert.NoError(t, err)
		assert.Len(t, logHandler.HandleCalls(), 1)
		assert.Len(t, lpbRepo.UpdateCalls(), 1)

		// Verify the filter query used the correct range
		filterCalls := ethClient.FilterLogsCalls()
		require.Len(t, filterCalls, 1)
		assert.Equal(t, big.NewInt(100), filterCalls[0].FilterQuery.FromBlock)
		assert.Equal(t, big.NewInt(109), filterCalls[0].FilterQuery.ToBlock)
	})

	t.Run("skips handler call when batch has no logs", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())

		currentBlock := big.NewInt(100)
		latestBlock := uint64(200)

		ethClient := &LogEthereumClientMock{
			BlockNumberFunc: func(_ context.Context) (uint64, error) {
				return latestBlock, nil
			},
			FilterLogsFunc: func(_ context.Context, _ ethereum.FilterQuery) ([]ethTypes.Log, error) {
				return []ethTypes.Log{}, nil
			},
		}
		lpbRepo := &LogLastProcessedBlockRepositoryMock{
			GetFunc: func(_ context.Context, _ types.LastProcessedBlockDocument) (*big.Int, error) {
				return currentBlock, nil
			},
			UpdateFunc: func(_ context.Context, _ types.LastProcessedBlockDocument, blockNumber *big.Int) error {
				assert.Equal(t, big.NewInt(110), blockNumber)
				cancel()
				return nil
			},
		}
		logHandler := &LogHandlerMock{}

		config := listener.LogListenerConfig{
			Component:     types.DocumentIdLastProcessedBlockPrivateHub,
			StartingBlock: currentBlock,
			BatchSize:     10,
			Addresses:     []common.Address{common.HexToAddress("0x1")},
		}

		logListener, err := listener.NewLogListener(config, logHandler, ethClient, lpbRepo)
		require.NoError(t, err)

		err = logListener.Run(ctx)

		assert.NoError(t, err)
		assert.Len(t, logHandler.HandleCalls(), 0)
		assert.Len(t, lpbRepo.UpdateCalls(), 1)
	})

	t.Run("retries on FilterLogs error", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())

		currentBlock := big.NewInt(100)
		latestBlock := uint64(200)

		filterLogsCallCount := 0
		ethClient := &LogEthereumClientMock{
			BlockNumberFunc: func(_ context.Context) (uint64, error) {
				return latestBlock, nil
			},
			FilterLogsFunc: func(_ context.Context, _ ethereum.FilterQuery) ([]ethTypes.Log, error) {
				filterLogsCallCount++
				if filterLogsCallCount == 1 {
					return nil, errors.New("transient network error")
				}
				return []ethTypes.Log{{BlockNumber: 100}}, nil
			},
		}
		lpbRepo := &LogLastProcessedBlockRepositoryMock{
			GetFunc: func(_ context.Context, _ types.LastProcessedBlockDocument) (*big.Int, error) {
				return currentBlock, nil
			},
			UpdateFunc: func(_ context.Context, _ types.LastProcessedBlockDocument, blockNumber *big.Int) error {
				assert.Equal(t, big.NewInt(110), blockNumber)
				cancel()
				return nil
			},
		}
		handleCallCount := 0
		logHandler := &LogHandlerMock{
			HandleFunc: func(_ context.Context, logs []ethTypes.Log) error {
				handleCallCount++
				require.GreaterOrEqual(t, len(logs), 1)
				return nil
			},
		}

		config := listener.LogListenerConfig{
			Component:     types.LastProcessedBlockDocumentPublicChain,
			StartingBlock: currentBlock,
			BatchSize:     10,
			Addresses:     []common.Address{common.HexToAddress("0x1")},
		}

		logListener, err := listener.NewLogListener(config, logHandler, ethClient, lpbRepo)
		require.NoError(t, err)

		err = logListener.Run(ctx)

		assert.NoError(t, err)
		assert.Equal(t, 2, filterLogsCallCount)
		assert.Equal(t, 1, handleCallCount)
		assert.Equal(t, 1, len(lpbRepo.UpdateCalls()))
	})

	t.Run("retries on handler error", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())

		currentBlock := big.NewInt(100)
		latestBlock := uint64(200)

		ethClient := &LogEthereumClientMock{
			BlockNumberFunc: func(_ context.Context) (uint64, error) {
				return latestBlock, nil
			},
			FilterLogsFunc: func(_ context.Context, _ ethereum.FilterQuery) ([]ethTypes.Log, error) {
				return []ethTypes.Log{{BlockNumber: 100}}, nil
			},
		}
		lpbRepo := &LogLastProcessedBlockRepositoryMock{
			GetFunc: func(_ context.Context, _ types.LastProcessedBlockDocument) (*big.Int, error) {
				return currentBlock, nil
			},
			UpdateFunc: func(_ context.Context, _ types.LastProcessedBlockDocument, blockNumber *big.Int) error {
				assert.Equal(t, big.NewInt(110), blockNumber)
				cancel()
				return nil
			},
		}
		handleCallCount := 0
		logHandler := &LogHandlerMock{
			HandleFunc: func(_ context.Context, logs []ethTypes.Log) error {
				handleCallCount++
				require.GreaterOrEqual(t, len(logs), 1)
				if handleCallCount == 1 {
					return errors.New("transient handler error")
				}
				return nil
			},
		}

		config := listener.LogListenerConfig{
			Component:     types.LastProcessedBlockDocumentPublicChain,
			StartingBlock: currentBlock,
			BatchSize:     10,
			Addresses:     []common.Address{common.HexToAddress("0x1")},
		}

		logListener, err := listener.NewLogListener(config, logHandler, ethClient, lpbRepo)
		require.NoError(t, err)

		err = logListener.Run(ctx)

		assert.NoError(t, err)
		assert.Equal(t, 2, handleCallCount)
		assert.Equal(t, 1, len(lpbRepo.UpdateCalls()))
	})
}

func TestRun_BatchingBehavior(t *testing.T) {
	testtools.SilenceLogger()

	/*
		These tests verify the listener's batching logic. The listener has two modes:
		- Batch mode: when far from tip (>10 blocks), processes blocks in batches
		- Near-tip mode: when within 10 blocks of tip, processes one block at a time

		Each test case specifies totalBlocks and batchSize, then verifies the expected
		sequence of batch sizes delivered to the handler. The trailing ones represent
		the near-tip single-block processing.
	*/
	testCases := []struct {
		name            string
		batchSize       int
		totalBlocks     int
		expectedBatches []int
	}{
		{
			name:            "batch size 10, total blocks 20: one full batch then 10 near-tip",
			batchSize:       10,
			totalBlocks:     20,
			expectedBatches: []int{10, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
		},
		{
			name:            "batch size 10, total blocks 9: all near-tip",
			batchSize:       10,
			totalBlocks:     9,
			expectedBatches: []int{1, 1, 1, 1, 1, 1, 1, 1, 1},
		},
		{
			name:            "batch size 10, total blocks 11: all near-tip",
			batchSize:       10,
			totalBlocks:     11,
			expectedBatches: []int{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
		},
		{
			name:            "batch size 10, total blocks 1: single near-tip block",
			batchSize:       10,
			totalBlocks:     1,
			expectedBatches: []int{1},
		},
		{
			name:            "batch size 1, total blocks 10: all single blocks",
			batchSize:       1,
			totalBlocks:     10,
			expectedBatches: []int{1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
		},
		{
			name:            "batch size 500, total blocks 5030: ten full batches + partial + near-tip",
			batchSize:       500,
			totalBlocks:     5030,
			expectedBatches: []int{500, 500, 500, 500, 500, 500, 500, 500, 500, 500, 20, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
		},
		{
			name:            "batch size 143, total blocks 1471: uneven batches + partial + near-tip",
			batchSize:       143,
			totalBlocks:     1471,
			expectedBatches: []int{143, 143, 143, 143, 143, 143, 143, 143, 143, 143, 31, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
		},
		{
			name:            "batch size 10, total blocks 55: normal case with partial batch",
			batchSize:       10,
			totalBlocks:     55,
			expectedBatches: []int{10, 10, 10, 10, 5, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
			defer cancel()

			currentBlock := big.NewInt(1)
			latestBlock := uint64(tc.totalBlocks)

			// Track which block ranges have been queried to return correct log counts
			ethClient := &LogEthereumClientMock{
				BlockNumberFunc: func(_ context.Context) (uint64, error) {
					return latestBlock, nil
				},
				FilterLogsFunc: func(_ context.Context, q ethereum.FilterQuery) ([]ethTypes.Log, error) {
					from := q.FromBlock.Uint64()
					to := q.ToBlock.Uint64()
					count := int(to - from + 1)
					logs := make([]ethTypes.Log, count)
					for i := range count {
						logs[i] = ethTypes.Log{BlockNumber: from + uint64(i)}
					}
					return logs, nil
				},
			}

			wantBlockNumber := new(big.Int).Set(currentBlock)
			updateCallCount := 0
			var lpbRepo *LogLastProcessedBlockRepositoryMock
			lpbRepo = &LogLastProcessedBlockRepositoryMock{
				GetFunc: func(_ context.Context, _ types.LastProcessedBlockDocument) (*big.Int, error) {
					return currentBlock, nil
				},
				CreateFunc: func(_ context.Context, _ types.LastProcessedBlockDocument, _ *big.Int) error {
					return nil
				},
				UpdateFunc: func(_ context.Context, _ types.LastProcessedBlockDocument, blockNumber *big.Int) error {
					batchIdx := len(lpbRepo.UpdateCalls()) - 1
					wantBlockNumber = new(big.Int).Add(wantBlockNumber, big.NewInt(int64(tc.expectedBatches[batchIdx])))
					assert.Equal(t, wantBlockNumber, blockNumber, "batch %d: wrong last processed block", batchIdx)
					updateCallCount++
					return nil
				},
			}
			handleCallCount := 0
			var logHandler *LogHandlerMock
			logHandler = &LogHandlerMock{
				HandleFunc: func(_ context.Context, logs []ethTypes.Log) error {
					batchIdx := len(logHandler.HandleCalls()) - 1
					require.Equal(t, tc.expectedBatches[batchIdx], len(logs), "batch %d: wrong log count", batchIdx)
					handleCallCount++

					if handleCallCount == len(tc.expectedBatches) {
						cancel()
					}
					return nil
				},
			}

			config := listener.LogListenerConfig{
				Component:     types.DocumentIdLastProcessedBlockPrivateHub,
				StartingBlock: currentBlock,
				BatchSize:     tc.batchSize,
				Addresses:     []common.Address{common.HexToAddress("0x1")},
			}

			logListener, err := listener.NewLogListener(config, logHandler, ethClient, lpbRepo)
			require.NoError(t, err)

			err = logListener.Run(ctx)

			assert.NoError(t, err)
			assert.Equal(t, len(tc.expectedBatches), handleCallCount)
		})
	}
}
