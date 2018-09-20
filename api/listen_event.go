package api

import (
	"github.com/mesg-foundation/core/database/services"
	"github.com/mesg-foundation/core/event"
	"github.com/mesg-foundation/core/pubsub"
	"github.com/mesg-foundation/core/service"
	"github.com/mesg-foundation/core/x/xstrings"
)

// ListenEventFilter is a filter func for filtering events.
type ListenEventFilter func(*EventListener)

// ListenEventKeyFilter returns an eventKey filter.
func ListenEventKeyFilter(eventKey string) ListenEventFilter {
	return func(ln *EventListener) {
		ln.eventKey = eventKey
	}
}

// ListenEvent listens events matches with eventFilter on serviceID.
func (a *API) ListenEvent(serviceID string, filters ...ListenEventFilter) (*EventListener, error) {
	l := newEventListener(a, filters...)
	return l, l.listen(serviceID)
}

// EventListener provides functionalities to listen MESG events.
type EventListener struct {
	// Events receives matching events.
	Events chan *event.Event

	// Err filled when event subscription finished with a failure.
	Err chan error

	// cancel stops listening for new events.
	cancel chan struct{}

	// listening indicates if listening started
	listening chan struct{}

	// filters.
	eventKey string

	api *API
}

// newEventListener creates a new EventListener with given api and filters.
func newEventListener(api *API, filters ...ListenEventFilter) *EventListener {
	ln := &EventListener{
		Events:    make(chan *event.Event, 0),
		Err:       make(chan error, 1),
		cancel:    make(chan struct{}, 0),
		listening: make(chan struct{}, 0),
		api:       api,
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
	s, err := services.Get(serviceID)
	if err != nil {
		return err
	}
	s, err = service.FromService(s, service.ContainerOption(l.api.container))
	if err != nil {
		return err
	}
	if err := l.validateEventKey(s); err != nil {
		return err
	}
	go l.listenLoop(s)
	<-l.listening
	return nil
}

func (l *EventListener) listenLoop(service *service.Service) {
	channel := service.EventSubscriptionChannel()
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
	_, err := service.GetEvent(l.eventKey)
	return err
}

func (l *EventListener) isSubscribedEvent(e *event.Event) bool {
	return xstrings.SliceContains([]string{"", "*", e.Key}, l.eventKey)
}
