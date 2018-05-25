package daemon

import (
	"testing"

	"github.com/mesg-foundation/core/docker"
	"github.com/stvp/assert"
)

func TestStop(t *testing.T) {
	Start()
	err := Stop()
	assert.Nil(t, err)
}

func TestStoptNetwork(t *testing.T) {
	Start()
	err := Stop()
	assert.Nil(t, err)
	err = <-WaitForFullStop()
	assert.Nil(t, err)
	network, err := docker.FindNetwork(NamespaceNetwork())
	assert.Nil(t, err)
	assert.Nil(t, network)
}
