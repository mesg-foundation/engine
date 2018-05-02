package service

import (
	"context"
	"testing"

	"github.com/mesg-foundation/core/pubsub"
	"github.com/mesg-foundation/core/service"
	"github.com/stvp/assert"
)

var serveremit = new(Server)

func TestEmit(t *testing.T) {
	service := service.Service{
		Name: "TestEmit",
		Dependencies: map[string]*service.Dependency{
			"test": &service.Dependency{
				Image: "nginx",
			},
		},
	}

	subscription := pubsub.Subscribe(service.EventSubscriptionChannel())

	go serveremit.EmitEvent(context.Background(), &EmitEventRequest{
		Service:   &service,
		EventKey:  "test",
		EventData: "",
	})

	res := <-subscription
	assert.NotNil(t, res)
}
