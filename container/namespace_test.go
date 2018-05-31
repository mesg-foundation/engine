package container

import (
	"strings"
	"testing"

	"github.com/stvp/assert"
)

func TestNamespace(t *testing.T) {
	namespace := Namespace([]string{"test"})
	assert.Equal(t, namespace, strings.Join([]string{namespacePrefix, "test"}, namespaceSeparator))
}

func TestNamespaceReplaceSpace(t *testing.T) {
	namespace := Namespace([]string{"test foo"})
	assert.Equal(t, namespace, strings.Join([]string{namespacePrefix, "test-foo"}, namespaceSeparator))
}

func TestServiceTag(t *testing.T) {
	tag := ServiceTag([]string{"test"})
	assert.Equal(t, tag, strings.Join([]string{serviceTagPrefix, "test"}, namespaceSeparator))
}

func TestServiceTagReplaceSpace(t *testing.T) {
	tag := ServiceTag([]string{"test foo"})
	assert.Equal(t, tag, strings.Join([]string{serviceTagPrefix, "test-foo"}, namespaceSeparator))
}
