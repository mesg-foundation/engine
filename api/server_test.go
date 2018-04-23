package api

import (
	"testing"
	"time"

	"github.com/stvp/assert"
)

const (
	waitForServe = 500 * time.Millisecond
)

func TestServerDefaultConfig(t *testing.T) {
	s := Server{}
	assert.NotEqual(t, s.address(), "", "address should not be empty")
	assert.NotEqual(t, s.network(), "", "network should not be empty")
}

func TestServerServe(t *testing.T) {
	s := Server{}
	go func() {
		time.Sleep(waitForServe)
		s.Stop()
	}()
	err := s.Serve()
	assert.Nil(t, err)
}

func TestServerServeAlreadyListening(t *testing.T) {
	s := Server{}
	go s.Serve()
	time.Sleep(waitForServe)
	err := s.Serve()
	defer s.Stop()
	assert.NotNil(t, err)
}
