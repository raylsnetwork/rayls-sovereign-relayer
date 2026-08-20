package txbatchclient

import (
	"encoding/json"
	"errors"
	"net"
	"net/url"
)

func isConnectionError(err error) bool {
	var netErr net.Error
	if errors.As(err, &netErr) {
		return true // It's a net.Error
	}

	var urlErr *url.Error
	return errors.As(err, &urlErr)
}

func isUnmarshalTypeError(err error) bool {
	var unmarshalErr *json.UnmarshalTypeError
	return errors.As(err, &unmarshalErr)
}
