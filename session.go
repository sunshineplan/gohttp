package gohttp

import (
	"net/http"
	"net/http/cookiejar"
	"net/url"

	"golang.org/x/net/publicsuffix"
)

// Session contains http.Client and http.Header
type Session struct {
	client *http.Client
	Header http.Header
}

// NewSession return a Session with default setting
func NewSession() *Session {
	jar, _ := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
	return &Session{
		client: &http.Client{
			Transport: &http.Transport{Proxy: nil},
			Jar:       jar,
		},
		Header: make(http.Header),
	}
}

// SetProxy sets Session client transport proxy
func (s *Session) SetProxy(proxy string) error {
	proxyURL, err := url.Parse(proxy)
	if err != nil {
		return err
	}
	s.client.Transport = &http.Transport{Proxy: http.ProxyURL(proxyURL)}
	return nil
}

// Cookies returns the cookies to send in a request for the given URL
func (s *Session) Cookies(u *url.URL) []*http.Cookie {
	return s.client.Jar.Cookies(u)
}

// SetCookie handles the receipt of the cookie in a reply for the given URL
func (s *Session) SetCookie(u *url.URL, name, value string) {
	s.SetCookies(u, []*http.Cookie{&http.Cookie{Name: name, Value: value}})
}

// SetCookies handles the receipt of the cookies in a reply for the given URL
func (s *Session) SetCookies(u *url.URL, cookies []*http.Cookie) {
	s.client.Jar.SetCookies(u, cookies)
}

// Get does Session get
func (s *Session) Get(url string, headers H) *Response {
	for k, v := range headers {
		s.Header.Set(k, v)
	}
	return doRequest("GET", url, s.Header, nil, s.client)
}

// Head does Session head
func (s *Session) Head(url string, headers H) *Response {
	for k, v := range headers {
		s.Header.Set(k, v)
	}
	return doRequest("HEAD", url, s.Header, nil, s.client)
}

// Post does Session post
func (s *Session) Post(url string, headers H, data interface{}) *Response {
	for k, v := range headers {
		s.Header.Set(k, v)
	}
	return doRequest("POST", url, s.Header, data, s.client)
}
