package daemon

import (
	"testing"
	"time"

	"github.com/mesg-foundation/core/docker"
	"github.com/stvp/assert"
)

func TestStop(t *testing.T) {
	Start()
	<-WaitForRunning()
	err := Stop()
	assert.Nil(t, err)
}

func TestStoptNetwork(t *testing.T) {
	Start()
	<-WaitForRunning()
	Stop()
	<-WaitForStopped()
	time.Sleep(3 * time.Second) //TODO: that's very ugly
	network, err := docker.FindNetwork(NamespaceNetwork())
	assert.Nil(t, err)
	assert.Nil(t, network)
}
