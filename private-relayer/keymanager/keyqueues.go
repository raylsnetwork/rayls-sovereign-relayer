// Package keymanager owns cross-chain signing-key management for the
// private relayer.
package keymanager

import "github.com/raylsnetwork/rayls-privacy-relayer-api/keyqueue"

// KeyQueues bundles the per-chain private-key queues used by the source
// and destination relayer halves. PrivateHub signs hub-side txs (the
// source relayer's executor target), PrivateNode signs node-side txs
// (the destination relayer's executor target), and DvpOperator signs
// DVP-operator-scoped txs.
//
// Fields are exported so the cmd/run binary and integration tests can
// construct the struct directly.
type KeyQueues struct {
	PrivateHub  *keyqueue.RaylsSignPrivateKeyManager
	PrivateNode *keyqueue.RaylsSignPrivateKeyManager
	DvpOperator *keyqueue.RaylsSignPrivateKeyManager
}
