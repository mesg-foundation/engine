package service

import (
	"context"
	"testing"

	"github.com/mesg-foundation/core/service"
	"github.com/mesg-foundation/core/types"
	"github.com/stvp/assert"
)

var serverstart = new(Server)

func TestStartService(t *testing.T) {
	protoService := types.ProtoService{
		Name: "TestStartService",
		Dependencies: map[string]*types.ProtoDependency{
			"test": &types.ProtoDependency{
				Image: "nginx",
			},
		},
	}
	reply, err := serverstart.Start(context.Background(), &types.StartServiceRequest{
		Service: &protoService,
	})
	service := service.New(&protoService)
	assert.Equal(t, service.IsRunning(), true)
	assert.Nil(t, err)
	assert.NotNil(t, reply)
	service.Stop()
}
