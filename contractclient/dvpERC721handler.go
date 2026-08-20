package contractclient

import (
	"context"
	"encoding/json"
	"math/big"
	"strconv"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/raylsnetwork/rayls-sovereign-relayer/contracts/RaylsErc721DvpHandler"
	"github.com/raylsnetwork/rayls-sovereign-relayer/conv"
	"github.com/raylsnetwork/rayls-sovereign-relayer/types"
	"github.com/raylsnetwork/rayls-sovereign-relayer/withstack"
)

const dvpERC721TxReceiptTimeout = 60 * time.Second

type DvpERC721HandlerClient struct {
	contract *RaylsErc721DvpHandler.RaylsErc721DvpHandler
	executor Executor
}

func NewDvpERC721HandlerClient(
	executor Executor,
) *DvpERC721HandlerClient {
	return &DvpERC721HandlerClient{
		contract: RaylsErc721DvpHandler.NewRaylsErc721DvpHandler(),
		executor: executor,
	}
}

func (c *DvpERC721HandlerClient) GetTotalSupply(ctx context.Context, tokenAddress common.Address) ([]*big.Int, error) {
	calldata := c.contract.PackGetTotalSupplyAtPN()

	raw, err := c.executor.Call(ctx, tokenAddress, calldata)
	if err != nil {
		return nil, WrapInDvpERC721HandlerClientError("failed to get total supply", withstack.Wrap(err))
	}

	supply, err := c.contract.UnpackGetTotalSupplyAtPN(raw)
	if err != nil {
		return nil, WrapInDvpERC721HandlerClientError("failed to unpack total supply", withstack.Wrap(err))
	}
	return supply, nil
}

func (c *DvpERC721HandlerClient) GetExtraData(ctx context.Context, tokenAddress common.Address, tokenId *big.Int) ([]byte, error) {
	calldata := c.contract.PackGetNftExtradaData(tokenId)

	raw, err := c.executor.Call(ctx, tokenAddress, calldata)
	if err != nil {
		return nil, WrapInDvpERC721HandlerClientError("failed to get extra data", withstack.Wrap(err))
	}

	data, err := c.contract.UnpackGetNftExtradaData(raw)
	if err != nil {
		return nil, WrapInDvpERC721HandlerClientError("failed to unpack extra data", withstack.Wrap(err))
	}

	marshaled, err := json.Marshal(data)
	if err != nil {
		return nil, WrapInDvpERC721HandlerClientError("failed to marshal extra data to json", err)
	}
	return marshaled, nil
}

func (c *DvpERC721HandlerClient) Unlock(ctx context.Context, chainEvenID string, tokenAddress common.Address, tokenId *big.Int) error {
	ctxTx, cancel := context.WithTimeout(ctx, dvpERC721TxReceiptTimeout)
	defer cancel()

	calldata := c.contract.PackUnlockFromDvp(tokenId)

	_, err := c.executor.Execute(ctxTx, IDFor("dvpERC721handler.Unlock", chainEvenID), calldata, tokenAddress)
	if err != nil {
		return WrapInDvpERC721HandlerClientError("failed to unlock ERC721 from Enygma DvP", withstack.Wrap(err))
	}

	return nil
}

func (c *DvpERC721HandlerClient) MarkSwapCompleted(
	ctx context.Context,
	tokenAddress common.Address,
	toAddress common.Address,
	fromChainId *big.Int,
	toChainId *big.Int,
	tokenResourceId string,
	tokenId string,
	tokenAmount *big.Int,
	sharedId string,
) error {
	sharedIdBytes, err := conv.StringToBytes32(sharedId)
	if err != nil {
		return WrapInDvpERC721HandlerClientError("failed to convert shared id to bytes32", err)
	}
	tokenIdBig, err := conv.StringToBigInt(tokenId)
	if err != nil {
		return WrapInDvpERC721HandlerClientError("failed to convert token id to big int", err)
	}

	// Create params struct to avoid stack too deep problem.
	params := RaylsErc721DvpHandler.SharedObjectsDvpSwapCompletedParams{
		TokenId:            tokenIdBig,
		DestinationChainId: toChainId,
		DestinationOwner:   toAddress,
		SharedId:           sharedIdBytes,
	}

	ctxTx, cancel := context.WithTimeout(ctx, dvpERC721TxReceiptTimeout)
	defer cancel()

	calldata := c.contract.PackDvpSwapCompleted(params)

	_, err = c.executor.Execute(ctxTx, IDFor("dvpERC721handler.MarkSwapCompleted", sharedId), calldata, tokenAddress)
	if err != nil {
		return WrapInDvpERC721HandlerClientError("failed to complete dvp swap", withstack.Wrap(err))
	}

	return nil
}

func (c *DvpERC721HandlerClient) NotifySenderWithPNCommunicator(
	ctx context.Context,
	tokenAddress common.Address,
	sharedId string,
	status types.DvpCommunicatorStatus,
	message string,
) error {
	sharedIdBytes, err := conv.StringToBytes32(sharedId)
	if err != nil {
		return WrapInDvpERC721HandlerClientError("failed to convert shared id to bytes32", err)
	}

	calldata := c.contract.PackNotifySenderWithPNCommunicator(sharedIdBytes, uint8(status), message)
	if err != nil {
		return WrapInDvpERC721HandlerClientError("failed to notify sender with pl communicator", withstack.Wrap(err))
	}

	ctxTx, cancel := context.WithTimeout(ctx, dvpERC721TxReceiptTimeout)
	defer cancel()

	_, err = c.executor.Execute(ctxTx, IDFor("dvpERC721handler.NotifySender", sharedId, strconv.Itoa(int(status))), calldata, tokenAddress)
	if err != nil {
		return WrapInDvpERC721HandlerClientError("failed to get notify sender with pl communicator receipt", err)
	}

	return nil
}
