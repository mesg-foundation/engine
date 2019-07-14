package eventsdk

import (
	"github.com/cskr/pubsub"
	"github.com/mesg-foundation/engine/event"
	"github.com/mesg-foundation/engine/hash"
	instancesdk "github.com/mesg-foundation/engine/sdk/instance"
	servicesdk "github.com/mesg-foundation/engine/sdk/service"
)

const (
	// streamTopic is topic used to broadcast events.
	streamTopic = "event-stream"
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

// Create a MESG event eventKey with eventData for service token.
func (e *Event) Create(instanceHash hash.Hash, eventKey string, eventData map[string]interface{}) (*event.Event, error) {
	event := event.Create(instanceHash, eventKey, eventData)

	instance, err := e.instance.Get(event.InstanceHash)
	if err != nil {
		return nil, err
	}

	service, err := e.service.Get(instance.ServiceHash)
	if err != nil {
		return nil, err
	}

	if err := service.RequireEventData(event.Key, event.Data); err != nil {
		return nil, err
	}

	go e.ps.Pub(event, streamTopic)
	return event, nil
}

// CreateEngineEvent creates and publishes an engine event
func (e *Event) CreateEngineEvent(eventType event.EngineEventType, data map[string]interface{}) *event.Event {
	evt := event.EngineEvent(eventType, data)
	go e.ps.Pub(evt, streamTopic)
	return evt
}

// GetStream broadcasts all events.
func (e *Event) GetStream(f *Filter) *Listener {
	l := NewListener(e.ps, streamTopic, f)
	go l.Listen()
	return l
}
