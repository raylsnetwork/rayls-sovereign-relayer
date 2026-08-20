package service

import (
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"errors"
	"fmt"
	"log/slog"
	"math/big"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/v2"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/txpool"
	"github.com/ethereum/go-ethereum/core/types"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/raylsnetwork/rayls-sovereign-relayer/cts/etherror"
	"github.com/raylsnetwork/rayls-sovereign-relayer/cts/ethrpc"
	"github.com/raylsnetwork/rayls-sovereign-relayer/ethretry"
	"github.com/raylsnetwork/rayls-sovereign-relayer/faultinjector"
	"github.com/raylsnetwork/rayls-sovereign-relayer/withstack"
)

type txOpsTxSyncRepository interface {
	Create(ctx context.Context, t *SyncTransaction) error
	Get(ctx context.Context, id string) (*SyncTransaction, error)
	Save(ctx context.Context, t *SyncTransaction) error
}

type txOpsKeyQueue interface {
	Enqueue(key *ecdsa.PrivateKey)
	Dequeue(ctx context.Context) (*ecdsa.PrivateKey, error)
}

type txOpsAuthGen interface {
	Get(context.Context, *ecdsa.PrivateKey) (*bind.TransactOpts, error)
	// InvalidateNonceCache resets the per-key nonce cache so the next
	// `Get` re-queries the chain. Required after any error path that
	// aborts a tx between Get (which reserved the nonce) and the
	// broadcast actually being accepted by the chain — otherwise the
	// cache strides past nonces that were never consumed.
	InvalidateNonceCache(common.Address)
}

type txOpsEthClient interface {
	bind.ContractBackend
	bind.DeployBackend
	TransactionReceipt(ctx context.Context, txHash common.Hash) (*types.Receipt, error)
	TransactionByHash(ctx context.Context, hash common.Hash) (tx *types.Transaction, isPending bool, err error)
	CallContract(ctx context.Context, call ethereum.CallMsg, blockNumber *big.Int) ([]byte, error)
	ChainID(ctx context.Context) (*big.Int, error)
}

type txOpsBatcher interface {
	Send(ctx context.Context, inputs []ethrpc.TransactionInput) ([]ethrpc.SendResult, error)
	GetReceipts(ctx context.Context, hashes []common.Hash) ([]ethrpc.ReceiptResult, error)
}

// defaultWaitMinedTimeout caps a single bind.WaitMined poll loop. WaitMined
// blocks until the tx mines or ctx is done; if the tx never mines it would
// otherwise poll forever on whatever ctx the caller passed. This is a
// last-resort ceiling for callers whose ctx has no (or a very loose) deadline —
// the relayer's per-message handler already imposes a tighter bound and, since
// context.WithTimeout adopts the earliest deadline, that tighter caller bound
// still wins. Kept generous so normal mining latency never trips it.
const defaultWaitMinedTimeout = 5 * time.Minute

// TxOpsService signs calldata with a key from the queue, broadcasts it,
// waits for the receipt, and reports the outcome. It mimics
// contractclient.Executor semantics closely but exposes a revert path as data
// instead of a typed error so the cts gRPC surface can carry it in a oneof.
// Write operations are idempotent on the caller-supplied id, backed by the
// cts_sync_tx ledger (txSyncRepo).
type TxOpsService struct {
	gen        txOpsAuthGen
	queue      txOpsKeyQueue
	client     txOpsEthClient
	batcher    txOpsBatcher
	txSyncRepo txOpsTxSyncRepository

	// waitMinedTimeout bounds each WaitMined step so a tx that never mines
	// can't wedge the service indefinitely. See defaultWaitMinedTimeout.
	waitMinedTimeout time.Duration
}

func NewTxOpsService(gen txOpsAuthGen, queue txOpsKeyQueue, client txOpsEthClient, batcher txOpsBatcher, txSyncRepo txOpsTxSyncRepository) *TxOpsService {
	return NewTxOpsServiceWithWaitMinedTimeout(gen, queue, client, batcher, txSyncRepo, defaultWaitMinedTimeout)
}

// NewTxOpsServiceWithWaitMinedTimeout is like NewTxOpsService but lets callers
// (and tests) set the per-tx WaitMined ceiling. A non-positive timeout falls
// back to the default.
func NewTxOpsServiceWithWaitMinedTimeout(
	gen txOpsAuthGen,
	queue txOpsKeyQueue,
	client txOpsEthClient,
	batcher txOpsBatcher,
	txSyncRepo txOpsTxSyncRepository,
	waitMinedTimeout time.Duration,
) *TxOpsService {
	if waitMinedTimeout <= 0 {
		waitMinedTimeout = defaultWaitMinedTimeout
	}
	return &TxOpsService{
		gen:              gen,
		queue:            queue,
		client:           client,
		batcher:          batcher,
		txSyncRepo:       txSyncRepo,
		waitMinedTimeout: waitMinedTimeout,
	}
}

// waitMined runs bind.WaitMined under a ctx bounded by waitMinedTimeout. Because
// context.WithTimeout adopts the earliest of the parent and new deadlines, a
// tighter caller deadline (e.g. the relayer's per-message handler timeout) still
// takes precedence; this only bites when the caller's ctx would otherwise block
// far longer — or forever.
func (s *TxOpsService) waitMined(ctx context.Context, txHash common.Hash) (*ethtypes.Receipt, error) {
	ctx, cancel := context.WithTimeout(ctx, s.waitMinedTimeout)
	defer cancel()
	return bind.WaitMined(ctx, s.client, txHash)
}

// SignAndSend is idempotent on `id`: a fresh id signs, persists, broadcasts,
// and waits for the receipt; a seen id enters the recovery path against the
// persisted (hash, RLP). The operation is not bounded here — the caller's ctx
// (e.g. the relayer's gRPC deadline) is the outer budget, and each WaitMined
// step is individually capped by waitMinedTimeout (CTS_WAIT_MINED_TIMEOUT_SECS).
func (s *TxOpsService) SignAndSend(ctx context.Context, id string, data []byte, address common.Address) (SignAndSendResult, error) {
	logger := slog.With(slog.String("id", id))
	logger.Info("signing and sending transaction")

	syncTx, err := s.txSyncRepo.Get(ctx, id)
	if err != nil {
		if !errors.Is(err, ErrTxNotFound) {
			// unsupported error from the repo - propagate
			logger.Error("failed to get transaction from repo", slog.Any("err", err))
			return SignAndSendResult{}, fmt.Errorf("get sync transaction: %w", err)
		}

		logger.Debug("new transaction - begin fresh path")
		// transaction doesn't exist - begin clean generation
		key, err := s.queue.Dequeue(ctx)
		if err != nil {
			logger.Error("failed to dequeue key", slog.Any("err", err))
			return SignAndSendResult{}, fmt.Errorf("dequeue key: %w", err)
		}
		defer s.queue.Enqueue(key)

		return s.executeTransaction(ctx, logger, key, id, data, address)
	} else {
		logger.Debug("transaction already seen - begin recovery path")

		// Hard-terminal rows (mined/reverted) are read-only: a retry reads the
		// stored verdict back without touching the chain.
		if syncTx.State.IsTerminal() {
			logger.Info("returning stored verdict, skipping chain", slog.String("state", syncTx.State.String()))
			return resultFromSyncTx(syncTx), nil
		}

		// transaction already exists but no result in the repo
		// try resending it nonce deduplication protects us from double mint
		logger.Debug("resending transaction to node")
		err = s.client.SendTransaction(ctx, syncTx.Tx)

		if err == nil || etherror.Is(err, etherror.AlreadyKnownError) {
			// the tx was just created or it was already in the mempool
			// either way the node knows for it, so just wait for it to finalize
			logger.Warn("transaction confirmed in node - finalizing")

			// wait receipt and update tx in repo. finalizeTransaction logs its
			// own errors, so we only wrap here.
			result, err := s.finalizeTransaction(ctx, logger, syncTx)
			if err != nil {
				return SignAndSendResult{}, fmt.Errorf("finalize transaction: %w", err)
			}
			return result, nil
		} else if isNodeErr(err, core.ErrNonceTooLow) {
			// either the transaction was mined, or someone else used our nonce
			// again, try finalizing the transaction and see if we get an error

			// check if the node has received our transaction
			_, _, err := s.client.TransactionByHash(ctx, syncTx.Tx.Hash())
			if err != nil {
				if ctx.Err() != nil {
					// context got canceled/expired could be due to
					// a network error or just because - anyway, return
					logger.Error("context done while looking up transaction by hash", slog.Any("err", ctx.Err()))
					return SignAndSendResult{}, ctx.Err()
				}
				if errors.Is(err, ethereum.NotFound) {
					// someone else used our nonce - resign and resend
					// transaction doesn't exist - begin clean generation
					logger.Debug("original tx never reached node and old nonce taken - reexecuting")
					key, err := s.queue.Dequeue(ctx)
					if err != nil {
						logger.Error("failed to dequeue key", slog.Any("err", err))
						return SignAndSendResult{}, fmt.Errorf("dequeue key: %w", err)
					}
					defer s.queue.Enqueue(key)

					return s.reExecuteTransaction(ctx, logger, key, syncTx, data, address)
				}
				logger.Error("failed to look up transaction by hash", slog.Any("err", err))
				return SignAndSendResult{}, fmt.Errorf("transaction by hash: %w", err)
			}
			// transaction was found - finalize it
			logger.Debug("tx already mined in node - finalizing")
			return s.finalizeTransaction(ctx, logger, syncTx)
		} else {
			// unsupported error code - propagate
			logger.Error("failed to send transaction to node", slog.Any("err", err))
			return SignAndSendResult{}, fmt.Errorf("send transaction: %w", err)
		}
	}

}

// isNodeErr reports whether err is sentinel. It matches by identity first (for
// in-process backends and mocks that return the sentinel directly) and falls
// back to message text, because an error reconstructed from a JSON-RPC response
// carries the node's message string but not the typed sentinel value. Matching
// against sentinel.Error() keeps the comparison in sync with the sentinel text.
// Note: text matching is go-ethereum-family specific — re-validate for other
// EVM clients.
func isNodeErr(err, sentinel error) bool {
	if errors.Is(err, sentinel) {
		return true
	}
	// case-insensitive: nodes vary the casing (e.g. "Nonce too low" vs the
	// sentinel's "nonce too low").
	return err != nil && strings.Contains(strings.ToLower(err.Error()), strings.ToLower(sentinel.Error()))
}

// resultFromSyncTx builds the caller-facing result from a hard-terminal row's
// stored verdict, without touching the chain. Reverted rows surface their
// stored revert data; mined rows surface the stored receipt.
func resultFromSyncTx(syncTx *SyncTransaction) SignAndSendResult {
	if syncTx.State == StateReverted {
		return SignAndSendResult{Revert: &Revert{RevertData: syncTx.RevertData}}
	}
	return SignAndSendResult{Success: &SignAndSendSuccess{Receipt: syncTx.Receipt}}
}

func (s *TxOpsService) Deploy(ctx context.Context, bytecode []byte, constructor []byte) (DeployResult, error) {
	key, err := s.queue.Dequeue(ctx)
	if err != nil {
		return DeployResult{}, fmt.Errorf("dequeue key: %w", err)
	}
	defer s.queue.Enqueue(key)

	auth, err := s.gen.Get(ctx, key)
	if err != nil {
		return DeployResult{}, fmt.Errorf("generate transaction auth: %w", err)
	}
	// Invalidate the reserved nonce on any pre-mine error path.
	signerAddr := crypto.PubkeyToAddress(key.PublicKey)

	nonce, err := s.getNonceForAddress(ctx, auth.From)
	if err != nil {
		s.gen.InvalidateNonceCache(signerAddr)
		return DeployResult{}, err
	}
	auth.Nonce = new(big.Int).SetUint64(nonce)

	var (
		address common.Address
		tx      *ethtypes.Transaction
	)
	err = ethretry.WithRetry(ctx, func() error {
		var deployErr error
		address, tx, deployErr = bind.DeployContract(auth, bytecode, s.client, constructor)
		return deployErr
	})
	if err != nil {
		s.gen.InvalidateNonceCache(signerAddr)
		return DeployResult{}, fmt.Errorf("send contract for deployment: %w", err)
	}

	var receipt *ethtypes.Receipt
	err = ethretry.WithRetry(ctx, func() error {
		var wmErr error
		receipt, wmErr = s.waitMined(ctx, tx.Hash())
		return wmErr
	})
	if err != nil {
		s.gen.InvalidateNonceCache(signerAddr)
		return DeployResult{}, fmt.Errorf("wait mined: %w", err)
	}

	if receipt.Status == ethtypes.ReceiptStatusFailed {
		revertBytes, reasonErr := s.tryGetRevertReason(ctx, auth.From, tx, receipt.BlockNumber)
		if reasonErr != nil {
			return DeployResult{}, fmt.Errorf("tx reverted, failed to get revert reason: %w", reasonErr)
		}
		return DeployResult{Revert: &Revert{RevertData: revertBytes}}, nil
	}

	return DeployResult{
		Success: &DeploySuccess{
			Address: address,
			Receipt: receipt,
		},
	}, nil

}

// Call executes an eth_call against `address` with `data` as calldata. A key
// is dequeued purely to derive the `from` address — matching
// contractclient.Executor.Call at contractclient/executor.go:125-147 so
// access-controlled view functions still see a known caller. The key is
// re-enqueued on return. Reverts are surfaced as CallResult.Revert with raw
// bytes so callers can decode against their own binding error unpackers;
// infra failures are returned as errors.
func (s *TxOpsService) Call(ctx context.Context, data []byte, address common.Address) (CallResult, error) {
	key, err := s.queue.Dequeue(ctx)
	if err != nil {
		return CallResult{}, fmt.Errorf("dequeue key: %w", err)
	}
	defer s.queue.Enqueue(key)

	from := crypto.PubkeyToAddress(key.PublicKey)

	msg := ethereum.CallMsg{
		From: from,
		To:   &address,
		Data: data,
	}

	var out []byte
	err = ethretry.WithRetry(ctx, func() error {
		var callErr error
		out, callErr = s.client.CallContract(ctx, msg, nil)
		return callErr
	})
	if err != nil {
		// A view-function revert arrives as an RPC DataError carrying the
		// ABI-encoded revert bytes. Surface it as data, not as an infra error.
		// WithRetry only re-runs on transient transport / decode failures, so a
		// deterministic revert returns immediately without consuming retries —
		// matching the read-only CallContract in tryGetRevertReason.
		if revertBytes, ok := extractRevertData(err); ok {
			return CallResult{Revert: &Revert{RevertData: revertBytes}}, nil
		}
		return CallResult{}, fmt.Errorf("call contract: %w", err)
	}

	return CallResult{Value: out}, nil
}

// getNonceForAddress returns the pending nonce for the given address from the chain.
func (s *TxOpsService) getNonceForAddress(ctx context.Context, from common.Address) (uint64, error) {
	nonce, err := s.client.PendingNonceAt(ctx, from)
	if err != nil {
		return 0, withstack.Wrap(fmt.Errorf("pending nonce for %x: %w", from, err))
	}
	return nonce, nil
}

// tryGetRevertReason re-runs the transaction as a read-only eth_call at the
// mined block to extract revert data. Mirrors contractclient.Executor.tryGetRevertReason.
func (s *TxOpsService) tryGetRevertReason(
	ctx context.Context,
	from common.Address,
	tx *ethtypes.Transaction,
	blockNumber *big.Int,
) ([]byte, error) {
	msg := ethereum.CallMsg{
		From:     from,
		To:       tx.To(),
		Gas:      tx.Gas(),
		GasPrice: tx.GasPrice(),
		Value:    tx.Value(),
		Data:     tx.Data(),
	}

	err := ethretry.WithRetry(ctx, func() error {
		_, callErr := s.client.CallContract(ctx, msg, blockNumber)
		return callErr
	})
	if err == nil {
		// The tx reverted on-chain but replaying it as a call succeeded —
		// state drift between mine-time and simulate-time. Return no data;
		// the caller will still see the revert classification.
		return nil, nil
	}

	if data, ok := extractRevertData(err); ok {
		return data, nil
	}

	return nil, err
}

const receiptPollInterval = 500 * time.Millisecond

// BatchSignAndSend signs all items with a single dequeued key (incrementing
// nonces), broadcasts them as a batch, waits for every receipt, and returns
// per-msgID results. Function-level errors are whole-batch infra failures;
// per-tx outcomes (success, revert, error) are in the result map.
func (s *TxOpsService) BatchSignAndSend(ctx context.Context, items []BatchItem) (map[string]BatchItemResult, error) {
	inputs := make([]ethrpc.TransactionInput, len(items))
	idByIndex := make([]string, len(items))
	for i, item := range items {
		inputs[i] = ethrpc.TransactionInput{
			ID:      item.MsgID,
			Data:    item.Data,
			Address: item.Address,
		}
		idByIndex[i] = item.MsgID
	}

	// s.batcher.Send handles nonce-cache invalidation internally for
	// whole-batch failures, per-tx broadcast errors, and partial
	// sign-step drops — see ethrpc.RPCBatcher.Send. No explicit
	// invalidation needed here (unlike SignAndSend / Deploy below,
	// which go through AuthGen directly and must invalidate
	// explicitly on any pre-mine error path).
	sendResults, err := s.batcher.Send(ctx, inputs)
	if err != nil {
		return nil, fmt.Errorf("batch send: %w", err)
	}

	results := make(map[string]BatchItemResult, len(items))

	// hashToMsgID maps tx hash to msgID for receipt correlation.
	hashToMsgID := make(map[common.Hash]string)
	for i, sr := range sendResults {
		msgID := idByIndex[i]
		switch {
		case sr.Revert != nil:
			results[msgID] = BatchItemResult{Revert: &Revert{RevertData: sr.Revert.Data}}
		case sr.Err != nil:
			results[msgID] = BatchItemResult{Err: sr.Err}
		default:
			hashToMsgID[sr.Hash] = msgID
		}
	}

	// Poll receipts for all successfully broadcast transactions.
	for len(hashToMsgID) > 0 {
		hashes := make([]common.Hash, 0, len(hashToMsgID))
		for h := range hashToMsgID {
			hashes = append(hashes, h)
		}

		receiptResults, err := s.batcher.GetReceipts(ctx, hashes)
		if err != nil {
			return nil, fmt.Errorf("batch get receipts: %w", err)
		}

		for i, rr := range receiptResults {
			hash := hashes[i]
			msgID := hashToMsgID[hash]
			switch {
			case rr.Pending:
				continue
			case rr.Receipt != nil:
				results[msgID] = BatchItemResult{
					Success: &SignAndSendSuccess{Receipt: rr.Receipt},
				}
				delete(hashToMsgID, hash)
			case rr.Revert != nil:
				results[msgID] = BatchItemResult{
					Revert: &Revert{RevertData: rr.Revert.Data},
				}
				delete(hashToMsgID, hash)
			case rr.Err != nil:
				results[msgID] = BatchItemResult{Err: rr.Err}
				delete(hashToMsgID, hash)
			}
		}

		if len(hashToMsgID) > 0 {
			select {
			case <-ctx.Done():
				return nil, fmt.Errorf("waiting for receipts: %w", ctx.Err())
			case <-time.After(receiptPollInterval):
			}
		}
	}

	return results, nil
}

// extractRevertData pulls revert bytes from a go-ethereum RPC error. Returns
// (nil, false) if the error is not a DataError or the blob cannot be decoded.
func extractRevertData(err error) ([]byte, bool) {
	var de rpc.DataError
	if !errors.As(err, &de) {
		return nil, false
	}
	raw, ok := de.ErrorData().(string)
	if !ok {
		return nil, false
	}
	data, decodeErr := hex.DecodeString(strings.TrimPrefix(raw, "0x"))
	if decodeErr != nil {
		return nil, false
	}
	return data, true
}

func (s *TxOpsService) generateTransaction(ctx context.Context, key *ecdsa.PrivateKey, data []byte, address common.Address) (*types.Transaction, error) {
	// auth.Nonce is already set by gen.Get from the shared nonce cache
	// (AuthGen reserves the next nonce there). Do NOT override it with a raw
	// PendingNonceAt read: that bypasses the cache, lets this path drift out of
	// sync with the other components signing for the same key (RPCBatcher), and
	// reintroduces the colliding/gapped-nonce race the cache exists to prevent.
	auth, err := s.gen.Get(ctx, key)
	if err != nil {
		return nil, fmt.Errorf("generate transaction auth: %w", err)
	}
	// ONLY GENERATE TRANSACTION. DONT SEND IT
	auth.NoSend = true

	bound := bind.NewBoundContract(
		address,
		abi.ABI{},
		s.client,
		s.client,
		s.client,
	)
	tx, err := bound.RawTransact(auth, data)
	if err != nil {
		// gen.Get reserved a nonce but no tx reached the wire (sign failure or a
		// revert during gas estimation). Invalidate so the next reservation
		// realigns with the chain rather than striding past a nonce that will
		// never be consumed — an unfilled nonce stalls everything behind it.
		s.gen.InvalidateNonceCache(crypto.PubkeyToAddress(key.PublicKey))
		return nil, fmt.Errorf("send tx: %w", err)
	}

	return tx, nil
}

func (s *TxOpsService) executeTransaction(ctx context.Context, logger *slog.Logger, key *ecdsa.PrivateKey, id string, data []byte, address common.Address) (SignAndSendResult, error) {
	tx, err := s.generateTransaction(ctx, key, data, address)
	if err != nil {
		// try extracting revert data from gas estimation
		if errors.As(err, new(rpc.DataError)) {
			if revertBytes, ok := extractRevertData(err); ok {
				logger.Debug("transaction reverted during gas estimation")
				return SignAndSendResult{Revert: &Revert{RevertData: revertBytes}}, nil
			}
		}
		logger.Error("failed to generate transaction", slog.Any("err", err))
		return SignAndSendResult{}, fmt.Errorf("send tx: %w", err)
	}

	from := crypto.PubkeyToAddress(key.PublicKey)
	// correlate every downstream line to the actual signed tx on-chain.
	logger = logger.With(slog.String("txHash", tx.Hash().Hex()), slog.String("from", from.Hex()))

	syncTx := NewSyncTransaction(id, from, tx)
	err = s.txSyncRepo.Create(ctx, &syncTx)
	if err != nil {
		// TODO: check if someone already created the tx before us
		// either wait for result, or return error and make relayer resend
		// Signed at the reserved nonce but never broadcast — release it so the
		// next reservation realigns with the chain instead of leaving a hole.
		s.gen.InvalidateNonceCache(from)
		logger.Error("failed to persist transaction", slog.Any("err", err))
		return SignAndSendResult{}, fmt.Errorf("create transaction: %w", err)
	}

	if err = faultinjector.Check("cts.service.TxOpsService.SignAndSend.between_persist_and_send"); err != nil {
		logger.Warn("FI triggered between persist and send")

		s.gen.InvalidateNonceCache(from)
		return SignAndSendResult{}, fmt.Errorf("fault injector at between_persist_and_send: %w", ErrRetriable)
	}

	// persisted before broadcast — only now put the signed tx on the wire.
	if err := s.client.SendTransaction(ctx, tx); err != nil && !errors.Is(err, txpool.ErrAlreadyKnown) {
		// The node rejected the tx, so the reserved nonce never reached the
		// wire. Invalidate so we don't strand every later nonce behind it.
		s.gen.InvalidateNonceCache(from)
		logger.Error("failed to broadcast transaction", slog.Any("err", err))
		return SignAndSendResult{}, fmt.Errorf("send transaction: %w", err)
	}

	if err := faultinjector.Check("cts.service.TxOpsService.SignAndSend.before_wait_mined"); err != nil {
		return SignAndSendResult{}, fmt.Errorf("fault injector at before_wait_mined: %w", err)
	}

	// wait receipt and update tx in repo. finalizeTransaction logs its own
	// errors, so we only wrap here.
	result, err := s.finalizeTransaction(ctx, logger, &syncTx)
	if err != nil {
		return SignAndSendResult{}, fmt.Errorf("finalize transaction: %w", err)
	}
	return result, nil
}

func (s *TxOpsService) finalizeTransaction(ctx context.Context, logger *slog.Logger, syncTx *SyncTransaction) (SignAndSendResult, error) {
	var result SignAndSendResult
	// wait for the transaction to get mined, bounded by waitMinedTimeout
	// (CTS_WAIT_MINED_TIMEOUT_SECS); a tighter caller deadline still wins.
	receipt, err := s.waitMined(ctx, syncTx.Tx.Hash())
	if err != nil {
		// A deadline here is expected backpressure: the tx hasn't mined within
		// the budget, so the caller retries and re-broadcasts the same bytes.
		// Anything else is a genuine failure.
		if errors.Is(err, context.DeadlineExceeded) {
			logger.Warn("timed out waiting for receipt - will retry", slog.Any("err", err))
		} else {
			logger.Error("failed waiting for transaction to be mined", slog.Any("err", err))
		}
		return SignAndSendResult{}, fmt.Errorf("wait mined: %w", err)
	}

	if receipt.Status == types.ReceiptStatusFailed {
		// On a failed revert-reason lookup we still resolve the row as reverted,
		// persisting empty revert data as the "reason unavailable" flag.
		revertData, err := s.tryGetRevertReason(ctx, syncTx.From, syncTx.Tx, receipt.BlockNumber)
		if err != nil {
			logger.Warn("failed to get revert reason - persisting empty revert data", slog.Any("err", err))
			revertData = []byte{}
		}
		result.Revert = &Revert{RevertData: revertData}
		if err := syncTx.ResolveReverted(receipt, revertData); err != nil {
			logger.Error("failed to resolve transaction as reverted", slog.Any("err", err))
			return SignAndSendResult{}, fmt.Errorf("resolve reverted: %w", err)
		}
		logger.Info("transaction reverted")
	} else {
		result.Success = &SignAndSendSuccess{Receipt: receipt}
		if err := syncTx.ResolveMined(receipt); err != nil {
			logger.Error("failed to resolve transaction as mined", slog.Any("err", err))
			return SignAndSendResult{}, fmt.Errorf("resolve mined: %w", err)
		}
		logger.Info("transaction mined")
	}

	err = s.txSyncRepo.Save(ctx, syncTx)
	if err != nil {
		logger.Error("failed to persist transaction result", slog.Any("err", err))
		return SignAndSendResult{}, fmt.Errorf("tx repo save: %w", err)
	}

	return result, nil
}

// reExecuteTransaction re-signs the logical message at a fresh nonce after the
// original nonce slot was consumed by another caller (nonce-too-low +
// TxByHash NotFound). It rebinds the loaded row to the new signed tx — keeping
// the id and the row version so the optimistic-concurrency check still holds —
// overwrites the persisted (hash, RLP), then broadcasts and finalizes.
func (s *TxOpsService) reExecuteTransaction(ctx context.Context, logger *slog.Logger, key *ecdsa.PrivateKey, syncTx *SyncTransaction, data []byte, address common.Address) (SignAndSendResult, error) {
	tx, err := s.generateTransaction(ctx, key, data, address)
	if err != nil {
		// try extracting revert data from gas estimation
		if errors.As(err, new(rpc.DataError)) {
			if revertBytes, ok := extractRevertData(err); ok {
				logger.Debug("transaction reverted during gas estimation")
				return SignAndSendResult{Revert: &Revert{RevertData: revertBytes}}, nil
			}
		}
		logger.Error("failed to generate transaction", slog.Any("err", err))
		return SignAndSendResult{}, fmt.Errorf("send tx: %w", err)
	}

	// re-signing at a fresh nonce is rare and operationally significant: it
	// means the original nonce slot was consumed by another caller (a crash
	// dropped the key lease). Surface old->new hash at warn level.
	logger.Warn("re-signing transaction at a fresh nonce",
		slog.String("oldTxHash", syncTx.Tx.Hash().Hex()),
		slog.String("newTxHash", tx.Hash().Hex()),
	)
	logger = logger.With(slog.String("txHash", tx.Hash().Hex()))

	// The fresh nonce just reserved for the re-signed tx is burned on any path
	// that fails before it reaches the wire.
	from := crypto.PubkeyToAddress(key.PublicKey)

	// rebind the loaded row to the freshly-signed tx and overwrite the
	// persisted (hash, RLP) before broadcasting.
	if err := syncTx.Rebind(tx); err != nil {
		s.gen.InvalidateNonceCache(from)
		logger.Error("failed to rebind transaction", slog.Any("err", err))
		return SignAndSendResult{}, fmt.Errorf("rebind tx: %w", err)
	}
	if err := s.txSyncRepo.Save(ctx, syncTx); err != nil {
		s.gen.InvalidateNonceCache(from)
		logger.Error("failed to overwrite transaction", slog.Any("err", err))
		return SignAndSendResult{}, fmt.Errorf("overwrite transaction: %w", err)
	}

	// persisted before broadcast — only now put the re-signed tx on the wire.
	if err := s.client.SendTransaction(ctx, tx); err != nil && !errors.Is(err, txpool.ErrAlreadyKnown) {
		s.gen.InvalidateNonceCache(from)
		logger.Error("failed to broadcast transaction", slog.Any("err", err))
		return SignAndSendResult{}, fmt.Errorf("send transaction: %w", err)
	}

	// finalizeTransaction logs its own errors, so we only wrap here.
	result, err := s.finalizeTransaction(ctx, logger, syncTx)
	if err != nil {
		// Broadcast may have succeeded, but we couldn't confirm — invalidate
		// so a stale local cache doesn't drift out of step with the chain.
		s.gen.InvalidateNonceCache(from)
		return SignAndSendResult{}, fmt.Errorf("finalize transaction: %w", err)
	}
	return result, nil
}
