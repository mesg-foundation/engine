package service

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGenerateId(t *testing.T) {
	service := Service{
		Name: "TestGenerateId",
	}
	hash := service.Hash()
	require.Equal(t, "v1_7cbae42a4e1c847ce6d03ac81af2b533", string(hash))
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
	require.NotEqual(t, string(hash1), string(hash2))
}
