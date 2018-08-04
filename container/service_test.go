package container

import (
	"testing"

	"github.com/stvp/assert"
)

func startTestService(name []string) (string, error) {
	return StartService(ServiceOptions{
		Image:     "nginx",
		Namespace: name,
	})
}

func TestStartService(t *testing.T) {
	namespace := []string{"TestStartService"}
	serviceID, err := startTestService(namespace)
	defer StopService(namespace)
	assert.Nil(t, err)
	assert.NotEqual(t, "", serviceID)
}

func TestStartService2Times(t *testing.T) {
	namespace := []string{"TestStartService2Times"}
	startTestService(namespace)
	defer StopService(namespace)
	serviceID, err := startTestService(namespace)
	assert.NotNil(t, err)
	assert.Equal(t, "", serviceID)
}

func TestStopService(t *testing.T) {
	namespace := []string{"TestStopService"}
	startTestService(namespace)
	err := StopService(namespace)
	assert.Nil(t, err)
}

func TestStopNotExistingService(t *testing.T) {
	namespace := []string{"TestStopNotExistingService"}
	err := StopService(namespace)
	assert.Nil(t, err)
}

func TestServiceStatusNeverStarted(t *testing.T) {
	namespace := []string{"TestServiceStatusNeverStarted"}
	status, err := ServiceStatus(namespace)
	assert.Nil(t, err)
	assert.NotEqual(t, RUNNING, status)
	assert.Equal(t, STOPPED, status)
}

func TestServiceStatusRunning(t *testing.T) {
	namespace := []string{"TestServiceStatusRunning"}
	startTestService(namespace)
	defer StopService(namespace)
	status, err := ServiceStatus(namespace)
	assert.Nil(t, err)
	assert.Equal(t, status, RUNNING)
	assert.NotEqual(t, status, STOPPED)
}

func TestServiceStatusStopped(t *testing.T) {
	namespace := []string{"TestServiceStatusStopped"}
	startTestService(namespace)
	StopService(namespace)
	status, err := ServiceStatus(namespace)
	assert.Nil(t, err)
	assert.Equal(t, status, STOPPED)
	assert.NotEqual(t, status, RUNNING)
}

func TestFindServiceNotExisting(t *testing.T) {
	_, err := FindService([]string{"TestFindServiceNotExisting"})
	assert.NotNil(t, err)
}

func TestFindService(t *testing.T) {
	namespace := []string{"TestFindService"}
	startTestService(namespace)
	defer StopService(namespace)
	service, err := FindService(namespace)
	assert.Nil(t, err)
	assert.NotEqual(t, "", service.ID)
}

func TestFindServiceCloseName(t *testing.T) {
	namespace := []string{"TestFindServiceCloseName", "name"}
	namespace1 := []string{"TestFindServiceCloseName", "name2"}
	startTestService(namespace)
	defer StopService(namespace)
	startTestService(namespace1)
	defer StopService(namespace1)
	service, err := FindService(namespace)
	assert.Nil(t, err)
	assert.NotEqual(t, "", service.ID)
}

func TestFindServiceStopped(t *testing.T) {
	namespace := []string{"TestFindServiceStopped"}
	startTestService(namespace)
	StopService(namespace)
	_, err := FindService(namespace)
	assert.NotNil(t, err)
}

func TestListServices(t *testing.T) {
	StartService(ServiceOptions{
		Image:     "nginx",
		Namespace: []string{"TestListServices"},
		Labels: map[string]string{
			"label_name": "value_1",
		},
	})
	StartService(ServiceOptions{
		Image:     "nginx",
		Namespace: []string{"TestListServiceswithValue2"},
		Labels: map[string]string{
			"label_name_2": "value_2",
		},
	})
	defer StopService([]string{"TestListServices"})
	defer StopService([]string{"TestListServiceswithValue2"})
	services, err := ListServices("label_name")
	assert.Nil(t, err)
	assert.Equal(t, 1, len(services))
	assert.Equal(t, Namespace([]string{"TestListServices"}), services[0].Spec.Name)
}

func TestServiceLogs(t *testing.T) {
	namespace := []string{"TestServiceLogs"}
	startTestService(namespace)
	defer StopService(namespace)
	reader, err := ServiceLogs(namespace)
	assert.Nil(t, err)
	assert.NotNil(t, reader)
}
