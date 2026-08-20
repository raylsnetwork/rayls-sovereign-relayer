package client

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"time"

	"github.com/raylsnetwork/rayls-sovereign-relayer/backoff"
	"github.com/raylsnetwork/rayls-sovereign-relayer/cryptography"
	"github.com/raylsnetwork/rayls-sovereign-relayer/withstack"
)

var ErrExceededRetryCount = errors.New("retry count exceeded")

type HTTPClient interface {
	Do(*http.Request) (*http.Response, error)
}

// SimpleClient wraps a standard *http.Client and implements HTTPClient,
// ensuring errors from Do() are wrapped with stack traces at the external call origin.
type SimpleClient struct {
	client *http.Client
}

func NewSimpleClient(client *http.Client) *SimpleClient {
	return &SimpleClient{client: client}
}

func (c *SimpleClient) Do(req *http.Request) (*http.Response, error) {
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, withstack.Wrap(fmt.Errorf("executing HTTP request: %w", err))
	}
	return resp, nil
}

type RetryClient struct {
	retryCount int
	backoff    backoff.Strategy
	client     HTTPClient
}

func NewRetryClient(retryCount int, backoffDuration time.Duration, client HTTPClient) *RetryClient {
	// Use linear backoff matching the previous cumulative delay pattern
	lin, _ := backoff.NewLinear(backoffDuration, backoffDuration, 0)

	return &RetryClient{
		retryCount: retryCount,
		backoff:    lin,
		client:     client,
	}
}

func (c *RetryClient) Do(req *http.Request) (*http.Response, error) {
	ctx := req.Context() //nolint:contextcheck // context derived from HTTP request
	if ctx == nil {
		ctx = context.Background()
	}

	var resp *http.Response

	err := c.backoff.Do(ctx, c.retryCount, func() error {
		return c.attemptRequest(req, &resp)
	})
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrExceededRetryCount, err)
	}

	return resp, nil
}

// attemptRequest executes a single request attempt with proper cloning and error handling
func (c *RetryClient) attemptRequest(req *http.Request, respPtr **http.Response) error {
	reqCpy := cloneRequest(req)

	r, err := c.client.Do(reqCpy)
	if err == nil && !isStatusRetriable(r.StatusCode) {
		*respPtr = r
		return nil
	}

	// Log the retry reason
	logRetryReason(err, r)

	if err != nil {
		return withstack.Wrap(fmt.Errorf("executing HTTP request: %w", err))
	}
	// Retriable status code
	return withstack.Wrap(fmt.Errorf("retriable HTTP status: %s", r.Status))
}

// cloneRequest creates a copy of the HTTP request for retry attempts
func cloneRequest(req *http.Request) *http.Request {
	reqCpy := new(http.Request)
	*reqCpy = *req

	if req.Body != nil {
		body, err := io.ReadAll(req.Body)
		if err != nil {
			// If we can't read the body, return original request
			return req
		}
		bodyCpy := bytes.Clone(body)

		req.Body = io.NopCloser(bytes.NewReader(body))
		reqCpy.Body = io.NopCloser(bytes.NewBuffer(bodyCpy))
	}

	return reqCpy
}

func logRetryReason(err error, resp *http.Response) {
	if err != nil {
		slog.Error("Request failed, retrying...",
			slog.Any("error", err),
			slog.Any("stack trace", withstack.Wrap(err)))
	} else if resp != nil {
		slog.Error("Request failed, retrying...",
			slog.String("status code", resp.Status))
	}
}

func isStatusRetriable(status int) bool {
	switch status {
	case http.StatusRequestTimeout:
		fallthrough
	case http.StatusInternalServerError:
		fallthrough
	case http.StatusBadGateway:
		fallthrough
	case http.StatusGatewayTimeout:
		fallthrough
	case http.StatusServiceUnavailable:
		return true
	}

	return false
}

func NewEncryptedAPIKeyClient(key, secret string, client HTTPClient) HTTPClient {
	apiKeyCli := NewAPIKeyClient(key, client)
	return NewAESGCMHTTPClient(secret, apiKeyCli)
}

type APIKeyClient struct {
	key string

	client HTTPClient
}

func NewAPIKeyClient(key string, client HTTPClient) *APIKeyClient {
	return &APIKeyClient{
		key:    key,
		client: client,
	}
}

func (c *APIKeyClient) Do(req *http.Request) (*http.Response, error) {
	req.Header.Set("Authorization", c.key)

	return c.client.Do(req)
}

type AESGCMHTTPClient struct {
	secret string

	client HTTPClient
}

func NewAESGCMHTTPClient(secret string, client HTTPClient) *AESGCMHTTPClient {
	return &AESGCMHTTPClient{
		secret: secret,
		client: client,
	}
}

func (c *AESGCMHTTPClient) Do(req *http.Request) (*http.Response, error) {
	req.Header.Set("Content-Type", "application/octet-stream")

	// Encrypt request
	var encryptedRequestContents []byte
	if req.Body != nil {
		requestContents, err := io.ReadAll(req.Body)
		if err != nil {
			return nil, withstack.Wrap(fmt.Errorf("failed to read request body: %w", err))
		}

		encryptedRequestContents, err = cryptography.EncryptData(c.secret, requestContents)
		if err != nil {
			return nil, withstack.Wrap(fmt.Errorf("failed to encrypt request data: %w", err))
		}
		updateRequestBody(req, encryptedRequestContents)
	}

	// Call client Do method
	response, err := c.client.Do(req)
	if err != nil {
		return response, withstack.Wrap(fmt.Errorf("executing encrypted HTTP request: %w", err))
	}

	// Decrypt response
	var responseContents []byte
	if response.Body != nil {
		var encryptedResponseContents []byte
		encryptedResponseContents, err = io.ReadAll(response.Body)
		if err != nil {
			return nil, withstack.Wrap(fmt.Errorf("failed to read response body: %w", err))
		}

		responseContents, err = cryptography.DecryptData(c.secret, encryptedResponseContents)
		if err != nil {
			return nil, withstack.Wrap(fmt.Errorf("failed to decrypt response data: %w", err))
		}

		updateResponseBody(response, responseContents)
	}

	return response, nil
}

func updateRequestBody(req *http.Request, body []byte) {
	bodyReaderCloser := io.NopCloser(bytes.NewReader(body))

	req.Body = bodyReaderCloser
	req.ContentLength = int64(len(body))
	req.GetBody = func() (io.ReadCloser, error) {
		r := bytes.NewReader(body)
		return io.NopCloser(r), nil
	}
}

func updateResponseBody(resp *http.Response, body []byte) {
	bodyReaderCloser := io.NopCloser(bytes.NewReader(body))

	resp.Body = bodyReaderCloser
	resp.ContentLength = int64(len(body))
}
