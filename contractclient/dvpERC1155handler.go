package contractclient

import (
	"context"
	"encoding/json"
	"math/big"
	"strconv"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/raylsnetwork/rayls-sovereign-relayer/contracts/RaylsErc1155DvpHandler"
	"github.com/raylsnetwork/rayls-sovereign-relayer/conv"
	"github.com/raylsnetwork/rayls-sovereign-relayer/dvp"
	"github.com/raylsnetwork/rayls-sovereign-relayer/types"
	"github.com/raylsnetwork/rayls-sovereign-relayer/withstack"
)

const dvpTxTimeout = 60 * time.Second

type DvpERC1155HandlerClient struct {
	contract *RaylsErc1155DvpHandler.RaylsErc1155DvpHandler
	executor Executor
}

func NewDvpERC1155HandlerClient(
	executor Executor,
) *DvpERC1155HandlerClient {
	return &DvpERC1155HandlerClient{
		contract: RaylsErc1155DvpHandler.NewRaylsErc1155DvpHandler(),
		executor: executor,
	}
}

func (c *DvpERC1155HandlerClient) GetAllTokenIdsWithSupply(ctx context.Context, tokenAddress common.Address) ([]dvp.DvpERC1155Supply, error) {
	calldata := c.contract.PackGetAllTokenIdsWithSupply()

	raw, err := c.executor.Call(ctx, tokenAddress, calldata)
	if err != nil {
		return nil, WrapInDvpERC1155HandlerClientError(
			"failed to get all token ids with supply",
			withstack.Wrap(err),
		)
	}

	result, err := c.contract.UnpackGetAllTokenIdsWithSupply(raw)
	if err != nil {
		return nil, WrapInDvpERC1155HandlerClientError(
			"failed to unpack all token ids with supply",
			withstack.Wrap(err),
		)
	}

	supplies := make([]dvp.DvpERC1155Supply, len(result))
	for i, supply := range result {
		supplies[i] = dvp.DvpERC1155Supply{
			Id:     supply.Id,
			Amount: supply.Amount,
		}
	}
	return supplies, nil
}

func (c *DvpERC1155HandlerClient) GetTokenExtraData(ctx context.Context, tokenAddress common.Address, tokenId *big.Int) ([]byte, error) {
	calldata := c.contract.PackGetTokenExtraData(tokenId)

	raw, err := c.executor.Call(ctx, tokenAddress, calldata)
	if err != nil {
		return nil, WrapInDvpERC1155HandlerClientError(
			"failed to get token extra data",
			withstack.Wrap(err),
		)
	}

	data, err := c.contract.UnpackGetTokenExtraData(raw)
	if err != nil {
		return nil, WrapInDvpERC1155HandlerClientError(
			"failed to unpack token extra data",
			withstack.Wrap(err),
		)
	}

	marshaledData, err := json.Marshal(data)
	if err != nil {
		return nil, WrapInDvpERC1155HandlerClientError("failed to marshal token extra data", err)
	}
	return marshaledData, nil
}

func (c *DvpERC1155HandlerClient) Unlock(
	ctx context.Context,
	chainEventID string,
	tokenAddress common.Address,
	tokenId *big.Int,
	tokenAmount *big.Int,
	to common.Address,
) error {
	calldata := c.contract.PackUnlockFromDvp(tokenId, tokenAmount, to)

	ctxTx, cancel := context.WithTimeout(ctx, dvpTxTimeout)
	defer cancel()

	_, err := c.executor.Execute(ctxTx, IDFor("dvpERC1155handler.Unlock", chainEventID), calldata, tokenAddress)
	if err != nil {
		return WrapInDvpERC1155HandlerClientError(
			"failed to unlock ERC1155 from Enygma DvP",
			withstack.Wrap(err),
		)
	}

	return nil
}

func (c *DvpERC1155HandlerClient) MarkSwapCompleted(
	ctx context.Context,
	tokenAddress common.Address,
	fromAddress common.Address,
	toAddress common.Address,
	fromChainId *big.Int,
	toChainId *big.Int,
	tokenResourceId string,
	tokenId string,
	tokenAmount *big.Int,
	data []byte,
	sharedId string,
) error {
	sharedIdBytes, err := conv.StringToBytes32(sharedId)
	if err != nil {
		return WrapInDvpERC1155HandlerClientError("failed to convert shared id to bytes32", err)
	}
	tokenIdBig, err := conv.StringToBigInt(tokenId)
	if err != nil {
		return WrapInDvpERC1155HandlerClientError("failed to convert token id to big int", err)
	}

	// Create params struct to avoid stack too deep problem
	params := RaylsErc1155DvpHandler.SharedObjectsDvpSwapCompletedParams{
		TokenId:            tokenIdBig,
		DestinationChainId: toChainId,
		DestinationOwner:   toAddress,
		SharedId:           sharedIdBytes,
	}

	ctxTx, cancel := context.WithTimeout(ctx, dvpTxTimeout)
	defer cancel()

	calldata := c.contract.PackDvpSwapCompleted(params, fromAddress, tokenAmount, data)

	_, err = c.executor.Execute(ctxTx, IDFor("dvpERC1155handler.MarkSwapCompleted", sharedId), calldata, tokenAddress)
	if err != nil {
		return WrapInDvpERC1155HandlerClientError("failed to complete dvp swap", withstack.Wrap(err))
	}

	return nil
}

func (c *DvpERC1155HandlerClient) NotifySenderWithPNCommunicator(
	ctx context.Context,
	tokenAddress common.Address,
	sharedId string,
	status types.DvpCommunicatorStatus,
	message string,
) error {
	sharedIdBytes, err := conv.StringToBytes32(sharedId)
	if err != nil {
		return WrapInDvpERC1155HandlerClientError("failed to convert shared id to bytes32", err)
	}

	calldata := c.contract.PackNotifySenderWithPNCommunicator(
		sharedIdBytes,
		uint8(status),
		message,
	)

	ctxTx, cancel := context.WithTimeout(ctx, dvpTxTimeout)
	defer cancel()

	_, err = c.executor.Execute(ctxTx, IDFor("dvpERC1155handler.NotifySender", sharedId, strconv.Itoa(int(status))), calldata, tokenAddress)
	if err != nil {
		return WrapInDvpERC1155HandlerClientError("failed to notify sender with pl communicator", withstack.Wrap(err))
	}

	return nil
}
