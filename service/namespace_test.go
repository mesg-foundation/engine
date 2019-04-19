package service

import (
	"crypto/sha1"
	"encoding/hex"
	"testing"

	"github.com/mesg-foundation/core/utils/hash"
	"github.com/stretchr/testify/require"
)

func TestServiceNamespace(t *testing.T) {
	service := &Service{
		Hash: "1",
		Name: "TestServiceNamespace",
	}
	namespace := service.namespace()
	sum := sha1.Sum([]byte(service.Hash))
	require.Equal(t, namespace, []string{hex.EncodeToString(sum[:])})
}

func TestDependencyNamespace(t *testing.T) {
	service := &Service{
		Hash: "1",
		Name: "TestDependencyNamespace",
		Dependencies: []*Dependency{
			{
				Key:   "test",
				Image: "http-server",
			},
		},
	}
	dep := service.Dependencies[0]
	sum := sha1.Sum([]byte(service.Hash))
	require.Equal(t, dep.namespace(service.namespace()), []string{hex.EncodeToString(sum[:]), "test"})
}

func TestEventSubscriptionChannel(t *testing.T) {
	service := &Service{Hash: "1"}
	require.Equal(t, service.EventSubTopic(), hash.Calculate(append(
		service.namespace(),
		eventTopic,
	)))
}

func TestTaskSubscriptionChannel(t *testing.T) {
	service := &Service{Hash: "1"}
	require.Equal(t, service.ExecutionSubTopic(), hash.Calculate(append(
		service.namespace(),
		executionTopic,
	)))
}
