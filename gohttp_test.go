package gohttp

import (
	"net/url"
	"testing"
)

func TestGoHTTP(t *testing.T) {
	r := Post("https://httpbin.org/post", nil, url.Values{"hello": []string{"world"}})
	var postResp struct {
		Form struct{ Hello string }
	}
	if err := r.JSON(&postResp); err != nil {
		t.Error(err)
	}
	if h := postResp.Form.Hello; h != "world" {
		t.Errorf("expected %q; got %q", "world", h)
	}

	s := NewSession()
	s.Header.Set("hello", "world")
	s.Get("https://httpbin.org/cookies/set/name/value", nil)
	r = s.Get("https://httpbin.org/get", nil)
	var getResp struct {
		Headers struct{ Hello, Cookie string }
	}
	if err := r.JSON(&getResp); err != nil {
		t.Error(err)
	}
	if h := getResp.Headers.Hello; h != "world" {
		t.Errorf("expected %q; got %q", "world", h)
	}
	if h := getResp.Headers.Cookie; h != "name=value" {
		t.Errorf("expected %q; got %q", "name=value", h)
	}
}
