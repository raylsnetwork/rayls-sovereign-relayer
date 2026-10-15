// Code generated via abigen V2 - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package IEnygmaDvpIntegration

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

// IEnygmaDvpIntegrationWithdrawOrDepositProof is an auto generated low-level Go binding around an user-defined struct.
type IEnygmaDvpIntegrationWithdrawOrDepositProof struct {
	PiA          [2]*big.Int
	PiB          [2][2]*big.Int
	PiC          [2]*big.Int
	PublicSignal []*big.Int
}

// IEnygmaDvpIntegrationMetaData contains all meta data concerning the IEnygmaDvpIntegration contract.
var IEnygmaDvpIntegrationMetaData = bind.MetaData{
	ABI: "[{\"type\":\"function\",\"name\":\"addDepositToDvpVerifier\",\"inputs\":[{\"name\":\"depositToDvpVerifier\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"k\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"addDvp\",\"inputs\":[{\"name\":\"dvpAddress\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"addWithdrawFromDvpVerifier\",\"inputs\":[{\"name\":\"withdrawFromDvpVerifier\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"k\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"consolidateFunds\",\"inputs\":[{\"name\":\"transaction\",\"type\":\"tuple\",\"internalType\":\"structIDvp.ProofReceipt\",\"components\":[{\"name\":\"proof\",\"type\":\"tuple\",\"internalType\":\"structIDvp.SnarkProof\",\"components\":[{\"name\":\"a\",\"type\":\"tuple\",\"internalType\":\"structIDvp.G1Point\",\"components\":[{\"name\":\"x\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"b\",\"type\":\"tuple\",\"internalType\":\"structIDvp.G2Point\",\"components\":[{\"name\":\"x\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"},{\"name\":\"y\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"}]},{\"name\":\"c\",\"type\":\"tuple\",\"internalType\":\"structIDvp.G1Point\",\"components\":[{\"name\":\"x\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]}]},{\"name\":\"treeNumbers\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"},{\"name\":\"message\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"merkleRoots\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"},{\"name\":\"commitments\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"},{\"name\":\"nullifiers\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"},{\"name\":\"revertCommitment\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"depositToDvp\",\"inputs\":[{\"name\":\"proof\",\"type\":\"tuple\",\"internalType\":\"structIEnygmaDvpIntegration.WithdrawOrDepositProof\",\"components\":[{\"name\":\"pi_a\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"},{\"name\":\"pi_b\",\"type\":\"uint256[2][2]\",\"internalType\":\"uint256[2][2]\"},{\"name\":\"pi_c\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"},{\"name\":\"public_signal\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"}]},{\"name\":\"encryptedMessages\",\"type\":\"bytes[]\",\"internalType\":\"bytes[]\"},{\"name\":\"encryptedBurnUpdate\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"getDepositToDvpVerifierAddress\",\"inputs\":[{\"name\":\"k\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getDvpAddress\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getVaultId\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getWithdrawFromDvpVerifierAddress\",\"inputs\":[{\"name\":\"k\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"setVaultId\",\"inputs\":[{\"name\":\"vaultId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"withdrawFromDvp\",\"inputs\":[{\"name\":\"proof\",\"type\":\"tuple\",\"internalType\":\"structIEnygmaDvpIntegration.WithdrawOrDepositProof\",\"components\":[{\"name\":\"pi_a\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"},{\"name\":\"pi_b\",\"type\":\"uint256[2][2]\",\"internalType\":\"uint256[2][2]\"},{\"name\":\"pi_c\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"},{\"name\":\"public_signal\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"}]},{\"name\":\"encryptedMessages\",\"type\":\"bytes[]\",\"internalType\":\"bytes[]\"},{\"name\":\"transaction\",\"type\":\"tuple\",\"internalType\":\"structIDvp.ProofReceipt\",\"components\":[{\"name\":\"proof\",\"type\":\"tuple\",\"internalType\":\"structIDvp.SnarkProof\",\"components\":[{\"name\":\"a\",\"type\":\"tuple\",\"internalType\":\"structIDvp.G1Point\",\"components\":[{\"name\":\"x\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"b\",\"type\":\"tuple\",\"internalType\":\"structIDvp.G2Point\",\"components\":[{\"name\":\"x\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"},{\"name\":\"y\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"}]},{\"name\":\"c\",\"type\":\"tuple\",\"internalType\":\"structIDvp.G1Point\",\"components\":[{\"name\":\"x\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]}]},{\"name\":\"treeNumbers\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"},{\"name\":\"message\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"merkleRoots\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"},{\"name\":\"commitments\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"},{\"name\":\"nullifiers\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"},{\"name\":\"revertCommitment\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"encryptedMintUpdate\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"event\",\"name\":\"ConsolidateFundsSuccessful\",\"inputs\":[{\"name\":\"transaction\",\"type\":\"tuple\",\"indexed\":false,\"internalType\":\"structIDvp.ProofReceipt\",\"components\":[{\"name\":\"proof\",\"type\":\"tuple\",\"internalType\":\"structIDvp.SnarkProof\",\"components\":[{\"name\":\"a\",\"type\":\"tuple\",\"internalType\":\"structIDvp.G1Point\",\"components\":[{\"name\":\"x\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"b\",\"type\":\"tuple\",\"internalType\":\"structIDvp.G2Point\",\"components\":[{\"name\":\"x\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"},{\"name\":\"y\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"}]},{\"name\":\"c\",\"type\":\"tuple\",\"internalType\":\"structIDvp.G1Point\",\"components\":[{\"name\":\"x\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]}]},{\"name\":\"treeNumbers\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"},{\"name\":\"message\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"merkleRoots\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"},{\"name\":\"commitments\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"},{\"name\":\"nullifiers\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"},{\"name\":\"revertCommitment\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"enygmaContractAddress\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"DepositToDvpSuccessful\",\"inputs\":[{\"name\":\"commitment\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"enygmaContractAddress\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"DvpRegistered\",\"inputs\":[{\"name\":\"dvpAddress\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"VaultIdSet\",\"inputs\":[{\"name\":\"vaultId\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"VerifierDepositToDvpRegistered\",\"inputs\":[{\"name\":\"verifierAddress\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"k\",\"type\":\"uint8\",\"indexed\":false,\"internalType\":\"uint8\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"VerifierWithdrawFromDvpRegistered\",\"inputs\":[{\"name\":\"verifierAddress\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"k\",\"type\":\"uint8\",\"indexed\":false,\"internalType\":\"uint8\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"WithdrawFromDvpSuccessful\",\"inputs\":[{\"name\":\"transaction\",\"type\":\"tuple\",\"indexed\":false,\"internalType\":\"structIDvp.ProofReceipt\",\"components\":[{\"name\":\"proof\",\"type\":\"tuple\",\"internalType\":\"structIDvp.SnarkProof\",\"components\":[{\"name\":\"a\",\"type\":\"tuple\",\"internalType\":\"structIDvp.G1Point\",\"components\":[{\"name\":\"x\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"b\",\"type\":\"tuple\",\"internalType\":\"structIDvp.G2Point\",\"components\":[{\"name\":\"x\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"},{\"name\":\"y\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"}]},{\"name\":\"c\",\"type\":\"tuple\",\"internalType\":\"structIDvp.G1Point\",\"components\":[{\"name\":\"x\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]}]},{\"name\":\"treeNumbers\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"},{\"name\":\"message\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"merkleRoots\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"},{\"name\":\"commitments\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"},{\"name\":\"nullifiers\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"},{\"name\":\"revertCommitment\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"enygmaContractAddress\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false}]",
	ID:  "IEnygmaDvpIntegration",
}

// IEnygmaDvpIntegration is an auto generated Go binding around an Ethereum contract.
type IEnygmaDvpIntegration struct {
	abi abi.ABI
}

// NewIEnygmaDvpIntegration creates a new instance of IEnygmaDvpIntegration.
func NewIEnygmaDvpIntegration() *IEnygmaDvpIntegration {
	parsed, err := IEnygmaDvpIntegrationMetaData.ParseABI()
	if err != nil {
		panic(errors.New("invalid ABI: " + err.Error()))
	}
	return &IEnygmaDvpIntegration{abi: *parsed}
}

// Instance creates a wrapper for a deployed contract instance at the given address.
// Use this to create the instance object passed to abigen v2 library functions Call, Transact, etc.
func (c *IEnygmaDvpIntegration) Instance(backend bind.ContractBackend, addr common.Address) *bind.BoundContract {
	return bind.NewBoundContract(addr, c.abi, backend, backend, backend)
}

// PackAddDepositToDvpVerifier is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x104c9774.
//
// Solidity: function addDepositToDvpVerifier(address depositToDvpVerifier, uint8 k) returns(bool)
func (iEnygmaDvpIntegration *IEnygmaDvpIntegration) PackAddDepositToDvpVerifier(depositToDvpVerifier common.Address, k uint8) []byte {
	enc, err := iEnygmaDvpIntegration.abi.Pack("addDepositToDvpVerifier", depositToDvpVerifier, k)
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackAddDepositToDvpVerifier is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x104c9774.
//
// Solidity: function addDepositToDvpVerifier(address depositToDvpVerifier, uint8 k) returns(bool)
func (iEnygmaDvpIntegration *IEnygmaDvpIntegration) UnpackAddDepositToDvpVerifier(data []byte) (bool, error) {
	out, err := iEnygmaDvpIntegration.abi.Unpack("addDepositToDvpVerifier", data)
	if err != nil {
		return *new(bool), err
	}
	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	return out0, err
}

// PackAddDvp is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x860d3ece.
//
// Solidity: function addDvp(address dvpAddress) returns(bool)
func (iEnygmaDvpIntegration *IEnygmaDvpIntegration) PackAddDvp(dvpAddress common.Address) []byte {
	enc, err := iEnygmaDvpIntegration.abi.Pack("addDvp", dvpAddress)
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackAddDvp is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x860d3ece.
//
// Solidity: function addDvp(address dvpAddress) returns(bool)
func (iEnygmaDvpIntegration *IEnygmaDvpIntegration) UnpackAddDvp(data []byte) (bool, error) {
	out, err := iEnygmaDvpIntegration.abi.Unpack("addDvp", data)
	if err != nil {
		return *new(bool), err
	}
	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	return out0, err
}

// PackAddWithdrawFromDvpVerifier is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x5cf983db.
//
// Solidity: function addWithdrawFromDvpVerifier(address withdrawFromDvpVerifier, uint8 k) returns(bool)
func (iEnygmaDvpIntegration *IEnygmaDvpIntegration) PackAddWithdrawFromDvpVerifier(withdrawFromDvpVerifier common.Address, k uint8) []byte {
	enc, err := iEnygmaDvpIntegration.abi.Pack("addWithdrawFromDvpVerifier", withdrawFromDvpVerifier, k)
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackAddWithdrawFromDvpVerifier is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x5cf983db.
//
// Solidity: function addWithdrawFromDvpVerifier(address withdrawFromDvpVerifier, uint8 k) returns(bool)
func (iEnygmaDvpIntegration *IEnygmaDvpIntegration) UnpackAddWithdrawFromDvpVerifier(data []byte) (bool, error) {
	out, err := iEnygmaDvpIntegration.abi.Unpack("addWithdrawFromDvpVerifier", data)
	if err != nil {
		return *new(bool), err
	}
	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	return out0, err
}

// PackConsolidateFunds is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x8d54447b.
//
// Solidity: function consolidateFunds((((uint256,uint256),(uint256[2],uint256[2]),(uint256,uint256)),uint256[],uint256,uint256[],uint256[],uint256[],uint256) transaction) returns(bool)
func (iEnygmaDvpIntegration *IEnygmaDvpIntegration) PackConsolidateFunds(transaction IDvpProofReceipt) []byte {
	enc, err := iEnygmaDvpIntegration.abi.Pack("consolidateFunds", transaction)
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackConsolidateFunds is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x8d54447b.
//
// Solidity: function consolidateFunds((((uint256,uint256),(uint256[2],uint256[2]),(uint256,uint256)),uint256[],uint256,uint256[],uint256[],uint256[],uint256) transaction) returns(bool)
func (iEnygmaDvpIntegration *IEnygmaDvpIntegration) UnpackConsolidateFunds(data []byte) (bool, error) {
	out, err := iEnygmaDvpIntegration.abi.Unpack("consolidateFunds", data)
	if err != nil {
		return *new(bool), err
	}
	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	return out0, err
}

// PackDepositToDvp is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xecc920c5.
//
// Solidity: function depositToDvp((uint256[2],uint256[2][2],uint256[2],uint256[]) proof, bytes[] encryptedMessages, bytes encryptedBurnUpdate) returns(bool)
func (iEnygmaDvpIntegration *IEnygmaDvpIntegration) PackDepositToDvp(proof IEnygmaDvpIntegrationWithdrawOrDepositProof, encryptedMessages [][]byte, encryptedBurnUpdate []byte) []byte {
	enc, err := iEnygmaDvpIntegration.abi.Pack("depositToDvp", proof, encryptedMessages, encryptedBurnUpdate)
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackDepositToDvp is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xecc920c5.
//
// Solidity: function depositToDvp((uint256[2],uint256[2][2],uint256[2],uint256[]) proof, bytes[] encryptedMessages, bytes encryptedBurnUpdate) returns(bool)
func (iEnygmaDvpIntegration *IEnygmaDvpIntegration) UnpackDepositToDvp(data []byte) (bool, error) {
	out, err := iEnygmaDvpIntegration.abi.Unpack("depositToDvp", data)
	if err != nil {
		return *new(bool), err
	}
	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	return out0, err
}

// PackGetDepositToDvpVerifierAddress is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xffa26f38.
//
// Solidity: function getDepositToDvpVerifierAddress(uint8 k) view returns(address)
func (iEnygmaDvpIntegration *IEnygmaDvpIntegration) PackGetDepositToDvpVerifierAddress(k uint8) []byte {
	enc, err := iEnygmaDvpIntegration.abi.Pack("getDepositToDvpVerifierAddress", k)
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackGetDepositToDvpVerifierAddress is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xffa26f38.
//
// Solidity: function getDepositToDvpVerifierAddress(uint8 k) view returns(address)
func (iEnygmaDvpIntegration *IEnygmaDvpIntegration) UnpackGetDepositToDvpVerifierAddress(data []byte) (common.Address, error) {
	out, err := iEnygmaDvpIntegration.abi.Unpack("getDepositToDvpVerifierAddress", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, err
}

// PackGetDvpAddress is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x7000a7c9.
//
// Solidity: function getDvpAddress() view returns(address)
func (iEnygmaDvpIntegration *IEnygmaDvpIntegration) PackGetDvpAddress() []byte {
	enc, err := iEnygmaDvpIntegration.abi.Pack("getDvpAddress")
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackGetDvpAddress is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x7000a7c9.
//
// Solidity: function getDvpAddress() view returns(address)
func (iEnygmaDvpIntegration *IEnygmaDvpIntegration) UnpackGetDvpAddress(data []byte) (common.Address, error) {
	out, err := iEnygmaDvpIntegration.abi.Unpack("getDvpAddress", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, err
}

// PackGetVaultId is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x4e2285f1.
//
// Solidity: function getVaultId() view returns(uint256)
func (iEnygmaDvpIntegration *IEnygmaDvpIntegration) PackGetVaultId() []byte {
	enc, err := iEnygmaDvpIntegration.abi.Pack("getVaultId")
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackGetVaultId is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x4e2285f1.
//
// Solidity: function getVaultId() view returns(uint256)
func (iEnygmaDvpIntegration *IEnygmaDvpIntegration) UnpackGetVaultId(data []byte) (*big.Int, error) {
	out, err := iEnygmaDvpIntegration.abi.Unpack("getVaultId", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, err
}

// PackGetWithdrawFromDvpVerifierAddress is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x9ffba8d8.
//
// Solidity: function getWithdrawFromDvpVerifierAddress(uint8 k) view returns(address)
func (iEnygmaDvpIntegration *IEnygmaDvpIntegration) PackGetWithdrawFromDvpVerifierAddress(k uint8) []byte {
	enc, err := iEnygmaDvpIntegration.abi.Pack("getWithdrawFromDvpVerifierAddress", k)
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackGetWithdrawFromDvpVerifierAddress is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x9ffba8d8.
//
// Solidity: function getWithdrawFromDvpVerifierAddress(uint8 k) view returns(address)
func (iEnygmaDvpIntegration *IEnygmaDvpIntegration) UnpackGetWithdrawFromDvpVerifierAddress(data []byte) (common.Address, error) {
	out, err := iEnygmaDvpIntegration.abi.Unpack("getWithdrawFromDvpVerifierAddress", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, err
}

// PackSetVaultId is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xb18c312d.
//
// Solidity: function setVaultId(uint256 vaultId) returns(bool)
func (iEnygmaDvpIntegration *IEnygmaDvpIntegration) PackSetVaultId(vaultId *big.Int) []byte {
	enc, err := iEnygmaDvpIntegration.abi.Pack("setVaultId", vaultId)
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackSetVaultId is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xb18c312d.
//
// Solidity: function setVaultId(uint256 vaultId) returns(bool)
func (iEnygmaDvpIntegration *IEnygmaDvpIntegration) UnpackSetVaultId(data []byte) (bool, error) {
	out, err := iEnygmaDvpIntegration.abi.Unpack("setVaultId", data)
	if err != nil {
		return *new(bool), err
	}
	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	return out0, err
}

// PackWithdrawFromDvp is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xbb295d70.
//
// Solidity: function withdrawFromDvp((uint256[2],uint256[2][2],uint256[2],uint256[]) proof, bytes[] encryptedMessages, (((uint256,uint256),(uint256[2],uint256[2]),(uint256,uint256)),uint256[],uint256,uint256[],uint256[],uint256[],uint256) transaction, bytes encryptedMintUpdate) returns(bool)
func (iEnygmaDvpIntegration *IEnygmaDvpIntegration) PackWithdrawFromDvp(proof IEnygmaDvpIntegrationWithdrawOrDepositProof, encryptedMessages [][]byte, transaction IDvpProofReceipt, encryptedMintUpdate []byte) []byte {
	enc, err := iEnygmaDvpIntegration.abi.Pack("withdrawFromDvp", proof, encryptedMessages, transaction, encryptedMintUpdate)
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackWithdrawFromDvp is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xbb295d70.
//
// Solidity: function withdrawFromDvp((uint256[2],uint256[2][2],uint256[2],uint256[]) proof, bytes[] encryptedMessages, (((uint256,uint256),(uint256[2],uint256[2]),(uint256,uint256)),uint256[],uint256,uint256[],uint256[],uint256[],uint256) transaction, bytes encryptedMintUpdate) returns(bool)
func (iEnygmaDvpIntegration *IEnygmaDvpIntegration) UnpackWithdrawFromDvp(data []byte) (bool, error) {
	out, err := iEnygmaDvpIntegration.abi.Unpack("withdrawFromDvp", data)
	if err != nil {
		return *new(bool), err
	}
	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	return out0, err
}

// IEnygmaDvpIntegrationConsolidateFundsSuccessful represents a ConsolidateFundsSuccessful event raised by the IEnygmaDvpIntegration contract.
type IEnygmaDvpIntegrationConsolidateFundsSuccessful struct {
	Transaction           IDvpProofReceipt
	EnygmaContractAddress common.Address
	Raw                   *types.Log // Blockchain specific contextual infos
}

const IEnygmaDvpIntegrationConsolidateFundsSuccessfulEventName = "ConsolidateFundsSuccessful"

// ContractEventName returns the user-defined event name.
func (IEnygmaDvpIntegrationConsolidateFundsSuccessful) ContractEventName() string {
	return IEnygmaDvpIntegrationConsolidateFundsSuccessfulEventName
}

// UnpackConsolidateFundsSuccessfulEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event ConsolidateFundsSuccessful((((uint256,uint256),(uint256[2],uint256[2]),(uint256,uint256)),uint256[],uint256,uint256[],uint256[],uint256[],uint256) transaction, address enygmaContractAddress)
func (iEnygmaDvpIntegration *IEnygmaDvpIntegration) UnpackConsolidateFundsSuccessfulEvent(log *types.Log) (*IEnygmaDvpIntegrationConsolidateFundsSuccessful, error) {
	event := "ConsolidateFundsSuccessful"
	if log.Topics[0] != iEnygmaDvpIntegration.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(IEnygmaDvpIntegrationConsolidateFundsSuccessful)
	if len(log.Data) > 0 {
		if err := iEnygmaDvpIntegration.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range iEnygmaDvpIntegration.abi.Events[event].Inputs {
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

// IEnygmaDvpIntegrationDepositToDvpSuccessful represents a DepositToDvpSuccessful event raised by the IEnygmaDvpIntegration contract.
type IEnygmaDvpIntegrationDepositToDvpSuccessful struct {
	Commitment            *big.Int
	EnygmaContractAddress common.Address
	Raw                   *types.Log // Blockchain specific contextual infos
}

const IEnygmaDvpIntegrationDepositToDvpSuccessfulEventName = "DepositToDvpSuccessful"

// ContractEventName returns the user-defined event name.
func (IEnygmaDvpIntegrationDepositToDvpSuccessful) ContractEventName() string {
	return IEnygmaDvpIntegrationDepositToDvpSuccessfulEventName
}

// UnpackDepositToDvpSuccessfulEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event DepositToDvpSuccessful(uint256 indexed commitment, address enygmaContractAddress)
func (iEnygmaDvpIntegration *IEnygmaDvpIntegration) UnpackDepositToDvpSuccessfulEvent(log *types.Log) (*IEnygmaDvpIntegrationDepositToDvpSuccessful, error) {
	event := "DepositToDvpSuccessful"
	if log.Topics[0] != iEnygmaDvpIntegration.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(IEnygmaDvpIntegrationDepositToDvpSuccessful)
	if len(log.Data) > 0 {
		if err := iEnygmaDvpIntegration.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range iEnygmaDvpIntegration.abi.Events[event].Inputs {
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

// IEnygmaDvpIntegrationDvpRegistered represents a DvpRegistered event raised by the IEnygmaDvpIntegration contract.
type IEnygmaDvpIntegrationDvpRegistered struct {
	DvpAddress common.Address
	Raw        *types.Log // Blockchain specific contextual infos
}

const IEnygmaDvpIntegrationDvpRegisteredEventName = "DvpRegistered"

// ContractEventName returns the user-defined event name.
func (IEnygmaDvpIntegrationDvpRegistered) ContractEventName() string {
	return IEnygmaDvpIntegrationDvpRegisteredEventName
}

// UnpackDvpRegisteredEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event DvpRegistered(address indexed dvpAddress)
func (iEnygmaDvpIntegration *IEnygmaDvpIntegration) UnpackDvpRegisteredEvent(log *types.Log) (*IEnygmaDvpIntegrationDvpRegistered, error) {
	event := "DvpRegistered"
	if log.Topics[0] != iEnygmaDvpIntegration.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(IEnygmaDvpIntegrationDvpRegistered)
	if len(log.Data) > 0 {
		if err := iEnygmaDvpIntegration.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range iEnygmaDvpIntegration.abi.Events[event].Inputs {
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

// IEnygmaDvpIntegrationVaultIdSet represents a VaultIdSet event raised by the IEnygmaDvpIntegration contract.
type IEnygmaDvpIntegrationVaultIdSet struct {
	VaultId *big.Int
	Raw     *types.Log // Blockchain specific contextual infos
}

const IEnygmaDvpIntegrationVaultIdSetEventName = "VaultIdSet"

// ContractEventName returns the user-defined event name.
func (IEnygmaDvpIntegrationVaultIdSet) ContractEventName() string {
	return IEnygmaDvpIntegrationVaultIdSetEventName
}

// UnpackVaultIdSetEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event VaultIdSet(uint256 vaultId)
func (iEnygmaDvpIntegration *IEnygmaDvpIntegration) UnpackVaultIdSetEvent(log *types.Log) (*IEnygmaDvpIntegrationVaultIdSet, error) {
	event := "VaultIdSet"
	if log.Topics[0] != iEnygmaDvpIntegration.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(IEnygmaDvpIntegrationVaultIdSet)
	if len(log.Data) > 0 {
		if err := iEnygmaDvpIntegration.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range iEnygmaDvpIntegration.abi.Events[event].Inputs {
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

// IEnygmaDvpIntegrationVerifierDepositToDvpRegistered represents a VerifierDepositToDvpRegistered event raised by the IEnygmaDvpIntegration contract.
type IEnygmaDvpIntegrationVerifierDepositToDvpRegistered struct {
	VerifierAddress common.Address
	K               uint8
	Raw             *types.Log // Blockchain specific contextual infos
}

const IEnygmaDvpIntegrationVerifierDepositToDvpRegisteredEventName = "VerifierDepositToDvpRegistered"

// ContractEventName returns the user-defined event name.
func (IEnygmaDvpIntegrationVerifierDepositToDvpRegistered) ContractEventName() string {
	return IEnygmaDvpIntegrationVerifierDepositToDvpRegisteredEventName
}

// UnpackVerifierDepositToDvpRegisteredEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event VerifierDepositToDvpRegistered(address indexed verifierAddress, uint8 k)
func (iEnygmaDvpIntegration *IEnygmaDvpIntegration) UnpackVerifierDepositToDvpRegisteredEvent(log *types.Log) (*IEnygmaDvpIntegrationVerifierDepositToDvpRegistered, error) {
	event := "VerifierDepositToDvpRegistered"
	if log.Topics[0] != iEnygmaDvpIntegration.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(IEnygmaDvpIntegrationVerifierDepositToDvpRegistered)
	if len(log.Data) > 0 {
		if err := iEnygmaDvpIntegration.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range iEnygmaDvpIntegration.abi.Events[event].Inputs {
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

// IEnygmaDvpIntegrationVerifierWithdrawFromDvpRegistered represents a VerifierWithdrawFromDvpRegistered event raised by the IEnygmaDvpIntegration contract.
type IEnygmaDvpIntegrationVerifierWithdrawFromDvpRegistered struct {
	VerifierAddress common.Address
	K               uint8
	Raw             *types.Log // Blockchain specific contextual infos
}

const IEnygmaDvpIntegrationVerifierWithdrawFromDvpRegisteredEventName = "VerifierWithdrawFromDvpRegistered"

// ContractEventName returns the user-defined event name.
func (IEnygmaDvpIntegrationVerifierWithdrawFromDvpRegistered) ContractEventName() string {
	return IEnygmaDvpIntegrationVerifierWithdrawFromDvpRegisteredEventName
}

// UnpackVerifierWithdrawFromDvpRegisteredEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event VerifierWithdrawFromDvpRegistered(address indexed verifierAddress, uint8 k)
func (iEnygmaDvpIntegration *IEnygmaDvpIntegration) UnpackVerifierWithdrawFromDvpRegisteredEvent(log *types.Log) (*IEnygmaDvpIntegrationVerifierWithdrawFromDvpRegistered, error) {
	event := "VerifierWithdrawFromDvpRegistered"
	if log.Topics[0] != iEnygmaDvpIntegration.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(IEnygmaDvpIntegrationVerifierWithdrawFromDvpRegistered)
	if len(log.Data) > 0 {
		if err := iEnygmaDvpIntegration.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range iEnygmaDvpIntegration.abi.Events[event].Inputs {
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

// IEnygmaDvpIntegrationWithdrawFromDvpSuccessful represents a WithdrawFromDvpSuccessful event raised by the IEnygmaDvpIntegration contract.
type IEnygmaDvpIntegrationWithdrawFromDvpSuccessful struct {
	Transaction           IDvpProofReceipt
	EnygmaContractAddress common.Address
	Raw                   *types.Log // Blockchain specific contextual infos
}

const IEnygmaDvpIntegrationWithdrawFromDvpSuccessfulEventName = "WithdrawFromDvpSuccessful"

// ContractEventName returns the user-defined event name.
func (IEnygmaDvpIntegrationWithdrawFromDvpSuccessful) ContractEventName() string {
	return IEnygmaDvpIntegrationWithdrawFromDvpSuccessfulEventName
}

// UnpackWithdrawFromDvpSuccessfulEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event WithdrawFromDvpSuccessful((((uint256,uint256),(uint256[2],uint256[2]),(uint256,uint256)),uint256[],uint256,uint256[],uint256[],uint256[],uint256) transaction, address enygmaContractAddress)
func (iEnygmaDvpIntegration *IEnygmaDvpIntegration) UnpackWithdrawFromDvpSuccessfulEvent(log *types.Log) (*IEnygmaDvpIntegrationWithdrawFromDvpSuccessful, error) {
	event := "WithdrawFromDvpSuccessful"
	if log.Topics[0] != iEnygmaDvpIntegration.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(IEnygmaDvpIntegrationWithdrawFromDvpSuccessful)
	if len(log.Data) > 0 {
		if err := iEnygmaDvpIntegration.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range iEnygmaDvpIntegration.abi.Events[event].Inputs {
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
