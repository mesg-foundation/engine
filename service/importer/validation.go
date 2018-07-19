package importer

import (
	"os"
)

// IsValid validates a service at a source
func IsValid(source string) (isValid bool, err error) {
	validation, err := Validate(source)
	if err != nil {
		return
	}
	isValid = validation.IsValid()
	return
}

// Validate validates a service at a source
func Validate(source string) (validation *ValidationResult, err error) {
	return validateFromPath(source)
}

// validateFromPath validates a service at a given path
func validateFromPath(path string) (validation *ValidationResult, err error) {
	// Service file
	validation = &ValidationResult{}
	data, errServicefile := readServiceFile(path)
	validation.ServiceFileExist = errServicefile == nil || os.IsNotExist(errServicefile) == false
	if errServicefile != nil && os.IsNotExist(errServicefile) == false {
		err = errServicefile
		return
	}
	validation.ServiceFileWarnings, err = validateServiceFile(data)
	if err != nil {
		return
	}

	// Dockerfile
	_, errDockerfile := readDockerfile(path)
	validation.DockerfileExist = errDockerfile == nil || os.IsNotExist(errDockerfile) == false
	if errDockerfile != nil && os.IsNotExist(errDockerfile) == false {
		err = errDockerfile
		return
	}

	return
}
