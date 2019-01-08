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
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/mesg-foundation/core/service/importer/assets"
	"github.com/xeipuuv/gojsonschema"
	yaml "gopkg.in/yaml.v2"
)

// ConfigurationDependencyKey is the reserved key of the service's configuration in the dependencies array.
const ConfigurationDependencyKey = "service"

func readServiceFile(path string) ([]byte, error) {
	file := filepath.Join(path, "mesg.yml")
	return ioutil.ReadFile(file)
}

// validateServiceFile returns a list of warnings.
func validateServiceFile(data []byte) ([]string, error) {
	var body interface{}
	if err := yaml.Unmarshal(data, &body); err != nil {
		return nil, fmt.Errorf("parse mesg.yml error: %s", err)
	}
	body = convert(body)
	result, err := validateServiceFileSchema(body, "service/importer/assets/schema.json")
	if err != nil {
		return nil, err
	}
	var warnings []string
	if !result.Valid() {
		for _, warning := range result.Errors() {
			warnings = append(warnings, warning.String())
		}
	}
	if depKeyWarning := validateServiceFileDependencyKey(body); depKeyWarning != "" {
		warnings = append(warnings, depKeyWarning)
	}
	return warnings, nil
}

func validateServiceFileDependencyKey(data interface{}) (warning string) {
	s, _ := data.(map[string]interface{})
	dep, _ := s["dependencies"].(map[string]interface{})
	if dep[ConfigurationDependencyKey] != nil {
		return fmt.Sprintf("cannot use %q as dependency key", ConfigurationDependencyKey)
	}
	return ""
}

func validateServiceFileSchema(data interface{}, schemaPath string) (*gojsonschema.Result, error) {
	schemaData, err := assets.Asset(schemaPath)
	if err != nil {
		return nil, err
	}
	schema := gojsonschema.NewBytesLoader(schemaData)
	loaded := gojsonschema.NewGoLoader(data)
	return gojsonschema.Validate(schema, loaded)
}

func convert(i interface{}) interface{} {
	x, ok := i.(map[interface{}]interface{})
	if !ok {
		return i
	}
	m2 := map[string]interface{}{}
	for k, v := range x {
		m2[k.(string)] = convert(v)
	}
	return m2
}
