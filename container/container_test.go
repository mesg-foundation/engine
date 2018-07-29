package container

import (
	"testing"

	"github.com/stvp/assert"
)

func TestContainer(t *testing.T) {
	c, err := New()
	assert.Nil(t, err)
	assert.NotNil(t, c)
}

func TestCreateSwarmIfNeeded(t *testing.T) {
	c, err := New()
	assert.Nil(t, err)
	err = c.createSwarmIfNeeded()
	assert.Nil(t, err)
}

func TestFindContainerNotExisting(t *testing.T) {
	c, err := New()
	assert.Nil(t, err)
	_, err = c.FindContainer([]string{"TestFindContainerNotExisting"})
	assert.NotNil(t, err)
}

func TestFindContainer(t *testing.T) {
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

func TestFindContainerStopped(t *testing.T) {
	c, err := New()
	assert.Nil(t, err)
	namespace := []string{"TestFindContainerStopped"}
	startTestService(namespace)
	c.StopService(namespace)
	_, err = c.FindContainer(namespace)
	assert.NotNil(t, err)
}

func TestContainerStatusNeverStarted(t *testing.T) {
	c, err := New()
	assert.Nil(t, err)
	namespace := []string{"TestContainerStatusNeverStarted"}
	status, err := c.Status(namespace)
	assert.Nil(t, err)
	assert.Equal(t, status, STOPPED)
}

func TestContainerStatusRunning(t *testing.T) {
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

func TestContainerStatusStopped(t *testing.T) {
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
