package docker

import (
	"testing"

	"github.com/stvp/assert"
)

func TestCreateNetwork(t *testing.T) {
	network, err := CreateNetwork("TestCreateNetwork")
	assert.Nil(t, err)
	assert.NotNil(t, network)
	DeleteNetwork("TestCreateNetwork")
}

func TestDeleteNetwork(t *testing.T) {
	CreateNetwork("TestDeleteNetwork")
	err := DeleteNetwork("TestDeleteNetwork")
	assert.Nil(t, err)
	network, err := FindNetwork("TestFindNetwork")
	assert.Nil(t, err)
	assert.Nil(t, network)
}

func TestFindNetwork(t *testing.T) {
	CreateNetwork("TestFindNetwork")
	network, err := FindNetwork("TestFindNetwork")
	assert.Nil(t, err)
	assert.NotNil(t, network)
	DeleteNetwork("TestFindNetwork")
}

func TestFindMissingNetwork(t *testing.T) {
	network, err := FindNetwork("TestFindMissingNetwork")
	assert.Nil(t, err)
	assert.Nil(t, network)
}
