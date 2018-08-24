package service

import (
	"context"
	"testing"

	"github.com/mesg-foundation/core/service"
	"github.com/stretchr/testify/require"
)

func TestEmit(t *testing.T) {
	var (
		path      = "./service-test-event"
		eventKey  = "request"
		eventData = `{"data":{}}`
		server    = newServer(t)
	)

	s, validationErr, err := server.api.DeployService(serviceTar(t, path))
	require.Zero(t, validationErr)
	require.NoError(t, err)
	defer server.api.DeleteService(s.Id)

	ln, err := server.api.ListenEvent(s.Id)
	require.NoError(t, err)
	defer ln.Close()

	_, err = server.EmitEvent(context.Background(), &EmitEventRequest{
		Token:     s.Id,
		EventKey:  eventKey,
		EventData: eventData,
	})
	require.NoError(t, err)

	select {
	case err := <-ln.Err:
		t.Error(err)

	case event := <-ln.Events:
		require.Equal(t, eventKey, event.Key)
		require.Equal(t, eventData, jsonMarshal(t, event.Data))
	}

}

func TestEmitNoData(t *testing.T) {
	var (
		path     = "./service-test-event"
		eventKey = "request"
		server   = newServer(t)
	)

	s, validationErr, err := server.api.DeployService(serviceTar(t, path))
	require.Zero(t, validationErr)
	require.NoError(t, err)
	defer server.api.DeleteService(s.Id)

	_, err = server.EmitEvent(context.Background(), &EmitEventRequest{
		Token:    s.Id,
		EventKey: eventKey,
	})
	require.Equal(t, err.Error(), "unexpected end of JSON input")
}

func TestEmitWrongData(t *testing.T) {
	var (
		path     = "./service-test-event"
		eventKey = "request"
		server   = newServer(t)
	)

	s, validationErr, err := server.api.DeployService(serviceTar(t, path))
	require.Zero(t, validationErr)
	require.NoError(t, err)
	defer server.api.DeleteService(s.Id)

	_, err = server.EmitEvent(context.Background(), &EmitEventRequest{
		Token:     s.Id,
		EventKey:  eventKey,
		EventData: "",
	})
	require.Equal(t, err.Error(), "unexpected end of JSON input")
}

func TestEmitWrongEvent(t *testing.T) {
	var (
		path     = "./service-test-event"
		eventKey = "test"
		server   = newServer(t)
	)

	s, validationErr, err := server.api.DeployService(serviceTar(t, path))
	require.Zero(t, validationErr)
	require.NoError(t, err)
	defer server.api.DeleteService(s.Id)

	_, err = server.EmitEvent(context.Background(), &EmitEventRequest{
		Token:     s.Id,
		EventKey:  eventKey,
		EventData: "{}",
	})
	require.Error(t, err)
	_, notFound := err.(*service.EventNotFoundError)
	require.True(t, notFound)
}

func TestServiceNotExists(t *testing.T) {
	server := newServer(t)
	_, err := server.EmitEvent(context.Background(), &EmitEventRequest{
		Token:     "TestServiceNotExists",
		EventKey:  "test",
		EventData: "{}",
	})
	require.NotNil(t, err)
}
