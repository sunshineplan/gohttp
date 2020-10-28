package gohttp

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestSession(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.SetCookie(w, &http.Cookie{Name: "hello", Value: "world"})
		fmt.Fprint(w, "Hello, world!")
	}))
	defer ts.Close()
	tsURL, _ := url.Parse(ts.URL)

	s := NewSession()
	s.Header.Set("hello", "world")
	s.SetCookie(tsURL, "one", "first")
	s.SetCookie(tsURL, "two", "second")
	resp := s.Get(ts.URL, H{"another": "header"})
	if resp.Error != nil {
		t.Error(resp.Error)
	}
	defer resp.Close()
	if h := resp.Request.Header.Get("hello"); h != "world" {
		t.Errorf("expected hello header %q; got %q", "world", h)
	}
	if h := resp.Request.Header.Get("another"); h != "header" {
		t.Errorf("expected hello header %q; got %q", "header", h)
	}
	if c := resp.Cookies[0]; c.String() != "hello=world" {
		t.Errorf("expected set cookie %q; got %q", "hello=world", c)
	}
	if c, _ := resp.Request.Cookie("one"); c.String() != "one=first" {
		t.Errorf("expected cookie %q; got %q", "one=first", c)
	}
	if c, _ := resp.Request.Cookie("two"); c.String() != "two=second" {
		t.Errorf("expected cookie %q; got %q", "two=second", c)
	}

	resp = s.Get(ts.URL, nil)
	if resp.Error != nil {
		t.Error(resp.Error)
	}
	defer resp.Close()
	if c, _ := resp.Request.Cookie("hello"); c.String() != "hello=world" {
		t.Errorf("expected cookie %q; got %q", "hello=world", c)
	}
}
