package publisher

import (
	"fmt"

	"github.com/cskr/pubsub"
	"github.com/mesg-foundation/engine/cosmos"
	"github.com/mesg-foundation/engine/event"
	"github.com/mesg-foundation/engine/hash"
	"github.com/mesg-foundation/engine/instance"
	"github.com/mesg-foundation/engine/protobuf/types"
	"github.com/mesg-foundation/engine/service"
	instancemodule "github.com/mesg-foundation/engine/x/instance"
	servicemodule "github.com/mesg-foundation/engine/x/service"
)

const (
	// streamTopic is topic used to broadcast events.
	streamTopic = "event-stream"
)

// EventPublisher exposes event APIs of MESG.
type EventPublisher struct {
	ps  *pubsub.PubSub
	rpc *cosmos.RPC
}

// New creates a new Event SDK with given options.
func New(rpc *cosmos.RPC) *EventPublisher {
	return &EventPublisher{
		ps:  pubsub.New(0),
		rpc: rpc,
	}
}

// Publish a MESG event eventKey with eventData for service token.
func (ep *EventPublisher) Publish(instanceHash hash.Hash, eventKey string, eventData *types.Struct) (*event.Event, error) {
	var i *instance.Instance
	route := fmt.Sprintf("custom/%s/%s/%s", instancemodule.QuerierRoute, instancemodule.QueryGet, instanceHash)
	if err := ep.rpc.QueryJSON(route, nil, &i); err != nil {
		return nil, err
	}

	var s *service.Service
	route = fmt.Sprintf("custom/%s/%s/%s", servicemodule.QuerierRoute, servicemodule.QueryGet, i.ServiceHash)
	if err := ep.rpc.QueryJSON(route, nil, &s); err != nil {
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
