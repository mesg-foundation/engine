package network

import (
	"testing"

	"github.com/stvp/assert"
)

func TestCreateNetwork(t *testing.T) {
	network, err := Create("TestCreateNetwork")
	assert.Nil(t, err)
	assert.NotNil(t, network)
	Delete("TestCreateNetwork")
}
