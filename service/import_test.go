package service

import (
	"testing"

	"github.com/stvp/assert"
)

func TestImportWrongFile(t *testing.T) {
	_, err := ImportFromPath("/path-to-non-existing-file")
	assert.NotNil(t, err)
}

func TestImportMalFormattedFile(t *testing.T) {
	_, err := ImportFromPath("./tests/service-mal-formatted")
	assert.NotNil(t, err)
}

func TestImportNonValidFile(t *testing.T) {
	_, err := ImportFromPath("./tests/service-file-invalid")
	assert.NotNil(t, err)
}

func TestImportValidFile(t *testing.T) {
	service, err := ImportFromPath("./tests/service-minimal-valid")
	assert.Nil(t, err)
	assert.NotNil(t, service)
	assert.Equal(t, service.GetName(), "minimal-valid")
}
