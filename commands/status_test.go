package commands

import (
	"testing"

	"github.com/mesg-foundation/core/container"
)

func TestStatusCmdRunE(t *testing.T) {
	m := newMockExecutor()
	c := newStatusCmd(m)

	m.On("Status").Return(container.RUNNING, nil)
	c.cmd.Execute()

	m.AssertExpectations(t)
}
