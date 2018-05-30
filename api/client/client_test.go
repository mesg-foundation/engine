package client

import (
	"testing"

	"github.com/stvp/assert"
)

func TestAPI(t *testing.T) {
	api, err := API()
	assert.Nil(t, err)
	assert.NotNil(t, api)
}

func TestGetClient(t *testing.T) {
	c, err := getClient()
	assert.Nil(t, err)
	assert.NotNil(t, c)
	assert.NotNil(t, _client)
}
