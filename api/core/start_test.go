package core

/*
import (
	"context"
	"testing"

	"github.com/mesg-foundation/core/database/services"

	"github.com/mesg-foundation/core/service"
	"github.com/stvp/assert"
)

var serverstart = new(Server)

func TestStartService(t *testing.T) {
	deployment, _ := serverstart.DeployService(context.Background(), &DeployServiceRequest{
		Service: &service.Service{
			Name: "TestStartService",
			Dependencies: map[string]*service.Dependency{
				"test": &service.Dependency{
					Image: "nginx",
				},
			},
		},
	})
	service, _ := services.Get(deployment.ServiceID)
	reply, err := serverstart.StartService(context.Background(), &StartServiceRequest{
		ServiceID: deployment.ServiceID,
	})
	assert.Nil(t, err)
	assert.True(t, service.IsRunning())
	assert.NotNil(t, reply)
	service.Stop()
	services.Delete(deployment.ServiceID)
}
*/
