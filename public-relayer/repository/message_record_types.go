// Decommissioning Teleport (vanilla, atomic).

package repository

import "time"

type MessageRecord struct {
	ID          string    `db:"id"`
	Status      int       `db:"status"`
	ForwardHash string    `db:"forward_hash"`
	RevertHash  string    `db:"revert_hash"`
	Error       string    `db:"error"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}
