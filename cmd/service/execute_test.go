package service

import (
	"testing"

	"github.com/mesg-foundation/core/service"
	"github.com/stvp/assert"
)

func TestReadJSONFile(t *testing.T) {
	d, e := readJSONFile("")
	assert.Nil(t, e)
	assert.Equal(t, "{}", d)

	d, e = readJSONFile("./doesntexistsfile")
	assert.NotNil(t, e)

	d, e = readJSONFile("./tests/validData.json")
	assert.Nil(t, e)
	assert.Equal(t, "{\"foo\":\"bar\"}", d)
}

func TestTaskKeysFromService(t *testing.T) {
	keys := taskKeysFromService(&service.Service{
		Tasks: map[string]*service.Task{
			"task1": {},
		},
	})
	assert.Equal(t, "task1", keys[0])
}
