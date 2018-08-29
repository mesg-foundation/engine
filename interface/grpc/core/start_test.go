package core

import (
	"context"
	"testing"

	"github.com/mesg-foundation/core/service"
	"github.com/stretchr/testify/require"
)

func TestStartService(t *testing.T) {
	var server = newServer(t)

	// we use a test service without tasks definition here otherwise we need to
	// spin up the gRPC server in order to prevent service exit with a failure
	// because it'll try to listen for tasks.
	s, validationErr, err := server.api.DeployService(serviceTar(t, eventServicePath))
	require.Zero(t, validationErr)
	require.NoError(t, err)
	defer server.api.DeleteService(s.Id)

	_, err = server.StartService(context.Background(), &StartServiceRequest{
		ServiceID: s.Id,
	})
	require.NoError(t, err)
	defer server.api.StopService(s.Id)

	s.Id = "" // TODO(ilgooz) remove this when Service type created by hand.
	status, err := s.Status()
	require.NoError(t, err)
	require.Equal(t, service.RUNNING, status)
}
