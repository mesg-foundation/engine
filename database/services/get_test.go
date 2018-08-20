package services

import (
	"testing"

	"github.com/mesg-foundation/core/service"
	"github.com/stretchr/testify/require"
)

func TestGet(t *testing.T) {
	hash, _ := Save(&service.Service{
		Name: "TestGet",
	})
	defer Delete(hash)
	service, err := Get(hash)
	require.Nil(t, err)
	require.Equal(t, service.Name, "TestGet")
}

func TestGetMissing(t *testing.T) {
	emptyService := service.Service{}
	service, err := Get("hash_that_doesnt_exists")
	require.Equal(t, err, NotFound{Hash: "hash_that_doesnt_exists"})
	require.Equal(t, service, emptyService)
}
