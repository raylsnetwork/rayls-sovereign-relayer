package service_test

import (
	"testing"

	"github.com/raylsnetwork/rayls-privacy-relayer-api/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// signatureVerifier is a test double that verifies signatures are received exactly once
// with the correct data and types
type signatureVerifier struct {
	t                   *testing.T
	wantLockData        []byte
	wantRevertData      []byte
	seenUnlockSignature bool
	seenRevertSignature bool
	unlockSignatureData []byte
	revertSignatureData []byte
}

func newSignatureVerifier(t *testing.T, wantLockData, wantRevertData []byte) *signatureVerifier {
	return &signatureVerifier{
		t:              t,
		wantLockData:   wantLockData,
		wantRevertData: wantRevertData,
	}
}

// verify is called when BatchCreate is invoked and verifies the signatures
func (v *signatureVerifier) verify(signatures []types.CalldataSignature) error {
	v.t.Helper()

	// Should receive exactly 2 signatures
	require.Len(v.t, signatures, 2, "should receive exactly 2 signatures")

	for _, sig := range signatures {
		switch sig.SignatureType {
		case types.UnlockOnDestinationSide:
			assert.False(v.t, v.seenUnlockSignature, "UnlockOnDestinationSide signature should only be seen once")
			v.seenUnlockSignature = true
			v.unlockSignatureData = sig.Signature
			assert.Equal(v.t, v.wantLockData, sig.Signature, "UnlockOnDestinationSide signature should have lock data")
		case types.RevertOnDestinationSide:
			assert.False(v.t, v.seenRevertSignature, "RevertOnDestinationSide signature should only be seen once")
			v.seenRevertSignature = true
			v.revertSignatureData = sig.Signature
			assert.Equal(
				v.t,
				v.wantRevertData,
				sig.Signature,
				"RevertOnDestinationSide signature should have revert data",
			)
		default:
			assert.Fail(v.t, "unexpected signature type", "got type: %v", sig.SignatureType)
		}
	}

	return nil
}

// assertAllSignaturesSeen verifies that both signatures were received
func (v *signatureVerifier) assertAllSignaturesSeen() {
	v.t.Helper()
	assert.True(v.t, v.seenUnlockSignature, "UnlockOnDestinationSide signature was not received")
	assert.True(v.t, v.seenRevertSignature, "RevertOnDestinationSide signature was not received")
}
