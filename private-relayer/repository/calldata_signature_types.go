// Decommissioning Teleport (vanilla, atomic).

package repository

// Deprecated: Decommissioning Teleport (vanilla, atomic).
const CalldataSignatureCollectionName = "calldata_signatures"

// Deprecated: Decommissioning Teleport (vanilla, atomic).
type CalldataSignature struct {
	SharedId                string `db:"shared_id"`
	Status                  uint8  `db:"status"`
	Signature               []byte `db:"signature"`
	ResourceId              []byte `db:"resource_id"`
	SignatureExecuteChainId string `db:"signature_execute_chain_id"`
	DestinationChainId      string `db:"destination_chain_id"`
	SignatureType           uint8  `db:"signature_type"`
}
