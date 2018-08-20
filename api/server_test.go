package api

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

const (
	waitForServe = 500 * time.Millisecond
)

func TestServerServe(t *testing.T) {
	s := Server{
		Network: "unix",
		Address: "TestServerServe.sock",
	}
	go func() {
		time.Sleep(waitForServe)
		s.Stop()
	}()
	err := s.Serve()
	require.Nil(t, err)
}

func TestServerServeNoAddress(t *testing.T) {
	s := Server{}
	go func() {
		time.Sleep(waitForServe)
		s.Stop()
	}()
	err := s.Serve()
	require.NotNil(t, err)
}

func TestServerServeAlreadyListening(t *testing.T) {
	s := Server{
		Network: "unix",
		Address: "TestServerServeAlreadyListening.sock",
	}
	go s.Serve()
	time.Sleep(waitForServe)
	err := s.Serve()
	defer s.Stop()
	require.NotNil(t, err)
}
