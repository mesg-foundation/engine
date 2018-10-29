package service

import (
	"context"
	"testing"

	"github.com/mesg-foundation/core/protobuf/serviceapi"
	"github.com/mesg-foundation/core/service"
	"github.com/stretchr/testify/require"
)

func TestEmit(t *testing.T) {
	var (
		eventKey       = "request"
		eventData      = `{"data":{}}`
		server, closer = newServer(t)
	)
	defer closer()

	s, validationErr, err := server.api.DeployService(serviceTar(t, eventServicePath))
	require.Zero(t, validationErr)
	require.NoError(t, err)
	defer server.api.DeleteService(s.ID)

	ln, err := server.api.ListenEvent(s.ID)
	require.NoError(t, err)
	defer ln.Close()

	_, err = server.EmitEvent(context.Background(), &serviceapi.EmitEventRequest{
		Token:     s.ID,
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
		eventKey       = "request"
		server, closer = newServer(t)
	)
	defer closer()

	s, validationErr, err := server.api.DeployService(serviceTar(t, eventServicePath))
	require.Zero(t, validationErr)
	require.NoError(t, err)
	defer server.api.DeleteService(s.ID)

	_, err = server.EmitEvent(context.Background(), &serviceapi.EmitEventRequest{
		Token:    s.ID,
		EventKey: eventKey,
	})
	require.Equal(t, err.Error(), "unexpected end of JSON input")
}

func TestEmitWrongData(t *testing.T) {
	var (
		eventKey       = "request"
		server, closer = newServer(t)
	)
	defer closer()

	s, validationErr, err := server.api.DeployService(serviceTar(t, eventServicePath))
	require.Zero(t, validationErr)
	require.NoError(t, err)
	defer server.api.DeleteService(s.ID)

	_, err = server.EmitEvent(context.Background(), &serviceapi.EmitEventRequest{
		Token:     s.ID,
		EventKey:  eventKey,
		EventData: "",
	})
	require.Equal(t, err.Error(), "unexpected end of JSON input")
}

func TestEmitWrongEvent(t *testing.T) {
	var (
		eventKey       = "test"
		server, closer = newServer(t)
	)
	defer closer()

	s, validationErr, err := server.api.DeployService(serviceTar(t, eventServicePath))
	require.Zero(t, validationErr)
	require.NoError(t, err)
	defer server.api.DeleteService(s.ID)

	_, err = server.EmitEvent(context.Background(), &serviceapi.EmitEventRequest{
		Token:     s.ID,
		EventKey:  eventKey,
		EventData: "{}",
	})
	require.Error(t, err)
	notFoundErr, ok := err.(*service.EventNotFoundError)
	require.True(t, ok)
	require.Equal(t, eventKey, notFoundErr.EventKey)
	require.Equal(t, s.Name, notFoundErr.ServiceName)
}

func TestEmitInvalidData(t *testing.T) {
	var (
		eventKey       = "request"
		eventData      = `{"body":{}}`
		server, closer = newServer(t)
	)
	defer closer()

	s, validationErr, err := server.api.DeployService(serviceTar(t, eventServicePath))
	require.Zero(t, validationErr)
	require.NoError(t, err)
	defer server.api.DeleteService(s.ID)

	_, err = server.EmitEvent(context.Background(), &serviceapi.EmitEventRequest{
		Token:     s.ID,
		EventKey:  eventKey,
		EventData: eventData,
	})
	require.Error(t, err)
	invalidErr, ok := err.(*service.InvalidEventDataError)
	require.True(t, ok)
	require.Equal(t, eventKey, invalidErr.EventKey)
	require.Equal(t, s.Name, invalidErr.ServiceName)
}

func TestServiceNotExists(t *testing.T) {
	server, closer := newServer(t)
	defer closer()

	_, err := server.EmitEvent(context.Background(), &serviceapi.EmitEventRequest{
		Token:     "TestServiceNotExists",
		EventKey:  "test",
		EventData: "{}",
	})
	require.Error(t, err)
}
