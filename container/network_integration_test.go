// +build integration

package container

import (
	"testing"

	"github.com/stvp/assert"
)

func TestIntegrationCreateNetwork(t *testing.T) {
	c, err := New()
	assert.Nil(t, err)
	networkID, err := c.CreateNetwork([]string{"TestCreateNetwork"})
	defer c.DeleteNetwork([]string{"TestCreateNetwork"})
	assert.Nil(t, err)
	assert.NotEqual(t, "", networkID)
}

func TestIntegrationCreateAlreadyExistingNetwork(t *testing.T) {
	c, err := New()
	assert.Nil(t, err)
	c.CreateNetwork([]string{"TestCreateAlreadyExistingNetwork"})
	networkID, err := c.CreateNetwork([]string{"TestCreateAlreadyExistingNetwork"})
	defer c.DeleteNetwork([]string{"TestCreateAlreadyExistingNetwork"})
	assert.Nil(t, err)
	assert.NotEqual(t, "", networkID)
}

func TestIntegrationDeleteNetwork(t *testing.T) {
	c, err := New()
	assert.Nil(t, err)
	c.CreateNetwork([]string{"TestDeleteNetwork"})
	err = c.DeleteNetwork([]string{"TestDeleteNetwork"})
	assert.Nil(t, err)
}

func TestIntegrationDeleteNotExistingNetwork(t *testing.T) {
	c, err := New()
	assert.Nil(t, err)
	err = c.DeleteNetwork([]string{"TestDeleteNotExistingNetwork"})
	assert.Nil(t, err)
}

func TestIntegrationFindNetwork(t *testing.T) {
	c, err := New()
	assert.Nil(t, err)
	c.CreateNetwork([]string{"TestFindNetwork"})
	defer c.DeleteNetwork([]string{"TestFindNetwork"})
	network, err := c.FindNetwork([]string{"TestFindNetwork"})
	assert.Nil(t, err)
	assert.NotEqual(t, "", network.ID)
}

func TestIntegrationFindNotExistingNetwork(t *testing.T) {
	c, err := New()
	assert.Nil(t, err)
	_, err = c.FindNetwork([]string{"TestFindNotExistingNetwork"})
	assert.NotNil(t, err)
}

func TestIntegrationFindDeletedNetwork(t *testing.T) {
	c, err := New()
	assert.Nil(t, err)
	c.CreateNetwork([]string{"TestFindDeletedNetwork"})
	c.DeleteNetwork([]string{"TestFindDeletedNetwork"})
	_, err = c.FindNetwork([]string{"TestFindDeletedNetwork"})
	assert.NotNil(t, err)
}
