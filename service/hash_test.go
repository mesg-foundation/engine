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
	assert.Equal(t, string(hash), "v1_1096bf901b57b1a7e647a85da0b50ad2")
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
