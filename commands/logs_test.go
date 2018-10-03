package commands

import (
	"bytes"
	"io/ioutil"
	"testing"

	"github.com/docker/docker/pkg/stdcopy"
	"github.com/mesg-foundation/core/container"
)

func TestLogsCmdRunE(t *testing.T) {
	m := newMockExecutor()
	c := newLogsCmd(m)

	// create reader for docker stdcopy
	buf := new(bytes.Buffer)
	w := stdcopy.NewStdWriter(buf, stdcopy.Stdout)
	w.Write([]byte("core: 2018-01-01 log\n"))

	m.On("Status").Return(container.RUNNING, nil)
	m.On("Logs").Return(ioutil.NopCloser(buf), nil)
	c.cmd.Execute()

	m.AssertExpectations(t)
}
