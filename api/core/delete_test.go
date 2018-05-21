package core

import (
	"context"
	"testing"

	"github.com/mesg-foundation/core/database/services"
	"github.com/mesg-foundation/core/service"
	"github.com/stvp/assert"
)

var serverdelete = new(Server)

func TestDeleteService(t *testing.T) {
	emptyService := service.Service{}
	service := service.Service{
		Name: "TestDeleteService",
		Dependencies: map[string]*service.Dependency{
			"test": &service.Dependency{
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
	assert.Nil(t, err)
	assert.NotNil(t, reply)
	x, _ := services.Get(deployment.ServiceID)
	assert.Equal(t, x, emptyService)
}
