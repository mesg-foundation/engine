package services

import (
	"testing"

	"github.com/mesg-foundation/core/service"
	"github.com/stretchr/testify/require"
)

func TestSaveAndPresentInDB(t *testing.T) {
	service := &service.Service{
		Name: "TestSaveAndPresentInDB",
	}
	Save(service)
	defer Delete(service.ID)
	srv, _ := Get(service.ID)
	require.Equal(t, srv.Name, "TestSaveAndPresentInDB")
}
