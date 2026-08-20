package contractclient

import (
	"context"
	"encoding/json"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/raylsnetwork/rayls-sovereign-relayer/contracts/DvpErc1155PNH"
	"github.com/raylsnetwork/rayls-sovereign-relayer/dvp"
	"github.com/raylsnetwork/rayls-sovereign-relayer/withstack"
)

type dvpERC1155Encryptor interface {
	EncryptDvpExtraData(ctx context.Context, rawDataBytes []byte) ([]byte, error)
}

type DvpERC1155Client struct {
	contract *DvpErc1155PNH.DvpErc1155PNH
	encr     dvpERC1155Encryptor
	executor BatchExecutor
}

func NewDvpERC1155Client(
	executor BatchExecutor,
	encr dvpERC1155Encryptor,
) *DvpERC1155Client {
	return &DvpERC1155Client{
		contract: DvpErc1155PNH.NewDvpErc1155PNH(),
		encr:     encr,
		executor: executor,
	}
}

func (c *DvpERC1155Client) Burn(
	ctx context.Context,
	chainEventID string,
	tokenAddress common.Address,
	tokenOwner common.Address,
	tokenId *big.Int,
	tokenAmount *big.Int,
) error {
	calldata := c.contract.PackBurn(tokenOwner, tokenId, tokenAmount)

	_, err := c.executor.Execute(ctx, IDFor("dvpERC1155.Burn", chainEventID), calldata, tokenAddress)
	if err != nil {
		return WrapInDvpERC1155ClientError("failed to burn ERC1155", withstack.Wrap(err))
	}

	return nil
}

func (c *DvpERC1155Client) Approve(ctx context.Context, chainEventID string, tokenAddress common.Address, to common.Address) error {
	calldata := c.contract.PackSetApprovalForAll(to, true)

	_, err := c.executor.Execute(ctx, IDFor("dvpERC1155.Approve", chainEventID), calldata, tokenAddress)
	if err != nil {
		return WrapInDvpERC1155ClientError("failed to approve ERC1155", withstack.Wrap(err))
	}

	return nil
}

func (c *DvpERC1155Client) UpdateExtraData(
	ctx context.Context,
	chainEventID string,
	tokenAddress common.Address,
	tokenId *big.Int,
	tokenAmount *big.Int,
	chainId *big.Int,
	extraDataBytes []byte,
	newOwner common.Address,
) error {
	encrData, err := c.encr.EncryptDvpExtraData(ctx, extraDataBytes)
	if err != nil {
		return WrapInDvpERC1155ClientError("failed to encrypt dvp ERC1155 extra data", err)
	}

	var extraData []DvpErc1155PNH.SharedObjectsDvp1155ExtraData
	err = json.Unmarshal(encrData, &extraData)
	if err != nil {
		return WrapInDvpERC1155ClientError("failed to unmarshal extra data", err)
	}

	tokenIds := []*big.Int{tokenId}
	tokenAmounts := []*big.Int{tokenAmount}
	calldata := c.contract.PackUpdateInfosAfterDvpWithdraw(tokenIds, tokenAmounts, chainId, extraData, newOwner)

	_, err = c.executor.Execute(ctx, IDFor("dvpERC1155.UpdateExtraData", chainEventID), calldata, tokenAddress)
	if err != nil {
		return WrapInDvpERC1155ClientError("failed to update ERC1155 extra data", withstack.Wrap(err))
	}

	return nil
}

func (c *DvpERC1155Client) MintBatch(ctx context.Context, mintDatas []*dvp.DvpERC1155MintData) (map[string]BatchResult, error) {
	items := make([]BatchInput, 0, len(mintDatas))
	for _, d := range mintDatas {
		encrData, err := c.encr.EncryptDvpExtraData(ctx, d.ExtraData)
		if err != nil {
			return nil, WrapInDvpERC1155ClientError("failed to encrypt ERC1155 extra data", err)
		}

		var extraData []DvpErc1155PNH.SharedObjectsDvp1155ExtraData
		if err := json.Unmarshal(encrData, &extraData); err != nil {
			return nil, WrapInDvpERC1155ClientError("failed to unmarshal extra data", err)
		}

		calldata := c.contract.PackMint(d.To, d.TokenID, d.TokenAmount, d.Data, d.ChainID, extraData)
		items = append(items, BatchInput{
			MsgID:   d.GetID(),
			Data:    calldata,
			Address: d.TokenAddress,
		})
	}

	results, err := c.executor.BatchExecute(ctx, items)
	if err != nil {
		return nil, WrapInDvpERC1155ClientError("failed to batch mint ERC1155", withstack.Wrap(err))
	}

	return results, nil
}
