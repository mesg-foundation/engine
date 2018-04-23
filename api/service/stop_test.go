package service

import (
	"context"
	"testing"

	"github.com/mesg-foundation/application/service"
	"github.com/mesg-foundation/application/types"
	"github.com/stvp/assert"
)

var serverstop = new(Server)

func TestStopService(t *testing.T) {
	protoService := types.ProtoService{
		Name: "TestStopService",
		Dependencies: map[string]*types.ProtoDependency{
			"test": &types.ProtoDependency{
				Image: "nginx",
			},
		},
	}
	service := service.New(&protoService)
	service.Start()
	reply, err := serverstop.Stop(context.Background(), &types.StopServiceRequest{
		Service: &protoService,
	})
	assert.Equal(t, service.IsRunning(), false)
	assert.Nil(t, err)
	assert.NotNil(t, reply)
}
