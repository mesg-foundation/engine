package quickstart

import "github.com/ilgooz/mesg-go/application"

func (q *QuickStart) whenRequest() (*application.Stream, error) {
	return q.app.
		WhenEvent(q.config.WebhookServiceID, application.EventFilterOption("request")).
		Map(sendgridRequest{
			Email:          q.config.Email,
			SendgridAPIKey: q.config.SendgridKey,
		}).
		Execute(q.config.DiscordInvServiceID, "send")
}

func (q *QuickStart) whenDiscordSend() (*application.Stream, error) {
	return q.app.
		WhenResult(q.config.DiscordInvServiceID, application.TaskFilterOption("send")).
		FilterFunc(func(r *application.Result) bool {
			var resp interface{}
			return r.Decode(&resp) == nil
		}).
		MapFunc(func(r *application.Result) application.Data {
			var resp interface{}
			r.Decode(&resp)
			return logRequest{
				ServiceID: q.config.DiscordInvServiceID,
				Data:      resp,
			}
		}).
		Execute(q.config.LogServiceID, "log")
}

type sendgridRequest struct {
	Email          string `json:"email"`
	SendgridAPIKey string `json:"sendgridAPIKey"`
}

type logRequest struct {
	ServiceID string      `json:"serviceID"`
	Data      interface{} `json:"data"`
}
