package core

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetService(t *testing.T) {
	var server = newServer(t)

	s, validationErr, err := server.api.DeployService(serviceTar(t, taskServicePath))
	require.Zero(t, validationErr)
	require.NoError(t, err)
	defer server.api.DeleteService(s.Id)

	reply, err := server.GetService(context.Background(), &GetServiceRequest{
		ServiceID: s.Id,
	})
	require.Nil(t, err)
	require.NotNil(t, reply)
	require.Equal(t, reply.Service.Name, "Task")
}
