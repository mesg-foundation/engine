package client

import (
	"context"
	"testing"

	"github.com/mesg-foundation/core/service"
	"github.com/stvp/assert"
)

var serverstart = new(Server)

func TestStartService(t *testing.T) {
	service := service.Service{
		Name: "TestStartService",
		Dependencies: map[string]*service.Dependency{
			"test": &service.Dependency{
				Image: "nginx",
			},
		},
	}
	reply, err := serverstart.StartService(context.Background(), &ServiceRequest{
		Service: &service,
	})
	assert.Equal(t, service.IsRunning(), true)
	assert.Nil(t, err)
	assert.NotNil(t, reply)
	service.Stop()
}
