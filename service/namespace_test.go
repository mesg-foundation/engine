package service

import (
	"testing"

	"github.com/mesg-foundation/core/utils/hash"
	"github.com/stretchr/testify/require"
)

func TestServiceNamespace(t *testing.T) {
	service := &Service{Name: "TestServiceNamespace"}
	namespace := service.namespace()
	require.Equal(t, namespace, []string{service.Hash()})
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
	require.Equal(t, namespace, []string{service.Hash(), "test"})
}

func TestEventSubscriptionChannel(t *testing.T) {
	service := &Service{Name: "TestEventSubscriptionChannel"}
	require.Equal(t, service.EventSubscriptionChannel(), hash.Calculate(append(
		service.namespace(),
		eventChannel,
	)))
}

func TestTaskSubscriptionChannel(t *testing.T) {
	service := &Service{Name: "TaskSubscriptionChannel"}
	require.Equal(t, service.TaskSubscriptionChannel(), hash.Calculate(append(
		service.namespace(),
		taskChannel,
	)))
}

func TestResultSubscriptionChannel(t *testing.T) {
	service := &Service{Name: "ResultSubscriptionChannel"}
	require.Equal(t, service.ResultSubscriptionChannel(), hash.Calculate(append(
		service.namespace(),
		resultChannel,
	)))
}
