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

	err := Save(service)
	defer Delete(service.Id)
	require.Nil(t, err)
	require.Equal(t, calculatedHash, service.Id)
}

func TestSaveAndPresentInDB(t *testing.T) {
	service := &service.Service{
		Name: "TestSaveAndPresentInDB",
	}
	Save(service)
	defer Delete(service.Id)
	srv, _ := Get(service.Id)
	require.Equal(t, srv.Name, "TestSaveAndPresentInDB")
}
