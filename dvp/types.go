package dvp

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/raylsnetwork/rayls-sovereign-relayer/types"
)

type Proof struct {
	A            [2]*big.Int
	B            [2][2]*big.Int
	C            [2]*big.Int
	PublicSignal []*big.Int
}

type ProofReceipt struct {
	Proof            *Proof
	TreeNumbers      []*big.Int
	Message          *big.Int
	MerkleRoots      []*big.Int
	Commitments      []*big.Int
	Nullifiers       []*big.Int
	RevertCommitment *big.Int
}

type ERC1155JoinSplitProofRequest struct {
	NftCommitment  *big.Int
	ValuesIn       []*big.Int
	SaltsIn        []*big.Int
	KeyPairsIn     []*types.PaymentSpendKey
	ValuesOut      []*big.Int
	SaltsOut       []*big.Int
	PubKeysOut     []*big.Int
	MerkleDepth    int
	MerkleProofs   []*types.MerkleProof
	MerkleRoots    []*big.Int
	TreeNumbers    []int
	ERC1155Address string
	ERC1155TokenId *big.Int
	RevertSalt     *big.Int
}

// EnygmaJSProofInput holds the inputs for an Enygma join-split proof that
// spends Deposits into a PaymentAmount output for DestinationPaymentPublicKey
// (salted with DestinationSalt), returning any change to the sender.
type EnygmaJSProofInput struct {
	SourceViewPublicKey         []byte
	NftCommitment               *big.Int
	DestinationPaymentPublicKey *big.Int
	DestinationSalt             *big.Int
	PaymentAmount               *big.Int
	TokenAddress                string
	Deposits                    []*types.DvpDeposit
}

type EnygmaJoinSplitProofRequest struct {
	NftCommitment *big.Int
	ValuesIn      []*big.Int
	SaltsIn       []*big.Int
	KeyPairsIn    []*types.PaymentSpendKey
	ValuesOut     []*big.Int
	SaltsOut      []*big.Int
	PubKeysOut    []*big.Int
	MerkleDepth   int
	MerkleProofs  []*types.MerkleProof
	MerkleRoots   []*big.Int
	TreeNumbers   []int
	ERC20Address  string
	RevertSalt    *big.Int
}

type ERC721OwnershipProofRequest struct {
	PaymentCommitment *big.Int
	UID               *big.Int
	KeyPairIn         *types.PaymentSpendKey
	SaltIn            *big.Int
	PubKeyOut         *big.Int
	SaltOut           *big.Int
	MerkleDepth       int
	MerkleProof       *types.MerkleProof
	MerkleRoot        *big.Int
	TreeNumber        int
	RevertSalt        *big.Int
}

type DvpERC721MintData struct {
	To           common.Address
	TokenAddress common.Address
	TokenID      *big.Int
	ChainID      *big.Int
	ExtraData    []byte
}

type DvpExtraData struct {
	Key      string
	Value    string
	IsPublic bool
}

// Used by tx generator
func (d *DvpERC721MintData) GetID() string {
	return d.TokenID.String()
}

type DvpERC1155MintData struct {
	To           common.Address
	TokenAddress common.Address
	TokenID      *big.Int
	TokenAmount  *big.Int
	ChainID      *big.Int
	Data         []byte
	ExtraData    []byte
}

// Used by tx generator
func (d *DvpERC1155MintData) GetID() string {
	return d.TokenID.String()
}

type DvpERC1155Supply struct {
	Id     *big.Int
	Amount *big.Int
}
