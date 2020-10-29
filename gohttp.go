package gohttp

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

var defaultAgent = "Chrome"
var defaultClient = &http.Client{Transport: &http.Transport{Proxy: nil}}

// H represents the key-value pairs in an HTTP header.
type H map[string]string

// SetAgent sets default user agent string.
func SetAgent(agent string) {
	if agent != "" {
		defaultAgent = agent
	}
}

func buildRequest(method, URL string, data interface{}) (*http.Request, error) {
	switch data.(type) {
	case nil:
		return http.NewRequest(method, URL, nil)
	case io.Reader:
		return http.NewRequest(method, URL, data.(io.Reader))
	case url.Values:
		req, err := http.NewRequest(method, URL, strings.NewReader(data.(url.Values).Encode()))
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
	req.Header.Set("User-Agent", defaultAgent)
	for k, v := range header {
		req.Header[k] = v
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
	}
}

// Close closes the response body.
func (r *Response) Close() {
	if r.Error == nil {
		r.Body.Close()
	}
}

// JSON parses the response body as JSON-encoded data
// and stores the result in the value pointed to by data.
func (r *Response) JSON(data interface{}) error {
	if r.Error != nil {
		return r.Error
	}
	defer r.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(body, &data); err != nil {
		return err
	}
	return nil
}

// Bytes returns a slice of byte of the response body.
func (r *Response) Bytes() []byte {
	if r.Error != nil {
		return nil
	}
	defer r.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil
	}
	return body
}

// String returns the contents of the response body as a string.
func (r *Response) String() string {
	return string(r.Bytes())
}
