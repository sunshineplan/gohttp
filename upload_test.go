package gohttp

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPanicF(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("gave nil panic; want panic")
		}
	}()
	F("", "I am Not A File")
}

func TestBuildMultipart(t *testing.T) {
	f := &File{ReadCloser: errReader(0)}
	if _, _, err := buildMultipart(nil, f); err == nil {
		t.Error("gave nil error; want error")
	}
	f.Fieldname = "test"
	f.Filename = "test"
	if _, _, err := buildMultipart(nil, f); err == nil {
		t.Error("gave nil error; want error")
	}
}

func TestUpload(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c, _ := ioutil.ReadAll(r.Body)
		fmt.Fprint(w, string(c))
	}))
	defer ts.Close()
	resp := Upload(ts.URL, nil, nil, &File{ReadCloser: errReader(0)})
	if resp.Error == nil {
		t.Error("gave nil error; want error")
	}
	resp = Upload(ts.URL, H{"header": "value"}, nil, F("readme", "README.md"))
	if resp.Error != nil {
		t.Error(resp.Error)
	}
}
