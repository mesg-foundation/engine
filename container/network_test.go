package container

import (
	"testing"

	"github.com/stvp/assert"
)

func TestCreateNetwork(t *testing.T) {
	networkID, err := CreateNetwork([]string{"TestCreateNetwork"})
	assert.Nil(t, err)
	assert.NotEqual(t, "", networkID)
	DeleteNetwork([]string{"TestCreateNetwork"})
}

func TestCreateAlreadyExistingNetwork(t *testing.T) {
	CreateNetwork([]string{"TestCreateAlreadyExistingNetwork"})
	networkID, err := CreateNetwork([]string{"TestCreateAlreadyExistingNetwork"})
	assert.Nil(t, err)
	assert.NotEqual(t, "", networkID)
	DeleteNetwork([]string{"TestCreateAlreadyExistingNetwork"})
}

func TestDeleteNetwork(t *testing.T) {
	CreateNetwork([]string{"TestDeleteNetwork"})
	err := DeleteNetwork([]string{"TestDeleteNetwork"})
	assert.Nil(t, err)
}

func TestDeleteNotExistingNetwork(t *testing.T) {
	err := DeleteNetwork([]string{"TestDeleteNotExistingNetwork"})
	assert.Nil(t, err)
}

func TestFindNetwork(t *testing.T) {
	CreateNetwork([]string{"TestFindNetwork"})
	network, err := FindNetwork([]string{"TestFindNetwork"})
	assert.Nil(t, err)
	assert.NotNil(t, network)
	DeleteNetwork([]string{"TestFindNetwork"})
}

func TestFindNotExistingNetwork(t *testing.T) {
	_, err := FindNetwork([]string{"TestFindNotExistingNetwork"})
	assert.NotNil(t, err)
}

func TestFindDeletedNetwork(t *testing.T) {
	CreateNetwork([]string{"TestFindDeletedNetwork"})
	DeleteNetwork([]string{"TestFindDeletedNetwork"})
	_, err := FindNetwork([]string{"TestFindDeletedNetwork"})
	assert.NotNil(t, err)
}
