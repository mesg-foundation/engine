// +build integration

package container

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIntegrationCreateNetwork(t *testing.T) {
	c, err := New()
	require.NoError(t, err)
	networkID, err := c.CreateNetwork([]string{"TestCreateNetwork"})
	defer c.DeleteNetwork([]string{"TestCreateNetwork"})
	require.NoError(t, err)
	require.NotEqual(t, "", networkID)
}

func TestIntegrationCreateAlreadyExistingNetwork(t *testing.T) {
	c, err := New()
	require.NoError(t, err)
	c.CreateNetwork([]string{"TestCreateAlreadyExistingNetwork"})
	networkID, err := c.CreateNetwork([]string{"TestCreateAlreadyExistingNetwork"})
	defer c.DeleteNetwork([]string{"TestCreateAlreadyExistingNetwork"})
	require.NoError(t, err)
	require.NotEqual(t, "", networkID)
}

func TestIntegrationDeleteNetwork(t *testing.T) {
	c, err := New()
	require.NoError(t, err)
	c.CreateNetwork([]string{"TestDeleteNetwork"})
	err = c.DeleteNetwork([]string{"TestDeleteNetwork"})
	require.NoError(t, err)
}

func TestIntegrationDeleteNotExistingNetwork(t *testing.T) {
	c, err := New()
	require.NoError(t, err)
	err = c.DeleteNetwork([]string{"TestDeleteNotExistingNetwork"})
	require.NoError(t, err)
}

func TestIntegrationFindNetwork(t *testing.T) {
	c, err := New()
	require.NoError(t, err)
	c.CreateNetwork([]string{"TestFindNetwork"})
	defer c.DeleteNetwork([]string{"TestFindNetwork"})
	network, err := c.FindNetwork([]string{"TestFindNetwork"})
	require.NoError(t, err)
	require.NotEqual(t, "", network.ID)
}

func TestIntegrationFindNotExistingNetwork(t *testing.T) {
	c, err := New()
	require.NoError(t, err)
	_, err = c.FindNetwork([]string{"TestFindNotExistingNetwork"})
	require.Error(t, err)
}

func TestIntegrationFindDeletedNetwork(t *testing.T) {
	c, err := New()
	require.NoError(t, err)
	c.CreateNetwork([]string{"TestFindDeletedNetwork"})
	c.DeleteNetwork([]string{"TestFindDeletedNetwork"})
	_, err = c.FindNetwork([]string{"TestFindDeletedNetwork"})
	require.Error(t, err)
}
