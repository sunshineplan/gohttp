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

func newSession(client *http.Client) *Session {
	return &Session{client, make(http.Header)}
}

// NewSession creates and initializes a new Session using initial contents.
func NewSession() *Session {
	c := new(http.Client)
	c.Jar, _ = cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
	return newSession(c)
}

func (s *Session) setProxy(fn func(*http.Request) (*url.URL, error)) {
	var tr *http.Transport
	var ok bool
	if s.client.Transport == nil {
		if tr, ok = http.DefaultTransport.(*http.Transport); ok {
			tr.Proxy = fn
		}
	} else {
		if tr, ok = s.client.Transport.(*http.Transport); ok {
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
	s.client.Timeout = d
}

// SetClient sets default client.
func (s *Session) SetClient(c *http.Client) {
	s.client = c
}

// Cookies returns the cookies to send in a request for the given URL.
func (s *Session) Cookies(u *url.URL) []*http.Cookie {
	return s.client.Jar.Cookies(u)
}

// SetCookie handles the receipt of the cookie in a reply for the given URL.
func (s *Session) SetCookie(u *url.URL, name, value string) {
	s.SetCookies(u, []*http.Cookie{{Name: name, Value: value}})
}

// SetCookies handles the receipt of the cookies in a reply for the given URL.
func (s *Session) SetCookies(u *url.URL, cookies []*http.Cookie) {
	s.client.Jar.SetCookies(u, cookies)
}
