// Code generated via abigen V2 - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package EndpointV1

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

// BatchMessage is an auto generated low-level Go binding around an user-defined struct.
type BatchMessage struct {
	ToChainId *big.Int
	To        common.Address
	Data      RaylsMessage
	MessageId [32]byte
}

// BridgedTransferMetadata is an auto generated low-level Go binding around an user-defined struct.
type BridgedTransferMetadata struct {
	AssetType    uint8
	Id           *big.Int
	From         common.Address
	To           common.Address
	TokenAddress common.Address
	Amount       *big.Int
}

// DestinationPayloadRequest is an auto generated low-level Go binding around an user-defined struct.
type DestinationPayloadRequest struct {
	DstChainId  *big.Int
	Destination common.Address
	Payload     []byte
}

// NewResourceMetadata is an auto generated low-level Go binding around an user-defined struct.
type NewResourceMetadata struct {
	Valid              bool
	ResourceDeployType uint8
	Bytecode           []byte
	FactoryTemplate    uint8
	InitializerParams  []byte
}

// RaylsMessage is an auto generated low-level Go binding around an user-defined struct.
type RaylsMessage struct {
	MessageMetadata RaylsMessageMetadata
	Payload         []byte
}

// RaylsMessageMetadata is an auto generated low-level Go binding around an user-defined struct.
type RaylsMessageMetadata struct {
	Valid                     bool
	Nonce                     *big.Int
	NewResourceMetadata       NewResourceMetadata
	ResourceId                [32]byte
	LockData                  []byte
	RevertPayloadDataSender   []byte
	RevertPayloadDataReceiver []byte
	TransferMetadata          BridgedTransferMetadata
	IgnoresNonce              bool
}

// ResourceIdCompletePayloadRequest is an auto generated low-level Go binding around an user-defined struct.
type ResourceIdCompletePayloadRequest struct {
	DstChainId         *big.Int
	ResourceId         [32]byte
	Payload            []byte
	LockData           []byte
	RevertDataSender   []byte
	RevertDataReceiver []byte
	TransferMetadata   BridgedTransferMetadata
}

// ResourceIdPayloadRequest is an auto generated low-level Go binding around an user-defined struct.
type ResourceIdPayloadRequest struct {
	DstChainId *big.Int
	ResourceId [32]byte
	Payload    []byte
}

// EndpointV1MetaData contains all meta data concerning the EndpointV1 contract.
var EndpointV1MetaData = bind.MetaData{
	ABI: "[{\"type\":\"function\",\"name\":\"UPGRADE_INTERFACE_VERSION\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"authority\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"batchMessageSender\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIBatchMessageSender\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"chainId\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"configureContracts\",\"inputs\":[{\"name\":\"_contractFactory\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_participantStorageReplica\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_tokenRegistry\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"configureEndpoint\",\"inputs\":[{\"name\":\"_contractFactory\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_participantStorageReplica\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_tokenRegistry\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_resourceManager\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_messageSender\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_messageReceiver\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_batchMessageSender\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"configureModules\",\"inputs\":[{\"name\":\"_resourceManager\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_messageSender\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_messageReceiver\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_batchMessageSender\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"contractVersion\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"getAddressByResourceId\",\"inputs\":[{\"name\":\"_resourceId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getChainId\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getInboundNonce\",\"inputs\":[{\"name\":\"_srcChainId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getMaxBatchMessages\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getOutboundNonce\",\"inputs\":[{\"name\":\"_dstChainId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getPrivateHubAddress\",\"inputs\":[{\"name\":\"_contractName\",\"type\":\"string\",\"internalType\":\"string\"}],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getPrivateHubId\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getTokenRegistry\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getUserGovernanceAddress\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"initialize\",\"inputs\":[{\"name\":\"_chainId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"_privateHubId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"_maxBatchMessages\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"authority_\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"messageReceiver\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIMessageReceiver\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"messageSender\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIMessageSender\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"privateHubAddress\",\"inputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"privateHubId\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"proxiableUUID\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"receivePayload\",\"inputs\":[{\"name\":\"_srcChainId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"_srcAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_dstAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_raylsMessage\",\"type\":\"tuple\",\"internalType\":\"structRaylsMessage\",\"components\":[{\"name\":\"messageMetadata\",\"type\":\"tuple\",\"internalType\":\"structRaylsMessageMetadata\",\"components\":[{\"name\":\"valid\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"nonce\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"newResourceMetadata\",\"type\":\"tuple\",\"internalType\":\"structNewResourceMetadata\",\"components\":[{\"name\":\"valid\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"resourceDeployType\",\"type\":\"uint8\",\"internalType\":\"enumResourceDeployType\"},{\"name\":\"bytecode\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"factoryTemplate\",\"type\":\"uint8\",\"internalType\":\"enumRaylsBridgeableERC\"},{\"name\":\"initializerParams\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"name\":\"resourceId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"lockData\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"revertPayloadDataSender\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"revertPayloadDataReceiver\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"transferMetadata\",\"type\":\"tuple\",\"internalType\":\"structBridgedTransferMetadata\",\"components\":[{\"name\":\"assetType\",\"type\":\"uint8\",\"internalType\":\"enumRaylsBridgeableERC\"},{\"name\":\"id\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"from\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"tokenAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"ignoresNonce\",\"type\":\"bool\",\"internalType\":\"bool\"}]},{\"name\":\"payload\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"name\":\"_messageId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"registerPrivateHubAddress\",\"inputs\":[{\"name\":\"_contractName\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"_contractAddressOnPrivateHub\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"registerResourceId\",\"inputs\":[{\"name\":\"_resourceId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"_implementationAddress\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"requestNewRaylsViewKeys\",\"inputs\":[{\"name\":\"blockNumber\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"resourceManager\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIResourceManager\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"send\",\"inputs\":[{\"name\":\"_dstChainId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"_destination\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_payload\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"messageId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"send\",\"inputs\":[{\"name\":\"_dstChainId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"_destination\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_payload\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"transferMetadata\",\"type\":\"tuple\",\"internalType\":\"structBridgedTransferMetadata\",\"components\":[{\"name\":\"assetType\",\"type\":\"uint8\",\"internalType\":\"enumRaylsBridgeableERC\"},{\"name\":\"id\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"from\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"tokenAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]}],\"outputs\":[{\"name\":\"messageId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"sendBatch\",\"inputs\":[{\"name\":\"_destinationPayloadRequests\",\"type\":\"tuple[]\",\"internalType\":\"structDestinationPayloadRequest[]\",\"components\":[{\"name\":\"_dstChainId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"_destination\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_payload\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}],\"outputs\":[{\"name\":\"batchId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"sendBatchToResourceId\",\"inputs\":[{\"name\":\"_resourceIdPayloadRequests\",\"type\":\"tuple[]\",\"internalType\":\"structResourceIdCompletePayloadRequest[]\",\"components\":[{\"name\":\"_dstChainId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"_resourceId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"_payload\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"_lockData\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"_revertDataSender\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"_revertDataReceiver\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"transferMetadata\",\"type\":\"tuple\",\"internalType\":\"structBridgedTransferMetadata\",\"components\":[{\"name\":\"assetType\",\"type\":\"uint8\",\"internalType\":\"enumRaylsBridgeableERC\"},{\"name\":\"id\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"from\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"tokenAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]}]}],\"outputs\":[{\"name\":\"batchId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"sendBatchToResourceId\",\"inputs\":[{\"name\":\"_resourceIdPayloadRequests\",\"type\":\"tuple[]\",\"internalType\":\"structResourceIdPayloadRequest[]\",\"components\":[{\"name\":\"_dstChainId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"_resourceId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"_payload\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}],\"outputs\":[{\"name\":\"batchId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"sendToResourceId\",\"inputs\":[{\"name\":\"_dstChainId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"_resourceId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"_payload\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"messageId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"sendToResourceId\",\"inputs\":[{\"name\":\"_dstChainId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"_resourceId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"_payload\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"_lockData\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"_revertDataSender\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"_revertDataReceiver\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"transferMetadata\",\"type\":\"tuple\",\"internalType\":\"structBridgedTransferMetadata\",\"components\":[{\"name\":\"assetType\",\"type\":\"uint8\",\"internalType\":\"enumRaylsBridgeableERC\"},{\"name\":\"id\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"from\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"tokenAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]}],\"outputs\":[{\"name\":\"messageId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"setMaxBatchMessages\",\"inputs\":[{\"name\":\"_maxBatchMessages\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setUserGovernance\",\"inputs\":[{\"name\":\"_userGovernance\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"upgradeToAndCall\",\"inputs\":[{\"name\":\"newImplementation\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"data\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"version\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"pure\"},{\"type\":\"event\",\"name\":\"AuthorityUpdated\",\"inputs\":[{\"name\":\"oldAuthority\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newAuthority\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"EndpointConfigured\",\"inputs\":[{\"name\":\"contractFactory\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"participantStorageReplica\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"tokenRegistry\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"resourceManager\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"messageSender\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"messageReceiver\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"batchMessageSender\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Initialized\",\"inputs\":[{\"name\":\"version\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"MessageBatchDispatched\",\"inputs\":[{\"name\":\"batchId\",\"type\":\"bytes32\",\"indexed\":false,\"internalType\":\"bytes32\"},{\"name\":\"from\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"messages\",\"type\":\"tuple[]\",\"indexed\":false,\"internalType\":\"structBatchMessage[]\",\"components\":[{\"name\":\"toChainId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"to\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"data\",\"type\":\"tuple\",\"internalType\":\"structRaylsMessage\",\"components\":[{\"name\":\"messageMetadata\",\"type\":\"tuple\",\"internalType\":\"structRaylsMessageMetadata\",\"components\":[{\"name\":\"valid\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"nonce\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"newResourceMetadata\",\"type\":\"tuple\",\"internalType\":\"structNewResourceMetadata\",\"components\":[{\"name\":\"valid\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"resourceDeployType\",\"type\":\"uint8\",\"internalType\":\"enumResourceDeployType\"},{\"name\":\"bytecode\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"factoryTemplate\",\"type\":\"uint8\",\"internalType\":\"enumRaylsBridgeableERC\"},{\"name\":\"initializerParams\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"name\":\"resourceId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"lockData\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"revertPayloadDataSender\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"revertPayloadDataReceiver\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"transferMetadata\",\"type\":\"tuple\",\"internalType\":\"structBridgedTransferMetadata\",\"components\":[{\"name\":\"assetType\",\"type\":\"uint8\",\"internalType\":\"enumRaylsBridgeableERC\"},{\"name\":\"id\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"from\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"tokenAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"ignoresNonce\",\"type\":\"bool\",\"internalType\":\"bool\"}]},{\"name\":\"payload\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"name\":\"messageId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"MessageDispatched\",\"inputs\":[{\"name\":\"messageId\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"from\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"toChainId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"to\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"data\",\"type\":\"tuple\",\"indexed\":false,\"internalType\":\"structRaylsMessage\",\"components\":[{\"name\":\"messageMetadata\",\"type\":\"tuple\",\"internalType\":\"structRaylsMessageMetadata\",\"components\":[{\"name\":\"valid\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"nonce\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"newResourceMetadata\",\"type\":\"tuple\",\"internalType\":\"structNewResourceMetadata\",\"components\":[{\"name\":\"valid\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"resourceDeployType\",\"type\":\"uint8\",\"internalType\":\"enumResourceDeployType\"},{\"name\":\"bytecode\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"factoryTemplate\",\"type\":\"uint8\",\"internalType\":\"enumRaylsBridgeableERC\"},{\"name\":\"initializerParams\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"name\":\"resourceId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"lockData\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"revertPayloadDataSender\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"revertPayloadDataReceiver\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"transferMetadata\",\"type\":\"tuple\",\"internalType\":\"structBridgedTransferMetadata\",\"components\":[{\"name\":\"assetType\",\"type\":\"uint8\",\"internalType\":\"enumRaylsBridgeableERC\"},{\"name\":\"id\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"from\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"tokenAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"ignoresNonce\",\"type\":\"bool\",\"internalType\":\"bool\"}]},{\"name\":\"payload\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ModulesConfigured\",\"inputs\":[{\"name\":\"resourceManager\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"messageSender\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"messageReceiver\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"batchMessageSender\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"PrivateHubAddressRegistered\",\"inputs\":[{\"name\":\"contractName\",\"type\":\"string\",\"indexed\":true,\"internalType\":\"string\"},{\"name\":\"contractAddress\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ResourceIdRegistered\",\"inputs\":[{\"name\":\"resourceId\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"implementationAddress\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"UpdateRaylsViewKeysRequest\",\"inputs\":[{\"name\":\"blockNumber\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Upgraded\",\"inputs\":[{\"name\":\"implementation\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"AddressEmptyCode\",\"inputs\":[{\"name\":\"target\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC1967InvalidImplementation\",\"inputs\":[{\"name\":\"implementation\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC1967NonPayable\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"Endpoint__EmptyContractName\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"Endpoint__InvalidPrivateHubAddress\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"FailedCall\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidInitialization\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NotInitializing\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"RaylsAccessManaged__ContractPaused\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"RaylsAccessManaged__InvalidAuthority\",\"inputs\":[{\"name\":\"authority\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"RaylsAccessManaged__MustSchedule\",\"inputs\":[{\"name\":\"caller\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"delay\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"type\":\"error\",\"name\":\"RaylsAccessManaged__Unauthorized\",\"inputs\":[{\"name\":\"caller\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"UUPSUnauthorizedCallContext\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"UUPSUnsupportedProxiableUUID\",\"inputs\":[{\"name\":\"slot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]}]",
	ID:  "EndpointV1",
	Bin: "0x60a06040523060805234801561001457600080fd5b5060805161415761003e600039600081816122fe01528181612327015261245e01526141576000f3fe6080604052600436106101f05760003560e01c80636892f9541161010e578063c166a9a31161009b578063c166a9a31461055c578063d2eb60901461057c578063d4c085291461059c578063d4f951c7146105bc578063d67bdd25146105dc578063de16f0b2146105fc578063e2cae9f51461061c578063f1e229741461063c578063fced8ee814610651578063fe12d9c31461066457600080fd5b80636892f9541461043b57806371f10b8b1461045b578063910e23471461047b5780639a8a05921461049b5780639b76c3f2146104b1578063a0a8e460146104d1578063ad3cb1cc146104e5578063b2958e9a14610516578063b8b8713814610529578063bf7e214f1461054757600080fd5b80633408e4701161018c5780633408e47014610330578063407fa89e1461034557806340b33bfa1461036557806341d71744146103785780634f02198c1461038b5780634f1ef2861461039e57806352d1902d146103b157806354fd4d50146103c657806355615960146103fb5780635afce8f21461041b57600080fd5b8063057838bd146101f55780630b39a9511461022557806311f50c8514610244578063150b375f1461026457806316e72d5f146102775780631a8d9216146102b85780631f8d519d146102da578063202e2374146102fa57806326b3260014610310575b600080fd5b34801561020157600080fd5b506005546001600160a01b03165b60405161021c919061291e565b60405180910390f35b34801561023157600080fd5b506001545b60405190815260200161021c565b34801561025057600080fd5b5061020f61025f366004612932565b610684565b6102366102723660046129a8565b6106f8565b34801561028357600080fd5b5061020f610292366004612b9a565b80516020818301810180516008825292820191909301209152546001600160a01b031681565b3480156102c457600080fd5b506102d86102d3366004612bce565b6107d1565b005b3480156102e657600080fd5b506102d86102f5366004612beb565b610809565b34801561030657600080fd5b5061023660015481565b34801561031c57600080fd5b506102d861032b366004612932565b61091d565b34801561033c57600080fd5b50600054610236565b34801561035157600080fd5b50610236610360366004612c70565b610969565b610236610373366004612cb1565b610bdd565b610236610386366004612d76565b610f59565b610236610399366004612c70565b611013565b6102d86103ac366004612e4d565b611335565b3480156103bd57600080fd5b50610236611354565b3480156103d257600080fd5b50604080518082019091526003815262322e3560e81b60208201525b60405161021c9190612eec565b34801561040757600080fd5b506102d8610416366004612eff565b611371565b34801561042757600080fd5b506102d8610436366004612f50565b61145c565b34801561044757600080fd5b5061020f610456366004612b9a565b6114b1565b34801561046757600080fd5b5060095461020f906001600160a01b031681565b34801561048757600080fd5b506102d8610496366004612932565b611556565b3480156104a757600080fd5b5061023660005481565b3480156104bd57600080fd5b50600c5461020f906001600160a01b031681565b3480156104dd57600080fd5b506001610236565b3480156104f157600080fd5b506103ee604051806040016040528060058152602001640352e302e360dc1b81525081565b610236610524366004612c70565b611571565b34801561053557600080fd5b506006546001600160a01b031661020f565b34801561055357600080fd5b5061020f611761565b34801561056857600080fd5b506102d8610577366004612f9b565b611770565b34801561058857600080fd5b506102d8610597366004613031565b611864565b3480156105a857600080fd5b506102366105b7366004612932565b611915565b3480156105c857600080fd5b506102d86105d7366004613082565b6119ae565b3480156105e857600080fd5b50600a5461020f906001600160a01b031681565b34801561060857600080fd5b50610236610617366004612932565b611ad6565b34801561062857600080fd5b50600b5461020f906001600160a01b031681565b34801561064857600080fd5b50600754610236565b61023661065f3660046130a7565b611b32565b34801561067057600080fd5b506102d861067f3660046131ee565b611c01565b6009546040516311f50c8560e01b8152600481018390526000916001600160a01b0316906311f50c8590602401602060405180830381865afa1580156106ce573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906106f29190613398565b92915050565b6000610710336000356001600160e01b031916611ca5565b610718612854565b610720612884565b6107c4604051806101200160405280898152602001886001600160a01b0316815260200187878080601f0160208091040260200160405190810160405280939291908181526020018383808284376000920182905250938552505050602080830182905260408084018890528051808301825283815260608501528051808301825283815260808501528051918201905290815260a082015260c001839052611df0565b925050505b949350505050565b6107e7336000356001600160e01b031916611ca5565b600680546001600160a01b0319166001600160a01b0392909216919091179055565b6000610813611fd6565b805490915060ff600160401b82041615906001600160401b031660008115801561083a5750825b90506000826001600160401b031660011480156108565750303b155b905081158015610864575080155b156108825760405163f92ee8a960e01b815260040160405180910390fd5b845467ffffffffffffffff1916600117855583156108ac57845460ff60401b1916600160401b1785555b6108b4611fff565b6001889055600089905560078790556108cc86612009565b831561091257845460ff60401b19168555604051600181527fc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d29060200160405180910390a15b505050505050505050565b610933336000356001600160e01b031916611ca5565b6040518181527fc21ad5212c3a721b356a7ecc7259eb485f1d3c6b20050fee7413ae000e5f19659060200160405180910390a150565b6000610981336000356001600160e01b031916611ca5565b6007548210610a155760405162461bcd60e51b815260206004820152604f60248201527f456e64706f696e7456313a20546865206d6178206e756d626572206f6620747260448201527f616e73616374696f6e7320616c6c6f77656420696e206120626174636820686160648201526e1cc81899595b88195e18d959591959608a1b608482015260a4015b60405180910390fd5b6000826001600160401b03811115610a2f57610a2f612a03565b604051908082528060200260200182016040528015610a6857816020015b610a556128bb565b815260200190600190039081610a4d5790505b50905060005b83811015610bd357610a7e612854565b610a86612884565b604051806101200160405280888886818110610aa457610aa46133b5565b9050602002810190610ab691906133cb565b358152602001888886818110610ace57610ace6133b5565b9050602002810190610ae091906133cb565b610af1906040810190602001612bce565b6001600160a01b03168152602001888886818110610b1157610b116133b5565b9050602002810190610b2391906133cb565b610b319060408101906133eb565b8080601f0160208091040260200160405190810160405280939291908181526020018383808284376000920182905250938552505050602080830182905260408084018790528051808301825283815260608501528051808301825283815260808501528051918201905290815260a082015260c0018290528451859085908110610bbe57610bbe6133b5565b60209081029190910101525050600101610a6e565b506107c98161204a565b6000610bf5336000356001600160e01b031916611ca5565b610bfd612884565b610c05612854565b86600003610eb4576000600460009054906101000a90046001600160a01b03166001600160a01b031663195ec9ee6040518163ffffffff1660e01b8152600401600060405180830381865afa158015610c62573d6000803e3d6000fd5b505050506040513d6000823e601f3d908101601f19168201604052610c8a91908101906134c2565b90506000805b8251811015610ce857600054838281518110610cae57610cae6133b5565b60200260200101516000015103610ce057828181518110610cd157610cd16133b5565b602002602001015160e0015191505b600101610c90565b5080610d4f5760405162461bcd60e51b815260206004820152603060248201527f5061727469636970616e74206973206e6f7420616c6c6f77656420746f20627260448201526f6f616463617374206d6573736167657360801b6064820152608401610a0c565b8015610eab5760005b8251811015610ea9576000838281518110610d7557610d756133b5565b60200260200101516000015190506000548114158015610d9757506103e78114155b15610ea05760006001858481518110610db257610db26133b5565b6020026020010151604001516003811115610dcf57610dcf613620565b1490508015610e9e57610e9b60405180610120016040528084815260200160006001600160a01b031681526020018c8c8080601f016020809104026020016040519081016040528093929190818152602001838380828437600081840152601f19601f8201169050808301925050505050505081526020018d815260200188815260200160405180602001604052806000815250815260200160405180602001604052806000815250815260200160405180602001604052806000815250815260200189815250611df0565b97505b505b50600101610d58565b505b505050506107c9565b6107c460405180610120016040528089815260200160006001600160a01b0316815260200187878080601f016020809104026020016040519081016040528093929190818152602001838380828437600092018290525093855250505060208083018b905260408084018790528051808301825283815260608501528051808301825283815260808501528051918201905290815260a082015260c001849052611df0565b6000610f71336000356001600160e01b031916611ca5565b610f79612854565b6110056040518061012001604052808c815260200160006001600160a01b031681526020018a8a8080601f016020809104026020016040519081016040528093929190818152602001838380828437600092019190915250505090825250602081018c905260408101849052606081018990526080810188905260a0810187905260c001859052611df0565b9a9950505050505050505050565b600061102b336000356001600160e01b031916611ca5565b600754821061104c5760405162461bcd60e51b8152600401610a0c90613636565b6000826001600160401b0381111561106657611066612a03565b60405190808252806020026020018201604052801561109f57816020015b61108c6128bb565b8152602001906001900390816110845790505b50905060005b83811015610bd3576110b5612854565b6040518061012001604052808787858181106110d3576110d36133b5565b90506020028101906110e5919061369f565b35815260006020820152604001878785818110611104576111046133b5565b9050602002810190611116919061369f565b6111249060408101906133eb565b8080601f016020809104026020016040519081016040528093929190818152602001838380828437600092019190915250505090825250602001878785818110611170576111706133b5565b9050602002810190611182919061369f565b6020013581526020018281526020018787858181106111a3576111a36133b5565b90506020028101906111b5919061369f565b6111c39060608101906133eb565b8080601f01602080910402602001604051908101604052809392919081815260200183838082843760009201919091525050509082525060200187878581811061120f5761120f6133b5565b9050602002810190611221919061369f565b61122f9060808101906133eb565b8080601f01602080910402602001604051908101604052809392919081815260200183838082843760009201919091525050509082525060200187878581811061127b5761127b6133b5565b905060200281019061128d919061369f565b61129b9060a08101906133eb565b8080601f0160208091040260200160405190810160405280939291908181526020018383808284376000920191909152505050908252506020018787858181106112e7576112e76133b5565b90506020028101906112f9919061369f565b60c00180360381019061130c91906136b6565b815250838381518110611321576113216133b5565b6020908102919091010152506001016110a5565b61133d6122f3565b61134682612381565b611350828261239a565b5050565b600061135e612453565b5060008051602061410283398151915290565b611387336000356001600160e01b031916611ca5565b6001600160a01b0381166113ae57604051639e0f82d960e01b815260040160405180910390fd5b81516000036113d057604051634b1b08d760e11b815260040160405180910390fd5b806001600160a01b0316826040516113e891906136d2565b604051908190038120907f3cd326fcdccf15c7b507c1d59fa37a9f36d686fefe59709a315bca9701aab48190600090a38060088360405161142991906136d2565b90815260405190819003602001902080546001600160a01b03929092166001600160a01b03199092169190911790555050565b611472336000356001600160e01b031916611ca5565b600380546001600160a01b039485166001600160a01b031991821617909155600480549385169382169390931790925560058054919093169116179055565b6000806001600160a01b03166008836040516114cd91906136d2565b90815260405160209181900382018120546001600160a01b03169290921415916114f9918591016136e4565b604051602081830303815290604052906115265760405162461bcd60e51b8152600401610a0c9190612eec565b5060088260405161153791906136d2565b908152604051908190036020019020546001600160a01b031692915050565b61156c336000356001600160e01b031916611ca5565b600755565b6000611589336000356001600160e01b031916611ca5565b60075482106115aa5760405162461bcd60e51b8152600401610a0c90613636565b6000826001600160401b038111156115c4576115c4612a03565b6040519080825280602002602001820160405280156115fd57816020015b6115ea6128bb565b8152602001906001900390816115e25790505b50905060005b83811015610bd357611613612854565b61161b612884565b604051806101200160405280888886818110611639576116396133b5565b905060200281019061164b91906133cb565b3581526000602082015260400188888681811061166a5761166a6133b5565b905060200281019061167c91906133cb565b61168a9060408101906133eb565b8080601f0160208091040260200160405190810160405280939291908181526020018383808284376000920191909152505050908252506020018888868181106116d6576116d66133b5565b90506020028101906116e891906133cb565b6020013581526020018381526020016040518060200160405280600081525081526020016040518060200160405280600081525081526020016040518060200160405280600081525081526020018281525084848151811061174c5761174c6133b5565b60209081029190910101525050600101611603565b600061176b61249c565b905090565b611786336000356001600160e01b031916611ca5565b600380546001600160a01b038981166001600160a01b03199283168117909355600480548a83169084168117909155600580548a84169085168117909155600980548a85169086168117909155600a80548a86169087168117909155600b80548a87169088168117909155600c8054968a16969097168617909655604080519788526020880194909452928601919091526060850152608084015260a083019190915260c08201527f2ad87a510d1161ee24e65e0d5b00f09ae5da635623267b550211319fecd829399060e00160405180910390a150505050505050565b61187a336000356001600160e01b031916611ca5565b600980546001600160a01b038681166001600160a01b03199283168117909355600a80548783169084168117909155600b80548784169085168117909155600c805493871693909416831790935560408051948552602085019190915283019190915260608201527fbc8b4d996f7bd7240f8b484b7c312fdf18fbb0e45105cc792927c7aaa4df5d5a9060800160405180910390a150505050565b600a546000906001600160a01b03166119405760405162461bcd60e51b8152600401610a0c90613744565b600a5460405163d4c0852960e01b8152600481018490526001600160a01b039091169063d4c08529906024015b602060405180830381865afa15801561198a573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906106f29190613793565b6119c4336000356001600160e01b031916611ca5565b6009546001600160a01b0316611a365760405162461bcd60e51b815260206004820152603160248201527f456e64706f696e7456313a205265736f757263654d616e61676572206d6f64756044820152701b19481b9bdd0818dbdb999a59dd5c9959607a1b6064820152608401610a0c565b60095460405163d4f951c760e01b8152600481018490526001600160a01b0383811660248301529091169063d4f951c790604401600060405180830381600087803b158015611a8457600080fd5b505af1158015611a98573d6000803e3d6000fd5b50506040516001600160a01b03841692508491507f28edd4827c66a9d7ff43b4da996da795733e847edd27ce220c091f5b1fd2cab390600090a35050565b600b546000906001600160a01b0316611b015760405162461bcd60e51b8152600401610a0c906137ac565b600b54604051636f0b785960e11b8152600481018490526001600160a01b039091169063de16f0b29060240161196d565b6000611b4a336000356001600160e01b031916611ca5565b611b52612854565b611bf6604051806101200160405280898152602001886001600160a01b0316815260200187878080601f0160208091040260200160405190810160405280939291908181526020018383808284376000920182905250938552505050602080830182905260408084018790528051808301825283815260608501528051808301825283815260808501528051918201905290815260a082015260c001859052611df0565b979650505050505050565b611c17336000356001600160e01b031916611ca5565b600b546001600160a01b0316611c3f5760405162461bcd60e51b8152600401610a0c906137ac565b600b5460405163fe12d9c360e01b81526001600160a01b039091169063fe12d9c390611c7790889088908890889088906004016139b2565b600060405180830381600087803b158015611c9157600080fd5b505af1158015610912573d6000803e3d6000fd5b6000611caf6124b5565b80549091506001600160a01b031680611cde576000604051638944034760e01b8152600401610a0c919061291e565b60405163b700961360e01b81526001600160a01b0385811660048301523060248301526001600160e01b031985166044830152600091829182919085169063b700961390606401606060405180830381865afa158015611d42573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190611d6691906139f6565b92509250925082611de7578015611d905760405163cc9855ad60e01b815260040160405180910390fd5b63ffffffff821615611dcc5760405163a426878960e01b81526001600160a01b038816600482015263ffffffff83166024820152604401610a0c565b86604051632ecd3d0360e21b8152600401610a0c919061291e565b50505050505050565b600a546000906001600160a01b0316611e1b5760405162461bcd60e51b8152600401610a0c90613744565b600b546001600160a01b0316611e435760405162461bcd60e51b8152600401610a0c906137ac565b600a54604051631372a16b60e21b815260009182916001600160a01b0390911690634dca85ac90611e7a9087903390600401613b0e565b6000604051808303816000875af1158015611e99573d6000803e3d6000fd5b505050506040513d6000823e601f3d908101601f19168201604052611ec19190810190613da7565b5091509150600054846000015103611f4957600b54600054602086015160405163fe12d9c360e01b81526001600160a01b039093169263fe12d9c392611f12929091339190889088906004016139b2565b600060405180830381600087803b158015611f2c57600080fd5b505af1158015611f40573d6000803e3d6000fd5b50505050611fcf565b6000611f6082338760000151886020015187612517565b9050818114611fcd5760405162461bcd60e51b815260206004820152603360248201527f456e64706f696e7456313a2064697370617463684d6573736167652073686f756044820152721b19081c995d1d5c9b881b595cdcd859d95259606a1b6064820152608401610a0c565b505b9392505050565b6000807ff0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a006106f2565b612007612568565b565b60006120136124b5565b80549091506001600160a01b0316156120415781604051638944034760e01b8152600401610a0c919061291e565b6113508261258d565b600c546000906001600160a01b03166120c25760405162461bcd60e51b815260206004820152603460248201527f456e64706f696e7456313a2042617463684d65737361676553656e646572206d6044820152731bd91d5b19481b9bdd0818dbdb999a59dd5c995960621b6064820152608401610a0c565b600b546001600160a01b03166120ea5760405162461bcd60e51b8152600401610a0c906137ac565b600c546000805460405163ae51c66360e01b8152919283926001600160a01b039091169163ae51c66391612125918891339190600401613df6565b6000604051808303816000875af1158015612144573d6000803e3d6000fd5b505050506040513d6000823e601f3d908101601f1916820160405261216c9190810190613edc565b509150915060005b825181101561227557600054838281518110612192576121926133b5565b6020026020010151600001510361226d57600b5460005484516001600160a01b039092169163fe12d9c3919033908790869081106121d2576121d26133b5565b6020026020010151602001518786815181106121f0576121f06133b5565b60200260200101516040015188878151811061220e5761220e6133b5565b6020026020010151606001516040518663ffffffff1660e01b815260040161223a9594939291906139b2565b600060405180830381600087803b15801561225457600080fd5b505af1158015612268573d6000803e3d6000fd5b505050505b600101612174565b50600061228382338561261d565b9050818114611fcd5760405162461bcd60e51b815260206004820152603660248201527f456e64706f696e7456313a2064697370617463684d6573736167654261746368604482015275081cda1bdd5b19081c995d1d5c9b8818985d18da125960521b6064820152608401610a0c565b306001600160a01b037f000000000000000000000000000000000000000000000000000000000000000016148061236357507f00000000000000000000000000000000000000000000000000000000000000006001600160a01b0316612357612662565b6001600160a01b031614155b156120075760405163703e46dd60e11b815260040160405180910390fd5b612397336000356001600160e01b031916611ca5565b50565b816001600160a01b03166352d1902d6040518163ffffffff1660e01b8152600401602060405180830381865afa9250505080156123f4575060408051601f3d908101601f191682019092526123f191810190613793565b60015b6124135781604051634c9c8ce360e01b8152600401610a0c919061291e565b600080516020614102833981519152811461244457604051632a87526960e21b815260048101829052602401610a0c565b61244e8383612678565b505050565b306001600160a01b037f000000000000000000000000000000000000000000000000000000000000000016146120075760405163703e46dd60e11b815260040160405180910390fd5b60006124a66124b5565b546001600160a01b0316919050565b60008060ff196124e660017f2d9a26166e6ed6927c3e14d89c4d572a357f984344fcaed3c38c3ff5cdb83f35614014565b6040516020016124f891815260200190565b60408051601f1981840301815291905280516020909101201692915050565b600083856001600160a01b0316877fa45ed7e52c0a5c8eaddf0cda36a0b3d1e4f928a09a993ea8c9195bf8653510608686604051612556929190614035565b60405180910390a45093949350505050565b6125706126ce565b61200757604051631afcd79f60e31b815260040160405180910390fd5b6001600160a01b0381166125b65780604051638944034760e01b8152600401610a0c919061291e565b60006125c06124b5565b80546040519192506001600160a01b03808516929116907fa3396fd7f6e0a21b50e5089d2da70d5ac0a3bbbd1f617a93f134b7638998019890600090a380546001600160a01b0319166001600160a01b0392909216919091179055565b60007f70d44ffd024502b749906da247dea9d3c21ae5f0754ece8ca7512dde985c1aca84848460405161265293929190614059565b60405180910390a1509192915050565b60006000805160206141028339815191526124a6565b612681826126e8565b6040516001600160a01b038316907fbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b90600090a28051156126c65761244e8282612744565b6113506127ba565b60006126d8611fd6565b54600160401b900460ff16919050565b806001600160a01b03163b6000036127155780604051634c9c8ce360e01b8152600401610a0c919061291e565b60008051602061410283398151915280546001600160a01b0319166001600160a01b0392909216919091179055565b6060600080846001600160a01b03168460405161276191906136d2565b600060405180830381855af49150503d806000811461279c576040519150601f19603f3d011682016040523d82523d6000602084013e6127a1565b606091505b50915091506127b18583836127d9565b95945050505050565b34156120075760405163b398979f60e01b815260040160405180910390fd5b6060826127ee576127e98261282c565b611fcf565b815115801561280557506001600160a01b0384163b155b156128255783604051639996b31560e01b8152600401610a0c919061291e565b5080611fcf565b80511561283b57805160208201fd5b60405163d6bda27560e01b815260040160405180910390fd5b6040805160a081019091526000808252602082019081526060602082015260400160008152602001606081525090565b6040805160c08101909152806000815260006020820181905260408201819052606082018190526080820181905260a09091015290565b6040518061012001604052806000815260200160006001600160a01b0316815260200160608152602001600080191681526020016128f7612854565b8152602001606081526020016060815260200160608152602001612919612884565b905290565b6001600160a01b0391909116815260200190565b60006020828403121561294457600080fd5b5035919050565b6001600160a01b038116811461239757600080fd5b60008083601f84011261297257600080fd5b5081356001600160401b0381111561298957600080fd5b6020830191508360208285010111156129a157600080fd5b9250929050565b600080600080606085870312156129be57600080fd5b8435935060208501356129d08161294b565b925060408501356001600160401b038111156129eb57600080fd5b6129f787828801612960565b95989497509550505050565b634e487b7160e01b600052604160045260246000fd5b60405160c081016001600160401b0381118282101715612a3b57612a3b612a03565b60405290565b60405160a081016001600160401b0381118282101715612a3b57612a3b612a03565b604080519081016001600160401b0381118282101715612a3b57612a3b612a03565b60405161012081016001600160401b0381118282101715612a3b57612a3b612a03565b60405161010081016001600160401b0381118282101715612a3b57612a3b612a03565b604051608081016001600160401b0381118282101715612a3b57612a3b612a03565b604051601f8201601f191681016001600160401b0381118282101715612b1557612b15612a03565b604052919050565b60006001600160401b03821115612b3657612b36612a03565b50601f01601f191660200190565b600082601f830112612b5557600080fd5b8135612b68612b6382612b1d565b612aed565b818152846020838601011115612b7d57600080fd5b816020850160208301376000918101602001919091529392505050565b600060208284031215612bac57600080fd5b81356001600160401b03811115612bc257600080fd5b6107c984828501612b44565b600060208284031215612be057600080fd5b8135611fcf8161294b565b60008060008060808587031215612c0157600080fd5b8435935060208501359250604085013591506060850135612c218161294b565b939692955090935050565b60008083601f840112612c3e57600080fd5b5081356001600160401b03811115612c5557600080fd5b6020830191508360208260051b85010111156129a157600080fd5b60008060208385031215612c8357600080fd5b82356001600160401b03811115612c9957600080fd5b612ca585828601612c2c565b90969095509350505050565b60008060008060608587031215612cc757600080fd5b843593506020850135925060408501356001600160401b038111156129eb57600080fd5b600d811061239757600080fd5b600060c08284031215612d0a57600080fd5b612d12612a19565b90508135612d1f81612ceb565b8152602082810135908201526040820135612d398161294b565b60408201526060820135612d4c8161294b565b60608201526080820135612d5f8161294b565b8060808301525060a082013560a082015292915050565b600080600080600080600080610180898b031215612d9357600080fd5b883597506020890135965060408901356001600160401b0380821115612db857600080fd5b612dc48c838d01612960565b909850965060608b0135915080821115612ddd57600080fd5b612de98c838d01612b44565b955060808b0135915080821115612dff57600080fd5b612e0b8c838d01612b44565b945060a08b0135915080821115612e2157600080fd5b50612e2e8b828c01612b44565b925050612e3e8a60c08b01612cf8565b90509295985092959890939650565b60008060408385031215612e6057600080fd5b8235612e6b8161294b565b915060208301356001600160401b03811115612e8657600080fd5b612e9285828601612b44565b9150509250929050565b60005b83811015612eb7578181015183820152602001612e9f565b50506000910152565b60008151808452612ed8816020860160208601612e9c565b601f01601f19169290920160200192915050565b602081526000611fcf6020830184612ec0565b60008060408385031215612f1257600080fd5b82356001600160401b03811115612f2857600080fd5b612f3485828601612b44565b9250506020830135612f458161294b565b809150509250929050565b600080600060608486031215612f6557600080fd5b8335612f708161294b565b92506020840135612f808161294b565b91506040840135612f908161294b565b809150509250925092565b600080600080600080600060e0888a031215612fb657600080fd5b8735612fc18161294b565b96506020880135612fd18161294b565b95506040880135612fe18161294b565b94506060880135612ff18161294b565b935060808801356130018161294b565b925060a08801356130118161294b565b915060c08801356130218161294b565b8091505092959891949750929550565b6000806000806080858703121561304757600080fd5b84356130528161294b565b935060208501356130628161294b565b925060408501356130728161294b565b91506060850135612c218161294b565b6000806040838503121561309557600080fd5b823591506020830135612f458161294b565b600080600080600061012086880312156130c057600080fd5b8535945060208601356130d28161294b565b935060408601356001600160401b038111156130ed57600080fd5b6130f988828901612960565b909450925061310d90508760608801612cf8565b90509295509295909350565b801515811461239757600080fd5b803561313281613119565b919050565b6002811061239757600080fd5b600060a0828403121561315657600080fd5b61315e612a41565b9050813561316b81613119565b8152602082013561317b81613137565b602082015260408201356001600160401b038082111561319a57600080fd5b6131a685838601612b44565b6040840152606084013591506131bb82612ceb565b81606084015260808401359150808211156131d557600080fd5b506131e284828501612b44565b60808301525092915050565b600080600080600060a0868803121561320657600080fd5b8535945060208601356132188161294b565b935060408601356132288161294b565b925060608601356001600160401b038082111561324457600080fd5b908701906040828a03121561325857600080fd5b613260612a63565b82358281111561326f57600080fd5b83016101c0818c03121561328257600080fd5b61328a612a85565b61329382613127565b8152602082013560208201526040820135848111156132b157600080fd5b6132bd8d828501613144565b604083015250606082013560608201526080820135848111156132df57600080fd5b6132eb8d828501612b44565b60808301525060a08201358481111561330357600080fd5b61330f8d828501612b44565b60a08301525060c08201358481111561332757600080fd5b6133338d828501612b44565b60c0830152506133468c60e08401612cf8565b60e08201526133586101a08301613127565b61010082015282525060208301358281111561337357600080fd5b61337f8b828601612b44565b6020830152509699959850939660800135949350505050565b6000602082840312156133aa57600080fd5b8151611fcf8161294b565b634e487b7160e01b600052603260045260246000fd5b60008235605e198336030181126133e157600080fd5b9190910192915050565b6000808335601e1984360301811261340257600080fd5b8301803591506001600160401b0382111561341c57600080fd5b6020019150368190038213156129a157600080fd5b60006001600160401b0382111561344a5761344a612a03565b5060051b60200190565b80516003811061313257600080fd5b80516004811061313257600080fd5b600082601f83011261348357600080fd5b8151613491612b6382612b1d565b8181528460208386010111156134a657600080fd5b6107c9826020830160208701612e9c565b805161313281613119565b600060208083850312156134d557600080fd5b82516001600160401b03808211156134ec57600080fd5b818501915085601f83011261350057600080fd5b815161350e612b6382613431565b81815260059190911b8301840190848101908883111561352d57600080fd5b8585015b838110156136135780518581111561354857600080fd5b8601610100818c03601f190181131561356057600080fd5b613568612aa8565b89830151815261357a60408401613454565b8a82015261358a60608401613463565b60408201526080830151888111156135a25760008081fd5b6135b08e8c83870101613472565b60608301525060a080840151898111156135ca5760008081fd5b6135d88f8d83880101613472565b60808401525060c0808501518284015260e091508185015181840152506136008385016134b7565b9082015285525050918601918601613531565b5098975050505050505050565b634e487b7160e01b600052602160045260246000fd5b60208082526043908201527f546865206d6178206e756d626572206f66207472616e73616374696f6e73206160408201527f6c6c6f77656420696e206120626174636820686173206265656e20657863656560608201526219195960ea1b608082015260a00190565b6000823561017e198336030181126133e157600080fd5b600060c082840312156136c857600080fd5b611fcf8383612cf8565b600082516133e1818460208701612e9c565b75507269766174652068756220636f6e7472616374202760501b815260008251613715816016850160208701612e9c565b7709c81b9bdd081b585c1c1959081bdb88195b991c1bda5b9d60421b6016939091019283015250602e01919050565b6020808252602f908201527f456e64706f696e7456313a204d65737361676553656e646572206d6f64756c6560408201526e081b9bdd0818dbdb999a59dd5c9959608a1b606082015260800190565b6000602082840312156137a557600080fd5b5051919050565b60208082526031908201527f456e64706f696e7456313a204d6573736167655265636569766572206d6f64756040820152701b19481b9bdd0818dbdb999a59dd5c9959607a1b606082015260800190565b600d811061380d5761380d613620565b9052565b805115158252600060208201516002811061382e5761382e613620565b80602085015250604082015160a0604085015261384e60a0850182612ec0565b9050606083015161386260608601826137fd565b50608083015184820360808601526127b18282612ec0565b6138858282516137fd565b60208101516020830152604081015160018060a01b038082166040850152806060840151166060850152806080840151166080850152505060a081015160a08301525050565b60008151604084526138e260408501825115159052565b6020810151606085015260408101516101c06080860152613907610200860182613811565b9050606082015160a08601526080820151603f19808784030160c088015261392f8383612ec0565b925060a08401519150808784030160e088015261394c8383612ec0565b925060c0840151915061010081888503018189015261396b8484612ec0565b935060e0850151925061398261012089018461387a565b93909301518015156101e088015292506139999050565b6020840151915084810360208601526127b18183612ec0565b8581526001600160a01b0385811660208301528416604082015260a0606082018190526000906139e4908301856138cb565b90508260808301529695505050505050565b600080600060608486031215613a0b57600080fd5b8351613a1681613119565b602085015190935063ffffffff81168114613a3057600080fd5b6040850151909250612f9081613119565b60006101c0825184526020830151613a6460208601826001600160a01b03169052565b506040830151816040860152613a7c82860182612ec0565b9150506060830151606085015260808301518482036080860152613aa08282613811565b91505060a083015184820360a0860152613aba8282612ec0565b91505060c083015184820360c0860152613ad48282612ec0565b91505060e083015184820360e0860152613aee8282612ec0565b91505061010080840151613b048287018261387a565b5090949350505050565b604081526000613b216040830185613a41565b905060018060a01b03831660208301529392505050565b600060a08284031215613b4a57600080fd5b613b52612a41565b90508151613b5f81613119565b81526020820151613b6f81613137565b602082015260408201516001600160401b0380821115613b8e57600080fd5b613b9a85838601613472565b604084015260608401519150613baf82612ceb565b8160608401526080840151915080821115613bc957600080fd5b506131e284828501613472565b600060c08284031215613be857600080fd5b613bf0612a19565b90508151613bfd81612ceb565b8152602082810151908201526040820151613c178161294b565b60408201526060820151613c2a8161294b565b60608201526080820151613c3d8161294b565b8060808301525060a082015160a082015292915050565b600060408284031215613c6657600080fd5b613c6e612a63565b905081516001600160401b0380821115613c8757600080fd5b908301906101c08286031215613c9c57600080fd5b613ca4612a85565b613cad836134b7565b815260208301516020820152604083015182811115613ccb57600080fd5b613cd787828601613b38565b60408301525060608301516060820152608083015182811115613cf957600080fd5b613d0587828601613472565b60808301525060a083015182811115613d1d57600080fd5b613d2987828601613472565b60a08301525060c083015182811115613d4157600080fd5b613d4d87828601613472565b60c083015250613d608660e08501613bd6565b60e0820152613d726101a084016134b7565b61010082015283526020840151915080821115613d8e57600080fd5b50613d9b84828501613472565b60208301525092915050565b600080600060608486031215613dbc57600080fd5b83516001600160401b03811115613dd257600080fd5b613dde86828701613c54565b93505060208401519150604084015190509250925092565b6000606082016060835280865180835260808501915060808160051b8601019250602080890160005b83811015613e4d57607f19888703018552613e3b868351613a41565b95509382019390820190600101613e1f565b505050506001600160a01b039590951660208401526040909201929092529392505050565b600082601f830112613e8357600080fd5b81516020613e93612b6383613431565b8083825260208201915060208460051b870101935086841115613eb557600080fd5b602086015b84811015613ed15780518352918301918301613eba565b509695505050505050565b600080600060608486031215613ef157600080fd5b83516001600160401b0380821115613f0857600080fd5b818601915086601f830112613f1c57600080fd5b81516020613f2c612b6383613431565b82815260059290921b8401810191818101908a841115613f4b57600080fd5b8286015b84811015613fdc57805186811115613f6657600080fd5b87016080818e03601f19011215613f7c57600080fd5b613f84612acb565b8582015181526040820151613f988161294b565b81870152606082015188811115613faf5760008081fd5b613fbd8f8883860101613c54565b6040830152506080919091015160608201528352918301918301613f4f565b509189015160408a01519298509650909350505080821115613ffd57600080fd5b5061400a86828701613e72565b9150509250925092565b818103818111156106f257634e487b7160e01b600052601160045260246000fd5b6001600160a01b03831681526040602082018190526000906107c9908301846138cb565b60006060808301868452602060018060a01b03808816828701526040606060408801528388518086526080955060808901915060808160051b8a0101858b0160005b838110156140ee57607f198c840301855281518051845287898201511689850152868101518a888601526140d18b8601826138cb565b918c0151948c01949094529488019492509087019060010161409b565b50909d9c5050505050505050505050505056fe360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbca264697066735822122098f7f1c0716730bbef097edfb5bad9cf21854c5bd1854a349af9000d63952ea564736f6c63430008180033",
}

// EndpointV1 is an auto generated Go binding around an Ethereum contract.
type EndpointV1 struct {
	abi abi.ABI
}

// NewEndpointV1 creates a new instance of EndpointV1.
func NewEndpointV1() *EndpointV1 {
	parsed, err := EndpointV1MetaData.ParseABI()
	if err != nil {
		panic(errors.New("invalid ABI: " + err.Error()))
	}
	return &EndpointV1{abi: *parsed}
}

// Instance creates a wrapper for a deployed contract instance at the given address.
// Use this to create the instance object passed to abigen v2 library functions Call, Transact, etc.
func (c *EndpointV1) Instance(backend bind.ContractBackend, addr common.Address) *bind.BoundContract {
	return bind.NewBoundContract(addr, c.abi, backend, backend, backend)
}

// PackUPGRADEINTERFACEVERSION is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xad3cb1cc.
//
// Solidity: function UPGRADE_INTERFACE_VERSION() view returns(string)
func (endpointV1 *EndpointV1) PackUPGRADEINTERFACEVERSION() []byte {
	enc, err := endpointV1.abi.Pack("UPGRADE_INTERFACE_VERSION")
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackUPGRADEINTERFACEVERSION is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xad3cb1cc.
//
// Solidity: function UPGRADE_INTERFACE_VERSION() view returns(string)
func (endpointV1 *EndpointV1) UnpackUPGRADEINTERFACEVERSION(data []byte) (string, error) {
	out, err := endpointV1.abi.Unpack("UPGRADE_INTERFACE_VERSION", data)
	if err != nil {
		return *new(string), err
	}
	out0 := *abi.ConvertType(out[0], new(string)).(*string)
	return out0, err
}

// PackAuthority is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xbf7e214f.
//
// Solidity: function authority() view returns(address)
func (endpointV1 *EndpointV1) PackAuthority() []byte {
	enc, err := endpointV1.abi.Pack("authority")
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackAuthority is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xbf7e214f.
//
// Solidity: function authority() view returns(address)
func (endpointV1 *EndpointV1) UnpackAuthority(data []byte) (common.Address, error) {
	out, err := endpointV1.abi.Unpack("authority", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, err
}

// PackBatchMessageSender is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x9b76c3f2.
//
// Solidity: function batchMessageSender() view returns(address)
func (endpointV1 *EndpointV1) PackBatchMessageSender() []byte {
	enc, err := endpointV1.abi.Pack("batchMessageSender")
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackBatchMessageSender is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x9b76c3f2.
//
// Solidity: function batchMessageSender() view returns(address)
func (endpointV1 *EndpointV1) UnpackBatchMessageSender(data []byte) (common.Address, error) {
	out, err := endpointV1.abi.Unpack("batchMessageSender", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, err
}

// PackChainId is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x9a8a0592.
//
// Solidity: function chainId() view returns(uint256)
func (endpointV1 *EndpointV1) PackChainId() []byte {
	enc, err := endpointV1.abi.Pack("chainId")
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackChainId is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x9a8a0592.
//
// Solidity: function chainId() view returns(uint256)
func (endpointV1 *EndpointV1) UnpackChainId(data []byte) (*big.Int, error) {
	out, err := endpointV1.abi.Unpack("chainId", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, err
}

// PackConfigureContracts is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x5afce8f2.
//
// Solidity: function configureContracts(address _contractFactory, address _participantStorageReplica, address _tokenRegistry) returns()
func (endpointV1 *EndpointV1) PackConfigureContracts(contractFactory common.Address, participantStorageReplica common.Address, tokenRegistry common.Address) []byte {
	enc, err := endpointV1.abi.Pack("configureContracts", contractFactory, participantStorageReplica, tokenRegistry)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackConfigureEndpoint is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xc166a9a3.
//
// Solidity: function configureEndpoint(address _contractFactory, address _participantStorageReplica, address _tokenRegistry, address _resourceManager, address _messageSender, address _messageReceiver, address _batchMessageSender) returns()
func (endpointV1 *EndpointV1) PackConfigureEndpoint(contractFactory common.Address, participantStorageReplica common.Address, tokenRegistry common.Address, resourceManager common.Address, messageSender common.Address, messageReceiver common.Address, batchMessageSender common.Address) []byte {
	enc, err := endpointV1.abi.Pack("configureEndpoint", contractFactory, participantStorageReplica, tokenRegistry, resourceManager, messageSender, messageReceiver, batchMessageSender)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackConfigureModules is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xd2eb6090.
//
// Solidity: function configureModules(address _resourceManager, address _messageSender, address _messageReceiver, address _batchMessageSender) returns()
func (endpointV1 *EndpointV1) PackConfigureModules(resourceManager common.Address, messageSender common.Address, messageReceiver common.Address, batchMessageSender common.Address) []byte {
	enc, err := endpointV1.abi.Pack("configureModules", resourceManager, messageSender, messageReceiver, batchMessageSender)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackContractVersion is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xa0a8e460.
//
// Solidity: function contractVersion() pure returns(uint256)
func (endpointV1 *EndpointV1) PackContractVersion() []byte {
	enc, err := endpointV1.abi.Pack("contractVersion")
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackContractVersion is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xa0a8e460.
//
// Solidity: function contractVersion() pure returns(uint256)
func (endpointV1 *EndpointV1) UnpackContractVersion(data []byte) (*big.Int, error) {
	out, err := endpointV1.abi.Unpack("contractVersion", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, err
}

// PackGetAddressByResourceId is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x11f50c85.
//
// Solidity: function getAddressByResourceId(bytes32 _resourceId) view returns(address)
func (endpointV1 *EndpointV1) PackGetAddressByResourceId(resourceId [32]byte) []byte {
	enc, err := endpointV1.abi.Pack("getAddressByResourceId", resourceId)
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackGetAddressByResourceId is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x11f50c85.
//
// Solidity: function getAddressByResourceId(bytes32 _resourceId) view returns(address)
func (endpointV1 *EndpointV1) UnpackGetAddressByResourceId(data []byte) (common.Address, error) {
	out, err := endpointV1.abi.Unpack("getAddressByResourceId", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, err
}

// PackGetChainId is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x3408e470.
//
// Solidity: function getChainId() view returns(uint256)
func (endpointV1 *EndpointV1) PackGetChainId() []byte {
	enc, err := endpointV1.abi.Pack("getChainId")
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackGetChainId is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x3408e470.
//
// Solidity: function getChainId() view returns(uint256)
func (endpointV1 *EndpointV1) UnpackGetChainId(data []byte) (*big.Int, error) {
	out, err := endpointV1.abi.Unpack("getChainId", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, err
}

// PackGetInboundNonce is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xde16f0b2.
//
// Solidity: function getInboundNonce(uint256 _srcChainId) view returns(uint256)
func (endpointV1 *EndpointV1) PackGetInboundNonce(srcChainId *big.Int) []byte {
	enc, err := endpointV1.abi.Pack("getInboundNonce", srcChainId)
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackGetInboundNonce is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xde16f0b2.
//
// Solidity: function getInboundNonce(uint256 _srcChainId) view returns(uint256)
func (endpointV1 *EndpointV1) UnpackGetInboundNonce(data []byte) (*big.Int, error) {
	out, err := endpointV1.abi.Unpack("getInboundNonce", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, err
}

// PackGetMaxBatchMessages is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xf1e22974.
//
// Solidity: function getMaxBatchMessages() view returns(uint256)
func (endpointV1 *EndpointV1) PackGetMaxBatchMessages() []byte {
	enc, err := endpointV1.abi.Pack("getMaxBatchMessages")
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackGetMaxBatchMessages is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xf1e22974.
//
// Solidity: function getMaxBatchMessages() view returns(uint256)
func (endpointV1 *EndpointV1) UnpackGetMaxBatchMessages(data []byte) (*big.Int, error) {
	out, err := endpointV1.abi.Unpack("getMaxBatchMessages", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, err
}

// PackGetOutboundNonce is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xd4c08529.
//
// Solidity: function getOutboundNonce(uint256 _dstChainId) view returns(uint256)
func (endpointV1 *EndpointV1) PackGetOutboundNonce(dstChainId *big.Int) []byte {
	enc, err := endpointV1.abi.Pack("getOutboundNonce", dstChainId)
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackGetOutboundNonce is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xd4c08529.
//
// Solidity: function getOutboundNonce(uint256 _dstChainId) view returns(uint256)
func (endpointV1 *EndpointV1) UnpackGetOutboundNonce(data []byte) (*big.Int, error) {
	out, err := endpointV1.abi.Unpack("getOutboundNonce", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, err
}

// PackGetPrivateHubAddress is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x6892f954.
//
// Solidity: function getPrivateHubAddress(string _contractName) view returns(address)
func (endpointV1 *EndpointV1) PackGetPrivateHubAddress(contractName string) []byte {
	enc, err := endpointV1.abi.Pack("getPrivateHubAddress", contractName)
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackGetPrivateHubAddress is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x6892f954.
//
// Solidity: function getPrivateHubAddress(string _contractName) view returns(address)
func (endpointV1 *EndpointV1) UnpackGetPrivateHubAddress(data []byte) (common.Address, error) {
	out, err := endpointV1.abi.Unpack("getPrivateHubAddress", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, err
}

// PackGetPrivateHubId is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x0b39a951.
//
// Solidity: function getPrivateHubId() view returns(uint256)
func (endpointV1 *EndpointV1) PackGetPrivateHubId() []byte {
	enc, err := endpointV1.abi.Pack("getPrivateHubId")
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackGetPrivateHubId is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x0b39a951.
//
// Solidity: function getPrivateHubId() view returns(uint256)
func (endpointV1 *EndpointV1) UnpackGetPrivateHubId(data []byte) (*big.Int, error) {
	out, err := endpointV1.abi.Unpack("getPrivateHubId", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, err
}

// PackGetTokenRegistry is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x057838bd.
//
// Solidity: function getTokenRegistry() view returns(address)
func (endpointV1 *EndpointV1) PackGetTokenRegistry() []byte {
	enc, err := endpointV1.abi.Pack("getTokenRegistry")
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackGetTokenRegistry is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x057838bd.
//
// Solidity: function getTokenRegistry() view returns(address)
func (endpointV1 *EndpointV1) UnpackGetTokenRegistry(data []byte) (common.Address, error) {
	out, err := endpointV1.abi.Unpack("getTokenRegistry", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, err
}

// PackGetUserGovernanceAddress is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xb8b87138.
//
// Solidity: function getUserGovernanceAddress() view returns(address)
func (endpointV1 *EndpointV1) PackGetUserGovernanceAddress() []byte {
	enc, err := endpointV1.abi.Pack("getUserGovernanceAddress")
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackGetUserGovernanceAddress is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xb8b87138.
//
// Solidity: function getUserGovernanceAddress() view returns(address)
func (endpointV1 *EndpointV1) UnpackGetUserGovernanceAddress(data []byte) (common.Address, error) {
	out, err := endpointV1.abi.Unpack("getUserGovernanceAddress", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, err
}

// PackInitialize is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x1f8d519d.
//
// Solidity: function initialize(uint256 _chainId, uint256 _privateHubId, uint256 _maxBatchMessages, address authority_) returns()
func (endpointV1 *EndpointV1) PackInitialize(chainId *big.Int, privateHubId *big.Int, maxBatchMessages *big.Int, authority common.Address) []byte {
	enc, err := endpointV1.abi.Pack("initialize", chainId, privateHubId, maxBatchMessages, authority)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackMessageReceiver is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xe2cae9f5.
//
// Solidity: function messageReceiver() view returns(address)
func (endpointV1 *EndpointV1) PackMessageReceiver() []byte {
	enc, err := endpointV1.abi.Pack("messageReceiver")
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackMessageReceiver is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xe2cae9f5.
//
// Solidity: function messageReceiver() view returns(address)
func (endpointV1 *EndpointV1) UnpackMessageReceiver(data []byte) (common.Address, error) {
	out, err := endpointV1.abi.Unpack("messageReceiver", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, err
}

// PackMessageSender is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xd67bdd25.
//
// Solidity: function messageSender() view returns(address)
func (endpointV1 *EndpointV1) PackMessageSender() []byte {
	enc, err := endpointV1.abi.Pack("messageSender")
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackMessageSender is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xd67bdd25.
//
// Solidity: function messageSender() view returns(address)
func (endpointV1 *EndpointV1) UnpackMessageSender(data []byte) (common.Address, error) {
	out, err := endpointV1.abi.Unpack("messageSender", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, err
}

// PackPrivateHubAddress is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x16e72d5f.
//
// Solidity: function privateHubAddress(string ) view returns(address)
func (endpointV1 *EndpointV1) PackPrivateHubAddress(arg0 string) []byte {
	enc, err := endpointV1.abi.Pack("privateHubAddress", arg0)
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackPrivateHubAddress is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x16e72d5f.
//
// Solidity: function privateHubAddress(string ) view returns(address)
func (endpointV1 *EndpointV1) UnpackPrivateHubAddress(data []byte) (common.Address, error) {
	out, err := endpointV1.abi.Unpack("privateHubAddress", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, err
}

// PackPrivateHubId is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x202e2374.
//
// Solidity: function privateHubId() view returns(uint256)
func (endpointV1 *EndpointV1) PackPrivateHubId() []byte {
	enc, err := endpointV1.abi.Pack("privateHubId")
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackPrivateHubId is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x202e2374.
//
// Solidity: function privateHubId() view returns(uint256)
func (endpointV1 *EndpointV1) UnpackPrivateHubId(data []byte) (*big.Int, error) {
	out, err := endpointV1.abi.Unpack("privateHubId", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, err
}

// PackProxiableUUID is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (endpointV1 *EndpointV1) PackProxiableUUID() []byte {
	enc, err := endpointV1.abi.Pack("proxiableUUID")
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackProxiableUUID is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (endpointV1 *EndpointV1) UnpackProxiableUUID(data []byte) ([32]byte, error) {
	out, err := endpointV1.abi.Unpack("proxiableUUID", data)
	if err != nil {
		return *new([32]byte), err
	}
	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)
	return out0, err
}

// PackReceivePayload is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xfe12d9c3.
//
// Solidity: function receivePayload(uint256 _srcChainId, address _srcAddress, address _dstAddress, ((bool,uint256,(bool,uint8,bytes,uint8,bytes),bytes32,bytes,bytes,bytes,(uint8,uint256,address,address,address,uint256),bool),bytes) _raylsMessage, bytes32 _messageId) returns()
func (endpointV1 *EndpointV1) PackReceivePayload(srcChainId *big.Int, srcAddress common.Address, dstAddress common.Address, raylsMessage RaylsMessage, messageId [32]byte) []byte {
	enc, err := endpointV1.abi.Pack("receivePayload", srcChainId, srcAddress, dstAddress, raylsMessage, messageId)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackRegisterPrivateHubAddress is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x55615960.
//
// Solidity: function registerPrivateHubAddress(string _contractName, address _contractAddressOnPrivateHub) returns()
func (endpointV1 *EndpointV1) PackRegisterPrivateHubAddress(contractName string, contractAddressOnPrivateHub common.Address) []byte {
	enc, err := endpointV1.abi.Pack("registerPrivateHubAddress", contractName, contractAddressOnPrivateHub)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackRegisterResourceId is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xd4f951c7.
//
// Solidity: function registerResourceId(bytes32 _resourceId, address _implementationAddress) returns()
func (endpointV1 *EndpointV1) PackRegisterResourceId(resourceId [32]byte, implementationAddress common.Address) []byte {
	enc, err := endpointV1.abi.Pack("registerResourceId", resourceId, implementationAddress)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackRequestNewRaylsViewKeys is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x26b32600.
//
// Solidity: function requestNewRaylsViewKeys(uint256 blockNumber) returns()
func (endpointV1 *EndpointV1) PackRequestNewRaylsViewKeys(blockNumber *big.Int) []byte {
	enc, err := endpointV1.abi.Pack("requestNewRaylsViewKeys", blockNumber)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackResourceManager is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x71f10b8b.
//
// Solidity: function resourceManager() view returns(address)
func (endpointV1 *EndpointV1) PackResourceManager() []byte {
	enc, err := endpointV1.abi.Pack("resourceManager")
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackResourceManager is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x71f10b8b.
//
// Solidity: function resourceManager() view returns(address)
func (endpointV1 *EndpointV1) UnpackResourceManager(data []byte) (common.Address, error) {
	out, err := endpointV1.abi.Unpack("resourceManager", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, err
}

// PackSend is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x150b375f.
//
// Solidity: function send(uint256 _dstChainId, address _destination, bytes _payload) payable returns(bytes32 messageId)
func (endpointV1 *EndpointV1) PackSend(dstChainId *big.Int, destination common.Address, payload []byte) []byte {
	enc, err := endpointV1.abi.Pack("send", dstChainId, destination, payload)
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackSend is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x150b375f.
//
// Solidity: function send(uint256 _dstChainId, address _destination, bytes _payload) payable returns(bytes32 messageId)
func (endpointV1 *EndpointV1) UnpackSend(data []byte) ([32]byte, error) {
	out, err := endpointV1.abi.Unpack("send", data)
	if err != nil {
		return *new([32]byte), err
	}
	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)
	return out0, err
}

// PackSend0 is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xfced8ee8.
//
// Solidity: function send(uint256 _dstChainId, address _destination, bytes _payload, (uint8,uint256,address,address,address,uint256) transferMetadata) payable returns(bytes32 messageId)
func (endpointV1 *EndpointV1) PackSend0(dstChainId *big.Int, destination common.Address, payload []byte, transferMetadata BridgedTransferMetadata) []byte {
	enc, err := endpointV1.abi.Pack("send0", dstChainId, destination, payload, transferMetadata)
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackSend0 is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xfced8ee8.
//
// Solidity: function send(uint256 _dstChainId, address _destination, bytes _payload, (uint8,uint256,address,address,address,uint256) transferMetadata) payable returns(bytes32 messageId)
func (endpointV1 *EndpointV1) UnpackSend0(data []byte) ([32]byte, error) {
	out, err := endpointV1.abi.Unpack("send0", data)
	if err != nil {
		return *new([32]byte), err
	}
	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)
	return out0, err
}

// PackSendBatch is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x407fa89e.
//
// Solidity: function sendBatch((uint256,address,bytes)[] _destinationPayloadRequests) returns(bytes32 batchId)
func (endpointV1 *EndpointV1) PackSendBatch(destinationPayloadRequests []DestinationPayloadRequest) []byte {
	enc, err := endpointV1.abi.Pack("sendBatch", destinationPayloadRequests)
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackSendBatch is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x407fa89e.
//
// Solidity: function sendBatch((uint256,address,bytes)[] _destinationPayloadRequests) returns(bytes32 batchId)
func (endpointV1 *EndpointV1) UnpackSendBatch(data []byte) ([32]byte, error) {
	out, err := endpointV1.abi.Unpack("sendBatch", data)
	if err != nil {
		return *new([32]byte), err
	}
	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)
	return out0, err
}

// PackSendBatchToResourceId is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x4f02198c.
//
// Solidity: function sendBatchToResourceId((uint256,bytes32,bytes,bytes,bytes,bytes,(uint8,uint256,address,address,address,uint256))[] _resourceIdPayloadRequests) payable returns(bytes32 batchId)
func (endpointV1 *EndpointV1) PackSendBatchToResourceId(resourceIdPayloadRequests []ResourceIdCompletePayloadRequest) []byte {
	enc, err := endpointV1.abi.Pack("sendBatchToResourceId", resourceIdPayloadRequests)
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackSendBatchToResourceId is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x4f02198c.
//
// Solidity: function sendBatchToResourceId((uint256,bytes32,bytes,bytes,bytes,bytes,(uint8,uint256,address,address,address,uint256))[] _resourceIdPayloadRequests) payable returns(bytes32 batchId)
func (endpointV1 *EndpointV1) UnpackSendBatchToResourceId(data []byte) ([32]byte, error) {
	out, err := endpointV1.abi.Unpack("sendBatchToResourceId", data)
	if err != nil {
		return *new([32]byte), err
	}
	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)
	return out0, err
}

// PackSendBatchToResourceId0 is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xb2958e9a.
//
// Solidity: function sendBatchToResourceId((uint256,bytes32,bytes)[] _resourceIdPayloadRequests) payable returns(bytes32 batchId)
func (endpointV1 *EndpointV1) PackSendBatchToResourceId0(resourceIdPayloadRequests []ResourceIdPayloadRequest) []byte {
	enc, err := endpointV1.abi.Pack("sendBatchToResourceId0", resourceIdPayloadRequests)
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackSendBatchToResourceId0 is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xb2958e9a.
//
// Solidity: function sendBatchToResourceId((uint256,bytes32,bytes)[] _resourceIdPayloadRequests) payable returns(bytes32 batchId)
func (endpointV1 *EndpointV1) UnpackSendBatchToResourceId0(data []byte) ([32]byte, error) {
	out, err := endpointV1.abi.Unpack("sendBatchToResourceId0", data)
	if err != nil {
		return *new([32]byte), err
	}
	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)
	return out0, err
}

// PackSendToResourceId is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x40b33bfa.
//
// Solidity: function sendToResourceId(uint256 _dstChainId, bytes32 _resourceId, bytes _payload) payable returns(bytes32 messageId)
func (endpointV1 *EndpointV1) PackSendToResourceId(dstChainId *big.Int, resourceId [32]byte, payload []byte) []byte {
	enc, err := endpointV1.abi.Pack("sendToResourceId", dstChainId, resourceId, payload)
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackSendToResourceId is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x40b33bfa.
//
// Solidity: function sendToResourceId(uint256 _dstChainId, bytes32 _resourceId, bytes _payload) payable returns(bytes32 messageId)
func (endpointV1 *EndpointV1) UnpackSendToResourceId(data []byte) ([32]byte, error) {
	out, err := endpointV1.abi.Unpack("sendToResourceId", data)
	if err != nil {
		return *new([32]byte), err
	}
	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)
	return out0, err
}

// PackSendToResourceId0 is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x41d71744.
//
// Solidity: function sendToResourceId(uint256 _dstChainId, bytes32 _resourceId, bytes _payload, bytes _lockData, bytes _revertDataSender, bytes _revertDataReceiver, (uint8,uint256,address,address,address,uint256) transferMetadata) payable returns(bytes32 messageId)
func (endpointV1 *EndpointV1) PackSendToResourceId0(dstChainId *big.Int, resourceId [32]byte, payload []byte, lockData []byte, revertDataSender []byte, revertDataReceiver []byte, transferMetadata BridgedTransferMetadata) []byte {
	enc, err := endpointV1.abi.Pack("sendToResourceId0", dstChainId, resourceId, payload, lockData, revertDataSender, revertDataReceiver, transferMetadata)
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackSendToResourceId0 is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x41d71744.
//
// Solidity: function sendToResourceId(uint256 _dstChainId, bytes32 _resourceId, bytes _payload, bytes _lockData, bytes _revertDataSender, bytes _revertDataReceiver, (uint8,uint256,address,address,address,uint256) transferMetadata) payable returns(bytes32 messageId)
func (endpointV1 *EndpointV1) UnpackSendToResourceId0(data []byte) ([32]byte, error) {
	out, err := endpointV1.abi.Unpack("sendToResourceId0", data)
	if err != nil {
		return *new([32]byte), err
	}
	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)
	return out0, err
}

// PackSetMaxBatchMessages is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x910e2347.
//
// Solidity: function setMaxBatchMessages(uint256 _maxBatchMessages) returns()
func (endpointV1 *EndpointV1) PackSetMaxBatchMessages(maxBatchMessages *big.Int) []byte {
	enc, err := endpointV1.abi.Pack("setMaxBatchMessages", maxBatchMessages)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackSetUserGovernance is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x1a8d9216.
//
// Solidity: function setUserGovernance(address _userGovernance) returns()
func (endpointV1 *EndpointV1) PackSetUserGovernance(userGovernance common.Address) []byte {
	enc, err := endpointV1.abi.Pack("setUserGovernance", userGovernance)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackUpgradeToAndCall is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (endpointV1 *EndpointV1) PackUpgradeToAndCall(newImplementation common.Address, data []byte) []byte {
	enc, err := endpointV1.abi.Pack("upgradeToAndCall", newImplementation, data)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackVersion is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x54fd4d50.
//
// Solidity: function version() pure returns(string)
func (endpointV1 *EndpointV1) PackVersion() []byte {
	enc, err := endpointV1.abi.Pack("version")
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackVersion is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x54fd4d50.
//
// Solidity: function version() pure returns(string)
func (endpointV1 *EndpointV1) UnpackVersion(data []byte) (string, error) {
	out, err := endpointV1.abi.Unpack("version", data)
	if err != nil {
		return *new(string), err
	}
	out0 := *abi.ConvertType(out[0], new(string)).(*string)
	return out0, err
}

// EndpointV1AuthorityUpdated represents a AuthorityUpdated event raised by the EndpointV1 contract.
type EndpointV1AuthorityUpdated struct {
	OldAuthority common.Address
	NewAuthority common.Address
	Raw          *types.Log // Blockchain specific contextual infos
}

const EndpointV1AuthorityUpdatedEventName = "AuthorityUpdated"

// ContractEventName returns the user-defined event name.
func (EndpointV1AuthorityUpdated) ContractEventName() string {
	return EndpointV1AuthorityUpdatedEventName
}

// UnpackAuthorityUpdatedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event AuthorityUpdated(address indexed oldAuthority, address indexed newAuthority)
func (endpointV1 *EndpointV1) UnpackAuthorityUpdatedEvent(log *types.Log) (*EndpointV1AuthorityUpdated, error) {
	event := "AuthorityUpdated"
	if log.Topics[0] != endpointV1.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(EndpointV1AuthorityUpdated)
	if len(log.Data) > 0 {
		if err := endpointV1.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range endpointV1.abi.Events[event].Inputs {
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

// EndpointV1EndpointConfigured represents a EndpointConfigured event raised by the EndpointV1 contract.
type EndpointV1EndpointConfigured struct {
	ContractFactory           common.Address
	ParticipantStorageReplica common.Address
	TokenRegistry             common.Address
	ResourceManager           common.Address
	MessageSender             common.Address
	MessageReceiver           common.Address
	BatchMessageSender        common.Address
	Raw                       *types.Log // Blockchain specific contextual infos
}

const EndpointV1EndpointConfiguredEventName = "EndpointConfigured"

// ContractEventName returns the user-defined event name.
func (EndpointV1EndpointConfigured) ContractEventName() string {
	return EndpointV1EndpointConfiguredEventName
}

// UnpackEndpointConfiguredEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event EndpointConfigured(address contractFactory, address participantStorageReplica, address tokenRegistry, address resourceManager, address messageSender, address messageReceiver, address batchMessageSender)
func (endpointV1 *EndpointV1) UnpackEndpointConfiguredEvent(log *types.Log) (*EndpointV1EndpointConfigured, error) {
	event := "EndpointConfigured"
	if log.Topics[0] != endpointV1.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(EndpointV1EndpointConfigured)
	if len(log.Data) > 0 {
		if err := endpointV1.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range endpointV1.abi.Events[event].Inputs {
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

// EndpointV1Initialized represents a Initialized event raised by the EndpointV1 contract.
type EndpointV1Initialized struct {
	Version uint64
	Raw     *types.Log // Blockchain specific contextual infos
}

const EndpointV1InitializedEventName = "Initialized"

// ContractEventName returns the user-defined event name.
func (EndpointV1Initialized) ContractEventName() string {
	return EndpointV1InitializedEventName
}

// UnpackInitializedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event Initialized(uint64 version)
func (endpointV1 *EndpointV1) UnpackInitializedEvent(log *types.Log) (*EndpointV1Initialized, error) {
	event := "Initialized"
	if log.Topics[0] != endpointV1.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(EndpointV1Initialized)
	if len(log.Data) > 0 {
		if err := endpointV1.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range endpointV1.abi.Events[event].Inputs {
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

// EndpointV1MessageBatchDispatched represents a MessageBatchDispatched event raised by the EndpointV1 contract.
type EndpointV1MessageBatchDispatched struct {
	BatchId  [32]byte
	From     common.Address
	Messages []BatchMessage
	Raw      *types.Log // Blockchain specific contextual infos
}

const EndpointV1MessageBatchDispatchedEventName = "MessageBatchDispatched"

// ContractEventName returns the user-defined event name.
func (EndpointV1MessageBatchDispatched) ContractEventName() string {
	return EndpointV1MessageBatchDispatchedEventName
}

// UnpackMessageBatchDispatchedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event MessageBatchDispatched(bytes32 batchId, address from, (uint256,address,((bool,uint256,(bool,uint8,bytes,uint8,bytes),bytes32,bytes,bytes,bytes,(uint8,uint256,address,address,address,uint256),bool),bytes),bytes32)[] messages)
func (endpointV1 *EndpointV1) UnpackMessageBatchDispatchedEvent(log *types.Log) (*EndpointV1MessageBatchDispatched, error) {
	event := "MessageBatchDispatched"
	if log.Topics[0] != endpointV1.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(EndpointV1MessageBatchDispatched)
	if len(log.Data) > 0 {
		if err := endpointV1.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range endpointV1.abi.Events[event].Inputs {
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

// EndpointV1MessageDispatched represents a MessageDispatched event raised by the EndpointV1 contract.
type EndpointV1MessageDispatched struct {
	MessageId [32]byte
	From      common.Address
	ToChainId *big.Int
	To        common.Address
	Data      RaylsMessage
	Raw       *types.Log // Blockchain specific contextual infos
}

const EndpointV1MessageDispatchedEventName = "MessageDispatched"

// ContractEventName returns the user-defined event name.
func (EndpointV1MessageDispatched) ContractEventName() string {
	return EndpointV1MessageDispatchedEventName
}

// UnpackMessageDispatchedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event MessageDispatched(bytes32 indexed messageId, address indexed from, uint256 indexed toChainId, address to, ((bool,uint256,(bool,uint8,bytes,uint8,bytes),bytes32,bytes,bytes,bytes,(uint8,uint256,address,address,address,uint256),bool),bytes) data)
func (endpointV1 *EndpointV1) UnpackMessageDispatchedEvent(log *types.Log) (*EndpointV1MessageDispatched, error) {
	event := "MessageDispatched"
	if log.Topics[0] != endpointV1.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(EndpointV1MessageDispatched)
	if len(log.Data) > 0 {
		if err := endpointV1.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range endpointV1.abi.Events[event].Inputs {
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

// EndpointV1ModulesConfigured represents a ModulesConfigured event raised by the EndpointV1 contract.
type EndpointV1ModulesConfigured struct {
	ResourceManager    common.Address
	MessageSender      common.Address
	MessageReceiver    common.Address
	BatchMessageSender common.Address
	Raw                *types.Log // Blockchain specific contextual infos
}

const EndpointV1ModulesConfiguredEventName = "ModulesConfigured"

// ContractEventName returns the user-defined event name.
func (EndpointV1ModulesConfigured) ContractEventName() string {
	return EndpointV1ModulesConfiguredEventName
}

// UnpackModulesConfiguredEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event ModulesConfigured(address resourceManager, address messageSender, address messageReceiver, address batchMessageSender)
func (endpointV1 *EndpointV1) UnpackModulesConfiguredEvent(log *types.Log) (*EndpointV1ModulesConfigured, error) {
	event := "ModulesConfigured"
	if log.Topics[0] != endpointV1.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(EndpointV1ModulesConfigured)
	if len(log.Data) > 0 {
		if err := endpointV1.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range endpointV1.abi.Events[event].Inputs {
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

// EndpointV1PrivateHubAddressRegistered represents a PrivateHubAddressRegistered event raised by the EndpointV1 contract.
type EndpointV1PrivateHubAddressRegistered struct {
	ContractName    common.Hash
	ContractAddress common.Address
	Raw             *types.Log // Blockchain specific contextual infos
}

const EndpointV1PrivateHubAddressRegisteredEventName = "PrivateHubAddressRegistered"

// ContractEventName returns the user-defined event name.
func (EndpointV1PrivateHubAddressRegistered) ContractEventName() string {
	return EndpointV1PrivateHubAddressRegisteredEventName
}

// UnpackPrivateHubAddressRegisteredEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event PrivateHubAddressRegistered(string indexed contractName, address indexed contractAddress)
func (endpointV1 *EndpointV1) UnpackPrivateHubAddressRegisteredEvent(log *types.Log) (*EndpointV1PrivateHubAddressRegistered, error) {
	event := "PrivateHubAddressRegistered"
	if log.Topics[0] != endpointV1.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(EndpointV1PrivateHubAddressRegistered)
	if len(log.Data) > 0 {
		if err := endpointV1.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range endpointV1.abi.Events[event].Inputs {
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

// EndpointV1ResourceIdRegistered represents a ResourceIdRegistered event raised by the EndpointV1 contract.
type EndpointV1ResourceIdRegistered struct {
	ResourceId            [32]byte
	ImplementationAddress common.Address
	Raw                   *types.Log // Blockchain specific contextual infos
}

const EndpointV1ResourceIdRegisteredEventName = "ResourceIdRegistered"

// ContractEventName returns the user-defined event name.
func (EndpointV1ResourceIdRegistered) ContractEventName() string {
	return EndpointV1ResourceIdRegisteredEventName
}

// UnpackResourceIdRegisteredEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event ResourceIdRegistered(bytes32 indexed resourceId, address indexed implementationAddress)
func (endpointV1 *EndpointV1) UnpackResourceIdRegisteredEvent(log *types.Log) (*EndpointV1ResourceIdRegistered, error) {
	event := "ResourceIdRegistered"
	if log.Topics[0] != endpointV1.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(EndpointV1ResourceIdRegistered)
	if len(log.Data) > 0 {
		if err := endpointV1.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range endpointV1.abi.Events[event].Inputs {
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

// EndpointV1UpdateRaylsViewKeysRequest represents a UpdateRaylsViewKeysRequest event raised by the EndpointV1 contract.
type EndpointV1UpdateRaylsViewKeysRequest struct {
	BlockNumber *big.Int
	Raw         *types.Log // Blockchain specific contextual infos
}

const EndpointV1UpdateRaylsViewKeysRequestEventName = "UpdateRaylsViewKeysRequest"

// ContractEventName returns the user-defined event name.
func (EndpointV1UpdateRaylsViewKeysRequest) ContractEventName() string {
	return EndpointV1UpdateRaylsViewKeysRequestEventName
}

// UnpackUpdateRaylsViewKeysRequestEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event UpdateRaylsViewKeysRequest(uint256 blockNumber)
func (endpointV1 *EndpointV1) UnpackUpdateRaylsViewKeysRequestEvent(log *types.Log) (*EndpointV1UpdateRaylsViewKeysRequest, error) {
	event := "UpdateRaylsViewKeysRequest"
	if log.Topics[0] != endpointV1.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(EndpointV1UpdateRaylsViewKeysRequest)
	if len(log.Data) > 0 {
		if err := endpointV1.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range endpointV1.abi.Events[event].Inputs {
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

// EndpointV1Upgraded represents a Upgraded event raised by the EndpointV1 contract.
type EndpointV1Upgraded struct {
	Implementation common.Address
	Raw            *types.Log // Blockchain specific contextual infos
}

const EndpointV1UpgradedEventName = "Upgraded"

// ContractEventName returns the user-defined event name.
func (EndpointV1Upgraded) ContractEventName() string {
	return EndpointV1UpgradedEventName
}

// UnpackUpgradedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event Upgraded(address indexed implementation)
func (endpointV1 *EndpointV1) UnpackUpgradedEvent(log *types.Log) (*EndpointV1Upgraded, error) {
	event := "Upgraded"
	if log.Topics[0] != endpointV1.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(EndpointV1Upgraded)
	if len(log.Data) > 0 {
		if err := endpointV1.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range endpointV1.abi.Events[event].Inputs {
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
func (endpointV1 *EndpointV1) UnpackError(raw []byte) (any, error) {
	if bytes.Equal(raw[:4], endpointV1.abi.Errors["AddressEmptyCode"].ID.Bytes()[:4]) {
		return endpointV1.UnpackAddressEmptyCodeError(raw[4:])
	}
	if bytes.Equal(raw[:4], endpointV1.abi.Errors["ERC1967InvalidImplementation"].ID.Bytes()[:4]) {
		return endpointV1.UnpackERC1967InvalidImplementationError(raw[4:])
	}
	if bytes.Equal(raw[:4], endpointV1.abi.Errors["ERC1967NonPayable"].ID.Bytes()[:4]) {
		return endpointV1.UnpackERC1967NonPayableError(raw[4:])
	}
	if bytes.Equal(raw[:4], endpointV1.abi.Errors["EndpointEmptyContractName"].ID.Bytes()[:4]) {
		return endpointV1.UnpackEndpointEmptyContractNameError(raw[4:])
	}
	if bytes.Equal(raw[:4], endpointV1.abi.Errors["EndpointInvalidPrivateHubAddress"].ID.Bytes()[:4]) {
		return endpointV1.UnpackEndpointInvalidPrivateHubAddressError(raw[4:])
	}
	if bytes.Equal(raw[:4], endpointV1.abi.Errors["FailedCall"].ID.Bytes()[:4]) {
		return endpointV1.UnpackFailedCallError(raw[4:])
	}
	if bytes.Equal(raw[:4], endpointV1.abi.Errors["InvalidInitialization"].ID.Bytes()[:4]) {
		return endpointV1.UnpackInvalidInitializationError(raw[4:])
	}
	if bytes.Equal(raw[:4], endpointV1.abi.Errors["NotInitializing"].ID.Bytes()[:4]) {
		return endpointV1.UnpackNotInitializingError(raw[4:])
	}
	if bytes.Equal(raw[:4], endpointV1.abi.Errors["RaylsAccessManagedContractPaused"].ID.Bytes()[:4]) {
		return endpointV1.UnpackRaylsAccessManagedContractPausedError(raw[4:])
	}
	if bytes.Equal(raw[:4], endpointV1.abi.Errors["RaylsAccessManagedInvalidAuthority"].ID.Bytes()[:4]) {
		return endpointV1.UnpackRaylsAccessManagedInvalidAuthorityError(raw[4:])
	}
	if bytes.Equal(raw[:4], endpointV1.abi.Errors["RaylsAccessManagedMustSchedule"].ID.Bytes()[:4]) {
		return endpointV1.UnpackRaylsAccessManagedMustScheduleError(raw[4:])
	}
	if bytes.Equal(raw[:4], endpointV1.abi.Errors["RaylsAccessManagedUnauthorized"].ID.Bytes()[:4]) {
		return endpointV1.UnpackRaylsAccessManagedUnauthorizedError(raw[4:])
	}
	if bytes.Equal(raw[:4], endpointV1.abi.Errors["UUPSUnauthorizedCallContext"].ID.Bytes()[:4]) {
		return endpointV1.UnpackUUPSUnauthorizedCallContextError(raw[4:])
	}
	if bytes.Equal(raw[:4], endpointV1.abi.Errors["UUPSUnsupportedProxiableUUID"].ID.Bytes()[:4]) {
		return endpointV1.UnpackUUPSUnsupportedProxiableUUIDError(raw[4:])
	}
	return nil, errors.New("Unknown error")
}

// EndpointV1AddressEmptyCode represents a AddressEmptyCode error raised by the EndpointV1 contract.
type EndpointV1AddressEmptyCode struct {
	Target common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error AddressEmptyCode(address target)
func EndpointV1AddressEmptyCodeErrorID() common.Hash {
	return common.HexToHash("0x9996b315c842ff135b8fc4a08ad5df1c344efbc03d2687aecc0678050d2aac89")
}

// UnpackAddressEmptyCodeError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error AddressEmptyCode(address target)
func (endpointV1 *EndpointV1) UnpackAddressEmptyCodeError(raw []byte) (*EndpointV1AddressEmptyCode, error) {
	out := new(EndpointV1AddressEmptyCode)
	if err := endpointV1.abi.UnpackIntoInterface(out, "AddressEmptyCode", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// EndpointV1ERC1967InvalidImplementation represents a ERC1967InvalidImplementation error raised by the EndpointV1 contract.
type EndpointV1ERC1967InvalidImplementation struct {
	Implementation common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error ERC1967InvalidImplementation(address implementation)
func EndpointV1ERC1967InvalidImplementationErrorID() common.Hash {
	return common.HexToHash("0x4c9c8ce3ceb3130f17f7cdba48d89b5b0129f266a8bac114e6e315a41879b617")
}

// UnpackERC1967InvalidImplementationError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error ERC1967InvalidImplementation(address implementation)
func (endpointV1 *EndpointV1) UnpackERC1967InvalidImplementationError(raw []byte) (*EndpointV1ERC1967InvalidImplementation, error) {
	out := new(EndpointV1ERC1967InvalidImplementation)
	if err := endpointV1.abi.UnpackIntoInterface(out, "ERC1967InvalidImplementation", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// EndpointV1ERC1967NonPayable represents a ERC1967NonPayable error raised by the EndpointV1 contract.
type EndpointV1ERC1967NonPayable struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error ERC1967NonPayable()
func EndpointV1ERC1967NonPayableErrorID() common.Hash {
	return common.HexToHash("0xb398979fa84f543c8e222f17890372c487baf85e062276c127fef521eea7224b")
}

// UnpackERC1967NonPayableError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error ERC1967NonPayable()
func (endpointV1 *EndpointV1) UnpackERC1967NonPayableError(raw []byte) (*EndpointV1ERC1967NonPayable, error) {
	out := new(EndpointV1ERC1967NonPayable)
	if err := endpointV1.abi.UnpackIntoInterface(out, "ERC1967NonPayable", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// EndpointV1EndpointEmptyContractName represents a Endpoint__EmptyContractName error raised by the EndpointV1 contract.
type EndpointV1EndpointEmptyContractName struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error Endpoint__EmptyContractName()
func EndpointV1EndpointEmptyContractNameErrorID() common.Hash {
	return common.HexToHash("0x963611aec76eb385db81575a316d9b7021418f010be553eee787ef8489c6cb7b")
}

// UnpackEndpointEmptyContractNameError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error Endpoint__EmptyContractName()
func (endpointV1 *EndpointV1) UnpackEndpointEmptyContractNameError(raw []byte) (*EndpointV1EndpointEmptyContractName, error) {
	out := new(EndpointV1EndpointEmptyContractName)
	if err := endpointV1.abi.UnpackIntoInterface(out, "EndpointEmptyContractName", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// EndpointV1EndpointInvalidPrivateHubAddress represents a Endpoint__InvalidPrivateHubAddress error raised by the EndpointV1 contract.
type EndpointV1EndpointInvalidPrivateHubAddress struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error Endpoint__InvalidPrivateHubAddress()
func EndpointV1EndpointInvalidPrivateHubAddressErrorID() common.Hash {
	return common.HexToHash("0x9e0f82d90ea7e60e65e16ea4883d6ab53b4a5373e62ea653802542f70fc4c97b")
}

// UnpackEndpointInvalidPrivateHubAddressError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error Endpoint__InvalidPrivateHubAddress()
func (endpointV1 *EndpointV1) UnpackEndpointInvalidPrivateHubAddressError(raw []byte) (*EndpointV1EndpointInvalidPrivateHubAddress, error) {
	out := new(EndpointV1EndpointInvalidPrivateHubAddress)
	if err := endpointV1.abi.UnpackIntoInterface(out, "EndpointInvalidPrivateHubAddress", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// EndpointV1FailedCall represents a FailedCall error raised by the EndpointV1 contract.
type EndpointV1FailedCall struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error FailedCall()
func EndpointV1FailedCallErrorID() common.Hash {
	return common.HexToHash("0xd6bda27508c0fb6d8a39b4b122878dab26f731a7d4e4abe711dd3731899052a4")
}

// UnpackFailedCallError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error FailedCall()
func (endpointV1 *EndpointV1) UnpackFailedCallError(raw []byte) (*EndpointV1FailedCall, error) {
	out := new(EndpointV1FailedCall)
	if err := endpointV1.abi.UnpackIntoInterface(out, "FailedCall", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// EndpointV1InvalidInitialization represents a InvalidInitialization error raised by the EndpointV1 contract.
type EndpointV1InvalidInitialization struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error InvalidInitialization()
func EndpointV1InvalidInitializationErrorID() common.Hash {
	return common.HexToHash("0xf92ee8a957075833165f68c320933b1a1294aafc84ee6e0dd3fb178008f9aaf5")
}

// UnpackInvalidInitializationError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error InvalidInitialization()
func (endpointV1 *EndpointV1) UnpackInvalidInitializationError(raw []byte) (*EndpointV1InvalidInitialization, error) {
	out := new(EndpointV1InvalidInitialization)
	if err := endpointV1.abi.UnpackIntoInterface(out, "InvalidInitialization", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// EndpointV1NotInitializing represents a NotInitializing error raised by the EndpointV1 contract.
type EndpointV1NotInitializing struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error NotInitializing()
func EndpointV1NotInitializingErrorID() common.Hash {
	return common.HexToHash("0xd7e6bcf8597daa127dc9f0048d2f08d5ef140a2cb659feabd700beff1f7a8302")
}

// UnpackNotInitializingError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error NotInitializing()
func (endpointV1 *EndpointV1) UnpackNotInitializingError(raw []byte) (*EndpointV1NotInitializing, error) {
	out := new(EndpointV1NotInitializing)
	if err := endpointV1.abi.UnpackIntoInterface(out, "NotInitializing", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// EndpointV1RaylsAccessManagedContractPaused represents a RaylsAccessManaged__ContractPaused error raised by the EndpointV1 contract.
type EndpointV1RaylsAccessManagedContractPaused struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error RaylsAccessManaged__ContractPaused()
func EndpointV1RaylsAccessManagedContractPausedErrorID() common.Hash {
	return common.HexToHash("0xcc9855adb54c3eae80b5cb29441d3eb81c3fcdc928abeaba0ac7979b3072c3ab")
}

// UnpackRaylsAccessManagedContractPausedError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error RaylsAccessManaged__ContractPaused()
func (endpointV1 *EndpointV1) UnpackRaylsAccessManagedContractPausedError(raw []byte) (*EndpointV1RaylsAccessManagedContractPaused, error) {
	out := new(EndpointV1RaylsAccessManagedContractPaused)
	if err := endpointV1.abi.UnpackIntoInterface(out, "RaylsAccessManagedContractPaused", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// EndpointV1RaylsAccessManagedInvalidAuthority represents a RaylsAccessManaged__InvalidAuthority error raised by the EndpointV1 contract.
type EndpointV1RaylsAccessManagedInvalidAuthority struct {
	Authority common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error RaylsAccessManaged__InvalidAuthority(address authority)
func EndpointV1RaylsAccessManagedInvalidAuthorityErrorID() common.Hash {
	return common.HexToHash("0x894403477abbba3df1c9b757907cda2d4b4d9ba52242d5e6bdae2dc6a343662b")
}

// UnpackRaylsAccessManagedInvalidAuthorityError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error RaylsAccessManaged__InvalidAuthority(address authority)
func (endpointV1 *EndpointV1) UnpackRaylsAccessManagedInvalidAuthorityError(raw []byte) (*EndpointV1RaylsAccessManagedInvalidAuthority, error) {
	out := new(EndpointV1RaylsAccessManagedInvalidAuthority)
	if err := endpointV1.abi.UnpackIntoInterface(out, "RaylsAccessManagedInvalidAuthority", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// EndpointV1RaylsAccessManagedMustSchedule represents a RaylsAccessManaged__MustSchedule error raised by the EndpointV1 contract.
type EndpointV1RaylsAccessManagedMustSchedule struct {
	Caller common.Address
	Delay  uint32
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error RaylsAccessManaged__MustSchedule(address caller, uint32 delay)
func EndpointV1RaylsAccessManagedMustScheduleErrorID() common.Hash {
	return common.HexToHash("0xa42687899abecbc3886d9231154d03927d3e61df458dda1e50129311f272684f")
}

// UnpackRaylsAccessManagedMustScheduleError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error RaylsAccessManaged__MustSchedule(address caller, uint32 delay)
func (endpointV1 *EndpointV1) UnpackRaylsAccessManagedMustScheduleError(raw []byte) (*EndpointV1RaylsAccessManagedMustSchedule, error) {
	out := new(EndpointV1RaylsAccessManagedMustSchedule)
	if err := endpointV1.abi.UnpackIntoInterface(out, "RaylsAccessManagedMustSchedule", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// EndpointV1RaylsAccessManagedUnauthorized represents a RaylsAccessManaged__Unauthorized error raised by the EndpointV1 contract.
type EndpointV1RaylsAccessManagedUnauthorized struct {
	Caller common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error RaylsAccessManaged__Unauthorized(address caller)
func EndpointV1RaylsAccessManagedUnauthorizedErrorID() common.Hash {
	return common.HexToHash("0xbb34f40c15462e9dd86b312158515c7ee80e2406cfa90a92a6fb6fad1bc48a48")
}

// UnpackRaylsAccessManagedUnauthorizedError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error RaylsAccessManaged__Unauthorized(address caller)
func (endpointV1 *EndpointV1) UnpackRaylsAccessManagedUnauthorizedError(raw []byte) (*EndpointV1RaylsAccessManagedUnauthorized, error) {
	out := new(EndpointV1RaylsAccessManagedUnauthorized)
	if err := endpointV1.abi.UnpackIntoInterface(out, "RaylsAccessManagedUnauthorized", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// EndpointV1UUPSUnauthorizedCallContext represents a UUPSUnauthorizedCallContext error raised by the EndpointV1 contract.
type EndpointV1UUPSUnauthorizedCallContext struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error UUPSUnauthorizedCallContext()
func EndpointV1UUPSUnauthorizedCallContextErrorID() common.Hash {
	return common.HexToHash("0xe07c8dba242a06571ac65fe4bbe20522c9fb111cb33599b799ff8039c1ed18f4")
}

// UnpackUUPSUnauthorizedCallContextError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error UUPSUnauthorizedCallContext()
func (endpointV1 *EndpointV1) UnpackUUPSUnauthorizedCallContextError(raw []byte) (*EndpointV1UUPSUnauthorizedCallContext, error) {
	out := new(EndpointV1UUPSUnauthorizedCallContext)
	if err := endpointV1.abi.UnpackIntoInterface(out, "UUPSUnauthorizedCallContext", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// EndpointV1UUPSUnsupportedProxiableUUID represents a UUPSUnsupportedProxiableUUID error raised by the EndpointV1 contract.
type EndpointV1UUPSUnsupportedProxiableUUID struct {
	Slot [32]byte
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error UUPSUnsupportedProxiableUUID(bytes32 slot)
func EndpointV1UUPSUnsupportedProxiableUUIDErrorID() common.Hash {
	return common.HexToHash("0xaa1d49a4c084bfa9aeeee2a0be65267a7f19ba7e1476b114dac513d2c14cb563")
}

// UnpackUUPSUnsupportedProxiableUUIDError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error UUPSUnsupportedProxiableUUID(bytes32 slot)
func (endpointV1 *EndpointV1) UnpackUUPSUnsupportedProxiableUUIDError(raw []byte) (*EndpointV1UUPSUnsupportedProxiableUUID, error) {
	out := new(EndpointV1UUPSUnsupportedProxiableUUID)
	if err := endpointV1.abi.UnpackIntoInterface(out, "UUPSUnsupportedProxiableUUID", raw); err != nil {
		return nil, err
	}
	return out, nil
}
