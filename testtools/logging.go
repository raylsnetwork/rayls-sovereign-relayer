package testtools

import (
	"io"
	"log/slog"
)

func SilenceLogger() {
	slog.SetDefault(slog.New(slog.NewTextHandler(io.Discard, nil)))
}
