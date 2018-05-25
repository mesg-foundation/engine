package core

import (
	"context"
	"testing"

	"github.com/mesg-foundation/core/daemon"
	"github.com/mesg-foundation/core/service"
	"github.com/stvp/assert"
)

var serverstop = new(Server)

func TestStopService(t *testing.T) {
	daemon.Start()
	service := service.Service{
		Name: "TestStopService",
		Dependencies: map[string]*service.Dependency{
			"test": &service.Dependency{
				Image: "nginx",
			},
		},
	}
	service.Start()
	reply, err := serverstop.StopService(context.Background(), &StopServiceRequest{
		Service: &service,
	})
	assert.Nil(t, err)
	assert.NotNil(t, reply)
	stopped, err := service.IsStopped()
	assert.Nil(t, err)
	assert.Equal(t, true, stopped)
}
