package service

import (
	"testing"

	"github.com/mesg-foundation/core/docker"
	"github.com/stvp/assert"
)

func TestList(t *testing.T) {
	service, _, _ := startTestService("TestList", "test")
	defer service.Stop()
	services, err := List()
	assert.Nil(t, err)
	assert.NotNil(t, services)
	assert.Equal(t, len(services), 1)
	assert.Equal(t, services[0].Spec.Name, docker.Namespace([]string{"TestList", "test"}))
}

// func TestListMultipleDependencies(t *testing.T) {
// 	service := &Service{
// 		Name: "TestListMultipleDependencies",
// 		Dependencies: map[string]*Dependency{
// 			"test": &Dependency{
// 				Image: "nginx",
// 			},
// 			"test2": &Dependency{
// 				Image: "nginx",
// 			},
// 		},
// 	}
// 	service.Start()
// 	defer service.Stop()
// 	services, err := List()
// 	assert.Nil(t, err)
// 	assert.Equal(t, len(services), 1)
// 	assert.Equal(t, services[0].Spec.Name, docker.Namespace([]string{"TestList", "test"}))
// }
