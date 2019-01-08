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
	"testing"

	"github.com/stretchr/testify/require"
	yaml "gopkg.in/yaml.v2"
)

// Test readServiceFile function

func TestReadServiceFile(t *testing.T) {
	data, err := readServiceFile("./tests/service-valid")
	require.NoError(t, err)
	require.True(t, (len(data) > 0))
}

func TestReadServiceFileDoesNotExist(t *testing.T) {
	data, err := readServiceFile("./tests/service-file-missing")
	require.Error(t, err)
	require.True(t, (len(data) == 0))
}

// Test validateServiceFileSchema function

func TestValidateServiceFileSchema(t *testing.T) {
	data, _ := readServiceFile("./tests/service-valid")
	var body interface{}
	_ = yaml.Unmarshal(data, &body)
	body = convert(body)
	result, err := validateServiceFileSchema(body, "service/importer/assets/schema.json")
	require.NoError(t, err)
	require.True(t, result.Valid())
}

func TestValidateServiceFileSchemaNotExisting(t *testing.T) {
	data, _ := readServiceFile("./tests/service-valid")
	var body interface{}
	_ = yaml.Unmarshal(data, &body)
	body = convert(body)
	_, err := validateServiceFileSchema(body, "service/assets/not_existing")
	require.Error(t, err)
}

// Test validateServiceFile function

func TestValidateServiceFile(t *testing.T) {
	data, _ := readServiceFile("./tests/service-valid")
	warnings, err := validateServiceFile(data)
	require.NoError(t, err)
	require.True(t, (len(warnings) == 0))
}

func TestValidateServiceFileMalFormatted(t *testing.T) {
	data, _ := readServiceFile("./tests/service-file-mal-formatted")
	warnings, err := validateServiceFile(data)
	require.Error(t, err)
	require.True(t, (len(warnings) == 0))
}

func TestValidateServiceFileWithErrors(t *testing.T) {
	data, _ := readServiceFile("./tests/service-file-invalid")
	warnings, err := validateServiceFile(data)
	require.NoError(t, err)
	require.Equal(t, 1, len(warnings))
}

func TestValidateServiceFileWithMultipleErrors(t *testing.T) {
	data, _ := readServiceFile("./tests/service-multiple-errors")
	warnings, err := validateServiceFile(data)
	require.NoError(t, err)
	require.Equal(t, 2, len(warnings))
}

func TestValidateServiceSidTooLong(t *testing.T) {
	data, _ := readServiceFile("./tests/sid-too-long")
	warnings, err := validateServiceFile(data)
	require.NoError(t, err)
	require.Contains(t, warnings[0], "sid: String length must be less than or equal to 39")
}
