package logrouter_test

import (
	"context"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/raylsnetwork/rayls-sovereign-relayer/backoff"
	"github.com/raylsnetwork/rayls-sovereign-relayer/private-relayer/dest/logrouter"
)

// logRouterStub provides default failing implementations for all LogRouter dependencies.
// This allows tests to override only the dependencies they care about, with all others
// failing if accidentally called.
type logRouterStub struct {
	t *testing.T

	Config logrouter.LogRouterConfig

	EndpointMQ     RouterMQFunc
	TeleportMQ     RouterMQFunc
	EnygmaMQ       RouterMQFunc
	DvpMQ          RouterMQFunc
	AuditManagerMQ RouterMQFunc
	Backoff        backoff.Strategy
}

type RouterMQFunc func(context.Context, logrouter.Block) error

// LogRouterStubOption represents an option for configuring the logRouterStub.
type LogRouterStubOption func(*logRouterStub)

// newLogRouterStub creates a new stub with all dependencies set to fail by default.
func newLogRouterStub(t *testing.T, opts ...LogRouterStubOption) *logRouterStub {
	t.Helper()

	stub := &logRouterStub{t: t}

	// Set default config with zero addresses
	stub.Config = logrouter.LogRouterConfig{
		EndpointAddress: common.Address{},
		TeleportAddress: common.Address{},
		EnygmaAddress:   common.Address{},
		DvpAddress:      common.Address{},
	}

	// Set default failing implementations
	stub.EndpointMQ = func(ctx context.Context, block logrouter.Block) error {
		t.Fatal("shouldn't call EndpointRouterMQ.Push")
		return nil
	}

	stub.TeleportMQ = func(ctx context.Context, block logrouter.Block) error {
		t.Fatal("shouldn't call TeleportRouterMQ.Push")
		return nil
	}

	stub.EnygmaMQ = func(ctx context.Context, block logrouter.Block) error {
		t.Fatal("shouldn't call EnygmaMQ.Push")
		return nil
	}

	stub.DvpMQ = func(ctx context.Context, block logrouter.Block) error {
		t.Fatal("shouldn't call DvpMQ.Push")
		return nil
	}

	stub.AuditManagerMQ = func(ctx context.Context, block logrouter.Block) error {
		t.Fatal("shouldn't call AuditManagerMQ.Push")
		return nil
	}

	// Set default immediate backoff for tests
	stub.Backoff = &immediateBackoff{}

	// Apply options
	for _, opt := range opts {
		opt(stub)
	}

	return stub
}

// immediateBackoff executes the function immediately without delays
type immediateBackoff struct{}

func (b *immediateBackoff) Next(attempt int) time.Duration {
	return 0
}

func (b *immediateBackoff) Do(ctx context.Context, maxAttempts int, fn func() error) error {
	var err error
	for i := 0; i < maxAttempts; i++ {
		if err = fn(); err == nil {
			return nil
		}
	}
	return err
}

// newLogRouter creates a LogRouter using the stub's configured dependencies.
func (s *logRouterStub) newLogRouter() *logrouter.LogRouter {
	s.t.Helper()

	endpointMQ := &RouterMQMock{
		PushFunc: s.EndpointMQ,
	}

	teleportMQ := &RouterMQMock{
		PushFunc: s.TeleportMQ,
	}

	enygmaMQ := &RouterMQMock{
		PushFunc: s.EnygmaMQ,
	}

	dvpMQ := &RouterMQMock{
		PushFunc: s.DvpMQ,
	}

	auditManagerMQ := &RouterMQMock{
		PushFunc: s.AuditManagerMQ,
	}

	return logrouter.NewWithCustomBackoff(s.Config, endpointMQ, teleportMQ, enygmaMQ, dvpMQ, auditManagerMQ, s.Backoff)
}

// withEndpointMQ configures the EndpointMQ dependency.
func withEndpointMQ(fn RouterMQFunc) LogRouterStubOption {
	return func(s *logRouterStub) {
		s.EndpointMQ = fn
	}
}

// withTeleportMQ configures the TeleportMQ dependency.
func withTeleportMQ(fn RouterMQFunc) LogRouterStubOption {
	return func(s *logRouterStub) {
		s.TeleportMQ = fn
	}
}

// withEnygmaMQ configures the EnygmaMQ dependency.
func withEnygmaMQ(fn RouterMQFunc) LogRouterStubOption {
	return func(s *logRouterStub) {
		s.EnygmaMQ = fn
	}
}

// withDvpMQ configures the DvpMQ dependency.
func withDvpMQ(fn RouterMQFunc) LogRouterStubOption {
	return func(s *logRouterStub) {
		s.DvpMQ = fn
	}
}

// withEndpointAddress configures the EndpointAddress in the config.
func withEndpointAddress(address common.Address) LogRouterStubOption {
	return func(s *logRouterStub) {
		s.Config.EndpointAddress = address
	}
}

// withTeleportAddress configures the TeleportAddress in the config.
func withTeleportAddress(address common.Address) LogRouterStubOption {
	return func(s *logRouterStub) {
		s.Config.TeleportAddress = address
	}
}

// withEnygmaAddress configures the EnygmaAddress in the config.
func withEnygmaAddress(address common.Address) LogRouterStubOption {
	return func(s *logRouterStub) {
		s.Config.EnygmaAddress = address
	}
}

// withDvpAddress configures the DvpAddress in the config.
func withDvpAddress(address common.Address) LogRouterStubOption {
	return func(s *logRouterStub) {
		s.Config.DvpAddress = address
	}
}

// withConfig sets the entire config at once.
func withConfig(config logrouter.LogRouterConfig) LogRouterStubOption {
	return func(s *logRouterStub) {
		s.Config = config
	}
}

// noopEndpointMQ returns a no-op implementation that doesn't fail.
func noopEndpointMQ() RouterMQFunc {
	return func(ctx context.Context, block logrouter.Block) error {
		return nil
	}
}

// noopTeleportMQ returns a no-op implementation that doesn't fail.
func noopTeleportMQ() RouterMQFunc {
	return func(ctx context.Context, block logrouter.Block) error {
		return nil
	}
}
