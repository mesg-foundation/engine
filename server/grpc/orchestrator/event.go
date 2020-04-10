package orchestrator

import (
	"github.com/mesg-foundation/engine/event"
	"github.com/mesg-foundation/engine/event/publisher"
	"github.com/mesg-foundation/engine/protobuf/acknowledgement"
)

type eventServer struct {
	eventPublisher *publisher.EventPublisher
	auth           *Authorizer
}

// NewEventServer creates a new Event Server.
func NewEventServer(eventPublisher *publisher.EventPublisher, auth *Authorizer) EventServer {
	return &eventServer{
		eventPublisher: eventPublisher,
		auth:           auth,
	}
}

// Stream returns stream of events.
func (s *eventServer) Stream(req *EventStreamRequest, stream Event_StreamServer) error {
	// check authorization
	if err := s.auth.IsAuthorized(stream.Context(), req); err != nil {
		return err
	}

	var f *event.Filter
	if req.Filter != nil {
		f = &event.Filter{
			Hash:         req.Filter.Hash,
			InstanceHash: req.Filter.InstanceHash,
			Key:          req.Filter.Key,
		}
	}
	eventStream := s.eventPublisher.GetStream(f)
	defer eventStream.Close()

	// send header to notify client that the stream is ready.
	if err := acknowledgement.SetStreamReady(stream); err != nil {
		return err
	}

	for {
		select {
		case event := <-eventStream.C:
			if err := stream.Send(event); err != nil {
				return err
			}
		case <-stream.Context().Done():
			return stream.Context().Err()
		}
	}
}
