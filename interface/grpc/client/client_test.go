package client

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAPI(t *testing.T) {
	api, err := API()
	require.NoError(t, err)
	require.NotNil(t, api)
}

func TestGetClient(t *testing.T) {
	c, err := getClient()
	require.NoError(t, err)
	require.NotNil(t, c)
	require.NotNil(t, _client)
}
