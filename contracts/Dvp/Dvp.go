// Code generated via abigen V2 - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package Dvp

import (
	"bytes"
	"errors"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/v2"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = bytes.Equal
	_ = errors.New
	_ = big.NewInt
	_ = common.Big1
	_ = types.BloomLookup
	_ = abi.ConvertType
)

// IDvpG1Point is an auto generated low-level Go binding around an user-defined struct.
type IDvpG1Point struct {
	X *big.Int
	Y *big.Int
}

// IDvpG2Point is an auto generated low-level Go binding around an user-defined struct.
type IDvpG2Point struct {
	X [2]*big.Int
	Y [2]*big.Int
}

// IDvpProofReceipt is an auto generated low-level Go binding around an user-defined struct.
type IDvpProofReceipt struct {
	Proof            IDvpSnarkProof
	TreeNumbers      []*big.Int
	Message          *big.Int
	MerkleRoots      []*big.Int
	Commitments      []*big.Int
	Nullifiers       []*big.Int
	RevertCommitment *big.Int
}

// IDvpSnarkProof is an auto generated low-level Go binding around an user-defined struct.
type IDvpSnarkProof struct {
	A IDvpG1Point
	B IDvpG2Point
	C IDvpG1Point
}

// DvpMetaData contains all meta data concerning the Dvp contract.
var DvpMetaData = bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"hashPoseidonContractAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"enygmaFactoryAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"dvpErc721FactoryAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"dvpErc1155FactoryAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"dvpTeleportAddr\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"authority_\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"VK_ID_AUCTION_BID\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"VK_ID_AUCTION_INIT\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"VK_ID_AUCTION_NOT_WINNING_BID\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"VK_ID_AUCTION_PRIVATE_OPENING\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"VK_ID_BROKER_REGISTRATION\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"VK_ID_LEGIT_BROKER\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"addEnygmaDvpIntegrationAddress\",\"inputs\":[{\"name\":\"enygmaDvpIntegrationAddress\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"addTokenToGroup\",\"inputs\":[{\"name\":\"vaultId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"uniqueIdParams\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"},{\"name\":\"groupId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"addVaultToGroup\",\"inputs\":[{\"name\":\"vaultId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"groupId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"authority\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"cancelSwap\",\"inputs\":[{\"name\":\"sharedId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"preimage\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"checkAndRegisterChallenge\",\"inputs\":[{\"name\":\"challenge_\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"completeSwap\",\"inputs\":[{\"name\":\"sharedId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"tokenAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"proofType\",\"type\":\"uint8\",\"internalType\":\"enumIDvp.SwapProofType\"},{\"name\":\"proof\",\"type\":\"tuple\",\"internalType\":\"structIDvp.ProofReceipt\",\"components\":[{\"name\":\"proof\",\"type\":\"tuple\",\"internalType\":\"structIDvp.SnarkProof\",\"components\":[{\"name\":\"a\",\"type\":\"tuple\",\"internalType\":\"structIDvp.G1Point\",\"components\":[{\"name\":\"x\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"b\",\"type\":\"tuple\",\"internalType\":\"structIDvp.G2Point\",\"components\":[{\"name\":\"x\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"},{\"name\":\"y\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"}]},{\"name\":\"c\",\"type\":\"tuple\",\"internalType\":\"structIDvp.G1Point\",\"components\":[{\"name\":\"x\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]}]},{\"name\":\"treeNumbers\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"},{\"name\":\"message\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"merkleRoots\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"},{\"name\":\"commitments\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"},{\"name\":\"nullifiers\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"},{\"name\":\"revertCommitment\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"encryptedData\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"depositERC1155\",\"inputs\":[{\"name\":\"contractAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"amountOrOne\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"data\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"publicKey\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"salt\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"encryptedBalanceUpdate\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"depositERC721\",\"inputs\":[{\"name\":\"contractAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"nftId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"publicKey\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"salt\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"encryptedBalanceUpdate\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"depositEnygma\",\"inputs\":[{\"name\":\"vaultId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"hashCommitment\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"dvpTeleportAddress\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"expireSwap\",\"inputs\":[{\"name\":\"sharedId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"getGroupPairId\",\"inputs\":[{\"name\":\"groupId1\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"groupId2\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"getVaultIdByAddress\",\"inputs\":[{\"name\":\"contractAddress\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"hashContractAddress\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"initializeDvp\",\"inputs\":[{\"name\":\"verifierAddress\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"initiateSwap\",\"inputs\":[{\"name\":\"sharedId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"encryptedData\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"ciphertext\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"tokenAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"proofType\",\"type\":\"uint8\",\"internalType\":\"enumIDvp.SwapProofType\"},{\"name\":\"proof\",\"type\":\"tuple\",\"internalType\":\"structIDvp.ProofReceipt\",\"components\":[{\"name\":\"proof\",\"type\":\"tuple\",\"internalType\":\"structIDvp.SnarkProof\",\"components\":[{\"name\":\"a\",\"type\":\"tuple\",\"internalType\":\"structIDvp.G1Point\",\"components\":[{\"name\":\"x\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"b\",\"type\":\"tuple\",\"internalType\":\"structIDvp.G2Point\",\"components\":[{\"name\":\"x\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"},{\"name\":\"y\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"}]},{\"name\":\"c\",\"type\":\"tuple\",\"internalType\":\"structIDvp.G1Point\",\"components\":[{\"name\":\"x\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]}]},{\"name\":\"treeNumbers\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"},{\"name\":\"message\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"merkleRoots\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"},{\"name\":\"commitments\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"},{\"name\":\"nullifiers\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"},{\"name\":\"revertCommitment\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"validityTime\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"passphrase\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"dvpId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"isAuditorRegistered\",\"inputs\":[{\"name\":\"publicKey\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"isMemberOfFromProofReceipt\",\"inputs\":[{\"name\":\"vaultId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"receipt\",\"type\":\"tuple\",\"internalType\":\"structIDvp.ProofReceipt\",\"components\":[{\"name\":\"proof\",\"type\":\"tuple\",\"internalType\":\"structIDvp.SnarkProof\",\"components\":[{\"name\":\"a\",\"type\":\"tuple\",\"internalType\":\"structIDvp.G1Point\",\"components\":[{\"name\":\"x\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"b\",\"type\":\"tuple\",\"internalType\":\"structIDvp.G2Point\",\"components\":[{\"name\":\"x\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"},{\"name\":\"y\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"}]},{\"name\":\"c\",\"type\":\"tuple\",\"internalType\":\"structIDvp.G1Point\",\"components\":[{\"name\":\"x\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]}]},{\"name\":\"treeNumbers\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"},{\"name\":\"message\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"merkleRoots\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"},{\"name\":\"commitments\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"},{\"name\":\"nullifiers\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"},{\"name\":\"revertCommitment\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"groupId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"isSwapExpired\",\"inputs\":[{\"name\":\"sharedId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"isTokenMemberOf\",\"inputs\":[{\"name\":\"vaultId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"uniqueIdParams\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"},{\"name\":\"groupId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"isValidExchangeGroupPair\",\"inputs\":[{\"name\":\"groupId1\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"groupId2\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"isValidSwapGroupPair\",\"inputs\":[{\"name\":\"groupId1\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"groupId2\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"isVaultMemberOf\",\"inputs\":[{\"name\":\"vaultId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"groupId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"mixFunds\",\"inputs\":[{\"name\":\"vaultId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"_tx\",\"type\":\"tuple\",\"internalType\":\"structIDvp.ProofReceipt\",\"components\":[{\"name\":\"proof\",\"type\":\"tuple\",\"internalType\":\"structIDvp.SnarkProof\",\"components\":[{\"name\":\"a\",\"type\":\"tuple\",\"internalType\":\"structIDvp.G1Point\",\"components\":[{\"name\":\"x\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"b\",\"type\":\"tuple\",\"internalType\":\"structIDvp.G2Point\",\"components\":[{\"name\":\"x\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"},{\"name\":\"y\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"}]},{\"name\":\"c\",\"type\":\"tuple\",\"internalType\":\"structIDvp.G1Point\",\"components\":[{\"name\":\"x\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]}]},{\"name\":\"treeNumbers\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"},{\"name\":\"message\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"merkleRoots\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"},{\"name\":\"commitments\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"},{\"name\":\"nullifiers\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"},{\"name\":\"revertCommitment\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"mixFundsERC1155\",\"inputs\":[{\"name\":\"contractAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_tx\",\"type\":\"tuple\",\"internalType\":\"structIDvp.ProofReceipt\",\"components\":[{\"name\":\"proof\",\"type\":\"tuple\",\"internalType\":\"structIDvp.SnarkProof\",\"components\":[{\"name\":\"a\",\"type\":\"tuple\",\"internalType\":\"structIDvp.G1Point\",\"components\":[{\"name\":\"x\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"b\",\"type\":\"tuple\",\"internalType\":\"structIDvp.G2Point\",\"components\":[{\"name\":\"x\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"},{\"name\":\"y\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"}]},{\"name\":\"c\",\"type\":\"tuple\",\"internalType\":\"structIDvp.G1Point\",\"components\":[{\"name\":\"x\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]}]},{\"name\":\"treeNumbers\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"},{\"name\":\"message\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"merkleRoots\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"},{\"name\":\"commitments\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"},{\"name\":\"nullifiers\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"},{\"name\":\"revertCommitment\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"name\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"registerAssetGroup\",\"inputs\":[{\"name\":\"assetGroupContractAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"assetGroupName\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"isAssetGroupFungible\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"treeDepth\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"registerAuditor\",\"inputs\":[{\"name\":\"auditorOffchainId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"auditorGroupId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"auditorPublicKey\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"registerBroker\",\"inputs\":[{\"name\":\"brokerRegistrationProof\",\"type\":\"tuple\",\"internalType\":\"structIDvp.ProofReceipt\",\"components\":[{\"name\":\"proof\",\"type\":\"tuple\",\"internalType\":\"structIDvp.SnarkProof\",\"components\":[{\"name\":\"a\",\"type\":\"tuple\",\"internalType\":\"structIDvp.G1Point\",\"components\":[{\"name\":\"x\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"b\",\"type\":\"tuple\",\"internalType\":\"structIDvp.G2Point\",\"components\":[{\"name\":\"x\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"},{\"name\":\"y\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"}]},{\"name\":\"c\",\"type\":\"tuple\",\"internalType\":\"structIDvp.G1Point\",\"components\":[{\"name\":\"x\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]}]},{\"name\":\"treeNumbers\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"},{\"name\":\"message\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"merkleRoots\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"},{\"name\":\"commitments\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"},{\"name\":\"nullifiers\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"},{\"name\":\"revertCommitment\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"registerExchangeGroupPair\",\"inputs\":[{\"name\":\"groupId1\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"groupId2\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"registerSwapGroupPair\",\"inputs\":[{\"name\":\"groupId1\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"groupId2\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"registerVault\",\"inputs\":[{\"name\":\"vaultContractAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"assetContractAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"treeDepth\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"registerZkAuction\",\"inputs\":[{\"name\":\"zkAuctionContractAddress\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"unregisterAuditor\",\"inputs\":[{\"name\":\"auditorOnchainId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"vaultById\",\"inputs\":[{\"name\":\"vaultId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"verifierContractAddress\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"verifyLegitBrokerReceipt\",\"inputs\":[{\"name\":\"receipt\",\"type\":\"tuple\",\"internalType\":\"structIDvp.ProofReceipt\",\"components\":[{\"name\":\"proof\",\"type\":\"tuple\",\"internalType\":\"structIDvp.SnarkProof\",\"components\":[{\"name\":\"a\",\"type\":\"tuple\",\"internalType\":\"structIDvp.G1Point\",\"components\":[{\"name\":\"x\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"b\",\"type\":\"tuple\",\"internalType\":\"structIDvp.G2Point\",\"components\":[{\"name\":\"x\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"},{\"name\":\"y\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"}]},{\"name\":\"c\",\"type\":\"tuple\",\"internalType\":\"structIDvp.G1Point\",\"components\":[{\"name\":\"x\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]}]},{\"name\":\"treeNumbers\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"},{\"name\":\"message\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"merkleRoots\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"},{\"name\":\"commitments\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"},{\"name\":\"nullifiers\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"},{\"name\":\"revertCommitment\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"withdrawERC1155\",\"inputs\":[{\"name\":\"contractAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"recipient\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"salt\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"proofTx\",\"type\":\"tuple\",\"internalType\":\"structIDvp.ProofReceipt\",\"components\":[{\"name\":\"proof\",\"type\":\"tuple\",\"internalType\":\"structIDvp.SnarkProof\",\"components\":[{\"name\":\"a\",\"type\":\"tuple\",\"internalType\":\"structIDvp.G1Point\",\"components\":[{\"name\":\"x\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"b\",\"type\":\"tuple\",\"internalType\":\"structIDvp.G2Point\",\"components\":[{\"name\":\"x\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"},{\"name\":\"y\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"}]},{\"name\":\"c\",\"type\":\"tuple\",\"internalType\":\"structIDvp.G1Point\",\"components\":[{\"name\":\"x\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]}]},{\"name\":\"treeNumbers\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"},{\"name\":\"message\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"merkleRoots\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"},{\"name\":\"commitments\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"},{\"name\":\"nullifiers\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"},{\"name\":\"revertCommitment\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"encryptedBalanceUpdate\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"withdrawERC721\",\"inputs\":[{\"name\":\"contractAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"nftId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"recipient\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"salt\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"proofTx\",\"type\":\"tuple\",\"internalType\":\"structIDvp.ProofReceipt\",\"components\":[{\"name\":\"proof\",\"type\":\"tuple\",\"internalType\":\"structIDvp.SnarkProof\",\"components\":[{\"name\":\"a\",\"type\":\"tuple\",\"internalType\":\"structIDvp.G1Point\",\"components\":[{\"name\":\"x\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"b\",\"type\":\"tuple\",\"internalType\":\"structIDvp.G2Point\",\"components\":[{\"name\":\"x\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"},{\"name\":\"y\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"}]},{\"name\":\"c\",\"type\":\"tuple\",\"internalType\":\"structIDvp.G1Point\",\"components\":[{\"name\":\"x\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]}]},{\"name\":\"treeNumbers\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"},{\"name\":\"message\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"merkleRoots\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"},{\"name\":\"commitments\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"},{\"name\":\"nullifiers\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"},{\"name\":\"revertCommitment\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"encryptedBalanceUpdate\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"withdrawEnygma\",\"inputs\":[{\"name\":\"vaultId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"_tx\",\"type\":\"tuple\",\"internalType\":\"structIDvp.ProofReceipt\",\"components\":[{\"name\":\"proof\",\"type\":\"tuple\",\"internalType\":\"structIDvp.SnarkProof\",\"components\":[{\"name\":\"a\",\"type\":\"tuple\",\"internalType\":\"structIDvp.G1Point\",\"components\":[{\"name\":\"x\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"b\",\"type\":\"tuple\",\"internalType\":\"structIDvp.G2Point\",\"components\":[{\"name\":\"x\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"},{\"name\":\"y\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"}]},{\"name\":\"c\",\"type\":\"tuple\",\"internalType\":\"structIDvp.G1Point\",\"components\":[{\"name\":\"x\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]}]},{\"name\":\"treeNumbers\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"},{\"name\":\"message\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"merkleRoots\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"},{\"name\":\"commitments\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"},{\"name\":\"nullifiers\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"},{\"name\":\"revertCommitment\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"event\",\"name\":\"AuditorRegistered\",\"inputs\":[{\"name\":\"onchainId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"offchainId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"groupId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"publicKey\",\"type\":\"uint256[2]\",\"indexed\":false,\"internalType\":\"uint256[2]\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"AuditorUnregistered\",\"inputs\":[{\"name\":\"onchainId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"AuthorityUpdated\",\"inputs\":[{\"name\":\"oldAuthority\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newAuthority\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"BrokerRegistered\",\"inputs\":[{\"name\":\"vaultId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"blindedBrokerPublicKey\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"CoinLocked\",\"inputs\":[{\"name\":\"assetId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"treeNumber\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"nullifier\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"CoinUnlocked\",\"inputs\":[{\"name\":\"assetId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"treeNumber\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"nullifier\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"LegitBrokerReceipt\",\"inputs\":[{\"name\":\"beacon\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"blindedBrokerPublicKey\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"PendingProofAddedToVault\",\"inputs\":[{\"name\":\"vaultId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"groupId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"targetReceiptId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"pendingProof\",\"type\":\"tuple\",\"indexed\":false,\"internalType\":\"structIDvp.ProofReceipt\",\"components\":[{\"name\":\"proof\",\"type\":\"tuple\",\"internalType\":\"structIDvp.SnarkProof\",\"components\":[{\"name\":\"a\",\"type\":\"tuple\",\"internalType\":\"structIDvp.G1Point\",\"components\":[{\"name\":\"x\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"b\",\"type\":\"tuple\",\"internalType\":\"structIDvp.G2Point\",\"components\":[{\"name\":\"x\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"},{\"name\":\"y\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"}]},{\"name\":\"c\",\"type\":\"tuple\",\"internalType\":\"structIDvp.G1Point\",\"components\":[{\"name\":\"x\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]}]},{\"name\":\"treeNumbers\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"},{\"name\":\"message\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"merkleRoots\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"},{\"name\":\"commitments\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"},{\"name\":\"nullifiers\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"},{\"name\":\"revertCommitment\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Settled\",\"inputs\":[{\"name\":\"receiptId1\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"receiptId2\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"TokenAddedToGroup\",\"inputs\":[{\"name\":\"vaultId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"tokenUniqueId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"groupId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"VaultAddedToGroup\",\"inputs\":[{\"name\":\"vaultId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"groupId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"VerifyOwnershipReceipt\",\"inputs\":[{\"name\":\"challenge\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"assetId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"tokenId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"amountOrOne\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"assetAddress\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"senderAddress\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"AuctionAlreadyExists\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"AuctionIdMismatch\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"AuctionStateMismatch\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"AuditorNotRegistered\",\"inputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"BidStateMismatch\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"BlindedBidMismatch\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"Dvp__AuditorAlreadyRegistered\",\"inputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"Dvp__BrokerAlreadyRegistered\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"Dvp__GroupFungibilityMismatch\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"Dvp__GroupIdOutOfRange\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"Dvp__GroupPairAlreadyRegistered\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"Dvp__InvalidDeliveryMessage\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"Dvp__InvalidPassphrase\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"Dvp__InvalidPaymentMessage\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"Dvp__InvalidRevertCommitment\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"Dvp__InvalidStatementSize\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"Dvp__InvalidSwapGroupPair\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"Dvp__SwapAlreadyExists\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"Dvp__SwapNotExpired\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"Dvp__SwapNotFound\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"Dvp__SwapNotPending\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"FungibleDeliveryVault\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"GroupMembershipMismatch\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidChallenge\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidExchangeGroupPair\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidMerkleRoot\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidNullifier\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidNumberOfInputs\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidNumberOfOutputs\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidOpening\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidPartialProofReceipt\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"JoinSplitWithSameCommitments\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NonFungiblePaymentVault\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NotImplemented\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NotWinningBidsCountMismatch\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"RaylsAccessManaged__ContractPaused\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"RaylsAccessManaged__InvalidAuthority\",\"inputs\":[{\"name\":\"authority\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"RaylsAccessManaged__MustSchedule\",\"inputs\":[{\"name\":\"caller\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"delay\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"type\":\"error\",\"name\":\"RaylsAccessManaged__Unauthorized\",\"inputs\":[{\"name\":\"caller\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ReentrancyGuardReentrantCall\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"RottenChallenge\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"WinningBidOpeningMismatch\",\"inputs\":[]}]",
	ID:  "Dvp",
	Bin: "0x60806040523480156200001157600080fd5b5060405162005efb38038062005efb8339810160408190526200003491620007ac565b60016000819055604080518082019091526012815271111594081cdb585c9d0818dbdb9d1c9858dd60721b6020820152620000709082620008d4565b50600280546001600160a01b038089166001600160a01b0319928316179092556003805492851692909116919091179055620000ac816200068a565b60408051600c8082526101a0820190925260009160208201610180803683370190505090506385e1969b60e01b81600081518110620000ef57620000ef620009a0565b6001600160e01b0319909216602092830291909101909101528051630c7b50db60e41b9082906001908110620001295762000129620009a0565b6001600160e01b0319909216602092830291909101909101528051633e25ed8760e11b9082906002908110620001635762000163620009a0565b6001600160e01b031990921660209283029190910190910152805163a30886bf60e01b90829060039081106200019d576200019d620009a0565b6001600160e01b03199092166020928302919091019091015280516394f1bd0760e01b9082906004908110620001d757620001d7620009a0565b6001600160e01b03199092166020928302919091019091015280516304473f2760e01b9082906005908110620002115762000211620009a0565b6001600160e01b031990921660209283029190910190910152805163c927402160e01b90829060069081106200024b576200024b620009a0565b6001600160e01b031990921660209283029190910190910152805163977701cf60e01b9082906007908110620002855762000285620009a0565b6001600160e01b0319909216602092830291909101909101528051635c5924e360e11b9082906008908110620002bf57620002bf620009a0565b6001600160e01b03199092166020928302919091019091015280516331ff026160e21b9082906009908110620002f957620002f9620009a0565b6001600160e01b03199092166020928302919091019091015280516315b0669f60e01b908290600a908110620003335762000333620009a0565b6001600160e01b03199092166020928302919091019091015280516326a5cc4360e11b908290600b9081106200036d576200036d620009a0565b6001600160e01b031992909216602092830291909101820152604080516003808252608082019092526000929091908201606080368337019050509050631a95319160e01b81600081518110620003c857620003c8620009a0565b6001600160e01b0319909216602092830291909101909101528051638a622d4960e01b9082906001908110620004025762000402620009a0565b6001600160e01b0319909216602092830291909101909101528051637997ad5b60e01b90829060029081106200043c576200043c620009a0565b6001600160e01b031992909216602092830291909101820152604080516003808252608082019092526000929091908201606080368337019050509050633cd6f55e60e01b81600081518110620004975762000497620009a0565b6001600160e01b03199092166020928302919091019091015280516315b0669f60e01b9082906001908110620004d157620004d1620009a0565b6001600160e01b03199092166020928302919091019091015280516326a5cc4360e11b90829060029081106200050b576200050b620009a0565b6001600160e01b0319929092166020928302919091019091015260408051600280825260608201909252600091816020015b60408051808201909152606080825260208201528152602001906001900390816200053d57505060408051608081018252600e9181019182526d22a72ca3a6a0afa1a922a0aa27a960911b6060820152908152602081018590528151919250908290600090620005b157620005b1620009a0565b602002602001018190525060405180604001604052806040518060400160405280600a81526020016910d3d25397d59055531560b21b81525081526020018381525081600181518110620006095762000609620009a0565b60209081029190910101526040516337af400760e11b81526001600160a01b03861690636f5e800e906200064690309088908690600401620009fe565b600060405180830381600087803b1580156200066157600080fd5b505af115801562000676573d6000803e3d6000fd5b505050505050505050505050505062000b0d565b6001600160a01b038116620006c157604051638944034760e01b81526001600160a01b038216600482015260240160405180910390fd5b6000620006cd6200072a565b80546040519192506001600160a01b03808516929116907fa3396fd7f6e0a21b50e5089d2da70d5ac0a3bbbd1f617a93f134b7638998019890600090a380546001600160a01b0319166001600160a01b0392909216919091179055565b60008060ff196200075d60017f2d9a26166e6ed6927c3e14d89c4d572a357f984344fcaed3c38c3ff5cdb83f3562000ae5565b6040516020016200077091815260200190565b60408051601f1981840301815291905280516020909101201692915050565b80516001600160a01b0381168114620007a757600080fd5b919050565b60008060008060008060c08789031215620007c657600080fd5b620007d1876200078f565b9550620007e1602088016200078f565b9450620007f1604088016200078f565b935062000801606088016200078f565b925062000811608088016200078f565b91506200082160a088016200078f565b90509295509295509295565b634e487b7160e01b600052604160045260246000fd5b600181811c908216806200085857607f821691505b6020821081036200087957634e487b7160e01b600052602260045260246000fd5b50919050565b601f821115620008cf576000816000526020600020601f850160051c81016020861015620008aa5750805b601f850160051c820191505b81811015620008cb57828155600101620008b6565b5050505b505050565b81516001600160401b03811115620008f057620008f06200082d565b620009088162000901845462000843565b846200087f565b602080601f831160018114620009405760008415620009275750858301515b600019600386901b1c1916600185901b178555620008cb565b600085815260208120601f198616915b82811015620009715788860151825594840194600190910190840162000950565b5085821015620009905787850151600019600388901b60f8161c191681555b5050505050600190811b01905550565b634e487b7160e01b600052603260045260246000fd5b60008151808452602080850194506020840160005b83811015620009f35781516001600160e01b03191687529582019590820190600101620009cb565b509495945050505050565b6001600160a01b0384168152606060208083018290526000919062000a2684830187620009b6565b604085820360408701528187518084528484019150848160051b850101858a016000805b8481101562000ad257601f1980898603018752835180518987528051808b890152855b8181101562000a8d578d81840101518f828b0101528d8101905062000a6d565b508781018e01869052918c0151601f90920190921686018681038d018c880152915062000abd828d0182620009b6565b978b0197955050509188019160010162000a4a565b50919d9c50505050505050505050505050565b8181038181111562000b0757634e487b7160e01b600052601160045260246000fd5b92915050565b6153de8062000b1d6000396000f3fe608060405234801561001057600080fd5b50600436106102545760003560e01c806385e1969b11610147578063b8b249c6116100be578063b8b249c61461050c578063bd9098e71461051f578063bf7e214f14610532578063c7b50db01461053a578063c7fc09841461054d578063c927402114610560578063d195a53f14610573578063df657fdf14610586578063e9c0c82f1461058e578063ef07a7bf146105a1578063f3d5b093146105b4578063f588b157146105bc57600080fd5b806385e1969b146104355780638a622d49146104485780638b016cf91461045b5780638ee9cb451461046e57806394f1bd0714610481578063961cafa814610494578063977701cf146104a75780639ccafcdc146104ba5780639f42f00d146104cb578063a30886bf146104de578063ab3b0387146104f1578063b02047d8146104f957600080fd5b806332f19a5c116101db57806332f19a5c146103465780633cd6f55e1461036f57806348b1ee1f146103825780634d4b9886146103955780635964460f146103a85780635f45ea54146103bb57806360f6b132146103ce5780636c9d68b7146103e157806378eefeec146103f45780637997ad5b146104075780637c4bdb0e1461041a578063816c33f01461042d57600080fd5b80628f47d114610259578063010f81271461027c57806304473f271461029f57806306fdde03146102c057806313295898146102d557806314488fe4146102ea57806315b0669f146102fd5780631727c0b0146103105780631a953191146103185780632afaaa941461032b5780632c56c23414610333575b600080fd5b6002546001600160a01b03165b6040516102739190614376565b60405180910390f35b61028f61028a3660046143f2565b6105cd565b6040519015158152602001610273565b6102b26102ad366004614464565b6108c4565b604051908152602001610273565b6102c8610af4565b60405161027391906144eb565b6102e86102e336600461481b565b610b86565b005b6102e86102f836600461489f565b61106b565b61028f61030b3660046148c1565b611206565b6102b2600781565b61028f61032636600461489f565b611362565b6102b2600681565b61028f610341366004614910565b6114b7565b61026661035436600461495f565b6000908152600660205260409020546001600160a01b031690565b61028f61037d36600461495f565b611516565b6102e861039036600461495f565b611583565b61028f6103a336600461489f565b6115c7565b6102b26103b6366004614998565b611695565b61028f6103c9366004614a73565b611b4d565b61028f6103dc36600461489f565b611d5c565b61028f6103ef366004614b15565b611ddf565b61028f610402366004614ba0565b611f3d565b61028f610415366004614c59565b61224a565b61028f610428366004614c97565b6122a3565b6102b2600981565b61028f610443366004614d14565b612380565b61028f610456366004614c59565b6124b2565b6102b261046936600461489f565b6125c5565b61028f61047c366004614d31565b6125f1565b61028f61048f36600461489f565b612655565b61028f6104a23660046148c1565b61282d565b61028f6104b5366004614d4d565b6129a2565b6003546001600160a01b0316610266565b61028f6104d936600461489f565b612abd565b61028f6104ec36600461489f565b612ae4565b6102b2600b81565b61028f61050736600461489f565b612cd8565b61028f61051a36600461495f565b612cff565b61028f61052d366004614d83565b612da0565b610266612fa4565b61028f610548366004614d14565b612fbd565b61028f61055b366004614b15565b612ffa565b61028f61056e366004614d14565b6131ed565b6102b2610581366004614d14565b6132a0565b6102b2600881565b61028f61059c36600461495f565b613308565b61028f6105af366004614e0d565b61336b565b6102b2600c81565b6004546001600160a01b0316610266565b60006105d76133f5565b60006105e2886132a0565b6000818152600660205260409020549091506001600160a01b0316806106235760405162461bcd60e51b815260040161061a90614e4b565b60405180910390fd5b6000816001600160a01b03166332c9b1836040518163ffffffff1660e01b8152600401602060405180830381865afa158015610663573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906106879190614e80565b6040516323b872dd60e01b81523360048201526001600160a01b038481166024830152604482018c9052919250908216906323b872dd90606401600060405180830381600087803b1580156106db57600080fd5b505af11580156106ef573d6000803e3d6000fd5b5060009250600391506106ff9050565b604051908082528060200260200182016040528015610728578160200160208202803683370190505b509050898160008151811061073f5761073f614e9d565b602002602001018181525050888160018151811061075f5761075f614e9d565b602002602001018181525050878160028151811061077f5761077f614e9d565b602090810291909101015260405163598b8e7160e01b81526000906001600160a01b0385169063598b8e71906107b9908590600401614eef565b6020604051808303816000875af11580156107d8573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906107fc9190614f02565b9050806108435760405162461bcd60e51b8152602060048201526015602482015274115c98cdcc8c4819195c1bdcda5d0819985a5b1959605a1b604482015260640161061a565b6003546040516378029b4960e11b81526001600160a01b039091169063f005369290610875908b908b90600401614f48565b600060405180830381600087803b15801561088f57600080fd5b505af11580156108a3573d6000803e3d6000fd5b505050506001955050505050506108ba6001600055565b9695505050505050565b60006108dc336000356001600160e01b03191661341f565b600060085460016108ed9190614f72565b600081815260066020818152604080842080546001600160a01b0319166001600160a01b038c81169190911782558a81168087526007855283872088905595879052939092529054600254600480546005549451636e10f4cb60e01b81529182018890526024820196909652604481018a90529084166064820152938316608485015290821660a48401529293509190911690636e10f4cb9060c4016020604051808303816000875af11580156109a8573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906109cc9190614f02565b5060006109d7612fa4565b604051630f317a8560e41b815260206004820152600a60248201526910d3d25397d59055531560b21b60448201529091506001600160a01b038216906325c471a090829063f317a85090606401602060405180830381865afa158015610a41573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610a659190614f85565b6040516001600160e01b031960e084901b1681526001600160401b0390911660048201526001600160a01b038916602482015260006044820152606401600060405180830381600087803b158015610abc57600080fd5b505af1158015610ad0573d6000803e3d6000fd5b505060088054925090506000610ae583614fa2565b90915550919695505050505050565b606060018054610b0390614fbb565b80601f0160208091040260200160405190810160405280929190818152602001828054610b2f90614fbb565b8015610b7c5780601f10610b5157610100808354040283529160200191610b7c565b820191906000526020600020905b815481529060010190602001808311610b5f57829003601f168201915b5050505050905090565b610b8e6133f5565b6000610bbc84604001518560800151600081518110610baf57610baf614e9d565b60200260200101516125c5565b60008181526012602090815260408083208b8452601390925290912054919250908214610bfc57604051635004560b60e11b815260040160405180910390fd5b610c0588613308565b15610c1a57610c138861356a565b5050611059565b60016012820154600160401b900460ff166004811115610c3c57610c3c614fef565b14610c5a57604051631e3f269560e11b815260040160405180910390fd5b6000610c65886132a0565b905060005b600f830154811015610d53576001830154600090815260066020526040902054600b840180546001600160a01b03909216916346d07259919084908110610cb357610cb3614e9d565b906000526020600020015485600301600c018481548110610cd657610cd6614e9d565b90600052602060002001546040518363ffffffff1660e01b8152600401610d07929190918252602082015260400190565b6020604051808303816000875af1158015610d26573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610d4a9190614f02565b50600101610c6a565b506002820154600090610d689060ff1661363f565b90506000610d758961363f565b60128501805491925060029160ff60401b1916600160401b83021790555060408051610180810182526003860180546101408301908152600488015461016084015260e08301908152835160808101808652610fec95859392859261010086019260058e01918391820190839060029082845b815481526020019060010190808311610de857505050918352505060408051808201918290526020909201919060028481019182845b815481526020019060010190808311610e1e575050505050815250508152602001600682016040518060400160405290816000820154815260200160018201548152505081525050815260200160088201805480602002602001604051908101604052809291908181526020018280548015610eb957602002820191906000526020600020905b815481526020019060010190808311610ea5575b5050505050815260200160098201548152602001600a8201805480602002602001604051908101604052809291908181526020018280548015610f1b57602002820191906000526020600020905b815481526020019060010190808311610f07575b50505050508152602001600b8201805480602002602001604051908101604052809291908181526020018280548015610f7357602002820191906000526020600020905b815481526020019060010190808311610f5f575b50505050508152602001600c8201805480602002602001604051908101604052809291908181526020018280548015610fcb57602002820191906000526020600020905b815481526020019060010190808311610fb7575b50505050508152602001600d82015481525050898660010154868686613661565b50600354604051630f2b65dd60e21b81526001600160a01b0390911690633cad977490611021908e908b908b90600401615005565b600060405180830381600087803b15801561103b57600080fd5b505af115801561104f573d6000803e3d6000fd5b5050505050505050505b6110636001600055565b505050505050565b6110736133f5565b600082815260136020908152604080832054808452601290925290912060016012820154600160401b900460ff1660048111156110b2576110b2614fef565b146110d057604051631e3f269560e11b815260040160405180910390fd5b6002546040805180820182528581526020810186905290516314d2f97b60e11b81526000926001600160a01b0316916329a5f2f6916111129190600401615042565b602060405180830381865afa15801561112f573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906111539190615050565b90508160110154811461117957604051633eaf74ed60e21b815260040160405180910390fd5b60128201805460ff60401b1916600360401b17905561119782613ae7565b6003546040516345c73e6b60e01b8152600481018790526001600160a01b03909116906345c73e6b90602401600060405180830381600087803b1580156111dd57600080fd5b505af11580156111f1573d6000803e3d6000fd5b505050505050506112026001600055565b5050565b600061121e336000356001600160e01b03191661341f565b6000828152600c6020908152604080832054878452600690925280832054905163c695cea960e01b81526001600160a01b03928316939190921691829063c695cea99061126f908990600401614eef565b602060405180830381865afa15801561128c573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906112b09190615050565b604051635334229560e01b815260048101899052602481018290529091506001600160a01b038416906353342295906044016020604051808303816000875af1158015611301573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906113259190614f02565b508481887f0d4ed26c1e8d7e7f1953cc13b4883a9dc52f61c721ba4703fe07c994af193b6c60405160405180910390a45060019695505050505050565b600061137a336000356001600160e01b03191661341f565b6000838152600660205260409020546001600160a01b0316806113af5760405162461bcd60e51b815260040161061a90615069565b6040805160018082528183019092526000916020808301908036833701905050905083816000815181106113e5576113e5614e9d565b602090810291909101015260405163598b8e7160e01b81526000906001600160a01b0384169063598b8e719061141f908590600401614eef565b6020604051808303816000875af115801561143e573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906114629190614f02565b9050806114a95760405162461bcd60e51b8152602060048201526015602482015274115b9e59db584819195c1bdcda5d0819985a5b1959605a1b604482015260640161061a565b600193505050505b92915050565b60006114c16133f5565b60006114cc846132a0565b9050611509818460405180604001604052806018815260200177115490cc4c4d4d481b5a5e08199d5b991cc819985a5b195960421b815250613d55565b9150506114b16001600055565b600061152e336000356001600160e01b03191661341f565b60008281526009602052604090205460ff16801561155f57604051632440cbf960e11b815260040160405180910390fd5b50506000818152600960205260409020805460ff191660019081179091555b919050565b61158b6133f5565b61159481613308565b6115b1576040516355803c3760e01b815260040160405180910390fd5b6115ba8161356a565b6115c46001600055565b50565b60006115df336000356001600160e01b03191661341f565b6000828152600c602052604090819020549051630de4559960e41b8152600481018590526001600160a01b0390911690819063de455990906024016020604051808303816000875af1158015611639573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061165d9190614f02565b50604051839085907f68c771c05d01f5c53189e3a23bc1fac5a3453efcd0992bddf8f390b581c2389c90600090a35060019392505050565b600061169f6133f5565b60008b815260136020526040902054156116cc576040516330c445f360e01b815260040160405180910390fd5b60006116d7876132a0565b60008181526006602052604090819020549051637bd1fb1560e11b81529192506001600160a01b031690819063f7a3f62a9061171790899060040161518b565b602060405180830381865afa158015611734573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906117589190614f02565b508560c0015160000361177e5760405163045aaba960e41b815260040160405180910390fd5b8360000361179f57604051633eaf74ed60e21b815260040160405180910390fd5b60005b8660a001515181101561186c57816001600160a01b0316637628c043886020015183815181106117d4576117d4614e9d565b60200260200101518960a0015184815181106117f2576117f2614e9d565b60200260200101516040518363ffffffff1660e01b8152600401611820929190918252602082015260400190565b6020604051808303816000875af115801561183f573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906118639190614f02565b506001016117a2565b50611899866080015160008151811061188757611887614e9d565b602002602001015187604001516125c5565b92506000808481526012602081905260409091200154600160401b900460ff1660048111156118ca576118ca614fef565b146118e8576040516330c445f360e01b815260040160405180910390fd5b60008381526012602052604090208d815560018082018490556002820180548a9260ff1990911690838181111561192157611921614fef565b02179055508651805180516003840190815560209182015160048501559082015180518a93918391600587019061195b90829060026142a5565b50602082015161197190600280840191906142a5565b50505060409190910151805160068301556020908101516007909201919091558281015180516119a792600885019201906142e3565b5060408201516009820155606082015180516119cd91600a8401916020909101906142e3565b50608082015180516119e991600b8401916020909101906142e3565b5060a08201518051611a0591600c8401916020909101906142e3565b5060c09190910151600d9091015560118101859055611a24864261519e565b6012820180546001600160401b039290921667ffffffffffffffff1983168117825560019268ffffffffffffffffff191617600160401b83021790555083601360008360000154815260200190815260200160002081905550600360009054906101000a90046001600160a01b03166001600160a01b031662e5996882600001548f8f8f8f8d60800151600081518110611ac057611ac0614e9d565b602090810291909101015160128901546040516001600160e01b031960e08a901b168152611b00979695949392916001600160401b0316906004016151c5565b600060405180830381600087803b158015611b1a57600080fd5b505af1158015611b2e573d6000803e3d6000fd5b50505050505050611b3f6001600055565b9a9950505050505050505050565b6000611b576133f5565b6000611b62896132a0565b6000818152600660205260409020549091506001600160a01b031680611b9a5760405162461bcd60e51b815260040161061a90614e4b565b6040805160038082526080820190925260009160208201606080368337019050509050600181600081518110611bd257611bd2614e9d565b6020026020010181815250508981600181518110611bf257611bf2614e9d565b6020026020010181815250508781600281518110611c1257611c12614e9d565b6020908102919091010152604051635702fc4160e01b81526000906001600160a01b03841690635702fc4190611c509085908e908d90600401615216565b6020604051808303816000875af1158015611c6f573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190611c939190614f02565b905080611cdb5760405162461bcd60e51b8152602060048201526016602482015275115c98cdcc8c481dda5d1a191c985dc819985a5b195960521b604482015260640161061a565b6003546040516378029b4960e11b81526001600160a01b039091169063f005369290611d0d908a908a90600401614f48565b600060405180830381600087803b158015611d2757600080fd5b505af1158015611d3b573d6000803e3d6000fd5b505050506001945050505050611d516001600055565b979650505050505050565b6000818152600c6020526040808220549051635949f44f60e01b8152600481018590526001600160a01b03909116908190635949f44f90602401602060405180830381865afa158015611db3573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190611dd79190614f02565b949350505050565b6000611de96133f5565b6040820151602083015151600090611e02576000611e22565b8360200151600081518110611e1957611e19614e9d565b60200260200101515b60408051600280825260608201835292935060009290916020830190803683370190505090508281600081518110611e5c57611e5c614e9d565b6020026020010181815250508181600181518110611e7c57611e7c614e9d565b6020908102919091010152600480548651604051634c948fb360e11b81526001600160a01b03909216926399291f6692611ebd92600c92909187910161524a565b602060405180830381865afa158015611eda573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190611efe9190614f02565b50604051829084907f78dd8c3ba670aacf126fdf991710d8ae9559838134b7e2f38d45911b7bd4237790600090a36001935050505061157e6001600055565b6000611f476133f5565b6000611f528a6132a0565b6000818152600660205260409020549091506001600160a01b031680611f8a5760405162461bcd60e51b815260040161061a90615272565b6000816001600160a01b03166332c9b1836040518163ffffffff1660e01b8152600401602060405180830381865afa158015611fca573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190611fee9190614e80565b9050806001600160a01b031663f242432a33848e8e8e6040518663ffffffff1660e01b81526004016120249594939291906152a9565b600060405180830381600087803b15801561203e57600080fd5b505af1158015612052573d6000803e3d6000fd5b5060009250600491506120629050565b60405190808252806020026020018201604052801561208b578160200160208202803683370190505b5090508a816000815181106120a2576120a2614e9d565b6020026020010181815250508b816001815181106120c2576120c2614e9d565b60200260200101818152505088816002815181106120e2576120e2614e9d565b602002602001018181525050878160038151811061210257612102614e9d565b602090810291909101015260405163598b8e7160e01b81526000906001600160a01b0385169063598b8e719061213c908590600401614eef565b6020604051808303816000875af115801561215b573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061217f9190614f02565b9050806121c75760405162461bcd60e51b8152602060048201526016602482015275115c98cc4c4d4d4819195c1bdcda5d0819985a5b195960521b604482015260640161061a565b6003546040516378029b4960e11b81526001600160a01b039091169063f0053692906121f9908b908b90600401614f48565b600060405180830381600087803b15801561221357600080fd5b505af1158015612227573d6000803e3d6000fd5b5050505060019550505050505061223e6001600055565b98975050505050505050565b6000612262336000356001600160e01b03191661341f565b61229c838360405180604001604052806017815260200176115b9e59db58481b5a5e08199d5b991cc819985a5b1959604a1b815250613d55565b9392505050565b60006122bb336000356001600160e01b03191661341f565b6122c36133f5565b600d546000818152600c60205260409081902080546001600160a01b0319166001600160a01b038916908117909155905163787b254360e01b815263787b2543906123189084908990899089906004016152e3565b6020604051808303816000875af1158015612337573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061235b9190614f02565b50600d805490600061236c83614fa2565b91905055506001915050611dd76001600055565b6000612398336000356001600160e01b03191661341f565b60006123a2612fa4565b604051630f317a8560e41b815260206004820152600e60248201526d22a72ca3a6a0afa1a922a0aa27a960911b60448201529091506001600160a01b038216906312c1e07290829063f317a85090606401602060405180830381865afa158015612410573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906124349190614f85565b6040516001600160e01b031960e084901b1681526001600160401b0390911660048201526001600160a01b038616602482015230604482015260006064820152608401600060405180830381600087803b15801561249157600080fd5b505af11580156124a5573d6000803e3d6000fd5b5060019695505050505050565b60006124ca336000356001600160e01b03191661341f565b6000838152600660205260409020546001600160a01b0316806124ff5760405162461bcd60e51b815260040161061a90615069565b6040805160008082526020820192839052635702fc4160e01b909252906001600160a01b038316635702fc4161253a84848960248301615216565b6020604051808303816000875af1158015612559573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061257d9190614f02565b9050806114a95760405162461bcd60e51b8152602060048201526016602482015275115b9e59db58481dda5d1a191c985dc819985a5b195960521b604482015260640161061a565b604080516020808201949094528082019290925280518083038201815260609092019052805191012090565b8051602080830151604080518084019490945283810191909152805180840382018152606090930181528251928201929092206000818152601090925291812060020154909190156126465750600192915050565b50600092915050565b50919050565b600061266d336000356001600160e01b03191661341f565b600d54831015806126805750600d548210155b1561269e5760405163aa7aacfd60e01b815260040160405180910390fd5b6000838152600c602090815260409182902054825163544e172d60e11b815292516001600160a01b0390911692839263a89c2e5a926004808401938290030181865afa1580156126f2573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906127169190614f02565b6127335760405163662dd32360e01b815260040160405180910390fd5b6000838152600c602090815260409182902054825163544e172d60e11b815292516001600160a01b0390911692839263a89c2e5a926004808401938290030181865afa158015612787573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906127ab9190614f02565b6127c85760405163662dd32360e01b815260040160405180910390fd5b60006127d486866125c5565b6000818152600b602052604090205490915060ff161561280757604051630838955960e41b815260040160405180910390fd5b6000818152600b60205260409020805460ff191660011790556001935050505092915050565b60408051600580825260c082019092526000918291906020820160a08036833701905050905060015b845181101561289e5784818151811061287157612871614e9d565b602002602001015182828151811061288b5761288b614e9d565b6020908102919091010152600101612856565b506000838152600c6020908152604080832054888452600690925280832054905163c695cea960e01b81526001600160a01b0392831693929091169063c695cea9906128ee908690600401614eef565b602060405180830381865afa15801561290b573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061292f9190615050565b60405163dbc6179760e01b815260048101899052602481018290529091506001600160a01b0383169063dbc6179790604401602060405180830381865afa15801561297e573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190611d519190614f02565b60006129ba336000356001600160e01b03191661341f565b81516020808401516040805180840194909452838101919091528051808403820181526060909301815282519282019290922060008181526010909252918120600201549003612a8f57612a0c61431d565b8581526020808201868152604080840187815260008681526010909452922083518155905160018201559051829190612a4b90600280840191906142a5565b509050508486837ffcd227d6493fed2db050968ecc93ebcf6bde9054733712fc92263419b075554c87604051612a819190615042565b60405180910390a450612ab2565b604051631e5cec8960e21b8152600481018690526024810185905260440161061a565b506001949350505050565b600080612aca84846125c5565b6000908152600a602052604090205460ff16949350505050565b6000612afc336000356001600160e01b03191661341f565b600d5483101580612b0f5750600d548210155b15612b2d5760405163aa7aacfd60e01b815260040160405180910390fd5b6000838152600c602090815260409182902054825163544e172d60e11b815292516001600160a01b0390911692839263a89c2e5a926004808401938290030181865afa158015612b81573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190612ba59190614f02565b612bc25760405163662dd32360e01b815260040160405180910390fd5b6000838152600c602090815260409182902054825163544e172d60e11b815292516001600160a01b0390911692839263a89c2e5a926004808401938290030181865afa158015612c16573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190612c3a9190614f02565b15612c585760405163662dd32360e01b815260040160405180910390fd5b604080516020808201889052818301879052825180830384018152606090920183528151918101919091206000818152600a9092529190205460ff1615612cb257604051630838955960e41b815260040160405180910390fd5b6000818152600a60205260409020805460ff191660011790556001935050505092915050565b600080612ce584846125c5565b6000908152600b602052604090205460ff16949350505050565b6000612d17336000356001600160e01b03191661341f565b6000828152601060205260408120600201549003612d4b57604051634947608360e01b81526004810183905260240161061a565b60008281526010602052604080822082815560018101839055600281018390556003018290555183917f368b4b23dfa2c5bf6dd7ee4303c8bed712aba2a8ce553f7f9781dc2056e199cd91a25b506001919050565b6000612daa6133f5565b6000612db58a6132a0565b6000818152600660205260409020549091506001600160a01b031680612ded5760405162461bcd60e51b815260040161061a90615272565b60408051600380825260808201909252600091602082016060803683370190505090508981600081518110612e2457612e24614e9d565b6020026020010181815250508a81600181518110612e4457612e44614e9d565b6020026020010181815250508781600281518110612e6457612e64614e9d565b6020908102919091010152604051635702fc4160e01b81526000906001600160a01b03841690635702fc4190612ea29085908e908d90600401615216565b6020604051808303816000875af1158015612ec1573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190612ee59190614f02565b905080612f2e5760405162461bcd60e51b8152602060048201526017602482015276115c98cc4c4d4d481dda5d1a191c985dc819985a5b1959604a1b604482015260640161061a565b6003546040516378029b4960e11b81526001600160a01b039091169063f005369290612f60908a908a90600401614f48565b600060405180830381600087803b158015612f7a57600080fd5b505af1158015612f8e573d6000803e3d6000fd5b50505050600194505050505061223e6001600055565b6000612fae613e44565b546001600160a01b0316919050565b6000612fd5336000356001600160e01b03191661341f565b50600480546001600160a01b0383166001600160a01b03199091161790556001919050565b6000613012336000356001600160e01b03191661341f565b61301b82613ea6565b5060008083602001515111613031576000613051565b826020015160008151811061304857613048614e9d565b60200260200101515b905060008084608001515111613068576000613088565b836080015160008151811061307f5761307f614e9d565b60200260200101515b6000818152600f60205260408120600c0154919250036131ca576000818152600f602090815260409091208551805180518355830151600183015591820151805187939183916002808401916130e0918391906142a5565b5060208201516130f690600280840191906142a5565b505050604091909101518051600683015560209081015160079092019190915582810151805161312c92600885019201906142e3565b50604082015160098201556060820151805161315291600a8401916020909101906142e3565b506080820151805161316e91600b8401916020909101906142e3565b5060a0820151805161318a91600c8401916020909101906142e3565b5060c09190910151600d90910155604051819083907fc721fb519482fab66604e05e7929747c245d19169d1046e78a02f7cbf99cda0590600090a36131e3565b604051639f8ef50d60e01b815260040160405180910390fd5b5060019392505050565b6000613205336000356001600160e01b03191661341f565b600580546001600160a01b0319166001600160a01b038481169182179092556002546004805460405163fab6874960e01b81529285169183019190915290921660248301529063fab68749906044016020604051808303816000875af1158015613273573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906132979190614f02565b50600192915050565b6001600160a01b0381166000908152600760205260408120548082036114b15760405162461bcd60e51b815260206004820152601f60248201527f436f6e74726163742061646472657373206e6f74207265676973746572656400604482015260640161061a565b60008181526013602090815260408083205480845260129092528220826012820154600160401b900460ff16600481111561334557613345614fef565b03613354575060009392505050565b601201546001600160401b03164210159392505050565b6000818152600c602052604080822054905163221ea59760e21b81526001600160a01b0390911690819063887a965c906133ab908890889060040161530f565b602060405180830381865afa1580156133c8573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906133ec9190614f02565b95945050505050565b60026000540361341857604051633ee5aeb560e01b815260040160405180910390fd5b6002600055565b6000613429613e44565b80549091506001600160a01b031680613458576000604051638944034760e01b815260040161061a9190614376565b60405163b700961360e01b81526001600160a01b0385811660048301523060248301526001600160e01b031985166044830152600091829182919085169063b700961390606401606060405180830381865afa1580156134bc573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906134e09190615328565b9250925092508261356157801561350a5760405163cc9855ad60e01b815260040160405180910390fd5b63ffffffff8216156135465760405163a426878960e01b81526001600160a01b038816600482015263ffffffff8316602482015260440161061a565b86604051632ecd3d0360e21b815260040161061a9190614376565b50505050505050565b600081815260136020908152604080832054808452601290925290912060016012820154600160401b900460ff1660048111156135a9576135a9614fef565b146135c757604051631e3f269560e11b815260040160405180910390fd5b60128101805460ff60401b1916600160421b1790556135e581613ae7565b6003546040516346bd060d60e11b8152600481018590526001600160a01b0390911690638d7a0c1a90602401600060405180830381600087803b15801561362b57600080fd5b505af1158015613561573d6000803e3d6000fd5b60008082600181111561365457613654614fef565b03612d9857506000919050565b6000856080015160008151811061367a5761367a614e9d565b60200260200101518760400151146136a5576040516307138dd360e21b815260040160405180910390fd5b86608001516000815181106136bc576136bc614e9d565b60200260200101518660400151146136e757604051634606462960e01b815260040160405180910390fd5b6000838152600c60205260409081902054905163221ea59760e21b81526001600160a01b039091169063887a965c906137269088908b9060040161530f565b602060405180830381865afa158015613743573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906137679190614f02565b613784576040516302c723f960e11b815260040160405180910390fd5b6000828152600c60205260409081902054905163221ea59760e21b81526001600160a01b039091169063887a965c906137c39087908a9060040161530f565b602060405180830381865afa1580156137e0573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906138049190614f02565b613821576040516302c723f960e11b815260040160405180910390fd5b60008481526006602052604090819020549051637bd1fb1560e11b81526001600160a01b039091169063f7a3f62a9061385e90899060040161518b565b602060405180830381865afa15801561387b573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061389f9190614f02565b506000858152600660205260409081902054905163c41daf5f60e01b81526001600160a01b039091169063c41daf5f906138dd908a9060040161518b565b6020604051808303816000875af11580156138fc573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906139209190614f02565b506000848152600660205260409081902054905163c41daf5f60e01b81526001600160a01b039091169063c41daf5f9061395e90899060040161518b565b6020604051808303816000875af115801561397d573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906139a19190614f02565b5060008481526006602052604090819020549051630c67ff8f60e21b81526001600160a01b039091169063319ffe3c906139df90899060040161518b565b6020604051808303816000875af11580156139fe573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190613a229190614f02565b5060008581526006602052604090819020549051630c67ff8f60e21b81526001600160a01b039091169063319ffe3c90613a60908a9060040161518b565b6020604051808303816000875af1158015613a7f573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190613aa39190614f02565b50866040015186604001517ff5b268a3ff315cc44ccceeef86259c9e8eef81ceecb14001543809115380dd6260405160405180910390a35060019695505050505050565b6001818101546000908152600660205260408082205481518481528083019092526001600160a01b03169281602001602082028036833701905050905082600301600d015481600081518110613b3f57613b3f614e9d565b6020908102919091010152604051630f695fef60e11b81526001600160a01b03831690631ed2bfde90613b76908490600401614eef565b6020604051808303816000875af1158015613b95573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190613bb99190614f02565b5060005b600f840154811015613d4f57826001600160a01b03166346d07259856003016008018381548110613bf057613bf0614e9d565b906000526020600020015486600301600c018481548110613c1357613c13614e9d565b90600052602060002001546040518363ffffffff1660e01b8152600401613c44929190918252602082015260400190565b6020604051808303816000875af1158015613c63573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190613c879190614f02565b50826001600160a01b0316637b22a8ac856003016008018381548110613caf57613caf614e9d565b906000526020600020015486600301600c018481548110613cd257613cd2614e9d565b90600052602060002001546040518363ffffffff1660e01b8152600401613d03929190918252602082015260400190565b6020604051808303816000875af1158015613d22573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190613d469190614f02565b50600101613bbd565b50505050565b6000838152600660205260408120546001600160a01b031680613db15760405162461bcd60e51b815260206004820152601460248201527315985d5b1d081b9bdd081c9959da5cdd195c995960621b604482015260640161061a565b60405163e73c500960e01b81526000906001600160a01b0383169063e73c500990613de090889060040161518b565b6020604051808303816000875af1158015613dff573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190613e239190614f02565b905083816124a55760405162461bcd60e51b815260040161061a91906144eb565b60008060ff19613e7560017f2d9a26166e6ed6927c3e14d89c4d572a357f984344fcaed3c38c3ff5cdb83f3561537e565b604051602001613e8791815260200190565b60408051601f1981840301815291905280516020909101201692915050565b60a081015151608082015151600091906002821015613ed85760405163191f39cb60e21b815260040160405180910390fd5b8015613ef757604051630243591560e21b815260040160405180910390fd5b60008085602001515111613f0c576000613f2c565b8460200151600081518110613f2357613f23614e9d565b60200260200101515b90506000600186602001515111613f44576000613f64565b8560200151600181518110613f5b57613f5b614e9d565b60200260200101515b6000818152600c60205260409081902054905163221ea59760e21b81529192506001600160a01b03169063887a965c90613fa49085908a9060040161530f565b602060405180830381865afa158015613fc1573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190613fe59190614f02565b614002576040516302c723f960e11b815260040160405180910390fd5b60008281526006602052604090819020549051639d3f69a160e01b81526001600160a01b0390911690639d3f69a19061403f90899060040161518b565b6020604051808303816000875af115801561405e573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906140829190614f02565b506000614090856003615391565b61409b906001614f72565b6140a6906002614f72565b6001600160401b038111156140bd576140bd61450d565b6040519080825280602002602001820160405280156140e6578160200160208202803683370190505b50905086604001518160008151811061410157614101614e9d565b6020908102919091010152600160005b86811015614168578860200151818151811061412f5761412f614e9d565b602002602001015183838061414390614fa2565b94508151811061415557614155614e9d565b6020908102919091010152600101614111565b5060005b868110156141c3578860600151818151811061418a5761418a614e9d565b602002602001015183838061419e90614fa2565b9450815181106141b0576141b0614e9d565b602090810291909101015260010161416c565b5060005b8681101561421e578860a0015181815181106141e5576141e5614e9d565b60200260200101518383806141f990614fa2565b94508151811061420b5761420b614e9d565b60209081029190910101526001016141c7565b50600480548951604051634c948fb360e11b81526001600160a01b03909216926399291f669261425592600b92909188910161524a565b602060405180830381865afa158015614272573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906142969190614f02565b50600198975050505050505050565b82600281019282156142d3579160200282015b828111156142d35782518255916020019190600101906142b8565b506142df929150614343565b5090565b8280548282559060005260206000209081019282156142d357916020028201828111156142d35782518255916020019190600101906142b8565b6040518060600160405280600081526020016000815260200161433e614358565b905290565b5b808211156142df5760008155600101614344565b60405180604001604052806002906020820280368337509192915050565b6001600160a01b0391909116815260200190565b6001600160a01b03811681146115c457600080fd5b803561157e8161438a565b60008083601f8401126143bc57600080fd5b5081356001600160401b038111156143d357600080fd5b6020830191508360208285010111156143eb57600080fd5b9250929050565b60008060008060008060a0878903121561440b57600080fd5b86356144168161438a565b955060208701359450604087013593506060870135925060808701356001600160401b0381111561444657600080fd5b61445289828a016143aa565b979a9699509497509295939492505050565b60008060006060848603121561447957600080fd5b83356144848161438a565b925060208401356144948161438a565b929592945050506040919091013590565b6000815180845260005b818110156144cb576020818501810151868301820152016144af565b506000602082860101526020601f19601f83011685010191505092915050565b60208152600061229c60208301846144a5565b80356002811061157e57600080fd5b634e487b7160e01b600052604160045260246000fd5b604080519081016001600160401b03811182821017156145455761454561450d565b60405290565b60405160e081016001600160401b03811182821017156145455761454561450d565b604051601f8201601f191681016001600160401b03811182821017156145955761459561450d565b604052919050565b6000604082840312156145af57600080fd5b6145b7614523565b9050813581526020820135602082015292915050565b600082601f8301126145de57600080fd5b6145e6614523565b8060408401858111156145f857600080fd5b845b818110156146125780358452602093840193016145fa565b509095945050505050565b600081830361010081121561463157600080fd5b604051606081018181106001600160401b03821117156146535761465361450d565b604052915081614663858561459d565b81526080603f198301121561467757600080fd5b61467f614523565b915061468e85604086016145cd565b825261469d85608086016145cd565b60208301528160208201526146b58560c0860161459d565b6040820152505092915050565b600082601f8301126146d357600080fd5b813560206001600160401b038211156146ee576146ee61450d565b8160051b6146fd82820161456d565b928352848101820192828101908785111561471757600080fd5b83870192505b84831015611d515782358252918301919083019061471d565b60006101c0828403121561474957600080fd5b61475161454b565b905061475d838361461d565b81526101008201356001600160401b038082111561477a57600080fd5b614786858386016146c2565b602084015261012084013560408401526101408401359150808211156147ab57600080fd5b6147b7858386016146c2565b60608401526101608401359150808211156147d157600080fd5b6147dd858386016146c2565b60808401526101808401359150808211156147f757600080fd5b50614804848285016146c2565b60a0830152506101a082013560c082015292915050565b60008060008060008060a0878903121561483457600080fd5b8635955060208701356148468161438a565b9450614854604088016144fe565b935060608701356001600160401b038082111561487057600080fd5b61487c8a838b01614736565b9450608089013591508082111561489257600080fd5b5061445289828a016143aa565b600080604083850312156148b257600080fd5b50508035926020909101359150565b6000806000606084860312156148d657600080fd5b8335925060208401356001600160401b038111156148f357600080fd5b6148ff868287016146c2565b925050604084013590509250925092565b6000806040838503121561492357600080fd5b823561492e8161438a565b915060208301356001600160401b0381111561494957600080fd5b61495585828601614736565b9150509250929050565b60006020828403121561497157600080fd5b5035919050565b6001600160401b03811681146115c457600080fd5b803561157e81614978565b6000806000806000806000806000806101008b8d0312156149b857600080fd5b8a35995060208b01356001600160401b03808211156149d657600080fd5b6149e28e838f016143aa565b909b50995060408d01359150808211156149fb57600080fd5b614a078e838f016143aa565b9099509750879150614a1b60608e0161439f565b9650614a2960808e016144fe565b955060a08d0135915080821115614a3f57600080fd5b50614a4c8d828e01614736565b935050614a5b60c08c0161498d565b915060e08b013590509295989b9194979a5092959850565b600080600080600080600060c0888a031215614a8e57600080fd5b8735614a998161438a565b9650602088013595506040880135614ab08161438a565b94506060880135935060808801356001600160401b0380821115614ad357600080fd5b614adf8b838c01614736565b945060a08a0135915080821115614af557600080fd5b50614b028a828b016143aa565b989b979a50959850939692959293505050565b600060208284031215614b2757600080fd5b81356001600160401b03811115614b3d57600080fd5b611dd784828501614736565b60006001600160401b03831115614b6257614b6261450d565b614b75601f8401601f191660200161456d565b9050828152838383011115614b8957600080fd5b828260208301376000602084830101529392505050565b60008060008060008060008060e0898b031215614bbc57600080fd5b8835614bc78161438a565b9750602089013596506040890135955060608901356001600160401b0380821115614bf157600080fd5b818b0191508b601f830112614c0557600080fd5b614c148c833560208501614b49565b965060808b0135955060a08b0135945060c08b0135915080821115614c3857600080fd5b50614c458b828c016143aa565b999c989b5096995094979396929594505050565b60008060408385031215614c6c57600080fd5b8235915060208301356001600160401b0381111561494957600080fd5b80151581146115c457600080fd5b60008060008060808587031215614cad57600080fd5b8435614cb88161438a565b935060208501356001600160401b03811115614cd357600080fd5b8501601f81018713614ce457600080fd5b614cf387823560208401614b49565b9350506040850135614d0481614c89565b9396929550929360600135925050565b600060208284031215614d2657600080fd5b813561229c8161438a565b600060408284031215614d4357600080fd5b61229c83836145cd565b600080600060808486031215614d6257600080fd5b8335925060208401359150614d7a85604086016145cd565b90509250925092565b60008060008060008060008060e0898b031215614d9f57600080fd5b8835614daa8161438a565b975060208901359650604089013595506060890135614dc88161438a565b94506080890135935060a08901356001600160401b0380821115614deb57600080fd5b614df78c838d01614736565b945060c08b0135915080821115614c3857600080fd5b600080600060608486031215614e2257600080fd5b8335925060208401356001600160401b03811115614e3f57600080fd5b6148ff86828701614736565b6020808252601b908201527a115c98cdcc8c481d985d5b1d081b9bdd081c9959da5cdd195c9959602a1b604082015260600190565b600060208284031215614e9257600080fd5b815161229c8161438a565b634e487b7160e01b600052603260045260246000fd5b60008151808452602080850194506020840160005b83811015614ee457815187529582019590820190600101614ec8565b509495945050505050565b60208152600061229c6020830184614eb3565b600060208284031215614f1457600080fd5b815161229c81614c89565b81835281816020850137506000828201602090810191909152601f909101601f19169091010190565b602081526000611dd7602083018486614f1f565b634e487b7160e01b600052601160045260246000fd5b808201808211156114b1576114b1614f5c565b600060208284031215614f9757600080fd5b815161229c81614978565b600060018201614fb457614fb4614f5c565b5060010190565b600181811c90821680614fcf57607f821691505b60208210810361264f57634e487b7160e01b600052602260045260246000fd5b634e487b7160e01b600052602160045260246000fd5b8381526040602082015260006133ec604083018486614f1f565b8060005b6002811015613d4f578151845260209384019390910190600101615023565b604081016114b1828461501f565b60006020828403121561506257600080fd5b5051919050565b6020808252601b908201527a115b9e59db58481d985d5b1d081b9bdd081c9959da5cdd195c9959602a1b604082015260600190565b6150b382825180518252602090810151910152565b60208101516150c660408401825161501f565b602001516150d7608084018261501f565b5060400151805160c08301526020015160e090910152565b60006101c06150ff84845161509e565b60208301518161010086015261511782860182614eb3565b9150506040830151610120850152606083015184820361014086015261513d8282614eb3565b91505060808301518482036101608601526151588282614eb3565b91505060a08301518482036101808601526151738282614eb3565b91505060c08301516101a08501528091505092915050565b60208152600061229c60208301846150ef565b6001600160401b038181168382160190808211156151be576151be614f5c565b5092915050565b87815260a0602082015260006151df60a08301888a614f1f565b82810360408401526151f2818789614f1f565b9150508360608301526001600160401b038316608083015298975050505050505050565b6060815260006152296060830186614eb3565b6001600160a01b038516602084015282810360408401526108ba81856150ef565b600061014085835261525f602084018661509e565b806101208401526108ba81840185614eb3565b6020808252601c908201527f45726331313535207661756c74206e6f74207265676973746572656400000000604082015260600190565b6001600160a01b03868116825285166020820152604081018490526060810183905260a060808201819052600090611d51908301846144a5565b8481526080602082015260006152fc60808301866144a5565b9315156040830152506060015292915050565b828152604060208201526000611dd760408301846150ef565b60008060006060848603121561533d57600080fd5b835161534881614c89565b602085015190935063ffffffff8116811461536257600080fd5b604085015190925061537381614c89565b809150509250925092565b818103818111156114b1576114b1614f5c565b80820281158282048414176114b1576114b1614f5c56fea264697066735822122071a3424f7a40533edf0b7f08c2fb8d518207f7034b32e8fa605a7c6ff1c5e72264736f6c63430008180033",
}

// Dvp is an auto generated Go binding around an Ethereum contract.
type Dvp struct {
	abi abi.ABI
}

// NewDvp creates a new instance of Dvp.
func NewDvp() *Dvp {
	parsed, err := DvpMetaData.ParseABI()
	if err != nil {
		panic(errors.New("invalid ABI: " + err.Error()))
	}
	return &Dvp{abi: *parsed}
}

// Instance creates a wrapper for a deployed contract instance at the given address.
// Use this to create the instance object passed to abigen v2 library functions Call, Transact, etc.
func (c *Dvp) Instance(backend bind.ContractBackend, addr common.Address) *bind.BoundContract {
	return bind.NewBoundContract(addr, c.abi, backend, backend, backend)
}

// PackConstructor is the Go binding used to pack the parameters required for
// contract deployment.
//
// Solidity: constructor(address hashPoseidonContractAddress, address enygmaFactoryAddress, address dvpErc721FactoryAddress, address dvpErc1155FactoryAddress, address dvpTeleportAddr, address authority_) returns()
func (dvp *Dvp) PackConstructor(hashPoseidonContractAddress common.Address, enygmaFactoryAddress common.Address, dvpErc721FactoryAddress common.Address, dvpErc1155FactoryAddress common.Address, dvpTeleportAddr common.Address, authority_ common.Address) []byte {
	enc, err := dvp.abi.Pack("", hashPoseidonContractAddress, enygmaFactoryAddress, dvpErc721FactoryAddress, dvpErc1155FactoryAddress, dvpTeleportAddr, authority_)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackVKIDAUCTIONBID is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x1727c0b0.
//
// Solidity: function VK_ID_AUCTION_BID() view returns(uint256)
func (dvp *Dvp) PackVKIDAUCTIONBID() []byte {
	enc, err := dvp.abi.Pack("VK_ID_AUCTION_BID")
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackVKIDAUCTIONBID is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x1727c0b0.
//
// Solidity: function VK_ID_AUCTION_BID() view returns(uint256)
func (dvp *Dvp) UnpackVKIDAUCTIONBID(data []byte) (*big.Int, error) {
	out, err := dvp.abi.Unpack("VK_ID_AUCTION_BID", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, err
}

// PackVKIDAUCTIONINIT is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x2afaaa94.
//
// Solidity: function VK_ID_AUCTION_INIT() view returns(uint256)
func (dvp *Dvp) PackVKIDAUCTIONINIT() []byte {
	enc, err := dvp.abi.Pack("VK_ID_AUCTION_INIT")
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackVKIDAUCTIONINIT is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x2afaaa94.
//
// Solidity: function VK_ID_AUCTION_INIT() view returns(uint256)
func (dvp *Dvp) UnpackVKIDAUCTIONINIT(data []byte) (*big.Int, error) {
	out, err := dvp.abi.Unpack("VK_ID_AUCTION_INIT", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, err
}

// PackVKIDAUCTIONNOTWINNINGBID is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x816c33f0.
//
// Solidity: function VK_ID_AUCTION_NOT_WINNING_BID() view returns(uint256)
func (dvp *Dvp) PackVKIDAUCTIONNOTWINNINGBID() []byte {
	enc, err := dvp.abi.Pack("VK_ID_AUCTION_NOT_WINNING_BID")
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackVKIDAUCTIONNOTWINNINGBID is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x816c33f0.
//
// Solidity: function VK_ID_AUCTION_NOT_WINNING_BID() view returns(uint256)
func (dvp *Dvp) UnpackVKIDAUCTIONNOTWINNINGBID(data []byte) (*big.Int, error) {
	out, err := dvp.abi.Unpack("VK_ID_AUCTION_NOT_WINNING_BID", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, err
}

// PackVKIDAUCTIONPRIVATEOPENING is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xdf657fdf.
//
// Solidity: function VK_ID_AUCTION_PRIVATE_OPENING() view returns(uint256)
func (dvp *Dvp) PackVKIDAUCTIONPRIVATEOPENING() []byte {
	enc, err := dvp.abi.Pack("VK_ID_AUCTION_PRIVATE_OPENING")
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackVKIDAUCTIONPRIVATEOPENING is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xdf657fdf.
//
// Solidity: function VK_ID_AUCTION_PRIVATE_OPENING() view returns(uint256)
func (dvp *Dvp) UnpackVKIDAUCTIONPRIVATEOPENING(data []byte) (*big.Int, error) {
	out, err := dvp.abi.Unpack("VK_ID_AUCTION_PRIVATE_OPENING", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, err
}

// PackVKIDBROKERREGISTRATION is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xab3b0387.
//
// Solidity: function VK_ID_BROKER_REGISTRATION() view returns(uint256)
func (dvp *Dvp) PackVKIDBROKERREGISTRATION() []byte {
	enc, err := dvp.abi.Pack("VK_ID_BROKER_REGISTRATION")
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackVKIDBROKERREGISTRATION is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xab3b0387.
//
// Solidity: function VK_ID_BROKER_REGISTRATION() view returns(uint256)
func (dvp *Dvp) UnpackVKIDBROKERREGISTRATION(data []byte) (*big.Int, error) {
	out, err := dvp.abi.Unpack("VK_ID_BROKER_REGISTRATION", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, err
}

// PackVKIDLEGITBROKER is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xf3d5b093.
//
// Solidity: function VK_ID_LEGIT_BROKER() view returns(uint256)
func (dvp *Dvp) PackVKIDLEGITBROKER() []byte {
	enc, err := dvp.abi.Pack("VK_ID_LEGIT_BROKER")
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackVKIDLEGITBROKER is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xf3d5b093.
//
// Solidity: function VK_ID_LEGIT_BROKER() view returns(uint256)
func (dvp *Dvp) UnpackVKIDLEGITBROKER(data []byte) (*big.Int, error) {
	out, err := dvp.abi.Unpack("VK_ID_LEGIT_BROKER", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, err
}

// PackAddEnygmaDvpIntegrationAddress is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x85e1969b.
//
// Solidity: function addEnygmaDvpIntegrationAddress(address enygmaDvpIntegrationAddress) returns(bool)
func (dvp *Dvp) PackAddEnygmaDvpIntegrationAddress(enygmaDvpIntegrationAddress common.Address) []byte {
	enc, err := dvp.abi.Pack("addEnygmaDvpIntegrationAddress", enygmaDvpIntegrationAddress)
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackAddEnygmaDvpIntegrationAddress is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x85e1969b.
//
// Solidity: function addEnygmaDvpIntegrationAddress(address enygmaDvpIntegrationAddress) returns(bool)
func (dvp *Dvp) UnpackAddEnygmaDvpIntegrationAddress(data []byte) (bool, error) {
	out, err := dvp.abi.Unpack("addEnygmaDvpIntegrationAddress", data)
	if err != nil {
		return *new(bool), err
	}
	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	return out0, err
}

// PackAddTokenToGroup is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x15b0669f.
//
// Solidity: function addTokenToGroup(uint256 vaultId, uint256[] uniqueIdParams, uint256 groupId) returns(bool)
func (dvp *Dvp) PackAddTokenToGroup(vaultId *big.Int, uniqueIdParams []*big.Int, groupId *big.Int) []byte {
	enc, err := dvp.abi.Pack("addTokenToGroup", vaultId, uniqueIdParams, groupId)
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackAddTokenToGroup is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x15b0669f.
//
// Solidity: function addTokenToGroup(uint256 vaultId, uint256[] uniqueIdParams, uint256 groupId) returns(bool)
func (dvp *Dvp) UnpackAddTokenToGroup(data []byte) (bool, error) {
	out, err := dvp.abi.Unpack("addTokenToGroup", data)
	if err != nil {
		return *new(bool), err
	}
	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	return out0, err
}

// PackAddVaultToGroup is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x4d4b9886.
//
// Solidity: function addVaultToGroup(uint256 vaultId, uint256 groupId) returns(bool)
func (dvp *Dvp) PackAddVaultToGroup(vaultId *big.Int, groupId *big.Int) []byte {
	enc, err := dvp.abi.Pack("addVaultToGroup", vaultId, groupId)
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackAddVaultToGroup is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x4d4b9886.
//
// Solidity: function addVaultToGroup(uint256 vaultId, uint256 groupId) returns(bool)
func (dvp *Dvp) UnpackAddVaultToGroup(data []byte) (bool, error) {
	out, err := dvp.abi.Unpack("addVaultToGroup", data)
	if err != nil {
		return *new(bool), err
	}
	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	return out0, err
}

// PackAuthority is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xbf7e214f.
//
// Solidity: function authority() view returns(address)
func (dvp *Dvp) PackAuthority() []byte {
	enc, err := dvp.abi.Pack("authority")
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackAuthority is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xbf7e214f.
//
// Solidity: function authority() view returns(address)
func (dvp *Dvp) UnpackAuthority(data []byte) (common.Address, error) {
	out, err := dvp.abi.Unpack("authority", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, err
}

// PackCancelSwap is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x14488fe4.
//
// Solidity: function cancelSwap(bytes32 sharedId, uint256 preimage) returns()
func (dvp *Dvp) PackCancelSwap(sharedId [32]byte, preimage *big.Int) []byte {
	enc, err := dvp.abi.Pack("cancelSwap", sharedId, preimage)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackCheckAndRegisterChallenge is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x3cd6f55e.
//
// Solidity: function checkAndRegisterChallenge(uint256 challenge_) returns(bool)
func (dvp *Dvp) PackCheckAndRegisterChallenge(challenge *big.Int) []byte {
	enc, err := dvp.abi.Pack("checkAndRegisterChallenge", challenge)
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackCheckAndRegisterChallenge is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x3cd6f55e.
//
// Solidity: function checkAndRegisterChallenge(uint256 challenge_) returns(bool)
func (dvp *Dvp) UnpackCheckAndRegisterChallenge(data []byte) (bool, error) {
	out, err := dvp.abi.Unpack("checkAndRegisterChallenge", data)
	if err != nil {
		return *new(bool), err
	}
	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	return out0, err
}

// PackCompleteSwap is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x13295898.
//
// Solidity: function completeSwap(bytes32 sharedId, address tokenAddress, uint8 proofType, (((uint256,uint256),(uint256[2],uint256[2]),(uint256,uint256)),uint256[],uint256,uint256[],uint256[],uint256[],uint256) proof, bytes encryptedData) returns()
func (dvp *Dvp) PackCompleteSwap(sharedId [32]byte, tokenAddress common.Address, proofType uint8, proof IDvpProofReceipt, encryptedData []byte) []byte {
	enc, err := dvp.abi.Pack("completeSwap", sharedId, tokenAddress, proofType, proof, encryptedData)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackDepositERC1155 is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x78eefeec.
//
// Solidity: function depositERC1155(address contractAddress, uint256 tokenId, uint256 amountOrOne, bytes data, uint256 publicKey, uint256 salt, bytes encryptedBalanceUpdate) returns(bool)
func (dvp *Dvp) PackDepositERC1155(contractAddress common.Address, tokenId *big.Int, amountOrOne *big.Int, data []byte, publicKey *big.Int, salt *big.Int, encryptedBalanceUpdate []byte) []byte {
	enc, err := dvp.abi.Pack("depositERC1155", contractAddress, tokenId, amountOrOne, data, publicKey, salt, encryptedBalanceUpdate)
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackDepositERC1155 is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x78eefeec.
//
// Solidity: function depositERC1155(address contractAddress, uint256 tokenId, uint256 amountOrOne, bytes data, uint256 publicKey, uint256 salt, bytes encryptedBalanceUpdate) returns(bool)
func (dvp *Dvp) UnpackDepositERC1155(data []byte) (bool, error) {
	out, err := dvp.abi.Unpack("depositERC1155", data)
	if err != nil {
		return *new(bool), err
	}
	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	return out0, err
}

// PackDepositERC721 is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x010f8127.
//
// Solidity: function depositERC721(address contractAddress, uint256 nftId, uint256 publicKey, uint256 salt, bytes encryptedBalanceUpdate) returns(bool)
func (dvp *Dvp) PackDepositERC721(contractAddress common.Address, nftId *big.Int, publicKey *big.Int, salt *big.Int, encryptedBalanceUpdate []byte) []byte {
	enc, err := dvp.abi.Pack("depositERC721", contractAddress, nftId, publicKey, salt, encryptedBalanceUpdate)
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackDepositERC721 is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x010f8127.
//
// Solidity: function depositERC721(address contractAddress, uint256 nftId, uint256 publicKey, uint256 salt, bytes encryptedBalanceUpdate) returns(bool)
func (dvp *Dvp) UnpackDepositERC721(data []byte) (bool, error) {
	out, err := dvp.abi.Unpack("depositERC721", data)
	if err != nil {
		return *new(bool), err
	}
	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	return out0, err
}

// PackDepositEnygma is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x1a953191.
//
// Solidity: function depositEnygma(uint256 vaultId, uint256 hashCommitment) returns(bool)
func (dvp *Dvp) PackDepositEnygma(vaultId *big.Int, hashCommitment *big.Int) []byte {
	enc, err := dvp.abi.Pack("depositEnygma", vaultId, hashCommitment)
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackDepositEnygma is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x1a953191.
//
// Solidity: function depositEnygma(uint256 vaultId, uint256 hashCommitment) returns(bool)
func (dvp *Dvp) UnpackDepositEnygma(data []byte) (bool, error) {
	out, err := dvp.abi.Unpack("depositEnygma", data)
	if err != nil {
		return *new(bool), err
	}
	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	return out0, err
}

// PackDvpTeleportAddress is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x9ccafcdc.
//
// Solidity: function dvpTeleportAddress() view returns(address)
func (dvp *Dvp) PackDvpTeleportAddress() []byte {
	enc, err := dvp.abi.Pack("dvpTeleportAddress")
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackDvpTeleportAddress is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x9ccafcdc.
//
// Solidity: function dvpTeleportAddress() view returns(address)
func (dvp *Dvp) UnpackDvpTeleportAddress(data []byte) (common.Address, error) {
	out, err := dvp.abi.Unpack("dvpTeleportAddress", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, err
}

// PackExpireSwap is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x48b1ee1f.
//
// Solidity: function expireSwap(bytes32 sharedId) returns()
func (dvp *Dvp) PackExpireSwap(sharedId [32]byte) []byte {
	enc, err := dvp.abi.Pack("expireSwap", sharedId)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackGetGroupPairId is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x8b016cf9.
//
// Solidity: function getGroupPairId(uint256 groupId1, uint256 groupId2) pure returns(uint256)
func (dvp *Dvp) PackGetGroupPairId(groupId1 *big.Int, groupId2 *big.Int) []byte {
	enc, err := dvp.abi.Pack("getGroupPairId", groupId1, groupId2)
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackGetGroupPairId is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x8b016cf9.
//
// Solidity: function getGroupPairId(uint256 groupId1, uint256 groupId2) pure returns(uint256)
func (dvp *Dvp) UnpackGetGroupPairId(data []byte) (*big.Int, error) {
	out, err := dvp.abi.Unpack("getGroupPairId", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, err
}

// PackGetVaultIdByAddress is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xd195a53f.
//
// Solidity: function getVaultIdByAddress(address contractAddress) view returns(uint256)
func (dvp *Dvp) PackGetVaultIdByAddress(contractAddress common.Address) []byte {
	enc, err := dvp.abi.Pack("getVaultIdByAddress", contractAddress)
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackGetVaultIdByAddress is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xd195a53f.
//
// Solidity: function getVaultIdByAddress(address contractAddress) view returns(uint256)
func (dvp *Dvp) UnpackGetVaultIdByAddress(data []byte) (*big.Int, error) {
	out, err := dvp.abi.Unpack("getVaultIdByAddress", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, err
}

// PackHashContractAddress is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x008f47d1.
//
// Solidity: function hashContractAddress() view returns(address)
func (dvp *Dvp) PackHashContractAddress() []byte {
	enc, err := dvp.abi.Pack("hashContractAddress")
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackHashContractAddress is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x008f47d1.
//
// Solidity: function hashContractAddress() view returns(address)
func (dvp *Dvp) UnpackHashContractAddress(data []byte) (common.Address, error) {
	out, err := dvp.abi.Unpack("hashContractAddress", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, err
}

// PackInitializeDvp is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xc7b50db0.
//
// Solidity: function initializeDvp(address verifierAddress) returns(bool)
func (dvp *Dvp) PackInitializeDvp(verifierAddress common.Address) []byte {
	enc, err := dvp.abi.Pack("initializeDvp", verifierAddress)
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackInitializeDvp is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xc7b50db0.
//
// Solidity: function initializeDvp(address verifierAddress) returns(bool)
func (dvp *Dvp) UnpackInitializeDvp(data []byte) (bool, error) {
	out, err := dvp.abi.Unpack("initializeDvp", data)
	if err != nil {
		return *new(bool), err
	}
	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	return out0, err
}

// PackInitiateSwap is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x5964460f.
//
// Solidity: function initiateSwap(bytes32 sharedId, bytes encryptedData, bytes ciphertext, address tokenAddress, uint8 proofType, (((uint256,uint256),(uint256[2],uint256[2]),(uint256,uint256)),uint256[],uint256,uint256[],uint256[],uint256[],uint256) proof, uint64 validityTime, uint256 passphrase) returns(bytes32 dvpId)
func (dvp *Dvp) PackInitiateSwap(sharedId [32]byte, encryptedData []byte, ciphertext []byte, tokenAddress common.Address, proofType uint8, proof IDvpProofReceipt, validityTime uint64, passphrase *big.Int) []byte {
	enc, err := dvp.abi.Pack("initiateSwap", sharedId, encryptedData, ciphertext, tokenAddress, proofType, proof, validityTime, passphrase)
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackInitiateSwap is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x5964460f.
//
// Solidity: function initiateSwap(bytes32 sharedId, bytes encryptedData, bytes ciphertext, address tokenAddress, uint8 proofType, (((uint256,uint256),(uint256[2],uint256[2]),(uint256,uint256)),uint256[],uint256,uint256[],uint256[],uint256[],uint256) proof, uint64 validityTime, uint256 passphrase) returns(bytes32 dvpId)
func (dvp *Dvp) UnpackInitiateSwap(data []byte) ([32]byte, error) {
	out, err := dvp.abi.Unpack("initiateSwap", data)
	if err != nil {
		return *new([32]byte), err
	}
	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)
	return out0, err
}

// PackIsAuditorRegistered is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x8ee9cb45.
//
// Solidity: function isAuditorRegistered(uint256[2] publicKey) view returns(bool)
func (dvp *Dvp) PackIsAuditorRegistered(publicKey [2]*big.Int) []byte {
	enc, err := dvp.abi.Pack("isAuditorRegistered", publicKey)
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackIsAuditorRegistered is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x8ee9cb45.
//
// Solidity: function isAuditorRegistered(uint256[2] publicKey) view returns(bool)
func (dvp *Dvp) UnpackIsAuditorRegistered(data []byte) (bool, error) {
	out, err := dvp.abi.Unpack("isAuditorRegistered", data)
	if err != nil {
		return *new(bool), err
	}
	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	return out0, err
}

// PackIsMemberOfFromProofReceipt is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xef07a7bf.
//
// Solidity: function isMemberOfFromProofReceipt(uint256 vaultId, (((uint256,uint256),(uint256[2],uint256[2]),(uint256,uint256)),uint256[],uint256,uint256[],uint256[],uint256[],uint256) receipt, uint256 groupId) view returns(bool)
func (dvp *Dvp) PackIsMemberOfFromProofReceipt(vaultId *big.Int, receipt IDvpProofReceipt, groupId *big.Int) []byte {
	enc, err := dvp.abi.Pack("isMemberOfFromProofReceipt", vaultId, receipt, groupId)
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackIsMemberOfFromProofReceipt is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xef07a7bf.
//
// Solidity: function isMemberOfFromProofReceipt(uint256 vaultId, (((uint256,uint256),(uint256[2],uint256[2]),(uint256,uint256)),uint256[],uint256,uint256[],uint256[],uint256[],uint256) receipt, uint256 groupId) view returns(bool)
func (dvp *Dvp) UnpackIsMemberOfFromProofReceipt(data []byte) (bool, error) {
	out, err := dvp.abi.Unpack("isMemberOfFromProofReceipt", data)
	if err != nil {
		return *new(bool), err
	}
	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	return out0, err
}

// PackIsSwapExpired is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xe9c0c82f.
//
// Solidity: function isSwapExpired(bytes32 sharedId) view returns(bool)
func (dvp *Dvp) PackIsSwapExpired(sharedId [32]byte) []byte {
	enc, err := dvp.abi.Pack("isSwapExpired", sharedId)
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackIsSwapExpired is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xe9c0c82f.
//
// Solidity: function isSwapExpired(bytes32 sharedId) view returns(bool)
func (dvp *Dvp) UnpackIsSwapExpired(data []byte) (bool, error) {
	out, err := dvp.abi.Unpack("isSwapExpired", data)
	if err != nil {
		return *new(bool), err
	}
	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	return out0, err
}

// PackIsTokenMemberOf is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x961cafa8.
//
// Solidity: function isTokenMemberOf(uint256 vaultId, uint256[] uniqueIdParams, uint256 groupId) view returns(bool)
func (dvp *Dvp) PackIsTokenMemberOf(vaultId *big.Int, uniqueIdParams []*big.Int, groupId *big.Int) []byte {
	enc, err := dvp.abi.Pack("isTokenMemberOf", vaultId, uniqueIdParams, groupId)
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackIsTokenMemberOf is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x961cafa8.
//
// Solidity: function isTokenMemberOf(uint256 vaultId, uint256[] uniqueIdParams, uint256 groupId) view returns(bool)
func (dvp *Dvp) UnpackIsTokenMemberOf(data []byte) (bool, error) {
	out, err := dvp.abi.Unpack("isTokenMemberOf", data)
	if err != nil {
		return *new(bool), err
	}
	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	return out0, err
}

// PackIsValidExchangeGroupPair is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xb02047d8.
//
// Solidity: function isValidExchangeGroupPair(uint256 groupId1, uint256 groupId2) view returns(bool)
func (dvp *Dvp) PackIsValidExchangeGroupPair(groupId1 *big.Int, groupId2 *big.Int) []byte {
	enc, err := dvp.abi.Pack("isValidExchangeGroupPair", groupId1, groupId2)
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackIsValidExchangeGroupPair is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xb02047d8.
//
// Solidity: function isValidExchangeGroupPair(uint256 groupId1, uint256 groupId2) view returns(bool)
func (dvp *Dvp) UnpackIsValidExchangeGroupPair(data []byte) (bool, error) {
	out, err := dvp.abi.Unpack("isValidExchangeGroupPair", data)
	if err != nil {
		return *new(bool), err
	}
	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	return out0, err
}

// PackIsValidSwapGroupPair is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x9f42f00d.
//
// Solidity: function isValidSwapGroupPair(uint256 groupId1, uint256 groupId2) view returns(bool)
func (dvp *Dvp) PackIsValidSwapGroupPair(groupId1 *big.Int, groupId2 *big.Int) []byte {
	enc, err := dvp.abi.Pack("isValidSwapGroupPair", groupId1, groupId2)
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackIsValidSwapGroupPair is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x9f42f00d.
//
// Solidity: function isValidSwapGroupPair(uint256 groupId1, uint256 groupId2) view returns(bool)
func (dvp *Dvp) UnpackIsValidSwapGroupPair(data []byte) (bool, error) {
	out, err := dvp.abi.Unpack("isValidSwapGroupPair", data)
	if err != nil {
		return *new(bool), err
	}
	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	return out0, err
}

// PackIsVaultMemberOf is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x60f6b132.
//
// Solidity: function isVaultMemberOf(uint256 vaultId, uint256 groupId) view returns(bool)
func (dvp *Dvp) PackIsVaultMemberOf(vaultId *big.Int, groupId *big.Int) []byte {
	enc, err := dvp.abi.Pack("isVaultMemberOf", vaultId, groupId)
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackIsVaultMemberOf is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x60f6b132.
//
// Solidity: function isVaultMemberOf(uint256 vaultId, uint256 groupId) view returns(bool)
func (dvp *Dvp) UnpackIsVaultMemberOf(data []byte) (bool, error) {
	out, err := dvp.abi.Unpack("isVaultMemberOf", data)
	if err != nil {
		return *new(bool), err
	}
	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	return out0, err
}

// PackMixFunds is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x7997ad5b.
//
// Solidity: function mixFunds(uint256 vaultId, (((uint256,uint256),(uint256[2],uint256[2]),(uint256,uint256)),uint256[],uint256,uint256[],uint256[],uint256[],uint256) _tx) returns(bool)
func (dvp *Dvp) PackMixFunds(vaultId *big.Int, tx IDvpProofReceipt) []byte {
	enc, err := dvp.abi.Pack("mixFunds", vaultId, tx)
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackMixFunds is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x7997ad5b.
//
// Solidity: function mixFunds(uint256 vaultId, (((uint256,uint256),(uint256[2],uint256[2]),(uint256,uint256)),uint256[],uint256,uint256[],uint256[],uint256[],uint256) _tx) returns(bool)
func (dvp *Dvp) UnpackMixFunds(data []byte) (bool, error) {
	out, err := dvp.abi.Unpack("mixFunds", data)
	if err != nil {
		return *new(bool), err
	}
	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	return out0, err
}

// PackMixFundsERC1155 is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x2c56c234.
//
// Solidity: function mixFundsERC1155(address contractAddress, (((uint256,uint256),(uint256[2],uint256[2]),(uint256,uint256)),uint256[],uint256,uint256[],uint256[],uint256[],uint256) _tx) returns(bool)
func (dvp *Dvp) PackMixFundsERC1155(contractAddress common.Address, tx IDvpProofReceipt) []byte {
	enc, err := dvp.abi.Pack("mixFundsERC1155", contractAddress, tx)
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackMixFundsERC1155 is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x2c56c234.
//
// Solidity: function mixFundsERC1155(address contractAddress, (((uint256,uint256),(uint256[2],uint256[2]),(uint256,uint256)),uint256[],uint256,uint256[],uint256[],uint256[],uint256) _tx) returns(bool)
func (dvp *Dvp) UnpackMixFundsERC1155(data []byte) (bool, error) {
	out, err := dvp.abi.Unpack("mixFundsERC1155", data)
	if err != nil {
		return *new(bool), err
	}
	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	return out0, err
}

// PackName is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (dvp *Dvp) PackName() []byte {
	enc, err := dvp.abi.Pack("name")
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackName is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (dvp *Dvp) UnpackName(data []byte) (string, error) {
	out, err := dvp.abi.Unpack("name", data)
	if err != nil {
		return *new(string), err
	}
	out0 := *abi.ConvertType(out[0], new(string)).(*string)
	return out0, err
}

// PackRegisterAssetGroup is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x7c4bdb0e.
//
// Solidity: function registerAssetGroup(address assetGroupContractAddress, string assetGroupName, bool isAssetGroupFungible, uint256 treeDepth) returns(bool)
func (dvp *Dvp) PackRegisterAssetGroup(assetGroupContractAddress common.Address, assetGroupName string, isAssetGroupFungible bool, treeDepth *big.Int) []byte {
	enc, err := dvp.abi.Pack("registerAssetGroup", assetGroupContractAddress, assetGroupName, isAssetGroupFungible, treeDepth)
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackRegisterAssetGroup is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x7c4bdb0e.
//
// Solidity: function registerAssetGroup(address assetGroupContractAddress, string assetGroupName, bool isAssetGroupFungible, uint256 treeDepth) returns(bool)
func (dvp *Dvp) UnpackRegisterAssetGroup(data []byte) (bool, error) {
	out, err := dvp.abi.Unpack("registerAssetGroup", data)
	if err != nil {
		return *new(bool), err
	}
	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	return out0, err
}

// PackRegisterAuditor is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x977701cf.
//
// Solidity: function registerAuditor(uint256 auditorOffchainId, uint256 auditorGroupId, uint256[2] auditorPublicKey) returns(bool)
func (dvp *Dvp) PackRegisterAuditor(auditorOffchainId *big.Int, auditorGroupId *big.Int, auditorPublicKey [2]*big.Int) []byte {
	enc, err := dvp.abi.Pack("registerAuditor", auditorOffchainId, auditorGroupId, auditorPublicKey)
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackRegisterAuditor is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x977701cf.
//
// Solidity: function registerAuditor(uint256 auditorOffchainId, uint256 auditorGroupId, uint256[2] auditorPublicKey) returns(bool)
func (dvp *Dvp) UnpackRegisterAuditor(data []byte) (bool, error) {
	out, err := dvp.abi.Unpack("registerAuditor", data)
	if err != nil {
		return *new(bool), err
	}
	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	return out0, err
}

// PackRegisterBroker is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xc7fc0984.
//
// Solidity: function registerBroker((((uint256,uint256),(uint256[2],uint256[2]),(uint256,uint256)),uint256[],uint256,uint256[],uint256[],uint256[],uint256) brokerRegistrationProof) returns(bool)
func (dvp *Dvp) PackRegisterBroker(brokerRegistrationProof IDvpProofReceipt) []byte {
	enc, err := dvp.abi.Pack("registerBroker", brokerRegistrationProof)
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackRegisterBroker is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xc7fc0984.
//
// Solidity: function registerBroker((((uint256,uint256),(uint256[2],uint256[2]),(uint256,uint256)),uint256[],uint256,uint256[],uint256[],uint256[],uint256) brokerRegistrationProof) returns(bool)
func (dvp *Dvp) UnpackRegisterBroker(data []byte) (bool, error) {
	out, err := dvp.abi.Unpack("registerBroker", data)
	if err != nil {
		return *new(bool), err
	}
	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	return out0, err
}

// PackRegisterExchangeGroupPair is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x94f1bd07.
//
// Solidity: function registerExchangeGroupPair(uint256 groupId1, uint256 groupId2) returns(bool)
func (dvp *Dvp) PackRegisterExchangeGroupPair(groupId1 *big.Int, groupId2 *big.Int) []byte {
	enc, err := dvp.abi.Pack("registerExchangeGroupPair", groupId1, groupId2)
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackRegisterExchangeGroupPair is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x94f1bd07.
//
// Solidity: function registerExchangeGroupPair(uint256 groupId1, uint256 groupId2) returns(bool)
func (dvp *Dvp) UnpackRegisterExchangeGroupPair(data []byte) (bool, error) {
	out, err := dvp.abi.Unpack("registerExchangeGroupPair", data)
	if err != nil {
		return *new(bool), err
	}
	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	return out0, err
}

// PackRegisterSwapGroupPair is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xa30886bf.
//
// Solidity: function registerSwapGroupPair(uint256 groupId1, uint256 groupId2) returns(bool)
func (dvp *Dvp) PackRegisterSwapGroupPair(groupId1 *big.Int, groupId2 *big.Int) []byte {
	enc, err := dvp.abi.Pack("registerSwapGroupPair", groupId1, groupId2)
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackRegisterSwapGroupPair is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xa30886bf.
//
// Solidity: function registerSwapGroupPair(uint256 groupId1, uint256 groupId2) returns(bool)
func (dvp *Dvp) UnpackRegisterSwapGroupPair(data []byte) (bool, error) {
	out, err := dvp.abi.Unpack("registerSwapGroupPair", data)
	if err != nil {
		return *new(bool), err
	}
	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	return out0, err
}

// PackRegisterVault is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x04473f27.
//
// Solidity: function registerVault(address vaultContractAddress, address assetContractAddress, uint256 treeDepth) returns(uint256)
func (dvp *Dvp) PackRegisterVault(vaultContractAddress common.Address, assetContractAddress common.Address, treeDepth *big.Int) []byte {
	enc, err := dvp.abi.Pack("registerVault", vaultContractAddress, assetContractAddress, treeDepth)
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackRegisterVault is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x04473f27.
//
// Solidity: function registerVault(address vaultContractAddress, address assetContractAddress, uint256 treeDepth) returns(uint256)
func (dvp *Dvp) UnpackRegisterVault(data []byte) (*big.Int, error) {
	out, err := dvp.abi.Unpack("registerVault", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, err
}

// PackRegisterZkAuction is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xc9274021.
//
// Solidity: function registerZkAuction(address zkAuctionContractAddress) returns(bool)
func (dvp *Dvp) PackRegisterZkAuction(zkAuctionContractAddress common.Address) []byte {
	enc, err := dvp.abi.Pack("registerZkAuction", zkAuctionContractAddress)
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackRegisterZkAuction is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xc9274021.
//
// Solidity: function registerZkAuction(address zkAuctionContractAddress) returns(bool)
func (dvp *Dvp) UnpackRegisterZkAuction(data []byte) (bool, error) {
	out, err := dvp.abi.Unpack("registerZkAuction", data)
	if err != nil {
		return *new(bool), err
	}
	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	return out0, err
}

// PackUnregisterAuditor is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xb8b249c6.
//
// Solidity: function unregisterAuditor(uint256 auditorOnchainId) returns(bool)
func (dvp *Dvp) PackUnregisterAuditor(auditorOnchainId *big.Int) []byte {
	enc, err := dvp.abi.Pack("unregisterAuditor", auditorOnchainId)
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackUnregisterAuditor is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xb8b249c6.
//
// Solidity: function unregisterAuditor(uint256 auditorOnchainId) returns(bool)
func (dvp *Dvp) UnpackUnregisterAuditor(data []byte) (bool, error) {
	out, err := dvp.abi.Unpack("unregisterAuditor", data)
	if err != nil {
		return *new(bool), err
	}
	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	return out0, err
}

// PackVaultById is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x32f19a5c.
//
// Solidity: function vaultById(uint256 vaultId) view returns(address)
func (dvp *Dvp) PackVaultById(vaultId *big.Int) []byte {
	enc, err := dvp.abi.Pack("vaultById", vaultId)
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackVaultById is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x32f19a5c.
//
// Solidity: function vaultById(uint256 vaultId) view returns(address)
func (dvp *Dvp) UnpackVaultById(data []byte) (common.Address, error) {
	out, err := dvp.abi.Unpack("vaultById", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, err
}

// PackVerifierContractAddress is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xf588b157.
//
// Solidity: function verifierContractAddress() view returns(address)
func (dvp *Dvp) PackVerifierContractAddress() []byte {
	enc, err := dvp.abi.Pack("verifierContractAddress")
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackVerifierContractAddress is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xf588b157.
//
// Solidity: function verifierContractAddress() view returns(address)
func (dvp *Dvp) UnpackVerifierContractAddress(data []byte) (common.Address, error) {
	out, err := dvp.abi.Unpack("verifierContractAddress", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, err
}

// PackVerifyLegitBrokerReceipt is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x6c9d68b7.
//
// Solidity: function verifyLegitBrokerReceipt((((uint256,uint256),(uint256[2],uint256[2]),(uint256,uint256)),uint256[],uint256,uint256[],uint256[],uint256[],uint256) receipt) returns(bool)
func (dvp *Dvp) PackVerifyLegitBrokerReceipt(receipt IDvpProofReceipt) []byte {
	enc, err := dvp.abi.Pack("verifyLegitBrokerReceipt", receipt)
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackVerifyLegitBrokerReceipt is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x6c9d68b7.
//
// Solidity: function verifyLegitBrokerReceipt((((uint256,uint256),(uint256[2],uint256[2]),(uint256,uint256)),uint256[],uint256,uint256[],uint256[],uint256[],uint256) receipt) returns(bool)
func (dvp *Dvp) UnpackVerifyLegitBrokerReceipt(data []byte) (bool, error) {
	out, err := dvp.abi.Unpack("verifyLegitBrokerReceipt", data)
	if err != nil {
		return *new(bool), err
	}
	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	return out0, err
}

// PackWithdrawERC1155 is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xbd9098e7.
//
// Solidity: function withdrawERC1155(address contractAddress, uint256 tokenId, uint256 amount, address recipient, uint256 salt, (((uint256,uint256),(uint256[2],uint256[2]),(uint256,uint256)),uint256[],uint256,uint256[],uint256[],uint256[],uint256) proofTx, bytes encryptedBalanceUpdate) returns(bool)
func (dvp *Dvp) PackWithdrawERC1155(contractAddress common.Address, tokenId *big.Int, amount *big.Int, recipient common.Address, salt *big.Int, proofTx IDvpProofReceipt, encryptedBalanceUpdate []byte) []byte {
	enc, err := dvp.abi.Pack("withdrawERC1155", contractAddress, tokenId, amount, recipient, salt, proofTx, encryptedBalanceUpdate)
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackWithdrawERC1155 is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xbd9098e7.
//
// Solidity: function withdrawERC1155(address contractAddress, uint256 tokenId, uint256 amount, address recipient, uint256 salt, (((uint256,uint256),(uint256[2],uint256[2]),(uint256,uint256)),uint256[],uint256,uint256[],uint256[],uint256[],uint256) proofTx, bytes encryptedBalanceUpdate) returns(bool)
func (dvp *Dvp) UnpackWithdrawERC1155(data []byte) (bool, error) {
	out, err := dvp.abi.Unpack("withdrawERC1155", data)
	if err != nil {
		return *new(bool), err
	}
	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	return out0, err
}

// PackWithdrawERC721 is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x5f45ea54.
//
// Solidity: function withdrawERC721(address contractAddress, uint256 nftId, address recipient, uint256 salt, (((uint256,uint256),(uint256[2],uint256[2]),(uint256,uint256)),uint256[],uint256,uint256[],uint256[],uint256[],uint256) proofTx, bytes encryptedBalanceUpdate) returns(bool)
func (dvp *Dvp) PackWithdrawERC721(contractAddress common.Address, nftId *big.Int, recipient common.Address, salt *big.Int, proofTx IDvpProofReceipt, encryptedBalanceUpdate []byte) []byte {
	enc, err := dvp.abi.Pack("withdrawERC721", contractAddress, nftId, recipient, salt, proofTx, encryptedBalanceUpdate)
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackWithdrawERC721 is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x5f45ea54.
//
// Solidity: function withdrawERC721(address contractAddress, uint256 nftId, address recipient, uint256 salt, (((uint256,uint256),(uint256[2],uint256[2]),(uint256,uint256)),uint256[],uint256,uint256[],uint256[],uint256[],uint256) proofTx, bytes encryptedBalanceUpdate) returns(bool)
func (dvp *Dvp) UnpackWithdrawERC721(data []byte) (bool, error) {
	out, err := dvp.abi.Unpack("withdrawERC721", data)
	if err != nil {
		return *new(bool), err
	}
	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	return out0, err
}

// PackWithdrawEnygma is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x8a622d49.
//
// Solidity: function withdrawEnygma(uint256 vaultId, (((uint256,uint256),(uint256[2],uint256[2]),(uint256,uint256)),uint256[],uint256,uint256[],uint256[],uint256[],uint256) _tx) returns(bool)
func (dvp *Dvp) PackWithdrawEnygma(vaultId *big.Int, tx IDvpProofReceipt) []byte {
	enc, err := dvp.abi.Pack("withdrawEnygma", vaultId, tx)
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackWithdrawEnygma is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x8a622d49.
//
// Solidity: function withdrawEnygma(uint256 vaultId, (((uint256,uint256),(uint256[2],uint256[2]),(uint256,uint256)),uint256[],uint256,uint256[],uint256[],uint256[],uint256) _tx) returns(bool)
func (dvp *Dvp) UnpackWithdrawEnygma(data []byte) (bool, error) {
	out, err := dvp.abi.Unpack("withdrawEnygma", data)
	if err != nil {
		return *new(bool), err
	}
	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	return out0, err
}

// DvpAuditorRegistered represents a AuditorRegistered event raised by the Dvp contract.
type DvpAuditorRegistered struct {
	OnchainId  *big.Int
	OffchainId *big.Int
	GroupId    *big.Int
	PublicKey  [2]*big.Int
	Raw        *types.Log // Blockchain specific contextual infos
}

const DvpAuditorRegisteredEventName = "AuditorRegistered"

// ContractEventName returns the user-defined event name.
func (DvpAuditorRegistered) ContractEventName() string {
	return DvpAuditorRegisteredEventName
}

// UnpackAuditorRegisteredEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event AuditorRegistered(uint256 indexed onchainId, uint256 indexed offchainId, uint256 indexed groupId, uint256[2] publicKey)
func (dvp *Dvp) UnpackAuditorRegisteredEvent(log *types.Log) (*DvpAuditorRegistered, error) {
	event := "AuditorRegistered"
	if log.Topics[0] != dvp.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(DvpAuditorRegistered)
	if len(log.Data) > 0 {
		if err := dvp.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range dvp.abi.Events[event].Inputs {
		if arg.Indexed {
			indexed = append(indexed, arg)
		}
	}
	if err := abi.ParseTopics(out, indexed, log.Topics[1:]); err != nil {
		return nil, err
	}
	out.Raw = log
	return out, nil
}

// DvpAuditorUnregistered represents a AuditorUnregistered event raised by the Dvp contract.
type DvpAuditorUnregistered struct {
	OnchainId *big.Int
	Raw       *types.Log // Blockchain specific contextual infos
}

const DvpAuditorUnregisteredEventName = "AuditorUnregistered"

// ContractEventName returns the user-defined event name.
func (DvpAuditorUnregistered) ContractEventName() string {
	return DvpAuditorUnregisteredEventName
}

// UnpackAuditorUnregisteredEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event AuditorUnregistered(uint256 indexed onchainId)
func (dvp *Dvp) UnpackAuditorUnregisteredEvent(log *types.Log) (*DvpAuditorUnregistered, error) {
	event := "AuditorUnregistered"
	if log.Topics[0] != dvp.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(DvpAuditorUnregistered)
	if len(log.Data) > 0 {
		if err := dvp.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range dvp.abi.Events[event].Inputs {
		if arg.Indexed {
			indexed = append(indexed, arg)
		}
	}
	if err := abi.ParseTopics(out, indexed, log.Topics[1:]); err != nil {
		return nil, err
	}
	out.Raw = log
	return out, nil
}

// DvpAuthorityUpdated represents a AuthorityUpdated event raised by the Dvp contract.
type DvpAuthorityUpdated struct {
	OldAuthority common.Address
	NewAuthority common.Address
	Raw          *types.Log // Blockchain specific contextual infos
}

const DvpAuthorityUpdatedEventName = "AuthorityUpdated"

// ContractEventName returns the user-defined event name.
func (DvpAuthorityUpdated) ContractEventName() string {
	return DvpAuthorityUpdatedEventName
}

// UnpackAuthorityUpdatedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event AuthorityUpdated(address indexed oldAuthority, address indexed newAuthority)
func (dvp *Dvp) UnpackAuthorityUpdatedEvent(log *types.Log) (*DvpAuthorityUpdated, error) {
	event := "AuthorityUpdated"
	if log.Topics[0] != dvp.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(DvpAuthorityUpdated)
	if len(log.Data) > 0 {
		if err := dvp.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range dvp.abi.Events[event].Inputs {
		if arg.Indexed {
			indexed = append(indexed, arg)
		}
	}
	if err := abi.ParseTopics(out, indexed, log.Topics[1:]); err != nil {
		return nil, err
	}
	out.Raw = log
	return out, nil
}

// DvpBrokerRegistered represents a BrokerRegistered event raised by the Dvp contract.
type DvpBrokerRegistered struct {
	VaultId                *big.Int
	BlindedBrokerPublicKey *big.Int
	Raw                    *types.Log // Blockchain specific contextual infos
}

const DvpBrokerRegisteredEventName = "BrokerRegistered"

// ContractEventName returns the user-defined event name.
func (DvpBrokerRegistered) ContractEventName() string {
	return DvpBrokerRegisteredEventName
}

// UnpackBrokerRegisteredEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event BrokerRegistered(uint256 indexed vaultId, uint256 indexed blindedBrokerPublicKey)
func (dvp *Dvp) UnpackBrokerRegisteredEvent(log *types.Log) (*DvpBrokerRegistered, error) {
	event := "BrokerRegistered"
	if log.Topics[0] != dvp.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(DvpBrokerRegistered)
	if len(log.Data) > 0 {
		if err := dvp.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range dvp.abi.Events[event].Inputs {
		if arg.Indexed {
			indexed = append(indexed, arg)
		}
	}
	if err := abi.ParseTopics(out, indexed, log.Topics[1:]); err != nil {
		return nil, err
	}
	out.Raw = log
	return out, nil
}

// DvpCoinLocked represents a CoinLocked event raised by the Dvp contract.
type DvpCoinLocked struct {
	AssetId    *big.Int
	TreeNumber *big.Int
	Nullifier  *big.Int
	Raw        *types.Log // Blockchain specific contextual infos
}

const DvpCoinLockedEventName = "CoinLocked"

// ContractEventName returns the user-defined event name.
func (DvpCoinLocked) ContractEventName() string {
	return DvpCoinLockedEventName
}

// UnpackCoinLockedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event CoinLocked(uint256 indexed assetId, uint256 indexed treeNumber, uint256 indexed nullifier)
func (dvp *Dvp) UnpackCoinLockedEvent(log *types.Log) (*DvpCoinLocked, error) {
	event := "CoinLocked"
	if log.Topics[0] != dvp.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(DvpCoinLocked)
	if len(log.Data) > 0 {
		if err := dvp.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range dvp.abi.Events[event].Inputs {
		if arg.Indexed {
			indexed = append(indexed, arg)
		}
	}
	if err := abi.ParseTopics(out, indexed, log.Topics[1:]); err != nil {
		return nil, err
	}
	out.Raw = log
	return out, nil
}

// DvpCoinUnlocked represents a CoinUnlocked event raised by the Dvp contract.
type DvpCoinUnlocked struct {
	AssetId    *big.Int
	TreeNumber *big.Int
	Nullifier  *big.Int
	Raw        *types.Log // Blockchain specific contextual infos
}

const DvpCoinUnlockedEventName = "CoinUnlocked"

// ContractEventName returns the user-defined event name.
func (DvpCoinUnlocked) ContractEventName() string {
	return DvpCoinUnlockedEventName
}

// UnpackCoinUnlockedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event CoinUnlocked(uint256 indexed assetId, uint256 indexed treeNumber, uint256 indexed nullifier)
func (dvp *Dvp) UnpackCoinUnlockedEvent(log *types.Log) (*DvpCoinUnlocked, error) {
	event := "CoinUnlocked"
	if log.Topics[0] != dvp.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(DvpCoinUnlocked)
	if len(log.Data) > 0 {
		if err := dvp.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range dvp.abi.Events[event].Inputs {
		if arg.Indexed {
			indexed = append(indexed, arg)
		}
	}
	if err := abi.ParseTopics(out, indexed, log.Topics[1:]); err != nil {
		return nil, err
	}
	out.Raw = log
	return out, nil
}

// DvpLegitBrokerReceipt represents a LegitBrokerReceipt event raised by the Dvp contract.
type DvpLegitBrokerReceipt struct {
	Beacon                 *big.Int
	BlindedBrokerPublicKey *big.Int
	Raw                    *types.Log // Blockchain specific contextual infos
}

const DvpLegitBrokerReceiptEventName = "LegitBrokerReceipt"

// ContractEventName returns the user-defined event name.
func (DvpLegitBrokerReceipt) ContractEventName() string {
	return DvpLegitBrokerReceiptEventName
}

// UnpackLegitBrokerReceiptEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event LegitBrokerReceipt(uint256 indexed beacon, uint256 indexed blindedBrokerPublicKey)
func (dvp *Dvp) UnpackLegitBrokerReceiptEvent(log *types.Log) (*DvpLegitBrokerReceipt, error) {
	event := "LegitBrokerReceipt"
	if log.Topics[0] != dvp.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(DvpLegitBrokerReceipt)
	if len(log.Data) > 0 {
		if err := dvp.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range dvp.abi.Events[event].Inputs {
		if arg.Indexed {
			indexed = append(indexed, arg)
		}
	}
	if err := abi.ParseTopics(out, indexed, log.Topics[1:]); err != nil {
		return nil, err
	}
	out.Raw = log
	return out, nil
}

// DvpPendingProofAddedToVault represents a PendingProofAddedToVault event raised by the Dvp contract.
type DvpPendingProofAddedToVault struct {
	VaultId         *big.Int
	GroupId         *big.Int
	TargetReceiptId *big.Int
	PendingProof    IDvpProofReceipt
	Raw             *types.Log // Blockchain specific contextual infos
}

const DvpPendingProofAddedToVaultEventName = "PendingProofAddedToVault"

// ContractEventName returns the user-defined event name.
func (DvpPendingProofAddedToVault) ContractEventName() string {
	return DvpPendingProofAddedToVaultEventName
}

// UnpackPendingProofAddedToVaultEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event PendingProofAddedToVault(uint256 indexed vaultId, uint256 indexed groupId, uint256 indexed targetReceiptId, (((uint256,uint256),(uint256[2],uint256[2]),(uint256,uint256)),uint256[],uint256,uint256[],uint256[],uint256[],uint256) pendingProof)
func (dvp *Dvp) UnpackPendingProofAddedToVaultEvent(log *types.Log) (*DvpPendingProofAddedToVault, error) {
	event := "PendingProofAddedToVault"
	if log.Topics[0] != dvp.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(DvpPendingProofAddedToVault)
	if len(log.Data) > 0 {
		if err := dvp.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range dvp.abi.Events[event].Inputs {
		if arg.Indexed {
			indexed = append(indexed, arg)
		}
	}
	if err := abi.ParseTopics(out, indexed, log.Topics[1:]); err != nil {
		return nil, err
	}
	out.Raw = log
	return out, nil
}

// DvpSettled represents a Settled event raised by the Dvp contract.
type DvpSettled struct {
	ReceiptId1 *big.Int
	ReceiptId2 *big.Int
	Raw        *types.Log // Blockchain specific contextual infos
}

const DvpSettledEventName = "Settled"

// ContractEventName returns the user-defined event name.
func (DvpSettled) ContractEventName() string {
	return DvpSettledEventName
}

// UnpackSettledEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event Settled(uint256 indexed receiptId1, uint256 indexed receiptId2)
func (dvp *Dvp) UnpackSettledEvent(log *types.Log) (*DvpSettled, error) {
	event := "Settled"
	if log.Topics[0] != dvp.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(DvpSettled)
	if len(log.Data) > 0 {
		if err := dvp.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range dvp.abi.Events[event].Inputs {
		if arg.Indexed {
			indexed = append(indexed, arg)
		}
	}
	if err := abi.ParseTopics(out, indexed, log.Topics[1:]); err != nil {
		return nil, err
	}
	out.Raw = log
	return out, nil
}

// DvpTokenAddedToGroup represents a TokenAddedToGroup event raised by the Dvp contract.
type DvpTokenAddedToGroup struct {
	VaultId       *big.Int
	TokenUniqueId *big.Int
	GroupId       *big.Int
	Raw           *types.Log // Blockchain specific contextual infos
}

const DvpTokenAddedToGroupEventName = "TokenAddedToGroup"

// ContractEventName returns the user-defined event name.
func (DvpTokenAddedToGroup) ContractEventName() string {
	return DvpTokenAddedToGroupEventName
}

// UnpackTokenAddedToGroupEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event TokenAddedToGroup(uint256 indexed vaultId, uint256 indexed tokenUniqueId, uint256 indexed groupId)
func (dvp *Dvp) UnpackTokenAddedToGroupEvent(log *types.Log) (*DvpTokenAddedToGroup, error) {
	event := "TokenAddedToGroup"
	if log.Topics[0] != dvp.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(DvpTokenAddedToGroup)
	if len(log.Data) > 0 {
		if err := dvp.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range dvp.abi.Events[event].Inputs {
		if arg.Indexed {
			indexed = append(indexed, arg)
		}
	}
	if err := abi.ParseTopics(out, indexed, log.Topics[1:]); err != nil {
		return nil, err
	}
	out.Raw = log
	return out, nil
}

// DvpVaultAddedToGroup represents a VaultAddedToGroup event raised by the Dvp contract.
type DvpVaultAddedToGroup struct {
	VaultId *big.Int
	GroupId *big.Int
	Raw     *types.Log // Blockchain specific contextual infos
}

const DvpVaultAddedToGroupEventName = "VaultAddedToGroup"

// ContractEventName returns the user-defined event name.
func (DvpVaultAddedToGroup) ContractEventName() string {
	return DvpVaultAddedToGroupEventName
}

// UnpackVaultAddedToGroupEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event VaultAddedToGroup(uint256 indexed vaultId, uint256 indexed groupId)
func (dvp *Dvp) UnpackVaultAddedToGroupEvent(log *types.Log) (*DvpVaultAddedToGroup, error) {
	event := "VaultAddedToGroup"
	if log.Topics[0] != dvp.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(DvpVaultAddedToGroup)
	if len(log.Data) > 0 {
		if err := dvp.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range dvp.abi.Events[event].Inputs {
		if arg.Indexed {
			indexed = append(indexed, arg)
		}
	}
	if err := abi.ParseTopics(out, indexed, log.Topics[1:]); err != nil {
		return nil, err
	}
	out.Raw = log
	return out, nil
}

// DvpVerifyOwnershipReceipt represents a VerifyOwnershipReceipt event raised by the Dvp contract.
type DvpVerifyOwnershipReceipt struct {
	Challenge     *big.Int
	AssetId       *big.Int
	TokenId       *big.Int
	AmountOrOne   *big.Int
	AssetAddress  common.Address
	SenderAddress common.Address
	Raw           *types.Log // Blockchain specific contextual infos
}

const DvpVerifyOwnershipReceiptEventName = "VerifyOwnershipReceipt"

// ContractEventName returns the user-defined event name.
func (DvpVerifyOwnershipReceipt) ContractEventName() string {
	return DvpVerifyOwnershipReceiptEventName
}

// UnpackVerifyOwnershipReceiptEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event VerifyOwnershipReceipt(uint256 indexed challenge, uint256 indexed assetId, uint256 indexed tokenId, uint256 amountOrOne, address assetAddress, address senderAddress)
func (dvp *Dvp) UnpackVerifyOwnershipReceiptEvent(log *types.Log) (*DvpVerifyOwnershipReceipt, error) {
	event := "VerifyOwnershipReceipt"
	if log.Topics[0] != dvp.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(DvpVerifyOwnershipReceipt)
	if len(log.Data) > 0 {
		if err := dvp.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range dvp.abi.Events[event].Inputs {
		if arg.Indexed {
			indexed = append(indexed, arg)
		}
	}
	if err := abi.ParseTopics(out, indexed, log.Topics[1:]); err != nil {
		return nil, err
	}
	out.Raw = log
	return out, nil
}

// UnpackError attempts to decode the provided error data using user-defined
// error definitions.
func (dvp *Dvp) UnpackError(raw []byte) (any, error) {
	if bytes.Equal(raw[:4], dvp.abi.Errors["AuctionAlreadyExists"].ID.Bytes()[:4]) {
		return dvp.UnpackAuctionAlreadyExistsError(raw[4:])
	}
	if bytes.Equal(raw[:4], dvp.abi.Errors["AuctionIdMismatch"].ID.Bytes()[:4]) {
		return dvp.UnpackAuctionIdMismatchError(raw[4:])
	}
	if bytes.Equal(raw[:4], dvp.abi.Errors["AuctionStateMismatch"].ID.Bytes()[:4]) {
		return dvp.UnpackAuctionStateMismatchError(raw[4:])
	}
	if bytes.Equal(raw[:4], dvp.abi.Errors["AuditorNotRegistered"].ID.Bytes()[:4]) {
		return dvp.UnpackAuditorNotRegisteredError(raw[4:])
	}
	if bytes.Equal(raw[:4], dvp.abi.Errors["BidStateMismatch"].ID.Bytes()[:4]) {
		return dvp.UnpackBidStateMismatchError(raw[4:])
	}
	if bytes.Equal(raw[:4], dvp.abi.Errors["BlindedBidMismatch"].ID.Bytes()[:4]) {
		return dvp.UnpackBlindedBidMismatchError(raw[4:])
	}
	if bytes.Equal(raw[:4], dvp.abi.Errors["DvpAuditorAlreadyRegistered"].ID.Bytes()[:4]) {
		return dvp.UnpackDvpAuditorAlreadyRegisteredError(raw[4:])
	}
	if bytes.Equal(raw[:4], dvp.abi.Errors["DvpBrokerAlreadyRegistered"].ID.Bytes()[:4]) {
		return dvp.UnpackDvpBrokerAlreadyRegisteredError(raw[4:])
	}
	if bytes.Equal(raw[:4], dvp.abi.Errors["DvpGroupFungibilityMismatch"].ID.Bytes()[:4]) {
		return dvp.UnpackDvpGroupFungibilityMismatchError(raw[4:])
	}
	if bytes.Equal(raw[:4], dvp.abi.Errors["DvpGroupIdOutOfRange"].ID.Bytes()[:4]) {
		return dvp.UnpackDvpGroupIdOutOfRangeError(raw[4:])
	}
	if bytes.Equal(raw[:4], dvp.abi.Errors["DvpGroupPairAlreadyRegistered"].ID.Bytes()[:4]) {
		return dvp.UnpackDvpGroupPairAlreadyRegisteredError(raw[4:])
	}
	if bytes.Equal(raw[:4], dvp.abi.Errors["DvpInvalidDeliveryMessage"].ID.Bytes()[:4]) {
		return dvp.UnpackDvpInvalidDeliveryMessageError(raw[4:])
	}
	if bytes.Equal(raw[:4], dvp.abi.Errors["DvpInvalidPassphrase"].ID.Bytes()[:4]) {
		return dvp.UnpackDvpInvalidPassphraseError(raw[4:])
	}
	if bytes.Equal(raw[:4], dvp.abi.Errors["DvpInvalidPaymentMessage"].ID.Bytes()[:4]) {
		return dvp.UnpackDvpInvalidPaymentMessageError(raw[4:])
	}
	if bytes.Equal(raw[:4], dvp.abi.Errors["DvpInvalidRevertCommitment"].ID.Bytes()[:4]) {
		return dvp.UnpackDvpInvalidRevertCommitmentError(raw[4:])
	}
	if bytes.Equal(raw[:4], dvp.abi.Errors["DvpInvalidStatementSize"].ID.Bytes()[:4]) {
		return dvp.UnpackDvpInvalidStatementSizeError(raw[4:])
	}
	if bytes.Equal(raw[:4], dvp.abi.Errors["DvpInvalidSwapGroupPair"].ID.Bytes()[:4]) {
		return dvp.UnpackDvpInvalidSwapGroupPairError(raw[4:])
	}
	if bytes.Equal(raw[:4], dvp.abi.Errors["DvpSwapAlreadyExists"].ID.Bytes()[:4]) {
		return dvp.UnpackDvpSwapAlreadyExistsError(raw[4:])
	}
	if bytes.Equal(raw[:4], dvp.abi.Errors["DvpSwapNotExpired"].ID.Bytes()[:4]) {
		return dvp.UnpackDvpSwapNotExpiredError(raw[4:])
	}
	if bytes.Equal(raw[:4], dvp.abi.Errors["DvpSwapNotFound"].ID.Bytes()[:4]) {
		return dvp.UnpackDvpSwapNotFoundError(raw[4:])
	}
	if bytes.Equal(raw[:4], dvp.abi.Errors["DvpSwapNotPending"].ID.Bytes()[:4]) {
		return dvp.UnpackDvpSwapNotPendingError(raw[4:])
	}
	if bytes.Equal(raw[:4], dvp.abi.Errors["FungibleDeliveryVault"].ID.Bytes()[:4]) {
		return dvp.UnpackFungibleDeliveryVaultError(raw[4:])
	}
	if bytes.Equal(raw[:4], dvp.abi.Errors["GroupMembershipMismatch"].ID.Bytes()[:4]) {
		return dvp.UnpackGroupMembershipMismatchError(raw[4:])
	}
	if bytes.Equal(raw[:4], dvp.abi.Errors["InvalidChallenge"].ID.Bytes()[:4]) {
		return dvp.UnpackInvalidChallengeError(raw[4:])
	}
	if bytes.Equal(raw[:4], dvp.abi.Errors["InvalidExchangeGroupPair"].ID.Bytes()[:4]) {
		return dvp.UnpackInvalidExchangeGroupPairError(raw[4:])
	}
	if bytes.Equal(raw[:4], dvp.abi.Errors["InvalidMerkleRoot"].ID.Bytes()[:4]) {
		return dvp.UnpackInvalidMerkleRootError(raw[4:])
	}
	if bytes.Equal(raw[:4], dvp.abi.Errors["InvalidNullifier"].ID.Bytes()[:4]) {
		return dvp.UnpackInvalidNullifierError(raw[4:])
	}
	if bytes.Equal(raw[:4], dvp.abi.Errors["InvalidNumberOfInputs"].ID.Bytes()[:4]) {
		return dvp.UnpackInvalidNumberOfInputsError(raw[4:])
	}
	if bytes.Equal(raw[:4], dvp.abi.Errors["InvalidNumberOfOutputs"].ID.Bytes()[:4]) {
		return dvp.UnpackInvalidNumberOfOutputsError(raw[4:])
	}
	if bytes.Equal(raw[:4], dvp.abi.Errors["InvalidOpening"].ID.Bytes()[:4]) {
		return dvp.UnpackInvalidOpeningError(raw[4:])
	}
	if bytes.Equal(raw[:4], dvp.abi.Errors["InvalidPartialProofReceipt"].ID.Bytes()[:4]) {
		return dvp.UnpackInvalidPartialProofReceiptError(raw[4:])
	}
	if bytes.Equal(raw[:4], dvp.abi.Errors["JoinSplitWithSameCommitments"].ID.Bytes()[:4]) {
		return dvp.UnpackJoinSplitWithSameCommitmentsError(raw[4:])
	}
	if bytes.Equal(raw[:4], dvp.abi.Errors["NonFungiblePaymentVault"].ID.Bytes()[:4]) {
		return dvp.UnpackNonFungiblePaymentVaultError(raw[4:])
	}
	if bytes.Equal(raw[:4], dvp.abi.Errors["NotImplemented"].ID.Bytes()[:4]) {
		return dvp.UnpackNotImplementedError(raw[4:])
	}
	if bytes.Equal(raw[:4], dvp.abi.Errors["NotWinningBidsCountMismatch"].ID.Bytes()[:4]) {
		return dvp.UnpackNotWinningBidsCountMismatchError(raw[4:])
	}
	if bytes.Equal(raw[:4], dvp.abi.Errors["RaylsAccessManagedContractPaused"].ID.Bytes()[:4]) {
		return dvp.UnpackRaylsAccessManagedContractPausedError(raw[4:])
	}
	if bytes.Equal(raw[:4], dvp.abi.Errors["RaylsAccessManagedInvalidAuthority"].ID.Bytes()[:4]) {
		return dvp.UnpackRaylsAccessManagedInvalidAuthorityError(raw[4:])
	}
	if bytes.Equal(raw[:4], dvp.abi.Errors["RaylsAccessManagedMustSchedule"].ID.Bytes()[:4]) {
		return dvp.UnpackRaylsAccessManagedMustScheduleError(raw[4:])
	}
	if bytes.Equal(raw[:4], dvp.abi.Errors["RaylsAccessManagedUnauthorized"].ID.Bytes()[:4]) {
		return dvp.UnpackRaylsAccessManagedUnauthorizedError(raw[4:])
	}
	if bytes.Equal(raw[:4], dvp.abi.Errors["ReentrancyGuardReentrantCall"].ID.Bytes()[:4]) {
		return dvp.UnpackReentrancyGuardReentrantCallError(raw[4:])
	}
	if bytes.Equal(raw[:4], dvp.abi.Errors["RottenChallenge"].ID.Bytes()[:4]) {
		return dvp.UnpackRottenChallengeError(raw[4:])
	}
	if bytes.Equal(raw[:4], dvp.abi.Errors["WinningBidOpeningMismatch"].ID.Bytes()[:4]) {
		return dvp.UnpackWinningBidOpeningMismatchError(raw[4:])
	}
	return nil, errors.New("Unknown error")
}

// DvpAuctionAlreadyExists represents a AuctionAlreadyExists error raised by the Dvp contract.
type DvpAuctionAlreadyExists struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error AuctionAlreadyExists()
func DvpAuctionAlreadyExistsErrorID() common.Hash {
	return common.HexToHash("0x04581cc83ec7c72cf1d5c596a9f6c3926f871b05e6a5259c8597e43a67357935")
}

// UnpackAuctionAlreadyExistsError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error AuctionAlreadyExists()
func (dvp *Dvp) UnpackAuctionAlreadyExistsError(raw []byte) (*DvpAuctionAlreadyExists, error) {
	out := new(DvpAuctionAlreadyExists)
	if err := dvp.abi.UnpackIntoInterface(out, "AuctionAlreadyExists", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// DvpAuctionIdMismatch represents a AuctionIdMismatch error raised by the Dvp contract.
type DvpAuctionIdMismatch struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error AuctionIdMismatch()
func DvpAuctionIdMismatchErrorID() common.Hash {
	return common.HexToHash("0x7264443b4a13b57ff1536458b77853110de47d3a2dafbec7c11ec292247280ee")
}

// UnpackAuctionIdMismatchError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error AuctionIdMismatch()
func (dvp *Dvp) UnpackAuctionIdMismatchError(raw []byte) (*DvpAuctionIdMismatch, error) {
	out := new(DvpAuctionIdMismatch)
	if err := dvp.abi.UnpackIntoInterface(out, "AuctionIdMismatch", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// DvpAuctionStateMismatch represents a AuctionStateMismatch error raised by the Dvp contract.
type DvpAuctionStateMismatch struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error AuctionStateMismatch()
func DvpAuctionStateMismatchErrorID() common.Hash {
	return common.HexToHash("0x3cc0342109499668631b1026af4cfec71300612a312929186232904abeff7df8")
}

// UnpackAuctionStateMismatchError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error AuctionStateMismatch()
func (dvp *Dvp) UnpackAuctionStateMismatchError(raw []byte) (*DvpAuctionStateMismatch, error) {
	out := new(DvpAuctionStateMismatch)
	if err := dvp.abi.UnpackIntoInterface(out, "AuctionStateMismatch", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// DvpAuditorNotRegistered represents a AuditorNotRegistered error raised by the Dvp contract.
type DvpAuditorNotRegistered struct {
	Arg0 *big.Int
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error AuditorNotRegistered(uint256 arg0)
func DvpAuditorNotRegisteredErrorID() common.Hash {
	return common.HexToHash("0x494760833fffd72e41b1a2c27fb37ce3aa452d773acc77a2dcb27f45055e27cb")
}

// UnpackAuditorNotRegisteredError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error AuditorNotRegistered(uint256 arg0)
func (dvp *Dvp) UnpackAuditorNotRegisteredError(raw []byte) (*DvpAuditorNotRegistered, error) {
	out := new(DvpAuditorNotRegistered)
	if err := dvp.abi.UnpackIntoInterface(out, "AuditorNotRegistered", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// DvpBidStateMismatch represents a BidStateMismatch error raised by the Dvp contract.
type DvpBidStateMismatch struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error BidStateMismatch()
func DvpBidStateMismatchErrorID() common.Hash {
	return common.HexToHash("0xffa541cbdd98412877647f85804b1330702bb0dde3f287cb555b0dfc22bc17f5")
}

// UnpackBidStateMismatchError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error BidStateMismatch()
func (dvp *Dvp) UnpackBidStateMismatchError(raw []byte) (*DvpBidStateMismatch, error) {
	out := new(DvpBidStateMismatch)
	if err := dvp.abi.UnpackIntoInterface(out, "BidStateMismatch", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// DvpBlindedBidMismatch represents a BlindedBidMismatch error raised by the Dvp contract.
type DvpBlindedBidMismatch struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error BlindedBidMismatch()
func DvpBlindedBidMismatchErrorID() common.Hash {
	return common.HexToHash("0x120ef2708e148eeffb08e72e5df99bc72a556d31cf6b6a99824d8dd47bbae439")
}

// UnpackBlindedBidMismatchError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error BlindedBidMismatch()
func (dvp *Dvp) UnpackBlindedBidMismatchError(raw []byte) (*DvpBlindedBidMismatch, error) {
	out := new(DvpBlindedBidMismatch)
	if err := dvp.abi.UnpackIntoInterface(out, "BlindedBidMismatch", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// DvpDvpAuditorAlreadyRegistered represents a Dvp__AuditorAlreadyRegistered error raised by the Dvp contract.
type DvpDvpAuditorAlreadyRegistered struct {
	Arg0 *big.Int
	Arg1 *big.Int
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error Dvp__AuditorAlreadyRegistered(uint256 arg0, uint256 arg1)
func DvpDvpAuditorAlreadyRegisteredErrorID() common.Hash {
	return common.HexToHash("0x7973b2249376125582bdd344d75d2524f1a99dc00c2cfd4e932a851b12782889")
}

// UnpackDvpAuditorAlreadyRegisteredError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error Dvp__AuditorAlreadyRegistered(uint256 arg0, uint256 arg1)
func (dvp *Dvp) UnpackDvpAuditorAlreadyRegisteredError(raw []byte) (*DvpDvpAuditorAlreadyRegistered, error) {
	out := new(DvpDvpAuditorAlreadyRegistered)
	if err := dvp.abi.UnpackIntoInterface(out, "DvpAuditorAlreadyRegistered", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// DvpDvpBrokerAlreadyRegistered represents a Dvp__BrokerAlreadyRegistered error raised by the Dvp contract.
type DvpDvpBrokerAlreadyRegistered struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error Dvp__BrokerAlreadyRegistered()
func DvpDvpBrokerAlreadyRegisteredErrorID() common.Hash {
	return common.HexToHash("0x9f8ef50d831bfcbffd27006c3da717578edaf3c8891533e1e433842ee9bb6d20")
}

// UnpackDvpBrokerAlreadyRegisteredError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error Dvp__BrokerAlreadyRegistered()
func (dvp *Dvp) UnpackDvpBrokerAlreadyRegisteredError(raw []byte) (*DvpDvpBrokerAlreadyRegistered, error) {
	out := new(DvpDvpBrokerAlreadyRegistered)
	if err := dvp.abi.UnpackIntoInterface(out, "DvpBrokerAlreadyRegistered", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// DvpDvpGroupFungibilityMismatch represents a Dvp__GroupFungibilityMismatch error raised by the Dvp contract.
type DvpDvpGroupFungibilityMismatch struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error Dvp__GroupFungibilityMismatch()
func DvpDvpGroupFungibilityMismatchErrorID() common.Hash {
	return common.HexToHash("0x662dd3233dad5c43f1871cae461a9f4650a9289c107e6a207b606ec0c5fbec9e")
}

// UnpackDvpGroupFungibilityMismatchError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error Dvp__GroupFungibilityMismatch()
func (dvp *Dvp) UnpackDvpGroupFungibilityMismatchError(raw []byte) (*DvpDvpGroupFungibilityMismatch, error) {
	out := new(DvpDvpGroupFungibilityMismatch)
	if err := dvp.abi.UnpackIntoInterface(out, "DvpGroupFungibilityMismatch", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// DvpDvpGroupIdOutOfRange represents a Dvp__GroupIdOutOfRange error raised by the Dvp contract.
type DvpDvpGroupIdOutOfRange struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error Dvp__GroupIdOutOfRange()
func DvpDvpGroupIdOutOfRangeErrorID() common.Hash {
	return common.HexToHash("0xaa7aacfdcbe96aac75feca9096ebd0190431203b5a096f3a93fb2b96476627ab")
}

// UnpackDvpGroupIdOutOfRangeError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error Dvp__GroupIdOutOfRange()
func (dvp *Dvp) UnpackDvpGroupIdOutOfRangeError(raw []byte) (*DvpDvpGroupIdOutOfRange, error) {
	out := new(DvpDvpGroupIdOutOfRange)
	if err := dvp.abi.UnpackIntoInterface(out, "DvpGroupIdOutOfRange", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// DvpDvpGroupPairAlreadyRegistered represents a Dvp__GroupPairAlreadyRegistered error raised by the Dvp contract.
type DvpDvpGroupPairAlreadyRegistered struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error Dvp__GroupPairAlreadyRegistered()
func DvpDvpGroupPairAlreadyRegisteredErrorID() common.Hash {
	return common.HexToHash("0x83895590ce62ad428548b6d993f915d9784733eedb679002b7243d03c3a7fbc5")
}

// UnpackDvpGroupPairAlreadyRegisteredError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error Dvp__GroupPairAlreadyRegistered()
func (dvp *Dvp) UnpackDvpGroupPairAlreadyRegisteredError(raw []byte) (*DvpDvpGroupPairAlreadyRegistered, error) {
	out := new(DvpDvpGroupPairAlreadyRegistered)
	if err := dvp.abi.UnpackIntoInterface(out, "DvpGroupPairAlreadyRegistered", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// DvpDvpInvalidDeliveryMessage represents a Dvp__InvalidDeliveryMessage error raised by the Dvp contract.
type DvpDvpInvalidDeliveryMessage struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error Dvp__InvalidDeliveryMessage()
func DvpDvpInvalidDeliveryMessageErrorID() common.Hash {
	return common.HexToHash("0x46064629ed0060b0c16339aee205e164e7b299f9d2b73ac1f4712f1a6758ec39")
}

// UnpackDvpInvalidDeliveryMessageError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error Dvp__InvalidDeliveryMessage()
func (dvp *Dvp) UnpackDvpInvalidDeliveryMessageError(raw []byte) (*DvpDvpInvalidDeliveryMessage, error) {
	out := new(DvpDvpInvalidDeliveryMessage)
	if err := dvp.abi.UnpackIntoInterface(out, "DvpInvalidDeliveryMessage", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// DvpDvpInvalidPassphrase represents a Dvp__InvalidPassphrase error raised by the Dvp contract.
type DvpDvpInvalidPassphrase struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error Dvp__InvalidPassphrase()
func DvpDvpInvalidPassphraseErrorID() common.Hash {
	return common.HexToHash("0xfabdd3b43d83367736f0b99fdc949827a5c61741c556fef579e80a303c2b1e77")
}

// UnpackDvpInvalidPassphraseError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error Dvp__InvalidPassphrase()
func (dvp *Dvp) UnpackDvpInvalidPassphraseError(raw []byte) (*DvpDvpInvalidPassphrase, error) {
	out := new(DvpDvpInvalidPassphrase)
	if err := dvp.abi.UnpackIntoInterface(out, "DvpInvalidPassphrase", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// DvpDvpInvalidPaymentMessage represents a Dvp__InvalidPaymentMessage error raised by the Dvp contract.
type DvpDvpInvalidPaymentMessage struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error Dvp__InvalidPaymentMessage()
func DvpDvpInvalidPaymentMessageErrorID() common.Hash {
	return common.HexToHash("0x1c4e374c59e40a10779571fc04075f25d370fa6de5cf9faa16ab9b1cb7ec19c5")
}

// UnpackDvpInvalidPaymentMessageError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error Dvp__InvalidPaymentMessage()
func (dvp *Dvp) UnpackDvpInvalidPaymentMessageError(raw []byte) (*DvpDvpInvalidPaymentMessage, error) {
	out := new(DvpDvpInvalidPaymentMessage)
	if err := dvp.abi.UnpackIntoInterface(out, "DvpInvalidPaymentMessage", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// DvpDvpInvalidRevertCommitment represents a Dvp__InvalidRevertCommitment error raised by the Dvp contract.
type DvpDvpInvalidRevertCommitment struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error Dvp__InvalidRevertCommitment()
func DvpDvpInvalidRevertCommitmentErrorID() common.Hash {
	return common.HexToHash("0x45aaba9046858075aea49a8b91bb31af13d7d2695571a26cba7a8a375bf13f77")
}

// UnpackDvpInvalidRevertCommitmentError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error Dvp__InvalidRevertCommitment()
func (dvp *Dvp) UnpackDvpInvalidRevertCommitmentError(raw []byte) (*DvpDvpInvalidRevertCommitment, error) {
	out := new(DvpDvpInvalidRevertCommitment)
	if err := dvp.abi.UnpackIntoInterface(out, "DvpInvalidRevertCommitment", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// DvpDvpInvalidStatementSize represents a Dvp__InvalidStatementSize error raised by the Dvp contract.
type DvpDvpInvalidStatementSize struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error Dvp__InvalidStatementSize()
func DvpDvpInvalidStatementSizeErrorID() common.Hash {
	return common.HexToHash("0x1d0cd2d0d196621c45428c29e25beb6299a0f1fbf925f9c895b348dd1ea46c34")
}

// UnpackDvpInvalidStatementSizeError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error Dvp__InvalidStatementSize()
func (dvp *Dvp) UnpackDvpInvalidStatementSizeError(raw []byte) (*DvpDvpInvalidStatementSize, error) {
	out := new(DvpDvpInvalidStatementSize)
	if err := dvp.abi.UnpackIntoInterface(out, "DvpInvalidStatementSize", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// DvpDvpInvalidSwapGroupPair represents a Dvp__InvalidSwapGroupPair error raised by the Dvp contract.
type DvpDvpInvalidSwapGroupPair struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error Dvp__InvalidSwapGroupPair()
func DvpDvpInvalidSwapGroupPairErrorID() common.Hash {
	return common.HexToHash("0x995807fc7a2631f023bb01e8792debc6406c0b443ddad11209a5c56d9f820f96")
}

// UnpackDvpInvalidSwapGroupPairError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error Dvp__InvalidSwapGroupPair()
func (dvp *Dvp) UnpackDvpInvalidSwapGroupPairError(raw []byte) (*DvpDvpInvalidSwapGroupPair, error) {
	out := new(DvpDvpInvalidSwapGroupPair)
	if err := dvp.abi.UnpackIntoInterface(out, "DvpInvalidSwapGroupPair", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// DvpDvpSwapAlreadyExists represents a Dvp__SwapAlreadyExists error raised by the Dvp contract.
type DvpDvpSwapAlreadyExists struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error Dvp__SwapAlreadyExists()
func DvpDvpSwapAlreadyExistsErrorID() common.Hash {
	return common.HexToHash("0x30c445f3fecb28984e889f6d825ba1940fbe5d069642372aafa2f85c4476f920")
}

// UnpackDvpSwapAlreadyExistsError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error Dvp__SwapAlreadyExists()
func (dvp *Dvp) UnpackDvpSwapAlreadyExistsError(raw []byte) (*DvpDvpSwapAlreadyExists, error) {
	out := new(DvpDvpSwapAlreadyExists)
	if err := dvp.abi.UnpackIntoInterface(out, "DvpSwapAlreadyExists", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// DvpDvpSwapNotExpired represents a Dvp__SwapNotExpired error raised by the Dvp contract.
type DvpDvpSwapNotExpired struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error Dvp__SwapNotExpired()
func DvpDvpSwapNotExpiredErrorID() common.Hash {
	return common.HexToHash("0x55803c37830b3d88fb092c753a09242b1b7e5d3a06c9b5ebbddeb4d2412ef845")
}

// UnpackDvpSwapNotExpiredError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error Dvp__SwapNotExpired()
func (dvp *Dvp) UnpackDvpSwapNotExpiredError(raw []byte) (*DvpDvpSwapNotExpired, error) {
	out := new(DvpDvpSwapNotExpired)
	if err := dvp.abi.UnpackIntoInterface(out, "DvpSwapNotExpired", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// DvpDvpSwapNotFound represents a Dvp__SwapNotFound error raised by the Dvp contract.
type DvpDvpSwapNotFound struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error Dvp__SwapNotFound()
func DvpDvpSwapNotFoundErrorID() common.Hash {
	return common.HexToHash("0xa008ac16307cf19cc440a07e5a65fad2bf5498235563b5e8737d444a0d619cd8")
}

// UnpackDvpSwapNotFoundError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error Dvp__SwapNotFound()
func (dvp *Dvp) UnpackDvpSwapNotFoundError(raw []byte) (*DvpDvpSwapNotFound, error) {
	out := new(DvpDvpSwapNotFound)
	if err := dvp.abi.UnpackIntoInterface(out, "DvpSwapNotFound", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// DvpDvpSwapNotPending represents a Dvp__SwapNotPending error raised by the Dvp contract.
type DvpDvpSwapNotPending struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error Dvp__SwapNotPending()
func DvpDvpSwapNotPendingErrorID() common.Hash {
	return common.HexToHash("0x3c7e4d2accb3cd547bec78b3d5148ea98ed2cf4df826d03ed56d650ee8f8adb1")
}

// UnpackDvpSwapNotPendingError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error Dvp__SwapNotPending()
func (dvp *Dvp) UnpackDvpSwapNotPendingError(raw []byte) (*DvpDvpSwapNotPending, error) {
	out := new(DvpDvpSwapNotPending)
	if err := dvp.abi.UnpackIntoInterface(out, "DvpSwapNotPending", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// DvpFungibleDeliveryVault represents a FungibleDeliveryVault error raised by the Dvp contract.
type DvpFungibleDeliveryVault struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error FungibleDeliveryVault()
func DvpFungibleDeliveryVaultErrorID() common.Hash {
	return common.HexToHash("0x670dd170e8c63c1a37ac6492a48a663c8be7e40756cad87030356b49d86d69e4")
}

// UnpackFungibleDeliveryVaultError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error FungibleDeliveryVault()
func (dvp *Dvp) UnpackFungibleDeliveryVaultError(raw []byte) (*DvpFungibleDeliveryVault, error) {
	out := new(DvpFungibleDeliveryVault)
	if err := dvp.abi.UnpackIntoInterface(out, "FungibleDeliveryVault", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// DvpGroupMembershipMismatch represents a GroupMembershipMismatch error raised by the Dvp contract.
type DvpGroupMembershipMismatch struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error GroupMembershipMismatch()
func DvpGroupMembershipMismatchErrorID() common.Hash {
	return common.HexToHash("0x058e47f20808110d7ca0669b8748b563d35894272ffc849a7d7ca3406ce3092f")
}

// UnpackGroupMembershipMismatchError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error GroupMembershipMismatch()
func (dvp *Dvp) UnpackGroupMembershipMismatchError(raw []byte) (*DvpGroupMembershipMismatch, error) {
	out := new(DvpGroupMembershipMismatch)
	if err := dvp.abi.UnpackIntoInterface(out, "GroupMembershipMismatch", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// DvpInvalidChallenge represents a InvalidChallenge error raised by the Dvp contract.
type DvpInvalidChallenge struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error InvalidChallenge()
func DvpInvalidChallengeErrorID() common.Hash {
	return common.HexToHash("0x92fa6fd36f01a57f4a4ab338b70872272b2991f95c9b087ddc91ca1ad302d4cf")
}

// UnpackInvalidChallengeError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error InvalidChallenge()
func (dvp *Dvp) UnpackInvalidChallengeError(raw []byte) (*DvpInvalidChallenge, error) {
	out := new(DvpInvalidChallenge)
	if err := dvp.abi.UnpackIntoInterface(out, "InvalidChallenge", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// DvpInvalidExchangeGroupPair represents a InvalidExchangeGroupPair error raised by the Dvp contract.
type DvpInvalidExchangeGroupPair struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error InvalidExchangeGroupPair()
func DvpInvalidExchangeGroupPairErrorID() common.Hash {
	return common.HexToHash("0xc24f9db6c39f44fa513ae2063bf654968cbbf9d50bcdac406caef2331c649d35")
}

// UnpackInvalidExchangeGroupPairError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error InvalidExchangeGroupPair()
func (dvp *Dvp) UnpackInvalidExchangeGroupPairError(raw []byte) (*DvpInvalidExchangeGroupPair, error) {
	out := new(DvpInvalidExchangeGroupPair)
	if err := dvp.abi.UnpackIntoInterface(out, "InvalidExchangeGroupPair", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// DvpInvalidMerkleRoot represents a InvalidMerkleRoot error raised by the Dvp contract.
type DvpInvalidMerkleRoot struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error InvalidMerkleRoot()
func DvpInvalidMerkleRootErrorID() common.Hash {
	return common.HexToHash("0x9dd854d332457a6415758d7c8a6d51a1b77786d562179d17107542ac2036054c")
}

// UnpackInvalidMerkleRootError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error InvalidMerkleRoot()
func (dvp *Dvp) UnpackInvalidMerkleRootError(raw []byte) (*DvpInvalidMerkleRoot, error) {
	out := new(DvpInvalidMerkleRoot)
	if err := dvp.abi.UnpackIntoInterface(out, "InvalidMerkleRoot", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// DvpInvalidNullifier represents a InvalidNullifier error raised by the Dvp contract.
type DvpInvalidNullifier struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error InvalidNullifier()
func DvpInvalidNullifierErrorID() common.Hash {
	return common.HexToHash("0x5d904cb2bdfbf6fcbc7ff7dc828b2fbb9b0811066733b81d4c5b5dbbee4d9c02")
}

// UnpackInvalidNullifierError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error InvalidNullifier()
func (dvp *Dvp) UnpackInvalidNullifierError(raw []byte) (*DvpInvalidNullifier, error) {
	out := new(DvpInvalidNullifier)
	if err := dvp.abi.UnpackIntoInterface(out, "InvalidNullifier", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// DvpInvalidNumberOfInputs represents a InvalidNumberOfInputs error raised by the Dvp contract.
type DvpInvalidNumberOfInputs struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error InvalidNumberOfInputs()
func DvpInvalidNumberOfInputsErrorID() common.Hash {
	return common.HexToHash("0x647ce72c2dbe9461ad375c21c29142e71c880eef510ef7bc14b76c980b7cab2d")
}

// UnpackInvalidNumberOfInputsError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error InvalidNumberOfInputs()
func (dvp *Dvp) UnpackInvalidNumberOfInputsError(raw []byte) (*DvpInvalidNumberOfInputs, error) {
	out := new(DvpInvalidNumberOfInputs)
	if err := dvp.abi.UnpackIntoInterface(out, "InvalidNumberOfInputs", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// DvpInvalidNumberOfOutputs represents a InvalidNumberOfOutputs error raised by the Dvp contract.
type DvpInvalidNumberOfOutputs struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error InvalidNumberOfOutputs()
func DvpInvalidNumberOfOutputsErrorID() common.Hash {
	return common.HexToHash("0x090d645497072b613292820cba5ba16a96f86aed66cf18e31a3e8f9611450cf1")
}

// UnpackInvalidNumberOfOutputsError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error InvalidNumberOfOutputs()
func (dvp *Dvp) UnpackInvalidNumberOfOutputsError(raw []byte) (*DvpInvalidNumberOfOutputs, error) {
	out := new(DvpInvalidNumberOfOutputs)
	if err := dvp.abi.UnpackIntoInterface(out, "InvalidNumberOfOutputs", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// DvpInvalidOpening represents a InvalidOpening error raised by the Dvp contract.
type DvpInvalidOpening struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error InvalidOpening()
func DvpInvalidOpeningErrorID() common.Hash {
	return common.HexToHash("0x399d2d9b53a26dc663ba55834acf43b5446f5f8b6dae32f76627a2cd208dab66")
}

// UnpackInvalidOpeningError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error InvalidOpening()
func (dvp *Dvp) UnpackInvalidOpeningError(raw []byte) (*DvpInvalidOpening, error) {
	out := new(DvpInvalidOpening)
	if err := dvp.abi.UnpackIntoInterface(out, "InvalidOpening", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// DvpInvalidPartialProofReceipt represents a InvalidPartialProofReceipt error raised by the Dvp contract.
type DvpInvalidPartialProofReceipt struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error InvalidPartialProofReceipt()
func DvpInvalidPartialProofReceiptErrorID() common.Hash {
	return common.HexToHash("0xc32e0148d6c0187addcd3c05cbf2261c2b6b9ee6ad91620b3e38860d584dd5ed")
}

// UnpackInvalidPartialProofReceiptError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error InvalidPartialProofReceipt()
func (dvp *Dvp) UnpackInvalidPartialProofReceiptError(raw []byte) (*DvpInvalidPartialProofReceipt, error) {
	out := new(DvpInvalidPartialProofReceipt)
	if err := dvp.abi.UnpackIntoInterface(out, "InvalidPartialProofReceipt", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// DvpJoinSplitWithSameCommitments represents a JoinSplitWithSameCommitments error raised by the Dvp contract.
type DvpJoinSplitWithSameCommitments struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error JoinSplitWithSameCommitments()
func DvpJoinSplitWithSameCommitmentsErrorID() common.Hash {
	return common.HexToHash("0xfaaa05e9fe3dff29ac37289237f6421de615c729d53e369f47f92874974df800")
}

// UnpackJoinSplitWithSameCommitmentsError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error JoinSplitWithSameCommitments()
func (dvp *Dvp) UnpackJoinSplitWithSameCommitmentsError(raw []byte) (*DvpJoinSplitWithSameCommitments, error) {
	out := new(DvpJoinSplitWithSameCommitments)
	if err := dvp.abi.UnpackIntoInterface(out, "JoinSplitWithSameCommitments", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// DvpNonFungiblePaymentVault represents a NonFungiblePaymentVault error raised by the Dvp contract.
type DvpNonFungiblePaymentVault struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error NonFungiblePaymentVault()
func DvpNonFungiblePaymentVaultErrorID() common.Hash {
	return common.HexToHash("0x1b778a5412d5f73980fe8cbf5a205f53d46e423f8cde365095cb9cdf6881ac01")
}

// UnpackNonFungiblePaymentVaultError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error NonFungiblePaymentVault()
func (dvp *Dvp) UnpackNonFungiblePaymentVaultError(raw []byte) (*DvpNonFungiblePaymentVault, error) {
	out := new(DvpNonFungiblePaymentVault)
	if err := dvp.abi.UnpackIntoInterface(out, "NonFungiblePaymentVault", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// DvpNotImplemented represents a NotImplemented error raised by the Dvp contract.
type DvpNotImplemented struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error NotImplemented()
func DvpNotImplementedErrorID() common.Hash {
	return common.HexToHash("0xd6234725c2592490a5b2926ed2315070d2f568d079cb53600cc4c507f13f8289")
}

// UnpackNotImplementedError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error NotImplemented()
func (dvp *Dvp) UnpackNotImplementedError(raw []byte) (*DvpNotImplemented, error) {
	out := new(DvpNotImplemented)
	if err := dvp.abi.UnpackIntoInterface(out, "NotImplemented", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// DvpNotWinningBidsCountMismatch represents a NotWinningBidsCountMismatch error raised by the Dvp contract.
type DvpNotWinningBidsCountMismatch struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error NotWinningBidsCountMismatch()
func DvpNotWinningBidsCountMismatchErrorID() common.Hash {
	return common.HexToHash("0x6bf5c54afbe241bd255f0021e6833d48a7446a270a574d1c53980d4b862f66c4")
}

// UnpackNotWinningBidsCountMismatchError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error NotWinningBidsCountMismatch()
func (dvp *Dvp) UnpackNotWinningBidsCountMismatchError(raw []byte) (*DvpNotWinningBidsCountMismatch, error) {
	out := new(DvpNotWinningBidsCountMismatch)
	if err := dvp.abi.UnpackIntoInterface(out, "NotWinningBidsCountMismatch", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// DvpRaylsAccessManagedContractPaused represents a RaylsAccessManaged__ContractPaused error raised by the Dvp contract.
type DvpRaylsAccessManagedContractPaused struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error RaylsAccessManaged__ContractPaused()
func DvpRaylsAccessManagedContractPausedErrorID() common.Hash {
	return common.HexToHash("0xcc9855adb54c3eae80b5cb29441d3eb81c3fcdc928abeaba0ac7979b3072c3ab")
}

// UnpackRaylsAccessManagedContractPausedError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error RaylsAccessManaged__ContractPaused()
func (dvp *Dvp) UnpackRaylsAccessManagedContractPausedError(raw []byte) (*DvpRaylsAccessManagedContractPaused, error) {
	out := new(DvpRaylsAccessManagedContractPaused)
	if err := dvp.abi.UnpackIntoInterface(out, "RaylsAccessManagedContractPaused", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// DvpRaylsAccessManagedInvalidAuthority represents a RaylsAccessManaged__InvalidAuthority error raised by the Dvp contract.
type DvpRaylsAccessManagedInvalidAuthority struct {
	Authority common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error RaylsAccessManaged__InvalidAuthority(address authority)
func DvpRaylsAccessManagedInvalidAuthorityErrorID() common.Hash {
	return common.HexToHash("0x894403477abbba3df1c9b757907cda2d4b4d9ba52242d5e6bdae2dc6a343662b")
}

// UnpackRaylsAccessManagedInvalidAuthorityError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error RaylsAccessManaged__InvalidAuthority(address authority)
func (dvp *Dvp) UnpackRaylsAccessManagedInvalidAuthorityError(raw []byte) (*DvpRaylsAccessManagedInvalidAuthority, error) {
	out := new(DvpRaylsAccessManagedInvalidAuthority)
	if err := dvp.abi.UnpackIntoInterface(out, "RaylsAccessManagedInvalidAuthority", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// DvpRaylsAccessManagedMustSchedule represents a RaylsAccessManaged__MustSchedule error raised by the Dvp contract.
type DvpRaylsAccessManagedMustSchedule struct {
	Caller common.Address
	Delay  uint32
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error RaylsAccessManaged__MustSchedule(address caller, uint32 delay)
func DvpRaylsAccessManagedMustScheduleErrorID() common.Hash {
	return common.HexToHash("0xa42687899abecbc3886d9231154d03927d3e61df458dda1e50129311f272684f")
}

// UnpackRaylsAccessManagedMustScheduleError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error RaylsAccessManaged__MustSchedule(address caller, uint32 delay)
func (dvp *Dvp) UnpackRaylsAccessManagedMustScheduleError(raw []byte) (*DvpRaylsAccessManagedMustSchedule, error) {
	out := new(DvpRaylsAccessManagedMustSchedule)
	if err := dvp.abi.UnpackIntoInterface(out, "RaylsAccessManagedMustSchedule", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// DvpRaylsAccessManagedUnauthorized represents a RaylsAccessManaged__Unauthorized error raised by the Dvp contract.
type DvpRaylsAccessManagedUnauthorized struct {
	Caller common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error RaylsAccessManaged__Unauthorized(address caller)
func DvpRaylsAccessManagedUnauthorizedErrorID() common.Hash {
	return common.HexToHash("0xbb34f40c15462e9dd86b312158515c7ee80e2406cfa90a92a6fb6fad1bc48a48")
}

// UnpackRaylsAccessManagedUnauthorizedError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error RaylsAccessManaged__Unauthorized(address caller)
func (dvp *Dvp) UnpackRaylsAccessManagedUnauthorizedError(raw []byte) (*DvpRaylsAccessManagedUnauthorized, error) {
	out := new(DvpRaylsAccessManagedUnauthorized)
	if err := dvp.abi.UnpackIntoInterface(out, "RaylsAccessManagedUnauthorized", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// DvpReentrancyGuardReentrantCall represents a ReentrancyGuardReentrantCall error raised by the Dvp contract.
type DvpReentrancyGuardReentrantCall struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error ReentrancyGuardReentrantCall()
func DvpReentrancyGuardReentrantCallErrorID() common.Hash {
	return common.HexToHash("0x3ee5aeb571de7fc460830b4d0017439a1ca56fb0bc39062227ade4fe4a24c1ca")
}

// UnpackReentrancyGuardReentrantCallError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error ReentrancyGuardReentrantCall()
func (dvp *Dvp) UnpackReentrancyGuardReentrantCallError(raw []byte) (*DvpReentrancyGuardReentrantCall, error) {
	out := new(DvpReentrancyGuardReentrantCall)
	if err := dvp.abi.UnpackIntoInterface(out, "ReentrancyGuardReentrantCall", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// DvpRottenChallenge represents a RottenChallenge error raised by the Dvp contract.
type DvpRottenChallenge struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error RottenChallenge()
func DvpRottenChallengeErrorID() common.Hash {
	return common.HexToHash("0x488197f25c1730adcb3c8fbef3e84006f8ac7d3139a2cef9d7168de60335b45a")
}

// UnpackRottenChallengeError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error RottenChallenge()
func (dvp *Dvp) UnpackRottenChallengeError(raw []byte) (*DvpRottenChallenge, error) {
	out := new(DvpRottenChallenge)
	if err := dvp.abi.UnpackIntoInterface(out, "RottenChallenge", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// DvpWinningBidOpeningMismatch represents a WinningBidOpeningMismatch error raised by the Dvp contract.
type DvpWinningBidOpeningMismatch struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error WinningBidOpeningMismatch()
func DvpWinningBidOpeningMismatchErrorID() common.Hash {
	return common.HexToHash("0xa828f42cf0867032dfeefe98dd83e2fb43071f40635efbaa478928e07b20bee5")
}

// UnpackWinningBidOpeningMismatchError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error WinningBidOpeningMismatch()
func (dvp *Dvp) UnpackWinningBidOpeningMismatchError(raw []byte) (*DvpWinningBidOpeningMismatch, error) {
	out := new(DvpWinningBidOpeningMismatch)
	if err := dvp.abi.UnpackIntoInterface(out, "WinningBidOpeningMismatch", raw); err != nil {
		return nil, err
	}
	return out, nil
}
