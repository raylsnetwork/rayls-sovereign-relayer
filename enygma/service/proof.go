package service

import (
	"context"
	"fmt"
	"log/slog"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/iden3/go-iden3-crypto/babyjub"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/conv"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/cryptography"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/cts/gen/proto/keys"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/enygma"
	telemetry "github.com/raylsnetwork/rayls-privacy-relayer-api/otel"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/types"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
)

// poseidonHashRandomSeed is the seed value used to compute the HashRandom constant via Poseidon hash.
const poseidonHashRandomSeed = 21

// var _ proofTracer = (*adapters.OTelTracer)(nil)

type proofTracer interface {
	Start(ctx context.Context, spanName string, opts ...trace.SpanStartOption) (context.Context, trace.Span)
}

// Ports for ProofService
type proofKOSClient interface {
	GenerateEnygmaSharedSecrets(
		ctx context.Context,
		in *keys.GenerateEnygmaSharedSecretsRequest,
		opts ...grpc.CallOption,
	) (*keys.GenerateEnygmaSharedSecretsResponse, error)
	GetPaymentSpendKey(
		ctx context.Context,
		in *keys.GetPaymentSpendKeyRequest,
		opts ...grpc.CallOption,
	) (*keys.PaymentSpendKeyResponse, error)
}

type proofAPIClient interface {
	CreateTransferProof(anonymityIndex int, request types.TransferProofRequest) (types.EnygmaProofResponse, error)
	CreateDepositProof(anonymityIndex int, request types.DepositProofRequest) (types.EnygmaProofResponse, error)
	CreateWithdrawProof(anonymityIndex int, request types.WithdrawProofRequest) (types.EnygmaProofResponse, error)
}

type proofEnygmaRepository interface {
	GetEnygmaByResourceId(ctx context.Context, resourceId string) (types.Enygma, error)
}

type ProofEnygmaClient interface {
	GetPublicValuesFinalised(ctx context.Context, tokenAddress common.Address) (*types.EnygmaPublicValues, error)
}

type EnygmaProofService struct {
	plChainId        *big.Int
	kosClient        proofKOSClient
	proofClient      proofAPIClient
	enygmaRepository proofEnygmaRepository
	enygmaClient     ProofEnygmaClient
	tracer           proofTracer
}

func NewEnygmaProofService(
	plChainId *big.Int,
	kosClient proofKOSClient,
	proofClient proofAPIClient,
	enygmaRepository proofEnygmaRepository,
	enygmaClient ProofEnygmaClient,
	tracer proofTracer,
) *EnygmaProofService {
	return &EnygmaProofService{
		plChainId:        plChainId,
		kosClient:        kosClient,
		proofClient:      proofClient,
		enygmaRepository: enygmaRepository,
		enygmaClient:     enygmaClient,
		tracer:           tracer,
	}
}

func (s *EnygmaProofService) GenerateTransferProof(
	ctx context.Context,
	params enygma.TransferProofParams,
) (*types.EnygmaProofResponse, []*types.Point, []*big.Int, []*big.Int, error) {
	ctx, span := s.tracer.Start(ctx, telemetry.SPAN_GENERATE_TRANSFER_PROOF)
	defer span.End()

	span.SetAttributes(
		attribute.Int(telemetry.ATTR_ANONYMITY_INDEX, params.AnonymityIndex),
		attribute.String(telemetry.ATTR_SENDER_AMOUNT, params.SenderAmount.String()),
		attribute.String(telemetry.ATTR_SENDER_CHAIN_ID, s.plChainId.String()),
	)

	if err := s.validateAnonymityIndex(params.AnonymityIndex, params.Batches); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, telemetry.STATUS_ERROR_ANONYMITY_INDEX_VALIDATION_FAILED)
		return nil, nil, nil, nil, fmt.Errorf("validating transfer proof anonymity index: %w", err)
	}

	commonParams := commonProofParams{
		ResourceId:     params.ResourceId,
		AnonymityIndex: params.AnonymityIndex,
		SenderChainId:  s.plChainId,
		SenderAmount:   params.SenderAmount,
		BlockNumber:    params.BlockNumber,
		Batches:        params.Batches,
		TokenAddress:   params.TokenAddress,
	}

	commonData, err := s.prepareCommonData(ctx, commonParams, ProofTypeTransfer)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, telemetry.STATUS_ERROR_FAILED_TO_PREPARE_COMMON_PROOF_DATA)
		return nil, nil, nil, nil, fmt.Errorf("preparing transfer proof common data: %w", err)
	}

	commonProofRequest := s.buildCommonProofRequest(commonParams, commonData)

	// TODO: consider moving tracing inside CreateTransferProof
	_, proofSpan := s.tracer.Start(ctx, telemetry.SPAN_CREATE_TRANSFER_PROOF)
	proof, err := s.proofClient.CreateTransferProof(params.AnonymityIndex, types.TransferProofRequest{
		CommonProofRequest: commonProofRequest,
	})
	proofSpan.End()
	if err != nil {
		proofSpan.RecordError(err)
		proofSpan.SetStatus(codes.Error, telemetry.STATUS_ERROR_PROOF_CREATION_FAILED)

		span.RecordError(err)
		span.SetStatus(codes.Error, telemetry.STATUS_ERROR_FAILED_TO_CREATE_TRANSFER_PROOF)
		return nil, nil, nil, nil, fmt.Errorf("creating transfer proof: %w", err)
	}
	proofSpan.SetStatus(codes.Ok, "")

	span.SetStatus(codes.Ok, telemetry.STATUS_SUCCESS_TRANSFER_PROOF_GENERATED)
	return &proof, commonData.DestinationNewCommitments, commonData.DestinationRandomFactors, commonData.MessageTags, nil
}

func (s *EnygmaProofService) GenerateDepositProof(
	ctx context.Context,
	params enygma.DepositProofParams,
) (*types.EnygmaProofResponse, []*types.Point, []*big.Int, []*big.Int, error) {
	if err := s.validateAnonymityIndex(params.AnonymityIndex, params.Batches); err != nil {
		return nil, nil, nil, nil, fmt.Errorf("validating deposit proof anonymity index: %w", err)
	}

	commonParams := commonProofParams{
		ResourceId:     params.ResourceId,
		AnonymityIndex: params.AnonymityIndex,
		SenderChainId:  s.plChainId,
		SenderAmount:   params.SenderAmount,
		BlockNumber:    params.BlockNumber,
		Batches:        params.Batches,
		TokenAddress:   params.TokenAddress,
	}

	commonData, err := s.prepareCommonData(ctx, commonParams, ProofTypeDeposit)
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("preparing deposit proof common data: %w", err)
	}

	commonProofRequest := s.buildCommonProofRequest(commonParams, commonData)

	proof, err := s.proofClient.CreateDepositProof(params.AnonymityIndex, types.DepositProofRequest{
		CommonProofRequest: commonProofRequest,
		TokenAddress:       params.TokenAddress,
		DepositCommitment:  params.DepositCommitment,
		DepositSalt:        params.DepositSalt,
		DepositPublicKey:   params.DepositPublicKey,
	})
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("creating deposit proof: %w", err)
	}

	return &proof, commonData.DestinationNewCommitments, commonData.DestinationRandomFactors, commonData.MessageTags, nil
}

func (s *EnygmaProofService) GenerateWithdrawProof(
	ctx context.Context,
	params enygma.WithdrawProofParams,
) (*types.EnygmaProofResponse, []*types.Point, []*big.Int, []*big.Int, error) {
	if err := s.validateAnonymityIndex(params.AnonymityIndex, params.Batches); err != nil {
		return nil, nil, nil, nil, fmt.Errorf("validating withdraw proof anonymity index: %w", err)
	}

	commonParams := commonProofParams{
		ResourceId:     params.ResourceId,
		AnonymityIndex: params.AnonymityIndex,
		SenderChainId:  s.plChainId,
		SenderAmount:   params.SenderAmount,
		BlockNumber:    params.BlockNumber,
		Batches:        params.Batches,
		TokenAddress:   params.TokenAddress,
	}

	commonData, err := s.prepareCommonData(ctx, commonParams, ProofTypeWithdraw)
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("preparing withdraw proof common data: %w", err)
	}
	commonProofRequest := s.buildCommonProofRequest(commonParams, commonData)

	proof, err := s.proofClient.CreateWithdrawProof(params.AnonymityIndex, types.WithdrawProofRequest{
		CommonProofRequest: commonProofRequest,
		TokenAddress:       params.TokenAddress,
		DepositCommitments: params.DepositCommitments,
		DepositSecretKeys:  params.DepositSecretKeys,
		DepositAmounts:     params.DepositAmounts,
		DepositSalts:       params.DepositSalts,
	})
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("creating withdraw proof: %w", err)
	}

	return &proof, commonData.DestinationNewCommitments, commonData.DestinationRandomFactors, commonData.MessageTags, nil
}

type commonProofParams struct {
	ResourceId     string
	AnonymityIndex int
	SenderChainId  *big.Int
	SenderAmount   *big.Int
	BlockNumber    *big.Int
	Batches        []*types.EnygmaTransferBatch
	TokenAddress   common.Address
}

type commonProofData struct {
	SenderPreviousBalance          *big.Int
	SenderPreviousRandomFactor     *big.Int
	SenderSecretKey                *big.Int
	Nullifier                      *big.Int
	DestinationAmounts             []*big.Int
	DestinationChainIDs            []*big.Int
	DestinationSharedSecrets       []*big.Int // k array: secrets[senderIdx] = Poseidon(previousR, sk), others from shared secrets
	DestinationRandomFactors       []*big.Int
	DestinationPublicKeys          []*big.Int
	DestinationNewCommitments      []*types.Point
	DestinationPreviousCommitments []*types.Point
	ArrayHashSecrets               []*big.Int // k array: arrayHashSecret[i] = Poseidon(secrets[i], secrets[i])
	MessageTags                    []*big.Int
}

type ProofType int

const (
	ProofTypeTransfer ProofType = iota
	ProofTypeDeposit
	ProofTypeWithdraw
)

// Private helper methods
func (s *EnygmaProofService) validateAnonymityIndex(anonymityIndex int, batches []*types.EnygmaTransferBatch) error {
	uniqueChainIDs := make(map[string]bool)
	for _, batch := range batches {
		uniqueChainIDs[batch.ToChainID.String()] = true
	}

	// Validate that unique chain IDs count does not exceed anonymity index
	if len(uniqueChainIDs) > anonymityIndex {
		return fmt.Errorf(
			"transaction not supported yet. Only transactions from 1->1 up to 1->k-1 are supported (unique destination chains: %d, anonymity index: %d)",
			len(uniqueChainIDs),
			anonymityIndex,
		)
	}
	return nil
}

func (s *EnygmaProofService) prepareCommonData(
	ctx context.Context,
	params commonProofParams,
	proofType ProofType,
) (*commonProofData, error) {
	destinationChainIDs := make([]*big.Int, len(params.Batches))
	for i, batch := range params.Batches {
		destinationChainIDs[i] = batch.ToChainID
	}

	// Log destination chain IDs for debugging message tag matching
	destChainIDsStr := make([]string, len(destinationChainIDs))
	for i, id := range destinationChainIDs {
		destChainIDsStr[i] = id.String()
	}
	slog.Debug("prepareCommonData: preparing proof data",
		slog.String("senderChainId", params.SenderChainId.String()),
		slog.String("blockNumber", params.BlockNumber.String()),
		slog.Any("destinationChainIDs", destChainIDsStr),
	)

	// TODO: consider moving tracing inside GetPaymentSpendKey
	ctx, keySpan := s.tracer.Start(ctx, telemetry.SPAN_GET_ENYGMA_KEY)
	resp, err := s.kosClient.GetPaymentSpendKey(ctx, &keys.GetPaymentSpendKeyRequest{})
	keySpan.End()
	if err != nil {
		keySpan.RecordError(err)
		keySpan.SetStatus(codes.Error, telemetry.STATUS_ERROR_FAILED_TO_GET_ENYGMA_KEY)
		return nil, fmt.Errorf("retrieving payment spend key: %w", err)
	}
	keySpan.SetStatus(codes.Ok, "")

	enygmaState, err := s.enygmaRepository.GetEnygmaByResourceId(ctx, params.ResourceId)
	if err != nil {
		return nil, fmt.Errorf("fetching enygma state by resource ID: %w", err)
	}

	publicValues, err := s.enygmaClient.GetPublicValuesFinalised(ctx, params.TokenAddress)
	if err != nil {
		return nil, fmt.Errorf("retrieving finalized public values: %w", err)
	}

	// Order of public values must be the same as the order of destinationChainIDs
	destPublicKeys := make([]*big.Int, len(destinationChainIDs))
	prevCommitments := make([]*types.Point, len(destinationChainIDs))
	for i, chainId := range destinationChainIDs {
		chainIdStr := chainId.String()
		pk := publicValues.PublicKeys[chainIdStr]
		if pk == nil {
			return nil, fmt.Errorf("missing public key for destination chain %s in enygma public values", chainIdStr)
		}
		destPublicKeys[i] = pk

		cm := publicValues.Commitments[chainIdStr]
		if cm == nil {
			return nil, fmt.Errorf("missing commitment for destination chain %s in enygma public values", chainIdStr)
		}
		prevCommitments[i] = cm
	}

	// Find sender index for validation
	senderIndex := -1
	for i, chainID := range destinationChainIDs {
		if chainID.Cmp(params.SenderChainId) == 0 {
			senderIndex = i
			break
		}
	}

	// Validate that our DB state matches the on-chain commitment before proceeding
	// This prevents wasted work sending proofs that will fail validation
	if senderIndex >= 0 && prevCommitments[senderIndex] != nil {
		expectedCommitment := cryptography.PedersenCommitmentEnygma(
			enygmaState.FinalizedBalance,
			enygmaState.FinalizedR,
		)
		onChainCommitment := prevCommitments[senderIndex]

		if expectedCommitment.X.Cmp(onChainCommitment.X) != 0 || expectedCommitment.Y.Cmp(onChainCommitment.Y) != 0 {
			slog.Warn("balance commitment mismatch between SC and DB - sync service needs to update DB state",
				slog.String("resourceId", params.ResourceId),
				slog.String("dbBalance", enygmaState.FinalizedBalance.String()),
				slog.String("expectedCommitmentX", expectedCommitment.X.String()),
				slog.String("expectedCommitmentY", expectedCommitment.Y.String()),
				slog.String("onChainCommitmentX", onChainCommitment.X.String()),
				slog.String("onChainCommitmentY", onChainCommitment.Y.String()),
			)
			return nil, fmt.Errorf("balance commitment mismatch between SC and DB")
		}
	}

	destAmounts := s.sumAmountsByChainId(params.SenderChainId, params.SenderAmount, params.Batches, proofType)

	sharedSecretsResp, err := s.kosClient.GenerateEnygmaSharedSecrets(
		ctx,
		&keys.GenerateEnygmaSharedSecretsRequest{
			ChainIDs:      conv.BigIntsToUint64s(destinationChainIDs),
			BlockNumber:   params.BlockNumber.Uint64(),
			SenderChainID: params.SenderChainId.Uint64(),
			PrevRFactor:   enygmaState.FinalizedR.Bytes(),
			ResourceID:    common.FromHex(params.ResourceId),
		},
	)
	if err != nil {
		return nil, fmt.Errorf("generating enygma shared secrets: %w", err)
	}

	secrets := bytesSliceToBigInts(sharedSecretsResp.GetSecrets())
	arrayHashSecrets := bytesSliceToBigInts(sharedSecretsResp.GetHashSecrets())
	messageTags := bytesSliceToBigInts(sharedSecretsResp.GetMessageTags())

	// Use secrets (not sharedSecrets) - circuit uses secrets[i] for ALL participants including sender
	randomFactors := s.genRandomFactors(params.SenderChainId, params.BlockNumber, destinationChainIDs, secrets)
	newCommitments := s.genCommitments(
		params.SenderChainId,
		params.SenderAmount,
		destinationChainIDs,
		destAmounts,
		randomFactors,
		proofType,
	)

	// Generate nullifier: Poseidon(arrayHashSecret[senderIdx], blockNumber)
	// Circuit uses raw Poseidon output (no ModHintBabyJubJub), which is naturally in JubJubPrimeGroup
	nullifierIdx := senderIndex
	nullifierInputs := []*big.Int{arrayHashSecrets[nullifierIdx], params.BlockNumber}
	nullifier, err := cryptography.GetPoseidonHashModNumber(nullifierInputs, cryptography.JubJubPrimeGroup)
	if err != nil {
		return nil, fmt.Errorf("computing poseidon hash for nullifier: %w", err)
	}

	return &commonProofData{
		SenderSecretKey:                new(big.Int).SetBytes(resp.SecretKey),
		SenderPreviousBalance:          enygmaState.FinalizedBalance,
		SenderPreviousRandomFactor:     enygmaState.FinalizedR,
		DestinationSharedSecrets:       secrets,
		DestinationRandomFactors:       randomFactors,
		DestinationAmounts:             destAmounts,
		DestinationChainIDs:            destinationChainIDs,
		DestinationPublicKeys:          destPublicKeys,
		DestinationPreviousCommitments: prevCommitments,
		DestinationNewCommitments:      newCommitments,
		Nullifier:                      nullifier,
		ArrayHashSecrets:               arrayHashSecrets,
		MessageTags:                    messageTags,
	}, nil
}

func (s *EnygmaProofService) buildCommonProofRequest(
	params commonProofParams,
	commonData *commonProofData,
) *types.CommonProofRequest {
	return &types.CommonProofRequest{
		ResourceId:                 params.ResourceId,
		SenderSecretKey:            commonData.SenderSecretKey,
		SenderBalance:              commonData.SenderPreviousBalance,
		SenderRandomFactor:         commonData.SenderPreviousRandomFactor,
		SenderChainId:              params.SenderChainId,
		SenderAmount:               params.SenderAmount,
		DestinationChainIDs:        commonData.DestinationChainIDs,
		DestinationAmounts:         commonData.DestinationAmounts,
		DestinationPublicKeys:      commonData.DestinationPublicKeys,
		DestinationPreviousCommits: commonData.DestinationPreviousCommitments,
		DestinationNewCommits:      commonData.DestinationNewCommitments,
		DestinationRandomFactors:   commonData.DestinationRandomFactors,
		DestinationSharedSecrets:   commonData.DestinationSharedSecrets,
		Nullifier:                  commonData.Nullifier,
		BlockNumber:                params.BlockNumber,
		ArrayHashSecrets:           commonData.ArrayHashSecrets,
		MessageTags:                commonData.MessageTags,
	}
}

// genRandomFactors calculates random values for Enygma proof generation
// Matches circuit logic:
// HashRandom = Poseidon(21) - RAW, not modded
// For receivers: RandomFactor[i] = neg(Poseidon(HashRandom, secrets[i], BlockNumber) mod r) = r - hash
// For sender: RandomFactor[sender] = sum of all receiver hashes (before negation)
// Constraint: sender_r + sum(all_other_r) = 0 mod r
func (s *EnygmaProofService) genRandomFactors(
	senderChainID *big.Int,
	blockNumber *big.Int,
	chainIDs []*big.Int,
	secrets []*big.Int,
) []*big.Int {
	rValues := make([]*big.Int, len(chainIDs))
	receiverHashesModP := make([]*big.Int, len(chainIDs))
	sumOfReceiverHashes := big.NewInt(0)

	// Compute HashRandom = Poseidon(21) - RAW, not reduced mod r
	// Circuit uses raw hash as input to next Poseidon, only final result is mod r
	hashRandom, err := cryptography.GetPoseidonHash([]*big.Int{big.NewInt(poseidonHashRandomSeed)})
	if err != nil {
		slog.Error("genRandomFactors: failed to compute hashRandom", slog.Any("error", err))
		hashRandom = big.NewInt(0)
	}

	// First pass: compute all hashes and sum of receiver hashes
	for i, chainID := range chainIDs {
		// RandomFactor = Poseidon(HashRandom, secrets[i], BlockNumber) mod r
		secret := secrets[i]
		if secret == nil {
			slog.Error("genRandomFactors: secret is nil", slog.Int("i", i))
			secret = big.NewInt(0)
		}
		inputs := []*big.Int{hashRandom, secret, blockNumber}
		hash, err := cryptography.GetPoseidonHashModNumber(inputs, cryptography.JubJubPrimeSubGroup)
		if err != nil {
			slog.Error("genRandomFactors: GetPoseidonHashModNumber failed", slog.Any("error", err))
			hash = big.NewInt(0)
		}
		receiverHashesModP[i] = hash

		// Add to sum only if this is a receiver (not sender)
		if chainID.Cmp(senderChainID) != 0 {
			sumOfReceiverHashes = cryptography.AddMod(sumOfReceiverHashes, hash, cryptography.JubJubPrimeSubGroup)
		}
	}

	// Second pass: assign the correct random factors based on role
	for i, chainID := range chainIDs {
		if chainID.Cmp(senderChainID) == 0 {
			// For sender: r_Sender = Sum(receiver hashes) mod p
			rValues[i] = sumOfReceiverHashes
		} else {
			// For receivers: r_Receiver_i = neg(hash_i mod p) = p - hash_i
			rValues[i] = cryptography.SubMod(
				cryptography.JubJubPrimeSubGroup,
				receiverHashesModP[i],
				cryptography.JubJubPrimeSubGroup,
			)
		}
	}

	return rValues
}

// SumAmountsByChainId calculates transfer amounts for each destination chain
func (s *EnygmaProofService) sumAmountsByChainId(
	senderChainId *big.Int,
	senderAmount *big.Int,
	batches []*types.EnygmaTransferBatch,
	proofType ProofType,
) []*big.Int {
	amounts := make([]*big.Int, 0, len(batches))

	// Loop through the destination chain IDs to ensure the right order of the values.
	for _, batch := range batches {
		amount := big.NewInt(0)

		for _, tx := range batch.Transactions {
			amount.Add(amount, tx.ToAmount)
		}

		if batch.ToChainID.Cmp(senderChainId) == 0 {
			if proofType == ProofTypeTransfer || proofType == ProofTypeDeposit {
				// Calculate (r - senderAmount) mod r to handle edge case where senderAmount = 0
				amount = cryptography.SubMod(
					cryptography.JubJubPrimeSubGroup,
					senderAmount,
					cryptography.JubJubPrimeSubGroup,
				)
			}
			if proofType == ProofTypeWithdraw {
				amount = senderAmount
			}
		}

		amounts = append(amounts, amount)
	}

	return amounts
}

// genCommitments creates cryptographic commitments for the transaction
func (s *EnygmaProofService) genCommitments(
	senderID *big.Int,
	senderAmount *big.Int,
	destinationChainIDs []*big.Int,
	destinationAmounts []*big.Int,
	randomFactors []*big.Int,
	proofType ProofType,
) []*types.Point {
	// TODO: use Point instead of babyjub.Point
	jubjubCommitments := make([]*babyjub.Point, 0, len(destinationChainIDs))

	for i, chainID := range destinationChainIDs {
		txCommit := cryptography.PedersenCommitmentEnygma(destinationAmounts[i], randomFactors[i])

		if chainID.Cmp(senderID) == 0 {
			if proofType == ProofTypeTransfer || proofType == ProofTypeDeposit {
				txCommit = cryptography.PedersenCommitmentEnygma(
					cryptography.GetNegative(senderAmount),
					randomFactors[i],
				)
			}
			if proofType == ProofTypeWithdraw {
				txCommit = cryptography.PedersenCommitmentEnygma(senderAmount, randomFactors[i])
			}
		}

		jubjubCommitments = append(jubjubCommitments, txCommit)
	}

	commitments := make([]*types.Point, len(destinationChainIDs))
	for i := 0; i < len(destinationChainIDs); i++ {
		commitments[i] = &types.Point{X: jubjubCommitments[i].X, Y: jubjubCommitments[i].Y}
	}

	return commitments
}

func bytesSliceToBigInts(slices [][]byte) []*big.Int {
	result := make([]*big.Int, len(slices))
	for i, b := range slices {
		result[i] = new(big.Int).SetBytes(b)
	}
	return result
}
