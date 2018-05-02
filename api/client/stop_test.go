package client

import (
	"context"
	"testing"

	"github.com/mesg-foundation/core/service"
	"github.com/stvp/assert"
)

var serverstop = new(Server)

func TestStopService(t *testing.T) {
	service := service.Service{
		Name: "TestStopService",
		Dependencies: map[string]*service.Dependency{
			"test": &service.Dependency{
				Image: "nginx",
			},
		},
	}
	service.Start()
	reply, err := serverstop.StopService(context.Background(), &ServiceRequest{
		Service: &service,
	})
	assert.Equal(t, service.IsRunning(), false)
	assert.Nil(t, err)
	assert.NotNil(t, reply)
}
