package commands

import "testing"

func TestStopCmdRunE(t *testing.T) {
	m := newMockExecutor()
	c := newStopCmd(m)

	m.On("Stop").Return(nil)
	c.cmd.Execute()

	m.AssertExpectations(t)
}
