package quickstart

import (
	"io"
	"log"
	"os"

	"github.com/ilgooz/mesg-go/application"
)

type Option func(*QuickStart)

type QuickStart struct {
	app       *application.Application
	config    Config
	log       *log.Logger
	logOutput io.Writer
}

type Config struct {
	WebhookServiceID    string
	DiscordInvServiceID string
	LogServiceID        string
	SendgridKey         string
	Email               string
}

func New(app *application.Application, config Config, options ...Option) *QuickStart {
	q := &QuickStart{
		app:       app,
		config:    config,
		logOutput: os.Stdout,
	}
	for _, option := range options {
		option(q)
	}
	q.log = log.New(q.logOutput, "quick-start", log.LstdFlags)
	return q
}

// LogOutputOption uses out as a log destination.
func LogOutputOption(out io.Writer) Option {
	return func(q *QuickStart) {
		q.logOutput = out
	}
}

func (q *QuickStart) Start() error {
	defer q.app.Close()

	requestStream, err := q.whenRequest()
	if err != nil {
		return err
	}

	discordSendStream, err := q.whenDiscordSend()
	if err != nil {
		return err
	}

	for {
		select {
		case err := <-requestStream.Err:
			if err != nil {
				return err
			}

		case err := <-discordSendStream.Err:
			if err != nil {
				return err
			}

		case execution := <-requestStream.Executions:
			if execution.Err != nil {
				q.log.Println(execution.Err)
			}

		case execution := <-discordSendStream.Executions:
			if execution.Err != nil {
				q.log.Println(execution.Err)
			}
		}
	}
}

func (q *QuickStart) Close() error {
	return q.app.Close()
}
