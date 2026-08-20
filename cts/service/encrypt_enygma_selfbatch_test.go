package service

import (
	"context"
	"errors"
	"math/big"
	"testing"

	"github.com/raylsnetwork/rayls-privacy-relayer-api/cts/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// selfBatchPassthroughEncryptor is a no-op at-rest encryptor so secrets round-trip unchanged in this
// unit test.
type selfBatchPassthroughEncryptor struct{}

func (selfBatchPassthroughEncryptor) Encrypt(data []byte) ([]byte, error) { return data, nil }
func (selfBatchPassthroughEncryptor) Decrypt(data []byte) ([]byte, error) { return data, nil }

// selfBatchKeysStub returns DISTINCT secrets per accessor so we can reproduce the seal/decrypt
// secret-source divergence on the sender's own-chain ("self"/change) batch:
//   - GetSharedSecret            → the static own-chain PAIRWISE row (what plain GCMEncrypt used)
//   - GetEnygmaSharedSelfSecret  → the rotating per-resource SELF-secret (what the on-chain tag,
//     the ZK proof, and GCMDecryptWithEnygmaMessageTag use)
//   - GetEnygmaSharedSecrets     → the anonymity-set view: the own-chain entry IS the self-secret
//     (mirrors KeysService.GetEnygmaSharedSecrets, which excludes the own pairwise row and appends
//     the self-secret for myChainID)
type selfBatchKeysStub struct {
	myChainID    *big.Int
	pairwiseSelf []byte
	selfSecret   []byte
}

func (s *selfBatchKeysStub) GetSharedSecret(_ context.Context, _ *big.Int, _ uint64) (domain.SharedSecret, error) {
	return domain.SharedSecret{ChainId: s.myChainID, Secret: s.pairwiseSelf, InitialBlock: 1}, nil
}

func (s *selfBatchKeysStub) GetAllSharedSecrets(_ context.Context, _ uint64) ([]domain.SharedSecret, error) {
	return []domain.SharedSecret{{ChainId: s.myChainID, Secret: s.pairwiseSelf}}, nil
}

func (s *selfBatchKeysStub) GetEnygmaSharedSecrets(_ context.Context, _ []*big.Int, _ *big.Int, _ uint64, _ []byte) ([]domain.SharedSecret, error) {
	// Own-chain entry resolves to the self-secret (as in production), so the decrypt derives the
	// GCM key from the self-secret for the self batch.
	return []domain.SharedSecret{{ChainId: s.myChainID, Secret: s.selfSecret}}, nil
}

func (s *selfBatchKeysStub) GetEnygmaSharedSelfSecret(_ context.Context, _ *big.Int, _ uint64, _ []byte) (domain.SharedSecret, error) {
	return domain.SharedSecret{ChainId: s.myChainID, Secret: s.selfSecret}, nil
}

func (s *selfBatchKeysStub) GetRaylsViewKeyPair(_ context.Context, _ uint64) (domain.RaylsViewKeyPair, error) {
	return domain.RaylsViewKeyPair{}, nil
}

// TestGCMEncryptEnygma_SelfBatch_RoundTripsUnderEnygmaDecrypt locks the fix for the Enygma self-batch
// AEAD divergence. EnygmaV1's anonymity-set batching always includes the sender's own chain as a
// destination (the change output). Its on-chain message tag — and the dedicated Enygma decrypt — use
// the rotating self-secret; the plain GCMEncrypt seal used the static own-chain pairwise row. Result:
// the own-chain batch failed AEAD on decrypt (OUTCOME_TAMPERED), so it was dropped in normal operation
// and made resync recovery deadlock (it can never re-decrypt the sender's own self/change batch).
//
// GCMEncryptEnygma seals the own-chain batch with the self-secret, restoring
// seal-key == tag-secret == decrypt-key. This test proves:
//  1. a self batch sealed by GCMEncryptEnygma decrypts cleanly under GCMDecryptWithEnygmaMessageTag, and
//  2. the same batch sealed by the old GCMEncrypt path fails AEAD (ErrAuthFailed) — the bug.
//
// Both assertions exercise the real production decrypt path; the on-chain message tag is always the
// self-secret's tag (what GenerateEnygmaSharedSecrets stamps), independent of the seal key.
func TestGCMEncryptEnygma_SelfBatch_RoundTripsUnderEnygmaDecrypt(t *testing.T) {
	t.Parallel()

	myChainID := big.NewInt(12346)
	keysSvc := &selfBatchKeysStub{
		myChainID: myChainID,
		// 31-byte secrets: the on-chain message tag interprets the secret as a Poseidon field
		// element (BN254, ~254-bit), so values must stay under the field modulus.
		pairwiseSelf: []byte("own-chain-pairwise-row-31bytes!"), // 31B — what the buggy seal used
		selfSecret:   []byte("rotating-self-secret-31-bytes!!"), // 31B — what tag + decrypt use
	}
	svc := NewEncryptService(myChainID, keysSvc, selfBatchPassthroughEncryptor{})

	ctx := context.Background()
	const blockNumber = uint64(15558)
	plaintext := []byte("enygma self/change batch payload")
	resourceID := []byte{0x5b, 0xfa, 0x74, 0xc7, 0x43, 0x91, 0x40, 0x28}
	anonymitySet := []*big.Int{myChainID}

	// The on-chain message tag for the sender's own slot is generated from the self-secret
	// (GenerateEnygmaSharedSecrets), regardless of which key seals the ciphertext.
	messageTag, err := generateEnygmaMessageTag(keysSvc.selfSecret, blockNumber)
	require.NoError(t, err)

	// (1) FIX: self batch (chainID == myChainID) sealed via the self-aware GCMEncryptEnygma.
	fixedCipher, err := svc.GCMEncryptEnygma(ctx, plaintext, myChainID.Uint64(), blockNumber, resourceID)
	require.NoError(t, err)

	decrypted, err := svc.GCMDecryptWithEnygmaMessageTag(ctx, fixedCipher, messageTag, blockNumber, anonymitySet, myChainID, resourceID)
	require.NoError(t, err, "self batch sealed with the self-secret must decrypt cleanly under the Enygma decrypt")
	assert.Equal(t, plaintext, decrypted, "round-trip must preserve the self/change batch payload")

	// (2) BUG: the same self batch sealed by the old plain GCMEncrypt (own-chain pairwise row) must
	//     fail AEAD under the Enygma decrypt — the tag still matches (self-secret), but the GCM key
	//     (self-secret) cannot open a ciphertext sealed with the pairwise row.
	buggyCipher, err := svc.GCMEncrypt(ctx, plaintext, myChainID.Uint64(), blockNumber)
	require.NoError(t, err)

	_, err = svc.GCMDecryptWithEnygmaMessageTag(ctx, buggyCipher, messageTag, blockNumber, anonymitySet, myChainID, resourceID)
	require.Error(t, err, "self batch sealed with the pairwise row must fail the self-secret decrypt (the bug)")
	assert.True(t, errors.Is(err, ErrAuthFailed), "tag matches but AEAD fails ⇒ ErrAuthFailed (OUTCOME_TAMPERED); got %v", err)
}
