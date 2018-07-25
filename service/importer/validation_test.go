package importer

import (
	"testing"

	"github.com/stvp/assert"
)

// Test Validate function

func TestValidate(t *testing.T) {
	validation, err := Validate("./tests/service-valid")
	assert.Nil(t, err)
	assert.True(t, validation.IsValid())
	assert.True(t, validation.ServiceFileExist)
	assert.Equal(t, 0, len(validation.ServiceFileWarnings))
	assert.True(t, validation.DockerfileExist)
}

func TestValidateDockerfileIsMissing(t *testing.T) {
	validation, err := Validate("./tests/service-docker-missing")
	assert.Nil(t, err)
	assert.False(t, validation.IsValid())
	assert.True(t, validation.ServiceFileExist)
	assert.Equal(t, 0, len(validation.ServiceFileWarnings))
	assert.False(t, validation.DockerfileExist)
}

func TestValidateFromMissingServiceFile(t *testing.T) {
	validation, err := Validate("./tests/service-file-missing")
	assert.Nil(t, err)
	assert.False(t, validation.IsValid())
	assert.False(t, validation.ServiceFileExist)
	assert.Equal(t, 1, len(validation.ServiceFileWarnings))
	assert.True(t, validation.DockerfileExist)
}

func TestValidateFromNonExistingPath(t *testing.T) {
	validation, err := Validate("./tests/service-non-existing")
	assert.Nil(t, err)
	assert.False(t, validation.IsValid())
	assert.False(t, validation.ServiceFileExist)
	assert.Equal(t, 1, len(validation.ServiceFileWarnings))
	assert.False(t, validation.DockerfileExist)
}

func TestValidateFromMalFormattedServiceFile(t *testing.T) {
	_, err := Validate("./tests/service-file-mal-formatted")
	assert.NotNil(t, err)
}

func TestValidateFromInvalidServiceFile(t *testing.T) {
	validation, err := Validate("./tests/service-file-invalid")
	assert.Nil(t, err)
	assert.False(t, validation.IsValid())
	assert.True(t, validation.ServiceFileExist)
	assert.Equal(t, 1, len(validation.ServiceFileWarnings))
	assert.True(t, validation.DockerfileExist)
}

// Test IsValid function

func TestIsValid(t *testing.T) {
	isValid, err := IsValid("./tests/service-valid")
	assert.Nil(t, err)
	assert.True(t, isValid)
}

func TestIsValidMalFormattedServiceFile(t *testing.T) {
	_, err := IsValid("./tests/service-file-mal-formatted")
	assert.NotNil(t, err)
}
