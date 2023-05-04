# GoHTTP

[![GoDev](https://img.shields.io/static/v1?label=godev&message=reference&color=00add8)][godev]
[![Go](https://github.com/sunshineplan/gohttp/workflows/Test/badge.svg)][actions]
[![CoverageStatus](https://coveralls.io/repos/github/sunshineplan/gohttp/badge.svg?branch=main&service=github)][coveralls]
[![GoReportCard](https://goreportcard.com/badge/github.com/sunshineplan/gohttp)][goreportcard]

[godev]: https://pkg.go.dev/github.com/sunshineplan/gohttp "GoDev"
[actions]: https://github.com/sunshineplan/gohttp/actions "GitHub Actions Page"
[coveralls]: https://coveralls.io/github/sunshineplan/gohttp?branch=main "Coverage Status"
[goreportcard]: https://goreportcard.com/report/github.com/sunshineplan/gohttp "Go Report Card"

Package GoHTTP is an elegant and simple HTTP library for Go.

## Installation

    go get -u github.com/sunshineplan/gohttp

## Documentation

https://pkg.go.dev/github.com/sunshineplan/gohttp

## License

[The MIT License (MIT)](https://raw.githubusercontent.com/sunshineplan/gohttp/main/LICENSE)

## Usage examples

A few usage examples can be found below. See the documentation for the full list of supported functions.

### HTTP request

```go
// HTTP GET request
r, _ := gohttp.Get("https://api.github.com/user", gohttp.H{"Authorization": "token"})
fmt.Print(r.StatusCode) // 200
fmt.Print(r.Header.Get("content-type")) // application/json; charset=utf-8
fmt.Print(r.String()) // {"type":"User"...

// HTTP POST request
r, _ = gohttp.Post("https://httpbin.org/post", nil, url.Values{"hello": []string{"world"}})
var data struct { Form struct{ Hello string } }
r.JSON(&data)
fmt.Println(data.Form.Hello) // world

// Upload File
r, _ := gohttp.Upload("https://httpbin.org/post", nil, nil, gohttp.F("readme", "README.md"))
var resp struct {
    Files   struct{ Readme string }
    Headers struct {
        ContentType string `json:"Content-Type"`
    }
}
r.JSON(&resp)
fmt.Println(strings.Split(resp.Files.Readme, "\r\n")[0])     // # GoHTTP
fmt.Println(strings.Split(resp.Headers.ContentType, ";")[0]) // multipart/form-data
```

### Session

```go
// Session provides cookie persistence and configuration
s := gohttp.NewSession()
s.Header.Set("hello", "world")
s.Get("https://httpbin.org/cookies/set/name/value", nil)
r, _ := s.Get("https://httpbin.org/get", nil)
var data struct { Headers struct{ Hello, Cookie string } }
r.JSON(&data)
fmt.Println(data.Headers.Hello)  // world
fmt.Println(data.Headers.Cookie) // name=value
```
