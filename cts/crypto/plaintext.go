package crypto

type PlaintextEncryptor struct{}

func (e *PlaintextEncryptor) Encrypt(plaintext []byte) ([]byte, error) {
	return plaintext, nil
}

func (e *PlaintextEncryptor) Decrypt(ciphertext []byte) ([]byte, error) {
	return ciphertext, nil
}
