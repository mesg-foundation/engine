package client

import (
	"testing"

	"github.com/stvp/assert"
)

func TestAPI(t *testing.T) {
	api := API()
	assert.NotNil(t, api)
}

func TestGetClient(t *testing.T) {
	c := getClient()
	assert.NotNil(t, c)
	assert.NotNil(t, _client)
}
