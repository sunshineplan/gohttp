package gohttp

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestSetProxy(t *testing.T) {
	if err := SetProxy("://localhost"); err == nil {
		t.Error("gave nil error; want some error")
	}

	s := NewSession()

	if err := s.SetProxy("http://localhost"); err != nil {
		t.Error(err)
	}

	if err := s.SetProxy("://localhost"); err == nil {
		t.Error("gave nil error; want some error")
	}
}

func TestSetClient(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Fatal("gave panic; want no panic")
		}
	}()
	SetClient(http.DefaultClient)
}

func TestGetAndHead(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Hello, world!")
	}))
	defer ts.Close()

	SetNoProxy()
	resp, err := Get(ts.URL, H{"hello": "world"})
	if err != nil {
		t.Fatal(err)
	}
	if resp.Request.Method != "GET" {
		t.Errorf("expected method %q; got %q", "GET", resp.Request.Method)
	}
	if resp.Request.URL.String() != ts.URL {
		t.Errorf("expected URL %q; got %q", ts.URL, resp.Request.URL.String())
	}
	if h := resp.Request.Header.Get("hello"); h != "world" {
		t.Errorf("expected hello header %q; got %q", "world", h)
	}
	if ua := resp.Request.Header.Get("user-agent"); ua != "Go-HTTP-Client" {
		t.Errorf("expected user agent %q; got %q", "Go-HTTP-Client", ua)
	}
	if s := resp.String(); s != "Hello, world!" {
		t.Error("Incorrect get response body:", s)
	}

	resp, err = Head(ts.URL, H{"hello": "world"})
	if err != nil {
		t.Fatal(err)
	}
	if resp.Request.Method != "HEAD" {
		t.Errorf("expected method %q; got %q", "HEAD", resp.Request.Method)
	}
	if l := resp.Request.ContentLength; l != 0 {
		t.Error("Incorrect head response body:", l)
	}
}

func TestPost(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c, _ := io.ReadAll(r.Body)
		fmt.Fprint(w, string(c))
	}))
	defer ts.Close()

	SetAgent("test")
	resp, err := Post(ts.URL, H{"hello": "world"}, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Close()
	if resp.Request.Method != "POST" {
		t.Errorf("expected method %q; got %q", "POST", resp.Request.Method)
	}
	if resp.Request.URL.String() != ts.URL {
		t.Errorf("expected URL %q; got %q", ts.URL, resp.Request.URL.String())
	}
	if h := resp.Request.Header.Get("hello"); h != "world" {
		t.Errorf("expected hello header %q; got %q", "world", h)
	}
	if ua := resp.Request.Header.Get("user-agent"); ua != "test" {
		t.Errorf("expected user agent %q; got %q", "test", ua)
	}
	if l := resp.Request.ContentLength; l != 0 {
		t.Error("Incorrect response body:", l)
	}

	resp, err = Post(ts.URL, nil, url.Values{"test": []string{"test"}})
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Close()
	if ct := resp.Request.Header.Get("Content-Type"); ct != "application/x-www-form-urlencoded" {
		t.Errorf("expected Content-Type header %q; got %q", "application/x-www-form-urlencoded", ct)
	}
	if s := resp.String(); s != "test=test" {
		t.Errorf("expected response body %q; got %q", "test=test", s)
	}

	resp, err = Post(ts.URL, nil, map[string]any{"test": "test"})
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Close()
	if ct := resp.Request.Header.Get("Content-Type"); ct != "application/json" {
		t.Errorf("expected Content-Type header %q; got %q", "application/json", ct)
	}
	var json struct{ Test string }
	if err := resp.JSON(&json); err != nil {
		t.Error(err)
	}
	if json != struct{ Test string }{Test: "test"} {
		t.Error("Incorrect response body:", json)
	}
}

func TestUpload(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c, _ := io.ReadAll(r.Body)
		fmt.Fprint(w, string(c))
	}))
	defer ts.Close()

	if _, err := Upload(ts.URL, nil, nil, &File{ReadCloser: errReader(0)}); err == nil {
		t.Error("gave nil error; want error")
	}

	if _, err := Upload(ts.URL, H{"header": "value"}, nil, F("readme", "README.md")); err != nil {
		t.Error(err)
	}
}
