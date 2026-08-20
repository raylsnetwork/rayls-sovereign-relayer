package merkle

import (
	"context"
	"fmt"
	"math/big"
	"strings"

	"github.com/raylsnetwork/rayls-privacy-relayer-api/cryptography"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/types"

	// "github.com/iden3/go-iden3-crypto/poseidon"
	"golang.org/x/crypto/sha3"
)

// var SNARK_SCALAR_FIELD, _ = new(big.Int).SetString("21888242871839275222246405745257275088548364400416034343698204186575808495617", 10)
var SNARK_SCALAR_FIELD = cryptography.JubJubPrimeGroup

type merkleTreeRepository interface {
	GetByNumberAndTokenAddress(ctx context.Context, treeNumber int, tokenAddress string) (*types.MerkleTree, error)
	CreateMerkleTree(ctx context.Context, tree *types.MerkleTree) error
	InsertLeaves(ctx context.Context, tokenAddress string, treeNumber int, leaves []string) error
}

type MerkleService struct {
	depth          int
	treeRepository merkleTreeRepository
}

func NewMerkleService(depth int, treeRepository merkleTreeRepository) *MerkleService {
	return &MerkleService{
		depth:          depth,
		treeRepository: treeRepository,
	}
}

// MerkleProof holds the data for a Merkle proof.
type MerkleProof struct {
	Element  *big.Int   `json:"element"`
	Elements []*big.Int `json:"elements"`
	Indices  *big.Int   `json:"indices"`
	Root     *big.Int   `json:"root"`
}

type MerkleTree struct {
	types.MerkleTree
	Zeros []*big.Int   // per-level "zero" values
	Tree  [][]*big.Int // levels 0..Depth (with level Depth being the root level)
}

func initMerkleTree(depth int, number int) (*MerkleTree, error) {
	mt := &MerkleTree{}
	mt.Number = number
	mt.Depth = depth

	zeros, err := getZeroValueLevels(depth)
	if err != nil {
		return nil, fmt.Errorf("initializing zero levels: %w", err)
	}
	mt.Zeros = zeros

	// Pre-allocate empty slices for the higher levels.
	mt.Tree = make([][]*big.Int, depth+1)
	for i := 0; i < depth; i++ {
		mt.Tree[i] = make([]*big.Int, 0)
	}
	// The root level is initialized as the hash of the zero value of the level below.
	rootHash, err := hashLeftRight(mt.Zeros[depth-1], mt.Zeros[depth-1])
	if err != nil {
		return nil, fmt.Errorf("computing initial root: %w", err)
	}
	mt.Tree[depth] = []*big.Int{rootHash}

	return mt, nil
}

func NewMerkleTree(treeData *types.MerkleTree) (*MerkleTree, error) {
	mt, err := initMerkleTree(treeData.Depth, treeData.Number)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize merkle tree: %w", err)
	}

	// Convert stored leaves (as []string) into []*big.Int.
	leaves, err := stringSliceToBigIntSlice(treeData.Leaves)
	if err != nil {
		return nil, fmt.Errorf("failed to parse leaves: %w", err)
	}
	mt.Tree[0] = leaves

	// Recompute intermediate levels.
	if err := mt.RebuildSparseTree(); err != nil {
		return nil, fmt.Errorf("failed to rebuild sparse tree: %w", err)
	}

	return mt, nil
}

// RebuildSparseTree recomputes the tree from the leaves.
func (mt *MerkleTree) RebuildSparseTree() error {
	for level := 0; level < mt.Depth; level++ {
		// Clear next level.
		mt.Tree[level+1] = make([]*big.Int, 0)
		nodes := mt.Tree[level]
		for pos := 0; pos < len(nodes); pos += 2 {
			left := nodes[pos]
			var right *big.Int
			if pos+1 < len(nodes) {
				right = nodes[pos+1]
			} else {
				right = mt.Zeros[level]
			}
			parent, err := hashLeftRight(left, right)
			if err != nil {
				return fmt.Errorf("hashing level %d pos %d: %w", level, pos, err)
			}
			mt.Tree[level+1] = append(mt.Tree[level+1], parent)
		}
	}
	return nil
}

// GenerateProof returns a Merkle proof for the given element.
func (mt *MerkleTree) GenerateProof(element *big.Int) (*MerkleProof, error) {
	index := indexOf(mt.Tree[0], element)

	if index == -1 {
		return nil, fmt.Errorf("couldn't find %s in the MerkleTree number: %d", element.String(), mt.Number)
	}

	activeTree := mt.Tree
	activeRoot := mt.Root()
	proofElements := make([]*big.Int, 0)
	var indicesBits []string

	// Walk up the tree collecting sibling nodes.
	for level := 0; level < mt.Depth; level++ {
		if index%2 == 0 {
			var sibling *big.Int
			if index+1 < len(activeTree[level]) {
				sibling = activeTree[level][index+1]
			} else {
				sibling = mt.Zeros[level]
			}
			proofElements = append(proofElements, sibling)
			indicesBits = append(indicesBits, "0")
		} else {
			proofElements = append(proofElements, activeTree[level][index-1])
			indicesBits = append(indicesBits, "1")
		}
		index /= 2
	}

	// Reverse bits and parse as a binary number.
	reverseSlice(indicesBits)
	indicesStr := strings.Join(indicesBits, "")
	proofIndices := new(big.Int)
	proofIndices, ok := proofIndices.SetString(indicesStr, 2)
	if !ok {
		return nil, fmt.Errorf("failed to parse indices binary string")
	}
	return &MerkleProof{
		Element:  element,
		Elements: proofElements,
		Indices:  proofIndices,
		Root:     activeRoot,
	}, nil
}

// Root returns the current root of the tree.
func (mt *MerkleTree) Root() *big.Int {
	if len(mt.Tree) > mt.Depth && len(mt.Tree[mt.Depth]) > 0 {
		return mt.Tree[mt.Depth][0]
	}
	return nil
}

func hashLeftRight(left, right *big.Int) (*big.Int, error) {
	inputs := []*big.Int{left, right}
	result, err := cryptography.GetPoseidonHashModNumber(inputs, cryptography.JubJubPrimeGroup)
	if err != nil {
		return nil, fmt.Errorf("hashLeftRight: %w", err)
	}
	return result, nil
}

func getZeroValue() *big.Int {
	hasher := sha3.NewLegacyKeccak256()
	hasher.Write([]byte("Dvp"))
	digest := hasher.Sum(nil)
	zero := new(big.Int).SetBytes(digest)
	// Replace "fieldPrime" with the correct field modulus from the Poseidon library.
	return new(big.Int).Mod(zero, SNARK_SCALAR_FIELD)
}

// getZeroValueLevels returns a slice of per-level zero values.
func getZeroValueLevels(depth int) ([]*big.Int, error) {
	levels := make([]*big.Int, depth)
	levels[0] = getZeroValue()
	for i := 1; i < depth; i++ {
		h, err := hashLeftRight(levels[i-1], levels[i-1])
		if err != nil {
			return nil, fmt.Errorf("computing zero level %d: %w", i, err)
		}
		levels[i] = h
	}
	return levels, nil
}

// indexOf returns the index of element in slice, or -1 if not found.
func indexOf(slice []*big.Int, element *big.Int) int {
	for i, v := range slice {
		if v.Cmp(element) == 0 {
			return i
		}
	}
	return -1
}

// reverseSlice reverses a slice of strings in-place.
func reverseSlice(s []string) {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
}

// stringSliceToBigIntSlice converts a slice of strings to a slice of *big.Int.
func stringSliceToBigIntSlice(slice []string) ([]*big.Int, error) {
	result := make([]*big.Int, len(slice))
	for i, s := range slice {
		num, ok := new(big.Int).SetString(s, 10)
		if !ok {
			return nil, fmt.Errorf("failed to parse big.Int from string: %s", s)
		}
		result[i] = num
	}
	return result, nil
}

func bigIntSliceToStringSlice(slice []*big.Int) []string {
	result := make([]string, len(slice))
	for i, num := range slice {
		result[i] = num.String()
	}
	return result
}

func (s *MerkleService) PopulateMerkleDbTree(
	ctx context.Context,
	tokenAddress string,
	bigIntTokenType *big.Int,
	bigIntTreeNumber *big.Int,
	bigIntLeaves []*big.Int,
) error {
	//nolint:gosec // DvpTokenType is a small enum, within int range
	tokenType := types.DvpTokenType(bigIntTokenType.Uint64())
	treeNumber := int(bigIntTreeNumber.Int64())
	leaves := bigIntSliceToStringSlice(bigIntLeaves)

	mt, err := s.treeRepository.GetByNumberAndTokenAddress(ctx, treeNumber, tokenAddress)
	if err != nil {
		return fmt.Errorf("failed to get latest tree: %w", err)
	}

	// If mt is nil, create a new tree.
	if mt == nil {
		mt = &types.MerkleTree{
			Type:         tokenType,
			TokenAddress: tokenAddress,
			Number:       treeNumber,
			Depth:        s.depth,
			Leaves:       []string{},
		}
		mt.Leaves = append(mt.Leaves, leaves...)
		err = s.treeRepository.CreateMerkleTree(ctx, mt)
		if err != nil {
			return fmt.Errorf("failed to create new tree: %w", err)
		}
		return nil
	}

	// Save the updated tree in the repository.
	err = s.treeRepository.InsertLeaves(ctx, mt.TokenAddress, mt.Number, leaves)
	if err != nil {
		return fmt.Errorf("failed to insert leaves: %w", err)
	}
	return nil
}

func (s *MerkleService) GenerateMerkleProof(
	ctx context.Context,
	commitment *big.Int,
	treeNumber int,
	tokenAddress string,
) (*MerkleProof, error) {
	if tokenAddress == "" {
		return nil, fmt.Errorf("invalid tokenAddress: deposit has empty tokenAddress (likely not set in database)")
	}
	dbTree, err := s.treeRepository.GetByNumberAndTokenAddress(ctx, treeNumber, tokenAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to get merkle tree for treeNumber=%d, tokenAddress=%s: %w", treeNumber, tokenAddress, err)
	}
	if dbTree == nil {
		return nil, fmt.Errorf("merkle tree not found for treeNumber=%d, tokenAddress=%s", treeNumber, tokenAddress)
	}
	merkleTree, err := NewMerkleTree(dbTree)
	if err != nil {
		return nil, fmt.Errorf("failed to construct merkle tree for treeNumber=%d, tokenAddress=%s: %w", treeNumber, tokenAddress, err)
	}
	element := new(big.Int)
	element, ok := element.SetString(commitment.String(), 10)
	if !ok {
		return nil, fmt.Errorf("failed to convert commitment to big.Int: %s", commitment.String())
	}
	proof, err := merkleTree.GenerateProof(element)
	if err != nil {
		return nil, fmt.Errorf("failed to generate merkle proof for commitment=%s in treeNumber=%d: %w", commitment.String(), treeNumber, err)
	}
	return proof, nil
}
