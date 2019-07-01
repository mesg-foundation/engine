// +build integration

package container

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIntegrationCreateNetwork(t *testing.T) {
	c, err := New(nstestprefix)
	require.NoError(t, err)
	networkID, err := c.CreateNetwork("TestCreateNetwork")
	defer c.DeleteNetwork("TestCreateNetwork")
	require.NoError(t, err)
	require.NotEqual(t, "", networkID)
}

func TestIntegrationCreateAlreadyExistingNetwork(t *testing.T) {
	c, err := New(nstestprefix)
	require.NoError(t, err)
	c.CreateNetwork("TestCreateAlreadyExistingNetwork")
	networkID, err := c.CreateNetwork("TestCreateAlreadyExistingNetwork")
	defer c.DeleteNetwork("TestCreateAlreadyExistingNetwork")
	require.NoError(t, err)
	require.NotEqual(t, "", networkID)
}

func TestIntegrationDeleteNetwork(t *testing.T) {
	c, err := New(nstestprefix)
	require.NoError(t, err)
	c.CreateNetwork("TestDeleteNetwork")
	err = c.DeleteNetwork("TestDeleteNetwork")
	require.NoError(t, err)
}

func TestIntegrationDeleteNotExistingNetwork(t *testing.T) {
	c, err := New(nstestprefix)
	require.NoError(t, err)
	err = c.DeleteNetwork("TestDeleteNotExistingNetwork")
	require.NoError(t, err)
}

func TestIntegrationFindNetwork(t *testing.T) {
	c, err := New(nstestprefix)
	require.NoError(t, err)
	c.CreateNetwork("TestFindNetwork")
	defer c.DeleteNetwork("TestFindNetwork")
	network, err := c.FindNetwork("TestFindNetwork")
	require.NoError(t, err)
	require.NotEqual(t, "", network.ID)
}

func TestIntegrationFindNotExistingNetwork(t *testing.T) {
	c, err := New(nstestprefix)
	require.NoError(t, err)
	_, err = c.FindNetwork("TestFindNotExistingNetwork")
	require.Error(t, err)
}

func TestIntegrationFindDeletedNetwork(t *testing.T) {
	c, err := New(nstestprefix)
	require.NoError(t, err)
	c.CreateNetwork("TestFindDeletedNetwork")
	c.DeleteNetwork("TestFindDeletedNetwork")
	_, err = c.FindNetwork("TestFindDeletedNetwork")
	require.Error(t, err)
}
