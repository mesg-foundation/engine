package service

import (
	"testing"

	"github.com/mesg-foundation/core/api"
	"github.com/stretchr/testify/require"
)

func newServer(t *testing.T) *Server {
	a, err := api.New()
	require.Nil(t, err)

	server, err := NewServer(APIOption(a))
	require.Nil(t, err)

	return server
}
