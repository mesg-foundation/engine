// +build integration

package container

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIntegrationWaitForStatusRunning(t *testing.T) {
	c, err := New(nstestprefix)
	require.NoError(t, err)
	namespace := "TestWaitForStatusRunning"
	startTestService(namespace)
	defer c.StopService(namespace)
	err = c.waitForStatus(namespace, RUNNING)
	require.NoError(t, err)
}

func TestIntegrationWaitForStatusStopped(t *testing.T) {
	c, err := New(nstestprefix)
	require.NoError(t, err)
	namespace := "TestWaitForStatusStopped"
	startTestService(namespace)
	c.waitForStatus(namespace, RUNNING)
	c.StopService(namespace)
	err = c.waitForStatus(namespace, STOPPED)
	require.NoError(t, err)
}

func TestIntegrationWaitForStatusTaskError(t *testing.T) {
	c, err := New(nstestprefix)
	require.NoError(t, err)
	namespace := "TestWaitForStatusTaskError"
	c.StartService(ServiceOptions{
		Image:     "awgdaywudaywudwa",
		Namespace: namespace,
	})
	defer c.StopService(namespace)
	err = c.waitForStatus(namespace, RUNNING)
	require.Error(t, err)
	require.Contains(t, err.Error(), "No such image: awgdaywudaywudwa:latest")
}
