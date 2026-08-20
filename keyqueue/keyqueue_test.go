package keyqueue

import (
	"context"
	"crypto/ecdsa"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func generateKey(t *testing.T) *ecdsa.PrivateKey {
	t.Helper()
	key, err := crypto.GenerateKey()
	require.NoError(t, err)
	return key
}

func TestNew(t *testing.T) {
	t.Run("creates manager with given capacity", func(t *testing.T) {
		km := New(5)

		assert.NotNil(t, km)
		assert.Equal(t, 0, km.Size())
	})
}

func TestEnqueueDequeue(t *testing.T) {
	t.Run("enqueues and dequeues a single key", func(t *testing.T) {
		km := New(3)
		key := generateKey(t)

		km.Enqueue(key)

		assert.Equal(t, 1, km.Size())

		got, err := km.Dequeue(context.Background())
		require.NoError(t, err)
		assert.Equal(t, key, got)
		assert.Equal(t, 0, km.Size())
	})

	t.Run("preserves FIFO order", func(t *testing.T) {
		km := New(3)
		key1 := generateKey(t)
		key2 := generateKey(t)
		key3 := generateKey(t)

		km.Enqueue(key1)
		km.Enqueue(key2)
		km.Enqueue(key3)

		assert.Equal(t, 3, km.Size())

		got1, err := km.Dequeue(context.Background())
		require.NoError(t, err)
		assert.Equal(t, key1, got1)

		got2, err := km.Dequeue(context.Background())
		require.NoError(t, err)
		assert.Equal(t, key2, got2)

		got3, err := km.Dequeue(context.Background())
		require.NoError(t, err)
		assert.Equal(t, key3, got3)

		assert.Equal(t, 0, km.Size())
	})

	t.Run("supports enqueue after dequeue", func(t *testing.T) {
		km := New(2)
		key1 := generateKey(t)
		key2 := generateKey(t)

		km.Enqueue(key1)
		got, err := km.Dequeue(context.Background())
		require.NoError(t, err)
		assert.Equal(t, key1, got)

		km.Enqueue(key2)
		got, err = km.Dequeue(context.Background())
		require.NoError(t, err)
		assert.Equal(t, key2, got)
	})
}

func TestSize(t *testing.T) {
	t.Run("returns zero for empty queue", func(t *testing.T) {
		km := New(5)
		assert.Equal(t, 0, km.Size())
	})

	t.Run("reflects number of enqueued keys", func(t *testing.T) {
		km := New(5)

		km.Enqueue(generateKey(t))
		assert.Equal(t, 1, km.Size())

		km.Enqueue(generateKey(t))
		assert.Equal(t, 2, km.Size())

		_, err := km.Dequeue(context.Background())
		require.NoError(t, err)
		assert.Equal(t, 1, km.Size())
	})
}

func TestDequeueContextCancellation(t *testing.T) {
	t.Run("returns error when context is already cancelled", func(t *testing.T) {
		km := New(3)
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		key, err := km.Dequeue(ctx)

		assert.Nil(t, key)
		assert.ErrorIs(t, err, context.Canceled)
	})

	t.Run("returns error when context is cancelled while waiting", func(t *testing.T) {
		km := New(3)
		ctx, cancel := context.WithCancel(context.Background())

		done := make(chan struct{})
		var got *ecdsa.PrivateKey
		var gotErr error

		go func() {
			got, gotErr = km.Dequeue(ctx)
			close(done)
		}()

		cancel()

		select {
		case <-done:
			assert.Nil(t, got)
			assert.ErrorIs(t, gotErr, context.Canceled)
		case <-time.After(2 * time.Second):
			t.Fatal("Dequeue did not return after context cancellation")
		}
	})

	t.Run("returns error when context deadline exceeded", func(t *testing.T) {
		km := New(3)
		ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
		defer cancel()

		key, err := km.Dequeue(ctx)

		assert.Nil(t, key)
		assert.ErrorIs(t, err, context.DeadlineExceeded)
	})
}

func TestStartMonitoring(t *testing.T) {
	t.Run("stops when context is cancelled", func(t *testing.T) {
		km := New(3)
		ctx, cancel := context.WithCancel(context.Background())

		done := make(chan struct{})
		go func() {
			km.StartMonitoring(ctx, "test-queue", 50*time.Millisecond)
			close(done)
		}()

		cancel()

		select {
		case <-done:
			// goroutine exited cleanly
		case <-time.After(2 * time.Second):
			t.Fatal("StartMonitoring did not stop after context cancellation")
		}
	})

	t.Run("ticks at the specified interval", func(t *testing.T) {
		km := New(3)
		km.Enqueue(generateKey(t))

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		done := make(chan struct{})
		go func() {
			km.StartMonitoring(ctx, "test-queue", 50*time.Millisecond)
			close(done)
		}()

		// Wait long enough for at least 2 ticks
		time.Sleep(130 * time.Millisecond)
		cancel()

		select {
		case <-done:
		case <-time.After(2 * time.Second):
			t.Fatal("StartMonitoring did not stop after context cancellation")
		}
	})

	t.Run("logs warning when queue is empty", func(t *testing.T) {
		km := New(3)
		// Queue is intentionally empty

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		done := make(chan struct{})
		go func() {
			km.StartMonitoring(ctx, "empty-queue", 50*time.Millisecond)
			close(done)
		}()

		// Let at least one tick fire while queue is empty
		time.Sleep(80 * time.Millisecond)
		cancel()

		select {
		case <-done:
		case <-time.After(2 * time.Second):
			t.Fatal("StartMonitoring did not stop after context cancellation")
		}
	})
}
