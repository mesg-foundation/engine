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
	hash, _ := services.Save(&service.Service{
		Name: "TestGetService",
	})
	defer services.Delete(hash)
	reply, err := servergetservice.GetService(context.Background(), &GetServiceRequest{
		ServiceID: hash,
	})
	require.Nil(t, err)
	require.NotNil(t, reply)
	require.Equal(t, reply.Service.Name, "TestGetService")
}
