package gohttp

import "net/http"

// Get issues a GET to the specified URL with headers.
func Get(url string, headers H) *Response {
	return GetWithClient(url, headers, defaultClient)
}

// GetWithClient issues a GET to the specified URL with headers and client.
func GetWithClient(url string, headers H, client *http.Client) *Response {
	return doRequest("GET", url, buildHeader(headers), nil, client)
}

// Head issues a HEAD to the specified URL with headers.
func Head(url string, headers H) *Response {
	return HeadWithClient(url, headers, defaultClient)
}

// HeadWithClient issues a HEAD to the specified URL with headers and client.
func HeadWithClient(url string, headers H, client *http.Client) *Response {
	return doRequest("HEAD", url, buildHeader(headers), nil, client)
}

// Post issues a POST to the specified URL with headers.
// Post data should be one of nil, io.Reader, url.Values, string map or struct.
func Post(url string, headers H, data any) *Response {
	return PostWithClient(url, headers, data, defaultClient)
}

// PostWithClient issues a POST to the specified URL with headers and client.
func PostWithClient(url string, headers H, data any, client *http.Client) *Response {
	return doRequest("POST", url, buildHeader(headers), data, client)
}
