package service

import (
	"testing"

	"github.com/stvp/assert"
)

func TestGetJSON(t *testing.T) {
	d, e := getJSON("")
	assert.Nil(t, e)
	assert.Equal(t, "{}", d)

	d, e = getJSON("./doesntexistsfile")
	assert.NotNil(t, e)

	d, e = getJSON("./tests/validData.json")
	assert.Nil(t, e)
	assert.Equal(t, "{\"foo\":\"bar\"}", d)
}
