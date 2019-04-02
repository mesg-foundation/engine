package service

import (
	"crypto/sha1"
	"encoding/hex"
	"testing"

	"github.com/mesg-foundation/core/utils/hash"
	"github.com/mr-tron/base58"
	"github.com/stretchr/testify/require"
)

func TestServiceNamespace(t *testing.T) {
	service, _ := FromService(&Service{
		Hash: "1",
		Name: "TestServiceNamespace",
	})
	namespace := service.namespace()
	h, err := base58.Decode(service.Hash)
	require.NoError(t, err)
	sum := sha1.Sum(h)
	require.Equal(t, namespace, []string{hex.EncodeToString(sum[:])})
}

func TestDependencyNamespace(t *testing.T) {
	service, _ := FromService(&Service{
		Hash: "1",
		Name: "TestDependencyNamespace",
		Dependencies: []*Dependency{
			{
				Key:   "test",
				Image: "http-server",
			},
		},
	})
	dep := service.Dependencies[0]
	h, err := base58.Decode(service.Hash)
	require.NoError(t, err)
	sum := sha1.Sum(h)
	require.Equal(t, dep.namespace(), []string{hex.EncodeToString(sum[:]), "test"})
}

func TestEventSubscriptionChannel(t *testing.T) {
	service, _ := FromService(&Service{
		Hash: "1",
		Name: "TestEventSubscriptionChannel",
	})
	require.Equal(t, service.EventSubscriptionChannel(), hash.Calculate(append(
		service.namespace(),
		eventChannel,
	)))
}

func TestTaskSubscriptionChannel(t *testing.T) {
	service, _ := FromService(&Service{
		Hash: "1",
		Name: "TaskSubscriptionChannel",
	})
	require.Equal(t, service.TaskSubscriptionChannel(), hash.Calculate(append(
		service.namespace(),
		taskChannel,
	)))
}

func TestResultSubscriptionChannel(t *testing.T) {
	service, _ := FromService(&Service{
		Hash: "1",
		Name: "ResultSubscriptionChannel",
	})
	require.Equal(t, service.ResultSubscriptionChannel(), hash.Calculate(append(
		service.namespace(),
		resultChannel,
	)))
}
