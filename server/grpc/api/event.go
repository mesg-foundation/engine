package api

import (
	"context"
	"errors"
	"fmt"

	"github.com/mesg-foundation/engine/event"
	"github.com/mesg-foundation/engine/event/publisher"
	"github.com/mesg-foundation/engine/protobuf/acknowledgement"
	"github.com/mesg-foundation/engine/protobuf/api"
)

// EventServer serve event functions.
type EventServer struct {
	ep *publisher.EventPublisher
}

// NewEventServer creates a new EventServer.
func NewEventServer(ep *publisher.EventPublisher) *EventServer {
	return &EventServer{ep: ep}
}

// Create creates a new event.
func (s *EventServer) Create(ctx context.Context, req *api.CreateEventRequest) (*api.CreateEventResponse, error) {
	if req.Key == "" {
		return nil, errors.New("create event: key missing")
	}

	event, err := s.ep.Publish(req.InstanceHash, req.Key, req.Data)
	if err != nil {
		return nil, fmt.Errorf("create event: data %s", err)
	}

	return &api.CreateEventResponse{Hash: event.Hash}, nil
}

// Stream returns stream of events.
func (s *EventServer) Stream(req *api.StreamEventRequest, resp api.Event_StreamServer) error {
	var f *event.Filter
	if req.Filter != nil {
		f = &event.Filter{
			Hash:         req.Filter.Hash,
			InstanceHash: req.Filter.InstanceHash,
			Key:          req.Filter.Key,
		}
	}
	stream := s.ep.GetStream(f)
	defer stream.Close()

	// send header to notify client that the stream is ready.
	if err := acknowledgement.SetStreamReady(resp); err != nil {
		return err
	}

	for event := range stream.C {
		if err := resp.Send(event); err != nil {
			return err
		}
	}

	return nil
}
