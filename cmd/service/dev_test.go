package service

import (
	"context"
	"testing"

	"github.com/mesg-foundation/core/api/core"
	"github.com/mesg-foundation/core/database/services"
	"github.com/stvp/assert"
)

func TestCreateService(t *testing.T) {
	path := "./tests/test"
	serviceID, err := createService(path)
	defer cli().StopService(context.Background(), &core.StopServiceRequest{
		ServiceID: serviceID,
	})
	defer services.Delete(serviceID)
	assert.Nil(t, err)
	assert.NotNil(t, serviceID)
}
