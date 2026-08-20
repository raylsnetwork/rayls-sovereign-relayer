package service

import "errors"

// ErrNotForRecipient signals that the encrypted payload was not addressed to
// this CTS instance. Returned in two situations:
//   - Message-tag path: none of the shared secrets we hold could produce the
//     supplied fingerprint/message tag, so the message is provably not for us.
//   - Salt-based DVP path: the salt we derived doesn't open the AEAD envelope,
//     which on this path is the only "not for me" signal we have.
//
// Callers should treat this as a normal protocol outcome — in any N-participant
// network the typical decrypt attempt is "not for me" for N-1 participants.
var ErrNotForRecipient = errors.New("message not for this recipient")

// ErrAuthFailed signals AEAD authentication failure in a context where it is
// genuinely suspicious — i.e. the message-tag matched one of our keys but the
// ciphertext still failed to verify. This implies a tag collision, a stale-key
// replay, or actual tampering, and warrants a louder log / alerting hook.
//
// In the salt-based DVP path the same underlying cryptographic event is
// reported as ErrNotForRecipient instead — see the doc on the service-level
// decrypt methods for the mapping rationale.
var ErrAuthFailed = errors.New("authenticated decryption failed")
