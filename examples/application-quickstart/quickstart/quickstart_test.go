package quickstart

import (
	"io/ioutil"
	"sync"
	"testing"

	"github.com/ilgooz/mesg-go/application"
	"github.com/ilgooz/mesg-go/application/applicationtest"
	"github.com/stvp/assert"
)

var config = Config{
	WebhookServiceID:    "x1",
	DiscordInvServiceID: "x2",
	LogServiceID:        "x3",
	SendgridKey:         "k1",
	Email:               "e1",
}

func newApplicationAndServer(t *testing.T) (*application.Application, *applicationtest.Server) {
	testServer := applicationtest.NewServer()
	application, err := application.New(
		application.DialOption(testServer.Socket()),
	)
	assert.Nil(t, err)
	assert.NotNil(t, application)
	return application, testServer
}

func TestWhenRequest(t *testing.T) {
	app, server := newApplicationAndServer(t)
	go server.Start()

	l := New(app, config, LogOutputOption(ioutil.Discard))
	go l.Start()

	assert.Nil(t, server.EmitEvent(config.WebhookServiceID, "request", nil))
	le := server.LastExecute()
	assert.Equal(t, config.DiscordInvServiceID, le.ServiceID())
	assert.Equal(t, "send", le.Task())

	var data sendgridRequest
	assert.Nil(t, le.Decode(&data))
	assert.Equal(t, config.SendgridKey, data.SendgridAPIKey)
	assert.Equal(t, config.Email, data.Email)
}

type logData struct {
	Info string `json:"info"`
}

func TestWhenResult(t *testing.T) {
	ldata := logData{"awesome log data"}

	app, server := newApplicationAndServer(t)
	go server.Start()

	l := New(app, config, LogOutputOption(ioutil.Discard))
	go l.Start()

	assert.Nil(t, server.EmitResult(config.DiscordInvServiceID, "send", "success", ldata))
	le := server.LastExecute()
	assert.Equal(t, config.LogServiceID, le.ServiceID())
	assert.Equal(t, "log", le.Task())

	var data logRequest
	assert.Nil(t, le.Decode(&data))
	assert.Equal(t, config.DiscordInvServiceID, data.ServiceID)
	assert.Equal(t, ldata.Info, data.Data.(map[string]interface{})["info"])
}
func TestClose(t *testing.T) {
	app, server := newApplicationAndServer(t)
	go server.Start()

	l := New(app, config, LogOutputOption(ioutil.Discard))

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		assert.NotNil(t, l.Start())
	}()

	assert.Nil(t, server.EmitEvent(config.WebhookServiceID, "request", nil))
	server.LastExecute()

	assert.Nil(t, l.Close())
	wg.Wait()
}
