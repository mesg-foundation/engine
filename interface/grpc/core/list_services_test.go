package core

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestListServices(t *testing.T) {
	server := newServer(t)

	reply, err := server.ListServices(context.Background(), &ListServicesRequest{})
	require.NoError(t, err)

	services, err := server.api.ListServices()
	require.NoError(t, err)

	require.Equal(t, toProtoServices(services), reply.Services)
}
