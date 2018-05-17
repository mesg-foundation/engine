package core

import (
	"context"
	"testing"

	"github.com/mesg-foundation/core/service"
	"github.com/stvp/assert"
)

var serverstop = new(Server)

var (
	testDaemonIP      = "localhost" // TODO: should be remove when a better implementation is up
	testSharedNetwork = "ingress"   // TODO: should be remove when a better implementation is up
)

func TestStopService(t *testing.T) {
	service := service.Service{
		Name: "TestStopService",
		Dependencies: map[string]*service.Dependency{
			"test": &service.Dependency{
				Image: "nginx",
			},
		},
	}
	service.Start(testDaemonIP, testSharedNetwork)
	reply, err := serverstop.StopService(context.Background(), &StopServiceRequest{
		Service: &service,
	})
	assert.Equal(t, service.IsRunning(), false)
	assert.Nil(t, err)
	assert.NotNil(t, reply)
}
