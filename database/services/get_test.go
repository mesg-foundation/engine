package services

import (
	"testing"

	"github.com/mesg-foundation/core/service"
	"github.com/stretchr/testify/require"
)

func TestGet(t *testing.T) {
	service := &service.Service{
		Name: "TestGet",
	}
	Save(service)
	defer Delete(service.ID)
	srv, err := Get(service.ID)
	require.Nil(t, err)
	require.Equal(t, srv.Name, "TestGet")
}

func TestGetMissing(t *testing.T) {
	emptyService := service.Service{}
	service, err := Get("hash_that_doesnt_exists")
	require.Equal(t, err, NotFound{Hash: "hash_that_doesnt_exists"})
	require.Equal(t, service, emptyService)
}
