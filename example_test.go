package gohttp

import (
	"fmt"
	"log"
	"net/url"
)

func Example() {
	r := Post("https://httpbin.org/post", nil, url.Values{"hello": []string{"world"}})
	var postResp struct {
		Form struct{ Hello string }
	}
	if err := r.JSON(&postResp); err != nil {
		log.Fatal(err)
	}
	fmt.Println(postResp.Form.Hello)
	// Output: world
}

func ExampleSession() {
	s := NewSession()
	s.Header.Set("hello", "world")
	s.Get("https://httpbin.org/cookies/set/name/value", nil)
	r := s.Get("https://httpbin.org/get", nil)
	var getResp struct {
		Headers struct{ Hello, Cookie string }
	}
	if err := r.JSON(&getResp); err != nil {
		log.Fatal(err)
	}
	fmt.Println(getResp.Headers.Hello, getResp.Headers.Cookie)
	// Output: world name=value
}
