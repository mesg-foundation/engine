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
	calculatedHash := service.Hash()
	hash, err := Save(service)
	defer Delete(hash)
	assert.Nil(t, err)
	assert.Equal(t, hash, calculatedHash)
}

func TestSaveAndPresentInDB(t *testing.T) {
	hash, _ := Save(&service.Service{
		Name: "TestSaveAndPresentInDB",
	})
	defer Delete(hash)
	service, _ := Get(hash)
	assert.Equal(t, service.Name, "TestSaveAndPresentInDB")
}
