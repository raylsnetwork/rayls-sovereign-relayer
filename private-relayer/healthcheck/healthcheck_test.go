package healthcheck_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/raylsnetwork/rayls-sovereign-relayer/private-relayer/healthcheck"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newHealthyMocks() (*DatabaseClientMock, *EthereumClientMock, *EthereumClientMock, *NATSClientMock) {
	dbMock := &DatabaseClientMock{
		PingFunc: func(ctx context.Context) error {
			return nil
		},
	}
	hubMock := &EthereumClientMock{
		BlockNumberFunc: func(ctx context.Context) (uint64, error) {
			return 100, nil
		},
	}
	nodeMock := &EthereumClientMock{
		BlockNumberFunc: func(ctx context.Context) (uint64, error) {
			return 200, nil
		},
	}
	natsMock := &NATSClientMock{
		IsConnectedFunc: func() bool {
			return true
		},
	}
	return dbMock, hubMock, nodeMock, natsMock
}

func newHealthcheck(
	natsMock *NATSClientMock,
	dbMock *DatabaseClientMock,
	hubMock *EthereumClientMock,
	nodeMock *EthereumClientMock,
) *healthcheck.Healthcheck {
	return healthcheck.New(
		healthcheck.Config{Path: "/healthcheck", Addr: "0"},
		natsMock, dbMock, hubMock, nodeMock,
	)
}

func TestServeHTTP(t *testing.T) {
	t.Run("returns 405 for POST request", func(t *testing.T) {
		dbMock, hubMock, nodeMock, natsMock := newHealthyMocks()
		h := newHealthcheck(natsMock, dbMock, hubMock, nodeMock)

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/healthcheck", nil)
		h.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusMethodNotAllowed, rec.Code)
		assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))
	})

	t.Run("returns 405 for PUT request", func(t *testing.T) {
		dbMock, hubMock, nodeMock, natsMock := newHealthyMocks()
		h := newHealthcheck(natsMock, dbMock, hubMock, nodeMock)

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPut, "/healthcheck", nil)
		h.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusMethodNotAllowed, rec.Code)
	})

	t.Run("returns 405 for DELETE request", func(t *testing.T) {
		dbMock, hubMock, nodeMock, natsMock := newHealthyMocks()
		h := newHealthcheck(natsMock, dbMock, hubMock, nodeMock)

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodDelete, "/healthcheck", nil)
		h.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusMethodNotAllowed, rec.Code)
	})

	t.Run("returns 200 with healthy status for GET", func(t *testing.T) {
		dbMock, hubMock, nodeMock, natsMock := newHealthyMocks()
		h := newHealthcheck(natsMock, dbMock, hubMock, nodeMock)

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/healthcheck", nil)
		h.ServeHTTP(rec, req)

		require.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))

		var body healthcheck.HealthResponse
		err := json.NewDecoder(rec.Body).Decode(&body)
		require.NoError(t, err)
		assert.Equal(t, "healthy", body.Status)
		assert.NotEmpty(t, body.Timestamp)
	})

	t.Run("returns 200 for HEAD request", func(t *testing.T) {
		dbMock, hubMock, nodeMock, natsMock := newHealthyMocks()
		h := newHealthcheck(natsMock, dbMock, hubMock, nodeMock)

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodHead, "/healthcheck", nil)
		h.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))
	})

	t.Run("returns 503 when database ping fails", func(t *testing.T) {
		dbMock, hubMock, nodeMock, natsMock := newHealthyMocks()
		dbMock.PingFunc = func(ctx context.Context) error {
			return errors.New("connection refused")
		}
		h := newHealthcheck(natsMock, dbMock, hubMock, nodeMock)

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/healthcheck", nil)
		h.ServeHTTP(rec, req)

		require.Equal(t, http.StatusServiceUnavailable, rec.Code)

		var body healthcheck.ErrorResponse
		err := json.NewDecoder(rec.Body).Decode(&body)
		require.NoError(t, err)
		assert.Equal(t, "unhealthy", body.Status)
		assert.Contains(t, body.Error, "failed to ping database")
	})

	t.Run("returns 503 when private hub BlockNumber fails", func(t *testing.T) {
		dbMock, hubMock, nodeMock, natsMock := newHealthyMocks()
		hubMock.BlockNumberFunc = func(ctx context.Context) (uint64, error) {
			return 0, errors.New("hub unreachable")
		}
		h := newHealthcheck(natsMock, dbMock, hubMock, nodeMock)

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/healthcheck", nil)
		h.ServeHTTP(rec, req)

		require.Equal(t, http.StatusServiceUnavailable, rec.Code)

		var body healthcheck.ErrorResponse
		err := json.NewDecoder(rec.Body).Decode(&body)
		require.NoError(t, err)
		assert.Equal(t, "unhealthy", body.Status)
		assert.Contains(t, body.Error, "failed to ping private hub client")
	})

	t.Run("returns 503 when private node BlockNumber fails", func(t *testing.T) {
		dbMock, hubMock, nodeMock, natsMock := newHealthyMocks()
		nodeMock.BlockNumberFunc = func(ctx context.Context) (uint64, error) {
			return 0, errors.New("node unreachable")
		}
		h := newHealthcheck(natsMock, dbMock, hubMock, nodeMock)

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/healthcheck", nil)
		h.ServeHTTP(rec, req)

		require.Equal(t, http.StatusServiceUnavailable, rec.Code)

		var body healthcheck.ErrorResponse
		err := json.NewDecoder(rec.Body).Decode(&body)
		require.NoError(t, err)
		assert.Equal(t, "unhealthy", body.Status)
		assert.Contains(t, body.Error, "failed to ping private node client")
	})

	t.Run("returns 503 when NATS is not connected", func(t *testing.T) {
		dbMock, hubMock, nodeMock, natsMock := newHealthyMocks()
		natsMock.IsConnectedFunc = func() bool {
			return false
		}
		h := newHealthcheck(natsMock, dbMock, hubMock, nodeMock)

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/healthcheck", nil)
		h.ServeHTTP(rec, req)

		require.Equal(t, http.StatusServiceUnavailable, rec.Code)

		var body healthcheck.ErrorResponse
		err := json.NewDecoder(rec.Body).Decode(&body)
		require.NoError(t, err)
		assert.Equal(t, "unhealthy", body.Status)
		assert.Contains(t, body.Error, "NATS connection is not active")
	})

	t.Run("short-circuits on first dependency failure", func(t *testing.T) {
		dbMock, hubMock, nodeMock, natsMock := newHealthyMocks()
		dbMock.PingFunc = func(ctx context.Context) error {
			return errors.New("db down")
		}
		h := newHealthcheck(natsMock, dbMock, hubMock, nodeMock)

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/healthcheck", nil)
		h.ServeHTTP(rec, req)

		require.Equal(t, http.StatusServiceUnavailable, rec.Code)

		assert.Len(t, hubMock.BlockNumberCalls(), 0)
		assert.Len(t, nodeMock.BlockNumberCalls(), 0)
		assert.Len(t, natsMock.IsConnectedCalls(), 0)
	})
}
