package testutil

type CeaserEncryptor struct {
	shift int
}

func NewCeaserEncryptor(shift int) *CeaserEncryptor {
	return &CeaserEncryptor{shift: shift}
}

func (e *CeaserEncryptor) Encrypt(plaintext []byte) ([]byte, error) {
	return applyCeaserCipher(plaintext, e.shift), nil
}

func (e *CeaserEncryptor) Decrypt(ciphertext []byte) ([]byte, error) {
	return applyCeaserCipher(ciphertext, -e.shift), nil
}

// byteRange is the total number of possible byte values (0-255).
const byteRange = 256

func applyCeaserCipher(input []byte, shift int) []byte {
	output := make([]byte, len(input))
	for i, char := range input {
		output[i] = byte((int(char) + shift + byteRange) % byteRange)
	}
	return output
}
