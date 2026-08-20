// FIXME(integration): does not compile under `-tags integration` — sets `PublicKey`, which
// types.DvpDeposit no longer has (now `Salt`). Pre-existing drift; left unfixed by request.

package testdata

import (
	"math/big"
	"time"

	"github.com/raylsnetwork/rayls-sovereign-relayer/private-relayer/repository"
	"github.com/raylsnetwork/rayls-sovereign-relayer/types"
)

var depositCreatedAt = time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

var (
	DvpDeposit1 = types.DvpDeposit{
		UserAddress:  "0xUser1",
		PublicKey:    big.NewInt(111),
		TokenAmount:  big.NewInt(1000),
		TokenAddress: "0xTokenA",
		TokenType:    types.DvpERC20,
		TokenID:      "",
		TreeNumber:   1,
		Commitment:   big.NewInt(9001),
		Nullifier:    big.NewInt(0),
		Status:       types.DvpDepositPending,
		CreatedAt:    depositCreatedAt,
	}
	DvpDeposit2 = types.DvpDeposit{
		UserAddress:  "0xUser1",
		PublicKey:    big.NewInt(222),
		TokenAmount:  big.NewInt(500),
		TokenAddress: "0xTokenA",
		TokenType:    types.DvpERC20,
		TokenID:      "",
		TreeNumber:   1,
		Commitment:   big.NewInt(9002),
		Nullifier:    big.NewInt(5555),
		Status:       types.DvpDepositUnspent,
		CreatedAt:    depositCreatedAt,
	}
	DvpDeposit3NFT = types.DvpDeposit{
		UserAddress:  "0xUser2",
		PublicKey:    big.NewInt(333),
		TokenAmount:  big.NewInt(1),
		TokenAddress: "0xTokenB",
		TokenType:    types.DvpERC721,
		TokenID:      "nft-token-id-1",
		TreeNumber:   0,
		Commitment:   big.NewInt(9003),
		Nullifier:    big.NewInt(0),
		Status:       types.DvpDepositUnspent,
		CreatedAt:    depositCreatedAt,
	}
)

var (
	ModelDvpDeposit1 = repository.DvpDeposit{
		UserAddress:  DvpDeposit1.UserAddress,
		PublicKey:    DvpDeposit1.PublicKey.String(),
		TokenAmount:  DvpDeposit1.TokenAmount.String(),
		TokenAddress: DvpDeposit1.TokenAddress,
		TokenType:    DvpDeposit1.TokenType,
		TokenID:      DvpDeposit1.TokenID,
		TreeNumber:   DvpDeposit1.TreeNumber,
		Commitment:   DvpDeposit1.Commitment.String(),
		Nullifier:    "0",
		Status:       DvpDeposit1.Status,
		CreatedAt:    DvpDeposit1.CreatedAt,
	}
	ModelDvpDeposit2 = repository.DvpDeposit{
		UserAddress:  DvpDeposit2.UserAddress,
		PublicKey:    DvpDeposit2.PublicKey.String(),
		TokenAmount:  DvpDeposit2.TokenAmount.String(),
		TokenAddress: DvpDeposit2.TokenAddress,
		TokenType:    DvpDeposit2.TokenType,
		TokenID:      DvpDeposit2.TokenID,
		TreeNumber:   DvpDeposit2.TreeNumber,
		Commitment:   DvpDeposit2.Commitment.String(),
		Nullifier:    DvpDeposit2.Nullifier.String(),
		Status:       DvpDeposit2.Status,
		CreatedAt:    DvpDeposit2.CreatedAt,
	}
	ModelDvpDeposit3NFT = repository.DvpDeposit{
		UserAddress:  DvpDeposit3NFT.UserAddress,
		PublicKey:    DvpDeposit3NFT.PublicKey.String(),
		TokenAmount:  DvpDeposit3NFT.TokenAmount.String(),
		TokenAddress: DvpDeposit3NFT.TokenAddress,
		TokenType:    DvpDeposit3NFT.TokenType,
		TokenID:      DvpDeposit3NFT.TokenID,
		TreeNumber:   DvpDeposit3NFT.TreeNumber,
		Commitment:   DvpDeposit3NFT.Commitment.String(),
		Nullifier:    "0",
		Status:       DvpDeposit3NFT.Status,
		CreatedAt:    DvpDeposit3NFT.CreatedAt,
	}
)
