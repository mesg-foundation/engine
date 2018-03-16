package service

import (
	"testing"

	"github.com/stvp/assert"
)

func TestStartService(t *testing.T) {
	service := &Service{
		Name: "test",
		Dependencies: map[string]*Dependency{
			"test": &Dependency{
				Image: "nginx",
			},
		},
	}
	err := service.Start()
	assert.Nil(t, err)
}
