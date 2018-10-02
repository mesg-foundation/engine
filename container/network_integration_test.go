// +build integration

package container

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIntegrationCreateNetwork(t *testing.T) {
	c, err := New()
	require.Nil(t, err)
	networkID, err := c.CreateNetwork([]string{"TestCreateNetwork"})
	defer c.DeleteNetwork([]string{"TestCreateNetwork"}, EventRemove)
	require.Nil(t, err)
	require.NotEqual(t, "", networkID)
}

func TestIntegrationCreateAlreadyExistingNetwork(t *testing.T) {
	c, err := New()
	require.Nil(t, err)
	c.CreateNetwork([]string{"TestCreateAlreadyExistingNetwork"})
	networkID, err := c.CreateNetwork([]string{"TestCreateAlreadyExistingNetwork"})
	defer c.DeleteNetwork([]string{"TestCreateAlreadyExistingNetwork"}, EventRemove)
	require.Nil(t, err)
	require.NotEqual(t, "", networkID)
}

func TestIntegrationDeleteNetwork(t *testing.T) {
	c, err := New()
	require.Nil(t, err)
	c.CreateNetwork([]string{"TestDeleteNetwork"})
	err = c.DeleteNetwork([]string{"TestDeleteNetwork"}, EventRemove)
	require.Nil(t, err)
}

func TestIntegrationDeleteNotExistingNetwork(t *testing.T) {
	c, err := New()
	require.Nil(t, err)
	err = c.DeleteNetwork([]string{"TestDeleteNotExistingNetwork"}, EventRemove)
	require.Nil(t, err)
}

func TestIntegrationFindNetwork(t *testing.T) {
	c, err := New()
	require.Nil(t, err)
	c.CreateNetwork([]string{"TestFindNetwork"})
	defer c.DeleteNetwork([]string{"TestFindNetwork"}, EventRemove)
	network, err := c.FindNetwork([]string{"TestFindNetwork"})
	require.Nil(t, err)
	require.NotEqual(t, "", network.ID)
}

func TestIntegrationFindNotExistingNetwork(t *testing.T) {
	c, err := New()
	require.Nil(t, err)
	_, err = c.FindNetwork([]string{"TestFindNotExistingNetwork"})
	require.NotNil(t, err)
}

func TestIntegrationFindDeletedNetwork(t *testing.T) {
	c, err := New()
	require.Nil(t, err)
	c.CreateNetwork([]string{"TestFindDeletedNetwork"})
	c.DeleteNetwork([]string{"TestFindDeletedNetwork"}, EventRemove)
	_, err = c.FindNetwork([]string{"TestFindDeletedNetwork"})
	require.NotNil(t, err)
}
