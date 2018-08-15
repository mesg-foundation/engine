package importer

import (
	"testing"

	"github.com/stretchr/testify/require"
	yaml "gopkg.in/yaml.v2"
)

// Test readServiceFile function

func TestReadServiceFile(t *testing.T) {
	data, err := readServiceFile("./tests/service-valid")
	require.Nil(t, err)
	require.True(t, (len(data) > 0))
}

func TestReadServiceFileDoesNotExist(t *testing.T) {
	data, err := readServiceFile("./tests/service-file-missing")
	require.NotNil(t, err)
	require.True(t, (len(data) == 0))
}

// Test validateServiceFileSchema function

func TestValidateServiceFileSchema(t *testing.T) {
	data, _ := readServiceFile("./tests/service-valid")
	var body interface{}
	_ = yaml.Unmarshal(data, &body)
	body = convert(body)
	result, err := validateServiceFileSchema(body, "service/importer/assets/schema.json")
	require.Nil(t, err)
	require.True(t, result.Valid())
}

func TestValidateServiceFileSchemaNotExisting(t *testing.T) {
	data, _ := readServiceFile("./tests/service-valid")
	var body interface{}
	_ = yaml.Unmarshal(data, &body)
	body = convert(body)
	_, err := validateServiceFileSchema(body, "service/assets/not_existing")
	require.NotNil(t, err)
}

// Test validateServiceFile function

func TestValidateServiceFile(t *testing.T) {
	data, _ := readServiceFile("./tests/service-valid")
	warnings, err := validateServiceFile(data)
	require.Nil(t, err)
	require.True(t, (len(warnings) == 0))
}

func TestValidateServiceFileMalFormatted(t *testing.T) {
	data, _ := readServiceFile("./tests/service-file-mal-formatted")
	warnings, err := validateServiceFile(data)
	require.NotNil(t, err)
	require.True(t, (len(warnings) == 0))
}

func TestValidateServiceFileWithErrors(t *testing.T) {
	data, _ := readServiceFile("./tests/service-file-invalid")
	warnings, err := validateServiceFile(data)
	require.Nil(t, err)
	require.Equal(t, 1, len(warnings))
}

func TestValidateServiceFileWithMultipleErrors(t *testing.T) {
	data, _ := readServiceFile("./tests/service-multiple-errors")
	warnings, err := validateServiceFile(data)
	require.Nil(t, err)
	require.Equal(t, 2, len(warnings))
}
