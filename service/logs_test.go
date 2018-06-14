package service

import (
	"testing"

	"github.com/stvp/assert"
)

func TestLogs(t *testing.T) {
	service := &Service{
		Name: "TestLogs",
		Dependencies: map[string]*Dependency{
			"test": &Dependency{
				Image: "nginx",
			},
			"test2": &Dependency{
				Image: "nginx",
			},
		},
	}
	service.Start()
	readers, err := service.Logs("*")
	assert.Nil(t, err)
	assert.Equal(t, 2, len(readers))
	defer service.Stop()
}

func TestLogsOnlyOneDependency(t *testing.T) {
	service := &Service{
		Name: "TestLogsOnlyOneDependency",
		Dependencies: map[string]*Dependency{
			"test": &Dependency{
				Image: "nginx",
			},
			"test2": &Dependency{
				Image: "nginx",
			},
		},
	}
	service.Start()
	readers, err := service.Logs("test2")
	assert.Nil(t, err)
	assert.Equal(t, 1, len(readers))
	defer service.Stop()
}
