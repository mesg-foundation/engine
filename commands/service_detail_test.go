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
	"encoding/json"
	"strings"
	"testing"

	"github.com/mesg-foundation/core/protobuf/coreapi"
	"github.com/mesg-foundation/core/utils/pretty"
	"github.com/stretchr/testify/require"
)

func TestServiceDetail(t *testing.T) {
	var (
		id      = "1"
		service = &coreapi.Service{Hash: "2", Name: "3", Events: []*coreapi.Event{{Key: "4"}}}
		m       = newMockExecutor()
		c       = newServiceDetailCmd(m)
	)
	c.cmd.SetArgs([]string{id})

	m.On("ServiceByID", id).Return(service, nil)

	closeStd := captureStd(t)
	c.cmd.Execute()
	stdout, _ := closeStd()
	r := bufio.NewReader(strings.NewReader(stdout))

	data, _ := json.Marshal(service)

	require.Equal(t, "Loading the service...", string(readLine(t, r)))
	require.Equal(t, string(pretty.ColorizeJSON(pretty.FgCyan, nil, true, data)),
		strings.TrimSpace(string(readLine(t, r))))

	m.AssertExpectations(t)
}
