package service

import (
	"sort"
	"testing"

	"github.com/mesg-foundation/core/service"
	"github.com/stvp/assert"
)

func TestServiceStatusString(t *testing.T) {
	s := service.Service{Name: "TestServiceStatusString"}
	status := serviceStatus{
		service: &s,
		status:  service.RUNNING,
	}
	assert.Contains(t, "[Running]", status.String())
	assert.Contains(t, s.Hash(), status.String())
	assert.Contains(t, s.Name, status.String())
}

func TestSort(t *testing.T) {
	status := []serviceStatus{
		serviceStatus{status: service.PARTIAL, service: &service.Service{Name: "Partial"}},
		serviceStatus{status: service.RUNNING, service: &service.Service{Name: "Running"}},
		serviceStatus{status: service.STOPPED, service: &service.Service{Name: "Stopped"}},
	}
	sort.Sort(byStatus(status))
	assert.Equal(t, status[0].status, service.RUNNING)
	assert.Equal(t, status[1].status, service.PARTIAL)
	assert.Equal(t, status[2].status, service.STOPPED)
}

func TestServicesWithStatus(t *testing.T) {
	services := append([]*service.Service{}, &service.Service{Name: "TestServicesWithStatus"})
	status, err := servicesWithStatus(services)
	assert.Nil(t, err)
	assert.Equal(t, status[0].status, service.STOPPED)
}
