package commands

import (
	"bufio"
	"strings"
	"testing"

	"github.com/mesg-foundation/core/commands/provider"
	"github.com/mesg-foundation/core/x/xos"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestServiceDeployCmdFlags(t *testing.T) {
	c := newServiceDeployCmd(nil)

	flags := c.cmd.Flags()
	flags.Set("env", "a=1,b=2")
	require.Equal(t, map[string]string{"a": "1", "b": "2"}, c.env)
}

func TestServiceDeploy(t *testing.T) {
	var (
		url = "1"
		id  = "2"
		env = []string{"A=3", "B=4"}
		m   = newMockExecutor()
		c   = newServiceDeployCmd(m)
	)
	c.cmd.SetArgs([]string{url})
	c.cmd.Flags().Set("env", strings.Join(env, ","))

	m.On("ServiceDeploy", url, xos.EnvSliceToMap(env), mock.Anything).
		Return(id, nil, nil).
		Run(func(args mock.Arguments) {
			statuses := args.Get(2).(chan provider.DeployStatus)
			statuses <- provider.DeployStatus{"5", provider.DonePositive}
			statuses <- provider.DeployStatus{"6", provider.DoneNegative}
			close(statuses)
		})

	closeStd := captureStd(t)
	c.cmd.Execute()
	stdout, _ := closeStd()
	r := bufio.NewReader(strings.NewReader(stdout))

	require.Equal(t, `✔ 5
⨯ 6
✔ Service deployed with hash: 2
To start it, run the command:
	mesg-core service start 2`, strings.TrimSpace(string(readAll(t, r))))

	m.AssertExpectations(t)
}
