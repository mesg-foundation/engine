package core

import (
	"context"
	"testing"

	"github.com/mesg-foundation/core/database/services"

	"github.com/mesg-foundation/core/service"
	"github.com/stvp/assert"
)

var serverservice = new(Server)

func TestService(t *testing.T) {
	hash, _ := services.Save(&service.Service{
		Name: "TestService",
	})
	defer services.Delete(hash)
	reply, err := serverservice.GetService(context.Background(), &GetServiceRequest{
		ServiceID: hash,
	})
	assert.Nil(t, err)
	assert.NotNil(t, reply)
	assert.Equal(t, reply.Service.Name, "TestService")
}
