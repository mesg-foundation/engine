package importer

import (
	"fmt"
	"os"
)

// IsValid validates a service at a source.
func IsValid(source string) (bool, error) {
	validation, err := Validate(source)
	if err != nil {
		return false, err
	}
	return validation.IsValid(), nil
}

// Validate validates a service at a source.
func Validate(source string) (*ValidationResult, error) {
	return validateFromPath(source)
}

// validateFromPath validates a service at a given path.
func validateFromPath(path string) (*ValidationResult, error) {
	fi, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	if !fi.IsDir() {
		return nil, fmt.Errorf("%s is not a directory", path)
	}

	// Service file
	data, err := readServiceFile(path)
	serviceFileExist := err == nil || !os.IsNotExist(err)
	if err != nil && !os.IsNotExist(err) {
		return nil, err
	}
	serviceFileWarnings := validateServiceFile(data)

	// Dockerfile
	_, err = readDockerfile(path)
	dockerfileExist := err == nil || !os.IsNotExist(err)
	if err != nil && !os.IsNotExist(err) {
		return nil, err
	}

	return &ValidationResult{
		ServiceFileExist:    serviceFileExist,
		ServiceFileWarnings: serviceFileWarnings,
		DockerfileExist:     dockerfileExist,
	}, nil
}
