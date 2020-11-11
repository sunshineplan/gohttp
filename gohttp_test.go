package gohttp

import (
	"bytes"
	"compress/flate"
	"compress/gzip"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
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
	if b, _ := ioutil.ReadAll(r.Body); string(b) != "test" {
		t.Errorf("expected request body %q; got %q", "test", string(b))
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

func TestBytes(t *testing.T) {
	r := &Response{Error: errors.New("test")}
	if b := r.Bytes(); b != nil {
		t.Error("gave non nil bytes; want nil")
	}
	r = &Response{Body: errReader(0)}
	if b := r.Bytes(); b != nil {
		t.Error("gave non nil bytes; want nil")
	}
	r = &Response{
		Header: http.Header{"Content-Encoding": []string{"gzip"}},
		Body:   errReader(0)}
	if b := r.Bytes(); b != nil {
		t.Error("gave non nil bytes; want nil")
	}

	var buf bytes.Buffer
	zw := gzip.NewWriter(&buf)
	zw.Write([]byte("test"))
	zw.Close()
	r = &Response{
		Header: http.Header{"Content-Encoding": []string{"gzip"}},
		Body:   ioutil.NopCloser(&buf)}
	if b := r.String(); b != "test" {
		t.Errorf("expected %q; got %q", "test", b)
	}
	if b := r.String(); b != "test" {
		t.Errorf("expected %q; got %q", "test", b)
	}

	fw, _ := flate.NewWriter(&buf, -1)
	fw.Write([]byte("deflate"))
	fw.Close()
	r = &Response{
		Header: http.Header{"Content-Encoding": []string{"deflate"}},
		Body:   ioutil.NopCloser(&buf)}
	if b := r.String(); b != "deflate" {
		t.Errorf("expected %q; got %q", "deflate", b)
	}
}

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

func TestSave(t *testing.T) {
	r := buildResponse(nil, errors.New("test"))
	if err := r.Save("error"); err == nil {
		t.Error("gave nil error; want error")
	}
	r = &Response{Body: ioutil.NopCloser(bytes.NewBufferString("test"))}
	if err := r.Save(""); err == nil {
		t.Error("gave nil error; want error")
	}
	f, _ := ioutil.TempFile("", "test")
	f.Close()
	defer os.Remove(f.Name())
	if err := r.Save(f.Name()); err != nil {
		t.Error(err)
	}
	b, err := ioutil.ReadFile(f.Name())
	if err != nil {
		t.Error(err)
	}
	if s := string(b); s != "test" {
		t.Errorf("expected %q; got %q", "test", s)
	}
}
