package daemon

import (
	"testing"

	"github.com/mesg-foundation/core/docker"
	"github.com/stvp/assert"
)

func TestStart(t *testing.T) {
	<-testForceAndWaitForFullStop()
	service, err := Start()
	assert.Nil(t, err)
	assert.NotNil(t, service)
}

func TestStartNetwork(t *testing.T) {
	Start()
	network, err := docker.FindNetwork(NamespaceNetwork())
	assert.Nil(t, err)
	assert.NotNil(t, network)
}
