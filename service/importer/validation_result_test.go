package importer

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// Test ValidationResult.IsValid function

func TestValidationResultIsValid(t *testing.T) {
	v := &ValidationResult{
		ServiceFileExist:    true,
		ServiceFileWarnings: []string{},
		DockerfileExist:     true,
	}
	require.True(t, v.IsValid())
}

func TestValidationResultIsValidServiceFileDoesNotExist(t *testing.T) {
	v := &ValidationResult{
		ServiceFileExist:    false,
		ServiceFileWarnings: []string{},
		DockerfileExist:     true,
	}
	require.False(t, v.IsValid())
}

func TestValidationResultIsValidServiceFileWarnings(t *testing.T) {
	v := &ValidationResult{
		ServiceFileExist:    true,
		ServiceFileWarnings: []string{"Warning"},
		DockerfileExist:     true,
	}
	require.False(t, v.IsValid())
}

func TestValidationResultIsValidDockfileDoesNotExist(t *testing.T) {
	v := &ValidationResult{
		ServiceFileExist:    true,
		ServiceFileWarnings: []string{},
		DockerfileExist:     false,
	}
	require.False(t, v.IsValid())
}
