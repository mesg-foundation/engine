package core

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/mesg-foundation/core/config"
	"github.com/mesg-foundation/core/execution"
	"github.com/mesg-foundation/core/protobuf/acknowledgement"
	"github.com/mesg-foundation/core/protobuf/coreapi"
	"github.com/mesg-foundation/core/sdk"
	eventsdk "github.com/mesg-foundation/core/sdk/event"
	executionsdk "github.com/mesg-foundation/core/sdk/execution"
	"github.com/mesg-foundation/core/server/grpc/api"
	"github.com/mesg-foundation/core/service"
	"github.com/mesg-foundation/core/version"
	"github.com/mesg-foundation/core/x/xerrors"
	"github.com/mr-tron/base58"
)

// Server is the type to aggregate all the APIs.
type Server struct {
	sdk *sdk.SDK
}

// NewServer creates a new Server.
func NewServer(sdk *sdk.SDK) *Server {
	return &Server{sdk: sdk}
}

// GetService returns service serviceID.
func (s *Server) GetService(ctx context.Context, request *coreapi.GetServiceRequest) (*coreapi.GetServiceReply, error) {
	hash, err := base58.Decode(request.ServiceID)
	if err != nil {
		return nil, err
	}

	ss, err := s.sdk.GetService(hash)
	if err != nil {
		return nil, err
	}
	status, err := s.sdk.Status(ss)
	if err != nil {
		return nil, err
	}
	details := &coreapi.Service{
		Definition: api.ToProtoService(ss),
		Status:     toProtoServiceStatusType(status),
	}
	return &coreapi.GetServiceReply{Service: details}, nil
}

// ListServices lists services.
func (s *Server) ListServices(ctx context.Context, request *coreapi.ListServicesRequest) (*coreapi.ListServicesReply, error) {
	services, err := s.sdk.ListServices()
	if err != nil {
		return nil, err
	}

	var (
		protoServices []*coreapi.Service
		mp            sync.Mutex

		servicesLen = len(services)
		errC        = make(chan error, servicesLen)
		wg          sync.WaitGroup
	)

	wg.Add(servicesLen)
	for _, ss := range services {
		go func(ss *service.Service) {
			defer wg.Done()
			status, err := s.sdk.Status(ss)
			if err != nil {
				errC <- err
				return
			}
			details := &coreapi.Service{
				Definition: api.ToProtoService(ss),
				Status:     toProtoServiceStatusType(status),
			}
			mp.Lock()
			protoServices = append(protoServices, details)
			mp.Unlock()
		}(ss)
	}

	wg.Wait()
	close(errC)

	var errs xerrors.Errors
	for err := range errC {
		errs = append(errs, err)
	}

	return &coreapi.ListServicesReply{Services: protoServices}, errs.ErrorOrNil()
}

// StartService starts a service.
func (s *Server) StartService(ctx context.Context, request *coreapi.StartServiceRequest) (*coreapi.StartServiceReply, error) {
	hash, err := base58.Decode(request.ServiceID)
	if err != nil {
		return nil, err
	}
	return &coreapi.StartServiceReply{}, s.sdk.StartService(hash)
}

// StopService stops a service.
func (s *Server) StopService(ctx context.Context, request *coreapi.StopServiceRequest) (*coreapi.StopServiceReply, error) {
	hash, err := base58.Decode(request.ServiceID)
	if err != nil {
		return nil, err
	}
	return &coreapi.StopServiceReply{}, s.sdk.StopService(hash)
}

// DeleteService stops and deletes service serviceID.
func (s *Server) DeleteService(ctx context.Context, request *coreapi.DeleteServiceRequest) (*coreapi.DeleteServiceReply, error) {
	hash, err := base58.Decode(request.ServiceID)
	if err != nil {
		return nil, err
	}
	return &coreapi.DeleteServiceReply{}, s.sdk.DeleteService(hash, request.DeleteData)
}

// ListenEvent listens events matches with eventFilter on serviceID.
func (s *Server) ListenEvent(request *coreapi.ListenEventRequest, stream coreapi.Core_ListenEventServer) error {
	hash, err := base58.Decode(request.ServiceID)
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

// ListenResult listens for results from a services.
func (s *Server) ListenResult(request *coreapi.ListenResultRequest, stream coreapi.Core_ListenResultServer) error {
	hash, err := base58.Decode(request.ServiceID)
	if err != nil {
		return err
	}

	filter := &executionsdk.Filter{
		Statuses: []execution.Status{
			execution.Completed,
			execution.Failed,
		},
		TaskKey: request.TaskFilter,
		Tags:    request.TagFilters,
	}

	ln, err := s.sdk.Execution.Listen(hash, filter)
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

		case execution := <-ln.C:
			outputs, err := json.Marshal(execution.Outputs)
			if err != nil {
				return err
			}
			if err := stream.Send(&coreapi.ResultData{
				ExecutionHash: hex.EncodeToString(execution.Hash),
				TaskKey:       execution.TaskKey,
				OutputData:    string(outputs),
				ExecutionTags: execution.Tags,
				Error:         execution.Error,
			}); err != nil {
				return err
			}
		}
	}
}

// ExecuteTask executes a task for a given service.
func (s *Server) ExecuteTask(ctx context.Context, request *coreapi.ExecuteTaskRequest) (*coreapi.ExecuteTaskReply, error) {
	hash, err := base58.Decode(request.ServiceID)
	if err != nil {
		return nil, err
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
			Hash: base58.Encode(s.Hash),
			Url:  s.URL,
			Key:  s.Key,
		}
	}
	return &coreapi.InfoReply{
		Version:  version.Version,
		Services: services,
	}, nil
}
