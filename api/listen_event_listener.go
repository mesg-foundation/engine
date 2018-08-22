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

	api *API
}

// newEventListener creates a new EventListener with given api.
func newEventListener(api *API) *EventListener {
	return &EventListener{
		Events: make(chan *event.Event, 0),
		Err:    make(chan error, 1),
		cancel: make(chan struct{}, 0),
		api:    api,
	}
}

// Close stops listening for events.
func (e *EventListener) Close() error {
	close(e.cancel)
	return nil
}

// listen listens events matches with eventFilter on serviceID.
func (e *EventListener) listen(serviceID, eventFilter string) error {
	service, err := services.Get(serviceID)
	if err != nil {
		return err
	}
	if err := e.validateEventKey(&service, eventFilter); err != nil {
		return err
	}
	go e.listenLoop(&service, eventFilter)
	return nil
}

func (e *EventListener) listenLoop(service *service.Service, eventFilter string) {
	channel := service.EventSubscriptionChannel()
	subscription := pubsub.Subscribe(channel)
	defer pubsub.Unsubscribe(channel, subscription)

	for {
		select {
		case <-e.cancel:
			return

		// TODO use e.Err when subscription fails.
		// currently we don't need this but when pubsub refactored, we'll
		// need to pass an error to Err chan.
		case data := <-subscription:
			event := data.(*event.Event)
			if e.isSubscribedEvent(eventFilter, event) {
				e.Events <- event
			}
		}
	}
}

func (e *EventListener) validateEventKey(service *service.Service, eventKey string) error {
	if eventKey == "" || eventKey == "*" {
		return nil
	}
	if _, ok := service.Events[eventKey]; ok {
		return nil
	}
	return fmt.Errorf("Event %q doesn't exist in this service", eventKey)
}

func (e *EventListener) isSubscribedEvent(eventFilter string, ev *event.Event) bool {
	return array.IncludedIn([]string{"", "*", ev.Key}, eventFilter)
}
