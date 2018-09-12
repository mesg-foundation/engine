package core

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDeleteService(t *testing.T) {
	var (
		path   = filepath.Join("..", "..", "..", "service-test", "task")
		server = newServer(t)
	)

	s, validationErr, err := server.api.DeployService(serviceTar(t, path))
	require.Zero(t, validationErr)
	require.NoError(t, err)

	reply, err := server.DeleteService(context.Background(), &DeleteServiceRequest{
		ServiceID: s.ID,
	})
	require.Nil(t, err)
	require.NotNil(t, reply)

	_, err = server.api.GetService(s.ID)
	require.Error(t, err)
}
