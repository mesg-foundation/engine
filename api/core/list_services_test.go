package core

import (
	"context"
	"testing"

	"github.com/mesg-foundation/core/database/services"
	"github.com/stretchr/testify/require"
)

var serverlistservices = new(Server)

func TestListServices(t *testing.T) {
	servicesFromAPI, err := serverlistservices.ListServices(context.Background(), &ListServicesRequest{})
	servicesFromDB, _ := services.All()
	require.Nil(t, err)
	require.Equal(t, len(servicesFromAPI.Services), len(servicesFromDB))
}
