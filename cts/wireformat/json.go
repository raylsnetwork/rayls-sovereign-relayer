// Package wireformat holds the JSON envelope helpers that translate between
// typed Go values and the raw `bytes plaintext` field on CTS Encrypt/Decrypt
// proto messages.
package wireformat

import "encoding/json"

// MarshalPlaintext serializes v as JSON for the plaintext field of a CTS
// encrypt request. Failure is effectively unreachable for the typed DTOs the
// relayer marshals today; the error is surfaced so that a future broken type
// (chan/func fields, NaN floats, a MarshalJSON returning an error) shows up
// loudly at the call site instead of producing a nil-plaintext gRPC
// InvalidArgument from CTS.
func MarshalPlaintext(v any) ([]byte, error) {
	return json.Marshal(v)
}

// Plaintexter is satisfied by every generated CTS Decrypt*Response proto type
// (DecryptResponse, DecryptWithoutFPResponse, DecryptWithoutFPWithSSResponse,
// DecryptEnygmaTransferBatchResponse) via their GetPlaintext accessor.
type Plaintexter interface {
	GetPlaintext() []byte
}

// UnmarshalPlaintext JSON-decodes resp.GetPlaintext() into T.
func UnmarshalPlaintext[T any](resp Plaintexter) (T, error) {
	var v T
	err := json.Unmarshal(resp.GetPlaintext(), &v)
	return v, err
}
