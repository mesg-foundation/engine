package core

import (
	"context"
	"encoding/json"
	"sync"

	"github.com/mesg-foundation/core/api"
	"github.com/mesg-foundation/core/config"
	"github.com/mesg-foundation/core/protobuf/acknowledgement"
	"github.com/mesg-foundation/core/protobuf/coreapi"
	"github.com/mesg-foundation/core/service"
	"github.com/mesg-foundation/core/version"
	"github.com/mesg-foundation/core/x/xerrors"
)

// Server is the type to aggregate all the APIs.
type Server struct {
	api *api.API
}

// NewServer creates a new Server.
func NewServer(api *api.API) *Server {
	return &Server{api: api}
}

// GetService returns service serviceID.
func (s *Server) GetService(ctx context.Context, request *coreapi.GetServiceRequest) (*coreapi.GetServiceReply, error) {
	ss, err := s.api.GetService(request.ServiceID)
	if err != nil {
		return nil, err
	}
	status, err := ss.Status()
	if err != nil {
		return nil, err
	}
	details := &coreapi.Service{
		Definition: toProtoService(ss),
		Status:     toProtoServiceStatusType(status),
	}
	return &coreapi.GetServiceReply{Service: details}, nil
}

// ListServices lists services.
func (s *Server) ListServices(ctx context.Context, request *coreapi.ListServicesRequest) (*coreapi.ListServicesReply, error) {
	services, err := s.api.ListServices()
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
	for _, s := range services {
		go func(s *service.Service) {
			defer wg.Done()
			status, err := s.Status()
			if err != nil {
				errC <- err
				return
			}
			details := &coreapi.Service{
				Definition: toProtoService(s),
				Status:     toProtoServiceStatusType(status),
			}
			mp.Lock()
			protoServices = append(protoServices, details)
			mp.Unlock()
		}(s)
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
	return &coreapi.StartServiceReply{}, s.api.StartService(request.ServiceID)
}

// StopService stops a service.
func (s *Server) StopService(ctx context.Context, request *coreapi.StopServiceRequest) (*coreapi.StopServiceReply, error) {
	return &coreapi.StopServiceReply{}, s.api.StopService(request.ServiceID)
}

// DeleteService stops and deletes service serviceID.
func (s *Server) DeleteService(ctx context.Context, request *coreapi.DeleteServiceRequest) (*coreapi.DeleteServiceReply, error) {
	return &coreapi.DeleteServiceReply{}, s.api.DeleteService(request.ServiceID, request.DeleteData)
}

// ListenEvent listens events matches with eventFilter on serviceID.
func (s *Server) ListenEvent(request *coreapi.ListenEventRequest, stream coreapi.Core_ListenEventServer) error {
	ln, err := s.api.ListenEvent(request.ServiceID, api.ListenEventKeyFilter(request.EventFilter))
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

		case err := <-ln.Err:
			return err

		case ev := <-ln.Events:
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
	ln, err := s.api.ListenResult(request.ServiceID,
		api.ListenResultTaskFilter(request.TaskFilter),
		api.ListenResultOutputFilter(request.OutputFilter),
		api.ListenResultTagFilters(request.TagFilters))
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

		case err := <-ln.Err:
			return err

		case execution := <-ln.Executions:
			outputs, err := json.Marshal(execution.OutputData)
			if err != nil {
				return err
			}
			if err := stream.Send(&coreapi.ResultData{
				ExecutionID:   execution.ID,
				TaskKey:       execution.TaskKey,
				OutputKey:     execution.OutputKey,
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
	var inputs map[string]interface{}
	if err := json.Unmarshal([]byte(request.InputData), &inputs); err != nil {
		return nil, err
	}

	executionID, err := s.api.ExecuteTask(request.ServiceID, request.TaskKey, inputs, request.ExecutionTags)
	return &coreapi.ExecuteTaskReply{
		ExecutionID: executionID,
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
			Hash: s.Hash,
			Url:  s.URL,
			Key:  s.Key,
		}
	}
	return &coreapi.InfoReply{
		Address:  c.Client.Address,
		Image:    c.Core.Image,
		Version:  version.Version,
		Services: services,
	}, nil
}
