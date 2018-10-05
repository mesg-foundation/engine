// +build integration

package core

import (
	"testing"

	"github.com/mesg-foundation/core/api"
	"github.com/stretchr/testify/require"
)

func TestIntegrationDeployService(t *testing.T) {
	url := "https://github.com/mesg-foundation/service-webhook"

	server, closer := newServer(t)
	defer closer()
	stream := newTestDeployStream(url)

	require.Nil(t, server.DeployService(stream))
	defer server.api.DeleteService(stream.serviceID)

	require.Len(t, stream.serviceID, 40)
	require.Contains(t, stream.statuses, api.DeployStatus{
		Message: "Image built with success",
		Type:    api.DonePositive,
	})
}
