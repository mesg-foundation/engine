package core

import (
	"context"
	"testing"

	"github.com/mesg-foundation/core/database/services"
	"github.com/mesg-foundation/core/service"
	"github.com/stretchr/testify/require"
)

var serverdelete = new(Server)

func TestDeleteService(t *testing.T) {
	emptyService := service.Service{}

	url := "https://github.com/mesg-foundation/service-webhook"

	server := newServer(t)
	stream := newTestDeployStream(url)
	server.DeployService(stream)

	reply, err := serverdelete.DeleteService(context.Background(), &DeleteServiceRequest{
		ServiceID: stream.serviceID,
	})
	require.Nil(t, err)
	require.NotNil(t, reply)
	x, _ := services.Get(stream.serviceID)
	require.Equal(t, emptyService, x)
}
