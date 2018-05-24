package docker

import (
	"testing"

	"github.com/stvp/assert"
)

func TestNetworkNamespace(t *testing.T) {
	namespace := networkNamespace([]string{"test"})
	assert.Equal(t, namespace, Namespace([]string{networkNamespacePrefix, "test"}))
}

func TestCreateNetwork(t *testing.T) {
	network, err := CreateNetwork([]string{"TestCreateNetwork"})
	assert.Nil(t, err)
	assert.NotNil(t, network)
	DeleteNetwork([]string{"TestCreateNetwork"})
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

func TestFindNotExistingNetwork(t *testing.T) {
	network, err := FindNetwork([]string{"TestFindNotExistingNetwork"})
	assert.Nil(t, err)
	assert.Nil(t, network)
}
