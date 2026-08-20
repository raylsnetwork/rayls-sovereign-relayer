package service

import "errors"

// Used to signal that the operation failed but can be retried
// later and potentially succeed. Intentionally kept vague so it
// can be used for different scenarios (e.g. DB connection error, FI point)
var ErrRetriable = errors.New("operation failed, retry may succeed")

// ErrTxAlreadyExists signals that a row with the supplied idempotency id
// already exists in the sync-tx repository. Returned by CreateTransaction;
// the caller treats this as a retry signal and falls through to the
// recovery path (load previously-stored result, or resume waiting on the
// previously-stored tx hash).
var ErrTxAlreadyExists = errors.New("tx already exists for id")

// ErrTxNotFound signals that no row exists for the supplied id. Returned
// by GetTransaction and GetResult.
var ErrTxNotFound = errors.New("tx not found for id")

// ErrNoResult signals that a row exists for the supplied id but its
// result has not been stored yet — i.e. the previous attempt crashed
// after CreateTransaction succeeded but before SetResult was reached.
// Returned by GetResult only; the caller resumes by fetching the receipt
// from the chain using the previously-stored hash.
var ErrNoResult = errors.New("tx exists but result not yet stored")

// ErrStaleTransition signals that Save found the row's version had advanced
// since it was loaded — a concurrent writer (crash recovery, a parallel
// same-id caller, the retention job) modified the row in between. The caller
// should re-Get, re-evaluate (the row may now be terminal), and retry.
var ErrStaleTransition = errors.New("stale transition: row version advanced since load")
