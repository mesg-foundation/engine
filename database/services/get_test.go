package services

import (
	"testing"

	"github.com/mesg-foundation/core/service"
	"github.com/stretchr/testify/require"
)

func TestGet(t *testing.T) {
	s := &service.Service{
		Name: "TestGet",
	}
	Save(s)
	defer Delete(s.ID)
	sr, err := Get(s.ID)
	require.Nil(t, err)
	require.Equal(t, sr.Name, "TestGet")
}

func TestGetMissing(t *testing.T) {
	s, err := Get("hash_that_doesnt_exists")
	require.Equal(t, err, NotFound{Hash: "hash_that_doesnt_exists"})
	require.Zero(t, s)
}
