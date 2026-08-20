package testdata

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"

	"github.com/raylsnetwork/rayls-sovereign-relayer/contracts/EnygmaPNEvents"
)

// EnygmaCreationLogOption represents a modification to an EnygmaCreation log fixture.
type EnygmaCreationLogOption func(*types.Log, *enygmaCreationData)

type enygmaCreationData struct {
	resourceID    [32]byte
	initialSupply *big.Int
}

// WithCreateResourceID sets the ResourceID field on an EnygmaCreation log when using options.
func WithCreateResourceID(resourceID [32]byte) EnygmaCreationLogOption {
	return func(log *types.Log, data *enygmaCreationData) {
		data.resourceID = resourceID
	}
}

// WithInitialSupply sets the InitialSupply field on an EnygmaCreation log when using options.
func WithCreateInitialSupply(initialSupply *big.Int) EnygmaCreationLogOption {
	return func(log *types.Log, data *enygmaCreationData) {
		data.initialSupply = initialSupply
	}
}

// WithCreateBlockNumber sets the BlockNumber field on an EnygmaCreation log when using options.
func WithCreateBlockNumber(blockNumber uint64) EnygmaCreationLogOption {
	return func(log *types.Log, data *enygmaCreationData) {
		log.BlockNumber = blockNumber
	}
}

// GetEnygmaCreationLog returns a sample EnygmaCreation event log with default values
func GetEnygmaCreationLog() types.Log {
	return NewEnygmaCreationLogWith()
}

// NewEnygmaCreationLogWith creates an EnygmaCreation log using default values
// and then applies any provided options.
func NewEnygmaCreationLogWith(opts ...EnygmaCreationLogOption) types.Log {
	// Default values
	data := enygmaCreationData{
		resourceID:    [32]byte{0x01, 0x02, 0x03},
		initialSupply: big.NewInt(1000000),
	}

	log := types.Log{
		Address:     common.HexToAddress("0x1234567890123456789012345678901234567890"),
		BlockNumber: 100,
		TxHash:      common.HexToHash("0xabcdef"),
		TxIndex:     1,
		BlockHash:   common.HexToHash("0x123456"),
		Index:       0,
	}

	// Apply options
	for _, opt := range opts {
		opt(&log, &data)
	}

	// Pack the data using ABI
	abi, _ := EnygmaPNEvents.EnygmaPNEventsMetaData.ParseABI()
	packedData, _ := abi.Events["EnygmaCreation"].Inputs.NonIndexed().Pack(data.resourceID, data.initialSupply)

	log.Topics = []common.Hash{
		crypto.Keccak256Hash([]byte("EnygmaCreation(bytes32,uint256)")),
	}
	log.Data = packedData

	return log
}

type EnygmaMintLogOption func(*types.Log, *enygmaMintData)

type enygmaMintData struct {
	resourceID [32]byte
	to         common.Address
	amount     *big.Int
}

func WithMintResourceID(resourceID [32]byte) EnygmaMintLogOption {
	return func(l *types.Log, data *enygmaMintData) {
		data.resourceID = resourceID
	}
}

func WithMintTxHash(hash common.Hash) EnygmaMintLogOption {
	return func(l *types.Log, emd *enygmaMintData) {
		l.TxHash = hash
	}
}

func WithMintToAddress(to common.Address) EnygmaMintLogOption {
	return func(l *types.Log, data *enygmaMintData) {
		data.to = to
	}
}

func WithMintAmount(amount *big.Int) EnygmaMintLogOption {
	return func(l *types.Log, data *enygmaMintData) {
		data.amount = amount
	}
}

func WithMintBlockNumber(blockNumber uint64) EnygmaMintLogOption {
	return func(l *types.Log, data *enygmaMintData) {
		l.BlockNumber = blockNumber
	}
}

func NewEnygmaMintLogWith(opts ...EnygmaMintLogOption) types.Log {
	data := enygmaMintData{
		resourceID: [32]byte{0x01, 0x02, 0x03},
		to:         common.HexToAddress("0x1111111111111111111111111111111111111111"),
		amount:     big.NewInt(500),
	}

	log := types.Log{
		Address: common.HexToAddress("0x1234567890123456789012345678901234567890"),
		Topics: []common.Hash{
			crypto.Keccak256Hash([]byte("EnygmaMint(bytes32,address,uint256)")),
		},
		BlockNumber: 101,
		TxHash:      common.HexToHash("0xabcdef01"),
		TxIndex:     2,
		BlockHash:   common.HexToHash("0x123457"),
		Index:       1,
	}

	for _, opt := range opts {
		opt(&log, &data)
	}

	abi, _ := EnygmaPNEvents.EnygmaPNEventsMetaData.ParseABI()
	packedData, _ := abi.Events["EnygmaMint"].Inputs.NonIndexed().Pack(data.resourceID, data.to, data.amount)
	log.Topics = []common.Hash{
		crypto.Keccak256Hash([]byte("EnygmaMint(bytes32,address,uint256)")),
	}
	log.Data = packedData

	return log
}

// GetEnygmaMintLog returns a sample EnygmaMint event log
func GetEnygmaMintLog() types.Log {
	return NewEnygmaMintLogWith()
}

type EnygmaBurnLogOption func(*types.Log, *enygmaBurnData)

type enygmaBurnData struct {
	resourceID [32]byte
	from       common.Address
	amount     *big.Int
}

func WithBurnResourceID(resourceID [32]byte) EnygmaBurnLogOption {
	return func(l *types.Log, data *enygmaBurnData) {
		data.resourceID = resourceID
	}
}

func WithBurnTxHash(hash common.Hash) EnygmaBurnLogOption {
	return func(l *types.Log, ebd *enygmaBurnData) {
		l.TxHash = hash
	}
}

func WithBurnFromAddress(from common.Address) EnygmaBurnLogOption {
	return func(l *types.Log, data *enygmaBurnData) {
		data.from = from
	}
}

func WithBurnAmount(amount *big.Int) EnygmaBurnLogOption {
	return func(l *types.Log, data *enygmaBurnData) {
		data.amount = amount
	}
}

func WithBurnBlockNumber(blockNumber uint64) EnygmaBurnLogOption {
	return func(l *types.Log, data *enygmaBurnData) {
		l.BlockNumber = blockNumber
	}
}

func NewEnygmaBurnLogWith(opts ...EnygmaBurnLogOption) types.Log {
	data := enygmaBurnData{
		resourceID: [32]byte{0x01, 0x02, 0x03},
		from:       common.HexToAddress("0x2222222222222222222222222222222222222222"),
		amount:     big.NewInt(200),
	}

	log := types.Log{
		Address: common.HexToAddress("0x1234567890123456789012345678901234567890"),
		Topics: []common.Hash{
			crypto.Keccak256Hash([]byte("EnygmaBurn(bytes32,address,uint256)")),
		},
		BlockNumber: 102,
		TxHash:      common.HexToHash("0xabcdef02"),
		TxIndex:     3,
		BlockHash:   common.HexToHash("0x123458"),
		Index:       2,
	}

	for _, opt := range opts {
		opt(&log, &data)
	}

	abi, _ := EnygmaPNEvents.EnygmaPNEventsMetaData.ParseABI()
	packedData, _ := abi.Events["EnygmaBurn"].Inputs.NonIndexed().Pack(data.resourceID, data.from, data.amount)
	log.Topics = []common.Hash{
		crypto.Keccak256Hash([]byte("EnygmaBurn(bytes32,address,uint256)")),
	}
	log.Data = packedData

	return log
}

// GetEnygmaBurnLog returns a sample EnygmaBurn event log
func GetEnygmaBurnLog() types.Log {
	return NewEnygmaBurnLogWith()
}

// EnygmaDepositToDvpLogOption represents a modification to an EnygmaDepositToDvp log fixture.
type EnygmaDepositToDvpLogOption func(*types.Log, *enygmaDepositToDvpData)

type enygmaDepositToDvpData struct {
	resourceID  [32]byte
	amount      *big.Int
	from        common.Address
	referenceID [32]byte
}

// WithDepositResourceID sets the ResourceID field on an EnygmaDepositToDvp log when using options.
func WithDepositResourceID(resourceID [32]byte) EnygmaDepositToDvpLogOption {
	return func(log *types.Log, data *enygmaDepositToDvpData) {
		data.resourceID = resourceID
	}
}

// WithDepositAmount sets the Amount field on an EnygmaDepositToDvp log when using options.
func WithDepositAmount(amount *big.Int) EnygmaDepositToDvpLogOption {
	return func(log *types.Log, data *enygmaDepositToDvpData) {
		data.amount = amount
	}
}

// WithDepositFrom sets the From address field on an EnygmaDepositToDvp log when using options.
func WithDepositFrom(from common.Address) EnygmaDepositToDvpLogOption {
	return func(log *types.Log, data *enygmaDepositToDvpData) {
		data.from = from
	}
}

// WithDepositReferenceID sets the ReferenceID field on an EnygmaDepositToDvp log when using options.
func WithDepositReferenceID(referenceID [32]byte) EnygmaDepositToDvpLogOption {
	return func(log *types.Log, data *enygmaDepositToDvpData) {
		data.referenceID = referenceID
	}
}

// WithDepositBlockNumber sets the BlockNumber field on an EnygmaDepositToDvp log when using options.
func WithDepositBlockNumber(blockNumber uint64) EnygmaDepositToDvpLogOption {
	return func(log *types.Log, data *enygmaDepositToDvpData) {
		log.BlockNumber = blockNumber
	}
}

// GetEnygmaDepositToDvpLog returns a sample EnygmaDepositToDvp event log with default values
func GetEnygmaDepositToDvpLog() types.Log {
	return NewEnygmaDepositToDvpLogWith()
}

// NewEnygmaDepositToDvpLogWith creates an EnygmaDepositToDvp log using default values
// and then applies any provided options.
func NewEnygmaDepositToDvpLogWith(opts ...EnygmaDepositToDvpLogOption) types.Log {
	// Default values
	data := enygmaDepositToDvpData{
		resourceID:  [32]byte{0x01, 0x02, 0x03},
		amount:      big.NewInt(300),
		from:        common.HexToAddress("0x3333333333333333333333333333333333333333"),
		referenceID: [32]byte{0x04, 0x05, 0x06},
	}

	log := types.Log{
		Address:     common.HexToAddress("0x1234567890123456789012345678901234567890"),
		BlockNumber: 103,
		TxHash:      common.HexToHash("0xabcdef03"),
		TxIndex:     4,
		BlockHash:   common.HexToHash("0x123459"),
		Index:       3,
	}

	// Apply options
	for _, opt := range opts {
		opt(&log, &data)
	}

	// Pack the data using ABI
	abi, _ := EnygmaPNEvents.EnygmaPNEventsMetaData.ParseABI()
	packedData, _ := abi.Events["EnygmaDepositToDvp"].Inputs.NonIndexed().Pack(data.resourceID, data.amount, data.from, data.referenceID)

	log.Topics = []common.Hash{
		crypto.Keccak256Hash([]byte("EnygmaDepositToDvp(bytes32,uint256,address,bytes32)")),
	}
	log.Data = packedData

	return log
}

// EnygmaWithdrawFromDvpLogOption represents a modification to an EnygmaWithdrawFromDvp log fixture.
type EnygmaWithdrawFromDvpLogOption func(*types.Log, *enygmaWithdrawFromDvpData)

type enygmaWithdrawFromDvpData struct {
	resourceID  [32]byte
	amount      *big.Int
	to          common.Address
	referenceID [32]byte
}

// WithWithdrawResourceID sets the ResourceID field on an EnygmaWithdrawFromDvp log when using options.
func WithWithdrawResourceID(resourceID [32]byte) EnygmaWithdrawFromDvpLogOption {
	return func(log *types.Log, data *enygmaWithdrawFromDvpData) {
		data.resourceID = resourceID
	}
}

// WithWithdrawAmount sets the Amount field on an EnygmaWithdrawFromDvp log when using options.
func WithWithdrawAmount(amount *big.Int) EnygmaWithdrawFromDvpLogOption {
	return func(log *types.Log, data *enygmaWithdrawFromDvpData) {
		data.amount = amount
	}
}

// WithWithdrawTo sets the To address field on an EnygmaWithdrawFromDvp log when using options.
func WithWithdrawTo(to common.Address) EnygmaWithdrawFromDvpLogOption {
	return func(log *types.Log, data *enygmaWithdrawFromDvpData) {
		data.to = to
	}
}

// WithWithdrawReferenceID sets the ReferenceID field on an EnygmaWithdrawFromDvp log when using options.
func WithWithdrawReferenceID(referenceID [32]byte) EnygmaWithdrawFromDvpLogOption {
	return func(log *types.Log, data *enygmaWithdrawFromDvpData) {
		data.referenceID = referenceID
	}
}

// WithWithdrawBlockNumber sets the BlockNumber field on an EnygmaWithdrawFromDvp log when using options.
func WithWithdrawBlockNumber(blockNumber uint64) EnygmaWithdrawFromDvpLogOption {
	return func(log *types.Log, data *enygmaWithdrawFromDvpData) {
		log.BlockNumber = blockNumber
	}
}

// GetEnygmaWithdrawFromDvpLog returns a sample EnygmaWithdrawFromDvp event log with default values
func GetEnygmaWithdrawFromDvpLog() types.Log {
	return NewEnygmaWithdrawFromDvpLogWith()
}

// NewEnygmaWithdrawFromDvpLogWith creates an EnygmaWithdrawFromDvp log using default values
// and then applies any provided options.
func NewEnygmaWithdrawFromDvpLogWith(opts ...EnygmaWithdrawFromDvpLogOption) types.Log {
	// Default values
	data := enygmaWithdrawFromDvpData{
		resourceID:  [32]byte{0x01, 0x02, 0x03},
		amount:      big.NewInt(150),
		to:          common.HexToAddress("0x4444444444444444444444444444444444444444"),
		referenceID: [32]byte{0x07, 0x08, 0x09},
	}

	log := types.Log{
		Address:     common.HexToAddress("0x1234567890123456789012345678901234567890"),
		BlockNumber: 104,
		TxHash:      common.HexToHash("0xabcdef04"),
		TxIndex:     5,
		BlockHash:   common.HexToHash("0x12345a"),
		Index:       4,
	}

	// Apply options
	for _, opt := range opts {
		opt(&log, &data)
	}

	// Pack the data using ABI
	abi, _ := EnygmaPNEvents.EnygmaPNEventsMetaData.ParseABI()
	packedData, _ := abi.Events["EnygmaWithdrawFromDvp"].Inputs.NonIndexed().Pack(data.resourceID, data.amount, data.to, data.referenceID)

	log.Topics = []common.Hash{
		crypto.Keccak256Hash([]byte("EnygmaWithdrawFromDvp(bytes32,uint256,address,bytes32)")),
	}
	log.Data = packedData

	return log
}

// GetDvp721CreationLog returns a sample Dvp721Creation event log
func GetDvp721CreationLog() types.Log {
	resourceId := [32]byte{0x0a, 0x0b, 0x0c}

	abi, _ := EnygmaPNEvents.EnygmaPNEventsMetaData.ParseABI()
	data, _ := abi.Events["Dvp721Creation"].Inputs.NonIndexed().Pack(resourceId)

	return types.Log{
		Address: common.HexToAddress("0x1234567890123456789012345678901234567890"),
		Topics: []common.Hash{
			crypto.Keccak256Hash([]byte("Dvp721Creation(bytes32)")),
		},
		Data:        data,
		BlockNumber: 105,
		TxHash:      common.HexToHash("0xabcdef05"),
		TxIndex:     6,
		BlockHash:   common.HexToHash("0x12345b"),
		Index:       5,
	}
}

// GetDvp721MintLog returns a sample Dvp721Mint event log
func GetDvp721MintLog() types.Log {
	resourceId := [32]byte{0x0a, 0x0b, 0x0c}
	nftId := big.NewInt(42)

	abi, _ := EnygmaPNEvents.EnygmaPNEventsMetaData.ParseABI()
	data, _ := abi.Events["Dvp721Mint"].Inputs.NonIndexed().Pack(resourceId, nftId)

	return types.Log{
		Address: common.HexToAddress("0x1234567890123456789012345678901234567890"),
		Topics: []common.Hash{
			crypto.Keccak256Hash([]byte("Dvp721Mint(bytes32,uint256)")),
		},
		Data:        data,
		BlockNumber: 106,
		TxHash:      common.HexToHash("0xabcdef06"),
		TxIndex:     7,
		BlockHash:   common.HexToHash("0x12345c"),
		Index:       6,
	}
}

// GetDvp721BurnLog returns a sample Dvp721Burn event log
func GetDvp721BurnLog() types.Log {
	resourceId := [32]byte{0x0a, 0x0b, 0x0c}
	nftId := big.NewInt(42)

	abi, _ := EnygmaPNEvents.EnygmaPNEventsMetaData.ParseABI()
	data, _ := abi.Events["Dvp721Burn"].Inputs.NonIndexed().Pack(resourceId, nftId)

	return types.Log{
		Address: common.HexToAddress("0x1234567890123456789012345678901234567890"),
		Topics: []common.Hash{
			crypto.Keccak256Hash([]byte("Dvp721Burn(bytes32,uint256)")),
		},
		Data:        data,
		BlockNumber: 107,
		TxHash:      common.HexToHash("0xabcdef07"),
		TxIndex:     8,
		BlockHash:   common.HexToHash("0x12345d"),
		Index:       7,
	}
}

// EnygmaSendTransferPNHLogOption represents a modification to an EnygmaSendTransferPNH log fixture.
type EnygmaSendTransferPNHLogOption func(*types.Log, *enygmaSendTransferCCData)

type enygmaSendTransferCCData struct {
	resourceID  [32]byte
	value       []*big.Int
	toChainId   []*big.Int
	to          []common.Address
	from        common.Address
	referenceID [32]byte
	// programData is parallel to to/value/toChainId: one array of steps per recipient.
	programData [][]EnygmaPNEvents.SharedObjectsEnygmaProgramData
}

// WithTransferResourceID sets the ResourceID field
func WithTransferResourceID(resourceID [32]byte) EnygmaSendTransferPNHLogOption {
	return func(log *types.Log, data *enygmaSendTransferCCData) {
		data.resourceID = resourceID
	}
}

// WithTransferValue sets the Value field
func WithTransferValue(value []*big.Int) EnygmaSendTransferPNHLogOption {
	return func(log *types.Log, data *enygmaSendTransferCCData) {
		data.value = value
	}
}

// WithTransferToChainId sets the ToChainId field
func WithTransferToChainId(toChainId []*big.Int) EnygmaSendTransferPNHLogOption {
	return func(log *types.Log, data *enygmaSendTransferCCData) {
		data.toChainId = toChainId
	}
}

// WithTransferTo sets the To field
func WithTransferTo(to []common.Address) EnygmaSendTransferPNHLogOption {
	return func(log *types.Log, data *enygmaSendTransferCCData) {
		data.to = to
	}
}

// WithTransferFrom sets the From field
func WithTransferFrom(from common.Address) EnygmaSendTransferPNHLogOption {
	return func(log *types.Log, data *enygmaSendTransferCCData) {
		data.from = from
	}
}

// WithTransferReferenceID sets the ReferenceID field
func WithTransferReferenceID(referenceID [32]byte) EnygmaSendTransferPNHLogOption {
	return func(log *types.Log, data *enygmaSendTransferCCData) {
		data.referenceID = referenceID
	}
}

// WithTransferProgramData sets the per-recipient programData arrays (parallel to to/value/toChainId).
func WithTransferProgramData(programData [][]EnygmaPNEvents.SharedObjectsEnygmaProgramData) EnygmaSendTransferPNHLogOption {
	return func(log *types.Log, data *enygmaSendTransferCCData) {
		data.programData = programData
	}
}

// WithTransferBlockNumber sets the BlockNumber field
func WithTransferBlockNumber(blockNumber uint64) EnygmaSendTransferPNHLogOption {
	return func(log *types.Log, data *enygmaSendTransferCCData) {
		log.BlockNumber = blockNumber
	}
}

// GetEnygmaSendTransferPNHLog returns a sample EnygmaSendTransferPNH event log
func GetEnygmaSendTransferPNHLog() types.Log {
	return NewEnygmaSendTransferPNHLogWith()
}

// NewEnygmaSendTransferPNHLogWith creates an EnygmaSendTransferPNH log using default values
// and then applies any provided options.
func NewEnygmaSendTransferPNHLogWith(opts ...EnygmaSendTransferPNHLogOption) types.Log {
	// Default values
	data := enygmaSendTransferCCData{
		resourceID: [32]byte{0x01, 0x02, 0x03},
		value:      []*big.Int{big.NewInt(100), big.NewInt(200)},
		toChainId:  []*big.Int{big.NewInt(1), big.NewInt(2)},
		to: []common.Address{
			common.HexToAddress("0x1111111111111111111111111111111111111111"),
			common.HexToAddress("0x2222222222222222222222222222222222222222"),
		},
		from:        common.HexToAddress("0x3333333333333333333333333333333333333333"),
		referenceID: [32]byte{0x04, 0x05, 0x06},
		// One single-element program-data array per recipient, mirroring a plain
		// transfer's [mintStep] stamped by the sender handler. Distinct Args per recipient
		// so tests can assert the parallel split preserved per-recipient order.
		programData: [][]EnygmaPNEvents.SharedObjectsEnygmaProgramData{
			{{ResourceId: [32]byte{0x0A}, ContractAddress: common.Address{}, Selector: [4]byte{0x11}, Args: []byte{0xDE, 0xAD}}},
			{{ResourceId: [32]byte{0x0B}, ContractAddress: common.Address{}, Selector: [4]byte{0x22}, Args: []byte{0xBE, 0xEF}}},
		},
	}

	log := types.Log{
		Address:     common.HexToAddress("0x1234567890123456789012345678901234567890"),
		BlockNumber: 108,
		TxHash:      common.HexToHash("0xabcdef08"),
		TxIndex:     9,
		BlockHash:   common.HexToHash("0x12345e"),
		Index:       8,
	}

	// Apply options
	for _, opt := range opts {
		opt(&log, &data)
	}

	// Pack the data using ABI
	abi, _ := EnygmaPNEvents.EnygmaPNEventsMetaData.ParseABI()
	packedData, _ := abi.Events["EnygmaSendTransferPNH"].Inputs.NonIndexed().Pack(
		data.resourceID, data.value, data.toChainId, data.to, data.from, data.referenceID, data.programData,
	)

	log.Topics = []common.Hash{
		crypto.Keccak256Hash([]byte("EnygmaSendTransferPNH(bytes32,uint256[],uint256[],address[],address,bytes32,(bytes32,address,bytes4,bytes)[][])")),
	}
	log.Data = packedData

	return log
}

// EnygmaSwapWithDvpForERC721LogOption represents a modification to an EnygmaSwapWithDvpForERC721 log fixture.
type EnygmaSwapWithDvpForERC721LogOption func(*types.Log, *enygmaSwapWithDvpForERC721Data)

type enygmaSwapWithDvpForERC721Data struct {
	resourceID    [32]byte
	nftId         *big.Int
	nftResourceId [32]byte
	enygmaAmount  *big.Int
	from          common.Address
	destChainId   *big.Int
	sharedId      [32]byte
	calldata      []byte
	validityTime  uint64
}

func WithSwap721ResourceID(resourceID [32]byte) EnygmaSwapWithDvpForERC721LogOption {
	return func(log *types.Log, data *enygmaSwapWithDvpForERC721Data) {
		data.resourceID = resourceID
	}
}

func WithSwap721NftResourceID(nftResourceId [32]byte) EnygmaSwapWithDvpForERC721LogOption {
	return func(log *types.Log, data *enygmaSwapWithDvpForERC721Data) {
		data.nftResourceId = nftResourceId
	}
}

func WithSwap721SharedID(sharedId [32]byte) EnygmaSwapWithDvpForERC721LogOption {
	return func(log *types.Log, data *enygmaSwapWithDvpForERC721Data) {
		data.sharedId = sharedId
	}
}

func WithSwap721NftId(nftId *big.Int) EnygmaSwapWithDvpForERC721LogOption {
	return func(log *types.Log, data *enygmaSwapWithDvpForERC721Data) {
		data.nftId = nftId
	}
}

func WithSwap721EnygmaAmount(enygmaAmount *big.Int) EnygmaSwapWithDvpForERC721LogOption {
	return func(log *types.Log, data *enygmaSwapWithDvpForERC721Data) {
		data.enygmaAmount = enygmaAmount
	}
}

func WithSwap721From(from common.Address) EnygmaSwapWithDvpForERC721LogOption {
	return func(log *types.Log, data *enygmaSwapWithDvpForERC721Data) {
		data.from = from
	}
}

func WithSwap721DestChainId(destChainId *big.Int) EnygmaSwapWithDvpForERC721LogOption {
	return func(log *types.Log, data *enygmaSwapWithDvpForERC721Data) {
		data.destChainId = destChainId
	}
}

func WithSwap721Calldata(calldata []byte) EnygmaSwapWithDvpForERC721LogOption {
	return func(log *types.Log, data *enygmaSwapWithDvpForERC721Data) {
		data.calldata = calldata
	}
}

func WithSwap721BlockNumber(blockNumber uint64) EnygmaSwapWithDvpForERC721LogOption {
	return func(log *types.Log, data *enygmaSwapWithDvpForERC721Data) {
		log.BlockNumber = blockNumber
	}
}

// GetEnygmaSwapWithDvpForERC721Log returns a sample EnygmaSwapWithDvpForERC721 event log
func GetEnygmaSwapWithDvpForERC721Log() types.Log {
	return NewEnygmaSwapWithDvpForERC721LogWith()
}

// NewEnygmaSwapWithDvpForERC721LogWith creates an EnygmaSwapWithDvpForERC721 log using default values
// and then applies any provided options.
func NewEnygmaSwapWithDvpForERC721LogWith(opts ...EnygmaSwapWithDvpForERC721LogOption) types.Log {
	// Default values
	data := enygmaSwapWithDvpForERC721Data{
		resourceID:    [32]byte{0x01, 0x02, 0x03},
		nftId:         big.NewInt(123),
		nftResourceId: [32]byte{0x0a, 0x0b, 0x0c},
		enygmaAmount:  big.NewInt(500),
		from:          common.HexToAddress("0x4444444444444444444444444444444444444444"),
		destChainId:   big.NewInt(2),
		sharedId:      [32]byte{0x0d, 0x0e, 0x0f},
		calldata:      []byte{0x01, 0x02, 0x03},
		validityTime:  1337,
	}

	log := types.Log{
		Address:     common.HexToAddress("0x1234567890123456789012345678901234567890"),
		BlockNumber: 109,
		TxHash:      common.HexToHash("0xabcdef09"),
		TxIndex:     10,
		BlockHash:   common.HexToHash("0x12345f"),
		Index:       9,
	}

	// Apply options
	for _, opt := range opts {
		opt(&log, &data)
	}

	// Pack the data using ABI
	abi, _ := EnygmaPNEvents.EnygmaPNEventsMetaData.ParseABI()
	packedData, _ := abi.Events["EnygmaSwapWithDvpForERC721"].Inputs.NonIndexed().Pack(
		data.resourceID, data.nftId, data.nftResourceId, data.enygmaAmount, data.from, data.destChainId, data.sharedId, data.calldata, data.validityTime,
	)

	log.Topics = []common.Hash{
		crypto.Keccak256Hash([]byte("EnygmaSwapWithDvpForERC721(bytes32,uint256,bytes32,uint256,address,uint256,bytes32,bytes,uint64)")),
	}
	log.Data = packedData

	return log
}

// EnygmaSwapWithDvpForERC1155LogOption represents a modification to an EnygmaSwapWithDvpForERC1155 log fixture.
type EnygmaSwapWithDvpForERC1155LogOption func(*types.Log, *enygmaSwapWithDvpForERC1155Data)

type enygmaSwapWithDvpForERC1155Data struct {
	resourceID     [32]byte
	nftId          *big.Int
	nftResourceId  [32]byte
	nftAmountOrOne *big.Int
	enygmaAmount   *big.Int
	from           common.Address
	destChainId    *big.Int
	sharedId       [32]byte
	calldata       []byte
	validityTime   uint64
}

func WithSwap1155ResourceID(resourceID [32]byte) EnygmaSwapWithDvpForERC1155LogOption {
	return func(log *types.Log, data *enygmaSwapWithDvpForERC1155Data) {
		data.resourceID = resourceID
	}
}

func WithSwap1155NftResourceID(nftResourceId [32]byte) EnygmaSwapWithDvpForERC1155LogOption {
	return func(log *types.Log, data *enygmaSwapWithDvpForERC1155Data) {
		data.nftResourceId = nftResourceId
	}
}

func WithSwap1155SharedID(sharedId [32]byte) EnygmaSwapWithDvpForERC1155LogOption {
	return func(log *types.Log, data *enygmaSwapWithDvpForERC1155Data) {
		data.sharedId = sharedId
	}
}

func WithSwap1155NftId(nftId *big.Int) EnygmaSwapWithDvpForERC1155LogOption {
	return func(log *types.Log, data *enygmaSwapWithDvpForERC1155Data) {
		data.nftId = nftId
	}
}

func WithSwap1155NftAmountOrOne(nftAmountOrOne *big.Int) EnygmaSwapWithDvpForERC1155LogOption {
	return func(log *types.Log, data *enygmaSwapWithDvpForERC1155Data) {
		data.nftAmountOrOne = nftAmountOrOne
	}
}

func WithSwap1155EnygmaAmount(enygmaAmount *big.Int) EnygmaSwapWithDvpForERC1155LogOption {
	return func(log *types.Log, data *enygmaSwapWithDvpForERC1155Data) {
		data.enygmaAmount = enygmaAmount
	}
}

func WithSwap1155From(from common.Address) EnygmaSwapWithDvpForERC1155LogOption {
	return func(log *types.Log, data *enygmaSwapWithDvpForERC1155Data) {
		data.from = from
	}
}

func WithSwap1155DestChainId(destChainId *big.Int) EnygmaSwapWithDvpForERC1155LogOption {
	return func(log *types.Log, data *enygmaSwapWithDvpForERC1155Data) {
		data.destChainId = destChainId
	}
}

func WithSwap1155Calldata(calldata []byte) EnygmaSwapWithDvpForERC1155LogOption {
	return func(log *types.Log, data *enygmaSwapWithDvpForERC1155Data) {
		data.calldata = calldata
	}
}

func WithSwap1155BlockNumber(blockNumber uint64) EnygmaSwapWithDvpForERC1155LogOption {
	return func(log *types.Log, data *enygmaSwapWithDvpForERC1155Data) {
		log.BlockNumber = blockNumber
	}
}

// GetEnygmaSwapWithDvpForERC1155Log returns a sample EnygmaSwapWithDvpForERC1155 event log
func GetEnygmaSwapWithDvpForERC1155Log() types.Log {
	return NewEnygmaSwapWithDvpForERC1155LogWith()
}

// NewEnygmaSwapWithDvpForERC1155LogWith creates an EnygmaSwapWithDvpForERC1155 log using default values
// and then applies any provided options.
func NewEnygmaSwapWithDvpForERC1155LogWith(opts ...EnygmaSwapWithDvpForERC1155LogOption) types.Log {
	// Default values
	data := enygmaSwapWithDvpForERC1155Data{
		resourceID:     [32]byte{0x01, 0x02, 0x03},
		nftId:          big.NewInt(456),
		nftResourceId:  [32]byte{0x0a, 0x0b, 0x0c},
		nftAmountOrOne: big.NewInt(10),
		enygmaAmount:   big.NewInt(750),
		from:           common.HexToAddress("0x5555555555555555555555555555555555555555"),
		destChainId:    big.NewInt(3),
		sharedId:       [32]byte{0x10, 0x11, 0x12},
		calldata:       []byte{0x04, 0x05, 0x06},
		validityTime:   1337,
	}

	log := types.Log{
		Address:     common.HexToAddress("0x1234567890123456789012345678901234567890"),
		BlockNumber: 110,
		TxHash:      common.HexToHash("0xabcdef0a"),
		TxIndex:     11,
		BlockHash:   common.HexToHash("0x123460"),
		Index:       10,
	}

	// Apply options
	for _, opt := range opts {
		opt(&log, &data)
	}

	// Pack the data using ABI
	abi, _ := EnygmaPNEvents.EnygmaPNEventsMetaData.ParseABI()
	packedData, _ := abi.Events["EnygmaSwapWithDvpForERC1155"].Inputs.NonIndexed().Pack(
		data.resourceID, data.nftId, data.nftResourceId, data.nftAmountOrOne, data.enygmaAmount, data.from, data.destChainId, data.sharedId, data.calldata, data.validityTime,
	)

	log.Topics = []common.Hash{
		crypto.Keccak256Hash([]byte("EnygmaSwapWithDvpForERC1155(bytes32,uint256,bytes32,uint256,uint256,address,uint256,bytes32,bytes,uint64)")),
	}
	log.Data = packedData

	return log
}

// GetDvp721DepositIntoDvpLog returns a sample Dvp721DepositIntoDvp event log
func GetDvp721DepositIntoDvpLog() types.Log {
	resourceId := [32]byte{0x0a, 0x0b, 0x0c}
	nftId := big.NewInt(789)
	from := common.HexToAddress("0x6666666666666666666666666666666666666666")

	abi, _ := EnygmaPNEvents.EnygmaPNEventsMetaData.ParseABI()
	data, _ := abi.Events["Dvp721DepositIntoDvp"].Inputs.NonIndexed().Pack(resourceId, nftId, from)

	return types.Log{
		Address: common.HexToAddress("0x1234567890123456789012345678901234567890"),
		Topics: []common.Hash{
			crypto.Keccak256Hash([]byte("Dvp721DepositIntoDvp(bytes32,uint256,address)")),
		},
		Data:        data,
		BlockNumber: 111,
		TxHash:      common.HexToHash("0xabcdef0b"),
		TxIndex:     12,
		BlockHash:   common.HexToHash("0x123461"),
		Index:       11,
	}
}

// GetDvp721WithdrawFromDvpLog returns a sample Dvp721WithdrawFromDvp event log
func GetDvp721WithdrawFromDvpLog() types.Log {
	resourceId := [32]byte{0x0a, 0x0b, 0x0c}
	nftId := big.NewInt(789)
	owner := common.HexToAddress("0x7777777777777777777777777777777777777777")

	abi, _ := EnygmaPNEvents.EnygmaPNEventsMetaData.ParseABI()
	data, _ := abi.Events["Dvp721WithdrawFromDvp"].Inputs.NonIndexed().Pack(resourceId, nftId, owner)

	return types.Log{
		Address: common.HexToAddress("0x1234567890123456789012345678901234567890"),
		Topics: []common.Hash{
			crypto.Keccak256Hash([]byte("Dvp721WithdrawFromDvp(bytes32,uint256,address)")),
		},
		Data:        data,
		BlockNumber: 112,
		TxHash:      common.HexToHash("0xabcdef0c"),
		TxIndex:     13,
		BlockHash:   common.HexToHash("0x123462"),
		Index:       12,
	}
}

// GetDvp721SwapForEnygmaLog returns a sample Dvp721SwapForEnygma event log
func GetDvp721SwapForEnygmaLog() types.Log {
	nftResourceId := [32]byte{0x0a, 0x0b, 0x0c}
	nftId := big.NewInt(999)
	enygmaAmount := big.NewInt(1000)
	enygmaResourceId := [32]byte{0x01, 0x02, 0x03}
	from := common.HexToAddress("0x8888888888888888888888888888888888888888")
	destChainId := big.NewInt(4)
	sharedId := [32]byte{0x13, 0x14, 0x15}
	calldata := []byte{0x07, 0x08, 0x09}
	validityTime := uint64(1337)

	abi, _ := EnygmaPNEvents.EnygmaPNEventsMetaData.ParseABI()
	data, _ := abi.Events["Dvp721SwapForEnygma"].Inputs.NonIndexed().Pack(
		nftResourceId, nftId, enygmaAmount, enygmaResourceId, from, destChainId, sharedId, calldata, validityTime,
	)

	return types.Log{
		Address: common.HexToAddress("0x1234567890123456789012345678901234567890"),
		Topics: []common.Hash{
			crypto.Keccak256Hash([]byte("Dvp721SwapForEnygma(bytes32,uint256,uint256,bytes32,address,uint256,bytes32,bytes,uint64)")),
		},
		Data:        data,
		BlockNumber: 113,
		TxHash:      common.HexToHash("0xabcdef0d"),
		TxIndex:     14,
		BlockHash:   common.HexToHash("0x123463"),
		Index:       13,
	}
}

// GetDvp1155CreationLog returns a sample Dvp1155Creation event log
func GetDvp1155CreationLog() types.Log {
	resourceId := [32]byte{0x0d, 0x0e, 0x0f}

	abi, _ := EnygmaPNEvents.EnygmaPNEventsMetaData.ParseABI()
	data, _ := abi.Events["Dvp1155Creation"].Inputs.NonIndexed().Pack(resourceId)

	return types.Log{
		Address: common.HexToAddress("0x1234567890123456789012345678901234567890"),
		Topics: []common.Hash{
			crypto.Keccak256Hash([]byte("Dvp1155Creation(bytes32)")),
		},
		Data:        data,
		BlockNumber: 114,
		TxHash:      common.HexToHash("0xabcdef0e"),
		TxIndex:     15,
		BlockHash:   common.HexToHash("0x123464"),
		Index:       14,
	}
}

// GetDvp1155MintLog returns a sample Dvp1155Mint event log
func GetDvp1155MintLog() types.Log {
	resourceId := [32]byte{0x0d, 0x0e, 0x0f}
	tokenId := big.NewInt(111)
	value := big.NewInt(50)
	data := []byte{0x0a, 0x0b}

	abi, _ := EnygmaPNEvents.EnygmaPNEventsMetaData.ParseABI()
	packedData, _ := abi.Events["Dvp1155Mint"].Inputs.NonIndexed().Pack(resourceId, tokenId, value, data)

	return types.Log{
		Address: common.HexToAddress("0x1234567890123456789012345678901234567890"),
		Topics: []common.Hash{
			crypto.Keccak256Hash([]byte("Dvp1155Mint(bytes32,uint256,uint256,bytes)")),
		},
		Data:        packedData,
		BlockNumber: 115,
		TxHash:      common.HexToHash("0xabcdef0f"),
		TxIndex:     16,
		BlockHash:   common.HexToHash("0x123465"),
		Index:       15,
	}
}

// GetDvp1155BurnLog returns a sample Dvp1155Burn event log
func GetDvp1155BurnLog() types.Log {
	resourceId := [32]byte{0x0d, 0x0e, 0x0f}
	to := common.HexToAddress("0x9999999999999999999999999999999999999999")
	tokenId := big.NewInt(111)
	value := big.NewInt(25)

	abi, _ := EnygmaPNEvents.EnygmaPNEventsMetaData.ParseABI()
	data, _ := abi.Events["Dvp1155Burn"].Inputs.NonIndexed().Pack(resourceId, to, tokenId, value)

	return types.Log{
		Address: common.HexToAddress("0x1234567890123456789012345678901234567890"),
		Topics: []common.Hash{
			crypto.Keccak256Hash([]byte("Dvp1155Burn(bytes32,address,uint256,uint256)")),
		},
		Data:        data,
		BlockNumber: 116,
		TxHash:      common.HexToHash("0xabcdef10"),
		TxIndex:     17,
		BlockHash:   common.HexToHash("0x123466"),
		Index:       16,
	}
}

// GetDvp1155DepositIntoDvpLog returns a sample Dvp1155DepositIntoDvp event log
func GetDvp1155DepositIntoDvpLog() types.Log {
	resourceId := [32]byte{0x0d, 0x0e, 0x0f}
	tokenId := big.NewInt(222)
	from := common.HexToAddress("0xaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa")
	value := big.NewInt(75)
	data := []byte{0x0c, 0x0d}

	abi, _ := EnygmaPNEvents.EnygmaPNEventsMetaData.ParseABI()
	packedData, _ := abi.Events["Dvp1155DepositIntoDvp"].Inputs.NonIndexed().Pack(resourceId, tokenId, from, value, data)

	return types.Log{
		Address: common.HexToAddress("0x1234567890123456789012345678901234567890"),
		Topics: []common.Hash{
			crypto.Keccak256Hash([]byte("Dvp1155DepositIntoDvp(bytes32,uint256,address,uint256,bytes)")),
		},
		Data:        packedData,
		BlockNumber: 117,
		TxHash:      common.HexToHash("0xabcdef11"),
		TxIndex:     18,
		BlockHash:   common.HexToHash("0x123467"),
		Index:       17,
	}
}
