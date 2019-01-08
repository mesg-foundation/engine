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
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/mesg-foundation/core/protobuf/coreapi"
	"github.com/stretchr/testify/require"
)

func TestServiceList(t *testing.T) {
	var (
		services = []*coreapi.Service{
			{Hash: "1", Name: "a", Status: coreapi.Service_RUNNING},
			{Hash: "2", Name: "b", Status: coreapi.Service_PARTIAL},
		}
		m = newMockExecutor()
		c = newServiceListCmd(m)
	)

	m.On("ServiceList").Return(services, nil)

	closeStd := captureStd(t)
	c.cmd.Execute()
	stdout, _ := closeStd()
	r := bufio.NewReader(strings.NewReader(stdout))

	for _, s := range []string{
		`Listing services\.\.\.`,
		`HASH\s+SID\s+NAME\s+STATUS`,
	} {
		matched, err := regexp.Match(fmt.Sprintf(`^\s*%s\s*$`, s), readLine(t, r))
		require.NoError(t, err)
		require.True(t, matched)
	}

	for _, s := range services {
		status := strings.ToLower(s.Status.String())
		pattern := fmt.Sprintf(`^\s*%s\s+%s\s+%s\s+%s\s*$`, s.Hash, s.Sid, s.Name, status)
		matched, err := regexp.Match(pattern, readLine(t, r))
		require.NoError(t, err)
		require.True(t, matched)
	}
}
