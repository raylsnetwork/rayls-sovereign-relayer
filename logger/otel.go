package logger

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/agoda-com/opentelemetry-go/otelslog"
	autosdk "github.com/agoda-com/opentelemetry-logs-go/autoconfigure/sdk/logs"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
)

func NewOtel(opts *otelslog.HandlerOptions) (slog.Handler, func(ctx context.Context) error, error) {
	resrc, err := newResource()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create OTeL resource: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	loggerProvider := autosdk.NewLoggerProvider(ctx, autosdk.WithResource(resrc))
	otelHandler := otelslog.NewOtelHandler(loggerProvider, opts)

	return otelHandler, loggerProvider.Shutdown, nil
}

// configure common attributes for all logs
func newResource() (*resource.Resource, error) {
	hostName, err := os.Hostname()
	if err != nil {
		return nil, fmt.Errorf("failed to get hostname: %w", err)
	}

	return resource.NewWithAttributes(
		semconv.SchemaURL,
		// semconv.ServiceName("passed by env var OTEL_SERVICE_NAME"),
		semconv.ServiceVersion("2.6"),
		semconv.HostName(hostName),
	), nil
}
