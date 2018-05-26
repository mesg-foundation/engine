package client

import (
	"errors"
	"testing"

	"github.com/mesg-foundation/core/api/core"
	"github.com/stvp/assert"
)

func TestExtractService(t *testing.T) {
	wf := When(&Event{Service: "XXX", Name: "TestExtractService"}).
		Then(&Task{Service: "YYY", Name: "TestExtractService"}).
		Then(&Task{Service: "ZZZ", Name: "TestExtractService"})
	services := wf.getServices()
	assert.Equal(t, len(services), 3)
	assert.Equal(t, services[0], "XXX")
	assert.Equal(t, services[1], "YYY")
	assert.Equal(t, services[2], "ZZZ")
}

func TestExtractServiceDuplicate(t *testing.T) {
	wf := When(&Event{Service: "XXX", Name: "TestExtractService"}).
		Then(&Task{Service: "YYY", Name: "TestExtractService"}).
		Then(&Task{Service: "YYY", Name: "TestExtractService2"})
	services := wf.getServices()
	assert.Equal(t, len(services), 2)
	assert.Equal(t, services[0], "XXX")
	assert.Equal(t, services[1], "YYY")
}

func TestValidEvent(t *testing.T) {
	wf := When(&Event{Service: "XXX", Name: "TestValidEvent"})
	assert.True(t, wf.validEvent(&core.EventData{EventKey: "TestValidEvent"}))
	assert.False(t, wf.validEvent(&core.EventData{EventKey: "invalid"}))
}

func TestIterateTask(t *testing.T) {
	wf := When(&Event{Service: "XXX", Name: "TestIterateTask"}).
		Then(&Task{Service: "YYY", Name: "TestIterateTask2"}).
		Then(&Task{Service: "YYY", Name: "TestIterateTask3"})
	var cpt = 0
	err := wf.iterateTask(func(task *Task) error {
		cpt++
		return nil
	})
	assert.Nil(t, err)
	assert.Equal(t, cpt, 2)
}

func TestIterateTaskError(t *testing.T) {
	wf := When(&Event{Service: "XXX", Name: "TestIterateTask"}).
		Then(&Task{Service: "YYY", Name: "TestIterateTask2"})
	err := wf.iterateTask(func(task *Task) error {
		return errors.New("error")
	})
	assert.NotNil(t, err)
}

func TestIterateService(t *testing.T) {
	wf := When(&Event{Service: "XXX", Name: "TestIterateService"}).
		Then(&Task{Service: "YYY", Name: "TestIterateService2"}).
		Then(&Task{Service: "ZZZ", Name: "TestIterateService3"})
	var cpt = 0
	err := wf.iterateService(func(service string) (interface{}, error) {
		cpt++
		return nil, nil
	})
	assert.Nil(t, err)
	assert.Equal(t, cpt, 3)
}

func TestIterateServiceWithErr(t *testing.T) {
	wf := When(&Event{Service: "XXX", Name: "TestIterateService"}).
		Then(&Task{Service: "YYY", Name: "TestIterateService2"}).
		Then(&Task{Service: "ZZZ", Name: "TestIterateService3"})
	err := wf.iterateService(func(service string) (interface{}, error) {
		return nil, errors.New("error")
	})
	assert.NotNil(t, err)
}
