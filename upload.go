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
	io.ReadCloser
	Fieldname string
	Filename  string
}

// F opens the file for creating File.
func F(fieldname, filename string) *File {
	f, err := os.Open(filename)
	if err != nil {
		panic(err)
	}

	return &File{
		ReadCloser: f,
		Fieldname:  fieldname,
		Filename:   filepath.Base(filename),
	}
}

func buildMultipart(params map[string]string, files ...*File) (io.Reader, string, error) {
	var body bytes.Buffer
	w := multipart.NewWriter(&body)
	defer w.Close()

	for _, file := range files {
		if file == nil {
			continue
		}
		defer file.Close()

		part, err := w.CreateFormFile(file.Fieldname, file.Filename)
		if err != nil {
			return nil, "", err
		}

		if _, err := io.Copy(part, file); err != nil {
			return nil, "", err
		}
	}

	for k, v := range params {
		w.WriteField(k, v)
	}

	return &body, w.FormDataContentType(), nil
}

// Upload issues a POST to the specified URL with a multipart document.
func Upload(url string, headers H, params map[string]string, files ...*File) *Response {
	return UploadWithClient(url, headers, params, files, defaultClient)
}

// UploadWithClient issues a POST to the specified URL with a multipart document and client.
func UploadWithClient(url string, headers H, params map[string]string, files []*File, client *http.Client) *Response {
	r, contentType, err := buildMultipart(params, files...)
	if err != nil {
		return &Response{Error: err}
	}
	h := H{"Content-Type": contentType}

	for k, v := range headers {
		h[k] = v
	}

	return PostWithClient(url, h, r, client)
}
