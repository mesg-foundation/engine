// +build integration

package core

import (
	"fmt"
	"testing"

	"github.com/cnf/structhash"
	"github.com/logrusorgru/aurora"
	"github.com/mesg-foundation/core/database/services"
	"github.com/mesg-foundation/core/mesg"
	"github.com/stretchr/testify/require"
)

func TestIntegrationDeployService(t *testing.T) {
	url := "https://github.com/mesg-foundation/service-webhook"

	server := newServer(t)
	stream := newTestDeployStream(url)

	require.Nil(t, server.DeployService(stream))
	defer services.Delete(stream.serviceID)

	require.Equal(t, 1, structhash.Version(stream.serviceID))
	require.Contains(t, stream.statuses, mesg.DeployStatus{
		Message: fmt.Sprintf("%s Completed.", aurora.Green("✔")),
		Type:    mesg.DONE,
	})
}
