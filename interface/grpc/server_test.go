package grpc

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

const (
	waitForServe = 500 * time.Millisecond
)

func TestServerServe(t *testing.T) {
	s := New("localhost:50052", nil, nil)
	go func() {
		time.Sleep(waitForServe)
		s.Close()
	}()
	err := s.Serve()
	require.NoError(t, err)
}

func TestServerServeNoAddress(t *testing.T) {
	s := Server{}
	go func() {
		time.Sleep(waitForServe)
		s.Close()
	}()
	err := s.Serve()
	require.Error(t, err)
}

func TestServerListenAfterClose(t *testing.T) {
	s := New("localhost:50052", nil, nil)
	go s.Serve()
	time.Sleep(waitForServe)
	s.Close()
	require.Equal(t, &alreadyClosedError{}, s.Serve())
}
