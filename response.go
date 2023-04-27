package gohttp

import (
	"compress/flate"
	"compress/gzip"
	"encoding/json"
	"io"
	"net/http"
	"os"

	"github.com/sunshineplan/gohttp/node"
	"golang.org/x/net/html/charset"
)

// Response represents the response from an HTTP request.
type Response struct {
	*http.Response
	Error  error
	cached bool
	bytes  []byte
}

func buildResponse(resp *http.Response, err error) *Response {
	if err != nil {
		return &Response{Response: new(http.Response), Error: err}
	}

	return &Response{Response: resp}
}

// Close closes the response body.
func (r *Response) Close() error {
	if r.Error == nil && r.Response != nil && r.Body != nil {
		return r.Body.Close()
	}
	return nil
}

// Bytes returns a slice of byte of the response body.
func (r *Response) Bytes() []byte {
	if r.Error != nil || r.Response == nil || r.Body == nil {
		return nil
	}
	if r.cached {
		return r.bytes
	}
	defer r.Body.Close()

	var reader io.Reader = r.Body
	switch r.Header.Get("Content-Encoding") {
	case "gzip":
		reader, r.Error = gzip.NewReader(reader)
		if r.Error != nil {
			return nil
		}
	case "deflate":
		reader = flate.NewReader(reader)
	}

	reader, r.Error = charset.NewReader(reader, r.Header.Get("Content-Type"))
	if r.Error != nil {
		return nil
	}
	r.bytes, r.Error = io.ReadAll(reader)
	if r.Error != nil {
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
	if r.Error != nil {
		return r.Error
	}

	return json.Unmarshal(r.Bytes(), data)
}

// Save saves the response data to file. It returns the number
// of bytes written and an error, if any.
func (r *Response) Save(file string) (int, error) {
	if r.Error != nil {
		return 0, r.Error
	}

	f, err := os.Create(file)
	if err != nil {
		return 0, err
	}
	defer f.Close()

	return f.Write(r.Bytes())
}

// Node returns the contents of the response body as a Node.
func (r *Response) Node() (node.Node, error) {
	return node.ParseHTML(r.String())
}
