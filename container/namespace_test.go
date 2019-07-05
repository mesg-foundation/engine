package container

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNamespace(t *testing.T) {
	c, _ := New("engine")
	namespace := c.Namespace("foo")
	require.Equal(t, namespace, "engine-foo")
}
