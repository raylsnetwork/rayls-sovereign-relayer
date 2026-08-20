package txbatchclient

import (
	"context"
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/withstack"
	"golang.org/x/exp/maps"
)

// batchCallRetryTimeout is the timeout for retrying batch RPC calls.
const batchCallRetryTimeout = 5 * time.Minute

type receipterEthereumClient interface {
	BatchCallContext(context.Context, []rpc.BatchElem) error
}

type TxReceipter struct {
	client receipterEthereumClient
}

func NewTxReceipter(client receipterEthereumClient) *TxReceipter {
	return &TxReceipter{
		client: client,
	}
}

func (s *TxReceipter) Get(ctx context.Context, txHashes []string) ([]ReceiptResult, error) {
	var receiptResults []ReceiptResult

	pendingHashes := map[string]interface{}{}
	for _, txHash := range txHashes {
		pendingHashes[txHash] = nil
	}

	for {
		select {
		case <-ctx.Done():
			return receiptResults, withstack.Wrap(ctx.Err())
		default:
			availableReceiptResults, err := s.getAvailable(ctx, maps.Keys(pendingHashes))
			if err != nil {
				return nil, fmt.Errorf("fetching available receipts: %w", err)
			}

			receiptResults = append(receiptResults, availableReceiptResults...)
			for _, res := range availableReceiptResults {
				delete(pendingHashes, res.Receipt.TxHash.String())
			}

			if len(pendingHashes) == 0 {
				return receiptResults, nil
			}

			select {
			case <-ctx.Done():
				return receiptResults, withstack.Wrap(ctx.Err())
			case <-time.After(500 * time.Millisecond):
			}
		}
	}
}

func (s *TxReceipter) GetSingle(ctx context.Context, txHash string) (ReceiptResult, error) {
	results, err := s.Get(ctx, []string{txHash})
	if err != nil {
		return ReceiptResult{}, fmt.Errorf("getting single receipt: %w", err)
	}

	if len(results) == 0 {
		return ReceiptResult{}, nil
	}

	return results[0], nil
}

func (s *TxReceipter) getAvailable(ctx context.Context, txHashes []string) ([]ReceiptResult, error) {
	var receiptResults []ReceiptResult

	requests := make([]rpc.BatchElem, 0, len(txHashes))
	for _, tx := range txHashes {
		requests = append(requests, rpc.BatchElem{
			Method: "eth_getTransactionReceipt",
			Args:   []interface{}{tx},
			Result: new(types.Receipt),
		})
	}

	retryCtx, cancelRetry := context.WithTimeout(ctx, batchCallRetryTimeout)
	defer cancelRetry()

	callErr := withRetry(retryCtx, func() error {
		opCtx, cancelOp := context.WithTimeout(retryCtx, time.Minute)
		defer cancelOp()
		return s.client.BatchCallContext(opCtx, requests)
	})
	if callErr != nil {
		return nil, withstack.Wrap(callErr)
	}

	for _, request := range requests {
		receipt, _ := request.Result.(*types.Receipt)
		if receipt.TxHash != (common.Hash{}) {
			var reqErr error
			if request.Error != nil {
				reqErr = withstack.Wrap(request.Error)
			}
			res := ReceiptResult{
				Receipt: receipt,
				Error:   reqErr,
			}
			receiptResults = append(receiptResults, res)
		}
	}

	return receiptResults, nil
}
