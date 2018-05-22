package docker

import (
	"testing"

	"github.com/stvp/assert"
)

func TestCreateNetwork(t *testing.T) {
	network, err := NetworkCreate("TestCreateNetwork")
	assert.Nil(t, err)
	assert.NotNil(t, network)
	NetworkDelete("TestCreateNetwork")
}

func TestDeleteNetwork(t *testing.T) {
	NetworkCreate("TestDeleteNetwork")
	err := NetworkDelete("TestDeleteNetwork")
	assert.Nil(t, err)
	network, err := NetworkFind("TestFindNetwork")
	assert.Nil(t, err)
	assert.Nil(t, network)
}

func TestFindNetwork(t *testing.T) {
	NetworkCreate("TestFindNetwork")
	network, err := NetworkFind("TestFindNetwork")
	assert.Nil(t, err)
	assert.NotNil(t, network)
	NetworkDelete("TestFindNetwork")
}

func TestFindMissingNetwork(t *testing.T) {
	network, err := NetworkFind("TestFindMissingNetwork")
	assert.Nil(t, err)
	assert.Nil(t, network)
}
