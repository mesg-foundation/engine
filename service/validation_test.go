package service

import (
	"testing"

	"github.com/stvp/assert"
)

func TestMinimalValidFile(t *testing.T) {
	valid, warnings, err := ValidServiceFile("./tests/minimal-valid.yml")
	assert.Nil(t, err)
	assert.Equal(t, valid, true)
	assert.Equal(t, len(warnings), 0)
}

func TestValidFile(t *testing.T) {
	valid, warnings, err := ValidServiceFile("./tests/valid.yml")
	assert.Nil(t, err)
	assert.Equal(t, valid, true)
	assert.Equal(t, len(warnings), 0)
}

func TestNonExistingFile(t *testing.T) {
	_, _, err := ValidServiceFile("./tests/non-existing-file.yml")
	assert.NotNil(t, err)
}

func TestMalFormattedFile(t *testing.T) {
	valid, warnings, err := ValidServiceFile("./tests/mal-formatted.yml")
	assert.Nil(t, err)
	assert.Equal(t, valid, false)
	assert.Equal(t, len(warnings), 1)
}

func TestInvalidFile(t *testing.T) {
	valid, warnings, err := ValidServiceFile("./tests/non-valid.yml")
	assert.Nil(t, err)
	assert.Equal(t, valid, false)
	assert.Equal(t, len(warnings), 1)
}

func TestValidPath(t *testing.T) {
	valid, warnings, err := ValidService("./tests/service-valid")
	assert.Nil(t, err)
	assert.Equal(t, valid, true)
	assert.Equal(t, len(warnings), 0)
}

func TestValidPathInvalidFile(t *testing.T) {
	valid, warnings, err := ValidService("./tests/service-file-invalid")
	assert.Nil(t, err)
	assert.Equal(t, valid, false)
	assert.Equal(t, len(warnings), 1)
}

func TestValidPathMissingYml(t *testing.T) {
	_, _, err := ValidService("./tests/service-file-missing")
	assert.NotNil(t, err)
}

func TestValidPathMissingDocker(t *testing.T) {
	_, _, err := ValidService("./tests/service-docker-missing")
	assert.NotNil(t, err)
}
