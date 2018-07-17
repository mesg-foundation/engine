package service

import (
	"testing"

	"github.com/stvp/assert"
)

func TestMinimalValidFile(t *testing.T) {
	validation, err := ValidateService("./tests/service-minimal-valid")
	assert.Nil(t, err)
	assert.True(t, validation.IsValid())
	assert.Equal(t, 0, len(validation.FileWarnings))
	assert.True(t, validation.DockerfileExist)
}

func TestValidFile(t *testing.T) {
	validation, err := ValidateService("./tests/service-valid")
	assert.Nil(t, err)
	assert.True(t, validation.IsValid())
	assert.Equal(t, 0, len(validation.FileWarnings))
	assert.True(t, validation.DockerfileExist)
}

func TestNonExistingPath(t *testing.T) {
	_, err := ValidateService("./tests/service-non-existing")
	assert.NotNil(t, err)
}

func TestMultipleErrors(t *testing.T) {
	validation, err := ValidateService("./tests/service-multiple-errors")
	assert.Nil(t, err)
	assert.False(t, validation.IsValid())
	assert.Equal(t, 2, len(validation.FileWarnings))
	assert.False(t, validation.DockerfileExist)
}

func TestMalFormattedFile(t *testing.T) {
	validation, err := ValidateService("./tests/service-file-mal-formatted")
	assert.Nil(t, err)
	assert.False(t, validation.IsValid())
	assert.Equal(t, 1, len(validation.FileWarnings))
	assert.True(t, validation.DockerfileExist)
}

func TestInvalidFile(t *testing.T) {
	validation, err := ValidateService("./tests/service-file-invalid")
	assert.Nil(t, err)
	assert.False(t, validation.IsValid())
	assert.Equal(t, 1, len(validation.FileWarnings))
	assert.True(t, validation.DockerfileExist)
}

func TestValidPathMissingYml(t *testing.T) {
	_, err := ValidateService("./tests/service-file-missing")
	assert.NotNil(t, err)
}

func TestValidPathMissingDocker(t *testing.T) {
	validation, err := ValidateService("./tests/service-docker-missing")
	assert.Nil(t, err)
	assert.False(t, validation.IsValid())
	assert.False(t, validation.DockerfileExist)
}
