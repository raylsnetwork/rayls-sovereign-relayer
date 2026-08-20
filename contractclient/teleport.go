// Decommissioning Teleport (vanilla, atomic): atomic TeleportClient methods below are deprecated; the type and StoreEncryptedDataBatch (shared) are retained.

package contractclient

import (
	"context"
	"errors"
	"log/slog"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/google/uuid"
	"github.com/raylsnetwork/rayls-sovereign-relayer/contracts/TeleportV1"
	sharedservice "github.com/raylsnetwork/rayls-sovereign-relayer/private-relayer/shared/service"
	"github.com/raylsnetwork/rayls-sovereign-relayer/txsim"
	"github.com/raylsnetwork/rayls-sovereign-relayer/types"
	"github.com/raylsnetwork/rayls-sovereign-relayer/withstack"
)

type encryptor interface {
	EncryptAdditionalData(context.Context, []types.AtomicTeleportAdditionalData) (string, error)
	EncryptMessages(context.Context, []types.DispatchedMessageToPrivateHub, *big.Int) (EncryptedMessages, error)
}

type simulator interface {
	DecodeRevertReason(context.Context, error) (txsim.ContractError, error)
}

type TeleportClient struct {
	address  common.Address
	contract *TeleportV1.TeleportV1
	executor Executor
	encr     encryptor
}

func NewTeleportClient(
	address common.Address,
	executor Executor,
	encr encryptor,
) *TeleportClient {
	return &TeleportClient{
		address:  address,
		contract: TeleportV1.NewTeleportV1(),
		executor: executor,
		encr:     encr,
	}
}

// Deprecated: Decommissioning Teleport (vanilla, atomic).
func (c *TeleportClient) SendAdditionalDataBatch(ctx context.Context, sharedIDs []string, data []types.AtomicTeleportAdditionalData) error {
	encrData, err := c.encr.EncryptAdditionalData(ctx, data)
	if err != nil {
		return WrapInTeleportClientError("failed to encrypt additional data", err)
	}

	calldata := c.contract.PackEmitAdditionalAtomicDataBatchFor(sharedIDs, encrData)

	// Best-effort key: opportunistic batch membership isn't stable across restarts;
	// per-message double-effect is guarded on-chain.
	_, err = c.executor.Execute(ctx, IDFor("teleport.SendAdditionalDataBatch", HashIDs(sharedIDs)), calldata, c.address)
	if err != nil {
		return WrapInTeleportClientError("failed to emit additional atomic data batch", withstack.Wrap(err))
	}
	return nil
}

// Deprecated: Decommissioning Teleport (vanilla, atomic).
func (c *TeleportClient) ExecuteAtomicMessageBatch(
	ctx context.Context,
	sharedIDs []string,
	data []types.AtomicTeleportAdditionalData,
) error {
	encrData, err := c.encr.EncryptAdditionalData(ctx, data)
	if err != nil {
		return WrapInTeleportClientError("failed to encrypt additional data", err)
	}

	calldata := c.contract.PackExecuteAtomicMessageBatch(sharedIDs, encrData)

	// Best-effort key: opportunistic batch membership isn't stable across restarts;
	// per-message double-effect is guarded on-chain.
	_, err = c.executor.Execute(ctx, IDFor("teleport.ExecuteAtomicMessageBatch", HashIDs(sharedIDs)), calldata, c.address)
	if err != nil {
		if IsRevertWithSelector(err, TeleportV1.TeleportV1TeleportV1MessageAlreadyRevertedErrorID()) {
			return sharedservice.ErrAlreadyReverted
		}
		return WrapInTeleportClientError("failed to execute atomic message batch", withstack.Wrap(err))
	}
	return nil
}

// Deprecated: Decommissioning Teleport (vanilla, atomic).
func (c *TeleportClient) RevertAtomicMessageBatch(ctx context.Context, sharedIDs []string, data []types.AtomicTeleportAdditionalData) error {
	encrData, err := c.encr.EncryptAdditionalData(ctx, data)
	if err != nil {
		return WrapInTeleportClientError("failed to encrypt additional data", err)
	}

	calldata := c.contract.PackRevertAtomicMessageBatch(sharedIDs, encrData)
	// Best-effort key: opportunistic batch membership isn't stable across restarts;
	// per-message double-effect is guarded on-chain.
	_, err = c.executor.Execute(ctx, IDFor("teleport.RevertAtomicMessageBatch", HashIDs(sharedIDs)), calldata, c.address)
	if err != nil {
		var rd *ErrorWithRevertData
		if errors.As(err, &rd) {
			slog.Warn("RevertAtomicMessageBatch", slog.Any("revertData", common.Bytes2Hex(rd.GetRevertData())))
		}
		if IsRevertWithSelector(err, TeleportV1.TeleportV1TeleportV1MessageAlreadyExecutedErrorID()) {
			return sharedservice.ErrAlreadyExecuted
		}
		return WrapInTeleportClientError("failed to revert atomic message batch", withstack.Wrap(err))
	}
	return nil
}

// Shared: StoreEncryptedDataBatch also serves the generic non-atomic relay; do not deprecate.
func (c *TeleportClient) StoreEncryptedDataBatch(
	ctx context.Context,
	sharedIDs []string,
	messages []types.DispatchedMessageToPrivateHub,
	chainID *big.Int,
) (common.Hash, error) {
	batchID, err := uuid.NewUUID()
	if err != nil {
		return common.Hash{}, WrapInTeleportClientError("failed to generate batch ID", withstack.Wrap(err))
	}

	encrMessages, err := c.encr.EncryptMessages(ctx, messages, chainID)
	if err != nil {
		return common.Hash{}, WrapInTeleportClientError("failed to encrypt messages", err)
	}

	dataBatch := TeleportV1.TeleportV1dataBatch{
		MessageTag: encrMessages.MessageTag,
		Data:       encrMessages.Data,
		BatchId:    batchID.String(),
		SharedIds:  sharedIDs,
	}

	// Best-effort key: opportunistic batch membership isn't stable across restarts;
	// per-message double-effect is guarded on-chain. Keyed on sharedIDs, not
	// calldata, since dataBatch embeds a non-deterministic batchID UUID.
	receipt, err := c.executor.Execute(ctx, IDFor("teleport.StoreEncryptedDataBatch", HashIDs(sharedIDs)), c.contract.PackStoreEncryptedDataBatch(dataBatch, encrMessages.BlockNumber), c.address)
	if err != nil {
		return common.Hash{}, WrapInTeleportClientError("failed to store encrypted data batch", withstack.Wrap(err))
	}

	return receipt.TxHash, nil
}

// Deprecated: Decommissioning Teleport (vanilla, atomic).
func (c *TeleportClient) GetAtomicMessageStatuses(ctx context.Context, sharedIDs []string) ([]types.AtomicStatusUpdateMessage, error) {
	calldata := c.contract.PackGetAtomicMessageStatuses(sharedIDs)

	raw, err := c.executor.Call(ctx, c.address, calldata)
	if err != nil {
		return nil, WrapInTeleportClientError("failed to get atomic message statuses", withstack.Wrap(err))
	}

	teleportSUMs, err := c.contract.UnpackGetAtomicMessageStatuses(raw)
	if err != nil {
		return nil, WrapInTeleportClientError("failed to unpack atomic message statuses", withstack.Wrap(err))
	}

	var sums []types.AtomicStatusUpdateMessage
	for _, teleportSUM := range teleportSUMs {
		status, err := getAtomicStatusFromString(teleportSUM.Status)
		if err != nil {
			return nil, ErrUnknownStatus
		}

		sums = append(sums, types.AtomicStatusUpdateMessage{
			SharedID: teleportSUM.MsgId,
			Status:   status,
		})
	}

	return sums, nil
}

// Decommissioning Teleport (vanilla, atomic).
func getAtomicStatusFromString(statusStr string) (types.AtomicStatus, error) {
	switch statusStr {
	case types.AtomicPending.String():
		return types.AtomicPendingStatus, nil
	case types.AtomicExecuted.String():
		return types.AtomicExecutedStatus, nil
	case types.AtomicRejected.String():
		return types.AtomicRejectedStatus, nil
	case types.AtomicReverted.String():
		return types.AtomicRevertedStatus, nil
	default:
		return 0, ErrUnknownStatus
	}
}
