package core

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"

	"github.com/mesg-foundation/core/config"
	"github.com/mesg-foundation/core/hash"
	"github.com/mesg-foundation/core/protobuf/acknowledgement"
	"github.com/mesg-foundation/core/protobuf/coreapi"
	"github.com/mesg-foundation/core/sdk"
	eventsdk "github.com/mesg-foundation/core/sdk/event"
	instancesdk "github.com/mesg-foundation/core/sdk/instance"
	"github.com/mesg-foundation/core/version"
)

// Server is the type to aggregate all the APIs.
type Server struct {
	sdk *sdk.SDK
}

// NewServer creates a new Server.
func NewServer(sdk *sdk.SDK) *Server {
	return &Server{sdk: sdk}
}

// ListenEvent listens events matches with eventFilter on serviceID.
func (s *Server) ListenEvent(request *coreapi.ListenEventRequest, stream coreapi.Core_ListenEventServer) error {
	hash, err := hash.Decode(request.ServiceID)
	if err != nil {
		return err
	}
	ln, err := s.sdk.Event.Listen(hash, &eventsdk.Filter{Key: request.EventFilter})
	if err != nil {
		return err
	}
	defer ln.Close()

	// send header to notify client that the stream is ready.
	if err := acknowledgement.SetStreamReady(stream); err != nil {
		return err
	}

	ctx := stream.Context()
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()

		case ev := <-ln.C:
			evData, err := json.Marshal(ev.Data)
			if err != nil {
				return err
			}

			if err := stream.Send(&coreapi.EventData{
				EventKey:  ev.Key,
				EventData: string(evData),
			}); err != nil {
				return err
			}
		}
	}
}

// ExecuteTask executes a task for a given service.
func (s *Server) ExecuteTask(ctx context.Context, request *coreapi.ExecuteTaskRequest) (*coreapi.ExecuteTaskReply, error) {
	hash, err := hash.Decode(request.ServiceID)
	if err != nil {
		return nil, err
	}

	instances, err := s.sdk.Instance.List(&instancesdk.Filter{ServiceHash: hash})
	if err != nil {
		return nil, err
	}
	if len(instances) != 1 {
		return nil, fmt.Errorf("you should have exactly one instance running to execute this service")
	}

	var inputs map[string]interface{}
	if err := json.Unmarshal([]byte(request.InputData), &inputs); err != nil {
		return nil, fmt.Errorf("cannot parse execution's inputs (JSON format): %s", err)
	}

	executionHash, err := s.sdk.Execution.Execute(hash, request.TaskKey, inputs, request.ExecutionTags)
	return &coreapi.ExecuteTaskReply{
		ExecutionHash: hex.EncodeToString(executionHash),
	}, err
}

// Info returns all necessary information from the core.
func (s *Server) Info(ctx context.Context, request *coreapi.InfoRequest) (*coreapi.InfoReply, error) {
	c, err := config.Global()
	if err != nil {
		return nil, err
	}
	services := make([]*coreapi.InfoReply_CoreService, len(c.Services()))
	for i, s := range c.Services() {
		services[i] = &coreapi.InfoReply_CoreService{
			Sid:  s.Sid,
			Hash: s.Hash.String(),
			Url:  s.URL,
			Key:  s.Key,
		}
	}
	return &coreapi.InfoReply{
		Version:  version.Version,
		Services: services,
	}, nil
}
