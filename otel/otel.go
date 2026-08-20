package otel

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/runtime"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/trace"

	"github.com/raylsnetwork/rayls-sovereign-relayer/withstack"
)

// installOTelErrorHandler routes OTel SDK errors through slog. Transient export
// failures (collector unreachable, timeout, DNS) are dropped entirely so dev
// runs without a collector stay quiet at every log level; other SDK errors
// still surface at warn so configuration problems remain visible.
func installOTelErrorHandler() {
	otel.SetErrorHandler(otel.ErrorHandlerFunc(func(err error) {
		msg := err.Error()
		if strings.Contains(msg, "connection refused") ||
			strings.Contains(msg, "context deadline exceeded") ||
			strings.Contains(msg, "no such host") {
			return
		}
		slog.Warn("OTel SDK error", slog.String("error", msg))
	}))
}

// setupOTelSDK bootstraps the OpenTelemetry pipeline.
// If it does not return an error, make sure to call shutdown for proper cleanup.
func SetupOTelSDK(ctx context.Context, isOtelDisabled bool) (shutdown func(context.Context) error, err error) {
	installOTelErrorHandler()

	var shutdownFuncs []func(context.Context) error

	// shutdown calls cleanup functions registered via shutdownFuncs.
	// The errors from the calls are joined.
	// Each registered cleanup will be invoked once.
	shutdown = func(ctx context.Context) error {
		var shutdownErr error
		for _, fn := range shutdownFuncs {
			shutdownErr = errors.Join(shutdownErr, fn(ctx))
		}
		shutdownFuncs = nil
		return shutdownErr
	}

	// handleErr calls shutdown for cleanup and makes sure that all errors are returned.
	handleErr := func(inErr error) {
		err = errors.Join(inErr, shutdown(ctx))
	}

	var prop propagation.TextMapPropagator
	if isOtelDisabled {
		// no-op Propagator
		prop = propagation.NewCompositeTextMapPropagator()
	} else {
		prop = propagation.NewCompositeTextMapPropagator(
			propagation.TraceContext{},
			propagation.Baggage{},
		)
	}
	otel.SetTextMapPropagator(prop)

	// Setup tracing
	var tracerProvider *trace.TracerProvider
	if isOtelDisabled {
		// No-op tracer provider when OTEL is disabled
		// https://github.com/open-telemetry/opentelemetry-go/blob/0fc383a3ac34cfae998eb660b13689bab8d0804c/trace/noop/noop.go#L35
		tracerProvider = trace.NewTracerProvider()
	} else {
		var traceExporter *otlptrace.Exporter
		traceExporter, err = otlptrace.New(ctx, otlptracehttp.NewClient())
		if err != nil {
			slog.ErrorContext(ctx, "failed to create trace exporter", slog.String("error", err.Error()))
			handleErr(err)
			return nil, withstack.Wrap(fmt.Errorf("creating OTel trace exporter: %w", err))
		}
		tracerProvider = trace.NewTracerProvider(trace.WithBatcher(traceExporter))
	}
	shutdownFuncs = append(shutdownFuncs, tracerProvider.Shutdown)
	otel.SetTracerProvider(tracerProvider)

	// Setup metrics
	var meterProvider *metric.MeterProvider
	if isOtelDisabled {
		// will be no-op because no readers are provided to the MeterProvider
		meterProvider = metric.NewMeterProvider()
	} else {
		var metricExporter *otlpmetrichttp.Exporter
		metricExporter, err = otlpmetrichttp.New(ctx)
		if err != nil {
			slog.ErrorContext(ctx, "failed to create metric exporter", slog.String("error", err.Error()))
			handleErr(err)
			return nil, withstack.Wrap(fmt.Errorf("creating OTel metric exporter: %w", err))
		}

		meterProvider = metric.NewMeterProvider(metric.WithReader(metric.NewPeriodicReader(metricExporter)))
	}
	shutdownFuncs = append(shutdownFuncs, meterProvider.Shutdown)
	otel.SetMeterProvider(meterProvider)

	// FOR REFERENCE ON GENERAL LOGGING. OUR INTEGRATION IS AT logger/slog.go and uses slog.Logger
	// logExporter, err := otlploghttp.New(ctx, otlploghttp.WithInsecure())
	// if err != nil {
	// 	return nil, err
	// }

	// loggerProvider := log.NewLoggerProvider(log.WithProcessor(log.NewBatchProcessor(logExporter)))
	// if err != nil {
	// 	handleErr(err)
	// 	return
	// }
	// shutdownFuncs = append(shutdownFuncs, loggerProvider.Shutdown)
	// global.SetLoggerProvider(loggerProvider)

	if !isOtelDisabled {
		err = runtime.Start(runtime.WithMinimumReadMemStatsInterval(time.Second))
		if err != nil {
			slog.ErrorContext(ctx, "otel runtime instrumentation failed:", slog.Any("error", err))
			handleErr(err)
			return nil, withstack.Wrap(fmt.Errorf("starting OTel runtime instrumentation: %w", err))
		}
	}

	return
}
