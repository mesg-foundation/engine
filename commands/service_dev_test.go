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

func TestServiceDevCmdFlags(t *testing.T) {
	c := newServiceDevCmd(nil)

	flags := c.cmd.Flags()
	require.Equal(t, flags.ShorthandLookup("e"), flags.Lookup("event-filter"))
	require.Equal(t, flags.ShorthandLookup("t"), flags.Lookup("task-filter"))
	require.Equal(t, flags.ShorthandLookup("o"), flags.Lookup("output-filter"))

	flags.Set("event-filter", "ef")
	require.Equal(t, "ef", c.eventFilter)

	flags.Set("task-filter", "tf")
	require.Equal(t, "tf", c.taskFilter)

	flags.Set("output-filter", "of")
	require.Equal(t, "of", c.outputFilter)
}
