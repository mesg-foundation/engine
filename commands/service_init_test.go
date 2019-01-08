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

func TestServiceInitCmdFlags(t *testing.T) {
	c := newServiceInitCmd(nil)

	flags := c.cmd.Flags()
	require.Equal(t, flags.ShorthandLookup("t"), flags.Lookup("template"))

	flags.Set("dir", "/")
	require.Equal(t, "/", c.dir)

	flags.Set("template", "github.com/mesg-foundation/awesome")
	require.Equal(t, "github.com/mesg-foundation/awesome", c.templateURL)
}
