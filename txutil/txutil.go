package txutil

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	ethTypes "github.com/ethereum/go-ethereum/core/types"

	"github.com/raylsnetwork/rayls-privacy-relayer-api/txbatchclient"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/txsim"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/withstack"
)

// MillisPerSecond is the conversion factor from milliseconds to seconds.
const MillisPerSecond = 1000.0

// Broadcaster sends a signed transaction to the network.
type Broadcaster interface {
	SendTransaction(ctx context.Context, tx *ethTypes.Transaction) error
}

// Receipter polls for transaction receipts.
type Receipter interface {
	GetSingle(ctx context.Context, txHash string) (txbatchclient.ReceiptResult, error)
}

// TransactionSimulator retrieves revert reasons for failed transactions.
type TransactionSimulator interface {
	GetRevertReason(ctx context.Context, txHash common.Hash) (txsim.ContractError, error)
}

// BroadcastSignedTx broadcasts a signed transaction, treating "already known" / "nonce too low"
// as non-errors for idempotent re-broadcast on crash recovery.
func BroadcastSignedTx(ctx context.Context, broadcaster Broadcaster, tx *ethTypes.Transaction) error {
	err := broadcaster.SendTransaction(ctx, tx)
	if err == nil {
		return nil
	}
	errMsg := strings.ToLower(err.Error())
	if strings.Contains(errMsg, "already known") ||
		strings.Contains(errMsg, "nonce too low") ||
		strings.Contains(errMsg, "known transaction") {
		slog.Debug("Transaction already broadcast (idempotent re-send)", slog.String("txHash", tx.Hash().Hex()))
		return nil
	}
	return withstack.Wrap(fmt.Errorf("failed to broadcast transaction: %w", err))
}

// WaitForReceipt polls for the receipt and validates the result, extracting the revert reason
// if the transaction was reverted.
func WaitForReceipt(ctx context.Context, receipter Receipter, txSimulator TransactionSimulator, txHash string) error {
	result, err := receipter.GetSingle(ctx, txHash)
	if err != nil {
		return withstack.Wrap(fmt.Errorf("failed to get receipt for tx %s: %w", txHash, err))
	}
	if result.Error != nil {
		return withstack.Wrap(fmt.Errorf("receipt error for tx %s: %w", txHash, result.Error))
	}
	if result.Receipt.Status == 0 {
		reason, _ := txSimulator.GetRevertReason(ctx, result.Receipt.TxHash)
		return withstack.Wrap(fmt.Errorf("transaction %s was reverted: %s", txHash, reason.String()))
	}
	return nil
}
