package testdata

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/raylsnetwork/rayls-privacy-relayer-api/contracts/DvpTeleport"
)

// TransferEncryptedDataLogOption represents a modification to a TransferEncryptedData log fixture.
type TransferEncryptedDataLogOption func(*types.Log, *transferEncryptedDataData)

type transferEncryptedDataData struct {
	print       string
	data        []byte
	blockNumber *big.Int
}

// WithDvpPrint sets the Print field.
func WithDvpPrint(print string) TransferEncryptedDataLogOption {
	return func(log *types.Log, data *transferEncryptedDataData) {
		data.print = print
	}
}

// WithDvpData sets the Data field.
func WithDvpData(d []byte) TransferEncryptedDataLogOption {
	return func(log *types.Log, data *transferEncryptedDataData) {
		data.data = d
	}
}

// WithDvpBlockNumber sets the BlockNumber field in the event data.
func WithDvpBlockNumber(blockNumber *big.Int) TransferEncryptedDataLogOption {
	return func(log *types.Log, data *transferEncryptedDataData) {
		data.blockNumber = blockNumber
	}
}

// WithDvpLogBlockNumber sets the log's BlockNumber.
func WithDvpLogBlockNumber(blockNumber uint64) TransferEncryptedDataLogOption {
	return func(log *types.Log, data *transferEncryptedDataData) {
		log.BlockNumber = blockNumber
	}
}

// WithDvpTxHash sets the TxHash field.
func WithDvpTxHash(hash common.Hash) TransferEncryptedDataLogOption {
	return func(log *types.Log, data *transferEncryptedDataData) {
		log.TxHash = hash
	}
}

// NewTransferEncryptedDataLogWith creates a TransferEncryptedData log with default values and applies options.
func NewTransferEncryptedDataLogWith(opts ...TransferEncryptedDataLogOption) types.Log {
	data := transferEncryptedDataData{
		print:       "dvp-transfer",
		data:        []byte{0xde, 0xad, 0xbe, 0xef},
		blockNumber: big.NewInt(100),
	}

	log := types.Log{
		Address:     common.HexToAddress("0x1234567890123456789012345678901234567890"),
		BlockNumber: 100,
		TxHash:      common.HexToHash("0xabcdef1234567890"),
		TxIndex:     1,
		BlockHash:   common.HexToHash("0x123456"),
		Index:       0,
	}

	for _, opt := range opts {
		opt(&log, &data)
	}

	abi, _ := DvpTeleport.DvpTeleportMetaData.ParseABI()
	packedData, _ := abi.Events["TransferEncryptedData"].Inputs.NonIndexed().Pack(
		data.print,
		data.data,
	)

	// Event: TransferEncryptedData(string print, bytes data, uint256 indexed blockNumber)
	// Event signature hash: 0x43c788ab3309319c9670ae5f123320c0ad30c5a96dc78a9b9bfccaf187242a16
	log.Topics = []common.Hash{
		common.HexToHash("0x43c788ab3309319c9670ae5f123320c0ad30c5a96dc78a9b9bfccaf187242a16"),
		common.BigToHash(data.blockNumber),
	}
	log.Data = packedData

	return log
}

// InitiateCalldataLogOption represents a modification to an InitiateCalldata log fixture.
type InitiateCalldataLogOption func(*types.Log, *initiateCalldataData)

type initiateCalldataData struct {
	sharedId [32]byte
}

// WithInitiateCalldataSharedId sets the SharedId field.
func WithInitiateCalldataSharedId(sharedId [32]byte) InitiateCalldataLogOption {
	return func(log *types.Log, data *initiateCalldataData) {
		data.sharedId = sharedId
	}
}

// WithInitiateCalldataLogBlockNumber sets the log's BlockNumber.
func WithInitiateCalldataLogBlockNumber(blockNumber uint64) InitiateCalldataLogOption {
	return func(log *types.Log, data *initiateCalldataData) {
		log.BlockNumber = blockNumber
	}
}

// NewInitiateCalldataLogWith creates an InitiateCalldata log with default values and applies options.
func NewInitiateCalldataLogWith(opts ...InitiateCalldataLogOption) types.Log {
	data := initiateCalldataData{
		sharedId: [32]byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08},
	}

	log := types.Log{
		Address:     common.HexToAddress("0x1234567890123456789012345678901234567890"),
		BlockNumber: 100,
		TxHash:      common.HexToHash("0xabcdef1234567890"),
		TxIndex:     1,
		BlockHash:   common.HexToHash("0x123456"),
		Index:       0,
	}

	for _, opt := range opts {
		opt(&log, &data)
	}

	// Event: InitiateCalldata(bytes32 indexed sharedId)
	// Event signature hash: 0x7390f63f046b620dbeb3de874ca62b42b9f821f6faf449f6324d0375808aec10
	log.Topics = []common.Hash{
		common.HexToHash("0x7390f63f046b620dbeb3de874ca62b42b9f821f6faf449f6324d0375808aec10"),
		common.BytesToHash(data.sharedId[:]),
	}
	log.Data = nil

	return log
}

// SwapStateChangedLogOption represents a modification to a SwapStateChanged log fixture.
type SwapStateChangedLogOption func(*types.Log, *swapStateChangedData)

type swapStateChangedData struct {
	sharedId [32]byte
	state    uint8
}

// WithSwapStateChangedSharedId sets the SharedId field.
func WithSwapStateChangedSharedId(sharedId [32]byte) SwapStateChangedLogOption {
	return func(log *types.Log, data *swapStateChangedData) {
		data.sharedId = sharedId
	}
}

// WithSwapStateChangedState sets the State field.
func WithSwapStateChangedState(state uint8) SwapStateChangedLogOption {
	return func(log *types.Log, data *swapStateChangedData) {
		data.state = state
	}
}

// WithSwapStateChangedLogBlockNumber sets the log's BlockNumber.
func WithSwapStateChangedLogBlockNumber(blockNumber uint64) SwapStateChangedLogOption {
	return func(log *types.Log, data *swapStateChangedData) {
		log.BlockNumber = blockNumber
	}
}

// NewSwapStateChangedLogWith creates a SwapStateChanged log with default values and applies options.
func NewSwapStateChangedLogWith(opts ...SwapStateChangedLogOption) types.Log {
	data := swapStateChangedData{
		sharedId: [32]byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08},
		state:    0,
	}

	log := types.Log{
		Address:     common.HexToAddress("0x1234567890123456789012345678901234567890"),
		BlockNumber: 100,
		TxHash:      common.HexToHash("0xabcdef1234567890"),
		TxIndex:     1,
		BlockHash:   common.HexToHash("0x123456"),
		Index:       0,
	}

	for _, opt := range opts {
		opt(&log, &data)
	}

	abi, _ := DvpTeleport.DvpTeleportMetaData.ParseABI()
	packedData, _ := abi.Events["SwapStateChanged"].Inputs.NonIndexed().Pack(
		data.state,
	)

	// Event: SwapStateChanged(bytes32 indexed sharedId, uint8 state)
	// Event signature hash: 0x1a29a2466d037dcb4e1cec4abdb8c609191cebb6093a9a341e412a08d2987188
	log.Topics = []common.Hash{
		common.HexToHash("0x1a29a2466d037dcb4e1cec4abdb8c609191cebb6093a9a341e412a08d2987188"),
		common.BytesToHash(data.sharedId[:]),
	}
	log.Data = packedData

	return log
}

// CalldataExecutedLogOption represents a modification to a CalldataExecuted log fixture.
type CalldataExecutedLogOption func(*types.Log, *calldataExecutedData)

type calldataExecutedData struct {
	sharedId [32]byte
}

// WithCalldataExecutedSharedId sets the SharedId field.
func WithCalldataExecutedSharedId(sharedId [32]byte) CalldataExecutedLogOption {
	return func(log *types.Log, data *calldataExecutedData) {
		data.sharedId = sharedId
	}
}

// WithCalldataExecutedLogBlockNumber sets the log's BlockNumber.
func WithCalldataExecutedLogBlockNumber(blockNumber uint64) CalldataExecutedLogOption {
	return func(log *types.Log, data *calldataExecutedData) {
		log.BlockNumber = blockNumber
	}
}

// NewCalldataExecutedLogWith creates a CalldataExecuted log with default values and applies options.
func NewCalldataExecutedLogWith(opts ...CalldataExecutedLogOption) types.Log {
	data := calldataExecutedData{
		sharedId: [32]byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08},
	}

	log := types.Log{
		Address:     common.HexToAddress("0x1234567890123456789012345678901234567890"),
		BlockNumber: 100,
		TxHash:      common.HexToHash("0xabcdef1234567890"),
		TxIndex:     1,
		BlockHash:   common.HexToHash("0x123456"),
		Index:       0,
	}

	for _, opt := range opts {
		opt(&log, &data)
	}

	// Event: CalldataExecuted(bytes32 indexed sharedId)
	// Event signature hash: 0xfa731f14885e8268a85e7e629527a1146a0a1fc5a8684fd35df634af45f642fc
	log.Topics = []common.Hash{
		common.HexToHash("0xfa731f14885e8268a85e7e629527a1146a0a1fc5a8684fd35df634af45f642fc"),
		common.BytesToHash(data.sharedId[:]),
	}
	log.Data = nil

	return log
}

// CommitmentsLogOption represents a modification to a Commitments log fixture.
type CommitmentsLogOption func(*types.Log, *commitmentsData)

type commitmentsData struct {
	tokenAddress common.Address
	tokenType    *big.Int
	treeNumber   *big.Int
	commitments  []*big.Int
}

// WithCommitmentsTokenAddress sets the TokenAddress field.
func WithCommitmentsTokenAddress(addr common.Address) CommitmentsLogOption {
	return func(log *types.Log, data *commitmentsData) {
		data.tokenAddress = addr
	}
}

// WithCommitmentsTokenType sets the TokenType field.
func WithCommitmentsTokenType(tokenType *big.Int) CommitmentsLogOption {
	return func(log *types.Log, data *commitmentsData) {
		data.tokenType = tokenType
	}
}

// WithCommitmentsTreeNumber sets the TreeNumber field.
func WithCommitmentsTreeNumber(treeNumber *big.Int) CommitmentsLogOption {
	return func(log *types.Log, data *commitmentsData) {
		data.treeNumber = treeNumber
	}
}

// WithCommitmentsValues sets the Commitments field.
func WithCommitmentsValues(commitments []*big.Int) CommitmentsLogOption {
	return func(log *types.Log, data *commitmentsData) {
		data.commitments = commitments
	}
}

// WithCommitmentsLogBlockNumber sets the log's BlockNumber.
func WithCommitmentsLogBlockNumber(blockNumber uint64) CommitmentsLogOption {
	return func(log *types.Log, data *commitmentsData) {
		log.BlockNumber = blockNumber
	}
}

// WithCommitmentsTxHash sets the TxHash field.
func WithCommitmentsTxHash(hash common.Hash) CommitmentsLogOption {
	return func(log *types.Log, data *commitmentsData) {
		log.TxHash = hash
	}
}

// WithCommitmentsLogIndex sets the log's Index field.
func WithCommitmentsLogIndex(index uint) CommitmentsLogOption {
	return func(log *types.Log, data *commitmentsData) {
		log.Index = index
	}
}

// NewCommitmentsLogWith creates a Commitments log with default values and applies options.
func NewCommitmentsLogWith(opts ...CommitmentsLogOption) types.Log {
	data := commitmentsData{
		tokenAddress: common.HexToAddress("0xaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"),
		tokenType:    big.NewInt(4),
		treeNumber:   big.NewInt(1),
		commitments:  []*big.Int{big.NewInt(111), big.NewInt(222)},
	}

	log := types.Log{
		Address:     common.HexToAddress("0x1234567890123456789012345678901234567890"),
		BlockNumber: 100,
		TxHash:      common.HexToHash("0xabcdef1234567890"),
		TxIndex:     1,
		BlockHash:   common.HexToHash("0x123456"),
		Index:       0,
	}

	for _, opt := range opts {
		opt(&log, &data)
	}

	abi, _ := DvpTeleport.DvpTeleportMetaData.ParseABI()
	packedData, _ := abi.Events["Commitments"].Inputs.NonIndexed().Pack(
		data.tokenType,
		data.commitments,
	)

	// Event: Commitments(address indexed tokenAddress, uint256 tokenType, uint256 indexed treeNumber, uint256[] commitments)
	// Event signature hash: 0x586ab18efd7bdd12e2b2064a7a74c9fabbb8c5a98d483ec5ec25899388231cc4
	log.Topics = []common.Hash{
		common.HexToHash("0x586ab18efd7bdd12e2b2064a7a74c9fabbb8c5a98d483ec5ec25899388231cc4"),
		common.BytesToHash(data.tokenAddress.Bytes()),
		common.BigToHash(data.treeNumber),
	}
	log.Data = packedData

	return log
}

// SwapCompletedLogOption represents a modification to a SwapCompleted log fixture.
type SwapCompletedLogOption func(*types.Log, *swapCompletedData)

type swapCompletedData struct {
	sharedId      [32]byte
	encryptedData []byte
}

// WithSwapCompletedSharedId sets the SharedId field.
func WithSwapCompletedSharedId(sharedId [32]byte) SwapCompletedLogOption {
	return func(log *types.Log, data *swapCompletedData) {
		data.sharedId = sharedId
	}
}

// WithSwapCompletedEncryptedData sets the EncryptedData field.
func WithSwapCompletedEncryptedData(encryptedData []byte) SwapCompletedLogOption {
	return func(log *types.Log, data *swapCompletedData) {
		data.encryptedData = encryptedData
	}
}

// WithSwapCompletedLogBlockNumber sets the log's BlockNumber.
func WithSwapCompletedLogBlockNumber(blockNumber uint64) SwapCompletedLogOption {
	return func(log *types.Log, data *swapCompletedData) {
		log.BlockNumber = blockNumber
	}
}

// WithSwapCompletedTxHash sets the TxHash field.
func WithSwapCompletedTxHash(hash common.Hash) SwapCompletedLogOption {
	return func(log *types.Log, data *swapCompletedData) {
		log.TxHash = hash
	}
}

// NewSwapCompletedLogWith creates a SwapCompleted log with default values and applies options.
func NewSwapCompletedLogWith(opts ...SwapCompletedLogOption) types.Log {
	data := swapCompletedData{
		sharedId:      [32]byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08},
		encryptedData: []byte{0xc0, 0xff, 0xee},
	}

	log := types.Log{
		Address:     common.HexToAddress("0x1234567890123456789012345678901234567890"),
		BlockNumber: 100,
		TxHash:      common.HexToHash("0xabcdef1234567890"),
		TxIndex:     1,
		BlockHash:   common.HexToHash("0x123456"),
		Index:       0,
	}

	for _, opt := range opts {
		opt(&log, &data)
	}

	abi, _ := DvpTeleport.DvpTeleportMetaData.ParseABI()
	packedData, _ := abi.Events["SwapCompleted"].Inputs.NonIndexed().Pack(
		data.encryptedData,
	)

	// Event: SwapCompleted(bytes32 indexed sharedId, bytes encryptedData)
	// Event signature hash: 0xc8c1414976453b7f8359322b6c644cbb667a30bab1ae7bc95e7b856d17d5ca8b
	log.Topics = []common.Hash{
		common.HexToHash("0xc8c1414976453b7f8359322b6c644cbb667a30bab1ae7bc95e7b856d17d5ca8b"),
		common.BytesToHash(data.sharedId[:]),
	}
	log.Data = packedData

	return log
}

// NullifiersLogOption represents a modification to a Nullifiers log fixture.
type NullifiersLogOption func(*types.Log, *nullifiersData)

type nullifiersData struct {
	tokenAddress common.Address
	tokenType    *big.Int
	nullifiers   []*big.Int
}

// WithNullifiersTokenAddress sets the TokenAddress field.
func WithNullifiersTokenAddress(addr common.Address) NullifiersLogOption {
	return func(log *types.Log, data *nullifiersData) {
		data.tokenAddress = addr
	}
}

// WithNullifiersTokenType sets the TokenType field.
func WithNullifiersTokenType(tokenType *big.Int) NullifiersLogOption {
	return func(log *types.Log, data *nullifiersData) {
		data.tokenType = tokenType
	}
}

// WithNullifiersValues sets the Nullifiers field.
func WithNullifiersValues(nullifiers []*big.Int) NullifiersLogOption {
	return func(log *types.Log, data *nullifiersData) {
		data.nullifiers = nullifiers
	}
}

// WithNullifiersLogBlockNumber sets the log's BlockNumber.
func WithNullifiersLogBlockNumber(blockNumber uint64) NullifiersLogOption {
	return func(log *types.Log, data *nullifiersData) {
		log.BlockNumber = blockNumber
	}
}

// WithNullifiersTxHash sets the TxHash field.
func WithNullifiersTxHash(hash common.Hash) NullifiersLogOption {
	return func(log *types.Log, data *nullifiersData) {
		log.TxHash = hash
	}
}

// NewNullifiersLogWith creates a Nullifiers log with default values and applies options.
func NewNullifiersLogWith(opts ...NullifiersLogOption) types.Log {
	data := nullifiersData{
		tokenAddress: common.HexToAddress("0xaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"),
		tokenType:    big.NewInt(4),
		nullifiers:   []*big.Int{big.NewInt(99999)},
	}

	log := types.Log{
		Address:     common.HexToAddress("0x1234567890123456789012345678901234567890"),
		BlockNumber: 100,
		TxHash:      common.HexToHash("0xabcdef1234567890"),
		TxIndex:     1,
		BlockHash:   common.HexToHash("0x123456"),
		Index:       0,
	}

	for _, opt := range opts {
		opt(&log, &data)
	}

	abi, _ := DvpTeleport.DvpTeleportMetaData.ParseABI()
	packedData, _ := abi.Events["Nullifiers"].Inputs.NonIndexed().Pack(
		data.tokenType,
		data.nullifiers,
	)

	// Event: Nullifiers(address indexed tokenAddress, uint256 tokenType, uint256[] nullifiers)
	// Event signature hash: 0x259f270c5a1d564ded25ccac004efc36cde097fc3840157651761576f40e9835
	log.Topics = []common.Hash{
		common.HexToHash("0x259f270c5a1d564ded25ccac004efc36cde097fc3840157651761576f40e9835"),
		common.BytesToHash(data.tokenAddress.Bytes()),
	}
	log.Data = packedData

	return log
}
