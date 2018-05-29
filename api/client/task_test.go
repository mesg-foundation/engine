package client

import (
	"testing"

	"github.com/mesg-foundation/core/api/core"
	"github.com/stvp/assert"
)

func TestProcessEventWithInvalidEventData(t *testing.T) {
	task := &Task{
		Name:    "TestProcessEventWithInvalidEventData",
		Service: "xxx",
	}
	data := &core.EventData{
		EventKey:  "EventX",
		EventData: "",
	}
	err := task.processEvent(data)
	assert.Equal(t, err.Error(), "unexpected end of JSON input")
}

func TestProcessEvent(t *testing.T) {
	task := &Task{
		Name:    "TestProcessEvent",
		Service: "xxx",
		Inputs: func(interface{}) interface{} {
			return "test"
		},
	}
	data := &core.EventData{
		EventKey:  "EventX",
		EventData: "{}",
	}
	err := task.processEvent(data)
	assert.Nil(t, err)
}
