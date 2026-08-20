package contractclient

import (
	"context"
	"crypto/ecdsa"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
)

// keyQueue represents a queue of private keys used for signing transactions.
// It is shared across clients that need to dequeue a key to build auth opts.
type keyQueue interface {
	Enqueue(key *ecdsa.PrivateKey)
	Dequeue(context.Context) (*ecdsa.PrivateKey, error)
}

// authGen abstracts generation of transaction options from a private key.
// It is shared between the Executor and legacy clients that still build auth directly.
type authGen interface {
	Get(context.Context, *ecdsa.PrivateKey) (*bind.TransactOpts, error)
}

