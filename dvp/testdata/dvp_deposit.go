package testdata

import (
	"math/big"

	"github.com/raylsnetwork/rayls-sovereign-relayer/types"
)

// DepositOption mutates a *types.DvpDeposit in place.
type DepositOption func(*types.DvpDeposit)

// NewDvpDeposit returns a deposit populated with sane defaults
// (Unspent ERC1155 fixture, 100 tokens, salt = 0xA1, commitment = 1000,
// nullifier = 0). Override fields with the WithDeposit* options.
func NewDvpDeposit(opts ...DepositOption) *types.DvpDeposit {
	d := &types.DvpDeposit{
		UserAddress:  "0xAlice",
		Salt:         big.NewInt(0xA1),
		TokenAmount:  big.NewInt(100),
		TokenAddress: "0xToken",
		TokenID:      "1",
		TokenType:    types.DvpERC1155,
		TreeNumber:   1,
		Commitment:   big.NewInt(1000),
		Nullifier:    big.NewInt(0),
		Status:       types.DvpDepositUnspent,
	}
	for _, opt := range opts {
		opt(d)
	}
	return d
}

// NewFungibleDeposit returns a deposit shaped for join-split proof inputs:
// matches NewDvpDeposit defaults but with caller-controlled amount, type,
// and a commitment derived as amount*10 (matches the historical helper).
func NewFungibleDeposit(amount int64, tokenType types.DvpTokenType, opts ...DepositOption) *types.DvpDeposit {
	merged := append([]DepositOption{
		WithDepositTokenAmount(big.NewInt(amount)),
		WithDepositTokenType(tokenType),
		WithDepositTokenID(""),
		WithDepositCommitment(big.NewInt(amount * 10)),
	}, opts...)
	return NewDvpDeposit(merged...)
}

// NewNFTDeposit returns a deposit shaped for ERC721 ownership-proof input:
// amount = 1, tokenID = nftID, commitment = 0xCC.
func NewNFTDeposit(nftID string, opts ...DepositOption) *types.DvpDeposit {
	merged := append([]DepositOption{
		WithDepositTokenAmount(big.NewInt(1)),
		WithDepositTokenType(types.DvpERC721),
		WithDepositTokenID(nftID),
		WithDepositTokenAddress("0xNFT"),
		WithDepositCommitment(big.NewInt(0xCC)),
		WithDepositSalt(big.NewInt(0xA2)),
	}, opts...)
	return NewDvpDeposit(merged...)
}

// WithDepositUserAddress sets DvpDeposit.UserAddress.
func WithDepositUserAddress(addr string) DepositOption {
	return func(d *types.DvpDeposit) { d.UserAddress = addr }
}

// WithDepositSalt sets DvpDeposit.Salt.
func WithDepositSalt(salt *big.Int) DepositOption {
	return func(d *types.DvpDeposit) { d.Salt = salt }
}

// WithDepositCommitment sets DvpDeposit.Commitment.
func WithDepositCommitment(commitment *big.Int) DepositOption {
	return func(d *types.DvpDeposit) { d.Commitment = commitment }
}

// WithDepositTokenAmount sets DvpDeposit.TokenAmount.
func WithDepositTokenAmount(amount *big.Int) DepositOption {
	return func(d *types.DvpDeposit) { d.TokenAmount = amount }
}

// WithDepositTokenAddress sets DvpDeposit.TokenAddress.
func WithDepositTokenAddress(addr string) DepositOption {
	return func(d *types.DvpDeposit) { d.TokenAddress = addr }
}

// WithDepositTokenID sets DvpDeposit.TokenID.
func WithDepositTokenID(id string) DepositOption {
	return func(d *types.DvpDeposit) { d.TokenID = id }
}

// WithDepositTokenType sets DvpDeposit.TokenType.
func WithDepositTokenType(t types.DvpTokenType) DepositOption {
	return func(d *types.DvpDeposit) { d.TokenType = t }
}

// WithDepositStatus sets DvpDeposit.Status.
func WithDepositStatus(s types.DvpDepositStatus) DepositOption {
	return func(d *types.DvpDeposit) { d.Status = s }
}

// WithDepositNullifier sets DvpDeposit.Nullifier.
func WithDepositNullifier(n *big.Int) DepositOption {
	return func(d *types.DvpDeposit) { d.Nullifier = n }
}

// WithDepositTreeNumber sets DvpDeposit.TreeNumber.
func WithDepositTreeNumber(n int) DepositOption {
	return func(d *types.DvpDeposit) { d.TreeNumber = n }
}
