// +build integration

package core

import (
	"fmt"
	"testing"

	"github.com/cnf/structhash"
	"github.com/logrusorgru/aurora"
	"github.com/mesg-foundation/core/api"
	"github.com/stretchr/testify/require"
)

func TestIntegrationDeployService(t *testing.T) {
	url := "https://github.com/mesg-foundation/service-webhook"

	server := newServer(t)
	stream := newTestDeployStream(url)

	require.Nil(t, server.DeployService(stream))
	defer server.api.DeleteService(stream.serviceID)

	require.Equal(t, 1, structhash.Version(stream.serviceID))
	require.Contains(t, stream.statuses, api.DeployStatus{
		Message: fmt.Sprintf("%s Completed.", aurora.Green("✔")),
		Type:    api.DONE,
	})
}
