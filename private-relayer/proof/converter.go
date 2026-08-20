package proof

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/raylsnetwork/rayls-sovereign-relayer/contracts/Proofs"
)

// ConvertEthHeader converts an Ethereum header to Proofs contract format.
// Post-London fields are optional in go-ethereum (nil on pre-fork chains) and
// default to zero when absent, which produces the correct RLP encoding for
// chains where those forks are not active.
func ConvertEthHeader(h *ethTypes.Header) Proofs.ProofsHeader {
	baseFeePerGas := h.BaseFee
	if baseFeePerGas == nil {
		baseFeePerGas = new(big.Int)
	}

	var withdrawalsRoot common.Hash
	if h.WithdrawalsHash != nil {
		withdrawalsRoot = *h.WithdrawalsHash
	}

	var blobGasUsed uint64
	if h.BlobGasUsed != nil {
		blobGasUsed = *h.BlobGasUsed
	}

	var excessBlobGas uint64
	if h.ExcessBlobGas != nil {
		excessBlobGas = *h.ExcessBlobGas
	}

	var parentBeaconBlockRoot common.Hash
	if h.ParentBeaconRoot != nil {
		parentBeaconBlockRoot = *h.ParentBeaconRoot
	}

	var requestsHash common.Hash
	if h.RequestsHash != nil {
		requestsHash = *h.RequestsHash
	}

	return Proofs.ProofsHeader{
		ParentHash:            h.ParentHash,
		UncleHash:             h.UncleHash,
		Coinbase:              h.Coinbase,
		Root:                  h.Root,
		TxHash:                h.TxHash,
		ReceiptHash:           h.ReceiptHash,
		Bloom:                 h.Bloom.Bytes(),
		Difficulty:            h.Difficulty,
		Number:                h.Number,
		GasLimit:              new(big.Int).SetUint64(h.GasLimit),
		GasUsed:               new(big.Int).SetUint64(h.GasUsed),
		Time:                  new(big.Int).SetUint64(h.Time),
		Extra:                 h.Extra,
		MixDigest:             h.MixDigest,
		Nonce:                 h.Nonce.Uint64(),
		BaseFeePerGas:         baseFeePerGas,
		WithdrawalsRoot:       withdrawalsRoot,
		BlobGasUsed:           blobGasUsed,
		ExcessBlobGas:         excessBlobGas,
		ParentBeaconBlockRoot: parentBeaconBlockRoot,
		RequestsHash:          requestsHash,
	}
}
