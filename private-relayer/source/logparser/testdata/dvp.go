package testdata

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"

	"github.com/raylsnetwork/rayls-sovereign-relayer/contracts/EnygmaPNEvents"
)

// Dvp721CreationLogOption represents a modification to a Dvp721Creation log fixture.
type Dvp721CreationLogOption func(*types.Log, *zkDvp721CreationData)

type zkDvp721CreationData struct {
	resourceID [32]byte
}

// WithDvp721CreationResourceID sets the ResourceID field on a Dvp721Creation log when using options.
func WithDvp721CreationResourceID(resourceID [32]byte) Dvp721CreationLogOption {
	return func(log *types.Log, data *zkDvp721CreationData) {
		data.resourceID = resourceID
	}
}

// WithDvp721CreationBlockNumber sets the BlockNumber field on a Dvp721Creation log when using options.
func WithDvp721CreationBlockNumber(blockNumber uint64) Dvp721CreationLogOption {
	return func(log *types.Log, data *zkDvp721CreationData) {
		log.BlockNumber = blockNumber
	}
}

// NewDvp721CreationLogWith creates a Dvp721Creation log using default values
// and then applies any provided options.
func NewDvp721CreationLogWith(opts ...Dvp721CreationLogOption) types.Log {
	// Default values
	data := zkDvp721CreationData{
		resourceID: [32]byte{0x0a, 0x0b, 0x0c},
	}

	log := types.Log{
		Address:     common.HexToAddress("0x1234567890123456789012345678901234567890"),
		BlockNumber: 105,
		TxHash:      common.HexToHash("0xabcdef05"),
		TxIndex:     6,
		BlockHash:   common.HexToHash("0x12345b"),
		Index:       5,
	}

	// Apply options
	for _, opt := range opts {
		opt(&log, &data)
	}

	// Pack the data using ABI
	abi, _ := EnygmaPNEvents.EnygmaPNEventsMetaData.ParseABI()
	packedData, _ := abi.Events["Dvp721Creation"].Inputs.NonIndexed().Pack(data.resourceID)

	log.Topics = []common.Hash{
		crypto.Keccak256Hash([]byte("Dvp721Creation(bytes32)")),
	}
	log.Data = packedData

	return log
}

// Dvp721MintLogOption represents a modification to a Dvp721Mint log fixture.
type Dvp721MintLogOption func(*types.Log, *zkDvp721MintData)

type zkDvp721MintData struct {
	resourceID [32]byte
	nftId      *big.Int
}

// WithDvp721MintResourceID sets the ResourceID field on a Dvp721Mint log when using options.
func WithDvp721MintResourceID(resourceID [32]byte) Dvp721MintLogOption {
	return func(log *types.Log, data *zkDvp721MintData) {
		data.resourceID = resourceID
	}
}

// WithDvp721MintNftId sets the NftId field on a Dvp721Mint log when using options.
func WithDvp721MintNftId(nftId *big.Int) Dvp721MintLogOption {
	return func(log *types.Log, data *zkDvp721MintData) {
		data.nftId = nftId
	}
}

// WithDvp721MintBlockNumber sets the BlockNumber field on a Dvp721Mint log when using options.
func WithDvp721MintBlockNumber(blockNumber uint64) Dvp721MintLogOption {
	return func(log *types.Log, data *zkDvp721MintData) {
		log.BlockNumber = blockNumber
	}
}

// NewDvp721MintLogWith creates a Dvp721Mint log using default values
// and then applies any provided options.
func NewDvp721MintLogWith(opts ...Dvp721MintLogOption) types.Log {
	// Default values
	data := zkDvp721MintData{
		resourceID: [32]byte{0x0a, 0x0b, 0x0c},
		nftId:      big.NewInt(42),
	}

	log := types.Log{
		Address:     common.HexToAddress("0x1234567890123456789012345678901234567890"),
		BlockNumber: 106,
		TxHash:      common.HexToHash("0xabcdef06"),
		TxIndex:     7,
		BlockHash:   common.HexToHash("0x12345c"),
		Index:       6,
	}

	// Apply options
	for _, opt := range opts {
		opt(&log, &data)
	}

	// Pack the data using ABI
	abi, _ := EnygmaPNEvents.EnygmaPNEventsMetaData.ParseABI()
	packedData, _ := abi.Events["Dvp721Mint"].Inputs.NonIndexed().Pack(data.resourceID, data.nftId)

	log.Topics = []common.Hash{
		crypto.Keccak256Hash([]byte("Dvp721Mint(bytes32,uint256)")),
	}
	log.Data = packedData

	return log
}

// Dvp721BurnLogOption represents a modification to a Dvp721Burn log fixture.
type Dvp721BurnLogOption func(*types.Log, *zkDvp721BurnData)

type zkDvp721BurnData struct {
	resourceID [32]byte
	nftId      *big.Int
}

// WithDvp721BurnResourceID sets the ResourceID field on a Dvp721Burn log when using options.
func WithDvp721BurnResourceID(resourceID [32]byte) Dvp721BurnLogOption {
	return func(log *types.Log, data *zkDvp721BurnData) {
		data.resourceID = resourceID
	}
}

// WithDvp721BurnNftId sets the NftId field on a Dvp721Burn log when using options.
func WithDvp721BurnNftId(nftId *big.Int) Dvp721BurnLogOption {
	return func(log *types.Log, data *zkDvp721BurnData) {
		data.nftId = nftId
	}
}

// WithDvp721BurnBlockNumber sets the BlockNumber field on a Dvp721Burn log when using options.
func WithDvp721BurnBlockNumber(blockNumber uint64) Dvp721BurnLogOption {
	return func(log *types.Log, data *zkDvp721BurnData) {
		log.BlockNumber = blockNumber
	}
}

// NewDvp721BurnLogWith creates a Dvp721Burn log using default values
// and then applies any provided options.
func NewDvp721BurnLogWith(opts ...Dvp721BurnLogOption) types.Log {
	// Default values
	data := zkDvp721BurnData{
		resourceID: [32]byte{0x0a, 0x0b, 0x0c},
		nftId:      big.NewInt(42),
	}

	log := types.Log{
		Address:     common.HexToAddress("0x1234567890123456789012345678901234567890"),
		BlockNumber: 107,
		TxHash:      common.HexToHash("0xabcdef07"),
		TxIndex:     8,
		BlockHash:   common.HexToHash("0x12345d"),
		Index:       7,
	}

	// Apply options
	for _, opt := range opts {
		opt(&log, &data)
	}

	// Pack the data using ABI
	abi, _ := EnygmaPNEvents.EnygmaPNEventsMetaData.ParseABI()
	packedData, _ := abi.Events["Dvp721Burn"].Inputs.NonIndexed().Pack(data.resourceID, data.nftId)

	log.Topics = []common.Hash{
		crypto.Keccak256Hash([]byte("Dvp721Burn(bytes32,uint256)")),
	}
	log.Data = packedData

	return log
}

// Dvp721DepositIntoDvp log constructor

type dvp721DepositIntoDvpData struct {
	resourceID [32]byte
	nftId      *big.Int
	from       common.Address
}

type Dvp721DepositIntoDvpLogOption func(*types.Log, *dvp721DepositIntoDvpData)

// WithDvp721DepositIntoDvpResourceID sets the ResourceID field on a Dvp721DepositIntoDvp log when using options.
func WithDvp721DepositIntoDvpResourceID(resourceID common.Hash) Dvp721DepositIntoDvpLogOption {
	return func(log *types.Log, data *dvp721DepositIntoDvpData) {
		copy(data.resourceID[:], resourceID[:])
	}
}

// WithDvp721DepositIntoDvpNftId sets the NftId field on a Dvp721DepositIntoDvp log when using options.
func WithDvp721DepositIntoDvpNftId(nftId *big.Int) Dvp721DepositIntoDvpLogOption {
	return func(log *types.Log, data *dvp721DepositIntoDvpData) {
		data.nftId = nftId
	}
}

// WithDvp721DepositIntoDvpFrom sets the From field on a Dvp721DepositIntoDvp log when using options.
func WithDvp721DepositIntoDvpFrom(from common.Address) Dvp721DepositIntoDvpLogOption {
	return func(log *types.Log, data *dvp721DepositIntoDvpData) {
		data.from = from
	}
}

// WithDvp721DepositIntoDvpBlockNumber sets the BlockNumber field on a Dvp721DepositIntoDvp log when using options.
func WithDvp721DepositIntoDvpBlockNumber(blockNumber uint64) Dvp721DepositIntoDvpLogOption {
	return func(log *types.Log, data *dvp721DepositIntoDvpData) {
		log.BlockNumber = blockNumber
	}
}

// NewDvp721DepositIntoDvpLogWith creates a Dvp721DepositIntoDvp log using default values
// and then applies any provided options.
func NewDvp721DepositIntoDvpLogWith(opts ...Dvp721DepositIntoDvpLogOption) types.Log {
	// Default values
	data := dvp721DepositIntoDvpData{
		resourceID: [32]byte{0x0a, 0x0b, 0x0c},
		nftId:      big.NewInt(42),
		from:       common.HexToAddress("0xABCDEF1234567890ABCDEF1234567890ABCDEF12"),
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
	packedData, _ := abi.Events["Dvp721DepositIntoDvp"].Inputs.NonIndexed().Pack(data.resourceID, data.nftId, data.from)

	log.Topics = []common.Hash{
		crypto.Keccak256Hash([]byte("Dvp721DepositIntoDvp(bytes32,uint256,address)")),
	}
	log.Data = packedData

	return log
}

// Dvp721WithdrawFromDvp log constructor

type dvp721WithdrawFromDvpData struct {
	resourceID [32]byte
	nftId      *big.Int
	owner      common.Address
}

type Dvp721WithdrawFromDvpLogOption func(*types.Log, *dvp721WithdrawFromDvpData)

// WithDvp721WithdrawFromDvpResourceID sets the ResourceID field on a Dvp721WithdrawFromDvp log when using options.
func WithDvp721WithdrawFromDvpResourceID(resourceID common.Hash) Dvp721WithdrawFromDvpLogOption {
	return func(log *types.Log, data *dvp721WithdrawFromDvpData) {
		copy(data.resourceID[:], resourceID[:])
	}
}

// WithDvp721WithdrawFromDvpNftId sets the NftId field on a Dvp721WithdrawFromDvp log when using options.
func WithDvp721WithdrawFromDvpNftId(nftId *big.Int) Dvp721WithdrawFromDvpLogOption {
	return func(log *types.Log, data *dvp721WithdrawFromDvpData) {
		data.nftId = nftId
	}
}

// WithDvp721WithdrawFromDvpOwner sets the Owner field on a Dvp721WithdrawFromDvp log when using options.
func WithDvp721WithdrawFromDvpOwner(owner common.Address) Dvp721WithdrawFromDvpLogOption {
	return func(log *types.Log, data *dvp721WithdrawFromDvpData) {
		data.owner = owner
	}
}

// WithDvp721WithdrawFromDvpBlockNumber sets the BlockNumber field on a Dvp721WithdrawFromDvp log when using options.
func WithDvp721WithdrawFromDvpBlockNumber(blockNumber uint64) Dvp721WithdrawFromDvpLogOption {
	return func(log *types.Log, data *dvp721WithdrawFromDvpData) {
		log.BlockNumber = blockNumber
	}
}

// NewDvp721WithdrawFromDvpLogWith creates a Dvp721WithdrawFromDvp log using default values
// and then applies any provided options.
func NewDvp721WithdrawFromDvpLogWith(opts ...Dvp721WithdrawFromDvpLogOption) types.Log {
	// Default values
	data := dvp721WithdrawFromDvpData{
		resourceID: [32]byte{0x0a, 0x0b, 0x0c},
		nftId:      big.NewInt(42),
		owner:      common.HexToAddress("0xABCDEF1234567890ABCDEF1234567890ABCDEF12"),
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
	packedData, _ := abi.Events["Dvp721WithdrawFromDvp"].Inputs.NonIndexed().Pack(data.resourceID, data.nftId, data.owner)

	log.Topics = []common.Hash{
		crypto.Keccak256Hash([]byte("Dvp721WithdrawFromDvp(bytes32,uint256,address)")),
	}
	log.Data = packedData

	return log
}

// Remaining Dvp event log constructors with minimal options (2-3 fields each for identification)

// Dvp721SwapForEnygma - ResourceID and NftId
type zkDvp721SwapForEnygmaData struct {
	sharedId         [32]byte
	destChainId      *big.Int
	from             common.Address
	nftResourceId    [32]byte
	nftId            *big.Int
	enygmaResourceId [32]byte
	enygmaAmount     *big.Int
	validityTime     uint64
}

type Dvp721SwapForEnygmaLogOption func(*types.Log, *zkDvp721SwapForEnygmaData)

func WithDvp721SwapForEnygmaNftResourceID(nftResourceID common.Hash) Dvp721SwapForEnygmaLogOption {
	return func(log *types.Log, data *zkDvp721SwapForEnygmaData) {
		copy(data.nftResourceId[:], nftResourceID[:])
	}
}

func WithDvp721SwapForEnygmaNftId(nftId *big.Int) Dvp721SwapForEnygmaLogOption {
	return func(log *types.Log, data *zkDvp721SwapForEnygmaData) {
		data.nftId = nftId
	}
}

func WithDvp721SwapForEnygmaBlockNumber(blockNumber uint64) Dvp721SwapForEnygmaLogOption {
	return func(log *types.Log, data *zkDvp721SwapForEnygmaData) {
		log.BlockNumber = blockNumber
	}
}

func NewDvp721SwapForEnygmaLogWith(opts ...Dvp721SwapForEnygmaLogOption) types.Log {
	data := zkDvp721SwapForEnygmaData{
		sharedId:         [32]byte{0x01, 0x02, 0x03},
		destChainId:      big.NewInt(1),
		from:             common.HexToAddress("0x1111111111111111111111111111111111111111"),
		nftResourceId:    [32]byte{0x0a, 0x0b, 0x0c},
		nftId:            big.NewInt(42),
		enygmaResourceId: [32]byte{0x0d, 0x0e, 0x0f},
		enygmaAmount:     big.NewInt(1000),
	}

	log := types.Log{
		Address:     common.HexToAddress("0x1234567890123456789012345678901234567890"),
		BlockNumber: 110,
		TxHash:      common.HexToHash("0xabcdef10"),
		TxIndex:     11,
		BlockHash:   common.HexToHash("0x123460"),
		Index:       10,
	}

	for _, opt := range opts {
		opt(&log, &data)
	}

	abi, _ := EnygmaPNEvents.EnygmaPNEventsMetaData.ParseABI()
	packedData, _ := abi.Events["Dvp721SwapForEnygma"].Inputs.NonIndexed().Pack(
		data.nftResourceId, data.nftId, data.enygmaAmount, data.enygmaResourceId, data.from, data.destChainId, data.sharedId, data.validityTime)

	log.Topics = []common.Hash{
		crypto.Keccak256Hash([]byte("Dvp721SwapForEnygma(bytes32,uint256,uint256,bytes32,address,uint256,bytes32,uint64)")),
	}
	log.Data = packedData
	return log
}

// Dvp1155Creation - ResourceID
type zkDvp1155CreationData struct {
	resourceID [32]byte
}

type Dvp1155CreationLogOption func(*types.Log, *zkDvp1155CreationData)

func WithDvp1155CreationResourceID(resourceID common.Hash) Dvp1155CreationLogOption {
	return func(log *types.Log, data *zkDvp1155CreationData) {
		copy(data.resourceID[:], resourceID[:])
	}
}

func WithDvp1155CreationBlockNumber(blockNumber uint64) Dvp1155CreationLogOption {
	return func(log *types.Log, data *zkDvp1155CreationData) {
		log.BlockNumber = blockNumber
	}
}

func NewDvp1155CreationLogWith(opts ...Dvp1155CreationLogOption) types.Log {
	data := zkDvp1155CreationData{
		resourceID: [32]byte{0x0a, 0x0b, 0x0c},
	}

	log := types.Log{
		Address:     common.HexToAddress("0x1234567890123456789012345678901234567890"),
		BlockNumber: 111,
		TxHash:      common.HexToHash("0xabcdef11"),
		TxIndex:     12,
		BlockHash:   common.HexToHash("0x123461"),
		Index:       11,
	}

	for _, opt := range opts {
		opt(&log, &data)
	}

	abi, _ := EnygmaPNEvents.EnygmaPNEventsMetaData.ParseABI()
	packedData, _ := abi.Events["Dvp1155Creation"].Inputs.NonIndexed().Pack(data.resourceID)

	log.Topics = []common.Hash{
		crypto.Keccak256Hash([]byte("Dvp1155Creation(bytes32)")),
	}
	log.Data = packedData
	return log
}

// Dvp1155Mint - ResourceID and TokenId
type zkDvp1155MintData struct {
	resourceID [32]byte
	tokenId    *big.Int
	value      *big.Int
	data       []byte
}

type Dvp1155MintLogOption func(*types.Log, *zkDvp1155MintData)

func WithDvp1155MintResourceID(resourceID common.Hash) Dvp1155MintLogOption {
	return func(log *types.Log, data *zkDvp1155MintData) {
		copy(data.resourceID[:], resourceID[:])
	}
}

func WithDvp1155MintTokenId(tokenId *big.Int) Dvp1155MintLogOption {
	return func(log *types.Log, data *zkDvp1155MintData) {
		data.tokenId = tokenId
	}
}

func WithDvp1155MintBlockNumber(blockNumber uint64) Dvp1155MintLogOption {
	return func(log *types.Log, data *zkDvp1155MintData) {
		log.BlockNumber = blockNumber
	}
}

func NewDvp1155MintLogWith(opts ...Dvp1155MintLogOption) types.Log {
	data := zkDvp1155MintData{
		resourceID: [32]byte{0x0a, 0x0b, 0x0c},
		tokenId:    big.NewInt(1),
		value:      big.NewInt(100),
		data:       []byte{},
	}

	log := types.Log{
		Address:     common.HexToAddress("0x1234567890123456789012345678901234567890"),
		BlockNumber: 112,
		TxHash:      common.HexToHash("0xabcdef12"),
		TxIndex:     13,
		BlockHash:   common.HexToHash("0x123462"),
		Index:       12,
	}

	for _, opt := range opts {
		opt(&log, &data)
	}

	abi, _ := EnygmaPNEvents.EnygmaPNEventsMetaData.ParseABI()
	packedData, _ := abi.Events["Dvp1155Mint"].Inputs.NonIndexed().Pack(data.resourceID, data.tokenId, data.value, data.data)

	log.Topics = []common.Hash{
		crypto.Keccak256Hash([]byte("Dvp1155Mint(bytes32,uint256,uint256,bytes)")),
	}
	log.Data = packedData
	return log
}

// Dvp1155Burn - ResourceID and TokenId
type zkDvp1155BurnData struct {
	resourceID [32]byte
	to         common.Address
	tokenId    *big.Int
	value      *big.Int
}

type Dvp1155BurnLogOption func(*types.Log, *zkDvp1155BurnData)

func WithDvp1155BurnResourceID(resourceID common.Hash) Dvp1155BurnLogOption {
	return func(log *types.Log, data *zkDvp1155BurnData) {
		copy(data.resourceID[:], resourceID[:])
	}
}

func WithDvp1155BurnTokenId(tokenId *big.Int) Dvp1155BurnLogOption {
	return func(log *types.Log, data *zkDvp1155BurnData) {
		data.tokenId = tokenId
	}
}

func WithDvp1155BurnBlockNumber(blockNumber uint64) Dvp1155BurnLogOption {
	return func(log *types.Log, data *zkDvp1155BurnData) {
		log.BlockNumber = blockNumber
	}
}

func NewDvp1155BurnLogWith(opts ...Dvp1155BurnLogOption) types.Log {
	data := zkDvp1155BurnData{
		resourceID: [32]byte{0x0a, 0x0b, 0x0c},
		to:         common.HexToAddress("0x1111111111111111111111111111111111111111"),
		tokenId:    big.NewInt(1),
		value:      big.NewInt(100),
	}

	log := types.Log{
		Address:     common.HexToAddress("0x1234567890123456789012345678901234567890"),
		BlockNumber: 113,
		TxHash:      common.HexToHash("0xabcdef13"),
		TxIndex:     14,
		BlockHash:   common.HexToHash("0x123463"),
		Index:       13,
	}

	for _, opt := range opts {
		opt(&log, &data)
	}

	abi, _ := EnygmaPNEvents.EnygmaPNEventsMetaData.ParseABI()
	packedData, _ := abi.Events["Dvp1155Burn"].Inputs.NonIndexed().Pack(data.resourceID, data.to, data.tokenId, data.value)

	log.Topics = []common.Hash{
		crypto.Keccak256Hash([]byte("Dvp1155Burn(bytes32,address,uint256,uint256)")),
	}
	log.Data = packedData
	return log
}

// Dvp1155DepositIntoDvp - ResourceID and TokenId
type dvp1155DepositIntoDvpData struct {
	resourceID [32]byte
	tokenId    *big.Int
	value      *big.Int
	data       []byte
	from       common.Address
}

type Dvp1155DepositIntoDvpLogOption func(*types.Log, *dvp1155DepositIntoDvpData)

func WithDvp1155DepositIntoDvpResourceID(resourceID common.Hash) Dvp1155DepositIntoDvpLogOption {
	return func(log *types.Log, data *dvp1155DepositIntoDvpData) {
		copy(data.resourceID[:], resourceID[:])
	}
}

func WithDvp1155DepositIntoDvpTokenId(tokenId *big.Int) Dvp1155DepositIntoDvpLogOption {
	return func(log *types.Log, data *dvp1155DepositIntoDvpData) {
		data.tokenId = tokenId
	}
}

func WithDvp1155DepositIntoDvpBlockNumber(blockNumber uint64) Dvp1155DepositIntoDvpLogOption {
	return func(log *types.Log, data *dvp1155DepositIntoDvpData) {
		log.BlockNumber = blockNumber
	}
}

func NewDvp1155DepositIntoDvpLogWith(opts ...Dvp1155DepositIntoDvpLogOption) types.Log {
	data := dvp1155DepositIntoDvpData{
		resourceID: [32]byte{0x0a, 0x0b, 0x0c},
		tokenId:    big.NewInt(1),
		value:      big.NewInt(100),
		data:       []byte{},
		from:       common.HexToAddress("0x1111111111111111111111111111111111111111"),
	}

	log := types.Log{
		Address:     common.HexToAddress("0x1234567890123456789012345678901234567890"),
		BlockNumber: 114,
		TxHash:      common.HexToHash("0xabcdef14"),
		TxIndex:     15,
		BlockHash:   common.HexToHash("0x123464"),
		Index:       14,
	}

	for _, opt := range opts {
		opt(&log, &data)
	}

	abi, _ := EnygmaPNEvents.EnygmaPNEventsMetaData.ParseABI()
	packedData, _ := abi.Events["Dvp1155DepositIntoDvp"].Inputs.NonIndexed().Pack(data.resourceID, data.tokenId, data.from, data.value, data.data)

	log.Topics = []common.Hash{
		crypto.Keccak256Hash([]byte("Dvp1155DepositIntoDvp(bytes32,uint256,address,uint256,bytes)")),
	}
	log.Data = packedData
	return log
}

// Dvp1155WithdrawFromDvp - ResourceID and TokenId
type dvp1155WithdrawFromDvpData struct {
	resourceID [32]byte
	tokenId    *big.Int
	value      *big.Int
	data       []byte
	owner      common.Address
}

type Dvp1155WithdrawFromDvpLogOption func(*types.Log, *dvp1155WithdrawFromDvpData)

func WithDvp1155WithdrawFromDvpResourceID(resourceID common.Hash) Dvp1155WithdrawFromDvpLogOption {
	return func(log *types.Log, data *dvp1155WithdrawFromDvpData) {
		copy(data.resourceID[:], resourceID[:])
	}
}

func WithDvp1155WithdrawFromDvpTokenId(tokenId *big.Int) Dvp1155WithdrawFromDvpLogOption {
	return func(log *types.Log, data *dvp1155WithdrawFromDvpData) {
		data.tokenId = tokenId
	}
}

func WithDvp1155WithdrawFromDvpBlockNumber(blockNumber uint64) Dvp1155WithdrawFromDvpLogOption {
	return func(log *types.Log, data *dvp1155WithdrawFromDvpData) {
		log.BlockNumber = blockNumber
	}
}

func NewDvp1155WithdrawFromDvpLogWith(opts ...Dvp1155WithdrawFromDvpLogOption) types.Log {
	data := dvp1155WithdrawFromDvpData{
		resourceID: [32]byte{0x0a, 0x0b, 0x0c},
		tokenId:    big.NewInt(1),
		value:      big.NewInt(100),
		data:       []byte{},
		owner:      common.HexToAddress("0x1111111111111111111111111111111111111111"),
	}

	log := types.Log{
		Address:     common.HexToAddress("0x1234567890123456789012345678901234567890"),
		BlockNumber: 115,
		TxHash:      common.HexToHash("0xabcdef15"),
		TxIndex:     16,
		BlockHash:   common.HexToHash("0x123465"),
		Index:       15,
	}

	for _, opt := range opts {
		opt(&log, &data)
	}

	abi, _ := EnygmaPNEvents.EnygmaPNEventsMetaData.ParseABI()
	packedData, _ := abi.Events["Dvp1155WithdrawFromDvp"].Inputs.NonIndexed().Pack(data.resourceID, data.tokenId, data.value, data.data, data.owner)

	log.Topics = []common.Hash{
		crypto.Keccak256Hash([]byte("Dvp1155WithdrawFromDvp(bytes32,uint256,uint256,bytes,address)")),
	}
	log.Data = packedData
	return log
}

// Dvp1155SwapForEnygma - TokenResourceID and TokenId
type zkDvp1155SwapForEnygmaData struct {
	tokenResourceId  [32]byte
	tokenId          *big.Int
	tokenValue       *big.Int
	tokenData        []byte
	enygmaAmount     *big.Int
	enygmaResourceId [32]byte
	from             common.Address
	destChainId      *big.Int
	sharedId         [32]byte
	validityTime     uint64
}

type Dvp1155SwapForEnygmaLogOption func(*types.Log, *zkDvp1155SwapForEnygmaData)

func WithDvp1155SwapForEnygmaTokenResourceID(tokenResourceID common.Hash) Dvp1155SwapForEnygmaLogOption {
	return func(log *types.Log, data *zkDvp1155SwapForEnygmaData) {
		copy(data.tokenResourceId[:], tokenResourceID[:])
	}
}

func WithDvp1155SwapForEnygmaTokenId(tokenId *big.Int) Dvp1155SwapForEnygmaLogOption {
	return func(log *types.Log, data *zkDvp1155SwapForEnygmaData) {
		data.tokenId = tokenId
	}
}

func WithDvp1155SwapForEnygmaBlockNumber(blockNumber uint64) Dvp1155SwapForEnygmaLogOption {
	return func(log *types.Log, data *zkDvp1155SwapForEnygmaData) {
		log.BlockNumber = blockNumber
	}
}

func NewDvp1155SwapForEnygmaLogWith(opts ...Dvp1155SwapForEnygmaLogOption) types.Log {
	data := zkDvp1155SwapForEnygmaData{
		tokenResourceId:  [32]byte{0x0a, 0x0b, 0x0c},
		tokenId:          big.NewInt(1),
		tokenValue:       big.NewInt(100),
		tokenData:        []byte{},
		enygmaAmount:     big.NewInt(1000),
		enygmaResourceId: [32]byte{0x0d, 0x0e, 0x0f},
		from:             common.HexToAddress("0x1111111111111111111111111111111111111111"),
		destChainId:      big.NewInt(1),
		sharedId:         [32]byte{0x01, 0x02, 0x03},
	}

	log := types.Log{
		Address:     common.HexToAddress("0x1234567890123456789012345678901234567890"),
		BlockNumber: 116,
		TxHash:      common.HexToHash("0xabcdef16"),
		TxIndex:     17,
		BlockHash:   common.HexToHash("0x123466"),
		Index:       16,
	}

	for _, opt := range opts {
		opt(&log, &data)
	}

	abi, _ := EnygmaPNEvents.EnygmaPNEventsMetaData.ParseABI()
	packedData, _ := abi.Events["Dvp1155SwapForEnygma"].Inputs.NonIndexed().Pack(
		data.tokenResourceId, data.tokenId, data.tokenValue, data.tokenData, data.enygmaAmount, data.enygmaResourceId, data.from, data.destChainId, data.sharedId, data.validityTime)

	log.Topics = []common.Hash{
		crypto.Keccak256Hash([]byte("Dvp1155SwapForEnygma(bytes32,uint256,uint256,bytes,uint256,bytes32,address,uint256,bytes32,uint64)")),
	}
	log.Data = packedData
	return log
}
