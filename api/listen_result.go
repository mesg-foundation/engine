package api

import (
	"github.com/mesg-foundation/core/execution"
	"github.com/mesg-foundation/core/pubsub"
	"github.com/mesg-foundation/core/service"
	"github.com/mesg-foundation/core/x/xstrings"
)

// ListenResult listens results matches that with filters on serviceID.
func (a *API) ListenResult(serviceID string, filters ...ListenResultFilter) (*ResultListener, error) {
	l := &ResultListener{
		Executions: make(chan *execution.Execution),
		Err:        make(chan error, 1),
		cancel:     make(chan struct{}),
		listening:  make(chan struct{}),
		api:        a,
	}
	for _, filter := range filters {
		filter(l)
	}
	return l, l.listen(serviceID)
}

// ResultListener provides functionalities to listen MESG results.
type ResultListener struct {
	// Executions receives matching executions for results.
	Executions chan *execution.Execution

	// Err filled when result subscription finished with a failure.
	Err chan error

	// cancel stops listening for new results.
	cancel chan struct{}

	// listening indicates if listening started
	listening chan struct{}

	// filters.
	taskKey    string
	outputKey  string
	tagFilters []string

	api *API
}

// ListenResultFilter is a filter func for filtering results.
type ListenResultFilter func(*ResultListener)

// ListenResultTaskFilter returns a taskKey filter.
func ListenResultTaskFilter(taskKey string) ListenResultFilter {
	return func(ln *ResultListener) {
		ln.taskKey = taskKey
	}
}

// ListenResultOutputFilter returns an outputKey filter.
func ListenResultOutputFilter(outputKey string) ListenResultFilter {
	return func(ln *ResultListener) {
		ln.outputKey = outputKey
	}
}

// ListenResultTagFilters returns a tags filter.
func ListenResultTagFilters(tags []string) ListenResultFilter {
	return func(ln *ResultListener) {
		ln.tagFilters = tags
	}
}

// Close stops listening for results.
func (l *ResultListener) Close() error {
	close(l.cancel)
	return nil
}

// listen listens results matches with filters on serviceID.
func (l *ResultListener) listen(serviceID string) error {
	s, err := l.api.db.Get(serviceID)
	if err != nil {
		return err
	}
	s, err = service.FromService(s, service.ContainerOption(l.api.container))
	if err != nil {
		return err
	}
	if err := l.validateTaskKey(s); err != nil {
		return err
	}
	go l.listenLoop(s)
	<-l.listening
	return nil
}

func (l *ResultListener) listenLoop(service *service.Service) {
	channel := service.ResultSubscriptionChannel()
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
			if l.isSubscribed(execution) {
				l.Executions <- execution
			}
		}
	}
}

func (l *ResultListener) validateTaskKey(service *service.Service) error {
	if l.taskKey == "" || l.taskKey == "*" {
		return nil
	}
	_, err := service.GetTask(l.taskKey)
	return err
}

func (l *ResultListener) isSubscribed(e *execution.Execution) bool {
	return l.isSubscribedToTags(e) &&
		l.isSubscribedToTask(e)
}

func (l *ResultListener) isSubscribedToTask(e *execution.Execution) bool {
	return xstrings.SliceContains([]string{"", "*", e.TaskKey}, l.taskKey)
}

func (l *ResultListener) isSubscribedToTags(e *execution.Execution) bool {
	for _, tag := range l.tagFilters {
		if !xstrings.SliceContains(e.Tags, tag) {
			return false
		}
	}
	return true
}
