package service

import (
	"strings"
	"testing"

	"github.com/stvp/assert"
)

func TestNamespace(t *testing.T) {
	service := &Service{Name: "test"}
	namespace := service.Namespace()
	assert.Equal(t, namespace, strings.Join([]string{NAMESPACE, "test"}, "-"))
}
