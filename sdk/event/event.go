package eventsdk

import (
	"github.com/cskr/pubsub"
	"github.com/mesg-foundation/core/event"
	"github.com/mesg-foundation/core/hash"
	instancesdk "github.com/mesg-foundation/core/sdk/instance"
	servicesdk "github.com/mesg-foundation/core/sdk/service"
)

const (
	// streamTopic is topic used to broadcast events.
	streamTopic = "event-stream"
	topic       = "Event"
)

// Event exposes event APIs of MESG.
type Event struct {
	ps       *pubsub.PubSub
	instance *instancesdk.Instance
	service  *servicesdk.Service
}

// New creates a new Event SDK with given options.
func New(ps *pubsub.PubSub, service *servicesdk.Service, instance *instancesdk.Instance) *Event {
	return &Event{
		ps:       ps,
		service:  service,
		instance: instance,
	}
}

// Emit emits a MESG event eventKey with eventData for service token.
func (e *Event) Emit(instanceHash hash.Hash, eventKey string, eventData map[string]interface{}) (*event.Event, error) {
	instance, err := e.instance.Get(instanceHash)
	if err != nil {
		return nil, err
	}

	service, err := e.service.Get(instance.ServiceHash)
	if err != nil {
		return nil, err
	}

	if err := service.RequireEventData(eventKey, eventData); err != nil {
		return nil, err
	}

	event := event.Create(instanceHash, eventKey, eventData)
	go e.ps.Pub(event, streamTopic)
	go e.ps.Pub(event, subTopic(instanceHash))
	return event, nil
}

// GetStream broadcasts all events.
func (e *Event) GetStream(f *Filter) *Listener {
	l := NewListener(e.ps, streamTopic, f)
	go l.Listen()
	return l
}

// Listen listens events matches with eventFilter on serviceID.
func (e *Event) Listen(serviceHash hash.Hash, f *Filter) (*Listener, error) {
	s, err := e.service.Get(serviceHash)
	if err != nil {
		return nil, err
	}

	if f.HasKey() {
		if _, err := s.GetEvent(f.Key); err != nil {
			return nil, err
		}
	}

	l := NewListener(e.ps, subTopic(s.Hash), f)
	go l.Listen()
	return l, nil
}

// subTopic returns the topic to listen for events from this service.
func subTopic(serviceHash hash.Hash) string {
	return serviceHash.String() + "." + topic
}
