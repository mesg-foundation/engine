package result

import (
	"testing"

	"github.com/mesg-foundation/core/types"
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

	subscription := getSubscription(&types.ListenResultRequest{
		Service: &protoService,
	})

	assert.NotNil(t, subscription)
}
