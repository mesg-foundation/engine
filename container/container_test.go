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
	waitForContainerStatus(namespace, RUNNING)
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
	waitForContainerStatus(namespace, RUNNING)
	status, err := Status(namespace)
	assert.Nil(t, err)
	assert.Equal(t, status, RUNNING)
}

func TestContainerStatusStopped(t *testing.T) {
	namespace := []string{"TestContainerStatusStopped"}
	startTestService(namespace)
	waitForContainerStatus(namespace, RUNNING)
	StopService(namespace)
	waitForContainerStatus(namespace, STOPPED)
	status, err := Status(namespace)
	assert.Nil(t, err)
	assert.Equal(t, status, STOPPED)
}

func TestWaitForContainerRunning(t *testing.T) {
	namespace := []string{"TestWaitForContainerRunning"}
	startTestService(namespace)
	defer StopService(namespace)
	err := waitForContainerStatus(namespace, RUNNING)
	assert.Nil(t, err)
}

func TestWaitForContainerStopped(t *testing.T) {
	namespace := []string{"TestWaitForContainerStopped"}
	startTestService(namespace)
	waitForContainerStatus(namespace, RUNNING)

	StopService(namespace)
	err := waitForContainerStatus(namespace, STOPPED)
	assert.Nil(t, err)
}
