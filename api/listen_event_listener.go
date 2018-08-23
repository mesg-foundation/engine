package api

import (
	"fmt"

	"github.com/mesg-foundation/core/database/services"
	"github.com/mesg-foundation/core/event"
	"github.com/mesg-foundation/core/pubsub"
	"github.com/mesg-foundation/core/service"
	"github.com/mesg-foundation/core/utils/array"
)

// EventListener provides functionalities to listen MESG events.
type EventListener struct {
	// Events receives matching events.
	Events chan *event.Event

	// Err filled when event subscription finished with a failure.
	Err chan error

	// cancel stops listening for new events.
	cancel chan struct{}

	// filters.
	eventKey string

	api *API
}

// newEventListener creates a new EventListener with given api and filters.
func newEventListener(api *API, filters ...ListenEventFilter) *EventListener {
	ln := &EventListener{
		Events: make(chan *event.Event, 0),
		Err:    make(chan error, 1),
		cancel: make(chan struct{}, 0),
		api:    api,
	}
	for _, filter := range filters {
		filter(ln)
	}
	return ln
}

// Close stops listening for events.
func (l *EventListener) Close() error {
	close(l.cancel)
	return nil
}

// listen listens events matches with eventFilter on serviceID.
func (l *EventListener) listen(serviceID string) error {
	service, err := services.Get(serviceID)
	if err != nil {
		return err
	}
	if err := l.validateEventKey(&service); err != nil {
		return err
	}
	go l.listenLoop(&service)
	return nil
}

func (l *EventListener) listenLoop(service *service.Service) {
	channel := service.EventSubscriptionChannel()
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
			event := data.(*event.Event)
			if l.isSubscribedEvent(event) {
				l.Events <- event
			}
		}
	}
}

func (l *EventListener) validateEventKey(service *service.Service) error {
	if l.eventKey == "" || l.eventKey == "*" {
		return nil
	}
	if _, ok := service.Events[l.eventKey]; ok {
		return nil
	}
	return fmt.Errorf("Event %q doesn't exist in this service", l.eventKey)
}

func (l *EventListener) isSubscribedEvent(e *event.Event) bool {
	return array.IncludedIn([]string{"", "*", e.Key}, l.eventKey)
}
