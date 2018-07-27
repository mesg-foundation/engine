package logger

import (
	"io"
	"log"
	"os"

	"github.com/ilgooz/mesg-go/service"
)

type Option func(*Logger)

// Logger is a logger service.
type Logger struct {
	service   *service.Service
	log       *log.Logger
	logOutput io.Writer
}

// New creates a new Logger runs over service s.
func New(service *service.Service, options ...Option) *Logger {
	l := &Logger{
		service:   service,
		logOutput: os.Stdout,
	}
	for _, option := range options {
		option(l)
	}
	l.log = log.New(l.logOutput, "logger", log.LstdFlags)
	return l
}

// LogOutputOption uses out as a log destination.
func LogOutputOption(out io.Writer) Option {
	return func(l *Logger) {
		l.logOutput = out
	}
}

// Start starts logger as a service.
func (l *Logger) Start() error {
	return l.service.Listen(
		service.NewTask("log", l.handler),
	)
}

// Close closes the service.
func (l *Logger) Close() error {
	return l.service.Close()
}
