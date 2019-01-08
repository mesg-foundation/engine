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

// Test From function

func TestFrom(t *testing.T) {
	service, err := From("./tests/service-minimal-valid")
	require.NoError(t, err)
	require.NotNil(t, service)
	require.Equal(t, service.Name, "minimal-valid")
}

func TestFromMalFormattedFile(t *testing.T) {
	_, err := From("./tests/service-file-mal-formatted")
	require.Error(t, err)
}

func TestFromValidationError(t *testing.T) {
	_, err := From("./tests/service-file-invalid")
	require.Error(t, err)
	_, typeCasting := err.(*ValidationError)
	require.True(t, typeCasting)
	require.Equal(t, (&ValidationError{}).Error(), err.Error())
}
