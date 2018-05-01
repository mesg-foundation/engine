package api

import (
	"testing"
	"time"

	"github.com/mesg-foundation/core/config"
	"github.com/stvp/assert"
)

const (
	waitForServe = 500 * time.Millisecond
)

func TestServerServe(t *testing.T) {
	s := Server{
		Network: config.Api.Server.Network(),
		Address: "TestServerServe.sock",
	}
	go func() {
		time.Sleep(waitForServe)
		s.Stop()
	}()
	err := s.Serve()
	assert.Nil(t, err)
}

func TestServerServeNoAddress(t *testing.T) {
	s := Server{}
	go func() {
		time.Sleep(waitForServe)
		s.Stop()
	}()
	err := s.Serve()
	assert.NotNil(t, err)
}

func TestServerServeAlreadyListening(t *testing.T) {
	s := Server{
		Network: config.Api.Server.Network(),
		Address: "TestServerServeAlreadyListening.sock",
	}
	go s.Serve()
	time.Sleep(waitForServe)
	err := s.Serve()
	defer s.Stop()
	assert.NotNil(t, err)
}
