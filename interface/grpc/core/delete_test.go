package core

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/mesg-foundation/core/protobuf/coreapi"
	"github.com/stretchr/testify/require"
)

func TestDeleteService(t *testing.T) {
	var (
		path           = filepath.Join("..", "..", "..", "service-test", "task")
		server, closer = newServer(t)
	)
	defer closer()

	s, validationErr, err := server.api.DeployService(serviceTar(t, path))
	require.Zero(t, validationErr)
	require.NoError(t, err)

	reply, err := server.DeleteService(context.Background(), &coreapi.DeleteServiceRequest{
		ServiceID: s.ID,
	})
	require.NoError(t, err)
	require.NotNil(t, reply)

	_, err = server.api.GetService(s.ID)
	require.Error(t, err)
}
