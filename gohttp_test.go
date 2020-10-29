package gohttp

import (
	"errors"
	"net/url"
	"testing"
)

func TestBuildResponse(t *testing.T) {
	r := buildResponse(nil, errors.New("test"))
	if err := r.JSON(nil); err == nil {
		t.Error("gave nil error; want error")
	}
	if b := r.Bytes(); b != nil {
		t.Error("gave non nil bytes; want nil")
	}
}

func TestBuildRequest(t *testing.T) {
	if _, err := buildRequest("bad method", "", url.Values{}); err == nil {
		t.Error("gave nil error; want error")
	}
	if _, err := buildRequest("bad method", "", "test"); err == nil {
		t.Error("gave nil error; want error")
	}
}

func TestDoRequest(t *testing.T) {
	if r := doRequest("bad method", "", nil, nil, nil); r.Error == nil {
		t.Error("gave nil error; want error")
	}
}
