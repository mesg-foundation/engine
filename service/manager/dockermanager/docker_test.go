package dockermanager

import (
	"crypto/sha1"
	"encoding/hex"
	"testing"

	"github.com/mesg-foundation/core/hash"
	"github.com/stretchr/testify/require"
)

func TestServiceNamespace(t *testing.T) {
	hash := hash.Int(1)
	namespace := serviceNamespace(hash)
	sum := sha1.Sum(hash)
	require.Equal(t, namespace, []string{hex.EncodeToString(sum[:])})
}

func TestDependencyNamespace(t *testing.T) {
	var (
		hash          = hash.Int(1)
		dependencyKey = "test"
	)
	sNamespace := serviceNamespace(hash)
	sum := sha1.Sum(hash)
	require.Equal(t, dependencyNamespace(sNamespace, dependencyKey), []string{hex.EncodeToString(sum[:]), dependencyKey})
}
