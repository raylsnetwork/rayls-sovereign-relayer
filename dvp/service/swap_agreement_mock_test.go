package service_test

import (
	"context"
	"sync"

	"github.com/ethereum/go-ethereum/common"

	"github.com/raylsnetwork/rayls-sovereign-relayer/types"
)

type swapAgreementEndpointClientMock struct {
	GetResourceAddressFunc func(ctx context.Context, resourceId string) (common.Address, error)

	calls struct {
		GetResourceAddress []struct {
			Ctx        context.Context
			ResourceId string
		}
	}
	lockGetResourceAddress sync.RWMutex
}

func (m *swapAgreementEndpointClientMock) GetResourceAddress(ctx context.Context, resourceId string) (common.Address, error) {
	m.lockGetResourceAddress.Lock()
	defer m.lockGetResourceAddress.Unlock()

	m.calls.GetResourceAddress = append(m.calls.GetResourceAddress, struct {
		Ctx        context.Context
		ResourceId string
	}{Ctx: ctx, ResourceId: resourceId})

	return m.GetResourceAddressFunc(ctx, resourceId)
}

func (m *swapAgreementEndpointClientMock) GetResourceAddressCalls() []struct {
	Ctx        context.Context
	ResourceId string
} {
	m.lockGetResourceAddress.RLock()
	defer m.lockGetResourceAddress.RUnlock()
	return m.calls.GetResourceAddress
}

type swapAgreementNotifyCall struct {
	Ctx          context.Context
	TokenAddress common.Address
	SharedId     string
	Status       types.DvpCommunicatorStatus
	Message      string
}

type swapAgreementEnygmaHandlerClientMock struct {
	NotifySenderWithPNCommunicatorFunc func(ctx context.Context, tokenAddress common.Address, sharedId string, status types.DvpCommunicatorStatus, message string) error

	calls struct {
		NotifySenderWithPNCommunicator []swapAgreementNotifyCall
	}
	lockNotifySenderWithPNCommunicator sync.RWMutex
}

func (m *swapAgreementEnygmaHandlerClientMock) NotifySenderWithPNCommunicator(ctx context.Context, tokenAddress common.Address, sharedId string, status types.DvpCommunicatorStatus, message string) error {
	m.lockNotifySenderWithPNCommunicator.Lock()
	defer m.lockNotifySenderWithPNCommunicator.Unlock()

	m.calls.NotifySenderWithPNCommunicator = append(m.calls.NotifySenderWithPNCommunicator, swapAgreementNotifyCall{
		Ctx:          ctx,
		TokenAddress: tokenAddress,
		SharedId:     sharedId,
		Status:       status,
		Message:      message,
	})

	return m.NotifySenderWithPNCommunicatorFunc(ctx, tokenAddress, sharedId, status, message)
}

func (m *swapAgreementEnygmaHandlerClientMock) NotifySenderWithPNCommunicatorCalls() []swapAgreementNotifyCall {
	m.lockNotifySenderWithPNCommunicator.RLock()
	defer m.lockNotifySenderWithPNCommunicator.RUnlock()
	return m.calls.NotifySenderWithPNCommunicator
}

type swapAgreementDvpERC721HandlerClientMock struct {
	NotifySenderWithPNCommunicatorFunc func(ctx context.Context, tokenAddress common.Address, sharedId string, status types.DvpCommunicatorStatus, message string) error

	calls struct {
		NotifySenderWithPNCommunicator []swapAgreementNotifyCall
	}
	lockNotifySenderWithPNCommunicator sync.RWMutex
}

func (m *swapAgreementDvpERC721HandlerClientMock) NotifySenderWithPNCommunicator(ctx context.Context, tokenAddress common.Address, sharedId string, status types.DvpCommunicatorStatus, message string) error {
	m.lockNotifySenderWithPNCommunicator.Lock()
	defer m.lockNotifySenderWithPNCommunicator.Unlock()

	m.calls.NotifySenderWithPNCommunicator = append(m.calls.NotifySenderWithPNCommunicator, swapAgreementNotifyCall{
		Ctx:          ctx,
		TokenAddress: tokenAddress,
		SharedId:     sharedId,
		Status:       status,
		Message:      message,
	})

	return m.NotifySenderWithPNCommunicatorFunc(ctx, tokenAddress, sharedId, status, message)
}

func (m *swapAgreementDvpERC721HandlerClientMock) NotifySenderWithPNCommunicatorCalls() []swapAgreementNotifyCall {
	m.lockNotifySenderWithPNCommunicator.RLock()
	defer m.lockNotifySenderWithPNCommunicator.RUnlock()
	return m.calls.NotifySenderWithPNCommunicator
}

type swapAgreementDvpERC1155HandlerClientMock struct {
	NotifySenderWithPNCommunicatorFunc func(ctx context.Context, tokenAddress common.Address, sharedId string, status types.DvpCommunicatorStatus, message string) error

	calls struct {
		NotifySenderWithPNCommunicator []swapAgreementNotifyCall
	}
	lockNotifySenderWithPNCommunicator sync.RWMutex
}

func (m *swapAgreementDvpERC1155HandlerClientMock) NotifySenderWithPNCommunicator(ctx context.Context, tokenAddress common.Address, sharedId string, status types.DvpCommunicatorStatus, message string) error {
	m.lockNotifySenderWithPNCommunicator.Lock()
	defer m.lockNotifySenderWithPNCommunicator.Unlock()

	m.calls.NotifySenderWithPNCommunicator = append(m.calls.NotifySenderWithPNCommunicator, swapAgreementNotifyCall{
		Ctx:          ctx,
		TokenAddress: tokenAddress,
		SharedId:     sharedId,
		Status:       status,
		Message:      message,
	})

	return m.NotifySenderWithPNCommunicatorFunc(ctx, tokenAddress, sharedId, status, message)
}

func (m *swapAgreementDvpERC1155HandlerClientMock) NotifySenderWithPNCommunicatorCalls() []swapAgreementNotifyCall {
	m.lockNotifySenderWithPNCommunicator.RLock()
	defer m.lockNotifySenderWithPNCommunicator.RUnlock()
	return m.calls.NotifySenderWithPNCommunicator
}
