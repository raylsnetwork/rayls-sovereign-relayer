package grpc

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/ethereum/go-ethereum/common"
	txopspb "github.com/raylsnetwork/rayls-sovereign-relayer/cts/gen/proto/txops"
	"github.com/raylsnetwork/rayls-sovereign-relayer/cts/service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type TxOpsService interface {
	SignAndSend(ctx context.Context, id string, data []byte, address common.Address) (service.SignAndSendResult, error)
	BatchSignAndSend(ctx context.Context, items []service.BatchItem) (map[string]service.BatchItemResult, error)
	Call(ctx context.Context, data []byte, address common.Address) (service.CallResult, error)
	Deploy(ctx context.Context, bytecode, constructor []byte) (service.DeployResult, error)
}

type TxOpsHandler struct {
	svc TxOpsService
}

func NewTxOpsHandler(svc TxOpsService) *TxOpsHandler {
	return &TxOpsHandler{svc: svc}
}

func (h *TxOpsHandler) SignAndSend(ctx context.Context, req *txopspb.SignAndSendRequest) (*txopspb.SignAndSendResponse, error) {
	if req.GetId() == "" {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}
	if len(req.GetData()) == 0 {
		return nil, status.Error(codes.InvalidArgument, "data is required")
	}
	if req.GetAddress() == "" {
		return nil, status.Error(codes.InvalidArgument, "address is required")
	}
	if !common.IsHexAddress(req.GetAddress()) {
		return nil, status.Error(codes.InvalidArgument, "address must be a valid hex address")
	}
	addr := common.HexToAddress(req.GetAddress())

	result, err := h.svc.SignAndSend(ctx, req.GetId(), req.GetData(), addr)
	if err != nil {
		statusCode := codes.Internal
		if ctx.Err() != nil || errors.Is(err, service.ErrRetriable) {
			// signal to the caller that the error is retriable
			statusCode = codes.Unavailable
		}
		return nil, status.Errorf(statusCode, "sign and send: %v", err)
	}

	switch {
	case result.Success != nil:
		receiptJSON, err := json.Marshal(result.Success.Receipt)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "sign and send: encode receipt: %v", err)
		}
		return &txopspb.SignAndSendResponse{
			Result: &txopspb.SignAndSendResponse_Receipt{Receipt: receiptJSON},
		}, nil
	case result.Revert != nil:
		return &txopspb.SignAndSendResponse{
			Result: &txopspb.SignAndSendResponse_RevertData{RevertData: result.Revert.RevertData},
		}, nil
	default:
		// Service contract violation: non-error result with neither arm set.
		return nil, status.Error(codes.Internal, "sign and send: empty result")
	}
}

func (h *TxOpsHandler) Deploy(ctx context.Context, req *txopspb.DeployRequest) (*txopspb.DeployResponse, error) {
	if len(req.GetBytecode()) == 0 {
		return nil, status.Error(codes.InvalidArgument, "bytecode is required")
	}

	result, err := h.svc.Deploy(ctx, req.GetBytecode(), req.GetConstructor())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "deploy: %v", err)
	}

	switch {
	case result.Success != nil:
		receiptJSON, err := json.Marshal(result.Success.Receipt)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "deploy: encode receipt: %v", err)
		}
		return &txopspb.DeployResponse{
			Result: &txopspb.DeployResponse_Success{
				Success: &txopspb.DeploySuccess{
					Receipt: receiptJSON,
					Address: result.Success.Address.Hex(),
				},
			},
		}, nil
	case result.Revert != nil:
		return &txopspb.DeployResponse{
			Result: &txopspb.DeployResponse_RevertData{RevertData: result.Revert.RevertData},
		}, nil
	default:
		return nil, status.Error(codes.Internal, "deploy: empty result")
	}
}

func (h *TxOpsHandler) Call(ctx context.Context, req *txopspb.CallRequest) (*txopspb.CallResponse, error) {
	if len(req.GetData()) == 0 {
		return nil, status.Error(codes.InvalidArgument, "data is required")
	}
	if req.GetAddress() == "" {
		return nil, status.Error(codes.InvalidArgument, "address is required")
	}
	if !common.IsHexAddress(req.GetAddress()) {
		return nil, status.Error(codes.InvalidArgument, "address must be a valid hex address")
	}
	addr := common.HexToAddress(req.GetAddress())

	result, err := h.svc.Call(ctx, req.GetData(), addr)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "call: %v", err)
	}

	// Note: Value is an empty-but-non-nil slice for functions that return
	// zero bytes, so checking Revert first disambiguates cleanly.
	switch {
	case result.Revert != nil:
		return &txopspb.CallResponse{
			Result: &txopspb.CallResponse_RevertData{RevertData: result.Revert.RevertData},
		}, nil
	case result.Value != nil:
		return &txopspb.CallResponse{
			Result: &txopspb.CallResponse_Value{Value: result.Value},
		}, nil
	default:
		// Service contract violation: non-error result with neither arm set.
		return nil, status.Error(codes.Internal, "call: empty result")
	}
}

func (h *TxOpsHandler) BatchSignAndSend(ctx context.Context, req *txopspb.BatchSignAndSendRequest) (*txopspb.BatchSignAndSendResponse, error) {
	if len(req.GetItems()) == 0 {
		return nil, status.Error(codes.InvalidArgument, "at least one item is required")
	}

	seen := make(map[string]struct{}, len(req.GetItems()))
	items := make([]service.BatchItem, 0, len(req.GetItems()))
	for _, item := range req.GetItems() {
		if item.GetMsgId() == "" {
			return nil, status.Error(codes.InvalidArgument, "msg_id is required for every item")
		}
		if _, dup := seen[item.GetMsgId()]; dup {
			return nil, status.Errorf(codes.InvalidArgument, "duplicate msg_id: %s", item.GetMsgId())
		}
		seen[item.GetMsgId()] = struct{}{}
		if len(item.GetData()) == 0 {
			return nil, status.Errorf(codes.InvalidArgument, "data is required for item %s", item.GetMsgId())
		}
		if item.GetAddress() == "" {
			return nil, status.Errorf(codes.InvalidArgument, "address is required for item %s", item.GetMsgId())
		}
		if !common.IsHexAddress(item.GetAddress()) {
			return nil, status.Errorf(codes.InvalidArgument, "address must be a valid hex address for item %s", item.GetMsgId())
		}
		items = append(items, service.BatchItem{
			MsgID:   item.GetMsgId(),
			Data:    item.GetData(),
			Address: common.HexToAddress(item.GetAddress()),
		})
	}

	results, err := h.svc.BatchSignAndSend(ctx, items)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "batch sign and send: %v", err)
	}

	pbResults := make(map[string]*txopspb.BatchSignAndSendItemResult, len(results))
	for msgID, r := range results {
		switch {
		case r.Success != nil:
			receiptJSON, err := json.Marshal(r.Success.Receipt)
			if err != nil {
				return nil, status.Errorf(codes.Internal, "batch sign and send: encode receipt for %s: %v", msgID, err)
			}
			pbResults[msgID] = &txopspb.BatchSignAndSendItemResult{
				Outcome: &txopspb.BatchSignAndSendItemResult_Receipt{Receipt: receiptJSON},
			}
		case r.Revert != nil:
			pbResults[msgID] = &txopspb.BatchSignAndSendItemResult{
				Outcome: &txopspb.BatchSignAndSendItemResult_RevertData{RevertData: r.Revert.RevertData},
			}
		case r.Err != nil:
			pbResults[msgID] = &txopspb.BatchSignAndSendItemResult{
				Outcome: &txopspb.BatchSignAndSendItemResult_Error{
					Error: &txopspb.BatchItemError{
						Message: r.Err.Error(),
					},
				},
			}
		}
	}

	return &txopspb.BatchSignAndSendResponse{Results: pbResults}, nil
}
