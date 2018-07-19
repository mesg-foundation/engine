package serialize

import (
	"os"
)

// IsValidFromPath validates a service at a given path
func IsValidFromPath(path string) (isValid bool, err error) {
	validation, err := ValidateFromPath(path)
	if err != nil {
		return
	}
	isValid = validation.IsValid()
	return
}

// ValidateFromPath validates a service at a given path
func ValidateFromPath(path string) (validation *ValidationResult, err error) {
	// Service file
	validation = &ValidationResult{
		ServiceFileExist: true,
		DockerfileExist:  true,
	}
	data, err := readServiceFile(path)
	if os.IsNotExist(err) {
		err = nil
		validation.ServiceFileExist = false
	}
	if err != nil {
		return
	}
	validation.ServiceFileWarnings, err = validateServiceFile(data)
	if err != nil {
		return
	}

	// Dockerfile
	_, err = readDockerfile(path)
	if os.IsNotExist(err) {
		err = nil
		validation.DockerfileExist = false
	}
	return
}
