// Decommissioning Teleport (vanilla, atomic): only EncryptAdditionalData below; Enygma/DVP encrypt methods are retained.

package contractclient

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	ps "github.com/raylsnetwork/rayls-sovereign-relayer/contracts/ParticipantStorageV1"
	"github.com/raylsnetwork/rayls-sovereign-relayer/cts/gen/proto/encrypt"
	"github.com/raylsnetwork/rayls-sovereign-relayer/cts/wireformat"
	"github.com/raylsnetwork/rayls-sovereign-relayer/types"
	"github.com/raylsnetwork/rayls-sovereign-relayer/withstack"
	"google.golang.org/grpc"
)

const encryptBlockTimeout = 30 * time.Second

type encryptEthereumClient interface {
	BlockByNumber(context.Context, *big.Int) (*ethTypes.Block, error)
}

type participantStorageClient interface {
	GetVenOperatorChainInfo(ctx context.Context, blockNumber *big.Int) (ps.ParticipantStructsPrivacyNodeViewData, error)
	GetMyChainInfo(ctx context.Context, blockNumber *big.Int) (ps.ParticipantStructsPrivacyNodeViewData, error)
	GetChainViewDataBatch(ctx context.Context, chainIDs []*big.Int) ([]ps.ParticipantStructsPrivacyNodeViewData, error)
	GetChainViewData(ctx context.Context, chainID, blockNumber *big.Int) (ps.ParticipantStructsPrivacyNodeViewData, error)
}

type ctsClient interface {
	Encrypt(
		ctx context.Context,
		req *encrypt.EncryptRequest,
		opts ...grpc.CallOption,
	) (*encrypt.EncryptResponse, error)
	EncryptWithoutFP(
		ctx context.Context,
		in *encrypt.EncryptWithoutFPRequest,
		opts ...grpc.CallOption,
	) (*encrypt.EncryptWithoutFPResponse, error)
	EncryptWithoutFPWithSS(
		ctx context.Context,
		in *encrypt.EncryptWithoutFPWithSSRequest,
		opts ...grpc.CallOption,
	) (*encrypt.EncryptWithoutFPWithSSResponse, error)
}

type Encryptor struct {
	cts        ctsClient
	psClient   participantStorageClient
	ethClient  encryptEthereumClient
	pnhChainID *big.Int
}

func NewEncryptor(cts ctsClient, psClient participantStorageClient, ethClient encryptEthereumClient, pnhChainID *big.Int) *Encryptor {
	return &Encryptor{
		cts:        cts,
		psClient:   psClient,
		ethClient:  ethClient,
		pnhChainID: pnhChainID,
	}
}

// Deprecated: Decommissioning Teleport (vanilla, atomic).
func (a *Encryptor) EncryptAdditionalData(ctx context.Context, data []types.AtomicTeleportAdditionalData) (string, error) {
	ctxBlock, cancelBlock := context.WithTimeout(ctx, encryptBlockTimeout)
	block, err := a.ethClient.BlockByNumber(ctxBlock, nil)
	cancelBlock() // Clean up immediately to prevent leak
	if err != nil {
		return "", WrapInEncryptorError("failed to get latest block", withstack.Wrap(err))
	}

	chainInfo, err := a.psClient.GetVenOperatorChainInfo(ctx, nil)
	if err != nil {
		return "", WrapInEncryptorError("failed to get VenOperator chain info", err)
	}

	plaintext, err := wireformat.MarshalPlaintext(data)
	if err != nil {
		return "", WrapInEncryptorError("failed to marshal additional data", withstack.Wrap(err))
	}

	encrData, err := a.cts.EncryptWithoutFP(
		ctx,
		&encrypt.EncryptWithoutFPRequest{
			Plaintext:   plaintext,
			ChainID:     chainInfo.ChainId.Uint64(),
			BlockNumber: block.NumberU64(),
		})
	if err != nil {
		return "", WrapInEncryptorError("failed to encrypt additional data", err)
	}

	return common.Bytes2Hex(encrData.GetEncryptedData()), nil
}

func (a *Encryptor) EncryptMessages(
	ctx context.Context,
	messages []types.DispatchedMessageToPrivateHub,
	chainID *big.Int,
) (EncryptedMessages, error) {
	ctxBlock, cancelBlock := context.WithTimeout(ctx, encryptBlockTimeout)
	defer cancelBlock()

	block, err := a.ethClient.BlockByNumber(ctxBlock, nil)
	if err != nil {
		return EncryptedMessages{}, WrapInEncryptorError("failed to get latest block", withstack.Wrap(err))
	}

	chainInfo, err := a.psClient.GetChainViewData(ctx, chainID, block.Number())
	if err != nil {
		return EncryptedMessages{}, WrapInEncryptorError("failed to get chain info", err)
	}

	plaintext, err := wireformat.MarshalPlaintext(messages)
	if err != nil {
		return EncryptedMessages{}, WrapInEncryptorError("failed to marshal messages", withstack.Wrap(err))
	}

	resp, err := a.cts.Encrypt(ctx, &encrypt.EncryptRequest{
		Plaintext:     plaintext,
		ChainID:       chainInfo.ChainId.Uint64(),
		BlockNumber:   block.Number().Uint64(),
		PrevBlockHash: block.ParentHash().String(),
	})
	if err != nil {
		return EncryptedMessages{}, WrapInEncryptorError("failed to encrypt messages", err)
	}

	return EncryptedMessages{
		BlockNumber: block.Number(),
		MessageTag:  resp.Fingerprint,
		Data:        resp.EncryptedData,
	}, nil
}

func (a *Encryptor) EncryptEnygmaTransferBatches(
	ctx context.Context,
	batches []*types.EnygmaTransferBatch,
	blockNumber *big.Int,
) ([][]byte, error) {
	chainIDs := make([]*big.Int, len(batches))
	for i, batch := range batches {
		chainIDs[i] = batch.ToChainID
	}

	chainInfos, err := a.psClient.GetChainViewDataBatch(ctx, chainIDs)
	if err != nil {
		return [][]byte{}, WrapInEncryptorError("failed to get chain infos", err)
	}

	infoByChainID := make(map[string]ps.ParticipantStructsPrivacyNodeViewData)
	for _, chainInfo := range chainInfos {
		infoByChainID[chainInfo.ChainId.String()] = chainInfo
	}

	encrBatches := make([][]byte, len(batches))
	for i, batch := range batches {
		plaintext, err := wireformat.MarshalPlaintext(batch)
		if err != nil {
			return [][]byte{}, WrapInEncryptorError("failed to marshal enygma transfer batch", withstack.Wrap(err))
		}

		resp, err := a.cts.EncryptWithoutFP(
			ctx,
			&encrypt.EncryptWithoutFPRequest{
				Plaintext:   plaintext,
				ChainID:     batch.ToChainID.Uint64(),
				BlockNumber: blockNumber.Uint64(),
				// ResourceID enables the CTS self-aware seal: the sender's own-chain (self/change)
				// batch (ToChainID == FromChainID) is sealed with the self-secret instead of the
				// own-chain pairwise row, so it can be AEAD-verified on decrypt and during resync.
				// common.FromHex matches how the self-secret is stored (proof.go: common.FromHex(resourceId)).
				ResourceID: common.FromHex(batch.ResourceId),
			},
		)
		if err != nil {
			return [][]byte{}, WrapInEncryptorError("failed to encrypt enygma transfer batch", err)
		}
		encrBatches[i] = resp.GetEncryptedData()
	}

	return encrBatches, nil
}

func (a *Encryptor) EncryptEnygmaTransferBatchCompleted(ctx context.Context, messages []types.EnygmaTransferCompleted) ([]byte, error) {
	block, err := a.ethClient.BlockByNumber(ctx, nil)
	if err != nil {
		return nil, WrapInEncryptorError("failed to get latest block", withstack.Wrap(err))
	}

	chainInfo, err := a.psClient.GetMyChainInfo(ctx, block.Number())
	if err != nil {
		return nil, WrapInEncryptorError("failed to get my chain info", err)
	}

	plaintext, err := wireformat.MarshalPlaintext(messages)
	if err != nil {
		return nil, WrapInEncryptorError("failed to marshal enygma transfer batch completed", withstack.Wrap(err))
	}

	resp, err := a.cts.EncryptWithoutFP(
		ctx,
		&encrypt.EncryptWithoutFPRequest{
			Plaintext:   plaintext,
			ChainID:     chainInfo.ChainId.Uint64(),
			BlockNumber: block.NumberU64(),
		},
	)
	if err != nil {
		return nil, WrapInEncryptorError("failed to encrypt enygma transfer batch completed", err)
	}

	return resp.EncryptedData, nil
}

func (a *Encryptor) EncryptDvpSwapMessage(
	ctx context.Context,
	salt *big.Int,
	message *types.DvpSwapMessage,
) ([]byte, error) {
	plaintext, err := wireformat.MarshalPlaintext(message)
	if err != nil {
		return nil, WrapInEncryptorError("failed to marshal dvp swap message", withstack.Wrap(err))
	}

	resp, err := a.cts.EncryptWithoutFPWithSS(
		ctx,
		&encrypt.EncryptWithoutFPWithSSRequest{
			Ss:        salt.Bytes(),
			Plaintext: plaintext,
		})
	if err != nil {
		return nil, WrapInEncryptorError("failed to encrypt dvp teleport message", err)
	}

	return resp.EncryptedData, nil
}

func (a *Encryptor) EncryptDvpExtraData(ctx context.Context, rawDataBytes []byte) ([]byte, error) {
	ctxBlock, cancelBlock := context.WithTimeout(ctx, encryptBlockTimeout)
	block, err := a.ethClient.BlockByNumber(ctxBlock, nil)
	cancelBlock() // Clean up immediately to prevent leak
	if err != nil {
		return []byte{}, WrapInEncryptorError("failed to get latest block", withstack.Wrap(err))
	}

	chainInfo, err := a.psClient.GetMyChainInfo(ctx, block.Number())
	if err != nil {
		return []byte{}, WrapInEncryptorError("failed to get my chain info", err)
	}

	type extraData struct {
		Key      string
		Value    string
		IsPublic bool
	}

	var encrData []extraData

	err = json.Unmarshal(rawDataBytes, &encrData)
	if err != nil {
		return []byte{}, WrapInEncryptorError("failed to unmarshal extra data", withstack.Wrap(err))
	}

	for i := range encrData {
		if !encrData[i].IsPublic {
			resp, err := a.cts.EncryptWithoutFP(ctx, &encrypt.EncryptWithoutFPRequest{
				Plaintext:   []byte(encrData[i].Value),
				ChainID:     chainInfo.ChainId.Uint64(),
				BlockNumber: block.NumberU64(),
			})
			if err != nil {
				return nil, WrapInEncryptorError("failed to encrypt dvp ERC721 extra data", err)
			}

			encrData[i].Value = base64.StdEncoding.EncodeToString(resp.EncryptedData)
		}
	}

	extraDataBytes, err := json.Marshal(encrData)
	if err != nil {
		return []byte{}, WrapInEncryptorError("failed to marshal extra data", withstack.Wrap(err))
	}

	return extraDataBytes, nil
}

// EncryptDvpBalanceUpdated creates encrypted balance update for Enygma DvP operations in order to notify Private Network Hub.
func (a *Encryptor) EncryptDvpBalanceUpdated(ctx context.Context, message types.DvpBalanceUpdated) ([]byte, error) {
	ctxBlock, cancelBlock := context.WithTimeout(ctx, encryptBlockTimeout)
	block, err := a.ethClient.BlockByNumber(ctxBlock, nil)
	cancelBlock() // Clean up immediately to prevent leak
	if err != nil {
		return nil, WrapInEncryptorError("failed to get latest block", withstack.Wrap(err))
	}

	operatorInfo, err := a.psClient.GetVenOperatorChainInfo(ctx, nil)
	if err != nil {
		return nil, WrapInEncryptorError("failed to get VenOperator chain info", err)
	}

	if message.DestinationChainId == nil {
		message.DestinationChainId = a.pnhChainID
	}

	plaintext, err := wireformat.MarshalPlaintext(message)
	if err != nil {
		return nil, WrapInEncryptorError("failed to marshal dvp balance updated", withstack.Wrap(err))
	}

	resp, err := a.cts.EncryptWithoutFP(
		ctx,
		&encrypt.EncryptWithoutFPRequest{
			Plaintext:   plaintext,
			ChainID:     operatorInfo.ChainId.Uint64(),
			BlockNumber: block.NumberU64(),
		},
	)
	if err != nil {
		return nil, WrapInEncryptorError("failed to encrypt dvp balance updated", err)
	}

	return resp.EncryptedData, nil
}
