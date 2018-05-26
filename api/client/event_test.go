package client

import (
	"testing"

	"github.com/stvp/assert"
)

func TestWhen(t *testing.T) {
	event := Event{
		Name:    "TestWhen",
		Service: "xxxx",
	}
	wf := When(&event)
	assert.NotNil(t, wf)
	assert.Equal(t, wf.Event.Name, event.Name)
	assert.Equal(t, wf.Event.Service, event.Service)
}
