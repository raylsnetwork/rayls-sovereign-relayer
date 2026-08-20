package contractclient

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/raylsnetwork/rayls-sovereign-relayer/contracts/DvpErc721PNH"
	"github.com/raylsnetwork/rayls-sovereign-relayer/dvp"
	"github.com/raylsnetwork/rayls-sovereign-relayer/withstack"
)

type dvpERC721Encryptor interface {
	EncryptDvpExtraData(ctx context.Context, rawDataBytes []byte) ([]byte, error)
}

type DvpERC721Client struct {
	contract *DvpErc721PNH.DvpErc721PNH
	encr     dvpERC721Encryptor
	executor BatchExecutor
}

func NewDvpERC721Client(
	executor BatchExecutor,
	encr dvpERC721Encryptor,
) *DvpERC721Client {
	return &DvpERC721Client{
		contract: DvpErc721PNH.NewDvpErc721PNH(),
		encr:     encr,
		executor: executor,
	}
}

func (c *DvpERC721Client) Burn(ctx context.Context, chainEventID string, tokenAddress common.Address, nftId *big.Int) error {
	calldata := c.contract.PackBurn(nftId)

	_, err := c.executor.Execute(ctx, IDFor("dvpERC721.Burn", chainEventID), calldata, tokenAddress)
	if err != nil {
		return WrapInDvpERC721ClientError("failed to burn ERC721", withstack.Wrap(err))
	}

	return nil
}

func (c *DvpERC721Client) Approve(ctx context.Context, chainEventID string, tokenAddress common.Address, to common.Address, nftId *big.Int) error {
	calldata := c.contract.PackApprove(to, nftId)

	_, err := c.executor.Execute(ctx, IDFor("dvpERC721.Approve", chainEventID), calldata, tokenAddress)
	if err != nil {
		var rd *ErrorWithRevertData
		if errors.As(err, &rd) {
			if reason, unpackErr := WithVanillaRevert(rd.GetRevertData(), c.contract.UnpackError); unpackErr == nil {
				return NewDvpERC721ClientError(fmt.Sprintf("approve failed with revert reason: %v", reason))
			}
		}
		return WrapInDvpERC721ClientError("failed to approve ERC721", withstack.Wrap(err))
	}

	return nil
}

func (c *DvpERC721Client) UpdateExtraData(
	ctx context.Context,
	chainEventID string,
	tokenAddress common.Address,
	nftId *big.Int,
	chainId *big.Int,
	extraDataBytes []byte,
	newOwner common.Address,
) error {
	encrData, err := c.encr.EncryptDvpExtraData(ctx, extraDataBytes)
	if err != nil {
		return WrapInDvpERC721ClientError("failed to encrypt ERC721 extra data", err)
	}

	var extraData []DvpErc721PNH.SharedObjectsDvp721ExtraData
	err = json.Unmarshal(encrData, &extraData)
	if err != nil {
		return WrapInDvpERC721ClientError("failed to unmarshal extra data", err)
	}

	calldata := c.contract.PackUpdateInfosAfterDvpWithdraw(nftId, chainId, extraData, newOwner)

	_, err = c.executor.Execute(ctx, IDFor("dvpERC721.UpdateExtraData", chainEventID), calldata, tokenAddress)
	if err != nil {
		return WrapInDvpERC721ClientError("failed to update ERC721 infos", withstack.Wrap(err))
	}

	return nil
}

func (c *DvpERC721Client) MintBatch(ctx context.Context, mintDatas []*dvp.DvpERC721MintData) (map[string]BatchResult, error) {
	items := make([]BatchInput, 0, len(mintDatas))
	for _, d := range mintDatas {
		encrData, err := c.encr.EncryptDvpExtraData(ctx, d.ExtraData)
		if err != nil {
			return nil, WrapInDvpERC721ClientError("failed to encrypt ERC721 extra data", err)
		}

		var extraData []DvpErc721PNH.SharedObjectsDvp721ExtraData
		if err := json.Unmarshal(encrData, &extraData); err != nil {
			return nil, WrapInDvpERC721ClientError("failed to unmarshal extra data", err)
		}

		calldata := c.contract.PackMint(d.To, d.TokenID, d.ChainID, extraData)
		items = append(items, BatchInput{
			MsgID:   d.GetID(),
			Data:    calldata,
			Address: d.TokenAddress,
		})
	}

	results, err := c.executor.BatchExecute(ctx, items)
	if err != nil {
		return nil, WrapInDvpERC721ClientError("failed to batch mint ERC721", withstack.Wrap(err))
	}

	return results, nil
}
