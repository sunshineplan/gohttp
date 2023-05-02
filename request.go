package gohttp

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strings"
)

func buildRequest(ctx context.Context, method, reqURL string, data any) (*http.Request, error) {
	var body io.Reader
	var contentType string

	switch data := data.(type) {
	case nil:
	case io.Reader:
		body = data
	case url.Values:
		body = strings.NewReader(data.Encode())
		contentType = "application/x-www-form-urlencoded"
	default:
		b, err := json.Marshal(data)
		if err != nil {
			return nil, err
		}
		body = bytes.NewBuffer(b)
		contentType = "application/json"
	}

	req, err := http.NewRequestWithContext(ctx, method, reqURL, body)
	if err != nil {
		return nil, err
	}
	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}

	return req, nil
}

func doRequest(ctx context.Context, method, url string, header http.Header, data any, client *http.Client) *Response {
	req, err := buildRequest(ctx, method, url, data)
	if err != nil {
		return &Response{Response: new(http.Response), Error: err}
	}

	for k, v := range defaultHeaders() {
		req.Header.Set(k, v)
	}

	for k, v := range header {
		req.Header[k] = v
	}

	return buildResponse(client.Do(req))
}
