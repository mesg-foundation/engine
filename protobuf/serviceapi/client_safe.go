package serviceapi

import (
	"context"
	"io"
	"time"

	"google.golang.org/grpc"
)

const reconnectDelay = 3 * time.Second

// ServiceClientSafe provides ServiceClient with stream reconnection.
type ServiceClientSafe struct {
	ServiceClient
}

// NewServiceClientSafe creates core client with reconnection.
func NewServiceClientSafe(cc *grpc.ClientConn) *ServiceClientSafe {
	return &ServiceClientSafe{
		ServiceClient: NewServiceClient(cc),
	}
}

// serviceListenTaskClientSafe is a client with reconnection.
type serviceListenTaskClientSafe struct {
	Service_ListenTaskClient

	client        *ServiceClientSafe
	data          chan *serviceListenTaskClientSafeResponse
	streamCreated chan struct{}

	ctx  context.Context
	in   *ListenTaskRequest
	opts []grpc.CallOption
}

// serviceListenTaskClientSafeResponse wraps ListenTask recv response.
type serviceListenTaskClientSafeResponse struct {
	taskData *TaskData
	err      error
}

// newServiceListenTaskClientSafe creates core ListenTask client.
func newServiceListenTaskClientSafe(client *ServiceClientSafe, ctx context.Context, in *ListenTaskRequest, opts ...grpc.CallOption) *serviceListenTaskClientSafe {
	s := &serviceListenTaskClientSafe{
		client:        client,
		data:          make(chan *serviceListenTaskClientSafeResponse),
		streamCreated: make(chan struct{}, 1),
		ctx:  ctx,
		in:   in,
		opts: opts,
	}
	go s.recvLoop()
	<-s.streamCreated
	return s
}

// recvLoop recives ListenTask response in loop and reconnect in on error.
func (s *serviceListenTaskClientSafe) recvLoop() {
	var err error
loop:
	for {
		// connect
		s.Service_ListenTaskClient, err = s.client.ServiceClient.ListenTask(s.ctx, s.in, s.opts...)
		s.streamCreated <- struct{}{}
		if err != nil {
			s.data <- &serviceListenTaskClientSafeResponse{nil, err}
			continue
		}

		// buffered channel, because it might happen that ctx.Done
		// will be notified first before stream.Recv in for loop.
		done := make(chan struct{}, 1)

		go func(c Service_ListenTaskClient) {
			select {
			case <-c.Context().Done():
				c.CloseSend()
			case <-done:
			}
		}(s.Service_ListenTaskClient)

		for {
			td, err := s.Service_ListenTaskClient.Recv()
			s.data <- &serviceListenTaskClientSafeResponse{td, err}
			if err != nil {
				done <- struct{}{}
				// in case of EOF end loop
				if err == io.EOF {
					break loop
				}
				break
			}
		}

		// sleep before reconnect
		time.Sleep(reconnectDelay)
	}
}

// Recv recives data from streams.
func (s *serviceListenTaskClientSafe) Recv() (*TaskData, error) {
	v := <-s.data
	return v.taskData, v.err
}

// ListenTask subscribes to a stream that listens for task to execute.
func (c *ServiceClientSafe) ListenTask(ctx context.Context, in *ListenTaskRequest, opts ...grpc.CallOption) (Service_ListenTaskClient, error) {
	return newServiceListenTaskClientSafe(c, ctx, in, opts...), nil
}
