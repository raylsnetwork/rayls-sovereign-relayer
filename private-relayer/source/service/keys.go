package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"math/big"

	"github.com/raylsnetwork/rayls-privacy-relayer-api/contractclient"
	ps "github.com/raylsnetwork/rayls-privacy-relayer-api/contracts/ParticipantStorageV1"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/cts/gen/proto/encrypt"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/cts/gen/proto/keys"
	"google.golang.org/grpc"
)

type KeysKOSClient interface {
	CreateViewKeyPair(ctx context.Context, in *keys.CreateViewKeyPairRequest, opts ...grpc.CallOption) (*keys.CreateViewKeyPairResponse, error)
	EncryptPrivateKey(
		ctx context.Context,
		in *encrypt.EncryptPrivateKeyRequest,
		opts ...grpc.CallOption,
	) (*encrypt.EncryptPrivateKeyResponse, error)
	GenerateKeyAgreement(
		ctx context.Context,
		in *keys.GenerateKeyAgreementRequest,
		opts ...grpc.CallOption,
	) (*keys.GenerateKeyAgreementResponse, error)
	CreateKeyAgreement(
		ctx context.Context,
		in *keys.CreateKeyAgreementRequest,
		opts ...grpc.CallOption,
	) (*keys.CreateKeyAgreementResponse, error)
}

type KeysParticipantStorageClient interface {
	GetVenOperatorChainInfo(ctx context.Context, blockNumber *big.Int) (ps.ParticipantStructsPrivacyNodeViewData, error)
	GetAllParticipantsViewData(ctx context.Context, blockNumber *big.Int) ([]ps.ParticipantStructsPrivacyNodeViewData, error)
	SetAuditInfo(ctx context.Context, chainID *big.Int, blockNumber *big.Int, publicKey string, encrPrivateKey, mac []byte) error
	SetChainViewData(ctx context.Context, chainID *big.Int, raylsViewPublicKey string, blockNumber *big.Int) error
	InitiateKeyAgreement(ctx context.Context, chainID *big.Int, ciphertext []byte, digest []byte, blockNumber *big.Int) error
}

type KeysService struct {
	myChainID *big.Int

	kosClient KeysKOSClient
	psClient  KeysParticipantStorageClient
}

func NewKeysService(myChainID *big.Int, kosClient KeysKOSClient, psClient KeysParticipantStorageClient) *KeysService {
	return &KeysService{
		myChainID: myChainID,

		kosClient: kosClient,
		psClient:  psClient,
	}
}

func (k *KeysService) UpdateRaylsViewKeys(ctx context.Context, blockNumber *big.Int) error {
	createResp, err := k.kosClient.CreateViewKeyPair(
		ctx,
		&keys.CreateViewKeyPairRequest{InitialBlock: blockNumber.Uint64()},
	)
	if err != nil {
		return fmt.Errorf("failed to create View keys for block number %d: %w", blockNumber.Int64(), err)
	}

	err = k.psClient.SetChainViewData(ctx, k.myChainID, createResp.GetPublicKey(), blockNumber)
	if err != nil {
		return fmt.Errorf("failed to set chain info: %w", err)
	}

	err = k.initiateNewKeyAgreements(ctx, blockNumber)
	if err != nil {
		return fmt.Errorf("failed to initiate new key agreements for block number %d: %w", blockNumber.Int64(), err)
	}

	venChainInfo, err := k.psClient.GetVenOperatorChainInfo(ctx, blockNumber)
	if err != nil {
		return fmt.Errorf("failed to get VEN chain info for block number %d: %w", blockNumber.Int64(), err)
	}

	resp, err := k.kosClient.EncryptPrivateKey(
		ctx,
		&encrypt.EncryptPrivateKeyRequest{
			ChainID:     venChainInfo.ChainId.Uint64(),
			BlockNumber: blockNumber.Uint64(),
		},
	)
	if err != nil {
		return fmt.Errorf("failed to encrypt private key for ven operator: %w", err)
	}

	err = k.psClient.SetAuditInfo(
		ctx,
		k.myChainID,
		new(big.Int).SetUint64(resp.InitialBlock),
		resp.PublicKeyHex,
		resp.EncryptedPrivateKey,
		resp.Mac,
	)
	if err != nil {
		return fmt.Errorf("failed to set audit info: %w", err)
	}

	return nil
}

func (k *KeysService) initiateNewKeyAgreements(ctx context.Context, blockNumber *big.Int) error {
	participants, err := k.psClient.GetAllParticipantsViewData(ctx, blockNumber)
	if err != nil {
		return fmt.Errorf("failed to get all chain view datas: %w", err)
	}

	for _, participant := range participants {
		slog.Info("Creating key agreement for participant", slog.String("chain_id", participant.ChainId.String()))
		genResp, err := k.kosClient.GenerateKeyAgreement(
			ctx,
			&keys.GenerateKeyAgreementRequest{
				ChainID:   participant.ChainId.Uint64(),
				PublicKey: []byte(participant.RaylsViewPublicKey),
			},
		)
		if err != nil {
			return fmt.Errorf("failed to generate key agreement: %w", err)
		}
		err = k.psClient.InitiateKeyAgreement(ctx, participant.ChainId, genResp.Ciphertext, genResp.Digest, blockNumber)
		switch {
		case err == nil:
			_, err = k.kosClient.CreateKeyAgreement(
				ctx,
				&keys.CreateKeyAgreementRequest{
					ChainID:      participant.ChainId.Uint64(),
					SharedSecret: genResp.SharedSecret,
					BlockNumber:  blockNumber.Uint64(),
				},
			)
			if err != nil {
				return fmt.Errorf("failed to create key agreement: %w", err)
			}
		case errors.Is(err, contractclient.ErrOutdatedKeyAgreement):
			slog.Warn(
				"Key agreement already exists for newer block number, skipping",
				slog.String("chain_id", participant.ChainId.String()),
			)
		default:
			return fmt.Errorf("failed to initiate key agreement: %w", err)
		}
	}

	return nil
}
