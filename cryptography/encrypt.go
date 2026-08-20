package cryptography

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"fmt"
	"io"
	"strings"

	"golang.org/x/crypto/hkdf"
	"golang.org/x/crypto/sha3"

	"github.com/raylsnetwork/rayls-privacy-relayer-api/withstack"
)

const minKeyLength = 32

// ErrAuthFailed signals that AEAD authentication failed — either the ciphertext
// was decrypted with the wrong key/salt, or it was tampered with after sealing.
// Go's stdlib uses an unexported sentinel for this; translating it here at the
// boundary keeps all "message authentication failed" string-matching in one place.
var ErrAuthFailed = errors.New("aead authentication failed")

// GenerateSecureRandomValue securely generates a random value of size in bytes
func GenerateSecureRandomValue(size int) ([]byte, error) {
	x := make([]byte, size)

	_, err := rand.Read(x)
	if err != nil {
		return nil, withstack.Wrap(fmt.Errorf("generating random bytes: %w", err))
	}

	return x, nil
}

func HashIt(data []byte) []byte {
	h := sha3.New256()
	h.Write(data)

	return h.Sum(nil)
}

func GetSharedMessageTag(sharedSecret []byte, randomness []byte) []byte {
	h := sha3.New256()
	h.Write(sharedSecret)
	h.Write(randomness)

	return h.Sum(nil)
}

// EncryptGCM encrypts a plaintext message using AES GCM ()
func EncryptGCM(key []byte, plaintext []byte, associatedData []byte) ([]byte, error) {
	// we use a 32 bytes key (which instantiates AES-256) for quantum security
	if len(key) < minKeyLength {
		key = HashIt(key)
	}

	// Generate a new AES cipher using the given key
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("failed to create AES cipher: %w", err)
	}

	// Generate a new GCM cipher
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("failed to create GCM cipher: %w", err)
	}

	// Generate a new nonce
	nonce, err := GenerateSecureRandomValue(gcm.NonceSize())
	if err != nil {
		return nil, fmt.Errorf("failed to generate nonce: %w", err)
	}

	// Seal the plaintext to the nonce and associated data
	ciphertext := gcm.Seal(nil, nonce, plaintext, associatedData)

	// Append associated data and nonce to ciphertext
	nctxt := make([]byte, 0, len(nonce)+len(ciphertext))
	nctxt = append(nctxt, nonce...)
	nctxt = append(nctxt, ciphertext...)
	payload := make([]byte, 0, len(associatedData)+len(nctxt))
	payload = append(payload, associatedData...)
	payload = append(payload, nctxt...)

	return payload, nil
}

// * // DecryptGCM receives a ciphertext and returns the plaintext (and potentially the associated data)
// this ciphertext has the following form: AD (16 bytes) || Nonce (12 bytes) || Encrypted Msg (remainder bytes)
func DecryptGCM(ciphertext []byte, key []byte) ([]byte, []byte, error) {
	// Generate a new AES cipher using the given key
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, nil, withstack.Wrap(fmt.Errorf("creating AES cipher: %w", err))
	}

	// Generate a new GCM cipher
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, nil, withstack.Wrap(fmt.Errorf("creating GCM: %w", err))
	}

	// Split associated data, nonce, and ciphertext
	associatedData := ciphertext[0:16]
	nonce := ciphertext[16 : 16+gcm.NonceSize()]
	ctxt := ciphertext[16+gcm.NonceSize():]

	// Open the ciphertext and return the plaintext
	plaintext, err := gcm.Open(nil, nonce, ctxt, associatedData)
	if err != nil {
		// crypto/cipher uses an unexported errOpen for AEAD verification failure.
		// Translate it to our typed sentinel here — this is the only place in the
		// codebase that matches against the stdlib's "message authentication failed"
		// string.
		if strings.Contains(err.Error(), "message authentication failed") {
			return nil, nil, ErrAuthFailed
		}
		return nil, nil, withstack.Wrap(fmt.Errorf("decrypting GCM: %w", err))
	}
	return associatedData, plaintext, nil
}

func DeriveSymmetricKey(sharedSecret []byte) ([]byte, error) {
	// public context for HKDF
	context := []byte("Rayls")

	reader := hkdf.New(sha3.New256, sharedSecret, nil, context)

	key := make([]byte, 32)

	if _, err := io.ReadFull(reader, key); err != nil {
		return nil, withstack.Wrap(fmt.Errorf("reading derived key from HKDF: %w", err))
	}

	return key, nil
}

// KMAC is a function to compute a keyed hash using SHA3-256
// receives the key (in []byte format) and the data (in []byte format) to be hashed
// returns the SHA3 message authentication code (in []byte format)
func KMAC(key []byte, data []byte) []byte {
	hash := sha3.New256()
	hash.Write(key)
	hash.Write(data)
	mac := hash.Sum(nil)

	return mac
}
