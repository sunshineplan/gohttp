package gohttp

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httputil"
	"regexp"
)

type debugger struct {
	rt       http.RoundTripper
	w        io.Writer
	reqBody  bool
	respBody bool
}

func (t *debugger) RoundTrip(req *http.Request) (*http.Response, error) {
	if t.rt == nil {
		t.rt = http.DefaultTransport
	}
	reqBody, err := httputil.DumpRequestOut(req, t.reqBody)
	if err != nil {
		return nil, err
	}
	t.Write("-> ", reqBody)
	res, err := t.rt.RoundTrip(req)
	if err != nil {
		return nil, err
	}
	respBody, err := httputil.DumpResponse(res, t.respBody)
	if err != nil {
		return nil, err
	}
	t.Write("<- ", respBody)
	return res, nil
}

var lineStart = regexp.MustCompile(`(?m)^`)

func (w *debugger) Write(prefix string, buf []byte) {
	w.w.Write(append(lineStart.ReplaceAll(bytes.TrimSpace(buf), []byte(prefix)), '\n', '\n'))
}
