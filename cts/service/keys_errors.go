package service

type KeysServiceError struct {
	msg string
}

func NewKeyServiceError(msg string) *KeysServiceError {
	return &KeysServiceError{
		msg: msg,
	}
}

func (e *KeysServiceError) Error() string {
	return e.msg
}

type KeyRepositoryError struct {
	msg string
}

func NewKeyRepositoryError(msg string) *KeyRepositoryError {
	return &KeyRepositoryError{
		msg: msg,
	}
}

func (e *KeyRepositoryError) Error() string {
	return e.msg
}

var (
	ErrNoRaylsViewKeysSet           = NewKeyServiceError("no rayls view keys have been set")
	ErrNoApplicableRaylsViewKeys    = NewKeyRepositoryError("no applicable rayls view key for block number")
	ErrNoRaylsSignKeys              = NewKeyServiceError("no rayls sign keys found")
	ErrNoPaymentSpendKeySet         = NewKeyServiceError("no payment spend key set")
	ErrRaylsSignKeysAlreadyExists   = NewKeyRepositoryError("rayls sign keys already exist")
	ErrNoApplicableSharedSecret     = NewKeyRepositoryError("no applicable shared secret for block number")
	ErrNoApplicableEnygmaSelfSecret = NewKeyRepositoryError("no applicable enygma self secret for block number")
)
