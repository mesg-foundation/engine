package commands

import (
	"bufio"
	"encoding/json"
	"strings"
	"testing"

	"github.com/mesg-foundation/core/protobuf/coreapi"
	"github.com/mesg-foundation/core/protobuf/definitions"
	"github.com/mesg-foundation/core/utils/pretty"
	"github.com/stretchr/testify/require"
)

func TestServiceDetail(t *testing.T) {
	var (
		id      = "1"
		service = &coreapi.Service{Definition: &definitions.Service{Hash: "2", Name: "3", Events: []*definitions.Event{{Key: "4"}}}}
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
