package gohttp

import (
	"bytes"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

// File contains the file part of a multipart message.
type File struct {
	io.Reader
	Fieldname string
	Filename  string
}

// F opens the file for creating File.
func F(fieldname, filename string) (file *File) {
	var err error
	file.Reader, err = os.Open(filename)
	if err != nil {
		panic(err)
	}
	file.Fieldname = fieldname
	file.Filename = filepath.Base(filename)
	return
}

// Upload issues a POST to the specified URL with a multipart document.
func Upload(url string, headers H, params map[string]string, files []*File) *Response {
	return UploadWithClient(url, headers, params, files, defaultClient)
}

// UploadWithClient issues a POST to the specified URL with a multipart document and client.
func UploadWithClient(url string, headers H, params map[string]string, files []*File, client *http.Client) *Response {
	var body bytes.Buffer
	w := multipart.NewWriter(&body)
	defer w.Close()

	for _, file := range files {
		if file == nil {
			continue
		}
		part, err := w.CreateFormFile(file.Fieldname, file.Filename)
		if err != nil {
			return &Response{Error: err}
		}
		if _, err := io.Copy(part, file); err != nil {
			return &Response{Error: err}
		}
	}
	for k, v := range params {
		w.WriteField(k, v)
	}

	h := H{"Content-Type": w.FormDataContentType()}
	for k, v := range headers {
		h[k] = v
	}
	return PostWithClient(url, h, &body, defaultClient)
}
