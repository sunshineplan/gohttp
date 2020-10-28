package gohttp

import "net/http"

// Get request
func Get(url string, headers H) *Response {
	return GetWithClient(url, headers, defaultClient)
}

// GetWithClient request with a custom client
func GetWithClient(url string, headers H, client *http.Client) *Response {
	h := make(http.Header)
	for k, v := range headers {
		h.Set(k, v)
	}
	return doRequest("GET", url, h, nil, client)
}

// Head request
func Head(url string, headers H) *Response {
	return HeadWithClient(url, headers, defaultClient)
}

// HeadWithClient request with a custom client
func HeadWithClient(url string, headers H, client *http.Client) *Response {
	h := make(http.Header)
	for k, v := range headers {
		h.Set(k, v)
	}
	return doRequest("HEAD", url, h, nil, client)
}

// Post request
func Post(url string, headers H, data interface{}) *Response {
	return PostWithClient(url, headers, data, defaultClient)
}

// PostWithClient request with a custom client
func PostWithClient(url string, headers H, data interface{}, client *http.Client) *Response {
	h := make(http.Header)
	for k, v := range headers {
		h.Set(k, v)
	}
	return doRequest("POST", url, h, data, client)
}
