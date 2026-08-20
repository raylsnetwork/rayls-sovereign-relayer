package client

import (
	"math/big"

	"github.com/raylsnetwork/rayls-sovereign-relayer/types"
)

type ProofAPITransferRequest struct {
	SenderID                  string     `json:"sender_id"`
	SecretKey                 string     `json:"secret_key"`
	SenderTxValue             string     `json:"sender_tx_value"`
	PublicKeys                []string   `json:"public_keys"`
	PreviousCommits           [][]string `json:"previous_commits"`
	PreviousSenderBalance     string     `json:"previous_sender_balance"`
	PreviousSenderRandomValue string     `json:"previous_sender_random_value"`
	TxCommits                 [][]string `json:"tx_commits"`
	TxValues                  []string   `json:"tx_values"`
	TxRandomValues            []string   `json:"tx_random_values"`
	SharedSecrets             []string   `json:"shared_secrets"`
	HashedSharedSecrets       []string   `json:"hashed_shared_secrets"`
	MessageTags               []string   `json:"message_tags"`
	BlockNumber               string     `json:"block_number"`
	Nullifier                 string     `json:"nullifier"`
	AnonymitySet              []string   `json:"anonymity_set"`
}

type ProofAPIWithdrawRequest struct {
	SenderID                  string     `json:"sender_id"`
	SecretKey                 string     `json:"secret_key"`
	SenderTxValue             string     `json:"sender_tx_value"`
	PublicKeys                []string   `json:"public_keys"`
	PreviousCommits           [][]string `json:"previous_commits"`
	PreviousSenderBalance     string     `json:"previous_sender_balance"`
	PreviousSenderRandomValue string     `json:"previous_sender_random_value"`
	TxCommits                 [][]string `json:"tx_commits"`
	TxValues                  []string   `json:"tx_values"`
	TxRandomValues            []string   `json:"tx_random_values"`
	SharedSecrets             []string   `json:"shared_secrets"`
	HashedSharedSecrets       []string   `json:"hashed_shared_secrets"`
	MessageTags               []string   `json:"message_tags"`
	BlockNumber               string     `json:"block_number"`
	Nullifier                 string     `json:"nullifier"`
	AnonymitySet              []string   `json:"anonymity_set"`
	Hashes                    []string   `json:"hashes"`
	SkDeposits                []string   `json:"sk_deposits"`
	VPerDeposit               []string   `json:"v_per_deposit"`
	Address                   string     `json:"address"`
	SaltsIn                   []string   `json:"saltsIn"`
}

type ProofAPIDepositRequest struct {
	SenderID                  string     `json:"sender_id"`
	SecretKey                 string     `json:"secret_key"`
	SenderTxValue             string     `json:"sender_tx_value"`
	PublicKeys                []string   `json:"public_keys"`
	PreviousCommits           [][]string `json:"previous_commits"`
	PreviousSenderBalance     string     `json:"previous_sender_balance"`
	PreviousSenderRandomValue string     `json:"previous_sender_random_value"`
	TxCommits                 [][]string `json:"tx_commits"`
	TxValues                  []string   `json:"tx_values"`
	TxRandomValues            []string   `json:"tx_random_values"`
	SharedSecrets             []string   `json:"shared_secrets"`
	HashedSharedSecrets       []string   `json:"hashed_shared_secrets"`
	MessageTags               []string   `json:"message_tags"`
	BlockNumber               string     `json:"block_number"`
	Nullifier                 string     `json:"nullifier"`
	AnonymitySet              []string   `json:"anonymity_set"`
	Hash                      string     `json:"hash"`
	Pk                        string     `json:"pk"`
	Address                   string     `json:"address"`
	SaltOut                   string     `json:"saltOut"`
}

// DVP
type ProofAPIDvpSpendKeys struct {
	PublicKey  string `json:"publicKey"`
	PrivateKey string `json:"privateKey"`
}

type ProofAPIDvpMerkleProof struct {
	Indices  string   `json:"indices"`
	Elements []string `json:"elements"`
}
type ProofAPIDvpERC1155JoinSplitRequestParams struct {
	NFTCommitment  string                   `json:"nftCommitment"`
	ValuesIn       []string                 `json:"valuesIn"`
	SaltsIn        []string                 `json:"saltsIn"`
	ValuesOut      []string                 `json:"valuesOut"`
	SaltsOut       []string                 `json:"saltsOut"`
	KeyPairsIn     []ProofAPIDvpSpendKeys   `json:"keyPairsIn"`
	PubKeysOut     []ProofAPIDvpSpendKeys   `json:"pubKeysOut"`
	MerkleDepth    int                      `json:"merkleDepth"`
	MerkleProofs   []ProofAPIDvpMerkleProof `json:"merkleProofs"`
	MerkleRoots    []string                 `json:"merkleRoots"`
	TreeNumbers    []int                    `json:"treeNumbers"`
	ERC1155Address string                   `json:"erc1155Address"`
	ERC1155TokenId string                   `json:"erc1155TokenId"`
	RevertSalt string `json:"revertSalt"`
}

type ProofAPIEnygmaJoinSplitRequestParams struct {
	NFTCommitment string                   `json:"nftCommitment"`
	ValuesIn      []string                 `json:"valuesIn"`
	SaltsIn       []string                 `json:"saltsIn"`
	ValuesOut     []string                 `json:"valuesOut"`
	SaltsOut      []string                 `json:"saltsOut"`
	KeyPairsIn    []ProofAPIDvpSpendKeys   `json:"keyPairsIn"`
	PubKeysOut    []ProofAPIDvpSpendKeys   `json:"pubKeysOut"`
	MerkleDepth   int                      `json:"merkleDepth"`
	MerkleProofs  []ProofAPIDvpMerkleProof `json:"merkleProofs"`
	MerkleRoots   []string                 `json:"merkleRoots"`
	TreeNumbers   []int                    `json:"treeNumbers"`
	ERC20Address  string                   `json:"erc20Address"`
	RevertSalt string `json:"revertSalt"`
}

type ProofAPIErc721OwnershipRequestParams struct {
	PaymentCommitment string                 `json:"paymentCommitment"`
	UID               string                 `json:"uid"`
	KeyPairIn         ProofAPIDvpSpendKeys   `json:"keyPairIn"`
	SaltIn            string                 `json:"saltIn"`
	PubKeyOut         ProofAPIDvpSpendKeys   `json:"pubKeyOut"`
	SaltOut           string                 `json:"saltOut"`
	MerkleDepth       int                    `json:"merkleDepth"`
	MerkleProof       ProofAPIDvpMerkleProof `json:"merkleProof"`
	MerkleRoot        string                 `json:"merkleRoot"`
	TreeNumber        int                    `json:"treeNumber"`
	RevertSalt string `json:"revertSalt"`
}

type ProofAPIResponse struct {
	Pi_A          []string   `json:"pi_a"`
	Pi_B          [][]string `json:"pi_b"`
	Pi_C          []string   `json:"pi_c"`
	Public_Signal []string   `json:"public_signal"`
}

// convertPointsToStringArrays converts Points to string arrays for API requests
func convertPointsToStringArrays(points []*types.Point) [][]string {
	var results [][]string

	for _, point := range points {
		if point == nil || point.X == nil || point.Y == nil {
			results = append(results, []string{"0", "0"})
			continue
		}
		result := make([]string, 0, 2)
		result = append(result, point.X.String())
		result = append(result, point.Y.String())
		results = append(results, result)
	}

	return results
}

// convertBigIntArrayToStringArray converts big.Int array to string array for API requests
func convertBigIntArrayToStringArray(bigInts []*big.Int) []string {
	results := make([]string, len(bigInts))
	for i, bigInt := range bigInts {
		results[i] = bigInt.String()
	}
	return results
}

func convertPaymentSpendKeysArrayToStringArray(keyPairs []*types.PaymentSpendKey) []ProofAPIDvpSpendKeys {
	results := make([]ProofAPIDvpSpendKeys, len(keyPairs))
	for i, keyPair := range keyPairs {
		results[i] = ProofAPIDvpSpendKeys{
			PublicKey:  keyPair.PublicKey.String(),
			PrivateKey: keyPair.SecretKey.String(),
		}
	}
	return results
}

func convertMerkleProofArrayToStringArray(merkleProofs []*types.MerkleProof) []ProofAPIDvpMerkleProof {
	results := make([]ProofAPIDvpMerkleProof, len(merkleProofs))
	for i, merkleProof := range merkleProofs {
		results[i] = ProofAPIDvpMerkleProof{
			Indices:  merkleProof.Indices.String(),
			Elements: convertBigIntArrayToStringArray(merkleProof.Elements),
		}
	}
	return results
}
