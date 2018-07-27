package applicationtest

import (
	"context"
	"errors"
	"sync"

	"github.com/mesg-foundation/core/api/core"
	uuid "github.com/satori/go.uuid"
)

// coreServer implements MESG's core/application server.
type coreServer struct {
	listenEventC  chan *core.ListenEventRequest
	listenResultC chan *core.ListenResultRequest
	serviceStartC chan *core.StartServiceRequest
	executeC      chan *core.ExecuteTaskRequest

	eventC  map[string]chan *core.EventData
	resultC map[string]chan *core.ResultData
	em      sync.Mutex

	closingC            map[string]chan struct{}
	nonExistentServices []string
}

func newCoreServer() *coreServer {
	return &coreServer{
		listenEventC:  make(chan *core.ListenEventRequest, 0),
		listenResultC: make(chan *core.ListenResultRequest, 0),
		serviceStartC: make(chan *core.StartServiceRequest, 0),
		executeC:      make(chan *core.ExecuteTaskRequest, 0),
		eventC:        make(map[string]chan *core.EventData, 0),
		resultC:       make(map[string]chan *core.ResultData, 0),
		closingC:      make(map[string]chan struct{}, 0),
	}
}

func (s *coreServer) DeleteService(ctx context.Context,
	request *core.DeleteServiceRequest) (reply *core.DeleteServiceReply, err error) {
	return &core.DeleteServiceReply{}, nil
}

func (s *coreServer) DeployService(ctx context.Context,
	request *core.DeployServiceRequest) (reply *core.DeployServiceReply, err error) {
	return &core.DeployServiceReply{}, nil
}

func (s *coreServer) ExecuteTask(ctx context.Context,
	request *core.ExecuteTaskRequest) (reply *core.ExecuteTaskReply, err error) {
	s.executeC <- request
	uuidV4, err := uuid.NewV4()
	id := uuidV4.String()
	return &core.ExecuteTaskReply{
		ExecutionID: id,
	}, err
}

func (s *coreServer) GetService(ctx context.Context,
	request *core.GetServiceRequest) (reply *core.GetServiceReply, err error) {
	return &core.GetServiceReply{}, nil
}

func (s *coreServer) ListServices(ctx context.Context,
	request *core.ListServicesRequest) (reply *core.ListServicesReply, err error) {
	return &core.ListServicesReply{}, nil
}

func (s *coreServer) ListenEvent(request *core.ListenEventRequest,
	stream core.Core_ListenEventServer) (err error) {
	s.listenEventC <- request

	s.em.Lock()
	if s.eventC[request.ServiceID] == nil {
		s.eventC[request.ServiceID] = make(chan *core.EventData, 0)
	}
	s.em.Unlock()

	for {
		select {
		case <-stream.Context().Done():
			return nil
		case event := <-s.eventC[request.ServiceID]:
			if err := stream.Send(event); err != nil {
				return err
			}
		}
	}
}

func (s *coreServer) ListenResult(request *core.ListenResultRequest,
	stream core.Core_ListenResultServer) (err error) {
	s.listenResultC <- request
	s.em.Lock()
	if s.resultC[request.ServiceID] == nil {
		s.resultC[request.ServiceID] = make(chan *core.ResultData, 0)
	}
	s.em.Unlock()

	for {
		select {
		case <-stream.Context().Done():
			return nil
		case result := <-s.resultC[request.ServiceID]:
			if err := stream.Send(result); err != nil {
				return err
			}
		}
	}
}

func (s *coreServer) StartService(ctx context.Context,
	request *core.StartServiceRequest) (reply *core.StartServiceReply, err error) {
	for _, id := range s.nonExistentServices {
		if request.ServiceID == id {
			return &core.StartServiceReply{}, errors.New("service does not exists")
		}
	}
	s.serviceStartC <- request
	return &core.StartServiceReply{}, nil
}

func (s *coreServer) StopService(ctx context.Context,
	request *core.StopServiceRequest) (reply *core.StopServiceReply, err error) {
	return &core.StopServiceReply{}, nil
}
