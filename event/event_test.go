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
	s := service.Service{
		Name: "TestCreateNotPresentEvent",
		Events: map[string]*service.Event{
			"test": {},
		},
	}
	var data map[string]interface{}
	_, err := Create(&s, "testinvalid", data)
	require.NotNil(t, err)
	_, notFound := err.(*service.EventNotFoundError)
	require.True(t, notFound)
}

func TestCreateInvalidData(t *testing.T) {
	s := service.Service{
		Name: "TestCreateInvalidData",
		Events: map[string]*service.Event{
			"test": {
				Data: map[string]*service.Parameter{
					"xxx": {},
				},
			},
		},
	}
	var data map[string]interface{}
	_, err := Create(&s, "test", data)
	require.NotNil(t, err)
	_, invalid := err.(*service.InvalidEventDataError)
	require.True(t, invalid)
}
