package quickstart

import (
	"io"
	"log"
	"os"

	"github.com/ilgooz/mesg-go/application"
)

type Option func(*QuickStart)

type QuickStart struct {
	// app is a MESG application
	app *application.Application

	config Config

	// streams holds MESG's application streams
	streams []*application.Stream
	errC    chan error

	log       *log.Logger
	logOutput io.Writer
}

// Config holds the application configuration.
type Config struct {
	WebhookServiceID    string
	DiscordInvServiceID string
	LogServiceID        string
	SendgridKey         string
	Email               string
}

// New creates a new QuickStart application.
func New(app *application.Application, config Config, options ...Option) *QuickStart {
	q := &QuickStart{
		app:       app,
		config:    config,
		errC:      make(chan error, 0),
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

// Start starts quickstart application.
func (q *QuickStart) Start() error {
	defer q.app.Close()

	q.monitor(q.whenRequest())
	q.monitor(q.whenDiscordSend())

	return q.wait()
}

func (q *QuickStart) monitor(stream *application.Stream, err error) {
	if err != nil {
		q.errC <- err
		return
	}
	q.streams = append(q.streams, stream)
}

func (q *QuickStart) wait() error {
	select {
	case err := <-q.errC:
		return err
	default:
	}

	errC := make(chan error, 0)

	for _, stream := range q.streams {
		go q.monitorStream(stream, errC)
	}

	return <-errC
}

func (q *QuickStart) monitorStream(stream *application.Stream, errC chan error) {
	for {
		select {
		case err := <-stream.Err:
			errC <- err
			return

		case execution := <-stream.Executions:
			if execution.Err != nil {
				q.log.Println(execution.Err)
			}
		}
	}
}

// Close gracefully closes quickstart application.
func (q *QuickStart) Close() error {
	return q.app.Close()
}
