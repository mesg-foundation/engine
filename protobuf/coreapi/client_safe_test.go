package coreapi

import (
	"context"
	"errors"
	"io"
	"net"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const address = ":50051"

func newCoreServer(t *testing.T, m *MockCoreServer) *grpc.Server {
	lis, err := net.Listen("tcp", address)
	if err != nil {
		t.Fatalf("listen: %v", err)
	}

	s := grpc.NewServer()
	RegisterCoreServer(s, m)
	reflection.Register(s)
	go func() { s.Serve(lis) }()
	return s
}

func newCoreClientSafe(t *testing.T) *CoreClientSafe {
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		t.Fatal(err)
	}

	return NewCoreClientSafe(conn)
}

func TestListenEventReconnect(t *testing.T) {
	m := &MockCoreServer{}
	m.On("ListenEvent", mock.Anything, mock.Anything).Return(errors.New("autoreconnect")).Once()
	m.On("ListenEvent", mock.Anything, mock.Anything).Return(nil).Once()

	s := newCoreServer(t, m)
	defer s.Stop()
	c := newCoreClientSafe(t)

	stream, err := c.ListenEvent(context.Background(), &ListenEventRequest{})
	if err != nil {
		t.Fatal(err)
	}

	_, err = stream.Recv()
	require.Contains(t, err.Error(), "autoreconnect")
	_, err = stream.Recv()
	require.Equal(t, io.EOF, err)
	m.AssertExpectations(t)
}

func TestListenResultReconnect(t *testing.T) {
	m := &MockCoreServer{}
	m.On("ListenResult", mock.Anything, mock.Anything).Return(errors.New("autoreconnect")).Once()
	m.On("ListenResult", mock.Anything, mock.Anything).Return(nil).Once()

	s := newCoreServer(t, m)
	defer s.Stop()
	c := newCoreClientSafe(t)

	stream, err := c.ListenResult(context.Background(), &ListenResultRequest{})
	if err != nil {
		t.Fatal(err)
	}

	_, err = stream.Recv()
	require.Contains(t, err.Error(), "autoreconnect")
	_, err = stream.Recv()
	require.Equal(t, io.EOF, err)
	m.AssertExpectations(t)
}

func TestServiceLogsReconnect(t *testing.T) {
	m := &MockCoreServer{}
	m.On("ServiceLogs", mock.Anything, mock.Anything).Return(errors.New("autoreconnect")).Once()
	m.On("ServiceLogs", mock.Anything, mock.Anything).Return(nil).Once()

	s := newCoreServer(t, m)
	defer s.Stop()
	c := newCoreClientSafe(t)

	stream, err := c.ServiceLogs(context.Background(), &ServiceLogsRequest{})
	if err != nil {
		t.Fatal(err)
	}

	_, err = stream.Recv()
	require.Contains(t, err.Error(), "autoreconnect")
	_, err = stream.Recv()
	require.Equal(t, io.EOF, err)
	m.AssertExpectations(t)
}
