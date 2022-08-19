package gohttp

import (
	"bytes"
	"compress/flate"
	"compress/gzip"
	"errors"
	"io"
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
	r = &Response{Response: &http.Response{Body: errReader(0)}}
	if b := r.Bytes(); b != nil {
		t.Error("gave non nil bytes; want nil")
	}
	r = &Response{
		Response: &http.Response{
			Header: http.Header{"Content-Encoding": []string{"gzip"}},
			Body:   errReader(0),
		},
	}
	if b := r.Bytes(); b != nil {
		t.Error("gave non nil bytes; want nil")
	}

	var buf bytes.Buffer
	zw := gzip.NewWriter(&buf)
	zw.Write([]byte("test"))
	zw.Close()
	r = &Response{
		Response: &http.Response{
			Header: http.Header{"Content-Encoding": []string{"gzip"}},
			Body:   io.NopCloser(&buf),
		},
	}
	if b := r.String(); b != "test" {
		t.Errorf("expected %q; got %q", "test", b)
	}
	if b := r.String(); b != "test" { // cached
		t.Errorf("expected %q; got %q", "test", b)
	}

	fw, _ := flate.NewWriter(&buf, -1)
	fw.Write([]byte("deflate"))
	fw.Close()
	r = &Response{
		Response: &http.Response{
			Header: http.Header{"Content-Encoding": []string{"deflate"}},
			Body:   io.NopCloser(&buf),
		},
	}
	if b := r.String(); b != "deflate" {
		t.Errorf("expected %q; got %q", "deflate", b)
	}
}

func TestJSON(t *testing.T) {
	r := buildResponse(nil, errors.New("test"))
	if err := r.JSON(nil); err == nil {
		t.Error("gave nil error; want error")
	}

	r = &Response{
		Response: &http.Response{Body: io.NopCloser(bytes.NewReader([]byte("-1")))},
	}
	var data uint
	if err := r.JSON(&data); err == nil {
		t.Error("gave nil error; want error")
	}

	r = &Response{
		Response: &http.Response{Body: errReader(0)},
	}
	if err := r.JSON(nil); err == nil {
		t.Error("gave nil error; want error")
	}
}

func TestSave(t *testing.T) {
	r := buildResponse(nil, errors.New("test"))
	if _, err := r.Save("error"); err == nil {
		t.Error("gave nil error; want error")
	}

	r = &Response{
		Response: &http.Response{Body: io.NopCloser(bytes.NewBufferString("test"))},
	}
	if _, err := r.Save(""); err == nil {
		t.Error("gave nil error; want error")
	}

	f, _ := os.CreateTemp("", "test")
	f.Close()
	defer os.Remove(f.Name())

	if n, err := r.Save(f.Name()); err != nil {
		t.Fatal(err)
	} else if n != 4 {
		t.Error(n)
	}
	b, err := os.ReadFile(f.Name())
	if err != nil {
		t.Fatal(err)
	}
	if s := string(b); s != "test" {
		t.Errorf("expected %q; got %q", "test", s)
	}
}

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

func TestSetClientPanic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("gave no panic; want panic")
		}
	}()
	SetClient(nil)
}
