package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/mesg-foundation/core/event"
	"github.com/mesg-foundation/core/hash"
	"github.com/mesg-foundation/core/protobuf/acknowledgement"
	"github.com/mesg-foundation/core/protobuf/api"
	"github.com/mesg-foundation/core/protobuf/definition"
	"github.com/mesg-foundation/core/sdk"
	eventsdk "github.com/mesg-foundation/core/sdk/event"
)

// EventServer serve event functions.
type EventServer struct {
	sdk *sdk.SDK
}

// NewEventServer creates a new EventServer.
func NewEventServer(sdk *sdk.SDK) *EventServer {
	return &EventServer{sdk: sdk}
}

// Create creates a new event.
func (s *EventServer) Create(ctx context.Context, req *api.CreateEventRequest) (*api.CreateEventResponse, error) {
	if req.Event.Hash != "" {
		return nil, errors.New("create event: hash not allowed")
	}

	if req.Event.Key == "" {
		return nil, errors.New("create event: key missing")
	}

	instanceHash, err := hash.Decode(req.Event.InstanceHash)
	if err != nil {
		return nil, fmt.Errorf("create event: instance %s", err)
	}

	var data map[string]interface{}
	if err := json.Unmarshal([]byte(req.Event.Data), &data); err != nil {
		return nil, fmt.Errorf("create event: data %s", err)
	}

	event, err := s.sdk.Event.Emit(instanceHash, req.Event.Key, data)
	if err != nil {
		return nil, fmt.Errorf("create event: data %s", err)
	}

	return &api.CreateEventResponse{Hash: event.Hash.String()}, nil
}

// Stream returns stream of events.
func (s *EventServer) Stream(req *api.StreamEventRequest, resp api.Event_StreamServer) error {
	var f *eventsdk.Filter
	if req.Filter != nil {
		instanceHash, err := hash.Decode(req.Filter.InstanceHash)
		if req.Filter.InstanceHash != "" && err != nil {
			return err
		}
		hash, err := hash.Decode(req.Filter.Hash)
		if req.Filter.Hash != "" && err != nil {
			return err
		}
		f = &eventsdk.Filter{
			Hash:         hash,
			InstanceHash: instanceHash,
			Key:          req.Filter.Key,
		}
	}
	stream := s.sdk.Event.GetStream(f)
	defer stream.Close()

	// send header to notify client that the stream is ready.
	if err := acknowledgement.SetStreamReady(resp); err != nil {
		return err
	}

	for exec := range stream.C {
		pexec, err := toProtoEvent(exec)
		if err != nil {
			return err
		}

		if err := resp.Send(pexec); err != nil {
			return err
		}
	}

	return nil
}

func toProtoEvent(e *event.Event) (*definition.Event, error) {
	data, err := json.Marshal(e.Data)
	if err != nil {
		return nil, err
	}

	return &definition.Event{
		Hash:         e.Hash.String(),
		InstanceHash: e.InstanceHash.String(),
		Key:          e.Key,
		Data:         string(data),
	}, nil
}
