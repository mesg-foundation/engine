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
	reply, err := serverstart.StartService(context.Background(), &StartServiceRequest{
		Service: &service,
	})
	assert.Nil(t, err)
	assert.True(t, service.IsRunning())
	assert.NotNil(t, reply)
	service.Stop()
}
