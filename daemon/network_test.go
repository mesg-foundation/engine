package daemon

import (
	"testing"

	"github.com/stvp/assert"
)

func TestIP(t *testing.T) {
	Start()
	defer Stop()

	daemonIP, err := IP()
	assert.Nil(t, err)
	assert.NotEqual(t, "", daemonIP)
}

func TestSharedNetwork(t *testing.T) {
	Start()
	defer Stop()

	network, err := SharedNetwork()
	assert.Nil(t, err)
	assert.NotNil(t, network)
}
