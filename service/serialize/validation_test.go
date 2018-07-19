package serialize

import (
	"testing"

	"github.com/stvp/assert"
)

// Test ValidateFromPath function

func TestValidateFromPath(t *testing.T) {
	validation, err := ValidateFromPath("./tests/service-valid")
	assert.Nil(t, err)
	assert.True(t, validation.IsValid())
	assert.True(t, validation.ServiceFileExist)
	assert.Equal(t, 0, len(validation.ServiceFileWarnings))
	assert.True(t, validation.DockerfileExist)
}

func TestValidateFromPathDockerfileIsMissing(t *testing.T) {
	validation, err := ValidateFromPath("./tests/service-docker-missing")
	assert.Nil(t, err)
	assert.False(t, validation.IsValid())
	assert.True(t, validation.ServiceFileExist)
	assert.Equal(t, 0, len(validation.ServiceFileWarnings))
	assert.False(t, validation.DockerfileExist)
}

func TestValidateFromMissingServiceFile(t *testing.T) {
	validation, err := ValidateFromPath("./tests/service-file-missing")
	assert.Nil(t, err)
	assert.False(t, validation.IsValid())
	assert.False(t, validation.ServiceFileExist)
	assert.Equal(t, 1, len(validation.ServiceFileWarnings))
	assert.True(t, validation.DockerfileExist)
}

func TestValidateFromNonExistingPath(t *testing.T) {
	validation, err := ValidateFromPath("./tests/service-non-existing")
	assert.Nil(t, err)
	assert.False(t, validation.IsValid())
	assert.False(t, validation.ServiceFileExist)
	assert.Equal(t, 1, len(validation.ServiceFileWarnings))
	assert.False(t, validation.DockerfileExist)
}

func TestValidateFromMalFormattedServiceFile(t *testing.T) {
	_, err := ValidateFromPath("./tests/service-file-mal-formatted")
	assert.NotNil(t, err)
}

func TestValidateFromInvalidServiceFile(t *testing.T) {
	validation, err := ValidateFromPath("./tests/service-file-invalid")
	assert.Nil(t, err)
	assert.False(t, validation.IsValid())
	assert.True(t, validation.ServiceFileExist)
	assert.Equal(t, 1, len(validation.ServiceFileWarnings))
	assert.True(t, validation.DockerfileExist)
}

// Test IsValidFromPath function

func TestIsValidFromPath(t *testing.T) {
	isValid, err := IsValidFromPath("./tests/service-valid")
	assert.Nil(t, err)
	assert.True(t, isValid)
}

func TestIsValidFromPathMalFormattedServiceFile(t *testing.T) {
	_, err := IsValidFromPath("./tests/service-file-mal-formatted")
	assert.NotNil(t, err)
}
