package contractclient

import (
	"context"
	"math/big"
	"strconv"

	"github.com/ethereum/go-ethereum/common"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/contracts/ProgrammabilityExecutorV1"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/contracts/RaylsEnygmaHandler"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/conv"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/types"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/withstack"
)

type EnygmaHandlerClient struct {
	contract *RaylsEnygmaHandler.RaylsEnygmaHandler
	executor BatchExecutor

	// Programmability: every inbound cross-chain transfer is dispatched as a
	// programData call to the destination PN's ProgrammabilityExecutor, which resolves the
	// target on-chain, gates it against the template registry, and target.call's crossMintStandard.
	programmabilityExecutor        *ProgrammabilityExecutorV1.ProgrammabilityExecutorV1
	programmabilityExecutorAddress common.Address
}

func NewEnygmaHandlerClient(
	executor BatchExecutor,
	programmabilityExecutorAddress common.Address,
) *EnygmaHandlerClient {
	return &EnygmaHandlerClient{
		contract:                       RaylsEnygmaHandler.NewRaylsEnygmaHandler(),
		executor:                       executor,
		programmabilityExecutor:        ProgrammabilityExecutorV1.NewProgrammabilityExecutorV1(),
		programmabilityExecutorAddress: programmabilityExecutorAddress,
	}
}

// ReceiveDestTransferBatch dispatches each recipient's programData array via one
// executeProgramData tx against the PN's ProgrammabilityExecutor. The blob array (built on
// the source side — a single [mintBlob] for plain transfers, or [mintBlob, userBlob...] for
// composed ones) is forwarded verbatim; the executor resolves each blob's target, gates it,
// and target.call's it (crossMintStandard for the mint blob). The token address is no longer
// the tx target — the executor resolves the token from each blob's resourceId.
//
// expectedMintTotal (per recipient: t.ToAmount) is the value the relayer received from the PNH
// batch for this transfer. The executor sums every crossMintStandard blob's _value and reverts
// the whole tx unless that sum equals this total — so the settlement mint cannot land an amount
// different from what the PNH authorized.
//
// INVARIANT: each loop iteration builds exactly one executeProgramData tx from one transfer
// tuple `t` — one source sender (t.FromAddress), one recipient, that tuple's own blob array.
// Tuples from different senders are aggregated only at the settlement-batch level and fanned
// back out to one tx per tuple here, never merged into a single blobs[] array. This is what
// makes a single per-call originSender / expectedMintTotal correct. Do NOT coalesce multiple
// tuples into one executeProgramData call without making originSender/expectedMintTotal per-blob.
func (c *EnygmaHandlerClient) ReceiveDestTransferBatch(ctx context.Context, transfers []*types.EnygmaCrossTransferData) (map[string]BatchResult, error) {
	items := make([]BatchInput, 0, len(transfers))
	for _, t := range transfers {
		// Owner-restricted programmability: the 3-arg `executeProgramData` overload attests the
		// originating EOA on the source chain (PNHTransfer.from, surfaced as
		// EnygmaTransferBatchTx.FromAddress here). The executor APPENDS this address as a trusted
		// 20-byte calldata tail on each userBlob's inner call (the userBlob's args carry ONLY the
		// target's leading params — no origin slot). The destination token reads the tail via
		// _getMsgSenderOnReceiveMethod and its in-body
		// hasContractScopedRole(TOKEN_OWNER, originSender, address(this)) check then rejects
		// non-owner attempts to mint via programmability. The settlement blob (crossMintStandard)
		// is exempt — it carries no origin and the executor dispatches it WITHOUT the tail.
		//
		// t.ProgramData carries the domain struct; the executor's Pack method wants the
		// identically-shaped struct from its own binding package (abigen emits a distinct named
		// type per package), so convert element-by-element.
		steps := make([]ProgrammabilityExecutorV1.SharedObjectsEnygmaProgramData, len(t.ProgramData))
		for j, s := range t.ProgramData {
			steps[j] = ProgrammabilityExecutorV1.SharedObjectsEnygmaProgramData{
				ResourceId:      s.ResourceId,
				ContractAddress: s.ContractAddress,
				Selector:        s.Selector,
				Args:            s.Args,
			}
		}
		calldata := c.programmabilityExecutor.PackExecuteProgramData(steps, t.ToAmount, t.FromAddress)
		items = append(items, BatchInput{
			MsgID:   t.GetID(),
			Data:    calldata,
			Address: c.programmabilityExecutorAddress,
		})
	}

	results, err := c.executor.BatchExecute(ctx, items)
	if err != nil {
		return nil, WrapInEnygmaHandlerClientError("failed to batch executeProgramData", withstack.Wrap(err))
	}

	return results, nil
}

func (c *EnygmaHandlerClient) RevertDestTransferBatch(ctx context.Context, tokenAddress common.Address, reverts []*types.EnygmaCrossTransferData) (map[string]BatchResult, error) {
	items := make([]BatchInput, 0, len(reverts))
	for _, r := range reverts {
		calldata := c.contract.PackCrossTransferRevertBatch(
			r.FromAddress,
			r.ToAddress,
			r.ToAmount,
			r.FromChainID,
			r.ReferenceId,
		)
		items = append(items, BatchInput{
			MsgID:   r.GetID(),
			Data:    calldata,
			Address: tokenAddress,
		})
	}

	results, err := c.executor.BatchExecute(ctx, items)
	if err != nil {
		return nil, WrapInEnygmaHandlerClientError("failed to batch revert", withstack.Wrap(err))
	}

	return results, nil
}

func (c *EnygmaHandlerClient) RevertSrcTransferBatch(ctx context.Context, revertTxs []*types.EnygmaTransferFailed) (map[string]BatchResult, error) {
	items := make([]BatchInput, 0, len(revertTxs))
	for _, tx := range revertTxs {
		// Issue #75 idempotency: the per-event ReferenceID is the dedup key
		// consumed by RaylsEnygmaHandler.crossRevertMint so duplicate retries
		// after a crash + NATS redelivery are a silent on-chain no-op.
		calldata := c.contract.PackCrossRevertMint(tx.Sender, tx.Amount, tx.Reason, tx.ReferenceID)
		items = append(items, BatchInput{
			MsgID:   tx.GetID(),
			Data:    calldata,
			Address: tx.EnygmaAddress,
		})
	}

	results, err := c.executor.BatchExecute(ctx, items)
	if err != nil {
		return nil, WrapInEnygmaHandlerClientError("failed to batch revert src transfer", withstack.Wrap(err))
	}

	return results, nil
}

func (c *EnygmaHandlerClient) RevertSrcSupplyBatch(ctx context.Context, revertTxs []*types.EnygmaSupplyUpdateFailed) (map[string]BatchResult, error) {
	items := make([]BatchInput, 0, len(revertTxs))
	for _, tx := range revertTxs {
		calldata := c.contract.PackSupplyUpdateRevert(tx.Amount, tx.To, tx.Type == types.EnygmaMint)
		items = append(items, BatchInput{
			MsgID:   tx.GetID(),
			Data:    calldata,
			Address: tx.EnygmaAddress,
		})
	}

	results, err := c.executor.BatchExecute(ctx, items)
	if err != nil {
		return nil, WrapInEnygmaHandlerClientError("failed to batch revert src supply", withstack.Wrap(err))
	}

	return results, nil
}

func (c *EnygmaHandlerClient) ReceiveWithdraw(
	ctx context.Context,
	chainEventID string,
	tokenAddress common.Address,
	to common.Address,
	value *big.Int,
	referenceId [32]byte,
) error {
	calldata := c.contract.PackReceiveWithdrawFromDvp(to, value, referenceId)

	_, err := c.executor.Execute(ctx, IDFor("enygmahandler.ReceiveWithdraw", chainEventID), calldata, tokenAddress)
	if err != nil {
		return WrapInEnygmaHandlerClientError("failed to receive withdraw from dvp", withstack.Wrap(err))
	}

	return nil
}

func (c *EnygmaHandlerClient) MarkSwapCompleted(
	ctx context.Context,
	tokenAddress common.Address,
	destinationChainId *big.Int,
	sharedId string,
) error {
	sharedIdBytes, err := conv.StringToBytes32(sharedId)
	if err != nil {
		return WrapInEnygmaHandlerClientError("failed to convert shared id to bytes32", err)
	}

	calldata := c.contract.PackDvpSwapCompleted(destinationChainId, sharedIdBytes)

	_, err = c.executor.Execute(ctx, IDFor("enygmahandler.MarkSwapCompleted", sharedId), calldata, tokenAddress)
	if err != nil {
		return WrapInEnygmaHandlerClientError("failed to complete dvp swap", withstack.Wrap(err))
	}

	return nil
}

func (c *EnygmaHandlerClient) NotifySenderWithPNCommunicator(
	ctx context.Context,
	tokenAddress common.Address,
	sharedId string,
	status types.DvpCommunicatorStatus,
	message string,
) error {
	sharedIdBytes, err := conv.StringToBytes32(sharedId)
	if err != nil {
		return WrapInEnygmaHandlerClientError("failed to convert shared id to bytes32", err)
	}

	calldata := c.contract.PackNotifySenderWithPNCommunicator(sharedIdBytes, uint8(status), message)

	_, err = c.executor.Execute(ctx, IDFor("enygmahandler.NotifySender", sharedId, strconv.Itoa(int(status))), calldata, tokenAddress)
	if err != nil {
		return WrapInEnygmaHandlerClientError("failed to notify sender with pl communicator", withstack.Wrap(err))
	}

	return nil
}
