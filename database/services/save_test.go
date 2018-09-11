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

	err := Save(service)
	defer Delete(service.ID)
	require.Nil(t, err)
	require.Equal(t, service.Hash(), service.ID)
}

func TestSaveAndPresentInDB(t *testing.T) {
	service := &service.Service{
		Name: "TestSaveAndPresentInDB",
	}
	Save(service)
	defer Delete(service.ID)
	srv, _ := Get(service.ID)
	require.Equal(t, srv.Name, "TestSaveAndPresentInDB")
}
