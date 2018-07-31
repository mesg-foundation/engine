// +build integration

package container

import (
	"testing"

	"github.com/stvp/assert"
)

func TestIntegrationCreateSwarmIfNeeded(t *testing.T) {
	c, err := New()
	assert.Nil(t, err)
	assert.Nil(t, c.createSwarmIfNeeded())
}

func TestIntegrationFindContainerNotExisting(t *testing.T) {
	c, err := New()
	assert.Nil(t, err)
	_, err = c.FindContainer([]string{"TestFindContainerNotExisting"})
	assert.NotNil(t, err)
}

func TestIntegrationFindContainer(t *testing.T) {
	c, err := New()
	assert.Nil(t, err)
	namespace := []string{"TestFindContainer"}
	startTestService(namespace)
	defer c.StopService(namespace)
	c.waitForStatus(namespace, RUNNING)
	container, err := c.FindContainer(namespace)
	assert.Nil(t, err)
	assert.NotEqual(t, "", container.ID)
}

func TestIntegrationFindContainerStopped(t *testing.T) {
	c, err := New()
	assert.Nil(t, err)
	namespace := []string{"TestFindContainerStopped"}
	startTestService(namespace)
	c.StopService(namespace)
	_, err = c.FindContainer(namespace)
	assert.NotNil(t, err)
}

func TestIntegrationContainerStatusNeverStarted(t *testing.T) {
	c, err := New()
	assert.Nil(t, err)
	namespace := []string{"TestContainerStatusNeverStarted"}
	status, err := c.Status(namespace)
	assert.Nil(t, err)
	assert.Equal(t, status, STOPPED)
}

func TestIntegrationContainerStatusRunning(t *testing.T) {
	c, err := New()
	assert.Nil(t, err)
	namespace := []string{"TestContainerStatusRunning"}
	startTestService(namespace)
	defer c.StopService(namespace)
	c.waitForStatus(namespace, RUNNING)
	status, err := c.Status(namespace)
	assert.Nil(t, err)
	assert.Equal(t, status, RUNNING)
}

func TestIntegrationContainerStatusStopped(t *testing.T) {
	c, err := New()
	assert.Nil(t, err)
	namespace := []string{"TestContainerStatusStopped"}
	startTestService(namespace)
	c.waitForStatus(namespace, RUNNING)
	c.StopService(namespace)
	c.waitForStatus(namespace, STOPPED)
	status, err := c.Status(namespace)
	assert.Nil(t, err)
	assert.Equal(t, status, STOPPED)
}
