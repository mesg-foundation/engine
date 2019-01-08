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
		url                     = "1"
		id                      = "2"
		env                     = []string{"A=3", "B=4"}
		m                       = newMockExecutor()
		c                       = newServiceDeployCmd(m)
		serviceDeployParameters = []interface{}{
			url,
			xos.EnvSliceToMap(env),
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
	c.cmd.Flags().Set("env", strings.Join(env, ","))

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
