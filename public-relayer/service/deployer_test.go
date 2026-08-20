// Decommissioning Teleport (vanilla, atomic).

package service_test

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/msgqueue"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/public-relayer/service"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/testtools"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func wrapDeployment(dpl service.Deployment) msgqueue.Message[service.Deployment] {
	return msgqueue.Message[service.Deployment]{
		V:   dpl,
		Ack: func(ctx context.Context) error { return nil },
	}
}

func wrapDeploymentWithAck(
	dpl service.Deployment,
	ack func(context.Context) error,
) msgqueue.Message[service.Deployment] {
	return msgqueue.Message[service.Deployment]{
		V:   dpl,
		Ack: ack,
	}
}

func wrapDeploymentWithAckCounter(dpl service.Deployment, counter *atomic.Int32) msgqueue.Message[service.Deployment] {
	return msgqueue.Message[service.Deployment]{
		V: dpl,
		Ack: func(ctx context.Context) error {
			counter.Add(1)
			return nil
		},
	}
}

func TestDeployer_Deploy(t *testing.T) {
	dpl := service.Deployment{
		PrivateAddress: common.HexToAddress("0xdeadbeef"),

		Standard: service.ERC20,

		Name:   "TestToken",
		Symbol: "ZK",
	}
	endpointAddress := common.HexToAddress("0xdeadc0de")
	publicAddress := common.HexToAddress("0xc0febabe")

	var respectsContextOnError bool

	testtools.SilenceLogger()

	t.Run("supports graceful shutdown", func(t *testing.T) {
		consumer := &DeploymentConsumerMock{
			NextFunc: func(ctx context.Context) (msgqueue.Message[service.Deployment], error) {
				<-ctx.Done()
				return msgqueue.Message[service.Deployment]{}, ctx.Err()
			},
		}
		deployerClient := &DeployerClientMock{
			DeployPublicChainERC20Func: func(_ context.Context, name, symbol string, privateAddress, raylsPublicRelayerEndpoint common.Address) (common.Address, error) {
				return publicAddress, nil
			},
		}
		tokentGovernanceClient := &TokenGovernanceClientMock{
			GetPublicAddressByPrivateAddressFunc: func(_ context.Context, privateAddress common.Address) (common.Address, error) {
				return common.Address{}, nil
			},
			UpdatePublicTokenAddressFunc: func(_ context.Context, privateAddress, publicAddress common.Address) error {
				return nil
			},
		}
		accessManagerClient := &AccessManagerClientMock{
			GrantAuthorizedSenderRoleFunc: func(ctx context.Context, tokenAddress common.Address) error {
				return nil
			},
		}

		svc := service.NewDeployerService(
			consumer,
			endpointAddress,
			deployerClient,
			tokentGovernanceClient,
			accessManagerClient,
		)

		hasGracefulShutdown := testtools.ShutdownFixture(t, svc.Deploy, 500*time.Millisecond)

		assert.True(t, hasGracefulShutdown)
	})

	t.Run("doesn't skip context check on error", func(t *testing.T) {
		consumer := &DeploymentConsumerMock{
			NextFunc: func(ctx context.Context) (msgqueue.Message[service.Deployment], error) {
				<-ctx.Done()
				return msgqueue.Message[service.Deployment]{}, ctx.Err()
			},
		}
		deployerClient := &DeployerClientMock{
			DeployPublicChainERC20Func: func(_ context.Context, name, symbol string, privateAddress, raylsPublicRelayerEndpoint common.Address) (common.Address, error) {
				return publicAddress, errors.New("example error")
			},
		}
		tokentGovernanceClient := &TokenGovernanceClientMock{
			GetPublicAddressByPrivateAddressFunc: func(_ context.Context, privateAddress common.Address) (common.Address, error) {
				return common.Address{}, nil
			},
			UpdatePublicTokenAddressFunc: func(_ context.Context, privateAddress, publicAddress common.Address) error {
				return nil
			},
		}
		accessManagerClient := &AccessManagerClientMock{
			GrantAuthorizedSenderRoleFunc: func(ctx context.Context, tokenAddress common.Address) error {
				return nil
			},
		}

		svc := service.NewDeployerService(
			consumer,
			endpointAddress,
			deployerClient,
			tokentGovernanceClient,
			accessManagerClient,
		)

		respectsContextOnError = testtools.ShutdownFixture(t, svc.Deploy, 500*time.Millisecond)

		assert.True(t, respectsContextOnError)
	})

	if !respectsContextOnError {
		t.Fatal("function doesn't observe context on error - cannot proceed with further tests")
	}

	t.Run("waits on Next between messages", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		var callCount atomic.Int32
		consumer := &DeploymentConsumerMock{
			NextFunc: func(ctx context.Context) (msgqueue.Message[service.Deployment], error) {
				if callCount.Add(1) > 1 {
					cancel()
					return msgqueue.Message[service.Deployment]{}, ctx.Err()
				}
				return wrapDeployment(dpl), nil
			},
		}
		deployerClient := &DeployerClientMock{
			DeployPublicChainERC20Func: func(_ context.Context, name, symbol string, privateAddress, raylsPublicRelayerEndpoint common.Address) (common.Address, error) {
				return publicAddress, nil
			},
		}
		tokentGovernanceClient := &TokenGovernanceClientMock{
			GetPublicAddressByPrivateAddressFunc: func(_ context.Context, privateAddress common.Address) (common.Address, error) {
				return common.Address{}, nil
			},
			UpdatePublicTokenAddressFunc: func(_ context.Context, privateAddress, publicAddress common.Address) error {
				return nil
			},
		}
		accessManagerClient := &AccessManagerClientMock{
			GrantAuthorizedSenderRoleFunc: func(ctx context.Context, tokenAddress common.Address) error {
				return nil
			},
		}

		svc := service.NewDeployerService(
			consumer,
			endpointAddress,
			deployerClient,
			tokentGovernanceClient,
			accessManagerClient,
		)

		err := svc.Deploy(ctx)
		require.NoError(t, err)

		assert.Len(t, consumer.NextCalls(), 2)
	})

	t.Run("acks deployment message and continues in case public address is already set", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		var ackCount atomic.Int32
		var callCount atomic.Int32
		consumer := &DeploymentConsumerMock{
			NextFunc: func(ctx context.Context) (msgqueue.Message[service.Deployment], error) {
				if callCount.Add(1) > 1 {
					cancel()
					return msgqueue.Message[service.Deployment]{}, ctx.Err()
				}
				return wrapDeploymentWithAckCounter(dpl, &ackCount), nil
			},
		}
		deployerClient := &DeployerClientMock{
			DeployPublicChainERC20Func: func(_ context.Context, name, symbol string, privateAddress, raylsPublicRelayerEndpoint common.Address) (common.Address, error) {
				assert.Fail(t, "shouldn't have deployed token")
				return publicAddress, nil
			},
		}
		tokentGovernanceClient := &TokenGovernanceClientMock{
			GetPublicAddressByPrivateAddressFunc: func(_ context.Context, privateAddress common.Address) (common.Address, error) {
				return publicAddress, nil
			},
			UpdatePublicTokenAddressFunc: func(_ context.Context, privateAddress, publicAddress common.Address) error {
				assert.Fail(t, "shouldn't have updated token public address")
				return nil
			},
		}
		accessManagerClient := &AccessManagerClientMock{
			GrantAuthorizedSenderRoleFunc: func(ctx context.Context, tokenAddress common.Address) error {
				assert.Fail(t, "shouldn't have granted AUTHORIZED_SENDER")
				return nil
			},
		}

		svc := service.NewDeployerService(
			consumer,
			endpointAddress,
			deployerClient,
			tokentGovernanceClient,
			accessManagerClient,
		)

		err := svc.Deploy(ctx)
		require.NoError(t, err)

		assert.Equal(t, int32(1), ackCount.Load())
	})

	t.Run("deploys contract and updates token public address", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		var callCount atomic.Int32
		consumer := &DeploymentConsumerMock{
			NextFunc: func(ctx context.Context) (msgqueue.Message[service.Deployment], error) {
				if callCount.Add(1) > 1 {
					cancel()
					return msgqueue.Message[service.Deployment]{}, ctx.Err()
				}
				return wrapDeployment(dpl), nil
			},
		}
		deployerClient := &DeployerClientMock{
			DeployPublicChainERC20Func: func(_ context.Context, gotName, gotSymbol string, gotPrivateTokenAddress, gotEndpointAddress common.Address) (common.Address, error) {
				assert.Equal(t, dpl.Name, gotName)
				assert.Equal(t, dpl.Symbol, gotSymbol)
				assert.Equal(t, dpl.PrivateAddress, gotPrivateTokenAddress)

				assert.Equal(t, endpointAddress, gotEndpointAddress)

				return publicAddress, nil
			},
		}
		tokentGovernanceClient := &TokenGovernanceClientMock{
			GetPublicAddressByPrivateAddressFunc: func(_ context.Context, gotPrivateAddress common.Address) (common.Address, error) {
				assert.Equal(t, dpl.PrivateAddress, gotPrivateAddress)

				return common.Address{}, nil
			},
			UpdatePublicTokenAddressFunc: func(_ context.Context, gotPrivateAddress, gotPublicAddress common.Address) error {
				assert.Equal(t, dpl.PrivateAddress, gotPrivateAddress)
				assert.Equal(t, publicAddress, gotPublicAddress)

				return nil
			},
		}
		accessManagerClient := &AccessManagerClientMock{
			GrantAuthorizedSenderRoleFunc: func(ctx context.Context, tokenAddress common.Address) error {
				return nil
			},
		}

		svc := service.NewDeployerService(
			consumer,
			endpointAddress,
			deployerClient,
			tokentGovernanceClient,
			accessManagerClient,
		)

		err := svc.Deploy(ctx)
		require.NoError(t, err)

		assert.Len(t, consumer.NextCalls(), 2)
	})

	t.Run("does not deploy when governance lookup fails", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		var ackCount atomic.Int32
		var callCount atomic.Int32
		consumer := &DeploymentConsumerMock{
			NextFunc: func(ctx context.Context) (msgqueue.Message[service.Deployment], error) {
				if callCount.Add(1) > 1 {
					cancel()
					return msgqueue.Message[service.Deployment]{}, ctx.Err()
				}
				return wrapDeploymentWithAckCounter(dpl, &ackCount), nil
			},
		}
		deployerClient := &DeployerClientMock{
			DeployPublicChainERC20Func: func(_ context.Context, name, symbol string, privateAddress, raylsPublicRelayerEndpoint common.Address) (common.Address, error) {
				assert.Fail(t, "shouldn't have deployed token when governance lookup failed")
				return publicAddress, nil
			},
		}
		tokentGovernanceClient := &TokenGovernanceClientMock{
			GetPublicAddressByPrivateAddressFunc: func(_ context.Context, privateAddress common.Address) (common.Address, error) {
				return common.Address{}, errors.New("governance-lookup-error")
			},
			UpdatePublicTokenAddressFunc: func(_ context.Context, privateAddress, publicAddress common.Address) error {
				assert.Fail(t, "shouldn't have updated token public address")
				return nil
			},
		}
		accessManagerClient := &AccessManagerClientMock{
			GrantAuthorizedSenderRoleFunc: func(ctx context.Context, tokenAddress common.Address) error {
				assert.Fail(t, "shouldn't have granted AUTHORIZED_SENDER")
				return nil
			},
		}

		svc := service.NewDeployerService(
			consumer,
			endpointAddress,
			deployerClient,
			tokentGovernanceClient,
			accessManagerClient,
		)

		err := svc.Deploy(ctx)
		require.NoError(t, err)

		assert.Equal(t, int32(0), ackCount.Load(), "message should not be acked so NATS redelivers")
		assert.Len(t, deployerClient.DeployPublicChainERC20Calls(), 0, "should not attempt deployment")
	})

	t.Run("continues on failure to deploy contract", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		var callCount atomic.Int32
		consumer := &DeploymentConsumerMock{
			NextFunc: func(ctx context.Context) (msgqueue.Message[service.Deployment], error) {
				if callCount.Add(1) > 1 {
					cancel()
					return msgqueue.Message[service.Deployment]{}, ctx.Err()
				}
				return wrapDeployment(dpl), nil
			},
		}
		deployerClient := &DeployerClientMock{
			DeployPublicChainERC20Func: func(_ context.Context, name, symbol string, privateAddress, raylsPublicRelayerEndpoint common.Address) (common.Address, error) {
				return publicAddress, errors.New("example-error")
			},
		}
		tokentGovernanceClient := &TokenGovernanceClientMock{
			GetPublicAddressByPrivateAddressFunc: func(_ context.Context, privateAddress common.Address) (common.Address, error) {
				return common.Address{}, nil
			},
			UpdatePublicTokenAddressFunc: func(_ context.Context, privateAddress, publicAddress common.Address) error {
				assert.Fail(t, "shouldn't have updated token public address")
				return nil
			},
		}
		accessManagerClient := &AccessManagerClientMock{
			GrantAuthorizedSenderRoleFunc: func(ctx context.Context, tokenAddress common.Address) error {
				assert.Fail(t, "shouldn't have granted AUTHORIZED_SENDER")
				return nil
			},
		}

		svc := service.NewDeployerService(
			consumer,
			endpointAddress,
			deployerClient,
			tokentGovernanceClient,
			accessManagerClient,
		)

		err := svc.Deploy(ctx)
		require.NoError(t, err)

		assert.Len(t, consumer.NextCalls(), 2)
	})

	t.Run("grants AUTHORIZED_SENDER to deployed token", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		var callCount atomic.Int32
		consumer := &DeploymentConsumerMock{
			NextFunc: func(ctx context.Context) (msgqueue.Message[service.Deployment], error) {
				if callCount.Add(1) > 1 {
					cancel()
					return msgqueue.Message[service.Deployment]{}, ctx.Err()
				}
				return wrapDeployment(dpl), nil
			},
		}
		deployerClient := &DeployerClientMock{
			DeployPublicChainERC20Func: func(_ context.Context, name, symbol string, privateAddress, raylsPublicRelayerEndpoint common.Address) (common.Address, error) {
				return publicAddress, nil
			},
		}
		tokentGovernanceClient := &TokenGovernanceClientMock{
			GetPublicAddressByPrivateAddressFunc: func(_ context.Context, privateAddress common.Address) (common.Address, error) {
				return common.Address{}, nil
			},
			UpdatePublicTokenAddressFunc: func(_ context.Context, privateAddress, publicAddress common.Address) error {
				return nil
			},
		}
		accessManagerClient := &AccessManagerClientMock{
			GrantAuthorizedSenderRoleFunc: func(ctx context.Context, gotTokenAddress common.Address) error {
				assert.Equal(t, publicAddress, gotTokenAddress)
				return nil
			},
		}

		svc := service.NewDeployerService(
			consumer,
			endpointAddress,
			deployerClient,
			tokentGovernanceClient,
			accessManagerClient,
		)

		err := svc.Deploy(ctx)
		require.NoError(t, err)

		assert.Len(t, accessManagerClient.GrantAuthorizedSenderRoleCalls(), 1)
	})

	t.Run("continues and acks message on failure to grant role", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		var ackCount atomic.Int32
		var callCount atomic.Int32
		consumer := &DeploymentConsumerMock{
			NextFunc: func(ctx context.Context) (msgqueue.Message[service.Deployment], error) {
				if callCount.Add(1) > 1 {
					cancel()
					return msgqueue.Message[service.Deployment]{}, ctx.Err()
				}
				return wrapDeploymentWithAckCounter(dpl, &ackCount), nil
			},
		}
		deployerClient := &DeployerClientMock{
			DeployPublicChainERC20Func: func(_ context.Context, name, symbol string, privateAddress, raylsPublicRelayerEndpoint common.Address) (common.Address, error) {
				return publicAddress, nil
			},
		}
		tokentGovernanceClient := &TokenGovernanceClientMock{
			GetPublicAddressByPrivateAddressFunc: func(_ context.Context, privateAddress common.Address) (common.Address, error) {
				return common.Address{}, nil
			},
			UpdatePublicTokenAddressFunc: func(_ context.Context, privateAddress, publicAddress common.Address) error {
				return nil
			},
		}
		accessManagerClient := &AccessManagerClientMock{
			GrantAuthorizedSenderRoleFunc: func(ctx context.Context, tokenAddress common.Address) error {
				return errors.New("authorization-failed")
			},
		}

		svc := service.NewDeployerService(
			consumer,
			endpointAddress,
			deployerClient,
			tokentGovernanceClient,
			accessManagerClient,
		)

		err := svc.Deploy(ctx)
		require.NoError(t, err)

		assert.Equal(t, int32(1), ackCount.Load())
	})
}
