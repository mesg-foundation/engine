package service

import (
	"context"
	"testing"

	"github.com/mesg-foundation/core/database/services"
	"github.com/mesg-foundation/core/pubsub"
	"github.com/mesg-foundation/core/service"
	"github.com/stretchr/testify/require"
)

func TestEmit(t *testing.T) {
	server := newServer(t)
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
	services.Save(&service)
	service.Id = "" // TODO(ilgooz) remove this when Service type created by hand.
	defer services.Delete(service.Id)

	subscription := pubsub.Subscribe(service.EventSubscriptionChannel())

	go server.EmitEvent(context.Background(), &EmitEventRequest{
		Token:     service.Hash(),
		EventKey:  "test",
		EventData: "{}",
	})

	res := <-subscription
	require.NotNil(t, res)
}

func TestEmitNoData(t *testing.T) {
	server := newServer(t)
	service := service.Service{}
	services.Save(&service)
	service.Id = "" // TODO(ilgooz) remove this when Service type created by hand.
	defer services.Delete(service.Id)
	_, err := server.EmitEvent(context.Background(), &EmitEventRequest{
		Token:    service.Hash(),
		EventKey: "test",
	})
	require.Equal(t, err.Error(), "unexpected end of JSON input")
}

func TestEmitWrongData(t *testing.T) {
	server := newServer(t)
	service := service.Service{}
	services.Save(&service)
	service.Id = "" // TODO(ilgooz) remove this when Service type created by hand.
	defer services.Delete(service.Id)
	_, err := server.EmitEvent(context.Background(), &EmitEventRequest{
		Token:     service.Hash(),
		EventKey:  "test",
		EventData: "",
	})
	require.Equal(t, err.Error(), "unexpected end of JSON input")
}

func TestEmitWrongEvent(t *testing.T) {
	server := newServer(t)
	srv := service.Service{Name: "TestEmitWrongEvent"}
	services.Save(&srv)
	srv.Id = "" // TODO(ilgooz) remove this when Service type created by hand.
	defer services.Delete(srv.Id)
	_, err := server.EmitEvent(context.Background(), &EmitEventRequest{
		Token:     srv.Hash(),
		EventKey:  "test",
		EventData: "{}",
	})
	require.NotNil(t, err)
	_, notFound := err.(*service.EventNotFoundError)
	require.True(t, notFound)
}

func TestServiceNotExists(t *testing.T) {
	server := newServer(t)
	service := service.Service{Name: "TestServiceNotExists"}
	_, err := server.EmitEvent(context.Background(), &EmitEventRequest{
		Token:     service.Hash(),
		EventKey:  "test",
		EventData: "{}",
	})
	require.NotNil(t, err)
}
