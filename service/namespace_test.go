package service

import (
	"testing"

	"github.com/stvp/assert"
)

func TestNamespace(t *testing.T) {
	service := &Service{Name: "test"}
	namespace := service.namespace()
	assert.Equal(t, namespace, "MESG-test")
}
