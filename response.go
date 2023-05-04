package gohttp

import (
	"bytes"
	"compress/flate"
	"compress/gzip"
	"encoding/json"
	"io"
	"net/http"
	"os"

	"golang.org/x/net/html/charset"
)

var _ io.ReadCloser = &Response{}

// Response represents the response from an HTTP request.
type Response struct {
	resp *http.Response
	body io.Reader

	// StatusCode represents the response status code.
	StatusCode int
	// Header maps header keys to values.
	Header http.Header
	// ContentLength records the length of the associated content.
	ContentLength int64

	buf    *bytes.Buffer
	cached bool
}

func buildResponse(resp *http.Response) (*Response, error) {
	var reader io.Reader = resp.Body
	switch resp.Header.Get("Content-Encoding") {
	case "gzip":
		var err error
		reader, err = gzip.NewReader(reader)
		if err != nil {
			return nil, err
		}
	case "deflate":
		reader = flate.NewReader(reader)
	}
	buf := new(bytes.Buffer)
	body, err := charset.NewReader(reader, resp.Header.Get("Content-Type"))
	if err != nil {
		if err == io.EOF {
			body = reader
		} else {
			return nil, err
		}
	}
	body = io.TeeReader(body, buf)
	return &Response{
		resp:          resp,
		body:          body,
		StatusCode:    resp.StatusCode,
		Header:        resp.Header,
		ContentLength: resp.ContentLength,
		buf:           buf,
	}, nil
}

// Read reads the response body.
func (r *Response) Read(p []byte) (int, error) {
	return r.body.Read(p)
}

// Close closes the response body.
func (r *Response) Close() error {
	return r.resp.Body.Close()
}

// Raw returns origin *http.Response.
func (r *Response) Raw() *http.Response {
	return r.resp
}

// Request is the request that was sent to obtain this Response.
func (r *Response) Request() *http.Request {
	return r.resp.Request
}

// Cookies parses and returns the cookies set in the Set-Cookie headers.
func (r *Response) Cookies() []*http.Cookie {
	return r.resp.Cookies()
}

// Bytes returns a slice of byte of the response body.
func (r *Response) Bytes() []byte {
	if r.cached {
		return r.buf.Bytes()
	}
	defer r.Close()

	if _, err := io.ReadAll(r); err != nil {
		return nil
	}
	r.cached = true

	return r.buf.Bytes()
}

// String returns the contents of the response body as a string.
func (r *Response) String() string {
	return string(r.Bytes())
}

// JSON parses the response body as JSON-encoded data
// and stores the result in the value pointed to by data.
func (r *Response) JSON(data any) error {
	return json.Unmarshal(r.Bytes(), data)
}

// Save saves the response data to file. It returns the number
// of bytes written and an error, if any.
func (r *Response) Save(file string) (int, error) {
	f, err := os.Create(file)
	if err != nil {
		return 0, err
	}
	defer f.Close()

	return f.Write(r.Bytes())
}
