package contractclient

import (
	"crypto/sha256"
	"encoding/hex"
	"sort"
	"strings"
)

// IDFor builds a deterministic idempotency id of the form
// "<op>:<part>:<part>". The op is a stable per-call-site operation name
// (e.g. "accessmanager.GrantAuthorizedSenderRole") and the parts are stable
// domain identifiers — never UUIDs, proofs, ciphertext, or anything else that
// changes between retries of the same logical operation. CTS uses this id to
// suppress duplicate SignAndSend calls, so two genuinely different operations
// must never produce the same id.
func IDFor(op string, parts ...string) string {
	if len(parts) == 0 {
		return op
	}
	return op + ":" + strings.Join(parts, ":")
}

// HashIDs deterministically folds a set of ids into a single hex token,
// independent of their order. Used by batch call sites whose identity is the
// set of sharedIDs they carry. The result is stable only if the batch
// membership is stable across retries.
func HashIDs(ids []string) string {
	sorted := append([]string(nil), ids...)
	sort.Strings(sorted)
	sum := sha256.Sum256([]byte(strings.Join(sorted, "\x00")))
	return hex.EncodeToString(sum[:])
}
