package gohttp

import (
	"compress/flate"
	"compress/gzip"
	"encoding/json"
	"io"
	"net/http"
	"os"

	"golang.org/x/net/html/charset"
)

// Response represents the response from an HTTP request.
type Response struct {
	*http.Response
	cached bool
	bytes  []byte
}

// Close closes the response body.
func (r *Response) Close() error {
	return r.Body.Close()
}

// Bytes returns a slice of byte of the response body.
func (r *Response) Bytes() []byte {
	if r.cached {
		return r.bytes
	}
	defer r.Close()

	var reader io.Reader = r.Body
	switch r.Header.Get("Content-Encoding") {
	case "gzip":
		var err error
		reader, err = gzip.NewReader(reader)
		if err != nil {
			return nil
		}
	case "deflate":
		reader = flate.NewReader(reader)
	}

	reader, err := charset.NewReader(reader, r.Header.Get("Content-Type"))
	if err != nil {
		return nil
	}
	r.bytes, err = io.ReadAll(reader)
	if err != nil {
		return nil
	}
	r.cached = true

	return r.bytes
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
