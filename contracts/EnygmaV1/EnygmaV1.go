// Code generated via abigen V2 - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package EnygmaV1

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

// EnygmaCreationParams is an auto generated low-level Go binding around an user-defined struct.
type EnygmaCreationParams struct {
	Name                       string
	Symbol                     string
	Decimals                   uint8
	ResourceId                 [32]byte
	Owner                      common.Address
	OwnerChainId               *big.Int
	ParticipantStorageContract common.Address
	Endpoint                   common.Address
	TokenRegistryContract      common.Address
	EnygmaTeleport             common.Address
	Factory                    common.Address
}

// IEnygmaV1EnygmaPointWithChainId is an auto generated low-level Go binding around an user-defined struct.
type IEnygmaV1EnygmaPointWithChainId struct {
	C1      *big.Int
	C2      *big.Int
	ChainId *big.Int
}

// IEnygmaV1EnygmaPublicKeyWithChainId is an auto generated low-level Go binding around an user-defined struct.
type IEnygmaV1EnygmaPublicKeyWithChainId struct {
	PublicKey *big.Int
	ChainId   *big.Int
}

// IEnygmaV1ExtractedProofData is an auto generated low-level Go binding around an user-defined struct.
type IEnygmaV1ExtractedProofData struct {
	K                uint8
	ProofLength      *big.Int
	ArrayHashSecrets []*big.Int
	PublicKeys       []*big.Int
	Balances         []IEnygmaV1Point
	Commitments      []IEnygmaV1Point
	Nullifier        *big.Int
	BlockNumber      *big.Int
	ChainIds         []*big.Int
	MessageTags      []*big.Int
}

// IEnygmaV1PendingMintOrBurn is an auto generated low-level Go binding around an user-defined struct.
type IEnygmaV1PendingMintOrBurn struct {
	PointToAddToBalance IEnygmaV1EnygmaPointWithChainId
	Amount              *big.Int
	BlockNumber         *big.Int
	TransactionType     uint8
}

// IEnygmaV1PendingTransaction is an auto generated low-level Go binding around an user-defined struct.
type IEnygmaV1PendingTransaction struct {
	PointsToAddToBalance []IEnygmaV1EnygmaPointWithChainId
	Nullifier            *big.Int
	TransactionType      uint8
}

// IEnygmaV1Point is an auto generated low-level Go binding around an user-defined struct.
type IEnygmaV1Point struct {
	C1 *big.Int
	C2 *big.Int
}

// IEnygmaV1SupplyUpdateTx is an auto generated low-level Go binding around an user-defined struct.
type IEnygmaV1SupplyUpdateTx struct {
	Amount *big.Int
	TxType uint8
}

// IEnygmaV1TransferProof is an auto generated low-level Go binding around an user-defined struct.
type IEnygmaV1TransferProof struct {
	PiA          [2]*big.Int
	PiB          [2][2]*big.Int
	PiC          [2]*big.Int
	PublicSignal []*big.Int
}

// EnygmaV1MetaData contains all meta data concerning the EnygmaV1 contract.
var EnygmaV1MetaData = bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"params\",\"type\":\"tuple\",\"internalType\":\"structEnygmaCreationParams\",\"components\":[{\"name\":\"name\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"symbol\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"decimals\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"resourceId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"ownerChainId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"participantStorageContract\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"endpoint\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"tokenRegistryContract\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"enygmaTeleport\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"factory\",\"type\":\"address\",\"internalType\":\"address\"}]}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"Name\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"Symbol\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"addTransferVerifier\",\"inputs\":[{\"name\":\"verifier\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"k\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"checkTotalSumOfBalances\",\"inputs\":[{\"name\":\"blockNumber\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"consumedNullifiers\",\"inputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"derivePk\",\"inputs\":[{\"name\":\"v\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"x2\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"y2\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"derivePkH\",\"inputs\":[{\"name\":\"r\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"x2\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"y2\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"dvpAddPendingTransaction\",\"inputs\":[{\"name\":\"proof\",\"type\":\"tuple\",\"internalType\":\"structIEnygmaV1.TransferProof\",\"components\":[{\"name\":\"pi_a\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"},{\"name\":\"pi_b\",\"type\":\"uint256[2][2]\",\"internalType\":\"uint256[2][2]\"},{\"name\":\"pi_c\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"},{\"name\":\"public_signal\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"}]},{\"name\":\"transactionType\",\"type\":\"uint8\",\"internalType\":\"enumIEnygmaV1.TxType\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"dvpChainId\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"dvpFinalisePendingTransactions\",\"inputs\":[{\"name\":\"currentBlockNumber\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"dvpIntegrationContractAddress\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"dvpSendEvents\",\"inputs\":[{\"name\":\"proofData\",\"type\":\"tuple\",\"internalType\":\"structIEnygmaV1.ExtractedProofData\",\"components\":[{\"name\":\"k\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"proofLength\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"arrayHashSecrets\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"},{\"name\":\"publicKeys\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"},{\"name\":\"balances\",\"type\":\"tuple[]\",\"internalType\":\"structIEnygmaV1.Point[]\",\"components\":[{\"name\":\"c1\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"c2\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"commitments\",\"type\":\"tuple[]\",\"internalType\":\"structIEnygmaV1.Point[]\",\"components\":[{\"name\":\"c1\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"c2\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"nullifier\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"blockNumber\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"chainIds\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"},{\"name\":\"messageTags\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"}]},{\"name\":\"encryptedMessages\",\"type\":\"bytes[]\",\"internalType\":\"bytes[]\"},{\"name\":\"encryptedUpdate\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"dvpSetLastblockNumPending\",\"inputs\":[{\"name\":\"newValue\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"dvpValidateTransferInputs\",\"inputs\":[{\"name\":\"proof\",\"type\":\"tuple\",\"internalType\":\"structIEnygmaV1.TransferProof\",\"components\":[{\"name\":\"pi_a\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"},{\"name\":\"pi_b\",\"type\":\"uint256[2][2]\",\"internalType\":\"uint256[2][2]\"},{\"name\":\"pi_c\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"},{\"name\":\"public_signal\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"}]}],\"outputs\":[],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"enygmaTeleport\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractEnygmaTeleport\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"factory\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getAddressByResourceId\",\"inputs\":[{\"name\":\"_resourceId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getBalanceByBlockNumber\",\"inputs\":[{\"name\":\"chainId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"blockNumber\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"x\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getBalanceFinalised\",\"inputs\":[{\"name\":\"chainId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"x\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getBalancePending\",\"inputs\":[{\"name\":\"chainId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"x\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getDvpIntegrationContractAddress\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getEndpointAddress\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getLastblockNumAtCurrentBlockNumber\",\"inputs\":[{\"name\":\"currentBlockNumber\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getNextBlockNumberToFinaliseAfter\",\"inputs\":[{\"name\":\"blockNumber\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getPendingMintsAndBurns\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"tuple[]\",\"internalType\":\"structIEnygmaV1.PendingMintOrBurn[]\",\"components\":[{\"name\":\"pointToAddToBalance\",\"type\":\"tuple\",\"internalType\":\"structIEnygmaV1.EnygmaPointWithChainId\",\"components\":[{\"name\":\"c1\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"c2\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"chainId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"blockNumber\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"transactionType\",\"type\":\"uint8\",\"internalType\":\"enumIEnygmaV1.TxType\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getPendingTransactions\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"tuple[]\",\"internalType\":\"structIEnygmaV1.PendingTransaction[]\",\"components\":[{\"name\":\"pointsToAddToBalance\",\"type\":\"tuple[]\",\"internalType\":\"structIEnygmaV1.EnygmaPointWithChainId[]\",\"components\":[{\"name\":\"c1\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"c2\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"chainId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"nullifier\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"transactionType\",\"type\":\"uint8\",\"internalType\":\"enumIEnygmaV1.TxType\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getPublicValuesByBlockNumber\",\"inputs\":[{\"name\":\"blockNumber\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple[]\",\"internalType\":\"structIEnygmaV1.EnygmaPointWithChainId[]\",\"components\":[{\"name\":\"c1\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"c2\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"chainId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"\",\"type\":\"tuple[]\",\"internalType\":\"structIEnygmaV1.EnygmaPublicKeyWithChainId[]\",\"components\":[{\"name\":\"publicKey\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"chainId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getPublicValuesFinalised\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"tuple[]\",\"internalType\":\"structIEnygmaV1.EnygmaPointWithChainId[]\",\"components\":[{\"name\":\"c1\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"c2\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"chainId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"\",\"type\":\"tuple[]\",\"internalType\":\"structIEnygmaV1.EnygmaPublicKeyWithChainId[]\",\"components\":[{\"name\":\"publicKey\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"chainId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getPublicValuesPending\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"tuple[]\",\"internalType\":\"structIEnygmaV1.EnygmaPointWithChainId[]\",\"components\":[{\"name\":\"c1\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"c2\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"chainId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"\",\"type\":\"tuple[]\",\"internalType\":\"structIEnygmaV1.EnygmaPublicKeyWithChainId[]\",\"components\":[{\"name\":\"publicKey\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"chainId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getTotalRegisteredBanks\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getTotalSupply\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getTransferVerifierAddress\",\"inputs\":[{\"name\":\"k\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"isNullifierUnspent\",\"inputs\":[{\"name\":\"nullifier\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"lastblockNum\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"lastblockNumAtCurrentBlockNumber\",\"inputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"lastblockNumPending\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"negateOnCurve\",\"inputs\":[{\"name\":\"x\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"nextBlockNumber\",\"inputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"ownerChainId\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"participantStorageContract\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"pedCom\",\"inputs\":[{\"name\":\"v\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"r\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"pendingBalancesTallied\",\"inputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"pendingMintsAndBurns\",\"inputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"pointToAddToBalance\",\"type\":\"tuple\",\"internalType\":\"structIEnygmaV1.EnygmaPointWithChainId\",\"components\":[{\"name\":\"c1\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"c2\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"chainId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"blockNumber\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"transactionType\",\"type\":\"uint8\",\"internalType\":\"enumIEnygmaV1.TxType\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"pendingTransactions\",\"inputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"nullifier\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"transactionType\",\"type\":\"uint8\",\"internalType\":\"enumIEnygmaV1.TxType\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"raylsNodeUserGovernance\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIUserGovernance\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"referenceBalance\",\"inputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"c1\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"c2\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"chainId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"resourceId\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"setDvpIntegrationContract\",\"inputs\":[{\"name\":\"_dvpIntegrationContractAddress\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setResourceId\",\"inputs\":[{\"name\":\"_resourceId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"tokenRegistryContract\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"totalSupply\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"totalSupplyX\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"totalSupplyY\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"transferBatch\",\"inputs\":[{\"name\":\"proof\",\"type\":\"tuple\",\"internalType\":\"structIEnygmaV1.TransferProof\",\"components\":[{\"name\":\"pi_a\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"},{\"name\":\"pi_b\",\"type\":\"uint256[2][2]\",\"internalType\":\"uint256[2][2]\"},{\"name\":\"pi_c\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"},{\"name\":\"public_signal\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"}]},{\"name\":\"encryptedMessages\",\"type\":\"bytes[]\",\"internalType\":\"bytes[]\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"transferVerifiers\",\"inputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"updateSupply\",\"inputs\":[{\"name\":\"_chainId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"_blockNumber\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"_update\",\"type\":\"tuple\",\"internalType\":\"structIEnygmaV1.SupplyUpdateTx\",\"components\":[{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"txType\",\"type\":\"uint8\",\"internalType\":\"enumIEnygmaV1.TxType\"}]}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"event\",\"name\":\"BalancesFinalised\",\"inputs\":[{\"name\":\"blockNumber\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"BurnSuccessful\",\"inputs\":[{\"name\":\"chainId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"burnValue\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"NullifierConsumed\",\"inputs\":[{\"name\":\"resourceId\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"nullifier\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"blockNumber\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"txType\",\"type\":\"uint8\",\"indexed\":false,\"internalType\":\"enumIEnygmaV1.TxType\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"SupplyMinted\",\"inputs\":[{\"name\":\"lastblockNum\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"toChainId\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"TokenRegistrationSubmitted\",\"inputs\":[{\"name\":\"tokenAddress\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"TransactionSuccessful\",\"inputs\":[{\"name\":\"senderAddress\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"VerifierRegistered\",\"inputs\":[{\"name\":\"verifierAddress\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"k\",\"type\":\"uint8\",\"indexed\":false,\"internalType\":\"uint8\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"EnygmaV1__OnlyDvpIntegrationAllowed\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"EnygmaV1__OnlyFactoryAllowed\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"RaylsApp__HubNotActive\",\"inputs\":[{\"name\":\"tokenAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"privacyNodeStatus\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"hubStatus\",\"type\":\"uint8\",\"internalType\":\"uint8\"}]},{\"type\":\"error\",\"name\":\"RaylsApp__PrivacyNodeFrozen\",\"inputs\":[{\"name\":\"tokenAddress\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"RaylsApp__PrivacyNodeNotActive\",\"inputs\":[{\"name\":\"tokenAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"privacyNodeStatus\",\"type\":\"uint8\",\"internalType\":\"uint8\"}]},{\"type\":\"error\",\"name\":\"RaylsApp__PublicChainNotActive\",\"inputs\":[{\"name\":\"tokenAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"privacyNodeStatus\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"publicChainStatus\",\"type\":\"uint8\",\"internalType\":\"uint8\"}]},{\"type\":\"error\",\"name\":\"RaylsApp__ResourceNotApproved\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"RaylsApp__TokenRegistryNotConfigured\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"RaylsApp__UnauthorizedTokenRegistry\",\"inputs\":[{\"name\":\"caller\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"tokenRegistry\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"RaylsApp__UserNotRegistered\",\"inputs\":[{\"name\":\"caller\",\"type\":\"address\",\"internalType\":\"address\"}]}]",
	ID:  "EnygmaV1",
	Bin: "0x60a06040526000600755600060085560006010553480156200002057600080fd5b50604051620063be380380620063be83398101604081905262000043916200044e565b60e0810151600080546001600160a01b0319166001600160a01b038316178155508151600491506200007690826200061e565b5060208101516005906200008b90826200061e565b506040818101516006805460ff191660ff90921691909117905543600c819055600d8190556000818152601860205291822055608080830151600a80546001600160a01b03199081166001600160a01b0393841617909155600984905560c0850151600e805483169184169190911790556007939093556001600855610100840151600f80548516918316919091179055606084015160035560a0840151600b5561012084015160118054909416908216179092556101408301519091169052620001556200015c565b50620007b0565b600e5460408051632d33587560e01b815290516001600160a01b03909216916000918391632d3358759160048082019286929091908290030181865afa158015620001ab573d6000803e3d6000fd5b505050506040513d6000823e601f3d908101601f19168201604052620001d59190810190620006ea565b9050620001ea6010546200023860201b60201c565b60005b8151811015620002335760008282815181106200020e576200020e6200079a565b6020026020010151905062000229816200023860201b60201c565b50600101620001ed565b505050565b600c546000908152601260209081526040808320848452909152902054158015620002805750600c546000908152601260209081526040808320848452909152902060010154155b15620002a957600c54600090815260126020908152604080832084845290915290206001908101555b600d546000908152601260209081526040808320848452909152902054158015620002f15750600d546000908152601260209081526040808320848452909152902060010154155b156200031a57600d54600090815260126020908152604080832084845290915290206001908101555b50565b634e487b7160e01b600052604160045260246000fd5b60405161016081016001600160401b03811182821017156200035957620003596200031d565b60405290565b604051601f8201601f191681016001600160401b03811182821017156200038a576200038a6200031d565b604052919050565b600082601f830112620003a457600080fd5b81516001600160401b03811115620003c057620003c06200031d565b6020620003d6601f8301601f191682016200035f565b8281528582848701011115620003eb57600080fd5b60005b838110156200040b578581018301518282018401528201620003ee565b506000928101909101919091529392505050565b805160ff811681146200043157600080fd5b919050565b80516001600160a01b03811681146200043157600080fd5b6000602082840312156200046157600080fd5b81516001600160401b03808211156200047957600080fd5b9083019061016082860312156200048f57600080fd5b6200049962000333565b825182811115620004a957600080fd5b620004b78782860162000392565b825250602083015182811115620004cd57600080fd5b620004db8782860162000392565b602083015250620004ef604084016200041f565b6040820152606083015160608201526200050c6080840162000436565b608082015260a083015160a08201526200052960c0840162000436565b60c08201526200053c60e0840162000436565b60e082015261010091506200055382840162000436565b8282015261012091506200056982840162000436565b8282015261014091506200057f82840162000436565b91810191909152949350505050565b600181811c90821680620005a357607f821691505b602082108103620005c457634e487b7160e01b600052602260045260246000fd5b50919050565b601f82111562000233576000816000526020600020601f850160051c81016020861015620005f55750805b601f850160051c820191505b81811015620006165782815560010162000601565b505050505050565b81516001600160401b038111156200063a576200063a6200031d565b62000652816200064b84546200058e565b84620005ca565b602080601f8311600181146200068a5760008415620006715750858301515b600019600386901b1c1916600185901b17855562000616565b600085815260208120601f198616915b82811015620006bb578886015182559484019460019091019084016200069a565b5085821015620006da5787850151600019600388901b60f8161c191681555b5050505050600190811b01905550565b60006020808385031215620006fe57600080fd5b82516001600160401b03808211156200071657600080fd5b818501915085601f8301126200072b57600080fd5b8151818111156200074057620007406200031d565b8060051b9150620007538483016200035f565b81815291830184019184810190888411156200076e57600080fd5b938501935b838510156200078e5784518252938501939085019062000773565b98975050505050505050565b634e487b7160e01b600052603260045260246000fd5b608051615be4620007da600039600081816106ed01528181610d1d015261199a0152615be46000f3fe608060405234801561001057600080fd5b50600436106102b85760003560e01c806371929e2a11610173578063a9bace48116100d4578063a9bace4814610677578063acf92e78146106cc578063c0a5452e146106d5578063c45a0155146106e8578063c4e41b221461070f578063ce630c1814610717578063ce884eb51461072a578063d11db83f1461073b578063d51ba5a414610750578063dc74664514610763578063e31fd42614610776578063ecccfd8e14610789578063febce3a0146107ac578063ff4d1815146107bf57600080fd5b806371929e2a14610556578063723dbbc41461055f5780637a856ee5146105725780637d894a16146105955780638052474d146105a85780638bcafc40146105b05780638c4f7355146105c357806392744b8d146105e65780639a41f721146105f9578063a01afbfb1461060c578063a0a0683e1461061f578063a204649314610648578063a655e99d1461065b578063a79d55e61461066e57600080fd5b80633796e2391161021d5780633796e2391461045a5780633bc39d5a1461046d5780633d1a03891461048057806348d9d307146104935780634c8b2126146104a85780635451e537146104c85780635e19f7cf146104db5780635e53dcb8146104e35780635f997c5b146104f657806361ac2243146104ff57806363a8374d1461051057806367511a4d14610531578063697b17031461053a5780636d10152c1461054d57600080fd5b80630274c133146102bd57806304e72b3c146102ed57806308d63e031461030257806311f50c851461031557806318160ddd146103285780631da744e61461033f578063213933e61461035f57806324481d921461037f57806325d1ef0b14610395578063270d5e2f146103c15780632938b2b5146103f55780633045aaf3146104155780633225e9e31461042a57806332530d3c14610452575b600080fd5b6002546102d0906001600160a01b031681565b6040516001600160a01b0390911681526020015b60405180910390f35b6103006102fb366004614ce4565b6107c8565b005b600f546102d0906001600160a01b031681565b6102d0610323366004614e83565b610864565b61033160095481565b6040519081526020016102e4565b61033161034d366004614e83565b60176020526000908152604090205481565b61033161036d366004614e83565b60009081526017602052604090205490565b6103876108d8565b6040516102e4929190614efb565b6102d06103a3366004614f5d565b60ff166000908152601360205260409020546001600160a01b031690565b6103e56103cf366004614e83565b60009081526019602052604090205460ff161590565b60405190151581526020016102e4565b610331610403366004614e83565b60009081526018602052604090205490565b61041d610b82565b6040516102e49190614fbe565b61043d610438366004614fd1565b610c14565b604080519283526020830191909152016102e4565b610331610c93565b6103e5610468366004615008565b610d10565b61043d61047b366004614e83565b610dc8565b61030061048e366004614e83565b610e4b565b61049b610e82565b6040516102e49190615075565b6103316104b6366004614e83565b60186020526000908152604090205481565b6103e56104d6366004614e83565b610f4e565b6103876110c4565b6103316104f1366004614e83565b61135f565b61033160035481565b601a546001600160a01b03166102d0565b61052361051e366004614e83565b61138e565b6040516102e49291906150ed565b61033160085481565b610300610548366004614e83565b6113c2565b61033160105481565b61033160075481565b61043d61056d366004614e83565b6113f2565b610585610580366004614e83565b611408565b6040516102e49493929190615101565b61043d6105a3366004614fd1565b611469565b61041d6114ac565b6103e56105be366004615247565b6114bb565b6103e56105d1366004614e83565b60196020526000908152604090205460ff1681565b6103876105f4366004614e83565b6116e3565b6103006106073660046152aa565b61198f565b61030061061a366004614e83565b611a61565b6102d061062d366004614e83565b6013602052600090815260409020546001600160a01b031681565b6011546102d0906001600160a01b031681565b6103006106693660046152c7565b611aac565b610331600c5481565b6106b1610685366004614fd1565b601260209081526000928352604080842090915290825290208054600182015460029092015490919083565b604080519384526020840192909252908201526060016102e4565b610331600d5481565b6103006106e3366004615312565b611af1565b6102d07f000000000000000000000000000000000000000000000000000000000000000081565b600954610331565b61043d610725366004614e83565b611b38565b6000546001600160a01b03166102d0565b610743611b44565b6040516102e49190615356565b61030061075e366004615415565b611c54565b61043d610771366004614e83565b612196565b601a546102d0906001600160a01b031681565b6103e5610797366004614e83565b60166020526000908152604090205460ff1681565b600e546102d0906001600160a01b031681565b610331600b5481565b601a546001600160a01b031633146107f35760405163014bdcb760e31b815260040160405180910390fd5b6107fd8383612219565b60115460405163c77b6c7f60e01b81526001600160a01b039091169063c77b6c7f9061082d908490600401614fbe565b600060405180830381600087803b15801561084757600080fd5b505af115801561085b573d6000803e3d6000fd5b50505050505050565b600080546040516311f50c8560e01b8152600481018490526001600160a01b03909116906311f50c8590602401602060405180830381865afa1580156108ae573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906108d29190615474565b92915050565b6060806000600e60009054906101000a90046001600160a01b031690506000816001600160a01b0316632d3358756040518163ffffffff1660e01b8152600401600060405180830381865afa158015610935573d6000803e3d6000fd5b505050506040513d6000823e601f3d908101601f1916820160405261095d9190810190615491565b80519091506000816001600160401b0381111561097c5761097c6149fa565b6040519080825280602002602001820160405280156109b557816020015b6109a2614891565b81526020019060019003908161099a5790505b50905060005b82811015610a345760008482815181106109d7576109d7615521565b602002602001015190506000806109ed83610dc8565b91509150604051806060016040528083815260200182815260200184815250858581518110610a1e57610a1e615521565b60209081029190910101525050506001016109bb565b506000846001600160a01b031663c6885dc76040518163ffffffff1660e01b8152600401600060405180830381865afa158015610a75573d6000803e3d6000fd5b505050506040513d6000823e601f3d908101601f19168201604052610a9d9190810190615537565b90506000836001600160401b03811115610ab957610ab96149fa565b604051908082528060200260200182016040528015610af257816020015b610adf6148b2565b815260200190600190039081610ad75790505b50905060005b84811015610b73576000838281518110610b1457610b14615521565b602002602001015190508060000151838381518110610b3557610b35615521565b602002602001015160000181815250508060400151838381518110610b5c57610b5c615521565b602090810291909101810151015250600101610af8565b50919791965090945050505050565b606060058054610b919061568a565b80601f0160208091040260200160405190810160405280929190818152602001828054610bbd9061568a565b8015610c0a5780601f10610bdf57610100808354040283529160200191610c0a565b820191906000526020600020905b815481529060010190602001808311610bed57829003601f168201915b5050505050905090565b60008181526012602090815260408083208584529091528120548190158015610c5757506000838152601260209081526040808320878452909152902060010154155b15610c685750600090506001610c8c565b50506000818152601260209081526040808320858452909152902080546001909101545b9250929050565b600e5460408051632d33587560e01b815290516000926001600160a01b03169183918391632d33587591600480830192869291908290030181865afa158015610ce0573d6000803e3d6000fd5b505050506040513d6000823e601f3d908101601f19168201604052610d089190810190615491565b519392505050565b6000336001600160a01b037f00000000000000000000000000000000000000000000000000000000000000001614610d5b576040516322a3b9b360e01b815260040160405180910390fd5b60ff821660008181526013602090815260409182902080546001600160a01b0319166001600160a01b038816908117909155915192835290917f8d4968a6ace762b6bad8887ba28d41b4673524e2af29f622792eb66e75378b45910160405180910390a250600192915050565b600c5460009081526012602090815260408083208484529091528120548190158015610e115750600c546000908152601260209081526040808320868452909152902060010154155b15610e225750600092600192509050565b5050600c5460009081526012602090815260408083209383529290522080546001909101549091565b601a546001600160a01b03163314610e765760405163014bdcb760e31b815260040160405180910390fd5b610e7f81612377565b50565b60606015805480602002602001604051908101604052809291908181526020016000905b82821015610f455760008481526020908190206040805160e08101825260068602909201805460808401908152600182015460a0850152600282015460c08501528352600381015493830193909352600483015490820152600580830154919291606084019160ff90911690811115610f2157610f2161503d565b6005811115610f3257610f3261503d565b8152505081526020019060010190610ea6565b50505050905090565b600080600080600e60009054906101000a90046001600160a01b031690506000816001600160a01b0316632d3358756040518163ffffffff1660e01b8152600401600060405180830381865afa158015610fac573d6000803e3d6000fd5b505050506040513d6000823e601f3d908101601f19168201604052610fd49190810190615491565b905060005b815181101561103a576000828281518110610ff657610ff6615521565b6020026020010151905060008061100d838b610c14565b9150915061101d888884846124b1565b90985096506110339250839150600190506156d4565b9050610fd9565b5060008061104a60105489610c14565b9150915061105a868684846124b1565b600754919750955086148015611071575084600854145b6110b65760405162461bcd60e51b81526020600482015260116024820152700acc2d8eacae640c8dedce840dac2e8c6d607b1b60448201526064015b60405180910390fd5b506001979650505050505050565b6060806000600e60009054906101000a90046001600160a01b031690506000816001600160a01b0316632d3358756040518163ffffffff1660e01b8152600401600060405180830381865afa158015611121573d6000803e3d6000fd5b505050506040513d6000823e601f3d908101601f191682016040526111499190810190615491565b80519091506000816001600160401b03811115611168576111686149fa565b6040519080825280602002602001820160405280156111a157816020015b61118e614891565b8152602001906001900390816111865790505b50905060005b828110156112205760008482815181106111c3576111c3615521565b602002602001015190506000806111d983612196565b9150915060405180606001604052808381526020018281526020018481525085858151811061120a5761120a615521565b60209081029190910101525050506001016111a7565b506000846001600160a01b031663c6885dc76040518163ffffffff1660e01b8152600401600060405180830381865afa158015611261573d6000803e3d6000fd5b505050506040513d6000823e601f3d908101601f191682016040526112899190810190615537565b90506000836001600160401b038111156112a5576112a56149fa565b6040519080825280602002602001820160405280156112de57816020015b6112cb6148b2565b8152602001906001900390816112c35790505b50905060005b84811015610b7357600083828151811061130057611300615521565b60200260200101519050806000015183838151811061132157611321615521565b60200260200101516000018181525050806040015183838151811061134857611348615521565b6020908102919091018101510152506001016112e4565b6000811561138657611381600083600080516020615b6f8339815191526125ff565b6108d2565b600092915050565b6014818154811061139e57600080fd5b60009182526020909120600390910201600181015460029091015490915060ff1682565b601a546001600160a01b031633146113ed5760405163014bdcb760e31b815260040160405180910390fd5b600d55565b6000806113fe8361263c565b9094909350915050565b6015818154811061141857600080fd5b60009182526020918290206040805160608101825260069093029091018054835260018101549383019390935260028301549082015260038201546004830154600590930154919350919060ff1684565b600080600080611478866113f2565b9150915060008061148887611b38565b9150915060008061149b868686866124b1565b909b909a5098505050505050505050565b606060048054610b919061568a565b600f5460008054604080516303408e4760e41b8152905192936001600160a01b0390811693859390911691633408e4709160048083019260209291908290030181865afa158015611510573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061153491906156e7565b600354604051630ce010b960e21b81526004810191909152602481018290529091506000906001600160a01b0384169063338042e490604401602060405180830381865afa15801561158a573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906115ae9190615700565b905080156116155760405162461bcd60e51b815260206004820152602e60248201527f546f6b656e20697320696e20667265657a652073746174757320666f7220746860448201526d1a5cc81c185c9d1a58da5c185b9d60921b60648201526084016110ad565b601a54600160a01b900460ff161561163f5760405162461bcd60e51b81526004016110ad90615722565b601a805460ff60a01b1916600160a01b179055600061165d876126d5565b905061166881612cd0565b80516116749088613376565b6116818160e00151612377565b61168c816003613813565b6116968187612219565b60e0810151600d5560405133907fe85c8c79cebe1b6656a265affa1c69c79539e5ae9a9c9229f5b5d8961978108090600090a25050601a805460ff60a01b19169055506001949350505050565b6060806000600e60009054906101000a90046001600160a01b031690506000816001600160a01b0316632d3358756040518163ffffffff1660e01b8152600401600060405180830381865afa158015611740573d6000803e3d6000fd5b505050506040513d6000823e601f3d908101601f191682016040526117689190810190615491565b80519091506000816001600160401b03811115611787576117876149fa565b6040519080825280602002602001820160405280156117c057816020015b6117ad614891565b8152602001906001900390816117a55790505b50905060005b828110156118405760008482815181106117e2576117e2615521565b602002602001015190506000806117f9838c610c14565b9150915060405180606001604052808381526020018281526020018481525085858151811061182a5761182a615521565b60209081029190910101525050506001016117c6565b506000846001600160a01b031663c6885dc76040518163ffffffff1660e01b8152600401600060405180830381865afa158015611881573d6000803e3d6000fd5b505050506040513d6000823e601f3d908101601f191682016040526118a99190810190615537565b90506000836001600160401b038111156118c5576118c56149fa565b6040519080825280602002602001820160405280156118fe57816020015b6118eb6148b2565b8152602001906001900390816118e35790505b50905060005b8481101561197f57600083828151811061192057611920615521565b60200260200101519050806000015183838151811061194157611941615521565b60200260200101516000018181525050806040015183838151811061196857611968615521565b602090810291909101810151015250600101611904565b5091989197509095505050505050565b336001600160a01b037f000000000000000000000000000000000000000000000000000000000000000016146119d8576040516322a3b9b360e01b815260040160405180910390fd5b6001600160a01b038116611a3f5760405162461bcd60e51b815260206004820152602860248201527f496e76616c69642044565020696e746567726174696f6e20636f6e7472616374604482015267206164647265737360c01b60648201526084016110ad565b601a80546001600160a01b0319166001600160a01b0392909216919091179055565b6000611a6b613a96565b9050336001600160a01b03821614611aa65760405162c2a4c960e31b81523360048201526001600160a01b03821660248201526044016110ad565b50600355565b601a546001600160a01b03163314611ad75760405163014bdcb760e31b815260040160405180910390fd5b6000611ae2826126d5565b9050611aed81612cd0565b5050565b601a546001600160a01b03163314611b1c5760405163014bdcb760e31b815260040160405180910390fd5b6000611b27836126d5565b9050611b338183613813565b505050565b6000806113fe83613b30565b60606014805480602002602001604051908101604052809291908181526020016000905b82821015610f45578382906000526020600020906003020160405180606001604052908160008201805480602002602001604051908101604052809291908181526020016000905b82821015611c005783829060005260206000209060030201604051806060016040529081600082015481526020016001820154815260200160028201548152505081526020019060010190611bb0565b5050509082525060018201546020820152600282015460409091019060ff166005811115611c3057611c3061503d565b6005811115611c4157611c4161503d565b8152505081526020019060010190611b68565b600e54600b54604051634647061560e11b815233600482015260248101919091526001600160a01b03909116906000908290638c8e0c2a90604401602060405180830381865afa158015611cac573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190611cd09190615700565b905080611d375760405162461bcd60e51b815260206004820152602f60248201527f4f6e6c7920497373756572204163636f756e7473206d617920706572666f726d60448201526e103a3437b9b29030b1ba34b7b7399760891b60648201526084016110ad565b601a54600160a01b900460ff1615611d615760405162461bcd60e51b81526004016110ad90615722565b601a805460ff60a01b1916600160a01b179055600c548411611dcf5760405162461bcd60e51b815260206004820152602160248201527f426c6f636b4e756d62657220697320616c72656164792066696e616c697365646044820152601760f91b60648201526084016110ad565b600d54841015611e185760405162461bcd60e51b815260206004820152601460248201527324b73b30b634b210213637b1b5a73ab6b132b91760611b60448201526064016110ad565b43841115611e815760405162461bcd60e51b815260206004820152603060248201527f426c6f636b4e756d62657220697320626967676572207468616e20437572726560448201526f373a10213637b1b590273ab6b132b91760811b60648201526084016110ad565b8251611ec55760405162461bcd60e51b81526020600482015260136024820152724e6f20616d6f756e7420746f2075706461746560681b60448201526064016110ad565b600183602001516005811115611edd57611edd61503d565b1480611efe5750600283602001516005811115611efc57611efc61503d565b145b611f3c5760405162461bcd60e51b815260206004820152600f60248201526e496e76616c6964207478207479706560881b60448201526064016110ad565b611f4584612377565b600080600185602001516005811115611f6057611f6061503d565b03611f76578451611f70906113f2565b90925090505b600285602001516005811115611f8e57611f8e61503d565b036120385784517f060c89ce5c263405370a08b6d0302b0bab3eedb83920ee0a677297dc392126f11015611ffb5760405162461bcd60e51b815260206004820152601460248201527304572726f723a206275726e56616c7565203e20560641b60448201526064016110ad565b84516120329061202b907f060c89ce5c263405370a08b6d0302b0bab3eedb83920ee0a677297dc392126f161576d565b6000611469565b90925090505b6015604051806080016040528060405180606001604052808681526020018581526020018b8152508152602001876000015181526020018881526020018760200151600581111561208b5761208b61503d565b9052815460018181018455600093845260209384902083518051600690940290910192835580850151838301556040908101516002840155938301516003830155928201516004820155606082015160058083018054949593949293909260ff19169184908111156120ff576120ff61503d565b02179055505050600d86905561211787838389613bbc565b60115460035460405163c651be9b60e01b81526001600160a01b039092169163c651be9b9161214e918a908a908d90600401615780565b600060405180830381600087803b15801561216857600080fd5b505af115801561217c573d6000803e3d6000fd5b5050601a805460ff60a01b19169055505050505050505050565b600d54600090815260126020908152604080832084845290915281205481901580156121df5750600d546000908152601260209081526040808320868452909152902060010154155b156121f05750600092600192509050565b5050600d5460009081526012602090815260408083209383529290522080546001909101549091565b816101000151518151146122835760405162461bcd60e51b815260206004820152602b60248201527f456e79676d6156313a20656e637279707465644d65737361676573206c656e6760448201526a0e8d040dad2e6dac2e8c6d60ab1b60648201526084016110ad565b60005b82610100015151811015611b335760115460035483516001600160a01b039092169163952acada91908590859081106122c1576122c1615521565b602002602001015186610120015185815181106122e0576122e0615521565b60200260200101518760e0015188610100015189604001518a6101000151898151811061230f5761230f615521565b60200260200101516040518863ffffffff1660e01b815260040161233997969594939291906157e8565b600060405180830381600087803b15801561235357600080fd5b505af1158015612367573d6000803e3d6000fd5b5050600190920191506122869050565b60008181526017602052604081205490036123ab57600c54600082815260176020526040902055600d546123ab9082613c15565b600c54811180156123c857506014541515806123c8575060155415155b80156123d55750600d5481115b80156123e55750600c54600d5410155b15610e7f576123f5600d54613d60565b600d80546000908152601660209081526040808320805460ff191660011790559254600c8054845260189092528383208190558082558252918120839055905461243e906116e3565b50601154600354600c54604051636924a82f60e01b81529394506001600160a01b0390921692636924a82f9261247b929187908790600401615848565b600060405180830381600087803b15801561249557600080fd5b505af11580156124a9573d6000803e3d6000fd5b505050505050565b600080851580156124c0575084155b156124cf5750829050816125f6565b831580156124db575084155b156124ea5750849050836125f6565b6000600080516020615b6f83398151915285880990506000600080516020615b6f83398151915285880990506000600080516020615b6f83398151915280838509620292f80990506000600080516020615b6f83398151915280898b09600080516020615b6f833981519152898d09089050600061258b84600080516020615b6f83398151915287620292fc09600080516020615b6f8339815191526125ff565b9050600080516020615b6f8339815191526125b7600080516020615b6f833981519152856001086141b0565b83099650600080516020615b6f8339815191526125ec6125e7600186600080516020615b6f8339815191526125ff565b6141b0565b8209955050505050505b94509492505050565b6000838381116126165761261383826156d4565b90505b828061262457612624615877565b6000612630868461576d565b089150505b9392505050565b600080827f2491aba8d3a191a76e35bc47bd9afe6cc88fee14d607cbe779f2349047d5c1577f2e07297f8d3c3d7818dbddfd24c35583f9a9d4ed0cb0c1d1348dd8f7f99152d78360015b84156126c85760018516156126a7576126a1828286866124b1565b90925090505b6126b184846141e3565b90945092506126c160028661588d565b9450612686565b9097909650945050505050565b61272e604051806101400160405280600060ff1681526020016000815260200160608152602001606081526020016060815260200160608152602001600081526020016000815260200160608152602001606081525090565b612737826141fd565b60ff168082526060830151516020830152806001600160401b03811115612760576127606149fa565b604051908082528060200260200182016040528015612789578160200160208202803683370190505b50604083015260005b818110156127e157836060015181815181106127b0576127b0615521565b6020026020010151836040015182815181106127ce576127ce615521565b6020908102919091010152600101612792565b50806001600160401b038111156127fa576127fa6149fa565b604051908082528060200260200182016040528015612823578160200160208202803683370190505b5060608301528060005b8281101561288557606085015161284482846156d4565b8151811061285457612854615521565b60200260200101518460600151828151811061287257612872615521565b602090810291909101015260010161282d565b50816001600160401b0381111561289e5761289e6149fa565b6040519080825280602002602001820160405280156128e357816020015b60408051808201909152600080825260208201528152602001906001900390816128bc5790505b50608084015260006128f68360026158af565b905060005b838110156129b2576040518060400160405280876060015183600261292091906158af565b61292a90866156d4565b8151811061293a5761293a615521565b60200260200101518152602001876060015183600261295991906158af565b61296390866156d4565b61296e9060016156d4565b8151811061297e5761297e615521565b60200260200101518152508560800151828151811061299f5761299f615521565b60209081029190910101526001016128fb565b50826001600160401b038111156129cb576129cb6149fa565b604051908082528060200260200182016040528015612a1057816020015b60408051808201909152600080825260208201528152602001906001900390816129e95790505b5060a08501526000612a238460046158af565b905060005b84811015612adf5760405180604001604052808860600151836002612a4d91906158af565b612a5790866156d4565b81518110612a6757612a67615521565b602002602001015181526020018860600151836002612a8691906158af565b612a9090866156d4565b612a9b9060016156d4565b81518110612aab57612aab615521565b60200260200101518152508660a001518281518110612acc57612acc615521565b6020908102919091010152600101612a28565b506000612aed8560066158af565b905086606001518181518110612b0557612b05615521565b602090810291909101015160c08701526060870151612b258260016156d4565b81518110612b3557612b35615521565b60200260200101518660e0018181525050846001600160401b03811115612b5e57612b5e6149fa565b604051908082528060200260200182016040528015612b87578160200160208202803683370190505b506101008701526000612b9b8660066158af565b612ba69060026156d4565b905060005b86811015612c04576060890151612bc282846156d4565b81518110612bd257612bd2615521565b60200260200101518861010001518281518110612bf157612bf1615521565b6020908102919091010152600101612bab565b50856001600160401b03811115612c1d57612c1d6149fa565b604051908082528060200260200182016040528015612c46578160200160208202803683370190505b506101208801526000612c5a8760076158af565b612c659060026156d4565b905060005b87811015612cc35760608a0151612c8182846156d4565b81518110612c9157612c91615521565b60200260200101518961012001518281518110612cb057612cb0615521565b6020908102919091010152600101612c6a565b5050505050505050919050565b8051600260ff82161015612d1c5760405162461bcd60e51b8152602060048201526013602482015272496e76616c69642076616c756520666f72206b60681b60448201526064016110ad565b612d278160086158c6565b612d329060026158e9565b60ff16826020015114612d575760405162461bcd60e51b81526004016110ad90615902565b8060ff168260a001515114612da95760405162461bcd60e51b81526020600482015260186024820152770aee4dedcce40c6dedadad2e8dacadce8e640d8cadccee8d60431b60448201526064016110ad565b8060ff168261010001515114612dff5760405162461bcd60e51b815260206004820152601b60248201527a0aee4dedcce4086d0c2d2dc92c8e64082e4e4c2f240d8cadccee8d602b1b60448201526064016110ad565b8060ff16826060015151148015612e1d57508060ff16826080015151145b8015612e3057508060ff16826040015151145b8015612e4457508060ff1682610120015151145b612e9c5760405162461bcd60e51b815260206004820152602360248201527f57726f6e67207075626c69635f7369676e616c206c656e67746820696e20707260448201526237b7b360e91b60648201526084016110ad565b60ff81166000908152601360205260409020546001600160a01b0316612f045760405162461bcd60e51b815260206004820152601c60248201527f5665726966696572206e6f742073657420666f7220676976656e206b0000000060448201526064016110ad565b600c548260e0015111612f6d5760405162461bcd60e51b815260206004820152602b60248201527f426c6f636b4e756d62657220696e2050726f6f662077617320616c726561647960448201526a103334b730b634b9b2b21760a91b60648201526084016110ad565b600d548260e001511015612fce5760405162461bcd60e51b815260206004820152602260248201527f496e76616c696420426c6f636b4e756d626572205573656420696e2050726f6f604482015261331760f11b60648201526084016110ad565b438260e0015111156130485760405162461bcd60e51b815260206004820152603e60248201527f426c6f636b4e756d626572205573656420696e2050726f6f662069732062696760448201527f676572207468616e2043757272656e7420426c6f636b204e756d6265722e000060648201526084016110ad565b60c082015160009081526019602052604090205460ff16156130c35760405162461bcd60e51b815260206004820152602e60248201527f4e756c6c696669657220616c7265616479207573656420696e2070656e64696e60448201526d33903a3930b739b0b1ba34b7b71760911b60648201526084016110ad565b6000806130d1600c546116e3565b9150915060005b8360ff1681101561336f57600061310d84876101000151848151811061310057613100615521565b6020026020010151614274565b9050600061313984886101000151858151811061312c5761312c615521565b6020026020010151614308565b90508451821061319a5760405162461bcd60e51b815260206004820152602660248201527f4d61746368696e672062616c616e636520666f7220636861696e4964206e6f7460448201526508199bdd5b9960d21b60648201526084016110ad565b835181106131fc5760405162461bcd60e51b815260206004820152602960248201527f4d61746368696e67207075626c6963206b657920666f7220636861696e4964206044820152681b9bdd08199bdd5b9960ba1b60648201526084016110ad565b83818151811061320e5761320e615521565b6020026020010151600001518760600151848151811061323057613230615521565b6020026020010151146132855760405162461bcd60e51b815260206004820152601c60248201527f496e76616c6964207075626c6963207369676e616c20666f7220706b0000000060448201526064016110ad565b84828151811061329757613297615521565b602002602001015160000151876080015184815181106132b9576132b9615521565b60200260200101516000015114801561330f57508482815181106132df576132df615521565b6020026020010151602001518760800151848151811061330157613301615521565b602002602001015160200151145b6133655760405162461bcd60e51b815260206004820152602160248201527f496e76616c6964207075626c6963207369676e616c20666f722062616c616e636044820152606560f81b60648201526084016110ad565b50506001016130d8565b5050505050565b8160ff166002036134d3578060600151516012146133d05760405162461bcd60e51b815260206004820152602480820152600080516020615b8f83398151915260448201526310359e9960e11b60648201526084016110ad565b60ff8216600090815260136020908152604091829020548351918401519284015160608501516001600160a01b0390921693634483e7219392909190613415906143a1565b6040518563ffffffff1660e01b81526004016134349493929190615990565b602060405180830381865afa158015613451573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906134759190615700565b611aed5760405162461bcd60e51b815260206004820152602960248201527f76657269667950726f6f662072657475726e65642066616c73653a20496e76616044820152683634b210383937b7b360b91b60648201526084016110ad565b8160ff1660030361359157806060015151601a1461352d5760405162461bcd60e51b815260206004820152602480820152600080516020615b8f833981519152604482015263206b3d3360e01b60648201526084016110ad565b60ff8216600090815260136020908152604091829020548351918401519284015160608501516001600160a01b039092169363514736cd939290919061357290614457565b6040518563ffffffff1660e01b815260040161343494939291906159ed565b8160ff1660040361364f578060600151516022146135eb5760405162461bcd60e51b815260206004820152602480820152600080516020615b8f833981519152604482015263081acf4d60e21b60648201526084016110ad565b60ff8216600090815260136020908152604091829020548351918401519284015160608501516001600160a01b0390921693638525a61d939290919061363090614507565b6040518563ffffffff1660e01b81526004016134349493929190615a3e565b8160ff1660050361370d57806060015151602a146136a95760405162461bcd60e51b815260206004820152602480820152600080516020615b8f833981519152604482015263206b3d3560e01b60648201526084016110ad565b60ff8216600090815260136020908152604091829020548351918401519284015160608501516001600160a01b039092169363e3ba9b6293929091906136ee906145b7565b6040518563ffffffff1660e01b81526004016134349493929190615a8f565b8160ff166006036137cb578060600151516032146137675760405162461bcd60e51b815260206004820152602480820152600080516020615b8f83398151915260448201526310359e9b60e11b60648201526084016110ad565b60ff8216600090815260136020908152604091829020548351918401519284015160608501516001600160a01b039092169363a63bed3793929091906137ac90614667565b6040518563ffffffff1660e01b81526004016134349493929190615ae0565b60405162461bcd60e51b815260206004820181905260248201527f546861742076616c7565206f66206b206973206e6f7420737570706f7274656460448201526064016110ad565b60c082015160009081526019602052604090205460ff16156138835760405162461bcd60e51b8152602060048201526024808201527f456e79676d6156313a206e756c6c696669657220616c726561647920636f6e736044820152631d5b595960e21b60648201526084016110ad565b60c08201805160009081526019602052604090819020805460ff1916600117905560e0840151915160035491519091907fcdae9b3d9cbf4b148a136938979c4aebd8a549a48f9fcf534e12883eb64fc373906138e0908690615b31565b60405180910390a46014805460018101825560009182526003027fce6d7b5282bd9a3661ae061feed1dbda4e52ab073b1f9285be6e155d9c38d4ec01905b8360a0015151811015613a5c578160000160405180606001604052808660a00151848151811061395057613950615521565b60200260200101516000015181526020018660a00151848151811061397757613977615521565b6020026020010151602001518152602001866101000151848151811061399f5761399f615521565b60209081029190910181015190915282546001818101855560009485529382902083516003909202019081559082015192810192909255604001516002909101556101008401518051613a549190839081106139fd576139fd615521565b60200260200101518560a001518381518110613a1b57613a1b615521565b6020026020010151600001518660a001518481518110613a3d57613a3d615521565b6020026020010151602001518760e0015187614717565b60010161391e565b5060c0830151600180830191909155600282018054849260ff1990911690836005811115613a8c57613a8c61503d565b0217905550505050565b600080546040516311f50c8560e01b8152600360048201526001600160a01b03909116906311f50c8590602401602060405180830381865afa158015613ae0573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190613b049190615474565b90506001600160a01b038116613b2d576040516336a41bd160e01b815260040160405180910390fd5b90565b600080827f16546696a66928d34f6be843f8a5afa2063161d92742811279454d60de5322527f109c1c7a758b3e8e54af1ce919fc24e1b986aab09a6b8082600f8694bb3c1b4b8360015b84156126c8576001851615613b9b57613b95828286866124b1565b90925090505b613ba584846141e3565b9094509250613bb560028661588d565b9450613b7a565b6000818152601260209081526040808320878452909152812080546001909101548291613bea9187876124b1565b6000948552601260209081526040808720998752989052969093209283555050600101929092555050565b600e5460408051632d33587560e01b815290516001600160a01b03909216916000918391632d3358759160048082019286929091908290030181865afa158015613c63573d6000803e3d6000fd5b505050506040513d6000823e601f3d908101601f19168201604052613c8b9190810190615491565b905060005b8151811015613d01576000828281518110613cad57613cad615521565b60200260200101519050613cc08161476f565b6000868152601260208181526040808420858552825280842080548a8652938352818520958552949091529091209081556001918201549082015501613c90565b50613d0d60105461476f565b505060009182526012602081815260408085206010805480885282855283882054968852948452828720948752848452828720959095559354855292815282842060019081015492909152919092200155565b613d6c601460006148cc565b6015546000906001600160401b03811115613d8957613d896149fa565b604051908082528060200260200182016040528015613db2578160200160208202803683370190505b5090506000805b601554811015613fd157600060158281548110613dd857613dd8615521565b60009182526020918290206040805160e0810182526006909302909101805460808401908152600182015460a0850152600282015460c08501528352600381015493830193909352600483015490820152600580830154919291606084019160ff90911690811115613e4c57613e4c61503d565b6005811115613e5d57613e5d61503d565b81525050905084816040015111613fc85760075460085482518051602090910151613e8a939291906124b1565b600855600755600181606001516005811115613ea857613ea861503d565b03613f1b57806020015160096000828254613ec391906156d4565b9091555050600c546020820151825160409081015190517feae287c62f1ff4911334dee03f631d5dded5284b1b03ea7bc1d6282916c7249f92613f0e92908252602082015260400190565b60405180910390a2613f9c565b600281606001516005811115613f3357613f3361503d565b03613f9c57806020015160096000828254613f4e919061576d565b925050819055508060000151604001517f262a9a1794440b6af993000f5805d7f51b5a19d4c32fcb10a1c5216beb0616f48260200151604051613f9391815260200190565b60405180910390a25b81848481518110613faf57613faf615521565b602090810291909101015282613fc481615b3f565b9350505b50600101613db9565b5060005b81811015614054576000838281518110613ff157613ff1615521565b602002602001015190506015818154811061400e5761400e615521565b600091825260208220600690910201818155600180820183905560028201839055600382018390556004820192909255600501805460ff19169055919091019050613fd5565b506015546000805b8281101561414b576015818154811061407757614077615521565b90600052602060002090600602016003015460001461414357601581815481106140a3576140a3615521565b9060005260206000209060060201601583815481106140c4576140c4615521565b60009182526020909120825460069092020190815560018083015481830155600280840154908301556003808401549083015560048084015490830155600580840154818401805460ff90921693909260ff1990921691849081111561412c5761412c61503d565b0217905550905050818061413f90615b3f565b9250505b60010161405c565b50805b828110156124a957601580548061416757614167615b58565b600082815260208120600660001990930192830201818155600181810183905560028201839055600382018390556004820192909255600501805460ff1916905591550161414e565b60006108d2826141cf6002600080516020615b6f83398151915261576d565b600080516020615b6f83398151915261484d565b6000806141f2848486866124b1565b915091509250929050565b60608101515160009060128190036142185750600292915050565b80601a036142295750600392915050565b8060220361423a5750600492915050565b80602a0361424b5750600592915050565b8060320361425c5750600692915050565b60405162461bcd60e51b81526004016110ad90615902565b6000805b83518110156142b4578284828151811061429457614294615521565b602002602001015160400151036142ac5790506108d2565b600101614278565b5060405162461bcd60e51b815260206004820152602260248201527f436861696e204944206e6f7420666f756e6420696e20706f696e747320617272604482015261617960f01b60648201526084016110ad565b6000805b8351811015614348578284828151811061432857614328615521565b602002602001015160200151036143405790506108d2565b60010161430c565b5060405162461bcd60e51b815260206004820152602760248201527f436861696e204944206e6f7420666f756e6420696e207075626c6963206b65796044820152667320617272617960c81b60648201526084016110ad565b6143a96148ed565b815160121461440c5760405162461bcd60e51b815260206004820152602960248201527f496e707574206172726179206d75737420686176652065786163746c7920313860448201526820656c656d656e747360b81b60648201526084016110ad565b60005b60128110156144515782818151811061442a5761442a615521565b602002602001015182826012811061444457614444615521565b602002015260010161440f565b50919050565b61445f61490c565b8151601a146144c25760405162461bcd60e51b815260206004820152602960248201527f496e707574206172726179206d75737420686176652065786163746c7920323660448201526820656c656d656e747360b81b60648201526084016110ad565b60005b601a811015614451578281815181106144e0576144e0615521565b60200260200101518282601a81106144fa576144fa615521565b60200201526001016144c5565b61450f61492b565b81516022146145725760405162461bcd60e51b815260206004820152602960248201527f496e707574206172726179206d75737420686176652065786163746c7920333460448201526820656c656d656e747360b81b60648201526084016110ad565b60005b60228110156144515782818151811061459057614590615521565b60200260200101518282602281106145aa576145aa615521565b6020020152600101614575565b6145bf61494a565b8151602a146146225760405162461bcd60e51b815260206004820152602960248201527f496e707574206172726179206d75737420686176652065786163746c7920343260448201526820656c656d656e747360b81b60648201526084016110ad565b60005b602a8110156144515782818151811061464057614640615521565b60200260200101518282602a811061465a5761465a615521565b6020020152600101614625565b61466f614969565b81516032146146d25760405162461bcd60e51b815260206004820152602960248201527f496e707574206172726179206d75737420686176652065786163746c7920353060448201526820656c656d656e747360b81b60648201526084016110ad565b60005b6032811015614451578281815181106146f0576146f0615521565b602002602001015182826032811061470a5761470a615521565b60200201526001016146d5565b61472385858585613bbc565b60048160058111156147375761473761503d565b1480614754575060058160058111156147525761475261503d565b145b1561336f5761336f6010546147688661135f565b8585613bbc565b600c5460009081526012602090815260408083208484529091529020541580156147b65750600c546000908152601260209081526040808320848452909152902060010154155b156147de57600c54600090815260126020908152604080832084845290915290206001908101555b600d5460009081526012602090815260408083208484529091529020541580156148255750600d546000908152601260209081526040808320848452909152902060010154155b15610e7f57600d54600090815260126020908152604080832093835292905220600190810155565b600060405160208152602080820152602060408201528460608201528360808201528260a082015260208160c08360055afa8080156102b857505051949350505050565b60405180606001604052806000815260200160008152602001600081525090565b604051806040016040528060008152602001600081525090565b5080546000825560030290600052602060002090810190610e7f9190614988565b6040518061024001604052806012906020820280368337509192915050565b604051806103400160405280601a906020820280368337509192915050565b6040518061044001604052806022906020820280368337509192915050565b604051806105400160405280602a906020820280368337509192915050565b6040518061064001604052806032906020820280368337509192915050565b808211156149b857600061499c82826149bc565b506000600182015560028101805460ff19169055600301614988565b5090565b5080546000825560030290600052602060002090810190610e7f91905b808211156149b85760008082556001820181905560028201556003016149d9565b634e487b7160e01b600052604160045260246000fd5b604080519081016001600160401b0381118282101715614a3257614a326149fa565b60405290565b60405161014081016001600160401b0381118282101715614a3257614a326149fa565b604051608081016001600160401b0381118282101715614a3257614a326149fa565b604051606081016001600160401b0381118282101715614a3257614a326149fa565b604051601f8201601f191681016001600160401b0381118282101715614ac757614ac76149fa565b604052919050565b803560ff81168114614ae057600080fd5b919050565b60006001600160401b03821115614afe57614afe6149fa565b5060051b60200190565b600082601f830112614b1957600080fd5b81356020614b2e614b2983614ae5565b614a9f565b8083825260208201915060208460051b870101935086841115614b5057600080fd5b602086015b84811015614b6c5780358352918301918301614b55565b509695505050505050565b600082601f830112614b8857600080fd5b81356020614b98614b2983614ae5565b82815260069290921b84018101918181019086841115614bb757600080fd5b8286015b84811015614b6c5760408189031215614bd45760008081fd5b614bdc614a10565b813581528482013585820152835291830191604001614bbb565b600082601f830112614c0757600080fd5b81356001600160401b03811115614c2057614c206149fa565b614c33601f8201601f1916602001614a9f565b818152846020838601011115614c4857600080fd5b816020850160208301376000918101602001919091529392505050565b600082601f830112614c7657600080fd5b81356020614c86614b2983614ae5565b82815260059290921b84018101918181019086841115614ca557600080fd5b8286015b84811015614b6c5780356001600160401b03811115614cc85760008081fd5b614cd68986838b0101614bf6565b845250918301918301614ca9565b600080600060608486031215614cf957600080fd5b83356001600160401b0380821115614d1057600080fd5b908501906101408288031215614d2557600080fd5b614d2d614a38565b614d3683614acf565b815260208301356020820152604083013582811115614d5457600080fd5b614d6089828601614b08565b604083015250606083013582811115614d7857600080fd5b614d8489828601614b08565b606083015250608083013582811115614d9c57600080fd5b614da889828601614b77565b60808301525060a083013582811115614dc057600080fd5b614dcc89828601614b77565b60a08301525060c083013560c082015260e083013560e08201526101008084013583811115614dfa57600080fd5b614e068a828701614b08565b8284015250506101208084013583811115614e2057600080fd5b614e2c8a828701614b08565b828401525050809550506020860135915080821115614e4a57600080fd5b614e5687838801614c65565b93506040860135915080821115614e6c57600080fd5b50614e7986828701614bf6565b9150509250925092565b600060208284031215614e9557600080fd5b5035919050565b8051825260208082015190830152604090810151910152565b60008151808452602080850194506020840160005b83811015614ef057614edd878351614e9c565b6060969096019590820190600101614eca565b509495945050505050565b60006040808352614f0f6040840186614eb5565b83810360208581019190915285518083528682019282019060005b81811015614f4f57845180518452840151848401529383019391850191600101614f2a565b509098975050505050505050565b600060208284031215614f6f57600080fd5b61263582614acf565b6000815180845260005b81811015614f9e57602081850181015186830182015201614f82565b506000602082860101526020601f19601f83011685010191505092915050565b6020815260006126356020830184614f78565b60008060408385031215614fe457600080fd5b50508035926020909101359150565b6001600160a01b0381168114610e7f57600080fd5b6000806040838503121561501b57600080fd5b823561502681614ff3565b915061503460208401614acf565b90509250929050565b634e487b7160e01b600052602160045260246000fd5b6006811061507157634e487b7160e01b600052602160045260246000fd5b9052565b602080825282518282018190526000919060409081850190868401855b828110156150e05781516150a7858251614e9c565b8087015160608681019190915286820151608087015201516150cc60a0860182615053565b5060c0939093019290850190600101615092565b5091979650505050505050565b828152604081016126356020830184615053565b60c0810161510f8287614e9c565b84606083015283608083015261512860a0830184615053565b95945050505050565b600082601f83011261514257600080fd5b61514a614a10565b80604084018581111561515c57600080fd5b845b8181101561517657803584526020938401930161515e565b509095945050505050565b6000610120828403121561519457600080fd5b61519c614a5b565b90506151a88383615131565b8152604083605f8401126151bb57600080fd5b6151c3614a10565b8060c08501868111156151d557600080fd5b604086015b818110156151fb576151ec8882615131565b845260209093019284016151da565b5081602086015261520c8782615131565b6040860152505050506101008201356001600160401b0381111561522f57600080fd5b61523b84828501614b08565b60608301525092915050565b6000806040838503121561525a57600080fd5b82356001600160401b038082111561527157600080fd5b61527d86838701615181565b9350602085013591508082111561529357600080fd5b506152a085828601614c65565b9150509250929050565b6000602082840312156152bc57600080fd5b813561263581614ff3565b6000602082840312156152d957600080fd5b81356001600160401b038111156152ef57600080fd5b6152fb84828501615181565b949350505050565b803560068110614ae057600080fd5b6000806040838503121561532557600080fd5b82356001600160401b0381111561533b57600080fd5b61534785828601615181565b92505061503460208401615303565b600060208083018184528085518083526040925060408601915060408160051b8701018488016000805b8481101561540657898403603f19018652825180516060808752815181880181905290918b0190859060808901905b808310156153d6576153c2828551614e9c565b928d019260019290920191908401906153af565b50848d0151898e0152938b0151936153f08c8a0186615053565b998c019997505050938901935050600101615380565b50919998505050505050505050565b6000806000838503608081121561542b57600080fd5b84359350602085013592506040603f198201121561544857600080fd5b50615451614a10565b6040850135815261546460608601615303565b6020820152809150509250925092565b60006020828403121561548657600080fd5b815161263581614ff3565b600060208083850312156154a457600080fd5b82516001600160401b038111156154ba57600080fd5b8301601f810185136154cb57600080fd5b80516154d9614b2982614ae5565b81815260059190911b820183019083810190878311156154f857600080fd5b928401925b82841015615516578351825292840192908401906154fd565b979650505050505050565b634e487b7160e01b600052603260045260246000fd5b6000602080838503121561554a57600080fd5b82516001600160401b038082111561556157600080fd5b818501915085601f83011261557557600080fd5b8151615583614b2982614ae5565b81815260059190911b830184019084810190888311156155a257600080fd5b8585015b8381101561567d578051858111156155bd57600080fd5b86016060818c03601f190112156155d357600080fd5b6155db614a7d565b8882015181526040820151878111156155f357600080fd5b8201603f81018d1361560457600080fd5b89810151615614614b2982614ae5565b81815260059190911b8201604001908b8101908f83111561563457600080fd5b6040840193505b8284101561565d57835161564e81614ff3565b8252928c0192908c019061563b565b848d015250505060609190910151604082015283529186019186016155a6565b5098975050505050505050565b600181811c9082168061569e57607f821691505b60208210810361445157634e487b7160e01b600052602260045260246000fd5b634e487b7160e01b600052601160045260246000fd5b808201808211156108d2576108d26156be565b6000602082840312156156f957600080fd5b5051919050565b60006020828403121561571257600080fd5b8151801515811461263557600080fd5b6020808252602b908201527f436f6e74726163742069732070726f63657373696e6720616e6f74686572207460408201526a3930b739b0b1ba34b7b71760a91b606082015260800190565b818103818111156108d2576108d26156be565b84815260208082018590528351604083015283015160a08201906157a76060840182615053565b5082608083015295945050505050565b60008151808452602080850194506020840160005b83811015614ef0578151875295820195908201906001016157cc565b87815260e06020820152600061580160e0830189614f78565b876040840152866060840152828103608084015261581f81876157b7565b905082810360a084015261583381866157b7565b9150508260c083015298975050505050505050565b84815283602082015282604082015260806060820152600061586d6080830184614eb5565b9695505050505050565b634e487b7160e01b600052601260045260246000fd5b6000826158aa57634e487b7160e01b600052601260045260246000fd5b500490565b80820281158282048414176108d2576108d26156be565b60ff81811683821602908116908181146158e2576158e26156be565b5092915050565b60ff81811683821601908111156108d2576108d26156be565b6020808252601c908201527f496e76616c6964207075626c69635f7369676e616c206c656e67746800000000604082015260600190565b8060005b600281101561595c57815184526020938401939091019060010161593d565b50505050565b8060005b600281101561595c5761597a848351615939565b6040939093019260209190910190600101615966565b610340810161599f8287615939565b6159ac6040830186615962565b6159b960c0830185615939565b61010082018360005b60128110156159e15781518352602092830192909101906001016159c2565b50505095945050505050565b61044081016159fc8287615939565b615a096040830186615962565b615a1660c0830185615939565b61010082018360005b601a8110156159e1578151835260209283019290910190600101615a1f565b6105408101615a4d8287615939565b615a5a6040830186615962565b615a6760c0830185615939565b61010082018360005b60228110156159e1578151835260209283019290910190600101615a70565b6106408101615a9e8287615939565b615aab6040830186615962565b615ab860c0830185615939565b61010082018360005b602a8110156159e1578151835260209283019290910190600101615ac1565b6107408101615aef8287615939565b615afc6040830186615962565b615b0960c0830185615939565b61010082018360005b60328110156159e1578151835260209283019290910190600101615b12565b602081016108d28284615053565b600060018201615b5157615b516156be565b5060010190565b634e487b7160e01b600052603160045260246000fdfe30644e72e131a029b85045b68181585d2833e84879b9709143e1f593f0000001496e76616c6964207075626c69635f7369676e616c206c656e67746820666f72a2646970667358221220d9b0ec10926ae4f2c2df663dd16d987769099b1181fd8e06cc2cd1d52850cfa064736f6c63430008180033",
}

// EnygmaV1 is an auto generated Go binding around an Ethereum contract.
type EnygmaV1 struct {
	abi abi.ABI
}

// NewEnygmaV1 creates a new instance of EnygmaV1.
func NewEnygmaV1() *EnygmaV1 {
	parsed, err := EnygmaV1MetaData.ParseABI()
	if err != nil {
		panic(errors.New("invalid ABI: " + err.Error()))
	}
	return &EnygmaV1{abi: *parsed}
}

// Instance creates a wrapper for a deployed contract instance at the given address.
// Use this to create the instance object passed to abigen v2 library functions Call, Transact, etc.
func (c *EnygmaV1) Instance(backend bind.ContractBackend, addr common.Address) *bind.BoundContract {
	return bind.NewBoundContract(addr, c.abi, backend, backend, backend)
}

// PackConstructor is the Go binding used to pack the parameters required for
// contract deployment.
//
// Solidity: constructor((string,string,uint8,bytes32,address,uint256,address,address,address,address,address) params) returns()
func (enygmaV1 *EnygmaV1) PackConstructor(params EnygmaCreationParams) []byte {
	enc, err := enygmaV1.abi.Pack("", params)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackName is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x8052474d.
//
// Solidity: function Name() view returns(string)
func (enygmaV1 *EnygmaV1) PackName() []byte {
	enc, err := enygmaV1.abi.Pack("Name")
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackName is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x8052474d.
//
// Solidity: function Name() view returns(string)
func (enygmaV1 *EnygmaV1) UnpackName(data []byte) (string, error) {
	out, err := enygmaV1.abi.Unpack("Name", data)
	if err != nil {
		return *new(string), err
	}
	out0 := *abi.ConvertType(out[0], new(string)).(*string)
	return out0, err
}

// PackSymbol is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x3045aaf3.
//
// Solidity: function Symbol() view returns(string)
func (enygmaV1 *EnygmaV1) PackSymbol() []byte {
	enc, err := enygmaV1.abi.Pack("Symbol")
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackSymbol is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x3045aaf3.
//
// Solidity: function Symbol() view returns(string)
func (enygmaV1 *EnygmaV1) UnpackSymbol(data []byte) (string, error) {
	out, err := enygmaV1.abi.Unpack("Symbol", data)
	if err != nil {
		return *new(string), err
	}
	out0 := *abi.ConvertType(out[0], new(string)).(*string)
	return out0, err
}

// PackAddTransferVerifier is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x3796e239.
//
// Solidity: function addTransferVerifier(address verifier, uint8 k) returns(bool)
func (enygmaV1 *EnygmaV1) PackAddTransferVerifier(verifier common.Address, k uint8) []byte {
	enc, err := enygmaV1.abi.Pack("addTransferVerifier", verifier, k)
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackAddTransferVerifier is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x3796e239.
//
// Solidity: function addTransferVerifier(address verifier, uint8 k) returns(bool)
func (enygmaV1 *EnygmaV1) UnpackAddTransferVerifier(data []byte) (bool, error) {
	out, err := enygmaV1.abi.Unpack("addTransferVerifier", data)
	if err != nil {
		return *new(bool), err
	}
	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	return out0, err
}

// PackCheckTotalSumOfBalances is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x5451e537.
//
// Solidity: function checkTotalSumOfBalances(uint256 blockNumber) view returns(bool)
func (enygmaV1 *EnygmaV1) PackCheckTotalSumOfBalances(blockNumber *big.Int) []byte {
	enc, err := enygmaV1.abi.Pack("checkTotalSumOfBalances", blockNumber)
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackCheckTotalSumOfBalances is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x5451e537.
//
// Solidity: function checkTotalSumOfBalances(uint256 blockNumber) view returns(bool)
func (enygmaV1 *EnygmaV1) UnpackCheckTotalSumOfBalances(data []byte) (bool, error) {
	out, err := enygmaV1.abi.Unpack("checkTotalSumOfBalances", data)
	if err != nil {
		return *new(bool), err
	}
	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	return out0, err
}

// PackConsumedNullifiers is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x8c4f7355.
//
// Solidity: function consumedNullifiers(uint256 ) view returns(bool)
func (enygmaV1 *EnygmaV1) PackConsumedNullifiers(arg0 *big.Int) []byte {
	enc, err := enygmaV1.abi.Pack("consumedNullifiers", arg0)
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackConsumedNullifiers is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x8c4f7355.
//
// Solidity: function consumedNullifiers(uint256 ) view returns(bool)
func (enygmaV1 *EnygmaV1) UnpackConsumedNullifiers(data []byte) (bool, error) {
	out, err := enygmaV1.abi.Unpack("consumedNullifiers", data)
	if err != nil {
		return *new(bool), err
	}
	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	return out0, err
}

// PackDerivePk is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x723dbbc4.
//
// Solidity: function derivePk(uint256 v) view returns(uint256 x2, uint256 y2)
func (enygmaV1 *EnygmaV1) PackDerivePk(v *big.Int) []byte {
	enc, err := enygmaV1.abi.Pack("derivePk", v)
	if err != nil {
		panic(err)
	}
	return enc
}

// DerivePkOutput serves as a container for the return parameters of contract
// method DerivePk.
type DerivePkOutput struct {
	X2 *big.Int
	Y2 *big.Int
}

// UnpackDerivePk is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x723dbbc4.
//
// Solidity: function derivePk(uint256 v) view returns(uint256 x2, uint256 y2)
func (enygmaV1 *EnygmaV1) UnpackDerivePk(data []byte) (DerivePkOutput, error) {
	out, err := enygmaV1.abi.Unpack("derivePk", data)
	outstruct := new(DerivePkOutput)
	if err != nil {
		return *outstruct, err
	}
	outstruct.X2 = abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	outstruct.Y2 = abi.ConvertType(out[1], new(big.Int)).(*big.Int)
	return *outstruct, err

}

// PackDerivePkH is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xce630c18.
//
// Solidity: function derivePkH(uint256 r) view returns(uint256 x2, uint256 y2)
func (enygmaV1 *EnygmaV1) PackDerivePkH(r *big.Int) []byte {
	enc, err := enygmaV1.abi.Pack("derivePkH", r)
	if err != nil {
		panic(err)
	}
	return enc
}

// DerivePkHOutput serves as a container for the return parameters of contract
// method DerivePkH.
type DerivePkHOutput struct {
	X2 *big.Int
	Y2 *big.Int
}

// UnpackDerivePkH is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xce630c18.
//
// Solidity: function derivePkH(uint256 r) view returns(uint256 x2, uint256 y2)
func (enygmaV1 *EnygmaV1) UnpackDerivePkH(data []byte) (DerivePkHOutput, error) {
	out, err := enygmaV1.abi.Unpack("derivePkH", data)
	outstruct := new(DerivePkHOutput)
	if err != nil {
		return *outstruct, err
	}
	outstruct.X2 = abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	outstruct.Y2 = abi.ConvertType(out[1], new(big.Int)).(*big.Int)
	return *outstruct, err

}

// PackDvpAddPendingTransaction is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xc0a5452e.
//
// Solidity: function dvpAddPendingTransaction((uint256[2],uint256[2][2],uint256[2],uint256[]) proof, uint8 transactionType) returns()
func (enygmaV1 *EnygmaV1) PackDvpAddPendingTransaction(proof IEnygmaV1TransferProof, transactionType uint8) []byte {
	enc, err := enygmaV1.abi.Pack("dvpAddPendingTransaction", proof, transactionType)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackDvpChainId is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x6d10152c.
//
// Solidity: function dvpChainId() view returns(uint256)
func (enygmaV1 *EnygmaV1) PackDvpChainId() []byte {
	enc, err := enygmaV1.abi.Pack("dvpChainId")
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackDvpChainId is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x6d10152c.
//
// Solidity: function dvpChainId() view returns(uint256)
func (enygmaV1 *EnygmaV1) UnpackDvpChainId(data []byte) (*big.Int, error) {
	out, err := enygmaV1.abi.Unpack("dvpChainId", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, err
}

// PackDvpFinalisePendingTransactions is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x3d1a0389.
//
// Solidity: function dvpFinalisePendingTransactions(uint256 currentBlockNumber) returns()
func (enygmaV1 *EnygmaV1) PackDvpFinalisePendingTransactions(currentBlockNumber *big.Int) []byte {
	enc, err := enygmaV1.abi.Pack("dvpFinalisePendingTransactions", currentBlockNumber)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackDvpIntegrationContractAddress is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xe31fd426.
//
// Solidity: function dvpIntegrationContractAddress() view returns(address)
func (enygmaV1 *EnygmaV1) PackDvpIntegrationContractAddress() []byte {
	enc, err := enygmaV1.abi.Pack("dvpIntegrationContractAddress")
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackDvpIntegrationContractAddress is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xe31fd426.
//
// Solidity: function dvpIntegrationContractAddress() view returns(address)
func (enygmaV1 *EnygmaV1) UnpackDvpIntegrationContractAddress(data []byte) (common.Address, error) {
	out, err := enygmaV1.abi.Unpack("dvpIntegrationContractAddress", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, err
}

// PackDvpSendEvents is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x04e72b3c.
//
// Solidity: function dvpSendEvents((uint8,uint256,uint256[],uint256[],(uint256,uint256)[],(uint256,uint256)[],uint256,uint256,uint256[],uint256[]) proofData, bytes[] encryptedMessages, bytes encryptedUpdate) returns()
func (enygmaV1 *EnygmaV1) PackDvpSendEvents(proofData IEnygmaV1ExtractedProofData, encryptedMessages [][]byte, encryptedUpdate []byte) []byte {
	enc, err := enygmaV1.abi.Pack("dvpSendEvents", proofData, encryptedMessages, encryptedUpdate)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackDvpSetLastblockNumPending is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x697b1703.
//
// Solidity: function dvpSetLastblockNumPending(uint256 newValue) returns()
func (enygmaV1 *EnygmaV1) PackDvpSetLastblockNumPending(newValue *big.Int) []byte {
	enc, err := enygmaV1.abi.Pack("dvpSetLastblockNumPending", newValue)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackDvpValidateTransferInputs is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xa655e99d.
//
// Solidity: function dvpValidateTransferInputs((uint256[2],uint256[2][2],uint256[2],uint256[]) proof) view returns()
func (enygmaV1 *EnygmaV1) PackDvpValidateTransferInputs(proof IEnygmaV1TransferProof) []byte {
	enc, err := enygmaV1.abi.Pack("dvpValidateTransferInputs", proof)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackEnygmaTeleport is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xa2046493.
//
// Solidity: function enygmaTeleport() view returns(address)
func (enygmaV1 *EnygmaV1) PackEnygmaTeleport() []byte {
	enc, err := enygmaV1.abi.Pack("enygmaTeleport")
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackEnygmaTeleport is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xa2046493.
//
// Solidity: function enygmaTeleport() view returns(address)
func (enygmaV1 *EnygmaV1) UnpackEnygmaTeleport(data []byte) (common.Address, error) {
	out, err := enygmaV1.abi.Unpack("enygmaTeleport", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, err
}

// PackFactory is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xc45a0155.
//
// Solidity: function factory() view returns(address)
func (enygmaV1 *EnygmaV1) PackFactory() []byte {
	enc, err := enygmaV1.abi.Pack("factory")
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackFactory is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xc45a0155.
//
// Solidity: function factory() view returns(address)
func (enygmaV1 *EnygmaV1) UnpackFactory(data []byte) (common.Address, error) {
	out, err := enygmaV1.abi.Unpack("factory", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, err
}

// PackGetAddressByResourceId is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x11f50c85.
//
// Solidity: function getAddressByResourceId(bytes32 _resourceId) view returns(address)
func (enygmaV1 *EnygmaV1) PackGetAddressByResourceId(resourceId [32]byte) []byte {
	enc, err := enygmaV1.abi.Pack("getAddressByResourceId", resourceId)
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackGetAddressByResourceId is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x11f50c85.
//
// Solidity: function getAddressByResourceId(bytes32 _resourceId) view returns(address)
func (enygmaV1 *EnygmaV1) UnpackGetAddressByResourceId(data []byte) (common.Address, error) {
	out, err := enygmaV1.abi.Unpack("getAddressByResourceId", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, err
}

// PackGetBalanceByBlockNumber is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x3225e9e3.
//
// Solidity: function getBalanceByBlockNumber(uint256 chainId, uint256 blockNumber) view returns(uint256 x, uint256 y)
func (enygmaV1 *EnygmaV1) PackGetBalanceByBlockNumber(chainId *big.Int, blockNumber *big.Int) []byte {
	enc, err := enygmaV1.abi.Pack("getBalanceByBlockNumber", chainId, blockNumber)
	if err != nil {
		panic(err)
	}
	return enc
}

// GetBalanceByBlockNumberOutput serves as a container for the return parameters of contract
// method GetBalanceByBlockNumber.
type GetBalanceByBlockNumberOutput struct {
	X *big.Int
	Y *big.Int
}

// UnpackGetBalanceByBlockNumber is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x3225e9e3.
//
// Solidity: function getBalanceByBlockNumber(uint256 chainId, uint256 blockNumber) view returns(uint256 x, uint256 y)
func (enygmaV1 *EnygmaV1) UnpackGetBalanceByBlockNumber(data []byte) (GetBalanceByBlockNumberOutput, error) {
	out, err := enygmaV1.abi.Unpack("getBalanceByBlockNumber", data)
	outstruct := new(GetBalanceByBlockNumberOutput)
	if err != nil {
		return *outstruct, err
	}
	outstruct.X = abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	outstruct.Y = abi.ConvertType(out[1], new(big.Int)).(*big.Int)
	return *outstruct, err

}

// PackGetBalanceFinalised is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x3bc39d5a.
//
// Solidity: function getBalanceFinalised(uint256 chainId) view returns(uint256 x, uint256 y)
func (enygmaV1 *EnygmaV1) PackGetBalanceFinalised(chainId *big.Int) []byte {
	enc, err := enygmaV1.abi.Pack("getBalanceFinalised", chainId)
	if err != nil {
		panic(err)
	}
	return enc
}

// GetBalanceFinalisedOutput serves as a container for the return parameters of contract
// method GetBalanceFinalised.
type GetBalanceFinalisedOutput struct {
	X *big.Int
	Y *big.Int
}

// UnpackGetBalanceFinalised is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x3bc39d5a.
//
// Solidity: function getBalanceFinalised(uint256 chainId) view returns(uint256 x, uint256 y)
func (enygmaV1 *EnygmaV1) UnpackGetBalanceFinalised(data []byte) (GetBalanceFinalisedOutput, error) {
	out, err := enygmaV1.abi.Unpack("getBalanceFinalised", data)
	outstruct := new(GetBalanceFinalisedOutput)
	if err != nil {
		return *outstruct, err
	}
	outstruct.X = abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	outstruct.Y = abi.ConvertType(out[1], new(big.Int)).(*big.Int)
	return *outstruct, err

}

// PackGetBalancePending is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xdc746645.
//
// Solidity: function getBalancePending(uint256 chainId) view returns(uint256 x, uint256 y)
func (enygmaV1 *EnygmaV1) PackGetBalancePending(chainId *big.Int) []byte {
	enc, err := enygmaV1.abi.Pack("getBalancePending", chainId)
	if err != nil {
		panic(err)
	}
	return enc
}

// GetBalancePendingOutput serves as a container for the return parameters of contract
// method GetBalancePending.
type GetBalancePendingOutput struct {
	X *big.Int
	Y *big.Int
}

// UnpackGetBalancePending is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xdc746645.
//
// Solidity: function getBalancePending(uint256 chainId) view returns(uint256 x, uint256 y)
func (enygmaV1 *EnygmaV1) UnpackGetBalancePending(data []byte) (GetBalancePendingOutput, error) {
	out, err := enygmaV1.abi.Unpack("getBalancePending", data)
	outstruct := new(GetBalancePendingOutput)
	if err != nil {
		return *outstruct, err
	}
	outstruct.X = abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	outstruct.Y = abi.ConvertType(out[1], new(big.Int)).(*big.Int)
	return *outstruct, err

}

// PackGetDvpIntegrationContractAddress is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x61ac2243.
//
// Solidity: function getDvpIntegrationContractAddress() view returns(address)
func (enygmaV1 *EnygmaV1) PackGetDvpIntegrationContractAddress() []byte {
	enc, err := enygmaV1.abi.Pack("getDvpIntegrationContractAddress")
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackGetDvpIntegrationContractAddress is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x61ac2243.
//
// Solidity: function getDvpIntegrationContractAddress() view returns(address)
func (enygmaV1 *EnygmaV1) UnpackGetDvpIntegrationContractAddress(data []byte) (common.Address, error) {
	out, err := enygmaV1.abi.Unpack("getDvpIntegrationContractAddress", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, err
}

// PackGetEndpointAddress is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xce884eb5.
//
// Solidity: function getEndpointAddress() view returns(address)
func (enygmaV1 *EnygmaV1) PackGetEndpointAddress() []byte {
	enc, err := enygmaV1.abi.Pack("getEndpointAddress")
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackGetEndpointAddress is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xce884eb5.
//
// Solidity: function getEndpointAddress() view returns(address)
func (enygmaV1 *EnygmaV1) UnpackGetEndpointAddress(data []byte) (common.Address, error) {
	out, err := enygmaV1.abi.Unpack("getEndpointAddress", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, err
}

// PackGetLastblockNumAtCurrentBlockNumber is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x213933e6.
//
// Solidity: function getLastblockNumAtCurrentBlockNumber(uint256 currentBlockNumber) view returns(uint256)
func (enygmaV1 *EnygmaV1) PackGetLastblockNumAtCurrentBlockNumber(currentBlockNumber *big.Int) []byte {
	enc, err := enygmaV1.abi.Pack("getLastblockNumAtCurrentBlockNumber", currentBlockNumber)
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackGetLastblockNumAtCurrentBlockNumber is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x213933e6.
//
// Solidity: function getLastblockNumAtCurrentBlockNumber(uint256 currentBlockNumber) view returns(uint256)
func (enygmaV1 *EnygmaV1) UnpackGetLastblockNumAtCurrentBlockNumber(data []byte) (*big.Int, error) {
	out, err := enygmaV1.abi.Unpack("getLastblockNumAtCurrentBlockNumber", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, err
}

// PackGetNextBlockNumberToFinaliseAfter is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x2938b2b5.
//
// Solidity: function getNextBlockNumberToFinaliseAfter(uint256 blockNumber) view returns(uint256)
func (enygmaV1 *EnygmaV1) PackGetNextBlockNumberToFinaliseAfter(blockNumber *big.Int) []byte {
	enc, err := enygmaV1.abi.Pack("getNextBlockNumberToFinaliseAfter", blockNumber)
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackGetNextBlockNumberToFinaliseAfter is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x2938b2b5.
//
// Solidity: function getNextBlockNumberToFinaliseAfter(uint256 blockNumber) view returns(uint256)
func (enygmaV1 *EnygmaV1) UnpackGetNextBlockNumberToFinaliseAfter(data []byte) (*big.Int, error) {
	out, err := enygmaV1.abi.Unpack("getNextBlockNumberToFinaliseAfter", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, err
}

// PackGetPendingMintsAndBurns is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x48d9d307.
//
// Solidity: function getPendingMintsAndBurns() view returns(((uint256,uint256,uint256),uint256,uint256,uint8)[])
func (enygmaV1 *EnygmaV1) PackGetPendingMintsAndBurns() []byte {
	enc, err := enygmaV1.abi.Pack("getPendingMintsAndBurns")
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackGetPendingMintsAndBurns is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x48d9d307.
//
// Solidity: function getPendingMintsAndBurns() view returns(((uint256,uint256,uint256),uint256,uint256,uint8)[])
func (enygmaV1 *EnygmaV1) UnpackGetPendingMintsAndBurns(data []byte) ([]IEnygmaV1PendingMintOrBurn, error) {
	out, err := enygmaV1.abi.Unpack("getPendingMintsAndBurns", data)
	if err != nil {
		return *new([]IEnygmaV1PendingMintOrBurn), err
	}
	out0 := *abi.ConvertType(out[0], new([]IEnygmaV1PendingMintOrBurn)).(*[]IEnygmaV1PendingMintOrBurn)
	return out0, err
}

// PackGetPendingTransactions is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xd11db83f.
//
// Solidity: function getPendingTransactions() view returns(((uint256,uint256,uint256)[],uint256,uint8)[])
func (enygmaV1 *EnygmaV1) PackGetPendingTransactions() []byte {
	enc, err := enygmaV1.abi.Pack("getPendingTransactions")
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackGetPendingTransactions is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xd11db83f.
//
// Solidity: function getPendingTransactions() view returns(((uint256,uint256,uint256)[],uint256,uint8)[])
func (enygmaV1 *EnygmaV1) UnpackGetPendingTransactions(data []byte) ([]IEnygmaV1PendingTransaction, error) {
	out, err := enygmaV1.abi.Unpack("getPendingTransactions", data)
	if err != nil {
		return *new([]IEnygmaV1PendingTransaction), err
	}
	out0 := *abi.ConvertType(out[0], new([]IEnygmaV1PendingTransaction)).(*[]IEnygmaV1PendingTransaction)
	return out0, err
}

// PackGetPublicValuesByBlockNumber is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x92744b8d.
//
// Solidity: function getPublicValuesByBlockNumber(uint256 blockNumber) view returns((uint256,uint256,uint256)[], (uint256,uint256)[])
func (enygmaV1 *EnygmaV1) PackGetPublicValuesByBlockNumber(blockNumber *big.Int) []byte {
	enc, err := enygmaV1.abi.Pack("getPublicValuesByBlockNumber", blockNumber)
	if err != nil {
		panic(err)
	}
	return enc
}

// GetPublicValuesByBlockNumberOutput serves as a container for the return parameters of contract
// method GetPublicValuesByBlockNumber.
type GetPublicValuesByBlockNumberOutput struct {
	Arg0 []IEnygmaV1EnygmaPointWithChainId
	Arg1 []IEnygmaV1EnygmaPublicKeyWithChainId
}

// UnpackGetPublicValuesByBlockNumber is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x92744b8d.
//
// Solidity: function getPublicValuesByBlockNumber(uint256 blockNumber) view returns((uint256,uint256,uint256)[], (uint256,uint256)[])
func (enygmaV1 *EnygmaV1) UnpackGetPublicValuesByBlockNumber(data []byte) (GetPublicValuesByBlockNumberOutput, error) {
	out, err := enygmaV1.abi.Unpack("getPublicValuesByBlockNumber", data)
	outstruct := new(GetPublicValuesByBlockNumberOutput)
	if err != nil {
		return *outstruct, err
	}
	outstruct.Arg0 = *abi.ConvertType(out[0], new([]IEnygmaV1EnygmaPointWithChainId)).(*[]IEnygmaV1EnygmaPointWithChainId)
	outstruct.Arg1 = *abi.ConvertType(out[1], new([]IEnygmaV1EnygmaPublicKeyWithChainId)).(*[]IEnygmaV1EnygmaPublicKeyWithChainId)
	return *outstruct, err

}

// PackGetPublicValuesFinalised is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x24481d92.
//
// Solidity: function getPublicValuesFinalised() view returns((uint256,uint256,uint256)[], (uint256,uint256)[])
func (enygmaV1 *EnygmaV1) PackGetPublicValuesFinalised() []byte {
	enc, err := enygmaV1.abi.Pack("getPublicValuesFinalised")
	if err != nil {
		panic(err)
	}
	return enc
}

// GetPublicValuesFinalisedOutput serves as a container for the return parameters of contract
// method GetPublicValuesFinalised.
type GetPublicValuesFinalisedOutput struct {
	Arg0 []IEnygmaV1EnygmaPointWithChainId
	Arg1 []IEnygmaV1EnygmaPublicKeyWithChainId
}

// UnpackGetPublicValuesFinalised is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x24481d92.
//
// Solidity: function getPublicValuesFinalised() view returns((uint256,uint256,uint256)[], (uint256,uint256)[])
func (enygmaV1 *EnygmaV1) UnpackGetPublicValuesFinalised(data []byte) (GetPublicValuesFinalisedOutput, error) {
	out, err := enygmaV1.abi.Unpack("getPublicValuesFinalised", data)
	outstruct := new(GetPublicValuesFinalisedOutput)
	if err != nil {
		return *outstruct, err
	}
	outstruct.Arg0 = *abi.ConvertType(out[0], new([]IEnygmaV1EnygmaPointWithChainId)).(*[]IEnygmaV1EnygmaPointWithChainId)
	outstruct.Arg1 = *abi.ConvertType(out[1], new([]IEnygmaV1EnygmaPublicKeyWithChainId)).(*[]IEnygmaV1EnygmaPublicKeyWithChainId)
	return *outstruct, err

}

// PackGetPublicValuesPending is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x5e19f7cf.
//
// Solidity: function getPublicValuesPending() view returns((uint256,uint256,uint256)[], (uint256,uint256)[])
func (enygmaV1 *EnygmaV1) PackGetPublicValuesPending() []byte {
	enc, err := enygmaV1.abi.Pack("getPublicValuesPending")
	if err != nil {
		panic(err)
	}
	return enc
}

// GetPublicValuesPendingOutput serves as a container for the return parameters of contract
// method GetPublicValuesPending.
type GetPublicValuesPendingOutput struct {
	Arg0 []IEnygmaV1EnygmaPointWithChainId
	Arg1 []IEnygmaV1EnygmaPublicKeyWithChainId
}

// UnpackGetPublicValuesPending is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x5e19f7cf.
//
// Solidity: function getPublicValuesPending() view returns((uint256,uint256,uint256)[], (uint256,uint256)[])
func (enygmaV1 *EnygmaV1) UnpackGetPublicValuesPending(data []byte) (GetPublicValuesPendingOutput, error) {
	out, err := enygmaV1.abi.Unpack("getPublicValuesPending", data)
	outstruct := new(GetPublicValuesPendingOutput)
	if err != nil {
		return *outstruct, err
	}
	outstruct.Arg0 = *abi.ConvertType(out[0], new([]IEnygmaV1EnygmaPointWithChainId)).(*[]IEnygmaV1EnygmaPointWithChainId)
	outstruct.Arg1 = *abi.ConvertType(out[1], new([]IEnygmaV1EnygmaPublicKeyWithChainId)).(*[]IEnygmaV1EnygmaPublicKeyWithChainId)
	return *outstruct, err

}

// PackGetTotalRegisteredBanks is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x32530d3c.
//
// Solidity: function getTotalRegisteredBanks() view returns(uint256)
func (enygmaV1 *EnygmaV1) PackGetTotalRegisteredBanks() []byte {
	enc, err := enygmaV1.abi.Pack("getTotalRegisteredBanks")
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackGetTotalRegisteredBanks is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x32530d3c.
//
// Solidity: function getTotalRegisteredBanks() view returns(uint256)
func (enygmaV1 *EnygmaV1) UnpackGetTotalRegisteredBanks(data []byte) (*big.Int, error) {
	out, err := enygmaV1.abi.Unpack("getTotalRegisteredBanks", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, err
}

// PackGetTotalSupply is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xc4e41b22.
//
// Solidity: function getTotalSupply() view returns(uint256)
func (enygmaV1 *EnygmaV1) PackGetTotalSupply() []byte {
	enc, err := enygmaV1.abi.Pack("getTotalSupply")
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackGetTotalSupply is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xc4e41b22.
//
// Solidity: function getTotalSupply() view returns(uint256)
func (enygmaV1 *EnygmaV1) UnpackGetTotalSupply(data []byte) (*big.Int, error) {
	out, err := enygmaV1.abi.Unpack("getTotalSupply", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, err
}

// PackGetTransferVerifierAddress is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x25d1ef0b.
//
// Solidity: function getTransferVerifierAddress(uint8 k) view returns(address)
func (enygmaV1 *EnygmaV1) PackGetTransferVerifierAddress(k uint8) []byte {
	enc, err := enygmaV1.abi.Pack("getTransferVerifierAddress", k)
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackGetTransferVerifierAddress is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x25d1ef0b.
//
// Solidity: function getTransferVerifierAddress(uint8 k) view returns(address)
func (enygmaV1 *EnygmaV1) UnpackGetTransferVerifierAddress(data []byte) (common.Address, error) {
	out, err := enygmaV1.abi.Unpack("getTransferVerifierAddress", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, err
}

// PackIsNullifierUnspent is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x270d5e2f.
//
// Solidity: function isNullifierUnspent(uint256 nullifier) view returns(bool)
func (enygmaV1 *EnygmaV1) PackIsNullifierUnspent(nullifier *big.Int) []byte {
	enc, err := enygmaV1.abi.Pack("isNullifierUnspent", nullifier)
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackIsNullifierUnspent is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x270d5e2f.
//
// Solidity: function isNullifierUnspent(uint256 nullifier) view returns(bool)
func (enygmaV1 *EnygmaV1) UnpackIsNullifierUnspent(data []byte) (bool, error) {
	out, err := enygmaV1.abi.Unpack("isNullifierUnspent", data)
	if err != nil {
		return *new(bool), err
	}
	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	return out0, err
}

// PackLastblockNum is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xa79d55e6.
//
// Solidity: function lastblockNum() view returns(uint256)
func (enygmaV1 *EnygmaV1) PackLastblockNum() []byte {
	enc, err := enygmaV1.abi.Pack("lastblockNum")
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackLastblockNum is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xa79d55e6.
//
// Solidity: function lastblockNum() view returns(uint256)
func (enygmaV1 *EnygmaV1) UnpackLastblockNum(data []byte) (*big.Int, error) {
	out, err := enygmaV1.abi.Unpack("lastblockNum", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, err
}

// PackLastblockNumAtCurrentBlockNumber is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x1da744e6.
//
// Solidity: function lastblockNumAtCurrentBlockNumber(uint256 ) view returns(uint256)
func (enygmaV1 *EnygmaV1) PackLastblockNumAtCurrentBlockNumber(arg0 *big.Int) []byte {
	enc, err := enygmaV1.abi.Pack("lastblockNumAtCurrentBlockNumber", arg0)
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackLastblockNumAtCurrentBlockNumber is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x1da744e6.
//
// Solidity: function lastblockNumAtCurrentBlockNumber(uint256 ) view returns(uint256)
func (enygmaV1 *EnygmaV1) UnpackLastblockNumAtCurrentBlockNumber(data []byte) (*big.Int, error) {
	out, err := enygmaV1.abi.Unpack("lastblockNumAtCurrentBlockNumber", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, err
}

// PackLastblockNumPending is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xacf92e78.
//
// Solidity: function lastblockNumPending() view returns(uint256)
func (enygmaV1 *EnygmaV1) PackLastblockNumPending() []byte {
	enc, err := enygmaV1.abi.Pack("lastblockNumPending")
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackLastblockNumPending is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xacf92e78.
//
// Solidity: function lastblockNumPending() view returns(uint256)
func (enygmaV1 *EnygmaV1) UnpackLastblockNumPending(data []byte) (*big.Int, error) {
	out, err := enygmaV1.abi.Unpack("lastblockNumPending", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, err
}

// PackNegateOnCurve is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x5e53dcb8.
//
// Solidity: function negateOnCurve(uint256 x) pure returns(uint256)
func (enygmaV1 *EnygmaV1) PackNegateOnCurve(x *big.Int) []byte {
	enc, err := enygmaV1.abi.Pack("negateOnCurve", x)
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackNegateOnCurve is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x5e53dcb8.
//
// Solidity: function negateOnCurve(uint256 x) pure returns(uint256)
func (enygmaV1 *EnygmaV1) UnpackNegateOnCurve(data []byte) (*big.Int, error) {
	out, err := enygmaV1.abi.Unpack("negateOnCurve", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, err
}

// PackNextBlockNumber is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x4c8b2126.
//
// Solidity: function nextBlockNumber(uint256 ) view returns(uint256)
func (enygmaV1 *EnygmaV1) PackNextBlockNumber(arg0 *big.Int) []byte {
	enc, err := enygmaV1.abi.Pack("nextBlockNumber", arg0)
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackNextBlockNumber is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x4c8b2126.
//
// Solidity: function nextBlockNumber(uint256 ) view returns(uint256)
func (enygmaV1 *EnygmaV1) UnpackNextBlockNumber(data []byte) (*big.Int, error) {
	out, err := enygmaV1.abi.Unpack("nextBlockNumber", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, err
}

// PackOwnerChainId is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xff4d1815.
//
// Solidity: function ownerChainId() view returns(uint256)
func (enygmaV1 *EnygmaV1) PackOwnerChainId() []byte {
	enc, err := enygmaV1.abi.Pack("ownerChainId")
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackOwnerChainId is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xff4d1815.
//
// Solidity: function ownerChainId() view returns(uint256)
func (enygmaV1 *EnygmaV1) UnpackOwnerChainId(data []byte) (*big.Int, error) {
	out, err := enygmaV1.abi.Unpack("ownerChainId", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, err
}

// PackParticipantStorageContract is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xfebce3a0.
//
// Solidity: function participantStorageContract() view returns(address)
func (enygmaV1 *EnygmaV1) PackParticipantStorageContract() []byte {
	enc, err := enygmaV1.abi.Pack("participantStorageContract")
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackParticipantStorageContract is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xfebce3a0.
//
// Solidity: function participantStorageContract() view returns(address)
func (enygmaV1 *EnygmaV1) UnpackParticipantStorageContract(data []byte) (common.Address, error) {
	out, err := enygmaV1.abi.Unpack("participantStorageContract", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, err
}

// PackPedCom is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x7d894a16.
//
// Solidity: function pedCom(uint256 v, uint256 r) view returns(uint256, uint256)
func (enygmaV1 *EnygmaV1) PackPedCom(v *big.Int, r *big.Int) []byte {
	enc, err := enygmaV1.abi.Pack("pedCom", v, r)
	if err != nil {
		panic(err)
	}
	return enc
}

// PedComOutput serves as a container for the return parameters of contract
// method PedCom.
type PedComOutput struct {
	Arg0 *big.Int
	Arg1 *big.Int
}

// UnpackPedCom is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x7d894a16.
//
// Solidity: function pedCom(uint256 v, uint256 r) view returns(uint256, uint256)
func (enygmaV1 *EnygmaV1) UnpackPedCom(data []byte) (PedComOutput, error) {
	out, err := enygmaV1.abi.Unpack("pedCom", data)
	outstruct := new(PedComOutput)
	if err != nil {
		return *outstruct, err
	}
	outstruct.Arg0 = abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	outstruct.Arg1 = abi.ConvertType(out[1], new(big.Int)).(*big.Int)
	return *outstruct, err

}

// PackPendingBalancesTallied is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xecccfd8e.
//
// Solidity: function pendingBalancesTallied(uint256 ) view returns(bool)
func (enygmaV1 *EnygmaV1) PackPendingBalancesTallied(arg0 *big.Int) []byte {
	enc, err := enygmaV1.abi.Pack("pendingBalancesTallied", arg0)
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackPendingBalancesTallied is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xecccfd8e.
//
// Solidity: function pendingBalancesTallied(uint256 ) view returns(bool)
func (enygmaV1 *EnygmaV1) UnpackPendingBalancesTallied(data []byte) (bool, error) {
	out, err := enygmaV1.abi.Unpack("pendingBalancesTallied", data)
	if err != nil {
		return *new(bool), err
	}
	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	return out0, err
}

// PackPendingMintsAndBurns is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x7a856ee5.
//
// Solidity: function pendingMintsAndBurns(uint256 ) view returns((uint256,uint256,uint256) pointToAddToBalance, uint256 amount, uint256 blockNumber, uint8 transactionType)
func (enygmaV1 *EnygmaV1) PackPendingMintsAndBurns(arg0 *big.Int) []byte {
	enc, err := enygmaV1.abi.Pack("pendingMintsAndBurns", arg0)
	if err != nil {
		panic(err)
	}
	return enc
}

// PendingMintsAndBurnsOutput serves as a container for the return parameters of contract
// method PendingMintsAndBurns.
type PendingMintsAndBurnsOutput struct {
	PointToAddToBalance IEnygmaV1EnygmaPointWithChainId
	Amount              *big.Int
	BlockNumber         *big.Int
	TransactionType     uint8
}

// UnpackPendingMintsAndBurns is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x7a856ee5.
//
// Solidity: function pendingMintsAndBurns(uint256 ) view returns((uint256,uint256,uint256) pointToAddToBalance, uint256 amount, uint256 blockNumber, uint8 transactionType)
func (enygmaV1 *EnygmaV1) UnpackPendingMintsAndBurns(data []byte) (PendingMintsAndBurnsOutput, error) {
	out, err := enygmaV1.abi.Unpack("pendingMintsAndBurns", data)
	outstruct := new(PendingMintsAndBurnsOutput)
	if err != nil {
		return *outstruct, err
	}
	outstruct.PointToAddToBalance = *abi.ConvertType(out[0], new(IEnygmaV1EnygmaPointWithChainId)).(*IEnygmaV1EnygmaPointWithChainId)
	outstruct.Amount = abi.ConvertType(out[1], new(big.Int)).(*big.Int)
	outstruct.BlockNumber = abi.ConvertType(out[2], new(big.Int)).(*big.Int)
	outstruct.TransactionType = *abi.ConvertType(out[3], new(uint8)).(*uint8)
	return *outstruct, err

}

// PackPendingTransactions is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x63a8374d.
//
// Solidity: function pendingTransactions(uint256 ) view returns(uint256 nullifier, uint8 transactionType)
func (enygmaV1 *EnygmaV1) PackPendingTransactions(arg0 *big.Int) []byte {
	enc, err := enygmaV1.abi.Pack("pendingTransactions", arg0)
	if err != nil {
		panic(err)
	}
	return enc
}

// PendingTransactionsOutput serves as a container for the return parameters of contract
// method PendingTransactions.
type PendingTransactionsOutput struct {
	Nullifier       *big.Int
	TransactionType uint8
}

// UnpackPendingTransactions is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x63a8374d.
//
// Solidity: function pendingTransactions(uint256 ) view returns(uint256 nullifier, uint8 transactionType)
func (enygmaV1 *EnygmaV1) UnpackPendingTransactions(data []byte) (PendingTransactionsOutput, error) {
	out, err := enygmaV1.abi.Unpack("pendingTransactions", data)
	outstruct := new(PendingTransactionsOutput)
	if err != nil {
		return *outstruct, err
	}
	outstruct.Nullifier = abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	outstruct.TransactionType = *abi.ConvertType(out[1], new(uint8)).(*uint8)
	return *outstruct, err

}

// PackRaylsNodeUserGovernance is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x0274c133.
//
// Solidity: function raylsNodeUserGovernance() view returns(address)
func (enygmaV1 *EnygmaV1) PackRaylsNodeUserGovernance() []byte {
	enc, err := enygmaV1.abi.Pack("raylsNodeUserGovernance")
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackRaylsNodeUserGovernance is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x0274c133.
//
// Solidity: function raylsNodeUserGovernance() view returns(address)
func (enygmaV1 *EnygmaV1) UnpackRaylsNodeUserGovernance(data []byte) (common.Address, error) {
	out, err := enygmaV1.abi.Unpack("raylsNodeUserGovernance", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, err
}

// PackReferenceBalance is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xa9bace48.
//
// Solidity: function referenceBalance(uint256 , uint256 ) view returns(uint256 c1, uint256 c2, uint256 chainId)
func (enygmaV1 *EnygmaV1) PackReferenceBalance(arg0 *big.Int, arg1 *big.Int) []byte {
	enc, err := enygmaV1.abi.Pack("referenceBalance", arg0, arg1)
	if err != nil {
		panic(err)
	}
	return enc
}

// ReferenceBalanceOutput serves as a container for the return parameters of contract
// method ReferenceBalance.
type ReferenceBalanceOutput struct {
	C1      *big.Int
	C2      *big.Int
	ChainId *big.Int
}

// UnpackReferenceBalance is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xa9bace48.
//
// Solidity: function referenceBalance(uint256 , uint256 ) view returns(uint256 c1, uint256 c2, uint256 chainId)
func (enygmaV1 *EnygmaV1) UnpackReferenceBalance(data []byte) (ReferenceBalanceOutput, error) {
	out, err := enygmaV1.abi.Unpack("referenceBalance", data)
	outstruct := new(ReferenceBalanceOutput)
	if err != nil {
		return *outstruct, err
	}
	outstruct.C1 = abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	outstruct.C2 = abi.ConvertType(out[1], new(big.Int)).(*big.Int)
	outstruct.ChainId = abi.ConvertType(out[2], new(big.Int)).(*big.Int)
	return *outstruct, err

}

// PackResourceId is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x5f997c5b.
//
// Solidity: function resourceId() view returns(bytes32)
func (enygmaV1 *EnygmaV1) PackResourceId() []byte {
	enc, err := enygmaV1.abi.Pack("resourceId")
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackResourceId is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x5f997c5b.
//
// Solidity: function resourceId() view returns(bytes32)
func (enygmaV1 *EnygmaV1) UnpackResourceId(data []byte) ([32]byte, error) {
	out, err := enygmaV1.abi.Unpack("resourceId", data)
	if err != nil {
		return *new([32]byte), err
	}
	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)
	return out0, err
}

// PackSetDvpIntegrationContract is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x9a41f721.
//
// Solidity: function setDvpIntegrationContract(address _dvpIntegrationContractAddress) returns()
func (enygmaV1 *EnygmaV1) PackSetDvpIntegrationContract(dvpIntegrationContractAddress common.Address) []byte {
	enc, err := enygmaV1.abi.Pack("setDvpIntegrationContract", dvpIntegrationContractAddress)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackSetResourceId is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xa01afbfb.
//
// Solidity: function setResourceId(bytes32 _resourceId) returns()
func (enygmaV1 *EnygmaV1) PackSetResourceId(resourceId [32]byte) []byte {
	enc, err := enygmaV1.abi.Pack("setResourceId", resourceId)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackTokenRegistryContract is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x08d63e03.
//
// Solidity: function tokenRegistryContract() view returns(address)
func (enygmaV1 *EnygmaV1) PackTokenRegistryContract() []byte {
	enc, err := enygmaV1.abi.Pack("tokenRegistryContract")
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackTokenRegistryContract is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x08d63e03.
//
// Solidity: function tokenRegistryContract() view returns(address)
func (enygmaV1 *EnygmaV1) UnpackTokenRegistryContract(data []byte) (common.Address, error) {
	out, err := enygmaV1.abi.Unpack("tokenRegistryContract", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, err
}

// PackTotalSupply is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (enygmaV1 *EnygmaV1) PackTotalSupply() []byte {
	enc, err := enygmaV1.abi.Pack("totalSupply")
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackTotalSupply is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (enygmaV1 *EnygmaV1) UnpackTotalSupply(data []byte) (*big.Int, error) {
	out, err := enygmaV1.abi.Unpack("totalSupply", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, err
}

// PackTotalSupplyX is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x71929e2a.
//
// Solidity: function totalSupplyX() view returns(uint256)
func (enygmaV1 *EnygmaV1) PackTotalSupplyX() []byte {
	enc, err := enygmaV1.abi.Pack("totalSupplyX")
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackTotalSupplyX is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x71929e2a.
//
// Solidity: function totalSupplyX() view returns(uint256)
func (enygmaV1 *EnygmaV1) UnpackTotalSupplyX(data []byte) (*big.Int, error) {
	out, err := enygmaV1.abi.Unpack("totalSupplyX", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, err
}

// PackTotalSupplyY is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x67511a4d.
//
// Solidity: function totalSupplyY() view returns(uint256)
func (enygmaV1 *EnygmaV1) PackTotalSupplyY() []byte {
	enc, err := enygmaV1.abi.Pack("totalSupplyY")
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackTotalSupplyY is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x67511a4d.
//
// Solidity: function totalSupplyY() view returns(uint256)
func (enygmaV1 *EnygmaV1) UnpackTotalSupplyY(data []byte) (*big.Int, error) {
	out, err := enygmaV1.abi.Unpack("totalSupplyY", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, err
}

// PackTransferBatch is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x8bcafc40.
//
// Solidity: function transferBatch((uint256[2],uint256[2][2],uint256[2],uint256[]) proof, bytes[] encryptedMessages) returns(bool)
func (enygmaV1 *EnygmaV1) PackTransferBatch(proof IEnygmaV1TransferProof, encryptedMessages [][]byte) []byte {
	enc, err := enygmaV1.abi.Pack("transferBatch", proof, encryptedMessages)
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackTransferBatch is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x8bcafc40.
//
// Solidity: function transferBatch((uint256[2],uint256[2][2],uint256[2],uint256[]) proof, bytes[] encryptedMessages) returns(bool)
func (enygmaV1 *EnygmaV1) UnpackTransferBatch(data []byte) (bool, error) {
	out, err := enygmaV1.abi.Unpack("transferBatch", data)
	if err != nil {
		return *new(bool), err
	}
	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	return out0, err
}

// PackTransferVerifiers is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xa0a0683e.
//
// Solidity: function transferVerifiers(uint256 ) view returns(address)
func (enygmaV1 *EnygmaV1) PackTransferVerifiers(arg0 *big.Int) []byte {
	enc, err := enygmaV1.abi.Pack("transferVerifiers", arg0)
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackTransferVerifiers is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xa0a0683e.
//
// Solidity: function transferVerifiers(uint256 ) view returns(address)
func (enygmaV1 *EnygmaV1) UnpackTransferVerifiers(data []byte) (common.Address, error) {
	out, err := enygmaV1.abi.Unpack("transferVerifiers", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, err
}

// PackUpdateSupply is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xd51ba5a4.
//
// Solidity: function updateSupply(uint256 _chainId, uint256 _blockNumber, (uint256,uint8) _update) returns()
func (enygmaV1 *EnygmaV1) PackUpdateSupply(chainId *big.Int, blockNumber *big.Int, update IEnygmaV1SupplyUpdateTx) []byte {
	enc, err := enygmaV1.abi.Pack("updateSupply", chainId, blockNumber, update)
	if err != nil {
		panic(err)
	}
	return enc
}

// EnygmaV1BalancesFinalised represents a BalancesFinalised event raised by the EnygmaV1 contract.
type EnygmaV1BalancesFinalised struct {
	BlockNumber *big.Int
	Raw         *types.Log // Blockchain specific contextual infos
}

const EnygmaV1BalancesFinalisedEventName = "BalancesFinalised"

// ContractEventName returns the user-defined event name.
func (EnygmaV1BalancesFinalised) ContractEventName() string {
	return EnygmaV1BalancesFinalisedEventName
}

// UnpackBalancesFinalisedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event BalancesFinalised(uint256 indexed blockNumber)
func (enygmaV1 *EnygmaV1) UnpackBalancesFinalisedEvent(log *types.Log) (*EnygmaV1BalancesFinalised, error) {
	event := "BalancesFinalised"
	if log.Topics[0] != enygmaV1.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(EnygmaV1BalancesFinalised)
	if len(log.Data) > 0 {
		if err := enygmaV1.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range enygmaV1.abi.Events[event].Inputs {
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

// EnygmaV1BurnSuccessful represents a BurnSuccessful event raised by the EnygmaV1 contract.
type EnygmaV1BurnSuccessful struct {
	ChainId   *big.Int
	BurnValue *big.Int
	Raw       *types.Log // Blockchain specific contextual infos
}

const EnygmaV1BurnSuccessfulEventName = "BurnSuccessful"

// ContractEventName returns the user-defined event name.
func (EnygmaV1BurnSuccessful) ContractEventName() string {
	return EnygmaV1BurnSuccessfulEventName
}

// UnpackBurnSuccessfulEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event BurnSuccessful(uint256 indexed chainId, uint256 burnValue)
func (enygmaV1 *EnygmaV1) UnpackBurnSuccessfulEvent(log *types.Log) (*EnygmaV1BurnSuccessful, error) {
	event := "BurnSuccessful"
	if log.Topics[0] != enygmaV1.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(EnygmaV1BurnSuccessful)
	if len(log.Data) > 0 {
		if err := enygmaV1.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range enygmaV1.abi.Events[event].Inputs {
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

// EnygmaV1NullifierConsumed represents a NullifierConsumed event raised by the EnygmaV1 contract.
type EnygmaV1NullifierConsumed struct {
	ResourceId  [32]byte
	Nullifier   *big.Int
	BlockNumber *big.Int
	TxType      uint8
	Raw         *types.Log // Blockchain specific contextual infos
}

const EnygmaV1NullifierConsumedEventName = "NullifierConsumed"

// ContractEventName returns the user-defined event name.
func (EnygmaV1NullifierConsumed) ContractEventName() string {
	return EnygmaV1NullifierConsumedEventName
}

// UnpackNullifierConsumedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event NullifierConsumed(bytes32 indexed resourceId, uint256 indexed nullifier, uint256 indexed blockNumber, uint8 txType)
func (enygmaV1 *EnygmaV1) UnpackNullifierConsumedEvent(log *types.Log) (*EnygmaV1NullifierConsumed, error) {
	event := "NullifierConsumed"
	if log.Topics[0] != enygmaV1.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(EnygmaV1NullifierConsumed)
	if len(log.Data) > 0 {
		if err := enygmaV1.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range enygmaV1.abi.Events[event].Inputs {
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

// EnygmaV1SupplyMinted represents a SupplyMinted event raised by the EnygmaV1 contract.
type EnygmaV1SupplyMinted struct {
	LastblockNum *big.Int
	Amount       *big.Int
	ToChainId    *big.Int
	Raw          *types.Log // Blockchain specific contextual infos
}

const EnygmaV1SupplyMintedEventName = "SupplyMinted"

// ContractEventName returns the user-defined event name.
func (EnygmaV1SupplyMinted) ContractEventName() string {
	return EnygmaV1SupplyMintedEventName
}

// UnpackSupplyMintedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event SupplyMinted(uint256 indexed lastblockNum, uint256 amount, uint256 toChainId)
func (enygmaV1 *EnygmaV1) UnpackSupplyMintedEvent(log *types.Log) (*EnygmaV1SupplyMinted, error) {
	event := "SupplyMinted"
	if log.Topics[0] != enygmaV1.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(EnygmaV1SupplyMinted)
	if len(log.Data) > 0 {
		if err := enygmaV1.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range enygmaV1.abi.Events[event].Inputs {
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

// EnygmaV1TokenRegistrationSubmitted represents a TokenRegistrationSubmitted event raised by the EnygmaV1 contract.
type EnygmaV1TokenRegistrationSubmitted struct {
	TokenAddress common.Address
	Raw          *types.Log // Blockchain specific contextual infos
}

const EnygmaV1TokenRegistrationSubmittedEventName = "TokenRegistrationSubmitted"

// ContractEventName returns the user-defined event name.
func (EnygmaV1TokenRegistrationSubmitted) ContractEventName() string {
	return EnygmaV1TokenRegistrationSubmittedEventName
}

// UnpackTokenRegistrationSubmittedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event TokenRegistrationSubmitted(address indexed tokenAddress)
func (enygmaV1 *EnygmaV1) UnpackTokenRegistrationSubmittedEvent(log *types.Log) (*EnygmaV1TokenRegistrationSubmitted, error) {
	event := "TokenRegistrationSubmitted"
	if log.Topics[0] != enygmaV1.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(EnygmaV1TokenRegistrationSubmitted)
	if len(log.Data) > 0 {
		if err := enygmaV1.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range enygmaV1.abi.Events[event].Inputs {
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

// EnygmaV1TransactionSuccessful represents a TransactionSuccessful event raised by the EnygmaV1 contract.
type EnygmaV1TransactionSuccessful struct {
	SenderAddress common.Address
	Raw           *types.Log // Blockchain specific contextual infos
}

const EnygmaV1TransactionSuccessfulEventName = "TransactionSuccessful"

// ContractEventName returns the user-defined event name.
func (EnygmaV1TransactionSuccessful) ContractEventName() string {
	return EnygmaV1TransactionSuccessfulEventName
}

// UnpackTransactionSuccessfulEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event TransactionSuccessful(address indexed senderAddress)
func (enygmaV1 *EnygmaV1) UnpackTransactionSuccessfulEvent(log *types.Log) (*EnygmaV1TransactionSuccessful, error) {
	event := "TransactionSuccessful"
	if log.Topics[0] != enygmaV1.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(EnygmaV1TransactionSuccessful)
	if len(log.Data) > 0 {
		if err := enygmaV1.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range enygmaV1.abi.Events[event].Inputs {
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

// EnygmaV1VerifierRegistered represents a VerifierRegistered event raised by the EnygmaV1 contract.
type EnygmaV1VerifierRegistered struct {
	VerifierAddress common.Address
	K               uint8
	Raw             *types.Log // Blockchain specific contextual infos
}

const EnygmaV1VerifierRegisteredEventName = "VerifierRegistered"

// ContractEventName returns the user-defined event name.
func (EnygmaV1VerifierRegistered) ContractEventName() string {
	return EnygmaV1VerifierRegisteredEventName
}

// UnpackVerifierRegisteredEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event VerifierRegistered(address indexed verifierAddress, uint8 k)
func (enygmaV1 *EnygmaV1) UnpackVerifierRegisteredEvent(log *types.Log) (*EnygmaV1VerifierRegistered, error) {
	event := "VerifierRegistered"
	if log.Topics[0] != enygmaV1.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(EnygmaV1VerifierRegistered)
	if len(log.Data) > 0 {
		if err := enygmaV1.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range enygmaV1.abi.Events[event].Inputs {
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
func (enygmaV1 *EnygmaV1) UnpackError(raw []byte) (any, error) {
	if bytes.Equal(raw[:4], enygmaV1.abi.Errors["EnygmaV1OnlyDvpIntegrationAllowed"].ID.Bytes()[:4]) {
		return enygmaV1.UnpackEnygmaV1OnlyDvpIntegrationAllowedError(raw[4:])
	}
	if bytes.Equal(raw[:4], enygmaV1.abi.Errors["EnygmaV1OnlyFactoryAllowed"].ID.Bytes()[:4]) {
		return enygmaV1.UnpackEnygmaV1OnlyFactoryAllowedError(raw[4:])
	}
	if bytes.Equal(raw[:4], enygmaV1.abi.Errors["RaylsAppHubNotActive"].ID.Bytes()[:4]) {
		return enygmaV1.UnpackRaylsAppHubNotActiveError(raw[4:])
	}
	if bytes.Equal(raw[:4], enygmaV1.abi.Errors["RaylsAppPrivacyNodeFrozen"].ID.Bytes()[:4]) {
		return enygmaV1.UnpackRaylsAppPrivacyNodeFrozenError(raw[4:])
	}
	if bytes.Equal(raw[:4], enygmaV1.abi.Errors["RaylsAppPrivacyNodeNotActive"].ID.Bytes()[:4]) {
		return enygmaV1.UnpackRaylsAppPrivacyNodeNotActiveError(raw[4:])
	}
	if bytes.Equal(raw[:4], enygmaV1.abi.Errors["RaylsAppPublicChainNotActive"].ID.Bytes()[:4]) {
		return enygmaV1.UnpackRaylsAppPublicChainNotActiveError(raw[4:])
	}
	if bytes.Equal(raw[:4], enygmaV1.abi.Errors["RaylsAppResourceNotApproved"].ID.Bytes()[:4]) {
		return enygmaV1.UnpackRaylsAppResourceNotApprovedError(raw[4:])
	}
	if bytes.Equal(raw[:4], enygmaV1.abi.Errors["RaylsAppTokenRegistryNotConfigured"].ID.Bytes()[:4]) {
		return enygmaV1.UnpackRaylsAppTokenRegistryNotConfiguredError(raw[4:])
	}
	if bytes.Equal(raw[:4], enygmaV1.abi.Errors["RaylsAppUnauthorizedTokenRegistry"].ID.Bytes()[:4]) {
		return enygmaV1.UnpackRaylsAppUnauthorizedTokenRegistryError(raw[4:])
	}
	if bytes.Equal(raw[:4], enygmaV1.abi.Errors["RaylsAppUserNotRegistered"].ID.Bytes()[:4]) {
		return enygmaV1.UnpackRaylsAppUserNotRegisteredError(raw[4:])
	}
	return nil, errors.New("Unknown error")
}

// EnygmaV1EnygmaV1OnlyDvpIntegrationAllowed represents a EnygmaV1__OnlyDvpIntegrationAllowed error raised by the EnygmaV1 contract.
type EnygmaV1EnygmaV1OnlyDvpIntegrationAllowed struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error EnygmaV1__OnlyDvpIntegrationAllowed()
func EnygmaV1EnygmaV1OnlyDvpIntegrationAllowedErrorID() common.Hash {
	return common.HexToHash("0x0a5ee5b8a62f48e8bb7d88c4f9fee802029e604b5305c91f80d4fe29d7074537")
}

// UnpackEnygmaV1OnlyDvpIntegrationAllowedError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error EnygmaV1__OnlyDvpIntegrationAllowed()
func (enygmaV1 *EnygmaV1) UnpackEnygmaV1OnlyDvpIntegrationAllowedError(raw []byte) (*EnygmaV1EnygmaV1OnlyDvpIntegrationAllowed, error) {
	out := new(EnygmaV1EnygmaV1OnlyDvpIntegrationAllowed)
	if err := enygmaV1.abi.UnpackIntoInterface(out, "EnygmaV1OnlyDvpIntegrationAllowed", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// EnygmaV1EnygmaV1OnlyFactoryAllowed represents a EnygmaV1__OnlyFactoryAllowed error raised by the EnygmaV1 contract.
type EnygmaV1EnygmaV1OnlyFactoryAllowed struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error EnygmaV1__OnlyFactoryAllowed()
func EnygmaV1EnygmaV1OnlyFactoryAllowedErrorID() common.Hash {
	return common.HexToHash("0x22a3b9b3b620e895ff9eb2b0bf415819ca7c289b0be4e243c2591a216bf672b4")
}

// UnpackEnygmaV1OnlyFactoryAllowedError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error EnygmaV1__OnlyFactoryAllowed()
func (enygmaV1 *EnygmaV1) UnpackEnygmaV1OnlyFactoryAllowedError(raw []byte) (*EnygmaV1EnygmaV1OnlyFactoryAllowed, error) {
	out := new(EnygmaV1EnygmaV1OnlyFactoryAllowed)
	if err := enygmaV1.abi.UnpackIntoInterface(out, "EnygmaV1OnlyFactoryAllowed", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// EnygmaV1RaylsAppHubNotActive represents a RaylsApp__HubNotActive error raised by the EnygmaV1 contract.
type EnygmaV1RaylsAppHubNotActive struct {
	TokenAddress      common.Address
	PrivacyNodeStatus uint8
	HubStatus         uint8
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error RaylsApp__HubNotActive(address tokenAddress, uint8 privacyNodeStatus, uint8 hubStatus)
func EnygmaV1RaylsAppHubNotActiveErrorID() common.Hash {
	return common.HexToHash("0xdc2ffb0fada912f0dd1b700d4ea9a9ce47e3ecdd1b7b155d2066b9a022a637c2")
}

// UnpackRaylsAppHubNotActiveError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error RaylsApp__HubNotActive(address tokenAddress, uint8 privacyNodeStatus, uint8 hubStatus)
func (enygmaV1 *EnygmaV1) UnpackRaylsAppHubNotActiveError(raw []byte) (*EnygmaV1RaylsAppHubNotActive, error) {
	out := new(EnygmaV1RaylsAppHubNotActive)
	if err := enygmaV1.abi.UnpackIntoInterface(out, "RaylsAppHubNotActive", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// EnygmaV1RaylsAppPrivacyNodeFrozen represents a RaylsApp__PrivacyNodeFrozen error raised by the EnygmaV1 contract.
type EnygmaV1RaylsAppPrivacyNodeFrozen struct {
	TokenAddress common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error RaylsApp__PrivacyNodeFrozen(address tokenAddress)
func EnygmaV1RaylsAppPrivacyNodeFrozenErrorID() common.Hash {
	return common.HexToHash("0xcecb8d3ce0d1417038942c9d252e856b5585275082aa5cdbca675fa64d7bfc24")
}

// UnpackRaylsAppPrivacyNodeFrozenError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error RaylsApp__PrivacyNodeFrozen(address tokenAddress)
func (enygmaV1 *EnygmaV1) UnpackRaylsAppPrivacyNodeFrozenError(raw []byte) (*EnygmaV1RaylsAppPrivacyNodeFrozen, error) {
	out := new(EnygmaV1RaylsAppPrivacyNodeFrozen)
	if err := enygmaV1.abi.UnpackIntoInterface(out, "RaylsAppPrivacyNodeFrozen", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// EnygmaV1RaylsAppPrivacyNodeNotActive represents a RaylsApp__PrivacyNodeNotActive error raised by the EnygmaV1 contract.
type EnygmaV1RaylsAppPrivacyNodeNotActive struct {
	TokenAddress      common.Address
	PrivacyNodeStatus uint8
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error RaylsApp__PrivacyNodeNotActive(address tokenAddress, uint8 privacyNodeStatus)
func EnygmaV1RaylsAppPrivacyNodeNotActiveErrorID() common.Hash {
	return common.HexToHash("0x44c58c43ed8f726e3330349bec7aa7300f000be36837ee0c2cf507d04511e1e8")
}

// UnpackRaylsAppPrivacyNodeNotActiveError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error RaylsApp__PrivacyNodeNotActive(address tokenAddress, uint8 privacyNodeStatus)
func (enygmaV1 *EnygmaV1) UnpackRaylsAppPrivacyNodeNotActiveError(raw []byte) (*EnygmaV1RaylsAppPrivacyNodeNotActive, error) {
	out := new(EnygmaV1RaylsAppPrivacyNodeNotActive)
	if err := enygmaV1.abi.UnpackIntoInterface(out, "RaylsAppPrivacyNodeNotActive", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// EnygmaV1RaylsAppPublicChainNotActive represents a RaylsApp__PublicChainNotActive error raised by the EnygmaV1 contract.
type EnygmaV1RaylsAppPublicChainNotActive struct {
	TokenAddress      common.Address
	PrivacyNodeStatus uint8
	PublicChainStatus uint8
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error RaylsApp__PublicChainNotActive(address tokenAddress, uint8 privacyNodeStatus, uint8 publicChainStatus)
func EnygmaV1RaylsAppPublicChainNotActiveErrorID() common.Hash {
	return common.HexToHash("0xd6e23bd403a5000c9afe5c2ed5202b3ff8e25d8c3644c1f51892016fb18e5ab9")
}

// UnpackRaylsAppPublicChainNotActiveError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error RaylsApp__PublicChainNotActive(address tokenAddress, uint8 privacyNodeStatus, uint8 publicChainStatus)
func (enygmaV1 *EnygmaV1) UnpackRaylsAppPublicChainNotActiveError(raw []byte) (*EnygmaV1RaylsAppPublicChainNotActive, error) {
	out := new(EnygmaV1RaylsAppPublicChainNotActive)
	if err := enygmaV1.abi.UnpackIntoInterface(out, "RaylsAppPublicChainNotActive", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// EnygmaV1RaylsAppResourceNotApproved represents a RaylsApp__ResourceNotApproved error raised by the EnygmaV1 contract.
type EnygmaV1RaylsAppResourceNotApproved struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error RaylsApp__ResourceNotApproved()
func EnygmaV1RaylsAppResourceNotApprovedErrorID() common.Hash {
	return common.HexToHash("0x970ad4f73c2c200faa068d3d920e2ef40fca6a5338655abcfb5212557edeed6b")
}

// UnpackRaylsAppResourceNotApprovedError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error RaylsApp__ResourceNotApproved()
func (enygmaV1 *EnygmaV1) UnpackRaylsAppResourceNotApprovedError(raw []byte) (*EnygmaV1RaylsAppResourceNotApproved, error) {
	out := new(EnygmaV1RaylsAppResourceNotApproved)
	if err := enygmaV1.abi.UnpackIntoInterface(out, "RaylsAppResourceNotApproved", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// EnygmaV1RaylsAppTokenRegistryNotConfigured represents a RaylsApp__TokenRegistryNotConfigured error raised by the EnygmaV1 contract.
type EnygmaV1RaylsAppTokenRegistryNotConfigured struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error RaylsApp__TokenRegistryNotConfigured()
func EnygmaV1RaylsAppTokenRegistryNotConfiguredErrorID() common.Hash {
	return common.HexToHash("0x36a41bd1f6f11cd28b716e935a926fb04f66e11a393b38a49bb660640f3b6dbf")
}

// UnpackRaylsAppTokenRegistryNotConfiguredError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error RaylsApp__TokenRegistryNotConfigured()
func (enygmaV1 *EnygmaV1) UnpackRaylsAppTokenRegistryNotConfiguredError(raw []byte) (*EnygmaV1RaylsAppTokenRegistryNotConfigured, error) {
	out := new(EnygmaV1RaylsAppTokenRegistryNotConfigured)
	if err := enygmaV1.abi.UnpackIntoInterface(out, "RaylsAppTokenRegistryNotConfigured", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// EnygmaV1RaylsAppUnauthorizedTokenRegistry represents a RaylsApp__UnauthorizedTokenRegistry error raised by the EnygmaV1 contract.
type EnygmaV1RaylsAppUnauthorizedTokenRegistry struct {
	Caller        common.Address
	TokenRegistry common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error RaylsApp__UnauthorizedTokenRegistry(address caller, address tokenRegistry)
func EnygmaV1RaylsAppUnauthorizedTokenRegistryErrorID() common.Hash {
	return common.HexToHash("0x061526480acdfaa09331b795496a6c50aaed25a45d9fca4c9d55fad56af8e09c")
}

// UnpackRaylsAppUnauthorizedTokenRegistryError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error RaylsApp__UnauthorizedTokenRegistry(address caller, address tokenRegistry)
func (enygmaV1 *EnygmaV1) UnpackRaylsAppUnauthorizedTokenRegistryError(raw []byte) (*EnygmaV1RaylsAppUnauthorizedTokenRegistry, error) {
	out := new(EnygmaV1RaylsAppUnauthorizedTokenRegistry)
	if err := enygmaV1.abi.UnpackIntoInterface(out, "RaylsAppUnauthorizedTokenRegistry", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// EnygmaV1RaylsAppUserNotRegistered represents a RaylsApp__UserNotRegistered error raised by the EnygmaV1 contract.
type EnygmaV1RaylsAppUserNotRegistered struct {
	Caller common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error RaylsApp__UserNotRegistered(address caller)
func EnygmaV1RaylsAppUserNotRegisteredErrorID() common.Hash {
	return common.HexToHash("0x4c1db902cce08bec31bedc484362fba54949899ac3c0bf0416f3c44af3284baa")
}

// UnpackRaylsAppUserNotRegisteredError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error RaylsApp__UserNotRegistered(address caller)
func (enygmaV1 *EnygmaV1) UnpackRaylsAppUserNotRegisteredError(raw []byte) (*EnygmaV1RaylsAppUserNotRegistered, error) {
	out := new(EnygmaV1RaylsAppUserNotRegistered)
	if err := enygmaV1.abi.UnpackIntoInterface(out, "RaylsAppUserNotRegistered", raw); err != nil {
		return nil, err
	}
	return out, nil
}
