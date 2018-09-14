package service

import (
	"testing"

	"github.com/mesg-foundation/core/utils/hash"
	"github.com/stretchr/testify/require"
)

func TestServiceNamespace(t *testing.T) {
	service, _ := FromService(&Service{Name: "TestServiceNamespace"}, ContainerOption(defaultContainer))
	namespace := service.namespace()
	require.Equal(t, namespace, []string{service.ID})
}

func TestDependencyNamespace(t *testing.T) {
	service, _ := FromService(&Service{
		Name: "TestDependencyNamespace",
		Dependencies: []*Dependency{
			{
				Key:   "test",
				Image: "http-server",
			},
		},
	}, ContainerOption(defaultContainer))
	dep := service.Dependencies[0]
	require.Equal(t, dep.namespace(), []string{service.ID, "test"})
}

func TestEventSubscriptionChannel(t *testing.T) {
	service, _ := FromService(&Service{Name: "TestEventSubscriptionChannel"}, ContainerOption(defaultContainer))
	require.Equal(t, service.EventSubscriptionChannel(), hash.Calculate(append(
		service.namespace(),
		eventChannel,
	)))
}

func TestTaskSubscriptionChannel(t *testing.T) {
	service, _ := FromService(&Service{Name: "TaskSubscriptionChannel"}, ContainerOption(defaultContainer))
	require.Equal(t, service.TaskSubscriptionChannel(), hash.Calculate(append(
		service.namespace(),
		taskChannel,
	)))
}

func TestResultSubscriptionChannel(t *testing.T) {
	service, _ := FromService(&Service{Name: "ResultSubscriptionChannel"}, ContainerOption(defaultContainer))
	require.Equal(t, service.ResultSubscriptionChannel(), hash.Calculate(append(
		service.namespace(),
		resultChannel,
	)))
}
