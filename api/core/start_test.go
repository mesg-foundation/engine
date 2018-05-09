package core

import (
	"context"
	"testing"
	"time"

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
	// TODO: Remove this sleep
	// Sleep added because the test are failling, probably a concurrence issue between
	// the test and the docker api. On a local machine it works fine but maybe the CI
	// have some delay when starting docker
	time.Sleep(2 * time.Second)
	assert.Nil(t, err)
	assert.True(t, service.IsRunning())
	assert.NotNil(t, reply)
	service.Stop()
}
