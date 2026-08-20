package proof_test

import (
	"math"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/private-relayer/proof"
	"github.com/stretchr/testify/assert"
)

func TestConvertEthHeader(t *testing.T) {
	t.Run("preserves all field values", func(t *testing.T) {
		parentHash := common.HexToHash("0xaaaa")
		uncleHash := common.HexToHash("0xbbbb")
		coinbase := common.HexToAddress("0xcccc")
		root := common.HexToHash("0xdddd")
		txHash := common.HexToHash("0xeeee")
		receiptHash := common.HexToHash("0xffff")
		difficulty := big.NewInt(12345)
		number := big.NewInt(100)
		mixDigest := common.HexToHash("0x1111")
		nonce := ethTypes.EncodeNonce(42)

		header := &ethTypes.Header{
			ParentHash:  parentHash,
			UncleHash:   uncleHash,
			Coinbase:    coinbase,
			Root:        root,
			TxHash:      txHash,
			ReceiptHash: receiptHash,
			Difficulty:  difficulty,
			Number:      number,
			GasLimit:    30_000_000,
			GasUsed:     21_000,
			Time:        1_700_000_000,
			Extra:       []byte("extra-data"),
			MixDigest:   mixDigest,
			Nonce:       nonce,
		}

		result := proof.ConvertEthHeader(header)

		assert.Equal(t, [32]byte(parentHash), result.ParentHash)
		assert.Equal(t, [32]byte(uncleHash), result.UncleHash)
		assert.Equal(t, coinbase, result.Coinbase)
		assert.Equal(t, [32]byte(root), result.Root)
		assert.Equal(t, [32]byte(txHash), result.TxHash)
		assert.Equal(t, [32]byte(receiptHash), result.ReceiptHash)
		assert.Equal(t, difficulty, result.Difficulty)
		assert.Equal(t, number, result.Number)
		assert.Equal(t, [32]byte(mixDigest), result.MixDigest)
		assert.Equal(t, []byte("extra-data"), result.Extra)
		assert.Equal(t, header.Bloom.Bytes(), result.Bloom)
		assert.Equal(t, big.NewInt(30_000_000), result.GasLimit)
		assert.Equal(t, big.NewInt(21_000), result.GasUsed)
		assert.Equal(t, big.NewInt(1_700_000_000), result.Time)
		assert.Equal(t, uint64(42), result.Nonce)
	})

	t.Run("preserves max uint64 values without overflow", func(t *testing.T) {
		maxUint64 := new(big.Int).SetUint64(math.MaxUint64)

		header := &ethTypes.Header{
			Difficulty: big.NewInt(0),
			Number:     big.NewInt(0),
			GasLimit:   math.MaxUint64,
			GasUsed:    math.MaxUint64,
			Time:       math.MaxUint64,
			Nonce:      ethTypes.EncodeNonce(math.MaxUint64),
		}

		result := proof.ConvertEthHeader(header)

		assert.Equal(t, maxUint64, result.GasLimit)
		assert.Equal(t, maxUint64, result.GasUsed)
		assert.Equal(t, maxUint64, result.Time)
		assert.Equal(t, uint64(math.MaxUint64), result.Nonce)
	})

	t.Run("preserves zero values without underflow", func(t *testing.T) {
		header := &ethTypes.Header{
			Difficulty: big.NewInt(0),
			Number:     big.NewInt(0),
			GasLimit:   0,
			GasUsed:    0,
			Time:       0,
			Nonce:      ethTypes.EncodeNonce(0),
			Extra:      nil,
		}

		result := proof.ConvertEthHeader(header)

		assert.Equal(t, big.NewInt(0), result.GasLimit)
		assert.Equal(t, big.NewInt(0), result.GasUsed)
		assert.Equal(t, big.NewInt(0), result.Time)
		assert.Equal(t, big.NewInt(0), result.Difficulty)
		assert.Equal(t, big.NewInt(0), result.Number)
		assert.Equal(t, uint64(0), result.Nonce)
		assert.Nil(t, result.Extra)
	})

	t.Run("preserves large big.Int values beyond uint64 range", func(t *testing.T) {
		// 2^128 - far exceeds uint64 range
		largeDifficulty := new(big.Int).Exp(big.NewInt(2), big.NewInt(128), nil)
		largeNumber := new(big.Int).Exp(big.NewInt(2), big.NewInt(128), nil)

		header := &ethTypes.Header{
			Difficulty: largeDifficulty,
			Number:     largeNumber,
			GasLimit:   0,
			GasUsed:    0,
			Time:       0,
		}

		result := proof.ConvertEthHeader(header)

		assert.Equal(t, largeDifficulty, result.Difficulty)
		assert.Equal(t, largeNumber, result.Number)
	})

	t.Run("preserves full bloom filter", func(t *testing.T) {
		var bloom ethTypes.Bloom
		for i := range bloom {
			bloom[i] = 0xFF
		}

		header := &ethTypes.Header{
			Bloom:      bloom,
			Difficulty: big.NewInt(0),
			Number:     big.NewInt(0),
		}

		result := proof.ConvertEthHeader(header)

		assert.Len(t, result.Bloom, 256)
		for i, b := range result.Bloom {
			assert.Equal(t, byte(0xFF), b, "bloom byte %d should be 0xFF", i)
		}
	})
}
