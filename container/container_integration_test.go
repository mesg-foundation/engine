// +build integration

package container

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIntegrationFindContainerNotExisting(t *testing.T) {
	c, err := New(nstestprefix)
	require.NoError(t, err)
	_, err = c.FindContainer("TestFindContainerNotExisting")
	require.Error(t, err)
}

func TestIntegrationFindContainer(t *testing.T) {
	c, err := New(nstestprefix)
	require.NoError(t, err)
	namespace := "TestFindContainer"
	startTestService(namespace)
	defer c.StopService(namespace)
	c.waitForStatus(namespace, RUNNING)
	container, err := c.FindContainer(namespace)
	require.NoError(t, err)
	require.NotEqual(t, "", container.ID)
}

func TestIntegrationFindContainerStopped(t *testing.T) {
	c, err := New(nstestprefix)
	require.NoError(t, err)
	namespace := "TestFindContainerStopped"
	startTestService(namespace)
	c.StopService(namespace)
	_, err = c.FindContainer(namespace)
	require.Error(t, err)
}

func TestIntegrationContainerStatusNeverStarted(t *testing.T) {
	c, err := New(nstestprefix)
	require.NoError(t, err)
	namespace := "TestContainerStatusNeverStarted"
	status, err := c.Status(namespace)
	require.NoError(t, err)
	require.Equal(t, status, STOPPED)
}

func TestIntegrationContainerStatusRunning(t *testing.T) {
	c, err := New(nstestprefix)
	require.NoError(t, err)
	namespace := "TestContainerStatusRunning"
	startTestService(namespace)
	defer c.StopService(namespace)
	c.waitForStatus(namespace, RUNNING)
	status, err := c.Status(namespace)
	require.NoError(t, err)
	require.Equal(t, status, RUNNING)
}

func TestIntegrationContainerStatusStopped(t *testing.T) {
	c, err := New(nstestprefix)
	require.NoError(t, err)
	namespace := "TestContainerStatusStopped"
	startTestService(namespace)
	c.waitForStatus(namespace, RUNNING)
	c.StopService(namespace)
	c.waitForStatus(namespace, STOPPED)
	status, err := c.Status(namespace)
	require.NoError(t, err)
	require.Equal(t, status, STOPPED)
}
