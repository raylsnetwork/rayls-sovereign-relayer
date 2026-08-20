// Package noncecache holds the per-address local nonce counter used by
// every component that signs transactions for the same keyQueue/identity
// (AuthGen, RPCBatcher, reaper resends, ...). Sharing a single instance
// across those components closes a cross-path race — without sharing,
// each component would maintain its own counter and concurrent paths
// could read the same chain-pending nonce and broadcast colliding txs.
package noncecache

import (
	"context"
	"fmt"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/raylsnetwork/rayls-sovereign-relayer/withstack"
)

// Seeder is the subset of *ethclient.Client (and friends) needed to seed
// an empty cache from the chain's pending state on first use.
type Seeder interface {
	PendingNonceAt(ctx context.Context, addr common.Address) (uint64, error)
}

// Cache is a goroutine-safe per-address nonce counter. The zero value is
// not usable — construct via New.
type Cache struct {
	mu     sync.Mutex
	nonces map[common.Address]uint64
	seeder Seeder
}

// New returns a fresh empty cache.
func New(seeder Seeder) *Cache {
	return &Cache{
		nonces: make(map[common.Address]uint64),
		seeder: seeder,
	}
}

// Reserve returns the starting nonce for `count` txs from `addr` and
// advances the cache by `count` atomically. On cache miss it seeds the
// counter from PendingNonceAt(addr) — the RPC happens outside the
// critical section, and a racy double-init resolves to max(existing,
// fresh) so the counter never goes backwards. The seed write and the
// reservation write happen in a single critical section so the common
// uncontended cache-miss path costs exactly one Lock/Unlock pair after
// the RPC, not two.
//
// Callers that fail to broadcast any of the reserved nonces must call
// Invalidate(addr) — otherwise the cache will stride past nonces that
// were never consumed, creating gaps that stall mining.
func (c *Cache) Reserve(ctx context.Context, addr common.Address, count int) (uint64, error) {
	if count <= 0 {
		return 0, fmt.Errorf("nonce reserve count must be > 0, got %d", count)
	}

	c.mu.Lock()
	cached, ok := c.nonces[addr]
	c.mu.Unlock()

	if !ok {
		// Slow path: seed from chain WITHOUT holding the mutex so a
		// hung RPC doesn't block other addresses from advancing.
		onChain, err := c.seeder.PendingNonceAt(ctx, addr)
		if err != nil {
			return 0, withstack.Wrap(fmt.Errorf("seeding nonce cache for %s: %w", addr.Hex(), err))
		}
		c.mu.Lock()
		// Defensive re-read: another goroutine may have seeded the
		// cache between our unlock above and re-lock here. Take
		// whichever is higher — the counter must never go backwards.
		if existing, set := c.nonces[addr]; set && existing > onChain {
			cached = existing
		} else {
			cached = onChain
		}
		start := cached
		c.nonces[addr] = start + uint64(count)
		c.mu.Unlock()
		return start, nil
	}

	// Fast path: cache already populated. Same defensive re-read in
	// case a different goroutine advanced the counter past us between
	// the initial read and this lock.
	//
	// Invariant relied on by the caller: all Reserve / Invalidate calls
	// for a given address are serialised through the keyQueue (only one
	// goroutine holds a given private key at a time), so no concurrent
	// Invalidate(addr) can fire while this fast path is running. Without
	// that invariant the `set` check on the re-read could miss an
	// invalidation that erased `addr` between the two locks, causing the
	// stale `cached` value to be written back and undo the invalidation.
	c.mu.Lock()
	if current, set := c.nonces[addr]; set && current > cached {
		cached = current
	}
	start := cached
	c.nonces[addr] = start + uint64(count)
	c.mu.Unlock()
	return start, nil
}

// Invalidate forces the next Reserve(addr) to re-query the chain. Call
// this whenever a previously-reserved nonce was NOT broadcast (sign
// failure, pre-flight revert, dropped batch, etc.) so the next caller
// realigns with the chain's actual pending state.
func (c *Cache) Invalidate(addr common.Address) {
	c.mu.Lock()
	delete(c.nonces, addr)
	c.mu.Unlock()
}
