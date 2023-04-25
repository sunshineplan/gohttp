package gohttp

import (
	"bytes"
	"io"
	"net/url"
	"testing"
)

func TestBuildRequest(t *testing.T) {
	if _, err := buildRequest("bad method", "", url.Values{}); err == nil {
		t.Error("gave nil error; want error")
	}

	if _, err := buildRequest("bad method", "", "test"); err == nil {
		t.Error("gave nil error; want error")
	}

	if _, err := buildRequest("bad method", "", make(chan int)); err == nil {
		t.Error("gave nil error; want error")
	}

	r, err := buildRequest("", "", bytes.NewBufferString("test"))
	if err != nil {
		t.Error(err)
	}

	if b, _ := io.ReadAll(r.Body); string(b) != "test" {
		t.Errorf("expected request body %q; got %q", "test", string(b))
	}
}

func TestDoRequest(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("gave no panic; want panic")
		}
	}()
	if r := doRequest("bad method", "", nil, nil, nil); r.Error == nil {
		t.Error("gave nil error; want error")
	}
	doRequest("GET", "", nil, nil, nil)
}
