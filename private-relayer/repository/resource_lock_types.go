package repository

import "time"

const ResourceLockCollectionName = "resource_locks"

type ResourceLock struct {
	ResourceId string    `db:"resource_id"`
	ExpiresAt  time.Time `db:"expires_at"`
}
