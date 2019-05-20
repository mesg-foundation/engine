// +build integration

package core

import (
	"testing"

	"github.com/mesg-foundation/core/api"
	"github.com/stretchr/testify/require"
)

func TestIntegrationDeployService(t *testing.T) {
	url := "git://github.com/mesg-foundation/service-webhook#single-outputs"

	server, closer := newServer(t)
	defer closer()
	stream := newTestDeployStream(url)

	require.Nil(t, server.DeployService(stream))
	defer server.api.DeleteService(stream.hash, false)

	require.Len(t, stream.sid, 7)
	require.NotEmpty(t, stream.hash)
	require.Contains(t, stream.statuses, api.DeployStatus{
		Message: "Image built with success",
		Type:    api.DonePositive,
	})
}
