package core

import (
	"context"
	"testing"

	"github.com/mesg-foundation/core/daemon"
	"github.com/mesg-foundation/core/service"
	"github.com/stvp/assert"
)

var serverstart = new(Server)

func TestStartService(t *testing.T) {
	daemon.Start()
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
	assert.NotNil(t, reply)
	running, err := service.IsRunning()
	assert.Nil(t, err)
	assert.Equal(t, true, running)
	service.Stop()
}
