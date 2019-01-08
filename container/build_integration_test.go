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

// +build integration

package container

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIntegrationBuild(t *testing.T) {
	c, err := New()
	require.NoError(t, err)
	tag, err := c.Build("test/")
	require.NoError(t, err)
	require.NotEqual(t, "", tag)
}

func TestIntegrationBuildNotWorking(t *testing.T) {
	c, err := New()
	require.NoError(t, err)
	tag, err := c.Build("test-not-valid/")
	require.Error(t, err)
	require.Equal(t, "", tag)
}

func TestIntegrationBuildWrongPath(t *testing.T) {
	c, err := New()
	require.NoError(t, err)
	_, err = c.Build("testss/")
	require.Error(t, err)
}
