package gohttp

import (
	"net/http"
	"time"
)

var (
	defaultAgent   = "Go-HTTP-Client"
	defaultSession = newSession(http.DefaultClient)
)

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
	return defaultSession.SetProxy(proxy)
}

// SetNoProxy sets default client use no proxy.
func SetNoProxy() {
	defaultSession.SetNoProxy()
}

// SetProxyFromEnvironment sets default client use environment proxy.
func SetProxyFromEnvironment() {
	defaultSession.SetProxyFromEnvironment()
}

// SetTimeout sets default timeout. Zero means no timeout.
func SetTimeout(d time.Duration) {
	defaultSession.SetTimeout(d)
}

// SetClient sets default client.
func SetClient(c *http.Client) {
	defaultSession.SetClient(c)
}

// Get issues a GET to the specified URL with headers.
func Get(url string, headers H) *Response {
	return defaultSession.Get(url, headers)
}

// Head issues a HEAD to the specified URL with headers.
func Head(url string, headers H) *Response {
	return defaultSession.Head(url, headers)
}

// Post issues a POST to the specified URL with headers.
// Post data should be one of nil, io.Reader, url.Values, string map or struct.
func Post(url string, headers H, data any) *Response {
	return defaultSession.Post(url, headers, data)
}

// Upload issues a POST to the specified URL with a multipart document.
func Upload(url string, headers H, params map[string]string, files ...*File) *Response {
	return defaultSession.Upload(url, headers, params, files...)
}

// Get issues a session GET to the specified URL with additional headers.
func (s *Session) Get(url string, headers H) *Response {
	h := s.Header.Clone()
	for k, v := range headers {
		h.Set(k, v)
	}
	return doRequest("GET", url, h, nil, s.Client)
}

// Head issues a session HEAD to the specified URL with additional headers.
func (s *Session) Head(url string, headers H) *Response {
	h := s.Header.Clone()
	for k, v := range headers {
		h.Set(k, v)
	}
	return doRequest("HEAD", url, h, nil, s.Client)
}

// Post issues a session POST to the specified URL with additional headers.
func (s *Session) Post(url string, headers H, data any) *Response {
	h := s.Header.Clone()
	for k, v := range headers {
		h.Set(k, v)
	}
	return doRequest("POST", url, h, data, s.Client)
}

// Upload issues a session POST to the specified URL with a multipart document and additional headers.
func (s *Session) Upload(url string, headers H, params map[string]string, files ...*File) *Response {
	data, contentType, err := buildMultipart(params, files...)
	if err != nil {
		return &Response{Response: new(http.Response), Error: err}
	}
	h := s.Header.Clone()
	h.Set("Content-Type", contentType)
	for k, v := range headers {
		h.Set(k, v)
	}
	return doRequest("POST", url, h, data, s.Client)
}

// KeepAlive repeatedly calls fn with a fixed interval delay between each call.
func (s *Session) KeepAlive(interval *time.Duration, fn func(*Session) error) (err error) {
	for ; err == nil; <-time.After(*interval) {
		err = fn(s)
	}
	return
}

// GetWithClient issues a GET to the specified URL with headers and client.
func GetWithClient(url string, headers H, client *http.Client) *Response {
	return newSession(client).Get(url, headers)
}

// HeadWithClient issues a HEAD to the specified URL with headers and client.
func HeadWithClient(url string, headers H, client *http.Client) *Response {
	return newSession(client).Head(url, headers)
}

// PostWithClient issues a POST to the specified URL with headers and client.
func PostWithClient(url string, headers H, data any, client *http.Client) *Response {
	return newSession(client).Post(url, headers, data)
}

// UploadWithClient issues a POST to the specified URL with a multipart document and client.
func UploadWithClient(url string, headers H, params map[string]string, files []*File, client *http.Client) *Response {
	return newSession(client).Upload(url, headers, params, files...)
}
