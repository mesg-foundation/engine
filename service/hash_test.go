package service

import (
	"testing"

	"github.com/stvp/assert"
)

func TestGenerateId(t *testing.T) {
	service := Service{
		Name: "TestGenerateId",
	}
	hash := service.Hash()
	assert.Equal(t, string(hash), "v1_c5cbd753f3d7f4f567fcd3c7d5576208")
}

func TestNoCollision(t *testing.T) {
	service1 := Service{
		Name: "TestNoCollision",
	}
	service2 := Service{
		Name: "TestNoCollision2",
	}
	hash1 := service1.Hash()
	hash2 := service2.Hash()
	assert.NotEqual(t, string(hash1), string(hash2))
}
