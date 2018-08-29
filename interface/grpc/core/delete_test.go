package core

import (
	"context"
	"fmt"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestDeleteService(t *testing.T) {
	var (
		path   = filepath.Join("..", "..", "..", "service-test", "task")
		server = newServer(t)
	)

	fmt.Println("before deploy", time.Now())
	s, validationErr, err := server.api.DeployService(serviceTar(t, path))
	require.Zero(t, validationErr)
	require.NoError(t, err)
	fmt.Println("after deploy", time.Now())

	reply, err := server.DeleteService(context.Background(), &DeleteServiceRequest{
		ServiceID: s.Id,
	})
	fmt.Println("after delete", time.Now())
	require.Nil(t, err)
	require.NotNil(t, reply)

	_, err = server.api.GetService(s.Id)
	fmt.Println("after get", time.Now())
	require.Error(t, err)
}
