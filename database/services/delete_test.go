package services

import (
	"testing"

	"github.com/mesg-foundation/core/service"
	"github.com/stvp/assert"
)

func TestDelete(t *testing.T) {
	hash, _ := Save(&service.Service{
		Name: "TestDelete",
	})
	err := Delete(hash)
	assert.Nil(t, err)
}
