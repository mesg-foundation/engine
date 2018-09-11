package core

import (
	"context"
	"testing"

	"github.com/mesg-foundation/core/service"
	"github.com/stretchr/testify/require"
)

func TestStopService(t *testing.T) {
	var (
		// we use a test service without tasks definition here otherwise we need to
		// spin up the gRPC server in order to prevent service exit with a failure
		// because it'll try to listen for tasks.
		path   = "../../../service-test/event"
		server = newServer(t)
	)

	s, validationErr, err := server.api.DeployService(serviceTar(t, path))
	require.Zero(t, validationErr)
	require.NoError(t, err)
	defer server.api.DeleteService(s.ID)

	require.NoError(t, server.api.StartService(s.ID))

	reply, err := server.StopService(context.Background(), &StopServiceRequest{
		ServiceID: s.ID,
	})

	status, err := s.Status()
	require.NoError(t, err)
	require.Equal(t, service.STOPPED, status)
	require.Nil(t, err)
	require.NotNil(t, reply)
}
