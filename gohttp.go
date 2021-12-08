package gohttp

import (
	"bytes"
	"compress/flate"
	"compress/gzip"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	urlpkg "net/url"
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

func setProxy(fn func(*http.Request) (*url.URL, error)) {
	var tr *http.Transport
	var ok bool
	if defaultClient.Transport == nil {
		if tr, ok = http.DefaultTransport.(*http.Transport); ok {
			tr.Proxy = fn
		}
	} else {
		if tr, ok = defaultClient.Transport.(*http.Transport); ok {
			tr.Proxy = fn
		}
	}
	if !ok {
		panic("Transport is not *http.Transport type")
	}
}

// SetProxy sets default client transport proxy.
func SetProxy(proxy string) error {
	proxyURL, err := url.Parse(proxy)
	if err != nil {
		return err
	}

	setProxy(http.ProxyURL(proxyURL))

	return nil
}

// SetNoProxy sets default client use no proxy.
func SetNoProxy() {
	setProxy(nil)
}

// SetProxyFromEnvironment sets default client use environment proxy.
func SetProxyFromEnvironment() {
	setProxy(http.ProxyFromEnvironment)
}

// SetClient sets default client.
func SetClient(c *http.Client) {
	if c != nil {
		defaultClient = c
	} else {
		panic("cannot set a nil client")
	}
}

func buildHeader(headers H) http.Header {
	h := make(http.Header)
	for k, v := range headers {
		h.Set(k, v)
	}
	return h
}

func buildRequest(method, url string, data interface{}) (*http.Request, error) {
	var body io.Reader
	var contentType string

	switch data := data.(type) {
	case nil:
	case io.Reader:
		body = data
	case urlpkg.Values:
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

	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}

	return req, nil
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
	*http.Response
	Error  error
	cached bool
	bytes  []byte
}

func buildResponse(resp *http.Response, err error) *Response {
	if err != nil {
		return &Response{Error: err}
	}

	return &Response{Response: resp}
}

// Close closes the response body.
func (r *Response) Close() error {
	if r.Error == nil {
		return r.Body.Close()
	}
	return nil
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
