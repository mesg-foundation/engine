package serviceapi

import (
	"context"
	"io"
	"sync"
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

	client *ServiceClientSafe
	data   chan *serviceListenTaskClientSafeResponse

	mx  sync.Mutex
	err error

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
		client: client,
		data:   make(chan *serviceListenTaskClientSafeResponse),
		ctx:    ctx,
		in:     in,
		opts:   opts,
	}
	s.mx.Lock()
	go s.recvLoop()
	return s
}

// recvLoop recives ListenTask response in loop and reconnect in on error.
func (s *serviceListenTaskClientSafe) recvLoop() {
loop:
	for {
		// connect
		stream, err := s.client.ServiceClient.ListenTask(s.ctx, s.in, s.opts...)
		if stream != nil {
			s.Service_ListenTaskClient = stream
		}
		s.err = err
		s.mx.Unlock()
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
		s.mx.Lock()
	}
}

// Recv recives data from streams.
func (s *serviceListenTaskClientSafe) Recv() (*TaskData, error) {
	v := <-s.data
	return v.taskData, v.err
}

// ListenTask subscribes to a stream that listens for task to execute.
func (c *ServiceClientSafe) ListenTask(ctx context.Context, in *ListenTaskRequest, opts ...grpc.CallOption) (Service_ListenTaskClient, error) {
	cs := newServiceListenTaskClientSafe(c, ctx, in, opts...)
	cs.mx.Lock()
	defer cs.mx.Unlock()
	return cs, cs.err
}
