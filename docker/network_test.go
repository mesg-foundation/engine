package docker

import (
	"testing"
	"time"

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

func TestWaitNetworkDeletion(t *testing.T) {
	namespace := []string{"TestWaitNetworkDeletion"}
	CreateNetwork(namespace)
	DeleteNetwork(namespace)
	err := <-WaitNetworkDeletion(namespace, 10*time.Second)
	assert.Nil(t, err)
}

func TestWaitNetworkDeletionTimeout(t *testing.T) {
	namespace := []string{"TestWaitNetworkDeletionTimeout"}
	CreateNetwork(namespace)
	defer DeleteNetwork(namespace)
	err := <-WaitNetworkDeletion(namespace, time.Second)
	assert.NotNil(t, err)
}
