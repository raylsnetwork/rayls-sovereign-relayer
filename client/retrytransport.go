package client

import (
	"log/slog"
	"net/http"
	"reflect"
	"time"

	"github.com/raylsnetwork/rayls-privacy-relayer-api/withstack"
)

type RetryTransport struct {
	http.RoundTripper
	InitialDelay   time.Duration
	DelayIncrement time.Duration
}

func (t *RetryTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	delay := t.InitialDelay

	for {
		resp, err := t.RoundTripper.RoundTrip(req)
		if err == nil {
			return resp, nil
		}

		delay += t.DelayIncrement
		slog.Warn(
			"Failed operation... retrying after delay",
			slog.String("delay", delay.String()),
			slog.Any("err", withstack.WrapWithDepth(err, 10)),
			slog.Any("type", reflect.TypeOf(err)),
		)

		select {
		case <-req.Context().Done():
			return nil, req.Context().Err()
		case <-time.After(delay):
		}
	}
}
