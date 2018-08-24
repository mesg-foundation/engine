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
	require.Equal(t, string(hash), "v1_cc79b9b84d6647ddb42e1a860c89a2a8")
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
