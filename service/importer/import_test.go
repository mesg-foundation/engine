package importer

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// Test From function

func TestFrom(t *testing.T) {
	service, err := From("./tests/service-minimal-valid")
	require.NoError(t, err)
	require.NotNil(t, service)
	require.Equal(t, service.Name, "minimal-valid")
}

func TestFromMalFormattedFile(t *testing.T) {
	_, err := From("./tests/service-file-mal-formatted")
	require.Error(t, err)
}

func TestFromValidationError(t *testing.T) {
	_, err := From("./tests/service-file-invalid")
	require.Error(t, err)
	_, typeCasting := err.(*ValidationError)
	require.True(t, typeCasting)
	require.Equal(t, (&ValidationError{}).Error(), err.Error())
}
