package importer

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// Test Validate function

func TestValidate(t *testing.T) {
	validation, err := Validate("./tests/service-valid")
	require.NoError(t, err)
	require.True(t, validation.IsValid())
	require.True(t, validation.ServiceFileExist)
	require.Len(t, validation.ServiceFileWarnings, 0)
	require.True(t, validation.DockerfileExist)
}

func TestValidateYMLNames(t *testing.T) {
	validation, err := Validate("./tests/service-names-valid")
	require.NoError(t, err)
	require.True(t, validation.IsValid())
	require.True(t, validation.ServiceFileExist)
	require.Len(t, validation.ServiceFileWarnings, 0)
}

func TestValidateDockerfileIsMissing(t *testing.T) {
	validation, err := Validate("./tests/service-docker-missing")
	require.NoError(t, err)
	require.False(t, validation.IsValid())
	require.True(t, validation.ServiceFileExist)
	require.Len(t, validation.ServiceFileWarnings, 0)
	require.False(t, validation.DockerfileExist)
}

func TestValidateFromMissingServiceFile(t *testing.T) {
	validation, err := Validate("./tests/service-file-missing")
	require.NoError(t, err)
	require.False(t, validation.IsValid())
	require.False(t, validation.ServiceFileExist)
	require.Len(t, validation.ServiceFileWarnings, 1)
	require.True(t, validation.DockerfileExist)
}

func TestValidateFromNonExistingPath(t *testing.T) {
	validation, err := Validate("./tests/service-non-existing")
	require.NoError(t, err)
	require.False(t, validation.IsValid())
	require.False(t, validation.ServiceFileExist)
	require.Len(t, validation.ServiceFileWarnings, 1)
	require.False(t, validation.DockerfileExist)
}

func TestValidateFromMalFormattedServiceFile(t *testing.T) {
	_, err := Validate("./tests/service-file-mal-formatted")
	require.Error(t, err)
}

func TestValidateFromInvalidServiceFile(t *testing.T) {
	validation, err := Validate("./tests/service-file-invalid")
	require.NoError(t, err)
	require.False(t, validation.IsValid())
	require.True(t, validation.ServiceFileExist)
	require.Len(t, validation.ServiceFileWarnings, 1)
	require.True(t, validation.DockerfileExist)
}

// Test IsValid function

func TestIsValid(t *testing.T) {
	isValid, err := IsValid("./tests/service-valid")
	require.NoError(t, err)
	require.True(t, isValid)
}

func TestIsValidMalFormattedServiceFile(t *testing.T) {
	_, err := IsValid("./tests/service-file-mal-formatted")
	require.Error(t, err)
}
