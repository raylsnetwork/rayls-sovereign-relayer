package testdata

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewProofReceipt_Defaults(t *testing.T) {
	r := NewProofReceipt()
	require.NotNil(t, r)
	require.NotNil(t, r.Proof)
	assert.Equal(t, big.NewInt(0xBEEF), r.RevertCommitment)
	assert.Len(t, r.Commitments, 2, "default proof receipt carries 2 output commitments (destination + change)")
	assert.Len(t, r.Nullifiers, 1)
	assert.Equal(t, big.NewInt(0xDEAD), r.Nullifiers[0])
}

func TestNewProofReceipt_Options(t *testing.T) {
	t.Run("WithRevertCommitment", func(t *testing.T) {
		r := NewProofReceipt(WithRevertCommitment(big.NewInt(0x1234)))
		assert.Equal(t, big.NewInt(0x1234), r.RevertCommitment)
	})
	t.Run("WithCommitments", func(t *testing.T) {
		r := NewProofReceipt(WithCommitments(big.NewInt(1), big.NewInt(2), big.NewInt(3)))
		assert.Equal(t, []*big.Int{big.NewInt(1), big.NewInt(2), big.NewInt(3)}, r.Commitments)
	})
	t.Run("WithNullifiers", func(t *testing.T) {
		r := NewProofReceipt(WithNullifiers(big.NewInt(7), big.NewInt(8)))
		assert.Equal(t, []*big.Int{big.NewInt(7), big.NewInt(8)}, r.Nullifiers)
	})
}

func TestNewProofSignal_DefaultJSLayout(t *testing.T) {
	signal := NewProofSignal()
	require.Len(t, signal, 7,
		"JS layout: [message, merkleRoot, nullifier, treeNumber, commitOut0, commitOut1, revertCommitment]")
	assert.Equal(t, big.NewInt(0xCAFE), signal[1], "merkleRoot at index 1")
	assert.Equal(t, big.NewInt(0xDEAD), signal[2], "nullifier at index 2")
	assert.Equal(t, big.NewInt(1), signal[3], "treeNumber at index 3")
	assert.Equal(t, big.NewInt(0xC0FFE), signal[4], "destination commitment at index 4")
	assert.Equal(t, big.NewInt(0xCAFEE), signal[5], "change commitment at index 5")
	assert.Equal(t, big.NewInt(0xBEEF), signal[6], "revert commitment is the last signal")
}

func TestNewProofSignal_OwnershipLayout(t *testing.T) {
	signal := NewProofSignal(WithOwnershipLayout())
	require.Len(t, signal, 6,
		"Ownership layout: [message, merkleRoot, nullifier, treeNumber, commitOut0, revertCommitment]")
	assert.Equal(t, big.NewInt(0xCC), signal[4], "single output commitment at index 4")
	assert.Equal(t, big.NewInt(0xBEEF), signal[5], "revert commitment is the last signal")
}

func TestNewProofSignal_CustomCommitments(t *testing.T) {
	signal := NewProofSignal(WithSignalCommitments(big.NewInt(11), big.NewInt(22), big.NewInt(33)))
	require.Len(t, signal, 8, "4 prefix + 3 commitments + 1 revert = 8")
	assert.Equal(t, big.NewInt(11), signal[4])
	assert.Equal(t, big.NewInt(33), signal[6])
	assert.Equal(t, big.NewInt(0xBEEF), signal[7])
}

func TestNewProofSignal_OverrideNullifierAndRevert(t *testing.T) {
	signal := NewProofSignal(
		WithSignalNullifier(big.NewInt(0xFEED)),
		WithSignalRevertCommitment(big.NewInt(0xBABE)),
	)
	assert.Equal(t, big.NewInt(0xFEED), signal[2])
	assert.Equal(t, big.NewInt(0xBABE), signal[len(signal)-1])
}
