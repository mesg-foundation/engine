package container

import (
	"testing"

	"github.com/stvp/assert"
)

func TestWaitForStatusRunning(t *testing.T) {
	c, err := New()
	assert.Nil(t, err)
	namespace := []string{"TestWaitForStatusRunning"}
	startTestService(namespace)
	defer c.StopService(namespace)
	err = c.waitForStatus(namespace, RUNNING)
	assert.Nil(t, err)
}

func TestWaitForStatusStopped(t *testing.T) {
	c, err := New()
	assert.Nil(t, err)
	namespace := []string{"TestWaitForStatusStopped"}
	startTestService(namespace)
	c.waitForStatus(namespace, RUNNING)
	c.StopService(namespace)
	err = c.waitForStatus(namespace, STOPPED)
	assert.Nil(t, err)
}

func TestWaitForStatusTaskError(t *testing.T) {
	c, err := New()
	assert.Nil(t, err)
	namespace := []string{"TestWaitForStatusTaskError"}
	c.StartService(ServiceOptions{
		Image:     "awgdaywudaywudwa",
		Namespace: namespace,
	})
	defer c.StopService(namespace)
	err = c.waitForStatus(namespace, RUNNING)
	assert.NotNil(t, err)
	assert.Equal(t, "No such image: awgdaywudaywudwa:latest", err.Error())
}
