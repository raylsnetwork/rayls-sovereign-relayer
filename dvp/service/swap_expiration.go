package service

import (
	"context"
	"log/slog"
	"time"

	"github.com/raylsnetwork/rayls-sovereign-relayer/types"
)

//go:generate moq --pkg service_test -out swap_expiration_mock_test.go . ExpirationSwapRepository ExpirationDvpClient

type ExpirationConfig struct {
	Interval time.Duration
}

type ExpirationSwapRepository interface {
	GetPendingSwaps(ctx context.Context) ([]*types.DvpSwap, error)
}

type ExpirationDvpClient interface {
	IsSwapExpired(ctx context.Context, sharedId string) (bool, error)
	ExpireSwap(ctx context.Context, sharedId string) error
}

type ExpirationService struct {
	ticker     *time.Ticker
	repository ExpirationSwapRepository
	dvpClient  ExpirationDvpClient
}

func NewExpirationService(
	config ExpirationConfig,
	repository ExpirationSwapRepository,
	dvpClient ExpirationDvpClient,
) *ExpirationService {
	return &ExpirationService{
		ticker:     time.NewTicker(config.Interval),
		repository: repository,
		dvpClient:  dvpClient,
	}
}

func (s *ExpirationService) Run(ctx context.Context) error {
	initialRun := make(chan struct{}, 1)
	initialRun <- struct{}{}

	for {
		select {
		case <-ctx.Done():
			s.ticker.Stop()
			return nil
		case <-initialRun:
		case <-s.ticker.C:
		}

		pendingSwaps, err := s.repository.GetPendingSwaps(ctx)
		if err != nil {
			slog.Error("Failed to get pending swaps", slog.Any("error", err))
			continue
		}

		if len(pendingSwaps) == 0 {
			slog.Debug("No pending swaps found")
			continue
		}

		for _, swap := range pendingSwaps {
			expired, err := s.dvpClient.IsSwapExpired(ctx, swap.SharedID)
			if err != nil {
				slog.Error("Failed to check if swap is expired in contract", slog.String("sharedId", swap.SharedID), slog.Any("error", err))
				continue
			}

			if !expired {
				continue
			}

			slog.Info("Expiring swap", slog.String("sharedId", swap.SharedID))

			err = s.dvpClient.ExpireSwap(ctx, swap.SharedID)
			if err != nil {
				slog.Error("Failed to expire swap", slog.String("sharedId", swap.SharedID), slog.Any("error", err))
			}
		}
	}
}
