// +build integration

package core

import (
	"fmt"
	"testing"

	"github.com/cnf/structhash"
	"github.com/logrusorgru/aurora"
	"github.com/mesg-foundation/core/database/services"
	"github.com/stretchr/testify/require"
)

func TestIntegrationDeployService(t *testing.T) {
	url := "https://github.com/mesg-foundation/service-webhook"

	server := newServer(t)
	stream := newTestDeployStream(url)

	require.Nil(t, server.DeployService(stream))
	require.Equal(t, 1, structhash.Version(stream.serviceID))
	require.True(t, stringSliceContains(stream.statuses, fmt.Sprintf("%s Completed.", aurora.Green("✔"))))
	services.Delete(stream.serviceID)
}
