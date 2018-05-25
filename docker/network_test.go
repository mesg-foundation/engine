package docker

import (
	"testing"
	"time"

	"github.com/stvp/assert"
)

func TestCreateNetworkOverlay(t *testing.T) {
	network, err := CreateNetwork([]string{"TestCreateNetworkOverlay"}, "overlay")
	assert.Nil(t, err)
	assert.NotNil(t, network)
	DeleteNetwork([]string{"TestCreateNetworkOverlay"})
}

func TestCreateNetworkBridge(t *testing.T) {
	network, err := CreateNetwork([]string{"TestCreateNetworkBridge"}, "bridge")
	assert.Nil(t, err)
	assert.NotNil(t, network)
	DeleteNetwork([]string{"TestCreateNetworkBridge"})
}

func TestCreateAlreadyExistingNetwork(t *testing.T) {
	CreateNetwork([]string{"TestCreateAlreadyExistingNetwork"}, "overlay")
	network, err := CreateNetwork([]string{"TestCreateAlreadyExistingNetwork"}, "overlay")
	assert.Nil(t, err)
	assert.NotNil(t, network)
	DeleteNetwork([]string{"TestCreateAlreadyExistingNetwork"})
}

func TestDeleteNetwork(t *testing.T) {
	CreateNetwork([]string{"TestDeleteNetwork"}, "overlay")
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
	CreateNetwork([]string{"TestFindNetwork"}, "overlay")
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

func TestWaitNetworkDeletion(t *testing.T) {
	namespace := []string{"TestWaitNetworkDeletion"}
	CreateNetwork(namespace, "overlay")
	DeleteNetwork(namespace)
	err := <-WaitNetworkDeletion(namespace, 10*time.Second)
	assert.Nil(t, err)
}

func TestWaitNetworkDeletionTimeout(t *testing.T) {
	namespace := []string{"TestWaitNetworkDeletionTimeout"}
	CreateNetwork(namespace, "overlay")
	defer DeleteNetwork(namespace)
	err := <-WaitNetworkDeletion(namespace, time.Second)
	assert.NotNil(t, err)
}

func TestAttachNetworkToContainer(t *testing.T) {
	namespace := []string{"TestAttachNetworkToContainer"}
	startTestService(namespace)
	<-WaitContainerStatus(namespace, RUNNING, time.Minute)
	CreateNetwork(namespace, "bridge")
	err := AttachNetworkToContainer(namespace, namespace)
	assert.Nil(t, err)
	IP, err := FindContainerIP(namespace, namespace)
	assert.Nil(t, err)
	assert.NotEqual(t, "", IP)
	StopService(namespace)
	<-WaitContainerStatus(namespace, STOPPED, time.Minute)
	err = DeleteNetwork(namespace)
	assert.Nil(t, err)
}
