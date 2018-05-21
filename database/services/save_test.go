package services

import (
	"testing"

	"github.com/mesg-foundation/core/service"
	"github.com/stvp/assert"
)

func TestSaveReturningHash(t *testing.T) {
	service := &service.Service{
		Name: "TestSaveReturningHash",
	}
	calculatedHash, _ := service.Hash()
	hash, err := Save(service)
	assert.Nil(t, err)
	assert.Equal(t, hash, calculatedHash)
	Delete(hash)
}

func TestSaveAndPresentInDB(t *testing.T) {
	hash, _ := Save(&service.Service{
		Name: "TestSaveAndPresentInDB",
	})
	service, _ := Get(hash)
	assert.Equal(t, service.Name, "TestSaveAndPresentInDB")
	Delete(hash)
}
