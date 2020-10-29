package gohttp

import (
	"bytes"
	"errors"
	"io/ioutil"
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
}

func TestDoRequest(t *testing.T) {
	if r := doRequest("bad method", "", nil, nil, nil); r.Error == nil {
		t.Error("gave nil error; want error")
	}
}

type errReader int

func (errReader) Read(p []byte) (n int, err error) {
	return 0, errors.New("test error")
}
func (errReader) Close() error { return nil }

func TestJSON(t *testing.T) {
	r := buildResponse(nil, errors.New("test"))
	if err := r.JSON(nil); err == nil {
		t.Error("gave nil error; want error")
	}
	r = &Response{Body: ioutil.NopCloser(bytes.NewReader([]byte("-1")))}
	var data uint
	if err := r.JSON(&data); err == nil {
		t.Error("gave nil error; want error")
	}
	r = &Response{Body: errReader(0)}
	if err := r.JSON(nil); err == nil {
		t.Error("gave nil error; want error")
	}
}

func TestBytes(t *testing.T) {
	r := &Response{Body: errReader(0)}
	if b := r.Bytes(); b != nil {
		t.Error("gave non nil bytes; want nil")
	}
}
