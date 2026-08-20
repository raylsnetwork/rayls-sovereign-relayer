package proofs

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
)

// Import must return an error on any input that fails to deserialize (nil,
// empty, or malformed bytes) so callers can reject it before use.
func TestProofDBImport_ErrorsOnBadInput(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		data []byte
	}{
		{"nil bytes", nil},
		{"empty bytes", []byte{}},
		{"not json", []byte("not-json")},
		{"truncated json", []byte(`{"aa":`)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := Import(tt.data)
			if err == nil {
				t.Fatalf("Import(%q) = nil error, want error", tt.data)
			}
			if got != nil {
				t.Fatalf("Import(%q) = %v, want nil ProofDB", tt.data, got)
			}
		})
	}
}

func TestProofDBImport_RoundTripsValidData(t *testing.T) {
	t.Parallel()
	orig := NewProofDB()
	if err := orig.Put([]byte{0x01, 0x02}, []byte("value")); err != nil {
		t.Fatalf("Put: %v", err)
	}
	exported, err := orig.Export()
	if err != nil {
		t.Fatalf("Export: %v", err)
	}
	got, err := Import(exported)
	if err != nil {
		t.Fatalf("Import(valid) = %v, want nil error", err)
	}
	if got == nil {
		t.Fatal("Import(valid) = nil, want non-nil")
	}
	v, err := got.Get([]byte{0x01, 0x02})
	if err != nil {
		t.Fatalf("Get after round-trip: %v", err)
	}
	if string(v) != "value" {
		t.Fatalf("Get = %q, want %q", v, "value")
	}
}

// A nil *ProofDB must surface an error from its read methods instead of
// panicking on the nil receiver (go-ethereum's trie.VerifyProof calls Get/Has).
func TestNilProofDB_ReadMethodsErrorNotPanic(t *testing.T) {
	t.Parallel()
	var db *ProofDB

	if _, err := db.Get([]byte{0x01}); err == nil {
		t.Error("(*ProofDB)(nil).Get = nil error, want error")
	}
	if _, err := db.Has([]byte{0x01}); err == nil {
		t.Error("(*ProofDB)(nil).Has = nil error, want error")
	}
}

func TestVerifyProof_NilProofDB_ErrorsNotPanic(t *testing.T) {
	t.Parallel()
	if _, err := VerifyProof(common.Hash{}, []byte{0x80}, nil); err == nil {
		t.Error("VerifyProof(nil proof) = nil error, want error")
	}
}

// Valid-but-empty JSON ("null"/"{}") deserializes to a non-nil, empty *ProofDB
// (no error), so VerifyProof against it must error, not panic.
func TestProofDBImport_EmptyJSON_ReturnsNonNilAndVerifyErrors(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		data []byte
	}{
		{"json null", []byte("null")},
		{"empty object", []byte("{}")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			db, err := Import(tt.data)
			if err != nil {
				t.Fatalf("Import(%q) = %v, want nil error", tt.data, err)
			}
			if db == nil {
				t.Fatalf("Import(%q) = nil, want non-nil empty ProofDB", tt.data)
			}
			if _, err := VerifyProof(common.Hash{}, []byte{0x80}, db); err == nil {
				t.Errorf("VerifyProof(empty ProofDB from %q) = nil error, want error", tt.data)
			}
		})
	}
}
