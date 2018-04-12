package service

import (
	"testing"

	"github.com/stvp/assert"
)

func TestMinimalValidFile(t *testing.T) {
	warnings, err := ValidService("./tests/service-minimal-valid")
	assert.Nil(t, err)
	assert.Equal(t, len(warnings), 0)
}

func TestValidFile(t *testing.T) {
	warnings, err := ValidService("./tests/service-valid")
	assert.Nil(t, err)
	assert.Equal(t, len(warnings), 0)
}

func TestNonExistingPath(t *testing.T) {
	warnings, err := ValidService("./tests/service-non-existing")
	assert.NotNil(t, err)
	assert.Equal(t, len(warnings), 0)
}

func TestMalFormattedFile(t *testing.T) {
	warnings, err := ValidService("./tests/service-file-mal-formatted")
	assert.NotNil(t, err)
	assert.Equal(t, len(warnings), 1)
}

func TestInvalidFile(t *testing.T) {
	warnings, err := ValidService("./tests/service-file-invalid")
	assert.NotNil(t, err)
	assert.Equal(t, len(warnings), 1)
}

func TestValidPathMissingYml(t *testing.T) {
	warnings, err := ValidService("./tests/service-file-missing")
	assert.NotNil(t, err)
	assert.Equal(t, len(warnings), 0)
}

func TestValidPathMissingDocker(t *testing.T) {
	warnings, err := ValidService("./tests/service-docker-missing")
	assert.NotNil(t, err)
	assert.Equal(t, len(warnings), 0)
}
