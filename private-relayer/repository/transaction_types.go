package repository

import (
	"time"

	"github.com/raylsnetwork/rayls-privacy-relayer-api/types"
)

const TransactionsCollectionName = "transactions"

type Transaction struct {
	ID                      string                   `db:"id"`
	TxHash                  string                   `db:"tx_hash"`
	TxHashDestination       string                   `db:"tx_hash_destination"`
	LogIndex                uint                     `db:"log_index"`
	SharedId                string                   `db:"shared_id"`
	State                   types.TransactionState   `db:"state"`
	Outcome                 types.TransactionOutcome `db:"outcome"`
	ProofInvalid            bool                     `db:"proof_invalid"`
	OriginatorChainId       string                   `db:"originator_chain_id"`
	DestinationChainId      string                   `db:"destination_chain_id"`
	MsgId                   []byte                   `db:"msg_id"`
	IsAtomic                bool                     `db:"is_atomic"`
	UpdatedAt               time.Time                `db:"updated_at"`
	CreatedAt               time.Time                `db:"created_at"`
	BatchId                 string                   `db:"batch_id"`
	BatchTxHashOnPrivateHub string                   `db:"batch_tx_hash_on_private_hub"`
	ResourceId              string                   `db:"resource_id"`
	FromContractAddress     string                   `db:"from_contract_address"`
	FromUserAddress         string                   `db:"from_user_address"`
	TransferMetadata_Id     string                   `db:"transfer_metadata_id"`
	TransferMetadata_Amount string                   `db:"transfer_metadata_amount"`
	BlockNumber             uint64                   `db:"block_number"`
	ParentHash              string                   `db:"parent_hash"`
}
