package importer

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// Test Validate function

func TestValidate(t *testing.T) {
	validation, err := Validate("./tests/service-valid")
	require.Nil(t, err)
	require.True(t, validation.IsValid())
	require.True(t, validation.ServiceFileExist)
	require.Equal(t, 0, len(validation.ServiceFileWarnings))
	require.True(t, validation.DockerfileExist)
}

func TestValidateDockerfileIsMissing(t *testing.T) {
	validation, err := Validate("./tests/service-docker-missing")
	require.Nil(t, err)
	require.False(t, validation.IsValid())
	require.True(t, validation.ServiceFileExist)
	require.Equal(t, 0, len(validation.ServiceFileWarnings))
	require.False(t, validation.DockerfileExist)
}

func TestValidateFromMissingServiceFile(t *testing.T) {
	validation, err := Validate("./tests/service-file-missing")
	require.Nil(t, err)
	require.False(t, validation.IsValid())
	require.False(t, validation.ServiceFileExist)
	require.Equal(t, 1, len(validation.ServiceFileWarnings))
	require.True(t, validation.DockerfileExist)
}

func TestValidateFromNonExistingPath(t *testing.T) {
	validation, err := Validate("./tests/service-non-existing")
	require.Nil(t, err)
	require.False(t, validation.IsValid())
	require.False(t, validation.ServiceFileExist)
	require.Equal(t, 1, len(validation.ServiceFileWarnings))
	require.False(t, validation.DockerfileExist)
}

func TestValidateFromMalFormattedServiceFile(t *testing.T) {
	_, err := Validate("./tests/service-file-mal-formatted")
	require.NotNil(t, err)
}

func TestValidateFromInvalidServiceFile(t *testing.T) {
	validation, err := Validate("./tests/service-file-invalid")
	require.Nil(t, err)
	require.False(t, validation.IsValid())
	require.True(t, validation.ServiceFileExist)
	require.Equal(t, 1, len(validation.ServiceFileWarnings))
	require.True(t, validation.DockerfileExist)
}

// Test IsValid function

func TestIsValid(t *testing.T) {
	isValid, err := IsValid("./tests/service-valid")
	require.Nil(t, err)
	require.True(t, isValid)
}

func TestIsValidMalFormattedServiceFile(t *testing.T) {
	_, err := IsValid("./tests/service-file-mal-formatted")
	require.NotNil(t, err)
}
