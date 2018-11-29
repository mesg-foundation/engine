package api

import (
	"github.com/mesg-foundation/core/execution"
	"github.com/mesg-foundation/core/pubsub"
	"github.com/mesg-foundation/core/service"
)

// ListenTask listens tasks on service token.
func (a *API) ListenTask(token string) (*TaskListener, error) {
	l := newTaskListener(a)
	return l, l.listen(token)
}

// TaskListener provides functionalities to listen MESG tasks.
type TaskListener struct {
	// Executions receives matching executions for tasks.
	Executions chan *execution.Execution

	// Err filled when task subscription finished with a failure.
	Err chan error

	// cancel stops listening for tasks.
	cancel chan struct{}

	// listening indicates if listening started
	listening chan struct{}

	api *API
}

// newTaskListener creates a new TaskListener with given api.
func newTaskListener(api *API) *TaskListener {
	return &TaskListener{
		Executions: make(chan *execution.Execution),
		Err:        make(chan error, 1),
		cancel:     make(chan struct{}),
		listening:  make(chan struct{}),
		api:        api,
	}
}

// Close stops listening for tasks.
func (l *TaskListener) Close() error {
	close(l.cancel)
	return nil
}

// listen listens tasks matches with service token.
func (l *TaskListener) listen(token string) error {
	s, err := l.api.db.Get(token)
	if err != nil {
		return err
	}
	s, err = service.FromService(s, service.ContainerOption(l.api.container))
	if err != nil {
		return err
	}
	go l.listenLoop(s)
	<-l.listening
	return nil
}

func (l *TaskListener) listenLoop(service *service.Service) {
	channel := service.TaskSubscriptionChannel()
	subscription := pubsub.Subscribe(channel)
	defer pubsub.Unsubscribe(channel, subscription)
	close(l.listening)

	for {
		select {
		case <-l.cancel:
			return

		// TODO use e.Err when subscription fails.
		// currently we don't need this but when pubsub refactored, we'll
		// need to pass an error to Err chan.
		case data := <-subscription:
			execution := data.(*execution.Execution)
			l.Executions <- execution
		}
	}
}
