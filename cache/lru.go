// Package cache provides a generic, thread-safe LRU (Least Recently Used) cache implementation.
package cache

import (
	"container/list"
	"sync"
)

// entry holds a key-value pair stored in the LRU list.
type entry[K comparable, V any] struct {
	key   K
	value V
}

// LRU is a fixed-capacity, thread-safe cache that evicts the least recently used entry when full.
// K must be comparable; V can be any type.
type LRU[K comparable, V any] struct {
	capacity int
	ll       *list.List          // doubly-linked list of *entry[K,V]
	cache    map[K]*list.Element // map of key to list element
	mu       sync.Mutex          // guards ll and cache
}

// New creates a new LRU cache with the given capacity. Panics if capacity < 1.
func NewLRU[K comparable, V any](capacity int) *LRU[K, V] {
	if capacity < 1 {
		panic("cache: capacity must be at least 1")
	}
	return &LRU[K, V]{
		capacity: capacity,
		ll:       list.New(),
		cache:    make(map[K]*list.Element, capacity),
	}
}

// Get looks up a key's value and marks it as most recently used.
// The found boolean indicates if the key was present.
func (c *LRU[K, V]) Get(key K) (value V, found bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if ele, ok := c.cache[key]; ok {
		c.ll.MoveToFront(ele)
		return ele.Value.(*entry[K, V]).value, true
	}
	// zero value if not found
	var zero V
	return zero, false
}

// Put inserts or updates a key-value pair, marking it as most recently used.
// If the cache is at capacity, it evicts the least recently used item.
func (c *LRU[K, V]) Put(key K, value V) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Update existing entry
	if ele, ok := c.cache[key]; ok {
		c.ll.MoveToFront(ele)
		ele.Value.(*entry[K, V]).value = value
		return
	}

	// Evict oldest if at capacity
	if c.ll.Len() >= c.capacity {
		c.removeOldest()
	}

	// Insert new entry at front
	ent := &entry[K, V]{key: key, value: value}
	ele := c.ll.PushFront(ent)
	c.cache[key] = ele
}

// Remove deletes a key from the cache if it exists.
func (c *LRU[K, V]) Remove(key K) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if ele, ok := c.cache[key]; ok {
		c.removeElement(ele)
	}
}

// Len returns the current number of items in the cache.
func (c *LRU[K, V]) Len() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.ll.Len()
}

// Clear resets the cache to empty state.
func (c *LRU[K, V]) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.ll.Init()
	for k := range c.cache {
		delete(c.cache, k)
	}
}

// removeOldest removes the least recently used element.
func (c *LRU[K, V]) removeOldest() {
	if ele := c.ll.Back(); ele != nil {
		c.removeElement(ele)
	}
}

// removeElement removes the specified list element and its map entry.
func (c *LRU[K, V]) removeElement(e *list.Element) {
	c.ll.Remove(e)
	ent := e.Value.(*entry[K, V])
	delete(c.cache, ent.key)
}
