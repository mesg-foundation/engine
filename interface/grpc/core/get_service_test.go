package core

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetService(t *testing.T) {
	server, closer := newServer(t)
	defer closer()

	s, validationErr, err := server.api.DeployService(serviceTar(t, taskServicePath))
	require.Zero(t, validationErr)
	require.NoError(t, err)
	defer server.api.DeleteService(s.ID)

	reply, err := server.GetService(context.Background(), &GetServiceRequest{
		ServiceID: s.ID,
	})
	require.Nil(t, err)
	require.NotNil(t, reply)
	require.Equal(t, reply.Service.Name, "Task")
}
