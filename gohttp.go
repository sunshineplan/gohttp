package gohttp

import (
	"context"
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

// Do sends an HTTP request and returns a response.
func Do(req *http.Request) (*Response, error) {
	return defaultSession.Do(req)
}

// Get issues a GET to the specified URL with headers.
func Get(url string, headers H) (*Response, error) {
	return defaultSession.Get(url, headers)
}

// GetWithContext issues a GET to the specified URL with context and headers.
func GetWithContext(ctx context.Context, url string, headers H) (*Response, error) {
	return defaultSession.GetWithContext(ctx, url, headers)
}

// Head issues a HEAD to the specified URL with headers.
func Head(url string, headers H) (*Response, error) {
	return defaultSession.Head(url, headers)
}

// HeadWithContext issues a HEAD to the specified URL with context and headers.
func HeadWithContext(ctx context.Context, url string, headers H) (*Response, error) {
	return defaultSession.HeadWithContext(ctx, url, headers)
}

// Post issues a POST to the specified URL with headers.
// Post data should be one of nil, io.Reader, url.Values, string map or struct.
func Post(url string, headers H, data any) (*Response, error) {
	return defaultSession.Post(url, headers, data)
}

// PostWithContext issues a POST to the specified URL with context and headers.
// Post data should be one of nil, io.Reader, url.Values, string map or struct.
func PostWithContext(ctx context.Context, url string, headers H, data any) (*Response, error) {
	return defaultSession.PostWithContext(ctx, url, headers, data)
}

// Upload issues a POST to the specified URL with a multipart document.
func Upload(url string, headers H, params map[string]string, files ...*File) (*Response, error) {
	return defaultSession.Upload(url, headers, params, files...)
}

// UploadWithContext issues a POST to the specified URL with context and a multipart document.
func UploadWithContext(ctx context.Context, url string, headers H, params map[string]string, files ...*File) (*Response, error) {
	return defaultSession.UploadWithContext(ctx, url, headers, params, files...)
}

// Do sends a session HTTP request and returns a response.
func (s *Session) Do(req *http.Request) (*Response, error) {
	req = req.Clone(req.Context())
	header := make(http.Header)
	for k, v := range defaultHeaders() {
		header.Set(k, v)
	}
	for k, v := range s.Header {
		header[k] = v
	}
	for k, v := range req.Header {
		header[k] = v
	}
	req.Header = header
	resp, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}
	return buildResponse(resp)
}

// Get issues a session GET to the specified URL with additional headers.
func (s *Session) Get(url string, headers H) (*Response, error) {
	return s.GetWithContext(context.Background(), url, headers)
}

// GetWithContext issues a session GET to the specified URL with context and additional headers.
func (s *Session) GetWithContext(ctx context.Context, url string, headers H) (*Response, error) {
	req, err := newRequest(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	return s.Do(req)
}

// Head issues a session HEAD to the specified URL with additional headers.
func (s *Session) Head(url string, headers H) (*Response, error) {
	return s.HeadWithContext(context.Background(), url, headers)
}

// HeadWithContext issues a session HEAD to the specified URL with context and additional headers.
func (s *Session) HeadWithContext(ctx context.Context, url string, headers H) (*Response, error) {
	req, err := newRequest(ctx, "HEAD", url, nil)
	if err != nil {
		return nil, err
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	return s.Do(req)
}

// Post issues a session POST to the specified URL with additional headers.
func (s *Session) Post(url string, headers H, data any) (*Response, error) {
	return s.PostWithContext(context.Background(), url, headers, data)
}

// PostWithContext issues a session POST to the specified URL with context and additional headers.
func (s *Session) PostWithContext(ctx context.Context, url string, headers H, data any) (*Response, error) {
	req, err := newRequest(ctx, "POST", url, data)
	if err != nil {
		return nil, err
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	return s.Do(req)
}

// Upload issues a session POST to the specified URL with a multipart document and additional headers.
func (s *Session) Upload(url string, headers H, params map[string]string, files ...*File) (*Response, error) {
	return s.UploadWithContext(context.Background(), url, headers, params, files...)
}

// UploadWithContext issues a session POST to the specified URL with context, a multipart document and additional headers.
func (s *Session) UploadWithContext(ctx context.Context, url string, headers H, params map[string]string, files ...*File) (*Response, error) {
	data, contentType, err := buildMultipart(params, files...)
	if err != nil {
		return nil, err
	}
	req, err := newRequest(ctx, "POST", url, data)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", contentType)
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	return s.Do(req)
}

// KeepAlive repeatedly calls fn with a fixed interval delay between each call.
func (s *Session) KeepAlive(interval *time.Duration, fn func(*Session) error) (err error) {
	for ; err == nil; <-time.After(*interval) {
		err = fn(s)
	}
	return
}

// GetWithClient issues a GET to the specified URL with context, headers and client.
func GetWithClient(ctx context.Context, url string, headers H, client *http.Client) (*Response, error) {
	return newSession(client).GetWithContext(ctx, url, headers)
}

// HeadWithClient issues a HEAD to the specified URL with context, headers and client.
func HeadWithClient(ctx context.Context, url string, headers H, client *http.Client) (*Response, error) {
	return newSession(client).HeadWithContext(ctx, url, headers)
}

// PostWithClient issues a POST to the specified URL with context, headers and client.
func PostWithClient(ctx context.Context, url string, headers H, data any, client *http.Client) (*Response, error) {
	return newSession(client).PostWithContext(ctx, url, headers, data)
}

// UploadWithClient issues a POST to the specified URL with context, a multipart document and client.
func UploadWithClient(ctx context.Context, url string, headers H, params map[string]string, files []*File, client *http.Client) (*Response, error) {
	return newSession(client).UploadWithContext(ctx, url, headers, params, files...)
}
