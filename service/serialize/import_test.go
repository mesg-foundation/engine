package serialize

import (
	"testing"

	"github.com/stvp/assert"
)

// Test FromPath function

func TestFromPath(t *testing.T) {
	service, err := FromPath("./tests/service-minimal-valid")
	assert.Nil(t, err)
	assert.NotNil(t, service)
	assert.Equal(t, service.Name, "minimal-valid")
}

func TestFromPathMalFormattedFile(t *testing.T) {
	_, err := FromPath("./tests/service-file-mal-formatted")
	assert.NotNil(t, err)
}

func TestFromPathValidationError(t *testing.T) {
	_, err := FromPath("./tests/service-file-invalid")
	assert.NotNil(t, err)
	_, typeCasting := err.(*ValidationError)
	assert.True(t, typeCasting)
	assert.Equal(t, (&ValidationError{}).Error(), err.Error())
}
