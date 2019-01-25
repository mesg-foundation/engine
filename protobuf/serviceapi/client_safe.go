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

// NewServiceClient creates core client with reconnection.
func NewServiceClientSafe(cc *grpc.ClientConn) *ServiceClientSafe {
	return &ServiceClientSafe{
		ServiceClient: NewServiceClient(cc),
	}
}

// core_ListenTaskClient is a client with reconnection.
type core_ListenTaskClient struct {
	Service_ListenTaskClient

	client *ServiceClientSafe
	c      chan *core_ListenTaskClientResponse

	ctx  context.Context
	in   *ListenTaskRequest
	opts []grpc.CallOption
}

// core_ListenTaskClientResponse wraps ListenTask recv response.
type core_ListenTaskClientResponse struct {
	taskData *TaskData
	err      error
}

// newService_ListenTaskClient creates core ListenTask client.
func newService_ListenTaskClient(client *ServiceClientSafe, ctx context.Context, in *ListenTaskRequest, opts ...grpc.CallOption) *core_ListenTaskClient {
	c := &core_ListenTaskClient{
		client: client,
		c:      make(chan *core_ListenTaskClientResponse),

		ctx:  ctx,
		in:   in,
		opts: opts,
	}
	go c.recvLoop()
	return c
}

// recvLoop recives ListenTask response in loop and reconnect in on error.
func (s *core_ListenTaskClient) recvLoop() {
	var err error
loop:
	for {
		// connect
		s.Service_ListenTaskClient, err = s.client.ServiceClient.ListenTask(s.ctx, s.in, s.opts...)
		if err != nil {
			s.c <- &core_ListenTaskClientResponse{nil, err}
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
			s.c <- &core_ListenTaskClientResponse{td, err}
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
func (s *core_ListenTaskClient) Recv() (*TaskData, error) {
	v := <-s.c
	return v.taskData, v.err
}

func (c *ServiceClientSafe) ListenTask(ctx context.Context, in *ListenTaskRequest, opts ...grpc.CallOption) (Service_ListenTaskClient, error) {
	return newService_ListenTaskClient(c, ctx, in, opts...), nil
}
