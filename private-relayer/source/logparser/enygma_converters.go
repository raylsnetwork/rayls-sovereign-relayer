package logparser

import (
	"encoding/hex"
	"fmt"
	"log/slog"
	"math/big"

	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/google/uuid"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/contracts/EnygmaPNEvents"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/private-relayer/source/service"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/types"
)

// convertEnygmaProgramData translates a contract-binding ProgramData step array into the
// domain representation, keeping the service/types layers free of abigen-generated structs.
func convertEnygmaProgramData(steps []EnygmaPNEvents.SharedObjectsEnygmaProgramData) []types.EnygmaProgramData {
	if steps == nil {
		return nil
	}
	domain := make([]types.EnygmaProgramData, len(steps))
	for i, s := range steps {
		domain[i] = types.EnygmaProgramData{
			ResourceId:      s.ResourceId,
			ContractAddress: s.ContractAddress,
			Selector:        s.Selector,
			Args:            s.Args,
		}
	}
	return domain
}

// Enygma Event Converters

// convertEnygmaCreation converts a contract EnygmaCreation event to service type
func convertEnygmaCreation(event *EnygmaPNEvents.EnygmaPNEventsEnygmaCreation) *service.EnygmaCreation {
	return &service.EnygmaCreation{
		InitialSupply: event.InitialSupply,
	}
}

// convertEnygmaSendTransferCC converts a contract EnygmaSendTransferCC event to service type
func convertEnygmaSendTransferCC(event *EnygmaPNEvents.EnygmaPNEventsEnygmaSendTransferPNH) []*service.EnygmaTransferTx {
	// Ensure that the number of chain IDs, addresses, values and program-data arrays match.
	// ProgramData is parallel to the other arrays: one array of blobs per recipient.
	if len(event.ToChainId) != len(event.To) ||
		len(event.ToChainId) != len(event.Value) ||
		len(event.ToChainId) != len(event.ProgramData) {
		slog.Warn("discarding malformed EnygmaSendTransferPNH event: parallel array length mismatch",
			slog.String("reference_id", hex.EncodeToString(event.ReferenceId[:])),
			slog.Int("to_chain_id", len(event.ToChainId)),
			slog.Int("to", len(event.To)),
			slog.Int("value", len(event.Value)),
			slog.Int("program_data", len(event.ProgramData)))
		return nil
	}

	var transfers []*service.EnygmaTransferTx

	// The event has multiple destinations. Each destination is a separate transaction.
	for i := range len(event.ToChainId) {
		transfers = append(transfers, &service.EnygmaTransferTx{
			DestIdx:     i,
			MessageId:   uuid.New().String(),
			ReferenceId: event.ReferenceId,
			FromAddress: event.From,
			ToChainId:   event.ToChainId[i],
			ToAmount:    event.Value[i],
			ToAddress:   event.To[i],
			ProgramData: convertEnygmaProgramData(event.ProgramData[i]),
		})
	}

	return transfers
}

// convertEnygmaMint converts a contract EnygmaMint event to service type
func convertEnygmaMint(event *EnygmaPNEvents.EnygmaPNEventsEnygmaMint) *service.EnygmaMint {
	return &service.EnygmaMint{
		To:     event.To,
		Amount: event.Amount,
		TxHash: event.Raw.TxHash,
	}
}

// convertEnygmaBurn converts a contract EnygmaBurn event to service type
func convertEnygmaBurn(event *EnygmaPNEvents.EnygmaPNEventsEnygmaBurn) *service.EnygmaBurn {
	return &service.EnygmaBurn{
		From:   event.From,
		Amount: event.Amount,
		TxHash: event.Raw.TxHash,
	}
}

// convertEnygmaDepositToDvp converts a contract EnygmaDepositToDvp event to service type
func convertEnygmaDepositToDvp(event *EnygmaPNEvents.EnygmaPNEventsEnygmaDepositToDvp) *service.EnygmaDepositToDvp {
	return &service.EnygmaDepositToDvp{
		Amount:        event.Amount,
		From:          event.From,
		ReferenceId:   event.ReferenceId,
		TxHash:        event.Raw.TxHash,
		TxBlockNumber: new(big.Int).SetUint64(event.Raw.BlockNumber),
	}
}

// convertEnygmaWithdrawFromDvp converts a contract EnygmaWithdrawFromDvp event to service type
func convertEnygmaWithdrawFromDvp(event *EnygmaPNEvents.EnygmaPNEventsEnygmaWithdrawFromDvp) *service.EnygmaWithdrawFromDvp {
	return &service.EnygmaWithdrawFromDvp{
		Amount:        event.Amount,
		To:            event.To,
		ReferenceId:   event.ReferenceId,
		TxHash:        event.Raw.TxHash,
		TxBlockNumber: new(big.Int).SetUint64(event.Raw.BlockNumber),
	}
}

// convertEnygmaSwapWithDvpForERC721 converts a contract EnygmaSwapWithDvpForERC721 event to service type
func convertEnygmaSwapWithDvpForERC721(
	event *EnygmaPNEvents.EnygmaPNEventsEnygmaSwapWithDvpForERC721,
) *service.DvpEnygmaSwapERC721 {
	return &service.DvpEnygmaSwapERC721{
		SharedId:      hex.EncodeToString(event.SharedId[:]),
		DestChainId:   event.DestChainId,
		From:          event.From,
		ResourceId:    hex.EncodeToString(event.ResourceId[:]),
		EnygmaAmount:  event.EnygmaAmount,
		NftResourceId: hex.EncodeToString(event.NftResourceId[:]),
		NftId:         event.NftId.String(),
		TxHash:        event.Raw.TxHash.Hex(),
		TxBlockNumber: new(big.Int).SetUint64(event.Raw.BlockNumber),
		ValidityTime:  event.ValidityTime,
	}
}

// convertEnygmaSwapWithDvpForERC1155 converts a contract EnygmaSwapWithDvpForERC1155 event to service type
func convertEnygmaSwapWithDvpForERC1155(
	event *EnygmaPNEvents.EnygmaPNEventsEnygmaSwapWithDvpForERC1155,
) *service.DvpEnygmaSwapERC1155 {
	return &service.DvpEnygmaSwapERC1155{
		SharedId:       hex.EncodeToString(event.SharedId[:]),
		DestChainId:    event.DestChainId,
		From:           event.From,
		ResourceId:     hex.EncodeToString(event.ResourceId[:]),
		EnygmaAmount:   event.EnygmaAmount,
		NftResourceId:  hex.EncodeToString(event.NftResourceId[:]),
		NftId:          event.NftId.String(),
		NftAmountOrOne: event.NftAmountOrOne,
		TxHash:         event.Raw.TxHash.Hex(),
		TxBlockNumber:  new(big.Int).SetUint64(event.Raw.BlockNumber),
		ValidityTime:   event.ValidityTime,
	}
}

// Dvp ERC721 Event Converters

// convertDvp721Creation converts a contract Dvp721Creation event to service type
func convertDvp721Creation(event *EnygmaPNEvents.EnygmaPNEventsDvp721Creation) *service.Dvp721Creation {
	return &service.Dvp721Creation{
		ChainEventID: getChainEventID(event.Raw),
		ResourceId:   hex.EncodeToString(event.ResourceId[:]),
	}
}

// convertDvp721Mint converts a contract Dvp721Mint event to service type
func convertDvp721Mint(event *EnygmaPNEvents.EnygmaPNEventsDvp721Mint) *service.Dvp721Mint {
	return &service.Dvp721Mint{
		ChainEventID: getChainEventID(event.Raw),
		ResourceId:   hex.EncodeToString(event.ResourceId[:]),
		NftId:        event.NftId,
	}
}

// convertDvp721Burn converts a contract Dvp721Burn event to service type
func convertDvp721Burn(event *EnygmaPNEvents.EnygmaPNEventsDvp721Burn) *service.Dvp721Burn {
	return &service.Dvp721Burn{
		ChainEventID: getChainEventID(event.Raw),
		ResourceId:   hex.EncodeToString(event.ResourceId[:]),
		NftId:        event.NftId,
	}
}

// convertDvp721DepositIntoDvp converts a contract Dvp721DepositIntoDvp event to service type
func convertDvp721DepositIntoDvp(
	event *EnygmaPNEvents.EnygmaPNEventsDvp721DepositIntoDvp,
) *service.Dvp721DepositIntoDvp {
	return &service.Dvp721DepositIntoDvp{
		ChainEventID:  getChainEventID(event.Raw),
		ResourceId:    hex.EncodeToString(event.ResourceId[:]),
		NftId:         event.NftId,
		From:          event.From,
		TxHash:        event.Raw.TxHash.Hex(),
		TxBlockNumber: new(big.Int).SetUint64(event.Raw.BlockNumber),
	}
}

// convertDvp721WithdrawFromDvp converts a contract Dvp721WithdrawFromDvp event to service type
func convertDvp721WithdrawFromDvp(
	event *EnygmaPNEvents.EnygmaPNEventsDvp721WithdrawFromDvp,
) *service.Dvp721WithdrawFromDvp {
	return &service.Dvp721WithdrawFromDvp{
		ChainEventID:  getChainEventID(event.Raw),
		ResourceId:    hex.EncodeToString(event.ResourceId[:]),
		NftId:         event.NftId,
		Owner:         event.Owner,
		TxHash:        event.Raw.TxHash.Hex(),
		TxBlockNumber: new(big.Int).SetUint64(event.Raw.BlockNumber),
	}
}

// convertDvp721SwapForEnygma converts a contract Dvp721SwapForEnygma event to service type
func convertDvp721SwapForEnygma(event *EnygmaPNEvents.EnygmaPNEventsDvp721SwapForEnygma) *service.Dvp721SwapForEnygma {
	return &service.Dvp721SwapForEnygma{
		SharedId:         hex.EncodeToString(event.SharedId[:]),
		DestChainId:      event.DestChainId,
		From:             event.From,
		NftResourceId:    hex.EncodeToString(event.NftResourceId[:]),
		NftId:            event.NftId.String(),
		EnygmaResourceId: hex.EncodeToString(event.EnygmaResourceId[:]),
		EnygmaAmount:     event.EnygmaAmount,
		TxHash:           event.Raw.TxHash.Hex(),
		TxBlockNumber:    new(big.Int).SetUint64(event.Raw.BlockNumber),
		ValidityTime:     event.ValidityTime,
	}
}

// Dvp ERC1155 Event Converters

// convertDvp1155Creation converts a contract Dvp1155Creation event to service type
func convertDvp1155Creation(event *EnygmaPNEvents.EnygmaPNEventsDvp1155Creation) *service.Dvp1155Creation {
	return &service.Dvp1155Creation{
		ChainEventID: getChainEventID(event.Raw),
		ResourceId:   hex.EncodeToString(event.ResourceId[:]),
	}
}

// convertDvp1155Mint converts a contract Dvp1155Mint event to service type
func convertDvp1155Mint(event *EnygmaPNEvents.EnygmaPNEventsDvp1155Mint) *service.Dvp1155Mint {
	return &service.Dvp1155Mint{
		ChainEventID: getChainEventID(event.Raw),
		ResourceId:   hex.EncodeToString(event.ResourceId[:]),
		TokenId:      event.TokenId,
		Value:        event.Value,
		Data:         event.Data,
	}
}

// convertDvp1155Burn converts a contract Dvp1155Burn event to service type
// Note: The 'To' field is omitted as it's not used in business logic
func convertDvp1155Burn(event *EnygmaPNEvents.EnygmaPNEventsDvp1155Burn) *service.Dvp1155Burn {
	return &service.Dvp1155Burn{
		ChainEventID: getChainEventID(event.Raw),
		ResourceId:   hex.EncodeToString(event.ResourceId[:]),
		TokenId:      event.TokenId,
		Value:        event.Value,
	}
}

// convertDvp1155DepositIntoDvp converts a contract Dvp1155DepositIntoDvp event to service type
func convertDvp1155DepositIntoDvp(
	event *EnygmaPNEvents.EnygmaPNEventsDvp1155DepositIntoDvp,
) *service.Dvp1155DepositIntoDvp {
	return &service.Dvp1155DepositIntoDvp{
		ChainEventID:  getChainEventID(event.Raw),
		ResourceId:    hex.EncodeToString(event.ResourceId[:]),
		TokenId:       event.TokenId,
		Value:         event.Value,
		Data:          event.Data,
		From:          event.From,
		TxHash:        event.Raw.TxHash.Hex(),
		TxBlockNumber: new(big.Int).SetUint64(event.Raw.BlockNumber),
	}
}

// convertDvp1155WithdrawFromDvp converts a contract Dvp1155WithdrawFromDvp event to service type
// Note: The 'Data' field is omitted as it's not used in business logic
func convertDvp1155WithdrawFromDvp(
	event *EnygmaPNEvents.EnygmaPNEventsDvp1155WithdrawFromDvp,
) *service.Dvp1155WithdrawFromDvp {
	return &service.Dvp1155WithdrawFromDvp{
		ChainEventID:  getChainEventID(event.Raw),
		ResourceId:    hex.EncodeToString(event.ResourceId[:]),
		TokenId:       event.TokenId,
		Value:         event.Value,
		Owner:         event.Owner,
		TxHash:        event.Raw.TxHash.Hex(),
		TxBlockNumber: new(big.Int).SetUint64(event.Raw.BlockNumber),
	}
}

// convertDvp1155SwapForEnygma converts a contract Dvp1155SwapForEnygma event to service type
// Note: The 'TokenData' field is omitted as it's not used in business logic
func convertDvp1155SwapForEnygma(
	event *EnygmaPNEvents.EnygmaPNEventsDvp1155SwapForEnygma,
) *service.Dvp1155SwapForEnygma {
	return &service.Dvp1155SwapForEnygma{
		SharedId:         hex.EncodeToString(event.SharedId[:]),
		DestChainId:      event.DestChainId,
		From:             event.From,
		TokenResourceId:  hex.EncodeToString(event.TokenResourceId[:]),
		TokenValue:       event.TokenValue,
		TokenId:          event.TokenId.String(),
		EnygmaResourceId: hex.EncodeToString(event.EnygmaResourceId[:]),
		EnygmaAmount:     event.EnygmaAmount,
		TxHash:           event.Raw.TxHash.Hex(),
		TxBlockNumber:    new(big.Int).SetUint64(event.Raw.BlockNumber),
		ValidityTime:     event.ValidityTime,
	}
}

func convertDvpSwapCancelled(event *EnygmaPNEvents.EnygmaPNEventsDvpSwapCancelled) *service.DvpSwapCancelled {
	return &service.DvpSwapCancelled{
		SharedId:           hex.EncodeToString(event.SharedId[:]),
		DestChainId:        event.DestChainId,
		TokenInResourceId:  hex.EncodeToString(event.TokenInResourceId[:]),
		TokenInAmount:      event.TokenInAmount,
		TokenInId:          event.TokenInId,
		TokenInStandard:    event.TokenInStandard,
		TokenOutResourceId: hex.EncodeToString(event.TokenOutResourceId[:]),
		TokenOutAmount:     event.TokenOutAmount,
		TokenOutId:         event.TokenOutId,
		TokenOutStandard:   event.TokenOutStandard,
	}
}

// Generate a unique ID by concatinating the block number,
// tx index and log index. Also known as log coordinates.
func getChainEventID(log *ethtypes.Log) string {
	return fmt.Sprintf("%d-%d-%d", log.BlockNumber, log.TxIndex, log.Index)
}
