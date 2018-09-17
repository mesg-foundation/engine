package core

import (
	"context"
	"testing"

	"github.com/mesg-foundation/core/protobuf/core"
	"github.com/stretchr/testify/require"
)

func TestListServices(t *testing.T) {
	server := newServer(t)

	reply, err := server.ListServices(context.Background(), &core.ListServicesRequest{})
	require.NoError(t, err)

	services, err := server.api.ListServices()
	require.NoError(t, err)

	require.Equal(t, toProtoServices(services), reply.Services)
}
