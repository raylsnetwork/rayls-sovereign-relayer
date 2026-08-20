package adapters

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

type OTelTracer struct {
	tracer trace.Tracer
}

func NewOTelTracer(serviceName string) *OTelTracer {
	return &OTelTracer{
		tracer: otel.Tracer(serviceName),
	}
}

func (t *OTelTracer) Start(
	ctx context.Context, spanName string, opts ...trace.SpanStartOption,
) (context.Context, trace.Span) {
	return t.tracer.Start(ctx, spanName, opts...)
}
