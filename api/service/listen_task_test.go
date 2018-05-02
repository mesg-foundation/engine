package service

import (
	"testing"

	"github.com/mesg-foundation/core/service"
	"github.com/stvp/assert"
)

var serverlisten = new(Server)

func TestGetSubscription(t *testing.T) {
	service := service.Service{
		Name: "TestGetSubscription",
		Dependencies: map[string]*service.Dependency{
			"test": &service.Dependency{
				Image: "nginx",
			},
		},
	}

	subscription := getSubscription(&ServiceRequest{
		Service: &service,
	})

	assert.NotNil(t, subscription)
}
