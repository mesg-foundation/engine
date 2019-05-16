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
	namespace := service.Namespace()
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
	require.Equal(t, dep.Namespace(service.Namespace()), []string{hex.EncodeToString(sum[:]), "test"})
}

func TestEventSubTopic(t *testing.T) {
	service := &Service{Hash: "1"}
	require.Equal(t, service.EventSubTopic(), hash.Calculate(append(
		service.Namespace(),
		eventTopic,
	)))
}

func TestExecutionSubTopic(t *testing.T) {
	service := &Service{Hash: "1"}
	require.Equal(t, service.ExecutionSubTopic(), hash.Calculate(append(
		service.Namespace(),
		executionTopic,
	)))
}
