package contractclient

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	txopspb "github.com/raylsnetwork/rayls-privacy-relayer-api/cts/gen/proto/txops"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// The proto field is named `receipt_rlp` for legacy reasons but
// receipts are JSON-encoded on the wire — RLP only round-trips
// consensus fields and drops TxHash / BlockHash / BlockNumber /
// TransactionIndex / GasUsed / ContractAddress, which the caller
// needs.
type TxOpsServiceClient interface {
	SignAndSend(ctx context.Context, in *txopspb.SignAndSendRequest, opts ...grpc.CallOption) (*txopspb.SignAndSendResponse, error)
	BatchSignAndSend(ctx context.Context, in *txopspb.BatchSignAndSendRequest, opts ...grpc.CallOption) (*txopspb.BatchSignAndSendResponse, error)
	Call(ctx context.Context, in *txopspb.CallRequest, opts ...grpc.CallOption) (*txopspb.CallResponse, error)
}

// BatchInput is one leg of a BatchExecute call.
type BatchInput struct {
	MsgID   string
	Data    []byte
	Address common.Address
}

// BatchResult is the per-item outcome of BatchExecute.
type BatchResult struct {
	Receipt *types.Receipt
	Err     error
}

type CTSExecutor struct {
	txOpsClient TxOpsServiceClient
}

func NewCTSExecutor(
	txOpsClient TxOpsServiceClient,
) *CTSExecutor {
	return &CTSExecutor{
		txOpsClient: txOpsClient,
	}
}

func (e *CTSExecutor) Execute(ctx context.Context, id string, calldata []byte, address common.Address) (*ethtypes.Receipt, error) {
	req := &txopspb.SignAndSendRequest{
		Id:      id,
		Data:    calldata,
		Address: address.Hex(),
	}

	var (
		resp *txopspb.SignAndSendResponse
		err  error
	)
	resp, err = e.txOpsClient.SignAndSend(ctx, req)
	if err != nil {
		return nil, wrapTxOpsStatusError("cts execute", err)
	}

	switch r := resp.GetResult().(type) {
	case *txopspb.SignAndSendResponse_Receipt:
		var receipt types.Receipt
		if err := json.Unmarshal(r.Receipt, &receipt); err != nil {
			return nil, fmt.Errorf("cts decode receipt: %w", err)
		}
		return &receipt, nil
	case *txopspb.SignAndSendResponse_RevertData:
		return nil, &ErrorWithRevertData{revertData: r.RevertData}
	default:
		return nil, fmt.Errorf("cts execute: empty response from cts")
	}
}

func (e *CTSExecutor) Call(ctx context.Context, address common.Address, calldata []byte) ([]byte, error) {
	req := &txopspb.CallRequest{
		Data:    calldata,
		Address: address.Hex(),
	}

	var (
		resp *txopspb.CallResponse
		err  error
	)
	resp, err = e.txOpsClient.Call(ctx, req)
	if err != nil {
		return nil, wrapTxOpsStatusError("cts call", err)
	}

	switch r := resp.GetResult().(type) {
	case *txopspb.CallResponse_Value:
		return r.Value, nil
	case *txopspb.CallResponse_RevertData:
		return nil, &ErrorWithRevertData{revertData: r.RevertData}
	default:
		return nil, fmt.Errorf("cts call: empty response from cts")
	}
}

func (e *CTSExecutor) BatchExecute(ctx context.Context, items []BatchInput) (map[string]BatchResult, error) {
	pbItems := make([]*txopspb.BatchSignAndSendItem, len(items))
	for i, item := range items {
		pbItems[i] = &txopspb.BatchSignAndSendItem{
			MsgId:   item.MsgID,
			Data:    item.Data,
			Address: item.Address.Hex(),
		}
	}

	resp, err := e.txOpsClient.BatchSignAndSend(ctx, &txopspb.BatchSignAndSendRequest{Items: pbItems})
	if err != nil {
		return nil, wrapTxOpsStatusError("cts batch execute", err)
	}

	results := make(map[string]BatchResult, len(resp.GetResults()))
	for msgID, r := range resp.GetResults() {
		switch outcome := r.GetOutcome().(type) {
		case *txopspb.BatchSignAndSendItemResult_Receipt:
			var receipt types.Receipt
			if err := json.Unmarshal(outcome.Receipt, &receipt); err != nil {
				return nil, fmt.Errorf("cts batch decode receipt for %s: %w", msgID, err)
			}
			results[msgID] = BatchResult{Receipt: &receipt}
		case *txopspb.BatchSignAndSendItemResult_RevertData:
			results[msgID] = BatchResult{Err: &ErrorWithRevertData{revertData: outcome.RevertData}}
		case *txopspb.BatchSignAndSendItemResult_Error:
			results[msgID] = BatchResult{Err: fmt.Errorf("cts batch item %s: %s", msgID, outcome.Error.GetMessage())}
		}
	}

	return results, nil
}

func wrapTxOpsStatusError(op string, err error) error {
	if st, ok := status.FromError(err); ok {
		return fmt.Errorf("%s: %s: %s", op, codeName(st.Code()), st.Message())
	}
	return fmt.Errorf("%s: %w", op, err)
}

func codeName(c codes.Code) string {
	return c.String()
}

type deployerClient interface {
	Deploy(ctx context.Context, in *txopspb.DeployRequest, opts ...grpc.CallOption) (*txopspb.DeployResponse, error)
}

type CTSDeployer struct {
	client deployerClient
}

func NewCTSDeployer(client deployerClient) *CTSDeployer {
	return &CTSDeployer{
		client: client,
	}
}

func (e *CTSDeployer) Deploy(ctx context.Context, bytecode []byte, constructor []byte) (common.Address, *types.Receipt, error) {
	req := &txopspb.DeployRequest{
		Bytecode:    bytecode,
		Constructor: constructor,
	}

	var (
		resp *txopspb.DeployResponse
		err  error
	)
	resp, err = e.client.Deploy(ctx, req)
	if err != nil {
		return common.Address{}, nil, wrapTxOpsStatusError("cts deploy", err)
	}

	switch r := resp.GetResult().(type) {
	case *txopspb.DeployResponse_Success:
		var receipt types.Receipt
		if err := json.Unmarshal(r.Success.GetReceipt(), &receipt); err != nil {
			return common.Address{}, nil, fmt.Errorf("cts deploy: decode receipt: %w", err)
		}
		if r.Success.GetAddress() == "" {
			return common.Address{}, nil, fmt.Errorf("cts deploy: empty address in response")
		}
		if !common.IsHexAddress(r.Success.GetAddress()) {
			return common.Address{}, nil, fmt.Errorf("cts deploy: invalid hex address %q", r.Success.GetAddress())
		}
		addr := common.HexToAddress(r.Success.GetAddress())
		return addr, &receipt, nil
	case *txopspb.DeployResponse_RevertData:
		return common.Address{}, nil, &ErrorWithRevertData{revertData: r.RevertData}
	default:
		return common.Address{}, nil, fmt.Errorf("cts deploy: empty response from cts")
	}
}
