package container

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNamespace(t *testing.T) {
	namespace := Namespace([]string{"test"})
	require.Equal(t, namespace, strings.Join([]string{namespacePrefix, "test"}, namespaceSeparator))
}

func TestNamespaceReplaceSpace(t *testing.T) {
	namespace := Namespace([]string{"test foo"})
	require.Equal(t, namespace, strings.Join([]string{namespacePrefix, "test-foo"}, namespaceSeparator))
}
