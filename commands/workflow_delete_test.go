package commands

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDeleteWorkflow(t *testing.T) {
	var (
		workflowID = "1"
		m          = newMockExecutor()
		c          = newDeleteWorkflowCmd(m)
	)

	m.On("DeleteWorkflow", workflowID).Return(nil)

	closeStd := captureStd(t)
	c.cmd.SetArgs([]string{workflowID})
	c.cmd.Execute()
	stdout, _ := closeStd()

	ss := strings.Split(stdout, "\n")
	require.Contains(t, ss[0], fmt.Sprintf("Deleting workflow %q", workflowID))
	require.Contains(t, ss[1], fmt.Sprintf("workflow %q deleted", workflowID))

	m.AssertExpectations(t)
}
