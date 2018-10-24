// Package workflow is Workflow System Service for managing and running workflows.
package workflow

import (
	"time"

	mesg "github.com/mesg-foundation/go-service"
)

const (
	createTaskKey = "create"
	deleteTaskKey = "delete"
)

// Workflow is Workflow System Service for managing and running workflows.
type Workflow struct {
	// timeout used to set timeouts for external requests.
	timeout time.Duration

	// s is mesg service.
	s *mesg.Service

	// cp is a core client provider.
	cp coreClientProvider
}

// New returns a new Workflow.
func New(options ...Option) (*Workflow, error) {
	r := &Workflow{
		timeout: 5 * time.Second,
	}
	for _, option := range options {
		option(r)
	}
	if r.s == nil {
		var err error
		r.s, err = mesg.New()
		if err != nil {
			return nil, err
		}
	}
	if r.cp == nil {
		r.cp = &defaultCoreClientProvider{timeout: r.timeout}
	}
	return r, nil
}

// Option is a configuration func for WSS.
type Option func(*Workflow)

// mesgOption returns an option for setting mesg service s.
func mesgOption(s *mesg.Service) Option {
	return func(r *Workflow) {
		r.s = s
	}
}

// coreClientProviderOption returns an option for setting core client provider.
func coreClientProviderOption(p coreClientProvider) Option {
	return func(w *Workflow) {
		w.cp = p
	}
}

// Start starts WSS.
func (w *Workflow) Start() error {
	return w.listenTasks()
}

func (w *Workflow) listenTasks() error {
	return w.s.Listen(
		mesg.Task(createTaskKey, w.createHandler),
		mesg.Task(deleteTaskKey, w.deleteHandler),
	)
}

// Close gracefully closes WSS.
func (w *Workflow) Close() error {
	return w.s.Close()
}
