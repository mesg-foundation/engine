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

package commands

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestStartCmdRunE(t *testing.T) {
	m := newMockExecutor()
	c := newStartCmd(m)

	closeStd := captureStd(t)

	m.On("Start").Return(nil)
	c.cmd.Execute()

	stdout, stderr := closeStd()
	require.Contains(t, stdout, "Starting Core")
	require.Contains(t, stdout, "Core started")
	require.Empty(t, stderr)

	m.AssertExpectations(t)
}

func TestStartCmdFlags(t *testing.T) {
	c := newStartCmd(nil)
	require.Equal(t, "text", c.lfv.String())
	require.Equal(t, "info", c.llv.String())

	require.NoError(t, c.cmd.Flags().Set("log-format", "json"))
	require.Equal(t, "json", c.lfv.String())

	require.NoError(t, c.cmd.Flags().Set("log-level", "debug"))
	require.Equal(t, "debug", c.llv.String())
}
