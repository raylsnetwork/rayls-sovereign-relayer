package enygma

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/raylsnetwork/rayls-sovereign-relayer/types"
)

// TransferProofParams contains parameters specific to transfer proof generation
type TransferProofParams struct {
	ResourceId     string
	AnonymityIndex int
	SenderAmount   *big.Int
	BlockNumber    *big.Int
	Batches        []*types.EnygmaTransferBatch
	TokenAddress   common.Address
}

// DepositProofParams contains parameters specific to deposit proof generation
type DepositProofParams struct {
	ResourceId        string
	AnonymityIndex    int
	SenderAmount      *big.Int
	BlockNumber       *big.Int
	Batches           []*types.EnygmaTransferBatch
	TokenAddress      common.Address
	DepositCommitment *big.Int
	DepositSalt       *big.Int
	DepositPublicKey  *big.Int
}

// WithdrawProofParams contains parameters specific to withdraw proof generation
type WithdrawProofParams struct {
	ResourceId         string
	AnonymityIndex     int
	PublicValues       *types.EnygmaPublicValues
	SenderAmount       *big.Int
	BlockNumber        *big.Int
	Batches            []*types.EnygmaTransferBatch
	TokenAddress       common.Address
	DepositCommitments []*big.Int
	DepositSecretKeys  []*big.Int
	DepositAmounts     []*big.Int
	DepositSalts       []*big.Int
}
