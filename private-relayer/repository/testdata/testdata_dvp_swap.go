package testdata

import (
	"math/big"
	"time"

	"github.com/raylsnetwork/rayls-sovereign-relayer/private-relayer/repository"
	"github.com/raylsnetwork/rayls-sovereign-relayer/types"
)

var (
	swapCreatedAt = time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	swapExpiredAt = time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC)
	swapFutureExp = time.Date(2099, 12, 31, 0, 0, 0, 0, time.UTC)
)

var (
	DvpSwap1 = types.DvpSwap{
		SharedID:               "swap-shared-id-1",
		From:                   "0xFromAddr1",
		To:                     "0xToAddr1",
		SourcePaymentPublicKey: big.NewInt(111),
		SourceChainID:          big.NewInt(1),
		DestPaymentPublicKey:   big.NewInt(222),
		DestChainID:            big.NewInt(2),
		TokenInAmount:          big.NewInt(1000),
		TokenInAddress:         "0xTokenIn1",
		TokenInResourceID:      "res-in-1",
		TokenInType:            types.DvpEnygma,
		TokenInID:              "",
		TokenOutAmount:         big.NewInt(2000),
		TokenOutAddress:        "0xTokenOut1",
		TokenOutResourceID:     "res-out-1",
		TokenOutType:           types.DvpERC721,
		TokenOutID:             "nft-1",
		Status:                 types.DvpSwapWaitingConfirmation,
		CreatedAt:              swapCreatedAt,
		ExpiresAt:              swapFutureExp,
		CancelledAt:            nil,
	}
	DvpSwapExpired = types.DvpSwap{
		SharedID:               "swap-shared-id-expired",
		From:                   "0xFromAddr2",
		To:                     "0xToAddr2",
		SourcePaymentPublicKey: nil,
		SourceChainID:          big.NewInt(1),
		DestPaymentPublicKey:   nil,
		DestChainID:            big.NewInt(2),
		TokenInAmount:          big.NewInt(500),
		TokenInAddress:         "0xTokenIn2",
		TokenInResourceID:      "res-in-2",
		TokenInType:            types.DvpEnygma,
		TokenInID:              "",
		TokenOutAmount:         big.NewInt(1000),
		TokenOutAddress:        "0xTokenOut2",
		TokenOutResourceID:     "res-out-2",
		TokenOutType:           types.DvpERC721,
		TokenOutID:             "nft-2",
		Status:                 types.DvpSwapWaitingConfirmation,
		CreatedAt:              swapCreatedAt,
		ExpiresAt:              swapExpiredAt,
		CancelledAt:            nil,
	}
)

var (
	ModelDvpSwap1 = repository.DvpSwap{
		SharedID:               DvpSwap1.SharedID,
		From:                   DvpSwap1.From,
		To:                     DvpSwap1.To,
		SourcePaymentPublicKey: DvpSwap1.SourcePaymentPublicKey.String(),
		SourceChainID:          DvpSwap1.SourceChainID.String(),
		DestPaymentPublicKey:   DvpSwap1.DestPaymentPublicKey.String(),
		DestChainID:            DvpSwap1.DestChainID.String(),
		TokenInAmount:          DvpSwap1.TokenInAmount.String(),
		TokenInAddress:         DvpSwap1.TokenInAddress,
		TokenInResourceID:      DvpSwap1.TokenInResourceID,
		TokenInType:            DvpSwap1.TokenInType,
		TokenInID:              DvpSwap1.TokenInID,
		TokenOutAmount:         DvpSwap1.TokenOutAmount.String(),
		TokenOutAddress:        DvpSwap1.TokenOutAddress,
		TokenOutResourceID:     DvpSwap1.TokenOutResourceID,
		TokenOutType:           DvpSwap1.TokenOutType,
		TokenOutID:             DvpSwap1.TokenOutID,
		Status:                 DvpSwap1.Status,
		CreatedAt:              DvpSwap1.CreatedAt,
		ExpiresAt:              DvpSwap1.ExpiresAt,
		CancelledAt:            nil,
	}
	ModelDvpSwapExpired = repository.DvpSwap{
		SharedID:               DvpSwapExpired.SharedID,
		From:                   DvpSwapExpired.From,
		To:                     DvpSwapExpired.To,
		SourcePaymentPublicKey: "",
		SourceChainID:          DvpSwapExpired.SourceChainID.String(),
		DestPaymentPublicKey:   "",
		DestChainID:            DvpSwapExpired.DestChainID.String(),
		TokenInAmount:          DvpSwapExpired.TokenInAmount.String(),
		TokenInAddress:         DvpSwapExpired.TokenInAddress,
		TokenInResourceID:      DvpSwapExpired.TokenInResourceID,
		TokenInType:            DvpSwapExpired.TokenInType,
		TokenInID:              DvpSwapExpired.TokenInID,
		TokenOutAmount:         DvpSwapExpired.TokenOutAmount.String(),
		TokenOutAddress:        DvpSwapExpired.TokenOutAddress,
		TokenOutResourceID:     DvpSwapExpired.TokenOutResourceID,
		TokenOutType:           DvpSwapExpired.TokenOutType,
		TokenOutID:             DvpSwapExpired.TokenOutID,
		Status:                 DvpSwapExpired.Status,
		CreatedAt:              DvpSwapExpired.CreatedAt,
		ExpiresAt:              DvpSwapExpired.ExpiresAt,
		CancelledAt:            nil,
	}
)
