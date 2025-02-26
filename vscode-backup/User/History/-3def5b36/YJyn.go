package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/google/go-querystring/query"
)

// NewJSONRequest returns an *http.Request with a json encoded body
func NewJSONRequest(ctx context.Context, method string, url string, reqObj interface{}) (*http.Request, error) {
	var body io.Reader
	if reqObj != nil {
		var buf bytes.Buffer
		if err := json.NewEncoder(&buf).Encode(reqObj); err != nil {
			return nil, fmt.Errorf("error encoding data: %s", err)
		}
		body = &buf
	}

	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, err
	}
	req.Header = http.Header{"Content-Type": []string{"application/json"}}
	return req, nil
}

// NewURLEncodeRequest returns an *http.Request with a url encoded body
func NewURLEncodeRequest(
	ctx context.Context, method string, url string, reqObj interface{}) (*http.Request, error) {
	qs, err := query.Values(reqObj)
	if err != nil {
		return nil, fmt.Errorf("error encoding data: %s", err)
	}
	var body io.Reader
	if method == http.MethodGet {
		url = url + "?" + qs.Encode()
	} else {
		body = strings.NewReader(qs.Encode())
	}
	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, err
	}
	req.Header = http.Header{"Content-Type": []string{"application/x-www-form-urlencoded"}}
	return req, nil
}
