package testdata

import (
	"math/big"
	"testing"

	"github.com/raylsnetwork/rayls-sovereign-relayer/types"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewDvpDeposit_Defaults(t *testing.T) {
	d := NewDvpDeposit()
	require.NotNil(t, d)
	assert.Equal(t, "0xAlice", d.UserAddress)
	assert.Equal(t, big.NewInt(0xA1), d.Salt)
	assert.Equal(t, types.DvpERC1155, d.TokenType)
	assert.Equal(t, types.DvpDepositUnspent, d.Status)
	assert.Equal(t, big.NewInt(100), d.TokenAmount)
	assert.Equal(t, big.NewInt(1000), d.Commitment)
	assert.Equal(t, 1, d.TreeNumber)
}

func TestNewFungibleDeposit(t *testing.T) {
	t.Run("amount and type override defaults; commitment = amount*10", func(t *testing.T) {
		d := NewFungibleDeposit(150, types.DvpEnygma)
		assert.Equal(t, big.NewInt(150), d.TokenAmount)
		assert.Equal(t, types.DvpEnygma, d.TokenType)
		assert.Equal(t, "", d.TokenID, "fungible deposits clear TokenID")
		assert.Equal(t, big.NewInt(1500), d.Commitment, "commitment defaults to amount*10 for collision-free fungible inputs")
	})

	t.Run("trailing options can override the derived commitment", func(t *testing.T) {
		d := NewFungibleDeposit(100, types.DvpERC1155, WithDepositCommitment(big.NewInt(0xC0FFE)))
		assert.Equal(t, big.NewInt(0xC0FFE), d.Commitment)
	})
}

func TestNewNFTDeposit(t *testing.T) {
	d := NewNFTDeposit("tok-7")
	assert.Equal(t, "tok-7", d.TokenID)
	assert.Equal(t, types.DvpERC721, d.TokenType)
	assert.Equal(t, "0xNFT", d.TokenAddress)
	assert.Equal(t, big.NewInt(1), d.TokenAmount, "NFT deposit always has amount=1")
	assert.Equal(t, big.NewInt(0xCC), d.Commitment)
}

func TestNewDvpDeposit_Options(t *testing.T) {
	t.Run("WithDepositUserAddress", func(t *testing.T) {
		assert.Equal(t, "0xCarol", NewDvpDeposit(WithDepositUserAddress("0xCarol")).UserAddress)
	})
	t.Run("WithDepositSalt", func(t *testing.T) {
		assert.Equal(t, big.NewInt(99), NewDvpDeposit(WithDepositSalt(big.NewInt(99))).Salt)
	})
	t.Run("WithDepositCommitment", func(t *testing.T) {
		assert.Equal(t, big.NewInt(7), NewDvpDeposit(WithDepositCommitment(big.NewInt(7))).Commitment)
	})
	t.Run("WithDepositStatus", func(t *testing.T) {
		assert.Equal(t, types.DvpDepositSpent, NewDvpDeposit(WithDepositStatus(types.DvpDepositSpent)).Status)
	})
	t.Run("WithDepositNullifier", func(t *testing.T) {
		assert.Equal(t, big.NewInt(0xBEE5), NewDvpDeposit(WithDepositNullifier(big.NewInt(0xBEE5))).Nullifier)
	})
	t.Run("WithDepositTreeNumber", func(t *testing.T) {
		assert.Equal(t, 42, NewDvpDeposit(WithDepositTreeNumber(42)).TreeNumber)
	})
}
