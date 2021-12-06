package gohttp

import (
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"time"

	"golang.org/x/net/publicsuffix"
)

var _ http.CookieJar = &Session{}

// Session provides cookie persistence and configuration.
type Session struct {
	client *http.Client
	Header http.Header
}

// NewSession creates and initializes a new Session using initial contents.
func NewSession() *Session {
	client := *defaultClient
	jar, _ := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})

	client.Jar = jar

	return &Session{
		client: &client,
		Header: make(http.Header),
	}
}

// SetProxy sets Session client transport proxy.
func (s *Session) SetProxy(proxy string) error {
	proxyURL, err := url.Parse(proxy)
	if err != nil {
		return err
	}

	s.client.Transport = &http.Transport{Proxy: http.ProxyURL(proxyURL)}

	return nil
}

// SetNoProxy sets Session client use no proxy.
func (s *Session) SetNoProxy() {
	s.client.Transport = &http.Transport{Proxy: nil}
}

// SetProxyFromEnvironment sets Session client use environment proxy.
func (s *Session) SetProxyFromEnvironment() {
	s.client.Transport = &http.Transport{Proxy: http.ProxyFromEnvironment}
}

// SetTimeout sets Session client timeout. Zero means no timeout.
func (s *Session) SetTimeout(d time.Duration) {
	s.client.Timeout = d
}

// SetClient sets default client.
func (s *Session) SetClient(c *http.Client) {
	if c != nil {
		s.client = c
	} else {
		panic("cannot set a nil client")
	}
}

// Cookies returns the cookies to send in a request for the given URL.
func (s *Session) Cookies(u *url.URL) []*http.Cookie {
	if u == nil {
		panic("url pointer is nil")
	}

	return s.client.Jar.Cookies(u)
}

// SetCookie handles the receipt of the cookie in a reply for the given URL.
func (s *Session) SetCookie(u *url.URL, name, value string) {
	s.SetCookies(u, []*http.Cookie{{Name: name, Value: value}})
}

// SetCookies handles the receipt of the cookies in a reply for the given URL.
func (s *Session) SetCookies(u *url.URL, cookies []*http.Cookie) {
	if u == nil {
		panic("url pointer is nil")
	}

	s.client.Jar.SetCookies(u, cookies)
}

// Get issues a session GET to the specified URL with additional headers.
func (s *Session) Get(url string, headers H) *Response {
	for k, v := range headers {
		s.Header.Set(k, v)
	}

	return doRequest("GET", url, s.Header, nil, s.client)
}

// Head issues a session HEAD to the specified URL with additional headers.
func (s *Session) Head(url string, headers H) *Response {
	for k, v := range headers {
		s.Header.Set(k, v)
	}

	return doRequest("HEAD", url, s.Header, nil, s.client)
}

// Post issues a session POST to the specified URL with additional headers.
func (s *Session) Post(url string, headers H, data interface{}) *Response {
	for k, v := range headers {
		s.Header.Set(k, v)
	}

	return doRequest("POST", url, s.Header, data, s.client)
}

// Upload issues a session POST to the specified URL with a multipart document and additional headers.
func (s *Session) Upload(url string, headers H, params map[string]string, files ...*File) *Response {
	data, contentType, err := buildMultipart(params, files...)
	if err != nil {
		return &Response{Error: err}
	}
	s.Header.Add("Content-Type", contentType)

	for k, v := range headers {
		s.Header.Set(k, v)
	}

	return doRequest("POST", url, s.Header, data, s.client)
}

// KeepAlive repeatedly calls fn with a fixed interval delay between each call.
func (s *Session) KeepAlive(interval time.Duration, fn func(*Session) error) (err error) {
	t := time.NewTicker(interval)
	defer t.Stop()

	for range t.C {
		if err = fn(s); err != nil {
			return
		}
	}

	return
}
