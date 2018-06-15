package container

import (
	"testing"

	"github.com/stvp/assert"
)

func TestFindContainerNotExisting(t *testing.T) {
	_, err := FindContainer([]string{"TestFindContainerNotExisting"})
	assert.NotNil(t, err)
}

func TestFindContainer(t *testing.T) {
	namespace := []string{"TestFindContainer"}
	startTestService(namespace)
	defer StopService(namespace)
	WaitForContainerStatus(namespace, RUNNING)
	container, err := FindContainer(namespace)
	assert.Nil(t, err)
	assert.NotEqual(t, "", container.ID)
}

func TestFindContainerStopped(t *testing.T) {
	namespace := []string{"TestFindContainerStopped"}
	startTestService(namespace)
	StopService(namespace)
	_, err := FindContainer(namespace)
	assert.NotNil(t, err)
}

func TestContainerStatusNeverStarted(t *testing.T) {
	namespace := []string{"TestContainerStatusNeverStarted"}
	status, err := Status(namespace)
	assert.Nil(t, err)
	assert.Equal(t, status, STOPPED)
}

func TestContainerStatusRunning(t *testing.T) {
	namespace := []string{"TestContainerStatusRunning"}
	startTestService(namespace)
	defer StopService(namespace)
	WaitForContainerStatus(namespace, RUNNING)
	status, err := Status(namespace)
	assert.Nil(t, err)
	assert.Equal(t, status, RUNNING)
}

func TestContainerStatusStopped(t *testing.T) {
	namespace := []string{"TestContainerStatusStopped"}
	startTestService(namespace)
	WaitForContainerStatus(namespace, RUNNING)
	StopService(namespace)
	WaitForContainerStatus(namespace, STOPPED)
	status, err := Status(namespace)
	assert.Nil(t, err)
	assert.Equal(t, status, STOPPED)
}
