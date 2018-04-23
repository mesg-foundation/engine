package event

import (
	"testing"

	"github.com/mesg-foundation/application/types"
	"github.com/stvp/assert"
)

var serverlisten = new(Server)

func TestGetSubscription(t *testing.T) {
	protoService := types.ProtoService{
		Name: "TestGetSubscription",
		Dependencies: map[string]*types.ProtoDependency{
			"test": &types.ProtoDependency{
				Image: "nginx",
			},
		},
	}

	subscription := getSubscription(&types.ListenEventRequest{
		Service: &protoService,
	})

	assert.NotNil(t, subscription)
}
