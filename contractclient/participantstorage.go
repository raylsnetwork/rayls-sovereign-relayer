package contractclient

import (
	"context"
	"math/big"
	"slices"

	"github.com/ethereum/go-ethereum/common"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/contracts/AuditManagerV1"
	ps "github.com/raylsnetwork/rayls-privacy-relayer-api/contracts/ParticipantStorageV1"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/withstack"
)

type ParticipantStorageClient struct {
	address         common.Address
	contract        *ps.ParticipantStorageV1
	auditManagerCtr *AuditManagerV1.AuditManagerV1

	executor Executor

	plChainID  *big.Int
	venChainID *big.Int
}

func NewParticipantStorageClient(
	address common.Address,
	plChainID, venChainID *big.Int,
	executor Executor,
) *ParticipantStorageClient {
	return &ParticipantStorageClient{
		address:         address,
		contract:        ps.NewParticipantStorageV1(),
		auditManagerCtr: AuditManagerV1.NewAuditManagerV1(),

		executor: executor,

		plChainID:  plChainID,
		venChainID: venChainID,
	}
}

func (c *ParticipantStorageClient) GetVenOperatorChainInfo(
	ctx context.Context,
	blockNumber *big.Int,
) (ps.ParticipantStructsPrivacyNodeViewData, error) {
	calldata := c.contract.PackGetChainViewData(c.venChainID)

	raw, err := c.executor.Call(ctx, c.address, calldata)
	if err != nil {
		return ps.ParticipantStructsPrivacyNodeViewData{}, WrapInParticipantStorageClientError(
			"failed to get VenOperator chain info",
			withstack.Wrap(err),
		)
	}

	venChainInfos, err := c.contract.UnpackGetChainViewData(raw)
	if err != nil {
		return ps.ParticipantStructsPrivacyNodeViewData{}, WrapInParticipantStorageClientError(
			"failed to unpack VenOperator chain info",
			withstack.Wrap(err),
		)
	}

	if len(venChainInfos) == 0 {
		return ps.ParticipantStructsPrivacyNodeViewData{}, ErrNoChainInfo
	}

	// If blockNumber is nil, return the latest chain info
	if blockNumber == nil {
		latestChainInfo := venChainInfos[0]
		for _, ci := range venChainInfos {
			if latestChainInfo.BlockNumber.Cmp(ci.BlockNumber) < 0 {
				latestChainInfo = ci
			}
		}
		return latestChainInfo, nil
	}

	// Otherwise, find the most recent info at or before the specified block number
	slices.SortFunc(venChainInfos, func(a, b ps.ParticipantStructsPrivacyNodeViewData) int {
		// Sort in descending order
		return b.BlockNumber.Cmp(a.BlockNumber)
	})

	for _, ci := range venChainInfos {
		if ci.BlockNumber.Cmp(blockNumber) <= 0 {
			return ci, nil
		}
	}

	return ps.ParticipantStructsPrivacyNodeViewData{}, ErrNoChainInfo
}

func (c *ParticipantStorageClient) GetMyChainInfo(
	ctx context.Context,
	blockNumber *big.Int,
) (ps.ParticipantStructsPrivacyNodeViewData, error) {
	calldata := c.contract.PackGetChainViewData(c.plChainID)

	raw, err := c.executor.Call(ctx, c.address, calldata)
	if err != nil {
		return ps.ParticipantStructsPrivacyNodeViewData{}, WrapInParticipantStorageClientError(
			"failed to get my chain info",
			withstack.Wrap(err),
		)
	}

	myChainInfos, err := c.contract.UnpackGetChainViewData(raw)
	if err != nil {
		return ps.ParticipantStructsPrivacyNodeViewData{}, WrapInParticipantStorageClientError(
			"failed to unpack my chain info",
			withstack.Wrap(err),
		)
	}

	slices.SortFunc(myChainInfos, func(a, b ps.ParticipantStructsPrivacyNodeViewData) int {
		// Sort in decending order
		return b.BlockNumber.Cmp(a.BlockNumber)
	})

	for _, ci := range myChainInfos {
		if ci.BlockNumber.Cmp(blockNumber) <= 0 {
			return ci, nil
		}
	}

	return ps.ParticipantStructsPrivacyNodeViewData{}, ErrNoChainInfo
}

func (c *ParticipantStorageClient) GetChainViewData(
	ctx context.Context,
	chainID, blockNumber *big.Int,
) (ps.ParticipantStructsPrivacyNodeViewData, error) {
	calldata := c.contract.PackGetChainViewData(chainID)

	raw, err := c.executor.Call(ctx, c.address, calldata)
	if err != nil {
		return ps.ParticipantStructsPrivacyNodeViewData{}, WrapInParticipantStorageClientError(
			"failed to get my chain info",
			withstack.Wrap(err),
		)
	}

	chainInfos, err := c.contract.UnpackGetChainViewData(raw)
	if err != nil {
		return ps.ParticipantStructsPrivacyNodeViewData{}, WrapInParticipantStorageClientError(
			"failed to unpack my chain info",
			withstack.Wrap(err),
		)
	}

	slices.SortFunc(chainInfos, func(a, b ps.ParticipantStructsPrivacyNodeViewData) int {
		// Sort in decending order
		return b.BlockNumber.Cmp(a.BlockNumber)
	})

	for _, ci := range chainInfos {
		if ci.BlockNumber.Cmp(blockNumber) <= 0 {
			return ci, nil
		}
	}

	return ps.ParticipantStructsPrivacyNodeViewData{}, ErrNoChainInfo
}

func (c *ParticipantStorageClient) GetChainViewDataBatch(
	ctx context.Context,
	chainIDs []*big.Int,
) ([]ps.ParticipantStructsPrivacyNodeViewData, error) {
	calldata := c.contract.PackGetParticipantDataBatch()

	raw, err := c.executor.Call(ctx, c.address, calldata)
	if err != nil {
		return []ps.ParticipantStructsPrivacyNodeViewData{}, WrapInParticipantStorageClientError(
			"failed to get participant data batch",
			withstack.Wrap(err),
		)
	}

	dataBatch, err := c.contract.UnpackGetParticipantDataBatch(raw)
	if err != nil {
		return []ps.ParticipantStructsPrivacyNodeViewData{}, WrapInParticipantStorageClientError(
			"failed to unpack participant data batch",
			withstack.Wrap(err),
		)
	}

	// Build a map for O(1) lookup by chain ID
	chainInfoByID := make(map[string]ps.ParticipantStructsPrivacyNodeViewData)
	for _, chainInfo := range dataBatch.PnViewData {
		chainInfoByID[chainInfo.ChainId.String()] = chainInfo
	}

	// Return chain infos in the same order as the requested chainIDs
	chainInfos := make([]ps.ParticipantStructsPrivacyNodeViewData, 0, len(chainIDs))
	for _, chainID := range chainIDs {
		if chainInfo, ok := chainInfoByID[chainID.String()]; ok {
			chainInfos = append(chainInfos, chainInfo)
		}
	}

	return chainInfos, nil
}

func (c *ParticipantStorageClient) GetAllParticipantsViewData(
	ctx context.Context,
	blockNumber *big.Int,
) ([]ps.ParticipantStructsPrivacyNodeViewData, error) {
	calldata := c.contract.PackGetParticipantDataBatch()

	raw, err := c.executor.Call(ctx, c.address, calldata)
	if err != nil {
		return []ps.ParticipantStructsPrivacyNodeViewData{}, WrapInParticipantStorageClientError(
			"failed to get participant data batch",
			withstack.Wrap(err),
		)
	}

	dataBatch, err := c.contract.UnpackGetParticipantDataBatch(raw)
	if err != nil {
		return []ps.ParticipantStructsPrivacyNodeViewData{}, WrapInParticipantStorageClientError(
			"failed to unpack participant data batch",
			withstack.Wrap(err),
		)
	}

	slices.SortFunc(dataBatch.PnViewData, func(a, b ps.ParticipantStructsPrivacyNodeViewData) int {
		// Sort in decending order
		return b.BlockNumber.Cmp(a.BlockNumber)
	})

	chainInfoByID := map[string]ps.ParticipantStructsPrivacyNodeViewData{}

	// Filter unique chain infos
	for _, ci := range dataBatch.PnViewData {
		if _, ok := chainInfoByID[ci.ChainId.String()]; ok {
			continue
		}
		if ci.BlockNumber.Cmp(blockNumber) <= 0 {
			chainInfoByID[ci.ChainId.String()] = ci
		}
	}

	chainInfos := make([]ps.ParticipantStructsPrivacyNodeViewData, 0, len(chainInfoByID))
	for _, ci := range chainInfoByID {
		chainInfos = append(chainInfos, ci)
	}

	return chainInfos, nil
}

func (c *ParticipantStorageClient) GetEnygmaParticipants(ctx context.Context) ([]*big.Int, error) {
	calldata := c.contract.PackGetEnygmaAllParticipantsChainIds()

	raw, err := c.executor.Call(ctx, c.address, calldata)
	if err != nil {
		return []*big.Int{}, WrapInParticipantStorageClientError(
			"failed to get enygma participants",
			withstack.Wrap(err),
		)
	}

	participants, err := c.contract.UnpackGetEnygmaAllParticipantsChainIds(raw)
	if err != nil {
		return []*big.Int{}, WrapInParticipantStorageClientError(
			"failed to unpack enygma participants",
			withstack.Wrap(err),
		)
	}
	return participants, nil
}

func (c *ParticipantStorageClient) GetAuditInfo(
	ctx context.Context,
	chainID, blockNumber *big.Int,
) (ps.ParticipantStructsAuditInfoData, error) {
	calldata := c.contract.PackGetAuditInfo(chainID)

	raw, err := c.executor.Call(ctx, c.address, calldata)
	if err != nil {
		return ps.ParticipantStructsAuditInfoData{}, WrapInParticipantStorageClientError(
			"failed to get audit info",
			withstack.Wrap(err),
		)
	}

	auditInfos, err := c.contract.UnpackGetAuditInfo(raw)
	if err != nil {
		return ps.ParticipantStructsAuditInfoData{}, WrapInParticipantStorageClientError(
			"failed to unpack audit info",
			withstack.Wrap(err),
		)
	}

	slices.SortFunc(auditInfos, func(a, b ps.ParticipantStructsAuditInfoData) int {
		// Sort in descending order
		return b.BlockNumber.Cmp(a.BlockNumber)
	})

	for _, info := range auditInfos {
		if info.BlockNumber.Cmp(blockNumber) <= 0 {
			return info, nil
		}
	}

	return ps.ParticipantStructsAuditInfoData{}, ErrNoAuditInfo
}

func (c *ParticipantStorageClient) GetPaymentSpendPublicKey(ctx context.Context, chainID *big.Int) (*big.Int, error) {
	calldata := c.contract.PackGetPaymentSpendPublicKeyByChainId(chainID)

	raw, err := c.executor.Call(ctx, c.address, calldata)
	if err != nil {
		return nil, WrapInParticipantStorageClientError("failed to get enygma public key", withstack.Wrap(err))
	}

	paymentPublicKey, err := c.contract.UnpackGetPaymentSpendPublicKeyByChainId(raw)
	if err != nil {
		return nil, WrapInParticipantStorageClientError("failed to unpack enygma public key", withstack.Wrap(err))
	}

	// Check if the key is zero (not set)
	if paymentPublicKey.Sign() == 0 {
		return nil, ErrNoPaymentSpendPublicKey
	}

	return paymentPublicKey, nil
}

func (c *ParticipantStorageClient) SetChainViewData(ctx context.Context, chainID *big.Int, raylsViewKey string, blockNumber *big.Int) error {
	calldata := c.contract.PackSetChainViewData(chainID, raylsViewKey, blockNumber)

	id := IDFor("participantStorage.SetChainViewData", chainID.String(), blockNumber.String())
	_, err := c.executor.Execute(ctx, id, calldata, c.address)
	if err != nil {
		return WrapInParticipantStorageClientError("failed to set chain enygma rayls view key", withstack.Wrap(err))
	}

	return nil
}

func (c *ParticipantStorageClient) SetAuditInfo(
	ctx context.Context,
	chainID *big.Int,
	blockNumber *big.Int,
	raylsViewKey string,
	encrPrivateKey, mac []byte,
) error {
	calldata := c.contract.PackSetAuditInfo(chainID, raylsViewKey, encrPrivateKey, mac, blockNumber)

	id := IDFor("participantStorage.SetAuditInfo", chainID.String(), blockNumber.String())
	_, err := c.executor.Execute(ctx, id, calldata, c.address)
	if err != nil {
		return WrapInParticipantStorageClientError("failed to set audit info", withstack.Wrap(err))
	}

	return nil
}

func (c *ParticipantStorageClient) SetPaymentSpendPublicKey(
	ctx context.Context,
	chainID *big.Int,
	paymentSpendPublicKey *big.Int,
	plAddresses []common.Address,
) error {
	calldata := c.contract.PackSetPaymentSpendPublicKey(chainID, paymentSpendPublicKey, plAddresses)

	id := IDFor("participantStorage.SetPaymentSpendPublicKey", chainID.String(), paymentSpendPublicKey.String())
	_, err := c.executor.Execute(ctx, id, calldata, c.address)
	if err != nil {
		return WrapInParticipantStorageClientError("failed to set enygma public key", withstack.Wrap(err))
	}

	return nil
}

func (c *ParticipantStorageClient) InitiateKeyAgreement(
	ctx context.Context,
	toChainID *big.Int,
	ciphertext []byte,
	digest []byte,
	blockNumber *big.Int,
) error {
	calldata := c.contract.PackInitiateKeyAgreement(c.plChainID, toChainID, ciphertext, digest, blockNumber)

	id := IDFor("participantStorage.InitiateKeyAgreement", toChainID.String(), blockNumber.String())
	_, err := c.executor.Execute(ctx, id, calldata, c.address)
	if err != nil {
		if IsRevertWithSelector(err, AuditManagerV1.AuditManagerV1AuditManagerV1BlockNumberLowerThanLatestKeyAgreementErrorID()) {
			return ErrOutdatedKeyAgreement
		}
		return WrapInParticipantStorageClientError("failed to initiate key agreement", withstack.Wrap(err))
	}

	return nil
}

func (c *ParticipantStorageClient) GetKeyAgreements(
	ctx context.Context,
	chainID *big.Int,
	blockNumber *big.Int,
) ([]ps.ParticipantStructsKeyAgreementData, error) {
	calldata := c.contract.PackGetKeyAgreements(chainID)

	raw, err := c.executor.Call(ctx, c.address, calldata)
	if err != nil {
		return []ps.ParticipantStructsKeyAgreementData{}, WrapInParticipantStorageClientError(
			"failed to get key agreements",
			withstack.Wrap(err),
		)
	}

	keyAgreements, err := c.contract.UnpackGetKeyAgreements(raw)
	if err != nil {
		return []ps.ParticipantStructsKeyAgreementData{}, WrapInParticipantStorageClientError(
			"failed to unpack key agreements",
			withstack.Wrap(err),
		)
	}

	slices.SortFunc(keyAgreements, func(a, b ps.ParticipantStructsKeyAgreementData) int {
		// Sort in decending order
		return b.BlockNumber.Cmp(a.BlockNumber)
	})

	agreementByChainID := map[string]ps.ParticipantStructsKeyAgreementData{}

	// Filter unique chain ids
	for _, agreement := range keyAgreements {
		if _, ok := agreementByChainID[agreement.ChainId.String()]; ok {
			continue
		}
		if agreement.BlockNumber.Cmp(blockNumber) <= 0 {
			agreementByChainID[agreement.ChainId.String()] = agreement
		}
	}

	keyAgreementsData := make([]ps.ParticipantStructsKeyAgreementData, 0, len(agreementByChainID))
	for _, agreement := range agreementByChainID {
		keyAgreementsData = append(keyAgreementsData, agreement)
	}

	return keyAgreementsData, nil
}
