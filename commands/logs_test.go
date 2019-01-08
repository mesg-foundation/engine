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
	"bytes"
	"io/ioutil"
	"testing"

	"github.com/docker/docker/pkg/stdcopy"
	"github.com/mesg-foundation/core/container"
	"github.com/stretchr/testify/require"
)

func TestLogsCmdRunE(t *testing.T) {
	var (
		m        = newMockExecutor()
		c        = newLogsCmd(m)
		closeStd = captureStd(t)
		buf      = new(bytes.Buffer)
		msgout   = []byte("core: 2018-01-01 log\n")
		msgerr   = []byte("core: 2018-01-01 errlog\n")
	)

	// create reader for docker stdcopy
	wout := stdcopy.NewStdWriter(buf, stdcopy.Stdout)
	wout.Write(msgout)

	werr := stdcopy.NewStdWriter(buf, stdcopy.Stderr)
	werr.Write(msgerr)

	m.On("Status").Return(container.RUNNING, nil)
	m.On("Logs").Return(ioutil.NopCloser(buf), nil)
	c.cmd.Execute()

	m.AssertExpectations(t)

	stdout, stderr := closeStd()
	require.Equal(t, string(msgout), stdout)
	require.Equal(t, string(msgerr), stderr)
}
