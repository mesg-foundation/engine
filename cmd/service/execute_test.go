package service

import (
	"testing"

	"github.com/mesg-foundation/core/service"
	"github.com/stretchr/testify/require"
)

func TestReadJSONFile(t *testing.T) {
	d, e := readJSONFile("")
	require.Nil(t, e)
	require.Equal(t, "{}", d)

	d, e = readJSONFile("./doesntexistsfile")
	require.NotNil(t, e)

	d, e = readJSONFile("./tests/validData.json")
	require.Nil(t, e)
	require.Equal(t, "{\"foo\":\"bar\"}", d)
}

func TestTaskKeysFromService(t *testing.T) {
	keys := taskKeysFromService(&service.Service{
		Tasks: map[string]*service.Task{
			"task1": {},
		},
	})
	require.Equal(t, "task1", keys[0])
}
