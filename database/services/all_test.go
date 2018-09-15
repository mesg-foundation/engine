package services

import (
	"testing"

	"github.com/mesg-foundation/core/service"
	"github.com/stretchr/testify/require"
)

func TestAll(t *testing.T) {
	service := &service.Service{Name: "Service"}
	Save(service)
	defer Delete(service.ID)
	services, err := All()
	require.Nil(t, err)

	found := false
	for _, s := range services {
		if s.Name == "Service" {
			found = true
			break
		}
	}
	require.True(t, found)
}
