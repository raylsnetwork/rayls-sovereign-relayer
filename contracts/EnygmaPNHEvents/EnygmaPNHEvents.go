// Code generated via abigen V2 - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package EnygmaPNHEvents

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

// EnygmaPNHEventsMetaData contains all meta data concerning the EnygmaPNHEvents contract.
var EnygmaPNHEventsMetaData = bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"_endpointAddress\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"authority\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"enygmaPnTransferCompleted\",\"inputs\":[{\"name\":\"encryptedMessage\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"mintCompleted\",\"inputs\":[{\"name\":\"encryptedMessage\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"resourceId\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"setEndpoint\",\"inputs\":[{\"name\":\"_endpoint\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"event\",\"name\":\"AuthorityUpdated\",\"inputs\":[{\"name\":\"oldAuthority\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newAuthority\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"EndpointUpdated\",\"inputs\":[{\"name\":\"oldEndpoint\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newEndpoint\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"EnygmaPnTransferCompleted\",\"inputs\":[{\"name\":\"encryptedMessage\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"MintCompleted\",\"inputs\":[{\"name\":\"encryptedMessage\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"EnygmaPNHEvents__InvalidEndpointAddress\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"EnygmaPNHEvents__UnauthorizedExecutor\",\"inputs\":[{\"name\":\"caller\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"RaylsAccessManaged__ContractPaused\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"RaylsAccessManaged__InvalidAuthority\",\"inputs\":[{\"name\":\"authority\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"RaylsAccessManaged__MustSchedule\",\"inputs\":[{\"name\":\"caller\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"delay\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"type\":\"error\",\"name\":\"RaylsAccessManaged__Unauthorized\",\"inputs\":[{\"name\":\"caller\",\"type\":\"address\",\"internalType\":\"address\"}]}]",
	ID:  "EnygmaPNHEvents",
	Bin: "0x60a060405234801561001057600080fd5b506040516107e83803806107e883398101604081905261002f916101f6565b6001600160a01b038116610056576040516303e8288360e31b815260040160405180910390fd5b600080546001600160a01b0319166001600160a01b0383169081178255600460808190526040805163bf7e214f60e01b8152905163bf7e214f928281019260209291908290030181865afa1580156100b2573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906100d691906101f6565b90506001600160a01b038116156100f0576100f0816100f7565b505061024d565b6001600160a01b03811661012d57604051638944034760e01b81526001600160a01b038216600482015260240160405180910390fd5b6000610137610194565b80546040519192506001600160a01b03808516929116907fa3396fd7f6e0a21b50e5089d2da70d5ac0a3bbbd1f617a93f134b7638998019890600090a380546001600160a01b0319166001600160a01b0392909216919091179055565b60008060ff196101c560017f2d9a26166e6ed6927c3e14d89c4d572a357f984344fcaed3c38c3ff5cdb83f35610226565b6040516020016101d791815260200190565b60408051601f1981840301815291905280516020909101201692915050565b60006020828403121561020857600080fd5b81516001600160a01b038116811461021f57600080fd5b9392505050565b8181038181111561024757634e487b7160e01b600052601160045260246000fd5b92915050565b6080516105816102676000396000607601526105816000f3fe608060405234801561001057600080fd5b50600436106100575760003560e01c80634f4aa3611461005c5780635f997c5b14610071578063bf7e214f146100ab578063d452bf06146100c0578063dbbb4155146100d3575b600080fd5b61006f61006a3660046103dc565b6100e6565b005b6100987f000000000000000000000000000000000000000000000000000000000000000081565b6040519081526020015b60405180910390f35b6100b3610139565b6040516100a2919061044e565b61006f6100ce3660046103dc565b610152565b61006f6100e1366004610462565b610199565b6100fc336000356001600160e01b031916610226565b7f7ddf329588eb36a7c55d811e526e12aa4be5727655de56ca03a89935040733ee828260405161012d929190610492565b60405180910390a15050565b600061014361037a565b546001600160a01b0316919050565b610168336000356001600160e01b031916610226565b7f8748893e53a0cead91711148bc3bfd8a6bbfead8486044464cc8c9f5239d2ff1828260405161012d929190610492565b6101af336000356001600160e01b031916610226565b6001600160a01b0381166101d6576040516303e8288360e31b815260040160405180910390fd5b600080546001600160a01b038381166001600160a01b0319831681178455604051919092169283917f241827194635b544bceb965d0e124305c0968051403ae83cb99d1ef86bd65f589190a35050565b600061023061037a565b80549091506001600160a01b031680610268576000604051638944034760e01b815260040161025f919061044e565b60405180910390fd5b60405163b700961360e01b81526001600160a01b0385811660048301523060248301526001600160e01b031985166044830152600091829182919085169063b700961390606401606060405180830381865afa1580156102cc573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906102f091906104d6565b9250925092508261037157801561031a5760405163cc9855ad60e01b815260040160405180910390fd5b63ffffffff8216156103565760405163a426878960e01b81526001600160a01b038816600482015263ffffffff8316602482015260440161025f565b86604051632ecd3d0360e21b815260040161025f919061044e565b50505050505050565b60008060ff196103ab60017f2d9a26166e6ed6927c3e14d89c4d572a357f984344fcaed3c38c3ff5cdb83f35610524565b6040516020016103bd91815260200190565b60408051601f1981840301815291905280516020909101201692915050565b600080602083850312156103ef57600080fd5b823567ffffffffffffffff8082111561040757600080fd5b818501915085601f83011261041b57600080fd5b81358181111561042a57600080fd5b86602082850101111561043c57600080fd5b60209290920196919550909350505050565b6001600160a01b0391909116815260200190565b60006020828403121561047457600080fd5b81356001600160a01b038116811461048b57600080fd5b9392505050565b60208152816020820152818360408301376000818301604090810191909152601f909201601f19160101919050565b805180151581146104d157600080fd5b919050565b6000806000606084860312156104eb57600080fd5b6104f4846104c1565b9250602084015163ffffffff8116811461050d57600080fd5b915061051b604085016104c1565b90509250925092565b8181038181111561054557634e487b7160e01b600052601160045260246000fd5b9291505056fea2646970667358221220fe4980fba331a6f82f8c808336cefd5118b962574b8df8d8baaec50801bdbe6964736f6c63430008180033",
}

// EnygmaPNHEvents is an auto generated Go binding around an Ethereum contract.
type EnygmaPNHEvents struct {
	abi abi.ABI
}

// NewEnygmaPNHEvents creates a new instance of EnygmaPNHEvents.
func NewEnygmaPNHEvents() *EnygmaPNHEvents {
	parsed, err := EnygmaPNHEventsMetaData.ParseABI()
	if err != nil {
		panic(errors.New("invalid ABI: " + err.Error()))
	}
	return &EnygmaPNHEvents{abi: *parsed}
}

// Instance creates a wrapper for a deployed contract instance at the given address.
// Use this to create the instance object passed to abigen v2 library functions Call, Transact, etc.
func (c *EnygmaPNHEvents) Instance(backend bind.ContractBackend, addr common.Address) *bind.BoundContract {
	return bind.NewBoundContract(addr, c.abi, backend, backend, backend)
}

// PackConstructor is the Go binding used to pack the parameters required for
// contract deployment.
//
// Solidity: constructor(address _endpointAddress) returns()
func (enygmaPNHEvents *EnygmaPNHEvents) PackConstructor(_endpointAddress common.Address) []byte {
	enc, err := enygmaPNHEvents.abi.Pack("", _endpointAddress)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackAuthority is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xbf7e214f.
//
// Solidity: function authority() view returns(address)
func (enygmaPNHEvents *EnygmaPNHEvents) PackAuthority() []byte {
	enc, err := enygmaPNHEvents.abi.Pack("authority")
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackAuthority is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xbf7e214f.
//
// Solidity: function authority() view returns(address)
func (enygmaPNHEvents *EnygmaPNHEvents) UnpackAuthority(data []byte) (common.Address, error) {
	out, err := enygmaPNHEvents.abi.Unpack("authority", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, err
}

// PackEnygmaPnTransferCompleted is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x4f4aa361.
//
// Solidity: function enygmaPnTransferCompleted(bytes encryptedMessage) returns()
func (enygmaPNHEvents *EnygmaPNHEvents) PackEnygmaPnTransferCompleted(encryptedMessage []byte) []byte {
	enc, err := enygmaPNHEvents.abi.Pack("enygmaPnTransferCompleted", encryptedMessage)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackMintCompleted is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xd452bf06.
//
// Solidity: function mintCompleted(bytes encryptedMessage) returns()
func (enygmaPNHEvents *EnygmaPNHEvents) PackMintCompleted(encryptedMessage []byte) []byte {
	enc, err := enygmaPNHEvents.abi.Pack("mintCompleted", encryptedMessage)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackResourceId is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x5f997c5b.
//
// Solidity: function resourceId() view returns(bytes32)
func (enygmaPNHEvents *EnygmaPNHEvents) PackResourceId() []byte {
	enc, err := enygmaPNHEvents.abi.Pack("resourceId")
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackResourceId is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x5f997c5b.
//
// Solidity: function resourceId() view returns(bytes32)
func (enygmaPNHEvents *EnygmaPNHEvents) UnpackResourceId(data []byte) ([32]byte, error) {
	out, err := enygmaPNHEvents.abi.Unpack("resourceId", data)
	if err != nil {
		return *new([32]byte), err
	}
	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)
	return out0, err
}

// PackSetEndpoint is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xdbbb4155.
//
// Solidity: function setEndpoint(address _endpoint) returns()
func (enygmaPNHEvents *EnygmaPNHEvents) PackSetEndpoint(endpoint common.Address) []byte {
	enc, err := enygmaPNHEvents.abi.Pack("setEndpoint", endpoint)
	if err != nil {
		panic(err)
	}
	return enc
}

// EnygmaPNHEventsAuthorityUpdated represents a AuthorityUpdated event raised by the EnygmaPNHEvents contract.
type EnygmaPNHEventsAuthorityUpdated struct {
	OldAuthority common.Address
	NewAuthority common.Address
	Raw          *types.Log // Blockchain specific contextual infos
}

const EnygmaPNHEventsAuthorityUpdatedEventName = "AuthorityUpdated"

// ContractEventName returns the user-defined event name.
func (EnygmaPNHEventsAuthorityUpdated) ContractEventName() string {
	return EnygmaPNHEventsAuthorityUpdatedEventName
}

// UnpackAuthorityUpdatedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event AuthorityUpdated(address indexed oldAuthority, address indexed newAuthority)
func (enygmaPNHEvents *EnygmaPNHEvents) UnpackAuthorityUpdatedEvent(log *types.Log) (*EnygmaPNHEventsAuthorityUpdated, error) {
	event := "AuthorityUpdated"
	if log.Topics[0] != enygmaPNHEvents.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(EnygmaPNHEventsAuthorityUpdated)
	if len(log.Data) > 0 {
		if err := enygmaPNHEvents.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range enygmaPNHEvents.abi.Events[event].Inputs {
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

// EnygmaPNHEventsEndpointUpdated represents a EndpointUpdated event raised by the EnygmaPNHEvents contract.
type EnygmaPNHEventsEndpointUpdated struct {
	OldEndpoint common.Address
	NewEndpoint common.Address
	Raw         *types.Log // Blockchain specific contextual infos
}

const EnygmaPNHEventsEndpointUpdatedEventName = "EndpointUpdated"

// ContractEventName returns the user-defined event name.
func (EnygmaPNHEventsEndpointUpdated) ContractEventName() string {
	return EnygmaPNHEventsEndpointUpdatedEventName
}

// UnpackEndpointUpdatedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event EndpointUpdated(address indexed oldEndpoint, address indexed newEndpoint)
func (enygmaPNHEvents *EnygmaPNHEvents) UnpackEndpointUpdatedEvent(log *types.Log) (*EnygmaPNHEventsEndpointUpdated, error) {
	event := "EndpointUpdated"
	if log.Topics[0] != enygmaPNHEvents.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(EnygmaPNHEventsEndpointUpdated)
	if len(log.Data) > 0 {
		if err := enygmaPNHEvents.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range enygmaPNHEvents.abi.Events[event].Inputs {
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

// EnygmaPNHEventsEnygmaPnTransferCompleted represents a EnygmaPnTransferCompleted event raised by the EnygmaPNHEvents contract.
type EnygmaPNHEventsEnygmaPnTransferCompleted struct {
	EncryptedMessage []byte
	Raw              *types.Log // Blockchain specific contextual infos
}

const EnygmaPNHEventsEnygmaPnTransferCompletedEventName = "EnygmaPnTransferCompleted"

// ContractEventName returns the user-defined event name.
func (EnygmaPNHEventsEnygmaPnTransferCompleted) ContractEventName() string {
	return EnygmaPNHEventsEnygmaPnTransferCompletedEventName
}

// UnpackEnygmaPnTransferCompletedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event EnygmaPnTransferCompleted(bytes encryptedMessage)
func (enygmaPNHEvents *EnygmaPNHEvents) UnpackEnygmaPnTransferCompletedEvent(log *types.Log) (*EnygmaPNHEventsEnygmaPnTransferCompleted, error) {
	event := "EnygmaPnTransferCompleted"
	if log.Topics[0] != enygmaPNHEvents.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(EnygmaPNHEventsEnygmaPnTransferCompleted)
	if len(log.Data) > 0 {
		if err := enygmaPNHEvents.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range enygmaPNHEvents.abi.Events[event].Inputs {
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

// EnygmaPNHEventsMintCompleted represents a MintCompleted event raised by the EnygmaPNHEvents contract.
type EnygmaPNHEventsMintCompleted struct {
	EncryptedMessage []byte
	Raw              *types.Log // Blockchain specific contextual infos
}

const EnygmaPNHEventsMintCompletedEventName = "MintCompleted"

// ContractEventName returns the user-defined event name.
func (EnygmaPNHEventsMintCompleted) ContractEventName() string {
	return EnygmaPNHEventsMintCompletedEventName
}

// UnpackMintCompletedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event MintCompleted(bytes encryptedMessage)
func (enygmaPNHEvents *EnygmaPNHEvents) UnpackMintCompletedEvent(log *types.Log) (*EnygmaPNHEventsMintCompleted, error) {
	event := "MintCompleted"
	if log.Topics[0] != enygmaPNHEvents.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(EnygmaPNHEventsMintCompleted)
	if len(log.Data) > 0 {
		if err := enygmaPNHEvents.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range enygmaPNHEvents.abi.Events[event].Inputs {
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
func (enygmaPNHEvents *EnygmaPNHEvents) UnpackError(raw []byte) (any, error) {
	if bytes.Equal(raw[:4], enygmaPNHEvents.abi.Errors["EnygmaPNHEventsInvalidEndpointAddress"].ID.Bytes()[:4]) {
		return enygmaPNHEvents.UnpackEnygmaPNHEventsInvalidEndpointAddressError(raw[4:])
	}
	if bytes.Equal(raw[:4], enygmaPNHEvents.abi.Errors["EnygmaPNHEventsUnauthorizedExecutor"].ID.Bytes()[:4]) {
		return enygmaPNHEvents.UnpackEnygmaPNHEventsUnauthorizedExecutorError(raw[4:])
	}
	if bytes.Equal(raw[:4], enygmaPNHEvents.abi.Errors["RaylsAccessManagedContractPaused"].ID.Bytes()[:4]) {
		return enygmaPNHEvents.UnpackRaylsAccessManagedContractPausedError(raw[4:])
	}
	if bytes.Equal(raw[:4], enygmaPNHEvents.abi.Errors["RaylsAccessManagedInvalidAuthority"].ID.Bytes()[:4]) {
		return enygmaPNHEvents.UnpackRaylsAccessManagedInvalidAuthorityError(raw[4:])
	}
	if bytes.Equal(raw[:4], enygmaPNHEvents.abi.Errors["RaylsAccessManagedMustSchedule"].ID.Bytes()[:4]) {
		return enygmaPNHEvents.UnpackRaylsAccessManagedMustScheduleError(raw[4:])
	}
	if bytes.Equal(raw[:4], enygmaPNHEvents.abi.Errors["RaylsAccessManagedUnauthorized"].ID.Bytes()[:4]) {
		return enygmaPNHEvents.UnpackRaylsAccessManagedUnauthorizedError(raw[4:])
	}
	return nil, errors.New("Unknown error")
}

// EnygmaPNHEventsEnygmaPNHEventsInvalidEndpointAddress represents a EnygmaPNHEvents__InvalidEndpointAddress error raised by the EnygmaPNHEvents contract.
type EnygmaPNHEventsEnygmaPNHEventsInvalidEndpointAddress struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error EnygmaPNHEvents__InvalidEndpointAddress()
func EnygmaPNHEventsEnygmaPNHEventsInvalidEndpointAddressErrorID() common.Hash {
	return common.HexToHash("0x1f414418f16e22a3086ab06481bad677e7ca267f74da945721e77991f55785dd")
}

// UnpackEnygmaPNHEventsInvalidEndpointAddressError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error EnygmaPNHEvents__InvalidEndpointAddress()
func (enygmaPNHEvents *EnygmaPNHEvents) UnpackEnygmaPNHEventsInvalidEndpointAddressError(raw []byte) (*EnygmaPNHEventsEnygmaPNHEventsInvalidEndpointAddress, error) {
	out := new(EnygmaPNHEventsEnygmaPNHEventsInvalidEndpointAddress)
	if err := enygmaPNHEvents.abi.UnpackIntoInterface(out, "EnygmaPNHEventsInvalidEndpointAddress", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// EnygmaPNHEventsEnygmaPNHEventsUnauthorizedExecutor represents a EnygmaPNHEvents__UnauthorizedExecutor error raised by the EnygmaPNHEvents contract.
type EnygmaPNHEventsEnygmaPNHEventsUnauthorizedExecutor struct {
	Caller common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error EnygmaPNHEvents__UnauthorizedExecutor(address caller)
func EnygmaPNHEventsEnygmaPNHEventsUnauthorizedExecutorErrorID() common.Hash {
	return common.HexToHash("0x929821f8767f3d8d4c9fc359fe9246617b5698a00dabb7989b62e99514f01cd3")
}

// UnpackEnygmaPNHEventsUnauthorizedExecutorError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error EnygmaPNHEvents__UnauthorizedExecutor(address caller)
func (enygmaPNHEvents *EnygmaPNHEvents) UnpackEnygmaPNHEventsUnauthorizedExecutorError(raw []byte) (*EnygmaPNHEventsEnygmaPNHEventsUnauthorizedExecutor, error) {
	out := new(EnygmaPNHEventsEnygmaPNHEventsUnauthorizedExecutor)
	if err := enygmaPNHEvents.abi.UnpackIntoInterface(out, "EnygmaPNHEventsUnauthorizedExecutor", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// EnygmaPNHEventsRaylsAccessManagedContractPaused represents a RaylsAccessManaged__ContractPaused error raised by the EnygmaPNHEvents contract.
type EnygmaPNHEventsRaylsAccessManagedContractPaused struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error RaylsAccessManaged__ContractPaused()
func EnygmaPNHEventsRaylsAccessManagedContractPausedErrorID() common.Hash {
	return common.HexToHash("0xcc9855adb54c3eae80b5cb29441d3eb81c3fcdc928abeaba0ac7979b3072c3ab")
}

// UnpackRaylsAccessManagedContractPausedError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error RaylsAccessManaged__ContractPaused()
func (enygmaPNHEvents *EnygmaPNHEvents) UnpackRaylsAccessManagedContractPausedError(raw []byte) (*EnygmaPNHEventsRaylsAccessManagedContractPaused, error) {
	out := new(EnygmaPNHEventsRaylsAccessManagedContractPaused)
	if err := enygmaPNHEvents.abi.UnpackIntoInterface(out, "RaylsAccessManagedContractPaused", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// EnygmaPNHEventsRaylsAccessManagedInvalidAuthority represents a RaylsAccessManaged__InvalidAuthority error raised by the EnygmaPNHEvents contract.
type EnygmaPNHEventsRaylsAccessManagedInvalidAuthority struct {
	Authority common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error RaylsAccessManaged__InvalidAuthority(address authority)
func EnygmaPNHEventsRaylsAccessManagedInvalidAuthorityErrorID() common.Hash {
	return common.HexToHash("0x894403477abbba3df1c9b757907cda2d4b4d9ba52242d5e6bdae2dc6a343662b")
}

// UnpackRaylsAccessManagedInvalidAuthorityError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error RaylsAccessManaged__InvalidAuthority(address authority)
func (enygmaPNHEvents *EnygmaPNHEvents) UnpackRaylsAccessManagedInvalidAuthorityError(raw []byte) (*EnygmaPNHEventsRaylsAccessManagedInvalidAuthority, error) {
	out := new(EnygmaPNHEventsRaylsAccessManagedInvalidAuthority)
	if err := enygmaPNHEvents.abi.UnpackIntoInterface(out, "RaylsAccessManagedInvalidAuthority", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// EnygmaPNHEventsRaylsAccessManagedMustSchedule represents a RaylsAccessManaged__MustSchedule error raised by the EnygmaPNHEvents contract.
type EnygmaPNHEventsRaylsAccessManagedMustSchedule struct {
	Caller common.Address
	Delay  uint32
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error RaylsAccessManaged__MustSchedule(address caller, uint32 delay)
func EnygmaPNHEventsRaylsAccessManagedMustScheduleErrorID() common.Hash {
	return common.HexToHash("0xa42687899abecbc3886d9231154d03927d3e61df458dda1e50129311f272684f")
}

// UnpackRaylsAccessManagedMustScheduleError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error RaylsAccessManaged__MustSchedule(address caller, uint32 delay)
func (enygmaPNHEvents *EnygmaPNHEvents) UnpackRaylsAccessManagedMustScheduleError(raw []byte) (*EnygmaPNHEventsRaylsAccessManagedMustSchedule, error) {
	out := new(EnygmaPNHEventsRaylsAccessManagedMustSchedule)
	if err := enygmaPNHEvents.abi.UnpackIntoInterface(out, "RaylsAccessManagedMustSchedule", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// EnygmaPNHEventsRaylsAccessManagedUnauthorized represents a RaylsAccessManaged__Unauthorized error raised by the EnygmaPNHEvents contract.
type EnygmaPNHEventsRaylsAccessManagedUnauthorized struct {
	Caller common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error RaylsAccessManaged__Unauthorized(address caller)
func EnygmaPNHEventsRaylsAccessManagedUnauthorizedErrorID() common.Hash {
	return common.HexToHash("0xbb34f40c15462e9dd86b312158515c7ee80e2406cfa90a92a6fb6fad1bc48a48")
}

// UnpackRaylsAccessManagedUnauthorizedError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error RaylsAccessManaged__Unauthorized(address caller)
func (enygmaPNHEvents *EnygmaPNHEvents) UnpackRaylsAccessManagedUnauthorizedError(raw []byte) (*EnygmaPNHEventsRaylsAccessManagedUnauthorized, error) {
	out := new(EnygmaPNHEventsRaylsAccessManagedUnauthorized)
	if err := enygmaPNHEvents.abi.UnpackIntoInterface(out, "RaylsAccessManagedUnauthorized", raw); err != nil {
		return nil, err
	}
	return out, nil
}
