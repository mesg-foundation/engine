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
	// st is storage for workflows.
	st Storage

	// timeout used to set timeouts for external requests.
	timeout time.Duration

	// s is mesg service.
	s *mesg.Service

	// cp is a core client provider.
	cp coreClientProvider

	// vm that runs workflows.
	vm *VM
}

// New returns a new Workflow with given storage st and options.
func New(coreAddr string, st Storage, options ...Option) (*Workflow, error) {
	w := &Workflow{
		timeout: 5 * time.Second,
		st:      st,
	}
	for _, option := range options {
		option(w)
	}
	if w.s == nil {
		s, err := mesg.New()
		if err != nil {
			return nil, err
		}
		w.s = s
	}

	var cp coreClientProvider
	if w.cp != nil {
		cp = w.cp
	} else {
		cp = &defaultCoreClientProvider{timeout: w.timeout}
	}

	core, err := cp.New(coreAddr)
	if err != nil {
		return nil, err
	}
	w.vm = newVM(core)
	return w, nil
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
	if err := w.runWorkflows(); err != nil {
		return err
	}
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
	if err := w.s.Close(); err != nil {
		return err
	}
	return w.vm.TerminateAll()
}

// runWorkflows runs all workflows in the storage.
func (w *Workflow) runWorkflows() error {
	workflows, err := w.st.All()
	if err != nil {
		return err
	}
	for _, workflow := range workflows {
		if err := w.vm.Run(workflow); err != nil {
			return err
		}
	}
	return nil
}
