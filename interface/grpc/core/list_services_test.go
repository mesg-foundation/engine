package core

import (
	"context"
	"testing"

	"github.com/mesg-foundation/core/protobuf/coreapi"
	"github.com/stretchr/testify/require"
)

func TestListServices(t *testing.T) {
	url := "https://github.com/mesg-foundation/service-webhook"
	server, closer := newServer(t)
	defer closer()

	stream := newTestDeployStream(url)
	require.NoError(t, server.DeployService(stream))
	defer server.api.DeleteService(stream.serviceID)

	reply, err := server.ListServices(context.Background(), &coreapi.ListServicesRequest{})
	require.NoError(t, err)

	services, err := server.api.ListServices()
	require.NoError(t, err)

	apiProtoServices := toProtoServices(services)

	require.Len(t, apiProtoServices, 1)
	require.Equal(t, reply.Services[0].ID, apiProtoServices[0].ID)
}
