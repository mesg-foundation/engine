package core

import (
	"context"
	"testing"

	"github.com/mesg-foundation/core/database/services"
	"github.com/stvp/assert"
)

var serverlistservices = new(Server)

func TestListServices(t *testing.T) {
	servicesFromAPI, err := serverlistservices.ListServices(context.Background(), &ListServicesRequest{})
	servicesFromDB, _ := services.All()
	assert.Nil(t, err)
	assert.Equal(t, len(servicesFromAPI.Services), len(servicesFromDB))
}
