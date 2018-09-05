package core

import (
	"context"
	"testing"

	"github.com/mesg-foundation/core/service"
	"github.com/stretchr/testify/require"
)

func TestStartService(t *testing.T) {
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

	_, err = server.StartService(context.Background(), &StartServiceRequest{
		ServiceID: s.ID,
	})
	require.NoError(t, err)
	defer server.api.StopService(s.ID)

	status, err := s.Status()
	require.NoError(t, err)
	require.Equal(t, service.RUNNING, status)
}
