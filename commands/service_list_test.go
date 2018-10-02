package commands

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/mesg-foundation/core/protobuf/coreapi"
	"github.com/stretchr/testify/require"
)

func TestServiceList(t *testing.T) {
	var (
		services = []*coreapi.Service{
			{ID: "1", Name: "a", Status: coreapi.Service_RUNNING},
			{ID: "2", Name: "b", Status: coreapi.Service_PARTIAL},
		}
		m  = &mockServiceExecutor{}
		st = newOutputStream(t)
	)

	c := newServiceListCmd(m, st)
	m.On("ServiceList").Return(services, nil)
	go c.cmd.Execute()

	matched, err := regexp.Match(`\s*^STATUS\s+SERVICE\s+NAME\s*$`, st.ReadLine())
	require.NoError(t, err)
	require.True(t, matched)

	for _, s := range services {
		var status string
		switch s.Status {
		case coreapi.Service_RUNNING:
			status = "running"
		case coreapi.Service_PARTIAL:
			status = "partial"
		}
		pattern := fmt.Sprintf(`\s*^%s\s+%s\s+%s\s*$`, status, s.ID, s.Name)
		matched, err = regexp.Match(pattern, st.ReadLine())
		require.NoError(t, err)
		require.True(t, matched)
	}
}
