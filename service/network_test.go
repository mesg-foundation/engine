package service

import (
	"testing"

	"github.com/stvp/assert"
)

func TestCreateNetwork(t *testing.T) {
	network, err := createNetwork("TestCreateNetwork")
	assert.Nil(t, err)
	assert.NotNil(t, network)
	deleteNetwork("TestCreateNetwork")
}

func TestFindNetwork(t *testing.T) {
	createNetwork("TestFindNetwork")
	network, err := findNetwork("TestFindNetwork")
	assert.Nil(t, err)
	assert.NotNil(t, network)
	deleteNetwork("TestFindNetwork")
}

func TestFindMissingNetwork(t *testing.T) {
	network, err := findNetwork("TestFindMissingNetwork")
	assert.Nil(t, err)
	assert.Nil(t, network)
}

func TestDeleteNetwork(t *testing.T) {
	createNetwork("TestDeleteNetwork")
	err := deleteNetwork("TestDeleteNetwork")
	assert.Nil(t, err)
	network, err := findNetwork("TestFindNetwork")
	assert.Nil(t, err)
	assert.Nil(t, network)
}
