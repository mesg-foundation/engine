package services

import (
	"testing"

	"github.com/mesg-foundation/core/service"
	"github.com/stvp/assert"
)

func TestGet(t *testing.T) {
	hash, _ := Save(&service.Service{
		Name: "TestGet",
	})
	service, err := Get(hash)
	assert.Nil(t, err)
	assert.Equal(t, service.Name, "TestGet")
	delete(hash)
}
