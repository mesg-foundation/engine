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

func TestServiceExecuteCmdFlags(t *testing.T) {
	c := newServiceExecuteCmd(nil)

	flags := c.cmd.Flags()
	require.Equal(t, flags.ShorthandLookup("t"), flags.Lookup("task"))
	require.Equal(t, flags.ShorthandLookup("d"), flags.Lookup("data"))
	require.Equal(t, flags.ShorthandLookup("j"), flags.Lookup("json"))

	flags.Set("task", "t")
	require.Equal(t, "t", c.taskKey)

	flags.Set("data", "k=v")
	require.Equal(t, map[string]string{"k": "v"}, c.executeData)

	flags.Set("json", "data.json")
	require.Equal(t, "data.json", c.jsonFile)
}
