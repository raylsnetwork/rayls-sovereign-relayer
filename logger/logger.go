package logger

import (
	"context"
	"log/slog"
	"os"

	"github.com/agoda-com/opentelemetry-go/otelslog"
)

var logger *slog.Logger

func InitLogger(handlerType, levelStr string, isOtelDisabled bool) (func(context.Context) error, error) {
	var (
		handler  slog.Handler
		shutdown = voidShutdown
	)

	level := getLogLevel(levelStr)
	if handlerType == "Text" {
		handler = NewColored(os.Stdout, &Options{Level: level})
	} else {
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: level})
	}

	if !isOtelDisabled {
		var (
			otelHandler slog.Handler
			err         error
		)

		otelHandler, shutdown, err = NewOtel(&otelslog.HandlerOptions{Level: level, AddBaggage: true})
		if err != nil {
			return shutdown, err
		}

		handler = &MultiHandler{
			level:    level,
			handlers: []slog.Handler{handler, otelHandler},
		}
	}

	logger = slog.New(handler)
	slog.SetDefault(logger)

	return shutdown, nil
}

func getLogLevel(levelStr string) slog.Level {
	switch levelStr {
	case slog.LevelDebug.String():
		return slog.LevelDebug
	case slog.LevelInfo.String():
		return slog.LevelInfo
	case slog.LevelWarn.String():
		return slog.LevelWarn
	case slog.LevelError.String():
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

func voidShutdown(context.Context) error {
	return nil
}
