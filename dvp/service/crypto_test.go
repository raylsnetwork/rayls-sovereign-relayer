package service_test

import (
	"math/big"
	"testing"

	"github.com/raylsnetwork/rayls-sovereign-relayer/cryptography"
	"github.com/raylsnetwork/rayls-sovereign-relayer/dvp/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCommitmentCalculator_CalculateNFTCommitment(t *testing.T) {
	t.Run("successfully calculates NFT commitment", func(t *testing.T) {
		calc := service.NewCommitmentCalculator()

		spendPK := big.NewInt(999)
		salt := big.NewInt(12345)

		result, _ := calc.CalculateNFTCommitment(spendPK, salt, "42", "0x1234567890123456789012345678901234567890")

		require.NotNil(t, result)
		assert.NotZero(t, result.Sign())
		assert.True(t, result.Cmp(cryptography.JubJubPrimeGroup) < 0)
	})

	t.Run("deterministic output", func(t *testing.T) {
		calc := service.NewCommitmentCalculator()

		spendPK := big.NewInt(999)
		salt := big.NewInt(12345)

		result1, _ := calc.CalculateNFTCommitment(spendPK, salt, "42", "0x1234567890123456789012345678901234567890")
		result2, _ := calc.CalculateNFTCommitment(spendPK, salt, "42", "0x1234567890123456789012345678901234567890")

		assert.Equal(t, result1, result2)
	})

	t.Run("different spendPK produces different commitment", func(t *testing.T) {
		calc := service.NewCommitmentCalculator()

		salt := big.NewInt(12345)
		nftID := "42"
		nftAddress := "0x1234567890123456789012345678901234567890"

		result1, _ := calc.CalculateNFTCommitment(big.NewInt(111), salt, nftID, nftAddress)
		result2, _ := calc.CalculateNFTCommitment(big.NewInt(222), salt, nftID, nftAddress)

		assert.NotEqual(t, result1, result2)
	})

	t.Run("different salt produces different commitment", func(t *testing.T) {
		calc := service.NewCommitmentCalculator()

		spendPK := big.NewInt(999)
		nftID := "42"
		nftAddress := "0x1234567890123456789012345678901234567890"

		result1, _ := calc.CalculateNFTCommitment(spendPK, big.NewInt(111), nftID, nftAddress)
		result2, _ := calc.CalculateNFTCommitment(spendPK, big.NewInt(222), nftID, nftAddress)

		assert.NotEqual(t, result1, result2)
	})

	t.Run("different nftID produces different commitment", func(t *testing.T) {
		calc := service.NewCommitmentCalculator()

		spendPK := big.NewInt(999)
		salt := big.NewInt(12345)
		nftAddress := "0x1234567890123456789012345678901234567890"

		result1, _ := calc.CalculateNFTCommitment(spendPK, salt, "1", nftAddress)
		result2, _ := calc.CalculateNFTCommitment(spendPK, salt, "2", nftAddress)

		assert.NotEqual(t, result1, result2)
	})

	t.Run("different nftAddress produces different commitment", func(t *testing.T) {
		calc := service.NewCommitmentCalculator()

		spendPK := big.NewInt(999)
		salt := big.NewInt(12345)
		nftID := "42"

		result1, _ := calc.CalculateNFTCommitment(spendPK, salt, nftID, "0x1111111111111111111111111111111111111111")
		result2, _ := calc.CalculateNFTCommitment(spendPK, salt, nftID, "0x2222222222222222222222222222222222222222")

		assert.NotEqual(t, result1, result2)
	})

	t.Run("result within field", func(t *testing.T) {
		calc := service.NewCommitmentCalculator()

		largeVal := new(big.Int).Sub(cryptography.JubJubPrimeGroup, big.NewInt(1))
		result, _ := calc.CalculateNFTCommitment(largeVal, largeVal, "999999999999", "0xffffffffffffffffffffffffffffffffffffffff")

		require.NotNil(t, result)
		assert.True(t, result.Cmp(cryptography.JubJubPrimeGroup) < 0)
		assert.True(t, result.Sign() >= 0)
	})
}

func TestCommitmentCalculator_CalculatePaymentCommitment(t *testing.T) {
	t.Run("successfully calculates payment commitment", func(t *testing.T) {
		calc := service.NewCommitmentCalculator()

		spendPK := big.NewInt(999)
		salt := big.NewInt(12345)
		amount := big.NewInt(1000)

		result, _ := calc.CalculatePaymentCommitment(spendPK, salt, amount, "0x1234567890123456789012345678901234567890")

		require.NotNil(t, result)
		assert.NotZero(t, result.Sign())
		assert.True(t, result.Cmp(cryptography.JubJubPrimeGroup) < 0)
	})

	t.Run("deterministic output", func(t *testing.T) {
		calc := service.NewCommitmentCalculator()

		spendPK := big.NewInt(999)
		salt := big.NewInt(12345)
		amount := big.NewInt(1000)
		tokenAddr := "0x1234567890123456789012345678901234567890"

		result1, _ := calc.CalculatePaymentCommitment(spendPK, salt, amount, tokenAddr)
		result2, _ := calc.CalculatePaymentCommitment(spendPK, salt, amount, tokenAddr)

		assert.Equal(t, result1, result2)
	})

	t.Run("different amounts produce different commitments", func(t *testing.T) {
		calc := service.NewCommitmentCalculator()

		spendPK := big.NewInt(999)
		salt := big.NewInt(12345)
		tokenAddr := "0x1234567890123456789012345678901234567890"

		result1, _ := calc.CalculatePaymentCommitment(spendPK, salt, big.NewInt(100), tokenAddr)
		result2, _ := calc.CalculatePaymentCommitment(spendPK, salt, big.NewInt(200), tokenAddr)

		assert.NotEqual(t, result1, result2)
	})

	t.Run("different salts produce different commitments", func(t *testing.T) {
		calc := service.NewCommitmentCalculator()

		spendPK := big.NewInt(999)
		amount := big.NewInt(1000)
		tokenAddr := "0x1234567890123456789012345678901234567890"

		result1, _ := calc.CalculatePaymentCommitment(spendPK, big.NewInt(1), amount, tokenAddr)
		result2, _ := calc.CalculatePaymentCommitment(spendPK, big.NewInt(2), amount, tokenAddr)

		assert.NotEqual(t, result1, result2)
	})

	t.Run("different token addresses produce different commitments", func(t *testing.T) {
		calc := service.NewCommitmentCalculator()

		spendPK := big.NewInt(999)
		salt := big.NewInt(12345)
		amount := big.NewInt(1000)

		result1, _ := calc.CalculatePaymentCommitment(spendPK, salt, amount, "0x1111111111111111111111111111111111111111")
		result2, _ := calc.CalculatePaymentCommitment(spendPK, salt, amount, "0x2222222222222222222222222222222222222222")

		assert.NotEqual(t, result1, result2)
	})

	t.Run("handles zero payment amount", func(t *testing.T) {
		calc := service.NewCommitmentCalculator()

		result, _ := calc.CalculatePaymentCommitment(big.NewInt(999), big.NewInt(1), big.NewInt(0), "0x1234567890123456789012345678901234567890")
		require.NotNil(t, result)
	})

	t.Run("handles large values", func(t *testing.T) {
		calc := service.NewCommitmentCalculator()

		largeVal := new(big.Int).Sub(cryptography.JubJubPrimeGroup, big.NewInt(1))
		result, _ := calc.CalculatePaymentCommitment(largeVal, largeVal, largeVal, "0xffffffffffffffffffffffffffffffffffffffff")

		require.NotNil(t, result)
		assert.True(t, result.Cmp(cryptography.JubJubPrimeGroup) < 0)
	})
}

func TestCommitmentCalculator_CalculateERC1155Commitment(t *testing.T) {
	t.Run("successfully calculates ERC1155 commitment", func(t *testing.T) {
		calc := service.NewCommitmentCalculator()

		spendPK := big.NewInt(999)
		salt := big.NewInt(12345)
		amount := big.NewInt(100)

		result, _ := calc.CalculateERC1155Commitment(spendPK, salt, "0x1234567890123456789012345678901234567890", "42", amount)

		require.NotNil(t, result)
		assert.NotZero(t, result.Sign())
		assert.True(t, result.Cmp(cryptography.JubJubPrimeGroup) < 0)
	})

	t.Run("deterministic output", func(t *testing.T) {
		calc := service.NewCommitmentCalculator()

		spendPK := big.NewInt(999)
		salt := big.NewInt(12345)
		amount := big.NewInt(100)
		tokenAddr := "0x1234567890123456789012345678901234567890"

		result1, _ := calc.CalculateERC1155Commitment(spendPK, salt, tokenAddr, "42", amount)
		result2, _ := calc.CalculateERC1155Commitment(spendPK, salt, tokenAddr, "42", amount)

		assert.Equal(t, result1, result2)
	})

	t.Run("different token IDs produce different commitments", func(t *testing.T) {
		calc := service.NewCommitmentCalculator()

		spendPK := big.NewInt(999)
		salt := big.NewInt(12345)
		tokenAddr := "0x1234567890123456789012345678901234567890"
		amount := big.NewInt(100)

		result1, _ := calc.CalculateERC1155Commitment(spendPK, salt, tokenAddr, "1", amount)
		result2, _ := calc.CalculateERC1155Commitment(spendPK, salt, tokenAddr, "2", amount)

		assert.NotEqual(t, result1, result2)
	})

	t.Run("different amounts produce different commitments", func(t *testing.T) {
		calc := service.NewCommitmentCalculator()

		spendPK := big.NewInt(999)
		salt := big.NewInt(12345)
		tokenAddr := "0x1234567890123456789012345678901234567890"

		result1, _ := calc.CalculateERC1155Commitment(spendPK, salt, tokenAddr, "42", big.NewInt(100))
		result2, _ := calc.CalculateERC1155Commitment(spendPK, salt, tokenAddr, "42", big.NewInt(200))

		assert.NotEqual(t, result1, result2)
	})

	t.Run("result within field", func(t *testing.T) {
		calc := service.NewCommitmentCalculator()

		largeVal := new(big.Int).Sub(cryptography.JubJubPrimeGroup, big.NewInt(1))
		result, _ := calc.CalculateERC1155Commitment(largeVal, largeVal, "0xffffffffffffffffffffffffffffffffffffffff", "999999999999", largeVal)

		require.NotNil(t, result)
		assert.True(t, result.Cmp(cryptography.JubJubPrimeGroup) < 0)
	})
}

func TestCommitmentCalculator_CalculateNullifier(t *testing.T) {
	t.Run("successfully calculates nullifier", func(t *testing.T) {
		calc := service.NewCommitmentCalculator()

		spendSK := big.NewInt(42)
		leafIndex := big.NewInt(7)

		result, err := calc.CalculateNullifier(spendSK, leafIndex)

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.NotZero(t, result.Sign())
		assert.True(t, result.Cmp(cryptography.JubJubPrimeGroup) < 0)
	})

	t.Run("deterministic output", func(t *testing.T) {
		calc := service.NewCommitmentCalculator()

		spendSK := big.NewInt(42)
		leafIndex := big.NewInt(7)

		result1, err := calc.CalculateNullifier(spendSK, leafIndex)
		require.NoError(t, err)
		result2, err := calc.CalculateNullifier(spendSK, leafIndex)
		require.NoError(t, err)

		assert.Equal(t, result1, result2)
	})

	t.Run("different spendSK produces different nullifier", func(t *testing.T) {
		calc := service.NewCommitmentCalculator()

		leafIndex := big.NewInt(7)

		result1, err := calc.CalculateNullifier(big.NewInt(1), leafIndex)
		require.NoError(t, err)
		result2, err := calc.CalculateNullifier(big.NewInt(2), leafIndex)
		require.NoError(t, err)

		assert.NotEqual(t, result1, result2)
	})

	t.Run("different leafIndex produces different nullifier", func(t *testing.T) {
		calc := service.NewCommitmentCalculator()

		spendSK := big.NewInt(42)

		result1, err := calc.CalculateNullifier(spendSK, big.NewInt(0))
		require.NoError(t, err)
		result2, err := calc.CalculateNullifier(spendSK, big.NewInt(1))
		require.NoError(t, err)

		assert.NotEqual(t, result1, result2)
	})

	t.Run("handles large values", func(t *testing.T) {
		calc := service.NewCommitmentCalculator()

		largeVal := new(big.Int).Sub(cryptography.JubJubPrimeGroup, big.NewInt(1))
		result, err := calc.CalculateNullifier(largeVal, largeVal)

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.True(t, result.Cmp(cryptography.JubJubPrimeGroup) < 0)
	})
}

func TestCommitmentCalculator_GetNFTUniqueID(t *testing.T) {
	t.Run("successfully calculates NFT unique ID", func(t *testing.T) {
		calc := service.NewCommitmentCalculator()

		result, err := calc.GetNFTUniqueID("0x1234567890123456789012345678901234567890", "123")

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.NotZero(t, result.Sign())
		assert.True(t, result.Cmp(cryptography.JubJubPrimeGroup) < 0)
	})

	t.Run("deterministic output", func(t *testing.T) {
		calc := service.NewCommitmentCalculator()

		result1, err1 := calc.GetNFTUniqueID("0xabcdefabcdefabcdefabcdefabcdefabcdefabcd", "999")
		require.NoError(t, err1)

		result2, err2 := calc.GetNFTUniqueID("0xabcdefabcdefabcdefabcdefabcdefabcdefabcd", "999")
		require.NoError(t, err2)

		assert.Equal(t, result1, result2)
	})

	t.Run("different addresses produce different unique IDs", func(t *testing.T) {
		calc := service.NewCommitmentCalculator()

		result1, err1 := calc.GetNFTUniqueID("0x1111111111111111111111111111111111111111", "123")
		require.NoError(t, err1)

		result2, err2 := calc.GetNFTUniqueID("0x2222222222222222222222222222222222222222", "123")
		require.NoError(t, err2)

		assert.NotEqual(t, result1, result2)
	})

	t.Run("different token IDs produce different unique IDs", func(t *testing.T) {
		calc := service.NewCommitmentCalculator()

		result1, err1 := calc.GetNFTUniqueID("0x1234567890123456789012345678901234567890", "111")
		require.NoError(t, err1)

		result2, err2 := calc.GetNFTUniqueID("0x1234567890123456789012345678901234567890", "222")
		require.NoError(t, err2)

		assert.NotEqual(t, result1, result2)
	})
}

// Benchmark tests
func BenchmarkCommitmentCalculator_CalculateNFTCommitment(b *testing.B) {
	calc := service.NewCommitmentCalculator()
	spendPK := big.NewInt(999)
	salt := big.NewInt(12345)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		calc.CalculateNFTCommitment(spendPK, salt, "42", "0x1234567890123456789012345678901234567890")
	}
}

func BenchmarkCommitmentCalculator_CalculatePaymentCommitment(b *testing.B) {
	calc := service.NewCommitmentCalculator()
	spendPK := big.NewInt(999)
	salt := big.NewInt(12345)
	amount := big.NewInt(1000)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		calc.CalculatePaymentCommitment(spendPK, salt, amount, "0x1234567890123456789012345678901234567890")
	}
}

func BenchmarkCommitmentCalculator_CalculateERC1155Commitment(b *testing.B) {
	calc := service.NewCommitmentCalculator()
	spendPK := big.NewInt(999)
	salt := big.NewInt(12345)
	amount := big.NewInt(100)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		calc.CalculateERC1155Commitment(spendPK, salt, "0x1234567890123456789012345678901234567890", "42", amount)
	}
}

func BenchmarkCommitmentCalculator_CalculateNullifier(b *testing.B) {
	calc := service.NewCommitmentCalculator()
	spendSK := big.NewInt(42)
	leafIndex := big.NewInt(7)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = calc.CalculateNullifier(spendSK, leafIndex)
	}
}

func BenchmarkCommitmentCalculator_GetNFTUniqueID(b *testing.B) {
	calc := service.NewCommitmentCalculator()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = calc.GetNFTUniqueID("0x1234567890123456789012345678901234567890", "123")
	}
}
