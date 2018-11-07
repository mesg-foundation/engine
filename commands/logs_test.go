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
