package service

import (
	"testing"

	"github.com/stvp/assert"
)

func TestLogs(t *testing.T) {
	service := &Service{
		Name: "TestLogs",
		Dependencies: map[string]*Dependency{
			"test": {
				Image: "nginx",
			},
			"test2": {
				Image: "nginx",
			},
		},
	}
	service.Start()
	defer service.Stop()
	readers, err := service.Logs("*")
	assert.Nil(t, err)
	assert.Equal(t, 2, len(readers))
}

func TestLogsOnlyOneDependency(t *testing.T) {
	service := &Service{
		Name: "TestLogsOnlyOneDependency",
		Dependencies: map[string]*Dependency{
			"test": {
				Image: "nginx",
			},
			"test2": {
				Image: "nginx",
			},
		},
	}
	service.Start()
	defer service.Stop()
	readers, err := service.Logs("test2")
	assert.Nil(t, err)
	assert.Equal(t, 1, len(readers))
}
