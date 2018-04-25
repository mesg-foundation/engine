package result

import (
	"context"
	"testing"

	"github.com/mesg-foundation/core/pubsub"
	"github.com/mesg-foundation/core/service"
	"github.com/mesg-foundation/core/types"
	"github.com/stvp/assert"
)

var serversubmit = new(Server)

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

	subscription := pubsub.Subscribe(service.ResultSubscriptionChannel())

	go serversubmit.Submit(context.Background(), &types.SubmitResultRequest{
		Service: &protoService,
		Output:  "test",
		Task:    "task test",
		Data:    "",
	})

	res := <-subscription
	assert.NotNil(t, res)
}
