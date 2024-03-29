package gohttp

import (
	"bytes"
	"io"
	"mime/multipart"
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
	data := new(bytes.Buffer)
	w := multipart.NewWriter(data)
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

	return data, w.FormDataContentType(), nil
}
