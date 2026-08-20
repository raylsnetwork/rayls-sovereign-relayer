// Decommissioning Teleport (vanilla, atomic).

package handler_test

import (
	"context"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/contracts/RNMessageDispatcherV1"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/contracts/PNTokenCoreV1"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/contracts/PNTokenRegistryV1"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/public-relayer/service"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/testtools"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPublicHandler(t *testing.T) {
	testtools.SilenceLogger()

	t.Run("pushes parsed MessageDispatched event to message queue", func(t *testing.T) {
		wantMessageID := [32]byte{0x01, 0x02, 0x03}
		wantFrom := common.HexToAddress("0xaaaa")
		wantToChainID := big.NewInt(42)
		wantTo := common.HexToAddress("0xbbbb")

		event := newMessageDispatchedEvent(
			withMessageID(wantMessageID),
			withFrom(wantFrom),
			withToChainID(wantToChainID),
			withTo(wantTo),
		)

		stub := newHandlerStub(
			t,
			withDispatcher(
				func(log *ethTypes.Log) (*RNMessageDispatcherV1.RNMessageDispatcherV1MessageDispatched, error) {
					return event, nil
				},
			),
			withMessagePublisher(func(ctx context.Context, msg service.Message) error {
				assert.Equal(t, common.BytesToHash(wantMessageID[:]), msg.ID)
				assert.Equal(t, wantFrom, msg.FromAddress)
				assert.Equal(t, wantToChainID, msg.ToChainID)
				assert.Equal(t, wantTo, msg.ToAddress)
				return nil
			}),
		)

		h := stub.newPublicHandler()
		err := h.Handle(context.TODO(), []ethTypes.Log{{}})
		require.NoError(t, err)
	})

	t.Run("processes multiple logs", func(t *testing.T) {
		eventA := newMessageDispatchedEvent(withMessageID([32]byte{0x01}))
		eventB := newMessageDispatchedEvent(withMessageID([32]byte{0x02}))

		callCount := 0
		stub := newHandlerStub(
			t,
			withDispatcher(
				func(log *ethTypes.Log) (*RNMessageDispatcherV1.RNMessageDispatcherV1MessageDispatched, error) {
					callCount++
					if callCount == 1 {
						return eventA, nil
					}
					return eventB, nil
				},
			),
			withMessagePublisher(func(ctx context.Context, msg service.Message) error {
				return nil
			}),
		)

		h := stub.newPublicHandler()
		err := h.Handle(context.TODO(), []ethTypes.Log{{}, {}})
		require.NoError(t, err)
		assert.Equal(t, 2, callCount)
	})

	t.Run("skips logs that do not match MessageDispatched", func(t *testing.T) {
		stub := newHandlerStub(t)
		// Default dispatcher returns errNotMatched, default publisher fails if called.
		// This verifies that unrecognized logs are silently skipped.

		h := stub.newPublicHandler()
		err := h.Handle(context.TODO(), []ethTypes.Log{{}})
		require.NoError(t, err)
	})

	t.Run("does not process PublicChainStatusUpdated events in public mode", func(t *testing.T) {
		// In public mode, token core / registry are nil. Only MessageDispatched events should be processed.
		stub := newHandlerStub(t)

		h := stub.newPublicHandler()
		err := h.Handle(context.TODO(), []ethTypes.Log{{}})
		require.NoError(t, err)
	})

	t.Run("returns nil on empty logs", func(t *testing.T) {
		stub := newHandlerStub(t)

		h := stub.newPublicHandler()
		err := h.Handle(context.TODO(), []ethTypes.Log{})
		require.NoError(t, err)
	})
}

func TestPrivateHandler(t *testing.T) {
	testtools.SilenceLogger()

	t.Run("pushes MessageDispatched to message queue", func(t *testing.T) {
		wantMessageID := [32]byte{0xdd, 0xee, 0xff}
		event := newMessageDispatchedEvent(withMessageID(wantMessageID))

		stub := newHandlerStub(
			t,
			withDispatcher(
				func(log *ethTypes.Log) (*RNMessageDispatcherV1.RNMessageDispatcherV1MessageDispatched, error) {
					return event, nil
				},
			),
			withMessagePublisher(func(ctx context.Context, msg service.Message) error {
				assert.Equal(t, common.BytesToHash(wantMessageID[:]), msg.ID)
				return nil
			}),
		)

		h := stub.newPrivateHandler()
		err := h.Handle(context.TODO(), []ethTypes.Log{{}})
		require.NoError(t, err)
	})

	t.Run("pushes ERC20 pending-deployment to deployment queue", func(t *testing.T) {
		wantAddr := common.HexToAddress("0xdeadbeef")
		event := newPublicChainStatusUpdatedEvent(withEventTokenAddress(wantAddr))
		token := newTokenStruct(
			withTokenAddress(wantAddr),
			withStandard(1), // ERC20
			withName("MyToken"),
			withSymbol("MT"),
			withURI("https://example.com"),
		)

		stub := newHandlerStub(t,
			withTokenCore(func(log *ethTypes.Log) (*PNTokenCoreV1.PNTokenCoreV1PublicChainStatusUpdated, error) {
				return event, nil
			}),
			withTokenRegistry(func(ctx context.Context, addr common.Address) (PNTokenRegistryV1.TokenStructsToken, error) {
				assert.Equal(t, wantAddr, addr)
				return token, nil
			}),
			withDeploymentPublisher(func(ctx context.Context, dpl service.Deployment) error {
				assert.Equal(t, wantAddr, dpl.PrivateAddress)
				assert.Equal(t, service.ERC20, dpl.Standard)
				assert.Equal(t, "MyToken", dpl.Name)
				assert.Equal(t, "MT", dpl.Symbol)
				assert.Equal(t, "https://example.com", dpl.URI)
				return nil
			}),
		)

		h := stub.newPrivateHandler()
		err := h.Handle(context.TODO(), []ethTypes.Log{{}})
		require.NoError(t, err)
	})

	t.Run("pushes ERC721 pending-deployment to deployment queue", func(t *testing.T) {
		stub := newHandlerStub(t,
			withTokenCore(func(log *ethTypes.Log) (*PNTokenCoreV1.PNTokenCoreV1PublicChainStatusUpdated, error) {
				return newPublicChainStatusUpdatedEvent(), nil
			}),
			withTokenRegistry(func(ctx context.Context, addr common.Address) (PNTokenRegistryV1.TokenStructsToken, error) {
				return newTokenStruct(withStandard(2)), nil // ERC721
			}),
			withDeploymentPublisher(func(ctx context.Context, dpl service.Deployment) error {
				assert.Equal(t, service.ERC721, dpl.Standard)
				return nil
			}),
		)

		h := stub.newPrivateHandler()
		err := h.Handle(context.TODO(), []ethTypes.Log{{}})
		require.NoError(t, err)
	})

	t.Run("pushes ERC1155 pending-deployment to deployment queue", func(t *testing.T) {
		stub := newHandlerStub(t,
			withTokenCore(func(log *ethTypes.Log) (*PNTokenCoreV1.PNTokenCoreV1PublicChainStatusUpdated, error) {
				return newPublicChainStatusUpdatedEvent(), nil
			}),
			withTokenRegistry(func(ctx context.Context, addr common.Address) (PNTokenRegistryV1.TokenStructsToken, error) {
				return newTokenStruct(withStandard(3)), nil // ERC1155
			}),
			withDeploymentPublisher(func(ctx context.Context, dpl service.Deployment) error {
				assert.Equal(t, service.ERC1155, dpl.Standard)
				return nil
			}),
		)

		h := stub.newPrivateHandler()
		err := h.Handle(context.TODO(), []ethTypes.Log{{}})
		require.NoError(t, err)
	})

	t.Run("ignores public-chain status transitions other than PENDING_DEPLOYMENT", func(t *testing.T) {
		// newStatus == DEPLOYED (2) is not the deploy cue; no metadata fetch, no deployment.
		stub := newHandlerStub(t,
			withTokenCore(func(log *ethTypes.Log) (*PNTokenCoreV1.PNTokenCoreV1PublicChainStatusUpdated, error) {
				return newPublicChainStatusUpdatedEvent(withNewStatus(2)), nil
			}),
			withTokenRegistry(func(ctx context.Context, addr common.Address) (PNTokenRegistryV1.TokenStructsToken, error) {
				assert.Fail(t, "should not read the registry for a non-pending status")
				return PNTokenRegistryV1.TokenStructsToken{}, nil
			}),
			// default DeploymentPublisher fails if called.
		)

		h := stub.newPrivateHandler()
		err := h.Handle(context.TODO(), []ethTypes.Log{{}})
		require.NoError(t, err)
	})

	t.Run("handles both MessageDispatched and PublicChainStatusUpdated in same log batch", func(t *testing.T) {
		dispatchedEvent := newMessageDispatchedEvent()

		messagePublished := false
		deploymentPublished := false

		stub := newHandlerStub(
			t,
			withDispatcher(
				func(log *ethTypes.Log) (*RNMessageDispatcherV1.RNMessageDispatcherV1MessageDispatched, error) {
					return dispatchedEvent, nil
				},
			),
			withTokenCore(func(log *ethTypes.Log) (*PNTokenCoreV1.PNTokenCoreV1PublicChainStatusUpdated, error) {
				return newPublicChainStatusUpdatedEvent(), nil
			}),
			withTokenRegistry(func(ctx context.Context, addr common.Address) (PNTokenRegistryV1.TokenStructsToken, error) {
				return newTokenStruct(), nil
			}),
			withMessagePublisher(func(ctx context.Context, msg service.Message) error {
				messagePublished = true
				return nil
			}),
			withDeploymentPublisher(func(ctx context.Context, dpl service.Deployment) error {
				deploymentPublished = true
				return nil
			}),
		)

		h := stub.newPrivateHandler()
		err := h.Handle(context.TODO(), []ethTypes.Log{{}})
		require.NoError(t, err)
		assert.True(t, messagePublished, "message should have been published")
		assert.True(t, deploymentPublished, "deployment should have been published")
	})

	t.Run("continues processing when message push fails", func(t *testing.T) {
		event := newMessageDispatchedEvent()

		stub := newHandlerStub(
			t,
			withDispatcher(
				func(log *ethTypes.Log) (*RNMessageDispatcherV1.RNMessageDispatcherV1MessageDispatched, error) {
					return event, nil
				},
			),
			withMessagePublisher(func(ctx context.Context, msg service.Message) error {
				return assert.AnError
			}),
		)

		h := stub.newPublicHandler()
		// Handler logs the error but returns nil — does not propagate push failures.
		err := h.Handle(context.TODO(), []ethTypes.Log{{}})
		require.NoError(t, err)
	})

	t.Run("continues processing when deployment push fails", func(t *testing.T) {
		stub := newHandlerStub(t,
			withTokenCore(func(log *ethTypes.Log) (*PNTokenCoreV1.PNTokenCoreV1PublicChainStatusUpdated, error) {
				return newPublicChainStatusUpdatedEvent(), nil
			}),
			withTokenRegistry(func(ctx context.Context, addr common.Address) (PNTokenRegistryV1.TokenStructsToken, error) {
				return newTokenStruct(), nil
			}),
			withDeploymentPublisher(func(ctx context.Context, dpl service.Deployment) error {
				return assert.AnError
			}),
		)

		h := stub.newPrivateHandler()
		err := h.Handle(context.TODO(), []ethTypes.Log{{}})
		require.NoError(t, err)
	})

	t.Run("skips deployment when registry lookup fails", func(t *testing.T) {
		stub := newHandlerStub(t,
			withTokenCore(func(log *ethTypes.Log) (*PNTokenCoreV1.PNTokenCoreV1PublicChainStatusUpdated, error) {
				return newPublicChainStatusUpdatedEvent(), nil
			}),
			withTokenRegistry(func(ctx context.Context, addr common.Address) (PNTokenRegistryV1.TokenStructsToken, error) {
				return PNTokenRegistryV1.TokenStructsToken{}, assert.AnError
			}),
			// default DeploymentPublisher fails if called.
		)

		h := stub.newPrivateHandler()
		err := h.Handle(context.TODO(), []ethTypes.Log{{}})
		require.NoError(t, err)
	})

	t.Run("handles unknown token standard gracefully", func(t *testing.T) {
		// Standard 0 (CUSTOM) is not mapped — getTokenStandard returns an error.
		// The handler logs the error and skips the deployment.
		stub := newHandlerStub(t,
			withTokenCore(func(log *ethTypes.Log) (*PNTokenCoreV1.PNTokenCoreV1PublicChainStatusUpdated, error) {
				return newPublicChainStatusUpdatedEvent(), nil
			}),
			withTokenRegistry(func(ctx context.Context, addr common.Address) (PNTokenRegistryV1.TokenStructsToken, error) {
				return newTokenStruct(withStandard(0)), nil
			}),
			withDeploymentPublisher(func(ctx context.Context, dpl service.Deployment) error {
				assert.Fail(t, "should not push deployment for unknown standard")
				return nil
			}),
		)

		h := stub.newPrivateHandler()
		err := h.Handle(context.TODO(), []ethTypes.Log{{}})
		require.NoError(t, err)
	})
}

func TestGetTokenStandard(t *testing.T) {
	testtools.SilenceLogger()

	tests := []struct {
		name     string
		standard uint8
		want     service.TokenStandard
		wantErr  bool
	}{
		{name: "ERC20", standard: 1, want: service.ERC20},
		{name: "ERC721", standard: 2, want: service.ERC721},
		{name: "ERC1155", standard: 3, want: service.ERC1155},
		{name: "CUSTOM (0) is unknown", standard: 0, wantErr: true},
		{name: "ENYGMA (4) is unknown", standard: 4, wantErr: true},
		{name: "ZKDVPERC721 (5) is unknown", standard: 5, wantErr: true},
		{name: "ZKDVPERC1155 (6) is unknown", standard: 6, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// getTokenStandard is unexported, so we test it indirectly through the handler.
			// We drive a PublicChainStatusUpdated event, have the registry return a token with
			// the given standard, and check what deployment gets pushed.
			var gotStandard service.TokenStandard
			pushed := false

			stub := newHandlerStub(
				t,
				withTokenCore(func(log *ethTypes.Log) (*PNTokenCoreV1.PNTokenCoreV1PublicChainStatusUpdated, error) {
					return newPublicChainStatusUpdatedEvent(), nil
				}),
				withTokenRegistry(func(ctx context.Context, addr common.Address) (PNTokenRegistryV1.TokenStructsToken, error) {
					return newTokenStruct(withStandard(tt.standard)), nil
				}),
				withDeploymentPublisher(func(ctx context.Context, dpl service.Deployment) error {
					gotStandard = dpl.Standard
					pushed = true
					return nil
				}),
			)

			h := stub.newPrivateHandler()
			err := h.Handle(context.TODO(), []ethTypes.Log{{}})
			require.NoError(t, err)

			if tt.wantErr {
				assert.False(t, pushed, "deployment should not have been pushed for unknown standard")
			} else {
				require.True(t, pushed, "deployment should have been pushed")
				assert.Equal(t, tt.want, gotStandard)
			}
		})
	}
}
