package core

import (
	"context"
	"testing"

	"github.com/mesg-foundation/core/protobuf/coreapi"
	"github.com/mesg-foundation/core/service"
	"github.com/stretchr/testify/require"
)

func TestStopService(t *testing.T) {
	server, closer := newServer(t)
	defer closer()

	// we use a test service without tasks definition here otherwise we need to
	// spin up the gRPC server in order to prevent service exit with a failure
	// because it'll try to listen for tasks.
	s, validationErr, err := server.api.DeployService(serviceTar(t, eventServicePath))
	require.Zero(t, validationErr)
	require.NoError(t, err)
	defer server.api.DeleteService(s.ID)

	require.NoError(t, server.api.StartService(s.ID))

	reply, err := server.StopService(context.Background(), &coreapi.StopServiceRequest{
		ServiceID: s.ID,
	})
	require.NoError(t, err)

	status, err := s.Status()
	require.NoError(t, err)
	require.Equal(t, service.STOPPED, status)
	require.NoError(t, err)
	require.NotNil(t, reply)
}
