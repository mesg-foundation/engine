// +build integration

package core

import (
	"context"
	"testing"

	"github.com/mesg-foundation/core/config"
	"github.com/mesg-foundation/core/protobuf/coreapi"
	"github.com/stretchr/testify/require"
)

func TestExecuteWithNonExistingService(t *testing.T) {
	server, closer := newServer(t)
	defer closer()

	_, err := server.ExecuteTask(context.Background(), &coreapi.ExecuteTaskRequest{
		ServiceID: "-",
		TaskKey:   "error",
		InputData: "{}",
	})
	require.Error(t, err)
}

func TestInfo(t *testing.T) {
	server, closer := newServer(t)
	defer closer()

	c, err := config.Global()
	require.NoError(t, err)
	reply, err := server.Info(context.Background(), &coreapi.InfoRequest{})
	require.NoError(t, err)
	require.NotNil(t, reply)
	for i, s := range reply.Services {
		require.Equal(t, s.Sid, c.Services()[i].Sid)
	}
}
