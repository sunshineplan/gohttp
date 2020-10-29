package gohttp

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestSession(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.SetCookie(w, &http.Cookie{Name: "hello", Value: "world"})
		c, _ := ioutil.ReadAll(r.Body)
		fmt.Fprint(w, string(c))
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
	if resp.Request.Method != "GET" {
		t.Errorf("expected method %q; got %q", "GET", resp.Request.Method)
	}
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

	resp = s.Head(ts.URL, H{"another": "header"})
	if resp.Error != nil {
		t.Error(resp.Error)
	}
	if resp.Request.Method != "HEAD" {
		t.Errorf("expected method %q; got %q", "HEAD", resp.Request.Method)
	}
	defer resp.Close()
	if c := s.Cookies(resp.Request.URL); len(c) != 3 {
		t.Errorf("expected cookies number %d; got %d", 3, len(c))
	}

	resp = s.Post(ts.URL, H{"another": "header"}, bytes.NewBufferString("Hello, world!"))
	if resp.Error != nil {
		t.Error(resp.Error)
	}
	if resp.Request.Method != "POST" {
		t.Errorf("expected method %q; got %q", "POST", resp.Request.Method)
	}
	defer resp.Close()
	if s := resp.String(); s != "Hello, world!" {
		t.Errorf("expected response body %q; got %q", "Hello, world!", s)
	}
}

func TestSetProxy(t *testing.T) {
	s := NewSession()
	if err := s.SetProxy("http://localhost"); err != nil {
		t.Error(err)
	}
	if err := s.SetProxy("://localhost"); err == nil {
		t.Error("gave nil error; want some error")
	}
}
