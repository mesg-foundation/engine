package core

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
				"test": {
					Image: "nginx",
				},
			},
		},
	})
	s, _ := services.Get(deployment.ServiceID)
	reply, err := serverstart.StartService(context.Background(), &StartServiceRequest{
		ServiceID: deployment.ServiceID,
	})
	assert.Nil(t, err)
	status, _ := s.Status()
	assert.Equal(t, service.RUNNING, status)
	assert.NotNil(t, reply)
	s.Stop()
	services.Delete(deployment.ServiceID)
}
