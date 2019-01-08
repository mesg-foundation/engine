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

	"github.com/spf13/pflag"
	"github.com/stretchr/testify/require"
)

func TestLogFormatValue(t *testing.T) {
	lfv := logFormatValue("")
	require.Equal(t, lfv.Type(), "string")

	pflag.Var(&lfv, "test-log-format", "")
	require.NoError(t, pflag.Set("test-log-format", "text"))
	require.NoError(t, pflag.Set("test-log-format", "json"))
	require.Error(t, pflag.Set("test-log-format", "unknown"))
}

func TestLogLevelValue(t *testing.T) {
	llv := logLevelValue("")
	require.Equal(t, llv.Type(), "string")

	pflag.Var(&llv, "test-log-level", "")
	require.NoError(t, pflag.Set("test-log-level", "debug"))
	require.NoError(t, pflag.Set("test-log-level", "info"))
	require.NoError(t, pflag.Set("test-log-level", "warn"))
	require.NoError(t, pflag.Set("test-log-level", "error"))
	require.NoError(t, pflag.Set("test-log-level", "fatal"))
	require.NoError(t, pflag.Set("test-log-level", "panic"))
	require.Error(t, pflag.Set("test-log-level", "unknown"))
}
