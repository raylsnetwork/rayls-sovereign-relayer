package grpc

import (
	"context"
	"errors"
	"math/big"

	encryptpb "github.com/raylsnetwork/rayls-privacy-relayer-api/cts/gen/proto/encrypt"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/cts/service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type EncryptService interface {
	GCMEncryptWithNormalMessageTag(
		ctx context.Context,
		blob []byte,
		chainID uint64,
		blockNumber uint64,
		prevBlockHash string) ([]byte, string, error)
	GCMDecryptWithNormalMessageTag(
		ctx context.Context,
		encryptedData []byte,
		fingerprint string,
		blockNumber uint64,
		prevBlockHash string) ([]byte, error)
	GCMEncrypt(
		ctx context.Context,
		blob []byte,
		chainID uint64,
		blockNumber uint64) ([]byte, error)
	GCMEncryptEnygma(
		ctx context.Context,
		blob []byte,
		chainID uint64,
		blockNumber uint64,
		resourceID []byte) ([]byte, error)
	GCMDecrypt(
		ctx context.Context,
		encryptedData []byte,
		chainID uint64,
		blockNumber uint64) ([]byte, error)
	GCMEncryptWithProvidedSS(blob []byte, ss []byte) ([]byte, error)
	GCMDecryptWithProvidedSS(encryptedData []byte, ss []byte) ([]byte, error)
	GCMDecryptWithEnygmaMessageTag(
		ctx context.Context,
		encryptedData []byte,
		messageTag []byte,
		blockNumber uint64,
		anonymitySet []*big.Int,
		chainID *big.Int,
		resourceID []byte) ([]byte, error)
	EncryptPrivateKey(
		ctx context.Context,
		chainID uint64,
		blockNumber uint64) (string, []byte, []byte, uint64, error)
}

type EncryptHandler struct {
	encryptpb.UnimplementedEncryptServiceServer

	svc EncryptService
}

func NewEncryptHandler(svc EncryptService) *EncryptHandler {
	return &EncryptHandler{svc: svc}
}

func (h *EncryptHandler) Encrypt(
	ctx context.Context,
	req *encryptpb.EncryptRequest,
) (*encryptpb.EncryptResponse, error) {
	if len(req.GetPlaintext()) == 0 {
		return nil, status.Error(codes.InvalidArgument, "plaintext is required")
	}
	if req.GetChainID() == 0 {
		return nil, status.Error(codes.InvalidArgument, "chain_id is required")
	}

	encryptedData, fp, err := h.svc.GCMEncryptWithNormalMessageTag(
		ctx,
		req.GetPlaintext(),
		req.GetChainID(),
		req.GetBlockNumber(),
		req.GetPrevBlockHash(),
	)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "encrypt: %v", err)
	}

	return &encryptpb.EncryptResponse{
		EncryptedData: encryptedData,
		Fingerprint:   fp,
	}, nil
}

func (h *EncryptHandler) Decrypt(
	ctx context.Context,
	req *encryptpb.DecryptRequest,
) (*encryptpb.DecryptResponse, error) {
	if len(req.GetEncryptedData()) == 0 {
		return nil, status.Error(codes.InvalidArgument, "encryptedData is required")
	}
	if req.GetFingerprint() == "" {
		return nil, status.Error(codes.InvalidArgument, "fingerprint is required")
	}

	plaintext, err := h.svc.GCMDecryptWithNormalMessageTag(
		ctx,
		req.GetEncryptedData(),
		req.GetFingerprint(),
		req.GetBlockNumber(),
		req.GetPrevBlockHash(),
	)
	if err != nil {
		if errors.Is(err, service.ErrNotForRecipient) {
			return &encryptpb.DecryptResponse{Outcome: encryptpb.DecryptOutcome_OUTCOME_NOT_FOR_RECIPIENT}, nil
		}
		if errors.Is(err, service.ErrAuthFailed) {
			return &encryptpb.DecryptResponse{Outcome: encryptpb.DecryptOutcome_OUTCOME_TAMPERED}, nil
		}
		return nil, status.Errorf(codes.Internal, "decrypt: %v", err)
	}

	return &encryptpb.DecryptResponse{
		Plaintext: plaintext,
		Outcome:   encryptpb.DecryptOutcome_OUTCOME_OK,
	}, nil
}

func (h *EncryptHandler) EncryptWithoutFP(
	ctx context.Context,
	req *encryptpb.EncryptWithoutFPRequest,
) (*encryptpb.EncryptWithoutFPResponse, error) {
	if len(req.GetPlaintext()) == 0 {
		return nil, status.Error(codes.InvalidArgument, "plaintext is required")
	}
	if req.GetChainID() == 0 {
		return nil, status.Error(codes.InvalidArgument, "chain_id is required")
	}

	var encryptedData []byte
	var err error
	if len(req.GetResourceID()) > 0 {
		// Enygma transfer-batch seal: use the self-aware secret selection so the sender's own-chain
		// (self/change) batch is sealed with the self-secret (matching the message tag + decrypt),
		// not the static own-chain pairwise row.
		encryptedData, err = h.svc.GCMEncryptEnygma(
			ctx,
			req.GetPlaintext(),
			req.GetChainID(),
			req.GetBlockNumber(),
			req.GetResourceID(),
		)
	} else {
		encryptedData, err = h.svc.GCMEncrypt(
			ctx,
			req.GetPlaintext(),
			req.GetChainID(),
			req.GetBlockNumber(),
		)
	}
	if err != nil {
		return nil, status.Errorf(codes.Internal, "encrypt: %v", err)
	}

	return &encryptpb.EncryptWithoutFPResponse{
		EncryptedData: encryptedData,
	}, nil
}

func (h *EncryptHandler) DecryptWithoutFP(
	ctx context.Context,
	req *encryptpb.DecryptWithoutFPRequest,
) (*encryptpb.DecryptWithoutFPResponse, error) {
	if len(req.GetEncryptedData()) == 0 {
		return nil, status.Error(codes.InvalidArgument, "encryptedData is required")
	}
	if req.GetChainID() == 0 {
		return nil, status.Error(codes.InvalidArgument, "chain_id is required")
	}

	plaintext, err := h.svc.GCMDecrypt(
		ctx,
		req.GetEncryptedData(),
		req.GetChainID(),
		req.GetBlockNumber(),
	)
	if err != nil {
		// No "not for me" outcome on this path — the caller explicitly addressed
		// a chain, so any AEAD failure is suspicious.
		if errors.Is(err, service.ErrAuthFailed) {
			return &encryptpb.DecryptWithoutFPResponse{Outcome: encryptpb.DecryptOutcome_OUTCOME_TAMPERED}, nil
		}
		return nil, status.Errorf(codes.Internal, "decrypt: %v", err)
	}

	return &encryptpb.DecryptWithoutFPResponse{
		Plaintext: plaintext,
		Outcome:   encryptpb.DecryptOutcome_OUTCOME_OK,
	}, nil
}

func (h *EncryptHandler) EncryptWithoutFPWithSS(
	ctx context.Context,
	req *encryptpb.EncryptWithoutFPWithSSRequest,
) (*encryptpb.EncryptWithoutFPWithSSResponse, error) {
	if len(req.GetPlaintext()) == 0 {
		return nil, status.Error(codes.InvalidArgument, "plaintext is required")
	}
	if len(req.GetSs()) == 0 {
		return nil, status.Error(codes.InvalidArgument, "ss is required")
	}
	encryptedData, err := h.svc.GCMEncryptWithProvidedSS(req.GetPlaintext(), req.GetSs())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "encrypt: %v", err)
	}
	return &encryptpb.EncryptWithoutFPWithSSResponse{EncryptedData: encryptedData}, nil
}

func (h *EncryptHandler) DecryptWithoutFPWithSS(
	ctx context.Context,
	req *encryptpb.DecryptWithoutFPWithSSRequest,
) (*encryptpb.DecryptWithoutFPWithSSResponse, error) {
	if len(req.GetEncryptedData()) == 0 {
		return nil, status.Error(codes.InvalidArgument, "encryptedData is required")
	}
	if len(req.GetSs()) == 0 {
		return nil, status.Error(codes.InvalidArgument, "ss is required")
	}
	plaintext, err := h.svc.GCMDecryptWithProvidedSS(req.GetEncryptedData(), req.GetSs())
	if err != nil {
		// Salt-based path: AEAD-fail is the only "not for me" signal we have.
		if errors.Is(err, service.ErrNotForRecipient) {
			return &encryptpb.DecryptWithoutFPWithSSResponse{Outcome: encryptpb.DecryptOutcome_OUTCOME_NOT_FOR_RECIPIENT}, nil
		}
		return nil, status.Errorf(codes.Internal, "decrypt: %v", err)
	}
	return &encryptpb.DecryptWithoutFPWithSSResponse{
		Plaintext: plaintext,
		Outcome:   encryptpb.DecryptOutcome_OUTCOME_OK,
	}, nil
}

func (h *EncryptHandler) DecryptEnygmaTransferBatch(
	ctx context.Context,
	req *encryptpb.DecryptEnygmaTransferBatchRequest,
) (*encryptpb.DecryptEnygmaTransferBatchResponse, error) {
	if len(req.GetEncryptedData()) == 0 {
		return nil, status.Error(codes.InvalidArgument, "encrypted_data is required")
	}
	if len(req.GetMessageTag()) == 0 {
		return nil, status.Error(codes.InvalidArgument, "message_tag is required")
	}
	if len(req.GetAnonymitySet()) == 0 {
		return nil, status.Error(codes.InvalidArgument, "anonymity_set is required")
	}
	if req.GetChainID() == 0 {
		return nil, status.Error(codes.InvalidArgument, "chain_id is required")
	}
	if len(req.GetResourceID()) == 0 {
		return nil, status.Error(codes.InvalidArgument, "resource_id is required")
	}

	anonymitySet := make([]*big.Int, len(req.GetAnonymitySet()))
	for i, id := range req.GetAnonymitySet() {
		anonymitySet[i] = new(big.Int).SetUint64(id)
	}

	plaintext, err := h.svc.GCMDecryptWithEnygmaMessageTag(
		ctx,
		req.GetEncryptedData(),
		req.GetMessageTag(),
		req.GetBlockNumber(),
		anonymitySet,
		new(big.Int).SetUint64(req.GetChainID()),
		req.GetResourceID(),
	)
	if err != nil {
		if errors.Is(err, service.ErrNotForRecipient) {
			return &encryptpb.DecryptEnygmaTransferBatchResponse{Outcome: encryptpb.DecryptOutcome_OUTCOME_NOT_FOR_RECIPIENT}, nil
		}
		if errors.Is(err, service.ErrAuthFailed) {
			return &encryptpb.DecryptEnygmaTransferBatchResponse{Outcome: encryptpb.DecryptOutcome_OUTCOME_TAMPERED}, nil
		}
		return nil, status.Errorf(codes.Internal, "decrypt enygma transfer batch: %v", err)
	}

	return &encryptpb.DecryptEnygmaTransferBatchResponse{
		Plaintext: plaintext,
		Outcome:   encryptpb.DecryptOutcome_OUTCOME_OK,
	}, nil
}

func (h *EncryptHandler) EncryptPrivateKey(
	ctx context.Context,
	req *encryptpb.EncryptPrivateKeyRequest,
) (*encryptpb.EncryptPrivateKeyResponse, error) {
	if req.GetChainID() == 0 {
		return nil, status.Error(codes.InvalidArgument, "chain_id is required")
	}

	publicKeyHex, encryptedPrivateKey, mac, initialBlock, err := h.svc.EncryptPrivateKey(
		ctx,
		req.GetChainID(),
		req.GetBlockNumber(),
	)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "encrypt private key: %v", err)
	}

	return &encryptpb.EncryptPrivateKeyResponse{
		PublicKeyHex:        publicKeyHex,
		EncryptedPrivateKey: encryptedPrivateKey,
		Mac:                 mac,
		InitialBlock:        initialBlock,
	}, nil
}
