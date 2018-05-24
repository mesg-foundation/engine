package docker

import (
	"testing"

	"github.com/stvp/assert"
)

func TestFindContainerNotExisting(t *testing.T) {
	namespace := []string{"TestFindContainerNotExisting"}
	container, err := FindContainer(namespace)
	assert.Nil(t, err)
	assert.Nil(t, container)
}

func TestFindContainer(t *testing.T) {
	namespace := []string{"TestFindContainer"}
	startTestService(namespace)
	defer StopService(namespace)
	err := WaitForContainer(namespace)
	assert.Nil(t, err)
	container, err := FindContainer(namespace)
	assert.Nil(t, err)
	assert.NotNil(t, container)
}

func TestFindContainerStopped(t *testing.T) {
	namespace := []string{"TestFindContainerStopped"}
	startTestService(namespace)
	StopService(namespace)
	container, err := FindContainer(namespace)
	assert.Nil(t, err)
	assert.Nil(t, container)
}

func TestContainerStatusNeverStarted(t *testing.T) {
	namespace := []string{"TestContainerStatusNeverStarted"}
	status, err := ContainerStatus(namespace)
	assert.Nil(t, err)
	assert.Equal(t, status, STOPPED)
}

func TestContainerStatusRunning(t *testing.T) {
	namespace := []string{"TestContainerStatusRunning"}
	startTestService(namespace)
	defer StopService(namespace)
	err := WaitForContainer(namespace)
	assert.Nil(t, err)
	status, err := ContainerStatus(namespace)
	assert.Nil(t, err)
	assert.Equal(t, status, RUNNING)
}

func TestContainerStatusStopped(t *testing.T) {
	namespace := []string{"TestContainerStatusStopped"}
	startTestService(namespace)
	StopService(namespace)
	status, err := ContainerStatus(namespace)
	assert.Nil(t, err)
	assert.Equal(t, status, STOPPED)
}
