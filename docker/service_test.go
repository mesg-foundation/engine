package docker

import (
	"testing"

	"github.com/docker/docker/api/types/swarm"
	"github.com/stvp/assert"
)

func startTestService(namespace []string) (dockerService *swarm.Service, err error) {
	return StartService(&ServiceOptions{
		Namespace: namespace,
		Image:     "nginx",
	})
}

func TestStartService(t *testing.T) {
	namespace := []string{"TestStartService"}
	service, err := startTestService(namespace)
	defer StopService(namespace)
	assert.Nil(t, err)
	assert.NotNil(t, service)
}

func TestStopService(t *testing.T) {
	namespace := []string{"TestStopService"}
	startTestService(namespace)
	err := StopService(namespace)
	assert.Nil(t, err)
}

func TestStopServiceAlreadyStopped(t *testing.T) {
	namespace := []string{"TestStopServiceAlreadyStopped"}
	err := StopService(namespace)
	assert.Nil(t, err)
}

func TestServiceStatusNeverStarted(t *testing.T) {
	namespace := []string{"TestServiceStatusNeverStarted"}
	status := ServiceStatus(namespace)
	assert.Equal(t, status, STOPPED)
	assert.Equal(t, IsServiceRunning(namespace), false)
	assert.Equal(t, IsServiceStopped(namespace), true)
}

func TestServiceStatusRunning(t *testing.T) {
	namespace := []string{"TestServiceStatusRunning"}
	startTestService(namespace)
	defer StopService(namespace)
	status := ServiceStatus(namespace)
	assert.Equal(t, status, RUNNING)
	assert.Equal(t, IsServiceRunning(namespace), true)
	assert.Equal(t, IsServiceStopped(namespace), false)
}

func TestServiceStatusStopped(t *testing.T) {
	namespace := []string{"TestServiceStatusStopped"}
	startTestService(namespace)
	StopService(namespace)
	status := ServiceStatus(namespace)
	assert.Equal(t, status, STOPPED)
	assert.Equal(t, IsServiceRunning(namespace), false)
	assert.Equal(t, IsServiceStopped(namespace), true)
}

func TestFindServiceNotExisting(t *testing.T) {
	namespace := []string{"TestFindServiceNotExisting"}
	service, err := FindService(namespace)
	assert.Nil(t, err)
	assert.Nil(t, service)
}

func TestFindService(t *testing.T) {
	namespace := []string{"TestFindService"}
	startTestService(namespace)
	defer StopService(namespace)
	service, err := FindService(namespace)
	assert.Nil(t, err)
	assert.NotNil(t, service)
}

func TestFindServiceStopped(t *testing.T) {
	namespace := []string{"TestFindServiceStopped"}
	startTestService(namespace)
	StopService(namespace)
	service, err := FindService(namespace)
	assert.Nil(t, err)
	assert.Nil(t, service)
}

func TestListServices(t *testing.T) {
	namespace := []string{"TestListServices"}
	namespace2 := []string{"TestListServiceswithValue2"}
	StartService(&ServiceOptions{
		Namespace: namespace,
		Image:     "nginx",
		Labels: map[string]string{
			"label_name": "value_1",
		},
	})
	StartService(&ServiceOptions{
		Namespace: namespace2,
		Image:     "nginx",
		Labels: map[string]string{
			"label_name_2": "value_2",
		},
	})
	defer StopService(namespace)
	defer StopService(namespace2)
	services, err := ListServices("label_name")
	assert.Nil(t, err)
	assert.Equal(t, 1, len(services))
	assert.Equal(t, Namespace(namespace), services[0].Spec.Name)
}
