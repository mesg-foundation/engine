package servicetest

import (
	"context"

	"github.com/mesg-foundation/core/protobuf/serviceapi"
	"google.golang.org/grpc"
)

// serviceServer implements MESG's service server.
type serviceServer struct {
	taskC   chan *serviceapi.TaskData
	emitC   chan *Event
	submitC chan *Execution
	token   string

	closingC chan struct{}
}

func newServiceServer() *serviceServer {
	return &serviceServer{
		emitC:    make(chan *Event, 1),
		taskC:    make(chan *serviceapi.TaskData, 0),
		submitC:  make(chan *Execution, 0),
		closingC: make(chan struct{}, 0),
	}
}

func (s *serviceServer) EmitEvent(context context.Context,
	request *serviceapi.EmitEventRequest) (reply *serviceapi.EmitEventReply, err error) {
	s.emitC <- &Event{
		name:  request.EventKey,
		data:  request.EventData,
		token: request.Token,
	}
	return &serviceapi.EmitEventReply{}, nil
}

func (s *serviceServer) ListenTask(request *serviceapi.ListenTaskRequest,
	stream serviceapi.Service_ListenTaskServer) (err error) {
	s.token = request.Token

	for {
		select {
		case <-stream.Context().Done():
			close(s.closingC)
			return nil
		case task := <-s.taskC:
			if err := stream.Send(task); err != nil {
				close(s.closingC)
				return err
			}
		}
	}
}

func (s *serviceServer) SubmitResult(context context.Context,
	request *serviceapi.SubmitResultRequest) (reply *serviceapi.SubmitResultReply, err error) {
	s.submitC <- &Execution{
		id:   request.ExecutionID,
		key:  request.OutputKey,
		data: request.OutputData,
	}
	return &serviceapi.SubmitResultReply{}, nil
}

type taskDataStream struct {
	taskC  chan *serviceapi.TaskData
	ctx    context.Context
	cancel context.CancelFunc
	grpc.ServerStream
}

func newTaskDataStream() *taskDataStream {
	ctx, cancel := context.WithCancel(context.Background())
	return &taskDataStream{
		taskC:  make(chan *serviceapi.TaskData, 0),
		ctx:    ctx,
		cancel: cancel,
	}
}

func (s taskDataStream) Send(data *serviceapi.TaskData) error {
	s.taskC <- data
	return nil
}

func (s taskDataStream) Context() context.Context {
	return s.ctx
}

func (s taskDataStream) close() {
	s.cancel()
}
