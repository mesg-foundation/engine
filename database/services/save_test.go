package services

import (
	"testing"

	"github.com/mesg-foundation/core/service"
	"github.com/stretchr/testify/require"
)

func TestSaveReturningHash(t *testing.T) {
	service := &service.Service{
		Name: "TestSaveReturningHash",
	}
	calculatedHash := service.Hash()
	hash, err := Save(service)
	defer Delete(hash)
	require.Nil(t, err)
	require.Equal(t, hash, calculatedHash)
}

func TestSaveAndPresentInDB(t *testing.T) {
	hash, _ := Save(&service.Service{
		Name: "TestSaveAndPresentInDB",
	})
	defer Delete(hash)
	service, _ := Get(hash)
	require.Equal(t, service.Name, "TestSaveAndPresentInDB")
}
