// Decommissioning Teleport (vanilla, atomic).

package repository

type RevertSignature struct {
	ID   string `db:"id"`
	Data []byte `db:"data"`
}
