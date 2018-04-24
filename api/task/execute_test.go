package task

import (
	"context"
	"testing"

	"github.com/mesg-foundation/core/pubsub"
	"github.com/mesg-foundation/core/service"
	"github.com/mesg-foundation/core/types"
	"github.com/stvp/assert"
)

var serverexecute = new(Server)

func TestExecute(t *testing.T) {
	protoService := types.ProtoService{
		Name: "TestExecute",
		Dependencies: map[string]*types.ProtoDependency{
			"test": &types.ProtoDependency{
				Image: "nginx",
			},
		},
	}
	service := service.New(&protoService)

	subscription := pubsub.Subscribe(service.TaskSubscriptionChannel())

	go serverexecute.Execute(context.Background(), &types.ExecuteTaskRequest{
		Service: &protoService,
		Task:    "Test",
		Data:    "",
	})

	res := <-subscription
	assert.NotNil(t, res)
}
