package container

import (
	"strings"
	"testing"

	"github.com/mesg-foundation/engine/version"
	"github.com/stretchr/testify/require"
)

func TestNamespace(t *testing.T) {
	c, _ := New()
	namespace := c.Namespace([]string{"test foo"})
	require.Equal(t, namespace, strings.Join([]string{version.Name, "test-foo"}, namespaceSeparator))
}
