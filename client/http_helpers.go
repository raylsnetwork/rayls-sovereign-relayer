package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/raylsnetwork/rayls-sovereign-relayer/withstack"
)

func buildRequestURL(basePath, endpoint string, params url.Values) string {
	requestURL, _ := url.Parse(basePath)
	requestURL.Path = endpoint
	requestURL.RawQuery = params.Encode()

	return requestURL.String()
}

func newJSONRequest(method string, url string, object interface{}) (*http.Request, error) {
	var body io.Reader

	if object != nil {
		jsonBytes, err := json.Marshal(object)
		if err != nil {
			return nil, withstack.Wrap(fmt.Errorf("failed to marshal request: %w", err))
		}
		body = bytes.NewBuffer(jsonBytes)
	}

	request, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, withstack.Wrap(fmt.Errorf("failed to create HTTP request: %w", err))
	}

	if object != nil {
		request.Header.Set("Content-Type", "application/json")
	}

	return request, nil
}

func handleInternalServerError(r *http.Response) error {
	defer func() { _ = r.Body.Close() }()
	var errResp struct {
		Error string `json:"error"`
	}
	if err := json.NewDecoder(r.Body).Decode(&errResp); err != nil {
		return fmt.Errorf("decoding internal server error response: %w", err)
	}
	return NewProofAPIError(errResp.Error)
}

func handleUnsupportedStatus(r *http.Response) error {
	defer func() { _ = r.Body.Close() }()
	bodyContent, _ := io.ReadAll(r.Body)
	return NewProofAPIError(string(bodyContent))
}
