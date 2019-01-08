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
)

// Test Validate function

func TestValidate(t *testing.T) {
	validation, err := Validate("./tests/service-valid")
	require.NoError(t, err)
	require.True(t, validation.IsValid())
	require.True(t, validation.ServiceFileExist)
	require.Len(t, validation.ServiceFileWarnings, 0)
	require.True(t, validation.DockerfileExist)
}

func TestValidateYMLNames(t *testing.T) {
	validation, err := Validate("./tests/service-names-valid")
	require.NoError(t, err)
	require.True(t, validation.IsValid())
	require.True(t, validation.ServiceFileExist)
	require.Len(t, validation.ServiceFileWarnings, 0)
}

func TestValidateDockerfileIsMissing(t *testing.T) {
	validation, err := Validate("./tests/service-docker-missing")
	require.NoError(t, err)
	require.False(t, validation.IsValid())
	require.True(t, validation.ServiceFileExist)
	require.Len(t, validation.ServiceFileWarnings, 0)
	require.False(t, validation.DockerfileExist)
}

func TestValidateFromMissingServiceFile(t *testing.T) {
	validation, err := Validate("./tests/service-file-missing")
	require.NoError(t, err)
	require.False(t, validation.IsValid())
	require.False(t, validation.ServiceFileExist)
	require.Len(t, validation.ServiceFileWarnings, 1)
	require.True(t, validation.DockerfileExist)
}

func TestValidateFromNonExistingPath(t *testing.T) {
	validation, err := Validate("./tests/service-non-existing")
	require.NoError(t, err)
	require.False(t, validation.IsValid())
	require.False(t, validation.ServiceFileExist)
	require.Len(t, validation.ServiceFileWarnings, 1)
	require.False(t, validation.DockerfileExist)
}

func TestValidateFromMalFormattedServiceFile(t *testing.T) {
	_, err := Validate("./tests/service-file-mal-formatted")
	require.Error(t, err)
}

func TestValidateFromInvalidServiceFile(t *testing.T) {
	validation, err := Validate("./tests/service-file-invalid")
	require.NoError(t, err)
	require.False(t, validation.IsValid())
	require.True(t, validation.ServiceFileExist)
	require.Len(t, validation.ServiceFileWarnings, 1)
	require.True(t, validation.DockerfileExist)
}

// Test IsValid function

func TestIsValid(t *testing.T) {
	isValid, err := IsValid("./tests/service-valid")
	require.NoError(t, err)
	require.True(t, isValid)
}

func TestIsValidMalFormattedServiceFile(t *testing.T) {
	_, err := IsValid("./tests/service-file-mal-formatted")
	require.Error(t, err)
}

func TestInvalidDependencyName(t *testing.T) {
	validation, err := Validate("./tests/service-invalid-dependency-name")
	require.NoError(t, err)
	require.Len(t, validation.ServiceFileWarnings, 1)
}
