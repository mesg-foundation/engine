package daemon

import (
	"testing"

	"github.com/mesg-foundation/core/docker"
	"github.com/stvp/assert"
)

func TestStop(t *testing.T) {
	<-WaitForFullyStop()
	_, err := Start()
	assert.Nil(t, err)
	<-WaitForRunning()
	err = Stop()
	assert.Nil(t, err)
}

func TestStoptNetwork(t *testing.T) {
	<-WaitForFullyStop()
	_, err := Start()
	assert.Nil(t, err)
	<-WaitForRunning()
	err = Stop()
	assert.Nil(t, err)
	<-WaitForFullyStop()
	network, err := docker.FindNetwork(NamespaceNetwork())
	assert.Nil(t, err)
	assert.Nil(t, network)
}
