package api

import (
	"testing"

	"github.com/stvp/assert"
)

func TestClientDefaultConfig(t *testing.T) {
	c := Client{}
	assert.NotEqual(t, c.target(), "", "target should not be empty")
}

func TestClientClose(t *testing.T) {
	c := Client{}
	_, err := c.conn()
	assert.Nil(t, err)
	c.Close()
}

func TestClientService(t *testing.T) {
	c := Client{}
	client, err := c.ServiceClient()
	assert.Nil(t, err)
	assert.NotNil(t, client)
}
