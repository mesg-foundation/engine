package services

import (
	"testing"

	"github.com/mesg-foundation/core/service"
	"github.com/stvp/assert"
)

func TestGenerateId(t *testing.T) {
	service := &service.Service{
		Name: "TestGenerateId",
	}
	hash, err := calculateHash(service)
	assert.Nil(t, err)
	assert.Equal(t, string(hash), "v1_1096bf901b57b1a7e647a85da0b50ad2")
}

func TestNoCollision(t *testing.T) {
	service1 := &service.Service{
		Name: "TestNoCollision",
	}
	service2 := &service.Service{
		Name: "TestNoCollision2",
	}
	hash1, _ := calculateHash(service1)
	hash2, _ := calculateHash(service2)
	assert.NotEqual(t, string(hash1), string(hash2))
}

func TestSaveReturningHash(t *testing.T) {
	service := &service.Service{
		Name: "TestSaveReturningHash",
	}
	calculatedHash, _ := calculateHash(service)
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
