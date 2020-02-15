package publisher

import (
	"github.com/cskr/pubsub"
	"github.com/mesg-foundation/engine/cosmos"
	"github.com/mesg-foundation/engine/event"
	"github.com/mesg-foundation/engine/hash"
	"github.com/mesg-foundation/engine/protobuf/types"
)

const (
	// streamTopic is topic used to broadcast events.
	streamTopic = "event-stream"
)

// EventPublisher exposes event APIs of MESG.
type EventPublisher struct {
	ps *pubsub.PubSub
	mc *cosmos.ModuleClient
}

// New creates a new Event SDK with given options.
func New(mc *cosmos.ModuleClient) *EventPublisher {
	return &EventPublisher{
		ps: pubsub.New(0),
		mc: mc,
	}
}

// Publish a MESG event eventKey with eventData for service token.
func (ep *EventPublisher) Publish(instanceHash hash.Hash, eventKey string, eventData *types.Struct) (*event.Event, error) {
	i, err := ep.mc.GetInstance(instanceHash)
	if err != nil {
		return nil, err
	}

	s, err := ep.mc.GetService(i.ServiceHash)
	if err != nil {
		return nil, err
	}

	if err := s.RequireEventData(eventKey, eventData); err != nil {
		return nil, err
	}

	e := event.New(instanceHash, eventKey, eventData)
	go ep.ps.Pub(e, streamTopic)
	return e, nil
}

// GetStream broadcasts all events.
func (ep *EventPublisher) GetStream(f *event.Filter) *event.Listener {
	l := event.NewListener(ep.ps, streamTopic, f)
	go l.Listen()
	return l
}
