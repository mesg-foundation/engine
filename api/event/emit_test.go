package event

import (
	"context"
	"testing"

	"github.com/mesg-foundation/core/pubsub"
	"github.com/mesg-foundation/core/service"
	"github.com/mesg-foundation/core/types"
	"github.com/stvp/assert"
)

var serveremit = new(Server)

func TestEmit(t *testing.T) {
	protoService := types.ProtoService{
		Name: "TestEmit",
		Dependencies: map[string]*types.ProtoDependency{
			"test": &types.ProtoDependency{
				Image: "nginx",
			},
		},
	}
	service := service.New(&protoService)

	subscription := pubsub.Subscribe(service.EventSubscriptionChannel())

	go serveremit.Emit(context.Background(), &types.EmitEventRequest{
		Service: &protoService,
		Data:    "",
	})

	res := <-subscription
	assert.NotNil(t, res)
}
