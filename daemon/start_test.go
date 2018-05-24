package daemon

import (
	"testing"

	"github.com/mesg-foundation/core/docker"
	"github.com/stvp/assert"
)

func TestStart(t *testing.T) {
	Stop()
	service, err := Start()
	defer Stop()
	assert.Nil(t, err)
	assert.NotNil(t, service)
}

func TestStartNetwork(t *testing.T) {
	Start()
	defer Stop()
	network, err := docker.FindNetwork(NamespaceNetwork())
	assert.Nil(t, err)
	assert.NotNil(t, network)
}
