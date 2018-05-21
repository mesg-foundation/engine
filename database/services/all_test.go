package services

import (
	"fmt"
	"testing"

	"github.com/mesg-foundation/core/service"
	"github.com/stvp/assert"
)

func TestAll(t *testing.T) {
	services, err := All()
	fmt.Println(services)
	assert.Nil(t, err)
	assert.Equal(t, len(services), 0)
}

func TestAfterSave(t *testing.T) {
	hash, _ := Save(&service.Service{Name: "TestAfterSave"})
	services, err := All()
	assert.Nil(t, err)
	assert.Equal(t, len(services), 1)
	Delete(hash)
}
