package network

import (
	"testing"

	"github.com/stvp/assert"
)

func TestFindNetwork(t *testing.T) {
	Create("TestFindNetwork")
	network, err := Find("TestFindNetwork")
	assert.Nil(t, err)
	assert.NotNil(t, network)
	Delete("TestFindNetwork")
}

func TestFindMissingNetwork(t *testing.T) {
	network, err := Find("TestFindMissingNetwork")
	assert.Nil(t, err)
	assert.Nil(t, network)
}
