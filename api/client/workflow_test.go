package client

import (
	"testing"

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
