package container

import (
	"testing"

	"github.com/stvp/assert"
)

func TestWaitForStatusRunning(t *testing.T) {
	namespace := []string{"TestWaitForStatusRunning"}
	startTestService(namespace)
	err := waitForStatus(namespace, RUNNING)
	assert.Nil(t, err)
	StopService(namespace)
}

func TestWaitForStatusStopped(t *testing.T) {
	namespace := []string{"TestWaitForStatusStopped"}
	startTestService(namespace)
	waitForStatus(namespace, RUNNING)
	StopService(namespace)
	err := waitForStatus(namespace, STOPPED)
	assert.Nil(t, err)
}

func TestWaitForStatusTaskError(t *testing.T) {
	namespace := []string{"TestWaitForStatusTaskError"}
	StartService(ServiceOptions{
		Image:     "awgdaywudaywudwa",
		Namespace: namespace,
	})
	err := waitForStatus(namespace, RUNNING)
	assert.NotNil(t, err)
	assert.Equal(t, "No such image: awgdaywudaywudwa:latest", err.Error())
	StopService(namespace)
}
