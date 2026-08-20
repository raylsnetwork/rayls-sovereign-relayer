package client

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"math/big"
	"net/http"
	"net/url"

	"github.com/raylsnetwork/rayls-sovereign-relayer/dvp"
	"github.com/raylsnetwork/rayls-sovereign-relayer/types"
	"github.com/raylsnetwork/rayls-sovereign-relayer/withstack"
)

const (
	// Proof input counts for different token type circuits.
	proofInputCountK6 = 6
	proofInputCountK5 = 5
	proofInputCountK4 = 4
	proofInputCountK3 = 3
)

type ProofAPIClient struct {
	baseURL                  string
	enygmaTransferK2         string
	enygmaTransferK3         string
	enygmaTransferK4         string
	enygmaTransferK5         string
	enygmaTransferK6         string
	enygmaDepositK2          string
	enygmaDepositK3          string
	enygmaDepositK4          string
	enygmaDepositK5          string
	enygmaDepositK6          string
	enygmaWithdrawK2         string
	enygmaWithdrawK3         string
	enygmaWithdrawK4         string
	enygmaWithdrawK5         string
	enygmaWithdrawK6         string
	enygmaJSEndpoint         string
	erc721OwnershipEndpoint  string
	erc1155JSEndpoint        string
	erc1155OwnershipEndpoint string
	// Proof generation is
	// computationally expensive, so this should be a non-retry client (e.g., SimpleClient)
	// to avoid re-triggering expensive computations on transient failures.
	httpClient HTTPClient
}

func NewProofAPIClient(baseURL string, httpClient HTTPClient) *ProofAPIClient {
	return &ProofAPIClient{
		baseURL:                  baseURL,
		enygmaTransferK2:         types.PROOF_API_ENYGMA_2,
		enygmaTransferK3:         types.PROOF_API_ENYGMA_3,
		enygmaTransferK4:         types.PROOF_API_ENYGMA_4,
		enygmaTransferK5:         types.PROOF_API_ENYGMA_5,
		enygmaTransferK6:         types.PROOF_API_ENYGMA_6,
		enygmaDepositK2:          types.PROOF_API_ENYGMA_DEPOSIT_2,
		enygmaDepositK3:          types.PROOF_API_ENYGMA_DEPOSIT_3,
		enygmaDepositK4:          types.PROOF_API_ENYGMA_DEPOSIT_4,
		enygmaDepositK5:          types.PROOF_API_ENYGMA_DEPOSIT_5,
		enygmaDepositK6:          types.PROOF_API_ENYGMA_DEPOSIT_6,
		enygmaWithdrawK2:         types.PROOF_API_ENYGMA_WITHDRAW_2,
		enygmaWithdrawK3:         types.PROOF_API_ENYGMA_WITHDRAW_3,
		enygmaWithdrawK4:         types.PROOF_API_ENYGMA_WITHDRAW_4,
		enygmaWithdrawK5:         types.PROOF_API_ENYGMA_WITHDRAW_5,
		enygmaWithdrawK6:         types.PROOF_API_ENYGMA_WITHDRAW_6,
		enygmaJSEndpoint:         types.PROOF_API_JOIN_SPLIT_ENYGMA,
		erc721OwnershipEndpoint:  types.PROOF_API_OWNERSHIP_721,
		erc1155JSEndpoint:        types.PROOF_API_JOIN_SPLIT_1155,
		erc1155OwnershipEndpoint: types.PROOF_API_OWNERSHIP_1155,
		httpClient:               httpClient,
	}
}

func (c *ProofAPIClient) CreateErc721OwnershipProof(req *dvp.ERC721OwnershipProofRequest) (dvp.Proof, error) {
	slog.Info("Creating ERC721 Ownership Proof")

	proofRequest := ProofAPIErc721OwnershipRequestParams{
		PaymentCommitment: req.PaymentCommitment.String(),
		UID:               req.UID.String(),
		KeyPairIn: ProofAPIDvpSpendKeys{
			PublicKey:  req.KeyPairIn.PublicKey.String(),
			PrivateKey: req.KeyPairIn.SecretKey.String(),
		},
		SaltIn: req.SaltIn.String(),
		PubKeyOut: ProofAPIDvpSpendKeys{
			PublicKey: req.PubKeyOut.String(),
		},
		SaltOut:       req.SaltOut.String(),
		RevertSalt:  req.RevertSalt.String(),
		MerkleDepth: req.MerkleDepth,
		MerkleProof: ProofAPIDvpMerkleProof{
			Indices:  req.MerkleProof.Indices.String(),
			Elements: convertBigIntArrayToStringArray(req.MerkleProof.Elements),
		},
		MerkleRoot: req.MerkleRoot.String(),
		TreeNumber: req.TreeNumber,
	}

	reqURL := buildRequestURL(c.baseURL, c.erc721OwnershipEndpoint, url.Values{})
	request, err := newJSONRequest(http.MethodPost, reqURL, proofRequest)
	if err != nil {
		slog.Error("Failed to create request", "error", err.Error())
		return dvp.Proof{}, fmt.Errorf("creating erc721 ownership proof request: %w", err)
	}

	response, err := c.httpClient.Do(request)
	if err != nil {
		slog.Error("HTTP request failed", "error", err.Error())
		return dvp.Proof{}, fmt.Errorf("executing erc721 ownership proof request: %w", err)
	}
	switch response.StatusCode {
	case http.StatusOK:
		return handleSuccessfulErc721OwnershipRequest(response)
	case http.StatusInternalServerError:
		err := handleInternalServerError(response)
		return dvp.Proof{}, fmt.Errorf("creating erc721 ownership proof: %w", err)
	default:
		err := handleUnsupportedStatus(response)
		return dvp.Proof{}, fmt.Errorf("creating erc721 ownership proof: %w", err)
	}
}

func (c *ProofAPIClient) CreateEnygmaJSProof(req *dvp.EnygmaJoinSplitProofRequest) (dvp.Proof, error) {
	slog.Info("Creating Enygma JS Proof")

	pubKeysOut := make([]ProofAPIDvpSpendKeys, len(req.PubKeysOut))
	for i, pubKey := range req.PubKeysOut {
		pubKeysOut[i] = ProofAPIDvpSpendKeys{PublicKey: pubKey.String()}
	}

	proofRequest := ProofAPIEnygmaJoinSplitRequestParams{
		NFTCommitment: req.NftCommitment.String(),
		TreeNumbers:   req.TreeNumbers,
		MerkleDepth:   req.MerkleDepth,
		ERC20Address:  req.ERC20Address,
		RevertSalt: req.RevertSalt.String(),
		ValuesIn:      convertBigIntArrayToStringArray(req.ValuesIn),
		SaltsIn:       convertBigIntArrayToStringArray(req.SaltsIn),
		ValuesOut:     convertBigIntArrayToStringArray(req.ValuesOut),
		SaltsOut:      convertBigIntArrayToStringArray(req.SaltsOut),
		KeyPairsIn:    convertPaymentSpendKeysArrayToStringArray(req.KeyPairsIn),
		PubKeysOut:    pubKeysOut,
		MerkleProofs:  convertMerkleProofArrayToStringArray(req.MerkleProofs),
		MerkleRoots:   convertBigIntArrayToStringArray(req.MerkleRoots),
	}

	reqURL := buildRequestURL(c.baseURL, c.enygmaJSEndpoint, url.Values{})
	request, err := newJSONRequest(http.MethodPost, reqURL, proofRequest)
	if err != nil {
		slog.Error("Failed to create request", "error", err.Error())
		return dvp.Proof{}, fmt.Errorf("creating enygma join-split proof request: %w", err)
	}

	response, err := c.httpClient.Do(request)
	if err != nil {
		slog.Error("HTTP request failed", "error", err.Error())
		return dvp.Proof{}, fmt.Errorf("executing enygma join-split proof request: %w", err)
	}
	switch response.StatusCode {
	case http.StatusOK:
		return handleSuccessfulEnygmaJSRequest(response)
	case http.StatusInternalServerError:
		err := handleInternalServerError(response)
		return dvp.Proof{}, fmt.Errorf("creating enygma join-split proof: %w", err)
	default:
		err := handleUnsupportedStatus(response)
		return dvp.Proof{}, fmt.Errorf("creating enygma join-split proof: %w", err)
	}
}

func (c *ProofAPIClient) CreateErc1155JSProof(req *dvp.ERC1155JoinSplitProofRequest) (dvp.Proof, error) {
	slog.Info("Creating ERC1155 JS Proof")

	pubKeysOut := make([]ProofAPIDvpSpendKeys, len(req.PubKeysOut))
	for i, pubKey := range req.PubKeysOut {
		pubKeysOut[i] = ProofAPIDvpSpendKeys{PublicKey: pubKey.String()}
	}

	proofRequest := ProofAPIDvpERC1155JoinSplitRequestParams{
		NFTCommitment:  req.NftCommitment.String(),
		TreeNumbers:    req.TreeNumbers,
		MerkleDepth:    req.MerkleDepth,
		ERC1155Address: req.ERC1155Address,
		ERC1155TokenId: req.ERC1155TokenId.String(),
		RevertSalt: req.RevertSalt.String(),
		ValuesIn:       convertBigIntArrayToStringArray(req.ValuesIn),
		SaltsIn:        convertBigIntArrayToStringArray(req.SaltsIn),
		ValuesOut:      convertBigIntArrayToStringArray(req.ValuesOut),
		SaltsOut:       convertBigIntArrayToStringArray(req.SaltsOut),
		KeyPairsIn:     convertPaymentSpendKeysArrayToStringArray(req.KeyPairsIn),
		PubKeysOut:     pubKeysOut,
		MerkleProofs:   convertMerkleProofArrayToStringArray(req.MerkleProofs),
		MerkleRoots:    convertBigIntArrayToStringArray(req.MerkleRoots),
	}

	reqURL := buildRequestURL(c.baseURL, c.erc1155JSEndpoint, url.Values{})
	request, err := newJSONRequest(http.MethodPost, reqURL, proofRequest)
	if err != nil {
		slog.Error("Failed to create request", "error", err.Error())
		return dvp.Proof{}, fmt.Errorf("creating erc1155 join-split proof request: %w", err)
	}

	response, err := c.httpClient.Do(request)
	if err != nil {
		slog.Error("HTTP request failed", "error", err.Error())
		return dvp.Proof{}, fmt.Errorf("executing erc1155 join-split proof request: %w", err)
	}
	switch response.StatusCode {
	case http.StatusOK:
		return handleSuccessfulErc1155JSRequest(response)
	case http.StatusInternalServerError:
		err := handleInternalServerError(response)
		return dvp.Proof{}, fmt.Errorf("creating erc1155 join-split proof: %w", err)
	default:
		err := handleUnsupportedStatus(response)
		return dvp.Proof{}, fmt.Errorf("creating erc1155 join-split proof: %w", err)
	}
}

func (c *ProofAPIClient) CreateTransferProof(k int, req types.TransferProofRequest) (types.EnygmaProofResponse, error) {
	slog.Info("Creating Transfer Proof")

	// Select endpoint based on k value
	var endpoint string
	switch k {
	case proofInputCountK6:
		endpoint = c.enygmaTransferK6
	case proofInputCountK5:
		endpoint = c.enygmaTransferK5
	case proofInputCountK4:
		endpoint = c.enygmaTransferK4
	case proofInputCountK3:
		endpoint = c.enygmaTransferK3
	case 2:
		endpoint = c.enygmaTransferK2
	default:
		return types.EnygmaProofResponse{}, fmt.Errorf("unsupported k value: %d, only k=2,3,4,5 and 6 are supported", k)
	}

	// Convert domain types to API types
	proofRequest := ProofAPITransferRequest{
		SenderID:                  req.SenderChainId.String(),
		SenderTxValue:             req.SenderAmount.String(),
		SecretKey:                 req.SenderSecretKey.String(),
		PreviousSenderBalance:     req.SenderBalance.String(),
		PreviousSenderRandomValue: req.SenderRandomFactor.String(),
		Nullifier:                 req.Nullifier.String(),
		BlockNumber:               req.BlockNumber.String(),
		PublicKeys:                convertBigIntArrayToStringArray(req.DestinationPublicKeys),
		PreviousCommits:           convertPointsToStringArrays(req.DestinationPreviousCommits),
		TxCommits:                 convertPointsToStringArrays(req.DestinationNewCommits),
		TxValues:                  convertBigIntArrayToStringArray(req.DestinationAmounts),
		TxRandomValues:            convertBigIntArrayToStringArray(req.DestinationRandomFactors),
		SharedSecrets:             convertBigIntArrayToStringArray(req.DestinationSharedSecrets),
		HashedSharedSecrets:       convertBigIntArrayToStringArray(req.ArrayHashSecrets),
		MessageTags:               convertBigIntArrayToStringArray(req.MessageTags),
		AnonymitySet:              convertBigIntArrayToStringArray(req.DestinationChainIDs),
	}

	reqURL := buildRequestURL(c.baseURL, endpoint, url.Values{})

	request, err := newJSONRequest(http.MethodPost, reqURL, proofRequest)
	if err != nil {
		slog.Error("Failed to create transfer proof request", "error", err.Error())
		return types.EnygmaProofResponse{}, fmt.Errorf("creating transfer proof request: %w", err)
	}

	response, err := c.httpClient.Do(request)
	if err != nil {
		slog.Error("HTTP request failed for transfer proof", "error", err.Error())
		return types.EnygmaProofResponse{}, fmt.Errorf("executing transfer proof request: %w", err)
	}

	switch response.StatusCode {
	case http.StatusOK:
		return handleSuccessfulTransferProofRequest(response)
	case http.StatusInternalServerError:
		err := handleInternalServerError(response)
		return types.EnygmaProofResponse{}, fmt.Errorf("creating transfer proof: %w", err)
	default:
		err := handleUnsupportedStatus(response)
		return types.EnygmaProofResponse{}, fmt.Errorf("creating transfer proof: %w", err)
	}
}

func (c *ProofAPIClient) CreateDepositProof(k int, req types.DepositProofRequest) (types.EnygmaProofResponse, error) {
	slog.Info("Creating Deposit Proof")

	// Select endpoint based on k value
	var endpoint string
	switch k {
	case proofInputCountK6:
		endpoint = c.enygmaDepositK6
	case proofInputCountK5:
		endpoint = c.enygmaDepositK5
	case proofInputCountK4:
		endpoint = c.enygmaDepositK4
	case proofInputCountK3:
		endpoint = c.enygmaDepositK3
	case 2:
		endpoint = c.enygmaDepositK2
	default:
		return types.EnygmaProofResponse{}, fmt.Errorf("unsupported k value: %d, only k=2,3,4,5 and 6 are supported", k)
	}

	// Convert domain types to API types
	proofRequest := ProofAPIDepositRequest{
		SenderID:                  req.SenderChainId.String(),
		SenderTxValue:             req.SenderAmount.String(),
		SecretKey:                 req.SenderSecretKey.String(),
		PreviousSenderBalance:     req.SenderBalance.String(),
		PreviousSenderRandomValue: req.SenderRandomFactor.String(),
		Nullifier:                 req.Nullifier.String(),
		BlockNumber:               req.BlockNumber.String(),
		PublicKeys:                convertBigIntArrayToStringArray(req.DestinationPublicKeys),
		PreviousCommits:           convertPointsToStringArrays(req.DestinationPreviousCommits),
		TxCommits:                 convertPointsToStringArrays(req.DestinationNewCommits),
		TxValues:                  convertBigIntArrayToStringArray(req.DestinationAmounts),
		TxRandomValues:            convertBigIntArrayToStringArray(req.DestinationRandomFactors),
		SharedSecrets:             convertBigIntArrayToStringArray(req.DestinationSharedSecrets),
		HashedSharedSecrets:       convertBigIntArrayToStringArray(req.ArrayHashSecrets),
		MessageTags:               convertBigIntArrayToStringArray(req.MessageTags),
		AnonymitySet:              convertBigIntArrayToStringArray(req.DestinationChainIDs),
		Hash:                      req.DepositCommitment.String(),
		Pk:                        req.DepositPublicKey.String(),
		Address:                   req.TokenAddress.Big().String(),
		SaltOut:                   req.DepositSalt.String(),
	}

	reqURL := buildRequestURL(c.baseURL, endpoint, url.Values{})

	request, err := newJSONRequest(http.MethodPost, reqURL, proofRequest)
	if err != nil {
		slog.Error("Failed to create deposit proof request", "error", err.Error())
		return types.EnygmaProofResponse{}, fmt.Errorf("creating deposit proof request: %w", err)
	}

	response, err := c.httpClient.Do(request)
	if err != nil {
		slog.Error("HTTP request failed for deposit proof", "error", err.Error())
		return types.EnygmaProofResponse{}, fmt.Errorf("executing deposit proof request: %w", err)
	}

	switch response.StatusCode {
	case http.StatusOK:
		return handleSuccessfulDepositProofRequest(response)
	case http.StatusInternalServerError:
		err := handleInternalServerError(response)
		return types.EnygmaProofResponse{}, fmt.Errorf("creating deposit proof: %w", err)
	default:
		err := handleUnsupportedStatus(response)
		return types.EnygmaProofResponse{}, fmt.Errorf("creating deposit proof: %w", err)
	}
}

func (c *ProofAPIClient) CreateWithdrawProof(k int, req types.WithdrawProofRequest) (types.EnygmaProofResponse, error) {
	slog.Info("Creating Withdraw Proof")

	// Select endpoint based on k value
	var endpoint string
	switch k {
	case proofInputCountK6:
		endpoint = c.enygmaWithdrawK6
	case proofInputCountK5:
		endpoint = c.enygmaWithdrawK5
	case proofInputCountK4:
		endpoint = c.enygmaWithdrawK4
	case proofInputCountK3:
		endpoint = c.enygmaWithdrawK3
	case 2:
		endpoint = c.enygmaWithdrawK2
	default:
		return types.EnygmaProofResponse{}, fmt.Errorf("unsupported k value: %d, only k=2,3,4,5 and 6 are supported", k)
	}

	// Convert domain types to API types
	proofRequest := ProofAPIWithdrawRequest{
		SenderID:                  req.SenderChainId.String(),
		SenderTxValue:             req.SenderAmount.String(),
		SecretKey:                 req.SenderSecretKey.String(),
		PreviousSenderBalance:     req.SenderBalance.String(),
		PreviousSenderRandomValue: req.SenderRandomFactor.String(),
		Nullifier:                 req.Nullifier.String(),
		BlockNumber:               req.BlockNumber.String(),
		PublicKeys:                convertBigIntArrayToStringArray(req.DestinationPublicKeys),
		PreviousCommits:           convertPointsToStringArrays(req.DestinationPreviousCommits),
		TxCommits:                 convertPointsToStringArrays(req.DestinationNewCommits),
		TxValues:                  convertBigIntArrayToStringArray(req.DestinationAmounts),
		TxRandomValues:            convertBigIntArrayToStringArray(req.DestinationRandomFactors),
		SharedSecrets:             convertBigIntArrayToStringArray(req.DestinationSharedSecrets),
		HashedSharedSecrets:       convertBigIntArrayToStringArray(req.ArrayHashSecrets),
		MessageTags:               convertBigIntArrayToStringArray(req.MessageTags),
		AnonymitySet:              convertBigIntArrayToStringArray(req.DestinationChainIDs),
		Hashes:                    convertBigIntArrayToStringArray(req.DepositCommitments),
		SkDeposits:                convertBigIntArrayToStringArray(req.DepositSecretKeys),
		VPerDeposit:               convertBigIntArrayToStringArray(req.DepositAmounts),
		Address:                   req.TokenAddress.Big().String(),
		SaltsIn:                   convertBigIntArrayToStringArray(req.DepositSalts),
	}

	reqURL := buildRequestURL(c.baseURL, endpoint, url.Values{})

	request, err := newJSONRequest(http.MethodPost, reqURL, proofRequest)
	if err != nil {
		slog.Error("Failed to create withdraw proof request", "error", err.Error())
		return types.EnygmaProofResponse{}, fmt.Errorf("creating withdraw proof request: %w", err)
	}

	response, err := c.httpClient.Do(request)
	if err != nil {
		slog.Error("HTTP request failed for withdraw proof", "error", err.Error())
		return types.EnygmaProofResponse{}, fmt.Errorf("executing withdraw proof request: %w", err)
	}

	switch response.StatusCode {
	case http.StatusOK:
		return handleSuccessfulWithdrawProofRequest(response)
	case http.StatusInternalServerError:
		err := handleInternalServerError(response)
		return types.EnygmaProofResponse{}, fmt.Errorf("creating withdraw proof: %w", err)
	default:
		err := handleUnsupportedStatus(response)
		return types.EnygmaProofResponse{}, fmt.Errorf("creating withdraw proof: %w", err)
	}
}

type proof struct {
	PiA          [2]*big.Int
	PiB          [2][2]*big.Int
	PiC          [2]*big.Int
	PublicSignal []*big.Int
}

func parseProofResponse(r *http.Response) (proof, error) {
	defer func() { _ = r.Body.Close() }()

	var response ProofAPIResponse
	err := json.NewDecoder(r.Body).Decode(&response)
	if err != nil {
		return proof{}, withstack.Wrap(fmt.Errorf("decoding proof API response: %w", err))
	}

	PiA0, _ := new(big.Int).SetString(response.Pi_A[0], 10) // base 10, no [2:]
	PiA1, _ := new(big.Int).SetString(response.Pi_A[1], 10)
	PiB00, _ := new(big.Int).SetString(response.Pi_B[0][0], 10)
	PiB01, _ := new(big.Int).SetString(response.Pi_B[0][1], 10)
	PiB10, _ := new(big.Int).SetString(response.Pi_B[1][0], 10)
	PiB11, _ := new(big.Int).SetString(response.Pi_B[1][1], 10)
	PiC0, _ := new(big.Int).SetString(response.Pi_C[0], 10)
	PiC1, _ := new(big.Int).SetString(response.Pi_C[1], 10)
	publicSignals := make([]*big.Int, len(response.Public_Signal))
	for i, v := range response.Public_Signal {
		publicSignals[i], _ = new(big.Int).SetString(v, 10)
	}

	return proof{
		PiA:          [2]*big.Int{PiA0, PiA1},
		PiB:          [2][2]*big.Int{{PiB00, PiB01}, {PiB10, PiB11}},
		PiC:          [2]*big.Int{PiC0, PiC1},
		PublicSignal: publicSignals,
	}, nil
}

func proofResponseToEnygmaProof(r *http.Response) (types.EnygmaProofResponse, error) {
	proofResponse, err := parseProofResponse(r)
	if err != nil {
		return types.EnygmaProofResponse{}, fmt.Errorf("parsing enygma proof response: %w", err)
	}

	return types.EnygmaProofResponse{
		PiA:          proofResponse.PiA,
		PiB:          proofResponse.PiB,
		PiC:          proofResponse.PiC,
		PublicSignal: proofResponse.PublicSignal,
	}, nil
}

func proofResponseToDvpProof(r *http.Response) (dvp.Proof, error) {
	proofResponse, err := parseProofResponse(r)
	if err != nil {
		return dvp.Proof{}, fmt.Errorf("parsing dvp proof response: %w", err)
	}
	return dvp.Proof{
		A:            proofResponse.PiA,
		B:            proofResponse.PiB,
		C:            proofResponse.PiC,
		PublicSignal: proofResponse.PublicSignal,
	}, nil
}

func handleSuccessfulEnygmaJSRequest(r *http.Response) (dvp.Proof, error) {
	return proofResponseToDvpProof(r)
}

func handleSuccessfulErc721OwnershipRequest(r *http.Response) (dvp.Proof, error) {
	return proofResponseToDvpProof(r)
}

func handleSuccessfulErc1155JSRequest(r *http.Response) (dvp.Proof, error) {
	return proofResponseToDvpProof(r)
}

func handleSuccessfulTransferProofRequest(r *http.Response) (types.EnygmaProofResponse, error) {
	return proofResponseToEnygmaProof(r)
}

func handleSuccessfulDepositProofRequest(r *http.Response) (types.EnygmaProofResponse, error) {
	return proofResponseToEnygmaProof(r)
}

func handleSuccessfulWithdrawProofRequest(r *http.Response) (types.EnygmaProofResponse, error) {
	return proofResponseToEnygmaProof(r)
}
