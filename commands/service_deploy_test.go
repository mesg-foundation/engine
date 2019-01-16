package commands

import (
	"bufio"
	"strings"
	"testing"

	"github.com/mesg-foundation/core/commands/provider"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestServiceDeployCmdFlags(t *testing.T) {
	c := newServiceDeployCmd(nil)

	flags := c.cmd.Flags()
	flags.Set("env", "a=1")
	flags.Set("env", "b=2")
	require.Equal(t, map[string]string{"a": "1", "b": "2"}, c.env)
}

func TestServiceDeploy(t *testing.T) {
	var (
		url                     = "1"
		id                      = "2"
		env                     = map[string]string{"A": "3", "B": "4"}
		m                       = newMockExecutor()
		c                       = newServiceDeployCmd(m)
		serviceDeployParameters = []interface{}{
			url,
			env,
			mock.Anything,
		}
		serviceDeployRunFunction = func(args mock.Arguments) {
			statuses := args.Get(2).(chan provider.DeployStatus)
			statuses <- provider.DeployStatus{
				Message: "5",
				Type:    provider.DonePositive,
			}
			statuses <- provider.DeployStatus{
				Message: "6",
				Type:    provider.DoneNegative,
			}
			close(statuses)
		}
	)
	c.cmd.SetArgs([]string{url})
	c.env = env

	m.On("ServiceDeploy", serviceDeployParameters...).Return(id, nil, nil).Run(serviceDeployRunFunction)

	closeStd := captureStd(t)
	c.cmd.Execute()
	stdout, _ := closeStd()
	r := bufio.NewReader(strings.NewReader(stdout))

	require.Equal(t, "✔ 5", string(readLine(t, r)))
	require.Equal(t, "⨯ 6", string(readLine(t, r)))
	require.Equal(t, "✔ Service deployed with hash: 2", string(readLine(t, r)))
	require.Equal(t, "To start it, run the command:", string(readLine(t, r)))
	require.Equal(t, "	mesg-core service start 2", string(readLine(t, r)))

	m.AssertExpectations(t)
}
