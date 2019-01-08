// Copyright 2018 MESG Foundation
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package importer

import (
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
	// Service file
	data, err := readServiceFile(path)
	serviceFileExist := err == nil || !os.IsNotExist(err)
	if err != nil && !os.IsNotExist(err) {
		return nil, err
	}
	serviceFileWarnings, err := validateServiceFile(data)
	if err != nil {
		return nil, err
	}

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
