package gohttp

import (
	"bytes"
	"compress/flate"
	"compress/gzip"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
)

var defaultAgent = "Go-HTTP-Client"
var defaultClient = http.DefaultClient

// H represents the key-value pairs in an HTTP header.
type H map[string]string

func defaultHeaders() H {
	return H{
		"User-Agent":      defaultAgent,
		"Accept-Encoding": "gzip, deflate",
		"Accept":          "*/*",
		"Connection":      "keep-alive",
	}
}

// SetAgent sets default user agent string.
func SetAgent(agent string) {
	if agent != "" {
		defaultAgent = agent
	}
}

// SetProxy sets default client transport proxy.
func SetProxy(proxy string) error {
	proxyURL, err := url.Parse(proxy)
	if err != nil {
		return err
	}

	defaultClient.Transport = &http.Transport{Proxy: http.ProxyURL(proxyURL)}

	return nil
}

// SetNoProxy sets default client use no proxy.
func SetNoProxy() {
	defaultClient.Transport = &http.Transport{Proxy: nil}
}

// SetProxyFromEnvironment sets default client use environment proxy.
func SetProxyFromEnvironment() {
	defaultClient.Transport = &http.Transport{Proxy: http.ProxyFromEnvironment}
}

// SetClient sets default client.
func SetClient(c *http.Client) {
	if c != nil {
		defaultClient = c
	} else {
		panic("cannot set a nil client")
	}
}

func buildRequest(method, URL string, data interface{}) (*http.Request, error) {
	switch data := data.(type) {
	case nil:
		return http.NewRequest(method, URL, nil)

	case io.Reader:
		return http.NewRequest(method, URL, data)

	case url.Values:
		req, err := http.NewRequest(method, URL, strings.NewReader(data.Encode()))
		if err != nil {
			return nil, err
		}
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		return req, nil

	default:
		jsonData, err := json.Marshal(data)
		if err != nil {
			return nil, err
		}

		req, err := http.NewRequest(method, URL, bytes.NewBuffer(jsonData))
		if err != nil {
			return nil, err
		}
		req.Header.Set("Content-Type", "application/json")

		return req, nil
	}
}

func doRequest(method, url string, header http.Header, data interface{}, client *http.Client) *Response {
	req, err := buildRequest(method, url, data)
	if err != nil {
		return &Response{Error: err}
	}

	for k, v := range defaultHeaders() {
		req.Header.Set(k, v)
	}

	for k, v := range header {
		req.Header[k] = v
	}

	if client == nil {
		panic("client pointer is nil")
	}

	return buildResponse(client.Do(req))
}

// Response represents the response from an HTTP request.
type Response struct {
	Error      error
	Body       io.ReadCloser
	StatusCode int
	Header     http.Header
	Cookies    []*http.Cookie
	Request    *http.Request
	Response   *http.Response
	cached     bool
	bytes      []byte
}

func buildResponse(resp *http.Response, err error) *Response {
	if err != nil {
		return &Response{Error: err}
	}

	return &Response{
		Error:      nil,
		Body:       resp.Body,
		StatusCode: resp.StatusCode,
		Header:     resp.Header,
		Cookies:    resp.Cookies(),
		Request:    resp.Request,
		Response:   resp,
	}
}

// Close closes the response body.
func (r *Response) Close() {
	if r.Error == nil {
		r.Body.Close()
	}
}

// Bytes returns a slice of byte of the response body.
func (r *Response) Bytes() []byte {
	if r.Error != nil {
		return nil
	}
	if r.cached {
		return r.bytes
	}
	defer r.Close()

	reader := r.Body
	switch r.Header.Get("Content-Encoding") {
	case "gzip":
		reader, r.Error = gzip.NewReader(reader)
		if r.Error != nil {
			return nil
		}
	case "deflate":
		reader = flate.NewReader(reader)
	}

	r.bytes, r.Error = io.ReadAll(reader)
	if r.Error != nil {
		return nil
	}
	r.cached = true

	return r.bytes
}

// String returns the contents of the response body as a string.
func (r *Response) String() string {
	return string(r.Bytes())
}

// JSON parses the response body as JSON-encoded data
// and stores the result in the value pointed to by data.
func (r *Response) JSON(data interface{}) error {
	if r.Error != nil {
		return r.Error
	}

	return json.Unmarshal(r.Bytes(), &data)
}

// Save saves the response data to file.
func (r *Response) Save(file string) error {
	if r.Error != nil {
		return r.Error
	}

	f, err := os.Create(file)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.Write(r.Bytes())

	return err
}
