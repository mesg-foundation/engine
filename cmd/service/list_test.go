package service

import (
	"sort"
	"testing"

	"github.com/mesg-foundation/core/service"
	"github.com/stretchr/testify/require"
)

func TestServiceStatusString(t *testing.T) {
	s := service.Service{Name: "TestServiceStatusString"}
	status := serviceStatus{
		service: &s,
		status:  service.RUNNING,
	}
	require.Contains(t, status.String(), "[Running]")
	require.Contains(t, status.String(), s.ID)
	require.Contains(t, status.String(), s.Name)
}

func TestSort(t *testing.T) {
	status := []serviceStatus{
		{status: service.PARTIAL, service: &service.Service{Name: "Partial"}},
		{status: service.RUNNING, service: &service.Service{Name: "Running"}},
		{status: service.STOPPED, service: &service.Service{Name: "Stopped"}},
	}
	sort.Sort(byStatus(status))
	require.Equal(t, status[0].status, service.RUNNING)
	require.Equal(t, status[1].status, service.PARTIAL)
	require.Equal(t, status[2].status, service.STOPPED)
}

func TestServicesWithStatus(t *testing.T) {
	services := append([]*service.Service{}, &service.Service{Name: "TestServicesWithStatus"})
	status, err := servicesWithStatus(services)
	require.Nil(t, err)
	require.Equal(t, status[0].status, service.STOPPED)
}
