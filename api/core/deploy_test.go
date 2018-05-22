package core

import (
	"context"
	"testing"

	"github.com/mesg-foundation/core/database/services"
	"github.com/mesg-foundation/core/service"
	"github.com/stvp/assert"
)

var serverdeploy = new(Server)

func TestDeployService(t *testing.T) {
	service := service.Service{
		Name: "TestDeployService",
		Dependencies: map[string]*service.Dependency{
			"test": &service.Dependency{
				Image: "nginx",
			},
		},
	}
	deployment, err := serverdeploy.DeployService(context.Background(), &DeployServiceRequest{
		Service: &service,
	})
	assert.Nil(t, err)
	assert.NotNil(t, deployment.ServiceID)
	services.Delete(deployment.ServiceID)
}
