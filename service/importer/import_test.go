package importer

import (
	"testing"

	"github.com/stvp/assert"
)

// Test From function

func TestFrom(t *testing.T) {
	service, err := From("./tests/service-minimal-valid")
	assert.Nil(t, err)
	assert.NotNil(t, service)
	assert.Equal(t, service.Name, "minimal-valid")
}

func TestFromMalFormattedFile(t *testing.T) {
	_, err := From("./tests/service-file-mal-formatted")
	assert.NotNil(t, err)
}

func TestFromValidationError(t *testing.T) {
	_, err := From("./tests/service-file-invalid")
	assert.NotNil(t, err)
	_, typeCasting := err.(*ValidationError)
	assert.True(t, typeCasting)
	assert.Equal(t, (&ValidationError{}).Error(), err.Error())
}
