package service

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/raylsnetwork/rayls-sovereign-relayer/conv"
	"github.com/raylsnetwork/rayls-sovereign-relayer/cryptography"
)

type CommitmentCalculator struct{}

func NewCommitmentCalculator() *CommitmentCalculator {
	return &CommitmentCalculator{}
}

// CalculateNFTCommitment computes H(H(spendPK, salt), uId) mod JubJubPrimeGroup,
// where uId = H(nftAddress, nftID) is produced by GetNFTUniqueID.
// Uses 2-input Poseidon hashing with mod reduction at each step to match the circuit.
func (s *CommitmentCalculator) CalculateNFTCommitment(
	spendPK, salt *big.Int, nftID string, nftAddress string,
) (*big.Int, error) {
	uId, err := s.GetNFTUniqueID(nftAddress, nftID)
	if err != nil {
		return nil, fmt.Errorf("nft commitment step 1: %w", err)
	}
	spendHash, err := cryptography.GetPoseidonHashModNumber([]*big.Int{spendPK, salt}, cryptography.JubJubPrimeGroup)
	if err != nil {
		return nil, fmt.Errorf("nft commitment step 2: %w", err)
	}
	result, err := cryptography.GetPoseidonHashModNumber([]*big.Int{spendHash, uId}, cryptography.JubJubPrimeGroup)
	if err != nil {
		return nil, fmt.Errorf("nft commitment step 3: %w", err)
	}
	return result, nil
}

// CalculatePaymentCommitment computes H(H(H(spendPK, salt), paymentAmount), tokenAddress) mod JubJubPrimeGroup.
// Uses chained 2-input Poseidon hashing with mod reduction at each step to match the circuit.
func (s *CommitmentCalculator) CalculatePaymentCommitment(
	spendPK, salt, paymentAmount *big.Int, tokenAddress string,
) (*big.Int, error) {
	tokenAddressBig := common.HexToAddress(tokenAddress).Big()

	step1, err := cryptography.GetPoseidonHashModNumber([]*big.Int{spendPK, salt}, cryptography.JubJubPrimeGroup)
	if err != nil {
		return nil, fmt.Errorf("payment commitment step 1: %w", err)
	}
	step2, err := cryptography.GetPoseidonHashModNumber([]*big.Int{step1, paymentAmount}, cryptography.JubJubPrimeGroup)
	if err != nil {
		return nil, fmt.Errorf("payment commitment step 2: %w", err)
	}
	result, err := cryptography.GetPoseidonHashModNumber([]*big.Int{step2, tokenAddressBig}, cryptography.JubJubPrimeGroup)
	if err != nil {
		return nil, fmt.Errorf("payment commitment step 3: %w", err)
	}
	return result, nil
}

// CalculateERC1155Commitment computes H(H(H(H(spendPK, salt), tokenAddress), tokenID), tokenAmount) mod JubJubPrimeGroup.
// Uses chained 2-input Poseidon hashing with mod reduction at each step to match the circuit.
func (s *CommitmentCalculator) CalculateERC1155Commitment(
	spendPK, salt *big.Int, tokenAddress string, tokenID string, tokenAmount *big.Int,
) (*big.Int, error) {
	tokenAddressBig := common.HexToAddress(tokenAddress).Big()
	tokenIDBig, ok := new(big.Int).SetString(tokenID, 10)
	if !ok {
		return nil, fmt.Errorf("parse erc1155 token id: %s", tokenID)
	}

	step1, err := cryptography.GetPoseidonHashModNumber([]*big.Int{spendPK, salt}, cryptography.JubJubPrimeGroup)
	if err != nil {
		return nil, fmt.Errorf("erc1155 commitment step 1: %w", err)
	}
	step2, err := cryptography.GetPoseidonHashModNumber([]*big.Int{step1, tokenAddressBig}, cryptography.JubJubPrimeGroup)
	if err != nil {
		return nil, fmt.Errorf("erc1155 commitment step 2: %w", err)
	}
	step3, err := cryptography.GetPoseidonHashModNumber([]*big.Int{step2, tokenIDBig}, cryptography.JubJubPrimeGroup)
	if err != nil {
		return nil, fmt.Errorf("erc1155 commitment step 3: %w", err)
	}
	result, err := cryptography.GetPoseidonHashModNumber([]*big.Int{step3, tokenAmount}, cryptography.JubJubPrimeGroup)
	if err != nil {
		return nil, fmt.Errorf("erc1155 commitment step 4: %w", err)
	}
	return result, nil
}

// CalculateNullifier computes H(spendSK, leafIndex) mod JubJubPrimeGroup.
func (s *CommitmentCalculator) CalculateNullifier(spendSK, leafIndex *big.Int) (*big.Int, error) {
	inputs := []*big.Int{spendSK, leafIndex}
	result, err := cryptography.GetPoseidonHashModNumber(inputs, cryptography.JubJubPrimeGroup)
	if err != nil {
		return nil, fmt.Errorf("calculate nullifier: %w", err)
	}
	return result, nil
}

func (s *CommitmentCalculator) GetNFTUniqueID(nftAddress string, nftID string) (*big.Int, error) {
	tokenAddress := common.HexToAddress(nftAddress)
	tokenID, _ := conv.StringToBigInt(nftID)

	// Generate unique ID using Poseidon hash of [assetContractAddress, nftId]
	nftUIDInputs := []*big.Int{tokenAddress.Big(), tokenID}
	nftUID, err := cryptography.GetPoseidonHashModNumber(nftUIDInputs, cryptography.JubJubPrimeGroup)
	if err != nil {
		return nil, fmt.Errorf("failed to generate NFT unique ID: %w", err)
	}

	return nftUID, nil
}
