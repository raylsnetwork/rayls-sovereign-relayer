package contractclient

import (
	"context"
	"encoding/hex"
	"errors"
	"log/slog"
	"math/big"

	"github.com/raylsnetwork/rayls-privacy-relayer-api/contracts/Dvp"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/conv"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/dvp"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/types"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/withstack"

	"github.com/ethereum/go-ethereum/common"
)

type dvpEncryptor interface {
	EncryptDvpSwapMessage(ctx context.Context, salt *big.Int, message *types.DvpSwapMessage) ([]byte, error)
}
type DvpClient struct {
	address          common.Address
	contract         *Dvp.Dvp
	executor         Executor
	operatorExecutor Executor
	encryptor        dvpEncryptor
}

func NewDvpClient(
	address common.Address,
	executor Executor,
	operatorExecutor Executor,
	encryptor dvpEncryptor,
) *DvpClient {
	return &DvpClient{
		address:          address,
		contract:         Dvp.NewDvp(),
		executor:         executor,
		operatorExecutor: operatorExecutor,
		encryptor:        encryptor,
	}
}

// ERC721
func (c *DvpClient) DepositERC721(ctx context.Context, chainEventID string, contractAddress common.Address, nftId *big.Int, publicKey *big.Int, salt *big.Int, encryptedBalanceUpdate []byte) error {
	calldata := c.contract.PackDepositERC721(contractAddress, nftId, publicKey, salt, encryptedBalanceUpdate)

	_, err := c.operatorExecutor.Execute(ctx, IDFor("dvp.DepositERC721", chainEventID), calldata, c.address)
	if err != nil {
		return WrapInDvpClientError("failed to deposit ERC721", withstack.Wrap(err))
	}

	return nil
}

func (c *DvpClient) WithdrawERC721(ctx context.Context, chainEventID string, contractAddress common.Address, nftId *big.Int, to common.Address, salt *big.Int, proof *dvp.ProofReceipt, encryptedBalanceUpdate []byte) error {
	// For ERC721 ownership proof: 1 input (nullifier), 1 output (commitment)
	contractProof := proofReceiptToContractProofReceipt(proof)

	calldata := c.contract.PackWithdrawERC721(contractAddress, nftId, to, salt, contractProof, encryptedBalanceUpdate)

	_, err := c.operatorExecutor.Execute(ctx, IDFor("dvp.WithdrawERC721", chainEventID), calldata, c.address)
	if err != nil {
		return WrapInDvpClientError("failed to withdraw ERC721", withstack.Wrap(err))
	}
	return nil
}

// ERC1155
func (c *DvpClient) MixFundsERC1155(ctx context.Context, chainEventID string, contractAddress common.Address, proofReceipt *dvp.ProofReceipt) error {
	contractJSProof := proofReceiptToContractProofReceipt(proofReceipt)

	calldata := c.contract.PackMixFundsERC1155(contractAddress, contractJSProof)

	_, err := c.executor.Execute(ctx, IDFor("dvp.MixFundsERC1155", chainEventID), calldata, c.address)
	if err != nil {
		return WrapInDvpClientError("failed to mix funds", withstack.Wrap(err))
	}

	return nil
}

func (c *DvpClient) DepositERC1155(ctx context.Context, chainEventID string, contractAddress common.Address, tokenId *big.Int, tokenAmount *big.Int, data []byte, publicKey *big.Int, salt *big.Int, encryptedBalanceUpdate []byte) error {
	calldata := c.contract.PackDepositERC1155(contractAddress, tokenId, tokenAmount, data, publicKey, salt, encryptedBalanceUpdate)

	_, err := c.operatorExecutor.Execute(ctx, IDFor("dvp.DepositERC1155", chainEventID), calldata, c.address)
	if err != nil {
		var rd *ErrorWithRevertData
		if errors.As(err, &rd) {
			slog.Error("failed to deposit ERC1155",
				slog.String("revert_data", hex.EncodeToString(rd.GetRevertData())),
			)
		}
		return WrapInDvpClientError("failed to deposit ERC1155", withstack.Wrap(err))
	}

	return nil
}

func (c *DvpClient) WithdrawERC1155(ctx context.Context, chainEventID string, contractAddress common.Address, tokenId *big.Int, tokenAmount *big.Int, to common.Address, salt *big.Int, proof *dvp.ProofReceipt, encryptedBalanceUpdate []byte) error {
	contractProof := proofReceiptToContractProofReceipt(proof)

	calldata := c.contract.PackWithdrawERC1155(contractAddress, tokenId, tokenAmount, to, salt, contractProof, encryptedBalanceUpdate)

	_, err := c.operatorExecutor.Execute(ctx, IDFor("dvp.WithdrawERC1155", chainEventID), calldata, c.address)
	if err != nil {
		return WrapInDvpClientError("failed to withdraw ERC1155", withstack.Wrap(err))
	}
	return nil
}

func (c *DvpClient) getSwapProofType(tokenType types.DvpTokenType) DvpSwapProofType {
	if tokenType == types.DvpEnygma || tokenType == types.DvpERC1155 {
		return DvpSwapPaymentProofType
	}

	return DvpSwapDeliveryProofType
}

func (c *DvpClient) InitiateSwap(ctx context.Context, salt *big.Int, ciphertext []byte, msg *types.DvpSwapMessage, proof *dvp.ProofReceipt, validityTime uint64, passphrase *big.Int) error {
	encrMsg, err := c.encryptor.EncryptDvpSwapMessage(ctx, salt, msg)
	if err != nil {
		return WrapInDvpClientError("failed to encrypt dvp swap message", err)
	}

	contractProof := proofReceiptToContractProofReceipt(proof)
	contractProofType := c.getSwapProofType(msg.TokenInType)

	sharedId, err := conv.StringToBytes32(msg.SharedId)
	if err != nil {
		return WrapInDvpClientError("failed to convert shared chainEventID to bytes32", err)
	}

	calldata := c.contract.PackInitiateSwap(sharedId, encrMsg, ciphertext, common.HexToAddress(msg.TokenInAddress), uint8(contractProofType), contractProof, validityTime, passphrase)

	_, err = c.executor.Execute(ctx, IDFor("dvp.InitiateSwap", msg.SharedId), calldata, c.address)
	if err != nil {
		if IsRevertWithSelector(err, Dvp.DvpDvpSwapAlreadyExistsErrorID()) {
			return dvp.ErrSwapAlreadyInitiated
		}
		return WrapInDvpClientError("failed to initiate swap", withstack.Wrap(err))
	}

	return nil
}

func (c *DvpClient) CompleteSwap(ctx context.Context, salt *big.Int, msg *types.DvpSwapMessage, proof *dvp.ProofReceipt) error {
	encrMsg, err := c.encryptor.EncryptDvpSwapMessage(ctx, salt, msg)
	if err != nil {
		return WrapInDvpClientError("failed to encrypt dvp swap message", err)
	}

	contractProof := proofReceiptToContractProofReceipt(proof)
	contractProofType := c.getSwapProofType(msg.TokenInType)

	sharedId, err := conv.StringToBytes32(msg.SharedId)
	if err != nil {
		return WrapInDvpClientError("failed to convert shared chainEventID to bytes32", err)
	}

	calldata := c.contract.PackCompleteSwap(sharedId, common.HexToAddress(msg.TokenInAddress), uint8(contractProofType), contractProof, encrMsg)

	_, err = c.executor.Execute(ctx, IDFor("dvp.CompleteSwap", msg.SharedId), calldata, c.address)
	if err != nil {
		return WrapInDvpClientError("failed to execute swap completion", withstack.Wrap(err))
	}

	return nil
}

func (c *DvpClient) IsSwapExpired(ctx context.Context, sharedId string) (bool, error) {
	sharedIdBytes, err := conv.StringToBytes32(sharedId)
	if err != nil {
		return false, WrapInDvpClientError("failed to convert shared chainEventID to bytes32", err)
	}

	calldata := c.contract.PackIsSwapExpired(sharedIdBytes)

	data, err := c.executor.Call(ctx, c.address, calldata)
	if err != nil {
		return false, WrapInDvpClientError("failed to check if swap is expired", withstack.Wrap(err))
	}

	expired, err := c.contract.UnpackIsSwapExpired(data)
	if err != nil {
		return false, WrapInDvpClientError("failed to unpack is expired", withstack.Wrap(err))
	}

	return expired, nil
}

func (c *DvpClient) CancelSwap(ctx context.Context, sharedId string, preimage *big.Int) error {
	sharedIdBytes, err := conv.StringToBytes32(sharedId)
	if err != nil {
		return WrapInDvpClientError("failed to convert shared chainEventID to bytes32", err)
	}

	calldata := c.contract.PackCancelSwap(sharedIdBytes, preimage)

	_, err = c.executor.Execute(ctx, IDFor("dvp.CancelSwap", sharedId), calldata, c.address)
	if err != nil {
		return WrapInDvpClientError("failed to execute swap cancellation", withstack.Wrap(err))
	}
	return nil
}

func (c *DvpClient) ExpireSwap(ctx context.Context, sharedId string) error {
	sharedIdBytes, err := conv.StringToBytes32(sharedId)
	if err != nil {
		return WrapInDvpClientError("failed to convert shared chainEventID to bytes32", err)
	}

	calldata := c.contract.PackExpireSwap(sharedIdBytes)

	_, err = c.executor.Execute(ctx, IDFor("dvp.ExpireSwap", sharedId), calldata, c.address)
	if err != nil {
		return WrapInDvpClientError("failed to execute swap expiration", withstack.Wrap(err))
	}
	return nil
}

func proofReceiptToContractProofReceipt(proofReceipt *dvp.ProofReceipt) Dvp.IDvpProofReceipt {
	return Dvp.IDvpProofReceipt{
		Proof: Dvp.IDvpSnarkProof{
			A: Dvp.IDvpG1Point{
				X: proofReceipt.Proof.A[0],
				Y: proofReceipt.Proof.A[1],
			},
			B: Dvp.IDvpG2Point{
				X: [2]*big.Int{
					proofReceipt.Proof.B[0][0],
					proofReceipt.Proof.B[0][1],
				},
				Y: [2]*big.Int{
					proofReceipt.Proof.B[1][0],
					proofReceipt.Proof.B[1][1],
				},
			},
			C: Dvp.IDvpG1Point{
				X: proofReceipt.Proof.C[0],
				Y: proofReceipt.Proof.C[1],
			},
		},
		TreeNumbers:      proofReceipt.TreeNumbers,
		Message:          proofReceipt.Message,
		MerkleRoots:      proofReceipt.MerkleRoots,
		Commitments:      proofReceipt.Commitments,
		Nullifiers:       proofReceipt.Nullifiers,
		RevertCommitment: proofReceipt.RevertCommitment,
	}
}
