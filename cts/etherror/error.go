package etherror

import (
	"strings"
)

// EthError is a logical node error together with every concrete wire string the
// nodes we talk to are known to emit for it.
//
// The canonical variants are validated against the live nodes at startup by
// VerifyLedgerErrorStrings (see probe.go);
type EthError struct {
	name     string
	variants []string
}

func (e EthError) Name() string { return e.name }

func (e EthError) Variants() []string { return e.variants }

func (e EthError) Matches(err error) bool {
	if err == nil {
		return false
	}
	msg := strings.ToLower(err.Error())
	for _, v := range e.variants {
		if strings.Contains(msg, strings.ToLower(v)) {
			return true
		}
	}
	return false
}

var (
	// AlreadyKnownError: the node already holds this exact tx (same hash) in its
	// mempool — an idempotent re-broadcast of identical signed bytes.
	AlreadyKnownError = EthError{
		name: "already known",
		variants: []string{
			"Known transaction",
			"already known",
			"already imported",
		},
	}
	// NonceTooLowError: the chain has already moved past this nonce.
	NonceTooLowError = EthError{
		name: "nonce too low",
		variants: []string{
			"Nonce too low",
			"nonce too low",
		},
	}
)

// Is reports whether err matches any known variant of target. nil never matches.
func Is(err error, target EthError) bool {
	return target.Matches(err)
}
