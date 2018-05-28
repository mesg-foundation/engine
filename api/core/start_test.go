package core

import (
	"context"
	"testing"

	"github.com/mesg-foundation/core/daemon"
	"github.com/mesg-foundation/core/database/services"
	"github.com/mesg-foundation/core/service"
	"github.com/stvp/assert"
)

var serverstart = new(Server)

func TestStartService(t *testing.T) {
	daemon.Start()
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
	assert.NotNil(t, reply)
	running, err := service.IsRunning()
	assert.Nil(t, err)
	assert.Equal(t, true, running)
	service.Stop()
	services.Delete(deployment.ServiceID)
}
