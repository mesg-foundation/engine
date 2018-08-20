package core

import (
	"context"
	"testing"

	"github.com/mesg-foundation/core/database/services"
	"github.com/mesg-foundation/core/service"
	"github.com/stretchr/testify/require"
)

var servergetservice = new(Server)

func TestGetService(t *testing.T) {
	service := &service.Service{
		Name: "TestGetService",
	}
	services.Save(service)
	defer services.Delete(service.Id)
	reply, err := servergetservice.GetService(context.Background(), &GetServiceRequest{
		ServiceID: service.Id,
	})
	require.Nil(t, err)
	require.NotNil(t, reply)
	require.Equal(t, reply.Service.Name, "TestGetService")
}
