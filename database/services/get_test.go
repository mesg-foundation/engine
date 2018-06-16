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
	Delete(hash)
}

func TestGetMissing(t *testing.T) {
	emptyService := service.Service{}
	service, err := Get("hash_that_doesnt_exists")
	assert.Equal(t, err, NotFound{Hash: "hash_that_doesnt_exists"})
	assert.Equal(t, service, emptyService)
}
