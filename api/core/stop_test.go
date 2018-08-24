package core

import (
	"context"
	"testing"

	"github.com/mesg-foundation/core/database/services"
	"github.com/mesg-foundation/core/service"
	"github.com/stretchr/testify/require"
)

var serverstop = new(Server)

func TestStopService(t *testing.T) {
	url := "https://github.com/mesg-foundation/service-webhook"

	server := newServer(t)
	stream := newTestDeployStream(url)
	server.DeployService(stream)

	s, _ := services.Get(stream.serviceID)
	s.Start()
	reply, err := serverstop.StopService(context.Background(), &StopServiceRequest{
		ServiceID: stream.serviceID,
	})
	status, _ := s.Status()
	require.Equal(t, service.STOPPED, status)
	require.Nil(t, err)
	require.NotNil(t, reply)
	services.Delete(stream.serviceID)
}
