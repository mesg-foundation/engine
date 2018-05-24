package services

import (
	"testing"

	"github.com/mesg-foundation/core/service"
	"github.com/stvp/assert"
)

func TestAll(t *testing.T) {
	hash, _ := Save(&service.Service{Name: "Service1"})
	services, err := All()
	founded := false
	for _, s := range services {
		if s.Name == "Service1" {
			founded = true
			break
		}
	}
	assert.Nil(t, err)
	assert.True(t, founded)
	Delete(hash)
}
