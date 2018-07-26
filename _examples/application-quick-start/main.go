// Package main is an application that uses functionalities from following services:
// https://github.com/mesg-foundation/service-webhook
// https://github.com/mesg-foundation/service-discord-invitation
package main

import (
	"log"

	"github.com/ilgooz/mesg-go/application"
)

var (
	webhookServiceID    = "v1_3b56fdfc5bb5f50def11efaa6045046c"
	discordInvServiceID = "v1_e4e8d488444bbd9389ad9eac150669be"
	logServiceID        = "v1_0df94b3537841789f1dbe8bc8690a4c6"
	sendgridKey         = "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"
	email               = "ilkergoktugozturk@gmail.com"
)

func main() {
	app, err := application.New()
	if err != nil {
		log.Fatal(err)
	}

	resultStream, err := app.
		WhenResult(discordInvServiceID, application.TaskFilterOption("send")).
		FilterFunc(func(r *application.Result) bool {
			var resp interface{}
			return r.Decode(&resp) == nil
		}).
		MapFunc(func(r *application.Result) application.Data {
			var resp interface{}
			r.Decode(&resp)
			return logRequest{
				ServiceID: discordInvServiceID,
				Data:      resp,
			}
		}).
		Execute(logServiceID, "log")

	if err != nil {
		log.Fatal(err)
	}

	eventStream, err := app.
		WhenEvent(webhookServiceID, application.EventFilterOption("request")).
		Map(sendgridRequest{
			Email:          email,
			SendgridAPIKey: sendgridKey,
		}).
		Execute(discordInvServiceID, "send")

	if err != nil {
		log.Fatal(err)
	}

	for {
		select {
		case err := <-eventStream.Err:
			if err != nil {
				log.Fatal(err)
			}

		case err := <-resultStream.Err:
			if err != nil {
				log.Fatal(err)
			}

		case <-resultStream.Executions:
		}
	}
}

type sendgridRequest struct {
	Email          string `json:"email"`
	SendgridAPIKey string `json:"sendgridAPIKey"`
}

type logRequest struct {
	ServiceID string      `json:"serviceID"`
	Data      interface{} `json:"data"`
}
