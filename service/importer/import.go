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
	yaml "gopkg.in/yaml.v2"
)

// From imports a service from a source.
func From(source string) (*ServiceDefinition, error) {
	return fromPath(source)
}

// fromPath imports a service from a path.
func fromPath(path string) (*ServiceDefinition, error) {
	isValid, err := IsValid(path)
	if err != nil {
		return nil, err
	}
	if !isValid {
		return nil, &ValidationError{}
	}
	data, err := readServiceFile(path)
	if err != nil {
		return nil, err
	}
	var importedService ServiceDefinition
	return &importedService, yaml.UnmarshalStrict(data, &importedService)
}
