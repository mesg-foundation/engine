package mesg

import (
	"log"
	"testing"
)

func ExampleService(t *testing.T) {
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
