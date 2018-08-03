package service

import (
	"testing"

	"github.com/mesg-foundation/core/utils/hash"
	"github.com/stvp/assert"
)

func TestServiceNamespace(t *testing.T) {
	service := &Service{Name: "TestServiceNamespace"}
	namespace := service.namespace()
	assert.Equal(t, namespace, []string{service.Hash()})
}

func TestDependencyNamespace(t *testing.T) {
	service := &Service{
		Name: "TestDependencyNamespace",
		Dependencies: map[string]*Dependency{
			"test": {
				Image: "nginx",
			},
		},
	}
	dep := service.DependenciesFromService()[0]
	namespace := dep.namespace()
	assert.Equal(t, namespace, []string{service.Hash(), "test"})
}

func TestEventSubscriptionChannel(t *testing.T) {
	service := &Service{Name: "TestEventSubscriptionChannel"}
	assert.Equal(t, service.EventSubscriptionChannel(), hash.Calculate(append(
		service.namespace(),
		eventChannel,
	)))
}

func TestTaskSubscriptionChannel(t *testing.T) {
	service := &Service{Name: "TaskSubscriptionChannel"}
	assert.Equal(t, service.TaskSubscriptionChannel(), hash.Calculate(append(
		service.namespace(),
		taskChannel,
	)))
}

func TestResultSubscriptionChannel(t *testing.T) {
	service := &Service{Name: "ResultSubscriptionChannel"}
	assert.Equal(t, service.ResultSubscriptionChannel(), hash.Calculate(append(
		service.namespace(),
		resultChannel,
	)))
}
