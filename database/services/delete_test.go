package services

import (
	"testing"

	"github.com/mesg-foundation/core/service"
	"github.com/stretchr/testify/require"
)

func TestDelete(t *testing.T) {
	service := &service.Service{
		Name: "TestDelete",
	}
	Save(service)
	err := Delete(service.Id)
	require.Nil(t, err)
}
