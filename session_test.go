package gohttp

import (
	"bytes"
	"fmt"
	"io"
	"mime"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestSession(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.SetCookie(w, &http.Cookie{Name: "hello", Value: "world"})
		c, _ := io.ReadAll(r.Body)
		fmt.Fprint(w, string(c))
	}))
	defer ts.Close()
	tsURL, _ := url.Parse(ts.URL)

	s := NewSession()
	s.Header.Set("hello", "world")
	s.SetCookie(tsURL, "one", "first")
	s.SetCookie(tsURL, "two", "second")
	resp, err := s.Get(ts.URL, H{"another": "header"})
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Close()

	if resp.Request().Method != "GET" {
		t.Errorf("expected method %q; got %q", "GET", resp.Request().Method)
	}
	if h := resp.Request().Header.Get("hello"); h != "world" {
		t.Errorf("expected hello header %q; got %q", "world", h)
	}
	if h := resp.Request().Header.Get("another"); h != "header" {
		t.Errorf("expected hello header %q; got %q", "header", h)
	}
	if c := resp.Cookies()[0]; c.String() != "hello=world" {
		t.Errorf("expected set cookie %q; got %q", "hello=world", c)
	}
	if c, _ := resp.Request().Cookie("one"); c.String() != "one=first" {
		t.Errorf("expected cookie %q; got %q", "one=first", c)
	}
	if c, _ := resp.Request().Cookie("two"); c.String() != "two=second" {
		t.Errorf("expected cookie %q; got %q", "two=second", c)
	}

	resp, err = s.Head(ts.URL, H{"another": "header"})
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Close()
	if resp.Request().Method != "HEAD" {
		t.Errorf("expected method %q; got %q", "HEAD", resp.Request().Method)
	}
	if c := s.Cookies(resp.Request().URL); len(c) != 3 {
		t.Errorf("expected cookies number %d; got %d", 3, len(c))
	}

	resp, err = s.Post(ts.URL, H{"post": "header"}, bytes.NewBufferString("Hello, world!"))
	if err != nil {
		t.Fatal(err)
	}
	if resp.Request().Method != "POST" {
		t.Errorf("expected method %q; got %q", "POST", resp.Request().Method)
	}
	if s := resp.String(); s != "Hello, world!" {
		t.Errorf("expected response body %q; got %q", "Hello, world!", s)
	}

	if _, err := s.Upload(ts.URL, nil, nil, &File{ReadCloser: errReader(0)}); err == nil {
		t.Error("gave nil error; want error")
	}

	f, _ := os.CreateTemp("", "test")
	f.WriteString("tempfile")
	f.Close()
	defer os.Remove(f.Name())

	resp, err = s.Upload(ts.URL, H{"upload": "header"}, map[string]string{"param": "test"}, F("file1", f.Name()), nil, F("file2", f.Name()))
	if err != nil {
		t.Fatal(err)
	}
	if resp.Request().Method != "POST" {
		t.Errorf("expected method %q; got %q", "POST", resp.Request().Method)
	}
	_, params, err := mime.ParseMediaType(resp.Request().Header.Get("Content-Type"))
	if err != nil {
		t.Error(err)
	}
	mr := multipart.NewReader(resp, params["boundary"])
	for {
		p, err := mr.NextPart()
		if err == io.EOF {
			return
		}

		if err != nil {
			t.Error(err)
		}

		switch p.FormName() {
		case "param":
			b, err := io.ReadAll(p)
			if err != nil {
				t.Fatal(err)
			}

			if s := string(b); s != "test" {
				t.Errorf("expected %q; got %q", "test", s)
			}
		case "file1", "file2":
			if fn := p.FileName(); fn != filepath.Base(f.Name()) {
				t.Errorf("expected %q; got %q", filepath.Base(f.Name()), fn)
			}

			b, err := io.ReadAll(p)
			if err != nil {
				t.Fatal(err)
			}

			if s := string(b); s != "tempfile" {
				t.Errorf("expected %q; got %q", "tempfile", s)
			}
		}
	}
}

func TestSessionTimeout(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(time.Second)
		fmt.Fprint(w, "sleep 1 second")
	}))
	defer ts.Close()

	s := NewSession()
	s.SetTimeout(100 * time.Millisecond)
	if _, err := s.Get(ts.URL, nil); err == nil {
		t.Fatal("gave nil error; want error")
	}
}

func TestCookies(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("gave no panic; want panic")
		}
	}()
	NewSession().Cookies(nil)
}

func TestSetCookie(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("gave no panic; want panic")
		}
	}()
	NewSession().SetCookie(nil, "", "")
}
