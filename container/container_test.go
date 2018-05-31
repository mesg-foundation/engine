package container

import (
	"errors"
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
	<-WaitContainerStatus(namespace, RUNNING, time.Minute)
	container, err := FindContainer(namespace)
	assert.Nil(t, err)
	assert.NotNil(t, container)
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
	<-WaitContainerStatus(namespace, RUNNING, time.Minute)
	status, err := ContainerStatus(namespace)
	assert.Nil(t, err)
	assert.Equal(t, status, RUNNING)
}

func TestContainerStatusStopped(t *testing.T) {
	namespace := []string{"TestContainerStatusStopped"}
	startTestService(namespace)
	<-WaitContainerStatus(namespace, RUNNING, time.Minute)
	fmt.Println("wait for running")
	StopService(namespace)
	<-WaitContainerStatus(namespace, STOPPED, time.Minute)
	fmt.Println("wait for stop")
	status, err := ContainerStatus(namespace)
	assert.Nil(t, err)
	assert.Equal(t, status, STOPPED)
}

func TestWaitForContainerRunning(t *testing.T) {
	namespace := []string{"TestWaitForContainerRunning"}
	startTestService(namespace)
	defer StopService(namespace)
	err := <-WaitContainerStatus(namespace, RUNNING, time.Minute)
	assert.Nil(t, err)
}

func TestWaitForContainerTimeout(t *testing.T) {
	namespace := []string{"TestWaitForContainerTimeout"}
	startTestService(namespace)
	defer StopService(namespace)
	err := <-WaitContainerStatus(namespace, RUNNING, time.Nanosecond)
	assert.NotNil(t, err)
}

func TestWaitForContainerStopped(t *testing.T) {
	namespace := []string{"TestWaitForContainerStopped"}
	startTestService(namespace)
	<-WaitContainerStatus(namespace, RUNNING, time.Minute)

	StopService(namespace)
	err := <-WaitContainerStatus(namespace, STOPPED, time.Minute)
	assert.Nil(t, err)
}

// WaitContainerStatus wait for the container to have the provided status until it reach the timeout
func WaitContainerStatus(namespace []string, status StatusType, timeout time.Duration) (wait chan error) {
	start := time.Now()
	wait = make(chan error, 1)
	go func() {
		for {
			currentStatus, err := ContainerStatus(namespace)
			if err != nil {
				wait <- err
				return
			}
			if currentStatus == status {
				close(wait)
				return
			}
			diff := time.Now().Sub(start)
			if diff.Nanoseconds() >= int64(timeout) {
				wait <- errors.New("Wait too long for the container, timeout reached")
				return
			}
			time.Sleep(500 * time.Millisecond)
		}
	}()
	return
}
