package service

import (
	"context"
	"testing"

	"github.com/mesg-foundation/core/database/services"
	"github.com/mesg-foundation/core/pubsub"
	"github.com/mesg-foundation/core/service"
	"github.com/stretchr/testify/require"
)

var serveremit = new(Server)

func TestEmit(t *testing.T) {
	service := service.Service{
		Name: "TestEmit",
		Events: map[string]*service.Event{
			"test": {},
		},
		Dependencies: map[string]*service.Dependency{
			"test": {
				Image: "nginx",
			},
		},
	}
	hash, _ := services.Save(&service)
	defer services.Delete(hash)

	subscription := pubsub.Subscribe(service.EventSubscriptionChannel())

	go serveremit.EmitEvent(context.Background(), &EmitEventRequest{
		Token:     service.Hash(),
		EventKey:  "test",
		EventData: "{}",
	})

	res := <-subscription
	require.NotNil(t, res)
}

func TestEmitNoData(t *testing.T) {
	service := service.Service{}
	hash, _ := services.Save(&service)
	defer services.Delete(hash)
	_, err := serveremit.EmitEvent(context.Background(), &EmitEventRequest{
		Token:    service.Hash(),
		EventKey: "test",
	})
	require.Equal(t, err.Error(), "unexpected end of JSON input")
}

func TestEmitWrongData(t *testing.T) {
	service := service.Service{}
	hash, _ := services.Save(&service)
	defer services.Delete(hash)
	_, err := serveremit.EmitEvent(context.Background(), &EmitEventRequest{
		Token:     service.Hash(),
		EventKey:  "test",
		EventData: "",
	})
	require.Equal(t, err.Error(), "unexpected end of JSON input")
}

func TestEmitWrongEvent(t *testing.T) {
	srv := service.Service{Name: "TestEmitWrongEvent"}
	hash, _ := services.Save(&srv)
	defer services.Delete(hash)
	_, err := serveremit.EmitEvent(context.Background(), &EmitEventRequest{
		Token:     srv.Hash(),
		EventKey:  "test",
		EventData: "{}",
	})
	require.NotNil(t, err)
	_, notFound := err.(*service.EventNotFoundError)
	require.True(t, notFound)
}

func TestServiceNotExists(t *testing.T) {
	service := service.Service{Name: "TestServiceNotExists"}
	_, err := serveremit.EmitEvent(context.Background(), &EmitEventRequest{
		Token:     service.Hash(),
		EventKey:  "test",
		EventData: "{}",
	})
	require.NotNil(t, err)
}
