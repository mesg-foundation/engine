package container

import (
	"fmt"
	"testing"
	"time"

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
	WaitForContainerStatus(namespace, RUNNING, time.Minute)
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
	status, err := ContainerStatus(namespace)
	assert.Nil(t, err)
	assert.Equal(t, status, STOPPED)
}

func TestContainerStatusRunning(t *testing.T) {
	namespace := []string{"TestContainerStatusRunning"}
	startTestService(namespace)
	defer StopService(namespace)
	WaitForContainerStatus(namespace, RUNNING, time.Minute)
	status, err := ContainerStatus(namespace)
	assert.Nil(t, err)
	assert.Equal(t, status, RUNNING)
}

func TestContainerStatusStopped(t *testing.T) {
	namespace := []string{"TestContainerStatusStopped"}
	startTestService(namespace)
	WaitForContainerStatus(namespace, RUNNING, time.Minute)
	fmt.Println("wait for running")
	StopService(namespace)
	WaitForContainerStatus(namespace, STOPPED, time.Minute)
	fmt.Println("wait for stop")
	status, err := ContainerStatus(namespace)
	assert.Nil(t, err)
	assert.Equal(t, status, STOPPED)
}

func TestWaitForContainerRunning(t *testing.T) {
	namespace := []string{"TestWaitForContainerRunning"}
	startTestService(namespace)
	defer StopService(namespace)
	err := WaitForContainerStatus(namespace, RUNNING, time.Minute)
	assert.Nil(t, err)
}

func TestWaitForContainerTimeout(t *testing.T) {
	namespace := []string{"TestWaitForContainerTimeout"}
	startTestService(namespace)
	defer StopService(namespace)
	err := WaitForContainerStatus(namespace, RUNNING, time.Nanosecond)
	assert.NotNil(t, err)
	_, ok := err.(*TimeoutError)
	assert.True(t, ok)
}

func TestWaitForContainerStopped(t *testing.T) {
	namespace := []string{"TestWaitForContainerStopped"}
	startTestService(namespace)
	WaitForContainerStatus(namespace, RUNNING, time.Minute)

	StopService(namespace)
	err := WaitForContainerStatus(namespace, STOPPED, time.Minute)
	assert.Nil(t, err)
}
