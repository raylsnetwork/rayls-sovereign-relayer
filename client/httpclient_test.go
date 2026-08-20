package client_test

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/raylsnetwork/rayls-privacy-relayer-api/client"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/cryptography"
	"github.com/raylsnetwork/rayls-privacy-relayer-api/testtools"
	"github.com/stretchr/testify/require"
)

type StubClient struct {
	spyRequest *http.Request

	responseBody []byte
}

func (s *StubClient) Do(req *http.Request) (*http.Response, error) {
	s.spyRequest = req

	response := &http.Response{
		Body: io.NopCloser(bytes.NewReader(s.responseBody)),
	}
	return response, nil
}

type StubClientWithFailCount struct {
	failCount int

	spyRequest   *http.Request
	spyCallCount int
}

func (s *StubClientWithFailCount) Do(req *http.Request) (*http.Response, error) {
	s.spyRequest = req

	s.spyCallCount++
	if s.spyCallCount <= s.failCount {
		return nil, errors.New("stub failure: simulated error")
	}

	resp := &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(strings.NewReader("success")),
		Header:     make(http.Header),
	}
	return resp, nil
}

type StubClientWithStatusCodeAndRetryCheck struct {
	statusCode int

	hasCalled  bool
	hasRetried bool
}

func (s *StubClientWithStatusCodeAndRetryCheck) Do(req *http.Request) (*http.Response, error) {
	if s.hasCalled {
		s.hasRetried = true
	}
	s.hasCalled = true

	resp := &http.Response{
		StatusCode: s.statusCode,
		Body:       io.NopCloser(strings.NewReader("success")),
		Header:     make(http.Header),
	}
	return resp, nil
}

func TestRetryClient(t *testing.T) {
	testtools.SilenceLogger()

	t.Run("forwards request to wrapped Do method", func(t *testing.T) {
		httpCli := &StubClient{}
		retryCli := client.NewRetryClient(1, time.Microsecond, httpCli)

		wantReq, _ := http.NewRequest(http.MethodPost, "example.url", nil)

		_, err := retryCli.Do(wantReq)
		require.Nil(t, err)

		require.Equal(t, wantReq, httpCli.spyRequest)
	})

	t.Run("returns ErrExceededRetryCount on exceeded retry count", func(t *testing.T) {
		wantRetryCount := 10

		httpCli := &StubClientWithFailCount{
			failCount: wantRetryCount,
		}
		retryCli := client.NewRetryClient(wantRetryCount, time.Microsecond, httpCli)

		wantReq, _ := http.NewRequest(http.MethodPost, "example.url", nil)

		_, err := retryCli.Do(wantReq)

		require.ErrorIs(t, err, client.ErrExceededRetryCount)
		require.Equal(t, wantRetryCount, httpCli.spyCallCount)
	})

	t.Run("retries request until success", func(t *testing.T) {
		wantRetryCount := 10

		httpCli := &StubClientWithFailCount{
			failCount: wantRetryCount - 1,
		}
		retryCli := client.NewRetryClient(wantRetryCount, time.Microsecond, httpCli)

		wantReq, _ := http.NewRequest(http.MethodPost, "example.url", nil)

		resp, err := retryCli.Do(wantReq)
		require.Nil(t, err)

		require.Equal(t, http.StatusOK, resp.StatusCode)
	})
}

func TestRetryClient_RetriableStatusCodes(t *testing.T) {
	testtools.SilenceLogger()
	cases := []struct {
		StatusCode  int
		ShouldRetry bool
	}{
		{
			StatusCode:  http.StatusRequestTimeout,
			ShouldRetry: true,
		},
		{
			StatusCode:  http.StatusInternalServerError,
			ShouldRetry: true,
		},
		{
			StatusCode:  http.StatusBadGateway,
			ShouldRetry: true,
		},
		{
			StatusCode:  http.StatusServiceUnavailable,
			ShouldRetry: true,
		},
		{
			StatusCode:  http.StatusGatewayTimeout,
			ShouldRetry: true,
		},
		{
			StatusCode:  http.StatusOK,
			ShouldRetry: false,
		},
		{
			StatusCode:  http.StatusNotFound,
			ShouldRetry: false,
		},
	}

	for _, c := range cases {
		t.Run(fmt.Sprintf("doesn't retry on status code %d", c.StatusCode), func(t *testing.T) {
			httpCli := &StubClientWithStatusCodeAndRetryCheck{
				statusCode: c.StatusCode,
			}
			retryCli := client.NewRetryClient(2, time.Microsecond, httpCli)
			wantReq, _ := http.NewRequest(http.MethodPost, "example.url", nil)

			_, err := retryCli.Do(wantReq)

			if c.ShouldRetry {
				require.ErrorIs(t, err, client.ErrExceededRetryCount)
				require.True(t, httpCli.hasRetried)
			} else {
				require.Nil(t, err)
				require.False(t, c.ShouldRetry)
			}
		})
	}
}

func TestAESGCMClient(t *testing.T) {
	t.Run("encrypts request body and sets content type", func(t *testing.T) {
		wantContentType := "application/octet-stream"
		wantRequestContents := "example body"

		secret := "example secret"
		httpCli := &StubClient{}
		secureCli := client.NewAESGCMHTTPClient(secret, httpCli)

		req, _ := http.NewRequest(http.MethodPost, "example.url", bytes.NewBuffer([]byte(wantRequestContents)))

		// Response is not set in the StubClient so the function will fail on the decrypt method.
		// We don't care about it here tho, as we are testing only the encryption flow.
		_, _ = secureCli.Do(req)

		gotContentType := httpCli.spyRequest.Header.Get("Content-Type")
		require.Equal(t, wantContentType, gotContentType)

		// Test that the GetBody() method is not nil.
		require.NotNil(t, req.GetBody)

		// Check that we replaced the buffer returned from GetBody() with the
		// same one with the same contents as the Body field,
		getBodyMethodReader, err := req.GetBody()
		require.Nil(t, err)

		getBodyMethodContents, err := io.ReadAll(getBodyMethodReader)
		require.Nil(t, err)

		bodyContents, err := io.ReadAll(req.Body)
		require.Nil(t, err)

		require.Equal(t, bodyContents, getBodyMethodContents)

		// Check that the ContentLength field is updated
		require.Equal(t, len(bodyContents), int(req.ContentLength))

		// Check that the request is correctly encrypted
		gotRequestContentsBytes, err := cryptography.DecryptData(secret, bodyContents)
		require.Nil(t, err)

		require.Equal(t, wantRequestContents, string(gotRequestContentsBytes))
	})

	t.Run("decrypts response body and overwrites it in response struct", func(t *testing.T) {
		wantResponseContents := "example body"

		secret := "example secret"
		encrResponseContents, err := cryptography.EncryptData(secret, []byte(wantResponseContents))
		require.Nil(t, err)

		httpCli := &StubClient{
			responseBody: []byte(encrResponseContents),
		}
		secureCli := client.NewAESGCMHTTPClient(secret, httpCli)

		req, _ := http.NewRequest(http.MethodPost, "example.url", bytes.NewBuffer([]byte("example request contents")))

		response, err := secureCli.Do(req)

		require.Nil(t, err)

		gotResponseContentsBytes, err := io.ReadAll(response.Body)
		require.Nil(t, err)

		// Check that the ContentLength field is updated
		require.Equal(t, len(gotResponseContentsBytes), int(response.ContentLength))

		// Check that the response is correctly decrypted
		require.Equal(t, wantResponseContents, string(gotResponseContentsBytes))
	})
}

func TestAPIKeyClient(t *testing.T) {
	t.Run("writes API Key to request header", func(t *testing.T) {
		wantKey := "example key"

		req, _ := http.NewRequest(http.MethodPost, "example.url", bytes.NewBuffer([]byte("example request contents")))

		httpCli := &StubClient{}
		secureCli := client.NewAPIKeyClient(wantKey, httpCli)

		_, err := secureCli.Do(req)

		require.Nil(t, err)

		gotKey := httpCli.spyRequest.Header.Get("Authorization")
		require.Equal(t, wantKey, gotKey)
	})
}
