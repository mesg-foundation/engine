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
	var client core.CoreClient
	err := task.processEvent(client, data)
	assert.Equal(t, err.Error(), "unexpected end of JSON input")
}
