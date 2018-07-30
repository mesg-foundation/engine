package client

import (
	"testing"

	"github.com/stvp/assert"
)

func TestServices(t *testing.T) {
	wf := &Workflow{
		OnEvent:  &Event{ServiceID: "xxx"},
		OnResult: &Result{ServiceID: "yyy"},
		Execute:  &Task{ServiceID: "zzz"},
	}
	services := wf.services()
	assert.Equal(t, len(services), 3)
	assert.Equal(t, services[0], "xxx")
	assert.Equal(t, services[1], "yyy")
	assert.Equal(t, services[2], "zzz")
}

func TestServicesDuplicate(t *testing.T) {
	wf := &Workflow{
		OnEvent:  &Event{ServiceID: "xxx"},
		OnResult: &Result{ServiceID: "yyy"},
		Execute:  &Task{ServiceID: "xxx"},
	}
	services := wf.services()
	assert.Equal(t, len(services), 2)
	assert.Equal(t, services[0], "xxx")
	assert.Equal(t, services[1], "yyy")
}
