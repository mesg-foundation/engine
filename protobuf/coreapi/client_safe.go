package coreapi

import (
	"context"
	"io"
	"time"

	"google.golang.org/grpc"
)

const reconnectDelay = 3 * time.Second

// CoreClientSafe provides CoreClient with stream reconnection.
type CoreClientSafe struct {
	CoreClient
}

// NewCoreClientSafe creates core client with reconnection.
func NewCoreClientSafe(cc *grpc.ClientConn) *CoreClientSafe {
	return &CoreClientSafe{
		CoreClient: NewCoreClient(cc),
	}
}

// coreListenEventClientSafe is a client with reconnection.
type coreListenEventClientSafe struct {
	Core_ListenEventClient

	client *CoreClientSafe
	c      chan *coreListenEventClientSafeResponse

	ctx  context.Context
	in   *ListenEventRequest
	opts []grpc.CallOption
}

// coreListenEventClientSafeResponse wraps ListenEvent recv response.
type coreListenEventClientSafeResponse struct {
	eventData *EventData
	err       error
}

// newCoreListenEventClientSafe creates core ListenEvent client.
func newCoreListenEventClientSafe(client *CoreClientSafe, ctx context.Context, in *ListenEventRequest, opts ...grpc.CallOption) *coreListenEventClientSafe {
	c := &coreListenEventClientSafe{
		client: client,
		c:      make(chan *coreListenEventClientSafeResponse),

		ctx:  ctx,
		in:   in,
		opts: opts,
	}
	waitStream := make(chan struct{}, 1)
	go c.recvLoop(waitStream)
	<-waitStream
	return c
}

// recvLoop receives ListenEvent response in loop and reconnect in on error.
func (s *coreListenEventClientSafe) recvLoop(waitStream chan struct{}) {
	var err error
loop:
	for {
		// connect
		s.Core_ListenEventClient, err = s.client.CoreClient.ListenEvent(s.ctx, s.in, s.opts...)
		waitStream <- struct{}{}
		if err != nil {
			s.c <- &coreListenEventClientSafeResponse{nil, err}
			continue
		}

		// buffered channel, because it might happen that ctx.Done
		// will be notified first before stream.Recv in for loop.
		done := make(chan struct{}, 1)

		go func(c Core_ListenEventClient) {
			select {
			case <-c.Context().Done():
				c.CloseSend()
			case <-done:
			}
		}(s.Core_ListenEventClient)

		for {
			td, err := s.Core_ListenEventClient.Recv()
			s.c <- &coreListenEventClientSafeResponse{td, err}
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

// Recv receives data from streams.
func (s *coreListenEventClientSafe) Recv() (*EventData, error) {
	v := <-s.c
	return v.eventData, v.err
}

// ListenEvent subscribes to a stream that listens for events from a service.
func (c *CoreClientSafe) ListenEvent(ctx context.Context, in *ListenEventRequest, opts ...grpc.CallOption) (Core_ListenEventClient, error) {
	return newCoreListenEventClientSafe(c, ctx, in, opts...), nil
}

// coreListenResultClientSafe is a client with reconnection.
type coreListenResultClientSafe struct {
	Core_ListenResultClient

	client *CoreClientSafe
	c      chan *coreListenResultClientSafeResponse

	ctx  context.Context
	in   *ListenResultRequest
	opts []grpc.CallOption
}

// coreListenResultClientSafeResponse wraps ListenResult recv response.
type coreListenResultClientSafeResponse struct {
	resultData *ResultData
	err        error
}

// newCoreListenResultClientSafe creates core ListenResult client.
func newCoreListenResultClientSafe(client *CoreClientSafe, ctx context.Context, in *ListenResultRequest, opts ...grpc.CallOption) *coreListenResultClientSafe {
	c := &coreListenResultClientSafe{
		client: client,
		c:      make(chan *coreListenResultClientSafeResponse),

		ctx:  ctx,
		in:   in,
		opts: opts,
	}
	waitStream := make(chan struct{}, 1)
	go c.recvLoop(waitStream)
	<-waitStream
	return c
}

// recvLoop receives ListenResult response in loop and reconnect in on error.
func (s *coreListenResultClientSafe) recvLoop(waitStream chan struct{}) {
	var err error
loop:
	for {
		// connect
		s.Core_ListenResultClient, err = s.client.CoreClient.ListenResult(s.ctx, s.in, s.opts...)
		waitStream <- struct{}{}
		if err != nil {
			s.c <- &coreListenResultClientSafeResponse{nil, err}
			continue
		}

		// buffered channel, because it might happen that ctx.Done
		// will be notified first before stream.Recv in for loop.
		done := make(chan struct{}, 1)

		go func(c Core_ListenResultClient) {
			select {
			case <-c.Context().Done():
				c.CloseSend()
			case <-done:
			}
		}(s.Core_ListenResultClient)

		for {
			td, err := s.Core_ListenResultClient.Recv()
			s.c <- &coreListenResultClientSafeResponse{td, err}
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

// Recv receives data from streams.
func (s *coreListenResultClientSafe) Recv() (*ResultData, error) {
	v := <-s.c
	return v.resultData, v.err
}

// ListenResult subscribes to a stream that listens for results from a service.
func (c *CoreClientSafe) ListenResult(ctx context.Context, in *ListenResultRequest, opts ...grpc.CallOption) (Core_ListenResultClient, error) {
	return newCoreListenResultClientSafe(c, ctx, in, opts...), nil
}

// coreServiceLogsClientSafe is a client with reconnection.
type coreServiceLogsClientSafe struct {
	Core_ServiceLogsClient

	client *CoreClientSafe
	c      chan *coreServiceLogsClientSafeResponse

	ctx  context.Context
	in   *ServiceLogsRequest
	opts []grpc.CallOption
}

// coreServiceLogsClientSafeResponse wraps ServiceLogs recv response.
type coreServiceLogsClientSafeResponse struct {
	logData *LogData
	err     error
}

// newCoreServiceLogsClientSafe creates core ServiceLogs client.
func newCoreServiceLogsClientSafe(client *CoreClientSafe, ctx context.Context, in *ServiceLogsRequest, opts ...grpc.CallOption) *coreServiceLogsClientSafe {
	c := &coreServiceLogsClientSafe{
		client: client,
		c:      make(chan *coreServiceLogsClientSafeResponse),

		ctx:  ctx,
		in:   in,
		opts: opts,
	}
	waitStream := make(chan struct{}, 1)
	go c.recvLoop(waitStream)
	<-waitStream
	return c
}

// recvLoop receives ServiceLogs response in loop and reconnect in on error.
func (s *coreServiceLogsClientSafe) recvLoop(waitStream chan struct{}) {
	var err error
loop:
	for {
		// connect
		s.Core_ServiceLogsClient, err = s.client.CoreClient.ServiceLogs(s.ctx, s.in, s.opts...)
		waitStream <- struct{}{}
		if err != nil {
			s.c <- &coreServiceLogsClientSafeResponse{nil, err}
			continue
		}

		// buffered channel, because it might happen that ctx.Done
		// will be notified first before stream.Recv in for loop.
		done := make(chan struct{}, 1)

		go func(c Core_ServiceLogsClient) {
			select {
			case <-c.Context().Done():
				c.CloseSend()
			case <-done:
			}
		}(s.Core_ServiceLogsClient)

		for {
			td, err := s.Core_ServiceLogsClient.Recv()
			s.c <- &coreServiceLogsClientSafeResponse{td, err}
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

// Recv receives data from streams.
func (s *coreServiceLogsClientSafe) Recv() (*LogData, error) {
	v := <-s.c
	return v.logData, v.err
}

// ServiceLogs subscribes to a stream that listens for logs from a service.
func (c *CoreClientSafe) ServiceLogs(ctx context.Context, in *ServiceLogsRequest, opts ...grpc.CallOption) (Core_ServiceLogsClient, error) {
	return newCoreServiceLogsClientSafe(c, ctx, in, opts...), nil
}
