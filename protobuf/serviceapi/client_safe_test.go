package serviceapi

import (
	"context"
	"errors"
	"io"
	"net"
	"testing"

	"github.com/mesg-foundation/core/protobuf/acknowledgement"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const address = ":50051"

func newServiceServer(t *testing.T, m *MockServiceServer) *grpc.Server {
	lis, err := net.Listen("tcp", address)
	if err != nil {
		t.Fatalf("listen: %v", err)
	}

	s := grpc.NewServer()
	RegisterServiceServer(s, m)
	reflection.Register(s)
	go func() { s.Serve(lis) }()
	return s
}

func newServiceClinetSafe(t *testing.T) *ServiceClientSafe {
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		t.Fatal(err)
	}

	return NewServiceClientSafe(conn, OnError)
}

func TestListenTaskReconnect(t *testing.T) {
	m := &MockServiceServer{}
	m.On("ListenTask", mock.Anything, mock.Anything).Return(func(_ *ListenTaskRequest, stream Service_ListenTaskServer) error {
		acknowledgement.SetStreamReady(stream)
		return errors.New("autoreconnect")
	}).Once()
	m.On("ListenTask", mock.Anything, mock.Anything).Return(func(_ *ListenTaskRequest, stream Service_ListenTaskServer) error {
		acknowledgement.SetStreamReady(stream)
		return nil
	}).Once()

	s := newServiceServer(t, m)
	defer s.Stop()
	c := newServiceClinetSafe(t)

	stream, err := c.ListenTask(context.Background(), &ListenTaskRequest{})
	if err != nil {
		t.Fatal(err)
	}

	_, err = stream.Recv()
	require.Contains(t, err.Error(), "autoreconnect")
	_, err = stream.Recv()
	require.Equal(t, io.EOF, err)
	m.AssertExpectations(t)
}
