package daemon

import (
	"testing"

	"github.com/stvp/assert"
)

func TestSharedNetwork(t *testing.T) {
	Start()
	network, err := SharedNetwork()
	assert.Nil(t, err)
	assert.NotNil(t, network)
}
