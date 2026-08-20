package txbatchclient

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/raylsnetwork/rayls-sovereign-relayer/withstack"
)

type MessageWithID interface {
	GetID() string
}

const nonceTimeout = 30 * time.Second

type ethereumClient interface {
	PendingNonceAt(context.Context, common.Address) (uint64, error)
}

type keyQueue interface {
	Enqueue(key *ecdsa.PrivateKey)
	Dequeue(context.Context) (*ecdsa.PrivateKey, error)
}

type transactionSender interface {
	Send(ctx context.Context, txs []*types.Transaction) ([]SendResult, error)
	SendAsync(ctx context.Context, txs []*types.Transaction) ([]SendAsyncResult, error)
}

type transactionGenerator[T MessageWithID] interface {
	Generate(data T, privateKey *ecdsa.PrivateKey, nonce uint64) (*types.Transaction, error)
}

type TypedTxSender[T MessageWithID] struct {
	client ethereumClient

	keyQueue keyQueue

	generator transactionGenerator[T]
	sender    transactionSender
}

func NewTypedTxSender[T MessageWithID](
	generator transactionGenerator[T],
	client ethereumClient,
	sender transactionSender,
	keyQueue keyQueue,
) *TypedTxSender[T] {
	return &TypedTxSender[T]{
		client: client,

		keyQueue: keyQueue,

		generator: generator,
		sender:    sender,
	}
}

func (t *TypedTxSender[T]) SendAsync(ctx context.Context, data []T) (map[string]SendAsyncResult, error) {
	key, err := t.keyQueue.Dequeue(ctx)
	if err != nil {
		return nil, fmt.Errorf("dequeuing signing key: %w", err)
	}
	defer t.keyQueue.Enqueue(key)

	nonce, err := getNonceForPrivateKey(ctx, t.client, key)
	if err != nil {
		return nil, fmt.Errorf("getting nonce for send async: %w", err)
	}

	var txs []*types.Transaction
	txHashToSharedIDMap := map[string]string{}
	for i, d := range data {
		var tx *types.Transaction
		tx, err = t.generator.Generate(
			d,
			key,
			nonce+uint64(i), //nolint:gosec // i is a slice index, always non-negative
		)
		if err != nil {
			return nil, fmt.Errorf("generating transaction: %w", err)
		}

		txs = append(txs, tx)
		txHashToSharedIDMap[tx.Hash().String()] = d.GetID()
	}

	sendResults, err := t.sender.SendAsync(ctx, txs)
	if err != nil {
		return nil, fmt.Errorf("sending async transactions: %w", err)
	}

	sharedIDToSendResultMap := map[string]SendAsyncResult{}
	for _, r := range sendResults {
		sharedID := txHashToSharedIDMap[r.Hash]
		sharedIDToSendResultMap[sharedID] = r
	}
	return sharedIDToSendResultMap, nil
}

// GenerateSigned signs all transactions for the given data using a key from the queue,
// returning them indexed by message ID. The signed transactions can be RLP-encoded and
// persisted for idempotent crash-recovery re-broadcast.
func (t *TypedTxSender[T]) GenerateSigned(ctx context.Context, data []T) (map[string]*types.Transaction, error) {
	key, err := t.keyQueue.Dequeue(ctx)
	if err != nil {
		return nil, fmt.Errorf("dequeuing signing key: %w", err)
	}
	defer t.keyQueue.Enqueue(key)

	nonce, err := getNonceForPrivateKey(ctx, t.client, key)
	if err != nil {
		return nil, fmt.Errorf("getting nonce for generate signed: %w", err)
	}

	result := make(map[string]*types.Transaction, len(data))
	for i, d := range data {
		tx, err := t.generator.Generate(d, key, nonce+uint64(i))
		if err != nil {
			return nil, fmt.Errorf("generating signed transaction: %w", err)
		}
		result[d.GetID()] = tx
	}
	return result, nil
}

func (t *TypedTxSender[T]) Send(ctx context.Context, data []T) (map[string]SendResult, error) {
	key, err := t.keyQueue.Dequeue(ctx)
	if err != nil {
		return nil, fmt.Errorf("dequeuing signing key: %w", err)
	}
	defer t.keyQueue.Enqueue(key)

	nonce, err := getNonceForPrivateKey(ctx, t.client, key)
	if err != nil {
		return nil, fmt.Errorf("getting nonce for send: %w", err)
	}

	var txs []*types.Transaction
	txHashToSharedIDMap := map[string]string{}
	for i, d := range data {
		var tx *types.Transaction
		tx, err = t.generator.Generate(
			d,
			key,
			nonce+uint64(i), //nolint:gosec // i is a slice index, always non-negative
		)
		if err != nil {
			return nil, fmt.Errorf("generating transaction: %w", err)
		}

		txs = append(txs, tx)
		txHashToSharedIDMap[tx.Hash().String()] = d.GetID()
	}

	ctxSend, cancelSend := context.WithTimeout(ctx, batchCallRetryTimeout)
	sendResults, err := t.sender.Send(ctxSend, txs)
	cancelSend() // Clean up immediately to prevent leak
	if err != nil {
		return nil, fmt.Errorf("sending transaction batch: %w", err)
	}

	sharedIDToSendResultMap := map[string]SendResult{}
	for _, r := range sendResults {
		sharedID := txHashToSharedIDMap[r.Hash]
		sharedIDToSendResultMap[sharedID] = r
	}
	return sharedIDToSendResultMap, nil
}

func getNonceForPrivateKey(ctx context.Context, client ethereumClient, privateKey *ecdsa.PrivateKey) (uint64, error) {
	fromAddress := crypto.PubkeyToAddress(privateKey.PublicKey)

	ctxNonce, cancelNonce := context.WithTimeout(ctx, nonceTimeout)
	nonce, err := client.PendingNonceAt(ctxNonce, fromAddress)
	cancelNonce() // Clean up immediately to prevent leak
	if err != nil {
		return 0, withstack.Wrap(fmt.Errorf("getting pending nonce for address [%x]: %w", fromAddress, err))
	}

	return nonce, nil
}
