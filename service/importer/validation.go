package importer

import (
	"os"
)

// IsValid validates a service at a source
func IsValid(source string) (bool, error) {
	validation, err := Validate(source)
	if err != nil {
		return false, err
	}
	return validation.IsValid(), nil
}

// Validate validates a service at a source
func Validate(source string) (*ValidationResult, error) {
	return validateFromPath(source)
}

// validateFromPath validates a service at a given path
func validateFromPath(path string) (*ValidationResult, error) {
	// Service file
	data, err := readServiceFile(path)
	serviceFileExist := err == nil || os.IsNotExist(err) == false
	if err != nil && os.IsNotExist(err) == false {
		return nil, err
	}
	serviceFileWarnings, err := validateServiceFile(data)
	if err != nil {
		return nil, err
	}

	// Dockerfile
	_, err = readDockerfile(path)
	dockerfileExist := err == nil || os.IsNotExist(err) == false
	if err != nil && os.IsNotExist(err) == false {
		return nil, err
	}

	return &ValidationResult{
		ServiceFileExist:    serviceFileExist,
		ServiceFileWarnings: serviceFileWarnings,
		DockerfileExist:     dockerfileExist,
	}, nil
}
