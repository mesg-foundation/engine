package api

import (
	"fmt"

	"github.com/mesg-foundation/core/database/services"
	"github.com/mesg-foundation/core/execution"
	"github.com/mesg-foundation/core/pubsub"
	"github.com/mesg-foundation/core/service"
	"github.com/mesg-foundation/core/utils/array"
)

// ResultListener provides functionalities to listen MESG results.
type ResultListener struct {
	// Executions receives matching executions for results.
	Executions chan *execution.Execution

	// Err filled when result subscription finished with a failure.
	Err chan error

	// cancel stops listening for new results.
	cancel chan struct{}

	// filters.
	taskKey    string
	outputKey  string
	tagFilters []string

	api *API
}

// newResultListener creates a new ResultListener with given api.
func newResultListener(api *API, filters ...ListenResultFilter) *ResultListener {
	ln := &ResultListener{
		Executions: make(chan *execution.Execution, 0),
		Err:        make(chan error, 1),
		cancel:     make(chan struct{}, 0),
		api:        api,
	}
	for _, filter := range filters {
		filter(ln)
	}
	return ln
}

// Close stops listening for results.
func (l *ResultListener) Close() error {
	close(l.cancel)
	return nil
}

// listen listens results matches with filters on serviceID.
func (l *ResultListener) listen(serviceID string) error {
	service, err := services.Get(serviceID)
	if err != nil {
		return err
	}
	if err := l.validateTask(&service); err != nil {
		return err
	}
	go l.listenLoop(&service)
	return nil
}

func (l *ResultListener) listenLoop(service *service.Service) {
	channel := service.ResultSubscriptionChannel()
	subscription := pubsub.Subscribe(channel)
	defer pubsub.Unsubscribe(channel, subscription)

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

func (l *ResultListener) validateTask(service *service.Service) error {
	if err := l.validateTaskKey(service); err != nil {
		return err
	}
	return l.validateOutputKey(service)
}

func (l *ResultListener) validateTaskKey(service *service.Service) error {
	if l.taskKey == "" || l.taskKey == "*" {
		return nil
	}
	if _, ok := service.Tasks[l.taskKey]; ok {
		return nil
	}
	return fmt.Errorf("Task %q doesn't exist in this service", l.taskKey)
}

func (l *ResultListener) validateOutputKey(service *service.Service) error {
	if l.outputKey == "" || l.outputKey == "*" {
		return nil
	}
	task, ok := service.Tasks[l.taskKey]
	if !ok {
		return fmt.Errorf("Task %q doesn't exist in this service", l.taskKey)
	}
	_, ok = task.Outputs[l.outputKey]
	if !ok {
		return fmt.Errorf("Output %q doesn't exist in the task %q of this service", l.outputKey, l.taskKey)
	}
	return nil
}

func (l *ResultListener) isSubscribed(e *execution.Execution) bool {
	return l.isSubscribedToTags(e) &&
		l.isSubscribedToTask(e) &&
		l.isSubscribedToOutput(e)
}

func (l *ResultListener) isSubscribedToTask(e *execution.Execution) bool {
	return array.IncludedIn([]string{"", "*", e.Task}, l.taskKey)
}

func (l *ResultListener) isSubscribedToOutput(e *execution.Execution) bool {
	return array.IncludedIn([]string{"", "*", e.Output}, l.outputKey)
}

func (l *ResultListener) isSubscribedToTags(e *execution.Execution) bool {
	for _, tag := range l.tagFilters {
		if !array.IncludedIn(e.Tags, tag) {
			return false
		}
	}
	return true
}
