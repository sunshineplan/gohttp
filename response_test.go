package gohttp

import (
	"bytes"
	"compress/flate"
	"compress/gzip"
	"errors"
	"io"
	"net/http"
	"os"
	"testing"
)

type errReader int

func (errReader) Read(p []byte) (n int, err error) {
	return 0, errors.New("test error")
}
func (errReader) Close() error { return nil }

func TestBytes(t *testing.T) {
	r, err := buildResponse(&http.Response{Body: errReader(0)})
	if err != nil {
		t.Fatal(err)
	}
	if b := r.Bytes(); b != nil {
		t.Error("gave non nil bytes; want nil")
	}
	if _, err := buildResponse(&http.Response{
		Header: http.Header{"Content-Encoding": []string{"gzip"}},
		Body:   errReader(0),
	}); err == nil {
		t.Error("gave nil error; want test error")
	}

	var buf bytes.Buffer
	zw := gzip.NewWriter(&buf)
	zw.Write([]byte("test"))
	zw.Close()
	r, err = buildResponse(&http.Response{
		Header: http.Header{"Content-Encoding": []string{"gzip"}},
		Body:   io.NopCloser(&buf),
	})
	if err != nil {
		t.Fatal(err)
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
	r, _ = buildResponse(&http.Response{
		Header: http.Header{"Content-Encoding": []string{"deflate"}},
		Body:   io.NopCloser(&buf),
	})
	if b := r.String(); b != "deflate" {
		t.Errorf("expected %q; got %q", "deflate", b)
	}
}

func TestJSON(t *testing.T) {
	r, err := buildResponse(&http.Response{Body: io.NopCloser(bytes.NewReader([]byte("-1")))})
	if err != nil {
		t.Fatal(err)
	}
	var data uint
	if err := r.JSON(&data); err == nil {
		t.Error("gave nil error; want error")
	}

	r, err = buildResponse(&http.Response{Body: errReader(0)})
	if err != nil {
		t.Fatal(err)
	}
	if err := r.JSON(nil); err == nil {
		t.Error("gave nil error; want error")
	}
}

func TestSave(t *testing.T) {
	r, _ := buildResponse(&http.Response{Body: io.NopCloser(bytes.NewBufferString("test"))})
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
