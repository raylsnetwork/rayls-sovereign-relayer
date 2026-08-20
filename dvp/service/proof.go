package service

//go:generate moq --skip-ensure --pkg service_test -out proof_mock_test.go . proofMerkleClient proofKeysClient proofAPIClient proofDepositRepository proofCommitmentCalculator proofTxManager

import (
	"context"
	"fmt"
	"math/big"

	"github.com/raylsnetwork/rayls-sovereign-relayer/client"
	"github.com/raylsnetwork/rayls-sovereign-relayer/cryptography"
	keyspb "github.com/raylsnetwork/rayls-sovereign-relayer/cts/gen/proto/keys"
	"github.com/raylsnetwork/rayls-sovereign-relayer/dvp"
	merkle "github.com/raylsnetwork/rayls-sovereign-relayer/merkle-service"
	"github.com/raylsnetwork/rayls-sovereign-relayer/private-relayer/repository"
	"github.com/raylsnetwork/rayls-sovereign-relayer/types"
	"google.golang.org/grpc"
)

type ProofServiceConfig struct {
	ChainID              *big.Int
	MerkleTreeDepth      int
	NumberOfJSParamsIn   int
	numberOfJSParamsOut  int
	numberOfOwnParamsIn  int
	numberOfOwnParamsOut int
}

// merkleClient interface defines methods for merkle proof generation
type proofMerkleClient interface {
	GenerateMerkleProof(ctx context.Context, commitment *big.Int, treeNumber int, tokenAddress string) (*merkle.MerkleProof, error)
}

// proofKeysClient defines methods for key operations against the CTS gRPC keys service.
type proofKeysClient interface {
	GetPaymentSpendKey(ctx context.Context, in *keyspb.GetPaymentSpendKeyRequest, opts ...grpc.CallOption) (*keyspb.PaymentSpendKeyResponse, error)
}

// paymentSpendKeyFromResponse converts the gRPC PaymentSpendKeyResponse (raw bytes)
// into the *big.Int-shaped types.PaymentSpendKey that the proof request types still expect.
func paymentSpendKeyFromResponse(resp *keyspb.PaymentSpendKeyResponse) types.PaymentSpendKey {
	return types.PaymentSpendKey{
		SecretKey: new(big.Int).SetBytes(resp.GetSecretKey()),
		PublicKey: new(big.Int).SetBytes(resp.GetPublicKey()),
	}
}

// proofClient interface defines methods for cryptographic proof generation
type proofAPIClient interface {
	CreateEnygmaJSProof(req *dvp.EnygmaJoinSplitProofRequest) (dvp.Proof, error)
	CreateErc1155JSProof(req *dvp.ERC1155JoinSplitProofRequest) (dvp.Proof, error)
	CreateErc721OwnershipProof(req *dvp.ERC721OwnershipProofRequest) (dvp.Proof, error)
}

// depositRepository interface defines methods for deposit persistence
type proofDepositRepository interface {
	CreateDeposit(ctx context.Context, deposit *types.DvpDeposit) error
	GetFungibleDeposits(
		ctx context.Context,
		tokenAddress string,
		userAddress string,
		tokenType types.DvpTokenType,
		status types.DvpDepositStatus,
	) ([]types.DvpDeposit, error)
	GetDepositsByToken(
		ctx context.Context,
		tokenAddress string,
		tokenId string,
		tokenType types.DvpTokenType,
		userAddress string,
		status types.DvpDepositStatus,
	) ([]types.DvpDeposit, error)
	BatchUpsertNullifiers(ctx context.Context, commitmentNullifierMap map[string]string) error
	UpsertDepositNullifier(ctx context.Context, commitment *big.Int, nullifier *big.Int) error
}

type proofCommitmentCalculator interface {
	GetNFTUniqueID(nftAddress string, nftID string) (*big.Int, error)
	CalculateNFTCommitment(spendPK, salt *big.Int, nftID string, nftAddress string) (*big.Int, error)
	CalculateERC1155Commitment(spendPK, salt *big.Int, tokenAddress string, tokenID string, tokenAmount *big.Int) (*big.Int, error)
	CalculatePaymentCommitment(spendPK, salt, paymentAmount *big.Int, tokenAddress string) (*big.Int, error)
}

type proofTxManager interface {
	WithTransaction(ctx context.Context, fn func(ctx context.Context) error) error
}

var (
	_ proofDepositRepository    = (*repository.DvpDepositRepository)(nil)
	_ proofCommitmentCalculator = (*CommitmentCalculator)(nil)
	_ proofMerkleClient         = (*merkle.MerkleService)(nil)
	_ proofAPIClient            = (*client.ProofAPIClient)(nil)
	_ proofTxManager            = (*repository.TransactionManager)(nil)
)

type ProofService struct {
	conf                 ProofServiceConfig
	merkleClient         proofMerkleClient
	keysClient           proofKeysClient
	proofClient          proofAPIClient
	commitmentCalculator proofCommitmentCalculator
	depositRepository    proofDepositRepository
	txManager            proofTxManager
}

func NewProofService(
	conf ProofServiceConfig,
	merkleClient proofMerkleClient,
	keysClient proofKeysClient,
	proofClient proofAPIClient,
	commitmentCalculator proofCommitmentCalculator,
	depositRepository proofDepositRepository,
	txManager proofTxManager,
) *ProofService {
	conf.numberOfJSParamsOut = 2  // Always 2 outputs(receive amount and change amount)
	conf.numberOfOwnParamsIn = 1  // Always 1
	conf.numberOfOwnParamsOut = 1 // Always 1

	return &ProofService{
		conf:                 conf,
		merkleClient:         merkleClient,
		keysClient:           keysClient,
		proofClient:          proofClient,
		commitmentCalculator: commitmentCalculator,
		depositRepository:    depositRepository,
		txManager:            txManager,
	}
}

func (s *ProofService) GenerateERC1155JSProof(
	ctx context.Context,
	sourceViewPublicKey []byte,
	commitmentMessage *big.Int,
	destinationPaymentPublicKey *big.Int,
	destinationSalt *big.Int,
	paymentAmount *big.Int,
	userAddress string,
	tokenAddress string,
	tokenId string,
	deposits []*types.DvpDeposit,
) (*dvp.ProofReceipt, error) {
	params, err := s.prepareCommonJSParams(ctx, deposits, paymentAmount, destinationPaymentPublicKey)
	if err != nil {
		return nil, fmt.Errorf("prepare common JS params for ERC1155: %w", err)
	}

	tokenID, ok := big.NewInt(0).SetString(tokenId, 10)
	if !ok {
		return nil, fmt.Errorf("failed to parse erc1155 token id: %s", tokenId)
	}

	changeSalt, _, err := cryptography.GenerateSalt(sourceViewPublicKey)
	if err != nil {
		return nil, fmt.Errorf("generate salt for ERC1155 change deposit: %w", err)
	}

	revertSalt, _, err := cryptography.GenerateSalt(sourceViewPublicKey)
	if err != nil {
		return nil, fmt.Errorf("generate salt for ERC1155 revert deposit: %w", err)
	}

	proof, err := s.proofClient.CreateErc1155JSProof(
		&dvp.ERC1155JoinSplitProofRequest{
			ERC1155Address: tokenAddress,
			ERC1155TokenId: tokenID,
			NftCommitment:  commitmentMessage,
			ValuesIn:       params.ValuesIn,
			SaltsIn:        params.SaltsIn,
			KeyPairsIn:     params.KeyPairsIn,
			ValuesOut:      params.ValuesOut,
			SaltsOut:       []*big.Int{destinationSalt, changeSalt},
			PubKeysOut:     params.PubKeysOut,
			MerkleDepth:    params.MerkleDepth,
			MerkleProofs:   params.MerkleProofs,
			MerkleRoots:    params.MerkleRoots,
			TreeNumbers:    params.TreeNumbers,
			RevertSalt:     revertSalt,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("create ERC1155 JS proof: %w", err)
	}

	receipt, err := s.proofToReceipt(&proof, s.conf.NumberOfJSParamsIn, s.conf.numberOfJSParamsOut)
	if err != nil {
		return nil, fmt.Errorf("parse ERC1155 JS proof receipt: %w", err)
	}

	if params.ChangeAmount.Cmp(big.NewInt(0)) > 0 {
		// TODO: Get the change commitment from the proof, instead of re-computing it.
		changeCommitment, err := s.commitmentCalculator.CalculateERC1155Commitment(
			params.ChangePublicKey,
			changeSalt,
			params.ChangeTokenAddress,
			tokenId,
			params.ChangeAmount,
		)
		if err != nil {
			return nil, fmt.Errorf("calculate ERC1155 change commitment: %w", err)
		}

		err = s.depositRepository.CreateDeposit(
			ctx,
			&types.DvpDeposit{
				UserAddress:  userAddress,
				Salt:         changeSalt,
				TokenAmount:  params.ChangeAmount,
				TokenAddress: tokenAddress,
				TokenType:    types.DvpERC1155,
				TokenID:      tokenId,
				Commitment:   changeCommitment,
				Status:       types.DvpDepositPending,
			},
		)
		if err != nil {
			return nil, fmt.Errorf("create ERC1155 change deposit: %w", err)
		}
	}

	err = s.depositRepository.CreateDeposit(
		ctx,
		&types.DvpDeposit{
			UserAddress:  userAddress,
			Salt:         revertSalt,
			TokenAmount:  params.TotalAmount,
			TokenAddress: tokenAddress,
			TokenType:    types.DvpERC1155,
			TokenID:      tokenId,
			Commitment:   receipt.RevertCommitment,
			Status:       types.DvpDepositPending,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("create ERC1155 revert deposit: %w", err)
	}

	return receipt, nil
}

func (s *ProofService) GenerateEnygmaJSProof(
	ctx context.Context,
	sourceViewPublicKey []byte,
	nftCommitment *big.Int,
	destinationPaymentPublicKey *big.Int,
	destinationSalt *big.Int,
	paymentAmount *big.Int,
	tokenAddress string,
	deposits []*types.DvpDeposit,
) (*dvp.ProofReceipt, error) {
	params, err := s.prepareCommonJSParams(ctx, deposits, paymentAmount, destinationPaymentPublicKey)
	if err != nil {
		return nil, fmt.Errorf("prepare common JS params for Enygma: %w", err)
	}

	changeSalt, _, err := cryptography.GenerateSalt(sourceViewPublicKey)
	if err != nil {
		return nil, fmt.Errorf("generate salt for Enygma change deposit: %w", err)
	}

	revertSalt, _, err := cryptography.GenerateSalt(sourceViewPublicKey)
	if err != nil {
		return nil, fmt.Errorf("generate salt for Enygma revert deposit: %w", err)
	}

	jsProof, err := s.proofClient.CreateEnygmaJSProof(
		&dvp.EnygmaJoinSplitProofRequest{
			NftCommitment: nftCommitment,
			ValuesIn:      params.ValuesIn,
			SaltsIn:       params.SaltsIn,
			KeyPairsIn:    params.KeyPairsIn,
			ValuesOut:     params.ValuesOut,
			SaltsOut:      []*big.Int{destinationSalt, changeSalt},
			PubKeysOut:    params.PubKeysOut,
			MerkleDepth:   params.MerkleDepth,
			MerkleProofs:  params.MerkleProofs,
			MerkleRoots:   params.MerkleRoots,
			TreeNumbers:   params.TreeNumbers,
			ERC20Address:  tokenAddress,
			RevertSalt:    revertSalt,
		})
	if err != nil {
		return nil, fmt.Errorf("create Enygma JS proof: %w", err)
	}

	receipt, err := s.proofToReceipt(&jsProof, s.conf.NumberOfJSParamsIn, s.conf.numberOfJSParamsOut)
	if err != nil {
		return nil, fmt.Errorf("parse Enygma JS proof receipt: %w", err)
	}

	// Create the change deposit if needed (for both consolidation and non-consolidation)
	if params.ChangeAmount.Cmp(big.NewInt(0)) > 0 {
		// TODO: Get the change commitment from the proof, instead of re-computing it.
		changeCommitment, err := s.commitmentCalculator.CalculatePaymentCommitment(
			params.ChangePublicKey,
			changeSalt,
			params.ChangeAmount,
			params.ChangeTokenAddress,
		)
		if err != nil {
			return nil, fmt.Errorf("calculate Enygma change commitment: %w", err)
		}

		err = s.depositRepository.CreateDeposit(
			ctx,
			&types.DvpDeposit{
				UserAddress:  deposits[0].UserAddress,
				Salt:         changeSalt,
				TokenAmount:  params.ChangeAmount,
				TokenAddress: tokenAddress,
				TokenType:    types.DvpEnygma,
				TokenID:      deposits[0].TokenID,
				Commitment:   changeCommitment,
				Status:       types.DvpDepositPending,
			},
		)
		if err != nil {
			return nil, fmt.Errorf("create Enygma change deposit: %w", err)
		}
	}

	if len(deposits) > 0 {
		err = s.depositRepository.CreateDeposit(
			ctx,
			&types.DvpDeposit{
				UserAddress:  deposits[0].UserAddress,
				Salt:         revertSalt,
				TokenAmount:  params.TotalAmount,
				TokenAddress: tokenAddress,
				TokenType:    types.DvpEnygma,
				TokenID:      deposits[0].TokenID,
				Commitment:   receipt.RevertCommitment,
				Status:       types.DvpDepositPending,
			},
		)
		if err != nil {
			return nil, fmt.Errorf("create Enygma revert deposit: %w", err)
		}
	}

	return receipt, nil
}

func (s *ProofService) GenerateOwnershipProof(
	ctx context.Context,
	sourceViewPublicKey []byte,
	paymentCommitment *big.Int,
	destinationPaymentPublicKey *big.Int,
	destinationSalt *big.Int,
	deposit *types.DvpDeposit,
) (*dvp.ProofReceipt, error) {
	nftUID, err := s.commitmentCalculator.GetNFTUniqueID(deposit.TokenAddress, deposit.TokenID)
	if err != nil {
		return nil, fmt.Errorf("get NFT unique ID: %w", err)
	}

	depositMerkleProof, err := s.merkleClient.GenerateMerkleProof(
		ctx,
		deposit.Commitment,
		int(deposit.TreeNumber),
		deposit.TokenAddress,
	)
	if err != nil {
		return nil, fmt.Errorf("generate merkle proof for ownership: %w", err)
	}

	spendKeyResp, err := s.keysClient.GetPaymentSpendKey(ctx, &keyspb.GetPaymentSpendKeyRequest{})
	if err != nil {
		return nil, fmt.Errorf("get payment spend key for ownership: %w", err)
	}
	spendKey := paymentSpendKeyFromResponse(spendKeyResp)

	revertSalt, _, err := cryptography.GenerateSalt(sourceViewPublicKey)
	if err != nil {
		return nil, fmt.Errorf("generate salt for ERC721 revert deposit: %w", err)
	}

	ownershipProof, err := s.proofClient.CreateErc721OwnershipProof(
		&dvp.ERC721OwnershipProofRequest{
			PaymentCommitment: paymentCommitment,
			UID:               nftUID,
			KeyPairIn:         &spendKey,
			SaltIn:            deposit.Salt,
			PubKeyOut:         destinationPaymentPublicKey,
			SaltOut:           destinationSalt,
			TreeNumber:        deposit.TreeNumber,
			MerkleDepth:       s.conf.MerkleTreeDepth,
			MerkleRoot:        depositMerkleProof.Root,
			MerkleProof: &types.MerkleProof{
				Indices:  depositMerkleProof.Indices,
				Elements: depositMerkleProof.Elements,
			},
			RevertSalt: revertSalt,
		})
	if err != nil {
		return nil, fmt.Errorf("create ERC721 ownership proof: %w", err)
	}

	receipt, err := s.proofToReceipt(&ownershipProof, s.conf.numberOfOwnParamsIn, s.conf.numberOfOwnParamsOut)
	if err != nil {
		return nil, fmt.Errorf("parse ERC721 ownership proof receipt: %w", err)
	}

	err = s.depositRepository.CreateDeposit(
		ctx,
		&types.DvpDeposit{
			UserAddress:  deposit.UserAddress,
			Salt:         revertSalt,
			TokenAmount:  deposit.TokenAmount,
			TokenAddress: deposit.TokenAddress,
			TokenType:    types.DvpERC721,
			TokenID:      deposit.TokenID,
			Commitment:   receipt.RevertCommitment,
			Status:       types.DvpDepositPending,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("create ERC721 revert deposit: %w", err)
	}

	return receipt, nil
}

func (s *ProofService) proofToReceipt(proof *dvp.Proof, numberOfParamsIn int, numberOfParamsOut int) (*dvp.ProofReceipt, error) {
	expectedLen := 2 + 3*numberOfParamsIn + numberOfParamsOut
	if len(proof.PublicSignal) < expectedLen {
		return nil, fmt.Errorf("public signal length %d is less than expected %d", len(proof.PublicSignal), expectedLen)
	}

	idx := 0
	message := proof.PublicSignal[idx]
	idx++

	merkleRootsFromSignal := proof.PublicSignal[idx : idx+numberOfParamsIn]
	idx += numberOfParamsIn

	nullifiersFromSignal := proof.PublicSignal[idx : idx+numberOfParamsIn]
	idx += numberOfParamsIn

	treeNumbersFromSignal := proof.PublicSignal[idx : idx+numberOfParamsIn]
	idx += numberOfParamsIn

	commitmentsFromSignal := proof.PublicSignal[idx : idx+numberOfParamsOut]
	idx += numberOfParamsOut

	revertCommitment := proof.PublicSignal[idx]

	receipt := &dvp.ProofReceipt{
		Proof:            proof,
		Message:          message,
		TreeNumbers:      treeNumbersFromSignal,
		MerkleRoots:      merkleRootsFromSignal,
		Nullifiers:       nullifiersFromSignal,
		Commitments:      commitmentsFromSignal,
		RevertCommitment: revertCommitment,
	}

	return receipt, nil
}

func (s *ProofService) prepareCommonJSParams(
	ctx context.Context,
	deposits []*types.DvpDeposit,
	paymentAmount *big.Int,
	destinationPublicKey *big.Int,
) (*commonJSProofParams, error) {
	var params commonJSProofParams
	params.MerkleDepth = s.conf.MerkleTreeDepth

	spendKeyResp, err := s.keysClient.GetPaymentSpendKey(ctx, &keyspb.GetPaymentSpendKeyRequest{})
	if err != nil {
		return nil, fmt.Errorf("get payment spend key: %w", err)
	}
	spendKey := paymentSpendKeyFromResponse(spendKeyResp)

	totalAmount := big.NewInt(0)

	for _, deposit := range deposits {
		depositProof, err := s.merkleClient.GenerateMerkleProof(
			ctx,
			deposit.Commitment,
			int(deposit.TreeNumber),
			deposit.TokenAddress,
		)
		if err != nil {
			return nil, fmt.Errorf("generate merkle proof for deposit: %w", err)
		}

		totalAmount = totalAmount.Add(totalAmount, deposit.TokenAmount)

		params.ValuesIn = append(params.ValuesIn, deposit.TokenAmount)
		params.SaltsIn = append(params.SaltsIn, deposit.Salt)
		params.KeyPairsIn = append(params.KeyPairsIn, &spendKey)
		params.TreeNumbers = append(params.TreeNumbers, deposit.TreeNumber)
		params.MerkleRoots = append(params.MerkleRoots, depositProof.Root)
		params.MerkleProofs = append(params.MerkleProofs, &types.MerkleProof{
			Elements: depositProof.Elements,
			Indices:  depositProof.Indices,
		})
	}

	changeAmount := new(big.Int).Sub(totalAmount, paymentAmount)
	if changeAmount.Sign() < 0 {
		return nil, fmt.Errorf("payment amount %s exceeds total deposit amount %s", paymentAmount, totalAmount)
	}

	params.TotalAmount = totalAmount
	params.ValuesOut = append(params.ValuesOut, paymentAmount)
	params.ValuesOut = append(params.ValuesOut, changeAmount)
	params.PubKeysOut = append(params.PubKeysOut, destinationPublicKey)
	params.PubKeysOut = append(params.PubKeysOut, spendKey.PublicKey)

	params.ChangeAmount = changeAmount
	params.ChangePublicKey = spendKey.PublicKey
	if len(deposits) > 0 {
		params.ChangeTokenAddress = deposits[0].TokenAddress
	}
	return &params, nil
}

type commonJSProofParams struct {
	ChangeAmount       *big.Int
	ChangePublicKey    *big.Int
	ChangeTokenAddress string
	TotalAmount        *big.Int
	ValuesIn           []*big.Int
	SaltsIn            []*big.Int
	KeyPairsIn         []*types.PaymentSpendKey
	ValuesOut          []*big.Int
	PubKeysOut         []*big.Int
	TreeNumbers        []int
	MerkleProofs       []*types.MerkleProof
	MerkleRoots        []*big.Int
	MerkleDepth        int
}

// GenerateEnygmaToERC721SwapProof builds the join-split proof that spends
// source Enygma deposits to mint the self-destination ERC721 commitment,
// then registers the pending ERC721 deposit and upserts the input nullifiers.
func (s *ProofService) GenerateEnygmaToERC721SwapProof(
	ctx context.Context,
	swap *types.DvpSwap,
	deposits []*types.DvpDeposit,
	sourceViewPublicKey []byte,
	selfSalt *big.Int,
	destSalt *big.Int,
	destSpendPubKey *big.Int,
) (*dvp.ProofReceipt, error) {
	spendKeyResp, err := s.keysClient.GetPaymentSpendKey(ctx, &keyspb.GetPaymentSpendKeyRequest{})
	if err != nil {
		return nil, fmt.Errorf("getting payment spend key: %w", err)
	}
	spendKey := paymentSpendKeyFromResponse(spendKeyResp)

	// Self destination commitment
	nftCommitment, err := s.commitmentCalculator.CalculateNFTCommitment(spendKey.PublicKey, selfSalt, swap.TokenOutID, swap.TokenOutAddress)
	if err != nil {
		return nil, fmt.Errorf("calculating NFT commitment: %w", err)
	}

	proof, err := s.GenerateEnygmaJSProof(ctx, sourceViewPublicKey, nftCommitment, destSpendPubKey, destSalt, swap.TokenInAmount, swap.TokenInAddress, deposits)
	if err != nil {
		return nil, fmt.Errorf("generating enygma JS proof: %w", err)
	}

	// Register the ERC721 token as a pending deposit.
	// It will be available for withdrawal once the swap is completed.
	err = s.depositRepository.CreateDeposit(ctx,
		&types.DvpDeposit{
			UserAddress:  swap.From,
			Salt:         selfSalt,
			TokenAmount:  swap.TokenOutAmount,
			TokenAddress: swap.TokenOutAddress,
			TokenType:    swap.TokenOutType,
			TokenID:      swap.TokenOutID,
			Commitment:   nftCommitment,
			Status:       types.DvpDepositPending,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("creating pending ERC721 deposit: %w", err)
	}

	commitmentNullifierMap := make(map[string]string)
	for i, deposit := range deposits {
		commitmentNullifierMap[deposit.Commitment.String()] = proof.Nullifiers[i].String()
	}

	err = s.depositRepository.BatchUpsertNullifiers(ctx, commitmentNullifierMap)
	if err != nil {
		return nil, fmt.Errorf("batch upserting nullifiers: %w", err)
	}

	return proof, nil
}

// GenerateEnygmaToERC1155SwapProof builds the join-split proof that spends
// source Enygma deposits to mint the self-destination ERC1155 commitment,
// then registers the pending ERC1155 deposit and upserts the input nullifiers.
func (s *ProofService) GenerateEnygmaToERC1155SwapProof(
	ctx context.Context,
	swap *types.DvpSwap,
	deposits []*types.DvpDeposit,
	sourceViewPublicKey []byte,
	selfSalt *big.Int,
	destSalt *big.Int,
	destSpendPubKey *big.Int,
) (*dvp.ProofReceipt, error) {
	spendKeyResp, err := s.keysClient.GetPaymentSpendKey(ctx, &keyspb.GetPaymentSpendKeyRequest{})
	if err != nil {
		return nil, fmt.Errorf("getting payment spend key: %w", err)
	}
	spendKey := paymentSpendKeyFromResponse(spendKeyResp)

	// Self-Destination commitment
	erc1155Commitment, err := s.commitmentCalculator.CalculateERC1155Commitment(spendKey.PublicKey, selfSalt, swap.TokenOutAddress, swap.TokenOutID, swap.TokenOutAmount)
	if err != nil {
		return nil, fmt.Errorf("calculating ERC1155 commitment: %w", err)
	}

	proof, err := s.GenerateEnygmaJSProof(ctx, sourceViewPublicKey, erc1155Commitment, destSpendPubKey, destSalt, swap.TokenInAmount, swap.TokenInAddress, deposits)
	if err != nil {
		return nil, fmt.Errorf("generating enygma JS proof: %w", err)
	}

	// Register the ERC1155 token as a pending deposit.
	err = s.depositRepository.CreateDeposit(ctx,
		&types.DvpDeposit{
			UserAddress:  swap.From,
			Salt:         selfSalt,
			TokenAmount:  swap.TokenOutAmount,
			TokenAddress: swap.TokenOutAddress,
			TokenType:    swap.TokenOutType,
			TokenID:      swap.TokenOutID,
			Commitment:   erc1155Commitment,
			Status:       types.DvpDepositPending,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("creating pending ERC1155 deposit: %w", err)
	}

	commitmentNullifierMap := make(map[string]string)
	for i, deposit := range deposits {
		commitmentNullifierMap[deposit.Commitment.String()] = proof.Nullifiers[i].String()
	}

	err = s.depositRepository.BatchUpsertNullifiers(ctx, commitmentNullifierMap)
	if err != nil {
		return nil, fmt.Errorf("batch upserting nullifiers: %w", err)
	}

	return proof, nil
}

// GenerateERC721ToEnygmaSwapProof builds the ownership proof that spends the
// source ERC721 deposit for the self-destination Enygma payment commitment.
// The input nullifier and pending Enygma deposit are persisted atomically.
func (s *ProofService) GenerateERC721ToEnygmaSwapProof(
	ctx context.Context,
	swap *types.DvpSwap,
	deposit *types.DvpDeposit,
	sourceViewPublicKey []byte,
	selfSalt *big.Int,
	destSalt *big.Int,
	destSpendPubKey *big.Int,
) (*dvp.ProofReceipt, error) {
	spendKeyResp, err := s.keysClient.GetPaymentSpendKey(ctx, &keyspb.GetPaymentSpendKeyRequest{})
	if err != nil {
		return nil, fmt.Errorf("getting payment spend key: %w", err)
	}
	spendKey := paymentSpendKeyFromResponse(spendKeyResp)

	// Self-Destination commitment
	paymentCommitment, err := s.commitmentCalculator.CalculatePaymentCommitment(spendKey.PublicKey, selfSalt, swap.TokenOutAmount, swap.TokenOutAddress)
	if err != nil {
		return nil, fmt.Errorf("calculating payment commitment: %w", err)
	}

	ownProof, err := s.GenerateOwnershipProof(ctx,
		sourceViewPublicKey,
		paymentCommitment,
		destSpendPubKey,
		destSalt,
		deposit,
	)
	if err != nil {
		return nil, fmt.Errorf("generating ownership proof: %w", err)
	}

	// Lock deposit, set nullifier, and create the pending Enygma deposit atomically
	// right before the teleport message. If the relayer crashes before this transaction,
	// no deposits are locked. If it crashes after, the teleport will be retried on
	// event redelivery and the duplicate-key check on the swap prevents double-processing.
	err = s.txManager.WithTransaction(ctx, func(txCtx context.Context) error {
		if txErr := s.depositRepository.UpsertDepositNullifier(txCtx, deposit.Commitment, ownProof.Nullifiers[0]); txErr != nil {
			return txErr
		}
		return s.depositRepository.CreateDeposit(txCtx,
			&types.DvpDeposit{
				UserAddress:  swap.From,
				Salt:         selfSalt,
				TokenAmount:  swap.TokenOutAmount,
				TokenAddress: swap.TokenOutAddress,
				TokenType:    swap.TokenOutType,
				TokenID:      swap.TokenOutID,
				Commitment:   paymentCommitment,
				Status:       types.DvpDepositPending,
			},
		)
	})
	if err != nil {
		return nil, fmt.Errorf("creating pending enygma deposit: %w", err)
	}

	return ownProof, nil
}

// GenerateERC1155ToEnygmaSwapProof builds the ERC1155 join-split proof that
// spends the source ERC1155 deposits for the self-destination Enygma payment
// commitment. Input nullifiers and the pending Enygma deposit are persisted
// atomically.
func (s *ProofService) GenerateERC1155ToEnygmaSwapProof(
	ctx context.Context,
	swap *types.DvpSwap,
	deposits []*types.DvpDeposit,
	sourceViewPublicKey []byte,
	selfSalt *big.Int,
	destSalt *big.Int,
	destSpendPubKey *big.Int,
) (*dvp.ProofReceipt, error) {
	spendKeyResp, err := s.keysClient.GetPaymentSpendKey(ctx, &keyspb.GetPaymentSpendKeyRequest{})
	if err != nil {
		return nil, fmt.Errorf("getting payment spend key: %w", err)
	}
	spendKey := paymentSpendKeyFromResponse(spendKeyResp)

	// Self-Destination commitment
	paymentCommitment, err := s.commitmentCalculator.CalculatePaymentCommitment(spendKey.PublicKey, selfSalt, swap.TokenOutAmount, swap.TokenOutAddress)
	if err != nil {
		return nil, fmt.Errorf("calculating payment commitment: %w", err)
	}

	jsProof, err := s.GenerateERC1155JSProof(ctx, sourceViewPublicKey, paymentCommitment, destSpendPubKey, destSalt, swap.TokenInAmount, swap.From, swap.TokenInAddress, swap.TokenInID, deposits)
	if err != nil {
		return nil, fmt.Errorf("generating ERC1155 JS proof: %w", err)
	}

	// Lock deposits, set nullifiers, and create the pending Enygma deposit atomically.
	commitmentNullifierMap := make(map[string]string)
	for i, deposit := range deposits {
		commitmentNullifierMap[deposit.Commitment.String()] = jsProof.Nullifiers[i].String()
	}
	err = s.txManager.WithTransaction(ctx, func(txCtx context.Context) error {
		if txErr := s.depositRepository.BatchUpsertNullifiers(txCtx, commitmentNullifierMap); txErr != nil {
			return txErr
		}
		return s.depositRepository.CreateDeposit(txCtx,
			&types.DvpDeposit{
				UserAddress:  swap.From,
				Salt:         selfSalt,
				TokenAmount:  swap.TokenOutAmount,
				TokenAddress: swap.TokenOutAddress,
				TokenType:    swap.TokenOutType,
				TokenID:      swap.TokenOutID,
				Commitment:   paymentCommitment,
				Status:       types.DvpDepositPending,
			},
		)
	})
	if err != nil {
		return nil, fmt.Errorf("locking deposits and creating pending enygma deposit: %w", err)
	}

	return jsProof, nil
}

// GenerateERC721WithdrawProof builds the ownership proof for an ERC721
// withdrawal. The destination is the DvP operator (no self-destination
// commitment to register), so no DB state changes are applied here — the
// caller is responsible for locking the input deposit and upserting its
// nullifier when it's ready to broadcast the withdrawal tx.
func (s *ProofService) GenerateERC721WithdrawProof(
	ctx context.Context,
	sourceViewPublicKey []byte,
	destSalt *big.Int,
	operatorPublicKey *big.Int,
	deposit *types.DvpDeposit,
) (*dvp.ProofReceipt, error) {
	return s.GenerateOwnershipProof(ctx,
		sourceViewPublicKey,
		big.NewInt(0),
		operatorPublicKey,
		destSalt,
		deposit,
	)
}

// GenerateERC1155WithdrawProof builds the ERC1155 join-split proof for an
// ERC1155 withdrawal. The destination is the DvP operator (paymentCommitment
// is zero). The caller is responsible for locking input deposits and
// upserting their nullifiers.
func (s *ProofService) GenerateERC1155WithdrawProof(
	ctx context.Context,
	sourceViewPublicKey []byte,
	destSalt *big.Int,
	operatorPublicKey *big.Int,
	userAddress string,
	tokenAddress string,
	tokenID string,
	tokenAmount *big.Int,
	deposits []*types.DvpDeposit,
) (*dvp.ProofReceipt, error) {
	return s.GenerateERC1155JSProof(ctx,
		sourceViewPublicKey,
		big.NewInt(0),
		operatorPublicKey,
		destSalt,
		tokenAmount,
		userAddress,
		tokenAddress,
		tokenID,
		deposits,
	)
}
