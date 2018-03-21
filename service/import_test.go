package service

import (
	"testing"

	"github.com/stvp/assert"
)

func TestImportWrongFile(t *testing.T) {
	_, err := ImportFromFile("/path-to-non-existing-file")
	assert.NotNil(t, err)
}

func TestImportMalFormattedFile(t *testing.T) {
	_, err := ImportFromFile("./tests/mal-formatted.yml")
	assert.NotNil(t, err)
}

func TestImportValidFile(t *testing.T) {
	service, err := ImportFromFile("./tests/minimal-valid.yml")
	assert.Nil(t, err)
	assert.NotNil(t, service)
	assert.Equal(t, service.Name, "minimal-valid")
}
