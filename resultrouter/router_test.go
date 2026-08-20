package resultrouter_test

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"
	"time"

	"github.com/raylsnetwork/rayls-privacy-relayer-api/msgqueue"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/resultrouter"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/testtools"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/testtools/spy"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func defaultConfig() resultrouter.Config {
	return resultrouter.Config{
		Identity:  "test-identity",
		BatchSize: 100,
		Interval:  time.Second,
	}
}

func emptyConsumer() *ConsumerMock {
	return &ConsumerMock{
		FetchFunc: func(ctx context.Context, count int) ([]msgqueue.Message[types.TxResult], error) {
			return nil, nil //nolint:nilnil // intentional empty fetch in test
		},
	}
}

func TestRouter(t *testing.T) {
	testtools.SilenceLogger()

	t.Run("supports graceful shutdown", func(t *testing.T) {
		r := resultrouter.New(defaultConfig(), emptyConsumer())

		hasGracefulShutdown := testtools.ShutdownFixture(t, r.Run, 10*time.Millisecond)
		assert.True(t, hasGracefulShutdown)
	})

	t.Run("respects context on fetch error", func(t *testing.T) {
		cons := &ConsumerMock{
			FetchFunc: func(ctx context.Context, count int) ([]msgqueue.Message[types.TxResult], error) {
				return nil, errors.New("transport blew up")
			},
		}
		r := resultrouter.New(defaultConfig(), cons)

		respects := testtools.ShutdownFixture(t, r.Run, time.Millisecond)
		assert.True(t, respects)
	})

	t.Run("returns to ticker on empty fetch", func(t *testing.T) {
		cons := emptyConsumer()
		handler := &HandlerMock{
			HandleFunc: func(ctx context.Context, results []types.TxResult) error {
				assert.Fail(t, "handler shouldn't be called when fetch is empty")
				return nil
			},
		}
		r := resultrouter.New(defaultConfig(), cons)
		r.Register("any.message.type", handler)

		ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond)
		defer cancel()
		require.NoError(t, r.Run(ctx))

		assert.GreaterOrEqual(t, len(cons.FetchCalls()), 1)
		assert.Empty(t, handler.HandleCalls())
	})

	t.Run("dispatches results to the handler matching their MessageType", func(t *testing.T) {
		ackA := spy.NewAck()
		ackB := spy.NewAck()
		batch := []msgqueue.Message[types.TxResult]{
			{V: types.TxResult{CorrelationID: "id-a", MessageType: "type.a", Kind: types.TxResultSuccess}, Ack: ackA.Fn()},
			{V: types.TxResult{CorrelationID: "id-b", MessageType: "type.a", Kind: types.TxResultSuccess}, Ack: ackB.Fn()},
		}

		cons := &ConsumerMock{
			FetchFunc: func(ctx context.Context, count int) ([]msgqueue.Message[types.TxResult], error) {
				if len(batch) == 0 {
					return nil, nil //nolint:nilnil
				}
				out := batch
				batch = nil
				return out, nil
			},
		}

		handler := &HandlerMock{
			HandleFunc: func(ctx context.Context, results []types.TxResult) error {
				require.Len(t, results, 2)
				assert.Equal(t, "id-a", results[0].CorrelationID)
				assert.Equal(t, "id-b", results[1].CorrelationID)
				return nil
			},
		}

		r := resultrouter.New(defaultConfig(), cons)
		r.Register("type.a", handler)

		ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond)
		defer cancel()
		require.NoError(t, r.Run(ctx))

		assert.Len(t, handler.HandleCalls(), 1)
		ackA.AssertCalled(t)
		ackB.AssertCalled(t)
	})

	t.Run("groups by MessageType when fetched batch is mixed", func(t *testing.T) {
		ack1 := spy.NewAck()
		ack2 := spy.NewAck()
		ack3 := spy.NewAck()
		batch := []msgqueue.Message[types.TxResult]{
			{V: types.TxResult{CorrelationID: "1", MessageType: "type.a"}, Ack: ack1.Fn()},
			{V: types.TxResult{CorrelationID: "2", MessageType: "type.b"}, Ack: ack2.Fn()},
			{V: types.TxResult{CorrelationID: "3", MessageType: "type.a"}, Ack: ack3.Fn()},
		}

		cons := &ConsumerMock{
			FetchFunc: func(ctx context.Context, count int) ([]msgqueue.Message[types.TxResult], error) {
				if len(batch) == 0 {
					return nil, nil //nolint:nilnil
				}
				out := batch
				batch = nil
				return out, nil
			},
		}

		var handlerAResults atomic.Value
		var handlerBResults atomic.Value

		handlerA := &HandlerMock{
			HandleFunc: func(ctx context.Context, results []types.TxResult) error {
				ids := make([]string, len(results))
				for i, r := range results {
					ids[i] = r.CorrelationID
				}
				handlerAResults.Store(ids)
				return nil
			},
		}
		handlerB := &HandlerMock{
			HandleFunc: func(ctx context.Context, results []types.TxResult) error {
				ids := make([]string, len(results))
				for i, r := range results {
					ids[i] = r.CorrelationID
				}
				handlerBResults.Store(ids)
				return nil
			},
		}

		r := resultrouter.New(defaultConfig(), cons)
		r.Register("type.a", handlerA)
		r.Register("type.b", handlerB)

		ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond)
		defer cancel()
		require.NoError(t, r.Run(ctx))

		assert.Equal(t, []string{"1", "3"}, handlerAResults.Load())
		assert.Equal(t, []string{"2"}, handlerBResults.Load())
		ack1.AssertCalled(t)
		ack2.AssertCalled(t)
		ack3.AssertCalled(t)
	})

	t.Run("does not ack when handler returns error", func(t *testing.T) {
		ack := spy.NewAck()
		batch := []msgqueue.Message[types.TxResult]{
			{V: types.TxResult{CorrelationID: "id-1", MessageType: "type.a"}, Ack: ack.Fn()},
		}

		cons := &ConsumerMock{
			FetchFunc: func(ctx context.Context, count int) ([]msgqueue.Message[types.TxResult], error) {
				if len(batch) == 0 {
					return nil, nil //nolint:nilnil
				}
				out := batch
				batch = nil
				return out, nil
			},
		}

		handler := &HandlerMock{
			HandleFunc: func(ctx context.Context, results []types.TxResult) error {
				return errors.New("handler failure")
			},
		}

		r := resultrouter.New(defaultConfig(), cons)
		r.Register("type.a", handler)

		ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond)
		defer cancel()
		require.NoError(t, r.Run(ctx))

		assert.Len(t, handler.HandleCalls(), 1)
		ack.AssertNotCalled(t, "should not ack when handler errors")
	})

	t.Run("acks groups independently — failure of one does not block another", func(t *testing.T) {
		ackA := spy.NewAck()
		ackB := spy.NewAck()
		batch := []msgqueue.Message[types.TxResult]{
			{V: types.TxResult{CorrelationID: "a", MessageType: "type.fail"}, Ack: ackA.Fn()},
			{V: types.TxResult{CorrelationID: "b", MessageType: "type.ok"}, Ack: ackB.Fn()},
		}

		cons := &ConsumerMock{
			FetchFunc: func(ctx context.Context, count int) ([]msgqueue.Message[types.TxResult], error) {
				if len(batch) == 0 {
					return nil, nil //nolint:nilnil
				}
				out := batch
				batch = nil
				return out, nil
			},
		}

		failing := &HandlerMock{
			HandleFunc: func(ctx context.Context, results []types.TxResult) error {
				return errors.New("boom")
			},
		}
		ok := &HandlerMock{
			HandleFunc: func(ctx context.Context, results []types.TxResult) error {
				return nil
			},
		}

		r := resultrouter.New(defaultConfig(), cons)
		r.Register("type.fail", failing)
		r.Register("type.ok", ok)

		ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond)
		defer cancel()
		require.NoError(t, r.Run(ctx))

		ackA.AssertNotCalled(t, "failed group should not ack")
		ackB.AssertCalled(t)
	})

	t.Run("acks unknown MessageType to keep queue moving", func(t *testing.T) {
		ack := spy.NewAck()
		batch := []msgqueue.Message[types.TxResult]{
			{V: types.TxResult{CorrelationID: "id-x", MessageType: "type.unknown"}, Ack: ack.Fn()},
		}

		cons := &ConsumerMock{
			FetchFunc: func(ctx context.Context, count int) ([]msgqueue.Message[types.TxResult], error) {
				if len(batch) == 0 {
					return nil, nil //nolint:nilnil
				}
				out := batch
				batch = nil
				return out, nil
			},
		}

		handler := &HandlerMock{
			HandleFunc: func(ctx context.Context, results []types.TxResult) error {
				assert.Fail(t, "registered handler should not be called for unknown MessageType")
				return nil
			},
		}

		r := resultrouter.New(defaultConfig(), cons)
		r.Register("type.something-else", handler)

		ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond)
		defer cancel()
		require.NoError(t, r.Run(ctx))

		assert.Empty(t, handler.HandleCalls())
		ack.AssertCalled(t)
	})

	t.Run("HandlerFunc adapter satisfies Handler", func(t *testing.T) {
		var seen []string
		f := resultrouter.HandlerFunc(func(ctx context.Context, results []types.TxResult) error {
			for _, r := range results {
				seen = append(seen, r.CorrelationID)
			}
			return nil
		})

		err := f.Handle(context.Background(), []types.TxResult{
			{CorrelationID: "x"}, {CorrelationID: "y"},
		})
		require.NoError(t, err)
		assert.Equal(t, []string{"x", "y"}, seen)
	})
}
