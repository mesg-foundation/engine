package docker

import (
	"strings"
	"testing"

	"github.com/stvp/assert"
)

func TestNamespace(t *testing.T) {
	namespace := Namespace([]string{"test"})
	assert.Equal(t, namespace, strings.Join([]string{namespacePrefix, "test"}, namespaceSeparator))
}
