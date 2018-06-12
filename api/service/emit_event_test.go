package service

import (
	"context"
	"testing"

	"github.com/mesg-foundation/core/database/services"

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
	hash, _ := services.Save(&service)
	defer services.Delete(hash)

	subscription := pubsub.Subscribe(service.EventSubscriptionChannel())

	go serveremit.EmitEvent(context.Background(), &EmitEventRequest{
		ServiceHash: service.Hash(),
		EventKey:    "test",
		EventData:   "{}",
	})

	res := <-subscription
	assert.NotNil(t, res)
}

func TestEmitNoData(t *testing.T) {
	service := service.Service{}
	hash, _ := services.Save(&service)
	defer services.Delete(hash)
	_, err := serveremit.EmitEvent(context.Background(), &EmitEventRequest{
		ServiceHash: service.Hash(),
		EventKey:    "test",
	})
	assert.Equal(t, err.Error(), "unexpected end of JSON input")
}

func TestEmitWrongData(t *testing.T) {
	service := service.Service{}
	hash, _ := services.Save(&service)
	defer services.Delete(hash)
	_, err := serveremit.EmitEvent(context.Background(), &EmitEventRequest{
		ServiceHash: service.Hash(),
		EventKey:    "test",
		EventData:   "",
	})
	assert.Equal(t, err.Error(), "unexpected end of JSON input")
}

func TestEmitWrongEvent(t *testing.T) {
	service := service.Service{Name: "TestEmitWrongEvent"}
	hash, _ := services.Save(&service)
	defer services.Delete(hash)
	_, err := serveremit.EmitEvent(context.Background(), &EmitEventRequest{
		ServiceHash: service.Hash(),
		EventKey:    "test",
		EventData:   "{}",
	})
	assert.Equal(t, err.Error(), "Event test doesn't exists in service TestEmitWrongEvent")
}

func TestServiceNotExists(t *testing.T) {
	service := service.Service{Name: "TestServiceNotExists"}
	_, err := serveremit.EmitEvent(context.Background(), &EmitEventRequest{
		ServiceHash: service.Hash(),
		EventKey:    "test",
		EventData:   "{}",
	})
	assert.NotNil(t, err)
}
