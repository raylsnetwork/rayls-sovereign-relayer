package listener

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/raylsnetwork/rayls-sovereign-relayer/types"
)

const logListenerTimeout = 30 * time.Second

type LogHandlerFunc func(ctx context.Context, logs []ethTypes.Log) error

func (f LogHandlerFunc) Handle(ctx context.Context, logs []ethTypes.Log) error {
	return f(ctx, logs)
}

type LogHandler interface {
	Handle(context.Context, []ethTypes.Log) error
}

//go:generate moq --pkg listener_test -out log_mock_test.go . LogEthereumClient LogLastProcessedBlockRepository LogHandler
type LogEthereumClient interface {
	BlockNumber(context.Context) (uint64, error)
	FilterLogs(context.Context, ethereum.FilterQuery) ([]ethTypes.Log, error)
	HeaderByNumber(context.Context, *big.Int) (*ethTypes.Header, error)
}

type LogLastProcessedBlockRepository interface {
	Get(context.Context, types.LastProcessedBlockDocument) (*big.Int, error)
	Create(context.Context, types.LastProcessedBlockDocument, *big.Int) error
	Update(context.Context, types.LastProcessedBlockDocument, *big.Int) error
}

type LogListenerConfig struct {
	Component types.LastProcessedBlockDocument

	StartingBlock *big.Int
	BatchSize     int

	Addresses []common.Address
}

type LogListener struct {
	component    types.LastProcessedBlockDocument
	currentBlock *big.Int
	batchSize    *big.Int
	addresses    []common.Address

	handler       LogHandler
	ledgerClient  LogEthereumClient
	lpbRepository LogLastProcessedBlockRepository

	ticker *time.Ticker
}

func NewLogListener(
	ctx context.Context,
	config LogListenerConfig,
	handler LogHandler,
	ledgerClient LogEthereumClient,
	lpbRepository LogLastProcessedBlockRepository,
) (*LogListener, error) {
	currentBlock, err := lpbRepository.Get(ctx, config.Component)
	if err != nil {
		if errors.Is(err, ErrLastProcessedBlockNotFound) {
			currentBlock = config.StartingBlock

			err = lpbRepository.Create(ctx, config.Component, currentBlock)
			if err != nil {
				return nil, WrapInListenerError("failed to set initial block for component", err)
			}
		} else {
			return nil, WrapInListenerError("failed to get current block for component", err)
		}
	}

	return &LogListener{
		component:    config.Component,
		currentBlock: currentBlock,
		batchSize:    big.NewInt(int64(config.BatchSize)),
		addresses:    config.Addresses,

		handler:       handler,
		ledgerClient:  ledgerClient,
		lpbRepository: lpbRepository,

		ticker: time.NewTicker(time.Second),
	}, nil
}

func (l *LogListener) Run(ctx context.Context) error {
	initialRun := make(chan struct{}, 1)
	initialRun <- struct{}{}

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-l.ticker.C:
		case <-initialRun:
		}

		timeoutCtx, timeoutCancel := context.WithTimeout(ctx, logListenerTimeout)
		latestBlockUint64, err := l.ledgerClient.BlockNumber(timeoutCtx)
		timeoutCancel() // Clean up immediately to prevent leak

		if err != nil {
			slog.Error("failed to get latest block number", slog.Any("error", err))
			continue // retry again
		}
		latestBlock := new(big.Int).SetUint64(latestBlockUint64)

		// Check that new blocks have been mined
		if l.currentBlock.Cmp(latestBlock) > 0 {
			continue // retry again
		}

		// Check if we're near the tip (within 10 blocks)
		tenBlocksFromTip := new(big.Int).Sub(latestBlock, big.NewInt(10))
		nearTip := l.currentBlock.Cmp(tenBlocksFromTip) >= 0

		if nearTip {
			// Process blocks one by one near the tip
			for blockNum := new(big.Int).Set(l.currentBlock); blockNum.Cmp(latestBlock) <= 0; blockNum.Add(blockNum, big.NewInt(1)) {
				slog.Debug(fmt.Sprintf("Processing %s block (near tip)", l.component), slog.Any("Block", blockNum))

				logs, err := l.processBlockWithBloomCheck(ctx, blockNum)
				if err != nil {
					slog.Error("failed to process block with bloom check", slog.Any("error", err))
					break // we break out of the the inner loop an retry in the outer loop
				}

				// Process and update for this single block
				if err := l.handleLogsAndUpdate(ctx, logs, blockNum); err != nil {
					slog.Error("failed to handle logs and update", slog.Any("error", err))
					break // we break out of the the inner loop an retry in the outer loop
				}
			}
		} else {
			// Batch process blocks far from tip
			fromBlock := l.currentBlock
			// As the FilterLogs function uses an inclusive range, we subtract
			// one from the batch size to get a "batch size" slice of logs.
			toBlock := new(big.Int).Add(fromBlock, new(big.Int).Sub(l.batchSize, big.NewInt(1)))

			// Limit the batch to the latest block
			if toBlock.Cmp(tenBlocksFromTip) > 0 {
				toBlock = tenBlocksFromTip
			}

			slog.Debug(
				fmt.Sprintf("Processing %s blocks (batch)", l.component),
				slog.Any("From", fromBlock),
				slog.Any("To", toBlock),
			)

			timeoutCtx, timeoutCancel := context.WithTimeout(ctx, logListenerTimeout)

			query := ethereum.FilterQuery{
				FromBlock: fromBlock,
				ToBlock:   toBlock,
				Addresses: l.addresses,
			}

			logs, err := l.ledgerClient.FilterLogs(timeoutCtx, query)
			timeoutCancel() // Clean up immediately to prevent leak

			if err != nil {
				slog.Error("failed to filter logs", slog.Any("error", err))
				continue // continue here to retry from the beginning so we do not stop the listener
			}

			// Process and update for the batch
			if err := l.handleLogsAndUpdate(ctx, logs, toBlock); err != nil {
				slog.Error("failed to handle logs and update", slog.Any("error", err))
				continue // continue here to retry from the beginning so we do not stop the listener
			}
		}
	}
}

func (l *LogListener) handleLogsAndUpdate(ctx context.Context, logs []ethTypes.Log, processedUpTo *big.Int) error {
	if len(logs) > 0 {
		if err := l.handler.Handle(ctx, logs); err != nil {
			return fmt.Errorf("failed to handle logs: %w", err)
		}
	}

	l.currentBlock = new(big.Int).Add(processedUpTo, big.NewInt(1))
	if err := l.lpbRepository.Update(ctx, l.component, l.currentBlock); err != nil {
		return fmt.Errorf("failed to update last processed block: %w", err)
	}

	return nil
}

func (l *LogListener) processBlockWithBloomCheck(ctx context.Context, blockNum *big.Int) ([]ethTypes.Log, error) {
	timeoutCtx, timeoutCancel := context.WithTimeout(ctx, logListenerTimeout)

	query := ethereum.FilterQuery{
		FromBlock: blockNum,
		ToBlock:   blockNum,
		Addresses: l.addresses,
	}

	logs, err := l.ledgerClient.FilterLogs(timeoutCtx, query)
	timeoutCancel() // Clean up immediately to prevent leak

	if err != nil {
		return nil, fmt.Errorf("failed to filter logs for block %s: %w", blockNum, err)
	}

	// If we got logs, return them
	if len(logs) > 0 {
		return logs, nil
	}

	// No logs, check bloom filter
	var header *ethTypes.Header
	for {
		timeoutCtx, timeoutCancel = context.WithTimeout(ctx, 10*time.Second)
		header, err = l.ledgerClient.HeaderByNumber(timeoutCtx, blockNum)
		timeoutCancel()

		if err == nil {
			break
		}

		slog.Error("Failed to get header, retrying", slog.Any("block", blockNum), slog.Any("error", err))
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-time.After(2 * time.Second):
		}
	}

	// Check if bloom indicates logs
	bloomHasLogs := false
	for _, addr := range l.addresses {
		if ethTypes.BloomLookup(header.Bloom, addr) {
			bloomHasLogs = true
			break
		}
	}

	if !bloomHasLogs {
		return []ethTypes.Log{}, nil
	}

	// Bloom says there should be logs, retry
	slog.Info("Bloom indicates logs but got none, retrying", slog.Any("block", blockNum))
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case <-time.After(2 * time.Second):
	}

	timeoutCtx, timeoutCancel = context.WithTimeout(ctx, logListenerTimeout)
	logs, err = l.ledgerClient.FilterLogs(timeoutCtx, query)
	timeoutCancel() // Clean up immediately to prevent leak

	if err != nil {
		return nil, fmt.Errorf("failed to filter logs on retry: %w", err)
	}

	if len(logs) == 0 {
		slog.Warn("Still no logs after retry despite bloom", slog.Any("block", blockNum))
	}

	return logs, nil
}
