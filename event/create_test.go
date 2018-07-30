package event

import (
	"testing"

	"github.com/mesg-foundation/core/service"
	"github.com/stvp/assert"
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
	assert.Nil(t, err)
	assert.Equal(t, &s, exec.Service)
	assert.Equal(t, data, exec.Data)
	assert.Equal(t, "test", exec.Key)
	assert.NotNil(t, exec.CreatedAt)
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
	assert.NotNil(t, err)
	_, notFound := err.(*service.EventNotFoundError)
	assert.True(t, notFound)
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
	assert.NotNil(t, err)
	_, invalid := err.(*service.InvalidEventDataError)
	assert.True(t, invalid)
}
