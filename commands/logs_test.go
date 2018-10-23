package commands

import (
	"bytes"
	"io/ioutil"
	"testing"

	"github.com/docker/docker/pkg/stdcopy"
	"github.com/mesg-foundation/core/container"
	"github.com/stretchr/testify/assert"
)

func TestLogsCmdRunE(t *testing.T) {
	var (
		m        = newMockExecutor()
		c        = newLogsCmd(m)
		closeStd = captureStd(t)
		buf      = new(bytes.Buffer)
		msg      = []byte("core: 2018-01-01 log\n")
	)

	// create reader for docker stdcopy
	w := stdcopy.NewStdWriter(buf, stdcopy.Stdout)
	w.Write(msg)

	m.On("Status").Return(container.RUNNING, nil)
	m.On("Logs").Return(ioutil.NopCloser(buf), nil)
	c.cmd.Execute()

	m.AssertExpectations(t)

	stdout, stderr := closeStd()
	assert.Zero(t, stderr)
	assert.Equal(t, stdout, string(msg))
}
