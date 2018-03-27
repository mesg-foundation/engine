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
	_, err := ImportFromFile("./tests/service-mal-formatted/mesg..yml")
	assert.NotNil(t, err)
}

func TestImportValidFile(t *testing.T) {
	service, err := ImportFromFile("./tests/service-minimal-valid/mesg.yml")
	assert.Nil(t, err)
	assert.NotNil(t, service)
	assert.Equal(t, service.Name, "minimal-valid")
}
