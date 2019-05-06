package importer

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// Test readServiceFile function

func TestReadServiceFile(t *testing.T) {
	data, err := readServiceFile("./tests/service-valid")
	require.NoError(t, err)
	require.True(t, (len(data) > 0))
}

func TestReadServiceFileDoesNotExist(t *testing.T) {
	data, err := readServiceFile("./tests/service-file-missing")
	require.Error(t, err)
	require.True(t, (len(data) == 0))
}

// Test validateServiceFile function

func TestValidateServiceFile(t *testing.T) {
	data, _ := readServiceFile("./tests/service-valid")
	warnings := validateServiceFile(data)
	require.True(t, (len(warnings) == 0))
}

func TestValidateServiceFileMalFormatted(t *testing.T) {
	data, _ := readServiceFile("./tests/service-file-mal-formatted")
	warnings := validateServiceFile(data)
	require.Equal(t, 1, len(warnings))
}

func TestValidateServiceFileWithErrors(t *testing.T) {
	data, _ := readServiceFile("./tests/service-file-invalid")
	warnings := validateServiceFile(data)
	require.Equal(t, 2, len(warnings))
}

func TestValidateServiceFileWithMultipleErrors(t *testing.T) {
	data, _ := readServiceFile("./tests/service-multiple-errors")
	warnings := validateServiceFile(data)
	require.Equal(t, 2, len(warnings))
}

func TestValidateServiceSidTooLong(t *testing.T) {
	data, _ := readServiceFile("./tests/sid-too-long")
	warnings := validateServiceFile(data)
	require.Equal(t, 1, len(warnings))
}

func TestValidateServiceDependencyImageMissing(t *testing.T) {
	data, _ := readServiceFile("./tests/dependency-image-missing")
	warnings := validateServiceFile(data)
	require.Equal(t, 1, len(warnings))
}

func TestValidateServiceDependenciesVolumesFrom(t *testing.T) {
	data, _ := readServiceFile("./tests/dependencies-volumesfrom")
	warnings := validateServiceFile(data)
	require.Equal(t, 2, len(warnings))
}
