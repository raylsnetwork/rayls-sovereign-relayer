package service

import (
	"context"
	"fmt"
	"log/slog"
	"math/big"
	"time"

	"github.com/raylsnetwork/rayls-privacy-relayer-api/types"
)

type WaitConfig struct {
	MaxRetries    int
	RetryInterval time.Duration
}

type waiterDepositRepository interface {
	GetDepositByCommitment(ctx context.Context, commitment *big.Int) (*types.DvpDeposit, error)
}

type DepositWaiter struct {
	config     WaitConfig
	repository waiterDepositRepository
}

func NewDepositWaiter(config WaitConfig, repository waiterDepositRepository) *DepositWaiter {
	return &DepositWaiter{
		config:     config,
		repository: repository,
	}
}

// WaitUntilDepositIsConfirmed polls the deposit repository until the deposit reaches Unspent status
// or max retries are exceeded. It returns an error if the deposit never becomes Unspent.
func (w *DepositWaiter) WaitUntilDepositIsConfirmed(ctx context.Context, deposit *types.DvpDeposit) error {
	retryCount := 0

	var err error
	for deposit.Status != types.DvpDepositUnspent {
		retryCount++

		slog.Debug(
			"Deposit not ready",
			slog.String("commitment", deposit.Commitment.String()),
			slog.Int("status", int(deposit.Status)),
		)
		if retryCount >= w.config.MaxRetries {
			return fmt.Errorf("deposit not ready after %d retries: status=%v", w.config.MaxRetries, deposit.Status)
		}

		time.Sleep(w.config.RetryInterval)

		deposit, err = w.repository.GetDepositByCommitment(ctx, deposit.Commitment)
		if err != nil {
			return fmt.Errorf("failed to fetch updated deposit status: %w", err)
		}
	}

	return nil
}
