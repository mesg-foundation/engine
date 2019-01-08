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

func TestRootServiceCmd(t *testing.T) {
	cmd := newRootServiceCmd(nil).cmd
	for _, tt := range []struct {
		use string
	}{
		{"deploy"},
		{"validate"},
		{"start"},
		{"stop"},
		{"detail"},
		{"list"},
		{"init"},
		{"delete"},
		{"logs"},
		{"gen-doc"},
		{"dev"},
		{"execute"},
	} {
		require.Truef(t, findCommandChildByUsePrefix(cmd, tt.use), "command %q not found", tt.use)
	}
}
