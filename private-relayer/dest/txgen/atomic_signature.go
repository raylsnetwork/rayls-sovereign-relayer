// Decommissioning Teleport (vanilla, atomic).

package txgen

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/contracts/EndpointV1"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/types"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/withstack"
)

// ethClientTimeout is the timeout for Ethereum client calls during initialization.
const ethClientTimeout = 30 * time.Second

// SignatureGenerator packs receivePayload calldata for atomic-signature
// transactions. It only produces the calldata bytes — signing and
// broadcasting are owned by CTS, which receives the calldata via the
// `cts.send.privatenode` NATS subject.
//
// Deprecated: Decommissioning Teleport (vanilla, atomic).
type SignatureGenerator struct {
	endpointAddress common.Address

	endpointABI *abi.ABI
	chainID     *big.Int
}

// Deprecated: Decommissioning Teleport (vanilla, atomic).
func NewSignatureGenerator(client EthereumClient, endpointAddress common.Address) (*SignatureGenerator, error) {
	endpointABI, err := EndpointV1.EndpointV1MetaData.ParseABI()
	if err != nil {
		return nil, withstack.Wrap(fmt.Errorf("error parsing endpoint ABI: %w", err))
	}

	ctxChainID, cancelChainID := context.WithTimeout(context.Background(), ethClientTimeout)
	chainID, err := client.ChainID(ctxChainID)
	cancelChainID()
	if err != nil {
		return nil, withstack.Wrap(fmt.Errorf("error getting chain ID: %w", err))
	}

	return &SignatureGenerator{
		endpointAddress: endpointAddress,
		endpointABI:     endpointABI,
		chainID:         chainID,
	}, nil
}

// Generate returns the receivePayload calldata for a single
// CalldataSignature. The returned bytes are the input to a tx targeting
// the endpoint address; CTS handles nonce, gas, signing, broadcast.
//
// Deprecated: Decommissioning Teleport (vanilla, atomic).
func (g *SignatureGenerator) Generate(signature types.CalldataSignature) ([]byte, error) {
	sharedIdBytes := []byte(signature.SharedId)
	msgId := crypto.Keccak256Hash(sharedIdBytes)

	var messageToSend EndpointV1.RaylsMessage
	messageToSend.Payload = signature.Signature
	messageToSend.MessageMetadata = EndpointV1.RaylsMessageMetadata{
		Valid:        true,
		IgnoresNonce: true,
		Nonce:        big.NewInt(0),
		ResourceId:   signature.ResourceId,
		NewResourceMetadata: EndpointV1.NewResourceMetadata{
			Valid:              false,
			ResourceDeployType: 0,
			Bytecode:           []byte{},
			FactoryTemplate:    0,
			InitializerParams:  []byte{},
		},
		LockData:                  []byte{},
		RevertPayloadDataSender:   []byte{},
		RevertPayloadDataReceiver: []byte{},
		TransferMetadata: EndpointV1.BridgedTransferMetadata{
			AssetType: uint8(0),
			Id:        big.NewInt(0),
			Amount:    big.NewInt(0),
		},
	}

	data, err := g.endpointABI.Pack(
		"receivePayload",
		g.chainID,
		g.endpointAddress,
		common.Address{0},
		messageToSend,
		msgId,
	)
	if err != nil {
		return nil, withstack.Wrap(fmt.Errorf("error packing the ABI: %w", err))
	}
	return data, nil
}
