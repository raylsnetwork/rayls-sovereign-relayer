package etherror

import (
	"context"
	"crypto/ecdsa"
	"crypto/rand"
	"fmt"
	"log/slog"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/v2"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"

	"github.com/raylsnetwork/rayls-sovereign-relayer/withstack"
)

const probeGasLimit = 30000

type ledgerProbeBackend interface {
	bind.ContractBackend
	bind.DeployBackend
	ChainID(context.Context) (*big.Int, error)
}

// VerifyLedgerErrorStrings asserts at startup that the live node's error
// messages match the constants in txops_errors.go.
func VerifyLedgerErrorStrings(ctx context.Context, client ledgerProbeBackend, key *ecdsa.PrivateKey) error {
	from := crypto.PubkeyToAddress(key.PublicKey)
	logger := slog.With(slog.String("probe", "ledger-error-strings"), slog.String("from", from.Hex()))

	chainID, err := client.ChainID(ctx)
	if err != nil {
		return withstack.Wrap(fmt.Errorf("probe: chain id: %w", err))
	}
	gasPrice, err := client.SuggestGasPrice(ctx)
	if err != nil {
		return withstack.Wrap(fmt.Errorf("probe: suggest gas price: %w", err))
	}
	nonce, err := client.PendingNonceAt(ctx, from)
	if err != nil {
		return withstack.Wrap(fmt.Errorf("probe: pending nonce for %s: %w", from.Hex(), err))
	}

	// Unique calldata so each boot's probe tx has distinct bytes.
	marker := make([]byte, 8)
	if _, err := rand.Read(marker); err != nil {
		return withstack.Wrap(fmt.Errorf("probe: random marker: %w", err))
	}

	signer := types.LatestSignerForChainID(chainID)
	auth := &bind.TransactOpts{
		From:     from,
		Nonce:    new(big.Int).SetUint64(nonce),
		GasPrice: gasPrice,
		GasLimit: probeGasLimit, // see probeGasLimit: MUST be set or gas estimation fails with ErrNoCode
		Value:    big.NewInt(0),
		NoSend:   true, // build + sign only; this probe controls every broadcast
		Context:  ctx,
		Signer: func(_ common.Address, tx *types.Transaction) (*types.Transaction, error) {
			return types.SignTx(tx, signer, key)
		},
	}

	// Self-transfer (to == from); marker rides as calldata and is ignored by the
	// EOA recipient. With Nonce/GasPrice/GasLimit all set, RawTransact issues no
	// network calls — it just builds and signs.
	bound := bind.NewBoundContract(from, abi.ABI{}, client, client, client)
	tx, err := bound.RawTransact(auth, marker)
	if err != nil {
		return withstack.Wrap(fmt.Errorf("probe: build probe tx: %w", err))
	}

	logger = logger.With(slog.String("txHash", tx.Hash().Hex()), slog.Uint64("nonce", nonce))
	logger.Info("verifying ledger error strings against live node")

	// 1. Initial broadcast — should be accepted into the mempool.
	bcastErr := client.SendTransaction(ctx, tx)
	if bcastErr != nil && !Is(bcastErr, AlreadyKnownError) {
		return withstack.Wrap(fmt.Errorf("probe: initial broadcast at nonce %d failed: %w", nonce, bcastErr))
	}

	// 2. Identical re-broadcast while still pending — must match a known
	//    "already known" variant.
	rebErr := client.SendTransaction(ctx, tx)
	switch {
	case Is(bcastErr, AlreadyKnownError):
		logger.Info("already-known string verified", slog.String("got", bcastErr.Error()))
	case Is(rebErr, AlreadyKnownError):
		logger.Info("already-known string verified", slog.String("got", rebErr.Error()))
	case rebErr == nil:
		logger.Info("re-broadcast accepted idempotently (nil); node emits no already-known string")
	default:
		return ledgerStringDivergence(AlreadyKnownError, rebErr)
	}

	// 3. Wait for the probe tx to mine so the account nonce advances past it.
	if _, err := bind.WaitMined(ctx, client, tx.Hash()); err != nil {
		return withstack.Wrap(fmt.Errorf("probe: waiting for probe tx %s to mine: %w", tx.Hash().Hex(), err))
	}

	// 4. Re-broadcast the now-stale tx — must match a known "nonce too low" variant.
	if staleErr := client.SendTransaction(ctx, tx); !Is(staleErr, NonceTooLowError) {
		return ledgerStringDivergence(NonceTooLowError, staleErr)
	} else {
		logger.Info("nonce-too-low string verified", slog.String("got", staleErr.Error()))
	}

	logger.Info("ledger error strings verified")
	return nil
}

// ledgerStringDivergence reports that the node returned an error matching none of
// target's known variants — the probe's whole reason to exist.
func ledgerStringDivergence(target EthError, got error) error {
	gotMsg := "<nil> (re-broadcast unexpectedly succeeded)"
	if got != nil {
		gotMsg = got.Error()
	}
	return fmt.Errorf(
		"ledger error string divergence for %q: node returned %q, which matches none of the known "+
			"variants %q — add this node's wire string to the EthError variants in cts/etherror/error.go",
		target.Name(), gotMsg, target.Variants(),
	)
}
