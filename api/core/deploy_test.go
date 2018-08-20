package core

import (
	"context"
	"testing"

	"github.com/mesg-foundation/core/database/services"
	"github.com/mesg-foundation/core/service"
	"github.com/stretchr/testify/require"
)

var serverdeploy = new(Server)

func TestDeployService(t *testing.T) {
	service := service.Service{
		Name: "TestDeployService",
		Dependencies: map[string]*service.Dependency{
			"test": {
				Image: "nginx",
			},
		},
	}
	deployment, err := serverdeploy.DeployService(context.Background(), &DeployServiceRequest{
		Service: &service,
	})
	require.Nil(t, err)
	require.NotNil(t, deployment)
	require.NotEqual(t, "", deployment.ServiceID)
	services.Delete(deployment.ServiceID)
}
