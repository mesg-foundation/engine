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

	client *ServiceClientSafe
	c      chan *serviceListenTaskClientSafeResponse

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
	c := &serviceListenTaskClientSafe{
		client: client,
		c:      make(chan *serviceListenTaskClientSafeResponse),

		ctx:  ctx,
		in:   in,
		opts: opts,
	}
	waitStream := make(chan struct{}, 1)
	go c.recvLoop(waitStream)
	<-waitStream
	return c
}

// recvLoop recives ListenTask response in loop and reconnect in on error.
func (s *serviceListenTaskClientSafe) recvLoop(waitStream chan struct{}) {
	var err error
loop:
	for {
		// connect
		s.Service_ListenTaskClient, err = s.client.ServiceClient.ListenTask(s.ctx, s.in, s.opts...)
		waitStream <- struct{}{}
		if err != nil {
			s.c <- &serviceListenTaskClientSafeResponse{nil, err}
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
			s.c <- &serviceListenTaskClientSafeResponse{td, err}
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
	v := <-s.c
	return v.taskData, v.err
}

// ListenTask subscribes to a stream that listens for task to execute.
func (c *ServiceClientSafe) ListenTask(ctx context.Context, in *ListenTaskRequest, opts ...grpc.CallOption) (Service_ListenTaskClient, error) {
	return newServiceListenTaskClientSafe(c, ctx, in, opts...), nil
}
