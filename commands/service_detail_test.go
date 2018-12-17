package commands

import (
	"bufio"
	"encoding/json"
	"strings"
	"testing"

	"github.com/mesg-foundation/core/protobuf/coreapi"
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

	data, _ := json.MarshalIndent(service, "", "  ")

	require.Equal(t, "Loading the service...", string(readLine(t, r)))
	require.Equal(t, string(data), strings.TrimSpace(string(readAll(t, r))))

	m.AssertExpectations(t)
}
