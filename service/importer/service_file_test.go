package importer

import (
	"testing"

	"github.com/stvp/assert"
	yaml "gopkg.in/yaml.v2"
)

// Test readServiceFile function

func TestReadServiceFile(t *testing.T) {
	data, err := readServiceFile("./tests/service-valid")
	assert.Nil(t, err)
	assert.True(t, (len(data) > 0))
}

func TestReadServiceFileDoesNotExist(t *testing.T) {
	data, err := readServiceFile("./tests/service-file-missing")
	assert.NotNil(t, err)
	assert.True(t, (len(data) == 0))
}

// Test validateServiceFileSchema function

func TestValidateServiceFileSchema(t *testing.T) {
	data, _ := readServiceFile("./tests/service-valid")
	var body interface{}
	_ = yaml.Unmarshal(data, &body)
	body = convert(body)
	result, err := validateServiceFileSchema(body, "service/importer/assets/schema.json")
	assert.Nil(t, err)
	assert.True(t, result.Valid())
}

func TestValidateServiceFileSchemaNotExisting(t *testing.T) {
	data, _ := readServiceFile("./tests/service-valid")
	var body interface{}
	_ = yaml.Unmarshal(data, &body)
	body = convert(body)
	_, err := validateServiceFileSchema(body, "service/assets/not_existing")
	assert.NotNil(t, err)
}

// Test validateServiceFile function

func TestValidateServiceFile(t *testing.T) {
	data, _ := readServiceFile("./tests/service-valid")
	warnings, err := validateServiceFile(data)
	assert.Nil(t, err)
	assert.True(t, (len(warnings) == 0))
}

func TestValidateServiceFileMalFormatted(t *testing.T) {
	data, _ := readServiceFile("./tests/service-file-mal-formatted")
	warnings, err := validateServiceFile(data)
	assert.NotNil(t, err)
	assert.True(t, (len(warnings) == 0))
}

func TestValidateServiceFileWithErrors(t *testing.T) {
	data, _ := readServiceFile("./tests/service-file-invalid")
	warnings, err := validateServiceFile(data)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(warnings))
}

func TestValidateServiceFileWithMultipleErrors(t *testing.T) {
	data, _ := readServiceFile("./tests/service-multiple-errors")
	warnings, err := validateServiceFile(data)
	assert.Nil(t, err)
	assert.Equal(t, 2, len(warnings))
}
