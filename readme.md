# mesg-go [![GoDoc](https://godoc.org/github.com/ilgooz/mesg-go?status.svg)](https://godoc.org/github.com/ilgooz/mesg-go) [![Go Report Card](https://goreportcard.com/badge/github.com/ilgooz/mesg-go)](https://goreportcard.com/report/github.com/ilgooz/mesg-go) [![Build Status](https://travis-ci.org/ilgooz/mesg-go.svg?branch=master)](https://travis-ci.org/ilgooz/mesg-go)
mesg-go is a service and application client for [mesg-core](https://github.com/mesg-foundation/core) with a high level API.
For more information please visit [mesg.com](https://mesg.com).

### Under Active Development
Breaking changes might be introduced untill the first release.

## Service example
```go
package main

import (
	"log"
)

func main() {
	srv, err := GetService()
	if err != nil {
		log.Fatal(err)
	}

	if err := srv.ListenTasks(
		NewTask("send", sendHandler),
	); err != nil {
		log.Fatal(err)
	}
}

func sendHandler(req *Request) {
	var data emailRequest
	if err := req.Get(&data); err != nil {
		log.Println(err)
		return
	}

	if err := req.Reply("success", successResponse{
		Code:    "202",
		Message: "the SendGrid message",
	}); err != nil {
		log.Println(err)
	}
}

type emailRequest struct {
	Email          string `json:"email"`
	SendgridAPIKey string `json:"sendgridAPIKey"`
}

type successResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type errorResponse struct {
	Message string `json:"message"`
}
```