package container

import (
	"strings"
	"testing"

	"github.com/mesg-foundation/engine/config"
	"github.com/stretchr/testify/require"
)

func TestNamespace(t *testing.T) {
	cfg, _ := config.Global()
	c, _ := New()
	namespace := c.Namespace([]string{"test"})
	require.Equal(t, namespace, strings.Join([]string{cfg.Name, "test"}, namespaceSeparator))
}

func TestNamespaceReplaceSpace(t *testing.T) {
	cfg, _ := config.Global()
	c, _ := New()
	namespace := c.Namespace([]string{"test foo"})
	require.Equal(t, namespace, strings.Join([]string{cfg.Name, "test-foo"}, namespaceSeparator))
}
