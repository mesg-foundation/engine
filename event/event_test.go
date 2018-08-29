package event

import (
	"testing"

	"github.com/mesg-foundation/core/service"
	"github.com/stretchr/testify/require"
)

func TestCreate(t *testing.T) {
	s := service.Service{
		Name: "TestCreate",
		Events: map[string]*service.Event{
			"test": {},
		},
	}
	var data map[string]interface{}
	exec, err := Create(&s, "test", data)
	require.Nil(t, err)
	require.Equal(t, &s, exec.Service)
	require.Equal(t, data, exec.Data)
	require.Equal(t, "test", exec.Key)
	require.NotNil(t, exec.CreatedAt)
}

func TestCreateNotPresentEvent(t *testing.T) {
	var (
		serviceName      = "TestCreateNotPresentEvent"
		eventName        = "test"
		invalidEventName = "testInvalid"
	)
	s := service.Service{
		Name: serviceName,
		Events: map[string]*service.Event{
			eventName: {},
		},
	}
	var data map[string]interface{}
	_, err := Create(&s, invalidEventName, data)
	require.Error(t, err)
	notFoundErr, ok := err.(*service.EventNotFoundError)
	require.True(t, ok)
	require.Equal(t, invalidEventName, notFoundErr.EventKey)
	require.Equal(t, serviceName, notFoundErr.ServiceName)
}

func TestCreateInvalidData(t *testing.T) {
	eventName := "test"
	s := service.Service{
		Name: "TestCreateInvalidData",
		Events: map[string]*service.Event{
			eventName: {
				Data: map[string]*service.Parameter{
					"xxx": {},
				},
			},
		},
	}
	var data map[string]interface{}
	_, err := Create(&s, "test", data)
	require.Error(t, err)
	invalidErr, ok := err.(*service.InvalidEventDataError)
	require.True(t, ok)
	require.Equal(t, eventName, invalidErr.EventKey)
}
