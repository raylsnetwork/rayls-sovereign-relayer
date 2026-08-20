package testdata

import (
	"math/big"

	"github.com/raylsnetwork/rayls-privacy-relayer-api/dvp"
)

// ProofReceiptOption mutates a *dvp.ProofReceipt in place.
type ProofReceiptOption func(*dvp.ProofReceipt)

// NewProofReceipt returns a ProofReceipt with sane defaults: empty Proof,
// two output commitments, one nullifier (= 0xDEAD), revert commitment =
// 0xBEEF. Override individual fields with the With* options.
func NewProofReceipt(opts ...ProofReceiptOption) *dvp.ProofReceipt {
	r := &dvp.ProofReceipt{
		Proof:            &dvp.Proof{},
		Message:          big.NewInt(0),
		MerkleRoots:      []*big.Int{},
		Commitments:      []*big.Int{big.NewInt(1), big.NewInt(2)},
		Nullifiers:       []*big.Int{big.NewInt(0xDEAD)},
		TreeNumbers:      []*big.Int{},
		RevertCommitment: big.NewInt(0xBEEF),
	}
	for _, opt := range opts {
		opt(r)
	}
	return r
}

// WithRevertCommitment sets ProofReceipt.RevertCommitment.
func WithRevertCommitment(c *big.Int) ProofReceiptOption {
	return func(r *dvp.ProofReceipt) { r.RevertCommitment = c }
}

// WithCommitments sets ProofReceipt.Commitments (the two output commitments).
func WithCommitments(commitments ...*big.Int) ProofReceiptOption {
	return func(r *dvp.ProofReceipt) { r.Commitments = commitments }
}

// WithNullifiers sets ProofReceipt.Nullifiers.
func WithNullifiers(nullifiers ...*big.Int) ProofReceiptOption {
	return func(r *dvp.ProofReceipt) { r.Nullifiers = nullifiers }
}

// NewProofSignal returns the public-signal []*big.Int slice produced by a
// gnark-API mock. Layout depends on the proof type:
//
//	JS proof   (numIn=1, numOut=2): [message, merkleRoot, nullifier, treeNumber, commitOut0, commitOut1, revertCommitment]
//	Ownership  (numIn=1, numOut=1): [message, merkleRoot, nullifier, treeNumber, commitOut0,            revertCommitment]
//
// proofToReceipt in dvp/service/proof.go consumes this slice; tests use
// the helpers below to build matching layouts.
func NewProofSignal(opts ...ProofSignalOption) []*big.Int {
	cfg := &proofSignalConfig{
		message:          big.NewInt(0),
		merkleRoot:       big.NewInt(0xCAFE),
		nullifier:        big.NewInt(0xDEAD),
		treeNumber:       big.NewInt(1),
		commitments:      []*big.Int{big.NewInt(0xC0FFE), big.NewInt(0xCAFEE)},
		revertCommitment: big.NewInt(0xBEEF),
	}
	for _, opt := range opts {
		opt(cfg)
	}
	signal := []*big.Int{cfg.message, cfg.merkleRoot, cfg.nullifier, cfg.treeNumber}
	signal = append(signal, cfg.commitments...)
	signal = append(signal, cfg.revertCommitment)
	return signal
}

type proofSignalConfig struct {
	message          *big.Int
	merkleRoot       *big.Int
	nullifier        *big.Int
	treeNumber       *big.Int
	commitments      []*big.Int
	revertCommitment *big.Int
}

// ProofSignalOption mutates an internal proof-signal config used by
// NewProofSignal.
type ProofSignalOption func(*proofSignalConfig)

// WithSignalNullifier sets the nullifier field of the public signal.
func WithSignalNullifier(n *big.Int) ProofSignalOption {
	return func(c *proofSignalConfig) { c.nullifier = n }
}

// WithSignalRevertCommitment sets the revert-commitment field of the public signal.
func WithSignalRevertCommitment(c *big.Int) ProofSignalOption {
	return func(cfg *proofSignalConfig) { cfg.revertCommitment = c }
}

// WithJSLayout configures the signal for a join-split proof: two output
// commitments (destination + change). This is the default layout, so
// callers usually omit this option.
func WithJSLayout() ProofSignalOption {
	return func(c *proofSignalConfig) {
		c.commitments = []*big.Int{big.NewInt(0xC0FFE), big.NewInt(0xCAFEE)}
	}
}

// WithOwnershipLayout configures the signal for an ERC721 ownership proof:
// a single output commitment (no change).
func WithOwnershipLayout() ProofSignalOption {
	return func(c *proofSignalConfig) {
		c.commitments = []*big.Int{big.NewInt(0xCC)}
	}
}

// WithSignalCommitments overrides the output commitments — useful when a
// test asserts that proofToReceipt slices a specific number of commitments
// from the public signal.
func WithSignalCommitments(commitments ...*big.Int) ProofSignalOption {
	return func(c *proofSignalConfig) { c.commitments = commitments }
}
