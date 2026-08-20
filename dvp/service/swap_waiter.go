package service

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/raylsnetwork/rayls-sovereign-relayer/types"
)

type waiterSwapRepository interface {
	GetSwapBySharedID(ctx context.Context, sharedId string) (*types.DvpSwap, error)
}

type SwapWaiter struct {
	config     WaitConfig
	repository waiterSwapRepository
}

func NewSwapWaiter(config WaitConfig, repository waiterSwapRepository) *SwapWaiter {
	return &SwapWaiter{
		config:     config,
		repository: repository,
	}
}

// WaitForSwapInitiation polls the swap repository until the swap is found or max retries are exceeded.
// It returns the swap if found, or an error if the swap was not found after waiting.
// This is used when a swap was initiated by the other side and we're waiting for it to appear in our DB.
func (w *SwapWaiter) WaitForSwapInitiation(ctx context.Context, sharedId string) (*types.DvpSwap, error) {
	slog.Debug(
		"Swap was initiated by the other side. Start polling our DB for the swap",
		slog.String("sharedId", sharedId),
	)

	retryCount := 0
	for retryCount < w.config.MaxRetries {
		swap, err := w.repository.GetSwapBySharedID(ctx, sharedId)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch swap from DB: %w", err)
		}

		if swap != nil {
			slog.Debug("Swap was found in our DB. Returning it...", slog.String("sharedId", sharedId))
			return swap, nil
		}

		retryCount++
		if retryCount < w.config.MaxRetries {
			slog.Debug(
				"Swap not found in our DB. Retrying...",
				slog.String("sharedId", sharedId),
				slog.Int("attempt", retryCount),
			)
			time.Sleep(w.config.RetryInterval)
		}
	}

	return nil, fmt.Errorf(
		"swap was initiated by the other side but not found in our DB after %d retries",
		w.config.MaxRetries,
	)
}
