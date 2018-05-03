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
		Events: map[string]*service.Event{
			"test": &service.Event{},
		},
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
		EventData: "{}",
	})

	res := <-subscription
	assert.NotNil(t, res)
}

func TestEmitNoData(t *testing.T) {
	_, err := serveremit.EmitEvent(context.Background(), &EmitEventRequest{
		Service:  &service.Service{},
		EventKey: "test",
	})
	assert.Equal(t, err.Error(), "unexpected end of JSON input")
}

func TestEmitWrongData(t *testing.T) {
	_, err := serveremit.EmitEvent(context.Background(), &EmitEventRequest{
		Service:   &service.Service{},
		EventKey:  "test",
		EventData: "",
	})
	assert.Equal(t, err.Error(), "unexpected end of JSON input")
}

func TestEmitWrongEvent(t *testing.T) {
	_, err := serveremit.EmitEvent(context.Background(), &EmitEventRequest{
		Service:   &service.Service{Name: "TestEmitWrongEvent"},
		EventKey:  "test",
		EventData: "{}",
	})
	assert.Equal(t, err.Error(), "Event test doesn't exists in service TestEmitWrongEvent")
}
