package txbatchclient

import (
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/withstack"
)

type TransactionGenerator[T any] interface {
	Generate(data T, privateKey *ecdsa.PrivateKey, nonce uint)
}

type senderEthereumClient interface {
	BatchCallContext(context.Context, []rpc.BatchElem) error
}

type txReceipter interface {
	Get(ctx context.Context, txHashes []string) ([]ReceiptResult, error)
}

type TxSender struct {
	client    senderEthereumClient
	receipter txReceipter
}

func NewTxSender(client senderEthereumClient, receipter txReceipter) *TxSender {
	return &TxSender{
		client:    client,
		receipter: receipter,
	}
}

func (s *TxSender) SendAsync(ctx context.Context, txs []*types.Transaction) ([]SendAsyncResult, error) {
	sendAsyncResults := make([]SendAsyncResult, 0, len(txs))

	requests := make([]rpc.BatchElem, 0, len(txs))
	for _, tx := range txs {
		txData, err := tx.MarshalBinary()
		if err != nil {
			sendAsyncResults = append(sendAsyncResults, SendAsyncResult{
				Error: fmt.Errorf("failed to marshal transaction: %w", err),
			})
			continue
		}

		txHex := hex.EncodeToString(txData)
		requests = append(requests, rpc.BatchElem{
			Method: "eth_sendRawTransaction",
			Args:   []interface{}{"0x" + txHex},
			Result: new(string),
		})
	}

	retryCtx, cancelRetry := context.WithTimeout(ctx, batchCallRetryTimeout)
	defer cancelRetry()

	err := withRetry(retryCtx, func() error {
		opCtx, cancelOp := context.WithTimeout(retryCtx, time.Minute)
		defer cancelOp()
		return s.client.BatchCallContext(opCtx, requests)
	})
	if err != nil {
		return nil, withstack.Wrap(fmt.Errorf("sending raw transactions batch: %w", err))
	}

	for _, req := range requests {
		if req.Result != nil {
			sendAsyncResults = append(sendAsyncResults, SendAsyncResult{
				Hash:  *req.Result.(*string),
				Error: req.Error,
			})
		} else {
			sendAsyncResults = append(sendAsyncResults, SendAsyncResult{
				Error: req.Error,
			})
		}
	}
	return sendAsyncResults, nil
}

func (s *TxSender) Send(ctx context.Context, txs []*types.Transaction) ([]SendResult, error) {
	sendResults := make([]SendResult, 0, len(txs))

	sendAsyncResults, err := s.SendAsync(ctx, txs)
	if err != nil {
		return nil, fmt.Errorf("sending transactions: %w", err)
	}

	var successfullHashes []string
	for _, asyncRes := range sendAsyncResults {
		if asyncRes.Error == nil {
			successfullHashes = append(successfullHashes, asyncRes.Hash)
		} else {
			sendResults = append(sendResults, SendResult{
				Hash:  asyncRes.Hash,
				Error: asyncRes.Error,
			})
		}
	}

	receiptResults, err := s.receipter.Get(ctx, successfullHashes)
	if err != nil {
		return sendResults, fmt.Errorf("getting transaction receipts: %w", err)
	}

	for _, receiptRes := range receiptResults {
		sendResults = append(sendResults, SendResult{
			Hash:    receiptRes.Receipt.TxHash.String(),
			Receipt: receiptRes.Receipt,
			Error:   receiptRes.Error,
		})
	}
	return sendResults, nil
}
