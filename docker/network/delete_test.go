package network

import (
	"testing"

	"github.com/stvp/assert"
)

func TestDeleteNetwork(t *testing.T) {
	Create("TestDeleteNetwork")
	err := Delete("TestDeleteNetwork")
	assert.Nil(t, err)
	network, err := Find("TestFindNetwork")
	assert.Nil(t, err)
	assert.Nil(t, network)
}
