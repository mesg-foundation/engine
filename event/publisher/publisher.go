package publisher

import (
	"context"

	"github.com/cskr/pubsub"
	"github.com/mesg-foundation/engine/event"
	"github.com/mesg-foundation/engine/hash"
	"github.com/mesg-foundation/engine/instance"
	"github.com/mesg-foundation/engine/protobuf/types"
	"github.com/mesg-foundation/engine/service"
)

// Store is the interface to implement to fetch data.
type Store interface {
	// FetchService returns a service from its hash.
	FetchService(ctx context.Context, hash hash.Hash) (*service.Service, error)
	// FetchInstance returns an instance from its hash.
	FetchInstance(ctx context.Context, hash hash.Hash) (*instance.Instance, error)
}

const (
	// streamTopic is topic used to broadcast events.
	streamTopic = "event-stream"
)

// EventPublisher exposes event APIs of MESG.
type EventPublisher struct {
	ps    *pubsub.PubSub
	store Store
}

// New creates a new Event SDK with given options.
func New(store Store) *EventPublisher {
	return &EventPublisher{
		ps:    pubsub.New(0),
		store: store,
	}
}

// Publish a MESG event eventKey with eventData for service token.
func (ep *EventPublisher) Publish(instanceHash hash.Hash, eventKey string, eventData *types.Struct) (*event.Event, error) {
	i, err := ep.store.FetchInstance(context.Background(), instanceHash)
	if err != nil {
		return nil, err
	}

	s, err := ep.store.FetchService(context.Background(), i.ServiceHash)
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
