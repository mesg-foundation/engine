package dockermanager

import (
	"crypto/sha1"
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestServiceNamespace(t *testing.T) {
	hash := "1"
	namespace := serviceNamespace(hash)
	sum := sha1.Sum([]byte(hash))
	require.Equal(t, namespace, []string{hex.EncodeToString(sum[:])})
}

func TestDependencyNamespace(t *testing.T) {
	var (
		hash          = "1"
		dependencyKey = "test"
	)
	sNamespace := serviceNamespace(hash)
	sum := sha1.Sum([]byte(hash))
	require.Equal(t, dependencyNamespace(sNamespace, dependencyKey), []string{hex.EncodeToString(sum[:]), dependencyKey})
}
