// Decommissioning Teleport (vanilla, atomic).

package repository

// Deprecated: Decommissioning Teleport (vanilla, atomic).
type AtomicSUM struct {
	SharedId    string `db:"shared_id"`
	Status      uint8  `db:"status"`
	IsProcessed bool   `db:"is_processed"`
}
