package api

import (
	"context"
	"errors"
	"fmt"

	structpb "github.com/golang/protobuf/ptypes/struct"
	"github.com/mesg-foundation/engine/event"
	"github.com/mesg-foundation/engine/hash"
	"github.com/mesg-foundation/engine/protobuf/acknowledgement"
	"github.com/mesg-foundation/engine/protobuf/api"
	"github.com/mesg-foundation/engine/protobuf/convert"
	"github.com/mesg-foundation/engine/protobuf/types"
	"github.com/mesg-foundation/engine/sdk"
	eventsdk "github.com/mesg-foundation/engine/sdk/event"
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
	instanceHash, err := hash.Decode(req.InstanceHash)
	if err != nil {
		return nil, fmt.Errorf("create event: instance %s", err)
	}

	if req.Key == "" {
		return nil, errors.New("create event: key missing")
	}

	data := make(map[string]interface{})
	if err := convert.Marshal(req.Data, &data); err != nil {
		return nil, err
	}

	event, err := s.sdk.Event.Create(instanceHash, req.Key, data)
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

	for event := range stream.C {
		e, err := toProtoEvent(event)
		if err != nil {
			return err
		}

		if err := resp.Send(e); err != nil {
			return err
		}
	}

	return nil
}

func toProtoEvent(e *event.Event) (*types.Event, error) {
	data := &structpb.Struct{}
	if err := convert.Unmarshal(e.Data, data); err != nil {
		return nil, err
	}

	return &types.Event{
		Hash:         e.Hash.String(),
		InstanceHash: e.InstanceHash.String(),
		Key:          e.Key,
		Data:         data,
	}, nil
}
