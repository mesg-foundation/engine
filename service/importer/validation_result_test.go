package importer

import (
	"testing"

	"github.com/stvp/assert"
)

// Test ValidationResult.IsValid function

func TestValidationResultIsValid(t *testing.T) {
	v := &ValidationResult{
		ServiceFileExist:    true,
		ServiceFileWarnings: []string{},
		DockerfileExist:     true,
	}
	assert.True(t, v.IsValid())
}

func TestValidationResultIsValidServiceFileDoesNotExist(t *testing.T) {
	v := &ValidationResult{
		ServiceFileExist:    false,
		ServiceFileWarnings: []string{},
		DockerfileExist:     true,
	}
	assert.False(t, v.IsValid())
}

func TestValidationResultIsValidServiceFileWarnings(t *testing.T) {
	v := &ValidationResult{
		ServiceFileExist:    true,
		ServiceFileWarnings: []string{"Warning"},
		DockerfileExist:     true,
	}
	assert.False(t, v.IsValid())
}

func TestValidationResultIsValidDockfileDoesNotExist(t *testing.T) {
	v := &ValidationResult{
		ServiceFileExist:    true,
		ServiceFileWarnings: []string{},
		DockerfileExist:     false,
	}
	assert.False(t, v.IsValid())
}
