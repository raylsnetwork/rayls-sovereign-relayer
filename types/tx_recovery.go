package types

// TxRecoveryData holds the crash-recovery state for a pending blockchain transaction.
type TxRecoveryData struct {
	ResourceID            string          // Resource identifier (token, swap, etc.)
	PrivateHubBlockNumber uint64          // Block number on the Private Network Hub
	FromChainID           string          // Source chain ID
	PrivateHubTxHash      string          // On-chain transaction hash (primary key)
	TxBytes               []byte          // RLP-encoded signed transaction to re-broadcast
	EventType             EnygmaEventType // Event type for categorization/querying
	TxNature              TxNature        // Domain category (enygma, dvp, other)
	Status                HistoryStatus   // Pending or Confirmed
}

// TxNature categorizes the domain of a recovery record.
type TxNature string

const (
	TxNatureEnygma TxNature = "enygma"
	TxNatureDvp    TxNature = "dvp"
	TxNatureOther  TxNature = "other"
)

// TxNatureFromEventType derives the TxNature from an EnygmaEventType.
func TxNatureFromEventType(eventType EnygmaEventType) TxNature {
	switch {
	case eventType <= EnygmaWithdrawFromDvp:
		return TxNatureEnygma
	case eventType >= DvpERC721Withdraw && eventType <= DvpSwapCompletion:
		return TxNatureDvp
	default:
		return TxNatureOther
	}
}
