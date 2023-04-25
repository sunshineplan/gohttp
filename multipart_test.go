package gohttp

import (
	"testing"
)

func TestPanicF(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("gave no panic; want panic")
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
