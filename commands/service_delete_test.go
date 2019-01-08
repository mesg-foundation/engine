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

func TestServiceDeleteCmdFlags(t *testing.T) {
	var (
		c     = newServiceDeleteCmd(nil)
		flags = c.cmd.Flags()
	)

	// check defaults.
	require.False(t, c.yes)
	require.False(t, c.all)
	require.False(t, c.keepData)

	require.Equal(t, flags.ShorthandLookup("y"), flags.Lookup("yes"))

	flags.Set("yes", "true")
	require.True(t, c.yes)

	flags.Set("all", "true")
	require.True(t, c.all)

	flags.Set("keep-data", "true")
	require.True(t, c.keepData)
}

func TestServiceDeletePreRunE(t *testing.T) {
	c := newServiceDeleteCmd(nil)

	c.discardOutput()
	require.Equal(t, errNoID, c.preRunE(c.cmd, nil))

	c.yes = true
	c.all = true
	require.NoError(t, c.preRunE(c.cmd, nil))
}

func TestServiceDeleteRunE(t *testing.T) {
	var (
		m = newMockExecutor()
		c = newServiceDeleteCmd(m)
	)

	tests := []struct {
		all      bool
		keepData bool
		arg      string
	}{
		{all: false, keepData: false, arg: "0"},
		{all: false, keepData: true, arg: "0"},
		{all: true, keepData: false},
		{all: true, keepData: true},
	}

	for _, tt := range tests {
		c.all = tt.all
		c.keepData = tt.keepData

		if tt.all {
			m.On("ServiceDeleteAll", !tt.keepData).Once().Return(nil)
		} else {
			m.On("ServiceDelete", !tt.keepData, tt.arg).Once().Return(nil)
		}
		require.NoError(t, c.runE(c.cmd, []string{tt.arg}))
	}
	m.AssertExpectations(t)
}
