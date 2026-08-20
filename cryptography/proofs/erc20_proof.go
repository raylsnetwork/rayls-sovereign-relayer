package proofs

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

// evmWordSize is the size of an EVM word in bytes (256 bits).
const evmWordSize = 32

func GetSlotForMapKey(keyInMap []byte, slotIndexForMap int) [32]byte {
	return crypto.Keccak256Hash(
		keyInMap,
		common.LeftPadBytes(big.NewInt(int64(slotIndexForMap)).Bytes(), evmWordSize),
	)
}

func GetSlotForERC20TokenHolder(slotIndexForHoldersMap int, tokenHolder common.Address) [32]byte {
	return GetSlotForMapKey(common.LeftPadBytes(tokenHolder[:], evmWordSize), slotIndexForHoldersMap)
}

func GetSlotForArrayItem(slotIndexForArray int, indexInArray int, itemSize int) [32]byte {
	bytes := crypto.Keccak256Hash(common.LeftPadBytes(big.NewInt(int64(slotIndexForArray)).Bytes(), evmWordSize))
	arrayPos := new(big.Int).SetBytes(bytes[:])
	itemPos := arrayPos.Add(arrayPos, big.NewInt(int64(indexInArray*itemSize)))
	var pos [32]byte
	copy(pos[:], itemPos.Bytes()[:evmWordSize])

	return pos
}
