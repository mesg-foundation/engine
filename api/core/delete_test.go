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
	service := service.Service{
		Name: "TestDeleteService",
		Dependencies: map[string]*service.Dependency{
			"test": {
				Image: "nginx",
			},
		},
	}
	deployment, _ := serverdelete.DeployService(context.Background(), &DeployServiceRequest{
		Service: &service,
	})
	reply, err := serverdelete.DeleteService(context.Background(), &DeleteServiceRequest{
		ServiceID: deployment.ServiceID,
	})
	require.Nil(t, err)
	require.NotNil(t, reply)
	x, _ := services.Get(deployment.ServiceID)
	require.Equal(t, x, emptyService)
}
