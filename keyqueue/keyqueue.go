package keyqueue

import (
	"context"
	"crypto/ecdsa"
	"log/slog"
	"time"
)

type RaylsSignPrivateKeyManager struct {
	keys chan *ecdsa.PrivateKey
}

func New(size int) *RaylsSignPrivateKeyManager {
	return &RaylsSignPrivateKeyManager{
		keys: make(chan *ecdsa.PrivateKey, size),
	}
}

func (kq *RaylsSignPrivateKeyManager) Enqueue(key *ecdsa.PrivateKey) {
	kq.keys <- key
}

func (kq *RaylsSignPrivateKeyManager) Dequeue(ctx context.Context) (*ecdsa.PrivateKey, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case key := <-kq.keys:
		return key, nil
	}
}

func (kq *RaylsSignPrivateKeyManager) Size() int {
	return len(kq.keys)
}

// StartMonitoring monitors the key queues and logs warnings when they become empty
func (kq *RaylsSignPrivateKeyManager) StartMonitoring(ctx context.Context, queueName string, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	slog.Info("Started key queue monitoring", "Queue", queueName)

	for {
		select {
		case <-ctx.Done():
			slog.Info("Stopping key queue monitoring", "Queue", queueName)
			return
		case <-ticker.C:
			keysSize := kq.Size()

			if keysSize == 0 {
				slog.Warn("⚠️  CRITICAL: Key queue is EMPTY - system may be unable to process transactions",
					"Queue", queueName,
					"Type", "keys",
					"Size", keysSize)
			}
			slog.Debug("Key queue status", "Queue", queueName, "KeysSize", keysSize)
		}
	}
}
