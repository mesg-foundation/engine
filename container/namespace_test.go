package container

import (
	"strings"
	"testing"

	"github.com/mesg-foundation/core/config"
	"github.com/stretchr/testify/require"
)

func TestNamespace(t *testing.T) {
	c, _ := config.Global()
	namespace := Namespace([]string{"test"})
	require.Equal(t, namespace, strings.Join([]string{c.Core.Name, "test"}, namespaceSeparator))
}

func TestNamespaceReplaceSpace(t *testing.T) {
	c, _ := config.Global()
	namespace := Namespace([]string{"test foo"})
	require.Equal(t, namespace, strings.Join([]string{c.Core.Name, "test-foo"}, namespaceSeparator))
}
