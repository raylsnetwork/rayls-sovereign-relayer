package contractclient

import (
	"context"
	"crypto/sha256"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/contracts/EnygmaDvpIntegration"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/dvp"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/types"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/withstack"
)

type dvpIntegrationEncryptor interface {
	EncryptDvpBalanceUpdated(ctx context.Context, message types.DvpBalanceUpdated) ([]byte, error)
	EncryptEnygmaTransferBatches(ctx context.Context, batches []*types.EnygmaTransferBatch, blockNumber *big.Int) ([][]byte, error)
}

type dvpIntegrationEthereumClient interface {
	BlockByNumber(context.Context, *big.Int) (*ethTypes.Block, error)
}

type DvpIntegrationClient struct {
	contract  *EnygmaDvpIntegration.EnygmaDvpIntegration
	encr      dvpIntegrationEncryptor
	executor  Executor
	ethClient dvpIntegrationEthereumClient
}

func NewDvpIntegrationClient(
	executor Executor,
	encr dvpIntegrationEncryptor,
	ethClient dvpIntegrationEthereumClient,
) *DvpIntegrationClient {
	return &DvpIntegrationClient{
		contract:  EnygmaDvpIntegration.NewEnygmaDvpIntegration(),
		encr:      encr,
		executor:  executor,
		ethClient: ethClient,
	}
}

// SignDeposit signs a deposit transaction without broadcasting it.
// The caller is responsible for broadcasting and waiting for the receipt.
func (c *DvpIntegrationClient) Deposit(
	ctx context.Context,
	chainEventID string,
	batches []*types.EnygmaTransferBatch,
	proof *types.EnygmaProofResponse,
	blockNumber *big.Int,
	chainId *big.Int,
	resourceId string,
	amount *big.Int,
	from common.Address,
	sourceTxHash common.Hash,
	dvpIntegrationAddress common.Address,
) error {
	encrBatches, err := c.encr.EncryptEnygmaTransferBatches(ctx, batches, blockNumber)
	if err != nil {
		return WrapInEnygmaDvpIntegrationClientError("failed to encrypt enygma deposit to dvp batches", err)
	}

	contractProof := EnygmaDvpIntegration.IEnygmaDvpIntegrationWithdrawOrDepositProof{
		PiA:          [2]*big.Int{proof.PiA[0], proof.PiA[1]},
		PiB:          [2][2]*big.Int{{proof.PiB[0][0], proof.PiB[0][1]}, {proof.PiB[1][0], proof.PiB[1][1]}},
		PiC:          [2]*big.Int{proof.PiC[0], proof.PiC[1]},
		PublicSignal: proof.PublicSignal,
	}

	block, err := c.ethClient.BlockByNumber(ctx, blockNumber)
	if err != nil {
		return WrapInEnygmaDvpIntegrationClientError("failed to get block for deposit timestamp", withstack.Wrap(err))
	}

	burnUpdate := types.DvpBalanceUpdated{
		ErcId:             nil,
		TokenType:         uint8(types.ENYGMA),
		ResourceId:        resourceId,
		SourceChainId:     chainId,
		SourceTxHash:      sourceTxHash.Hex(),
		SourceTxTimestamp: time.Unix(int64(block.Time()), 0).UTC(),
		From:              from.Hex(),
		To:                dvpIntegrationAddress.Hex(),
		Amount:            amount,
		UpdateType:        types.Burn,
	}

	encryptedBurnUpdate, err := c.encr.EncryptDvpBalanceUpdated(ctx, burnUpdate)
	if err != nil {
		return WrapInEnygmaDvpIntegrationClientError("failed to encrypt balance update", err)
	}

	calldata := c.contract.PackDepositToDvp(contractProof, encrBatches, encryptedBurnUpdate)

	_, err = c.executor.Execute(ctx, IDFor("dvpintegration.Deposit", chainEventID), calldata, dvpIntegrationAddress)
	if err != nil {
		return WrapInEnygmaDvpIntegrationClientError("failed to execute deposit to dvp", withstack.Wrap(err))
	}

	return nil
}

func (c *DvpIntegrationClient) Withdraw(
	ctx context.Context,
	chainEventID string,
	batches []*types.EnygmaTransferBatch,
	proof *types.EnygmaProofResponse,
	blockNumber *big.Int,
	jsProof *dvp.ProofReceipt,
	chainId *big.Int,
	resourceId string,
	amount *big.Int,
	from common.Address,
	sourceTxHash common.Hash,
	dvpIntegrationAddress common.Address,
) error {
	encrBatches, err := c.encr.EncryptEnygmaTransferBatches(ctx, batches, blockNumber)
	if err != nil {
		return WrapInEnygmaDvpIntegrationClientError("failed to encrypt enygma withdraw from dvp batches", err)
	}

	contractProof := EnygmaDvpIntegration.IEnygmaDvpIntegrationWithdrawOrDepositProof{
		PiA:          [2]*big.Int{proof.PiA[0], proof.PiA[1]},
		PiB:          [2][2]*big.Int{{proof.PiB[0][0], proof.PiB[0][1]}, {proof.PiB[1][0], proof.PiB[1][1]}},
		PiC:          [2]*big.Int{proof.PiC[0], proof.PiC[1]},
		PublicSignal: proof.PublicSignal,
	}

	dvpProof := EnygmaDvpIntegration.IDvpProofReceipt{
		Proof: EnygmaDvpIntegration.IDvpSnarkProof{
			A: EnygmaDvpIntegration.IDvpG1Point{
				X: jsProof.Proof.A[0],
				Y: jsProof.Proof.A[1],
			},
			B: EnygmaDvpIntegration.IDvpG2Point{
				X: [2]*big.Int{
					jsProof.Proof.B[0][0],
					jsProof.Proof.B[0][1],
				},
				Y: [2]*big.Int{
					jsProof.Proof.B[1][0],
					jsProof.Proof.B[1][1],
				},
			},
			C: EnygmaDvpIntegration.IDvpG1Point{
				X: jsProof.Proof.C[0],
				Y: jsProof.Proof.C[1],
			},
		},
		TreeNumbers:      jsProof.TreeNumbers,
		Message:          jsProof.Message,
		MerkleRoots:      jsProof.MerkleRoots,
		Nullifiers:       jsProof.Nullifiers,
		Commitments:      jsProof.Commitments,
		RevertCommitment: jsProof.RevertCommitment,
	}

	block, err := c.ethClient.BlockByNumber(ctx, blockNumber)
	if err != nil {
		return WrapInEnygmaDvpIntegrationClientError("failed to get block for withdrawal timestamp", withstack.Wrap(err))
	}

	mintUpdate := types.DvpBalanceUpdated{
		ErcId:             nil,
		TokenType:         uint8(types.ENYGMA),
		ResourceId:        resourceId,
		SourceChainId:     chainId,
		SourceTxHash:      sourceTxHash.Hex(),
		SourceTxTimestamp: time.Unix(int64(block.Time()), 0).UTC(),
		From:              from.Hex(),
		To:                dvpIntegrationAddress.Hex(),
		Amount:            amount,
		UpdateType:        types.Mint,
	}

	encryptedMintUpdate, err := c.encr.EncryptDvpBalanceUpdated(ctx, mintUpdate)
	if err != nil {
		return WrapInEnygmaDvpIntegrationClientError("failed to encrypt balance update", err)
	}

	calldata := c.contract.PackWithdrawFromDvp(contractProof, encrBatches, dvpProof, encryptedMintUpdate)

	_, err = c.executor.Execute(ctx, IDFor("dvpintegration.Withdraw", chainEventID), calldata, dvpIntegrationAddress)
	if err != nil {
		return WrapInEnygmaDvpIntegrationClientError("failed to execute withdraw from dvp", withstack.Wrap(err))
	}

	return nil
}

func (c *DvpIntegrationClient) ConsolidateFunds(
	ctx context.Context,
	chainEventID string,
	dvpIntegrationAddress common.Address,
	proofReceipt *dvp.ProofReceipt) error {
	contractJSProof := EnygmaDvpIntegration.IDvpProofReceipt{
		Proof: EnygmaDvpIntegration.IDvpSnarkProof{
			A: EnygmaDvpIntegration.IDvpG1Point{
				X: proofReceipt.Proof.A[0],
				Y: proofReceipt.Proof.A[1],
			},
			B: EnygmaDvpIntegration.IDvpG2Point{
				X: [2]*big.Int{
					proofReceipt.Proof.B[0][0],
					proofReceipt.Proof.B[0][1],
				},
				Y: [2]*big.Int{
					proofReceipt.Proof.B[1][0],
					proofReceipt.Proof.B[1][1],
				},
			},
			C: EnygmaDvpIntegration.IDvpG1Point{
				X: proofReceipt.Proof.C[0],
				Y: proofReceipt.Proof.C[1],
			},
		},
		Message:          proofReceipt.Message,
		TreeNumbers:      proofReceipt.TreeNumbers,
		MerkleRoots:      proofReceipt.MerkleRoots,
		Nullifiers:       proofReceipt.Nullifiers,
		Commitments:      proofReceipt.Commitments,
		RevertCommitment: proofReceipt.RevertCommitment,
	}

	calldata := c.contract.PackConsolidateFunds(contractJSProof)

	// ConsoldateFunds is called recursively for the same message ID.
	// Hash the calldata to get a unique identifier for this particular loop.
	calldataHash := sha256.Sum256(calldata)
	_, err := c.executor.Execute(ctx, IDFor("dvpintegration.ConsolidateFunds", chainEventID, common.Bytes2Hex(calldataHash[:])), calldata, dvpIntegrationAddress)
	if err != nil {
		return WrapInEnygmaDvpIntegrationClientError("failed to consolidate funds", withstack.Wrap(err))
	}

	return nil
}
