// +build integration

package container

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIntegrationCreateSwarmIfNeeded(t *testing.T) {
	c, err := New()
	require.Nil(t, err)
	require.Nil(t, c.createSwarmIfNeeded())
}

func TestIntegrationFindContainerNotExisting(t *testing.T) {
	c, err := New()
	require.Nil(t, err)
	_, err = c.FindContainer([]string{"TestFindContainerNotExisting"})
	require.NotNil(t, err)
}

func TestIntegrationFindContainer(t *testing.T) {
	c, err := New()
	require.Nil(t, err)
	namespace := []string{"TestFindContainer"}
	startTestService(namespace)
	defer c.StopService(namespace)
	c.waitForStatus(namespace, RUNNING)
	container, err := c.FindContainer(namespace)
	require.Nil(t, err)
	require.NotEqual(t, "", container.ID)
}

func TestIntegrationFindContainerStopped(t *testing.T) {
	c, err := New()
	require.Nil(t, err)
	namespace := []string{"TestFindContainerStopped"}
	startTestService(namespace)
	c.StopService(namespace)
	_, err = c.FindContainer(namespace)
	require.NotNil(t, err)
}

func TestIntegrationContainerStatusNeverStarted(t *testing.T) {
	c, err := New()
	require.Nil(t, err)
	namespace := []string{"TestContainerStatusNeverStarted"}
	status, err := c.Status(namespace)
	require.Nil(t, err)
	require.Equal(t, status, STOPPED)
}

func TestIntegrationContainerStatusRunning(t *testing.T) {
	c, err := New()
	require.Nil(t, err)
	namespace := []string{"TestContainerStatusRunning"}
	startTestService(namespace)
	defer c.StopService(namespace)
	c.waitForStatus(namespace, RUNNING)
	status, err := c.Status(namespace)
	require.Nil(t, err)
	require.Equal(t, status, RUNNING)
}

func TestIntegrationContainerStatusStopped(t *testing.T) {
	c, err := New()
	require.Nil(t, err)
	namespace := []string{"TestContainerStatusStopped"}
	startTestService(namespace)
	c.waitForStatus(namespace, RUNNING)
	c.StopService(namespace)
	c.waitForStatus(namespace, STOPPED)
	status, err := c.Status(namespace)
	require.Nil(t, err)
	require.Equal(t, status, STOPPED)
}
