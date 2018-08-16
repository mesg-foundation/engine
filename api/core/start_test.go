package core

import (
	"context"
	"testing"

	"github.com/mesg-foundation/core/database/services"
	"github.com/mesg-foundation/core/service"
	"github.com/stretchr/testify/require"
)

var serverstart = new(Server)

func TestStartService(t *testing.T) {
	url := "https://github.com/mesg-foundation/service-webhook"

	server := newServer(t)
	stream := newTestDeployStream(url)
	server.DeployService(stream)

	s, _ := services.Get(stream.serviceID)
	reply, err := serverstart.StartService(context.Background(), &StartServiceRequest{
		ServiceID: stream.serviceID,
	})
	require.Nil(t, err)
	status, _ := s.Status()
	require.Equal(t, service.RUNNING, status)
	require.NotNil(t, reply)
	s.Stop()
	services.Delete(stream.serviceID)
}
