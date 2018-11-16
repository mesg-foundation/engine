package commands

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCreateWorkflow(t *testing.T) {
	var (
		workflowFile = "../workflow-test/simple.yml"
		workflowName = "1"
		workflowID   = "2"
		m            = newMockExecutor()
		c            = newCreateWorkflowCmd(m)
	)

	m.On("CreateWorkflow", workflowFile, workflowName).Once().Return(workflowID, nil)

	closeStd := captureStd(t)
	c.cmd.SetArgs([]string{workflowFile})
	require.NoError(t, c.cmd.Flags().Set("name", workflowName))
	c.cmd.Execute()
	stdout, _ := closeStd()

	ss := strings.Split(stdout, "\n")
	require.Contains(t, ss[0], "Creating workflow...")
	require.Contains(t, ss[1], "Workflow is running")
	require.Contains(t, ss[2], "To see its logs, run the command:")
	require.Contains(t, ss[3], fmt.Sprintf("mesg-core workflow logs %s", workflowID))

	m.AssertExpectations(t)
}
