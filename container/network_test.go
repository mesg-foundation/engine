package container

import (
	"testing"

	"github.com/stvp/assert"
)

func TestCreateNetworkOverlay(t *testing.T) {
	network, err := CreateNetwork([]string{"TestCreateNetworkOverlay"})
	assert.Nil(t, err)
	assert.NotNil(t, network)
	DeleteNetwork([]string{"TestCreateNetworkOverlay"})
}

func TestCreateNetworkBridge(t *testing.T) {
	network, err := CreateNetwork([]string{"TestCreateNetworkBridge"})
	assert.Nil(t, err)
	assert.NotNil(t, network)
	DeleteNetwork([]string{"TestCreateNetworkBridge"})
}

func TestCreateAlreadyExistingNetwork(t *testing.T) {
	CreateNetwork([]string{"TestCreateAlreadyExistingNetwork"})
	network, err := CreateNetwork([]string{"TestCreateAlreadyExistingNetwork"})
	assert.Nil(t, err)
	assert.NotNil(t, network)
	DeleteNetwork([]string{"TestCreateAlreadyExistingNetwork"})
}

func TestDeleteNetwork(t *testing.T) {
	CreateNetwork([]string{"TestDeleteNetwork"})
	err := DeleteNetwork([]string{"TestDeleteNetwork"})
	assert.Nil(t, err)
	network, err := FindNetwork([]string{"TestFindNetwork"})
	assert.Nil(t, err)
	assert.Nil(t, network)
}

func TestDeleteNotExistingNetwork(t *testing.T) {
	err := DeleteNetwork([]string{"TestDeleteNotExistingNetwork"})
	assert.Nil(t, err)
	network, err := FindNetwork([]string{"TestDeleteNotExistingNetwork"})
	assert.Nil(t, err)
	assert.Nil(t, network)
}

func TestFindNetwork(t *testing.T) {
	CreateNetwork([]string{"TestFindNetwork"})
	network, err := FindNetwork([]string{"TestFindNetwork"})
	assert.Nil(t, err)
	assert.NotNil(t, network)
	DeleteNetwork([]string{"TestFindNetwork"})
}

func TestFindStrictNotExistingNetwork(t *testing.T) {
	network, err := FindNetworkStrict([]string{"TestFindStrictNotExistingNetwork"})
	assert.NotNil(t, err)
	assert.Nil(t, network)
}

func TestFindNotExistingNetwork(t *testing.T) {
	network, err := FindNetwork([]string{"TestFindNotExistingNetwork"})
	assert.Nil(t, err)
	assert.Nil(t, network)
}
