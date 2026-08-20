// Decommissioning Teleport (vanilla, atomic).

package service

import "github.com/ethereum/go-ethereum/common"

type TokenStandard int

// Must match RaylsNodeBridgeableERC in RNMessageLib.sol:
// 0: CUSTOM, 1: ERC20, 2: ERC721, 3: ERC1155, 4: ENYGMA, 5: ZKDVPERC721, 6: ZKDVPERC1155
const (
	CUSTOM TokenStandard = iota
	ERC20
	ERC721
	ERC1155
	ENYGMA
	ZKDVPERC721
	ZKDVPERC1155
)

type Deployment struct {
	PrivateAddress common.Address

	Standard TokenStandard

	URI    string
	Name   string
	Symbol string
}

func (d Deployment) GetID() string {
	return d.PrivateAddress.String()
}
