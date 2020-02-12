package eventsdk

import (
	"fmt"

	"github.com/cskr/pubsub"
	"github.com/mesg-foundation/engine/cosmos"
	"github.com/mesg-foundation/engine/event"
	"github.com/mesg-foundation/engine/hash"
	instancepb "github.com/mesg-foundation/engine/instance"
	"github.com/mesg-foundation/engine/protobuf/types"
	servicesdk "github.com/mesg-foundation/engine/sdk/service"
	"github.com/mesg-foundation/engine/x/instance"
)

const (
	// streamTopic is topic used to broadcast events.
	streamTopic = "event-stream"
)

// Event exposes event APIs of MESG.
type Event struct {
	ps      *pubsub.PubSub
	client  *cosmos.Client
	service *servicesdk.SDK
}

// New creates a new Event SDK with given options.
func New(ps *pubsub.PubSub, service *servicesdk.SDK, client *cosmos.Client) *Event {
	return &Event{
		ps:      ps,
		client:  client,
		service: service,
	}
}

// Create a MESG event eventKey with eventData for service token.
func (e *Event) Create(instanceHash hash.Hash, eventKey string, eventData *types.Struct) (*event.Event, error) {
	event := event.Create(instanceHash, eventKey, eventData)

	var inst instancepb.Instance
	if err := e.client.QueryJSON(fmt.Sprintf("custom/%s/%s/%s", instance.QuerierRoute, instance.QueryGetInstance, event.InstanceHash), nil, &inst); err != nil {
		return nil, err
	}

	service, err := e.service.Get(inst.ServiceHash)
	if err != nil {
		return nil, err
	}

	if err := service.RequireEventData(event.Key, event.Data); err != nil {
		return nil, err
	}

	go e.ps.Pub(event, streamTopic)
	return event, nil
}

// GetStream broadcasts all events.
func (e *Event) GetStream(f *Filter) *Listener {
	l := NewListener(e.ps, streamTopic, f)
	go l.Listen()
	return l
}
