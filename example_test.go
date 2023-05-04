package gohttp

import (
	"fmt"
	"log"
	"net/url"
	"strings"
)

func Example() {
	resp, err := Post("https://httpbin.org/post", nil, url.Values{"hello": []string{"world"}})
	if err != nil {
		log.Fatal(err)
	}
	var data struct {
		Form struct{ Hello string }
	}
	if err := resp.JSON(&data); err != nil {
		log.Fatal(err)
	}
	fmt.Println(data.Form.Hello)
	// world
}

func ExampleUpload() {
	resp, err := Upload("https://httpbin.org/post", nil, nil, F("readme", "README.md"))
	if err != nil {
		log.Fatal(err)
	}
	var data struct {
		Files   struct{ Readme string }
		Headers struct {
			ContentType string `json:"Content-Type"`
		}
	}
	if err := resp.JSON(&data); err != nil {
		log.Fatal(err)
	}
	fmt.Println(strings.Split(data.Headers.ContentType, ";")[0])
	// multipart/form-data
}

func ExampleSession() {
	s := NewSession()
	s.Header.Set("hello", "world")
	s.Get("https://httpbin.org/cookies/set/name/value", nil)
	resp, err := s.Get("https://httpbin.org/get", nil)
	if err != nil {
		log.Fatal(err)
	}
	var data struct {
		Headers struct{ Hello, Cookie string }
	}
	if err := resp.JSON(&data); err != nil {
		log.Fatal(err)
	}
	fmt.Println(data.Headers.Hello, data.Headers.Cookie)
	// world name=value
}
