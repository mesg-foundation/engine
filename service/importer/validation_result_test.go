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
