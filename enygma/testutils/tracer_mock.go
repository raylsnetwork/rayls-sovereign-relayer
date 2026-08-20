package testutils

import (
	"context"
	"sync"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type MockSpan struct {
	trace.Span // Embed to get any unexported methods
	EndCalled  bool
	lock       sync.Mutex
}

func (m *MockSpan) End(options ...trace.SpanEndOption) {
	m.lock.Lock()
	defer m.lock.Unlock()
	m.EndCalled = true
}

func (m *MockSpan) AddEvent(name string, options ...trace.EventOption)  {}
func (m *MockSpan) AddLink(link trace.Link)                             {}
func (m *MockSpan) IsRecording() bool                                   { return true }
func (m *MockSpan) RecordError(err error, options ...trace.EventOption) {}
func (m *MockSpan) SpanContext() trace.SpanContext                      { return trace.SpanContext{} }
func (m *MockSpan) SetAttributes(kv ...attribute.KeyValue)              {}
func (m *MockSpan) SetStatus(code codes.Code, description string)       {}
func (m *MockSpan) SetName(name string)                                 {}
func (m *MockSpan) TracerProvider() trace.TracerProvider                { return nil }

type MockTracer struct {
	StartFunc  func(ctx context.Context, spanName string, opts ...trace.SpanStartOption) (context.Context, trace.Span)
	StartCalls struct {
		sync.RWMutex
		Calls []struct {
			Ctx      context.Context
			SpanName string
			Opts     []trace.SpanStartOption
		}
	}
	MockSpans []*MockSpan
}

func (m *MockTracer) Start(
	ctx context.Context,
	spanName string,
	opts ...trace.SpanStartOption,
) (context.Context, trace.Span) {
	m.StartCalls.Lock()
	defer m.StartCalls.Unlock()
	m.StartCalls.Calls = append(m.StartCalls.Calls, struct {
		Ctx      context.Context
		SpanName string
		Opts     []trace.SpanStartOption
	}{
		Ctx:      ctx,
		SpanName: spanName,
		Opts:     opts,
	})

	mockSpan := &MockSpan{}
	m.MockSpans = append(m.MockSpans, mockSpan)

	if m.StartFunc != nil {
		return m.StartFunc(ctx, spanName, opts...)
	}
	return ctx, mockSpan
}
