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
	*http.Client
	Header http.Header
}

// NewSession creates and initializes a new Session using initial contents.
func NewSession() *Session {
	client := *defaultClient
	client.Jar, _ = cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})

	return &Session{
		Client: &client,
		Header: make(http.Header),
	}
}

func (s *Session) setProxy(fn func(*http.Request) (*url.URL, error)) {
	var tr *http.Transport
	var ok bool
	if s.Transport == nil {
		if tr, ok = http.DefaultTransport.(*http.Transport); ok {
			tr.Proxy = fn
		}
	} else {
		if tr, ok = s.Transport.(*http.Transport); ok {
			tr.Proxy = fn
		}
	}
	if !ok {
		panic("Transport is not *http.Transport type")
	}
}

// SetProxy sets Session client transport proxy.
func (s *Session) SetProxy(proxy string) error {
	proxyURL, err := url.Parse(proxy)
	if err != nil {
		return err
	}

	s.setProxy(http.ProxyURL(proxyURL))

	return nil
}

// SetNoProxy sets Session client use no proxy.
func (s *Session) SetNoProxy() {
	s.setProxy(nil)
}

// SetProxyFromEnvironment sets Session client use environment proxy.
func (s *Session) SetProxyFromEnvironment() {
	s.setProxy(http.ProxyFromEnvironment)
}

// SetTimeout sets Session client timeout. Zero means no timeout.
func (s *Session) SetTimeout(d time.Duration) {
	s.Timeout = d
}

// SetClient sets default client.
func (s *Session) SetClient(c *http.Client) {
	s.Client = c
}

// Cookies returns the cookies to send in a request for the given URL.
func (s *Session) Cookies(u *url.URL) []*http.Cookie {
	return s.Jar.Cookies(u)
}

// SetCookie handles the receipt of the cookie in a reply for the given URL.
func (s *Session) SetCookie(u *url.URL, name, value string) {
	s.SetCookies(u, []*http.Cookie{{Name: name, Value: value}})
}

// SetCookies handles the receipt of the cookies in a reply for the given URL.
func (s *Session) SetCookies(u *url.URL, cookies []*http.Cookie) {
	s.Jar.SetCookies(u, cookies)
}

// Get issues a session GET to the specified URL with additional headers.
func (s *Session) Get(url string, headers H) *Response {
	for k, v := range headers {
		s.Header.Set(k, v)
	}

	return doRequest("GET", url, s.Header, nil, s.Client)
}

// Head issues a session HEAD to the specified URL with additional headers.
func (s *Session) Head(url string, headers H) *Response {
	for k, v := range headers {
		s.Header.Set(k, v)
	}

	return doRequest("HEAD", url, s.Header, nil, s.Client)
}

// Post issues a session POST to the specified URL with additional headers.
func (s *Session) Post(url string, headers H, data any) *Response {
	for k, v := range headers {
		s.Header.Set(k, v)
	}

	return doRequest("POST", url, s.Header, data, s.Client)
}

// Upload issues a session POST to the specified URL with a multipart document and additional headers.
func (s *Session) Upload(url string, headers H, params map[string]string, files ...*File) *Response {
	data, contentType, err := buildMultipart(params, files...)
	if err != nil {
		return &Response{Response: new(http.Response), Error: err}
	}
	s.Header.Add("Content-Type", contentType)

	for k, v := range headers {
		s.Header.Set(k, v)
	}

	return doRequest("POST", url, s.Header, data, s.Client)
}

// KeepAlive repeatedly calls fn with a fixed interval delay between each call.
func (s *Session) KeepAlive(interval *time.Duration, fn func(*Session) error) (err error) {
	for ; err == nil; <-time.After(*interval) {
		err = fn(s)
	}
	return
}
